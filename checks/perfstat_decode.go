//go:build aix || linux
// +build aix linux

// The structures of libperfstat (AIX) are decoded and the values are calculated in this file without any
// AIX dependencies. It is also built on Linux, so the unit tests run in the CI.

package checks

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"time"
)

// Sizes of the perfstat structures (AIX 7.2 and newer). The offsets of /usr/include/libperfstat.h
// were verified on AIX 7.3 (big endian).
const (
	perfstatIdentifierLength = 64

	perfstatCpuTotalSize     = 696
	perfstatCpuTotalSize71   = 672 // AIX 7.1
	perfstatCpuSize          = 640
	perfstatMemoryTotalSize  = 176
	perfstatPartitionSize    = 688 // AIX 7.1 TL1 and newer
	perfstatPartitionSize71  = 608
	perfstatMaxLogicalCpus   = 4096
	perfstatPageSize         = 4096 // The memory values are in 4 KiB pages
	perfstatLoadAverageShift = 16   // SBITS
)

// perfstat_cpu_total_t
const (
	cpuTotalOffsetUser        = 80  // user, sys, idle, wait: ticks of the logical CPUs
	cpuTotalOffsetLoadAverage = 200 // loadavg[3]
)

// perfstat_cpu_t
const (
	cpuOffsetUser  = 64  // user, sys, idle, wait: ticks of the logical CPU
	cpuOffsetPuser = 264 // puser, psys, pidle, pwait: PURR of the logical CPU
)

// perfstat_partition_total_t
const (
	partitionOffsetType               = 64
	partitionOffsetOnlineCpus         = 80  // online virtual processors
	partitionOffsetOnlineMemory       = 96  // MB
	partitionOffsetEntitlement        = 120 // entitled_proc_capacity in 1/100 processor units
	partitionOffsetPoolCpus           = 156 // phys_cpus_pool: physical processors of the shared pool
	partitionOffsetPuser              = 160 // puser, psys, pidle, pwait: PURR of the partition
	partitionOffsetInvoluntaryCswitch = 208 // invol_virt_cswitch
	partitionOffsetVoluntaryCswitch   = 216 // vol_virt_cswitch
	partitionOffsetTimebase           = 224 // timebase_last
	partitionOffsetOnlineLogicalCpus  = 432
	partitionOffsetSmtThreads         = 436

	// Bit field perfstat_partition_type_t, the first field is the most significant bit
	partitionTypeSmtEnabled    = 1 << 30
	partitionTypeSharedEnabled = 1 << 26
	partitionTypeCapped        = 1 << 24
)

// perfstat_memory_total_t
const (
	memoryOffsetRealTotal  = 8
	memoryOffsetRealFree   = 16
	memoryOffsetRealPinned = 24
	memoryOffsetPgspIns    = 72
	memoryOffsetPgspOuts   = 80
	memoryOffsetNumperm    = 112
	memoryOffsetPgspTotal  = 120
	memoryOffsetPgspFree   = 128
)

// perfstatTimes contains user, sys, idle and wait as ticks or PURR (the order of perfstat)
type perfstatTimes struct {
	User uint64
	Sys  uint64
	Idle uint64
	Wait uint64
}

// aixCpuTimes is one sample of the CPU times of an AIX LPAR
type aixCpuTimes struct {
	Time time.Time
	// Ticks of all logical CPUs
	Logical perfstatTimes
	// PURR (Processor Utilization of Resources Register) of the partition, in timebase ticks
	Partition perfstatTimes
	Timebase  uint64
	// Entitled processor units, e.g. 0.2
	Entitlement float64
	// PartitionAvailable is false if perfstat_partition_total failed
	PartitionAvailable bool
	Shared             bool
	Capped             bool
	SmtEnabled         bool
	VirtualCpus        int32 // Online virtual processors (dedicated partitions: processors)
	PoolCpus           int32 // Physical processors of the shared processor pool
	LogicalCpus        uint32
	SmtThreads         uint32
	MemoryMB           uint64
	// Virtual processor context switches
	VoluntaryCswitches   uint64
	InvoluntaryCswitches uint64
	// PURR of the logical CPUs (SMT threads)
	Cpus []perfstatTimes
}

type aixMemory struct {
	RealTotal  uint64
	RealFree   uint64
	RealPinned uint64
	Numperm    uint64 // Pages used for the file cache
	PgspTotal  uint64
	PgspFree   uint64
	PgspIns    uint64
	PgspOuts   uint64
}

type aixMount struct {
	Device     string
	Mountpoint string
	Fstype     string
	Opts       []string
}

func perfstatUint64(buf []byte, offset int) uint64 {
	return binary.BigEndian.Uint64(buf[offset : offset+8])
}

func perfstatInt32(buf []byte, offset int) int32 {
	return int32(binary.BigEndian.Uint32(buf[offset : offset+4]))
}

func decodePerfstatTimes(buf []byte, offset int) perfstatTimes {
	return perfstatTimes{
		User: perfstatUint64(buf, offset),
		Sys:  perfstatUint64(buf, offset+8),
		Idle: perfstatUint64(buf, offset+16),
		Wait: perfstatUint64(buf, offset+24),
	}
}

func checkPerfstatSize(name string, buf []byte, size int) error {
	if len(buf) < size {
		return fmt.Errorf("%s: %d bytes are too short, expected %d", name, len(buf), size)
	}
	return nil
}

// decodeCpuTotalTicks returns the ticks of all logical CPUs of perfstat_cpu_total_t
func decodeCpuTotalTicks(buf []byte) (perfstatTimes, error) {
	if err := checkPerfstatSize("perfstat_cpu_total", buf, perfstatCpuTotalSize71); err != nil {
		return perfstatTimes{}, err
	}
	return decodePerfstatTimes(buf, cpuTotalOffsetUser), nil
}

// decodeLoadAverage returns the load average of 1, 5 and 15 minutes of perfstat_cpu_total_t
func decodeLoadAverage(buf []byte) ([3]float64, error) {
	var load [3]float64
	if err := checkPerfstatSize("perfstat_cpu_total", buf, perfstatCpuTotalSize71); err != nil {
		return load, err
	}
	for i := range load {
		value := perfstatUint64(buf, cpuTotalOffsetLoadAverage+i*8)
		load[i] = float64(value) / float64(uint64(1)<<perfstatLoadAverageShift)
	}
	return load, nil
}

// decodeCpus returns the PURR of every logical CPU of an array of perfstat_cpu_t
func decodeCpus(buf []byte, count int) ([]perfstatTimes, error) {
	if err := checkPerfstatSize("perfstat_cpu", buf, count*perfstatCpuSize); err != nil {
		return nil, err
	}
	cpus := make([]perfstatTimes, count)
	for i := range cpus {
		cpus[i] = decodePerfstatTimes(buf[i*perfstatCpuSize:], cpuOffsetPuser)
	}
	return cpus, nil
}

// decodePartition sets the partition values of perfstat_partition_total_t
func decodePartition(buf []byte, times *aixCpuTimes) error {
	if err := checkPerfstatSize("perfstat_partition_total", buf, perfstatPartitionSize71); err != nil {
		return err
	}
	partitionType := uint32(perfstatInt32(buf, partitionOffsetType))
	times.PartitionAvailable = true
	times.Shared = partitionType&partitionTypeSharedEnabled != 0
	times.Capped = partitionType&partitionTypeCapped != 0
	times.SmtEnabled = partitionType&partitionTypeSmtEnabled != 0
	times.VirtualCpus = perfstatInt32(buf, partitionOffsetOnlineCpus)
	times.MemoryMB = perfstatUint64(buf, partitionOffsetOnlineMemory)
	times.Entitlement = float64(perfstatInt32(buf, partitionOffsetEntitlement)) / 100
	times.PoolCpus = perfstatInt32(buf, partitionOffsetPoolCpus)
	times.Partition = decodePerfstatTimes(buf, partitionOffsetPuser)
	times.InvoluntaryCswitches = perfstatUint64(buf, partitionOffsetInvoluntaryCswitch)
	times.VoluntaryCswitches = perfstatUint64(buf, partitionOffsetVoluntaryCswitch)
	times.Timebase = perfstatUint64(buf, partitionOffsetTimebase)
	times.LogicalCpus = uint32(perfstatInt32(buf, partitionOffsetOnlineLogicalCpus))
	times.SmtThreads = uint32(perfstatInt32(buf, partitionOffsetSmtThreads))
	return nil
}

func decodeMemoryTotal(buf []byte) (aixMemory, error) {
	if err := checkPerfstatSize("perfstat_memory_total", buf, perfstatMemoryTotalSize); err != nil {
		return aixMemory{}, err
	}
	return aixMemory{
		RealTotal:  perfstatUint64(buf, memoryOffsetRealTotal),
		RealFree:   perfstatUint64(buf, memoryOffsetRealFree),
		RealPinned: perfstatUint64(buf, memoryOffsetRealPinned),
		Numperm:    perfstatUint64(buf, memoryOffsetNumperm),
		PgspTotal:  perfstatUint64(buf, memoryOffsetPgspTotal),
		PgspFree:   perfstatUint64(buf, memoryOffsetPgspFree),
		PgspIns:    perfstatUint64(buf, memoryOffsetPgspIns),
		PgspOuts:   perfstatUint64(buf, memoryOffsetPgspOuts),
	}, nil
}

// cpuDelta is the difference of two perfstatTimes as float values
type cpuDelta struct {
	User float64
	Sys  float64
	Idle float64
	Wait float64
}

func (d cpuDelta) total() float64 {
	return d.User + d.Sys + d.Idle + d.Wait
}

// percentOf returns the values as percent of total
func (d cpuDelta) percentOf(total float64) cpuDetails {
	if total <= 0 {
		return cpuDetails{Idle: 100}
	}
	return cpuDetails{
		User:   d.User / total * 100,
		System: d.Sys / total * 100,
		Idle:   d.Idle / total * 100,
		Iowait: d.Wait / total * 100,
	}
}

func counterDelta(prev, current uint64) float64 {
	if current > prev {
		return float64(current - prev)
	}
	return 0
}

func perfstatTimesDelta(prev, current perfstatTimes) cpuDelta {
	return cpuDelta{
		User: counterDelta(prev.User, current.User),
		Sys:  counterDelta(prev.Sys, current.Sys),
		Idle: counterDelta(prev.Idle, current.Idle),
		Wait: counterDelta(prev.Wait, current.Wait),
	}
}

// usagePercentage returns user + sys of details, limited to 0..100
func usagePercentage(details cpuDetails) float64 {
	return math.Min(100, math.Max(0, details.User+details.System))
}

// aixCpuResult calculates the CPU usage of an AIX LPAR between two samples
func aixCpuResult(prev, current *aixCpuTimes) *resultCpu {
	result := &resultCpu{}

	total := aixPartitionDetails(prev, current)
	result.DetailsTotal = &total
	result.PercentageTotal = usagePercentage(total)

	if usage, ok := aixSharedPartitionUsage(prev, current); ok {
		result.Partition = &cpuPartition{
			PhysicalCpus:          usage.PhysicalCpus,
			Entitlement:           current.Entitlement,
			EntitlementPercentage: usage.EntitlementPercentage,
		}
	}

	// The PURR of a logical CPU only increases while the hardware thread is dispatched, so the
	// values per logical CPU show the usage of the thread while it was running.
	result.PercentagePerCore = make([]float64, len(current.Cpus))
	result.DetailsPerCore = make([]cpuDetails, len(current.Cpus))
	for i := range current.Cpus {
		delta := perfstatTimesDelta(prev.Cpus[i], current.Cpus[i])
		details := delta.percentOf(delta.total())
		result.DetailsPerCore[i] = details
		result.PercentagePerCore[i] = usagePercentage(details)
	}
	return result
}

// aixPartitionDetails returns the CPU usage of the partition
func aixPartitionDetails(prev, current *aixCpuTimes) cpuDetails {
	purr := perfstatTimesDelta(prev.Partition, current.Partition)
	purrTotal := purr.total()
	if purrTotal <= 0 {
		// No PURR values, e.g. perfstat_partition_total failed
		ticks := perfstatTimesDelta(prev.Logical, current.Logical)
		return ticks.percentOf(ticks.total())
	}
	if !current.Shared {
		// Dedicated processors: the PURR of all threads covers the whole time
		return purr.percentOf(purrTotal)
	}
	return sharedPartitionDetails(prev, current, purr)
}

// sharedPartitionDetails calculates the usage of a shared processor partition like lparstat:
// the values are relative to the entitlement, the unused entitlement counts as idle and wait.
// An uncapped partition can use more than its entitlement, the values are relative to the used PURR then.
func sharedPartitionDetails(prev, current *aixCpuTimes, purr cpuDelta) cpuDetails {
	purrTotal := purr.total()
	timebase := counterDelta(prev.Timebase, current.Timebase)
	entitledPurr := timebase * current.Entitlement
	if entitledPurr < purrTotal {
		entitledPurr = purrTotal
	}
	unused := entitledPurr - purrTotal

	// The unused entitlement is split between idle and wait like the ticks of the logical CPUs
	ticks := perfstatTimesDelta(prev.Logical, current.Logical)
	waitAndIdle := ticks.Wait + ticks.Idle
	if waitAndIdle > 0 {
		purr.Wait = purr.Wait + unused*ticks.Wait/waitAndIdle
		purr.Idle = purr.Idle + unused*ticks.Idle/waitAndIdle
	} else {
		purr.Idle = purr.Idle + unused
	}
	return purr.percentOf(entitledPurr)
}

// aixPartitionUsage is the processor usage of a shared processor partition
type aixPartitionUsage struct {
	PhysicalCpus          float64 // physc: consumed physical processors
	EntitlementPercentage float64 // %entc: consumed physical processors in percent of the entitlement
	MaxPhysicalCpus       float64 // The physical processors the partition can use at most
	CapacityPercentage    float64 // physc in percent of MaxPhysicalCpus
}

// aixSharedPartitionUsage calculates the processor usage of a shared processor partition. Dedicated partitions
// always own their processors, so ok is false for them.
func aixSharedPartitionUsage(prev, current *aixCpuTimes) (aixPartitionUsage, bool) {
	if !prev.PartitionAvailable || !current.PartitionAvailable || !current.Shared || current.Entitlement <= 0 {
		return aixPartitionUsage{}, false
	}
	timebase := counterDelta(prev.Timebase, current.Timebase)
	if timebase <= 0 {
		return aixPartitionUsage{}, false
	}

	purr := perfstatTimesDelta(prev.Partition, current.Partition)
	usage := aixPartitionUsage{
		PhysicalCpus:    purr.total() / timebase,
		MaxPhysicalCpus: aixMaxPhysicalCpus(current),
	}
	usage.EntitlementPercentage = usage.PhysicalCpus / current.Entitlement * 100
	if usage.MaxPhysicalCpus > 0 {
		usage.CapacityPercentage = math.Min(100, usage.PhysicalCpus/usage.MaxPhysicalCpus*100)
	}
	return usage, true
}

// aixMaxPhysicalCpus returns the processor capacity a shared partition can use at most: a capped partition
// its entitlement, an uncapped partition one physical processor per virtual processor, limited by the size
// of the shared processor pool
func aixMaxPhysicalCpus(times *aixCpuTimes) float64 {
	if times.Capped {
		return times.Entitlement
	}
	maximum := float64(times.VirtualCpus)
	if times.PoolCpus > 0 && float64(times.PoolCpus) < maximum {
		maximum = float64(times.PoolCpus)
	}
	return maximum
}

// resultLpar is the configuration and the processor usage of an AIX LPAR (logical partition)
type resultLpar struct {
	Shared      bool    `json:"shared"`       // Shared processor partition, false for dedicated processors
	Capped      bool    `json:"capped"`       // A capped partition can not use more than its entitlement
	SmtThreads  uint32  `json:"smt_threads"`  // SMT threads per virtual processor
	LogicalCpus uint32  `json:"logical_cpus"` // Online logical CPUs (virtual processors * SMT threads)
	VirtualCpus int32   `json:"virtual_cpus"` // Online virtual processors (dedicated partitions: processors)
	Entitlement float64 `json:"entitlement"`  // Entitled processor units
	PoolCpus    int32   `json:"pool_cpus"`    // Physical processors of the shared processor pool
	MemoryMB    uint64  `json:"memory_mb"`    // Online memory

	// The usage values are only available for shared processor partitions (0 for dedicated partitions)
	PhysicalCpus          float64 `json:"physc"`                  // Consumed physical processors
	EntitlementPercentage float64 `json:"entitlement_percentage"` // physc in percent of the entitlement
	MaxPhysicalCpus       float64 `json:"max_physc"`              // Capped: entitlement, uncapped: virtual processors (at most the pool size)
	CapacityPercentage    float64 `json:"capacity_percentage"`    // physc in percent of max_physc

	// Virtual processor context switches per second. Involuntary context switches happen if the
	// hypervisor preempts a virtual processor, e.g. because the entitlement is used up.
	VirtualCswitches            float64 `json:"virtual_context_switches"`
	InvoluntaryVirtualCswitches float64 `json:"involuntary_virtual_context_switches"`
}

// aixLparResult returns the configuration of the partition and the usage between two samples
func aixLparResult(prev, current *aixCpuTimes) (*resultLpar, error) {
	if !current.PartitionAvailable {
		return nil, fmt.Errorf("perfstat_partition_total is not available")
	}

	result := &resultLpar{
		Shared:      current.Shared,
		Capped:      current.Capped,
		SmtThreads:  current.SmtThreads,
		LogicalCpus: current.LogicalCpus,
		VirtualCpus: current.VirtualCpus,
		Entitlement: current.Entitlement,
		PoolCpus:    current.PoolCpus,
		MemoryMB:    current.MemoryMB,
	}
	if !current.SmtEnabled {
		result.SmtThreads = 1
	}

	if usage, ok := aixSharedPartitionUsage(prev, current); ok {
		result.PhysicalCpus = usage.PhysicalCpus
		result.EntitlementPercentage = usage.EntitlementPercentage
		result.MaxPhysicalCpus = usage.MaxPhysicalCpus
		result.CapacityPercentage = usage.CapacityPercentage
	}

	seconds := current.Time.Sub(prev.Time).Seconds()
	if prev.PartitionAvailable && seconds > 0 {
		voluntary := counterDelta(prev.VoluntaryCswitches, current.VoluntaryCswitches)
		involuntary := counterDelta(prev.InvoluntaryCswitches, current.InvoluntaryCswitches)
		result.VirtualCswitches = (voluntary + involuntary) / seconds
		result.InvoluntaryVirtualCswitches = involuntary / seconds
	}
	return result, nil
}

// aixMemoryResult converts the memory values of perfstat. Used memory is the computational memory
// (like "percentage of memory used for computational pages" of vmstat -v), the file cache is available.
func aixMemoryResult(m aixMemory) *resultMemory {
	total := m.RealTotal * perfstatPageSize
	free := m.RealFree * perfstatPageSize
	cache := m.Numperm * perfstatPageSize

	used := uint64(0)
	if total > free+cache {
		used = total - free - cache
	}

	result := &resultMemory{
		Total:     total,
		Available: total - used,
		Used:      used,
		Free:      free,
		Wired:     m.RealPinned * perfstatPageSize,
	}
	if total > 0 {
		result.Percent = float64(used) / float64(total) * 100
	}
	return result
}

func aixSwapResult(m aixMemory) *resultSwap {
	total := m.PgspTotal * perfstatPageSize
	free := m.PgspFree * perfstatPageSize
	if free > total {
		free = total
	}
	result := &resultSwap{
		Total: total,
		Used:  total - free,
		Free:  free,
		Sin:   m.PgspIns * perfstatPageSize,
		Sout:  m.PgspOuts * perfstatPageSize,
	}
	if total > 0 {
		result.Percent = float64(result.Used) / float64(total) * 100
	}
	return result
}

// parseAixMount parses the output of the AIX mount command:
//
//	  node       mounted        mounted over    vfs       date        options
//	-------- ---------------  ---------------  ------ ------------ ---------------
//	         /dev/hd4         /                jfs2   Apr 30 13:17 rw,log=/dev/hd8
//	server   /export          /mnt             nfs3   Apr 30 13:20 rw,bg,hard
//
// Local file systems have no node. The date always consists of three fields.
func parseAixMount(output string) []aixMount {
	var mounts []aixMount
	for _, line := range strings.Split(output, "\n") {
		if mount, ok := parseAixMountLine(line); ok {
			mounts = append(mounts, mount)
		}
	}
	return mounts
}

func parseAixMountLine(line string) (aixMount, bool) {
	fields := strings.Fields(line)
	if len(fields) < 6 || strings.HasPrefix(fields[0], "---") || fields[0] == "node" {
		return aixMount{}, false
	}

	remote := !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t")
	mount := aixMount{}
	var rest []string
	if remote {
		if len(fields) < 7 {
			return aixMount{}, false
		}
		mount.Device = fields[0] + ":" + fields[1]
		mount.Mountpoint = fields[2]
		mount.Fstype = fields[3]
		rest = fields[4:]
	} else {
		mount.Device = fields[0]
		mount.Mountpoint = fields[1]
		mount.Fstype = fields[2]
		rest = fields[3:]
	}

	// rest: month, day, time and the optional options
	if len(rest) >= 4 {
		mount.Opts = strings.Split(rest[3], ",")
	}
	return mount, true
}
