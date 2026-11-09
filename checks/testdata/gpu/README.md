# Test data of the GPU check

The GPU check (`checks/gpu*.go`) is tested against dumps of real systems, so it can be developed without access to
the hardware. Before a dump is added, it must be checked for personal data: process names, user names and paths,
host names, serial numbers and the PDI of NVIDIA GPUs must be removed.

| Directory | Created with                                                                                                   | Used by                            |
|-----------|----------------------------------------------------------------------------------------------------------------|------------------------------------|
| `nvidia/` | `nvidia-smi -q -x`, part of the dumps of [oitc-nvidia-dump](https://github.com/openITCOCKPIT/oitc-nvidia-dump) | `gpu_test.go`, `gpu_linux_test.go` |
| `amd/`    | [amd-sysfs-dump.sh](amd-sysfs-dump.sh) (driver amdgpu, Linux)                                                  | `gpu_amd_linux_test.go`            |
| `intel/`  | [intel-sysfs-dump.sh](intel-sysfs-dump.sh) (drivers i915 and xe, Linux)                                        | `gpu_intel_linux_test.go`          |


The tests of AMD and Intel build a fake `/sys/class/drm` (and `/proc` for the DRM clients) from the dumps, see
`buildAmdSysfs()` and `buildIntelSysfs()`. Windows has no dumps: `gpu_windows_test.go` uses the values of a live
system and reads the GPUs of the system the test runs on.

## Creating dumps

The scripts only read files of sysfs and procfs, they do not change anything. They should run as root, as some files
(e.g. `gpu_metrics` and the `fdinfo` of the processes of other users) are only readable by root. Run them once when
the GPU is idle and once under load (e.g. `glmark2`, a video or a game):

```
sh amd-sysfs-dump.sh > amd-gpu-dump-idle.txt
sh intel-sysfs-dump.sh > intel-gpu-dump-idle.txt
```

`intel-sysfs-dump.sh` takes about 7 seconds, as it samples the counters twice. Since version 3 the first line
contains the version of the script, the tests support all versions of the dumps in `intel/`. Serial numbers and
unique ids are not dumped, the DRM clients are dumped without process names.

## Not implemented

Ideas which were left out

- **Total memory of Intel dGPUs (Arc)**: the used memory is the sum of the memory of the DRM clients, the total
  memory is only available with an ioctl of the DRM device (`DRM_I915_QUERY_MEMORY_REGIONS`,
  `DRM_XE_DEVICE_QUERY_MEM_REGIONS`), so `memory_total_bytes` is unknown. Intel GPUs are rare in data centers.
- **Power of AMD GPUs on Windows**: the AMD driver reports the power 0 to D3DKMT, also under full load. The power is
  only available with the AMD Display Library (ADL, `atiadlxx.dll`). On Linux the power is read from hwmon.
