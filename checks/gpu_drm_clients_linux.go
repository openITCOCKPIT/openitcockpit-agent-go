package checks

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// The DRM drivers report the busy time of the GPU engines and the memory per client (open file of /dev/dri/*) in
// /proc/<pid>/fdinfo/<fd>, like intel_gpu_top and nvtop use it:
// https://docs.kernel.org/gpu/drm-usage-stats.html
// The utilization of an engine is the busy time of all clients between two check runs. Clients which were started or
// stopped between the runs are not counted.

var (
	gpuProcPath = "/proc"
	gpuNow      = time.Now

	drmClientsMutex    sync.Mutex
	drmClientsPrevious *drmClientsSample
)

type drmClient struct {
	pdev        string            // PCI bus id, e.g. 0000:00:02.0
	pid         int               // Process which owns the file
	engines     map[string]uint64 // drm-engine-<engine>: busy time in ns (i915, amdgpu)
	cycles      map[string]uint64 // drm-cycles-<engine>: busy cycles (xe)
	totalCycles map[string]uint64 // drm-total-cycles-<engine>: GPU timestamp in cycles (xe)
	capacity    map[string]uint64 // drm-engine-capacity-<engine>: number of engines of this type, default 1
	vram        uint64            // Sum of drm-resident-<region> of the memory of the GPU (i915: local0, xe: vram0) in bytes
	hasVram     bool
}

type drmClientsSample struct {
	time    time.Time
	clients map[string]*drmClient // key: pdev/drm-client-id
}

// drmUsage is the usage of a GPU calculated from the DRM clients
type drmUsage struct {
	utilization *float64 // Busiest engine, like the Task Manager of Windows
	encoder     *float64
	decoder     *float64
	processes   int
	vramUsed    *uint64 // Memory of the GPU used by the clients (Intel), nil if no client reports it
}

// readDrmClients reads the DRM clients of the drivers from /proc/<pid>/fdinfo. Only the files which are links to
// /dev/dri/* are read. Several files of the same client (drm-client-id) are counted once.
func readDrmClients(drivers ...string) *drmClientsSample {
	sample := &drmClientsSample{time: gpuNow(), clients: map[string]*drmClient{}}
	processes, err := os.ReadDir(gpuProcPath)
	if err != nil {
		return sample
	}
	for _, process := range processes {
		pid, err := strconv.Atoi(process.Name())
		if err != nil {
			continue
		}
		fdPath := filepath.Join(gpuProcPath, process.Name(), "fd")
		fds, err := os.ReadDir(fdPath)
		if err != nil {
			continue
		}
		for _, fd := range fds {
			target, err := os.Readlink(filepath.Join(fdPath, fd.Name()))
			if err != nil || !strings.HasPrefix(target, "/dev/dri/") {
				continue
			}
			client, id := parseDrmFdinfo(filepath.Join(gpuProcPath, process.Name(), "fdinfo", fd.Name()), drivers)
			if client == nil {
				continue
			}
			client.pid = pid
			key := client.pdev + "/" + id
			if _, ok := sample.clients[key]; !ok {
				sample.clients[key] = client
			}
		}
	}
	return sample
}

// parseDrmFdinfo returns the client and its drm-client-id, nil if the file is no client of the drivers
func parseDrmFdinfo(path string, drivers []string) (*drmClient, string) {
	file, err := os.Open(path)
	if err != nil {
		return nil, ""
	}
	defer file.Close()

	client := &drmClient{engines: map[string]uint64{}, cycles: map[string]uint64{}, totalCycles: map[string]uint64{}, capacity: map[string]uint64{}}
	var driver, id string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key, value, ok := strings.Cut(scanner.Text(), ":")
		if !ok || !strings.HasPrefix(key, "drm-") {
			continue
		}
		value = strings.TrimSpace(value)
		switch {
		case key == "drm-driver":
			driver = value
		case key == "drm-pdev":
			client.pdev = value
		case key == "drm-client-id":
			id = value
		case strings.HasPrefix(key, "drm-engine-capacity-"):
			client.capacity[strings.TrimPrefix(key, "drm-engine-capacity-")] = parseDrmNumber(value)
		case strings.HasPrefix(key, "drm-engine-"):
			client.engines[strings.TrimPrefix(key, "drm-engine-")] = parseDrmNumber(value)
		case strings.HasPrefix(key, "drm-total-cycles-"):
			client.totalCycles[strings.TrimPrefix(key, "drm-total-cycles-")] = parseDrmNumber(value)
		case strings.HasPrefix(key, "drm-cycles-"):
			client.cycles[strings.TrimPrefix(key, "drm-cycles-")] = parseDrmNumber(value)
		case strings.HasPrefix(key, "drm-resident-") && isDrmVramRegion(strings.TrimPrefix(key, "drm-resident-")):
			client.vram += parseDrmBytes(value)
			client.hasVram = true
		}
	}
	if id == "" || client.pdev == "" {
		return nil, ""
	}
	for _, d := range drivers {
		if d == driver {
			return client, id
		}
	}
	return nil, ""
}

// parseDrmNumber returns the number of values like "277160 ns" or "2"
func parseDrmNumber(value string) uint64 {
	number, _, _ := strings.Cut(value, " ")
	n, _ := strconv.ParseUint(number, 10, 64)
	return n
}

// parseDrmBytes returns the bytes of memory values like "506440 KiB", "12 MiB" or "0"
func parseDrmBytes(value string) uint64 {
	number, unit, _ := strings.Cut(strings.TrimSpace(value), " ")
	n, _ := strconv.ParseUint(number, 10, 64)
	switch unit {
	case "KiB":
		return n * 1024
	case "MiB":
		return n * 1024 * 1024
	case "GiB":
		return n * 1024 * 1024 * 1024
	}
	return n
}

// isDrmVramRegion returns true for the memory regions of the GPU itself: local0 (i915), vram0 (xe). Stolen memory is
// part of the system memory.
func isDrmVramRegion(region string) bool {
	return strings.HasPrefix(region, "local") || strings.HasPrefix(region, "vram")
}

// drmUsageByPdev calculates the usage of the GPUs (by PCI bus id) between two samples. Without a previous sample
// only the number of processes is known.
func drmUsageByPdev(previous, current *drmClientsSample) map[string]*drmUsage {
	result := map[string]*drmUsage{}
	get := func(pdev string) *drmUsage {
		if result[pdev] == nil {
			result[pdev] = &drmUsage{}
		}
		return result[pdev]
	}

	pids := map[string]map[int]bool{}
	for _, client := range current.clients {
		if pids[client.pdev] == nil {
			pids[client.pdev] = map[int]bool{}
		}
		pids[client.pdev][client.pid] = true
	}
	for pdev, p := range pids {
		get(pdev).processes = len(p)
	}
	// Buffers which are shared between clients are counted for every client
	for _, client := range current.clients {
		if !client.hasVram {
			continue
		}
		usage := get(client.pdev)
		if usage.vramUsed == nil {
			usage.vramUsed = new(uint64)
		}
		*usage.vramUsed += client.vram
	}

	if previous == nil {
		return result
	}
	elapsed := current.time.Sub(previous.time)
	if elapsed <= 0 {
		return result
	}

	// Busy fraction of each engine type, summed over the clients
	type engineKey struct{ pdev, engine string }
	busy := map[engineKey]float64{}
	capacity := map[engineKey]uint64{}
	for key, client := range current.clients {
		before := previous.clients[key]
		if before == nil {
			continue
		}
		for engine, ns := range client.engines {
			if old, ok := before.engines[engine]; ok && ns >= old {
				busy[engineKey{client.pdev, engine}] += float64(ns-old) / float64(elapsed.Nanoseconds())
			}
		}
		for engine, cycles := range client.cycles {
			old, ok := before.cycles[engine]
			total, oldTotal := client.totalCycles[engine], before.totalCycles[engine]
			if ok && cycles >= old && total > oldTotal {
				busy[engineKey{client.pdev, engine}] += float64(cycles-old) / float64(total-oldTotal)
			}
		}
		for engine, c := range client.capacity {
			capacity[engineKey{client.pdev, engine}] = c
		}
	}

	maximum := func(current *float64, value float64) *float64 {
		value = min(value, 100)
		if current == nil || value > *current {
			return &value
		}
		return current
	}
	for key, fraction := range busy {
		c := capacity[key]
		if c == 0 {
			c = 1
		}
		percent := fraction / float64(c) * 100
		usage := get(key.pdev)
		usage.utilization = maximum(usage.utilization, percent)
		switch engine := key.engine; {
		case engine == "video" || engine == "vcs":
			// Intel: the video engines decode and encode
			usage.encoder = maximum(usage.encoder, percent)
			usage.decoder = maximum(usage.decoder, percent)
		case strings.HasPrefix(engine, "enc"):
			usage.encoder = maximum(usage.encoder, percent)
		case strings.HasPrefix(engine, "dec"):
			usage.decoder = maximum(usage.decoder, percent)
		}
	}
	return result
}

// readDrmUsage reads the DRM clients and returns the usage since the last call. hasPrevious is false on the first
// call, then the utilization is unknown.
func readDrmUsage(drivers ...string) (usage map[string]*drmUsage, hasPrevious bool) {
	drmClientsMutex.Lock()
	defer drmClientsMutex.Unlock()
	current := readDrmClients(drivers...)
	usage = drmUsageByPdev(drmClientsPrevious, current)
	hasPrevious = drmClientsPrevious != nil
	drmClientsPrevious = current
	return usage, hasPrevious
}
