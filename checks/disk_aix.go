package checks

import (
	"context"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// Pseudo file systems of AIX
var aixFstypeToIgnore = map[string]bool{
	"procfs": true, // /proc
	"ahafs":  true, // Autonomic Health Advisor File System (/aha)
	"namefs": true,
}

// diskUsageProbe reads statfs of the file systems, see diskProbe
var diskUsageProbe diskProbe[*unix.Statfs_t]

func readAixStatfs(mountpoint string) (*unix.Statfs_t, error) {
	stat := &unix.Statfs_t{}
	if err := unix.Statfs(mountpoint, stat); err != nil {
		return nil, err
	}
	return stat, nil
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckDisk) Run(ctx context.Context) (interface{}, error) {
	// The mount command only reads the mount table of the kernel, it does not access the file systems
	output, err := exec.CommandContext(ctx, "mount").Output()
	if err != nil {
		return nil, err
	}

	var mounts []aixMount
	var mountpoints []string
	for _, mount := range parseAixMount(string(output)) {
		if aixFstypeToIgnore[mount.Fstype] {
			continue
		}
		mounts = append(mounts, mount)
		mountpoints = append(mountpoints, mount.Mountpoint)
	}

	stats := diskUsageProbe.readAll(mountpoints, readAixStatfs, diskUsageTimeout)
	diskResults := make([]*resultDisk, 0, len(mounts))
	for i, mount := range mounts {
		result := &resultDisk{}
		result.Disk.Device = mount.Device
		result.Disk.Mountpoint = mount.Mountpoint
		result.Disk.Fstype = mount.Fstype
		result.Disk.Opts = mount.Opts

		switch {
		case stats[i].Hanging:
			log.Errorln("DiskCheck: ", mount.Mountpoint, " does not respond")
			result.Hanging = true
		case stats[i].Err != nil:
			log.Errorln("DiskCheck: Error for ", mount.Mountpoint, stats[i].Err)
		default:
			setAixDiskUsage(result, stats[i].Value)
		}
		diskResults = append(diskResults, result)
	}
	return diskResults, nil
}

func setAixDiskUsage(result *resultDisk, stat *unix.Statfs_t) {
	blockSize := stat.Bsize
	result.Usage.Total = stat.Blocks * blockSize
	result.Usage.Free = stat.Bavail * blockSize
	if stat.Blocks > stat.Bfree {
		result.Usage.Used = (stat.Blocks - stat.Bfree) * blockSize
	}
	// Like df: used / (used + available)
	if result.Usage.Used+result.Usage.Free > 0 {
		result.Usage.Percent = float64(result.Usage.Used) / float64(result.Usage.Used+result.Usage.Free) * 100
	}

	result.Usage.InodesTotal = stat.Files
	result.Usage.InodesFree = stat.Ffree
	if stat.Files > stat.Ffree {
		result.Usage.InodesUsed = stat.Files - stat.Ffree
	}
	if stat.Files > 0 {
		result.Usage.InodesPercent = float64(result.Usage.InodesUsed) / float64(stat.Files) * 100
	}
}
