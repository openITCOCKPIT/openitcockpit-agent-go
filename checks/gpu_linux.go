package checks

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

var (
	gpuPciDevicesPath = "/sys/bus/pci/devices"
	gpuKmsgPath       = "/dev/kmsg"
)

// countNvidiaPciDevices returns the number of NVIDIA GPUs (PCI class display controller) on the PCI bus, -1 if
// unknown. NVSwitches have another PCI class (bridge) and are not counted. The files vendor and class are described
// in https://www.kernel.org/doc/Documentation/ABI/testing/sysfs-bus-pci, the class codes in https://pci-ids.ucw.cz/read/PD/03
func countNvidiaPciDevices() int {
	devices, err := os.ReadDir(gpuPciDevicesPath)
	if err != nil {
		return -1
	}
	count := 0
	for _, device := range devices {
		vendor, err := os.ReadFile(filepath.Join(gpuPciDevicesPath, device.Name(), "vendor"))
		if err != nil || strings.TrimSpace(string(vendor)) != "0x10de" {
			continue
		}
		class, err := os.ReadFile(filepath.Join(gpuPciDevicesPath, device.Name(), "class"))
		if err != nil {
			continue
		}
		// 0x0300 VGA compatible controller, 0x0302 3D controller (data center GPUs)
		if c := strings.TrimSpace(string(class)); strings.HasPrefix(c, "0x0300") || strings.HasPrefix(c, "0x0302") {
			count++
		}
	}
	return count
}

// Xid errors of the NVIDIA driver in the kernel log. The meaning and the recommended action of every Xid code are
// listed in the Xid catalog of NVIDIA: https://docs.nvidia.com/deploy/xid-errors/analyzing-xid-catalog.html
// e.g.: NVRM: Xid (PCI:0000:3b:00): 79, pid=1234, name=python3, GPU has fallen off the bus.
// or:   nvidia-nvswitch3: SXid (PCI:0000:05:00.0): 12028, Non-fatal, Link 32 egress non-posted PRIV error (First)
var (
	xidPattern        = regexp.MustCompile(`\b(S?Xid) \(PCI:([0-9a-fA-F]+:[0-9a-fA-F]+:[0-9a-fA-F]+(?:\.[0-9a-fA-F])?)\): (\d+),\s*(.*)`)
	xidProcessPattern = regexp.MustCompile(`^pid=[^,]*, name=[^,]*,\s*`)
)

// gpuKernelLog collects the messages of the GPU drivers in the kernel log: the Xid errors of the NVIDIA driver and the
// GPU hangs and resets of the amdgpu driver
type gpuKernelLog struct {
	since  int64 // Unix timestamp, older messages are ignored
	xids   []*resultGpuXid
	events []*resultGpuDriverEvent // Oldest first
}

// readGpuDriverLogs reads the kernel log (/dev/kmsg) once for the Xid errors and the events of the amdgpu driver since
// the given time. The kernel log is a ring buffer, older messages may already be overwritten on systems with many
// kernel messages. If the kernel log can not be read, the error is only reported in the Xid state.
func readGpuDriverLogs(ctx context.Context, since time.Time) *gpuDriverLogs {
	log := &gpuKernelLog{since: since.Unix(), xids: []*resultGpuXid{}, events: []*resultGpuDriverEvent{}}
	if err := log.read(); err != nil {
		return &gpuDriverLogs{xidError: err}
	}
	return &gpuDriverLogs{xids: log.xids, events: newGpuDriverEventList(log.events)}
}

// read reads all records of /dev/kmsg. /dev/kmsg returns one record per read, -EAGAIN if no more records are
// available (O_NONBLOCK) and -EPIPE if a record was overwritten: https://www.kernel.org/doc/Documentation/ABI/testing/dev-kmsg
func (l *gpuKernelLog) read() error {
	fd, err := unix.Open(gpuKmsgPath, unix.O_RDONLY|unix.O_NONBLOCK|unix.O_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("could not read the kernel log %s: %w", gpuKmsgPath, err)
	}
	defer unix.Close(fd)

	var monotonic unix.Timespec
	if err := unix.ClockGettime(unix.CLOCK_MONOTONIC, &monotonic); err != nil {
		return fmt.Errorf("could not determine the boot time: %w", err)
	}
	bootTime := time.Now().Add(-time.Duration(monotonic.Nano()))

	buf := make([]byte, 16*1024)
	for {
		n, err := unix.Read(fd, buf)
		if err != nil {
			switch {
			case errors.Is(err, unix.EAGAIN):
				// No more messages
				return nil
			case errors.Is(err, unix.EPIPE), errors.Is(err, unix.EINTR):
				// The message was overwritten while reading
				continue
			}
			return fmt.Errorf("could not read the kernel log %s: %w", gpuKmsgPath, err)
		}
		if n == 0 {
			return nil
		}
		l.addRecord(string(buf[:n]), bootTime)
	}
}

// addRecord adds a record of /dev/kmsg if it is an Xid error or an event of the amdgpu driver
func (l *gpuKernelLog) addRecord(record string, bootTime time.Time) {
	timestamp, message, ok := parseKmsgRecord(record, bootTime)
	if !ok || timestamp < l.since {
		return
	}
	if xid := parseXidMessage(message); xid != nil {
		xid.Timestamp = timestamp
		l.xids = append(l.xids, xid)
		return
	}
	if event := parseAmdgpuMessage(message); event != nil {
		event.Timestamp = timestamp
		l.events = append(l.events, event)
	}
}

// parseKmsgRecord parses a record of /dev/kmsg: "priority,sequence,microseconds since boot,flags;message" and
// returns the Unix timestamp and the message
func parseKmsgRecord(record string, bootTime time.Time) (int64, string, bool) {
	header, message, found := strings.Cut(record, ";")
	if !found {
		return 0, "", false
	}
	// Continuation lines (dictionary) start with a space
	message, _, _ = strings.Cut(message, "\n")

	fields := strings.Split(header, ",")
	if len(fields) < 3 {
		return 0, "", false
	}
	microseconds, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return 0, "", false
	}
	return bootTime.Add(time.Duration(microseconds) * time.Microsecond).Unix(), message, true
}

// newGpuDriverEventList returns the number of events and the newest events
func newGpuDriverEventList(events []*resultGpuDriverEvent) *resultGpuDriverEventList {
	result := &resultGpuDriverEventList{Count: len(events), Events: []*resultGpuDriverEvent{}}
	for i := len(events) - 1; i >= 0 && len(result.Events) < gpuMaxXidErrors; i-- {
		result.Events = append(result.Events, events[i])
	}
	return result
}

// parseXidMessage parses the Xid and SXid messages of the NVIDIA driver. The name of the process is removed.
func parseXidMessage(message string) *resultGpuXid {
	if !strings.Contains(message, "Xid") {
		return nil
	}
	match := xidPattern.FindStringSubmatch(message)
	if match == nil {
		return nil
	}
	code, err := strconv.ParseInt(match[3], 10, 64)
	if err != nil {
		return nil
	}
	return &resultGpuXid{
		Type:     match[1],
		Code:     code,
		PciBusID: normalizePciBusID(match[2]),
		Message:  strings.TrimSpace(xidProcessPattern.ReplaceAllString(match[4], "")),
	}
}

// hasPlatformGpus returns true if a GPU of the drivers amdgpu, i915 or xe exists
func hasPlatformGpus() bool {
	return hasAmdGpus() || hasIntelGpus()
}

// collectPlatformGpus returns the GPUs of the drivers amdgpu, i915 and xe. The utilization of Intel GPUs and the
// number of processes are read from the DRM clients.
func collectPlatformGpus(ctx context.Context, existing []*resultGpu) []*resultGpu {
	usage, hasPrevious := readDrmUsage("amdgpu", "i915", "xe")
	gpus := collectAmdGpus()
	for _, gpu := range gpus {
		gpu.Processes = 0
		if u := usage[gpu.PciBusID]; u != nil {
			gpu.Processes = u.processes
		}
	}
	return append(gpus, collectIntelGpus(usage, hasPrevious)...)
}
