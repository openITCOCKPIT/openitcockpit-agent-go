package checks

import (
	"regexp"
	"strings"
)

// Messages of the amdgpu driver in the kernel log about GPU hangs and resets, like the TDR events on Windows. The
// messages are written by:
// - Job timeouts and ring resets: amdgpu_job_timedout() in
//   https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/drivers/gpu/drm/amd/amdgpu/amdgpu_job.c
// - GPU resets: amdgpu_device_gpu_recover() and amdgpu_do_asic_reset() in
//   https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/drivers/gpu/drm/amd/amdgpu/amdgpu_device.c
// - Fatal RAS errors of data center GPUs (Instinct):
//   https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/drivers/gpu/drm/amd/amdgpu/amdgpu_ras.c
// Not reported:
// - The other RAS errors, they are counted in sysfs (ecc_errors of the GPU)
// - Page faults ("[gfxhub] page fault (src_id:..."), they are usually caused by bugs of applications, like the Xid
//   errors 13 and 31 of NVIDIA, and are written with many lines
// - The process information after a timeout ("Process information: process ... pid ..."), it contains the names of
//   processes

const (
	gpuDriverEventError   = 1
	gpuDriverEventWarning = 2
)

var (
	// Prefixes of the messages:
	// - dev_err(), dev_info() and RAS_EVENT_LOG(): "amdgpu 0000:03:00.0: ", RAS_EVENT_LOG() also adds an event id "{12}"
	// - drm_err() (via dev_err()): "amdgpu 0000:03:00.0: [drm] *ERROR* "
	// - DRM_ERROR() of older kernels (e.g. the ring timeout up to 6.8): "[drm:amdgpu_job_timedout [amdgpu]] *ERROR* "
	// - DRM_INFO() and DRM_WARN() of older kernels (e.g. "VRAM is lost" up to 6.8): "[drm] "
	// Older kernels (e.g. 5.15 to 6.12) add "amdgpu: " after the device (dev_fmt of amdgpu.h).
	amdgpuMessagePattern = regexp.MustCompile(`^(?:amdgpu ([0-9a-fA-F]+:[0-9a-fA-F]+:[0-9a-fA-F]+\.[0-9a-fA-F]): |\[drm:\S+ \[amdgpu\]\] |\[drm\] )(?:amdgpu: )?(?:\[drm\] )?(?:\*ERROR\* )?(?:\{\d+\})?(.+)$`)

	amdgpuEventPatterns = []amdgpuEventPattern{
		// A job of a ring (e.g. gfx_0.0.0) did not finish in time, the GPU hangs
		{regexp.MustCompile(`^ring \S+ timeout, but soft recovered`), gpuDriverEventWarning},
		{regexp.MustCompile(`^ring \S+ timeout, signaled seq=`), gpuDriverEventError},
		// Reset of a single ring after a timeout (newer kernels and GPUs)
		{regexp.MustCompile(`^Ring \S+ reset succeeded`), gpuDriverEventWarning},
		{regexp.MustCompile(`^Ring \S+ reset failed`), gpuDriverEventError},
		// Reset of the whole GPU
		{regexp.MustCompile(`^GPU reset\(\d+\) succeeded`), gpuDriverEventWarning},
		{regexp.MustCompile(`^GPU reset\(\d+\) failed`), gpuDriverEventError},
		{regexp.MustCompile(`^GPU reset retry is beyond the retry limit`), gpuDriverEventError},
		{regexp.MustCompile(`^GPU (?:mode1|link|pre asic) reset failed`), gpuDriverEventError},
		{regexp.MustCompile(`^ASIC reset failed with error`), gpuDriverEventError},
		{regexp.MustCompile(`^VRAM is lost due to GPU reset`), gpuDriverEventError},
		// Fatal RAS error, the GPU is reset
		{regexp.MustCompile(`^uncorrectable hardware error\(ERREVENT_ATHUB_INTERRUPT\) detected`), gpuDriverEventError},
	}
)

type amdgpuEventPattern struct {
	pattern   *regexp.Regexp
	entryType int64
}

// parseAmdgpuMessage returns the event of a message of the amdgpu driver in the kernel log or nil if the message
// is not about a GPU hang or reset
func parseAmdgpuMessage(message string) *resultGpuDriverEvent {
	if !strings.HasPrefix(message, "amdgpu ") && !strings.HasPrefix(message, "[drm") {
		return nil
	}
	match := amdgpuMessagePattern.FindStringSubmatch(message)
	if match == nil {
		return nil
	}
	text := strings.TrimSpace(match[2])
	for _, event := range amdgpuEventPatterns {
		if !event.pattern.MatchString(text) {
			continue
		}
		pciBusID := ""
		if match[1] != "" {
			pciBusID = normalizePciBusID(match[1])
		}
		return &resultGpuDriverEvent{
			Source:    "amdgpu",
			EntryType: event.entryType,
			PciBusID:  pciBusID,
			Message:   text,
		}
	}
	return nil
}
