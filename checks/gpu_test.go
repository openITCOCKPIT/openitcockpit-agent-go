package checks

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func readNvidiaTestdata(t *testing.T, name string) ([]*resultGpu, *resultGpusNvidia) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "gpu", "nvidia", name))
	if err != nil {
		t.Fatal(err)
	}
	nvidia := &resultGpusNvidia{}
	gpus, err := parseNvidiaSmi(data, nvidia)
	if err != nil {
		t.Fatal(err)
	}
	return gpus, nvidia
}

func assertFloat(t *testing.T, name string, value *float64, expected float64) {
	t.Helper()
	if value == nil || *value != expected {
		t.Errorf("%s: expected %v, got %v", name, expected, deref(value))
	}
}

func assertUint(t *testing.T, name string, value *uint64, expected uint64) {
	t.Helper()
	if value == nil || *value != expected {
		t.Errorf("%s: expected %v, got %v", name, expected, deref(value))
	}
}

func assertBool(t *testing.T, name string, value *bool, expected bool) {
	t.Helper()
	if value == nil || *value != expected {
		t.Errorf("%s: expected %v, got %v", name, expected, deref(value))
	}
}

func assertNil[T any](t *testing.T, name string, value *T) {
	t.Helper()
	if value != nil {
		t.Errorf("%s: expected null, got %v", name, *value)
	}
}

func deref[T any](value *T) interface{} {
	if value == nil {
		return nil
	}
	return *value
}

// 4x H200 NVL with NVLink, driver 595 (Linux)
func TestParseNvidiaSmiH200(t *testing.T) {
	gpus, nvidia := readNvidiaTestdata(t, "h200-nvl-4x-linux-595.91.07.xml")
	if nvidia.DriverVersion != "595.91.07" || nvidia.CudaVersion != "13.2" || nvidia.AttachedGpus != 4 || len(gpus) != 4 {
		t.Fatalf("unexpected result: %+v, %d GPUs", nvidia, len(gpus))
	}
	gpu := gpus[0]
	if gpu.Vendor != "NVIDIA" || gpu.Name != "NVIDIA H200 NVL" || gpu.PciBusID != "0000:18:00.0" || gpu.Serial != "SERIAL0001" ||
		gpu.UUID != "GPU-00000000-0000-0000-0000-000000000001" || gpu.PerformanceState != "P0" || gpu.Integrated ||
		gpu.PersistenceMode != "Disabled" || gpu.MigMode != "Disabled" || gpu.RecoveryAction != "None" || gpu.Processes != 1 {
		t.Errorf("unexpected GPU: %+v", gpu)
	}
	assertUint(t, "memory_total_bytes", gpu.MemoryTotalBytes, 143771*1024*1024)
	assertUint(t, "memory_used_bytes", gpu.MemoryUsedBytes, 132004*1024*1024)
	assertFloat(t, "temperature_celsius", gpu.TemperatureCelsius, 35)
	assertFloat(t, "temperature_memory_celsius", gpu.TemperatureMemoryCelsius, 37)
	// gpu_temp_tlimit 52 C, the GPU slows down at -2 C
	assertFloat(t, "temperature_headroom_celsius", gpu.TemperatureHeadroomCelsius, 54)
	assertNil(t, "temperature_slowdown_celsius", gpu.TemperatureSlowdownCelsius)
	assertNil(t, "fan_speed_percent", gpu.FanSpeedPercent)
	assertFloat(t, "power_watts", gpu.PowerWatts, 99.2)
	assertFloat(t, "power_limit_watts", gpu.PowerLimitWatts, 600)
	assertFloat(t, "clock_graphics_mhz", gpu.ClockGraphicsMhz, 1785)
	assertBool(t, "ecc_enabled", gpu.EccEnabled, true)
	assertUint(t, "ecc_uncorrected_volatile", gpu.EccUncorrectedVolatile, 0)
	assertUint(t, "ecc_uncorrected_aggregate", gpu.EccUncorrectedAggregate, 0)
	assertNil(t, "retired_pages_double_bit", gpu.RetiredPagesDoubleBit)
	assertUint(t, "remapped_rows_uncorrectable", gpu.RemappedRowsUncorrectable, 0)
	assertBool(t, "remapped_rows_pending", gpu.RemappedRowsPending, false)
	assertBool(t, "remapped_rows_failure", gpu.RemappedRowsFailure, false)
	assertBool(t, "memory_repair_pending", gpu.MemoryRepairPending, false)
	assertBool(t, "memory_unrepairable", gpu.MemoryUnrepairable, false)
	assertFloat(t, "pcie_link_gen", gpu.PcieLinkGen, 5)
	assertFloat(t, "pcie_link_width", gpu.PcieLinkWidth, 16)
	if len(gpu.ThrottleReasons) != 0 || len(gpu.ClockEventReasons) != 0 || len(gpu.Messages) != 0 {
		t.Errorf("unexpected reasons or messages: %+v", gpu)
	}
	for i, gpu := range gpus {
		if gpu.UUID == "" || gpu.PciBusID == "" || gpu.MemoryTotalBytes == nil || gpu.TemperatureCelsius == nil || gpu.PowerWatts == nil {
			t.Errorf("GPU %d: missing values: %+v", i, gpu)
		}
	}
}

// DGX Spark (GB10) with unified memory: nvidia-smi reports no memory and no power limit
func TestParseNvidiaSmiDgxSpark(t *testing.T) {
	gpus, nvidia := readNvidiaTestdata(t, "dgx-spark-gb10-linux-595.99.02.xml")
	if nvidia.DriverVersion != "595.99.02" || len(gpus) != 1 {
		t.Fatalf("unexpected result: %+v, %d GPUs", nvidia, len(gpus))
	}
	gpu := gpus[0]
	if gpu.Name != "NVIDIA GB10" || gpu.PciBusID != "000f:01:00.0" || !gpu.Integrated || gpu.Serial != "" || gpu.Processes != 3 {
		t.Errorf("unexpected GPU: %+v", gpu)
	}
	// The memory of the processes: 19 + 6 + 27426 MiB
	assertNil(t, "memory_total_bytes", gpu.MemoryTotalBytes)
	assertUint(t, "memory_used_bytes", gpu.MemoryUsedBytes, 27451*1024*1024)
	assertFloat(t, "temperature_celsius", gpu.TemperatureCelsius, 43)
	assertFloat(t, "temperature_headroom_celsius", gpu.TemperatureHeadroomCelsius, 53)
	assertFloat(t, "power_watts", gpu.PowerWatts, 10.7)
	assertNil(t, "power_limit_watts", gpu.PowerLimitWatts)
	assertNil(t, "clock_memory_mhz", gpu.ClockMemoryMhz)
	assertNil(t, "ecc_enabled", gpu.EccEnabled)
	assertNil(t, "ecc_uncorrected_volatile", gpu.EccUncorrectedVolatile)
	assertNil(t, "memory_repair_pending", gpu.MemoryRepairPending)
	assertFloat(t, "pcie_link_width", gpu.PcieLinkWidth, 1)
}

// GeForce RTX 5090 on Windows (WDDM): the clock event reason "reliability" is active when idle
func TestParseNvidiaSmiRtx5090Windows(t *testing.T) {
	gpus, nvidia := readNvidiaTestdata(t, "rtx-5090-windows-617.14.xml")
	if nvidia.DriverVersion != "617.14" || len(gpus) != 1 {
		t.Fatalf("unexpected result: %+v, %d GPUs", nvidia, len(gpus))
	}
	gpu := gpus[0]
	if gpu.Name != "NVIDIA GeForce RTX 5090" || gpu.DriverModel != "WDDM" || gpu.Serial != "" || gpu.Integrated {
		t.Errorf("unexpected GPU: %+v", gpu)
	}
	if !reflect.DeepEqual(gpu.ClockEventReasons, []string{"reliability"}) || len(gpu.ThrottleReasons) != 0 {
		t.Errorf("unexpected clock event reasons: %v, throttle reasons: %v", gpu.ClockEventReasons, gpu.ThrottleReasons)
	}
	assertUint(t, "memory_total_bytes", gpu.MemoryTotalBytes, 32607*1024*1024)
	assertFloat(t, "fan_speed_percent", gpu.FanSpeedPercent, 0)
	assertFloat(t, "temperature_headroom_celsius", gpu.TemperatureHeadroomCelsius, 39)
	assertFloat(t, "power_watts", gpu.PowerWatts, 71.98)
	assertNil(t, "ecc_enabled", gpu.EccEnabled)
	assertNil(t, "remapped_rows_pending", gpu.RemappedRowsPending)
}

// GeForce RTX 5070 Ti on Windows: PCIe Gen 1 when idle
func TestParseNvidiaSmiRtx5070TiWindows(t *testing.T) {
	gpus, _ := readNvidiaTestdata(t, "rtx-5070-ti-windows-610.88.xml")
	if len(gpus) != 1 {
		t.Fatalf("expected 1 GPU, got %d", len(gpus))
	}
	gpu := gpus[0]
	if gpu.PerformanceState != "P8" || gpu.PciBusID != "0000:2d:00.0" {
		t.Errorf("unexpected GPU: %+v", gpu)
	}
	assertFloat(t, "pcie_link_gen", gpu.PcieLinkGen, 1)
	assertFloat(t, "pcie_link_gen_max", gpu.PcieLinkGenMax, 4)
	assertFloat(t, "power_limit_watts", gpu.PowerLimitWatts, 300)
}

// Elements of older drivers (before 535): power_readings, clocks_throttle_reasons, single and double bit ECC
// errors, absolute temperature thresholds
const testNvidiaSmiLegacy = `<?xml version="1.0" ?>
<!DOCTYPE nvidia_smi_log SYSTEM "nvsmi_device_v11.dtd">
<nvidia_smi_log>
	<driver_version>470.199.02</driver_version>
	<cuda_version>11.4</cuda_version>
	<attached_gpus>2</attached_gpus>
	<gpu id="00000000:3B:00.0">
		<product_name>Tesla V100-PCIE-32GB</product_name>
		<serial>0323418012345</serial>
		<uuid>GPU-11111111-2222-3333-4444-555555555555</uuid>
		<performance_state>P0</performance_state>
		<clocks_throttle_reasons>
			<clocks_throttle_reason_gpu_idle>Not Active</clocks_throttle_reason_gpu_idle>
			<clocks_throttle_reason_sw_power_cap>Active</clocks_throttle_reason_sw_power_cap>
			<clocks_throttle_reason_hw_slowdown>Active</clocks_throttle_reason_hw_slowdown>
			<clocks_throttle_reason_hw_thermal_slowdown>Not Active</clocks_throttle_reason_hw_thermal_slowdown>
		</clocks_throttle_reasons>
		<fb_memory_usage>
			<total>32510 MiB</total>
			<used>1024 MiB</used>
			<free>31486 MiB</free>
		</fb_memory_usage>
		<ecc_mode>
			<current_ecc>Enabled</current_ecc>
			<pending_ecc>Enabled</pending_ecc>
		</ecc_mode>
		<ecc_errors>
			<volatile>
				<single_bit>
					<device_memory>3</device_memory>
					<total>5</total>
				</single_bit>
				<double_bit>
					<device_memory>1</device_memory>
					<total>1</total>
				</double_bit>
			</volatile>
			<aggregate>
				<single_bit>
					<total>12</total>
				</single_bit>
				<double_bit>
					<total>2</total>
				</double_bit>
			</aggregate>
		</ecc_errors>
		<retired_pages>
			<multiple_single_bit_retirement>
				<retired_count>0</retired_count>
				<retired_pagelist>
				</retired_pagelist>
			</multiple_single_bit_retirement>
			<double_bit_retirement>
				<retired_count>2</retired_count>
				<retired_pagelist>
				</retired_pagelist>
			</double_bit_retirement>
			<pending_blacklist>Yes</pending_blacklist>
			<pending_retirement>Yes</pending_retirement>
		</retired_pages>
		<remapped_rows>N/A</remapped_rows>
		<temperature>
			<gpu_temp>84 C</gpu_temp>
			<gpu_temp_max_threshold>90 C</gpu_temp_max_threshold>
			<gpu_temp_slow_threshold>87 C</gpu_temp_slow_threshold>
			<gpu_temp_max_gpu_threshold>83 C</gpu_temp_max_gpu_threshold>
			<memory_temp>80 C</memory_temp>
		</temperature>
		<power_readings>
			<power_state>P0</power_state>
			<power_management>Supported</power_management>
			<power_draw>249.50 W</power_draw>
			<power_limit>250.00 W</power_limit>
			<default_power_limit>250.00 W</default_power_limit>
			<enforced_power_limit>250.00 W</enforced_power_limit>
		</power_readings>
		<fan_speed>N/A</fan_speed>
	</gpu>
	<gpu id="00000000:AF:00.0">
		<product_name>Tesla V100-PCIE-32GB</product_name>
		<uuid>GPU-66666666-2222-3333-4444-555555555555</uuid>
		<performance_state>[GPU requires reset]</performance_state>
		<fb_memory_usage>
			<total>[GPU requires reset]</total>
			<used>[GPU requires reset]</used>
		</fb_memory_usage>
		<temperature>
			<gpu_temp>[GPU requires reset]</gpu_temp>
		</temperature>
		<gpu_recovery_action>Reset</gpu_recovery_action>
	</gpu>
</nvidia_smi_log>
`

func TestParseNvidiaSmiLegacy(t *testing.T) {
	nvidia := &resultGpusNvidia{}
	gpus, err := parseNvidiaSmi([]byte(testNvidiaSmiLegacy), nvidia)
	if err != nil {
		t.Fatal(err)
	}
	if nvidia.DriverVersion != "470.199.02" || nvidia.AttachedGpus != 2 || len(gpus) != 2 {
		t.Fatalf("unexpected result: %+v, %d GPUs", nvidia, len(gpus))
	}
	gpu := gpus[0]
	if gpu.PciBusID != "0000:3b:00.0" || gpu.Serial != "0323418012345" {
		t.Errorf("unexpected GPU: %+v", gpu)
	}
	if !reflect.DeepEqual(gpu.ClockEventReasons, []string{"sw_power_cap", "hw_slowdown"}) || !reflect.DeepEqual(gpu.ThrottleReasons, []string{"hw_slowdown"}) {
		t.Errorf("unexpected clock event reasons: %v, throttle reasons: %v", gpu.ClockEventReasons, gpu.ThrottleReasons)
	}
	assertFloat(t, "temperature_slowdown_celsius", gpu.TemperatureSlowdownCelsius, 87)
	assertFloat(t, "temperature_shutdown_celsius", gpu.TemperatureShutdownCelsius, 90)
	assertFloat(t, "temperature_headroom_celsius", gpu.TemperatureHeadroomCelsius, 3)
	assertFloat(t, "power_watts", gpu.PowerWatts, 249.5)
	assertFloat(t, "power_limit_watts", gpu.PowerLimitWatts, 250)
	assertUint(t, "ecc_corrected_volatile", gpu.EccCorrectedVolatile, 5)
	assertUint(t, "ecc_uncorrected_volatile", gpu.EccUncorrectedVolatile, 1)
	assertUint(t, "ecc_corrected_aggregate", gpu.EccCorrectedAggregate, 12)
	assertUint(t, "ecc_uncorrected_aggregate", gpu.EccUncorrectedAggregate, 2)
	assertUint(t, "retired_pages_double_bit", gpu.RetiredPagesDoubleBit, 2)
	assertBool(t, "retired_pages_pending", gpu.RetiredPagesPending, true)
	assertNil(t, "remapped_rows_pending", gpu.RemappedRowsPending)
	assertNil(t, "fan_speed_percent", gpu.FanSpeedPercent)

	// A GPU which requires a reset is reported with a message instead of the values
	gpu = gpus[1]
	if gpu.RecoveryAction != "Reset" || gpu.PerformanceState != "" || !reflect.DeepEqual(gpu.Messages, []string{"GPU requires reset"}) {
		t.Errorf("unexpected GPU: %+v", gpu)
	}
	assertNil(t, "temperature_celsius", gpu.TemperatureCelsius)
	assertNil(t, "memory_total_bytes", gpu.MemoryTotalBytes)
}

func TestParseNvidiaSmiInvalid(t *testing.T) {
	for _, output := range []string{"", "NVIDIA-SMI has failed because it couldn't communicate with the NVIDIA driver.", "<foo></foo>"} {
		if _, err := parseNvidiaSmi([]byte(output), &resultGpusNvidia{}); err == nil {
			t.Errorf("%q: expected an error", output)
		}
	}
}

func TestParseNvidiaNumber(t *testing.T) {
	tests := map[string]interface{}{
		"35 C": 35.0, "-2 C": -2.0, "99.20 W": 99.2, "143771 MiB": 143771.0, "0 %": 0.0, "16x": 16.0,
		"N/A": nil, "[N/A]": nil, "": nil, "Requested functionality has been deprecated": nil,
	}
	for value, expected := range tests {
		if got := deref(parseNvidiaNumber(value)); got != expected {
			t.Errorf("%q: expected %v, got %v", value, expected, got)
		}
	}
}

func TestNormalizePciBusID(t *testing.T) {
	tests := map[string]string{
		"00000000:18:00.0": "0000:18:00.0",
		"0000000F:01:00.0": "000f:01:00.0",
		"0000:3b:00":       "0000:3b:00.0",
		"0000:05:00.0":     "0000:05:00.0",
		"foo":              "foo",
	}
	for id, expected := range tests {
		if got := normalizePciBusID(id); got != expected {
			t.Errorf("%q: expected %q, got %q", id, expected, got)
		}
	}
}

func TestAssignXidErrors(t *testing.T) {
	gpus := []*resultGpu{{PciBusID: "0000:18:00.0", XidErrors: []*resultGpuXid{}}, {PciBusID: "0000:3b:00.0", XidErrors: []*resultGpuXid{}}}
	xids := []*resultGpuXid{
		{Timestamp: 100, Type: "Xid", Code: 13, PciBusID: "0000:18:00.0"},
		{Timestamp: 200, Type: "Xid", Code: 79, PciBusID: "0000:af:00.0"},
		{Timestamp: 300, Type: "SXid", Code: 12028, PciBusID: "0000:18:00.0"},
	}
	for i := 0; i < gpuMaxXidErrors+5; i++ {
		xids = append(xids, &resultGpuXid{Timestamp: int64(1000 + i), Type: "Xid", Code: 31, PciBusID: "0000:3b:00.0"})
	}
	unassigned := assignXidErrors(gpus, xids)

	if gpus[0].XidCount != 1 || len(gpus[0].XidErrors) != 1 || gpus[0].XidErrors[0].Code != 13 {
		t.Errorf("unexpected Xid errors of GPU 0: %d %+v", gpus[0].XidCount, gpus[0].XidErrors)
	}
	// Only the newest Xid errors are reported, but all are counted
	if gpus[1].XidCount != gpuMaxXidErrors+5 || len(gpus[1].XidErrors) != gpuMaxXidErrors || gpus[1].XidErrors[0].Timestamp != int64(1000+gpuMaxXidErrors+4) {
		t.Errorf("unexpected Xid errors of GPU 1: %d %+v", gpus[1].XidCount, gpus[1].XidErrors[0])
	}
	// Newest first: the SXid of the NVSwitch and the Xid of the unknown GPU
	if len(unassigned) != 2 || unassigned[0].Code != 12028 || unassigned[1].Code != 79 {
		t.Errorf("unexpected unassigned Xid errors: %+v", unassigned)
	}
}
