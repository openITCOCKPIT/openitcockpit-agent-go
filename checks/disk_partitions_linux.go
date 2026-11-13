package checks

import (
	"context"

	"github.com/shirou/gopsutil/v4/disk"
)

// readPartitions lists the file systems of /proc/self/mountinfo
func readPartitions(_ string) ([]disk.PartitionStat, error) {
	return disk.PartitionsWithContext(context.Background(), true)
}

// notRespondingMounts is not required on Linux, statfs of a not responding file system blocks
func notRespondingMounts() map[string]bool {
	return nil
}
