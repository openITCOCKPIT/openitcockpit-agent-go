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
* Apple macOS 15 Sequoia or newer (Intel / Apple Silicon)
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
msiexec.exe /i openitcockpit-agent*.msi INSTALLDIR="C:\Program Files\it-novum\openitcockpit-agent\" /qn
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

| Platform              | Windows                | Linux | macOS | FreeBSD |
|-----------------------|------------------------|-------|-------|---------|
| 64 bit (amd64)        | ✅                      | ✅     | ✅     | ✅       |
| 32 bit (i386)         | ✅                      | ✅     | -     | ✅       |
| arm64 / Apple Silicon | Use the 32 bit version | ✅     | ✅     | ✅       |


Please see to Wiki how to [cross compile binaries](https://github.com/openITCOCKPIT/openitcockpit-agent-go/wiki/Build-binary#cross-compile) for different operating systems and CPU architectures.

## Full documentation
Do you want to build own binaries, learn more about cross compiling or how to start hacking the Agent?

Please see the [full documentation](https://github.com/openITCOCKPIT/openitcockpit-agent-go/wiki).

## Supported Checks

| Check              | Windows | Linux     | macOS                   | FreeBSD                 |
|--------------------|---------|-----------|-------------------------|-------------------------|
| CPU Load           | -       | ✅         | ✅                       | ✅                       |
| CPU Percent        | ✅       | ✅         | ✅                       | ✅                       |
| Custom Checks      | ✅       | ✅         | ✅                       | ✅                       |
| Disk IO            | ✅       | ✅         | ✅                       | ✅                       |
| Disk Usage         | ✅       | ✅         | ✅                       | ✅                       |
| Docker Container   | ✅       | ✅         | ✅                       | untested                |
| Memory             | ✅       | ✅         | ✅                       | ✅                       |
| Network Interfaces | ✅       | ✅         | ✅                       | ✅                       |
| Network IO         | ✅       | ✅         | ✅                       | ✅                       |
| NTP                | ✅       | ✅         | Only compare timestamps | Only compare timestamps |
| Processes          | ✅       | ✅         | ✅                       | ✅                       |
| Sensors            | -       | ✅         | ✅                       | -                       |
| Services           | ✅       | ✅ systemd | ✅ launchd               | -                       |
| Swap               | ✅       | ✅         | ✅                       | ✅                       |
| Windows Event log  | ✅       | -         | -                       | -                       |
| Users              | ✅       | ✅         | ✅                       | ✅                       |

### Advanced Checks

| Check       | Windows | Linux                                                                                              | macOS | FreeBSD |
|-------------|---------|----------------------------------------------------------------------------------------------------|-------|---------|
| Libvirt VMs | -       | ✅ [readme](https://docs.openitcockpit.io/en/agent/build-binary/#enable-libvirt-support-linux-only) | -     | -       |


## Software Inventory

This agent can collect information about installed packages, installed software, and available security updates for all major operating systems.
The agent will also detect, if a reboot of the system is required. Reboots are often necessary to apply the latest security updates.

### Linux

On linux systems, the agent will collect all installed packages, and list available updates and also (if supported) security updates.

| Package manager | Installed packages | Available Updates | Security Updates | Reebot Required                                                                                                     |   |
|-----------------|--------------------|-------------------|------------------|---------------------------------------------------------------------------------------------------------------------|---|
| `apt`           | ✅                  | ✅                 | ✅                | ✅ [Method](https://www.debian.org/doc/debian-policy/ch-opersys.html#signaling-that-a-reboot-is-required)            |   |
| `dnf`           | ✅                  | ✅                 | ✅                | ✅ [Method](https://dnf-plugins-core.readthedocs.io/en/latest/needs_restarting.html)                                 |   |
| `yum`           | ✅                  | ✅                 | ✅                | ✅ [Method](https://man7.org/linux/man-pages/man1/needs-restarting.1.html)                                           |   |
| `zypper`        | ✅                  | ✅                 | ✅                | ✅ [Method](https://support.scc.suse.com/s/kb/How-to-check-if-system-reboot-is-needed-after-patching?language=en_US) |   |
| `pacman`        | ✅                  | ✅                 | -                | _Not supported_                                                                                                     |   |
| `rpm`           | ✅                  | -                 | -                | -                                                                                                                   |   |



### Windows

| Available Windows Updates                 | Installed Software | Reebot Required |
|-------------------------------------------|--------------------|-----------------|
| Via PowerShell `Microsoft.Update.Session` | via Registry       | via Registry    |

On Windows Systems, only operating system related updates will be reported, as there is no package manager available.

### macOS

| Available macOS Updates | Installed Software    | Reebot Required |
|-------------------------|-----------------------|-----------------|
| ✅                       | via `system_profiler` | _Not supported_ |

On macOS Systems, only operating system related updates will be reported, as there is no package manager available.

## License
```
Copyright (C) 2021-2025  it-novum GmbH
Copyright (C) 2025-today Allgeier IT Services GmbH

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
