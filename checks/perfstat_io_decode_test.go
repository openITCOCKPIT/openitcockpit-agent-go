//go:build aix || linux
// +build aix linux

package checks

import (
	"encoding/binary"
	"testing"
)

func putPerfstatUint64(buf []byte, offset int, value uint64) {
	binary.BigEndian.PutUint64(buf[offset:], value)
}

// Structures perfstat_netinterface_t of the libperfstat function perfstat_netinterface() (no command, works without root)
func TestDecodeNetInterfaces(t *testing.T) {
	size := perfstatNetInterfaceSizes[0]
	buf := make([]byte, 2*size)
	copy(buf, "en0")
	buf[netInterfaceOffsetType] = 6
	putPerfstatUint64(buf, netInterfaceOffsetMtu, 1500)
	putPerfstatUint64(buf, netInterfaceOffsetIpackets, 10)
	putPerfstatUint64(buf, netInterfaceOffsetIbytes, 1000)
	putPerfstatUint64(buf, netInterfaceOffsetIerrors, 1)
	putPerfstatUint64(buf, netInterfaceOffsetOpackets, 20)
	putPerfstatUint64(buf, netInterfaceOffsetObytes, 2000)
	putPerfstatUint64(buf, netInterfaceOffsetOerrors, 2)
	putPerfstatUint64(buf, netInterfaceOffsetBitrate, 10000000000)
	putPerfstatUint64(buf, netInterfaceOffsetXmitdrops, 3)
	putPerfstatUint64(buf, netInterfaceOffsetIqdrops, 4)
	copy(buf[size:], "lo0")

	interfaces, err := decodeNetInterfaces(buf, size, 2)
	if err != nil {
		t.Fatal(err)
	}
	want := aixNetInterface{
		Name: "en0", Type: 6, Mtu: 1500, Bitrate: 10000000000,
		PacketsRecv: 10, BytesRecv: 1000, ErrorsIn: 1,
		PacketsSent: 20, BytesSent: 2000, ErrorsOut: 2,
		DropOut: 3, DropIn: 4,
	}
	if interfaces[0] != want {
		t.Errorf("got %+v, want %+v", interfaces[0], want)
	}
	if interfaces[1].Name != "lo0" || interfaces[1].BytesRecv != 0 {
		t.Errorf("unexpected interface %+v", interfaces[1])
	}

	counters := aixNetIOCounters(interfaces)
	if counters[0].Name != "en0" || counters[0].BytesSent != 2000 || counters[0].BytesRecv != 1000 || counters[0].Dropin != 4 || counters[0].Dropout != 3 {
		t.Errorf("unexpected counters %+v", counters[0])
	}

	// Version 1 has no if_iqdrops
	old, err := decodeNetInterfaces(buf[:perfstatNetInterfaceSizes[1]], perfstatNetInterfaceSizes[1], 1)
	if err != nil || old[0].DropIn != 0 || old[0].DropOut != 3 {
		t.Errorf("unexpected version 1 interface %+v %v", old, err)
	}

	if _, err := decodeNetInterfaces(buf[:size], size, 2); err == nil {
		t.Errorf("expected an error for a short buffer")
	}
}

// Structures perfstat_disk_t of the libperfstat function perfstat_disk() (no command, works without root)
func TestDecodeDisks(t *testing.T) {
	size := perfstatDiskSizes[0]
	buf := make([]byte, 2*size)
	copy(buf, "hdisk0")
	copy(buf[diskOffsetVgname:], "rootvg")
	putPerfstatUint64(buf, diskOffsetRxfers, 100)
	putPerfstatUint64(buf, diskOffsetXfers, 250)
	putPerfstatUint64(buf, diskOffsetRblks, 2048)
	putPerfstatUint64(buf, diskOffsetWblks, 4096)
	putPerfstatUint64(buf, diskOffsetTime, 150)
	putPerfstatUint64(buf, diskOffsetRserv, 512000*30)
	putPerfstatUint64(buf, diskOffsetWserv, 512000*70)

	copy(buf[size:], "hdisk9")
	binary.BigEndian.PutUint32(buf[size+diskOffsetDkType:], diskPseudoDiskBit)

	disks, err := decodeDisks(buf, size, 2)
	if err != nil {
		t.Fatal(err)
	}
	want := aixDisk{
		Name: "hdisk0", VolumeGroup: "rootvg",
		ReadTransfers: 100, Transfers: 250, ReadBlocks: 2048, WriteBlocks: 4096,
		BusyTicks: 150, ReadService: 512000 * 30, WriteService: 512000 * 70,
	}
	if disks[0] != want {
		t.Errorf("got %+v, want %+v", disks[0], want)
	}
	if disks[1].Name != "hdisk9" || !disks[1].Pseudo {
		t.Errorf("unexpected disk %+v", disks[1])
	}

	counters := aixDiskIOCounters(disks)
	if len(counters) != 1 {
		t.Fatalf("pseudo disks should be ignored: %+v", counters)
	}
	hdisk0 := counters["hdisk0"]
	if hdisk0.ReadCount != 100 || hdisk0.WriteCount != 150 || hdisk0.ReadBytes != 2048*512 || hdisk0.WriteBytes != 4096*512 {
		t.Errorf("unexpected counters %+v", hdisk0)
	}
	if !aixDisksHaveCounters(disks) || aixDisksHaveCounters([]aixDisk{{Name: "hdisk1"}}) {
		t.Errorf("unexpected result of aixDisksHaveCounters")
	}
	if hdisk0.ReadTime != 30 || hdisk0.WriteTime != 70 || hdisk0.IoTime != 1500 {
		t.Errorf("unexpected times %+v", hdisk0)
	}
}

// Output of "lsattr -E -l sys0 -a iostat -F value" (works without root)
func TestParseAixIostatAttribute(t *testing.T) {
	if enabled, err := parseAixIostatAttribute("true\n"); err != nil || !enabled {
		t.Errorf("got %v %v", enabled, err)
	}
	if enabled, err := parseAixIostatAttribute("false\n"); err != nil || enabled {
		t.Errorf("got %v %v", enabled, err)
	}
	if _, err := parseAixIostatAttribute("lsattr: 0514-519"); err == nil {
		t.Errorf("expected an error")
	}
}
