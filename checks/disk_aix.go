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

	mounts := parseAixMount(string(output))
	diskResults := make([]*resultDisk, 0, len(mounts))
	for _, mount := range mounts {
		if aixFstypeToIgnore[mount.Fstype] {
			continue
		}

		result := &resultDisk{}
		result.Disk.Device = mount.Device
		result.Disk.Mountpoint = mount.Mountpoint
		result.Disk.Fstype = mount.Fstype
		result.Disk.Opts = mount.Opts

		var stat unix.Statfs_t
		if err := unix.Statfs(mount.Mountpoint, &stat); err != nil {
			log.Errorln("DiskCheck: Error for ", mount.Mountpoint, err)
		} else {
			setAixDiskUsage(result, &stat)
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
