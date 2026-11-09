# openITCOCKPIT Monitoring Agent 3
Cross-Platform Monitoring Agent for openITCOCKPIT written in Go

[![Go Report Card](https://goreportcard.com/badge/github.com/openITCOCKPIT/openitcockpit-agent-go)](https://goreportcard.com/report/github.com/openITCOCKPIT/openitcockpit-agent-go)
[![Build Status](https://drone.openitcockpit.io/buildStatus/icon?job=openitcockpit-agent-go%2Fmain)](https://drone.openitcockpit.io/job/openitcockpit-agent-go/job/main/)

## Table of contents
* [Supported operating systems](#supported-operating-systems)
* [Installation](#installation)
  - [Debian and Ubuntu](#debian-and-ubuntu)
  - [Red Hat Linux / CentOS / openSUSE](#red-hat-linux--centos--opensuse)
  - [Arch Linux](#arch-linux)
  - [Windows](#windows)
  - [macOS](#macos)
  - [FreeBSD](#freebsd)
* [Supported Platforms](#supported-platforms)
* [Full documentation](#full-documentation)
* [License](#license)

## Supported operating systems

* Microsoft Windows Server 2016 or newer
* Microsoft Windows 10 or newer
* Apple macOS 26 Tahoe or newer for Apple Silicon Macs
  * macOS 15 Sequoia or newer for Intel Macs
* Linux (Everything from Debian 10 (Buster) / AlmaLinux 8 and newer should work fine)
* FreeBSD 15

### Maybe supported operating systems

> [!NOTE]
> As we move forward we have to drop support for older operating systems from time to time. This does not mean that the openITCOCKPIT Agent can not be used on older operating systems - we simply do not test theme anymore.

If the latest version is not working on your operating system, we recommend to go to the [releases page](https://github.com/openITCOCKPIT/openitcockpit-agent-go/releases) and select a version that fits to the era of your operating system.
Otherwise you can try build the agent from source.

The following operating systems will may work with the openITCOCKPIT Agent up to version 3.1.0:

* Microsoft Windows Server 2012
* Microsoft Windows 8 or newer
* Apple macOS 10.14 Mojave or newer (Intel / Apple Silicon)
* Linux (Everything from Debian 6.0 (Squeeze) / CentOS 6.6 and newer should work fine)

Please notice: Due to old versions of PowerShell on Windows 7 / Windows Server 2008 R2 you need to add add the required Firewall rules manually to Windows Firewall.
Windows 7 / Windows Server 2008 R2 is official not supported by the Agent - even if it probably works.

## Update from Agent 1.x to 3.x
Please see the [update guide from the documentation](https://github.com/openITCOCKPIT/openitcockpit-agent-go/wiki/Update-from-Agent-1.x-to-3.x) for details.

## Requirements
* openITCOCKPIT Version >= 4.2

## Installation

Please visit the [release page](https://github.com/openITCOCKPIT/openitcockpit-agent-go/releases) to download the latest or older versions.

### Debian and Ubuntu

#### Using the repository

```
curl https://packages.openitcockpit.io/repokey.txt | sudo apt-key add

echo "deb https://packages.openitcockpit.io/openitcockpit-agent/deb/stable deb main"  | sudo tee /etc/apt/sources.list.d/openitcockpit-agent.list
sudo apt-get update

sudo apt-get install openitcockpit-agent
```

#### Manually
Install
```
sudo apt-get install ./openitcockpit-agent_3.x.x_amd64.deb
```

Uninstall
```
sudo apt-get purge openitcockpit-agent
```

### Red Hat Linux / CentOS / openSUSE

#### Using the repository

```
cat <<EOT > /etc/yum.repos.d/openitcockpit-agent.repo
[openitcockpit-agent]
name=openITCOCKPIT Monitoring Agent
baseurl=https://packages.openitcockpit.io/openitcockpit-agent/rpm/stable/$basearch/
enabled=1
gpgcheck=1
gpgkey=https://packages.openitcockpit.io/repokey.txt
EOT

yum-config-manager --enable openitcockpit-agent

yum install openitcockpit-agent
```

#### Manually
Install
```
rpm -i openitcockpit-agent-3.x.x-x.x86_64.rpm
```

Uninstall
```
rpm -e openitcockpit-agent
```

### Arch Linux
Install
```
sudo pacman -U openitcockpit-agent-3.x.x-x-x86_64.pkg.tar.zst
```

Uninstall
```
sudo pacman -R openitcockpit-agent
```

### Windows
Install

**GUI**

Install with double clicking the msi installer file.

![openITCOCKPIT Monitoring Agent MSI installer](/docs/images/msi_installer_new.png)

**CLI**

Automated install

```
msiexec.exe /i openitcockpit-agent*.msi INSTALLDIR="C:\Program Files\openitcockpit\openitcockpit-agent\" /qn
```

Uninstall

Please use the Windows built-in graphical software manager to uninstall.

### macOS

**GUI**

Install with double clicking the pkg installer file.

![openITCOCKPIT Monitoring Agent PKG installer](/docs/images/pkg_install_macos3.png)

**CLI**

Install
```
sudo installer -pkg openitcockpit-agent-3.x.x-darwin-amd64.pkg -target / -verbose
```

Uninstall
```
sudo installer -pkg openitcockpit-agent-uninstaller-3.x.x-darwin-amd64.pkg -target / -verbose
```

### FreeBSD

Currently we do not provide pre-build packages for FreeBSD systems.

But you can build the binary from source like so:
```
GOOS=freebsd GOARCH=amd64 go build -o openitcockpit-agent
```


## Supported Platforms

| Platform              | Windows | Linux | macOS | FreeBSD |
|-----------------------|---------|-------|-------|---------|
| 64 bit (amd64)        | ✅       | ✅     | ✅     | ✅       |
| 32 bit (i386)         | ✅       | ✅     | -     | ✅       |
| arm64 / Apple Silicon | ✅       | ✅     | ✅     | ✅       |


Please see to Wiki how to [cross compile binaries](https://github.com/openITCOCKPIT/openitcockpit-agent-go/wiki/Build-binary#cross-compile) for different operating systems and CPU architectures.

## Full documentation
Do you want to build own binaries, learn more about cross compiling or how to start hacking the Agent?

Please see the [full documentation](https://github.com/openITCOCKPIT/openitcockpit-agent-go/wiki).

## Supported Checks

| Check                        | Windows                                    | Linux                                                                                                      | macOS                       | FreeBSD         |
|------------------------------|--------------------------------------------|------------------------------------------------------------------------------------------------------------|-----------------------------|-----------------|
| Ceph                         | -                                          | ✅                                                                                                         | -                           | -               |
| Certificates                 | ✅ file, TLS, certificate store            | ✅ file, TLS                                                                                               | ✅ file, TLS                | ✅ file, TLS    |
| CPU Load                     | -                                          | ✅                                                                                                         | ✅                          | ✅              |
| CPU Percent                  | ✅                                         | ✅                                                                                                         | ✅                          | ✅              |
| Custom Checks                | ✅                                         | ✅                                                                                                         | ✅                          | ✅              |
| Disk IO                      | ✅                                         | ✅                                                                                                         | ✅                          | ✅              |
| Disk Usage                   | ✅                                         | ✅ incl. inodes                                                                                            | ✅ incl. inodes             | ✅ incl. inodes |
| Docker Container             | ✅ incl. healthchecks                      | ✅ incl. healthchecks                                                                                      | ✅ incl. healthchecks       | untested        |
| Exim Mail Queue              | -                                          | ✅                                                                                                         | -                           | ✅              |
| Files & Directories          | ✅                                         | ✅                                                                                                         | ✅                          | ✅              |
| GPUs                         | ✅ NVIDIA, AMD, Intel, driver errors (TDR) | ✅ NVIDIA, AMD, Intel, Xid and amdgpu errors                                                               | -                           | -               |
| HTTP(S)                      | ✅                                         | ✅                                                                                                         | ✅                          | ✅              |
| IPMI (local BMC)             | ✅ in-band                                 | ✅ in-band                                                                                                 | -                           | -               |
| Logfiles                     | ✅ text files                              | ✅ text files, journald                                                                                    | ✅ text files, unified log  | ✅ text files   |
| LVM Thin Pools               | -                                          | ✅                                                                                                         | -                           | -               |
| Memory                       | ✅                                         | ✅                                                                                                         | ✅                          | ✅              |
| Missing Mounts (fstab)       | -                                          | ✅                                                                                                         | -                           | ✅              |
| Network Interfaces           | ✅                                         | ✅                                                                                                         | ✅                          | ✅              |
| Network IO                   | ✅                                         | ✅                                                                                                         | ✅                          | ✅              |
| Network Ports                | ✅                                         | ✅ incl. TCP errors                                                                                        | ✅                          | ✅              |
| NTP                          | ✅                                         | ✅                                                                                                         | ✅ SNTP                     | ✅ ntpd, SNTP   |
| Postfix Mail Queue           | -                                          | ✅                                                                                                         | -                           | ✅              |
| Processes                    | ✅                                         | ✅                                                                                                         | ✅                          | ✅              |
| Proxmox VE                   | -                                          | ✅                                                                                                         | -                           | -               |
| Redfish (BMC)                | ✅ via network                             | ✅ via network                                                                                             | ✅ via network              | ✅ via network  |
| S.M.A.R.T.                   | ✅ smartctl                                | ✅ smartctl                                                                                                | ✅ smartctl                 | ✅ smartctl     |
| Security                     | ✅ firewall, Defender/antivirus, BitLocker | ✅ SELinux, AppArmor, firewall                                                                             | ✅ FileVault, SIP, firewall | -               |
| Sensors                      | -                                          | ✅                                                                                                         | ✅                          | -               |
| Services                     | ✅                                         | ✅ systemd                                                                                                 | ✅ launchd                  | -               |
| Swap                         | ✅                                         | ✅                                                                                                         | ✅                          | ✅              |
| System Health                | -                                          | ✅ RAID, OOM kills, file descriptors, conntrack, PSI, bonding, filesystem errors, ECC memory errors (EDAC) | -                           | -               |
| Users                        | ✅                                         | ✅                                                                                                         | ✅                          | ✅              |
| Windows Event log            | ✅                                         | -                                                                                                          | -                           | -               |
| Windows Performance Counters | ✅                                         | -                                                                                                          | -                           | -               |
| Windows Scheduled Tasks      | ✅                                         | -                                                                                                          | -                           | -               |
| ZFS                          | -                                          | ✅                                                                                                         | -                           | ✅              |

Logfiles, Certificates, Files & Directories, HTTP(S), Redfish and Windows Performance Counters are configured in separate INI files.
See the examples in the [example](example/) directory. S.M.A.R.T. requires `smartctl` of smartmontools 7.0 or newer.
Ceph, Exim, LVM, Postfix (3.1 or newer), Proxmox VE, S.M.A.R.T. and ZFS are enabled automatically if the required tools are installed.
GPUs: NVIDIA GPUs are read with `nvidia-smi` (Linux and Windows), AMD and Intel GPUs from sysfs (Linux) or with the
interface of the Windows Task Manager (D3DKMT). Linux also reports the Xid errors of the NVIDIA driver and the GPU hangs
and resets of the amdgpu driver from the kernel log, Windows the driver errors (TDR) of the event log.
IPMI reads the sensors, the chassis status and the SEL of the local BMC without ipmitool (Linux: /dev/ipmi0, Windows: IPMI driver).
Network Ports reports the listening ports with their processes, the TCP connections per state and the usage of the ephemeral ports.
Redfish monitors the hardware of servers (Dell iDRAC, HPE iLO, Lenovo XClarity, Supermicro, ...) through the API of their BMCs,
so a single agent in the management network can monitor many servers.

### Advanced Checks

| Check       | Windows | Linux                                                                                              | macOS | FreeBSD |
|-------------|---------|----------------------------------------------------------------------------------------------------|-------|---------|
| Libvirt VMs | -       | ✅ [readme](https://docs.openitcockpit.io/en/agent/build-binary/#enable-libvirt-support-linux-only) | -     | -       |


## Software Inventory

The openITCOCKPIT Agent collects detailed information about installed packages, installed software, and available security updates across all major operating systems. It also detects when a system reboot is required. An important step to ensure that the latest security updates are fully applied.

### Linux


On Linux systems, the agent gathers a complete list of installed packages, identifies available updates, and, where supported, highlights security updates.

For each installed package the agent also reports the architecture and, with apt, dnf, yum, zypper and rpm, the source package and its version, e.g. `curl` for `libcurl4t64`. With pacman only the architecture is reported.

| Package Manager | Installed Packages | Available Updates | Security Updates | Reboot Required                                                                                                     |
|-----------------|--------------------|-------------------|------------------|---------------------------------------------------------------------------------------------------------------------|
| `apt`           | ✅                  | ✅                 | ✅                | ✅ [Method](https://www.debian.org/doc/debian-policy/ch-opersys.html#signaling-that-a-reboot-is-required)            |
| `dnf`           | ✅                  | ✅                 | ✅                | ✅ [Method](https://dnf-plugins-core.readthedocs.io/en/latest/needs_restarting.html)                                 |
| `yum`           | ✅                  | ✅                 | ✅                | ✅ [Method](https://man7.org/linux/man-pages/man1/needs-restarting.1.html)                                           |
| `zypper`        | ✅                  | ✅                 | ✅                | ✅ [Method](https://support.scc.suse.com/s/kb/How-to-check-if-system-reboot-is-needed-after-patching?language=en_US) |
| `pacman`        | ✅                  | ✅                 | -                | _Not supported_                                                                                                     |
| `rpm`           | ✅                  | -                 | -                | -                                                                                                                   |


### Windows

| Available Windows Updates                 | Installed Software | Reboot Required |
|-------------------------------------------|--------------------|-----------------|
| Via PowerShell `Microsoft.Update.Session` | via Registry       | via Registry    |

On Windows systems, only operating system updates are reported, as there is no unified package manager. The agent retrieves information about installed software and reboot requirements via the Windows Registry.


### macOS

| Available macOS Updates | Installed Software    | Reboot Required |
|-------------------------|-----------------------|-----------------|
| ✅                       | via `system_profiler` | _Not supported_ |


On macOS systems, only operating system updates are reported, as there is no package manager. Installed software information is collected using `system_profiler`. Reboot detection is not supported on macOS.


### Third Party Package Managers

| Operating System | Package Manager   |
|------------------|-------------------|
| macOS            | Homebrew (`brew`) |


Support for optional third party package managers can be enabled by the option `enable-alternative-packagemanagers = True` in the agent configuration file.
This package managers are not part of the official operating system and therefore this feature is not enabled by default.

## License
```
Copyright (C) 2021-2025  it-novum GmbH
Copyright (C) 2025-today AVENDIS GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
