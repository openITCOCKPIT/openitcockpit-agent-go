#!/bin/sh
# Collects the information of Intel GPUs (drivers i915 and xe) from sysfs, read only.
# Usage: sh intel-sysfs-dump.sh > intel-gpu-dump.txt   (as root, some files are only readable by root)
# Run it once when the GPU is idle and once under load (e.g. a video or glmark2).
# Serial numbers and unique ids are not dumped.

echo "# intel-sysfs-dump 4 $(date -Iseconds)"
echo "kernel: $(uname -r)"
grep PRETTY_NAME /etc/os-release 2>/dev/null
command -v lspci >/dev/null && lspci -nn | grep -Ei 'vga|display|3d'
grep -E '^(MemTotal|MemAvailable):' /proc/meminfo

dump_files() {
  # $1 = directory, $2... = file names
  dir=$1
  shift
  for file in "$@"; do
    if [ -r "$dir/$file" ] && [ -f "$dir/$file" ]; then
      echo "$file: $(tr '\n' ' ' < "$dir/$file" 2>/dev/null)"
    fi
  done
}

dump_dir() {
  # All readable regular files of a directory (not recursive)
  for file in "$1"/*; do
    [ -f "$file" ] && [ -r "$file" ] || continue
    case "$(basename "$file")" in *serial*|*unique_id*|*uevent*|*config*|*rom*|*resource*) continue ;; esac
    echo "$(basename "$file"): $(head -c 400 "$file" 2>/dev/null | tr '\n' ' ')"
  done
}

for card in /sys/class/drm/card*; do
  case "$card" in *-*) continue ;; esac  # card0-DP-1 etc. are connectors
  device="$card/device"
  [ -e "$device/vendor" ] || continue
  driver=$(basename "$(readlink -f "$device/driver")")
  case "$driver" in i915|xe) ;; *) continue ;; esac
  echo
  echo "== $card driver=$driver pci=$(basename "$(readlink -f "$device")")"
  dump_files "$device" vendor device subsystem_vendor subsystem_device revision class boot_vga \
    current_link_speed current_link_width max_link_speed max_link_width
  # Upstream PCI bridges (Arc GPUs report the internal link of the GPU, 2.5 GT/s x1)
  parent=$(dirname "$(readlink -f "$device")")
  for level in 1 2 3; do
    [ -r "$parent/vendor" ] || break
    echo "-- upstream $level $(basename "$parent")"
    dump_files "$parent" vendor device class current_link_speed current_link_width max_link_speed max_link_width
    parent=$(dirname "$parent")
  done
  # i915: card level frequencies and per gt power management (rc6, throttle reasons)
  dump_files "$card" gt_act_freq_mhz gt_cur_freq_mhz gt_min_freq_mhz gt_max_freq_mhz gt_boost_freq_mhz \
    gt_RP0_freq_mhz gt_RP1_freq_mhz gt_RPn_freq_mhz power/rc6_enable power/rc6_residency_ms
  for gt in "$card"/gt/gt*; do
    [ -d "$gt" ] || continue
    echo "-- $gt"
    dump_dir "$gt"
  done
  # xe: per tile and gt frequencies, throttle reasons and idle residency
  for tile in "$device"/tile*; do
    [ -d "$tile" ] || continue
    echo "-- $tile"
    dump_dir "$tile"
    for gt in "$tile"/gt*; do
      [ -d "$gt" ] || continue
      for dir in "$gt" "$gt"/freq0 "$gt"/freq0/throttle "$gt"/gtidle; do
        [ -d "$dir" ] || continue
        echo "-- $dir"
        dump_dir "$dir"
      done
    done
  done
  for hwmon in "$device"/hwmon/hwmon*; do
    [ -d "$hwmon" ] || continue
    echo "-- $hwmon name=$(cat "$hwmon/name" 2>/dev/null)"
    for file in "$hwmon"/*; do
      case "$(basename "$file")" in
        *_input|*_label|*_crit|*_emergency|*_max|*_rated_max|*_cap|*_max_interval) [ -r "$file" ] && echo "$(basename "$file"): $(cat "$file" 2>/dev/null)" ;;
      esac
    done
  done
  # Second sample after 2 seconds: utilization from the idle residency, power from the energy counter
  sleep 2
  echo "-- after 2 seconds"
  for f in "$card"/gt/gt*/rc6_residency_ms "$card"/power/rc6_residency_ms "$device"/tile*/gt*/gtidle/idle_residency_ms \
    "$device"/hwmon/hwmon*/energy*_input "$card"/gt_act_freq_mhz "$device"/tile*/gt*/freq0/act_freq; do
    [ -r "$f" ] && echo "$f: $(cat "$f")"
  done
done

# Engine busy time and memory of the DRM clients (only the drm-* lines, no process names), twice with 5 seconds
# between the samples to calculate the utilization. Must run as root to read the clients of all users.
dump_clients() {
  echo "time: $(date +%s.%N)"
  for fdinfo in /proc/[0-9]*/fdinfo/*; do
    grep -q '^drm-driver:[[:space:]]*\(i915\|xe\)' "$fdinfo" 2>/dev/null || continue
    grep -E '^drm-(driver|pdev|client-id|engine-|cycles-|total-cycles-|total-|resident-)' "$fdinfo" 2>/dev/null
    echo "--"
  done | awk '/^drm-client-id/ { if (seen[$2]++) skip=1; else skip=0 } !skip'
  echo "clients: $(grep -l '^drm-driver:[[:space:]]*\(i915\|xe\)' /proc/[0-9]*/fdinfo/* 2>/dev/null | wc -l) files"
}
echo
echo "== drm clients (fdinfo), uid $(id -u)"
dump_clients
sleep 5
echo
echo "== drm clients after 5 seconds (fdinfo)"
dump_clients
for f in /sys/class/drm/card*/gt_act_freq_mhz /sys/class/drm/card*/device/tile*/gt*/freq0/act_freq; do
  [ -r "$f" ] && echo "$f: $(cat "$f")"
done
