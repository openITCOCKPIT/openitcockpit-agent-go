package checks

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// AMD GPUs on Linux are read from the sysfs files of the amdgpu driver:
// - Device files (gpu_busy_percent, mem_busy_percent, gpu_metrics, the hwmon sensors):
//   https://docs.kernel.org/gpu/amdgpu/thermal.html
// - Memory, product_name, serial_number, unique_id: https://docs.kernel.org/gpu/amdgpu/driver-misc.html
// - RAS error counters: https://docs.kernel.org/gpu/amdgpu/ras.html
// - Generic hwmon interface: https://docs.kernel.org/hwmon/sysfs-interface.html
// - PCI files (vendor, device, class, current_link_speed): https://www.kernel.org/doc/Documentation/ABI/testing/sysfs-bus-pci

var (
	gpuDrmPath     = "/sys/class/drm"
	gpuPciIdsPaths = []string{"/usr/share/misc/pci.ids", "/usr/share/hwdata/pci.ids", "/usr/share/pci.ids"}
	// Product names of AMD GPUs by device and revision id of libdrm (package libdrm-common or libdrm), e.g. "747E, C8,
	// AMD Radeon RX 7800 XT": https://gitlab.freedesktop.org/mesa/drm/-/blob/main/data/amdgpu.ids
	gpuAmdgpuIdsPaths = []string{"/usr/share/libdrm/amdgpu.ids", "/opt/amdgpu/share/libdrm/amdgpu.ids"}

	drmCardPattern = regexp.MustCompile(`^card\d+$`)

	pciIdsMutex sync.Mutex
	pciIdsCache = map[string]string{}
)

// Bits of indep_throttle_status of gpu_metrics (ASIC independent), SMU_THROTTLER_*_BIT in
// https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/drivers/gpu/drm/amd/pm/swsmu/inc/amdgpu_smu.h
var amdThrottleBits = map[int]string{
	0: "ppt0", 1: "ppt1", 2: "ppt2", 3: "ppt3", 4: "spl", 5: "fppt", 6: "sppt", 7: "sppt_apu",
	16: "tdc_gfx", 17: "tdc_soc", 18: "tdc_mem", 19: "tdc_vdd", 20: "tdc_cvip", 21: "edc_cpu", 22: "edc_gfx", 23: "apcc",
	32: "temp_gpu", 33: "temp_core", 34: "temp_mem", 35: "temp_edge", 36: "temp_hotspot", 37: "temp_soc",
	38: "temp_vr_gfx", 39: "temp_vr_soc", 40: "temp_vr_mem0", 41: "temp_vr_mem1", 42: "temp_liquid0", 43: "temp_liquid1",
	44: "vrhot0", 45: "vrhot1", 46: "prochot_cpu", 47: "prochot_gfx", 56: "ppm", 57: "fit",
}

// Throttle reasons of AMD GPUs which indicate a problem. The power management firmware also limits the clocks by the
// power and temperature limits in normal operation, e.g. temp_hotspot is active under load at 50 °C, so these are
// only reported as clock event reasons. Too high temperatures are detected by temperature_headroom_celsius.
var amdProblemThrottleReasons = map[string]bool{
	"vrhot0": true, "vrhot1": true, "prochot_cpu": true, "prochot_gfx": true,
}

// amdThrottleReasons returns the names of all active bits of indep_throttle_status and of the bits which indicate a
// problem
func amdThrottleReasons(status uint64) (active []string, problems []string) {
	active, problems = []string{}, []string{}
	for bit := 0; bit < 64; bit++ {
		name, ok := amdThrottleBits[bit]
		if !ok || status&(1<<bit) == 0 {
			continue
		}
		active = append(active, name)
		if amdProblemThrottleReasons[name] {
			problems = append(problems, name)
		}
	}
	return active, problems
}

// drmGpuDevices returns the device directories of the GPUs which use one of the drivers, e.g.
// /sys/class/drm/card1/device. The card directory is the parent directory.
func drmGpuDevices(drivers ...string) []string {
	cards, err := os.ReadDir(gpuDrmPath)
	if err != nil {
		return nil
	}
	var devices []string
	for _, card := range cards {
		// card1-DP-1 etc. are connectors, renderD128 the render node of the same GPU
		if !drmCardPattern.MatchString(card.Name()) {
			continue
		}
		device := filepath.Join(gpuDrmPath, card.Name(), "device")
		driver, err := filepath.EvalSymlinks(filepath.Join(device, "driver"))
		if err != nil || !slices.Contains(drivers, filepath.Base(driver)) {
			continue
		}
		devices = append(devices, device)
	}
	sort.Strings(devices)
	return devices
}

// amdGpuDevices returns the device directories of the GPUs which use the amdgpu driver
func amdGpuDevices() []string {
	return drmGpuDevices("amdgpu")
}

func hasAmdGpus() bool {
	return len(amdGpuDevices()) > 0
}

func collectAmdGpus() []*resultGpu {
	gpus := []*resultGpu{}
	for _, device := range amdGpuDevices() {
		gpus = append(gpus, readAmdGpu(device))
	}
	return gpus
}

func readAmdGpu(device string) *resultGpu {
	gpu := &resultGpu{
		Vendor:            "AMD",
		Source:            "amdgpu",
		Serial:            readSysfsString(filepath.Join(device, "serial_number")),
		UUID:              readSysfsString(filepath.Join(device, "unique_id")),
		VbiosVersion:      readSysfsString(filepath.Join(device, "vbios_version")),
		Processes:         -1,
		ThrottleReasons:   []string{},
		ClockEventReasons: []string{},
		XidErrors:         []*resultGpuXid{},
		Messages:          []string{},
	}
	if path, err := filepath.EvalSymlinks(device); err == nil {
		gpu.PciBusID = filepath.Base(path)
	}
	gpu.Name = readSysfsString(filepath.Join(device, "product_name"))
	if gpu.Name == "" {
		deviceID := readSysfsString(filepath.Join(device, "device"))
		gpu.Name = amdGpuMarketingName(deviceID, readSysfsString(filepath.Join(device, "revision")))
		if gpu.Name == "" {
			gpu.Name = amdGpuName(deviceID)
		}
	}

	if metrics, err := os.ReadFile(filepath.Join(device, "gpu_metrics")); err == nil {
		integrated, throttleStatus, hasThrottleStatus := parseAmdGpuMetrics(metrics)
		gpu.Integrated = integrated
		if hasThrottleStatus {
			gpu.ClockEventReasons, gpu.ThrottleReasons = amdThrottleReasons(throttleStatus)
		}
	}

	hwmon := amdGpuHwmon(device)
	if !gpu.Integrated && isAmdApuHwmon(hwmon) {
		// Older APUs (e.g. Raven and Picasso, Ryzen 2000G and 3000G) do not support gpu_metrics
		gpu.Integrated = true
	}

	gpu.UtilizationPercent = readSysfsNumber(filepath.Join(device, "gpu_busy_percent"), 1)
	gpu.MemoryUtilizationPercent = readSysfsNumber(filepath.Join(device, "mem_busy_percent"), 1)

	vramTotal := readSysfsUint(filepath.Join(device, "mem_info_vram_total"))
	vramUsed := readSysfsUint(filepath.Join(device, "mem_info_vram_used"))
	if gpu.Integrated {
		// The VRAM of an APU is a small part of the system memory, the GPU mostly uses the system memory (GTT)
		gttUsed := readSysfsUint(filepath.Join(device, "mem_info_gtt_used"))
		if vramUsed != nil || gttUsed != nil {
			var used uint64
			if vramUsed != nil {
				used += *vramUsed
			}
			if gttUsed != nil {
				used += *gttUsed
			}
			gpu.MemoryUsedBytes = &used
		}
	} else {
		gpu.MemoryTotalBytes, gpu.MemoryUsedBytes = vramTotal, vramUsed
	}

	gpu.PcieLinkGen = pcieGeneration(readSysfsString(filepath.Join(device, "current_link_speed")))
	gpu.PcieLinkGenMax = pcieGeneration(readSysfsString(filepath.Join(device, "max_link_speed")))
	gpu.PcieLinkWidth = readSysfsNumber(filepath.Join(device, "current_link_width"), 1)
	gpu.PcieLinkWidthMax = readSysfsNumber(filepath.Join(device, "max_link_width"), 1)

	readAmdRas(device, gpu)

	if hwmon != "" {
		readAmdHwmon(hwmon, gpu)
	}
	return gpu
}

// readAmdHwmon reads the sensors of the amdgpu hwmon device (millidegrees, microwatts, Hz), see
// https://docs.kernel.org/gpu/amdgpu/thermal.html#hwmon-interfaces
func readAmdHwmon(hwmon string, gpu *resultGpu) {
	labels := map[string]string{}
	files, _ := filepath.Glob(filepath.Join(hwmon, "temp*_label"))
	for _, file := range files {
		labels[readSysfsString(file)] = strings.TrimSuffix(filepath.Base(file), "_label")
	}
	temperature := func(label, suffix string) *float64 {
		sensor, ok := labels[label]
		if !ok {
			return nil
		}
		return readSysfsNumber(filepath.Join(hwmon, sensor+"_"+suffix), 1000)
	}

	gpu.TemperatureCelsius = temperature("edge", "input")
	gpu.TemperatureHotspotCelsius = temperature("junction", "input")
	if gpu.TemperatureCelsius == nil {
		// Instinct MI300A only has the junction sensor, old GPUs a temp1 without label
		gpu.TemperatureCelsius = gpu.TemperatureHotspotCelsius
		if gpu.TemperatureCelsius == nil {
			gpu.TemperatureCelsius = readSysfsNumber(filepath.Join(hwmon, "temp1_input"), 1000)
		}
	}
	gpu.TemperatureMemoryCelsius = temperature("mem", "input")

	// The GPU slows down at the critical temperature of the junction (hotspot) and shuts down at the emergency
	// temperature ("temperature emergency max value (asic shutdown)"). Older GPUs only have the edge sensor, junction
	// and mem exist on SOC15 dGPUs (Vega and newer).
	current := gpu.TemperatureHotspotCelsius
	gpu.TemperatureSlowdownCelsius = temperature("junction", "crit")
	gpu.TemperatureShutdownCelsius = temperature("junction", "emergency")
	if current == nil || gpu.TemperatureSlowdownCelsius == nil {
		current = gpu.TemperatureCelsius
		gpu.TemperatureSlowdownCelsius = temperature("edge", "crit")
		gpu.TemperatureShutdownCelsius = temperature("edge", "emergency")
	}
	if current != nil && gpu.TemperatureSlowdownCelsius != nil {
		headroom := *gpu.TemperatureSlowdownCelsius - *current
		gpu.TemperatureHeadroomCelsius = &headroom
	}

	// Fans can stop at low temperatures (zero RPM), 0 % is no error. Without pwm1 the speed is calculated from the
	// RPM and the maximum RPM.
	gpu.FanSpeedRpm = readSysfsNumber(filepath.Join(hwmon, "fan1_input"), 1)
	if pwm := readSysfsNumber(filepath.Join(hwmon, "pwm1"), 1); pwm != nil {
		pwmMax := 255.0
		if value := readSysfsNumber(filepath.Join(hwmon, "pwm1_max"), 1); value != nil && *value > 0 {
			pwmMax = *value
		}
		percent := *pwm / pwmMax * 100
		gpu.FanSpeedPercent = &percent
	} else if rpmMax := readSysfsNumber(filepath.Join(hwmon, "fan1_max"), 1); gpu.FanSpeedRpm != nil && rpmMax != nil && *rpmMax > 0 {
		percent := min(*gpu.FanSpeedRpm / *rpmMax * 100, 100)
		gpu.FanSpeedPercent = &percent
	}

	// The power (PPT) of an APU is the power of the whole processor: "power1_average: average power used by the SoC
	// in microWatts. On APUs this includes the CPU."
	if !gpu.Integrated {
		gpu.PowerWatts = readSysfsNumber(filepath.Join(hwmon, "power1_average"), 1000000)
		if gpu.PowerWatts == nil {
			gpu.PowerWatts = readSysfsNumber(filepath.Join(hwmon, "power1_input"), 1000000)
		}
		gpu.PowerLimitWatts = readSysfsNumber(filepath.Join(hwmon, "power1_cap"), 1000000)
	}

	freqLabels := map[string]string{}
	files, _ = filepath.Glob(filepath.Join(hwmon, "freq*_label"))
	for _, file := range files {
		freqLabels[readSysfsString(file)] = strings.TrimSuffix(filepath.Base(file), "_label")
	}
	if sensor, ok := freqLabels["sclk"]; ok {
		gpu.ClockGraphicsMhz = readSysfsNumber(filepath.Join(hwmon, sensor+"_input"), 1000000)
	}
	if sensor, ok := freqLabels["mclk"]; ok {
		gpu.ClockMemoryMhz = readSysfsNumber(filepath.Join(hwmon, sensor+"_input"), 1000000)
	}
}

// readAmdRas reads the error counters of RAS (data center GPUs like the Instinct series): "ue: 0\nce: 0\n", see
// https://docs.kernel.org/gpu/amdgpu/ras.html
func readAmdRas(device string, gpu *resultGpu) {
	files, _ := filepath.Glob(filepath.Join(device, "ras", "*_err_count"))
	if len(files) == 0 {
		return
	}
	var corrected, uncorrected uint64
	found := false
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			key, value, ok := strings.Cut(line, ":")
			if !ok {
				continue
			}
			count, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
			if err != nil {
				continue
			}
			switch strings.TrimSpace(key) {
			case "ue":
				uncorrected += count
				found = true
			case "ce":
				corrected += count
				found = true
			}
		}
	}
	if found {
		enabled := true
		gpu.EccEnabled = &enabled
		gpu.EccCorrectedVolatile = &corrected
		gpu.EccUncorrectedVolatile = &uncorrected
	}
}

// parseAmdGpuMetrics reads the header of the binary gpu_metrics of the power management firmware. Format revision 1
// are dedicated GPUs (struct gpu_metrics_v1_x), 2 APUs (struct gpu_metrics_v2_x). The throttle status is only read
// from the layout of version 1.3 (RDNA 2 and 3). The structures are defined in
// https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/drivers/gpu/drm/amd/include/kgd_pp_interface.h
func parseAmdGpuMetrics(data []byte) (integrated bool, throttleStatus uint64, hasThrottleStatus bool) {
	if len(data) < 4 {
		return false, 0, false
	}
	size := int(binary.LittleEndian.Uint16(data[0:2]))
	formatRevision, contentRevision := data[2], data[3]
	if formatRevision == 2 {
		return true, 0, false
	}
	// struct gpu_metrics_v1_3 (120 bytes): indep_throttle_status is the last field (offset 112)
	if formatRevision == 1 && contentRevision == 3 && size == 120 && len(data) >= 120 {
		return false, binary.LittleEndian.Uint64(data[112:120]), true
	}
	return false, 0, false
}

// isAmdApuHwmon returns true if the hwmon device has the voltage of the northbridge (vddnb). The kernel only creates it
// for APUs: "only APUs other than gc 9,4,3 have vddnb" in drivers/gpu/drm/amd/pm/amdgpu_pm.c
// (https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/drivers/gpu/drm/amd/pm/amdgpu_pm.c)
func isAmdApuHwmon(hwmon string) bool {
	if hwmon == "" {
		return false
	}
	return readSysfsString(filepath.Join(hwmon, "in1_label")) == "vddnb"
}

func amdGpuHwmon(device string) string {
	hwmons, _ := filepath.Glob(filepath.Join(device, "hwmon", "hwmon*"))
	for _, hwmon := range hwmons {
		if readSysfsString(filepath.Join(hwmon, "name")) == "amdgpu" {
			return hwmon
		}
	}
	return ""
}

// amdGpuName returns the name of an AMD GPU, e.g. "AMD Navi 32 [Radeon RX 7700 XT / 7800 XT]"
func amdGpuName(deviceID string) string {
	return gpuPciName("1002", "AMD", deviceID)
}

// amdGpuMarketingName returns the product name of libdrm (amdgpu.ids) by device and revision id, "" if unknown
func amdGpuMarketingName(deviceID, revision string) string {
	deviceID = strings.ToUpper(strings.TrimPrefix(strings.ToLower(deviceID), "0x"))
	revision = strings.ToUpper(strings.TrimPrefix(strings.ToLower(revision), "0x"))
	if deviceID == "" || revision == "" {
		return ""
	}
	key := "amdgpu.ids:" + deviceID + ":" + revision

	pciIdsMutex.Lock()
	defer pciIdsMutex.Unlock()
	if name, ok := pciIdsCache[key]; ok {
		return name
	}
	name := lookupAmdgpuIds(deviceID, revision)
	pciIdsCache[key] = name
	return name
}

// lookupAmdgpuIds searches the device in amdgpu.ids of libdrm: "device_id,\trevision_id,\tproduct_name"
func lookupAmdgpuIds(deviceID, revision string) string {
	for _, path := range gpuAmdgpuIdsPaths {
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fields := strings.SplitN(scanner.Text(), ",", 3)
			if len(fields) != 3 || strings.HasPrefix(fields[0], "#") {
				continue
			}
			if strings.EqualFold(strings.TrimSpace(fields[0]), deviceID) && strings.EqualFold(strings.TrimSpace(fields[1]), revision) {
				return strings.TrimSpace(fields[2])
			}
		}
		return ""
	}
	return ""
}

// gpuPciName returns the name of the GPU of the PCI id database (https://pci-ids.ucw.cz/) with the vendor name as
// prefix, e.g. "AMD Navi 32 [Radeon RX 7700 XT / 7800 XT]"
func gpuPciName(vendorID, vendorName, deviceID string) string {
	deviceID = strings.ToLower(strings.TrimPrefix(deviceID, "0x"))
	if deviceID == "" {
		return vendorName + " GPU"
	}
	key := vendorID + ":" + deviceID

	pciIdsMutex.Lock()
	defer pciIdsMutex.Unlock()
	if name, ok := pciIdsCache[key]; ok {
		return name
	}
	name := lookupPciIds(vendorID, deviceID)
	if name == "" {
		name = fmt.Sprintf("%s GPU [%s]", vendorName, key)
	} else {
		name = vendorName + " " + name
	}
	pciIdsCache[key] = name
	return name
}

// lookupPciIds searches the device in the PCI id database (pci.ids of pciutils or hwdata)
func lookupPciIds(vendorID, deviceID string) string {
	for _, path := range gpuPciIdsPaths {
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		defer file.Close()

		inVendor := false
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" || line[0] == '#' {
				continue
			}
			if line[0] != '\t' {
				if inVendor {
					return ""
				}
				inVendor = strings.HasPrefix(line, vendorID+"  ")
				continue
			}
			// Devices start with one tab, subsystems with two
			if inVendor && !strings.HasPrefix(line, "\t\t") && strings.HasPrefix(line, "\t"+deviceID+"  ") {
				return strings.TrimSpace(line[len(deviceID)+3:])
			}
		}
		return ""
	}
	return ""
}

// pcieGeneration converts the link speed to the PCIe generation, e.g. "16.0 GT/s PCIe" to 4
func pcieGeneration(speed string) *float64 {
	fields := strings.Fields(speed)
	if len(fields) == 0 {
		return nil
	}
	transfers, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return nil
	}
	generations := []struct {
		transfers  float64
		generation float64
	}{{2.5, 1}, {5, 2}, {8, 3}, {16, 4}, {32, 5}, {64, 6}, {128, 7}}
	for _, g := range generations {
		if transfers <= g.transfers {
			generation := g.generation
			return &generation
		}
	}
	return nil
}

func readSysfsString(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// readSysfsNumber returns the value of the file divided by the divisor, e.g. millidegrees to degrees
func readSysfsNumber(path string, divisor float64) *float64 {
	value, err := strconv.ParseFloat(readSysfsString(path), 64)
	if err != nil {
		return nil
	}
	value /= divisor
	return &value
}

func readSysfsUint(path string) *uint64 {
	value, err := strconv.ParseUint(readSysfsString(path), 10, 64)
	if err != nil {
		return nil
	}
	return &value
}
