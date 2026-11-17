//go:build aix || linux
// +build aix linux

// Parsers for the output of AIX commands (lssrc, errpt, lsvg, lspath). The file is also built on Linux,
// so the unit tests run in the CI.

package checks

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// resultSrcSubsystem is a subsystem of the AIX System Resource Controller (lssrc -a)
type resultSrcSubsystem struct {
	Name   string `json:"name"`   // e.g.: sshd
	Group  string `json:"group"`  // e.g.: ssh, empty if the subsystem has no group
	Pid    int64  `json:"pid"`    // 0 if the subsystem is not running
	Status string `json:"status"` // e.g.: active, inoperative, stopping
}

// parseLssrc parses the output of lssrc -a:
//
//	Subsystem         Group            PID          Status
//	 sshd             ssh              8061234      active
//	 aso                               6619620      active
//	 dfpd             tcpip                         inoperative
//	 cdromd                                         inoperative
func parseLssrc(output string) []*resultSrcSubsystem {
	subsystems := []*resultSrcSubsystem{}
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 || len(fields) > 4 || fields[0] == "Subsystem" {
			continue
		}

		subsystem := &resultSrcSubsystem{
			Name:   fields[0],
			Status: fields[len(fields)-1],
		}
		// The columns between the name and the status are the optional group and the optional PID
		for _, field := range fields[1 : len(fields)-1] {
			pid, err := strconv.ParseInt(field, 10, 64)
			if err == nil {
				subsystem.Pid = pid
			} else {
				subsystem.Group = field
			}
		}
		subsystems = append(subsystems, subsystem)
	}
	return subsystems
}

// resultErrptEntry is an entry of the AIX error log (errpt)
type resultErrptEntry struct {
	Identifier  string `json:"identifier"`  // Error identifier, e.g.: DC32143C
	Timestamp   int64  `json:"timestamp"`   // Unix timestamp
	Type        string `json:"type"`        // P (permanent), T (temporary), I (informational), U (unknown)
	Class       string `json:"class"`       // H (hardware), S (software), O (operator), U (undetermined)
	Resource    string `json:"resource"`    // e.g.: ent4, hdisk0
	Description string `json:"description"` // e.g.: PING TO REMOTE HOST FAILED
}

type resultErrpt struct {
	Window    int64               `json:"window"`    // Time window of the entries in seconds
	Entries   []*resultErrptEntry `json:"entries"`   // Newest entry first
	Truncated bool                `json:"truncated"` // More entries than errptMaxEntries
}

// errptMaxEntries limits the size of the check result
const errptMaxEntries = 500

// errptTimeFormat is the format of the timestamps of errpt (MMDDhhmmYY) and of errpt -s
const errptTimeFormat = "0102150406"

// parseErrpt parses the summary of errpt, the timestamps are in the local time zone of the system:
//
//	IDENTIFIER TIMESTAMP  T C RESOURCE_NAME  DESCRIPTION
//	DC32143C   1103152026 P H ent4           PING TO REMOTE HOST FAILED
func parseErrpt(output string, location *time.Location) []*resultErrptEntry {
	entries := []*resultErrptEntry{}
	for _, line := range strings.Split(output, "\n") {
		entry, ok := parseErrptLine(line, location)
		if ok {
			entries = append(entries, entry)
		}
	}
	// errpt lists the newest entry first, the order is kept for entries with the same minute
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Timestamp > entries[j].Timestamp
	})
	return entries
}

func parseErrptLine(line string, location *time.Location) (*resultErrptEntry, bool) {
	fields := strings.Fields(line)
	if len(fields) < 5 || fields[0] == "IDENTIFIER" {
		return nil, false
	}
	timestamp, err := time.ParseInLocation(errptTimeFormat, fields[1], location)
	if err != nil {
		return nil, false
	}

	entry := &resultErrptEntry{
		Identifier: fields[0],
		Timestamp:  timestamp.Unix(),
		Type:       fields[2],
		Class:      fields[3],
		Resource:   fields[4],
	}
	if len(fields) > 5 {
		entry.Description = strings.Join(fields[5:], " ")
	}
	return entry, true
}

// errptResult limits the entries to errptMaxEntries
func errptResult(entries []*resultErrptEntry, window time.Duration) *resultErrpt {
	result := &resultErrpt{
		Window:  int64(window.Seconds()),
		Entries: entries,
	}
	if len(entries) > errptMaxEntries {
		result.Entries = entries[:errptMaxEntries]
		result.Truncated = true
	}
	return result
}

// resultVolumeGroup is an AIX LVM volume group
type resultVolumeGroup struct {
	Name   string `json:"name"`   // e.g.: rootvg
	Online bool   `json:"online"` // false if the volume group is varied off, e.g. on the standby node of a cluster

	// The following values are only available for online volume groups
	State       string                  `json:"state"`        // e.g.: active
	Permission  string                  `json:"permission"`   // e.g.: read/write
	PPSizeMB    int64                   `json:"pp_size_mb"`   // Size of a physical partition
	TotalPPs    int64                   `json:"total_pps"`    // Physical partitions
	FreePPs     int64                   `json:"free_pps"`     // Free physical partitions
	TotalPVs    int64                   `json:"total_pvs"`    // Physical volumes
	ActivePVs   int64                   `json:"active_pvs"`   // Active physical volumes
	StalePVs    int64                   `json:"stale_pvs"`    // Physical volumes with stale partitions
	StalePPs    int64                   `json:"stale_pps"`    // Stale physical partitions (mirror copies which are not in sync)
	Quorum      string                  `json:"quorum"`       // e.g.: 2 (Enabled)
	LVs         []*resultLogicalVolume  `json:"lvs"`          // Logical volumes
	PVs         []*resultPhysicalVolume `json:"pvs"`          // Physical volumes
	ReadErrors  []string                `json:"read_errors"`  // Errors of lsvg
	UsedPercent float64                 `json:"used_percent"` // Used physical partitions in percent
}

// resultLogicalVolume is a logical volume of lsvg -l
type resultLogicalVolume struct {
	Name       string `json:"name"`       // e.g.: hd4
	Type       string `json:"type"`       // e.g.: jfs2, paging, boot
	LPs        int64  `json:"lps"`        // Logical partitions
	PPs        int64  `json:"pps"`        // Physical partitions (LPs * number of copies)
	PVs        int64  `json:"pvs"`        // Physical volumes
	State      string `json:"state"`      // e.g.: open/syncd, closed/syncd, open/stale
	Mountpoint string `json:"mountpoint"` // e.g.: /, N/A
}

// resultPhysicalVolume is a physical volume of lsvg -p
type resultPhysicalVolume struct {
	Name     string `json:"name"`      // e.g.: hdisk0
	State    string `json:"state"`     // e.g.: active, missing, removed
	TotalPPs int64  `json:"total_pps"` // Physical partitions
	FreePPs  int64  `json:"free_pps"`  // Free physical partitions
}

// parseLsvgList parses the output of lsvg and lsvg -o (one volume group per line)
func parseLsvgList(output string) []string {
	names := []string{}
	for _, line := range strings.Split(output, "\n") {
		name := strings.TrimSpace(line)
		if name != "" && !strings.Contains(name, " ") {
			names = append(names, name)
		}
	}
	return names
}

// parseLsvg parses the output of lsvg <vg>. Every line contains two "KEY: value" pairs:
//
//	VG STATE:           active                   PP SIZE:        128 megabyte(s)
//	TOTAL PVs:          2                        VG DESCRIPTORS: 3
func parseLsvg(output string, vg *resultVolumeGroup) {
	values := parseLsvgValues(output)
	vg.State = values["VG STATE"]
	vg.Permission = values["VG PERMISSION"]
	vg.PPSizeMB = lsvgPPSizeMB(values["PP SIZE"])
	vg.TotalPPs = leadingInt(values["TOTAL PPs"])
	vg.FreePPs = leadingInt(values["FREE PPs"])
	vg.TotalPVs = leadingInt(values["TOTAL PVs"])
	vg.ActivePVs = leadingInt(values["ACTIVE PVs"])
	vg.StalePVs = leadingInt(values["STALE PVs"])
	vg.StalePPs = leadingInt(values["STALE PPs"])
	vg.Quorum = values["QUORUM"]
	if vg.TotalPPs > 0 {
		vg.UsedPercent = float64(vg.TotalPPs-vg.FreePPs) / float64(vg.TotalPPs) * 100
	}
}

// lsvgValueColumn is the column of the second "KEY: value" pair of lsvg <vg>
const lsvgValueColumn = 45

func parseLsvgValues(output string) map[string]string {
	values := map[string]string{}
	for _, line := range strings.Split(output, "\n") {
		if len(line) > lsvgValueColumn {
			addLsvgValue(values, line[:lsvgValueColumn])
			addLsvgValue(values, line[lsvgValueColumn:])
		} else {
			addLsvgValue(values, line)
		}
	}
	return values
}

func addLsvgValue(values map[string]string, pair string) {
	key, value, found := strings.Cut(pair, ":")
	if !found {
		return
	}
	key = strings.TrimSpace(key)
	if key != "" {
		values[key] = strings.TrimSpace(value)
	}
}

// lsvgPPSizeMB returns the size of a physical partition in MB, e.g. "128 megabyte(s)" or "1 gigabyte(s)"
func lsvgPPSizeMB(value string) int64 {
	size := leadingInt(value)
	if strings.Contains(value, "gigabyte") {
		return size * 1024
	}
	return size
}

// leadingInt returns the number at the beginning of a value like "128 megabyte(s)" or "79 (5056 megabytes)"
func leadingInt(value string) int64 {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return 0
	}
	number, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0
	}
	return number
}

// parseLsvgLogicalVolumes parses the output of lsvg -l <vg>:
//
//	rootvg:
//	LV NAME             TYPE       LPs     PPs     PVs  LV STATE      MOUNT POINT
//	hd4                 jfs2       8       16      2    open/syncd    /
func parseLsvgLogicalVolumes(output string) []*resultLogicalVolume {
	lvs := []*resultLogicalVolume{}
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 7 || fields[0] == "LV" {
			continue
		}
		lps, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			continue
		}
		lvs = append(lvs, &resultLogicalVolume{
			Name:       fields[0],
			Type:       fields[1],
			LPs:        lps,
			PPs:        leadingInt(fields[3]),
			PVs:        leadingInt(fields[4]),
			State:      fields[5],
			Mountpoint: strings.Join(fields[6:], " "),
		})
	}
	return lvs
}

// parseLsvgPhysicalVolumes parses the output of lsvg -p <vg>:
//
//	rootvg:
//	PV_NAME           PV STATE          TOTAL PPs   FREE PPs    FREE DISTRIBUTION
//	hdisk0            active            479         230         60..00..00..74..96
func parseLsvgPhysicalVolumes(output string) []*resultPhysicalVolume {
	pvs := []*resultPhysicalVolume{}
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 4 || fields[0] == "PV_NAME" {
			continue
		}
		total, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			continue
		}
		pvs = append(pvs, &resultPhysicalVolume{
			Name:     fields[0],
			State:    fields[1],
			TotalPPs: total,
			FreePPs:  leadingInt(fields[3]),
		})
	}
	return pvs
}

// resultMultipathDisk contains the paths of a disk (AIX MPIO, lspath)
type resultMultipathDisk struct {
	Name         string                 `json:"name"`          // e.g.: hdisk0
	Paths        []*resultMultipathPath `json:"paths"`         // All paths of the disk
	EnabledPaths int                    `json:"enabled_paths"` // Paths with the status Enabled
}

type resultMultipathPath struct {
	Parent string `json:"parent"` // Adapter, e.g.: fscsi0, vscsi1
	Status string `json:"status"` // Enabled, Disabled, Failed, Defined, Missing, Detected
}

// parseLspath parses the output of lspath (status, disk, parent):
//
//	Enabled hdisk0   vscsi0
//	Failed  hdisk0   vscsi1
func parseLspath(output string) []*resultMultipathDisk {
	disksByName := map[string]*resultMultipathDisk{}
	disks := []*resultMultipathDisk{}
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			continue
		}
		disk, exists := disksByName[fields[1]]
		if !exists {
			disk = &resultMultipathDisk{Name: fields[1], Paths: []*resultMultipathPath{}}
			disksByName[fields[1]] = disk
			disks = append(disks, disk)
		}
		disk.Paths = append(disk.Paths, &resultMultipathPath{Parent: fields[2], Status: fields[0]})
		if fields[0] == "Enabled" {
			disk.EnabledPaths++
		}
	}
	return disks
}

// commandError returns an error with the output of the command (e.g. the AIX message like 0516-010 ...)
func commandError(name string, err error, output string) error {
	output = strings.TrimSpace(output)
	if output == "" {
		return fmt.Errorf("%s: %w", name, err)
	}
	return fmt.Errorf("%s: %w: %s", name, err, output)
}
