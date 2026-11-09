package checks

import (
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Intel GPUs on Linux (drivers i915 and xe) are read from sysfs:
// - i915: frequencies and throttle reasons of the GT (cardN/gt_*_freq_mhz, cardN/gt/gtN/throttle_reason_*),
//   defined in drivers/gpu/drm/i915/gt/intel_gt_sysfs_pm.c
// - xe: frequencies and throttle reasons (device/tileN/gtN/freq0/...): https://docs.kernel.org/gpu/xe/xe_gt_freq.html
// - hwmon of dedicated GPUs (power, energy, temperatures, fans):
//   https://www.kernel.org/doc/Documentation/ABI/testing/sysfs-driver-intel-i915-hwmon
//   https://www.kernel.org/doc/Documentation/ABI/testing/sysfs-driver-intel-xe-hwmon
// - The utilization is calculated from the busy time of the engines of the DRM clients (gpu_drm_clients_linux.go).
//   The RC6 residency (time in the power saving state) is no utilization: under a light load the GPU does not enter
//   RC6 at all, although it runs at the minimum frequency.

// The integrated GPU of Intel processors always has this PCI location
const intelIntegratedGpuPciBusID = "0000:00:02.0"

// Throttle reasons of Intel GPUs which indicate a problem. The power limits (pl1, pl2, pl4) and current limits
// (vr_tdc, iccmax) also reduce the frequency in normal operation.
var intelProblemThrottleReasons = map[string]bool{
	"thermal": true, "prochot": true, "ratl": true, "vr_thermalert": true, "vr_thermal": true, "soc_thermal": true,
	"soc_avg_thermal": true, "mem_thermal": true, "psys_crit": true,
}

var (
	// Energy counter of the dedicated GPUs (hwmon energy1_input) of the previous check run by PCI bus id
	intelEnergyMutex    sync.Mutex
	intelEnergyPrevious = map[string]intelEnergySample{}
)

type intelEnergySample struct {
	time        time.Time
	microjoules float64
}

func intelGpuDevices() []string {
	return drmGpuDevices("i915", "xe")
}

func hasIntelGpus() bool {
	return len(intelGpuDevices()) > 0
}

func collectIntelGpus(usage map[string]*drmUsage, hasPrevious bool) []*resultGpu {
	gpus := []*resultGpu{}
	for _, device := range intelGpuDevices() {
		gpus = append(gpus, readIntelGpu(device, usage, hasPrevious))
	}
	return gpus
}

func readIntelGpu(device string, usage map[string]*drmUsage, hasPrevious bool) *resultGpu {
	card := filepath.Dir(device)
	driver := ""
	if path, err := filepath.EvalSymlinks(filepath.Join(device, "driver")); err == nil {
		driver = filepath.Base(path)
	}
	gpu := &resultGpu{
		Vendor:            "Intel",
		Source:            driver,
		Name:              gpuPciName("8086", "Intel", readSysfsString(filepath.Join(device, "device"))),
		Processes:         0,
		ThrottleReasons:   []string{},
		ClockEventReasons: []string{},
		XidErrors:         []*resultGpuXid{},
		Messages:          []string{},
	}
	if path, err := filepath.EvalSymlinks(device); err == nil {
		gpu.PciBusID = filepath.Base(path)
	}
	gpu.Integrated = gpu.PciBusID == intelIntegratedGpuPciBusID

	if u := usage[gpu.PciBusID]; u != nil {
		gpu.UtilizationPercent, gpu.EncoderUtilizationPercent, gpu.DecoderUtilizationPercent = u.utilization, u.encoder, u.decoder
		gpu.Processes = u.processes
		if !gpu.Integrated {
			// The memory of the GPU used by the clients (drm-resident-local0). The size of the memory is only
			// available with an ioctl (DRM_I915_QUERY_MEMORY_REGIONS), so memory_total_bytes stays unknown.
			gpu.MemoryUsedBytes = u.vramUsed
		}
	}
	if gpu.UtilizationPercent == nil && hasPrevious {
		// No client was running during both check runs
		if u := usage[gpu.PciBusID]; u == nil || u.processes == 0 {
			idle := 0.0
			gpu.UtilizationPercent = &idle
		}
	}

	switch driver {
	case "i915":
		// The actual frequency is 0 while the GPU is in the power saving state RC6
		gpu.ClockGraphicsMhz = readSysfsNumber(filepath.Join(card, "gt_act_freq_mhz"), 1)
		gpu.ClockGraphicsMaxMhz = readSysfsNumber(filepath.Join(card, "gt_RP0_freq_mhz"), 1)
		readIntelThrottleReasons(filepath.Join(card, "gt", "gt0"), "throttle_reason_", gpu)
	case "xe":
		freq := filepath.Join(device, "tile0", "gt0", "freq0")
		gpu.ClockGraphicsMhz = readSysfsNumber(filepath.Join(freq, "act_freq"), 1)
		gpu.ClockGraphicsMaxMhz = readSysfsNumber(filepath.Join(freq, "rp0_freq"), 1)
		readIntelThrottleReasons(filepath.Join(freq, "throttle"), "reason_", gpu)
	}

	if !gpu.Integrated {
		// The link of integrated GPUs is reported as "Unknown" with the width 0 or 255
		readIntelPcieLink(device, gpu)
	}

	// Only dedicated GPUs have a hwmon device, the device level hwmon has the name of the driver
	hwmons, _ := filepath.Glob(filepath.Join(device, "hwmon", "hwmon*"))
	for _, hwmon := range hwmons {
		if readSysfsString(filepath.Join(hwmon, "name")) == driver {
			readIntelHwmon(hwmon, driver, gpu)
			break
		}
	}
	return gpu
}

// readIntelPcieLink reads the PCIe link of a dedicated GPU. Arc GPUs contain PCIe bridges and report the internal link
// of the GPU (2.5 GT/s x1), the link of the card is the link of the upstream bridge (at most 2 levels up).
func readIntelPcieLink(device string, gpu *resultGpu) {
	path, err := filepath.EvalSymlinks(device)
	if err != nil {
		path = device
	}
	for level := 0; ; level++ {
		gpu.PcieLinkGen = pcieGeneration(readSysfsString(filepath.Join(path, "current_link_speed")))
		gpu.PcieLinkGenMax = pcieGeneration(readSysfsString(filepath.Join(path, "max_link_speed")))
		gpu.PcieLinkWidth = readSysfsNumber(filepath.Join(path, "current_link_width"), 1)
		gpu.PcieLinkWidthMax = readSysfsNumber(filepath.Join(path, "max_link_width"), 1)
		internal := gpu.PcieLinkGenMax != nil && *gpu.PcieLinkGenMax == 1 && gpu.PcieLinkWidthMax != nil && *gpu.PcieLinkWidthMax == 1
		parent := filepath.Dir(path)
		// 0x0604: PCI bridge
		if !internal || level >= 2 || readSysfsString(filepath.Join(parent, "vendor")) != "0x8086" ||
			!strings.HasPrefix(readSysfsString(filepath.Join(parent, "class")), "0x0604") {
			return
		}
		path = parent
	}
}

// readIntelThrottleReasons reads the active throttle reasons, e.g. throttle_reason_pl1 (i915) or reason_pl1 (xe). The
// files contain 1 if the reason is active.
func readIntelThrottleReasons(dir, prefix string, gpu *resultGpu) {
	files, _ := filepath.Glob(filepath.Join(dir, prefix+"*"))
	sort.Strings(files)
	for _, file := range files {
		name := strings.TrimPrefix(filepath.Base(file), prefix)
		// status is set if any reason is active
		if name == "status" || readSysfsString(file) != "1" {
			continue
		}
		gpu.ClockEventReasons = append(gpu.ClockEventReasons, name)
		if intelProblemThrottleReasons[name] {
			gpu.ThrottleReasons = append(gpu.ThrottleReasons, name)
		}
	}
}

// readIntelHwmon reads the sensors of dedicated GPUs (microjoules, microwatts, millidegrees). The power is calculated
// from the energy counter of the previous check run.
func readIntelHwmon(hwmon, driver string, gpu *resultGpu) {
	switch driver {
	case "i915":
		gpu.TemperatureCelsius = readSysfsNumber(filepath.Join(hwmon, "temp1_input"), 1000)
	case "xe":
		// temp2: package, temp3: VRAM
		gpu.TemperatureCelsius = readSysfsNumber(filepath.Join(hwmon, "temp2_input"), 1000)
		gpu.TemperatureMemoryCelsius = readSysfsNumber(filepath.Join(hwmon, "temp3_input"), 1000)
		gpu.TemperatureSlowdownCelsius = readSysfsNumber(filepath.Join(hwmon, "temp2_crit"), 1000)
		gpu.TemperatureShutdownCelsius = readSysfsNumber(filepath.Join(hwmon, "temp2_emergency"), 1000)
		if gpu.TemperatureCelsius != nil && gpu.TemperatureSlowdownCelsius != nil {
			headroom := *gpu.TemperatureSlowdownCelsius - *gpu.TemperatureCelsius
			gpu.TemperatureHeadroomCelsius = &headroom
		}
	}

	gpu.FanSpeedRpm = readSysfsNumber(filepath.Join(hwmon, "fan1_input"), 1)

	// Power limit: sustained limit (PL1) of the card (power1), on xe of the package (power2) if the card has none,
	// otherwise the default limit (rated_max). 0 means disabled.
	for _, file := range []string{"power1_max", "power2_max", "power1_rated_max", "power2_rated_max"} {
		if limit := readSysfsNumber(filepath.Join(hwmon, file), 1000000); limit != nil && *limit > 0 {
			gpu.PowerLimitWatts = limit
			break
		}
	}

	// Energy of the card (energy1), on xe of the package (energy2) if the card has none
	energy := readSysfsNumber(filepath.Join(hwmon, "energy1_input"), 1)
	if energy == nil {
		energy = readSysfsNumber(filepath.Join(hwmon, "energy2_input"), 1)
	}
	if energy == nil {
		return
	}
	now := gpuNow()
	intelEnergyMutex.Lock()
	defer intelEnergyMutex.Unlock()
	previous, ok := intelEnergyPrevious[gpu.PciBusID]
	intelEnergyPrevious[gpu.PciBusID] = intelEnergySample{time: now, microjoules: *energy}
	if elapsed := now.Sub(previous.time).Seconds(); ok && elapsed > 0 && *energy >= previous.microjoules {
		watts := (*energy - previous.microjoules) / 1000000 / elapsed
		gpu.PowerWatts = &watts
	}
}
