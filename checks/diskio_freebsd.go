package checks

import (
	"context"

	"github.com/shirou/gopsutil/v4/disk"
)

// diskIoCounters returns the counters of each disk
func diskIoCounters(ctx context.Context) (map[string]disk.IOCountersStat, error) {
	return disk.IOCountersWithContext(ctx)
}
