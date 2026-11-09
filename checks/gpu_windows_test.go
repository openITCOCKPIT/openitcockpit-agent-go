package checks

import (
	"context"
	"strings"
	"testing"
)

func TestParseGpuPerfCounters(t *testing.T) {
	engines := []*resultPerfCounterValue{
		// Two processes on the 3D engine of the RX 7800 XT
		{Instance: "pid_10444_luid_0x00000000_0x00013C9C_phys_0_eng_0_engtype_3D", Value: 30},
		{Instance: "pid_2345_luid_0x00000000_0x00013C9C_phys_0_eng_0_engtype_3D", Value: 25.5},
		{Instance: "pid_10444_luid_0x00000000_0x00013C9C_phys_0_eng_1_engtype_Copy", Value: 3},
		{Instance: "pid_10444_luid_0x00000000_0x00013C9C_phys_0_eng_2_engtype_Compute 0", Value: 10},
		{Instance: "pid_10444_luid_0x00000000_0x00013C9C_phys_0_eng_5_engtype_VideoDecode", Value: 12},
		{Instance: "pid_10444_luid_0x00000000_0x00013C9C_phys_0_eng_6_engtype_VideoEncode", Value: 7},
		// A second physical adapter (linked adapters) is ignored
		{Instance: "pid_10444_luid_0x00000000_0x00013C9C_phys_1_eng_0_engtype_3D", Value: 90},
		// The sum of all processes is limited to 100 %
		{Instance: "pid_1_luid_0x00000000_0x000174BB_phys_0_eng_0_engtype_3D", Value: 80},
		{Instance: "pid_2_luid_0x00000000_0x000174BB_phys_0_eng_0_engtype_3D", Value: 40},
		// AMD: one engine for encoding and decoding
		{Instance: "pid_2_luid_0x00000000_0x000174BB_phys_0_eng_10_engtype_Video Codec 0", Value: 15},
		{Instance: "invalid", Value: 50},
	}
	dedicated := []*resultPerfCounterValue{
		{Instance: "luid_0x00000000_0x00013C9C_phys_0", Value: 2740000000},
		{Instance: "luid_0x00000000_0x000174BB_phys_0", Value: 16000000},
	}
	shared := []*resultPerfCounterValue{
		{Instance: "luid_0x00000000_0x00013C9C_phys_0", Value: 70000000},
		{Instance: "luid_0x00000000_0x000174BB_phys_0", Value: 14000000},
	}
	counters := parseGpuPerfCounters(engines, dedicated, shared)
	if len(counters) != 2 {
		t.Fatalf("expected 2 adapters, got %d", len(counters))
	}
	dgpu := counters["0x00000000_0x00013c9c"]
	if dgpu == nil {
		t.Fatalf("missing adapter: %v", counters)
	}
	assertFloat(t, "utilization", dgpu.utilization, 55.5)
	assertFloat(t, "decoder", dgpu.decoder, 12)
	assertFloat(t, "encoder", dgpu.encoder, 7)
	assertUint(t, "dedicated", dgpu.dedicatedUsage, 2740000000)
	assertUint(t, "shared", dgpu.sharedUsage, 70000000)

	igpu := counters["0x00000000_0x000174bb"]
	assertFloat(t, "utilization", igpu.utilization, 100)
	assertFloat(t, "encoder", igpu.encoder, 15)
	assertFloat(t, "decoder", igpu.decoder, 15)
}

func TestD3dkmtAdapterToGpu(t *testing.T) {
	counters := map[string]*gpuPerfCounters{
		"0x00000000_0x00013c9c": {utilization: ptr(55.5), dedicatedUsage: ptr(uint64(2740000000)), sharedUsage: ptr(uint64(70000000))},
		"0x00000000_0x000174bb": {utilization: ptr(4.0), dedicatedUsage: ptr(uint64(16000000)), sharedUsage: ptr(uint64(14000000))},
	}
	// Values of the RX 7800 XT and the Raphael iGPU of the PC of the developer, without the AMD software
	dgpu := d3dkmtAdapterToGpu(&d3dkmtAdapter{
		luid: "0x00000000_0x00013C9C", name: "AMD Radeon RX 7800 XT", biosVersion: "113-EXT90455-100", adapterType: 0x230b,
		vendorID: 0x1002, deviceID: 0x747e, pciAddress: &d3dkmtAdapterAddress{BusNumber: 3},
		segments:      &d3dkmtSegmentSizeInfo{DedicatedVideoMemorySize: 16176 << 20, SharedSystemMemorySize: 15915 << 20},
		driverVersion: "32.0.32015.7029",
		perfData:      &d3dkmtAdapterPerfData{Temperature: 380, FanRPM: 0, Power: 85, MemoryFrequency: 96000000},
		perfCaps:      &d3dkmtAdapterPerfDataCaps{MaxFanRPM: 3100, TemperatureMax: 110},
		graphicsClock: ptr(36.0), graphicsMaxClk: ptr(2213.0),
	}, counters)
	if dgpu.Vendor != "AMD" || dgpu.Source != "d3dkmt" || dgpu.Name != "AMD Radeon RX 7800 XT" || dgpu.PciBusID != "0000:03:00.0" ||
		dgpu.Integrated || dgpu.DriverVersion != "32.0.32015.7029" || dgpu.Processes != -1 {
		t.Errorf("unexpected GPU: %+v", dgpu)
	}
	assertFloat(t, "temperature_celsius", dgpu.TemperatureCelsius, 38)
	// AMD reports the maximum temperature in degrees Celsius instead of tenths, no warning temperature
	assertFloat(t, "temperature_shutdown_celsius", dgpu.TemperatureShutdownCelsius, 110)
	assertNil(t, "temperature_slowdown_celsius", dgpu.TemperatureSlowdownCelsius)
	assertNil(t, "temperature_headroom_celsius", dgpu.TemperatureHeadroomCelsius)
	assertFloat(t, "fan_speed_percent", dgpu.FanSpeedPercent, 0)
	assertFloat(t, "fan_speed_rpm", dgpu.FanSpeedRpm, 0)
	assertFloat(t, "power_percent", dgpu.PowerPercent, 8.5)

	// The AMD driver reports the power 0 also under full load, 0 is unknown
	amd := d3dkmtAdapterToGpu(&d3dkmtAdapter{
		luid: "0x00000000_0x00013C9C", name: "AMD Radeon RX 7800 XT", vendorID: 0x1002,
		segments: &d3dkmtSegmentSizeInfo{DedicatedVideoMemorySize: 16176 << 20},
		perfData: &d3dkmtAdapterPerfData{Temperature: 620, Power: 0},
	}, counters)
	assertNil(t, "power_percent", amd.PowerPercent)
	assertNil(t, "power_watts", dgpu.PowerWatts)
	assertFloat(t, "clock_memory_mhz", dgpu.ClockMemoryMhz, 96)
	assertFloat(t, "clock_graphics_max_mhz", dgpu.ClockGraphicsMaxMhz, 2213)
	assertFloat(t, "utilization_percent", dgpu.UtilizationPercent, 55.5)
	assertUint(t, "memory_total_bytes", dgpu.MemoryTotalBytes, 16176<<20)
	assertUint(t, "memory_used_bytes", dgpu.MemoryUsedBytes, 2740000000)

	igpu := d3dkmtAdapterToGpu(&d3dkmtAdapter{
		luid: "0x00000000_0x000174BB", name: "AMD Radeon(TM) Graphics", adapterType: 0x2303, vendorID: 0x1002,
		pciAddress: &d3dkmtAdapterAddress{BusNumber: 17},
		segments:   &d3dkmtSegmentSizeInfo{DedicatedVideoMemorySize: 485 << 20, SharedSystemMemorySize: 15915 << 20},
		perfData:   &d3dkmtAdapterPerfData{Temperature: 440, Power: 300},
	}, counters)
	if !igpu.Integrated || igpu.PciBusID != "0000:11:00.0" {
		t.Errorf("unexpected GPU: %+v", igpu)
	}
	// Integrated GPUs: dedicated and shared memory, the total is the system memory, no power
	assertUint(t, "memory_used_bytes", igpu.MemoryUsedBytes, 30000000)
	assertNil(t, "memory_total_bytes", igpu.MemoryTotalBytes)
	assertNil(t, "power_percent", igpu.PowerPercent)
	assertNil(t, "fan_speed_percent", igpu.FanSpeedPercent)
	assertNil(t, "fan_speed_rpm", igpu.FanSpeedRpm)
	assertFloat(t, "temperature_celsius", igpu.TemperatureCelsius, 44)
}

func TestD3dkmtCapsTemperature(t *testing.T) {
	// Documented in tenths of a degree Celsius, but the AMD driver reports degrees Celsius
	tests := map[uint32]interface{}{0: nil, 110: 110.0, 950: 95.0, 1050: 105.0}
	for value, expected := range tests {
		if got := deref(d3dkmtCapsTemperature(value)); got != expected {
			t.Errorf("%d: expected %v, got %v", value, expected, got)
		}
	}
}

func TestIsD3dkmtGpu(t *testing.T) {
	tests := map[uint32]bool{
		0x230b: true,  // AMD Radeon RX 7800 XT
		0x2303: true,  // AMD Radeon iGPU (Raphael)
		0x0800: true,  // Compute only, e.g. data center GPUs without display
		0x0105: false, // Microsoft Basic Render Driver (software)
		0x0342: false, // Remote desktop display (indirect display) of an Azure VM
		0x010a: false, // Display adapter of a Hyper-V VM (Azure, Windows 11 ARM64), no rendering
		0x0000: false,
	}
	for adapterType, expected := range tests {
		if got := isD3dkmtGpu(adapterType); got != expected {
			t.Errorf("0x%x: expected %v, got %v", adapterType, expected, got)
		}
	}
}

func TestIsValidPciAddress(t *testing.T) {
	tests := []struct {
		address  d3dkmtAdapterAddress
		expected bool
	}{
		{d3dkmtAdapterAddress{BusNumber: 3}, true},
		{d3dkmtAdapterAddress{BusNumber: 0xff, DeviceNumber: 31, FunctionNumber: 7}, true},
		{d3dkmtAdapterAddress{BusNumber: 0xFFFFFFFF, DeviceNumber: 0xffff, FunctionNumber: 0xffff}, false},
		{d3dkmtAdapterAddress{BusNumber: 0, DeviceNumber: 0xffff, FunctionNumber: 0xffff}, false},
		{d3dkmtAdapterAddress{BusNumber: 0, DeviceNumber: 0x5620, FunctionNumber: 0xe0c7}, false},
	}
	for _, test := range tests {
		if got := isValidPciAddress(&test.address); got != test.expected {
			t.Errorf("%+v: expected %v, got %v", test.address, test.expected, got)
		}
	}
}

func ptr[T any](value T) *T {
	return &value
}

// Reads the GPUs of the system, e.g. of the developer PC or a VM without GPU
func TestCollectPlatformGpusWindows(t *testing.T) {
	if !hasPlatformGpus() {
		t.Skip("no GPU")
	}
	gpus := collectPlatformGpus(context.Background(), nil)
	if len(gpus) == 0 {
		t.Fatal("expected GPUs")
	}
	for _, gpu := range gpus {
		if gpu.Name == "" || gpu.Vendor == "" || gpu.Source != "d3dkmt" || strings.Contains(gpu.Name, "Basic Render") {
			t.Errorf("unexpected GPU: %+v", gpu)
		}
		t.Logf("%s (%s, %s): integrated=%v temperature=%v utilization=%v memory=%v/%v", gpu.Name, gpu.Vendor, gpu.PciBusID,
			gpu.Integrated, deref(gpu.TemperatureCelsius), deref(gpu.UtilizationPercent), deref(gpu.MemoryUsedBytes), deref(gpu.MemoryTotalBytes))
	}

	// GPUs of nvidia-smi are skipped
	skipped := collectPlatformGpus(context.Background(), []*resultGpu{{PciBusID: "00000000:" + strings.TrimPrefix(gpus[0].PciBusID, "0000:")}})
	if gpus[0].PciBusID != "" && len(skipped) != len(gpus)-1 {
		t.Errorf("expected %d GPUs, got %d", len(gpus)-1, len(skipped))
	}
}

func TestReadGpuDriverEvents(t *testing.T) {
	events := readGpuDriverEvents(context.Background())
	if events == nil || events.Error != "" || events.Events == nil || events.Count < len(events.Events) {
		t.Errorf("unexpected result: %+v", events)
	}
}
