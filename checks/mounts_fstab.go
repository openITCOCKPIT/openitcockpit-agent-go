//go:build linux || freebsd
// +build linux freebsd

package checks

import (
	"context"

	"github.com/shirou/gopsutil/v4/disk"
)

const fstabPath = "/etc/fstab"

// parseMountTable returns the filesystems of the fstab which should be mounted
func parseMountTable(data []byte) []*resultFstabEntry {
	return parseFstab(data)
}

func currentMountpoints(ctx context.Context) ([]string, error) {
	// all = true also returns pseudo filesystems like autofs (systemd automount) and network filesystems
	partitions, err := disk.PartitionsWithContext(ctx, true)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(partitions))
	for _, partition := range partitions {
		result = append(result, partition.Mountpoint)
	}
	return result, nil
}
