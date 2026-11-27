//go:build aix || linux
// +build aix linux

package checks

import (
	"reflect"
	"testing"
	"time"
)

// Output of "lssrc -a" (works without root)
func TestParseLssrc(t *testing.T) {
	output := `Subsystem         Group            PID          Status
 syslogd          ras              5833150      active
 aso                               6619620      active
 nfsd             nfs              6881618      active
 dfpd             tcpip                         inoperative
 cdromd                                         inoperative
 sendmail         mail                          stopping
`
	want := []*resultSrcSubsystem{
		{Name: "syslogd", Group: "ras", Pid: 5833150, Status: "active"},
		{Name: "aso", Pid: 6619620, Status: "active"},
		{Name: "nfsd", Group: "nfs", Pid: 6881618, Status: "active"},
		{Name: "dfpd", Group: "tcpip", Status: "inoperative"},
		{Name: "cdromd", Status: "inoperative"},
		{Name: "sendmail", Group: "mail", Status: "stopping"},
	}
	if got := parseLssrc(output); !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
	if got := parseLssrc(""); got == nil || len(got) != 0 {
		t.Errorf("expected an empty list, got %+v", got)
	}
}

// Output of "errpt -s <mmddHHMMyy>" (works without root, runs with TZ=UTC0 in the check)
func TestParseErrpt(t *testing.T) {
	location := time.FixedZone("CET", 3600)
	output := `IDENTIFIER TIMESTAMP  T C RESOURCE_NAME  DESCRIPTION
DC32143C   1103152026 P H ent4           PING TO REMOTE HOST FAILED
A6DF45AA   1103152126 I O RMCdaemon      The daemon is started.
C7E7362C   1014150226 T S cluster0       Node is heartbeating solely over disk or
EAA3D429   1102080026 U S LVDD           PHYSICAL PARTITION MARKED STALE
`
	entries := parseErrpt(output, location)
	if len(entries) != 4 {
		t.Fatalf("unexpected number of entries %d", len(entries))
	}

	// Newest entry first
	wantFirst := &resultErrptEntry{
		Identifier:  "A6DF45AA",
		Timestamp:   time.Date(2026, 11, 3, 15, 21, 0, 0, location).Unix(),
		Type:        "I",
		Class:       "O",
		Resource:    "RMCdaemon",
		Description: "The daemon is started.",
	}
	if !reflect.DeepEqual(entries[0], wantFirst) {
		t.Errorf("got %+v, want %+v", entries[0], wantFirst)
	}
	if entries[1].Identifier != "DC32143C" || entries[1].Description != "PING TO REMOTE HOST FAILED" || entries[1].Type != "P" || entries[1].Class != "H" {
		t.Errorf("unexpected second entry %+v", entries[1])
	}
	if entries[3].Identifier != "C7E7362C" || entries[3].Timestamp != time.Date(2026, 10, 14, 15, 2, 0, 0, location).Unix() {
		t.Errorf("unexpected oldest entry %+v", entries[3])
	}

	if got := parseErrpt("", location); got == nil || len(got) != 0 {
		t.Errorf("expected an empty list, got %+v", got)
	}
}

func TestErrptResultTruncated(t *testing.T) {
	entries := make([]*resultErrptEntry, errptMaxEntries+10)
	result := errptResult(entries, 24*time.Hour)
	if !result.Truncated || len(result.Entries) != errptMaxEntries || result.Window != 86400 {
		t.Errorf("unexpected result: truncated=%v entries=%d window=%d", result.Truncated, len(result.Entries), result.Window)
	}
	result = errptResult(entries[:3], time.Hour)
	if result.Truncated || len(result.Entries) != 3 || result.Window != 3600 {
		t.Errorf("unexpected result: truncated=%v entries=%d window=%d", result.Truncated, len(result.Entries), result.Window)
	}
}

// Output of "lsvg -L" (all volume groups) and "lsvg -o" (online volume groups), works without root
func TestParseLsvgList(t *testing.T) {
	got := parseLsvgList("rootvg\ndatavg\ncaavg_private\n\n")
	want := []string{"rootvg", "datavg", "caavg_private"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// Output of "lsvg -L rootvg" (works without root)
func TestParseLsvg(t *testing.T) {
	output := `VOLUME GROUP:       rootvg                   VG IDENTIFIER:  00f6db0a00004c000000016f5d1d2a3b
VG STATE:           active                   PP SIZE:        128 megabyte(s)
VG PERMISSION:      read/write               TOTAL PPs:      958 (122624 megabytes)
MAX LVs:            256                      FREE PPs:       460 (58880 megabytes)
LVs:                12                       USED PPs:       498 (63744 megabytes)
OPEN LVs:           10                       QUORUM:         2 (Enabled)
TOTAL PVs:          2                        VG DESCRIPTORS: 3
STALE PVs:          1                        STALE PPs:      7
ACTIVE PVs:         1                        AUTO ON:        no
MAX PPs per VG:     32512
MAX PPs per PV:     1016                     MAX PVs:        32
LTG size (Dynamic): 256 kilobyte(s)          AUTO SYNC:      no
`
	vg := &resultVolumeGroup{Name: "rootvg", Online: true}
	parseLsvg(output, vg)

	if vg.State != "active" || vg.Permission != "read/write" || vg.Quorum != "2 (Enabled)" {
		t.Errorf("unexpected state %+v", vg)
	}
	if vg.PPSizeMB != 128 || vg.TotalPPs != 958 || vg.FreePPs != 460 {
		t.Errorf("unexpected partitions %+v", vg)
	}
	if vg.TotalPVs != 2 || vg.ActivePVs != 1 || vg.StalePVs != 1 || vg.StalePPs != 7 {
		t.Errorf("unexpected physical volumes %+v", vg)
	}
	assertFloatNear(t, "used percent", vg.UsedPercent, 51.98)

	if size := lsvgPPSizeMB("1 gigabyte(s)"); size != 1024 {
		t.Errorf("unexpected PP size %d", size)
	}
}

// Output of "lsvg -L -l rootvg" (works without root)
func TestParseLsvgLogicalVolumes(t *testing.T) {
	output := `rootvg:
LV NAME             TYPE       LPs     PPs     PVs  LV STATE      MOUNT POINT
hd5                 boot       1       2       2    closed/syncd  N/A
hd4                 jfs2       8       16      2    open/syncd    /
hd2                 jfs2       39      78      2    open/stale    /usr
livedump            jfs2       2       4       2    open/syncd    /var/adm/ras/livedump
`
	want := []*resultLogicalVolume{
		{Name: "hd5", Type: "boot", LPs: 1, PPs: 2, PVs: 2, State: "closed/syncd", Mountpoint: "N/A"},
		{Name: "hd4", Type: "jfs2", LPs: 8, PPs: 16, PVs: 2, State: "open/syncd", Mountpoint: "/"},
		{Name: "hd2", Type: "jfs2", LPs: 39, PPs: 78, PVs: 2, State: "open/stale", Mountpoint: "/usr"},
		{Name: "livedump", Type: "jfs2", LPs: 2, PPs: 4, PVs: 2, State: "open/syncd", Mountpoint: "/var/adm/ras/livedump"},
	}
	if got := parseLsvgLogicalVolumes(output); !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

// Output of "lsvg -L -p rootvg" (works without root)
func TestParseLsvgPhysicalVolumes(t *testing.T) {
	output := `rootvg:
PV_NAME           PV STATE          TOTAL PPs   FREE PPs    FREE DISTRIBUTION
hdisk1            active            479         230         60..00..00..74..96
hdisk0            missing           479         230         60..00..00..74..96
`
	want := []*resultPhysicalVolume{
		{Name: "hdisk1", State: "active", TotalPPs: 479, FreePPs: 230},
		{Name: "hdisk0", State: "missing", TotalPPs: 479, FreePPs: 230},
	}
	if got := parseLsvgPhysicalVolumes(output); !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

// Output of "lspath" (works without root)
func TestParseLspath(t *testing.T) {
	output := `Enabled hdisk0   vscsi0
Enabled hdisk1   fscsi0
Enabled hdisk0   vscsi1
Failed  hdisk1   fscsi1
Missing hdisk1   fscsi2
`
	disks := parseLspath(output)
	if len(disks) != 2 {
		t.Fatalf("unexpected number of disks %d", len(disks))
	}
	want0 := &resultMultipathDisk{
		Name:         "hdisk0",
		Paths:        []*resultMultipathPath{{Parent: "vscsi0", Status: "Enabled"}, {Parent: "vscsi1", Status: "Enabled"}},
		EnabledPaths: 2,
	}
	if !reflect.DeepEqual(disks[0], want0) {
		t.Errorf("got %+v, want %+v", disks[0], want0)
	}
	if disks[1].Name != "hdisk1" || len(disks[1].Paths) != 3 || disks[1].EnabledPaths != 1 || disks[1].Paths[1].Status != "Failed" {
		t.Errorf("unexpected disk %+v", disks[1])
	}
}
