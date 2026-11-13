package checks

import (
	"context"
	"errors"

	"github.com/shirou/gopsutil/v4/disk"
)

// errAixIostatDisabled is returned if AIX does not count the disk I/O, the agent does not change the setting
var errAixIostatDisabled = errors.New("the disk I/O history of AIX is disabled (sys0 attribute iostat=false), enable it with: chdev -l sys0 -a iostat=true")

// diskIoCounters returns the counters of each physical disk (perfstat_disk).
// AIX counts the disk I/O if the I/O history is enabled (sys0 attribute iostat) or if a tool enabled the collection,
// e.g. topas_nmon, which records the performance data on most systems. If nothing counts, all counters are 0.
func diskIoCounters(ctx context.Context) (map[string]disk.IOCountersStat, error) {
	disks, err := perfstatDisks()
	if err != nil {
		return nil, err
	}
	if !aixDisksHaveCounters(disks) && !aixIostatEnabled(ctx) {
		return nil, errAixIostatDisabled
	}
	return aixDiskIOCounters(disks), nil
}

// aixIostatEnabled returns true if the sys0 attribute iostat is true or can not be read
func aixIostatEnabled(ctx context.Context) bool {
	output, err := runAixCommand(ctx, "lsattr", "-E", "-l", "sys0", "-a", "iostat", "-F", "value")
	if err != nil {
		return true
	}
	enabled, err := parseAixIostatAttribute(output)
	if err != nil {
		return true
	}
	return enabled
}
