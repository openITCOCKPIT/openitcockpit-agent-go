package packagemanager

import (
	"testing"
)

// Sample output of the command
// LANG=c zypper --no-refresh --quiet --xmlout list-updates
var zypperAvailableUpdatesOutput = `<?xml version='1.0'?>
<stream>
<update-status version="0.6">
<update-list>
<update kind="package" name="kernel-default" edition="6.12.0-160000.8.1" arch="aarch64" edition-old="6.12.0-160000.7.1"><summary>The Standard Kernel</summary><description>The standard kernel for both uniprocessor and multiprocessor systems.


Source Timestamp: 2025-12-11 09:18:13 +0000
GIT Revision: 5d31a95c3fc60422efb739def20feefbf5d4d151
GIT Branch: SL-16.0</description><license/><source url="https://download.opensuse.org/distribution/leap/16.0/repo/oss/aarch64" alias="https-download.opensuse.org-4bc031d9"/></update><update kind="package" name="kernel-default-extra" edition="6.12.0-160000.8.1" arch="aarch64" edition-old="6.12.0-160000.7.1"><summary>The Standard Kernel - Unsupported kernel modules</summary><description>The standard kernel for both uniprocessor and multiprocessor systems.

This package contains additional modules not supported by SUSE.


Source Timestamp: 2025-12-11 09:18:13 +0000
GIT Revision: 5d31a95c3fc60422efb739def20feefbf5d4d151
GIT Branch: SL-16.0</description><license/><source url="https://download.opensuse.org/distribution/leap/16.0/repo/oss/aarch64" alias="https-download.opensuse.org-4bc031d9"/></update><update kind="package" name="kernel-default-optional" edition="6.12.0-160000.8.1" arch="aarch64" edition-old="6.12.0-160000.7.1"><summary>The Standard Kernel - Optional kernel modules</summary><description>The standard kernel for both uniprocessor and multiprocessor systems.

This package contains optional modules only for openSUSE Leap.


Source Timestamp: 2025-12-11 09:18:13 +0000
GIT Revision: 5d31a95c3fc60422efb739def20feefbf5d4d151
GIT Branch: SL-16.0</description><license/><source url="https://download.opensuse.org/distribution/leap/16.0/repo/oss/aarch64" alias="https-download.opensuse.org-4bc031d9"/></update><update kind="package" name="qemu-guest-agent" edition="10.0.7-160000.1.1" arch="aarch64" edition-old="10.0.4-160000.1.1"><summary>Guest agent for QEMU</summary><description>This package contains the QEMU guest agent. It is installed in the linux guest
to provide information and control at the guest OS level.</description><license/><source url="https://download.opensuse.org/distribution/leap/16.0/repo/oss/aarch64" alias="https-download.opensuse.org-4bc031d9"/></update></update-list>
</update-status>
</stream>`

// Sample output of the command
// LANG=c zypper --no-refresh --quiet --xmlout list-patches --category security
var zypperAvailableSecurityPatchesOutput = `<?xml version='1.0'?>
<stream>
<update-status version="0.6">
<update-list>
<update kind="patch" name="openSUSE-Leap-16.0-112" edition="1" arch="noarch" status="needed" category="security" severity="important" pkgmanager="false" restart="false" interactive="false"><summary>Security update for qemu</summary><description>This update for qemu fixes the following issues:

Update to version 10.0.7.

Security issues fixed:

- CVE-2025-12464: stack-based buffer overflow in the e1000 network device operations can be exploited by a malicious
  guest user to crash the QEMU process on the host (bsc#1253002).
- CVE-2025-11234: use-after-free in WebSocket handshake operations can be exploited by a malicious client with network
  access to the VNC WebSocket port to cause a denial-of-service (bsc#1250984).

Other updates and bugfixes:

- Version 10.0.7:
  * kvm: Fix kvm_vm_ioctl() and kvm_device_ioctl() return value
  * docs/devel: Update URL for make-pullreq script
  * target/arm: Fix assert on BRA.
  * hw/aspeed/{xdma, rtc, sdhci}: Fix endianness to DEVICE_LITTLE_ENDIAN
  * hw/core/machine: Provide a description for aux-ram-share property
  * hw/pci: Make msix_init take a uint32_t for nentries
  * block/io_uring: avoid potentially getting stuck after resubmit at the end of ioq_submit()
  * block-backend: Fix race when resuming queued requests
  * ui/vnc: Fix qemu abort when query vnc info
  * chardev/char-pty: Do not ignore chr_write() failures
  * hw/display/exynos4210_fimd: Account for zero length in fimd_update_memory_section()
  * hw/arm/armv7m: Disable reentrancy guard for v7m_sysreg_ns_ops MRs
  * hw/arm/aspeed: Fix missing SPI IRQ connection causing DMA interrupt failure
  * migration: Fix transition to COLO state from precopy
  * Full backport list: https://lore.kernel.org/qemu-devel/1765037524.347582.2700543.nullmailer@tls.msk.ru/

- Version 10.0.6:
  * linux-user/microblaze: Fix little-endianness binary
  * target/hppa: correct size bit parity for fmpyadd
  * target/i386: user: do not set up a valid LDT on reset
  * async: access bottom half flags with qatomic_read
  * target/i386: fix x86_64 pushw op
  * i386/tcg/smm_helper: Properly apply DR values on SMM entry / exit
  * i386/cpu: Prevent delivering SIPI during SMM in TCG mode
  * i386/kvm: Expose ARCH_CAP_FB_CLEAR when invulnerable to MDS
  * target/i386: Fix CR2 handling for non-canonical addresses
  * block/curl.c: Use explicit long constants in curl_easy_setopt calls
  * pcie_sriov: Fix broken MMIO accesses from SR-IOV VFs
  * target/riscv: rvv: Fix vslide1[up|down].vx unexpected result when XLEN2 and SEWd
  * target/riscv: Fix ssamoswap error handling
  * Full backport list: https://lore.kernel.org/qemu-devel/1761022287.744330.6357.nullmailer@tls.msk.ru/

- Version 10.0.5:
  * tests/functional/test_aarch64_sbsaref_freebsd: Fix the URL of the ISO image
  * tests/functional/test_ppc_bamboo: Replace broken link with working assets
  * physmem: Destroy all CPU AddressSpaces on unrealize
  * memory: New AS helper to serialize destroy+free
  * include/system/memory.h: Clarify address_space_destroy() behaviour
  * migration: Fix state transition in postcopy_start() error handling
  * target/riscv: rvv: Modify minimum VLEN according to enabled vector extensions
  * target/riscv: rvv: Replace checking V by checking Zve32x
  * target/riscv: Fix endianness swap on compressed instructions
  * hw/riscv/riscv-iommu: Fixup PDT Nested Walk
  * Full backport list: https://lore.kernel.org/qemu-devel/1759986125.676506.643525.nullmailer@tls.msk.ru/

- [openSUSE][RPM]: really fix *-virtio-gpu-pci dependency on ARM (bsc#1254286).
- [openSUSE][RPM] spec: make glusterfs support conditional (bsc#1254494).
</description><license/><source url="https://download.opensuse.org/distribution/leap/16.0/repo/oss/aarch64" alias="https-download.opensuse.org-4bc031d9"/><issue-date time_t="1766061723" text="2025-12-18T12:42:03Z"/><issue-list><issue type="cve" id="CVE-2025-11234"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-11234</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-11234</href></issue><issue type="cve" id="CVE-2025-12464"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-12464</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-12464</href></issue><issue type="bugzilla" id="1254286"><title>cockpit - Virtual machines - Import VM - import qcow2 image failed and it fails to start because of missing &apos;qemu-hw-display-virtio-gpu&apos; and &apos;qemu-hw-display-virtio-gpu-pci&apos;</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1254286</href></issue><issue type="bugzilla" id="1230042"><title>[SLM6.2][PPC64LE] Fail to start QEMU - qemu-system-ppc64</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1230042</href></issue><issue type="bugzilla" id="1254494"><title>glusterfs: drop 32-bit architectures</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1254494</href></issue><issue type="bugzilla" id="1250984"><title>VUL-0: CVE-2025-11234: qemu: qemu-kvm: use-after-free in websocket handshake code can lead to denial of service</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1250984</href></issue><issue type="bugzilla" id="1253002"><title>VUL-0: CVE-2025-12464: qemu: net: pad packets to minimum length in qemu_receive_packet()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253002</href></issue></issue-list></update><update kind="patch" name="openSUSE-Leap-16.0-113" edition="1" arch="noarch" status="needed" category="security" severity="important" pkgmanager="false" restart="true" interactive="true"><summary>Security update for the Linux Kernel</summary><description>
The SUSE Linux Enterprise 16.0 kernel was updated to fix various security issues

The following security issues were fixed:

- CVE-2022-50253: bpf: make sure skb-&gt;len != 0 when redirecting to a tunneling device (bsc#1249912).
- CVE-2025-37916: pds_core: remove write-after-free of client_id (bsc#1243474).
- CVE-2025-38084: mm/hugetlb: unshare page tables during VMA split, not before (bsc#1245431 bsc#1245498).
- CVE-2025-38085: mm/hugetlb: fix huge_pmd_unshare() vs GUP-fast race (bsc#1245431 bsc#1245499).
- CVE-2025-38321: smb: Log an error when close_all_cached_dirs fails (bsc#1246328).
- CVE-2025-38728: smb3: fix for slab out of bounds on mount to ksmbd (bsc#1249256).
- CVE-2025-39805: net: macb: fix unregister_netdev call order in macb_remove() (bsc#1249982).
- CVE-2025-39819: fs/smb: Fix inconsistent refcnt update (bsc#1250176).
- CVE-2025-39822: io_uring/kbuf: fix signedness in this_len calculation (bsc#1250034).
- CVE-2025-39831: fbnic: Move phylink resume out of service_task and into open/close (bsc#1249977).
- CVE-2025-39859: ptp: ocp: fix use-after-free bugs causing by ptp_ocp_watchdog (bsc#1250252).
- CVE-2025-39897: net: xilinx: axienet: Add error handling for RX metadata pointer retrieval (bsc#1250746).
- CVE-2025-39917: bpf: Fix out-of-bounds dynptr write in bpf_crypto_crypt (bsc#1250723).
- CVE-2025-39944: octeontx2-pf: Fix use-after-free bugs in otx2_sync_tstamp() (bsc#1251120).
- CVE-2025-39961: iommu/amd/pgtbl: Fix possible race while increase page table level (bsc#1251817).
- CVE-2025-39980: nexthop: Forbid FDB status change while nexthop is in a group (bsc#1252063).
- CVE-2025-39990: bpf: Check the helper function is valid in get_helper_proto (bsc#1252054).
- CVE-2025-40001: scsi: mvsas: Fix use-after-free bugs in mvs_work_queue (bsc#1252303).
- CVE-2025-40003: net: mscc: ocelot: Fix use-after-free caused by cyclic delayed work (bsc#1252301).
- CVE-2025-40006: mm/hugetlb: fix folio is still mapped when deleted (bsc#1252342).
- CVE-2025-40021: tracing: dynevent: Add a missing lockdown check on dynevent (bsc#1252681).
- CVE-2025-40024: vhost: Take a reference on the task in struct vhost_task (bsc#1252686).
- CVE-2025-40027: net/9p: fix double req put in p9_fd_cancelled (bsc#1252763).
- CVE-2025-40031: tee: fix register_shm_helper() (bsc#1252779).
- CVE-2025-40033: remoteproc: pru: Fix potential NULL pointer dereference in pru_rproc_set_ctable() (bsc#1252824).
- CVE-2025-40038: KVM: SVM: Skip fastpath emulation on VM-Exit if next RIP isn&apos;t valid (bsc#1252817).
- CVE-2025-40047: io_uring/waitid: always prune wait queue entry in io_waitid_wait() (bsc#1252790).
- CVE-2025-40053: net: dlink: handle copy_thresh allocation failure (bsc#1252808).
- CVE-2025-40055: ocfs2: fix double free in user_cluster_connect() (bsc#1252821).
- CVE-2025-40059: coresight: Fix incorrect handling for return value of devm_kzalloc (bsc#1252809).
- CVE-2025-40064: smc: Fix use-after-free in __pnet_find_base_ndev() (bsc#1252845).
- CVE-2025-40070: pps: fix warning in pps_register_cdev when register device fail (bsc#1252836).
- CVE-2025-40074: tcp: convert to dev_net_rcu() (bsc#1252794).
- CVE-2025-40075: tcp_metrics: use dst_dev_net_rcu() (bsc#1252795).
- CVE-2025-40081: perf: arm_spe: Prevent overflow in PERF_IDX2OFF() (bsc#1252776).
- CVE-2025-40083: net/sched: sch_qfq: Fix null-deref in agg_dequeue (bsc#1252912).
- CVE-2025-40086: drm/xe: Don&apos;t allow evicting of BOs in same VM in array of VM binds (bsc#1252923).
- CVE-2025-40098: ALSA: hda: cs35l41: Fix NULL pointer dereference in cs35l41_get_acpi_mute_state() (bsc#1252917).
- CVE-2025-40101: btrfs: fix memory leaks when rejecting a non SINGLE data profile without an RST (bsc#1252901).
- CVE-2025-40102: KVM: arm64: Prevent access to vCPU events before init (bsc#1252919).
- CVE-2025-40105: vfs: Don&apos;t leak disconnected dentries on umount (bsc#1252928).
- CVE-2025-40133: mptcp: Call dst_release() in mptcp_active_enable() (bsc#1253328).
- CVE-2025-40134: dm: fix NULL pointer dereference in __dm_suspend() (bsc#1253386).
- CVE-2025-40135: ipv6: use RCU in ip6_xmit() (bsc#1253342).
- CVE-2025-40139: smc: Use __sk_dst_get() and dst_dev_rcu() in in smc_clc_prfx_set() (bsc#1253409).
- CVE-2025-40149: tls: Use __sk_dst_get() and dst_dev_rcu() in get_netdev_for_sock() (bsc#1253355).
- CVE-2025-40153: mm: hugetlb: avoid soft lockup when mprotect to large memory area (bsc#1253408).
- CVE-2025-40157: EDAC/i10nm: Skip DIMM enumeration on a disabled memory controller (bsc#1253423).
- CVE-2025-40158: ipv6: use RCU in ip6_output() (bsc#1253402).
- CVE-2025-40159: xsk: Harden userspace-supplied xdp_desc validation (bsc#1253403).
- CVE-2025-40168: smc: Use __sk_dst_get() and dst_dev_rcu() in smc_clc_prfx_match() (bsc#1253427).
- CVE-2025-40169: bpf: Reject negative offsets for ALU ops (bsc#1253416).
- CVE-2025-40173: net/ip6_tunnel: Prevent perpetual tunnel growth (bsc#1253421).
- CVE-2025-40175: idpf: cleanup remaining SKBs in PTP flows (bsc#1253426).
- CVE-2025-40176: tls: wait for pending async decryptions if tls_strp_msg_hold fails (bsc#1253425).
- CVE-2025-40178: pid: Add a judgment for ns null in pid_nr_ns (bsc#1253463).
- CVE-2025-40185: ice: ice_adapter: release xa entry on adapter allocation failure (bsc#1253394).
- CVE-2025-40201: kernel/sys.c: fix the racy usage of task_lock(tsk-&gt;group_leader) in sys_prlimit64() paths (bsc#1253455).
- CVE-2025-40203: listmount: don&apos;t call path_put() under namespace semaphore (bsc#1253457).

The following non security issues were fixed:

- ACPI: scan: Update honor list for RPMI System MSI (stable-fixes).
- ACPICA: Update dsmethod.c to get rid of unused variable warning (stable-fixes).
- Disable CONFIG_CPU5_WDT The cpu5wdt driver doesn&apos;t implement a
  proper watchdog interface and has many code issues. It only handles
  obscure and obsolete hardware. Stop building and supporting this driver
  (jsc#PED-14062).
- Fix &quot;drm/xe: Don&apos;t allow evicting of BOs in same VM in array of VM binds&quot; (bsc#1252923)
- KVM: SVM: Delete IRTE link from previous vCPU before setting new IRTE (git-fixes).
- KVM: SVM: Delete IRTE link from previous vCPU irrespective of new routing (git-fixes).
- KVM: SVM: Mark VMCB_LBR dirty when MSR_IA32_DEBUGCTLMSR is updated (git-fixes).
- KVM: s390: improve interrupt cpu for wakeup (bsc#1235463).
- KVM: s390: kABI backport for &apos;last_sleep_cpu&apos; (bsc#1252352).
- KVM: x86/mmu: Return -EAGAIN if userspace deletes/moves memslot during prefault (git-fixes).
- PCI/ERR: Update device error_state already after reset (stable-fixes).
- PM: EM: Slightly reduce em_check_capacity_update() overhead (stable-fixes).
- Revert &quot;net/mlx5e: Update and set Xon/Xoff upon MTU set&quot; (git-fixes).
- Revert &quot;net/mlx5e: Update and set Xon/Xoff upon port speed set&quot; (git-fixes).
- Update config files: enable zstd module decompression (jsc#PED-14115).
- bpf/selftests: Fix test_tcpnotify_user (bsc#1253635).
- btrfs: do not clear read-only when adding sprout device (bsc#1253238).
- btrfs: do not update last_log_commit when logging inode due to a new name (git-fixes).
- dm: fix queue start/stop imbalance under suspend/load/resume races (bsc#1253386)
- drm/amd/display: Add AVI infoframe copy in copy_stream_update_to_stream (stable-fixes).
- drm/amd/display: update color on atomic commit time (stable-fixes).
- drm/amd/display: update dpp/disp clock from smu clock table (stable-fixes).
- drm/radeon: delete radeon_fence_process in is_signaled, no deadlock (stable-fixes).
- hwmon: (lenovo-ec-sensors) Update P8 supprt (stable-fixes).
- media: amphion: Delete v4l2_fh synchronously in .release() (stable-fixes).
- mount: handle NULL values in mnt_ns_release() (bsc#1254308)
- net/smc: Remove validation of reserved bits in CLC Decline (bsc#1252357).
- net: phy: move realtek PHY driver to its own subdirectory (jsc#PED-14353).
- net: phy: realtek: add defines for shadowed c45 standard registers (jsc#PED-14353).
- net: phy: realtek: add helper RTL822X_VND2_C22_REG (jsc#PED-14353).
- net: phy: realtek: change order of calls in C22 read_status() (jsc#PED-14353).
- net: phy: realtek: clear 1000Base-T link partner advertisement (jsc#PED-14353).
- net: phy: realtek: improve mmd register access for internal PHY&apos;s (jsc#PED-14353).
- net: phy: realtek: read duplex and gbit master from PHYSR register (jsc#PED-14353).
- net: phy: realtek: switch from paged to MMD ops in rtl822x functions (jsc#PED-14353).
- net: phy: realtek: use string choices helpers (jsc#PED-14353).
- net: xilinx: axienet: Fix IRQ coalescing packet count overflow (bsc#1250746)
- net: xilinx: axienet: Fix RX skb ring management in DMAengine mode (bsc#1250746)
- net: xilinx: axienet: Fix Tx skb circular buffer occupancy check in dmaengine xmit (bsc#1250746)
- nvmet-auth: update sc_c in host response (git-fixes bsc#1249397).
- nvmet-auth: update sc_c in target host hash calculation (git-fixes).
- perf list: Add IBM z17 event descriptions (jsc#PED-13611).
- platform/x86:intel/pmc: Update Arrow Lake telemetry GUID (git-fixes).
- powercap: intel_rapl: Add support for Panther Lake platform (jsc#PED-13949).
- pwm: pca9685: Use bulk write to atomicially update registers (stable-fixes).
- r8169: add PHY c45 ops for MDIO_MMD_VENDOR2 registers (jsc#PED-14353).
- r8169: add support for Intel Killer E5000 (jsc#PED-14353).
- r8169: add support for RTL8125BP rev.b (jsc#PED-14353).
- r8169: add support for RTL8125D rev.b (jsc#PED-14353).
- r8169: adjust version numbering for RTL8126 (jsc#PED-14353).
- r8169: align RTL8125 EEE config with vendor driver (jsc#PED-14353).
- r8169: align RTL8125/RTL8126 PHY config with vendor driver (jsc#PED-14353).
- r8169: align RTL8126 EEE config with vendor driver (jsc#PED-14353).
- r8169: align WAKE_PHY handling with r8125/r8126 vendor drivers (jsc#PED-14353).
- r8169: avoid duplicated messages if loading firmware fails and switch to warn level (jsc#PED-14353).
- r8169: don&apos;t take RTNL lock in rtl_task() (jsc#PED-14353).
- r8169: enable EEE at 2.5G per default on RTL8125B (jsc#PED-14353).
- r8169: enable RTL8168H/RTL8168EP/RTL8168FP ASPM support (jsc#PED-14353).
- r8169: fix inconsistent indenting in rtl8169_get_eth_mac_stats (jsc#PED-14353).
- r8169: implement additional ethtool stats ops (jsc#PED-14353).
- r8169: improve __rtl8169_set_wol (jsc#PED-14353).
- r8169: improve initialization of RSS registers on RTL8125/RTL8126 (jsc#PED-14353).
- r8169: improve rtl_set_d3_pll_down (jsc#PED-14353).
- r8169: increase max jumbo packet size on RTL8125/RTL8126 (jsc#PED-14353).
- r8169: remove leftover locks after reverted change (jsc#PED-14353).
- r8169: remove original workaround for RTL8125 broken rx issue (jsc#PED-14353).
- r8169: remove rtl_dash_loop_wait_high/low (jsc#PED-14353).
- r8169: remove support for chip version 11 (jsc#PED-14353).
- r8169: remove unused flag RTL_FLAG_TASK_RESET_NO_QUEUE_WAKE (jsc#PED-14353).
- r8169: replace custom flag with disable_work() et al (jsc#PED-14353).
- r8169: switch away from deprecated pcim_iomap_table (jsc#PED-14353).
- r8169: use helper r8169_mod_reg8_cond to simplify rtl_jumbo_config (jsc#PED-14353).
- ring-buffer: Update pages_touched to reflect persistent buffer content (git-fixes).
- s390/mm: Fix __ptep_rdp() inline assembly (bsc#1253643).
- sched/fair: Get rid of sched_domains_curr_level hack for tl-&gt;cpumask() (bsc#1246843).
- sched/fair: Have SD_SERIALIZE affect newidle balancing (bsc#1248792).
- sched/fair: Proportional newidle balance (bsc#1248792).
- sched/fair: Proportional newidle balance -KABI (bsc#1248792).
- sched/fair: Revert max_newidle_lb_cost bump (bsc#1248792).
- sched/fair: Skip sched_balance_running cmpxchg when balance is not due (bsc#1248792).
- sched/fair: Small cleanup to sched_balance_newidle() (bsc#1248792).
- sched/fair: Small cleanup to update_newidle_cost() (bsc#1248792).
- scsi: lpfc: Add capability to register Platform Name ID to fabric (bsc#1254119).
- scsi: lpfc: Allow support for BB credit recovery in point-to-point topology (bsc#1254119).
- scsi: lpfc: Ensure unregistration of rpis for received PLOGIs (bsc#1254119).
- scsi: lpfc: Fix leaked ndlp krefs when in point-to-point topology (bsc#1254119).
- scsi: lpfc: Fix reusing an ndlp that is marked NLP_DROPPED during FLOGI (bsc#1254119).
- scsi: lpfc: Modify kref handling for Fabric Controller ndlps (bsc#1254119).
- scsi: lpfc: Remove redundant NULL ptr assignment in lpfc_els_free_iocb() (bsc#1254119).
- scsi: lpfc: Revise discovery related function headers and comments (bsc#1254119).
- scsi: lpfc: Update lpfc version to 14.4.0.12 (bsc#1254119).
- scsi: lpfc: Update various NPIV diagnostic log messaging (bsc#1254119).
- selftests/run_kselftest.sh: Add --skip argument option (bsc#1254221).
- smpboot: introduce SDTL_INIT() helper to tidy sched topology setup (bsc#1246843).
- soc/tegra: fuse: speedo-tegra210: Update speedo IDs (git-fixes).
- spi: tegra210-quad: Check hardware status on timeout (bsc#1253155)
- spi: tegra210-quad: Fix timeout handling (bsc#1253155)
- spi: tegra210-quad: Refactor error handling into helper functions (bsc#1253155)
- spi: tegra210-quad: Update dummy sequence configuration (git-fixes)
- tcp_bpf: Call sk_msg_free() when tcp_bpf_send_verdict() fails to allocate psock-&gt;cork (bsc#1250705).
- wifi: ath11k: Add quirk entries for Thinkpad T14s Gen3 AMD (bsc#1254181).
- wifi: mt76: do not add wcid entries to sta poll list during MCU reset (bsc#1254315).
- wifi: mt76: introduce mt792x_config_mac_addr_list routine (bsc#1254315).
- wifi: mt76: mt7925: Fix logical vs bitwise typo (bsc#1254315).
- wifi: mt76: mt7925: Remove unnecessary if-check (bsc#1254315).
- wifi: mt76: mt7925: Simplify HIF suspend handling to avoid suspend fail (bsc#1254315).
- wifi: mt76: mt7925: add EHT control support based on the CLC data (bsc#1254315).
- wifi: mt76: mt7925: add handler to hif suspend/resume event (bsc#1254315).
- wifi: mt76: mt7925: add pci restore for hibernate (bsc#1254315).
- wifi: mt76: mt7925: config the dwell time by firmware (bsc#1254315).
- wifi: mt76: mt7925: extend MCU support for testmode (bsc#1254315).
- wifi: mt76: mt7925: fix CLC command timeout when suspend/resume (bsc#1254315).
- wifi: mt76: mt7925: fix missing hdr_trans_tlv command for broadcast wtbl (bsc#1254315).
- wifi: mt76: mt7925: fix the unfinished command of regd_notifier before suspend (bsc#1254315).
- wifi: mt76: mt7925: refine the txpower initialization flow (bsc#1254315).
- wifi: mt76: mt7925: replace zero-length array with flexible-array member (bsc#1254315).
- wifi: mt76: mt7925: update the channel usage when the regd domain changed (bsc#1254315).
- wifi: mt76: mt7925e: fix too long of wifi resume time (bsc#1254315).
- x86/smpboot: avoid SMT domain attach/destroy if SMT is not enabled (bsc#1246843).
- x86/smpboot: moves x86_topology to static initialize and truncate (bsc#1246843).
- x86/smpboot: remove redundant CONFIG_SCHED_SMT (bsc#1246843).
</description><license/><source url="https://download.opensuse.org/distribution/leap/16.0/repo/oss/aarch64" alias="https-download.opensuse.org-4bc031d9"/><issue-date time_t="1766165896" text="2025-12-19T17:38:16Z"/><issue-list><issue type="bugzilla" id="1235463"><title>Partner-L3: The installation takes too long on &quot;Preparing disks&quot; process with thinksystem 940-16i raid card and u.3 nvme disk</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1235463</href></issue><issue type="bugzilla" id="1243474"><title>VUL-0: CVE-2025-37916: kernel: pds_core: remove write-after-free of client_id</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1243474</href></issue><issue type="bugzilla" id="1245193"><title>backport nvmet-loop fixes</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1245193</href></issue><issue type="bugzilla" id="1245431"><title>HugeTLB - fix unshare pages</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1245431</href></issue><issue type="bugzilla" id="1245498"><title>VUL-0: CVE-2025-38084: kernel: mm/hugetlb: unshare page tables during VMA split, not before</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1245498</href></issue><issue type="bugzilla" id="1245499"><title>VUL-0: CVE-2025-38085: kernel: mm/hugetlb: fix huge_pmd_unshare() vs GUP-fast race</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1245499</href></issue><issue type="bugzilla" id="1246328"><title>VUL-0: CVE-2025-38321: kernel: smb: Log an error when close_all_cached_dirs fails</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1246328</href></issue><issue type="bugzilla" id="1246843"><title>SLES16 RC1 [ Regression ] [ 6.12.0-160000.16-default ]: Qemu guest boot fails with continuous soft-lockups while trying to boot with 8 NUMA nodes.</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1246843</href></issue><issue type="bugzilla" id="1247500"><title>nvme over FC: kernel soft lockup on module removal</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1247500</href></issue><issue type="bugzilla" id="1248792"><title>[PBOnline POWER10 HANA2SPS086 &amp; 079] benchCppDistributedTransactionBarrier.py reports up to 40% performance deterioration on Power10</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1248792</href></issue><issue type="bugzilla" id="1249256"><title>VUL-0: CVE-2025-38728: kernel: smb3: fix for slab out of bounds on mount to ksmbd</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1249256</href></issue><issue type="bugzilla" id="1249397"><title>[NetApp SLES15 SP7 Bug]: sc_c field not updated in host response to controller challenge for secure concat</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1249397</href></issue><issue type="bugzilla" id="1249912"><title>VUL-0: CVE-2022-50253: kernel: bpf: make sure skb-&gt;len != 0 when redirecting to a tunneling device</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1249912</href></issue><issue type="bugzilla" id="1249977"><title>VUL-0: CVE-2025-39831: kernel: fbnic: Move phylink resume out of service_task and into open/close</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1249977</href></issue><issue type="bugzilla" id="1249982"><title>VUL-0: CVE-2025-39805: kernel: net: macb: fix unregister_netdev call order in macb_remove()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1249982</href></issue><issue type="bugzilla" id="1250034"><title>VUL-0: CVE-2025-39822: kernel: io_uring/kbuf: fix signedness in this_len calculation</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1250034</href></issue><issue type="bugzilla" id="1250176"><title>VUL-0: CVE-2025-39819: kernel: fs/smb: Fix inconsistent refcnt update</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1250176</href></issue><issue type="bugzilla" id="1250237"><title>nftables stack guard hit + kernel panic on synproxy in output chain</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1250237</href></issue><issue type="bugzilla" id="1250252"><title>VUL-0: CVE-2025-39859: kernel: ptp: ocp: fix use-after-free bugs causing by ptp_ocp_watchdog</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1250252</href></issue><issue type="bugzilla" id="1250705"><title>VUL-0: CVE-2025-39913: kernel: tcp_bpf: Call sk_msg_free() when tcp_bpf_send_verdict() fails to allocate psock-&gt;cork.</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1250705</href></issue><issue type="bugzilla" id="1250723"><title>VUL-0: CVE-2025-39917: kernel: bpf: Fix out-of-bounds dynptr write in bpf_crypto_crypt</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1250723</href></issue><issue type="bugzilla" id="1250746"><title>VUL-0: CVE-2025-39897: kernel: net: xilinx: axienet: Add error handling for RX metadata pointer retrieval</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1250746</href></issue><issue type="bugzilla" id="1251120"><title>VUL-0: CVE-2025-39944: kernel: octeontx2-pf: Fix use-after-free bugs in otx2_sync_tstamp()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1251120</href></issue><issue type="bugzilla" id="1251817"><title>VUL-0: CVE-2025-39961: kernel: iommu/amd/pgtbl: Fix possible race while increase page table level</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1251817</href></issue><issue type="bugzilla" id="1252054"><title>VUL-0: CVE-2025-39990: kernel: bpf: Check the helper function is valid in get_helper_proto</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252054</href></issue><issue type="bugzilla" id="1252063"><title>VUL-0: CVE-2025-39980: kernel: nexthop: Forbid FDB status change while nexthop is in a group</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252063</href></issue><issue type="bugzilla" id="1252301"><title>VUL-0: CVE-2025-40003: kernel: net: mscc: ocelot: Fix use-after-free caused by cyclic delayed work</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252301</href></issue><issue type="bugzilla" id="1252303"><title>VUL-0: CVE-2025-40001: kernel: scsi: mvsas: Fix use-after-free bugs in mvs_work_queue</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252303</href></issue><issue type="bugzilla" id="1252342"><title>VUL-0: CVE-2025-40006: kernel: mm/hugetlb: fix folio is still mapped when deleted</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252342</href></issue><issue type="bugzilla" id="1252352"><title>SLES 15 SP7 - KVM: s390: improve interrupt cpu for wakeup</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252352</href></issue><issue type="bugzilla" id="1252357"><title>SLES 16.0 - net/smc: Remove validation of reserved bits in CLC Decline msg</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252357</href></issue><issue type="bugzilla" id="1252681"><title>VUL-0: CVE-2025-40021: kernel: tracing: dynevent: Add a missing lockdown check on dynevent</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252681</href></issue><issue type="bugzilla" id="1252686"><title>VUL-0: CVE-2025-40024: kernel: vhost: Take a reference on the task in struct vhost_task.</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252686</href></issue><issue type="bugzilla" id="1252763"><title>VUL-0: CVE-2025-40027: kernel: net/9p: fix double req put in p9_fd_cancelled</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252763</href></issue><issue type="bugzilla" id="1252776"><title>VUL-0: CVE-2025-40081: kernel: perf: arm_spe: Prevent overflow in PERF_IDX2OFF()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252776</href></issue><issue type="bugzilla" id="1252779"><title>VUL-0: CVE-2025-40031: kernel: tee: fix register_shm_helper()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252779</href></issue><issue type="bugzilla" id="1252790"><title>VUL-0: CVE-2025-40047: kernel: io_uring/waitid: always prune wait queue entry in io_waitid_wait()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252790</href></issue><issue type="bugzilla" id="1252794"><title>VUL-0: CVE-2025-40074: kernel: ipv4: start using dst_dev_rcu()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252794</href></issue><issue type="bugzilla" id="1252795"><title>VUL-0: CVE-2025-40075: kernel: tcp_metrics: use dst_dev_net_rcu()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252795</href></issue><issue type="bugzilla" id="1252808"><title>VUL-0: CVE-2025-40053: kernel: net: dlink: handle copy_thresh allocation failure</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252808</href></issue><issue type="bugzilla" id="1252809"><title>VUL-0: CVE-2025-40059: kernel: coresight: Fix incorrect handling for return value of devm_kzalloc</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252809</href></issue><issue type="bugzilla" id="1252817"><title>VUL-0: CVE-2025-40038: kernel: KVM: SVM: Skip fastpath emulation on VM-Exit if next RIP isn&apos;t valid</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252817</href></issue><issue type="bugzilla" id="1252821"><title>VUL-0: CVE-2025-40055: kernel: ocfs2: fix double free in user_cluster_connect()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252821</href></issue><issue type="bugzilla" id="1252824"><title>VUL-0: CVE-2025-40033: kernel: remoteproc: pru: Fix potential NULL pointer dereference in pru_rproc_set_ctable()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252824</href></issue><issue type="bugzilla" id="1252836"><title>VUL-0: CVE-2025-40070: kernel: pps: fix warning in pps_register_cdev when register device fail</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252836</href></issue><issue type="bugzilla" id="1252845"><title>VUL-0: CVE-2025-40064: kernel: smc: Fix use-after-free in __pnet_find_base_ndev().</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252845</href></issue><issue type="bugzilla" id="1252901"><title>VUL-0: CVE-2025-40101: kernel: btrfs: fix memory leaks when rejecting a non SINGLE data profile without an RST</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252901</href></issue><issue type="bugzilla" id="1252912"><title>VUL-0: CVE-2025-40083: kernel: net/sched: sch_qfq: Fix null-deref in agg_dequeue</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252912</href></issue><issue type="bugzilla" id="1252917"><title>VUL-0: CVE-2025-40098: kernel: ALSA: hda: cs35l41: Fix NULL pointer dereference in cs35l41_get_acpi_mute_state()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252917</href></issue><issue type="bugzilla" id="1252919"><title>VUL-0: CVE-2025-40102: kernel: KVM: arm64: Prevent access to vCPU events before init</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252919</href></issue><issue type="bugzilla" id="1252923"><title>VUL-0: CVE-2025-40086: kernel: drm/xe: Don&apos;t allow evicting of BOs in same VM in array of VM binds</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252923</href></issue><issue type="bugzilla" id="1252928"><title>VUL-0: CVE-2025-40105: kernel: vfs: Don&apos;t leak disconnected dentries on umount</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252928</href></issue><issue type="bugzilla" id="1253018"><title>VUL-0: CVE-2025-40107: kernel: can: hi311x: fix null pointer dereference when resuming from sleep before interface was enabled</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253018</href></issue><issue type="bugzilla" id="1253155"><title>Nvidia Grace: Backport: spi: tegra210-quad: Improve timeout handling under high system load</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253155</href></issue><issue type="bugzilla" id="1253176"><title>VUL-0: CVE-2025-40109: kernel: crypto: rng - Ensure set_ent is always present</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253176</href></issue><issue type="bugzilla" id="1253238"><title>SLES16.0 [6.12.0-160000.18-default] FS/BTRFS: btrfs/282 btrfs/323 generic/363 tests fails with missing kernel fix</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253238</href></issue><issue type="bugzilla" id="1253275"><title>VUL-0: CVE-2025-40110: kernel: drm/vmwgfx: Fix a null-ptr access in the cursor snooper</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253275</href></issue><issue type="bugzilla" id="1253318"><title>VUL-0: CVE-2025-40115: kernel: scsi: mpt3sas: Fix crash in transport port remove by using ioc_info()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253318</href></issue><issue type="bugzilla" id="1253324"><title>VUL-0: CVE-2025-40116: kernel: usb: host: max3421-hcd: Fix error pointer dereference in probe cleanup</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253324</href></issue><issue type="bugzilla" id="1253328"><title>VUL-0: CVE-2025-40133: kernel: mptcp: Use __sk_dst_get() and dst_dev_rcu() in mptcp_active_enable().</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253328</href></issue><issue type="bugzilla" id="1253330"><title>VUL-0: CVE-2025-40132: kernel: ASoC: Intel: sof_sdw: Prevent jump to NULL add_sidecar callback</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253330</href></issue><issue type="bugzilla" id="1253342"><title>VUL-0: CVE-2025-40135: kernel: ipv6: use RCU in ip6_xmit()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253342</href></issue><issue type="bugzilla" id="1253348"><title>VUL-0: CVE-2025-40142: kernel: ALSA: pcm: Disable bottom softirqs as part of spin_lock_irq() on PREEMPT_RT</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253348</href></issue><issue type="bugzilla" id="1253349"><title>VUL-0: CVE-2025-40140: kernel: net: usb: Remove disruptive netif_wake_queue in rtl8150_set_multicast</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253349</href></issue><issue type="bugzilla" id="1253352"><title>VUL-0: CVE-2025-40141: kernel: Bluetooth: ISO: Fix possible UAF on iso_conn_free</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253352</href></issue><issue type="bugzilla" id="1253355"><title>VUL-0: CVE-2025-40149: kernel: tls: Use __sk_dst_get() and dst_dev_rcu() in get_netdev_for_sock().</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253355</href></issue><issue type="bugzilla" id="1253360"><title>VUL-0: CVE-2025-40120: kernel: net: usb: asix: hold PM usage ref to avoid PM/MDIO + RTNL deadlock</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253360</href></issue><issue type="bugzilla" id="1253362"><title>VUL-0: CVE-2025-40111: kernel: drm/vmwgfx: Fix Use-after-free in validation</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253362</href></issue><issue type="bugzilla" id="1253363"><title>VUL-0: CVE-2025-40118: kernel: scsi: pm80xx: Fix array-index-out-of-of-bounds on rmmod</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253363</href></issue><issue type="bugzilla" id="1253367"><title>VUL-0: CVE-2025-40121: kernel: ASoC: Intel: bytcr_rt5651: Fix invalid quirk input mapping</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253367</href></issue><issue type="bugzilla" id="1253369"><title>VUL-0: CVE-2025-40127: kernel: hwrng: ks-sa - fix division by zero in ks_sa_rng_init</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253369</href></issue><issue type="bugzilla" id="1253386"><title>VUL-0: CVE-2025-40134: kernel: dm: fix NULL pointer dereference in __dm_suspend()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253386</href></issue><issue type="bugzilla" id="1253394"><title>VUL-0: CVE-2025-40185: kernel: ice: ice_adapter: release xa entry on adapter allocation failure</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253394</href></issue><issue type="bugzilla" id="1253395"><title>VUL-0: CVE-2025-40207: kernel: media: v4l2-subdev: Fix alloc failure check in v4l2_subdev_call_state_try()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253395</href></issue><issue type="bugzilla" id="1253402"><title>VUL-0: CVE-2025-40158: kernel: ipv6: use RCU in ip6_output()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253402</href></issue><issue type="bugzilla" id="1253403"><title>VUL-0: CVE-2025-40159: kernel: xsk: Harden userspace-supplied xdp_desc validation</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253403</href></issue><issue type="bugzilla" id="1253405"><title>VUL-0: CVE-2025-40165: kernel: media: nxp: imx8-isi: m2m: Fix streaming cleanup on release</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253405</href></issue><issue type="bugzilla" id="1253407"><title>VUL-0: CVE-2025-40164: kernel: usbnet: Fix using smp_processor_id() in preemptible code warnings</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253407</href></issue><issue type="bugzilla" id="1253408"><title>VUL-0: CVE-2025-40153: kernel: mm: hugetlb: avoid soft lockup when mprotect to large memory area</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253408</href></issue><issue type="bugzilla" id="1253409"><title>VUL-0: CVE-2025-40139: kernel: smc: Use __sk_dst_get() and dst_dev_rcu() in in smc_clc_prfx_set().</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253409</href></issue><issue type="bugzilla" id="1253410"><title>VUL-0: CVE-2025-40161: kernel: mailbox: zynqmp-ipi: Fix SGI cleanup on unbind</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253410</href></issue><issue type="bugzilla" id="1253412"><title>VUL-0: CVE-2025-40171: kernel: nvmet-fc: move lsop put work to nvmet_fc_ls_req_op</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253412</href></issue><issue type="bugzilla" id="1253416"><title>VUL-0: CVE-2025-40169: kernel: bpf: Reject negative offsets for ALU ops</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253416</href></issue><issue type="bugzilla" id="1253421"><title>VUL-0: CVE-2025-40173: kernel: net/ip6_tunnel: Prevent perpetual tunnel growth</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253421</href></issue><issue type="bugzilla" id="1253422"><title>VUL-0: CVE-2025-40162: kernel: ASoC: amd/sdw_utils: avoid NULL deref when devm_kasprintf() fails</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253422</href></issue><issue type="bugzilla" id="1253423"><title>VUL-0: CVE-2025-40157: kernel: EDAC/i10nm: Skip DIMM enumeration on a disabled memory controller</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253423</href></issue><issue type="bugzilla" id="1253424"><title>VUL-0: CVE-2025-40172: kernel: accel/qaic: Treat remaining == 0 as error in find_and_map_user_pages()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253424</href></issue><issue type="bugzilla" id="1253425"><title>VUL-0: CVE-2025-40176: kernel: tls: wait for pending async decryptions if tls_strp_msg_hold fails</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253425</href></issue><issue type="bugzilla" id="1253426"><title>VUL-0: CVE-2025-40175: kernel: idpf: cleanup remaining SKBs in PTP flows</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253426</href></issue><issue type="bugzilla" id="1253427"><title>VUL-0: CVE-2025-40168: kernel: smc: Use __sk_dst_get() and dst_dev_rcu() in smc_clc_prfx_match().</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253427</href></issue><issue type="bugzilla" id="1253428"><title>VUL-0: CVE-2025-40156: kernel: PM / devfreq: mtk-cci: Fix potential error pointer dereference in probe()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253428</href></issue><issue type="bugzilla" id="1253431"><title>VUL-0: CVE-2025-40154: kernel: ASoC: Intel: bytcr_rt5640: Fix invalid quirk input mapping</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253431</href></issue><issue type="bugzilla" id="1253433"><title>VUL-0: CVE-2025-40166: kernel: drm/xe/guc: Check GuC running state before deregistering exec queue</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253433</href></issue><issue type="bugzilla" id="1253436"><title>VUL-0: CVE-2025-40204: kernel: sctp: Fix MAC comparison to be constant-time</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253436</href></issue><issue type="bugzilla" id="1253438"><title>VUL-0: CVE-2025-40186: kernel: tcp: Don&apos;t call reqsk_fastopen_remove() in tcp_conn_request().</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253438</href></issue><issue type="bugzilla" id="1253440"><title>VUL-0: CVE-2025-40180: kernel: mailbox: zynqmp-ipi: Fix out-of-bounds access in mailbox cleanup loop</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253440</href></issue><issue type="bugzilla" id="1253441"><title>VUL-0: CVE-2025-40183: kernel: bpf: Fix metadata_dst leak __bpf_redirect_neigh_v{4,6}</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253441</href></issue><issue type="bugzilla" id="1253443"><title>VUL-0: CVE-2025-40177: kernel: accel/qaic: Fix bootlog initialization ordering</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253443</href></issue><issue type="bugzilla" id="1253445"><title>VUL-0: CVE-2025-40194: kernel: cpufreq: intel_pstate: Fix object lifecycle issue in update_qos_request()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253445</href></issue><issue type="bugzilla" id="1253448"><title>VUL-0: CVE-2025-40200: kernel: Squashfs: reject negative file sizes in squashfs_read_inode()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253448</href></issue><issue type="bugzilla" id="1253449"><title>VUL-0: CVE-2025-40188: kernel: pwm: berlin: Fix wrong register in suspend/resume</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253449</href></issue><issue type="bugzilla" id="1253450"><title>VUL-0: CVE-2025-40197: kernel: media: mc: Clear minor number before put device</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253450</href></issue><issue type="bugzilla" id="1253451"><title>VUL-0: CVE-2025-40202: kernel: ipmi: Rework user message limit handling</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253451</href></issue><issue type="bugzilla" id="1253453"><title>VUL-0: CVE-2025-40198: kernel: ext4: avoid potential buffer over-read in parse_apply_sb_mount_options()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253453</href></issue><issue type="bugzilla" id="1253455"><title>VUL-0: CVE-2025-40201: kernel: kernel/sys.c: fix the racy usage of task_lock(tsk-&gt;group_leader) in sys_prlimit64() paths</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253455</href></issue><issue type="bugzilla" id="1253456"><title>VUL-0: CVE-2025-40205: kernel: btrfs: avoid potential out-of-bounds in btrfs_encode_fh()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253456</href></issue><issue type="bugzilla" id="1253457"><title>VUL-0: CVE-2025-40203: kernel: listmount: don&apos;t call path_put() under namespace semaphore</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253457</href></issue><issue type="bugzilla" id="1253463"><title>VUL-0: CVE-2025-40178: kernel: pid: Add a check for ns is null in pid_nr_ns</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253463</href></issue><issue type="bugzilla" id="1253472"><title>VUL-0: CVE-2025-40129: kernel: sunrpc: fix null pointer dereference on zero-length checksum</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253472</href></issue><issue type="bugzilla" id="1253622"><title>VUL-0: CVE-2025-40192: kernel: Revert &quot;ipmi: fix msg stack when IPMI is disconnected&quot;</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253622</href></issue><issue type="bugzilla" id="1253624"><title>VUL-0: CVE-2025-40196: kernel: fs: quota: create dedicated workqueue for quota_release_work</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253624</href></issue><issue type="bugzilla" id="1253635"><title>selftests: bpf: test_tcpnotify_user segfaults</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253635</href></issue><issue type="bugzilla" id="1253643"><title>SLES 16 SP0 - s390/mm: Fix __ptep_rdp() inline assembly</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253643</href></issue><issue type="bugzilla" id="1253647"><title>VUL-0: CVE-2025-40187: kernel: net/sctp: fix a null dereference in sctp_disposition sctp_sf_do_5_1D_ce()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253647</href></issue><issue type="bugzilla" id="1254119"><title>Update Broadcom Emulex lpfc driver for SL-16.1  to 14.4.0.12</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1254119</href></issue><issue type="bugzilla" id="1254181"><title>ath11k WiFi suspend/resume breakage on Lenovo Thinkpad T14s Gen 3 AMD</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1254181</href></issue><issue type="bugzilla" id="1254221"><title>selftests: run_kselftest.sh: openQA test fails when skipping tests from a large collection</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1254221</href></issue><issue type="bugzilla" id="1254308"><title>SLE-16.0 KOTD: NULL pointer dereference in listmount() syscall</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1254308</href></issue><issue type="bugzilla" id="1254315"><title>Lenovo Thinkpad P14s doesn&apos;t sleep</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1254315</href></issue><issue type="cve" id="CVE-2022-50253"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-50253</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-50253</href></issue><issue type="cve" id="CVE-2025-37916"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-37916</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-37916</href></issue><issue type="cve" id="CVE-2025-38084"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-38084</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-38084</href></issue><issue type="cve" id="CVE-2025-38085"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-38085</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-38085</href></issue><issue type="cve" id="CVE-2025-38321"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-38321</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-38321</href></issue><issue type="cve" id="CVE-2025-38728"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-38728</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-38728</href></issue><issue type="cve" id="CVE-2025-39805"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39805</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39805</href></issue><issue type="cve" id="CVE-2025-39819"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39819</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39819</href></issue><issue type="cve" id="CVE-2025-39822"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39822</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39822</href></issue><issue type="cve" id="CVE-2025-39831"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39831</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39831</href></issue><issue type="cve" id="CVE-2025-39859"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39859</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39859</href></issue><issue type="cve" id="CVE-2025-39897"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39897</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39897</href></issue><issue type="cve" id="CVE-2025-39917"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39917</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39917</href></issue><issue type="cve" id="CVE-2025-39944"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39944</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39944</href></issue><issue type="cve" id="CVE-2025-39961"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39961</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39961</href></issue><issue type="cve" id="CVE-2025-39980"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39980</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39980</href></issue><issue type="cve" id="CVE-2025-39990"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39990</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39990</href></issue><issue type="cve" id="CVE-2025-40001"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40001</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40001</href></issue><issue type="cve" id="CVE-2025-40003"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40003</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40003</href></issue><issue type="cve" id="CVE-2025-40006"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40006</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40006</href></issue><issue type="cve" id="CVE-2025-40021"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40021</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40021</href></issue><issue type="cve" id="CVE-2025-40024"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40024</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40024</href></issue><issue type="cve" id="CVE-2025-40027"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40027</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40027</href></issue><issue type="cve" id="CVE-2025-40031"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40031</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40031</href></issue><issue type="cve" id="CVE-2025-40033"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40033</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40033</href></issue><issue type="cve" id="CVE-2025-40038"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40038</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40038</href></issue><issue type="cve" id="CVE-2025-40047"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40047</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40047</href></issue><issue type="cve" id="CVE-2025-40053"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40053</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40053</href></issue><issue type="cve" id="CVE-2025-40055"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40055</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40055</href></issue><issue type="cve" id="CVE-2025-40059"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40059</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40059</href></issue><issue type="cve" id="CVE-2025-40064"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40064</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40064</href></issue><issue type="cve" id="CVE-2025-40070"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40070</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40070</href></issue><issue type="cve" id="CVE-2025-40074"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40074</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40074</href></issue><issue type="cve" id="CVE-2025-40075"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40075</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40075</href></issue><issue type="cve" id="CVE-2025-40081"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40081</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40081</href></issue><issue type="cve" id="CVE-2025-40083"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40083</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40083</href></issue><issue type="cve" id="CVE-2025-40086"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40086</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40086</href></issue><issue type="cve" id="CVE-2025-40098"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40098</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40098</href></issue><issue type="cve" id="CVE-2025-40101"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40101</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40101</href></issue><issue type="cve" id="CVE-2025-40102"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40102</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40102</href></issue><issue type="cve" id="CVE-2025-40105"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40105</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40105</href></issue><issue type="cve" id="CVE-2025-40107"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40107</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40107</href></issue><issue type="cve" id="CVE-2025-40109"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40109</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40109</href></issue><issue type="cve" id="CVE-2025-40110"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40110</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40110</href></issue><issue type="cve" id="CVE-2025-40111"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40111</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40111</href></issue><issue type="cve" id="CVE-2025-40115"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40115</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40115</href></issue><issue type="cve" id="CVE-2025-40116"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40116</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40116</href></issue><issue type="cve" id="CVE-2025-40118"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40118</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40118</href></issue><issue type="cve" id="CVE-2025-40120"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40120</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40120</href></issue><issue type="cve" id="CVE-2025-40121"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40121</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40121</href></issue><issue type="cve" id="CVE-2025-40127"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40127</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40127</href></issue><issue type="cve" id="CVE-2025-40129"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40129</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40129</href></issue><issue type="cve" id="CVE-2025-40132"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40132</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40132</href></issue><issue type="cve" id="CVE-2025-40133"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40133</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40133</href></issue><issue type="cve" id="CVE-2025-40134"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40134</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40134</href></issue><issue type="cve" id="CVE-2025-40135"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40135</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40135</href></issue><issue type="cve" id="CVE-2025-40139"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40139</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40139</href></issue><issue type="cve" id="CVE-2025-40140"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40140</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40140</href></issue><issue type="cve" id="CVE-2025-40141"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40141</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40141</href></issue><issue type="cve" id="CVE-2025-40142"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40142</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40142</href></issue><issue type="cve" id="CVE-2025-40149"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40149</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40149</href></issue><issue type="cve" id="CVE-2025-40153"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40153</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40153</href></issue><issue type="cve" id="CVE-2025-40154"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40154</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40154</href></issue><issue type="cve" id="CVE-2025-40156"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40156</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40156</href></issue><issue type="cve" id="CVE-2025-40157"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40157</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40157</href></issue><issue type="cve" id="CVE-2025-40158"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40158</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40158</href></issue><issue type="cve" id="CVE-2025-40159"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40159</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40159</href></issue><issue type="cve" id="CVE-2025-40161"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40161</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40161</href></issue><issue type="cve" id="CVE-2025-40162"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40162</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40162</href></issue><issue type="cve" id="CVE-2025-40164"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40164</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40164</href></issue><issue type="cve" id="CVE-2025-40165"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40165</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40165</href></issue><issue type="cve" id="CVE-2025-40166"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40166</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40166</href></issue><issue type="cve" id="CVE-2025-40168"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40168</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40168</href></issue><issue type="cve" id="CVE-2025-40169"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40169</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40169</href></issue><issue type="cve" id="CVE-2025-40171"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40171</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40171</href></issue><issue type="cve" id="CVE-2025-40172"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40172</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40172</href></issue><issue type="cve" id="CVE-2025-40173"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40173</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40173</href></issue><issue type="cve" id="CVE-2025-40175"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40175</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40175</href></issue><issue type="cve" id="CVE-2025-40176"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40176</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40176</href></issue><issue type="cve" id="CVE-2025-40177"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40177</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40177</href></issue><issue type="cve" id="CVE-2025-40178"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40178</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40178</href></issue><issue type="cve" id="CVE-2025-40180"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40180</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40180</href></issue><issue type="cve" id="CVE-2025-40183"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40183</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40183</href></issue><issue type="cve" id="CVE-2025-40185"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40185</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40185</href></issue><issue type="cve" id="CVE-2025-40186"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40186</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40186</href></issue><issue type="cve" id="CVE-2025-40187"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40187</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40187</href></issue><issue type="cve" id="CVE-2025-40188"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40188</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40188</href></issue><issue type="cve" id="CVE-2025-40192"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40192</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40192</href></issue><issue type="cve" id="CVE-2025-40194"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40194</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40194</href></issue><issue type="cve" id="CVE-2025-40196"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40196</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40196</href></issue><issue type="cve" id="CVE-2025-40197"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40197</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40197</href></issue><issue type="cve" id="CVE-2025-40198"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40198</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40198</href></issue><issue type="cve" id="CVE-2025-40200"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40200</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40200</href></issue><issue type="cve" id="CVE-2025-40201"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40201</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40201</href></issue><issue type="cve" id="CVE-2025-40202"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40202</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40202</href></issue><issue type="cve" id="CVE-2025-40203"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40203</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40203</href></issue><issue type="cve" id="CVE-2025-40204"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40204</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40204</href></issue><issue type="cve" id="CVE-2025-40205"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40205</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40205</href></issue><issue type="cve" id="CVE-2025-40206"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40206</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40206</href></issue><issue type="cve" id="CVE-2025-40207"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40207</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40207</href></issue><issue type="jira" id="PED-13611"><title>https://jira.suse.com/browse/PED-13611</title><href>https://jira.suse.com/browse/PED-13611</href></issue><issue type="jira" id="PED-13949"><title>https://jira.suse.com/browse/PED-13949</title><href>https://jira.suse.com/browse/PED-13949</href></issue><issue type="jira" id="PED-14062"><title>https://jira.suse.com/browse/PED-14062</title><href>https://jira.suse.com/browse/PED-14062</href></issue><issue type="jira" id="PED-14115"><title>https://jira.suse.com/browse/PED-14115</title><href>https://jira.suse.com/browse/PED-14115</href></issue><issue type="jira" id="PED-14353"><title>https://jira.suse.com/browse/PED-14353</title><href>https://jira.suse.com/browse/PED-14353</href></issue></issue-list></update><update kind="patch" name="openSUSE-Leap-16.0-112" edition="1" arch="noarch" status="needed" category="security" severity="important" pkgmanager="false" restart="false" interactive="false"><summary>Security update for qemu</summary><description>This update for qemu fixes the following issues:

Update to version 10.0.7.

Security issues fixed:

- CVE-2025-12464: stack-based buffer overflow in the e1000 network device operations can be exploited by a malicious
  guest user to crash the QEMU process on the host (bsc#1253002).
- CVE-2025-11234: use-after-free in WebSocket handshake operations can be exploited by a malicious client with network
  access to the VNC WebSocket port to cause a denial-of-service (bsc#1250984).

Other updates and bugfixes:

- Version 10.0.7:
  * kvm: Fix kvm_vm_ioctl() and kvm_device_ioctl() return value
  * docs/devel: Update URL for make-pullreq script
  * target/arm: Fix assert on BRA.
  * hw/aspeed/{xdma, rtc, sdhci}: Fix endianness to DEVICE_LITTLE_ENDIAN
  * hw/core/machine: Provide a description for aux-ram-share property
  * hw/pci: Make msix_init take a uint32_t for nentries
  * block/io_uring: avoid potentially getting stuck after resubmit at the end of ioq_submit()
  * block-backend: Fix race when resuming queued requests
  * ui/vnc: Fix qemu abort when query vnc info
  * chardev/char-pty: Do not ignore chr_write() failures
  * hw/display/exynos4210_fimd: Account for zero length in fimd_update_memory_section()
  * hw/arm/armv7m: Disable reentrancy guard for v7m_sysreg_ns_ops MRs
  * hw/arm/aspeed: Fix missing SPI IRQ connection causing DMA interrupt failure
  * migration: Fix transition to COLO state from precopy
  * Full backport list: https://lore.kernel.org/qemu-devel/1765037524.347582.2700543.nullmailer@tls.msk.ru/

- Version 10.0.6:
  * linux-user/microblaze: Fix little-endianness binary
  * target/hppa: correct size bit parity for fmpyadd
  * target/i386: user: do not set up a valid LDT on reset
  * async: access bottom half flags with qatomic_read
  * target/i386: fix x86_64 pushw op
  * i386/tcg/smm_helper: Properly apply DR values on SMM entry / exit
  * i386/cpu: Prevent delivering SIPI during SMM in TCG mode
  * i386/kvm: Expose ARCH_CAP_FB_CLEAR when invulnerable to MDS
  * target/i386: Fix CR2 handling for non-canonical addresses
  * block/curl.c: Use explicit long constants in curl_easy_setopt calls
  * pcie_sriov: Fix broken MMIO accesses from SR-IOV VFs
  * target/riscv: rvv: Fix vslide1[up|down].vx unexpected result when XLEN2 and SEWd
  * target/riscv: Fix ssamoswap error handling
  * Full backport list: https://lore.kernel.org/qemu-devel/1761022287.744330.6357.nullmailer@tls.msk.ru/

- Version 10.0.5:
  * tests/functional/test_aarch64_sbsaref_freebsd: Fix the URL of the ISO image
  * tests/functional/test_ppc_bamboo: Replace broken link with working assets
  * physmem: Destroy all CPU AddressSpaces on unrealize
  * memory: New AS helper to serialize destroy+free
  * include/system/memory.h: Clarify address_space_destroy() behaviour
  * migration: Fix state transition in postcopy_start() error handling
  * target/riscv: rvv: Modify minimum VLEN according to enabled vector extensions
  * target/riscv: rvv: Replace checking V by checking Zve32x
  * target/riscv: Fix endianness swap on compressed instructions
  * hw/riscv/riscv-iommu: Fixup PDT Nested Walk
  * Full backport list: https://lore.kernel.org/qemu-devel/1759986125.676506.643525.nullmailer@tls.msk.ru/

- [openSUSE][RPM]: really fix *-virtio-gpu-pci dependency on ARM (bsc#1254286).
- [openSUSE][RPM] spec: make glusterfs support conditional (bsc#1254494).
</description><license/><source url="http://cdn.opensuse.org/distribution/leap/16.0/repo/oss/aarch64" alias="openSUSE:repo-oss"/><issue-date time_t="1766061723" text="2025-12-18T12:42:03Z"/><issue-list><issue type="cve" id="CVE-2025-11234"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-11234</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-11234</href></issue><issue type="cve" id="CVE-2025-12464"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-12464</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-12464</href></issue><issue type="bugzilla" id="1254286"><title>cockpit - Virtual machines - Import VM - import qcow2 image failed and it fails to start because of missing &apos;qemu-hw-display-virtio-gpu&apos; and &apos;qemu-hw-display-virtio-gpu-pci&apos;</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1254286</href></issue><issue type="bugzilla" id="1230042"><title>[SLM6.2][PPC64LE] Fail to start QEMU - qemu-system-ppc64</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1230042</href></issue><issue type="bugzilla" id="1254494"><title>glusterfs: drop 32-bit architectures</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1254494</href></issue><issue type="bugzilla" id="1250984"><title>VUL-0: CVE-2025-11234: qemu: qemu-kvm: use-after-free in websocket handshake code can lead to denial of service</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1250984</href></issue><issue type="bugzilla" id="1253002"><title>VUL-0: CVE-2025-12464: qemu: net: pad packets to minimum length in qemu_receive_packet()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253002</href></issue></issue-list></update><update kind="patch" name="openSUSE-Leap-16.0-113" edition="1" arch="noarch" status="needed" category="security" severity="important" pkgmanager="false" restart="true" interactive="true"><summary>Security update for the Linux Kernel</summary><description>
The SUSE Linux Enterprise 16.0 kernel was updated to fix various security issues

The following security issues were fixed:

- CVE-2022-50253: bpf: make sure skb-&gt;len != 0 when redirecting to a tunneling device (bsc#1249912).
- CVE-2025-37916: pds_core: remove write-after-free of client_id (bsc#1243474).
- CVE-2025-38084: mm/hugetlb: unshare page tables during VMA split, not before (bsc#1245431 bsc#1245498).
- CVE-2025-38085: mm/hugetlb: fix huge_pmd_unshare() vs GUP-fast race (bsc#1245431 bsc#1245499).
- CVE-2025-38321: smb: Log an error when close_all_cached_dirs fails (bsc#1246328).
- CVE-2025-38728: smb3: fix for slab out of bounds on mount to ksmbd (bsc#1249256).
- CVE-2025-39805: net: macb: fix unregister_netdev call order in macb_remove() (bsc#1249982).
- CVE-2025-39819: fs/smb: Fix inconsistent refcnt update (bsc#1250176).
- CVE-2025-39822: io_uring/kbuf: fix signedness in this_len calculation (bsc#1250034).
- CVE-2025-39831: fbnic: Move phylink resume out of service_task and into open/close (bsc#1249977).
- CVE-2025-39859: ptp: ocp: fix use-after-free bugs causing by ptp_ocp_watchdog (bsc#1250252).
- CVE-2025-39897: net: xilinx: axienet: Add error handling for RX metadata pointer retrieval (bsc#1250746).
- CVE-2025-39917: bpf: Fix out-of-bounds dynptr write in bpf_crypto_crypt (bsc#1250723).
- CVE-2025-39944: octeontx2-pf: Fix use-after-free bugs in otx2_sync_tstamp() (bsc#1251120).
- CVE-2025-39961: iommu/amd/pgtbl: Fix possible race while increase page table level (bsc#1251817).
- CVE-2025-39980: nexthop: Forbid FDB status change while nexthop is in a group (bsc#1252063).
- CVE-2025-39990: bpf: Check the helper function is valid in get_helper_proto (bsc#1252054).
- CVE-2025-40001: scsi: mvsas: Fix use-after-free bugs in mvs_work_queue (bsc#1252303).
- CVE-2025-40003: net: mscc: ocelot: Fix use-after-free caused by cyclic delayed work (bsc#1252301).
- CVE-2025-40006: mm/hugetlb: fix folio is still mapped when deleted (bsc#1252342).
- CVE-2025-40021: tracing: dynevent: Add a missing lockdown check on dynevent (bsc#1252681).
- CVE-2025-40024: vhost: Take a reference on the task in struct vhost_task (bsc#1252686).
- CVE-2025-40027: net/9p: fix double req put in p9_fd_cancelled (bsc#1252763).
- CVE-2025-40031: tee: fix register_shm_helper() (bsc#1252779).
- CVE-2025-40033: remoteproc: pru: Fix potential NULL pointer dereference in pru_rproc_set_ctable() (bsc#1252824).
- CVE-2025-40038: KVM: SVM: Skip fastpath emulation on VM-Exit if next RIP isn&apos;t valid (bsc#1252817).
- CVE-2025-40047: io_uring/waitid: always prune wait queue entry in io_waitid_wait() (bsc#1252790).
- CVE-2025-40053: net: dlink: handle copy_thresh allocation failure (bsc#1252808).
- CVE-2025-40055: ocfs2: fix double free in user_cluster_connect() (bsc#1252821).
- CVE-2025-40059: coresight: Fix incorrect handling for return value of devm_kzalloc (bsc#1252809).
- CVE-2025-40064: smc: Fix use-after-free in __pnet_find_base_ndev() (bsc#1252845).
- CVE-2025-40070: pps: fix warning in pps_register_cdev when register device fail (bsc#1252836).
- CVE-2025-40074: tcp: convert to dev_net_rcu() (bsc#1252794).
- CVE-2025-40075: tcp_metrics: use dst_dev_net_rcu() (bsc#1252795).
- CVE-2025-40081: perf: arm_spe: Prevent overflow in PERF_IDX2OFF() (bsc#1252776).
- CVE-2025-40083: net/sched: sch_qfq: Fix null-deref in agg_dequeue (bsc#1252912).
- CVE-2025-40086: drm/xe: Don&apos;t allow evicting of BOs in same VM in array of VM binds (bsc#1252923).
- CVE-2025-40098: ALSA: hda: cs35l41: Fix NULL pointer dereference in cs35l41_get_acpi_mute_state() (bsc#1252917).
- CVE-2025-40101: btrfs: fix memory leaks when rejecting a non SINGLE data profile without an RST (bsc#1252901).
- CVE-2025-40102: KVM: arm64: Prevent access to vCPU events before init (bsc#1252919).
- CVE-2025-40105: vfs: Don&apos;t leak disconnected dentries on umount (bsc#1252928).
- CVE-2025-40133: mptcp: Call dst_release() in mptcp_active_enable() (bsc#1253328).
- CVE-2025-40134: dm: fix NULL pointer dereference in __dm_suspend() (bsc#1253386).
- CVE-2025-40135: ipv6: use RCU in ip6_xmit() (bsc#1253342).
- CVE-2025-40139: smc: Use __sk_dst_get() and dst_dev_rcu() in in smc_clc_prfx_set() (bsc#1253409).
- CVE-2025-40149: tls: Use __sk_dst_get() and dst_dev_rcu() in get_netdev_for_sock() (bsc#1253355).
- CVE-2025-40153: mm: hugetlb: avoid soft lockup when mprotect to large memory area (bsc#1253408).
- CVE-2025-40157: EDAC/i10nm: Skip DIMM enumeration on a disabled memory controller (bsc#1253423).
- CVE-2025-40158: ipv6: use RCU in ip6_output() (bsc#1253402).
- CVE-2025-40159: xsk: Harden userspace-supplied xdp_desc validation (bsc#1253403).
- CVE-2025-40168: smc: Use __sk_dst_get() and dst_dev_rcu() in smc_clc_prfx_match() (bsc#1253427).
- CVE-2025-40169: bpf: Reject negative offsets for ALU ops (bsc#1253416).
- CVE-2025-40173: net/ip6_tunnel: Prevent perpetual tunnel growth (bsc#1253421).
- CVE-2025-40175: idpf: cleanup remaining SKBs in PTP flows (bsc#1253426).
- CVE-2025-40176: tls: wait for pending async decryptions if tls_strp_msg_hold fails (bsc#1253425).
- CVE-2025-40178: pid: Add a judgment for ns null in pid_nr_ns (bsc#1253463).
- CVE-2025-40185: ice: ice_adapter: release xa entry on adapter allocation failure (bsc#1253394).
- CVE-2025-40201: kernel/sys.c: fix the racy usage of task_lock(tsk-&gt;group_leader) in sys_prlimit64() paths (bsc#1253455).
- CVE-2025-40203: listmount: don&apos;t call path_put() under namespace semaphore (bsc#1253457).

The following non security issues were fixed:

- ACPI: scan: Update honor list for RPMI System MSI (stable-fixes).
- ACPICA: Update dsmethod.c to get rid of unused variable warning (stable-fixes).
- Disable CONFIG_CPU5_WDT The cpu5wdt driver doesn&apos;t implement a
  proper watchdog interface and has many code issues. It only handles
  obscure and obsolete hardware. Stop building and supporting this driver
  (jsc#PED-14062).
- Fix &quot;drm/xe: Don&apos;t allow evicting of BOs in same VM in array of VM binds&quot; (bsc#1252923)
- KVM: SVM: Delete IRTE link from previous vCPU before setting new IRTE (git-fixes).
- KVM: SVM: Delete IRTE link from previous vCPU irrespective of new routing (git-fixes).
- KVM: SVM: Mark VMCB_LBR dirty when MSR_IA32_DEBUGCTLMSR is updated (git-fixes).
- KVM: s390: improve interrupt cpu for wakeup (bsc#1235463).
- KVM: s390: kABI backport for &apos;last_sleep_cpu&apos; (bsc#1252352).
- KVM: x86/mmu: Return -EAGAIN if userspace deletes/moves memslot during prefault (git-fixes).
- PCI/ERR: Update device error_state already after reset (stable-fixes).
- PM: EM: Slightly reduce em_check_capacity_update() overhead (stable-fixes).
- Revert &quot;net/mlx5e: Update and set Xon/Xoff upon MTU set&quot; (git-fixes).
- Revert &quot;net/mlx5e: Update and set Xon/Xoff upon port speed set&quot; (git-fixes).
- Update config files: enable zstd module decompression (jsc#PED-14115).
- bpf/selftests: Fix test_tcpnotify_user (bsc#1253635).
- btrfs: do not clear read-only when adding sprout device (bsc#1253238).
- btrfs: do not update last_log_commit when logging inode due to a new name (git-fixes).
- dm: fix queue start/stop imbalance under suspend/load/resume races (bsc#1253386)
- drm/amd/display: Add AVI infoframe copy in copy_stream_update_to_stream (stable-fixes).
- drm/amd/display: update color on atomic commit time (stable-fixes).
- drm/amd/display: update dpp/disp clock from smu clock table (stable-fixes).
- drm/radeon: delete radeon_fence_process in is_signaled, no deadlock (stable-fixes).
- hwmon: (lenovo-ec-sensors) Update P8 supprt (stable-fixes).
- media: amphion: Delete v4l2_fh synchronously in .release() (stable-fixes).
- mount: handle NULL values in mnt_ns_release() (bsc#1254308)
- net/smc: Remove validation of reserved bits in CLC Decline (bsc#1252357).
- net: phy: move realtek PHY driver to its own subdirectory (jsc#PED-14353).
- net: phy: realtek: add defines for shadowed c45 standard registers (jsc#PED-14353).
- net: phy: realtek: add helper RTL822X_VND2_C22_REG (jsc#PED-14353).
- net: phy: realtek: change order of calls in C22 read_status() (jsc#PED-14353).
- net: phy: realtek: clear 1000Base-T link partner advertisement (jsc#PED-14353).
- net: phy: realtek: improve mmd register access for internal PHY&apos;s (jsc#PED-14353).
- net: phy: realtek: read duplex and gbit master from PHYSR register (jsc#PED-14353).
- net: phy: realtek: switch from paged to MMD ops in rtl822x functions (jsc#PED-14353).
- net: phy: realtek: use string choices helpers (jsc#PED-14353).
- net: xilinx: axienet: Fix IRQ coalescing packet count overflow (bsc#1250746)
- net: xilinx: axienet: Fix RX skb ring management in DMAengine mode (bsc#1250746)
- net: xilinx: axienet: Fix Tx skb circular buffer occupancy check in dmaengine xmit (bsc#1250746)
- nvmet-auth: update sc_c in host response (git-fixes bsc#1249397).
- nvmet-auth: update sc_c in target host hash calculation (git-fixes).
- perf list: Add IBM z17 event descriptions (jsc#PED-13611).
- platform/x86:intel/pmc: Update Arrow Lake telemetry GUID (git-fixes).
- powercap: intel_rapl: Add support for Panther Lake platform (jsc#PED-13949).
- pwm: pca9685: Use bulk write to atomicially update registers (stable-fixes).
- r8169: add PHY c45 ops for MDIO_MMD_VENDOR2 registers (jsc#PED-14353).
- r8169: add support for Intel Killer E5000 (jsc#PED-14353).
- r8169: add support for RTL8125BP rev.b (jsc#PED-14353).
- r8169: add support for RTL8125D rev.b (jsc#PED-14353).
- r8169: adjust version numbering for RTL8126 (jsc#PED-14353).
- r8169: align RTL8125 EEE config with vendor driver (jsc#PED-14353).
- r8169: align RTL8125/RTL8126 PHY config with vendor driver (jsc#PED-14353).
- r8169: align RTL8126 EEE config with vendor driver (jsc#PED-14353).
- r8169: align WAKE_PHY handling with r8125/r8126 vendor drivers (jsc#PED-14353).
- r8169: avoid duplicated messages if loading firmware fails and switch to warn level (jsc#PED-14353).
- r8169: don&apos;t take RTNL lock in rtl_task() (jsc#PED-14353).
- r8169: enable EEE at 2.5G per default on RTL8125B (jsc#PED-14353).
- r8169: enable RTL8168H/RTL8168EP/RTL8168FP ASPM support (jsc#PED-14353).
- r8169: fix inconsistent indenting in rtl8169_get_eth_mac_stats (jsc#PED-14353).
- r8169: implement additional ethtool stats ops (jsc#PED-14353).
- r8169: improve __rtl8169_set_wol (jsc#PED-14353).
- r8169: improve initialization of RSS registers on RTL8125/RTL8126 (jsc#PED-14353).
- r8169: improve rtl_set_d3_pll_down (jsc#PED-14353).
- r8169: increase max jumbo packet size on RTL8125/RTL8126 (jsc#PED-14353).
- r8169: remove leftover locks after reverted change (jsc#PED-14353).
- r8169: remove original workaround for RTL8125 broken rx issue (jsc#PED-14353).
- r8169: remove rtl_dash_loop_wait_high/low (jsc#PED-14353).
- r8169: remove support for chip version 11 (jsc#PED-14353).
- r8169: remove unused flag RTL_FLAG_TASK_RESET_NO_QUEUE_WAKE (jsc#PED-14353).
- r8169: replace custom flag with disable_work() et al (jsc#PED-14353).
- r8169: switch away from deprecated pcim_iomap_table (jsc#PED-14353).
- r8169: use helper r8169_mod_reg8_cond to simplify rtl_jumbo_config (jsc#PED-14353).
- ring-buffer: Update pages_touched to reflect persistent buffer content (git-fixes).
- s390/mm: Fix __ptep_rdp() inline assembly (bsc#1253643).
- sched/fair: Get rid of sched_domains_curr_level hack for tl-&gt;cpumask() (bsc#1246843).
- sched/fair: Have SD_SERIALIZE affect newidle balancing (bsc#1248792).
- sched/fair: Proportional newidle balance (bsc#1248792).
- sched/fair: Proportional newidle balance -KABI (bsc#1248792).
- sched/fair: Revert max_newidle_lb_cost bump (bsc#1248792).
- sched/fair: Skip sched_balance_running cmpxchg when balance is not due (bsc#1248792).
- sched/fair: Small cleanup to sched_balance_newidle() (bsc#1248792).
- sched/fair: Small cleanup to update_newidle_cost() (bsc#1248792).
- scsi: lpfc: Add capability to register Platform Name ID to fabric (bsc#1254119).
- scsi: lpfc: Allow support for BB credit recovery in point-to-point topology (bsc#1254119).
- scsi: lpfc: Ensure unregistration of rpis for received PLOGIs (bsc#1254119).
- scsi: lpfc: Fix leaked ndlp krefs when in point-to-point topology (bsc#1254119).
- scsi: lpfc: Fix reusing an ndlp that is marked NLP_DROPPED during FLOGI (bsc#1254119).
- scsi: lpfc: Modify kref handling for Fabric Controller ndlps (bsc#1254119).
- scsi: lpfc: Remove redundant NULL ptr assignment in lpfc_els_free_iocb() (bsc#1254119).
- scsi: lpfc: Revise discovery related function headers and comments (bsc#1254119).
- scsi: lpfc: Update lpfc version to 14.4.0.12 (bsc#1254119).
- scsi: lpfc: Update various NPIV diagnostic log messaging (bsc#1254119).
- selftests/run_kselftest.sh: Add --skip argument option (bsc#1254221).
- smpboot: introduce SDTL_INIT() helper to tidy sched topology setup (bsc#1246843).
- soc/tegra: fuse: speedo-tegra210: Update speedo IDs (git-fixes).
- spi: tegra210-quad: Check hardware status on timeout (bsc#1253155)
- spi: tegra210-quad: Fix timeout handling (bsc#1253155)
- spi: tegra210-quad: Refactor error handling into helper functions (bsc#1253155)
- spi: tegra210-quad: Update dummy sequence configuration (git-fixes)
- tcp_bpf: Call sk_msg_free() when tcp_bpf_send_verdict() fails to allocate psock-&gt;cork (bsc#1250705).
- wifi: ath11k: Add quirk entries for Thinkpad T14s Gen3 AMD (bsc#1254181).
- wifi: mt76: do not add wcid entries to sta poll list during MCU reset (bsc#1254315).
- wifi: mt76: introduce mt792x_config_mac_addr_list routine (bsc#1254315).
- wifi: mt76: mt7925: Fix logical vs bitwise typo (bsc#1254315).
- wifi: mt76: mt7925: Remove unnecessary if-check (bsc#1254315).
- wifi: mt76: mt7925: Simplify HIF suspend handling to avoid suspend fail (bsc#1254315).
- wifi: mt76: mt7925: add EHT control support based on the CLC data (bsc#1254315).
- wifi: mt76: mt7925: add handler to hif suspend/resume event (bsc#1254315).
- wifi: mt76: mt7925: add pci restore for hibernate (bsc#1254315).
- wifi: mt76: mt7925: config the dwell time by firmware (bsc#1254315).
- wifi: mt76: mt7925: extend MCU support for testmode (bsc#1254315).
- wifi: mt76: mt7925: fix CLC command timeout when suspend/resume (bsc#1254315).
- wifi: mt76: mt7925: fix missing hdr_trans_tlv command for broadcast wtbl (bsc#1254315).
- wifi: mt76: mt7925: fix the unfinished command of regd_notifier before suspend (bsc#1254315).
- wifi: mt76: mt7925: refine the txpower initialization flow (bsc#1254315).
- wifi: mt76: mt7925: replace zero-length array with flexible-array member (bsc#1254315).
- wifi: mt76: mt7925: update the channel usage when the regd domain changed (bsc#1254315).
- wifi: mt76: mt7925e: fix too long of wifi resume time (bsc#1254315).
- x86/smpboot: avoid SMT domain attach/destroy if SMT is not enabled (bsc#1246843).
- x86/smpboot: moves x86_topology to static initialize and truncate (bsc#1246843).
- x86/smpboot: remove redundant CONFIG_SCHED_SMT (bsc#1246843).
</description><license/><source url="http://cdn.opensuse.org/distribution/leap/16.0/repo/oss/aarch64" alias="openSUSE:repo-oss"/><issue-date time_t="1766165896" text="2025-12-19T17:38:16Z"/><issue-list><issue type="bugzilla" id="1235463"><title>Partner-L3: The installation takes too long on &quot;Preparing disks&quot; process with thinksystem 940-16i raid card and u.3 nvme disk</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1235463</href></issue><issue type="bugzilla" id="1243474"><title>VUL-0: CVE-2025-37916: kernel: pds_core: remove write-after-free of client_id</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1243474</href></issue><issue type="bugzilla" id="1245193"><title>backport nvmet-loop fixes</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1245193</href></issue><issue type="bugzilla" id="1245431"><title>HugeTLB - fix unshare pages</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1245431</href></issue><issue type="bugzilla" id="1245498"><title>VUL-0: CVE-2025-38084: kernel: mm/hugetlb: unshare page tables during VMA split, not before</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1245498</href></issue><issue type="bugzilla" id="1245499"><title>VUL-0: CVE-2025-38085: kernel: mm/hugetlb: fix huge_pmd_unshare() vs GUP-fast race</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1245499</href></issue><issue type="bugzilla" id="1246328"><title>VUL-0: CVE-2025-38321: kernel: smb: Log an error when close_all_cached_dirs fails</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1246328</href></issue><issue type="bugzilla" id="1246843"><title>SLES16 RC1 [ Regression ] [ 6.12.0-160000.16-default ]: Qemu guest boot fails with continuous soft-lockups while trying to boot with 8 NUMA nodes.</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1246843</href></issue><issue type="bugzilla" id="1247500"><title>nvme over FC: kernel soft lockup on module removal</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1247500</href></issue><issue type="bugzilla" id="1248792"><title>[PBOnline POWER10 HANA2SPS086 &amp; 079] benchCppDistributedTransactionBarrier.py reports up to 40% performance deterioration on Power10</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1248792</href></issue><issue type="bugzilla" id="1249256"><title>VUL-0: CVE-2025-38728: kernel: smb3: fix for slab out of bounds on mount to ksmbd</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1249256</href></issue><issue type="bugzilla" id="1249397"><title>[NetApp SLES15 SP7 Bug]: sc_c field not updated in host response to controller challenge for secure concat</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1249397</href></issue><issue type="bugzilla" id="1249912"><title>VUL-0: CVE-2022-50253: kernel: bpf: make sure skb-&gt;len != 0 when redirecting to a tunneling device</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1249912</href></issue><issue type="bugzilla" id="1249977"><title>VUL-0: CVE-2025-39831: kernel: fbnic: Move phylink resume out of service_task and into open/close</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1249977</href></issue><issue type="bugzilla" id="1249982"><title>VUL-0: CVE-2025-39805: kernel: net: macb: fix unregister_netdev call order in macb_remove()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1249982</href></issue><issue type="bugzilla" id="1250034"><title>VUL-0: CVE-2025-39822: kernel: io_uring/kbuf: fix signedness in this_len calculation</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1250034</href></issue><issue type="bugzilla" id="1250176"><title>VUL-0: CVE-2025-39819: kernel: fs/smb: Fix inconsistent refcnt update</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1250176</href></issue><issue type="bugzilla" id="1250237"><title>nftables stack guard hit + kernel panic on synproxy in output chain</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1250237</href></issue><issue type="bugzilla" id="1250252"><title>VUL-0: CVE-2025-39859: kernel: ptp: ocp: fix use-after-free bugs causing by ptp_ocp_watchdog</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1250252</href></issue><issue type="bugzilla" id="1250705"><title>VUL-0: CVE-2025-39913: kernel: tcp_bpf: Call sk_msg_free() when tcp_bpf_send_verdict() fails to allocate psock-&gt;cork.</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1250705</href></issue><issue type="bugzilla" id="1250723"><title>VUL-0: CVE-2025-39917: kernel: bpf: Fix out-of-bounds dynptr write in bpf_crypto_crypt</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1250723</href></issue><issue type="bugzilla" id="1250746"><title>VUL-0: CVE-2025-39897: kernel: net: xilinx: axienet: Add error handling for RX metadata pointer retrieval</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1250746</href></issue><issue type="bugzilla" id="1251120"><title>VUL-0: CVE-2025-39944: kernel: octeontx2-pf: Fix use-after-free bugs in otx2_sync_tstamp()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1251120</href></issue><issue type="bugzilla" id="1251817"><title>VUL-0: CVE-2025-39961: kernel: iommu/amd/pgtbl: Fix possible race while increase page table level</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1251817</href></issue><issue type="bugzilla" id="1252054"><title>VUL-0: CVE-2025-39990: kernel: bpf: Check the helper function is valid in get_helper_proto</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252054</href></issue><issue type="bugzilla" id="1252063"><title>VUL-0: CVE-2025-39980: kernel: nexthop: Forbid FDB status change while nexthop is in a group</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252063</href></issue><issue type="bugzilla" id="1252301"><title>VUL-0: CVE-2025-40003: kernel: net: mscc: ocelot: Fix use-after-free caused by cyclic delayed work</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252301</href></issue><issue type="bugzilla" id="1252303"><title>VUL-0: CVE-2025-40001: kernel: scsi: mvsas: Fix use-after-free bugs in mvs_work_queue</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252303</href></issue><issue type="bugzilla" id="1252342"><title>VUL-0: CVE-2025-40006: kernel: mm/hugetlb: fix folio is still mapped when deleted</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252342</href></issue><issue type="bugzilla" id="1252352"><title>SLES 15 SP7 - KVM: s390: improve interrupt cpu for wakeup</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252352</href></issue><issue type="bugzilla" id="1252357"><title>SLES 16.0 - net/smc: Remove validation of reserved bits in CLC Decline msg</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252357</href></issue><issue type="bugzilla" id="1252681"><title>VUL-0: CVE-2025-40021: kernel: tracing: dynevent: Add a missing lockdown check on dynevent</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252681</href></issue><issue type="bugzilla" id="1252686"><title>VUL-0: CVE-2025-40024: kernel: vhost: Take a reference on the task in struct vhost_task.</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252686</href></issue><issue type="bugzilla" id="1252763"><title>VUL-0: CVE-2025-40027: kernel: net/9p: fix double req put in p9_fd_cancelled</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252763</href></issue><issue type="bugzilla" id="1252776"><title>VUL-0: CVE-2025-40081: kernel: perf: arm_spe: Prevent overflow in PERF_IDX2OFF()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252776</href></issue><issue type="bugzilla" id="1252779"><title>VUL-0: CVE-2025-40031: kernel: tee: fix register_shm_helper()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252779</href></issue><issue type="bugzilla" id="1252790"><title>VUL-0: CVE-2025-40047: kernel: io_uring/waitid: always prune wait queue entry in io_waitid_wait()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252790</href></issue><issue type="bugzilla" id="1252794"><title>VUL-0: CVE-2025-40074: kernel: ipv4: start using dst_dev_rcu()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252794</href></issue><issue type="bugzilla" id="1252795"><title>VUL-0: CVE-2025-40075: kernel: tcp_metrics: use dst_dev_net_rcu()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252795</href></issue><issue type="bugzilla" id="1252808"><title>VUL-0: CVE-2025-40053: kernel: net: dlink: handle copy_thresh allocation failure</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252808</href></issue><issue type="bugzilla" id="1252809"><title>VUL-0: CVE-2025-40059: kernel: coresight: Fix incorrect handling for return value of devm_kzalloc</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252809</href></issue><issue type="bugzilla" id="1252817"><title>VUL-0: CVE-2025-40038: kernel: KVM: SVM: Skip fastpath emulation on VM-Exit if next RIP isn&apos;t valid</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252817</href></issue><issue type="bugzilla" id="1252821"><title>VUL-0: CVE-2025-40055: kernel: ocfs2: fix double free in user_cluster_connect()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252821</href></issue><issue type="bugzilla" id="1252824"><title>VUL-0: CVE-2025-40033: kernel: remoteproc: pru: Fix potential NULL pointer dereference in pru_rproc_set_ctable()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252824</href></issue><issue type="bugzilla" id="1252836"><title>VUL-0: CVE-2025-40070: kernel: pps: fix warning in pps_register_cdev when register device fail</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252836</href></issue><issue type="bugzilla" id="1252845"><title>VUL-0: CVE-2025-40064: kernel: smc: Fix use-after-free in __pnet_find_base_ndev().</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252845</href></issue><issue type="bugzilla" id="1252901"><title>VUL-0: CVE-2025-40101: kernel: btrfs: fix memory leaks when rejecting a non SINGLE data profile without an RST</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252901</href></issue><issue type="bugzilla" id="1252912"><title>VUL-0: CVE-2025-40083: kernel: net/sched: sch_qfq: Fix null-deref in agg_dequeue</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252912</href></issue><issue type="bugzilla" id="1252917"><title>VUL-0: CVE-2025-40098: kernel: ALSA: hda: cs35l41: Fix NULL pointer dereference in cs35l41_get_acpi_mute_state()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252917</href></issue><issue type="bugzilla" id="1252919"><title>VUL-0: CVE-2025-40102: kernel: KVM: arm64: Prevent access to vCPU events before init</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252919</href></issue><issue type="bugzilla" id="1252923"><title>VUL-0: CVE-2025-40086: kernel: drm/xe: Don&apos;t allow evicting of BOs in same VM in array of VM binds</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252923</href></issue><issue type="bugzilla" id="1252928"><title>VUL-0: CVE-2025-40105: kernel: vfs: Don&apos;t leak disconnected dentries on umount</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1252928</href></issue><issue type="bugzilla" id="1253018"><title>VUL-0: CVE-2025-40107: kernel: can: hi311x: fix null pointer dereference when resuming from sleep before interface was enabled</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253018</href></issue><issue type="bugzilla" id="1253155"><title>Nvidia Grace: Backport: spi: tegra210-quad: Improve timeout handling under high system load</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253155</href></issue><issue type="bugzilla" id="1253176"><title>VUL-0: CVE-2025-40109: kernel: crypto: rng - Ensure set_ent is always present</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253176</href></issue><issue type="bugzilla" id="1253238"><title>SLES16.0 [6.12.0-160000.18-default] FS/BTRFS: btrfs/282 btrfs/323 generic/363 tests fails with missing kernel fix</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253238</href></issue><issue type="bugzilla" id="1253275"><title>VUL-0: CVE-2025-40110: kernel: drm/vmwgfx: Fix a null-ptr access in the cursor snooper</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253275</href></issue><issue type="bugzilla" id="1253318"><title>VUL-0: CVE-2025-40115: kernel: scsi: mpt3sas: Fix crash in transport port remove by using ioc_info()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253318</href></issue><issue type="bugzilla" id="1253324"><title>VUL-0: CVE-2025-40116: kernel: usb: host: max3421-hcd: Fix error pointer dereference in probe cleanup</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253324</href></issue><issue type="bugzilla" id="1253328"><title>VUL-0: CVE-2025-40133: kernel: mptcp: Use __sk_dst_get() and dst_dev_rcu() in mptcp_active_enable().</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253328</href></issue><issue type="bugzilla" id="1253330"><title>VUL-0: CVE-2025-40132: kernel: ASoC: Intel: sof_sdw: Prevent jump to NULL add_sidecar callback</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253330</href></issue><issue type="bugzilla" id="1253342"><title>VUL-0: CVE-2025-40135: kernel: ipv6: use RCU in ip6_xmit()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253342</href></issue><issue type="bugzilla" id="1253348"><title>VUL-0: CVE-2025-40142: kernel: ALSA: pcm: Disable bottom softirqs as part of spin_lock_irq() on PREEMPT_RT</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253348</href></issue><issue type="bugzilla" id="1253349"><title>VUL-0: CVE-2025-40140: kernel: net: usb: Remove disruptive netif_wake_queue in rtl8150_set_multicast</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253349</href></issue><issue type="bugzilla" id="1253352"><title>VUL-0: CVE-2025-40141: kernel: Bluetooth: ISO: Fix possible UAF on iso_conn_free</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253352</href></issue><issue type="bugzilla" id="1253355"><title>VUL-0: CVE-2025-40149: kernel: tls: Use __sk_dst_get() and dst_dev_rcu() in get_netdev_for_sock().</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253355</href></issue><issue type="bugzilla" id="1253360"><title>VUL-0: CVE-2025-40120: kernel: net: usb: asix: hold PM usage ref to avoid PM/MDIO + RTNL deadlock</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253360</href></issue><issue type="bugzilla" id="1253362"><title>VUL-0: CVE-2025-40111: kernel: drm/vmwgfx: Fix Use-after-free in validation</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253362</href></issue><issue type="bugzilla" id="1253363"><title>VUL-0: CVE-2025-40118: kernel: scsi: pm80xx: Fix array-index-out-of-of-bounds on rmmod</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253363</href></issue><issue type="bugzilla" id="1253367"><title>VUL-0: CVE-2025-40121: kernel: ASoC: Intel: bytcr_rt5651: Fix invalid quirk input mapping</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253367</href></issue><issue type="bugzilla" id="1253369"><title>VUL-0: CVE-2025-40127: kernel: hwrng: ks-sa - fix division by zero in ks_sa_rng_init</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253369</href></issue><issue type="bugzilla" id="1253386"><title>VUL-0: CVE-2025-40134: kernel: dm: fix NULL pointer dereference in __dm_suspend()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253386</href></issue><issue type="bugzilla" id="1253394"><title>VUL-0: CVE-2025-40185: kernel: ice: ice_adapter: release xa entry on adapter allocation failure</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253394</href></issue><issue type="bugzilla" id="1253395"><title>VUL-0: CVE-2025-40207: kernel: media: v4l2-subdev: Fix alloc failure check in v4l2_subdev_call_state_try()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253395</href></issue><issue type="bugzilla" id="1253402"><title>VUL-0: CVE-2025-40158: kernel: ipv6: use RCU in ip6_output()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253402</href></issue><issue type="bugzilla" id="1253403"><title>VUL-0: CVE-2025-40159: kernel: xsk: Harden userspace-supplied xdp_desc validation</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253403</href></issue><issue type="bugzilla" id="1253405"><title>VUL-0: CVE-2025-40165: kernel: media: nxp: imx8-isi: m2m: Fix streaming cleanup on release</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253405</href></issue><issue type="bugzilla" id="1253407"><title>VUL-0: CVE-2025-40164: kernel: usbnet: Fix using smp_processor_id() in preemptible code warnings</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253407</href></issue><issue type="bugzilla" id="1253408"><title>VUL-0: CVE-2025-40153: kernel: mm: hugetlb: avoid soft lockup when mprotect to large memory area</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253408</href></issue><issue type="bugzilla" id="1253409"><title>VUL-0: CVE-2025-40139: kernel: smc: Use __sk_dst_get() and dst_dev_rcu() in in smc_clc_prfx_set().</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253409</href></issue><issue type="bugzilla" id="1253410"><title>VUL-0: CVE-2025-40161: kernel: mailbox: zynqmp-ipi: Fix SGI cleanup on unbind</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253410</href></issue><issue type="bugzilla" id="1253412"><title>VUL-0: CVE-2025-40171: kernel: nvmet-fc: move lsop put work to nvmet_fc_ls_req_op</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253412</href></issue><issue type="bugzilla" id="1253416"><title>VUL-0: CVE-2025-40169: kernel: bpf: Reject negative offsets for ALU ops</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253416</href></issue><issue type="bugzilla" id="1253421"><title>VUL-0: CVE-2025-40173: kernel: net/ip6_tunnel: Prevent perpetual tunnel growth</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253421</href></issue><issue type="bugzilla" id="1253422"><title>VUL-0: CVE-2025-40162: kernel: ASoC: amd/sdw_utils: avoid NULL deref when devm_kasprintf() fails</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253422</href></issue><issue type="bugzilla" id="1253423"><title>VUL-0: CVE-2025-40157: kernel: EDAC/i10nm: Skip DIMM enumeration on a disabled memory controller</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253423</href></issue><issue type="bugzilla" id="1253424"><title>VUL-0: CVE-2025-40172: kernel: accel/qaic: Treat remaining == 0 as error in find_and_map_user_pages()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253424</href></issue><issue type="bugzilla" id="1253425"><title>VUL-0: CVE-2025-40176: kernel: tls: wait for pending async decryptions if tls_strp_msg_hold fails</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253425</href></issue><issue type="bugzilla" id="1253426"><title>VUL-0: CVE-2025-40175: kernel: idpf: cleanup remaining SKBs in PTP flows</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253426</href></issue><issue type="bugzilla" id="1253427"><title>VUL-0: CVE-2025-40168: kernel: smc: Use __sk_dst_get() and dst_dev_rcu() in smc_clc_prfx_match().</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253427</href></issue><issue type="bugzilla" id="1253428"><title>VUL-0: CVE-2025-40156: kernel: PM / devfreq: mtk-cci: Fix potential error pointer dereference in probe()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253428</href></issue><issue type="bugzilla" id="1253431"><title>VUL-0: CVE-2025-40154: kernel: ASoC: Intel: bytcr_rt5640: Fix invalid quirk input mapping</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253431</href></issue><issue type="bugzilla" id="1253433"><title>VUL-0: CVE-2025-40166: kernel: drm/xe/guc: Check GuC running state before deregistering exec queue</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253433</href></issue><issue type="bugzilla" id="1253436"><title>VUL-0: CVE-2025-40204: kernel: sctp: Fix MAC comparison to be constant-time</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253436</href></issue><issue type="bugzilla" id="1253438"><title>VUL-0: CVE-2025-40186: kernel: tcp: Don&apos;t call reqsk_fastopen_remove() in tcp_conn_request().</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253438</href></issue><issue type="bugzilla" id="1253440"><title>VUL-0: CVE-2025-40180: kernel: mailbox: zynqmp-ipi: Fix out-of-bounds access in mailbox cleanup loop</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253440</href></issue><issue type="bugzilla" id="1253441"><title>VUL-0: CVE-2025-40183: kernel: bpf: Fix metadata_dst leak __bpf_redirect_neigh_v{4,6}</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253441</href></issue><issue type="bugzilla" id="1253443"><title>VUL-0: CVE-2025-40177: kernel: accel/qaic: Fix bootlog initialization ordering</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253443</href></issue><issue type="bugzilla" id="1253445"><title>VUL-0: CVE-2025-40194: kernel: cpufreq: intel_pstate: Fix object lifecycle issue in update_qos_request()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253445</href></issue><issue type="bugzilla" id="1253448"><title>VUL-0: CVE-2025-40200: kernel: Squashfs: reject negative file sizes in squashfs_read_inode()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253448</href></issue><issue type="bugzilla" id="1253449"><title>VUL-0: CVE-2025-40188: kernel: pwm: berlin: Fix wrong register in suspend/resume</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253449</href></issue><issue type="bugzilla" id="1253450"><title>VUL-0: CVE-2025-40197: kernel: media: mc: Clear minor number before put device</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253450</href></issue><issue type="bugzilla" id="1253451"><title>VUL-0: CVE-2025-40202: kernel: ipmi: Rework user message limit handling</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253451</href></issue><issue type="bugzilla" id="1253453"><title>VUL-0: CVE-2025-40198: kernel: ext4: avoid potential buffer over-read in parse_apply_sb_mount_options()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253453</href></issue><issue type="bugzilla" id="1253455"><title>VUL-0: CVE-2025-40201: kernel: kernel/sys.c: fix the racy usage of task_lock(tsk-&gt;group_leader) in sys_prlimit64() paths</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253455</href></issue><issue type="bugzilla" id="1253456"><title>VUL-0: CVE-2025-40205: kernel: btrfs: avoid potential out-of-bounds in btrfs_encode_fh()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253456</href></issue><issue type="bugzilla" id="1253457"><title>VUL-0: CVE-2025-40203: kernel: listmount: don&apos;t call path_put() under namespace semaphore</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253457</href></issue><issue type="bugzilla" id="1253463"><title>VUL-0: CVE-2025-40178: kernel: pid: Add a check for ns is null in pid_nr_ns</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253463</href></issue><issue type="bugzilla" id="1253472"><title>VUL-0: CVE-2025-40129: kernel: sunrpc: fix null pointer dereference on zero-length checksum</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253472</href></issue><issue type="bugzilla" id="1253622"><title>VUL-0: CVE-2025-40192: kernel: Revert &quot;ipmi: fix msg stack when IPMI is disconnected&quot;</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253622</href></issue><issue type="bugzilla" id="1253624"><title>VUL-0: CVE-2025-40196: kernel: fs: quota: create dedicated workqueue for quota_release_work</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253624</href></issue><issue type="bugzilla" id="1253635"><title>selftests: bpf: test_tcpnotify_user segfaults</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253635</href></issue><issue type="bugzilla" id="1253643"><title>SLES 16 SP0 - s390/mm: Fix __ptep_rdp() inline assembly</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253643</href></issue><issue type="bugzilla" id="1253647"><title>VUL-0: CVE-2025-40187: kernel: net/sctp: fix a null dereference in sctp_disposition sctp_sf_do_5_1D_ce()</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1253647</href></issue><issue type="bugzilla" id="1254119"><title>Update Broadcom Emulex lpfc driver for SL-16.1  to 14.4.0.12</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1254119</href></issue><issue type="bugzilla" id="1254181"><title>ath11k WiFi suspend/resume breakage on Lenovo Thinkpad T14s Gen 3 AMD</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1254181</href></issue><issue type="bugzilla" id="1254221"><title>selftests: run_kselftest.sh: openQA test fails when skipping tests from a large collection</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1254221</href></issue><issue type="bugzilla" id="1254308"><title>SLE-16.0 KOTD: NULL pointer dereference in listmount() syscall</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1254308</href></issue><issue type="bugzilla" id="1254315"><title>Lenovo Thinkpad P14s doesn&apos;t sleep</title><href>https://bugzilla.suse.com/show_bug.cgi?id=1254315</href></issue><issue type="cve" id="CVE-2022-50253"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-50253</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-50253</href></issue><issue type="cve" id="CVE-2025-37916"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-37916</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-37916</href></issue><issue type="cve" id="CVE-2025-38084"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-38084</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-38084</href></issue><issue type="cve" id="CVE-2025-38085"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-38085</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-38085</href></issue><issue type="cve" id="CVE-2025-38321"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-38321</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-38321</href></issue><issue type="cve" id="CVE-2025-38728"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-38728</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-38728</href></issue><issue type="cve" id="CVE-2025-39805"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39805</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39805</href></issue><issue type="cve" id="CVE-2025-39819"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39819</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39819</href></issue><issue type="cve" id="CVE-2025-39822"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39822</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39822</href></issue><issue type="cve" id="CVE-2025-39831"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39831</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39831</href></issue><issue type="cve" id="CVE-2025-39859"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39859</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39859</href></issue><issue type="cve" id="CVE-2025-39897"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39897</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39897</href></issue><issue type="cve" id="CVE-2025-39917"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39917</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39917</href></issue><issue type="cve" id="CVE-2025-39944"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39944</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39944</href></issue><issue type="cve" id="CVE-2025-39961"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39961</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39961</href></issue><issue type="cve" id="CVE-2025-39980"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39980</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39980</href></issue><issue type="cve" id="CVE-2025-39990"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39990</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-39990</href></issue><issue type="cve" id="CVE-2025-40001"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40001</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40001</href></issue><issue type="cve" id="CVE-2025-40003"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40003</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40003</href></issue><issue type="cve" id="CVE-2025-40006"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40006</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40006</href></issue><issue type="cve" id="CVE-2025-40021"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40021</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40021</href></issue><issue type="cve" id="CVE-2025-40024"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40024</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40024</href></issue><issue type="cve" id="CVE-2025-40027"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40027</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40027</href></issue><issue type="cve" id="CVE-2025-40031"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40031</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40031</href></issue><issue type="cve" id="CVE-2025-40033"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40033</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40033</href></issue><issue type="cve" id="CVE-2025-40038"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40038</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40038</href></issue><issue type="cve" id="CVE-2025-40047"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40047</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40047</href></issue><issue type="cve" id="CVE-2025-40053"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40053</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40053</href></issue><issue type="cve" id="CVE-2025-40055"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40055</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40055</href></issue><issue type="cve" id="CVE-2025-40059"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40059</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40059</href></issue><issue type="cve" id="CVE-2025-40064"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40064</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40064</href></issue><issue type="cve" id="CVE-2025-40070"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40070</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40070</href></issue><issue type="cve" id="CVE-2025-40074"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40074</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40074</href></issue><issue type="cve" id="CVE-2025-40075"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40075</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40075</href></issue><issue type="cve" id="CVE-2025-40081"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40081</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40081</href></issue><issue type="cve" id="CVE-2025-40083"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40083</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40083</href></issue><issue type="cve" id="CVE-2025-40086"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40086</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40086</href></issue><issue type="cve" id="CVE-2025-40098"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40098</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40098</href></issue><issue type="cve" id="CVE-2025-40101"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40101</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40101</href></issue><issue type="cve" id="CVE-2025-40102"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40102</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40102</href></issue><issue type="cve" id="CVE-2025-40105"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40105</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40105</href></issue><issue type="cve" id="CVE-2025-40107"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40107</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40107</href></issue><issue type="cve" id="CVE-2025-40109"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40109</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40109</href></issue><issue type="cve" id="CVE-2025-40110"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40110</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40110</href></issue><issue type="cve" id="CVE-2025-40111"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40111</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40111</href></issue><issue type="cve" id="CVE-2025-40115"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40115</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40115</href></issue><issue type="cve" id="CVE-2025-40116"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40116</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40116</href></issue><issue type="cve" id="CVE-2025-40118"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40118</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40118</href></issue><issue type="cve" id="CVE-2025-40120"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40120</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40120</href></issue><issue type="cve" id="CVE-2025-40121"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40121</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40121</href></issue><issue type="cve" id="CVE-2025-40127"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40127</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40127</href></issue><issue type="cve" id="CVE-2025-40129"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40129</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40129</href></issue><issue type="cve" id="CVE-2025-40132"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40132</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40132</href></issue><issue type="cve" id="CVE-2025-40133"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40133</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40133</href></issue><issue type="cve" id="CVE-2025-40134"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40134</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40134</href></issue><issue type="cve" id="CVE-2025-40135"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40135</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40135</href></issue><issue type="cve" id="CVE-2025-40139"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40139</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40139</href></issue><issue type="cve" id="CVE-2025-40140"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40140</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40140</href></issue><issue type="cve" id="CVE-2025-40141"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40141</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40141</href></issue><issue type="cve" id="CVE-2025-40142"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40142</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40142</href></issue><issue type="cve" id="CVE-2025-40149"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40149</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40149</href></issue><issue type="cve" id="CVE-2025-40153"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40153</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40153</href></issue><issue type="cve" id="CVE-2025-40154"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40154</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40154</href></issue><issue type="cve" id="CVE-2025-40156"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40156</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40156</href></issue><issue type="cve" id="CVE-2025-40157"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40157</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40157</href></issue><issue type="cve" id="CVE-2025-40158"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40158</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40158</href></issue><issue type="cve" id="CVE-2025-40159"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40159</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40159</href></issue><issue type="cve" id="CVE-2025-40161"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40161</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40161</href></issue><issue type="cve" id="CVE-2025-40162"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40162</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40162</href></issue><issue type="cve" id="CVE-2025-40164"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40164</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40164</href></issue><issue type="cve" id="CVE-2025-40165"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40165</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40165</href></issue><issue type="cve" id="CVE-2025-40166"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40166</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40166</href></issue><issue type="cve" id="CVE-2025-40168"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40168</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40168</href></issue><issue type="cve" id="CVE-2025-40169"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40169</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40169</href></issue><issue type="cve" id="CVE-2025-40171"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40171</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40171</href></issue><issue type="cve" id="CVE-2025-40172"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40172</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40172</href></issue><issue type="cve" id="CVE-2025-40173"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40173</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40173</href></issue><issue type="cve" id="CVE-2025-40175"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40175</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40175</href></issue><issue type="cve" id="CVE-2025-40176"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40176</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40176</href></issue><issue type="cve" id="CVE-2025-40177"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40177</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40177</href></issue><issue type="cve" id="CVE-2025-40178"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40178</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40178</href></issue><issue type="cve" id="CVE-2025-40180"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40180</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40180</href></issue><issue type="cve" id="CVE-2025-40183"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40183</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40183</href></issue><issue type="cve" id="CVE-2025-40185"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40185</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40185</href></issue><issue type="cve" id="CVE-2025-40186"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40186</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40186</href></issue><issue type="cve" id="CVE-2025-40187"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40187</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40187</href></issue><issue type="cve" id="CVE-2025-40188"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40188</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40188</href></issue><issue type="cve" id="CVE-2025-40192"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40192</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40192</href></issue><issue type="cve" id="CVE-2025-40194"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40194</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40194</href></issue><issue type="cve" id="CVE-2025-40196"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40196</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40196</href></issue><issue type="cve" id="CVE-2025-40197"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40197</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40197</href></issue><issue type="cve" id="CVE-2025-40198"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40198</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40198</href></issue><issue type="cve" id="CVE-2025-40200"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40200</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40200</href></issue><issue type="cve" id="CVE-2025-40201"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40201</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40201</href></issue><issue type="cve" id="CVE-2025-40202"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40202</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40202</href></issue><issue type="cve" id="CVE-2025-40203"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40203</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40203</href></issue><issue type="cve" id="CVE-2025-40204"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40204</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40204</href></issue><issue type="cve" id="CVE-2025-40205"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40205</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40205</href></issue><issue type="cve" id="CVE-2025-40206"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40206</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40206</href></issue><issue type="cve" id="CVE-2025-40207"><title>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40207</title><href>http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-40207</href></issue><issue type="jira" id="PED-13611"><title>https://jira.suse.com/browse/PED-13611</title><href>https://jira.suse.com/browse/PED-13611</href></issue><issue type="jira" id="PED-13949"><title>https://jira.suse.com/browse/PED-13949</title><href>https://jira.suse.com/browse/PED-13949</href></issue><issue type="jira" id="PED-14062"><title>https://jira.suse.com/browse/PED-14062</title><href>https://jira.suse.com/browse/PED-14062</href></issue><issue type="jira" id="PED-14115"><title>https://jira.suse.com/browse/PED-14115</title><href>https://jira.suse.com/browse/PED-14115</href></issue><issue type="jira" id="PED-14353"><title>https://jira.suse.com/browse/PED-14353</title><href>https://jira.suse.com/browse/PED-14353</href></issue></issue-list></update></update-list>
</update-status>
</stream>`

func TestZypperManager_getInstalledPackagesWithCancel(t *testing.T) {

	zypper := ZypperManager{}

	if !zypper.IsAvailable() {
		t.Skip("zypper is not available, skipping test")
	}

	ctx := t.Context()
	packages, err := zypper.ListInstalledPackages(ctx)
	if err != nil {
		t.Fatalf("Error getting installed packages: %v", err)
	}

	if len(packages) == 0 {
		t.Fatal("Expected at least one installed package, got zero")
	}
}

func TestParseZypperUpgradeOutput(t *testing.T) {
	z := ZypperManager{}
	pkgs, err := z.parseZypperUpgradeOutput(zypperAvailableUpdatesOutput)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := []PackageUpdate{
		{
			Name:             "kernel-default",
			CurrentVersion:   "6.12.0-160000.7.1",
			AvailableVersion: "6.12.0-160000.8.1",
			IsSecurityUpdate: false,
		},
		{
			Name:             "kernel-default-extra",
			CurrentVersion:   "6.12.0-160000.7.1",
			AvailableVersion: "6.12.0-160000.8.1",
			IsSecurityUpdate: false,
		},
		{
			Name:             "kernel-default-optional",
			CurrentVersion:   "6.12.0-160000.7.1",
			AvailableVersion: "6.12.0-160000.8.1",
			IsSecurityUpdate: false,
		},
		{
			Name:             "qemu-guest-agent",
			CurrentVersion:   "10.0.4-160000.1.1",
			AvailableVersion: "10.0.7-160000.1.1",
			IsSecurityUpdate: false,
		},
	}

	if len(pkgs) != len(expected) {
		t.Fatalf("Expected %d packages, got %d", len(expected), len(pkgs))
	}

	for i, pkg := range pkgs {
		if pkg.Name != expected[i].Name {
			t.Errorf("Package %d: expected Name %q, got %q", i, expected[i].Name, pkg.Name)
		}
		if pkg.CurrentVersion != expected[i].CurrentVersion {
			t.Errorf("Package %d: expected CurrentVersion %q, got %q", i, expected[i].CurrentVersion, pkg.CurrentVersion)
		}
		if pkg.AvailableVersion != expected[i].AvailableVersion {
			t.Errorf("Package %d: expected AvailableVersion %q, got %q", i, expected[i].AvailableVersion, pkg.AvailableVersion)
		}
		if pkg.IsSecurityUpdate != expected[i].IsSecurityUpdate {
			t.Errorf("Package %d: expected IsSecurityUpdate %v, got %v", i, expected[i].IsSecurityUpdate, pkg.IsSecurityUpdate)
		}
	}
}

func TestParseZypperUpgradeOutput_InvalidXML(t *testing.T) {
	z := ZypperManager{}
	invalidOutput := `<stream><update-status><update-list><update></update-list></update-status></stream`
	_, err := z.parseZypperUpgradeOutput(invalidOutput)
	if err == nil {
		t.Error("Expected error for invalid XML, got nil")
	}
}

func TestParseZypperSecurityPatchesOutput(t *testing.T) {
	z := ZypperManager{}

	expected := []PackageUpdate{
		{
			Name:             "openSUSE-Leap-16.0-112",
			CurrentVersion:   "",
			AvailableVersion: "1",
			IsSecurityUpdate: true,
			IsPatch:          true,
		},
		{
			Name:             "openSUSE-Leap-16.0-113",
			CurrentVersion:   "",
			AvailableVersion: "1",
			IsSecurityUpdate: true,
			IsPatch:          true,
		},
		{
			Name:             "openSUSE-Leap-16.0-112",
			CurrentVersion:   "",
			AvailableVersion: "1",
			IsSecurityUpdate: true,
			IsPatch:          true,
		},
		{
			Name:             "openSUSE-Leap-16.0-113",
			CurrentVersion:   "",
			AvailableVersion: "1",
			IsSecurityUpdate: true,
			IsPatch:          true,
		},
	}

	result, err := z.parseZypperSecurityPatchesOutput(zypperAvailableSecurityPatchesOutput)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != len(expected) {
		t.Fatalf("Expected %d packages, got %d", len(expected), len(result))
	}

	for i, pkg := range result {
		if pkg.Name != expected[i].Name {
			t.Errorf("Package %d: expected Name %q, got %q", i, expected[i].Name, pkg.Name)
		}
		if pkg.CurrentVersion != expected[i].CurrentVersion {
			t.Errorf("Package %d: expected CurrentVersion %q, got %q", i, expected[i].CurrentVersion, pkg.CurrentVersion)
		}
		if pkg.AvailableVersion != expected[i].AvailableVersion {
			t.Errorf("Package %d: expected AvailableVersion %q, got %q", i, expected[i].AvailableVersion, pkg.AvailableVersion)
		}
		if pkg.IsSecurityUpdate != expected[i].IsSecurityUpdate {
			t.Errorf("Package %d: expected IsSecurityUpdate %v, got %v", i, expected[i].IsSecurityUpdate, pkg.IsSecurityUpdate)
		}
		if pkg.IsPatch != expected[i].IsPatch {
			t.Errorf("Package %d: expected IsPatch %v, got %v", i, expected[i].IsPatch, pkg.IsPatch)
		}
	}

}

func TestParseZypperSecurityPatchesOutput_EmptyList(t *testing.T) {
	z := ZypperManager{}

	xmlInput := `
<stream>
  <update-status>
	<update-list>
	</update-list>
  </update-status>
</stream>
`
	result, err := z.parseZypperSecurityPatchesOutput(xmlInput)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty list, got %+v", result)
	}

}

func TestParseZypperSecurityPatchesOutput_InvalidXML(t *testing.T) {
	z := ZypperManager{}

	xmlInput := `<stream><update-status><update-list><update></update-list></update-status></stream>`

	_, err := z.parseZypperSecurityPatchesOutput(xmlInput)
	if err == nil {
		t.Error("Expected error for invalid XML, got nil")
	}
}
