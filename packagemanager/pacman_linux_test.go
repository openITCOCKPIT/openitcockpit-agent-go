package packagemanager

import (
	"context"
	"testing"
)

// Sample output of the command
// LANG=c pacman -Qi
var samplePacmanQiOutput = `Name            : acl
Version         : 2.3.2-1
Description     : Access control list utilities, libraries and headers
Architecture    : x86_64
URL             : https://savannah.nongnu.org/projects/acl
Licenses        : LGPL
Groups          : None
Provides        : xfsacl  libacl.so=1-64
Depends On      : glibc
Optional Deps   : None
Required By     : coreutils  gettext  libarchive  sed  shadow  systemd  tar  vim
Optional For    : None
Conflicts With  : xfsacl
Replaces        : xfsacl
Installed Size  : 329.98 KiB
Packager        : Christian Hesse <eworm@archlinux.org>
Build Date      : Wed Jan 24 08:57:20 2024
Install Date    : Sun Oct 5 00:04:30 2025
Install Reason  : Installed as a dependency for another package
Install Script  : No
Validated By    : Signature

Name            : archlinux-keyring
Version         : 20250929-1
Description     : Arch Linux PGP keyring
Architecture    : any
URL             : https://gitlab.archlinux.org/archlinux/archlinux-keyring/
Licenses        : GPL-3.0-or-later
Groups          : None
Provides        : None
Depends On      : pacman
Optional Deps   : None
Required By     : base
Optional For    : None
Conflicts With  : None
Replaces        : None
Installed Size  : 1728.89 KiB
Packager        : Christian Hesse <eworm@archlinux.org>
Build Date      : Mon Sep 29 06:55:38 2025
Install Date    : Sun Oct 5 00:04:32 2025
Install Reason  : Installed as a dependency for another package
Install Script  : Yes
Validated By    : Signature

Name            : attr
Version         : 2.5.2-1
Description     : Extended attribute support library for ACL support
Architecture    : x86_64
URL             : https://savannah.nongnu.org/projects/attr
Licenses        : LGPL
Groups          : None
Provides        : xfsattr  libattr.so=1-64
Depends On      : glibc
Optional Deps   : None
Required By     : coreutils  gettext  shadow
Optional For    : None
Conflicts With  : xfsattr
Replaces        : xfsattr
Installed Size  : 216.00 KiB
Packager        : Christian Hesse <eworm@archlinux.org>
Build Date      : Sun Jan 14 14:58:15 2024
Install Date    : Sun Oct 5 00:04:30 2025
Install Reason  : Installed as a dependency for another package
Install Script  : No
Validated By    : Signature

Name            : audit
Version         : 4.0.5-1
Description     : Userspace components of the audit framework
Architecture    : x86_64
URL             : https://github.com/linux-audit/audit-userspace
Licenses        : GPL-2.0-or-later  LGPL-2.0-or-later
Groups          : None
Provides        : libaudit.so=1-64  libauparse.so=0-64
Depends On      : glibc  krb5  libkrb5.so=3-64  libgssapi_krb5.so=2-64  libcap-ng  libcap-ng.so=0-64
Optional Deps   : libldap: for audispd-zos-remote [installed]
                  sh: for augenrules [installed]
Required By     : dbus  dbus-broker  pam  shadow  systemd
Optional For    : None
Conflicts With  : None
Replaces        : None
Installed Size  : 1082.46 KiB
Packager        : David Runge <dvzrv@archlinux.org>
Build Date      : Wed Jun 4 08:08:34 2025
Install Date    : Sun Oct 5 00:04:30 2025
Install Reason  : Installed as a dependency for another package
Install Script  : No
Validated By    : Signature

Name            : base
Version         : 3-2
Description     : Minimal package set to define a basic Arch Linux installation
Architecture    : any
URL             : https://www.archlinux.org
Licenses        : GPL
Groups          : None
Provides        : None
Depends On      : filesystem  gcc-libs  glibc  bash  coreutils  file  findutils  gawk  grep  procps-ng  sed  tar  gettext  pciutils  psmisc  shadow  util-linux  bzip2  gzip  xz  licenses  pacman  archlinux-keyring  systemd  systemd-sysvcompat  iputils
                  iproute2
Optional Deps   : linux: bare metal support
Required By     : None
Optional For    : None
Conflicts With  : None
Replaces        : None
Installed Size  : 0.00 B
Packager        : Jan Alexander Steffens (heftig) <heftig@archlinux.org>
Build Date      : Sun Oct 8 01:15:41 2023
Install Date    : Sun Oct 5 00:04:32 2025
Install Reason  : Explicitly installed
Install Script  : No
Validated By    : Signature

Name            : bash
Version         : 5.3.3-2
Description     : The GNU Bourne Again shell
Architecture    : x86_64
URL             : https://www.gnu.org/software/bash/bash.html
Licenses        : GPL-3.0-or-later
Groups          : None
Provides        : sh
Depends On      : readline  libreadline.so=8-64  glibc  ncurses
Optional Deps   : bash-completion: for tab completion
Required By     : base  bzip2  ca-certificates-utils  e2fsprogs  gawk  gdbm  gettext  glib2  gnupg  gpm  gzip  icu  iptables  keyutils  krb5  libassuan  libgpg-error  libksba  libpcap  libxml2  mc  npth  pacman  pcre  systemd  which  xz
Optional For    : audit  ncurses  pcre2  tzdata  vim-runtime
Conflicts With  : None
Replaces        : None
Installed Size  : 9.56 MiB
Packager        : Tobias Powalowski <tpowa@archlinux.org>
Build Date      : Fri Aug 1 19:56:35 2025
Install Date    : Sun Oct 5 00:04:30 2025
Install Reason  : Installed as a dependency for another package
Install Script  : Yes
Validated By    : Signature

Name            : brotli
Version         : 1.1.0-3
Description     : Generic-purpose lossless compression algorithm
Architecture    : x86_64
URL             : https://github.com/google/brotli
Licenses        : MIT
Groups          : None
Provides        : libbrotlicommon.so=1-64  libbrotlidec.so=1-64  libbrotlienc.so=1-64
Depends On      : glibc
Optional Deps   : None
Required By     : curl  gnutls
Optional For    : None
Conflicts With  : None
Replaces        : None
Installed Size  : 1028.73 KiB
Packager        : Jelle van der Waa <jelle@archlinux.org>
Build Date      : Sun Nov 17 14:21:20 2024
Install Date    : Sun Oct 5 00:04:31 2025
Install Reason  : Installed as a dependency for another package
Install Script  : No
Validated By    : Signature

Name            : bzip2
Version         : 1.0.8-6
Description     : A high-quality data compression program
Architecture    : x86_64
URL             : https://sourceware.org/bzip2/
Licenses        : BSD
Groups          : None
Provides        : libbz2.so=1.0-64
Depends On      : glibc  sh
Optional Deps   : None
Required By     : base  file  gnupg  libarchive  libelf  pcre  pcre2
Optional For    : None
Conflicts With  : None
Replaces        : None
Installed Size  : 145.00 KiB
Packager        : Christian Hesse <eworm@archlinux.org>
Build Date      : Sun Mar 17 22:29:13 2024
Install Date    : Sun Oct 5 00:04:30 2025
Install Reason  : Installed as a dependency for another package
Install Script  : No
Validated By    : Signature

Name            : ca-certificates
Version         : 20240618-1
Description     : Common CA certificates - default providers
Architecture    : any
URL             : https://src.fedoraproject.org/rpms/ca-certificates
Licenses        : GPL-2.0-or-later
Groups          : None
Provides        : None
Depends On      : ca-certificates-mozilla
Optional Deps   : None
Required By     : curl
Optional For    : openssl
Conflicts With  : ca-certificates-cacert<=20140824-4
Replaces        : ca-certificates-cacert<=20140824-4
Installed Size  : 0.00 B
Packager        : Jan Alexander Steffens (heftig) <heftig@archlinux.org>
Build Date      : Tue Jun 18 18:36:40 2024
Install Date    : Sun Oct 5 00:04:31 2025
Install Reason  : Installed as a dependency for another package
Install Script  : No
Validated By    : Signature

Name            : ca-certificates-mozilla
Version         : 3.117-1
Description     : Mozilla's set of trusted CA certificates
Architecture    : x86_64
URL             : https://developer.mozilla.org/en-US/docs/Mozilla/Projects/NSS
Licenses        : MPL-2.0
Groups          : None
Provides        : None
Depends On      : ca-certificates-utils>=20181109-3
Optional Deps   : None
Required By     : ca-certificates
Optional For    : None
Conflicts With  : None
Replaces        : None
Installed Size  : 1110.87 KiB
Packager        : Jan Alexander Steffens (heftig) <heftig@archlinux.org>
Build Date      : Fri Oct 3 13:17:47 2025
Install Date    : Sun Oct 5 00:04:31 2025
Install Reason  : Installed as a dependency for another package
Install Script  : No
Validated By    : Signature

Name            : ca-certificates-utils
Version         : 20240618-1
Description     : Common CA certificates (utilities)
Architecture    : any
URL             : https://src.fedoraproject.org/rpms/ca-certificates
Licenses        : GPL-2.0-or-later
Groups          : None
Provides        : ca-certificates  ca-certificates-java
Depends On      : bash  coreutils  findutils  p11-kit
Optional Deps   : None
Required By     : ca-certificates-mozilla  curl
Optional For    : openssl
Conflicts With  : ca-certificates-java
Replaces        : ca-certificates-java
Installed Size  : 13.63 KiB
Packager        : Jan Alexander Steffens (heftig) <heftig@archlinux.org>
Build Date      : Tue Jun 18 18:36:40 2024
Install Date    : Sun Oct 5 00:04:31 2025
Install Reason  : Installed as a dependency for another package
Install Script  : Yes
Validated By    : Signature

Name            : coreutils
Version         : 9.8-2
Description     : The basic file, shell and text manipulation utilities of the GNU operating system
Architecture    : x86_64
URL             : https://www.gnu.org/software/coreutils/
Licenses        : GPL-3.0-or-later  GFDL-1.3-or-later
Groups          : None
Provides        : None
Depends On      : acl  attr  glibc  gmp  libcap  openssl
Optional Deps   : None
Required By     : base  ca-certificates-utils  gzip  p11-kit  pacman  util-linux
Optional For    : None
Conflicts With  : None
Replaces        : None
Installed Size  : 16.91 MiB
Packager        : Tobias Powalowski <tpowa@archlinux.org>
Build Date      : Thu Sep 25 11:39:58 2025
Install Date    : Sun Oct 5 00:04:30 2025
Install Reason  : Installed as a dependency for another package
Install Script  : No
Validated By    : Signature

Name            : cryptsetup
Version         : 2.8.1-1
Description     : Userspace setup tool for transparent encryption of block devices using dm-crypt
Architecture    : x86_64
URL             : https://gitlab.com/cryptsetup/cryptsetup/
Licenses        : GPL-2.0-or-later
Groups          : None
Provides        : libcryptsetup.so=12-64
Depends On      : device-mapper  libdevmapper.so=1.02-64  glibc  openssl  libcrypto.so=3-64  popt  util-linux-libs  libblkid.so=1-64  libuuid.so=1-64  json-c  libjson-c.so=5-64
Optional Deps   : None
Required By     : systemd
Optional For    : None
Conflicts With  : mkinitcpio<38-1
Replaces        : None
Installed Size  : 3.27 MiB
Packager        : Christian Hesse <eworm@archlinux.org>
Build Date      : Tue Aug 19 13:28:47 2025
Install Date    : Sun Oct 5 00:04:31 2025
Install Reason  : Installed as a dependency for another package
Install Script  : No
Validated By    : Signature

Name            : curl
Version         : 8.16.0-1
Description     : command line tool and library for transferring data with URLs
Architecture    : x86_64
URL             : https://curl.se/
Licenses        : MIT
Groups          : None
Provides        : libcurl.so=4-64
Depends On      : ca-certificates  brotli  libbrotlidec.so=1-64  krb5  libgssapi_krb5.so=2-64  libidn2  libidn2.so=0-64  libnghttp2  libnghttp2.so=14-64  libnghttp3  libnghttp3.so=9-64  libpsl  libpsl.so=5-64  libssh2  libssh2.so=1-64  zlib
                  libz.so=1-64  zstd  libzstd.so=1-64  openssl  libcrypto.so=3-64  libssl.so=3-64
Optional Deps   : None
Required By     : libelf  pacman  tpm2-tss
Optional For    : pciutils  systemd
Conflicts With  : wcurl
Replaces        : wcurl
Installed Size  : 2.07 MiB
Packager        : Christian Hesse <eworm@archlinux.org>
Build Date      : Wed Sep 10 06:04:11 2025
Install Date    : Sun Oct 5 00:04:31 2025
Install Reason  : Installed as a dependency for another package
Install Script  : No
Validated By    : Signature

Name            : dbus
Version         : 1.16.2-1
Description     : Freedesktop.org message bus system
Architecture    : x86_64
URL             : https://www.freedesktop.org/wiki/Software/dbus/
Licenses        : AFL-2.1 OR GPL-2.0-or-later
Groups          : None
Provides        : libdbus  libdbus-1.so=3-64
Depends On      : audit  expat  glibc  libcap-ng  systemd-libs  libaudit.so=1-64  libcap-ng.so=0-64  libexpat.so=1-64  libsystemd.so=0-64
Optional Deps   : None
Required By     : dbus-broker-units  libpcap  systemd
Optional For    : None
Conflicts With  : libdbus
Replaces        : libdbus
Installed Size  : 1004.89 KiB
Packager        : Jan Alexander Steffens (heftig) <heftig@archlinux.org>
Build Date      : Sun Mar 2 01:56:01 2025
Install Date    : Sun Oct 5 00:04:31 2025
Install Reason  : Installed as a dependency for another package
Install Script  : No
Validated By    : Signature
`

// LANG=c checkupdates
var checkupdatesSampleOutput = `archlinux-keyring 20251116-1 -> 20260107-2
audit 4.1.2-1 -> 4.1.2-2
binutils 2.45.1-1 -> 2.45.1+r35+g12d0a1dbc1b9-1
brotli 1.1.0-3 -> 1.2.0-1
curl 8.17.0-2 -> 8.18.0-2
gcc-libs 15.2.1+r301+gf24307422d1d-1 -> 15.2.1+r447+g6a64f6c3ebb8-1
gdbm 1.26-1 -> 1.26-2
glib2 2.86.3-1 -> 2.86.3-2
glibc 2.42+r33+gde1fe81f4714-1 -> 2.42+r47+ga1d3294a5bed-1
icu 78.1-1 -> 78.2-1
libarchive 3.8.4-1 -> 3.8.5-1
libcap-ng 0.8.5-3 -> 0.8.5-4
libpcap 1.10.5-3 -> 1.10.6-1
libseccomp 2.5.6-1 -> 2.6.0-1
libtasn1 4.20.0-1 -> 4.21.0-1
libxml2 2.15.1-4 -> 2.15.1-5
linux-api-headers 6.17-1 -> 6.18-1
sqlite 3.51.1-1 -> 3.51.2-1
systemd 259-1 -> 259-2
systemd-libs 259-1 -> 259-2
systemd-sysvcompat 259-1 -> 259-2
util-linux 2.41.3-1 -> 2.41.3-2
util-linux-libs 2.41.3-1 -> 2.41.3-2`

func TestParsePacmanQiOutput(t *testing.T) {
	mgr := PacmanManager{}
	packages, err := mgr.parsePacmanQiOutput(samplePacmanQiOutput)
	if err != nil {
		t.Fatalf("parsePacmanQiOutput returned error: %v", err)
	}

	// Check a few known packages from the sample output
	expected := []Package{
		{Name: "acl", Version: "2.3.2-1", Description: "Access control list utilities, libraries and headers"},
		{Name: "archlinux-keyring", Version: "20250929-1", Description: "Arch Linux PGP keyring"},
		{Name: "attr", Version: "2.5.2-1", Description: "Extended attribute support library for ACL support"},
		{Name: "audit", Version: "4.0.5-1", Description: "Userspace components of the audit framework"},
		{Name: "base", Version: "3-2", Description: "Minimal package set to define a basic Arch Linux installation"},
		{Name: "bash", Version: "5.3.3-2", Description: "The GNU Bourne Again shell"},
		{Name: "brotli", Version: "1.1.0-3", Description: "Generic-purpose lossless compression algorithm"},
		{Name: "bzip2", Version: "1.0.8-6", Description: "A high-quality data compression program"},
		{Name: "ca-certificates", Version: "20240618-1", Description: "Common CA certificates - default providers"},
		{Name: "ca-certificates-mozilla", Version: "3.117-1", Description: "Mozilla's set of trusted CA certificates"},
		{Name: "ca-certificates-utils", Version: "20240618-1", Description: "Common CA certificates (utilities)"},
		{Name: "coreutils", Version: "9.8-2", Description: "The basic file, shell and text manipulation utilities of the GNU operating system"},
	}

	if len(packages) < len(expected) {
		t.Errorf("Expected at least %d packages, got %d", len(expected), len(packages))
	}

	for i, exp := range expected {
		if i >= len(packages) {
			t.Errorf("Missing package at index %d: %+v", i, exp)
			continue
		}
		got := packages[i]
		if got.Name != exp.Name || got.Version != exp.Version || got.Description != exp.Description {
			t.Errorf("Package %d mismatch:\nGot:  %+v\nWant: %+v", i, got, exp)
		}
	}
}

func TestGetUpgradablePackagesForRc127(t *testing.T) {
	// This test will fail on Arch linux systems with installed pacman-contrib.
	pacman := PacmanManager{}

	expactedErrorMsg := "checkupdates command not found; please install the 'pacman-contrib' package"

	packages, err := pacman.getUpgradablePackages(context.Background())
	if err == nil {
		t.Fatalf("Expected error but got none, packages: %+v", packages)
	}

	if err.Error() != expactedErrorMsg {
		t.Fatalf("Unexpected error message:\nGot:  %s\nWant: %s", err.Error(), expactedErrorMsg)
	}
}

// Test for parsePacmanUpgradeOutput using checkupdatesSampleOutput
func TestParsePacmanUpgradeOutput(t *testing.T) {
	mgr := PacmanManager{}
	updates, err := mgr.parsePacmanUpgradeOutput(checkupdatesSampleOutput)
	if err != nil {
		t.Fatalf("parsePacmanUpgradeOutput returned error: %v", err)
	}

	expected := []PackageUpdate{
		{Name: "archlinux-keyring", CurrentVersion: "20251116-1", AvailableVersion: "20260107-2", IsSecurityUpdate: false},
		{Name: "audit", CurrentVersion: "4.1.2-1", AvailableVersion: "4.1.2-2", IsSecurityUpdate: false},
		{Name: "binutils", CurrentVersion: "2.45.1-1", AvailableVersion: "2.45.1+r35+g12d0a1dbc1b9-1", IsSecurityUpdate: false},
		{Name: "brotli", CurrentVersion: "1.1.0-3", AvailableVersion: "1.2.0-1", IsSecurityUpdate: false},
		{Name: "curl", CurrentVersion: "8.17.0-2", AvailableVersion: "8.18.0-2", IsSecurityUpdate: false},
		{Name: "gcc-libs", CurrentVersion: "15.2.1+r301+gf24307422d1d-1", AvailableVersion: "15.2.1+r447+g6a64f6c3ebb8-1", IsSecurityUpdate: false},
		{Name: "gdbm", CurrentVersion: "1.26-1", AvailableVersion: "1.26-2", IsSecurityUpdate: false},
		{Name: "glib2", CurrentVersion: "2.86.3-1", AvailableVersion: "2.86.3-2", IsSecurityUpdate: false},
		{Name: "glibc", CurrentVersion: "2.42+r33+gde1fe81f4714-1", AvailableVersion: "2.42+r47+ga1d3294a5bed-1", IsSecurityUpdate: false},
		{Name: "icu", CurrentVersion: "78.1-1", AvailableVersion: "78.2-1", IsSecurityUpdate: false},
		{Name: "libarchive", CurrentVersion: "3.8.4-1", AvailableVersion: "3.8.5-1", IsSecurityUpdate: false},
		{Name: "libcap-ng", CurrentVersion: "0.8.5-3", AvailableVersion: "0.8.5-4", IsSecurityUpdate: false},
		{Name: "libpcap", CurrentVersion: "1.10.5-3", AvailableVersion: "1.10.6-1", IsSecurityUpdate: false},
		{Name: "libseccomp", CurrentVersion: "2.5.6-1", AvailableVersion: "2.6.0-1", IsSecurityUpdate: false},
		{Name: "libtasn1", CurrentVersion: "4.20.0-1", AvailableVersion: "4.21.0-1", IsSecurityUpdate: false},
		{Name: "libxml2", CurrentVersion: "2.15.1-4", AvailableVersion: "2.15.1-5", IsSecurityUpdate: false},
		{Name: "linux-api-headers", CurrentVersion: "6.17-1", AvailableVersion: "6.18-1", IsSecurityUpdate: false},
		{Name: "sqlite", CurrentVersion: "3.51.1-1", AvailableVersion: "3.51.2-1", IsSecurityUpdate: false},
		{Name: "systemd", CurrentVersion: "259-1", AvailableVersion: "259-2", IsSecurityUpdate: false},
		{Name: "systemd-libs", CurrentVersion: "259-1", AvailableVersion: "259-2", IsSecurityUpdate: false},
		{Name: "systemd-sysvcompat", CurrentVersion: "259-1", AvailableVersion: "259-2", IsSecurityUpdate: false},
		{Name: "util-linux", CurrentVersion: "2.41.3-1", AvailableVersion: "2.41.3-2", IsSecurityUpdate: false},
		{Name: "util-linux-libs", CurrentVersion: "2.41.3-1", AvailableVersion: "2.41.3-2", IsSecurityUpdate: false},
	}

	if len(updates) != len(expected) {
		t.Errorf("Expected %d updates, got %d", len(expected), len(updates))
	}

	for i, exp := range expected {
		if i >= len(updates) {
			t.Errorf("Missing update at index %d: %+v", i, exp)
			continue
		}
		got := updates[i]
		if got != exp {
			t.Errorf("Update %d mismatch:\nGot:  %+v\nWant: %+v", i, got, exp)
		}
	}
}
