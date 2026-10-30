package checks

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	intelDumpCardPattern  = regexp.MustCompile(`^== /sys/class/drm/(card\d+) driver=(\S+) pci=(\S+)`)
	intelDumpDirPattern   = regexp.MustCompile(`^-- (/sys/class/drm/card\d+(?:/device)?)(\S*)(?: name=(\S+))?$`)
	intelDumpDeviceFiles  = []string{"vendor", "device", "subsystem_vendor", "subsystem_device", "revision", "class", "boot_vga", "current_link_speed", "current_link_width", "max_link_speed", "max_link_width"}
	intelDumpClientRecord = regexp.MustCompile(`(?m)^(drm-[\w-]+):\s*(.*)$`)
	// Version 4: upstream PCI bridges of the GPU, e.g. "-- upstream 1 0000:05:01.0"
	intelDumpUpstreamPattern = regexp.MustCompile(`(?m)^-- upstream (\d+) (\S+)$`)
)

// buildIntelSysfs creates a fake /sys/class/drm and /proc with the files of a dump of intel-sysfs-dump.sh. Every DRM
// client of the dump gets its own process with the file descriptor 3. Sample 1 uses the values of the second samples
// of the dump: "-- after 2 seconds" (e.g. the energy counter) and the second sample of the clients (version 2 and 3).
func buildIntelSysfs(t *testing.T, dump string, sample int) (drmPath, procPath string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "gpu", "intel", dump))
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

	sysfs, clients, _ := strings.Cut(string(data), "== drm clients (fdinfo)")
	// Version 2 and 3 of the dump script contain a second sample of the clients
	clients, secondClients, found := strings.Cut(clients, "== drm clients after")
	if sample == 1 && found {
		clients = secondClients
	}
	// The PCI devices are nested like in the real sysfs: the upstream bridges are the parent directories of the GPU
	// (only one GPU per dump)
	upstreamDirs := map[string]string{} // level -> directory
	devicesDir := filepath.Join(root, "devices")
	upstream := intelDumpUpstreamPattern.FindAllStringSubmatch(sysfs, -1)
	for i := len(upstream) - 1; i >= 0; i-- {
		devicesDir = filepath.Join(devicesDir, upstream[i][2])
		upstreamDirs[upstream[i][1]] = devicesDir
	}

	var card, device, dir string
	cardPaths := map[string][2]string{} // /sys/class/drm/cardN -> fake card and device directory
	for _, line := range strings.Split(sysfs, "\n") {
		if match := intelDumpCardPattern.FindStringSubmatch(line); match != nil {
			card = filepath.Join(root, "drm", match[1])
			device = filepath.Join(devicesDir, match[3])
			symlink(device, filepath.Join(card, "device"))
			symlink(filepath.Join(root, "drivers", match[2]), filepath.Join(device, "driver"))
			cardPaths["/sys/class/drm/"+match[1]] = [2]string{card, device}
			dir = ""
			continue
		}
		if strings.HasPrefix(line, "-- after") {
			// Second sample with absolute paths
			card = ""
			continue
		}
		if match := intelDumpUpstreamPattern.FindStringSubmatch(line); match != nil {
			dir = upstreamDirs[match[1]]
			continue
		}
		if path, value, found := strings.Cut(line, ": "); found && strings.HasPrefix(path, "/sys/class/drm/card") {
			if sample == 1 {
				for prefix, paths := range cardPaths {
					switch {
					case strings.HasPrefix(path, prefix+"/device/"):
						write(filepath.Join(paths[1], strings.TrimPrefix(path, prefix+"/device/")), strings.TrimSpace(value)+"\n")
					case strings.HasPrefix(path, prefix+"/"):
						write(filepath.Join(paths[0], strings.TrimPrefix(path, prefix+"/")), strings.TrimSpace(value)+"\n")
					}
				}
			}
			continue
		}
		if match := intelDumpDirPattern.FindStringSubmatch(line); match != nil && card != "" {
			base := card
			if strings.HasSuffix(match[1], "/device") {
				base = device
			}
			dir = filepath.Join(base, match[2])
			if match[3] != "" {
				write(filepath.Join(dir, "name"), match[3]+"\n")
			}
			continue
		}
		file, value, found := strings.Cut(line, ": ")
		if !found || card == "" || strings.HasPrefix(file, "#") {
			continue
		}
		switch {
		case dir != "":
			write(filepath.Join(dir, file), strings.TrimSpace(value)+"\n")
		case slices.Contains(intelDumpDeviceFiles, file):
			write(filepath.Join(device, file), strings.TrimSpace(value)+"\n")
		default:
			write(filepath.Join(card, file), strings.TrimSpace(value)+"\n")
		}
	}

	proc := filepath.Join(root, "proc")
	if err := os.MkdirAll(proc, 0700); err != nil {
		t.Fatal(err)
	}
	for i, record := range strings.Split(clients, "\n--") {
		if !intelDumpClientRecord.MatchString(record) {
			continue
		}
		pid := filepath.Join(proc, strconv.Itoa(1000+i))
		if err := os.MkdirAll(filepath.Join(pid, "fd"), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("/dev/dri/renderD128", filepath.Join(pid, "fd", "3")); err != nil {
			t.Fatal(err)
		}
		write(filepath.Join(pid, "fdinfo", "3"), strings.TrimSpace(record)+"\n")
	}
	return filepath.Join(root, "drm"), proc
}

// Raptor Lake-P (Iris Xe Graphics, i915) of a notebook, Debian 13, kernel 6.12: idle and 49 seconds later
func TestCollectIntelGpusI915(t *testing.T) {
	setIntelTestState(t)

	start := time.Date(2026, 10, 28, 19, 33, 48, 0, time.Local)
	gpuNow = func() time.Time { return start }
	gpuDrmPath, gpuProcPath = buildIntelSysfs(t, "raptorlake-p-i915-idle.txt", 0)
	if !hasPlatformGpus() {
		t.Fatal("expected an Intel GPU")
	}
	gpus := collectPlatformGpus(nil, nil)
	if len(gpus) != 1 {
		t.Fatalf("expected 1 GPU, got %d", len(gpus))
	}
	gpu := gpus[0]
	if gpu.Vendor != "Intel" || gpu.Source != "i915" || gpu.Name != "Intel Raptor Lake-P [Iris Xe Graphics]" ||
		gpu.PciBusID != "0000:00:02.0" || !gpu.Integrated || gpu.Processes != 27 {
		t.Errorf("unexpected GPU: %+v", gpu)
	}
	// The utilization requires a second sample, the GPU is in RC6 (frequency 0)
	assertNil(t, "utilization_percent", gpu.UtilizationPercent)
	assertFloat(t, "clock_graphics_mhz", gpu.ClockGraphicsMhz, 0)
	assertFloat(t, "clock_graphics_max_mhz", gpu.ClockGraphicsMaxMhz, 1250)
	assertNil(t, "pcie_link_width", gpu.PcieLinkWidth)
	assertNil(t, "temperature_celsius", gpu.TemperatureCelsius)
	assertNil(t, "power_watts", gpu.PowerWatts)
	if len(gpu.ClockEventReasons) != 0 || len(gpu.ThrottleReasons) != 0 {
		t.Errorf("unexpected throttle reasons: %v %v", gpu.ClockEventReasons, gpu.ThrottleReasons)
	}

	// 49 seconds later: the render engine was busy for 21273252 ns of all clients
	gpuNow = func() time.Time { return start.Add(49 * time.Second) }
	gpuDrmPath, gpuProcPath = buildIntelSysfs(t, "raptorlake-p-i915-load.txt", 0)
	gpus = collectPlatformGpus(nil, nil)
	if len(gpus) != 1 {
		t.Fatalf("expected 1 GPU, got %d", len(gpus))
	}
	gpu = gpus[0]
	if gpu.UtilizationPercent == nil || *gpu.UtilizationPercent < 0.0434 || *gpu.UtilizationPercent > 0.0435 {
		t.Errorf("unexpected utilization: %v", deref(gpu.UtilizationPercent))
	}
	assertFloat(t, "decoder_utilization_percent", gpu.DecoderUtilizationPercent, 0)
	assertFloat(t, "clock_graphics_mhz", gpu.ClockGraphicsMhz, 300)
}

func TestDrmUsageByPdev(t *testing.T) {
	start := time.Unix(1000, 0)
	client := func(pdev string, pid int, engines, cycles, totalCycles, capacity map[string]uint64) *drmClient {
		return &drmClient{pdev: pdev, pid: pid, engines: engines, cycles: cycles, totalCycles: totalCycles, capacity: capacity}
	}
	none := map[string]uint64{}
	previous := &drmClientsSample{time: start, clients: map[string]*drmClient{
		// i915: busy time in ns, two video engines
		"0000:00:02.0/1": client("0000:00:02.0", 10, map[string]uint64{"render": 0, "video": 0}, none, none, map[string]uint64{"video": 2}),
		// xe: busy cycles and the GPU timestamp in cycles
		"0000:03:00.0/1": client("0000:03:00.0", 20, none, map[string]uint64{"rcs": 1000}, map[string]uint64{"rcs": 10000}, none),
		"0000:03:00.0/2": client("0000:03:00.0", 21, none, map[string]uint64{"rcs": 0}, map[string]uint64{"rcs": 10000}, none),
		// amdgpu
		"0000:0a:00.0/1": client("0000:0a:00.0", 30, map[string]uint64{"gfx": 0, "enc": 0, "dec": 0}, none, none, none),
	}}
	current := &drmClientsSample{time: start.Add(2 * time.Second), clients: map[string]*drmClient{
		"0000:00:02.0/1": client("0000:00:02.0", 10, map[string]uint64{"render": 1e9, "video": 2e9}, none, none, map[string]uint64{"video": 2}),
		"0000:00:02.0/2": client("0000:00:02.0", 11, map[string]uint64{"render": 5e9}, none, none, none), // new client, not counted
		"0000:03:00.0/1": client("0000:03:00.0", 20, none, map[string]uint64{"rcs": 1500}, map[string]uint64{"rcs": 11000}, none),
		"0000:03:00.0/2": client("0000:03:00.0", 21, none, map[string]uint64{"rcs": 250}, map[string]uint64{"rcs": 11000}, none),
		"0000:0a:00.0/1": client("0000:0a:00.0", 30, map[string]uint64{"gfx": 1e9, "enc": 4e8, "dec": 0}, none, none, none),
	}}

	usage := drmUsageByPdev(previous, current)
	igpu := usage["0000:00:02.0"]
	if igpu == nil || igpu.processes != 2 {
		t.Fatalf("unexpected usage: %+v", igpu)
	}
	assertFloat(t, "i915 utilization", igpu.utilization, 50)
	assertFloat(t, "i915 encoder", igpu.encoder, 50)
	assertFloat(t, "i915 decoder", igpu.decoder, 50)

	xe := usage["0000:03:00.0"]
	assertFloat(t, "xe utilization", xe.utilization, 75)
	assertNil(t, "xe encoder", xe.encoder)

	amd := usage["0000:0a:00.0"]
	assertFloat(t, "amdgpu utilization", amd.utilization, 50)
	assertFloat(t, "amdgpu encoder", amd.encoder, 20)
	assertFloat(t, "amdgpu decoder", amd.decoder, 0)

	// Without a previous sample only the processes are known
	usage = drmUsageByPdev(nil, current)
	if usage["0000:00:02.0"].processes != 2 || usage["0000:00:02.0"].utilization != nil {
		t.Errorf("unexpected usage: %+v", usage["0000:00:02.0"])
	}
}

// Dedicated GPU with the xe driver (constructed from the ABI documentation, no real dump yet)
func TestReadIntelGpuXe(t *testing.T) {
	root := t.TempDir()
	device := filepath.Join(root, "devices", "0000:03:00.0")
	files := map[string]string{
		"device":                                  "0xe20b",
		"current_link_speed":                      "16.0 GT/s PCIe",
		"max_link_speed":                          "16.0 GT/s PCIe",
		"current_link_width":                      "8",
		"max_link_width":                          "8",
		"tile0/gt0/freq0/act_freq":                "2400",
		"tile0/gt0/freq0/rp0_freq":                "2850",
		"tile0/gt0/freq0/throttle/status":         "1",
		"tile0/gt0/freq0/throttle/reason_pl1":     "1",
		"tile0/gt0/freq0/throttle/reason_pl2":     "0",
		"tile0/gt0/freq0/throttle/reason_thermal": "1",
		"hwmon/hwmon5/name":                       "xe",
		"hwmon/hwmon5/temp2_input":                "65000",
		"hwmon/hwmon5/temp3_input":                "70000",
		"hwmon/hwmon5/temp2_crit":                 "95000",
		"hwmon/hwmon5/temp2_emergency":            "105000",
		"hwmon/hwmon5/power1_max":                 "190000000",
		"hwmon/hwmon5/energy1_input":              "1000000",
		"hwmon/hwmon5/fan1_input":                 "1200",
	}
	for name, content := range files {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(device, name)), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(device, name), []byte(content+"\n"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, "drivers", "xe"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "drivers", "xe"), filepath.Join(device, "driver")); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "drm", "card1"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(device, filepath.Join(root, "drm", "card1", "device")); err != nil {
		t.Fatal(err)
	}

	now := gpuNow
	start := time.Unix(1000, 0)
	t.Cleanup(func() {
		gpuNow = now
		intelEnergyMutex.Lock()
		delete(intelEnergyPrevious, "0000:03:00.0")
		intelEnergyMutex.Unlock()
	})
	gpuNow = func() time.Time { return start }
	usage := map[string]*drmUsage{"0000:03:00.0": {utilization: ptrFloat(42), processes: 3}}
	gpu := readIntelGpu(filepath.Join(root, "drm", "card1", "device"), usage, true)
	if gpu.Source != "xe" || gpu.Integrated || gpu.Processes != 3 {
		t.Errorf("unexpected GPU: %+v", gpu)
	}
	assertFloat(t, "utilization_percent", gpu.UtilizationPercent, 42)
	assertFloat(t, "clock_graphics_mhz", gpu.ClockGraphicsMhz, 2400)
	assertFloat(t, "clock_graphics_max_mhz", gpu.ClockGraphicsMaxMhz, 2850)
	assertFloat(t, "temperature_celsius", gpu.TemperatureCelsius, 65)
	assertFloat(t, "temperature_memory_celsius", gpu.TemperatureMemoryCelsius, 70)
	assertFloat(t, "temperature_headroom_celsius", gpu.TemperatureHeadroomCelsius, 30)
	assertFloat(t, "temperature_shutdown_celsius", gpu.TemperatureShutdownCelsius, 105)
	assertFloat(t, "power_limit_watts", gpu.PowerLimitWatts, 190)
	assertFloat(t, "fan_speed_rpm", gpu.FanSpeedRpm, 1200)
	assertFloat(t, "pcie_link_gen", gpu.PcieLinkGen, 4)
	assertFloat(t, "pcie_link_width", gpu.PcieLinkWidth, 8)
	// The power requires a second sample of the energy counter
	assertNil(t, "power_watts", gpu.PowerWatts)
	if strings.Join(gpu.ClockEventReasons, ",") != "pl1,thermal" || strings.Join(gpu.ThrottleReasons, ",") != "thermal" {
		t.Errorf("unexpected throttle reasons: %v %v", gpu.ClockEventReasons, gpu.ThrottleReasons)
	}

	// 20 J in 2 seconds
	if err := os.WriteFile(filepath.Join(device, "hwmon/hwmon5/energy1_input"), []byte("21000000\n"), 0600); err != nil {
		t.Fatal(err)
	}
	gpuNow = func() time.Time { return start.Add(2 * time.Second) }
	gpu = readIntelGpu(filepath.Join(root, "drm", "card1", "device"), usage, true)
	assertFloat(t, "power_watts", gpu.PowerWatts, 10)
}

func ptrFloat(value float64) *float64 {
	return &value
}

// Arc GPUs report the internal link of the GPU (2.5 GT/s x1), the link of the card is the link of the upstream bridge
func TestReadIntelPcieLink(t *testing.T) {
	root := t.TempDir()
	upstream := filepath.Join(root, "pci0000:00", "0000:00:01.0", "0000:01:00.0")
	downstream := filepath.Join(upstream, "0000:02:01.0")
	gpuDevice := filepath.Join(downstream, "0000:03:00.0")
	files := map[string]map[string]string{
		upstream:   {"vendor": "0x8086", "class": "0x060400", "current_link_speed": "16.0 GT/s PCIe", "max_link_speed": "16.0 GT/s PCIe", "current_link_width": "16", "max_link_width": "16"},
		downstream: {"vendor": "0x8086", "class": "0x060400", "current_link_speed": "2.5 GT/s PCIe", "max_link_speed": "2.5 GT/s PCIe", "current_link_width": "1", "max_link_width": "1"},
		gpuDevice:  {"vendor": "0x8086", "class": "0x030000", "current_link_speed": "2.5 GT/s PCIe", "max_link_speed": "2.5 GT/s PCIe", "current_link_width": "1", "max_link_width": "1"},
	}
	for dir, content := range files {
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatal(err)
		}
		for name, value := range content {
			if err := os.WriteFile(filepath.Join(dir, name), []byte(value+"\n"), 0600); err != nil {
				t.Fatal(err)
			}
		}
	}
	gpu := &resultGpu{}
	readIntelPcieLink(gpuDevice, gpu)
	assertFloat(t, "pcie_link_gen", gpu.PcieLinkGen, 4)
	assertFloat(t, "pcie_link_width_max", gpu.PcieLinkWidthMax, 16)

	// A GPU with a real link of 2.5 GT/s x1 behind a bridge of another vendor
	if err := os.WriteFile(filepath.Join(downstream, "vendor"), []byte("0x1022\n"), 0600); err != nil {
		t.Fatal(err)
	}
	gpu = &resultGpu{}
	readIntelPcieLink(gpuDevice, gpu)
	assertFloat(t, "pcie_link_gen", gpu.PcieLinkGen, 1)
	assertFloat(t, "pcie_link_width", gpu.PcieLinkWidth, 1)
}

// setIntelTestState isolates the global state of the GPU check for a test and returns a fake PCI id database
func setIntelTestState(t *testing.T) {
	t.Helper()
	drmPath, procPath, pciIdsPaths, now := gpuDrmPath, gpuProcPath, gpuPciIdsPaths, gpuNow
	reset := func() {
		drmClientsMutex.Lock()
		drmClientsPrevious = nil
		drmClientsMutex.Unlock()
		intelEnergyMutex.Lock()
		intelEnergyPrevious = map[string]intelEnergySample{}
		intelEnergyMutex.Unlock()
		pciIdsMutex.Lock()
		pciIdsCache = map[string]string{}
		pciIdsMutex.Unlock()
	}
	t.Cleanup(func() {
		gpuDrmPath, gpuProcPath, gpuPciIdsPaths, gpuNow = drmPath, procPath, pciIdsPaths, now
		reset()
	})
	reset()
	pciIds := filepath.Join(t.TempDir(), "pci.ids")
	if err := os.WriteFile(pciIds, []byte("8086  Intel Corporation\n\t56a5  DG2 [Arc A380]\n\ta7a1  Raptor Lake-P [Iris Xe Graphics]\n"), 0600); err != nil {
		t.Fatal(err)
	}
	gpuPciIdsPaths = []string{pciIds}
}

// Arc A380 (i915, Ubuntu 24.04, kernel 6.8) under load: ffmpeg uses the video engines. intel_gpu_top showed Video
// 63.11 % and VideoEnhance 20.78 % at the same time.
func TestCollectIntelGpusArcA380Load(t *testing.T) {
	setIntelTestState(t)
	// Times of the two samples of the clients in the dump (16.6 seconds, the dump script is slow)
	first, second := time.Unix(1793214888, 520801611), time.Unix(1793214905, 157295636)

	gpuNow = func() time.Time { return first }
	gpuDrmPath, gpuProcPath = buildIntelSysfs(t, "arc-a380-i915-load.txt", 0)
	gpus := collectPlatformGpus(nil, nil)
	if len(gpus) != 1 {
		t.Fatalf("expected 1 GPU, got %d", len(gpus))
	}
	gpu := gpus[0]
	if gpu.Name != "Intel DG2 [Arc A380]" || gpu.Source != "i915" || gpu.PciBusID != "0000:06:00.0" || gpu.Integrated || gpu.Processes != 1 {
		t.Errorf("unexpected GPU: %+v", gpu)
	}
	assertNil(t, "utilization_percent", gpu.UtilizationPercent)
	// drm-resident-local0 of the client
	assertUint(t, "memory_used_bytes", gpu.MemoryUsedBytes, 506440*1024)
	assertNil(t, "memory_total_bytes", gpu.MemoryTotalBytes)
	assertFloat(t, "power_limit_watts", gpu.PowerLimitWatts, 55)
	assertFloat(t, "clock_graphics_mhz", gpu.ClockGraphicsMhz, 1400)
	assertFloat(t, "clock_graphics_max_mhz", gpu.ClockGraphicsMaxMhz, 2450)
	// Temperature and fan only with kernel 6.12 and newer
	assertNil(t, "temperature_celsius", gpu.TemperatureCelsius)
	assertNil(t, "fan_speed_rpm", gpu.FanSpeedRpm)
	// The dump has no data of the upstream bridges, so the internal link of the GPU is reported
	assertFloat(t, "pcie_link_width", gpu.PcieLinkWidth, 1)

	gpuNow = func() time.Time { return second }
	gpuDrmPath, gpuProcPath = buildIntelSysfs(t, "arc-a380-i915-load.txt", 1)
	gpu = collectPlatformGpus(nil, nil)[0]
	// Video: (63085203168 - 41519085608) ns / (16.636 s * 2 engines) = 64.82 %
	for name, value := range map[string]*float64{"utilization_percent": gpu.UtilizationPercent, "encoder_utilization_percent": gpu.EncoderUtilizationPercent, "decoder_utilization_percent": gpu.DecoderUtilizationPercent} {
		if value == nil || *value < 64.81 || *value > 64.83 {
			t.Errorf("%s: expected 64.82, got %v", name, deref(value))
		}
	}
}

// Arc A380 idle without clients: 15.9 W from the energy counter
func TestCollectIntelGpusArcA380Idle(t *testing.T) {
	setIntelTestState(t)
	start := time.Unix(1793214374, 0)
	gpuNow = func() time.Time { return start }
	gpuDrmPath, gpuProcPath = buildIntelSysfs(t, "arc-a380-i915-idle.txt", 0)
	gpu := collectPlatformGpus(nil, nil)[0]
	assertNil(t, "power_watts", gpu.PowerWatts)
	assertFloat(t, "clock_graphics_mhz", gpu.ClockGraphicsMhz, 0)

	// The second sample of the energy counter 2 seconds later
	gpuNow = func() time.Time { return start.Add(2 * time.Second) }
	gpuDrmPath, gpuProcPath = buildIntelSysfs(t, "arc-a380-i915-idle.txt", 1)
	gpu = collectPlatformGpus(nil, nil)[0]
	// (245045861877 - 245014123596) µJ / 2 s
	if gpu.PowerWatts == nil || *gpu.PowerWatts < 15.86 || *gpu.PowerWatts > 15.88 {
		t.Errorf("unexpected power: %v", deref(gpu.PowerWatts))
	}
	// No client: idle, no memory known
	assertFloat(t, "utilization_percent", gpu.UtilizationPercent, 0)
	assertNil(t, "memory_used_bytes", gpu.MemoryUsedBytes)
	if gpu.Processes != 0 {
		t.Errorf("unexpected processes: %d", gpu.Processes)
	}
}

// Arc A380 (dump script version 4 with the upstream PCI bridges) under load (ffmpeg). The GPU and the first bridge
// report the internal link (2.5 GT/s x1), the second bridge the link of the card: PCIe 4.0 x8 card in a slot of the
// chipset with PCIe 3.0 x4.
func TestCollectIntelGpusArcA380Bridges(t *testing.T) {
	setIntelTestState(t)
	first, second := time.Unix(1793310590, 362447817), time.Unix(1793310606, 553547019)

	gpuNow = func() time.Time { return first }
	gpuDrmPath, gpuProcPath = buildIntelSysfs(t, "arc-a380-i915-bridges-load.txt", 0)
	gpu := collectPlatformGpus(nil, nil)[0]
	if gpu.PciBusID != "0000:06:00.0" || gpu.Integrated {
		t.Errorf("unexpected GPU: %+v", gpu)
	}
	assertFloat(t, "pcie_link_gen", gpu.PcieLinkGen, 3)
	assertFloat(t, "pcie_link_width", gpu.PcieLinkWidth, 4)
	assertFloat(t, "pcie_link_gen_max", gpu.PcieLinkGenMax, 4)
	assertFloat(t, "pcie_link_width_max", gpu.PcieLinkWidthMax, 8)
	assertUint(t, "memory_used_bytes", gpu.MemoryUsedBytes, 704344*1024)

	gpuNow = func() time.Time { return second }
	gpuDrmPath, gpuProcPath = buildIntelSysfs(t, "arc-a380-i915-bridges-load.txt", 1)
	gpu = collectPlatformGpus(nil, nil)[0]
	// Video: (107505285212 - 87303535852) ns / (16.191 s * 2 engines) = 62.38 %
	if gpu.UtilizationPercent == nil || *gpu.UtilizationPercent < 62.37 || *gpu.UtilizationPercent > 62.40 {
		t.Errorf("unexpected utilization: %v", deref(gpu.UtilizationPercent))
	}
	assertFloat(t, "clock_graphics_mhz", gpu.ClockGraphicsMhz, 1400)
}
