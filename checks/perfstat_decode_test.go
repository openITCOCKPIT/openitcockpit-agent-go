//go:build aix || linux
// +build aix linux

package checks

import (
	"encoding/binary"
	"math"
	"reflect"
	"testing"
)

func putPerfstatTimes(buf []byte, offset int, times perfstatTimes) {
	binary.BigEndian.PutUint64(buf[offset:], times.User)
	binary.BigEndian.PutUint64(buf[offset+8:], times.Sys)
	binary.BigEndian.PutUint64(buf[offset+16:], times.Idle)
	binary.BigEndian.PutUint64(buf[offset+24:], times.Wait)
}

func assertFloat(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 0.01 {
		t.Errorf("%s: got %.4f, want %.4f", name, got, want)
	}
}

// The values of a shared uncapped LPAR on AIX 7.3 (SMT8, entitlement 0.20) were read with perfstat and
// compared with lparstat in the same 5 seconds: 0.2 %user, 1.1 %sys, 0.0 %wait, 98.8 %idle
func TestAixCpuResultSharedPartition(t *testing.T) {
	prev := &aixCpuTimes{
		Logical:     perfstatTimes{User: 1000, Sys: 2000, Idle: 100000, Wait: 10},
		Partition:   perfstatTimes{User: 5406030855164, Sys: 20863889681320, Idle: 100944711474780, Wait: 51238074736},
		Timebase:    69238650772483126,
		Entitlement: 0.20,
		Shared:      true,
		Cpus:        []perfstatTimes{{User: 10, Sys: 10, Idle: 80}, {User: 5, Idle: 5}},
	}
	current := &aixCpuTimes{
		Logical:     perfstatTimes{User: 1000, Sys: 2000, Idle: 100000 + 3995, Wait: 10},
		Partition:   perfstatTimes{User: 5406030855164 + 940000, Sys: 20863889681320 + 5474000, Idle: 100944711474780 + 27169744, Wait: 51238074736},
		Timebase:    69238650772483126 + 2560001316,
		Entitlement: 0.20,
		Shared:      true,
		Cpus:        []perfstatTimes{{User: 30, Sys: 30, Idle: 140}, {User: 5, Idle: 5}},
	}

	result := aixCpuResult(prev, current)
	assertFloat(t, "user", result.DetailsTotal.User, 0.1836)
	assertFloat(t, "system", result.DetailsTotal.System, 1.0691)
	assertFloat(t, "iowait", result.DetailsTotal.Iowait, 0)
	assertFloat(t, "idle", result.DetailsTotal.Idle, 98.7473)
	assertFloat(t, "total", result.PercentageTotal, 1.2527)

	// Per logical CPU: the usage while the thread was dispatched
	if len(result.PercentagePerCore) != 2 || len(result.DetailsPerCore) != 2 {
		t.Fatalf("unexpected number of CPUs: %v", result.PercentagePerCore)
	}
	assertFloat(t, "cpu0", result.PercentagePerCore[0], 40)
	assertFloat(t, "cpu0 idle", result.DetailsPerCore[0].Idle, 60)
	// A thread which was not dispatched is idle
	assertFloat(t, "cpu1", result.PercentagePerCore[1], 0)
	assertFloat(t, "cpu1 idle", result.DetailsPerCore[1].Idle, 100)
}

// The unused entitlement is split between idle and wait like the ticks of the logical CPUs
func TestAixCpuResultSharedPartitionWait(t *testing.T) {
	prev := &aixCpuTimes{Entitlement: 1, Shared: true, Cpus: []perfstatTimes{{}}}
	current := &aixCpuTimes{
		Logical:     perfstatTimes{Idle: 75, Wait: 25},
		Partition:   perfstatTimes{User: 100, Sys: 100},
		Timebase:    1000,
		Entitlement: 1,
		Shared:      true,
		Cpus:        []perfstatTimes{{}},
	}

	result := aixCpuResult(prev, current)
	assertFloat(t, "user", result.DetailsTotal.User, 10)
	assertFloat(t, "system", result.DetailsTotal.System, 10)
	assertFloat(t, "iowait", result.DetailsTotal.Iowait, 20)
	assertFloat(t, "idle", result.DetailsTotal.Idle, 60)
	assertFloat(t, "total", result.PercentageTotal, 20)
}

// An uncapped partition which uses more than its entitlement is calculated relative to the used PURR
func TestAixCpuResultUncappedAboveEntitlement(t *testing.T) {
	prev := &aixCpuTimes{Entitlement: 0.5, Shared: true}
	current := &aixCpuTimes{
		Logical:     perfstatTimes{User: 90, Idle: 10},
		Partition:   perfstatTimes{User: 1800, Idle: 200},
		Timebase:    1000,
		Entitlement: 0.5,
		Shared:      true,
	}

	result := aixCpuResult(prev, current)
	assertFloat(t, "user", result.DetailsTotal.User, 90)
	assertFloat(t, "idle", result.DetailsTotal.Idle, 10)
	assertFloat(t, "total", result.PercentageTotal, 90)
}

func TestAixCpuResultDedicatedPartition(t *testing.T) {
	prev := &aixCpuTimes{Entitlement: 2}
	current := &aixCpuTimes{
		Logical:     perfstatTimes{User: 50, Idle: 50},
		Partition:   perfstatTimes{User: 300, Sys: 100, Idle: 1500, Wait: 100},
		Timebase:    1000,
		Entitlement: 2,
	}

	result := aixCpuResult(prev, current)
	assertFloat(t, "user", result.DetailsTotal.User, 15)
	assertFloat(t, "system", result.DetailsTotal.System, 5)
	assertFloat(t, "iowait", result.DetailsTotal.Iowait, 5)
	assertFloat(t, "idle", result.DetailsTotal.Idle, 75)
}

// Without the values of perfstat_partition_total the ticks of the logical CPUs are used
func TestAixCpuResultWithoutPartition(t *testing.T) {
	prev := &aixCpuTimes{Logical: perfstatTimes{User: 100}}
	current := &aixCpuTimes{Logical: perfstatTimes{User: 130, Sys: 10, Idle: 50, Wait: 10}}

	result := aixCpuResult(prev, current)
	assertFloat(t, "user", result.DetailsTotal.User, 30)
	assertFloat(t, "system", result.DetailsTotal.System, 10)
	assertFloat(t, "idle", result.DetailsTotal.Idle, 50)
	assertFloat(t, "total", result.PercentageTotal, 40)

	// Identical samples
	result = aixCpuResult(current, current)
	if math.IsNaN(result.PercentageTotal) || result.PercentageTotal != 0 || result.DetailsTotal.Idle != 100 {
		t.Errorf("unexpected result for identical samples: %v %+v", result.PercentageTotal, result.DetailsTotal)
	}
}

func TestDecodePartition(t *testing.T) {
	buf := make([]byte, perfstatPartitionSize)
	// Type of a shared uncapped partition with SMT
	binary.BigEndian.PutUint32(buf[partitionOffsetType:], 0xfeaa0000)
	binary.BigEndian.PutUint32(buf[partitionOffsetEntitlement:], 20)
	putPerfstatTimes(buf, partitionOffsetPuser, perfstatTimes{User: 1, Sys: 2, Idle: 3, Wait: 4})
	binary.BigEndian.PutUint64(buf[partitionOffsetTimebase:], 69238650772483126)

	times := &aixCpuTimes{}
	if err := decodePartition(buf, times); err != nil {
		t.Fatal(err)
	}
	if !times.Shared {
		t.Error("partition is not shared")
	}
	assertFloat(t, "entitlement", times.Entitlement, 0.2)
	if times.Partition != (perfstatTimes{User: 1, Sys: 2, Idle: 3, Wait: 4}) || times.Timebase != 69238650772483126 {
		t.Errorf("unexpected partition values: %+v %v", times.Partition, times.Timebase)
	}

	// Dedicated partition (shared_enabled is not set)
	binary.BigEndian.PutUint32(buf[partitionOffsetType:], 0xf2aa0000)
	if err := decodePartition(buf, times); err != nil || times.Shared {
		t.Errorf("dedicated partition: shared=%v err=%v", times.Shared, err)
	}

	if err := decodePartition(buf[:100], times); err == nil {
		t.Error("expected an error for a short buffer")
	}
}

func TestDecodeCpuTotal(t *testing.T) {
	buf := make([]byte, perfstatCpuTotalSize)
	putPerfstatTimes(buf, cpuTotalOffsetUser, perfstatTimes{User: 10, Sys: 20, Idle: 30, Wait: 40})
	binary.BigEndian.PutUint64(buf[cpuTotalOffsetLoadAverage:], 65536+32768) // 1.5
	binary.BigEndian.PutUint64(buf[cpuTotalOffsetLoadAverage+8:], 16384)     // 0.25
	binary.BigEndian.PutUint64(buf[cpuTotalOffsetLoadAverage+16:], 0)

	ticks, err := decodeCpuTotalTicks(buf)
	if err != nil || ticks != (perfstatTimes{User: 10, Sys: 20, Idle: 30, Wait: 40}) {
		t.Errorf("unexpected ticks %+v %v", ticks, err)
	}
	load, err := decodeLoadAverage(buf)
	if err != nil || load != [3]float64{1.5, 0.25, 0} {
		t.Errorf("unexpected load average %v %v", load, err)
	}
	if _, err := decodeLoadAverage(buf[:200]); err == nil {
		t.Error("expected an error for a short buffer")
	}
}

func TestDecodeCpus(t *testing.T) {
	buf := make([]byte, 2*perfstatCpuSize)
	putPerfstatTimes(buf, cpuOffsetPuser, perfstatTimes{User: 1, Sys: 2, Idle: 3, Wait: 4})
	putPerfstatTimes(buf, perfstatCpuSize+cpuOffsetPuser, perfstatTimes{User: 5, Sys: 6, Idle: 7, Wait: 8})
	// The ticks of the logical CPU must not be used
	putPerfstatTimes(buf, cpuOffsetUser, perfstatTimes{User: 99, Sys: 99, Idle: 99, Wait: 99})

	cpus, err := decodeCpus(buf, 2)
	if err != nil {
		t.Fatal(err)
	}
	want := []perfstatTimes{{User: 1, Sys: 2, Idle: 3, Wait: 4}, {User: 5, Sys: 6, Idle: 7, Wait: 8}}
	if !reflect.DeepEqual(cpus, want) {
		t.Errorf("got %+v, want %+v", cpus, want)
	}
	if _, err := decodeCpus(buf, 3); err == nil {
		t.Error("expected an error for a short buffer")
	}
}

// Values of AIX 7.3 with 6 GiB, vmstat -v showed 46 % computational pages
func TestAixMemoryAndSwapResult(t *testing.T) {
	buf := make([]byte, perfstatMemoryTotalSize)
	binary.BigEndian.PutUint64(buf[memoryOffsetRealTotal:], 1572864)
	binary.BigEndian.PutUint64(buf[memoryOffsetRealFree:], 300000)
	binary.BigEndian.PutUint64(buf[memoryOffsetRealPinned:], 400000)
	binary.BigEndian.PutUint64(buf[memoryOffsetNumperm:], 549346)
	binary.BigEndian.PutUint64(buf[memoryOffsetPgspTotal:], 262144)
	binary.BigEndian.PutUint64(buf[memoryOffsetPgspFree:], 196608)
	binary.BigEndian.PutUint64(buf[memoryOffsetPgspIns:], 10)
	binary.BigEndian.PutUint64(buf[memoryOffsetPgspOuts:], 20)

	memory, err := decodeMemoryTotal(buf)
	if err != nil {
		t.Fatal(err)
	}

	mem := aixMemoryResult(memory)
	if mem.Total != 6*1024*1024*1024 || mem.Free != 300000*4096 || mem.Wired != 400000*4096 {
		t.Errorf("unexpected memory result %+v", mem)
	}
	if mem.Used != (1572864-300000-549346)*4096 || mem.Available != mem.Total-mem.Used {
		t.Errorf("unexpected used memory %+v", mem)
	}
	assertFloat(t, "memory percent", mem.Percent, 46)

	swap := aixSwapResult(memory)
	if swap.Total != 1024*1024*1024 || swap.Used != 256*1024*1024 || swap.Free != 768*1024*1024 {
		t.Errorf("unexpected swap result %+v", swap)
	}
	if swap.Sin != 10*4096 || swap.Sout != 20*4096 {
		t.Errorf("unexpected swap in/out %+v", swap)
	}
	assertFloat(t, "swap percent", swap.Percent, 25)

	// No paging space
	if swap := aixSwapResult(aixMemory{}); swap.Percent != 0 || swap.Used != 0 {
		t.Errorf("unexpected swap result without paging space %+v", swap)
	}
}

func TestParseAixMount(t *testing.T) {
	output := `  node       mounted        mounted over    vfs       date        options
-------- ---------------  ---------------  ------ ------------ ---------------
         /dev/hd4         /                jfs2   Apr 30 13:17 rw,log=/dev/hd8
         /dev/livedump    /var/adm/ras/livedump jfs2   Apr 30 13:18 rw,log=/dev/hd8
         /proc            /proc            procfs Apr 30 13:18 rw
         /aha             /aha             ahafs  Apr 30 13:18 rw
nfsserver /export/data    /mnt/data        nfs3   May 02 08:01 bg,hard,intr,vers=3
         /dev/cd0         /cdrom           cdrfs  May 02 08:05
`
	want := []aixMount{
		{Device: "/dev/hd4", Mountpoint: "/", Fstype: "jfs2", Opts: []string{"rw", "log=/dev/hd8"}},
		{Device: "/dev/livedump", Mountpoint: "/var/adm/ras/livedump", Fstype: "jfs2", Opts: []string{"rw", "log=/dev/hd8"}},
		{Device: "/proc", Mountpoint: "/proc", Fstype: "procfs", Opts: []string{"rw"}},
		{Device: "/aha", Mountpoint: "/aha", Fstype: "ahafs", Opts: []string{"rw"}},
		{Device: "nfsserver:/export/data", Mountpoint: "/mnt/data", Fstype: "nfs3", Opts: []string{"bg", "hard", "intr", "vers=3"}},
		{Device: "/dev/cd0", Mountpoint: "/cdrom", Fstype: "cdrfs"},
	}

	mounts := parseAixMount(output)
	if !reflect.DeepEqual(mounts, want) {
		t.Errorf("got\n%+v\nwant\n%+v", mounts, want)
	}
}
