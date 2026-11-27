//go:build aix || linux
// +build aix linux

// Decoding of perfstat_netinterface_t and perfstat_disk_t (AIX) without any AIX dependencies.
// It is also built on Linux, so the unit tests run in the CI.

package checks

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/net"
)

// Sizes of the structure versions. perfstat rejects sizes it does not know, the sizes are tried in this order.
var (
	perfstatNetInterfaceSizes = []int{240, 224} // version 3 (with if_iqdrops and if_arpdrops), version 1
	perfstatDiskSizes         = []int{496, 480} // version 2 (with version and dk_type), version 1
)

// Offsets of perfstat_netinterface_t (64 bit, big endian)
const (
	netInterfaceOffsetName       = 0 // char[64]
	netInterfaceOffsetType       = 128
	netInterfaceOffsetMtu        = 136
	netInterfaceOffsetIpackets   = 144
	netInterfaceOffsetIbytes     = 152
	netInterfaceOffsetIerrors    = 160
	netInterfaceOffsetOpackets   = 168
	netInterfaceOffsetObytes     = 176
	netInterfaceOffsetOerrors    = 184
	netInterfaceOffsetCollisions = 192
	netInterfaceOffsetBitrate    = 200
	netInterfaceOffsetXmitdrops  = 208
	netInterfaceOffsetIqdrops    = 224 // version 3
)

// Offsets of perfstat_disk_t (64 bit, big endian)
const (
	diskOffsetName    = 0   // char[64]
	diskOffsetVgname  = 128 // char[64]
	diskOffsetRxfers  = 216 // xrate, number of transfers from the disk
	diskOffsetXfers   = 224 // transfers to and from the disk
	diskOffsetWblks   = 232
	diskOffsetRblks   = 240
	diskOffsetTime    = 256 // time the disk is active, in clock ticks
	diskOffsetRserv   = 344 // read service time, in timebase ticks
	diskOffsetWserv   = 384 // write service time, in timebase ticks
	diskOffsetDkType  = 488 // version 2: bit field, the first bit marks pseudo disks
	diskBlockSize     = 512 // rblks and wblks are 512 byte blocks
	diskPseudoDiskBit = 1 << 31
)

const (
	// aixClockTicksPerSecond is the rate of the clock ticks of AIX (HZ)
	aixClockTicksPerSecond = 100
	// aixTimebaseFrequency is the frequency of the timebase register of POWER7 and newer
	aixTimebaseFrequency = 512000000
)

// aixNetInterface are the counters of a network interface (perfstat_netinterface_t)
type aixNetInterface struct {
	Name        string
	Type        uint8 // see /usr/include/net/if_types.h, e.g. 6 = Ethernet, 24 = loopback
	Mtu         uint64
	Bitrate     uint64 // bit per second
	PacketsRecv uint64
	BytesRecv   uint64
	ErrorsIn    uint64
	PacketsSent uint64
	BytesSent   uint64
	ErrorsOut   uint64
	Collisions  uint64
	DropOut     uint64
	DropIn      uint64 // 0 for version 1
}

// aixDisk are the counters of a disk (perfstat_disk_t). AIX only counts if the I/O history is enabled
// (sys0 attribute iostat=true) or while a tool like iostat or topas_nmon is running.
// The units were verified with iostat -D on AIX 7.3: 512 byte blocks, busy time in clock ticks (100 Hz),
// service times in timebase ticks (512 MHz).
type aixDisk struct {
	Name          string
	VolumeGroup   string
	ReadTransfers uint64
	Transfers     uint64
	ReadBlocks    uint64
	WriteBlocks   uint64
	BusyTicks     uint64
	ReadService   uint64 // timebase ticks
	WriteService  uint64 // timebase ticks
	Pseudo        bool
}

// perfstatString returns a NUL terminated string of a structure
func perfstatString(buf []byte, offset int) string {
	value := buf[offset : offset+perfstatIdentifierLength]
	if i := bytes.IndexByte(value, 0); i >= 0 {
		value = value[:i]
	}
	return string(value)
}

// decodeNetInterfaces decodes count structures of the given size
func decodeNetInterfaces(buf []byte, size, count int) ([]aixNetInterface, error) {
	if err := checkPerfstatSize("perfstat_netinterface", buf, size*count); err != nil {
		return nil, err
	}
	interfaces := make([]aixNetInterface, 0, count)
	for i := 0; i < count; i++ {
		interfaces = append(interfaces, decodeNetInterface(buf[i*size:(i+1)*size]))
	}
	return interfaces, nil
}

func decodeNetInterface(buf []byte) aixNetInterface {
	netInterface := aixNetInterface{
		Name:        perfstatString(buf, netInterfaceOffsetName),
		Type:        buf[netInterfaceOffsetType],
		Mtu:         perfstatUint64(buf, netInterfaceOffsetMtu),
		Bitrate:     perfstatUint64(buf, netInterfaceOffsetBitrate),
		PacketsRecv: perfstatUint64(buf, netInterfaceOffsetIpackets),
		BytesRecv:   perfstatUint64(buf, netInterfaceOffsetIbytes),
		ErrorsIn:    perfstatUint64(buf, netInterfaceOffsetIerrors),
		PacketsSent: perfstatUint64(buf, netInterfaceOffsetOpackets),
		BytesSent:   perfstatUint64(buf, netInterfaceOffsetObytes),
		ErrorsOut:   perfstatUint64(buf, netInterfaceOffsetOerrors),
		Collisions:  perfstatUint64(buf, netInterfaceOffsetCollisions),
		DropOut:     perfstatUint64(buf, netInterfaceOffsetXmitdrops),
	}
	if len(buf) >= netInterfaceOffsetIqdrops+8 {
		netInterface.DropIn = perfstatUint64(buf, netInterfaceOffsetIqdrops)
	}
	return netInterface
}

// decodeDisks decodes count structures of the given size
func decodeDisks(buf []byte, size, count int) ([]aixDisk, error) {
	if err := checkPerfstatSize("perfstat_disk", buf, size*count); err != nil {
		return nil, err
	}
	disks := make([]aixDisk, 0, count)
	for i := 0; i < count; i++ {
		disks = append(disks, decodeDisk(buf[i*size:(i+1)*size]))
	}
	return disks, nil
}

func decodeDisk(buf []byte) aixDisk {
	disk := aixDisk{
		Name:          perfstatString(buf, diskOffsetName),
		VolumeGroup:   perfstatString(buf, diskOffsetVgname),
		ReadTransfers: perfstatUint64(buf, diskOffsetRxfers),
		Transfers:     perfstatUint64(buf, diskOffsetXfers),
		ReadBlocks:    perfstatUint64(buf, diskOffsetRblks),
		WriteBlocks:   perfstatUint64(buf, diskOffsetWblks),
		BusyTicks:     perfstatUint64(buf, diskOffsetTime),
		ReadService:   perfstatUint64(buf, diskOffsetRserv),
		WriteService:  perfstatUint64(buf, diskOffsetWserv),
	}
	if len(buf) >= diskOffsetDkType+4 {
		disk.Pseudo = uint32(perfstatInt32(buf, diskOffsetDkType))&diskPseudoDiskBit != 0
	}
	return disk
}

// aixNetIOCounters converts the counters of the network interfaces
func aixNetIOCounters(interfaces []aixNetInterface) []net.IOCountersStat {
	counters := make([]net.IOCountersStat, 0, len(interfaces))
	for _, netInterface := range interfaces {
		counters = append(counters, net.IOCountersStat{
			Name:        netInterface.Name,
			BytesSent:   netInterface.BytesSent,
			BytesRecv:   netInterface.BytesRecv,
			PacketsSent: netInterface.PacketsSent,
			PacketsRecv: netInterface.PacketsRecv,
			Errin:       netInterface.ErrorsIn,
			Errout:      netInterface.ErrorsOut,
			Dropin:      netInterface.DropIn,
			Dropout:     netInterface.DropOut,
		})
	}
	return counters
}

// aixDiskIOCounters converts the counters of the physical disks, the times are converted to milliseconds
func aixDiskIOCounters(disks []aixDisk) map[string]disk.IOCountersStat {
	counters := make(map[string]disk.IOCountersStat, len(disks))
	for _, aixDisk := range disks {
		if aixDisk.Pseudo {
			continue
		}
		writeTransfers := uint64(0)
		if aixDisk.Transfers > aixDisk.ReadTransfers {
			writeTransfers = aixDisk.Transfers - aixDisk.ReadTransfers
		}
		counters[aixDisk.Name] = disk.IOCountersStat{
			Name:       aixDisk.Name,
			ReadCount:  aixDisk.ReadTransfers,
			WriteCount: writeTransfers,
			ReadBytes:  aixDisk.ReadBlocks * diskBlockSize,
			WriteBytes: aixDisk.WriteBlocks * diskBlockSize,
			ReadTime:   timebaseToMilliseconds(aixDisk.ReadService),
			WriteTime:  timebaseToMilliseconds(aixDisk.WriteService),
			IoTime:     aixDisk.BusyTicks * 1000 / aixClockTicksPerSecond,
		}
	}
	return counters
}

// aixDisksHaveCounters returns true if any disk has counted a transfer
func aixDisksHaveCounters(disks []aixDisk) bool {
	for _, aixDisk := range disks {
		if aixDisk.Transfers > 0 {
			return true
		}
	}
	return false
}

func timebaseToMilliseconds(ticks uint64) uint64 {
	return ticks / (aixTimebaseFrequency / 1000)
}

// parseAixIostatAttribute returns true if AIX maintains the disk I/O history (lsattr -E -l sys0 -a iostat -F value)
func parseAixIostatAttribute(output string) (bool, error) {
	value := strings.TrimSpace(output)
	switch value {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}
	return false, fmt.Errorf("unexpected value of the sys0 attribute iostat: %q", value)
}
