//go:build aix || linux
// +build aix linux

package checks

import (
	"encoding/binary"
	"math"
	"reflect"
	"testing"
	"time"
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
	binary.BigEndian.PutUint32(buf[partitionOffsetOnlineCpus:], 1)
	binary.BigEndian.PutUint64(buf[partitionOffsetOnlineMemory:], 6144)
	binary.BigEndian.PutUint32(buf[partitionOffsetPoolCpus:], 8)
	binary.BigEndian.PutUint64(buf[partitionOffsetInvoluntaryCswitch:], 11)
	binary.BigEndian.PutUint64(buf[partitionOffsetVoluntaryCswitch:], 22)
	binary.BigEndian.PutUint32(buf[partitionOffsetOnlineLogicalCpus:], 8)
	binary.BigEndian.PutUint32(buf[partitionOffsetSmtThreads:], 8)

	times := &aixCpuTimes{}
	if err := decodePartition(buf, times); err != nil {
		t.Fatal(err)
	}
	if !times.PartitionAvailable || !times.Shared || times.Capped || !times.SmtEnabled {
		t.Errorf("unexpected partition type: %+v", times)
	}
	if times.VirtualCpus != 1 || times.PoolCpus != 8 || times.MemoryMB != 6144 || times.LogicalCpus != 8 || times.SmtThreads != 8 {
		t.Errorf("unexpected partition configuration: %+v", times)
	}
	if times.InvoluntaryCswitches != 11 || times.VoluntaryCswitches != 22 {
		t.Errorf("unexpected context switches: %+v", times)
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

	// Shared capped partition
	binary.BigEndian.PutUint32(buf[partitionOffsetType:], 0xffaa0000)
	if err := decodePartition(buf, times); err != nil || !times.Shared || !times.Capped {
		t.Errorf("capped partition: shared=%v capped=%v err=%v", times.Shared, times.Capped, err)
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

// sharedPartitionSamples returns two samples of the shared uncapped LPAR of AIX 7.3 (1 virtual processor,
// entitlement 0.20, pool of 8 processors), 5 seconds apart. purr is the consumed PURR of the interval.
func sharedPartitionSamples(purr uint64) (*aixCpuTimes, *aixCpuTimes) {
	start := time.Date(2026, 11, 13, 10, 4, 0, 0, time.UTC)
	prev := &aixCpuTimes{
		Time:                 start,
		Partition:            perfstatTimes{User: 1000, Sys: 2000, Idle: 3000},
		Timebase:             69238650772483126,
		Entitlement:          0.20,
		PartitionAvailable:   true,
		Shared:               true,
		SmtEnabled:           true,
		VirtualCpus:          1,
		PoolCpus:             8,
		LogicalCpus:          8,
		SmtThreads:           8,
		MemoryMB:             6144,
		VoluntaryCswitches:   1000,
		InvoluntaryCswitches: 100,
	}
	current := *prev
	current.Time = start.Add(5 * time.Second)
	current.Partition = perfstatTimes{User: 1000 + purr/2, Sys: 2000 + purr/4, Idle: 3000 + purr/4}
	current.Timebase = prev.Timebase + 2560001316
	current.VoluntaryCswitches = 1000 + 1400
	current.InvoluntaryCswitches = 100 + 100
	return prev, &current
}

// physc 0.013 and 6.6 %entc like lparstat
func TestAixLparResultSharedUncapped(t *testing.T) {
	prev, current := sharedPartitionSamples(33583744)

	result, err := aixLparResult(prev, current)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Shared || result.Capped || result.VirtualCpus != 1 || result.PoolCpus != 8 || result.SmtThreads != 8 || result.LogicalCpus != 8 || result.MemoryMB != 6144 {
		t.Errorf("unexpected configuration %+v", result)
	}
	assertFloat(t, "entitlement", result.Entitlement, 0.2)
	assertFloat(t, "physc", result.PhysicalCpus, 0.0131)
	assertFloat(t, "entc", result.EntitlementPercentage, 6.56)
	// An uncapped partition can use one physical processor per virtual processor
	assertFloat(t, "max physc", result.MaxPhysicalCpus, 1)
	assertFloat(t, "capacity", result.CapacityPercentage, 1.31)
	assertFloat(t, "virtual context switches", result.VirtualCswitches, 300)
	assertFloat(t, "involuntary virtual context switches", result.InvoluntaryVirtualCswitches, 20)
}

// One busy thread: physc 1.00 and 497.5 %entc like lparstat, the partition uses all of its virtual processors
func TestAixLparResultSharedUncappedBusy(t *testing.T) {
	prev, current := sharedPartitionSamples(2547201309)

	result, err := aixLparResult(prev, current)
	if err != nil {
		t.Fatal(err)
	}
	assertFloat(t, "physc", result.PhysicalCpus, 0.995)
	assertFloat(t, "entc", result.EntitlementPercentage, 497.5)
	assertFloat(t, "capacity", result.CapacityPercentage, 99.5)
}

// A capped partition can use at most its entitlement
func TestAixLparResultSharedCapped(t *testing.T) {
	prev, current := sharedPartitionSamples(2560001316 / 10) // physc 0.1
	prev.Capped = true
	current.Capped = true

	result, err := aixLparResult(prev, current)
	if err != nil {
		t.Fatal(err)
	}
	assertFloat(t, "physc", result.PhysicalCpus, 0.1)
	assertFloat(t, "max physc", result.MaxPhysicalCpus, 0.2)
	assertFloat(t, "capacity", result.CapacityPercentage, 50)
	assertFloat(t, "entc", result.EntitlementPercentage, 50)
}

// The pool limits an uncapped partition with more virtual processors than the pool has processors
func TestAixMaxPhysicalCpusLimitedByPool(t *testing.T) {
	times := &aixCpuTimes{VirtualCpus: 12, PoolCpus: 8, Entitlement: 2}
	assertFloat(t, "max physc", aixMaxPhysicalCpus(times), 8)
	times.PoolCpus = 0
	assertFloat(t, "max physc without pool", aixMaxPhysicalCpus(times), 12)
}

// Dedicated partitions own their processors, there are no physc values
func TestAixLparResultDedicated(t *testing.T) {
	prev, current := sharedPartitionSamples(2560001316)
	prev.Shared = false
	current.Shared = false
	current.Entitlement = 2

	result, err := aixLparResult(prev, current)
	if err != nil {
		t.Fatal(err)
	}
	if result.Shared || result.PhysicalCpus != 0 || result.EntitlementPercentage != 0 || result.CapacityPercentage != 0 {
		t.Errorf("unexpected usage of a dedicated partition %+v", result)
	}

	if cpu := aixCpuResult(prev, current); cpu.Partition != nil {
		t.Errorf("unexpected partition values of a dedicated partition in the CPU result %+v", cpu.Partition)
	}
}

func TestAixLparResultWithoutPartition(t *testing.T) {
	if _, err := aixLparResult(&aixCpuTimes{}, &aixCpuTimes{}); err == nil {
		t.Error("expected an error without partition values")
	}
}

func TestAixCpuResultPartition(t *testing.T) {
	prev, current := sharedPartitionSamples(33583744)

	result := aixCpuResult(prev, current)
	if result.Partition == nil {
		t.Fatal("no partition values")
	}
	assertFloat(t, "physc", result.Partition.PhysicalCpus, 0.0131)
	assertFloat(t, "entitlement", result.Partition.Entitlement, 0.2)
	assertFloat(t, "entc", result.Partition.EntitlementPercentage, 6.56)
}
