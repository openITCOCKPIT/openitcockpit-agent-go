//go:build aix || linux
// +build aix linux

// Parser for the paging spaces of AIX (lsps -a). The file is also built on Linux, so the unit tests run in the CI.

package checks

import (
	"strconv"
	"strings"
)

// resultPagingSpace is a paging space of AIX (lsps -a)
type resultPagingSpace struct {
	Name            string   `json:"name"`             // e.g.: hd6, paging00
	PhysicalVolumes []string `json:"physical_volumes"` // e.g.: ["hdisk0"], a paging space can span several physical volumes
	VolumeGroup     string   `json:"volume_group"`     // e.g.: rootvg
	Type            string   `json:"type"`             // lv (logical volume) or nfs
	Size            uint64   `json:"size"`             // Bytes
	Used            uint64   `json:"used"`             // Bytes (calculated from the percentage of lsps)
	Percent         float64  `json:"percent"`          // Used in percent
	Active          bool     `json:"active"`           // The paging space is in use
	Auto            bool     `json:"auto"`             // The paging space is activated at boot
}

// lspsColumns is the minimum number of columns of lsps -a (Chksum is missing on older AIX versions)
const lspsColumns = 8

// parseLsps parses the output of lsps -a. A paging space which spans several physical volumes
// has a line for each physical volume.
//
//	Page Space      Physical Volume   Volume Group    Size %Used   Active    Auto    Type   Chksum
//	hd6             hdisk0            rootvg        4096MB     2     yes     yes      lv       0
func parseLsps(output string) []*resultPagingSpace {
	pagingSpaces := []*resultPagingSpace{}
	byName := map[string]*resultPagingSpace{}
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) < lspsColumns || fields[0] == "Page" {
			continue
		}
		size, ok := parseLspsSize(fields[3])
		if !ok {
			continue
		}
		percent, err := strconv.ParseFloat(fields[4], 64)
		if err != nil {
			continue
		}

		pagingSpace, exists := byName[fields[0]]
		if !exists {
			pagingSpace = &resultPagingSpace{
				Name:            fields[0],
				PhysicalVolumes: []string{},
				VolumeGroup:     fields[2],
				Type:            fields[7],
				Active:          fields[5] == "yes",
				Auto:            fields[6] == "yes",
			}
			byName[fields[0]] = pagingSpace
			pagingSpaces = append(pagingSpaces, pagingSpace)
		}
		addPagingSpaceVolume(pagingSpace, fields[1], size, percent)
	}
	return pagingSpaces
}

// addPagingSpaceVolume adds the size and the usage of a physical volume to the paging space
func addPagingSpaceVolume(pagingSpace *resultPagingSpace, physicalVolume string, size uint64, percent float64) {
	pagingSpace.PhysicalVolumes = append(pagingSpace.PhysicalVolumes, physicalVolume)
	totalSize := pagingSpace.Size + size
	if totalSize > 0 {
		// Average of the percentages weighted by the size of the physical volumes
		pagingSpace.Percent = (pagingSpace.Percent*float64(pagingSpace.Size) + percent*float64(size)) / float64(totalSize)
	}
	pagingSpace.Size = totalSize
	pagingSpace.Used += uint64(float64(size) * percent / 100)
}

// parseLspsSize parses a size like 4096MB or 2GB and returns bytes
func parseLspsSize(value string) (uint64, bool) {
	units := map[string]uint64{
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
	}
	for suffix, factor := range units {
		if strings.HasSuffix(value, suffix) {
			number, err := strconv.ParseUint(strings.TrimSuffix(value, suffix), 10, 64)
			if err != nil {
				return 0, false
			}
			return number * factor, true
		}
	}
	return 0, false
}
