//go:build darwin || freebsd
// +build darwin freebsd

package checks

import (
	"github.com/shirou/gopsutil/v4/disk"
	"golang.org/x/sys/unix"
)

type mountFlagName struct {
	flag uint64
	name string
}

// readPartitions lists the file systems with getfsstat(MNT_NOWAIT). It returns the mount table of the kernel
// without asking the file systems, so it does not block if a network file system does not respond (MNT_WAIT would).
func readPartitions(_ string) ([]disk.PartitionStat, error) {
	stats, err := getfsstatNoWait()
	if err != nil {
		return nil, err
	}

	partitions := make([]disk.PartitionStat, 0, len(stats))
	for i := range stats {
		partitions = append(partitions, newPartition(&stats[i]))
	}
	return partitions, nil
}

func getfsstatNoWait() ([]unix.Statfs_t, error) {
	count, err := unix.Getfsstat(nil, unix.MNT_NOWAIT)
	if err != nil {
		return nil, err
	}
	stats := make([]unix.Statfs_t, count)
	count, err = unix.Getfsstat(stats, unix.MNT_NOWAIT)
	if err != nil {
		return nil, err
	}
	// The number of file systems can change between both calls
	return stats[:count], nil
}

func newPartition(stat *unix.Statfs_t) disk.PartitionStat {
	return disk.PartitionStat{
		Device:     unix.ByteSliceToString(stat.Mntfromname[:]),
		Mountpoint: unix.ByteSliceToString(stat.Mntonname[:]),
		Fstype:     unix.ByteSliceToString(stat.Fstypename[:]),
		Opts:       mountOptions(uint64(stat.Flags)),
	}
}

// mountOptions returns the mount options like gopsutil, e.g. ["rw", "nosuid", "local"]
func mountOptions(flags uint64) []string {
	options := []string{"rw"}
	if flags&unix.MNT_RDONLY != 0 {
		options = []string{"ro"}
	}
	for _, option := range mountFlagNames {
		if flags&option.flag != 0 {
			options = append(options, option.name)
		}
	}
	return options
}
