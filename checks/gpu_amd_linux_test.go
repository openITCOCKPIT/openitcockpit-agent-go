package checks

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var (
	amdDumpLspciPattern = regexp.MustCompile(`^([0-9a-f]{2}:[0-9a-f]{2}\.[0-9a-f]) .*\[1002:([0-9a-f]{4})\]`)
	amdDumpCardPattern  = regexp.MustCompile(`^== /sys/class/drm/(card\d+)/device driver=(\S+)`)
	amdDumpHwmonPattern = regexp.MustCompile(`^-- /sys/class/drm/card\d+/device/hwmon/(hwmon\d+) name=(\S+)`)
)

// buildAmdSysfs creates a fake /sys/class/drm with the files of a dump of amd-sysfs-dump.sh. The device directories
// are symlinks to the PCI devices, like in the real sysfs.
func buildAmdSysfs(t *testing.T, dump string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "gpu", "amd", dump))
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	write := func(path, content string) {
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	symlink := func(target, link string) {
		// The targets exist in the real sysfs
		if err := os.MkdirAll(target, 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Dir(link), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, link); err != nil {
			t.Fatal(err)
		}
	}

	type card struct {
		name, driver string
		lines        []string
	}
	pciBusIDs := map[string]string{}
	var cards []*card
	for _, line := range strings.Split(string(data), "\n") {
		if match := amdDumpLspciPattern.FindStringSubmatch(line); match != nil {
			pciBusIDs["0x"+match[2]] = "0000:" + match[1]
			continue
		}
		if match := amdDumpCardPattern.FindStringSubmatch(line); match != nil {
			cards = append(cards, &card{name: match[1], driver: match[2]})
			continue
		}
		if len(cards) > 0 {
			cards[len(cards)-1].lines = append(cards[len(cards)-1].lines, line)
		}
	}

	for _, c := range cards {
		var deviceID string
		for _, line := range c.lines {
			if strings.HasPrefix(line, "device: ") {
				deviceID = strings.TrimSpace(strings.TrimPrefix(line, "device: "))
			}
		}
		device := filepath.Join(root, "devices", pciBusIDs[deviceID])
		symlink(device, filepath.Join(root, "drm", c.name, "device"))
		symlink(filepath.Join(root, "drivers", c.driver), filepath.Join(device, "driver"))

		dir := device
		for _, line := range c.lines {
			if match := amdDumpHwmonPattern.FindStringSubmatch(line); match != nil {
				dir = filepath.Join(device, "hwmon", match[1])
				write(filepath.Join(dir, "name"), match[2]+"\n")
				continue
			}
			file, value, found := strings.Cut(line, ": ")
			if !found || value == "(present, not dumped)" {
				continue
			}
			if file == "gpu_metrics (base64)" {
				metrics, err := base64.StdEncoding.DecodeString(value)
				if err != nil {
					t.Fatal(err)
				}
				write(filepath.Join(dir, "gpu_metrics"), string(metrics))
				continue
			}
			write(filepath.Join(dir, file), strings.TrimSpace(value)+"\n")
		}
	}

	// Connectors, render nodes and GPUs of other drivers are ignored
	if err := os.MkdirAll(filepath.Join(root, "drm", "card1-DP-1"), 0700); err != nil {
		t.Fatal(err)
	}
	symlink(filepath.Join(root, "devices", pciBusIDs["0x747e"]), filepath.Join(root, "drm", "renderD128", "device"))
	// BMC graphics of servers (ASPEED)
	other := filepath.Join(root, "devices", "0000:02:00.0")
	write(filepath.Join(other, "vendor"), "0x1a03\n")
	symlink(other, filepath.Join(root, "drm", "card0", "device"))
	symlink(filepath.Join(root, "drivers", "ast"), filepath.Join(other, "driver"))

	return filepath.Join(root, "drm")
}

func setAmdSysfs(t *testing.T, dump string) {
	t.Helper()
	drmPath, pciIdsPaths, procPath, amdgpuIdsPaths := gpuDrmPath, gpuPciIdsPaths, gpuProcPath, gpuAmdgpuIdsPaths
	gpuProcPath = t.TempDir()
	drmClientsMutex.Lock()
	drmClientsPrevious = nil
	drmClientsMutex.Unlock()
	t.Cleanup(func() {
		gpuDrmPath, gpuPciIdsPaths, gpuProcPath, gpuAmdgpuIdsPaths = drmPath, pciIdsPaths, procPath, amdgpuIdsPaths
		pciIdsMutex.Lock()
		pciIdsCache = map[string]string{}
		pciIdsMutex.Unlock()
	})
	gpuDrmPath = buildAmdSysfs(t, dump)

	// The PCI id database only knows the RX 7800 XT, not the iGPU
	pciIds := filepath.Join(t.TempDir(), "pci.ids")
	content := "# pci.ids\n1002  Advanced Micro Devices, Inc. [AMD/ATI]\n\t7480  Navi 33 [Radeon RX 7600/7600 XT/7600M XT/7600S/7700S / PRO W7600]\n\t747e  Navi 32 [Radeon RX 7700 XT / 7800 XT]\n\t\t1002 0e3b  Radeon RX 7800 XT\n10de  NVIDIA Corporation\n\t164e  Not an AMD GPU\n"
	if err := os.WriteFile(pciIds, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	gpuPciIdsPaths = []string{filepath.Join(t.TempDir(), "missing"), pciIds}

	// The product names of libdrm only know the revision C8 of the device 747E
	amdgpuIds := filepath.Join(t.TempDir(), "amdgpu.ids")
	if err := os.WriteFile(amdgpuIds, []byte("# List of AMDGPU IDs\n\n1.0.0\n747E,\tC8,\tAMD Radeon RX 7800 XT\n747E,\tFF,\tAMD Radeon RX 7700 XT\n15D8,\tC8,\tAMD Radeon Vega 11 Graphics\n"), 0600); err != nil {
		t.Fatal(err)
	}
	gpuAmdgpuIdsPaths = []string{amdgpuIds}
	pciIdsMutex.Lock()
	pciIdsCache = map[string]string{}
	pciIdsMutex.Unlock()
}

// Radeon RX 7800 XT and Raphael iGPU of a Ryzen 7000, idle
func TestCollectAmdGpusIdle(t *testing.T) {
	setAmdSysfs(t, "rx7800xt-raphael-idle.txt")
	if !hasPlatformGpus() {
		t.Fatal("expected AMD GPUs")
	}
	gpus := collectAmdGpus()
	if len(gpus) != 2 {
		t.Fatalf("expected 2 GPUs, got %d", len(gpus))
	}

	gpu := gpus[0]
	if gpu.Vendor != "AMD" || gpu.Source != "amdgpu" || gpu.Name != "AMD Radeon RX 7800 XT" ||
		gpu.PciBusID != "0000:03:00.0" || gpu.Integrated || gpu.VbiosVersion != "113-EXT90455-100" || gpu.Processes != -1 {
		t.Errorf("unexpected GPU: %+v", gpu)
	}
	assertFloat(t, "utilization_percent", gpu.UtilizationPercent, 0)
	assertUint(t, "memory_total_bytes", gpu.MemoryTotalBytes, 17163091968)
	assertUint(t, "memory_used_bytes", gpu.MemoryUsedBytes, 927903744)
	assertFloat(t, "temperature_celsius", gpu.TemperatureCelsius, 36)
	assertFloat(t, "temperature_hotspot_celsius", gpu.TemperatureHotspotCelsius, 42)
	assertFloat(t, "temperature_memory_celsius", gpu.TemperatureMemoryCelsius, 62)
	assertFloat(t, "temperature_slowdown_celsius", gpu.TemperatureSlowdownCelsius, 110)
	assertFloat(t, "temperature_shutdown_celsius", gpu.TemperatureShutdownCelsius, 115)
	assertFloat(t, "temperature_headroom_celsius", gpu.TemperatureHeadroomCelsius, 68)
	// Zero RPM: the fan stops when idle
	assertFloat(t, "fan_speed_percent", gpu.FanSpeedPercent, 0)
	assertFloat(t, "fan_speed_rpm", gpu.FanSpeedRpm, 0)
	assertFloat(t, "power_watts", gpu.PowerWatts, 26)
	assertFloat(t, "power_limit_watts", gpu.PowerLimitWatts, 228)
	assertFloat(t, "clock_graphics_mhz", gpu.ClockGraphicsMhz, 36)
	assertFloat(t, "clock_memory_mhz", gpu.ClockMemoryMhz, 456)
	assertFloat(t, "pcie_link_gen", gpu.PcieLinkGen, 4)
	assertFloat(t, "pcie_link_width", gpu.PcieLinkWidth, 16)
	assertNil(t, "ecc_enabled", gpu.EccEnabled)
	if len(gpu.ClockEventReasons) != 0 || len(gpu.ThrottleReasons) != 0 {
		t.Errorf("unexpected clock event reasons: %v %v", gpu.ClockEventReasons, gpu.ThrottleReasons)
	}

	// The iGPU is neither part of amdgpu.ids nor of the PCI id database, it is detected as APU by the format of
	// gpu_metrics
	gpu = gpus[1]
	if gpu.Name != "AMD GPU [1002:164e]" || gpu.PciBusID != "0000:11:00.0" || !gpu.Integrated {
		t.Errorf("unexpected GPU: %+v", gpu)
	}
	// VRAM (carve-out of the system memory) and GTT
	assertUint(t, "memory_used_bytes", gpu.MemoryUsedBytes, 16687104+14143488)
	assertNil(t, "memory_total_bytes", gpu.MemoryTotalBytes)
	assertFloat(t, "temperature_celsius", gpu.TemperatureCelsius, 40)
	assertNil(t, "temperature_headroom_celsius", gpu.TemperatureHeadroomCelsius)
	// The power of an APU is the power of the whole processor
	assertNil(t, "power_watts", gpu.PowerWatts)
	assertNil(t, "fan_speed_percent", gpu.FanSpeedPercent)
	assertNil(t, "memory_utilization_percent", gpu.MemoryUtilizationPercent)
	assertNil(t, "clock_memory_mhz", gpu.ClockMemoryMhz)
	assertFloat(t, "clock_graphics_mhz", gpu.ClockGraphicsMhz, 600)
}

// Load on the RX 7800 XT: the firmware limits the clock by the hotspot temperature (at 50 °C), which is no problem
func TestCollectAmdGpusLoad(t *testing.T) {
	setAmdSysfs(t, "rx7800xt-raphael-load.txt")
	gpus := collectAmdGpus()
	if len(gpus) != 2 {
		t.Fatalf("expected 2 GPUs, got %d", len(gpus))
	}
	gpu := gpus[0]
	assertFloat(t, "utilization_percent", gpu.UtilizationPercent, 87)
	assertFloat(t, "memory_utilization_percent", gpu.MemoryUtilizationPercent, 4)
	assertFloat(t, "power_watts", gpu.PowerWatts, 75)
	assertFloat(t, "clock_graphics_mhz", gpu.ClockGraphicsMhz, 2510)
	assertFloat(t, "temperature_headroom_celsius", gpu.TemperatureHeadroomCelsius, 60)
	if !reflect.DeepEqual(gpu.ClockEventReasons, []string{"temp_hotspot"}) || len(gpu.ThrottleReasons) != 0 {
		t.Errorf("unexpected clock event reasons: %v %v", gpu.ClockEventReasons, gpu.ThrottleReasons)
	}
}

// Load on the iGPU
func TestCollectAmdGpusIgpuLoad(t *testing.T) {
	setAmdSysfs(t, "rx7800xt-raphael-igpu-load.txt")
	gpus := collectAmdGpus()
	if len(gpus) != 2 {
		t.Fatalf("expected 2 GPUs, got %d", len(gpus))
	}
	gpu := gpus[1]
	assertFloat(t, "utilization_percent", gpu.UtilizationPercent, 95)
	assertFloat(t, "clock_graphics_mhz", gpu.ClockGraphicsMhz, 2200)
	assertFloat(t, "temperature_celsius", gpu.TemperatureCelsius, 50)
}

func TestCheckGpuRunAmd(t *testing.T) {
	setAmdSysfs(t, "rx7800xt-raphael-idle.txt")
	defer func(path string) { gpuKmsgPath = path }(gpuKmsgPath)
	gpuKmsgPath = filepath.Join(t.TempDir(), "kmsg")

	r, err := (&CheckGpu{platform: true}).Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	result := r.(*resultGpus)
	if len(result.Gpus) != 2 || result.Nvidia != nil || result.Gpus[1].Index != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
	// The memory of an integrated GPU is the system memory
	if result.Gpus[1].MemoryTotalBytes == nil || *result.Gpus[1].MemoryTotalBytes == 0 {
		t.Errorf("expected the system memory: %v", deref(result.Gpus[1].MemoryTotalBytes))
	}
}

func TestReadAmdRas(t *testing.T) {
	device := t.TempDir()
	for name, content := range map[string]string{"umc_err_count": "ue: 1\nce: 12\n", "gfx_err_count": "ue: 0\nce: 3\nde: 0\n"} {
		if err := os.MkdirAll(filepath.Join(device, "ras"), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(device, "ras", name), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	gpu := &resultGpu{}
	readAmdRas(device, gpu)
	assertBool(t, "ecc_enabled", gpu.EccEnabled, true)
	assertUint(t, "ecc_corrected_volatile", gpu.EccCorrectedVolatile, 15)
	assertUint(t, "ecc_uncorrected_volatile", gpu.EccUncorrectedVolatile, 1)
}

func TestParseAmdGpuMetrics(t *testing.T) {
	metrics := make([]byte, 120)
	metrics[0], metrics[2], metrics[3] = 120, 1, 3
	metrics[112] = 1        // ppt0
	metrics[112+5] = 1 << 4 // bit 44: vrhot0
	integrated, status, ok := parseAmdGpuMetrics(metrics)
	if integrated || !ok || status != 1|1<<44 {
		t.Errorf("unexpected result: %v %x %v", integrated, status, ok)
	}
	// Other layouts are not read
	metrics[3] = 4
	if _, _, ok := parseAmdGpuMetrics(metrics); ok {
		t.Error("unexpected throttle status")
	}
	if integrated, _, ok := parseAmdGpuMetrics([]byte{120, 0, 2, 1}); !integrated || ok {
		t.Error("expected an APU")
	}
	if integrated, _, ok := parseAmdGpuMetrics(nil); integrated || ok {
		t.Error("unexpected result for empty metrics")
	}
}

func TestAmdThrottleReasons(t *testing.T) {
	// SMU_THROTTLER_*_BIT of amdgpu_smu.h: PPT0 0, TEMP_HOTSPOT 36, VRHOT0 44, PROCHOT_CPU 46, PROCHOT_GFX 47, FIT 57
	active, problems := amdThrottleReasons(1 | 1<<36 | 1<<44 | 1<<46 | 1<<47 | 1<<57 | 1<<63)
	if !reflect.DeepEqual(active, []string{"ppt0", "temp_hotspot", "vrhot0", "prochot_cpu", "prochot_gfx", "fit"}) ||
		!reflect.DeepEqual(problems, []string{"vrhot0", "prochot_cpu", "prochot_gfx"}) {
		t.Errorf("unexpected reasons: %v %v", active, problems)
	}
	if active, problems := amdThrottleReasons(0); len(active) != 0 || len(problems) != 0 || active == nil {
		t.Errorf("unexpected reasons: %v %v", active, problems)
	}
}

func TestPcieGeneration(t *testing.T) {
	tests := map[string]interface{}{
		"2.5 GT/s PCIe": 1.0, "5.0 GT/s PCIe": 2.0, "8.0 GT/s PCIe": 3.0, "16.0 GT/s PCIe": 4.0, "32.0 GT/s PCIe": 5.0,
		"64.0 GT/s PCIe": 6.0, "Unknown": nil, "": nil,
	}
	for speed, expected := range tests {
		if got := deref(pcieGeneration(speed)); got != expected {
			t.Errorf("%q: expected %v, got %v", speed, expected, got)
		}
	}
}

// A system with NVIDIA and AMD GPUs: nvidia-smi and the amdgpu driver report their own GPUs
func TestCheckGpuRunNvidiaAndAmd(t *testing.T) {
	setAmdSysfs(t, "rx7800xt-raphael-idle.txt")
	defer func(path string) { gpuKmsgPath = path }(gpuKmsgPath)
	gpuKmsgPath = filepath.Join(t.TempDir(), "kmsg")

	data, err := os.ReadFile(filepath.Join("testdata", "gpu", "nvidia", "h200-nvl-4x-linux-595.91.07.xml"))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "query.xml"), data, 0600); err != nil {
		t.Fatal(err)
	}
	nvidiaSmi := filepath.Join(dir, "nvidia-smi")
	if err := os.WriteFile(nvidiaSmi, []byte("#!/bin/sh\ncat "+filepath.Join(dir, "query.xml")+"\n"), 0700); err != nil {
		t.Fatal(err)
	}

	r, err := (&CheckGpu{nvidiaSmi: nvidiaSmi, platform: true}).Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	result := r.(*resultGpus)
	var vendors []string
	pciBusIDs := map[string]bool{}
	for i, gpu := range result.Gpus {
		vendors = append(vendors, gpu.Vendor+"/"+gpu.Source)
		if gpu.Index != i || pciBusIDs[gpu.PciBusID] {
			t.Errorf("unexpected index or duplicate PCI bus id: %+v", gpu)
		}
		pciBusIDs[gpu.PciBusID] = true
	}
	expected := []string{"NVIDIA/nvidia-smi", "NVIDIA/nvidia-smi", "NVIDIA/nvidia-smi", "NVIDIA/nvidia-smi", "AMD/amdgpu", "AMD/amdgpu"}
	if !reflect.DeepEqual(vendors, expected) || result.Nvidia == nil || result.Nvidia.Error != "" {
		t.Errorf("unexpected GPUs: %v %+v", vendors, result.Nvidia)
	}
}

// Instinct MI300A: only the junction sensor; GPU without pwm1: the fan speed is calculated from the RPM
func TestReadAmdHwmonFallbacks(t *testing.T) {
	hwmon := t.TempDir()
	for name, content := range map[string]string{
		"temp1_label": "junction", "temp1_input": "61000", "temp1_crit": "100000",
		"fan1_input": "1500", "fan1_max": "3000",
	} {
		if err := os.WriteFile(filepath.Join(hwmon, name), []byte(content+"\n"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	gpu := &resultGpu{}
	readAmdHwmon(hwmon, gpu)
	assertFloat(t, "temperature_celsius", gpu.TemperatureCelsius, 61)
	assertFloat(t, "temperature_hotspot_celsius", gpu.TemperatureHotspotCelsius, 61)
	assertFloat(t, "temperature_headroom_celsius", gpu.TemperatureHeadroomCelsius, 39)
	assertFloat(t, "fan_speed_rpm", gpu.FanSpeedRpm, 1500)
	assertFloat(t, "fan_speed_percent", gpu.FanSpeedPercent, 50)
}

func TestAmdGpuMarketingName(t *testing.T) {
	setAmdSysfs(t, "rx7800xt-raphael-idle.txt")
	tests := map[[2]string]string{
		{"0x747e", "0xc8"}: "AMD Radeon RX 7800 XT",
		{"0x747e", "0xff"}: "AMD Radeon RX 7700 XT",
		{"0x747e", "0x00"}: "",
		{"", "0xc8"}:       "",
	}
	for ids, expected := range tests {
		if name := amdGpuMarketingName(ids[0], ids[1]); name != expected {
			t.Errorf("%v: expected %q, got %q", ids, expected, name)
		}
	}
}

// Picasso APU (Ryzen 3000G, Radeon Vega 11) under load (glmark2): no gpu_metrics, the APU is detected by vddnb. Live
// test: the Ubuntu resources app showed 99 % utilization, 1.39 of 2.15 GB video memory, 1.40 GHz and 46 °C.
func TestCollectAmdGpusPicasso(t *testing.T) {
	setAmdSysfs(t, "picasso-vega11-load.txt")
	gpus := collectAmdGpus()
	if len(gpus) != 1 {
		t.Fatalf("expected 1 GPU, got %d", len(gpus))
	}
	gpu := gpus[0]
	if gpu.Name != "AMD Radeon Vega 11 Graphics" || gpu.PciBusID != "0000:27:00.0" || !gpu.Integrated || gpu.VbiosVersion != "113-PICASSO-118" {
		t.Errorf("unexpected GPU: %+v", gpu)
	}
	assertFloat(t, "utilization_percent", gpu.UtilizationPercent, 89)
	// VRAM (carve-out of the system memory) and GTT
	assertUint(t, "memory_used_bytes", gpu.MemoryUsedBytes, 1415872512+156766208)
	assertNil(t, "memory_total_bytes", gpu.MemoryTotalBytes)
	assertFloat(t, "temperature_celsius", gpu.TemperatureCelsius, 52)
	assertFloat(t, "clock_graphics_mhz", gpu.ClockGraphicsMhz, 1400)
	assertFloat(t, "pcie_link_gen", gpu.PcieLinkGen, 3)
	// No power (PPT without value), no fan
	assertNil(t, "power_watts", gpu.PowerWatts)
	assertNil(t, "fan_speed_percent", gpu.FanSpeedPercent)
	assertNil(t, "temperature_headroom_celsius", gpu.TemperatureHeadroomCelsius)
}
