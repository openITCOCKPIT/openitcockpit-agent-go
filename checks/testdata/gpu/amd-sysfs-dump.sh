#!/bin/sh
# Collects the information of AMD GPUs (driver amdgpu) from sysfs, read only.
# Usage: sh amd-sysfs-dump.sh > amd-gpu-dump.txt   (as root, some files are only readable by root)

echo "# amd-sysfs-dump $(date -Iseconds)"
echo "kernel: $(uname -r)"
grep PRETTY_NAME /etc/os-release 2>/dev/null
command -v lspci >/dev/null && lspci -nn | grep -Ei 'vga|display|3d'
grep -E '^(MemTotal|MemAvailable):' /proc/meminfo

for device in /sys/class/drm/card*/device; do
  [ -e "$device/vendor" ] || continue
  case "$device" in *-*) continue ;; esac  # card0-DP-1 etc. are connectors
  echo
  echo "== $device driver=$(basename "$(readlink -f "$device/driver")")"
  for file in vendor device subsystem_vendor subsystem_device revision boot_vga \
    gpu_busy_percent mem_busy_percent \
    mem_info_vram_total mem_info_vram_used mem_info_vis_vram_total mem_info_vis_vram_used \
    mem_info_gtt_total mem_info_gtt_used mem_info_vram_vendor \
    current_link_speed current_link_width max_link_speed max_link_width \
    power_dpm_state power_dpm_force_performance_level \
    pp_dpm_sclk pp_dpm_mclk pp_dpm_fclk pp_dpm_socclk pp_dpm_pcie \
    product_name product_number serial_number unique_id vbios_version \
    ras/features ras/umc_err_count ras/gfx_err_count ras/sdma_err_count ras/mmhub_err_count ras/pcie_bif_err_count; do
    if [ -r "$device/$file" ]; then
      value=$(tr '\n' ' ' < "$device/$file" 2>/dev/null)
      case "$file" in serial_number|unique_id) value="(present, not dumped)" ;; esac
      echo "$file: $value"
    fi
  done
  # Binary metrics of the power management firmware (format depends on the GPU generation)
  if [ -r "$device/gpu_metrics" ]; then
    echo "gpu_metrics (base64): $(base64 -w0 "$device/gpu_metrics" 2>/dev/null)"
  fi
  for hwmon in "$device"/hwmon/hwmon*; do
    [ -d "$hwmon" ] || continue
    echo "-- $hwmon name=$(cat "$hwmon/name" 2>/dev/null)"
    for file in "$hwmon"/*; do
      case "$(basename "$file")" in
        *_input|*_label|*_crit|*_crit_hyst|*_emergency|*_max|*_min|*_cap|*_cap_max|*_cap_min|*_average|*_enable|*_target|pwm1)
          [ -r "$file" ] && echo "$(basename "$file"): $(cat "$file" 2>/dev/null)" ;;
      esac
    done
  done
done
