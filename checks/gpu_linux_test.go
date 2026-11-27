package checks

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParseXidMessage(t *testing.T) {
	tests := []struct {
		message string
		xid     *resultGpuXid
	}{
		{
			"NVRM: Xid (PCI:0000:3b:00): 79, pid=1234, name=python3, GPU has fallen off the bus.",
			&resultGpuXid{Type: "Xid", Code: 79, PciBusID: "0000:3b:00.0", Message: "GPU has fallen off the bus."},
		},
		{
			"NVRM: Xid (PCI:0000:18:00): 119, pid='<unknown>', name=<unknown>, Timeout after 6s of waiting for RPC response from GPU0 GSP!",
			&resultGpuXid{Type: "Xid", Code: 119, PciBusID: "0000:18:00.0", Message: "Timeout after 6s of waiting for RPC response from GPU0 GSP!"},
		},
		{
			"NVRM: Xid (PCI:0000:01:00): 13, Graphics SM Warp Exception on (GPC 0, TPC 0, SM 0): Out Of Range Address",
			&resultGpuXid{Type: "Xid", Code: 13, PciBusID: "0000:01:00.0", Message: "Graphics SM Warp Exception on (GPC 0, TPC 0, SM 0): Out Of Range Address"},
		},
		{
			"nvidia-nvswitch3: SXid (PCI:0000:05:00.0): 12028, Non-fatal, Link 32 egress non-posted PRIV error (First)",
			&resultGpuXid{Type: "SXid", Code: 12028, PciBusID: "0000:05:00.0", Message: "Non-fatal, Link 32 egress non-posted PRIV error (First)"},
		},
		{"NVRM: loading NVIDIA UNIX Open Kernel Module for x86_64  595.91.07", nil},
		{"Xid is mentioned without PCI address", nil},
	}
	for _, test := range tests {
		xid := parseXidMessage(test.message)
		if (xid == nil) != (test.xid == nil) || (xid != nil && *xid != *test.xid) {
			t.Errorf("%q: expected %+v, got %+v", test.message, test.xid, xid)
		}
	}
}

func TestParseKmsgRecord(t *testing.T) {
	bootTime := time.Unix(1792483200, 0)
	record := "3,1234,90000000,-;NVRM: Xid (PCI:0000:3b:00): 48, pid=1234, name=python3, An uncorrectable double bit error (DBE) has been detected on GPU in the framebuffer at partition 6, subpartition 0.\n SUBSYSTEM=pci\n DEVICE=+pci:0000:3b:00.0\n"
	timestamp, message, ok := parseKmsgRecord(record, bootTime)
	if !ok || timestamp != 1792483290 ||
		message != "NVRM: Xid (PCI:0000:3b:00): 48, pid=1234, name=python3, An uncorrectable double bit error (DBE) has been detected on GPU in the framebuffer at partition 6, subpartition 0." {
		t.Errorf("unexpected record: %d %q", timestamp, message)
	}
	for _, record := range []string{"invalid", "3,1,x,-;NVRM: Xid (PCI:0000:3b:00): 48, foo", "3;foo"} {
		if _, _, ok := parseKmsgRecord(record, bootTime); ok {
			t.Errorf("%q: expected an invalid record", record)
		}
	}
}

func TestParseAmdgpuMessage(t *testing.T) {
	tests := []struct {
		message string
		event   *resultGpuDriverEvent
	}{
		// Kernel 7.0
		{
			"amdgpu 0000:03:00.0: ring gfx_0.0.0 timeout, signaled seq=1234, emitted seq=1236",
			&resultGpuDriverEvent{EntryType: 1, PciBusID: "0000:03:00.0", Message: "ring gfx_0.0.0 timeout, signaled seq=1234, emitted seq=1236"},
		},
		{
			"amdgpu 0000:03:00.0: ring sdma0 timeout, but soft recovered",
			&resultGpuDriverEvent{EntryType: 2, PciBusID: "0000:03:00.0", Message: "ring sdma0 timeout, but soft recovered"},
		},
		{
			"amdgpu 0000:03:00.0: Ring gfx_0.0.0 reset succeeded",
			&resultGpuDriverEvent{EntryType: 2, PciBusID: "0000:03:00.0", Message: "Ring gfx_0.0.0 reset succeeded"},
		},
		{
			"amdgpu 0000:03:00.0: Ring gfx_0.0.0 reset failed",
			&resultGpuDriverEvent{EntryType: 1, PciBusID: "0000:03:00.0", Message: "Ring gfx_0.0.0 reset failed"},
		},
		{
			"amdgpu 0000:03:00.0: GPU reset(2) failed with error -110",
			&resultGpuDriverEvent{EntryType: 1, PciBusID: "0000:03:00.0", Message: "GPU reset(2) failed with error -110"},
		},
		{
			"amdgpu 0000:03:00.0: VRAM is lost due to GPU reset!",
			&resultGpuDriverEvent{EntryType: 1, PciBusID: "0000:03:00.0", Message: "VRAM is lost due to GPU reset!"},
		},
		{
			"amdgpu 0000:c1:00.0: {3}uncorrectable hardware error(ERREVENT_ATHUB_INTERRUPT) detected!",
			&resultGpuDriverEvent{EntryType: 1, PciBusID: "0000:c1:00.0", Message: "uncorrectable hardware error(ERREVENT_ATHUB_INTERRUPT) detected!"},
		},
		{
			"amdgpu 0000:03:00.0: [drm] *ERROR* GPU mode1 reset failed",
			&resultGpuDriverEvent{EntryType: 1, PciBusID: "0000:03:00.0", Message: "GPU mode1 reset failed"},
		},
		// Kernel 5.15 to 6.12: "amdgpu: " after the device
		{
			"amdgpu 0000:0b:00.0: amdgpu: GPU reset(1) succeeded!",
			&resultGpuDriverEvent{EntryType: 2, PciBusID: "0000:0b:00.0", Message: "GPU reset(1) succeeded!"},
		},
		{
			"amdgpu 0000:0b:00.0: amdgpu: ASIC reset failed with error, -110 for drm dev, 0000:0b:00.0",
			&resultGpuDriverEvent{EntryType: 1, PciBusID: "0000:0b:00.0", Message: "ASIC reset failed with error, -110 for drm dev, 0000:0b:00.0"},
		},
		// Kernel 5.15 to 6.8: DRM_ERROR and DRM_INFO without the PCI address
		{
			"[drm:amdgpu_job_timedout [amdgpu]] *ERROR* ring gfx timeout, signaled seq=110381, emitted seq=110383",
			&resultGpuDriverEvent{EntryType: 1, Message: "ring gfx timeout, signaled seq=110381, emitted seq=110383"},
		},
		{
			"[drm] VRAM is lost due to GPU reset!",
			&resultGpuDriverEvent{EntryType: 1, Message: "VRAM is lost due to GPU reset!"},
		},
		// Not reported
		{"[drm:amdgpu_job_timedout [amdgpu]] *ERROR* Process information: process glmark2 pid 1234 thread glmark2:cs0 pid 1240", nil},
		{"amdgpu 0000:03:00.0:  Process glmark2 pid 1234 thread glmark2:cs0 pid 1240", nil},
		{"amdgpu 0000:03:00.0: amdgpu: [gfxhub] page fault (src_id:0 ring:24 vmid:3 pasid:32769)", nil},
		{"amdgpu 0000:03:00.0: amdgpu: Starting gfx_0.0.0 ring reset", nil},
		{"amdgpu 0000:03:00.0: amdgpu: SMU is initialized successfully!", nil},
		{"amdgpu 0000:c1:00.0: {2}5 new correctable hardware errors detected in umc block", nil},
		{"i915 0000:00:02.0: [drm] *ERROR* GPU HANG: ecode 12:1:85dffffb, in glmark2 [1234]", nil},
		{"[drm] Initialized amdgpu 3.64.0 for 0000:03:00.0 on minor 1", nil},
	}
	for _, test := range tests {
		event := parseAmdgpuMessage(test.message)
		if test.event != nil {
			test.event.Source = "amdgpu"
		}
		if (event == nil) != (test.event == nil) || (event != nil && *event != *test.event) {
			t.Errorf("%q: expected %+v, got %+v", test.message, test.event, event)
		}
	}
}

func TestGpuKernelLog(t *testing.T) {
	bootTime := time.Unix(1792483200, 0)
	log := &gpuKernelLog{since: 1792483250, xids: []*resultGpuXid{}, events: []*resultGpuDriverEvent{}}
	records := []string{
		// Older than the time window
		"3,1,10000000,-;NVRM: Xid (PCI:0000:3b:00): 79, pid=1234, name=python3, GPU has fallen off the bus.",
		"3,2,20000000,-;amdgpu 0000:03:00.0: ring gfx_0.0.0 timeout, signaled seq=1, emitted seq=3",
		"3,3,60000000,-;NVRM: Xid (PCI:0000:3b:00): 79, pid=1234, name=python3, GPU has fallen off the bus.",
		"3,4,70000000,-;amdgpu 0000:03:00.0: ring gfx_0.0.0 timeout, signaled seq=1, emitted seq=3\n SUBSYSTEM=pci\n DEVICE=+pci:0000:03:00.0\n",
		"3,5,70000001,-;amdgpu 0000:03:00.0:  Process glmark2 pid 1234 thread glmark2:cs0 pid 1240",
		"6,6,80000000,-;amdgpu 0000:03:00.0: GPU reset(1) succeeded!",
		"6,7,90000000,-;EXT4-fs (sda1): mounted filesystem",
	}
	for _, record := range records {
		log.addRecord(record, bootTime)
	}
	if len(log.xids) != 1 || log.xids[0].Code != 79 || log.xids[0].Timestamp != 1792483260 {
		t.Errorf("unexpected Xid errors: %+v", log.xids)
	}
	events := newGpuDriverEventList(log.events)
	if events.Count != 2 || len(events.Events) != 2 || events.Error != "" {
		t.Fatalf("unexpected events: %+v", events)
	}
	// Newest first
	if events.Events[0].Message != "GPU reset(1) succeeded!" || events.Events[0].Timestamp != 1792483280 ||
		events.Events[1].EntryType != 1 || events.Events[1].Timestamp != 1792483270 {
		t.Errorf("unexpected events: %+v %+v", events.Events[0], events.Events[1])
	}
}

func TestNewGpuDriverEventList(t *testing.T) {
	events := []*resultGpuDriverEvent{}
	for i := 0; i < 15; i++ {
		events = append(events, &resultGpuDriverEvent{Timestamp: int64(i)})
	}
	list := newGpuDriverEventList(events)
	if list.Count != 15 || len(list.Events) != gpuMaxXidErrors || list.Events[0].Timestamp != 14 || list.Events[9].Timestamp != 5 {
		t.Errorf("unexpected list: %d %d", list.Count, len(list.Events))
	}
	if list := newGpuDriverEventList([]*resultGpuDriverEvent{}); list.Count != 0 || list.Events == nil {
		t.Errorf("expected an empty list: %+v", list)
	}
}

func TestReadGpuDriverLogs(t *testing.T) {
	// Reads the real kernel log, the system has no GPU with errors
	logs := readGpuDriverLogs(context.Background(), time.Now().Add(-gpuXidWindow))
	if logs.xidError != nil {
		if errors.Is(logs.xidError, os.ErrPermission) || errors.Is(logs.xidError, os.ErrNotExist) {
			t.Skipf("the kernel log is not readable: %s", logs.xidError)
		}
		t.Fatal(logs.xidError)
	}
	if logs.xids == nil || logs.events == nil || logs.events.Events == nil {
		t.Errorf("expected empty lists: %+v", logs)
	}
}

func TestCountNvidiaPciDevices(t *testing.T) {
	dir := t.TempDir()
	devices := map[string][2]string{
		"0000:18:00.0": {"0x10de", "0x030200"}, // H200
		"0000:2d:00.0": {"0x10de", "0x030000"}, // GeForce
		"0000:2d:00.1": {"0x10de", "0x040300"}, // HDMI audio of the GeForce
		"0000:05:00.0": {"0x10de", "0x068000"}, // NVSwitch
		"0000:00:02.0": {"0x8086", "0x030000"}, // Intel iGPU
	}
	for name, values := range devices {
		if err := os.MkdirAll(filepath.Join(dir, name), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name, "vendor"), []byte(values[0]+"\n"), 0600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name, "class"), []byte(values[1]+"\n"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	defer func(path string) { gpuPciDevicesPath = path }(gpuPciDevicesPath)

	gpuPciDevicesPath = dir
	if count := countNvidiaPciDevices(); count != 2 {
		t.Errorf("expected 2 GPUs, got %d", count)
	}
	gpuPciDevicesPath = filepath.Join(dir, "missing")
	if count := countNvidiaPciDevices(); count != -1 {
		t.Errorf("expected -1, got %d", count)
	}
}

func TestCheckGpuRun(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "gpu", "nvidia", "h200-nvl-4x-linux-595.91.07.xml"))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "query.xml"), data, 0600); err != nil {
		t.Fatal(err)
	}
	nvidiaSmi := filepath.Join(dir, "nvidia-smi")
	writeScript := func(script string) {
		if err := os.WriteFile(nvidiaSmi, []byte(script), 0700); err != nil {
			t.Fatal(err)
		}
	}
	defer func(path string) { gpuKmsgPath = path }(gpuKmsgPath)
	gpuKmsgPath = filepath.Join(dir, "kmsg")

	writeScript("#!/bin/sh\n[ \"$*\" = \"-q -x\" ] || exit 1\ncat " + filepath.Join(dir, "query.xml") + "\n")
	r, err := (&CheckGpu{nvidiaSmi: nvidiaSmi}).Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	result := r.(*resultGpus)
	if len(result.Gpus) != 4 || result.Nvidia == nil || result.Nvidia.Error != "" || result.Nvidia.AttachedGpus != 4 {
		t.Fatalf("unexpected result: %+v %+v", result, result.Nvidia)
	}
	for i, gpu := range result.Gpus {
		if gpu.Index != i {
			t.Errorf("unexpected index of GPU %d: %d", i, gpu.Index)
		}
	}
	// The kernel log is not readable
	if result.Xid == nil || !strings.Contains(result.Xid.Error, "could not read the kernel log") || result.Xid.Unassigned == nil {
		t.Errorf("expected an error for the kernel log: %+v", result.Xid)
	}
	// The error is only reported once
	if result.DriverEvents != nil {
		t.Errorf("unexpected driver events: %+v", result.DriverEvents)
	}

	// The driver is not loaded
	writeScript("#!/bin/sh\necho \"NVIDIA-SMI has failed because it couldn't communicate with the NVIDIA driver. Make sure that the latest NVIDIA driver is installed and running.\"\nexit 9\n")
	r, err = (&CheckGpu{nvidiaSmi: nvidiaSmi}).Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	result = r.(*resultGpus)
	if len(result.Gpus) != 0 || result.Nvidia == nil || !strings.Contains(result.Nvidia.Error, "couldn't communicate with the NVIDIA driver") {
		t.Errorf("expected an error: %+v %+v", result, result.Nvidia)
	}

	// One GPU has fallen off the bus, nvidia-smi still reports the other GPUs
	writeScript("#!/bin/sh\ncat " + filepath.Join(dir, "query.xml") + "\necho 'Unable to determine the device handle for GPU0000:AF:00.0: Unknown Error' >&2\nexit 15\n")
	r, err = (&CheckGpu{nvidiaSmi: nvidiaSmi}).Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	result = r.(*resultGpus)
	if len(result.Gpus) != 4 || !strings.Contains(result.Nvidia.Error, "Unable to determine the device handle for GPU0000:AF:00.0") {
		t.Errorf("expected 4 GPUs and an error: %d %+v", len(result.Gpus), result.Nvidia)
	}
}
