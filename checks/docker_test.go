package checks

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/moby/moby/api/types/container"
)

func TestChecksCheckDocker(t *testing.T) {

	check := &CheckDocker{}

	cr, err := check.Run(context.Background())
	if err != nil {
		// Windows
		if strings.Contains(err.Error(), "error during connect") {
			fmt.Println("Docker not installed or running on this system ???")
			t.SkipNow()
		}
		// macOS and Linux
		if strings.Contains(err.Error(), "Is the docker daemon running") {
			fmt.Println("Docker not installed or running on this system ???")
			t.SkipNow()
		}
		if strings.Contains(err.Error(), "Got permission denied while trying to connect to the Docker daemon socket") {
			fmt.Println(err)
			t.SkipNow()
		}

		t.Fatal(err)

	}

	results, ok := cr.([]*resultDocker)
	if !ok {
		t.Fatal("False type")

	}

	// Only running containers use memory, stopped containers report 0
	for _, result := range results {
		if result.State == "running" && result.MemoryUsed <= 0 {
			t.Errorf("running container %s uses no memory", result.Name)
		}
		if result.CpuPercentage < 0 || result.MemoryPercentage < 0 {
			t.Errorf("container %s: invalid values %+v", result.Name, result)
		}
	}
}

func TestDockerRestartsLastHour(t *testing.T) {
	dockerRestartsMutex.Lock()
	dockerRestartSamples = map[string][]dockerRestartSample{}
	dockerRestartsMutex.Unlock()

	start := time.Date(2026, 10, 13, 10, 0, 0, 0, time.UTC)
	// Restarts before the agent started are not counted
	if restarts := dockerRestartsLastHour("abc", 5, start); restarts != 0 {
		t.Errorf("expected 0 restarts on the first sample, got %d", restarts)
	}
	if restarts := dockerRestartsLastHour("abc", 7, start.Add(30*time.Minute)); restarts != 2 {
		t.Errorf("expected 2 restarts, got %d", restarts)
	}
	if restarts := dockerRestartsLastHour("abc", 8, start.Add(50*time.Minute)); restarts != 3 {
		t.Errorf("expected 3 restarts, got %d", restarts)
	}
	// The first sample left the window, the sample of 10:30 is the new reference
	if restarts := dockerRestartsLastHour("abc", 8, start.Add(100*time.Minute)); restarts != 1 {
		t.Errorf("expected 1 restart within the last hour, got %d", restarts)
	}

	// Missing containers are kept for an hour (restart loops)
	forgetRemovedDockerContainers(map[string]bool{}, start.Add(110*time.Minute))
	dockerRestartsMutex.Lock()
	_, kept := dockerRestartSamples["abc"]
	dockerRestartsMutex.Unlock()
	if !kept {
		t.Error("samples of a missing container must be kept for an hour")
	}
	forgetRemovedDockerContainers(map[string]bool{}, start.Add(170*time.Minute))
	dockerRestartsMutex.Lock()
	_, kept = dockerRestartSamples["abc"]
	dockerRestartsMutex.Unlock()
	if kept {
		t.Error("samples of a removed container must be deleted after an hour")
	}
}

func TestDockerStatsOSType(t *testing.T) {
	tests := []struct {
		statsOSType  string
		daemonOSType string
		expected     string
	}{
		{"windows", "windows", "windows"}, // Docker 29 and newer send the OS with the stats
		{"linux", "windows", "linux"},
		{"", "windows", "windows"}, // Older daemons: the OS of the daemon from the ping
		{"", "linux", "linux"},
		{"", "", ""},
	}
	for _, test := range tests {
		stats := &container.StatsResponse{OSType: test.statsOSType}
		if osType := dockerStatsOSType(stats, test.daemonOSType); osType != test.expected {
			t.Errorf("stats %q, daemon %q: expected %q, got %q", test.statsOSType, test.daemonOSType, test.expected, osType)
		}
	}
}
