package checks

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

// NVIDIA GPUs are read with "nvidia-smi -q -x", the fields of the output are described in the documentation of
// nvidia-smi: https://docs.nvidia.com/deploy/nvidia-smi/index.html

// Default locations of nvidia-smi, which are not always part of the PATH of the agent service
var nvidiaSmiLocations = []string{
	"/usr/bin/nvidia-smi",
	"/usr/local/bin/nvidia-smi",
	`C:\Windows\System32\nvidia-smi.exe`,
	`C:\Program Files\NVIDIA Corporation\NVSMI\nvidia-smi.exe`,
}

// Clock event reasons which slow the GPU down because of a problem. The other reasons are normal, e.g. gpu_idle,
// sw_power_cap (power limit reached under load) or reliability (Blackwell, also active when idle). See "Clocks Event
// Reasons" in https://docs.nvidia.com/deploy/nvidia-smi/index.html, e.g. "HW Slowdown is engaged, reducing the core
// clocks by a factor of 2 or more".
var nvidiaThrottleReasons = map[string]bool{
	"hw_slowdown":             true,
	"hw_thermal_slowdown":     true,
	"hw_power_brake_slowdown": true,
	"sw_thermal_slowdown":     true,
}

func findNvidiaSmi() string {
	if path, err := exec.LookPath("nvidia-smi"); err == nil {
		return path
	}
	for _, path := range nvidiaSmiLocations {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// runNvidiaSmi returns the GPUs of "nvidia-smi -q -x". If some GPUs can not be read (e.g. a GPU has fallen off the
// bus), nvidia-smi exits with an error but still reports the other GPUs.
func runNvidiaSmi(ctx context.Context, nvidiaSmi string, nvidia *resultGpusNvidia) ([]*resultGpu, error) {
	ctx, cancel := context.WithTimeout(ctx, gpuCommandTimeout)
	defer cancel()

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, nvidiaSmi, "-q", "-x")
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	runErr := cmd.Run()

	gpus, parseErr := parseNvidiaSmi(stdout.Bytes(), nvidia)
	if runErr != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = strings.TrimSpace(stdout.String())
		}
		if parseErr != nil {
			gpus = nil
		}
		return gpus, fmt.Errorf("nvidia-smi -q -x failed: %s: %s", runErr, firstLine(message))
	}
	return gpus, parseErr
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return strings.TrimSpace(s[:i])
	}
	return s
}

// nvidiaXMLNode is a generic element of the XML output of nvidia-smi. The elements differ between the driver
// versions, so the values are looked up by their path with fallbacks for older drivers.
type nvidiaXMLNode struct {
	XMLName xml.Name
	Attrs   []xml.Attr       `xml:",any,attr"`
	Text    string           `xml:",chardata"`
	Nodes   []*nvidiaXMLNode `xml:",any"`
}

// find returns the first element of the slash separated path
func (n *nvidiaXMLNode) find(path string) *nvidiaXMLNode {
	current := n
	for _, name := range strings.Split(path, "/") {
		var next *nvidiaXMLNode
		for _, child := range current.Nodes {
			if child.XMLName.Local == name {
				next = child
				break
			}
		}
		if next == nil {
			return nil
		}
		current = next
	}
	return current
}

// findAll returns all elements with the name below the slash separated path
func (n *nvidiaXMLNode) findAll(path, name string) []*nvidiaXMLNode {
	parent := n.find(path)
	if parent == nil {
		return nil
	}
	var nodes []*nvidiaXMLNode
	for _, child := range parent.Nodes {
		if child.XMLName.Local == name {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// value returns the first available value of the paths, "" if no value is available
func (n *nvidiaXMLNode) value(paths ...string) string {
	for _, path := range paths {
		if node := n.find(path); node != nil && nvidiaAvailable(node.Text) {
			return strings.TrimSpace(node.Text)
		}
	}
	return ""
}

// number returns the first available number of the paths, e.g. 35 of "35 C" or 16 of "16x"
func (n *nvidiaXMLNode) number(paths ...string) *float64 {
	for _, path := range paths {
		if value := parseNvidiaNumber(n.value(path)); value != nil {
			return value
		}
	}
	return nil
}

// sum returns the sum of all available numbers of the paths, nil if no number is available
func (n *nvidiaXMLNode) sum(paths ...string) *uint64 {
	var sum uint64
	found := false
	for _, path := range paths {
		if value := parseNvidiaNumber(n.value(path)); value != nil && *value >= 0 {
			sum += uint64(*value)
			found = true
		}
	}
	if !found {
		return nil
	}
	return &sum
}

// mebibytes returns the first available value of the paths in bytes, e.g. of "143771 MiB"
func (n *nvidiaXMLNode) mebibytes(paths ...string) *uint64 {
	value := n.number(paths...)
	if value == nil || *value < 0 {
		return nil
	}
	result := uint64(*value) * 1024 * 1024
	return &result
}

// yesNo returns true for Yes, Enabled or Active, false for No, Disabled or Not Active and nil otherwise
func (n *nvidiaXMLNode) yesNo(paths ...string) *bool {
	value := strings.ToLower(n.value(paths...))
	var result bool
	switch value {
	case "yes", "enabled", "active":
		result = true
	case "no", "disabled", "not active":
		result = false
	default:
		return nil
	}
	return &result
}

func (n *nvidiaXMLNode) attr(name string) string {
	for _, attr := range n.Attrs {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}

// nvidiaAvailable returns false for the placeholders of nvidia-smi, e.g. N/A, [N/A], Not Supported, [GPU requires
// reset] or "Requested functionality has been deprecated"
func nvidiaAvailable(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" || strings.HasPrefix(value, "[") || strings.EqualFold(value, "N/A") {
		return false
	}
	lower := strings.ToLower(value)
	return !strings.Contains(lower, "not supported") && !strings.Contains(lower, "deprecated") &&
		!strings.Contains(lower, "unknown error") && !strings.Contains(lower, "not available")
}

// parseNvidiaNumber returns the number of values like "35 C", "-2 C", "99.20 W", "143771 MiB", "0 %" or "16x"
func parseNvidiaNumber(value string) *float64 {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return nil
	}
	number, err := strconv.ParseFloat(strings.TrimSuffix(fields[0], "x"), 64)
	if err != nil {
		return nil
	}
	return &number
}

// Messages of nvidia-smi in the values of a GPU which require an action
var nvidiaProblemValues = []string{"GPU requires reset", "GPU is lost", "Unknown Error"}

func parseNvidiaSmi(output []byte, nvidia *resultGpusNvidia) ([]*resultGpu, error) {
	var root nvidiaXMLNode
	if err := xml.Unmarshal(output, &root); err != nil {
		return nil, fmt.Errorf("could not parse output of nvidia-smi -q -x: %s", err)
	}
	if root.XMLName.Local != "nvidia_smi_log" {
		return nil, fmt.Errorf("unexpected output of nvidia-smi -q -x: %s", root.XMLName.Local)
	}

	nvidia.DriverVersion = root.value("driver_version")
	nvidia.CudaVersion = root.value("cuda_version")
	if attached := root.number("attached_gpus"); attached != nil {
		nvidia.AttachedGpus = int(*attached)
	}

	gpus := []*resultGpu{}
	for _, node := range root.Nodes {
		if node.XMLName.Local == "gpu" {
			gpus = append(gpus, parseNvidiaGpu(node, nvidia.DriverVersion))
		}
	}
	return gpus, nil
}

func parseNvidiaGpu(g *nvidiaXMLNode, driverVersion string) *resultGpu {
	gpu := &resultGpu{
		Vendor:            "NVIDIA",
		Source:            "nvidia-smi",
		Name:              g.value("product_name"),
		PciBusID:          normalizePciBusID(g.attr("id")),
		UUID:              g.value("uuid"),
		Serial:            g.value("serial"),
		VbiosVersion:      g.value("vbios_version"),
		DriverVersion:     driverVersion,
		PerformanceState:  g.value("performance_state"),
		PersistenceMode:   g.value("persistence_mode"),
		DriverModel:       g.value("driver_model/current_dm"),
		MigMode:           g.value("mig_mode/current_mig"),
		ComputeMode:       g.value("compute_mode"),
		FabricState:       g.value("fabric/state"),
		FabricStatus:      g.value("fabric/status"),
		RecoveryAction:    g.value("gpu_recovery_action"),
		ThrottleReasons:   []string{},
		ClockEventReasons: []string{},
		XidErrors:         []*resultGpuXid{},
		Messages:          []string{},
	}
	if strings.Trim(gpu.Serial, "0") == "" {
		// GeForce GPUs report the serial number 0
		gpu.Serial = ""
	}
	if gpu.PciBusID == "" {
		gpu.PciBusID = normalizePciBusID(g.value("pci/pci_bus_id"))
	}

	gpu.UtilizationPercent = g.number("utilization/gpu_util")
	gpu.MemoryUtilizationPercent = g.number("utilization/memory_util")
	gpu.EncoderUtilizationPercent = g.number("utilization/encoder_util")
	gpu.DecoderUtilizationPercent = g.number("utilization/decoder_util")

	gpu.MemoryTotalBytes = g.mebibytes("fb_memory_usage/total")
	gpu.MemoryUsedBytes = g.mebibytes("fb_memory_usage/used")

	processes := g.findAll("processes", "process_info")
	gpu.Processes = len(processes)

	if gpu.MemoryTotalBytes == nil {
		// Unified memory (DGX Spark, GB10): the GPU has no own memory, nvidia-smi only reports the memory of the
		// processes
		gpu.Integrated = true
		var used uint64
		for _, process := range processes {
			if memory := process.mebibytes("used_memory"); memory != nil {
				used += *memory
			}
		}
		gpu.MemoryUsedBytes = &used
	}

	parseNvidiaTemperature(g, gpu)

	gpu.FanSpeedPercent = g.number("fan_speed")
	gpu.PowerWatts = g.number(
		"gpu_power_readings/average_power_draw", "gpu_power_readings/instant_power_draw", "gpu_power_readings/power_draw",
		"power_readings/average_power_draw", "power_readings/instant_power_draw", "power_readings/power_draw",
	)
	gpu.PowerLimitWatts = g.number(
		"gpu_power_readings/current_power_limit", "gpu_power_readings/enforced_power_limit",
		"power_readings/enforced_power_limit", "power_readings/power_limit",
	)

	gpu.ClockGraphicsMhz = g.number("clocks/graphics_clock")
	gpu.ClockGraphicsMaxMhz = g.number("max_clocks/graphics_clock")
	gpu.ClockMemoryMhz = g.number("clocks/mem_clock")

	// Driver 535 renamed clocks_throttle_reasons to clocks_event_reasons
	for _, path := range []string{"clocks_event_reasons", "clocks_throttle_reasons"} {
		reasons := g.find(path)
		if reasons == nil {
			continue
		}
		for _, reason := range reasons.Nodes {
			if strings.TrimSpace(reason.Text) != "Active" {
				continue
			}
			name := strings.TrimPrefix(strings.TrimPrefix(reason.XMLName.Local, "clocks_event_reason_"), "clocks_throttle_reason_")
			gpu.ClockEventReasons = append(gpu.ClockEventReasons, name)
			if nvidiaThrottleReasons[name] {
				gpu.ThrottleReasons = append(gpu.ThrottleReasons, name)
			}
		}
		break
	}

	gpu.EccEnabled = g.yesNo("ecc_mode/current_ecc")
	// Ampere and newer report SRAM and DRAM errors, older GPUs single and double bit errors
	gpu.EccCorrectedVolatile = g.sum("ecc_errors/volatile/sram_correctable", "ecc_errors/volatile/dram_correctable",
		"ecc_errors/volatile/single_bit/total")
	gpu.EccUncorrectedVolatile = g.sum("ecc_errors/volatile/sram_uncorrectable", "ecc_errors/volatile/sram_uncorrectable_parity",
		"ecc_errors/volatile/sram_uncorrectable_secded", "ecc_errors/volatile/dram_uncorrectable", "ecc_errors/volatile/double_bit/total")
	gpu.EccCorrectedAggregate = g.sum("ecc_errors/aggregate/sram_correctable", "ecc_errors/aggregate/dram_correctable",
		"ecc_errors/aggregate/single_bit/total")
	gpu.EccUncorrectedAggregate = g.sum("ecc_errors/aggregate/sram_uncorrectable", "ecc_errors/aggregate/sram_uncorrectable_parity",
		"ecc_errors/aggregate/sram_uncorrectable_secded", "ecc_errors/aggregate/dram_uncorrectable", "ecc_errors/aggregate/double_bit/total")

	gpu.RetiredPagesSingleBit = g.sum("retired_pages/multiple_single_bit_retirement/retired_count")
	gpu.RetiredPagesDoubleBit = g.sum("retired_pages/double_bit_retirement/retired_count")
	gpu.RetiredPagesPending = g.yesNo("retired_pages/pending_retirement", "retired_pages/pending_blacklist")

	gpu.RemappedRowsCorrectable = g.sum("remapped_rows/remapped_row_corr")
	gpu.RemappedRowsUncorrectable = g.sum("remapped_rows/remapped_row_unc")
	gpu.RemappedRowsPending = g.yesNo("remapped_rows/remapped_row_pending")
	gpu.RemappedRowsFailure = g.yesNo("remapped_rows/remapped_row_failure")

	channelRepair, tpcRepair := g.yesNo("ecc_errors/channel_repair_pending"), g.yesNo("ecc_errors/tpc_repair_pending")
	if channelRepair != nil || tpcRepair != nil {
		pending := (channelRepair != nil && *channelRepair) || (tpcRepair != nil && *tpcRepair)
		gpu.MemoryRepairPending = &pending
	}
	gpu.MemoryUnrepairable = g.yesNo("ecc_errors/unrepairable_memory")

	gpu.PcieLinkGen = g.number("pci/pci_gpu_link_info/pcie_gen/current_link_gen")
	gpu.PcieLinkGenMax = g.number("pci/pci_gpu_link_info/pcie_gen/max_link_gen")
	gpu.PcieLinkWidth = g.number("pci/pci_gpu_link_info/link_widths/current_link_width")
	gpu.PcieLinkWidthMax = g.number("pci/pci_gpu_link_info/link_widths/max_link_width")

	addNvidiaProblemMessages(g, gpu)
	return gpu
}

// parseNvidiaTemperature reads the temperatures and thresholds. Older drivers report absolute thresholds
// (gpu_temp_slow_threshold), newer drivers only the distance to the T.Limit (gpu_temp_tlimit) and thresholds
// relative to it, e.g. gpu_temp_slow_tlimit_threshold = -2 C: the GPU slows down 2 degrees above the T.Limit.
// "The T.Limit sensor measures the current margin in degree Celsius to the maximum operating temperature",
// https://docs.nvidia.com/deploy/nvidia-smi/index.html
func parseNvidiaTemperature(g *nvidiaXMLNode, gpu *resultGpu) {
	gpu.TemperatureCelsius = g.number("temperature/gpu_temp")
	gpu.TemperatureMemoryCelsius = g.number("temperature/memory_temp")
	gpu.TemperatureSlowdownCelsius = g.number("temperature/gpu_temp_slow_threshold")
	gpu.TemperatureShutdownCelsius = g.number("temperature/gpu_temp_max_threshold")

	if tlimit := g.number("temperature/gpu_temp_tlimit"); tlimit != nil {
		headroom := *tlimit
		if slow := g.number("temperature/gpu_temp_slow_tlimit_threshold"); slow != nil {
			headroom -= *slow
		}
		gpu.TemperatureHeadroomCelsius = &headroom
	} else if gpu.TemperatureSlowdownCelsius != nil && gpu.TemperatureCelsius != nil {
		headroom := *gpu.TemperatureSlowdownCelsius - *gpu.TemperatureCelsius
		gpu.TemperatureHeadroomCelsius = &headroom
	}
}

// addNvidiaProblemMessages adds the messages of nvidia-smi like "[GPU requires reset]" which are reported instead of
// the values of a GPU
func addNvidiaProblemMessages(g *nvidiaXMLNode, gpu *resultGpu) {
	var walk func(node *nvidiaXMLNode)
	walk = func(node *nvidiaXMLNode) {
		value := strings.Trim(strings.TrimSpace(node.Text), "[]")
		for _, problem := range nvidiaProblemValues {
			if strings.EqualFold(value, problem) && !slices.Contains(gpu.Messages, problem) {
				gpu.Messages = append(gpu.Messages, problem)
			}
		}
		for _, child := range node.Nodes {
			walk(child)
		}
	}
	walk(g)
}
