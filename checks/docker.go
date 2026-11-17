package checks

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	log "github.com/sirupsen/logrus"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
)

// CheckDocker gathers information about Docker containers
type CheckDocker struct {
}

// Name will be used in the response as check name
func (c *CheckDocker) Name() string {
	return "docker"
}

type resultDocker struct {
	Id               string  `json:"id"`                // First 10 chars of the container id
	Name             string  `json:"name"`              // First name of the docker container
	Image            string  `json:"image"`             // Name of the used Docker image
	SizeRw           int64   `json:"size_rw"`           // Modified or created files (delta/diff to the docker image) in bytes
	SizeRootFs       int64   `json:"size_root_fs"`      // Total size of the containers file system in bytes
	State            string  `json:"state"`             // created, restarting, running, removing, paused, exited, dead // https://docs.docker.com/engine/api/v1.41/#operation/ContainerList
	Status           string  `json:"status"`            // Up 41 minutes
	NetworkRx        float64 `json:"network_rx"`        // received data in bytes
	NetworkTx        float64 `json:"network_tx"`        // sent data in bytes
	CpuPercentage    float64 `json:"cpu_percentage"`    // CPU usage of the container as percentage
	MemoryPercentage float64 `json:"memory_percentage"` // Memory usage of the container as percentage
	DiskRead         uint64  `json:"disk_read"`         // Linux and macOS only - Size of read bytes
	DiskWrite        uint64  `json:"disk_write"`        // Linux and macOS only - Size of written bytes
	MemoryUsed       float64 `json:"memory_used"`       // Used memory of the container in bytes (Windows, Linux and macOs)

	Health              string `json:"health"`                // healthy, unhealthy, starting or empty if the container has no healthcheck
	HealthFailingStreak int    `json:"health_failing_streak"` // Number of consecutive failed healthchecks
	HealthOutput        string `json:"health_output"`         // Output of the last healthcheck
	RestartCount        int    `json:"restart_count"`         // Restarts by the restart policy since the container was created
	RestartsLastHour    int    `json:"restarts_last_hour"`    // Restarts by the restart policy within the last hour, detected by the agent
	OOMKilled           bool   `json:"oom_killed"`            // The last run of the container was killed because of the memory limit
	ExitCode            int    `json:"exit_code"`             // Exit code of the last run
	StartedAt           int64  `json:"started_at"`            // Unix timestamp
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckDocker) Run(ctx context.Context) (interface{}, error) {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	daemonOSType, err := c.negotiateAPIVersion(ctx, cli)
	if err != nil {
		return nil, err
	}

	list, err := cli.ContainerList(ctx, client.ContainerListOptions{
		All: false, //Only show running containers
	})
	if err != nil {
		return nil, err
	}
	containers := list.Items

	existing := make(map[string]bool, len(containers))
	for _, ctr := range containers {
		existing[ctr.ID] = true
	}
	forgetRemovedDockerContainers(existing, time.Now())

	var wg sync.WaitGroup
	resultChan := make(chan *resultDocker, len(containers))
	errorChan := make(chan error, len(containers))

	dockerResults := make([]*resultDocker, 0)
	for _, ctr := range containers {
		// we use ctr as variable name (instead of container) to avoid shadowing the package name "container" / naming conflict

		wg.Add(1)
		go func(ctr container.Summary) {
			// This is a new "thread" / go func for each docker container
			defer wg.Done()

			// IncludePreviousSample: the daemon takes two samples, the CPU usage is calculated from PreCPUStats
			response, err := cli.ContainerStats(ctx, ctr.ID, client.ContainerStatsOptions{
				Stream:                false,
				IncludePreviousSample: true,
			})
			if err != nil {
				errorChan <- err
				return
			}

			defer response.Body.Close()
			responseJson := json.NewDecoder(response.Body)

			var stats *container.StatsResponse
			if err := responseJson.Decode(&stats); err != nil {
				errorChan <- err
				return
			}

			var bytesRead, bytesWrite uint64
			var memory, memoryPercentage, cpuPercentage, networkRx, networkTx float64

			// Memory and CPU usage needs to be calculated manually
			// https://github.com/docker/cli/blob/a4a07c643042f4e2a75bf872f38b134502848214/cli/command/container/stats_helpers.go#L79-L128
			if dockerStatsOSType(stats, daemonOSType) != "windows" {
				// Docker daemon is running in Linux or macOS
				previousCPU := float64(stats.PreCPUStats.CPUUsage.TotalUsage)
				previousSystem := float64(stats.PreCPUStats.SystemUsage)
				cpuPercentage = c.calcCpuPercentageUnix(previousCPU, previousSystem, stats)
				bytesRead, bytesWrite = c.calcDiskIoUnix(stats.BlkioStats)
				networkRx, networkTx = c.calcNetworkIo(stats.Networks)
				memory = c.calcMemoryUsageUnix(stats.MemoryStats)
				memoryLimit := float64(stats.MemoryStats.Limit)
				memoryPercentage = c.calcMemoryPercentageUnix(memoryLimit, memory)
			} else {
				// Docker daemon is running in Windows
				cpuPercentage = c.calcCpuPercentageWindows(stats)
				bytesRead = stats.StorageStats.ReadSizeBytes
				bytesWrite = stats.StorageStats.WriteSizeBytes
				networkRx, networkTx = c.calcNetworkIo(stats.Networks)
				memory = float64(stats.MemoryStats.PrivateWorkingSet)
			}

			containerResult := &resultDocker{
				Id:               ctr.ID[:10],
				Name:             ctr.Names[0],
				Image:            ctr.Image,
				SizeRw:           ctr.SizeRw,
				SizeRootFs:       ctr.SizeRootFs,
				State:            string(ctr.State), // running
				Status:           ctr.Status,        // Up 6 minutes
				CpuPercentage:    cpuPercentage,
				MemoryPercentage: memoryPercentage,
				NetworkRx:        networkRx,
				NetworkTx:        networkTx,
				DiskRead:         bytesRead,
				DiskWrite:        bytesWrite,
				MemoryUsed:       memory,
			}
			c.addInspectInfo(ctx, cli, ctr.ID, containerResult)

			// Return error back to check thread
			resultChan <- containerResult
		}(ctr)
	}

	// Check thread Wait for all results to come back and close the channels
	go func() {
		wg.Wait()
		close(errorChan)
		close(resultChan)
	}()

	for err := range errorChan {
		log.Errorln("Docker Check error: ", err)
	}

	for result := range resultChan {
		dockerResults = append(dockerResults, result)
	}

	return dockerResults, nil
}

// negotiateAPIVersion pings the Docker daemon to use the highest API version supported by the daemon and the client
// (API 1.40 / Docker 19.03 or newer) and returns the OS of the daemon (linux or windows).
// Without the negotiation the client would send its own API version, which older daemons reject with
// "client version ... is too new".
func (c *CheckDocker) negotiateAPIVersion(ctx context.Context, cli *client.Client) (string, error) {
	ping, err := cli.Ping(ctx, client.PingOptions{
		NegotiateAPIVersion: true,
	})
	if err != nil {
		return "", err
	}
	return ping.OSType, nil
}

// dockerStatsOSType returns the OS of the container. Docker 29 (API 1.52) and newer send it with the stats, for older
// daemons the OS of the daemon is used
func dockerStatsOSType(stats *container.StatsResponse, daemonOSType string) string {
	if stats.OSType != "" {
		return stats.OSType
	}
	return daemonOSType
}

// The restart counts of the containers, to detect restart loops. Stored outside of CheckDocker, as the checks get recreated on reload
var (
	dockerRestartsMutex  sync.Mutex
	dockerRestartSamples = map[string][]dockerRestartSample{}
)

type dockerRestartSample struct {
	time  time.Time
	count int
}

const (
	dockerRestartWindow      = time.Hour
	dockerHealthOutputLength = 1024
)

// addInspectInfo adds the healthcheck state and restart information of a container
func (c *CheckDocker) addInspectInfo(ctx context.Context, cli *client.Client, id string, result *resultDocker) {
	response, err := cli.ContainerInspect(ctx, id, client.ContainerInspectOptions{})
	if err != nil {
		log.Errorln("Docker Check error: ", err)
		return
	}
	inspect := response.Container

	result.RestartCount = inspect.RestartCount
	result.RestartsLastHour = dockerRestartsLastHour(id, inspect.RestartCount, time.Now())
	if inspect.State == nil {
		return
	}

	result.OOMKilled = inspect.State.OOMKilled
	result.ExitCode = inspect.State.ExitCode
	if startedAt, err := time.Parse(time.RFC3339Nano, inspect.State.StartedAt); err == nil && startedAt.Year() > 1 {
		result.StartedAt = startedAt.Unix()
	}
	if health := inspect.State.Health; health != nil {
		result.Health = string(health.Status)
		result.HealthFailingStreak = health.FailingStreak
		if len(health.Log) > 0 {
			last := health.Log[len(health.Log)-1]
			result.HealthOutput = truncateLogLineTo(strings.TrimSpace(last.Output), dockerHealthOutputLength)
		}
	}
}

// dockerRestartsLastHour returns the number of restarts within the last hour.
// The restart count is sampled with every check run, restarts before the agent started are not counted
func dockerRestartsLastHour(id string, restartCount int, now time.Time) int {
	dockerRestartsMutex.Lock()
	defer dockerRestartsMutex.Unlock()

	samples := append(dockerRestartSamples[id], dockerRestartSample{time: now, count: restartCount})
	windowStart := now.Add(-dockerRestartWindow)
	first := 0
	for first < len(samples)-1 && samples[first+1].time.Before(windowStart) {
		first++
	}
	samples = append([]dockerRestartSample(nil), samples[first:]...)
	dockerRestartSamples[id] = samples

	restarts := restartCount - samples[0].count
	if restarts < 0 {
		// The container was recreated with the same id (should not happen)
		return 0
	}
	return restarts
}

// forgetRemovedDockerContainers removes the restart samples of containers which were not seen within the last hour.
// A container in a restart loop can be missing in the container list for a short time
func forgetRemovedDockerContainers(existing map[string]bool, now time.Time) {
	dockerRestartsMutex.Lock()
	defer dockerRestartsMutex.Unlock()
	for id, samples := range dockerRestartSamples {
		if !existing[id] && (len(samples) == 0 || now.Sub(samples[len(samples)-1].time) > dockerRestartWindow) {
			delete(dockerRestartSamples, id)
		}
	}
}

// Configure the command or return false if the command was disabled
func (c *CheckDocker) Configure(config *config.Configuration) (bool, error) {
	return config.Docker, nil
}

func (c *CheckDocker) calcCpuPercentageUnix(previousCPU, previousSystem float64, stats *container.StatsResponse) float64 {
	// Credit to:
	// https://github.com/docker/cli/blob/e31e00585363e6a989e37fa92cf06481ea218344/cli/command/container/stats_helpers.go#L166-L183
	var cpuPercentage float64 = 0.0
	var cpuDelta float64 = float64(stats.CPUStats.CPUUsage.TotalUsage) - previousCPU
	var systemDelta float64 = float64(stats.CPUStats.SystemUsage) - previousSystem
	onlineCpus := stats.CPUStats.OnlineCPUs

	if onlineCpus == 0 {
		onlineCpus = uint32(len(stats.CPUStats.CPUUsage.PercpuUsage))
	}

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		//invalid operation: mismatched types float64 and uint32compiler
		cpuPercentage = (cpuDelta / systemDelta) * float64(onlineCpus) * 100.0
	}

	return cpuPercentage
}

func (c *CheckDocker) calcCpuPercentageWindows(stats *container.StatsResponse) float64 {
	// Credit to:
	// https://github.com/docker/cli/blob/e31e00585363e6a989e37fa92cf06481ea218344/cli/command/container/stats_helpers.go#L185-L199
	possibleIntervals := uint64(stats.Read.Sub(stats.PreRead).Nanoseconds())
	possibleIntervals = possibleIntervals / 100
	possibleIntervals = possibleIntervals * uint64(stats.NumProcs)

	// Used intervals
	usedIntervals := stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage

	// Percentage avoiding divide-by-zero
	var cpuPercentage float64 = 0.0
	if possibleIntervals > 0 {
		cpuPercentage = float64(usedIntervals) / float64(possibleIntervals) * 100.0
	}

	return cpuPercentage
}

func (c *CheckDocker) calcDiskIoUnix(diskio container.BlkioStats) (uint64, uint64) {
	// Credit to:
	// https://github.com/docker/cli/blob/e31e00585363e6a989e37fa92cf06481ea218344/cli/command/container/stats_helpers.go#L201-L215

	var bytesRead, bytesWrite uint64
	for _, blkioStatEntry := range diskio.IoServiceBytesRecursive {
		if len(blkioStatEntry.Op) == 0 {
			continue
		}
		switch blkioStatEntry.Op[0] {
		case 'r', 'R':
			bytesRead = bytesRead + blkioStatEntry.Value
		case 'w', 'W':
			bytesWrite = bytesWrite + blkioStatEntry.Value
		}
	}

	return bytesRead, bytesWrite
}

func (c *CheckDocker) calcNetworkIo(network map[string]container.NetworkStats) (float64, float64) {
	// Credit to:
	// https://github.com/docker/cli/blob/e31e00585363e6a989e37fa92cf06481ea218344/cli/command/container/stats_helpers.go#L217-L225
	var networkRx, networkTx float64

	for _, stats := range network {
		networkRx = networkRx + float64(stats.RxBytes)
		networkTx = networkTx + float64(stats.TxBytes)
	}
	return networkRx, networkTx
}

func (c *CheckDocker) calcMemoryUsageUnix(memory container.MemoryStats) float64 {
	// Credit to:
	// https://github.com/docker/cli/blob/e31e00585363e6a989e37fa92cf06481ea218344/cli/command/container/stats_helpers.go#L239-L249

	// cgroup v1
	// https://stackoverflow.com/a/2050629
	if inactiveFile, isCgroupV1 := memory.Stats["total_inactive_file"]; isCgroupV1 && inactiveFile < memory.Usage {
		return float64(memory.Usage - inactiveFile)
	}

	// cgroup v2
	if inactiveFile := memory.Stats["inactive_file"]; inactiveFile < memory.Usage {
		return float64(memory.Usage - inactiveFile)
	}

	return float64(memory.Usage)
}

func (c *CheckDocker) calcMemoryPercentageUnix(limit float64, noCache float64) float64 {
	// Credit to:
	// https://github.com/docker/cli/blob/e31e00585363e6a989e37fa92cf06481ea218344/cli/command/container/stats_helpers.go#L239-L249
	if limit != 0 {
		return noCache / limit * 100.0
	}
	return 0.0
}
