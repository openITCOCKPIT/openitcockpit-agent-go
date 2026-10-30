package checks

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
	"github.com/shirou/gopsutil/v4/mem"
)

const (
	gpuCommandTimeout = 30 * time.Second

	// Xid errors of the kernel log (Linux) are reported for this time window
	gpuXidWindow = 24 * time.Hour
	// Newest Xid errors which are reported per GPU, the number of all Xid errors of the time window is in xid_count
	gpuMaxXidErrors = 10
)

// CheckGpu gathers the state of the GPUs: NVIDIA with nvidia-smi, other GPUs with the amdgpu driver (Linux) or the
// generic interfaces of Windows (D3DKMT and performance counters)
type CheckGpu struct {
	nvidiaSmi string
	// GPUs of the platform specific source: amdgpu (Linux) or D3DKMT (Windows)
	platform bool
}

// Name will be used in the response as check name
func (c *CheckGpu) Name() string {
	return "gpus"
}

type resultGpus struct {
	Gpus         []*resultGpu              `json:"gpus"`
	Nvidia       *resultGpusNvidia         `json:"nvidia,omitempty"`
	Xid          *resultGpuXidState        `json:"xid,omitempty"`
	DriverEvents *resultGpuDriverEventList `json:"driver_events,omitempty"` // Windows event log, Linux kernel log (amdgpu)
}

type resultGpusNvidia struct {
	DriverVersion string `json:"driver_version"`
	CudaVersion   string `json:"cuda_version"`
	AttachedGpus  int    `json:"attached_gpus"` // Number of GPUs reported by nvidia-smi
	PciDevices    int    `json:"pci_devices"`   // Number of NVIDIA GPUs on the PCI bus (Linux), -1 if unknown. Higher than attached_gpus if GPUs have fallen off the bus
	Error         string `json:"error"`         // nvidia-smi failed
}

type resultGpuXidState struct {
	Error string `json:"error"` // The kernel log could not be read
	// Xid errors which could not be assigned to a GPU, e.g. of a GPU that has fallen off the bus, and SXid errors of NVSwitches
	Unassigned []*resultGpuXid `json:"unassigned"`
}

// resultGpuXid is an Xid error of the NVIDIA driver. The meaning and the recommended action of the codes are listed in
// https://docs.nvidia.com/deploy/xid-errors/analyzing-xid-catalog.html
type resultGpuXid struct {
	Timestamp int64  `json:"timestamp"`  // Unix timestamp
	Type      string `json:"type"`       // Xid (GPU) or SXid (NVSwitch)
	Code      int64  `json:"code"`       // e.g. 79 = GPU has fallen off the bus
	PciBusID  string `json:"pci_bus_id"` // e.g. 0000:3b:00.0
	Message   string `json:"message"`
}

// resultGpuDriverEventList are the events of the GPU drivers of the last 24 hours: the Windows event log or the
// messages of the amdgpu driver in the Linux kernel log (GPU hangs and resets)
type resultGpuDriverEventList struct {
	Error  string                  `json:"error"` // The event log could not be read
	Count  int                     `json:"count"`
	Events []*resultGpuDriverEvent `json:"events"` // Newest events
}

type resultGpuDriverEvent struct {
	Timestamp int64  `json:"timestamp"`  // Unix timestamp
	Source    string `json:"source"`     // Display (4101: the driver stopped responding and has recovered, TDR), nvlddmkm (NVIDIA), amdkmdag (AMD), igfx (Intel), amdgpu (Linux)
	EventID   int64  `json:"event_id"`   // 0 on Linux
	EntryType int64  `json:"entry_type"` // 1 = error, 2 = warning
	PciBusID  string `json:"pci_bus_id"` // Linux, e.g. 0000:03:00.0, empty if the message does not contain it
	Message   string `json:"message"`
}

// gpuDriverLogs are the errors of the GPU drivers which are read from the logs of the operating system
type gpuDriverLogs struct {
	xids     []*resultGpuXid // Linux: Xid errors of the NVIDIA driver, nil if not supported
	xidError error           // Linux: the kernel log could not be read
	events   *resultGpuDriverEventList
}

type resultGpu struct {
	// Position in the list: NVIDIA GPUs (nvidia-smi) first, then the GPUs of the platform (amdgpu or D3DKMT). It
	// changes if a GPU or a source is missing, so GPUs have to be identified by pci_bus_id or uuid.
	Index            int    `json:"index"`
	Vendor           string `json:"vendor"` // NVIDIA, AMD, Intel, ...
	Source           string `json:"source"` // nvidia-smi, amdgpu or d3dkmt
	Name             string `json:"name"`
	PciBusID         string `json:"pci_bus_id"` // e.g. 0000:18:00.0
	UUID             string `json:"uuid"`
	Serial           string `json:"serial"`
	VbiosVersion     string `json:"vbios_version"`
	DriverVersion    string `json:"driver_version"`
	Integrated       bool   `json:"integrated"`        // Integrated GPU or unified memory (e.g. DGX Spark), the GPU uses the system memory
	PerformanceState string `json:"performance_state"` // P0 (maximum performance) to P12 (minimum)

	// NVIDIA
	PersistenceMode string `json:"persistence_mode,omitempty"` // Enabled or Disabled (Linux)
	DriverModel     string `json:"driver_model,omitempty"`     // WDDM, TCC or MCDM (Windows)
	MigMode         string `json:"mig_mode,omitempty"`         // Enabled or Disabled
	ComputeMode     string `json:"compute_mode,omitempty"`     // Default, Exclusive_Process or Prohibited
	FabricState     string `json:"fabric_state,omitempty"`     // NVSwitch systems: Not Started, In Progress or Completed
	FabricStatus    string `json:"fabric_status,omitempty"`    // NVSwitch systems: Success or the error
	RecoveryAction  string `json:"recovery_action,omitempty"`  // None, Reset, Reboot, Drain P2P or Drain and Reset

	UtilizationPercent        *float64 `json:"utilization_percent"`
	MemoryUtilizationPercent  *float64 `json:"memory_utilization_percent"` // Time the memory was read or written
	EncoderUtilizationPercent *float64 `json:"encoder_utilization_percent"`
	DecoderUtilizationPercent *float64 `json:"decoder_utilization_percent"`

	// Integrated GPUs: memory_total_bytes is the system memory, memory_used_bytes the memory used by the GPU processes
	MemoryTotalBytes *uint64 `json:"memory_total_bytes"`
	MemoryUsedBytes  *uint64 `json:"memory_used_bytes"`

	TemperatureCelsius         *float64 `json:"temperature_celsius"`
	TemperatureHotspotCelsius  *float64 `json:"temperature_hotspot_celsius"` // AMD: junction
	TemperatureMemoryCelsius   *float64 `json:"temperature_memory_celsius"`
	TemperatureSlowdownCelsius *float64 `json:"temperature_slowdown_celsius"` // The GPU slows down at this temperature
	TemperatureShutdownCelsius *float64 `json:"temperature_shutdown_celsius"` // The GPU shuts down at this temperature
	// Degrees until the GPU slows down because of its temperature. Newer NVIDIA drivers only report this value
	// (T.Limit) and no absolute thresholds.
	TemperatureHeadroomCelsius *float64 `json:"temperature_headroom_celsius"`

	FanSpeedPercent *float64 `json:"fan_speed_percent"` // 0 is normal for fans which stop at low temperatures
	FanSpeedRpm     *float64 `json:"fan_speed_rpm"`     // AMD, Intel and Windows (D3DKMT)
	PowerWatts      *float64 `json:"power_watts"`       // AMD APUs only report the power of the whole processor, which is not reported
	PowerLimitWatts *float64 `json:"power_limit_watts"`
	PowerPercent    *float64 `json:"power_percent"` // Windows (D3DKMT): power in percent of the maximum, the power in watts is unknown

	ClockGraphicsMhz    *float64 `json:"clock_graphics_mhz"`
	ClockGraphicsMaxMhz *float64 `json:"clock_graphics_max_mhz"`
	ClockMemoryMhz      *float64 `json:"clock_memory_mhz"`

	// Active reasons which slow the GPU down because of a problem. NVIDIA: hw_slowdown, hw_thermal_slowdown,
	// hw_power_brake_slowdown, sw_thermal_slowdown. AMD: vrhot0, vrhot1, prochot_cpu, prochot_gfx
	ThrottleReasons []string `json:"throttle_reasons"`
	// All active clock event reasons, also the normal ones like gpu_idle or sw_power_cap (NVIDIA) or ppt0 and
	// temp_hotspot (AMD, the power management firmware also limits the clocks by the temperature in normal operation)
	ClockEventReasons []string `json:"clock_event_reasons"`

	EccEnabled              *bool   `json:"ecc_enabled"`
	EccCorrectedVolatile    *uint64 `json:"ecc_corrected_volatile"` // Since the last driver load
	EccUncorrectedVolatile  *uint64 `json:"ecc_uncorrected_volatile"`
	EccCorrectedAggregate   *uint64 `json:"ecc_corrected_aggregate"` // Lifetime of the GPU
	EccUncorrectedAggregate *uint64 `json:"ecc_uncorrected_aggregate"`

	RetiredPagesSingleBit *uint64 `json:"retired_pages_single_bit"` // Pascal and Volta
	RetiredPagesDoubleBit *uint64 `json:"retired_pages_double_bit"`
	RetiredPagesPending   *bool   `json:"retired_pages_pending"` // A reset is required to retire the pages

	RemappedRowsCorrectable   *uint64 `json:"remapped_rows_correctable"` // Ampere and newer
	RemappedRowsUncorrectable *uint64 `json:"remapped_rows_uncorrectable"`
	RemappedRowsPending       *bool   `json:"remapped_rows_pending"` // A reset is required to remap the rows
	RemappedRowsFailure       *bool   `json:"remapped_rows_failure"` // Remapping failed, the GPU should be replaced
	MemoryRepairPending       *bool   `json:"memory_repair_pending"` // Blackwell: channel or TPC repair pending, a reset is required
	MemoryUnrepairable        *bool   `json:"memory_unrepairable"`

	PcieLinkGen      *float64 `json:"pcie_link_gen"` // GPUs reduce the link speed when idle
	PcieLinkGenMax   *float64 `json:"pcie_link_gen_max"`
	PcieLinkWidth    *float64 `json:"pcie_link_width"`
	PcieLinkWidthMax *float64 `json:"pcie_link_width_max"`

	Processes int `json:"processes"` // Number of processes using the GPU, -1 if unknown

	XidCount  int             `json:"xid_count"`  // Xid errors of the last 24 hours
	XidErrors []*resultGpuXid `json:"xid_errors"` // Newest Xid errors of the last 24 hours

	Messages []string `json:"messages"` // e.g. GPU requires reset
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckGpu) Run(ctx context.Context) (interface{}, error) {
	result := &resultGpus{Gpus: []*resultGpu{}}

	if c.nvidiaSmi != "" {
		nvidia := &resultGpusNvidia{PciDevices: countNvidiaPciDevices()}
		gpus, err := runNvidiaSmi(ctx, c.nvidiaSmi, nvidia)
		if err != nil {
			nvidia.Error = err.Error()
		}
		result.Nvidia = nvidia
		result.Gpus = append(result.Gpus, gpus...)
	}

	if c.platform {
		// NVIDIA GPUs which are already reported by nvidia-smi are skipped
		result.Gpus = append(result.Gpus, collectPlatformGpus(ctx, result.Gpus)...)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	for _, gpu := range result.Gpus {
		if gpu.Integrated && gpu.MemoryTotalBytes == nil {
			if vm, err := mem.VirtualMemoryWithContext(ctx); err == nil {
				total := vm.Total
				gpu.MemoryTotalBytes = &total
			}
		}
	}

	logs := readGpuDriverLogs(ctx, time.Now().Add(-gpuXidWindow))
	result.DriverEvents = logs.events
	if logs.xidError != nil {
		result.Xid = &resultGpuXidState{Error: logs.xidError.Error(), Unassigned: []*resultGpuXid{}}
	} else if logs.xids != nil {
		result.Xid = &resultGpuXidState{Unassigned: assignXidErrors(result.Gpus, logs.xids)}
	}

	for i, gpu := range result.Gpus {
		gpu.Index = i
	}
	return result, nil
}

// Configure the command or return false if the command was disabled
func (c *CheckGpu) Configure(config *config.Configuration) (bool, error) {
	if !config.Gpus {
		return false, nil
	}
	// Only enable the check if a GPU driver is installed
	c.nvidiaSmi = findNvidiaSmi()
	c.platform = hasPlatformGpus()
	return c.nvidiaSmi != "" || c.platform, nil
}

// assignXidErrors adds the Xid errors to the GPUs with the same PCI bus id and returns the remaining ones
func assignXidErrors(gpus []*resultGpu, xids []*resultGpuXid) []*resultGpuXid {
	sort.SliceStable(xids, func(i, j int) bool { return xids[i].Timestamp > xids[j].Timestamp })

	byPciBusID := map[string]*resultGpu{}
	for _, gpu := range gpus {
		byPciBusID[gpu.PciBusID] = gpu
	}
	unassigned := []*resultGpuXid{}
	for _, xid := range xids {
		gpu := byPciBusID[xid.PciBusID]
		if gpu == nil || xid.Type != "Xid" {
			if len(unassigned) < gpuMaxXidErrors {
				unassigned = append(unassigned, xid)
			}
			continue
		}
		gpu.XidCount++
		if len(gpu.XidErrors) < gpuMaxXidErrors {
			gpu.XidErrors = append(gpu.XidErrors, xid)
		}
	}
	return unassigned
}

// normalizePciBusID converts the PCI bus ids of nvidia-smi (00000000:18:00.0) and of the kernel log (0000:18:00)
// to the format of sysfs (0000:18:00.0)
func normalizePciBusID(id string) string {
	var domain, bus, device, function uint64
	if n, _ := fmt.Sscanf(id, "%x:%x:%x.%x", &domain, &bus, &device, &function); n < 3 {
		return id
	}
	return fmt.Sprintf("%04x:%02x:%02x.%x", domain, bus, device, function)
}
