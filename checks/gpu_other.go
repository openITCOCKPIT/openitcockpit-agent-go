//go:build !linux && !windows

package checks

import (
	"context"
	"time"
)

// countNvidiaPciDevices is only implemented on Linux
func countNvidiaPciDevices() int {
	return -1
}

// readGpuDriverLogs is only implemented on Linux and Windows
func readGpuDriverLogs(ctx context.Context, since time.Time) *gpuDriverLogs {
	return &gpuDriverLogs{}
}

// Only nvidia-smi is supported on this platform
func hasPlatformGpus() bool {
	return false
}

func collectPlatformGpus(ctx context.Context, existing []*resultGpu) []*resultGpu {
	return nil
}
