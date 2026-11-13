//go:build aix || linux
// +build aix linux

package checks

import (
	"reflect"
	"testing"
)

func TestParseLsps(t *testing.T) {
	output := `Page Space      Physical Volume   Volume Group    Size %Used   Active    Auto    Type   Chksum
hd6             hdisk0            rootvg        4096MB     2     yes     yes      lv       0
paging00        hdisk1            pagingvg      2048MB    10     yes      no      lv       0
paging00        hdisk2            pagingvg      2048MB    30     yes      no      lv       0
paging01        hdisk3            pagingvg         1GB     0      no      no      lv       0
`
	pagingSpaces := parseLsps(output)
	if len(pagingSpaces) != 3 {
		t.Fatalf("unexpected paging spaces %+v", pagingSpaces)
	}
	mb := uint64(1024 * 1024)
	want := &resultPagingSpace{
		Name: "hd6", PhysicalVolumes: []string{"hdisk0"}, VolumeGroup: "rootvg", Type: "lv",
		Size: 4096 * mb, Used: uint64(float64(4096*mb) * 0.02), Percent: 2, Active: true, Auto: true,
	}
	if !reflect.DeepEqual(pagingSpaces[0], want) {
		t.Errorf("got %+v, want %+v", pagingSpaces[0], want)
	}

	// Spans two physical volumes: 10% of 2 GB + 30% of 2 GB = 20% of 4 GB
	paging00 := pagingSpaces[1]
	if !reflect.DeepEqual(paging00.PhysicalVolumes, []string{"hdisk1", "hdisk2"}) || paging00.Size != 4096*mb ||
		paging00.Percent < 19.99 || paging00.Percent > 20.01 || !paging00.Active || paging00.Auto {
		t.Errorf("unexpected paging space %+v", paging00)
	}

	paging01 := pagingSpaces[2]
	if paging01.Size != 1024*mb || paging01.Active || paging01.Used != 0 {
		t.Errorf("unexpected paging space %+v", paging01)
	}

	// Older AIX versions have no Chksum column
	old := parseLsps(`Page Space  Physical Volume   Volume Group    Size   %Used  Active  Auto  Type
hd6         hdisk0            rootvg         512MB     1     yes   yes    lv
`)
	if len(old) != 1 || old[0].Size != 512*mb || old[0].Type != "lv" {
		t.Errorf("unexpected paging spaces %+v", old)
	}
}
