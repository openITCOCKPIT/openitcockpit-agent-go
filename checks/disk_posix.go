//go:build linux || darwin || freebsd
// +build linux darwin freebsd

package checks

import (
	"context"
	"errors"
	"strings"

	"github.com/shirou/gopsutil/v4/disk"
	log "github.com/sirupsen/logrus"
)

var devToIgnore = map[string]bool{
	"sysfs":       true,
	"proc":        true,
	"udev":        true,
	"devpts":      true,
	"devfs":       true,
	"tmpfs":       true,
	"securityfs":  true,
	"cgroup":      true,
	"cgroup2":     true,
	"pstore":      true,
	"debugfs":     true,
	"hugetlbfs":   true,
	"systemd-1":   true,
	"mqueue":      true,
	"none":        true,
	"sunrpc":      true,
	"nfsd":        true,
	"nsfs":        true,
	"fusectl":     true,
	"configfs":    true,
	"overlay":     true,
	"shm":         true,
	"tracefs":     true,
	"binfmt_misc": true,
}

// diskPartitionsProbe lists the file systems (readPartitions). The listing should not block,
// the probe only makes sure that the check always returns.
var diskPartitionsProbe diskProbe[[]disk.PartitionStat]

// diskUsageProbe reads the usage (statfs) of the file systems
var diskUsageProbe diskProbe[*disk.UsageStat]

var errDiskPartitionsHanging = errors.New("listing the file systems did not finish in time, a network file system (e.g. NFS) does not respond")

func readDiskUsage(mountpoint string) (*disk.UsageStat, error) {
	return disk.UsageWithContext(context.Background(), mountpoint)
}

// isIgnoredPartition returns true for pseudo file systems and ThinLinc sessions
func isIgnoredPartition(device disk.PartitionStat) bool {
	if devToIgnore[device.Device] {
		return true
	}
	// Skip over thinlinc sessions
	// https://github.com/openITCOCKPIT/openitcockpit-agent-go/issues/93
	// Should we also add device.Fstype == "nfs" to the if condition?
	if strings.Contains(strings.ToLower(device.Mountpoint), "thinlinc/sessions") {
		return true
	}

	// This will also ignore the main disk for LXC containers :(
	//if strings.HasPrefix(device.Device, "/dev/loop") {
	//	continue
	//}
	return false
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckDisk) Run(_ context.Context) (interface{}, error) {
	partitions := diskPartitionsProbe.readAll([]string{"partitions"}, readPartitions, diskUsageTimeout)[0]
	if partitions.Hanging {
		return nil, errDiskPartitionsHanging
	}
	if partitions.Err != nil {
		return nil, partitions.Err
	}

	var devices []disk.PartitionStat
	var mountpoints []string
	for _, device := range partitions.Value {
		if isIgnoredPartition(device) {
			continue
		}
		devices = append(devices, device)
		mountpoints = append(mountpoints, device.Mountpoint)
	}

	usages := diskUsageProbe.readAll(mountpoints, readDiskUsage, diskUsageTimeout)
	notResponding := notRespondingMounts()
	diskResults := make([]*resultDisk, 0, len(devices))
	for i, device := range devices {
		result := &resultDisk{}
		result.Disk.Device = device.Device
		result.Disk.Mountpoint = device.Mountpoint
		result.Disk.Fstype = device.Fstype
		result.Disk.Opts = device.Opts
		setDiskUsage(result, usages[i], notResponding[device.Mountpoint])
		diskResults = append(diskResults, result)
	}

	return diskResults, nil
}

// setDiskUsage sets the usage or hanging, if statfs did not return in time or the file system is marked as not responding
func setDiskUsage(result *resultDisk, probe diskProbeResult[*disk.UsageStat], notResponding bool) {
	if probe.Hanging || notResponding {
		log.Errorln("DiskCheck: ", result.Disk.Mountpoint, " does not respond")
		result.Hanging = true
		return
	}
	if probe.Err != nil {
		log.Errorln("DiskCheck: Error for ", result.Disk.Mountpoint, probe.Err)
	}
	usage := probe.Value
	if usage == nil {
		return
	}
	result.Usage.Total = usage.Total
	result.Usage.Used = usage.Used
	result.Usage.Free = usage.Free
	result.Usage.Percent = usage.UsedPercent
	result.Usage.InodesTotal = usage.InodesTotal
	result.Usage.InodesUsed = usage.InodesUsed
	result.Usage.InodesFree = usage.InodesFree
	result.Usage.InodesPercent = usage.InodesUsedPercent
}
