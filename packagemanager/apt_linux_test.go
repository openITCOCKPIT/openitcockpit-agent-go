package packagemanager

import (
	"strings"
	"testing"
)

// Sample output of the command
// dpkg-query -W  "-f=\${Package} \${Version} \${Description}\nEND\n"
var dpkgQueryOutput = `apt 2.8.3 commandline package manager
 This package provides commandline tools for searching and
 managing as well as querying information about packages
 as a low-level access to all features of the libapt-pkg library.
 .
 These include:
  * apt-get for retrieval of packages and information about them
    from authenticated sources and for installation, upgrade and
    removal of packages together with their dependencies
  * apt-cache for querying available information about installed
    as well as installable packages
  * apt-cdrom to use removable media as a source for packages
  * apt-config as an interface to the configuration settings
  * apt-key as an interface to manage authentication keys
END
base-files 13ubuntu10.3 Debian base system miscellaneous files
 This package contains the basic filesystem hierarchy of a Debian system, and
 several important miscellaneous files, such as /etc/debian_version,
 /etc/host.conf, /etc/issue, /etc/motd, /etc/profile, and others,
 and the text of several common licenses in use on Debian systems.
END
sed 4.9-2build1 GNU stream editor for filtering/transforming text
 sed reads the specified files or the standard input if no
 files are specified, makes editing changes according to a
 list of commands, and writes the results to the standard
 output.
END
sensible-utils 0.0.22 Utilities for sensible alternative selection
 This package provides a number of small utilities which are used
 by programs to sensibly select and spawn an appropriate browser,
 editor, pager, or terminal emulator.
 .
 The specific utilities included are: sensible-browser sensible-editor
 sensible-pager sensible-terminal
END
sysvinit-utils 3.08-6ubuntu3 System-V-like utilities
 This package contains the important System-V-like utilities.
 .
 Specifically, this package includes:
 init-d-script, fstab-decode, killall5, pidof
 .
 It also contains the library scripts sourced by init-d-script and other
 initscripts that were formally in lsb-base.
END
tar 1.35+dfsg-3build1 GNU version of the tar archiving utility
 Tar is a program for packaging a set of files as a single archive in tar
 format.  The function it performs is conceptually similar to cpio, and to
 things like PKZIP in the DOS world.  It is heavily used by the Debian package
 management system, and is useful for performing system backups and exchanging
 sets of files with others.
END
ubuntu-keyring 2023.11.28.1 GnuPG keys of the Ubuntu archive
 The Ubuntu project digitally signs its Release files. This package
 contains the archive keys used for that.
END
unminimize 0.2.1 Un-minimize your minimial images or setup
 This package contains the unminimize script that helps the user in
 unminimizing their minimized images or cloud setup.
END
util-linux 2.39.3-9ubuntu6.3 miscellaneous system utilities
 This package contains a number of important utilities, most of which
 are oriented towards maintenance of your system. Some of the more
 important utilities included in this package allow you to view kernel
 messages, create new filesystems, view block device information,
 interface with real time clock, etc.
END
zlib1g 1:1.3.dfsg-3.1ubuntu2.1 compression library - runtime
 zlib is a library implementing the deflate compression method found
 in gzip and PKZIP.  This package includes the shared library.
END`

// LANG=c apt-get -s -o Debug::NoLocking=1 upgrade
var ubuntuNoUpdatesAvailableOutput = `Reading package lists... Done
Building dependency tree... Done
Reading state information... Done
Calculating upgrade... Done
0 upgraded, 0 newly installed, 0 to remove and 0 not upgraded.`

var ubuntuUpdatesAvailableOutput = `NOTE: This is only a simulation!
      apt-get needs root privileges for real execution.
      Keep also in mind that locking is deactivated,
      so don't depend on the relevance to the real current situation!
Reading package lists... Done
Building dependency tree... Done
Reading state information... Done
Calculating upgrade... Done
The following packages were automatically installed and are no longer required:
  libdrm-nouveau2 libdrm-radeon1 libgl1-amber-dri libglapi-amber libllvm19 libxcb-dri2-0
Use 'apt autoremove' to remove them.
The following packages will be upgraded:
  gir1.2-glib-2.0 libgirepository-2.0-0 libglib2.0-0t64 libglib2.0-bin libglib2.0-data libglib2.0-dev libglib2.0-dev-bin
7 upgraded, 0 newly installed, 0 to remove and 0 not upgraded.
Inst libglib2.0-data [2.80.0-6ubuntu3.5] (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [all])
Inst libglib2.0-dev [2.80.0-6ubuntu3.5] (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64]) []
Inst libglib2.0-bin [2.80.0-6ubuntu3.5] (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64]) []
Inst libgirepository-2.0-0 [2.80.0-6ubuntu3.5] (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64]) []
Inst libglib2.0-dev-bin [2.80.0-6ubuntu3.5] (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64]) []
Inst gir1.2-glib-2.0 [2.80.0-6ubuntu3.5] (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64]) []
Inst libglib2.0-0t64 [2.80.0-6ubuntu3.5] (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64])
Conf libglib2.0-data (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [all])
Conf libglib2.0-dev (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64])
Conf libglib2.0-bin (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64])
Conf libgirepository-2.0-0 (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64])
Conf libglib2.0-dev-bin (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64])
Conf gir1.2-glib-2.0 (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64])
Conf libglib2.0-0t64 (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64])`

var ubuntuUpdatesAvailableOutputLargeOutput = `Reading package lists... Done
Building dependency tree... Done
Reading state information... Done
Calculating upgrade... Done
The following packages have been kept back:
  linux-headers-generic linux-headers-virtual linux-image-virtual linux-virtual
The following packages will be upgraded:
  apparmor bsdextrautils bsdutils dhcpcd-base eject fdisk gir1.2-glib-2.0 landscape-common libapparmor1 libblkid1 libfdisk1 libglib2.0-0t64 libglib2.0-bin libglib2.0-data libmbim-glib4 libmbim-proxy libmbim-utils libmount1 libmysqlclient21
  libnetplan1 libnss-systemd libpam-systemd libsmartcols1 libsystemd-shared libsystemd0 libudev1 libuuid1 linux-libc-dev linux-tools-common mount mysql-client-8.0 mysql-client-core-8.0 mysql-server mysql-server-8.0 mysql-server-core-8.0
  netplan-generator netplan.io openitcockpit openitcockpit-agent openitcockpit-check-network-interfaces-plugin openitcockpit-check-proxmox-plugin openitcockpit-checkmk openitcockpit-common openitcockpit-event-collectd openitcockpit-frontend-angular
  openitcockpit-graphing openitcockpit-module-autoreport openitcockpit-module-checkmk openitcockpit-module-design openitcockpit-module-distribute openitcockpit-module-evc openitcockpit-module-grafana openitcockpit-module-import
  openitcockpit-module-map openitcockpit-module-network openitcockpit-module-prometheus openitcockpit-module-proxmox openitcockpit-module-sla openitcockpit-naemon openitcockpit-node openitcockpit-nsta openitcockpit-statusengine3-worker
  python3-netplan qemu-guest-agent systemd systemd-dev systemd-resolved systemd-sysv systemd-timesyncd udev util-linux uuid-runtime
72 upgraded, 0 newly installed, 0 to remove and 4 not upgraded.
Inst bsdutils [1:2.39.3-9ubuntu6.3] (1:2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Conf bsdutils (1:2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Inst util-linux [2.39.3-9ubuntu6.3] (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Conf util-linux (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Inst mount [2.39.3-9ubuntu6.3] (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Inst libnss-systemd [255.4-1ubuntu8.11] (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64]) []
Inst systemd-dev [255.4-1ubuntu8.11] (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [all]) []
Inst libblkid1 [2.39.3-9ubuntu6.3] (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64]) []
Conf libblkid1 (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64]) []
Inst systemd-timesyncd [255.4-1ubuntu8.11] (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64]) []
Inst systemd-resolved [255.4-1ubuntu8.11] (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64]) []
Inst libsystemd-shared [255.4-1ubuntu8.11] (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64]) [systemd:amd64 ]
Inst libsystemd0 [255.4-1ubuntu8.11] (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64]) [systemd:amd64 ]
Conf libsystemd0 (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64]) [systemd:amd64 ]
Inst systemd-sysv [255.4-1ubuntu8.11] (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64]) [systemd:amd64 ]
Inst libpam-systemd [255.4-1ubuntu8.11] (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64]) [systemd:amd64 ]
Inst systemd [255.4-1ubuntu8.11] (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64])
Inst udev [255.4-1ubuntu8.11] (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64]) []
Inst libudev1 [255.4-1ubuntu8.11] (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64])
Conf libudev1 (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64])
Inst libapparmor1 [4.0.1really4.0.1-0ubuntu0.24.04.4] (4.0.1really4.0.1-0ubuntu0.24.04.5 Ubuntu:24.04/noble-updates [amd64])
Inst libmount1 [2.39.3-9ubuntu6.3] (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Conf libmount1 (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Inst libuuid1 [2.39.3-9ubuntu6.3] (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Conf libuuid1 (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Inst libfdisk1 [2.39.3-9ubuntu6.3] (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Inst libsmartcols1 [2.39.3-9ubuntu6.3] (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Conf libsmartcols1 (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Inst uuid-runtime [2.39.3-9ubuntu6.3] (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Inst mysql-client-core-8.0 [8.0.44-0ubuntu0.24.04.1] (8.0.44-0ubuntu0.24.04.2 Ubuntu:24.04/noble-updates [amd64])
Inst mysql-client-8.0 [8.0.44-0ubuntu0.24.04.1] (8.0.44-0ubuntu0.24.04.2 Ubuntu:24.04/noble-updates [amd64])
Inst mysql-server-8.0 [8.0.44-0ubuntu0.24.04.1] (8.0.44-0ubuntu0.24.04.2 Ubuntu:24.04/noble-updates [amd64]) []
Inst mysql-server-core-8.0 [8.0.44-0ubuntu0.24.04.1] (8.0.44-0ubuntu0.24.04.2 Ubuntu:24.04/noble-updates [amd64])
Inst openitcockpit-checkmk [2.2.0-20251202095540generic] (2.2.0-20251211153142generic noble/stable noble:noble [amd64])
Inst openitcockpit-graphing [5.3.0-20251202095323generic] (5.3.1-20251211152816generic noble/stable noble:noble [amd64])
Inst openitcockpit-node [5.3.0-20251202100036generic] (5.3.1-20251211153608generic noble/stable noble:noble [amd64])
Inst gir1.2-glib-2.0 [2.80.0-6ubuntu3.4] (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64]) []
Inst libglib2.0-data [2.80.0-6ubuntu3.4] (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [all]) []
Inst libglib2.0-bin [2.80.0-6ubuntu3.4] (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64]) []
Inst libglib2.0-0t64 [2.80.0-6ubuntu3.4] (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64])
Inst qemu-guest-agent [1:8.2.2+ds-0ubuntu1.10] (1:8.2.2+ds-0ubuntu1.11 Ubuntu:24.04/noble-updates [amd64])
Inst dhcpcd-base [1:10.0.6-1ubuntu3.1] (1:10.0.6-1ubuntu3.2 Ubuntu:24.04/noble-updates [amd64])
Inst eject [2.39.3-9ubuntu6.3] (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Inst netplan-generator [1.1.2-2~ubuntu24.04.2] (1.1.2-8ubuntu1~24.04.1 Ubuntu:24.04/noble-updates [amd64]) []
Inst python3-netplan [1.1.2-2~ubuntu24.04.2] (1.1.2-8ubuntu1~24.04.1 Ubuntu:24.04/noble-updates [amd64]) []
Inst netplan.io [1.1.2-2~ubuntu24.04.2] (1.1.2-8ubuntu1~24.04.1 Ubuntu:24.04/noble-updates [amd64]) []
Inst libnetplan1 [1.1.2-2~ubuntu24.04.2] (1.1.2-8ubuntu1~24.04.1 Ubuntu:24.04/noble-updates [amd64])
Inst apparmor [4.0.1really4.0.1-0ubuntu0.24.04.4] (4.0.1really4.0.1-0ubuntu0.24.04.5 Ubuntu:24.04/noble-updates [amd64])
Inst bsdextrautils [2.39.3-9ubuntu6.3] (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Inst fdisk [2.39.3-9ubuntu6.3] (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Inst landscape-common [24.02-0ubuntu5.6] (24.02-0ubuntu5.7 Ubuntu:24.04/noble-updates [amd64])
Inst libmbim-proxy [1.31.2-0ubuntu3] (1.31.2-0ubuntu3.1 Ubuntu:24.04/noble-updates [amd64]) []
Inst libmbim-glib4 [1.31.2-0ubuntu3] (1.31.2-0ubuntu3.1 Ubuntu:24.04/noble-updates [amd64])
Inst libmbim-utils [1.31.2-0ubuntu3] (1.31.2-0ubuntu3.1 Ubuntu:24.04/noble-updates [amd64])
Inst libmysqlclient21 [8.0.44-0ubuntu0.24.04.1] (8.0.44-0ubuntu0.24.04.2 Ubuntu:24.04/noble-updates [amd64])
Inst linux-libc-dev [6.8.0-88.89] (6.8.0-90.91 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64])
Inst linux-tools-common [6.8.0-88.89] (6.8.0-90.91 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [all])
Inst mysql-server [8.0.44-0ubuntu0.24.04.1] (8.0.44-0ubuntu0.24.04.2 Ubuntu:24.04/noble-updates [all])
Inst openitcockpit-common [5.3.0-20251202095317noble] (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Inst openitcockpit-naemon [1.4.4-20251202102416noble] (1.4.4-20251211160843noble noble/stable noble:noble [amd64])
Inst openitcockpit-statusengine3-worker [3.8.2-20251202102859noble] (3.8.2-20251211161324noble noble/stable noble:noble [amd64])
Inst openitcockpit [5.3.0-20251202102709noble] (5.3.1-20251211161137noble noble/stable noble:noble [amd64])
Inst openitcockpit-agent [3.3.0] (3.4.0 deb/stable deb:deb [amd64])
Inst openitcockpit-check-network-interfaces-plugin [1.0.3-20251202105320noble] (1.0.3-20251211163620noble noble/stable noble:noble [amd64])
Inst openitcockpit-check-proxmox-plugin [1.0.2-20251202104251noble] (1.0.2-20251211162647noble noble/stable noble:noble [amd64])
Inst openitcockpit-event-collectd [5.3.0-20251202103744noble] (5.3.1-20251211162219noble noble/stable noble:noble [amd64])
Inst openitcockpit-frontend-angular [5.3.0-20251202100216generic] (5.3.1-20251211153746generic noble/stable noble:noble [amd64])
Inst openitcockpit-module-autoreport [5.3.0-20251202095317noble] (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Inst openitcockpit-module-checkmk [5.3.0-20251202095317noble] (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Inst openitcockpit-module-design [5.3.0-20251202095317noble] (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Inst openitcockpit-nsta [5.3.0-20251202102630noble] (5.3.1-20251211161100noble noble/stable noble:noble [amd64])
Inst openitcockpit-module-distribute [5.3.0-20251202095317noble] (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Inst openitcockpit-module-evc [5.3.0-20251202095317noble] (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Inst openitcockpit-module-grafana [5.3.0-20251202095317noble] (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Inst openitcockpit-module-import [5.3.0-20251202095317noble] (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Inst openitcockpit-module-map [5.3.0-20251202095317noble] (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Inst openitcockpit-module-network [5.3.0-20251202095317noble] (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Inst openitcockpit-module-prometheus [5.3.0-20251202095317noble] (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Inst openitcockpit-module-proxmox [5.3.0-20251202095317noble] (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Inst openitcockpit-module-sla [5.3.0-20251202095317noble] (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Conf mount (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Conf libnss-systemd (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64])
Conf systemd-dev (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [all])
Conf systemd-timesyncd (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64])
Conf systemd-resolved (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64])
Conf libsystemd-shared (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64])
Conf systemd-sysv (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64])
Conf libpam-systemd (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64])
Conf systemd (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64])
Conf udev (255.4-1ubuntu8.12 Ubuntu:24.04/noble-updates [amd64])
Conf libapparmor1 (4.0.1really4.0.1-0ubuntu0.24.04.5 Ubuntu:24.04/noble-updates [amd64])
Conf libfdisk1 (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Conf uuid-runtime (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Conf mysql-client-core-8.0 (8.0.44-0ubuntu0.24.04.2 Ubuntu:24.04/noble-updates [amd64])
Conf mysql-client-8.0 (8.0.44-0ubuntu0.24.04.2 Ubuntu:24.04/noble-updates [amd64])
Conf mysql-server-8.0 (8.0.44-0ubuntu0.24.04.2 Ubuntu:24.04/noble-updates [amd64])
Conf mysql-server-core-8.0 (8.0.44-0ubuntu0.24.04.2 Ubuntu:24.04/noble-updates [amd64])
Conf openitcockpit-checkmk (2.2.0-20251211153142generic noble/stable noble:noble [amd64])
Conf openitcockpit-graphing (5.3.1-20251211152816generic noble/stable noble:noble [amd64])
Conf openitcockpit-node (5.3.1-20251211153608generic noble/stable noble:noble [amd64])
Conf gir1.2-glib-2.0 (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64])
Conf libglib2.0-data (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [all])
Conf libglib2.0-bin (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64])
Conf libglib2.0-0t64 (2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64])
Conf qemu-guest-agent (1:8.2.2+ds-0ubuntu1.11 Ubuntu:24.04/noble-updates [amd64])
Conf dhcpcd-base (1:10.0.6-1ubuntu3.2 Ubuntu:24.04/noble-updates [amd64])
Conf eject (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Conf netplan-generator (1.1.2-8ubuntu1~24.04.1 Ubuntu:24.04/noble-updates [amd64])
Conf python3-netplan (1.1.2-8ubuntu1~24.04.1 Ubuntu:24.04/noble-updates [amd64])
Conf netplan.io (1.1.2-8ubuntu1~24.04.1 Ubuntu:24.04/noble-updates [amd64])
Conf libnetplan1 (1.1.2-8ubuntu1~24.04.1 Ubuntu:24.04/noble-updates [amd64])
Conf apparmor (4.0.1really4.0.1-0ubuntu0.24.04.5 Ubuntu:24.04/noble-updates [amd64])
Conf bsdextrautils (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Conf fdisk (2.39.3-9ubuntu6.4 Ubuntu:24.04/noble-updates [amd64])
Conf landscape-common (24.02-0ubuntu5.7 Ubuntu:24.04/noble-updates [amd64])
Conf libmbim-proxy (1.31.2-0ubuntu3.1 Ubuntu:24.04/noble-updates [amd64])
Conf libmbim-glib4 (1.31.2-0ubuntu3.1 Ubuntu:24.04/noble-updates [amd64])
Conf libmbim-utils (1.31.2-0ubuntu3.1 Ubuntu:24.04/noble-updates [amd64])
Conf libmysqlclient21 (8.0.44-0ubuntu0.24.04.2 Ubuntu:24.04/noble-updates [amd64])
Conf linux-libc-dev (6.8.0-90.91 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [amd64])
Conf linux-tools-common (6.8.0-90.91 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [all])
Conf mysql-server (8.0.44-0ubuntu0.24.04.2 Ubuntu:24.04/noble-updates [all])
Conf openitcockpit-common (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Conf openitcockpit-naemon (1.4.4-20251211160843noble noble/stable noble:noble [amd64])
Conf openitcockpit-statusengine3-worker (3.8.2-20251211161324noble noble/stable noble:noble [amd64])
Conf openitcockpit (5.3.1-20251211161137noble noble/stable noble:noble [amd64])
Conf openitcockpit-agent (3.4.0 deb/stable deb:deb [amd64])
Conf openitcockpit-check-network-interfaces-plugin (1.0.3-20251211163620noble noble/stable noble:noble [amd64])
Conf openitcockpit-check-proxmox-plugin (1.0.2-20251211162647noble noble/stable noble:noble [amd64])
Conf openitcockpit-event-collectd (5.3.1-20251211162219noble noble/stable noble:noble [amd64])
Conf openitcockpit-frontend-angular (5.3.1-20251211153746generic noble/stable noble:noble [amd64])
Conf openitcockpit-module-autoreport (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Conf openitcockpit-module-checkmk (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Conf openitcockpit-module-design (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Conf openitcockpit-nsta (5.3.1-20251211161100noble noble/stable noble:noble [amd64])
Conf openitcockpit-module-distribute (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Conf openitcockpit-module-evc (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Conf openitcockpit-module-grafana (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Conf openitcockpit-module-import (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Conf openitcockpit-module-map (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Conf openitcockpit-module-network (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Conf openitcockpit-module-prometheus (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Conf openitcockpit-module-proxmox (5.3.1-20251211154831noble noble/stable noble:noble [amd64])
Conf openitcockpit-module-sla (5.3.1-20251211154831noble noble/stable noble:noble [amd64])`

var debianUpdatesAvailableOutputLargeOutput = `Reading package lists... Done
Building dependency tree... Done
Reading state information... Done
Calculating upgrade... Done
The following packages will be upgraded:
  libsodium23 openitcockpit openitcockpit-agent openitcockpit-checkmk openitcockpit-common openitcockpit-event-collectd openitcockpit-frontend-angular openitcockpit-graphing openitcockpit-module-autoreport openitcockpit-module-change-calendar
  openitcockpit-module-checkmk openitcockpit-module-design openitcockpit-module-distribute openitcockpit-module-evc openitcockpit-module-grafana openitcockpit-module-import openitcockpit-module-nwc openitcockpit-module-prometheus
  openitcockpit-module-scm openitcockpit-module-snmp-trap openitcockpit-module-vmware-snapshot openitcockpit-naemon openitcockpit-node openitcockpit-nsta openitcockpit-statusengine3-worker php8.4-bcmath php8.4-bz2 php8.4-cli php8.4-common php8.4-curl
  php8.4-fpm php8.4-gd php8.4-intl php8.4-ldap php8.4-mbstring php8.4-mysql php8.4-opcache php8.4-phpdbg php8.4-readline php8.4-soap php8.4-xml php8.4-zip
42 upgraded, 0 newly installed, 0 to remove and 0 not upgraded.
Inst openitcockpit-checkmk [2.2.0-20251103055546generic] (2.2.0-20251212052813generic trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-graphing [5.2.0-20251103052518generic] (5.3.1-20251212052511generic trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-node [5.2.0-20251103061057generic] (5.3.1-20251212053245generic trixie/nightly trixie:trixie [amd64])
Inst libsodium23 [1.0.18-1+b2] (1.0.18-1+deb13u1 Debian-Security:13/stable-security [amd64])
Inst php8.4-zip [8.4.11-1] (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64]) []
Inst php8.4-xml [8.4.11-1] (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64]) []
Inst php8.4-soap [8.4.11-1] (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64]) []
Inst php8.4-readline [8.4.11-1] (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64]) []
Inst php8.4-opcache [8.4.11-1] (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64]) []
Inst php8.4-cli [8.4.11-1] (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64]) []
Inst php8.4-phpdbg [8.4.11-1] (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64]) []
Inst php8.4-mysql [8.4.11-1] (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64]) []
Inst php8.4-mbstring [8.4.11-1] (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64]) []
Inst php8.4-intl [8.4.11-1] (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64]) []
Inst php8.4-gd [8.4.11-1] (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64]) []
Inst php8.4-fpm [8.4.11-1] (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64]) []
Inst php8.4-curl [8.4.11-1] (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64]) []
Inst php8.4-bz2 [8.4.11-1] (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64]) []
Inst php8.4-bcmath [8.4.11-1] (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64]) []
Inst php8.4-ldap [8.4.11-1] (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64]) []
Inst php8.4-common [8.4.11-1] (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Inst openitcockpit-common [5.2.0-20251103052522trixie] (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-naemon [1.4.4-20251103054811trixie] (1.4.4-20251212054633trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-statusengine3-worker [3.8.2-20251103055402trixie] (3.8.2-20251212055047trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit [5.2.0-20251103055139trixie] (5.3.1-20251212054900trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-agent [3.2.1] (3.4.0 deb/stable deb:deb [amd64])
Inst openitcockpit-event-collectd [5.2.0-20251103060422trixie] (5.3.1-20251212055755trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-frontend-angular [5.2.0-20251103061442generic] (5.3.1-20251212053420generic trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-module-autoreport [5.2.0-20251103052522trixie] (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-module-change-calendar [5.2.0-20251103052522trixie] (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-module-checkmk [5.2.0-20251103052522trixie] (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-module-design [5.2.0-20251103052522trixie] (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-nsta [5.2.0-20251103055055trixie] (5.3.1-20251212054825trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-module-distribute [5.2.0-20251103052522trixie] (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-module-evc [5.2.0-20251103052522trixie] (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-module-grafana [5.2.0-20251103052522trixie] (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-module-import [5.2.0-20251103052522trixie] (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-module-nwc [5.2.0-20251103052522trixie] (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-module-prometheus [5.2.0-20251103052522trixie] (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-module-scm [5.2.0-20251103052522trixie] (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-module-snmp-trap [5.2.0-20251103052522trixie] (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Inst openitcockpit-module-vmware-snapshot [5.2.0-20251103052522trixie] (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-checkmk (2.2.0-20251212052813generic trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-graphing (5.3.1-20251212052511generic trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-node (5.3.1-20251212053245generic trixie/nightly trixie:trixie [amd64])
Conf libsodium23 (1.0.18-1+deb13u1 Debian-Security:13/stable-security [amd64])
Conf php8.4-zip (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Conf php8.4-xml (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Conf php8.4-soap (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Conf php8.4-readline (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Conf php8.4-opcache (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Conf php8.4-cli (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Conf php8.4-phpdbg (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Conf php8.4-mysql (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Conf php8.4-mbstring (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Conf php8.4-intl (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Conf php8.4-gd (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Conf php8.4-fpm (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Conf php8.4-curl (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Conf php8.4-bz2 (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Conf php8.4-bcmath (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Conf php8.4-ldap (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Conf php8.4-common (8.4.16-1~deb13u1 Debian-Security:13/stable-security [amd64])
Conf openitcockpit-common (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-naemon (1.4.4-20251212054633trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-statusengine3-worker (3.8.2-20251212055047trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit (5.3.1-20251212054900trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-agent (3.4.0 deb/stable deb:deb [amd64])
Conf openitcockpit-event-collectd (5.3.1-20251212055755trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-frontend-angular (5.3.1-20251212053420generic trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-module-autoreport (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-module-change-calendar (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-module-checkmk (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-module-design (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-nsta (5.3.1-20251212054825trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-module-distribute (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-module-evc (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-module-grafana (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-module-import (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-module-nwc (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-module-prometheus (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-module-scm (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-module-snmp-trap (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])
Conf openitcockpit-module-vmware-snapshot (5.3.1-20251212052513trixie trixie/nightly trixie:trixie [amd64])`

func TestAptManager_getInstalledPackagesWithCancel(t *testing.T) {

	apt := AptManager{}

	if !apt.IsAvailable() {
		t.Skip("Skipping test; apt is not available on this system.")
	}

	ctx := t.Context()
	output, err := apt.getInstalledPackagesWithCancel(ctx)
	if err != nil {
		t.Fatalf("Error getting installed packages: %v", err)
	}

	if !strings.HasPrefix(output, "adduser") {
		t.Errorf("Unexpected output, got: %s", output)
	}
}

func TestAptManager_parseDpkgOutput(t *testing.T) {

	apt := AptManager{}

	packages, err := apt.parseDpkgOutput(dpkgQueryOutput)
	if err != nil {
		t.Fatalf("Error parsing dpkg output: %v", err)
	}

	// Table-driven test cases
	tests := []struct {
		name        string
		pkgName     string
		wantVersion string
		wantDesc    string
	}{
		{
			name:        "apt package",
			pkgName:     "apt",
			wantVersion: "2.8.3",
			wantDesc:    "commandline package manager",
		},
		{
			name:        "sed package",
			pkgName:     "sed",
			wantVersion: "4.9-2build1",
			wantDesc:    "GNU stream editor for filtering/transforming text",
		},
		{
			name:        "util-linux package",
			pkgName:     "util-linux",
			wantVersion: "2.39.3-9ubuntu6.3",
			wantDesc:    "miscellaneous system utilities",
		},
		{
			name:        "zlib1g package",
			pkgName:     "zlib1g",
			wantVersion: "1:1.3.dfsg-3.1ubuntu2.1",
			wantDesc:    "compression library - runtime",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, pkg := range packages {
				// Assume pkg has Name, Version, Description fields
				if pkg.Name == tt.pkgName {
					found = true
					if pkg.Version != tt.wantVersion {
						t.Errorf("%s: got version %q, want %q", tt.pkgName, pkg.Version, tt.wantVersion)
					}
					if !strings.HasPrefix(pkg.Description, tt.wantDesc) {
						t.Errorf("%s: got description %q, want %q", tt.pkgName, pkg.Description, tt.wantDesc)
					}
				}
			}
			if !found {
				t.Errorf("Package %q not found in parsed output", tt.pkgName)
			}
		})
	}
}

func TestAptManager_parseAptUpgradeOutput_NoUpdates(t *testing.T) {
	apt := AptManager{}

	updates, err := apt.parseAptUpgradeOutput(ubuntuNoUpdatesAvailableOutput)
	if err != nil {
		t.Fatalf("Error parsing apt upgrade output: %v", err)
	}

	if len(updates) != 0 {
		t.Errorf("Expected no updates, got: %v", updates)
	}
}

func TestAptManager_parseAptUpgradeOutput_WithUpdates(t *testing.T) {
	apt := AptManager{}

	updates, err := apt.parseAptUpgradeOutput(ubuntuUpdatesAvailableOutput)
	if err != nil {
		t.Fatalf("Error parsing apt upgrade output: %v", err)
	}

	expectedUpdates := map[string]struct {
		currentVersion   string
		availableVersion string
		isSecurityUpdate bool
	}{
		"libglib2.0-data":       {"2.80.0-6ubuntu3.5", "2.80.0-6ubuntu3.6", true},
		"libglib2.0-dev":        {"2.80.0-6ubuntu3.5", "2.80.0-6ubuntu3.6", true},
		"libglib2.0-bin":        {"2.80.0-6ubuntu3.5", "2.80.0-6ubuntu3.6", true},
		"libgirepository-2.0-0": {"2.80.0-6ubuntu3.5", "2.80.0-6ubuntu3.6", true},
		"libglib2.0-dev-bin":    {"2.80.0-6ubuntu3.5", "2.80.0-6ubuntu3.6", true},
		"gir1.2-glib-2.0":       {"2.80.0-6ubuntu3.5", "2.80.0-6ubuntu3.6", true},
		"libglib2.0-0t64":       {"2.80.0-6ubuntu3.5", "2.80.0-6ubuntu3.6", true},
	}

	if len(updates) != len(expectedUpdates) {
		t.Fatalf("Expected %d updates, got %d", len(expectedUpdates), len(updates))
	}

	for _, update := range updates {
		expected, exists := expectedUpdates[update.Name]
		if !exists {
			t.Errorf("Unexpected update found: %s", update.Name)
			continue
		}

		if update.CurrentVersion != expected.currentVersion {
			t.Errorf("%s: got current version %q, want %q", update.Name, update.CurrentVersion, expected.currentVersion)
		}

		if update.AvailableVersion != expected.availableVersion {
			t.Errorf("%s: got available version %q, want %q", update.Name, update.AvailableVersion, expected.availableVersion)
		}
		if update.IsSecurityUpdate != expected.isSecurityUpdate {
			t.Errorf("%s: got isSecurityUpdate %v, want %v", update.Name, update.IsSecurityUpdate, expected.isSecurityUpdate)
		}
	}
}

func TestAptManager_parseAptUpgradeLargeOutput_WithUpdates(t *testing.T) {
	apt := AptManager{}
	updates, err := apt.parseAptUpgradeOutput(ubuntuUpdatesAvailableOutputLargeOutput)
	if err != nil {
		t.Fatalf("Error parsing apt upgrade output: %v", err)
	}

	if len(updates) != 72 {
		t.Fatalf("Expected 72 updates, got %d", len(updates))
	}

	// We expect 6 security updates in this large output
	securityUpdateCount := 0
	for _, update := range updates {
		if update.IsSecurityUpdate {
			securityUpdateCount++
		}
	}

	if securityUpdateCount != 6 {
		t.Errorf("Expected 6 security updates, got %d", securityUpdateCount)
	}
}

func TestAptManager_parseAptUpgradeLargeOutput_WithUpdatesDebian(t *testing.T) {
	apt := AptManager{}
	updates, err := apt.parseAptUpgradeOutput(debianUpdatesAvailableOutputLargeOutput)
	if err != nil {
		t.Fatalf("Error parsing apt upgrade output: %v", err)
	}

	if len(updates) != 42 {
		t.Fatalf("Expected 42 updates, got %d", len(updates))
	}

	// We expect 18 security updates in this large output
	securityUpdateCount := 0
	for _, update := range updates {
		if update.IsSecurityUpdate {
			securityUpdateCount++
		}
	}

	if securityUpdateCount != 18 {
		t.Errorf("Expected 18 security updates, got %d", securityUpdateCount)
	}
}
