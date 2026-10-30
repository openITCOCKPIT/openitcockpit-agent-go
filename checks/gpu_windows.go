package checks

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// GPUs on Windows: D3DKMT (gdi32.dll), the interface of the Task Manager, works for the GPUs of all vendors with a
// WDDM driver. NVIDIA GPUs in the TCC mode are no WDDM adapters and are only reported by nvidia-smi.
//
// Documentation of D3DKMT (header d3dkmthk.h of the Windows SDK/WDK):
// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/d3dkmthk/
// The structures below are named like in the header, e.g. d3dkmtAdapterPerfData = D3DKMT_ADAPTER_PERFDATA:
// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/d3dkmthk/ns-d3dkmthk-_d3dkmt_adapter_perfdata

var (
	// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/d3dkmthk/nf-d3dkmthk-d3dkmtenumadapters2
	// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/d3dkmthk/nf-d3dkmthk-d3dkmtqueryadapterinfo
	// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/d3dkmthk/nf-d3dkmthk-d3dkmtcloseadapter
	gdi32DLL                   = windows.NewLazySystemDLL("gdi32.dll")
	procD3DKMTEnumAdapters2    = gdi32DLL.NewProc("D3DKMTEnumAdapters2")
	procD3DKMTQueryAdapterInfo = gdi32DLL.NewProc("D3DKMTQueryAdapterInfo")
	procD3DKMTCloseAdapter     = gdi32DLL.NewProc("D3DKMTCloseAdapter")

	// The performance counters of the GPUs are kept open, as the utilization requires two samples
	gpuPdhMutex   sync.Mutex
	gpuPdhQueries = map[string]*pdhQuery{}

	// Instances of the performance counters "GPU Engine" and "GPU Adapter Memory". They contain the LUID of the
	// adapter (D3DKMT_ADAPTERINFO.AdapterLuid), the index of the physical adapter of linked adapters (phys), the engine
	// (eng) and its type (engtype). The format is not documented, it is the format of Windows 10 and 11, e.g.:
	// pid_10444_luid_0x00000000_0x00013C9C_phys_0_eng_0_engtype_3D
	gpuEngineInstancePattern = regexp.MustCompile(`(?i)luid_(0x[0-9a-f]{8}_0x[0-9a-f]{8})_phys_(\d+)_eng_(\d+)_engtype_(.*)$`)
	// e.g.: luid_0x00000000_0x00013C9C_phys_0
	gpuMemoryInstancePattern = regexp.MustCompile(`(?i)^luid_(0x[0-9a-f]{8}_0x[0-9a-f]{8})_phys_0$`)
)

const (
	// Values of KMTQUERYADAPTERINFOTYPE, the type of D3DKMTQueryAdapterInfo:
	// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/d3dkmthk/ne-d3dkmthk-_kmtqueryadapterinfotype
	// The documentation only lists the names. The values are the position in the enum, but 14 is not used:
	// KMTQAITYPE_DRIVERVERSION = 13 and KMTQAITYPE_ADAPTERTYPE = 15 are set explicitly in d3dkmthk.h, e.g. in the
	// header of Wine: https://gitlab.winehq.org/wine/wine/-/blob/master/include/ddk/d3dkmthk.h
	// All values were verified on Windows 11 (build 26200) with an AMD Radeon RX 7800 XT and a Radeon iGPU.
	kmtqaiGetSegmentSize           = 3  // KMTQAITYPE_GETSEGMENTSIZE: D3DKMT_SEGMENTSIZEINFO
	kmtqaiAdapterAddress           = 6  // KMTQAITYPE_ADAPTERADDRESS: D3DKMT_ADAPTERADDRESS
	kmtqaiAdapterRegistryInfo      = 8  // KMTQAITYPE_ADAPTERREGISTRYINFO: D3DKMT_ADAPTERREGISTRYINFO
	kmtqaiAdapterType              = 15 // KMTQAITYPE_ADAPTERTYPE: D3DKMT_ADAPTERTYPE
	kmtqaiPhysicalAdapterDeviceIDs = 31 // KMTQAITYPE_PHYSICALADAPTERDEVICEIDS: D3DKMT_QUERY_DEVICE_IDS
	kmtqaiKmdDriverVersion         = 49 // KMTQAITYPE_KMD_DRIVER_VERSION: D3DKMT_KMD_DRIVER_VERSION
	kmtqaiNodePerfData             = 61 // KMTQAITYPE_NODEPERFDATA: D3DKMT_NODE_PERFDATA
	kmtqaiAdapterPerfData          = 62 // KMTQAITYPE_ADAPTERPERFDATA: D3DKMT_ADAPTER_PERFDATA
	kmtqaiAdapterPerfDataCaps      = 63 // KMTQAITYPE_ADAPTERPERFDATA_CAPS: D3DKMT_ADAPTER_PERFDATACAPS

	// NTSTATUS of D3DKMTQueryAdapterInfo if the size of the structure does not match
	statusInvalidParameter = 0xC000000D

	// Bits of D3DKMT_ADAPTERTYPE, in the order of the bit field:
	// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/d3dkmthk/ns-d3dkmthk-_d3dkmt_adaptertype
	d3dkmtAdapterTypeSoftwareDevice   = 1 << 2 // e.g. the Microsoft Basic Render Driver
	d3dkmtAdapterTypeHybridIntegrated = 1 << 5 // Integrated GPU of a notebook with hybrid graphics
	d3dkmtAdapterTypeIndirectDisplay  = 1 << 6 // Virtual displays, e.g. of remote desktop software

	// GPUs with less dedicated memory are integrated GPUs, e.g. 128 MiB (Intel) or 512 MiB (AMD, the Radeon iGPU of a
	// Ryzen 7000 reports 485 MiB). Windows has no flag for integrated GPUs, D3DKMT only reports hybrid graphics of
	// notebooks (d3dkmtAdapterTypeHybridIntegrated).
	gpuIntegratedMaxDedicatedMemory = 1 << 30

	// The temperatures of D3DKMT_ADAPTER_PERFDATACAPS are documented in tenths of a degree Celsius, but the AMD driver
	// reports degrees Celsius (110 for the RX 7800 XT). Values below this limit are degrees Celsius.
	d3dkmtCapsTemperatureDeciLimit = 200

	// XPath query for the events of the GPU drivers in the System event log of the last 24 hours, see
	// https://learn.microsoft.com/en-us/windows/win32/wes/consuming-events (Level <= 3: critical, error and warning:
	// https://learn.microsoft.com/en-us/windows/win32/wes/eventmanifestschema-leveltype-complextype)
	// - Display 4101: "Display driver ... stopped responding and has successfully recovered", the timeout detection
	//   and recovery (TDR) of Windows. The message is described in
	//   https://learn.microsoft.com/en-us/windows-hardware/drivers/display/timeout-detection-and-recovery, the event ID
	//   is not documented.
	// - The kernel drivers of the vendors log with the name of their service: nvlddmkm (NVIDIA), amdkmdag, amdkmdap
	//   and amdwddmg (AMD), igfx and igfxn (Intel)
	// 86400000 ms = 24 hours
	gpuDriverEventsQuery = `*[System[((Provider[@Name='Display'] and EventID=4101) or ` +
		`(Provider[@Name='nvlddmkm' or @Name='amdkmdag' or @Name='amdkmdap' or @Name='amdwddmg' or @Name='igfx' or @Name='igfxn'] and Level <= 3)) ` +
		`and TimeCreated[timediff(@SystemTime) <= 86400000]]]`
	gpuMaxDriverEventRecords = 1000
)

// LUID: https://learn.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-luid
type d3dkmtLuid struct {
	LowPart  uint32
	HighPart int32
}

// String returns the LUID in the format of the instances of the performance counters
func (l d3dkmtLuid) String() string {
	return fmt.Sprintf("0x%08x_0x%08x", uint32(l.HighPart), l.LowPart)
}

// D3DKMT_ADAPTERINFO
type d3dkmtAdapterInfo struct {
	Adapter                   uint32
	AdapterLuid               d3dkmtLuid
	NumOfSources              uint32
	PrecisePresentRegionsPref int32
}

// D3DKMT_ENUMADAPTERS2. The pointers get the alignment of the architecture like in C. The uint64 fields of the
// other structures need an explicit padding, as Go aligns them to 4 bytes on 386, but C to 8 bytes (D3DKMT_ALIGN64).
type d3dkmtEnumAdapters2 struct {
	NumAdapters uint32
	Adapters    *d3dkmtAdapterInfo
}

// D3DKMT_QUERYADAPTERINFO
type d3dkmtQueryAdapterInfo struct {
	Adapter  uint32
	Type     int32
	Data     unsafe.Pointer
	DataSize uint32
}

// D3DKMT_ADAPTERREGISTRYINFO, the strings have the length MAX_PATH (260)
type d3dkmtAdapterRegistryInfo struct {
	AdapterString [260]uint16
	BiosString    [260]uint16
	DacType       [260]uint16
	ChipType      [260]uint16
}

// D3DKMT_SEGMENTSIZEINFO, in bytes
type d3dkmtSegmentSizeInfo struct {
	DedicatedVideoMemorySize  uint64
	DedicatedSystemMemorySize uint64
	SharedSystemMemorySize    uint64
}

// D3DKMT_ADAPTERADDRESS: PCI location, the PCI domain is not reported
type d3dkmtAdapterAddress struct {
	BusNumber      uint32
	DeviceNumber   uint32
	FunctionNumber uint32
}

// D3DKMT_QUERY_DEVICE_IDS with D3DKMT_DEVICE_IDS (PCI ids)
type d3dkmtQueryDeviceIDs struct {
	PhysicalAdapterIndex uint32
	VendorID             uint32
	DeviceID             uint32
	SubVendorID          uint32
	SubSystemID          uint32
	RevisionID           uint32
	BusType              uint32
}

// D3DKMT_ADAPTER_PERFDATA (Windows 10 1803, WDDM 2.4 and newer)
type d3dkmtAdapterPerfData struct {
	PhysicalAdapterIndex uint32
	_                    uint32
	MemoryFrequency      uint64 // Hz
	MaxMemoryFrequency   uint64 // Hz
	MaxMemoryFrequencyOC uint64 // Hz
	MemoryBandwidth      uint64 // Bytes per second
	PCIEBandwidth        uint64 // Bytes per second
	FanRPM               uint32
	Power                uint32 // Tenths of a percent
	Temperature          uint32 // Tenths of a degree Celsius
	PowerStateOverride   uint8
	_                    [3]byte
}

// D3DKMT_ADAPTER_PERFDATACAPS (Windows 10 1803, WDDM 2.4 and newer)
type d3dkmtAdapterPerfDataCaps struct {
	PhysicalAdapterIndex uint32
	_                    uint32
	MaxMemoryBandwidth   uint64
	MaxPCIEBandwidth     uint64
	MaxFanRPM            uint32
	TemperatureMax       uint32 // The temperature before damage levels occur, see d3dkmtCapsTemperatureDeciLimit
	TemperatureWarning   uint32 // The temperature where the throttling begins, see d3dkmtCapsTemperatureDeciLimit
	_                    uint32
}

// D3DKMT_NODE_PERFDATA. The documentation contains a field Reserved (uint64) at the end, but Windows 11 (build
// 26200) only accepts the size without this field (56 bytes) and returns STATUS_INVALID_PARAMETER for 64 bytes.
type d3dkmtNodePerfData struct {
	NodeOrdinal          uint32
	PhysicalAdapterIndex uint32
	Frequency            uint64 // Hz
	MaxFrequency         uint64 // Hz, without overclocking
	MaxFrequencyOC       uint64 // Hz
	Voltage              uint32 // mV
	VoltageMax           uint32 // mV
	VoltageMaxOC         uint32 // mV
	_                    uint32
	MaxTransitionLatency uint64 // 100 ns
}

// d3dkmtNodePerfDataReserved is the structure of the documentation with the field Reserved, in case newer Windows
// versions expect it
type d3dkmtNodePerfDataReserved struct {
	d3dkmtNodePerfData
	Reserved uint64
}

// d3dkmtAdapter are the values of an adapter, which are read with D3DKMT
type d3dkmtAdapter struct {
	luid           string
	name           string
	biosVersion    string
	adapterType    uint32
	vendorID       uint32
	deviceID       uint32
	pciAddress     *d3dkmtAdapterAddress
	segments       *d3dkmtSegmentSizeInfo
	driverVersion  string
	perfData       *d3dkmtAdapterPerfData
	perfCaps       *d3dkmtAdapterPerfDataCaps
	graphicsClock  *float64
	graphicsMaxClk *float64
}

var errD3dkmtInvalidParameter = errors.New("STATUS_INVALID_PARAMETER")

func d3dkmtQuery(adapter uint32, infoType int32, data unsafe.Pointer, size uintptr) error {
	query := d3dkmtQueryAdapterInfo{Adapter: adapter, Type: infoType, Data: data, DataSize: uint32(size)}
	status, _, _ := procD3DKMTQueryAdapterInfo.Call(uintptr(unsafe.Pointer(&query)))
	switch uint32(status) {
	case 0:
		return nil
	case statusInvalidParameter:
		return fmt.Errorf("D3DKMTQueryAdapterInfo %d: %w", infoType, errD3dkmtInvalidParameter)
	}
	return fmt.Errorf("D3DKMTQueryAdapterInfo %d: NTSTATUS 0x%08x", infoType, uint32(status))
}

// enumD3dkmtAdapters returns the hardware adapters (GPUs), software adapters like the Microsoft Basic Render Driver
// and indirect display adapters (e.g. of remote desktop software) are skipped
func enumD3dkmtAdapters() ([]*d3dkmtAdapter, error) {
	if err := gdi32DLL.Load(); err != nil {
		return nil, err
	}
	if err := procD3DKMTEnumAdapters2.Find(); err != nil {
		// Windows 8 and newer
		return nil, err
	}

	var enum d3dkmtEnumAdapters2
	if status, _, _ := procD3DKMTEnumAdapters2.Call(uintptr(unsafe.Pointer(&enum))); status != 0 {
		return nil, fmt.Errorf("D3DKMTEnumAdapters2: NTSTATUS 0x%08x", uint32(status))
	}
	if enum.NumAdapters == 0 {
		return nil, nil
	}
	infos := make([]d3dkmtAdapterInfo, enum.NumAdapters)
	enum.Adapters = &infos[0]
	if status, _, _ := procD3DKMTEnumAdapters2.Call(uintptr(unsafe.Pointer(&enum))); status != 0 {
		return nil, fmt.Errorf("D3DKMTEnumAdapters2: NTSTATUS 0x%08x", uint32(status))
	}

	var adapters []*d3dkmtAdapter
	for i := range infos[:enum.NumAdapters] {
		info := &infos[i]
		adapter := readD3dkmtAdapter(info)
		_, _, _ = procD3DKMTCloseAdapter.Call(uintptr(unsafe.Pointer(&info.Adapter)))
		if adapter.adapterType&(d3dkmtAdapterTypeSoftwareDevice|d3dkmtAdapterTypeIndirectDisplay) != 0 {
			continue
		}
		adapters = append(adapters, adapter)
	}
	return adapters, nil
}

func readD3dkmtAdapter(info *d3dkmtAdapterInfo) *d3dkmtAdapter {
	handle := info.Adapter
	adapter := &d3dkmtAdapter{luid: info.AdapterLuid.String()}

	_ = d3dkmtQuery(handle, kmtqaiAdapterType, unsafe.Pointer(&adapter.adapterType), unsafe.Sizeof(adapter.adapterType))

	var registry d3dkmtAdapterRegistryInfo
	if d3dkmtQuery(handle, kmtqaiAdapterRegistryInfo, unsafe.Pointer(&registry), unsafe.Sizeof(registry)) == nil {
		adapter.name = windows.UTF16ToString(registry.AdapterString[:])
		adapter.biosVersion = windows.UTF16ToString(registry.BiosString[:])
	}

	ids := d3dkmtQueryDeviceIDs{}
	if d3dkmtQuery(handle, kmtqaiPhysicalAdapterDeviceIDs, unsafe.Pointer(&ids), unsafe.Sizeof(ids)) == nil {
		adapter.vendorID, adapter.deviceID = ids.VendorID, ids.DeviceID
	}

	// Adapters without PCI location (e.g. the Microsoft Basic Render Driver) report the bus 0xFFFFFFFF
	var address d3dkmtAdapterAddress
	if d3dkmtQuery(handle, kmtqaiAdapterAddress, unsafe.Pointer(&address), unsafe.Sizeof(address)) == nil && address.BusNumber != 0xFFFFFFFF {
		adapter.pciAddress = &address
	}

	var segments d3dkmtSegmentSizeInfo
	if d3dkmtQuery(handle, kmtqaiGetSegmentSize, unsafe.Pointer(&segments), unsafe.Sizeof(segments)) == nil {
		adapter.segments = &segments
	}

	var version uint64
	if d3dkmtQuery(handle, kmtqaiKmdDriverVersion, unsafe.Pointer(&version), unsafe.Sizeof(version)) == nil && version != 0 {
		adapter.driverVersion = fmt.Sprintf("%d.%d.%d.%d", version>>48&0xFFFF, version>>32&0xFFFF, version>>16&0xFFFF, version&0xFFFF)
	}

	// Windows 10 1803 (WDDM 2.4) and newer
	perfData := d3dkmtAdapterPerfData{}
	if d3dkmtQuery(handle, kmtqaiAdapterPerfData, unsafe.Pointer(&perfData), unsafe.Sizeof(perfData)) == nil {
		adapter.perfData = &perfData
	}
	perfCaps := d3dkmtAdapterPerfDataCaps{}
	if d3dkmtQuery(handle, kmtqaiAdapterPerfDataCaps, unsafe.Pointer(&perfCaps), unsafe.Sizeof(perfCaps)) == nil {
		adapter.perfCaps = &perfCaps
	}
	// Node 0 is the 3D engine
	node := d3dkmtNodePerfDataReserved{}
	err := d3dkmtQuery(handle, kmtqaiNodePerfData, unsafe.Pointer(&node), unsafe.Sizeof(node.d3dkmtNodePerfData))
	if errors.Is(err, errD3dkmtInvalidParameter) {
		err = d3dkmtQuery(handle, kmtqaiNodePerfData, unsafe.Pointer(&node), unsafe.Sizeof(node))
	}
	if err == nil && node.MaxFrequency > 0 {
		// Hz to MHz
		clock, maxClock := float64(node.Frequency)/1e6, float64(node.MaxFrequency)/1e6
		adapter.graphicsClock, adapter.graphicsMaxClk = &clock, &maxClock
	}
	return adapter
}

// hasPlatformGpus returns true if Windows has a hardware GPU adapter
func hasPlatformGpus() bool {
	adapters, err := enumD3dkmtAdapters()
	return err == nil && len(adapters) > 0
}

// collectPlatformGpus returns the GPUs of D3DKMT. NVIDIA GPUs which were already reported by nvidia-smi are skipped.
func collectPlatformGpus(ctx context.Context, existing []*resultGpu) []*resultGpu {
	adapters, err := enumD3dkmtAdapters()
	if err != nil || len(adapters) == 0 {
		return nil
	}

	known := map[string]bool{}
	for _, gpu := range existing {
		// The PCI domain is not reported by D3DKMT
		if _, location, found := strings.Cut(gpu.PciBusID, ":"); found {
			known[location] = true
		}
	}

	counters := readGpuPerfCounters(ctx)
	gpus := []*resultGpu{}
	for _, adapter := range adapters {
		gpu := d3dkmtAdapterToGpu(adapter, counters)
		if location := strings.TrimPrefix(gpu.PciBusID, "0000:"); gpu.PciBusID != "" && known[location] {
			continue
		}
		gpus = append(gpus, gpu)
	}
	return gpus
}

func d3dkmtAdapterToGpu(adapter *d3dkmtAdapter, counters map[string]*gpuPerfCounters) *resultGpu {
	gpu := &resultGpu{
		Vendor:            gpuVendorName(adapter.vendorID),
		Source:            "d3dkmt",
		Name:              adapter.name,
		VbiosVersion:      adapter.biosVersion,
		DriverVersion:     adapter.driverVersion,
		Processes:         -1,
		ThrottleReasons:   []string{},
		ClockEventReasons: []string{},
		XidErrors:         []*resultGpuXid{},
		Messages:          []string{},
	}
	if adapter.pciAddress != nil {
		gpu.PciBusID = fmt.Sprintf("0000:%02x:%02x.%x", adapter.pciAddress.BusNumber, adapter.pciAddress.DeviceNumber, adapter.pciAddress.FunctionNumber)
	}
	gpu.Integrated = adapter.adapterType&d3dkmtAdapterTypeHybridIntegrated != 0 ||
		(adapter.segments != nil && adapter.segments.DedicatedVideoMemorySize < gpuIntegratedMaxDedicatedMemory)

	if perf := adapter.perfData; perf != nil {
		if perf.Temperature > 0 {
			// Tenths of a degree Celsius
			temperature := float64(perf.Temperature) / 10
			gpu.TemperatureCelsius = &temperature
		}
		if perf.MemoryFrequency > 0 {
			// Hz to MHz
			clock := float64(perf.MemoryFrequency) / 1e6
			gpu.ClockMemoryMhz = &clock
		}
		// The power of integrated GPUs is the power of the whole processor. The AMD driver does not report the power of
		// dedicated GPUs, it is 0 also under full load (Radeon RX 7800 XT), so 0 is unknown.
		if !gpu.Integrated && perf.Power > 0 {
			// Tenths of a percent
			power := float64(perf.Power) / 10
			gpu.PowerPercent = &power
		}
		if adapter.perfCaps != nil && adapter.perfCaps.MaxFanRPM > 0 {
			// Without the maximum the RPM is unknown, as GPUs without fan report 0
			rpm := float64(perf.FanRPM)
			gpu.FanSpeedRpm = &rpm
		}
		if caps := adapter.perfCaps; caps != nil && caps.MaxFanRPM > 0 {
			// Fans can stop at low temperatures (zero RPM), 0 % is no error
			fan := float64(perf.FanRPM) / float64(caps.MaxFanRPM) * 100
			gpu.FanSpeedPercent = &fan
		}
	}
	// The documentation only describes the temperature as "the main temperature sensors reading", so no headroom is
	// calculated, the thresholds may belong to another sensor (e.g. the hotspot of AMD GPUs)
	if caps := adapter.perfCaps; caps != nil {
		gpu.TemperatureSlowdownCelsius = d3dkmtCapsTemperature(caps.TemperatureWarning)
		gpu.TemperatureShutdownCelsius = d3dkmtCapsTemperature(caps.TemperatureMax)
	}
	gpu.ClockGraphicsMhz, gpu.ClockGraphicsMaxMhz = adapter.graphicsClock, adapter.graphicsMaxClk

	if c := counters[strings.ToLower(adapter.luid)]; c != nil {
		gpu.UtilizationPercent = c.utilization
		gpu.EncoderUtilizationPercent = c.encoder
		gpu.DecoderUtilizationPercent = c.decoder
		if gpu.Integrated {
			// Integrated GPUs mostly use the shared system memory
			if c.dedicatedUsage != nil || c.sharedUsage != nil {
				var used uint64
				if c.dedicatedUsage != nil {
					used += *c.dedicatedUsage
				}
				if c.sharedUsage != nil {
					used += *c.sharedUsage
				}
				gpu.MemoryUsedBytes = &used
			}
		} else {
			gpu.MemoryUsedBytes = c.dedicatedUsage
		}
	}
	if !gpu.Integrated && adapter.segments != nil {
		total := adapter.segments.DedicatedVideoMemorySize
		gpu.MemoryTotalBytes = &total
	}
	return gpu
}

// d3dkmtCapsTemperature converts a temperature of D3DKMT_ADAPTER_PERFDATACAPS to degrees Celsius
func d3dkmtCapsTemperature(value uint32) *float64 {
	if value == 0 {
		return nil
	}
	temperature := float64(value)
	if value >= d3dkmtCapsTemperatureDeciLimit {
		temperature /= 10
	}
	return &temperature
}

// gpuVendorName returns the name of the PCI vendor id, see https://pci-ids.ucw.cz/
func gpuVendorName(vendorID uint32) string {
	switch vendorID {
	case 0x10de:
		return "NVIDIA"
	case 0x1002, 0x1022:
		return "AMD"
	case 0x8086:
		return "Intel"
	case 0x5143:
		return "Qualcomm"
	case 0x1414:
		return "Microsoft"
	case 0:
		return ""
	}
	return fmt.Sprintf("0x%04x", vendorID)
}

// gpuPerfCounters are the values of the performance counters of a GPU adapter
type gpuPerfCounters struct {
	utilization    *float64
	encoder        *float64
	decoder        *float64
	dedicatedUsage *uint64
	sharedUsage    *uint64
}

// Performance counters of the GPUs, the same values as in the Task Manager:
// https://devblogs.microsoft.com/directx/gpus-in-the-task-manager/
const (
	gpuCounterEngine          = `\GPU Engine(*)\Utilization Percentage`
	gpuCounterDedicatedMemory = `\GPU Adapter Memory(*)\Dedicated Usage`
	gpuCounterSharedMemory    = `\GPU Adapter Memory(*)\Shared Usage`
)

// readGpuPerfCounters returns the utilization and memory of the GPUs by the LUID of the adapter
func readGpuPerfCounters(ctx context.Context) map[string]*gpuPerfCounters {
	gpuPdhMutex.Lock()
	defer gpuPdhMutex.Unlock()

	values := map[string][]*resultPerfCounterValue{}
	var newQueries []*pdhQuery
	for _, counter := range []string{gpuCounterEngine, gpuCounterDedicatedMemory, gpuCounterSharedMemory} {
		query, ok := gpuPdhQueries[counter]
		if !ok {
			var err error
			if query, err = openPdhQuery(counter); err != nil {
				continue
			}
			gpuPdhQueries[counter] = query
		}
		if err := query.collect(); err != nil {
			continue
		}
		if query.samples == 1 {
			newQueries = append(newQueries, query)
		}
	}
	// The utilization requires a second sample
	if len(newQueries) > 0 {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(pdhFirstSampleWait):
		}
		for _, query := range newQueries {
			_ = query.collect()
		}
	}
	for counter, query := range gpuPdhQueries {
		if v, err := query.values(); err == nil {
			values[counter] = v
		}
	}
	return parseGpuPerfCounters(values[gpuCounterEngine], values[gpuCounterDedicatedMemory], values[gpuCounterSharedMemory])
}

// parseGpuPerfCounters calculates the utilization like the Task Manager: the utilization of an engine is the sum of
// all processes, the utilization of the GPU is the utilization of the busiest engine ("we opted to pick the
// percentage utilization of the busiest engine as a representative of the overall GPU usage",
// https://devblogs.microsoft.com/directx/gpus-in-the-task-manager/)
func parseGpuPerfCounters(engines, dedicated, shared []*resultPerfCounterValue) map[string]*gpuPerfCounters {
	result := map[string]*gpuPerfCounters{}
	get := func(luid string) *gpuPerfCounters {
		luid = strings.ToLower(luid)
		if result[luid] == nil {
			result[luid] = &gpuPerfCounters{}
		}
		return result[luid]
	}

	type engine struct {
		luid, engineType string
	}
	engineUtilization := map[engine]float64{}
	for _, value := range engines {
		match := gpuEngineInstancePattern.FindStringSubmatch(value.Instance)
		if match == nil || match[2] != "0" {
			continue
		}
		engineUtilization[engine{luid: strings.ToLower(match[1]), engineType: match[3] + "_" + match[4]}] += value.Value
	}
	maximum := func(current *float64, value float64) *float64 {
		value = min(value, 100)
		if current == nil || value > *current {
			return &value
		}
		return current
	}
	for e, utilization := range engineUtilization {
		counters := get(e.luid)
		counters.utilization = maximum(counters.utilization, utilization)
		engineType := strings.ToLower(e.engineType)
		switch {
		case strings.Contains(engineType, "encode"):
			counters.encoder = maximum(counters.encoder, utilization)
		case strings.Contains(engineType, "decode"):
			counters.decoder = maximum(counters.decoder, utilization)
		case strings.Contains(engineType, "codec"):
			// AMD: one engine for encoding and decoding (VCN), e.g. "Video Codec 0"
			counters.encoder = maximum(counters.encoder, utilization)
			counters.decoder = maximum(counters.decoder, utilization)
		}
	}

	for _, value := range dedicated {
		if match := gpuMemoryInstancePattern.FindStringSubmatch(value.Instance); match != nil && value.Value >= 0 {
			bytes := uint64(value.Value)
			get(match[1]).dedicatedUsage = &bytes
		}
	}
	for _, value := range shared {
		if match := gpuMemoryInstancePattern.FindStringSubmatch(value.Instance); match != nil && value.Value >= 0 {
			bytes := uint64(value.Value)
			get(match[1]).sharedUsage = &bytes
		}
	}
	return result
}

// readGpuDriverEvents returns the TDR events and the errors of the GPU drivers of the last 24 hours
func readGpuDriverEvents(ctx context.Context) *resultGpuDriverEventList {
	result := &resultGpuDriverEventList{Events: []*resultGpuDriverEvent{}}
	if err := wevtapiDLL.Load(); err != nil {
		result.Error = err.Error()
		return result
	}
	renderContext, _, err := procEvtCreateRenderContext.Call(0, 0, evtRenderContextSystem)
	if renderContext == 0 {
		result.Error = fmt.Sprintf("EvtCreateRenderContext: %s", err)
		return result
	}
	defer evtClose(windows.Handle(renderContext))

	// The publisher metadata of the event log check is not thread safe
	evtStatesMutex.Lock()
	defer evtStatesMutex.Unlock()

	records, err := evtReadRecords(ctx, "System", gpuDriverEventsQuery, windows.Handle(renderContext), &evtRenderBuffer{}, gpuMaxDriverEventRecords)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	// Newest records first
	result.Count = len(records)
	for _, record := range records {
		if len(result.Events) >= gpuMaxXidErrors {
			break
		}
		result.Events = append(result.Events, &resultGpuDriverEvent{
			Timestamp: record.TimeGenerated,
			Source:    record.Source,
			EventID:   record.EventID,
			EntryType: record.EntryType,
			Message:   record.Message,
		})
	}
	return result
}

// countNvidiaPciDevices is only implemented on Linux
func countNvidiaPciDevices() int {
	return -1
}

// readGpuDriverLogs returns the events of the GPU drivers in the event log, the errors of the NVIDIA driver are
// reported as driver events (no Xid errors). The time window of 24 hours is part of gpuDriverEventsQuery.
func readGpuDriverLogs(ctx context.Context, since time.Time) *gpuDriverLogs {
	return &gpuDriverLogs{events: readGpuDriverEvents(ctx)}
}
