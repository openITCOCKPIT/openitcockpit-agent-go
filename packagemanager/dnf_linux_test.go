package packagemanager

import (
	"strings"
	"testing"
)

// Sample output of the command
// LANG=c dnf list installed
var dnfListInstalledOutput = `Failed to set locale, defaulting to C.UTF-8
Installed Packages
acl.x86_64                                                                                                                    2.3.1-4.el9                                                                                                            @baseos
almalinux-gpg-keys.x86_64                                                                                                     9.7-1.el9                                                                                                              @baseos
almalinux-release.x86_64                                                                                                      9.7-1.el9                                                                                                              @baseos
almalinux-repos.x86_64                                                                                                        9.7-1.el9                                                                                                              @baseos
alternatives.x86_64                                                                                                           1.24-2.el9                                                                                                             @baseos
attr.x86_64                                                                                                                   2.5.1-3.el9                                                                                                            @baseos
audit-libs.x86_64                                                                                                             3.1.5-7.el9                                                                                                            @baseos
basesystem.noarch                                                                                                             11-13.el9                                                                                                              @baseos
bash.x86_64                                                                                                                   5.1.8-9.el9                                                                                                            @baseos
binutils.x86_64                                                                                                               2.35.2-67.el9_7.1                                                                                                      @baseos
binutils-gold.x86_64                                                                                                          2.35.2-67.el9_7.1                                                                                                      @baseos
bzip2-libs.x86_64                                                                                                             1.0.8-10.el9_5                                                                                                         @baseos
ca-certificates.noarch                                                                                                        2025.2.80_v9.0.305-91.el9                                                                                              @baseos
coreutils-single.x86_64                                                                                                       8.32-39.el9                                                                                                            @baseos
cracklib.x86_64                                                                                                               2.9.6-27.el9                                                                                                           @baseos
cracklib-dicts.x86_64                                                                                                         2.9.6-27.el9                                                                                                           @baseos
crypto-policies.noarch                                                                                                        20250905-1.git377cc42.el9_7                                                                                            @baseos
crypto-policies-scripts.noarch                                                                                                20250905-1.git377cc42.el9_7                                                                                            @baseos
curl-minimal.x86_64                                                                                                           7.76.1-34.el9                                                                                                          @baseos`

var dnfListInstalledOutputForUpdateTest = `Failed to set locale, defaulting to C.UTF-8
Installed Packages
ImageMagick7-djvu.x86_64                                                                                                  1:7.1.2.8-1.el9.remi                                                                                             @remi-safe
ImageMagick7-heic.x86_64                                                                                                  1:7.1.2.8-1.el9.remi                                                                                             @remi-safe
ImageMagick7-libs.x86_64                                                                                                  1:7.1.2.8-1.el9.remi                                                                                             @remi-safe
LibRaw.x86_64                                                                                                             0.21.1-1.el9                                                                                                     @appstream
ModemManager-glib.x86_64                                                                                                  1.20.2-1.el9                                                                                                     @baseos
NetworkManager.x86_64                                                                                                     1:1.52.0-7.el9_6                                                                                                 @baseos
NetworkManager-libnm.x86_64                                                                                               1:1.52.0-7.el9_6                                                                                                 @baseos
avahi-glib.x86_64                                                                                                         0.8-22.el9_6.1                                                                                                   @appstream
avahi-libs.x86_64                                                                                                         0.8-22.el9_6.1                                                                                                   @baseos
basesystem.noarch                                                                                                         11-13.el9                                                                                                        @System
bash.x86_64                                                                                                               5.1.8-9.el9                                                                                                      @System
bash-completion.noarch                                                                                                    1:2.11-5.el9                                                                                                     @baseos
binutils.x86_64                                                                                                           2.35.2-63.el9                                                                                                    @baseos
binutils-gold.x86_64                                                                                                      2.35.2-63.el9                                                                                                    @baseos
bluez-libs.x86_64                                                                                                         5.72-4.el9                                                                                                       @baseos
boost-program-options.x86_64                                                                                              1.75.0-10.el9                                                                                                    @appstream
bubblewrap.x86_64                                                                                                         0.4.1-8.el9_5                                                                                                    @baseos
bzip2.x86_64                                                                                                              1.0.8-10.el9_5                                                                                                   @baseos`

// LANG=c dnf check-update
var dnfAvailableUpdatesOutput = `Failed to set locale, defaulting to C.UTF-8

ImageMagick7-djvu.x86_64                                                                                                    1:7.1.2.12-1.el9.remi                                                                                           remi-safe
ImageMagick7-heic.x86_64                                                                                                    1:7.1.2.12-1.el9.remi                                                                                           remi-safe
ImageMagick7-libs.x86_64                                                                                                    1:7.1.2.12-1.el9.remi                                                                                           remi-safe
NetworkManager.x86_64                                                                                                       1:1.54.0-3.el9_7                                                                                                baseos
NetworkManager-libnm.x86_64                                                                                                 1:1.54.0-3.el9_7                                                                                                baseos
NetworkManager-team.x86_64                                                                                                  1:1.54.0-3.el9_7                                                                                                baseos
NetworkManager-tui.x86_64                                                                                                   1:1.54.0-3.el9_7                                                                                                baseos
almalinux-gpg-keys.x86_64                                                                                                   9.7-1.el9                                                                                                       baseos
almalinux-release.x86_64                                                                                                    9.7-1.el9                                                                                                       baseos
almalinux-repos.x86_64                                                                                                      9.7-1.el9                                                                                                       baseos
alsa-lib.x86_64                                                                                                             1.2.14-1.el9                                                                                                    appstream
annobin.x86_64                                                                                                              12.98-1.el9                                                                                                     appstream
audit.x86_64                                                                                                                3.1.5-7.el9                                                                                                     baseos
audit-libs.x86_64                                                                                                           3.1.5-7.el9                                                                                                     baseos
avahi-glib.x86_64                                                                                                           0.8-23.el9                                                                                                      appstream
avahi-libs.x86_64                                                                                                           0.8-23.el9                                                                                                      baseos
binutils.x86_64                                                                                                             2.35.2-67.el9_7.1                                                                                               baseos
binutils-gold.x86_64                                                                                                        2.35.2-67.el9_7.1                                                                                               baseos
bluez-libs.x86_64                                                                                                           5.83-2.el9                                                                                                      baseos
boost-program-options.x86_64                                                                                                1.75.0-13.el9_7                                                                                                 appstream
bubblewrap.x86_64                                                                                                           0.6.3-1.el9                                                                                                     baseos
ca-certificates.noarch                                                                                                      2025.2.80_v9.0.305-91.el9                                                                                       baseos
chrony.x86_64                                                                                                               4.6.1-2.el9                                                                                                     baseos
cloud-init.noarch                                                                                                           24.4-7.el9.alma.1                                                                                               appstream
cockpit-bridge.noarch                                                                                                       344-1.el9                                                                                                       baseos
cockpit-system.noarch                                                                                                       344-1.el9                                                                                                       baseos
cockpit-ws.x86_64                                                                                                           344-1.el9                                                                                                       baseos
container-selinux.noarch                                                                                                    4:2.240.0-3.el9_7                                                                                               appstream
containerd.io.x86_64                                                                                                        2.2.1-1.el9                                                                                                     docker-ce-stable
cpp.x86_64                                                                                                                  11.5.0-11.el9.alma.1                                                                                            appstream
cronie.x86_64                                                                                                               1.5.7-15.el9                                                                                                    baseos
cronie-anacron.x86_64                                                                                                       1.5.7-15.el9                                                                                                    baseos
crypto-policies.noarch                                                                                                      20250905-1.git377cc42.el9_7                                                                                     baseos
crypto-policies-scripts.noarch                                                                                              20250905-1.git377cc42.el9_7                                                                                     baseos
cryptsetup-libs.x86_64                                                                                                      2.7.2-4.el9                                                                                                     baseos
cups-libs.x86_64                                                                                                            1:2.3.3op2-34.el9_7                                                                                             baseos
curl.x86_64                                                                                                                 7.76.1-34.el9                                                                                                   baseos
cyrus-sasl-lib.x86_64                                                                                                       2.1.27-22.el9                                                                                                   baseos
device-mapper.x86_64                                                                                                        9:1.02.206-2.el9_7.1                                                                                            baseos
device-mapper-libs.x86_64                                                                                                   9:1.02.206-2.el9_7.1                                                                                            baseos
dnf.noarch                                                                                                                  4.14.0-31.el9.alma.1                                                                                            baseos
dnf-data.noarch                                                                                                             4.14.0-31.el9.alma.1                                                                                            baseos
dnf-plugins-core.noarch                                                                                                     4.3.0-24.el9_7                                                                                                  baseos
docker-buildx-plugin.x86_64                                                                                                 0.30.1-1.el9                                                                                                    docker-ce-stable
docker-ce.x86_64                                                                                                            3:29.1.3-1.el9                                                                                                  docker-ce-stable
docker-ce-cli.x86_64                                                                                                        1:29.1.3-1.el9                                                                                                  docker-ce-stable
docker-ce-rootless-extras.x86_64                                                                                            29.1.3-1.el9                                                                                                    docker-ce-stable
docker-compose-plugin.x86_64                                                                                                5.0.1-1.el9                                                                                                     docker-ce-stable
dracut.x86_64                                                                                                               057-104.git20250919.el9_7                                                                                       baseos
dracut-config-generic.x86_64                                                                                                057-104.git20250919.el9_7                                                                                       baseos
dracut-network.x86_64                                                                                                       057-104.git20250919.el9_7                                                                                       baseos
dracut-squash.x86_64                                                                                                        057-104.git20250919.el9_7                                                                                       baseos
dwz.x86_64                                                                                                                  0.16-1.el9                                                                                                      appstream
e2fsprogs.x86_64                                                                                                            1.46.5-8.el9                                                                                                    baseos
e2fsprogs-libs.x86_64                                                                                                       1.46.5-8.el9                                                                                                    baseos
efi-filesystem.noarch                                                                                                       6-4.el9                                                                                                         baseos
efi-srpm-macros.noarch                                                                                                      6-4.el9                                                                                                         appstream
elfutils-debuginfod-client.x86_64                                                                                           0.193-1.el9.alma.1                                                                                              baseos
elfutils-default-yama-scope.noarch                                                                                          0.193-1.el9.alma.1                                                                                              baseos
elfutils-libelf.x86_64                                                                                                      0.193-1.el9.alma.1                                                                                              baseos
elfutils-libs.x86_64                                                                                                        0.193-1.el9.alma.1                                                                                              baseos
emacs-filesystem.noarch                                                                                                     1:27.2-18.el9                                                                                                   appstream
ethtool.x86_64                                                                                                              2:6.15-2.el9                                                                                                    baseos
expat.x86_64                                                                                                                2.5.0-5.el9_7.1                                                                                                 baseos
fuse-overlayfs.x86_64                                                                                                       1.15-1.el9                                                                                                      appstream
fwupd.x86_64                                                                                                                1.9.31-1.el9.alma.1                                                                                             baseos
fwupd-plugin-flashrom.x86_64                                                                                                1.9.31-1.el9.alma.1                                                                                             appstream
gcc.x86_64                                                                                                                  11.5.0-11.el9.alma.1                                                                                            appstream
gcc-c++.x86_64                                                                                                              11.5.0-11.el9.alma.1                                                                                            appstream
gcc-plugin-annobin.x86_64                                                                                                   11.5.0-11.el9.alma.1                                                                                            appstream
glib2.x86_64                                                                                                                2.68.4-18.el9_7                                                                                                 baseos
glibc.x86_64                                                                                                                2.34-231.el9_7.2                                                                                                baseos
glibc-common.x86_64                                                                                                         2.34-231.el9_7.2                                                                                                baseos
glibc-devel.x86_64                                                                                                          2.34-231.el9_7.2                                                                                                appstream
glibc-gconv-extra.x86_64                                                                                                    2.34-231.el9_7.2                                                                                                baseos
glibc-headers.x86_64                                                                                                        2.34-231.el9_7.2                                                                                                appstream
glibc-minimal-langpack.x86_64                                                                                               2.34-231.el9_7.2                                                                                                baseos
gnutls.x86_64                                                                                                               3.8.3-9.el9                                                                                                     baseos
go-srpm-macros.noarch                                                                                                       3.6.0-12.el9_7                                                                                                  appstream
grub2-common.noarch                                                                                                         1:2.06-114.el9_7.alma.1                                                                                         baseos
grub2-efi-x64.x86_64                                                                                                        1:2.06-114.el9_7.alma.1                                                                                         baseos
grub2-pc.x86_64                                                                                                             1:2.06-114.el9_7.alma.1                                                                                         baseos
grub2-pc-modules.noarch                                                                                                     1:2.06-114.el9_7.alma.1                                                                                         baseos
grub2-tools.x86_64                                                                                                          1:2.06-114.el9_7.alma.1                                                                                         baseos
grub2-tools-efi.x86_64                                                                                                      1:2.06-114.el9_7.alma.1                                                                                         baseos
grub2-tools-extra.x86_64                                                                                                    1:2.06-114.el9_7.alma.1                                                                                         baseos
grub2-tools-minimal.x86_64                                                                                                  1:2.06-114.el9_7.alma.1                                                                                         baseos
grubby.x86_64                                                                                                               8.40-68.el9                                                                                                     baseos
gsettings-desktop-schemas.x86_64                                                                                            40.0-8.el9_7                                                                                                    baseos
gtk-update-icon-cache.x86_64                                                                                                3.24.31-8.el9                                                                                                   appstream
gtk3.x86_64                                                                                                                 3.24.31-8.el9                                                                                                   appstream
highway.x86_64                                                                                                              1.3.0-1.el9                                                                                                     epel
httpd-filesystem.noarch                                                                                                     2.4.62-7.el9_7.3                                                                                                appstream
hwdata.noarch                                                                                                               0.348-9.20.el9                                                                                                  baseos
hwloc-libs.x86_64                                                                                                           2.4.1-6.el9_7                                                                                                   baseos
ima-evm-utils.x86_64                                                                                                        1.6.2-2.el9                                                                                                     baseos
iproute.x86_64                                                                                                              6.14.0-2.el9                                                                                                    baseos
iputils.x86_64                                                                                                              20210202-15.el9_7                                                                                               baseos
irqbalance.x86_64                                                                                                           2:1.9.4-4.el9                                                                                                   baseos
jasper-libs.x86_64                                                                                                          2.0.28-4.el9                                                                                                    appstream
jq.x86_64                                                                                                                   1.6-19.el9                                                                                                      baseos
kernel.x86_64                                                                                                               5.14.0-611.16.1.el9_7                                                                                           baseos
kernel-core.x86_64                                                                                                          5.14.0-611.16.1.el9_7                                                                                           baseos
kernel-headers.x86_64                                                                                                       5.14.0-611.16.1.el9_7                                                                                           appstream
kernel-modules.x86_64                                                                                                       5.14.0-611.16.1.el9_7                                                                                           baseos
kernel-modules-core.x86_64                                                                                                  5.14.0-611.16.1.el9_7                                                                                           baseos
kernel-srpm-macros.noarch                                                                                                   1.0-14.el9                                                                                                      appstream
kernel-tools.x86_64                                                                                                         5.14.0-611.16.1.el9_7                                                                                           baseos
kernel-tools-libs.x86_64                                                                                                    5.14.0-611.16.1.el9_7                                                                                           baseos
kexec-tools.x86_64                                                                                                          2.0.29-10.el9                                                                                                   baseos
kmod.x86_64                                                                                                                 28-11.el9                                                                                                       baseos
kmod-libs.x86_64                                                                                                            28-11.el9                                                                                                       baseos
kpartx.x86_64                                                                                                               0.8.7-39.el9                                                                                                    baseos
less.x86_64                                                                                                                 590-6.el9                                                                                                       baseos
libatomic.x86_64                                                                                                            11.5.0-11.el9.alma.1                                                                                            baseos
libbpf.x86_64                                                                                                               2:1.5.0-2.el9                                                                                                   baseos
libcap.x86_64                                                                                                               2.48-10.el9                                                                                                     baseos
libcom_err.x86_64                                                                                                           1.46.5-8.el9                                                                                                    baseos
libcurl.x86_64                                                                                                              7.76.1-34.el9                                                                                                   baseos
libdav1d.x86_64                                                                                                             1.5.2-1.el9                                                                                                     epel
libdnf.x86_64                                                                                                               0.69.0-16.el9.alma.1                                                                                            baseos
libgcc.x86_64                                                                                                               11.5.0-11.el9.alma.1                                                                                            baseos
libgexiv2.x86_64                                                                                                            0.14.3-1.el9                                                                                                    appstream
libgomp.x86_64                                                                                                              11.5.0-11.el9.alma.1                                                                                            baseos
libibverbs.x86_64                                                                                                           57.0-2.el9                                                                                                      baseos
libldb.x86_64                                                                                                               4.22.4-6.el9_7                                                                                                  baseos
libnfsidmap.x86_64                                                                                                          1:2.5.4-38.el9                                                                                                  baseos
librepo.x86_64                                                                                                              1.14.5-3.el9                                                                                                    baseos
libsepol.x86_64                                                                                                             3.6-3.el9                                                                                                       baseos
libsoup.x86_64                                                                                                              2.72.0-12.el9_7.1                                                                                               appstream
libss.x86_64                                                                                                                1.46.5-8.el9                                                                                                    baseos
libssh.x86_64                                                                                                               0.10.4-17.el9_7                                                                                                 baseos
libssh-config.noarch                                                                                                        0.10.4-17.el9_7                                                                                                 baseos
libsss_certmap.x86_64                                                                                                       2.9.7-4.el9_7.1                                                                                                 baseos
libsss_idmap.x86_64                                                                                                         2.9.7-4.el9_7.1                                                                                                 baseos
libsss_nss_idmap.x86_64                                                                                                     2.9.7-4.el9_7.1                                                                                                 baseos
libsss_sudo.x86_64                                                                                                          2.9.7-4.el9_7.1                                                                                                 baseos
libstdc++.x86_64                                                                                                            11.5.0-11.el9.alma.1                                                                                            baseos
libstdc++-devel.x86_64                                                                                                      11.5.0-11.el9.alma.1                                                                                            appstream
libsysfs.x86_64                                                                                                             2.1.1-11.el9                                                                                                    baseos
libtalloc.x86_64                                                                                                            2.4.3-1.el9                                                                                                     baseos
libtdb.x86_64                                                                                                               1.4.13-1.el9                                                                                                    baseos
libtevent.x86_64                                                                                                            0.16.2-1.el9                                                                                                    baseos
libtiff.x86_64                                                                                                              4.4.0-15.el9_7.2                                                                                                appstream
libtraceevent.x86_64                                                                                                        1.8.4-2.el9                                                                                                     baseos
libuser.x86_64                                                                                                              0.63-17.el9                                                                                                     baseos
libxml2.x86_64                                                                                                              2.9.13-14.el9_7                                                                                                 baseos
linux-firmware-whence.noarch                                                                                                20251111-155.1.el9_7                                                                                            baseos
llvm-libs.x86_64                                                                                                            20.1.8-3.el9                                                                                                    appstream
logrotate.x86_64                                                                                                            3.18.0-12.el9                                                                                                   baseos
lshw.x86_64                                                                                                                 B.02.20-2.el9                                                                                                   baseos
man-db.x86_64                                                                                                               2.9.3-9.el9                                                                                                     baseos
mesa-dri-drivers.x86_64                                                                                                     25.0.7-3.el9_7.alma.1                                                                                           appstream
mesa-filesystem.x86_64                                                                                                      25.0.7-3.el9_7.alma.1                                                                                           appstream
mesa-libEGL.x86_64                                                                                                          25.0.7-3.el9_7.alma.1                                                                                           appstream
mesa-libGL.x86_64                                                                                                           25.0.7-3.el9_7.alma.1                                                                                           appstream
mesa-libgbm.x86_64                                                                                                          25.0.7-3.el9_7.alma.1                                                                                           appstream
microcode_ctl.noarch                                                                                                        4:20250812-1.el9                                                                                                baseos
mokutil.x86_64                                                                                                              2:0.7.2-1.el9                                                                                                   baseos
mysql.x86_64                                                                                                                8.0.44-1.el9_7                                                                                                  appstream
mysql-common.x86_64                                                                                                         8.0.44-1.el9_7                                                                                                  appstream
mysql-errmsg.x86_64                                                                                                         8.0.44-1.el9_7                                                                                                  appstream
mysql-libs.x86_64                                                                                                           8.0.44-1.el9_7                                                                                                  appstream
mysql-server.x86_64                                                                                                         8.0.44-1.el9_7                                                                                                  appstream
ncurses.x86_64                                                                                                              6.2-12.20210508.el9                                                                                             baseos
ncurses-base.noarch                                                                                                         6.2-12.20210508.el9                                                                                             baseos
ncurses-libs.x86_64                                                                                                         6.2-12.20210508.el9                                                                                             baseos
nfs-utils.x86_64                                                                                                            1:2.5.4-38.el9                                                                                                  baseos
numactl-libs.x86_64                                                                                                         2.0.19-3.el9                                                                                                    baseos
openitcockpit.x86_64                                                                                                        5.3.1_20260105094727RHEL9-1.RHEL9                                                                               openitcockpit
openitcockpit-check-microsoft365-plugin.x86_64                                                                              1.0.2_20260105100508RHEL9-1.RHEL9                                                                               openitcockpit
openitcockpit-check-proxmox-plugin.x86_64                                                                                   1.0.2_20260105100425RHEL9-1.RHEL9                                                                               openitcockpit
openitcockpit-common.x86_64                                                                                                 5.3.1_20260105091347RHEL9-1.RHEL9                                                                               openitcockpit
openitcockpit-frontend-angular.x86_64                                                                                       5.3.1_20260105082225generic-1.generic                                                                           openitcockpit
openitcockpit-graphing.x86_64                                                                                               5.3.1_20260105081347generic-1.generic                                                                           openitcockpit
openitcockpit-module-distribute.x86_64                                                                                      5.3.1_20260105091347RHEL9-1.RHEL9                                                                               openitcockpit
openitcockpit-module-grafana.x86_64                                                                                         5.3.1_20260105091347RHEL9-1.RHEL9                                                                               openitcockpit
openitcockpit-module-jira.x86_64                                                                                            5.3.1_20260105091347RHEL9-1.RHEL9                                                                               openitcockpit
openitcockpit-naemon.x86_64                                                                                                 1.4.4_20260105094423RHEL9-1.RHEL9                                                                               openitcockpit
openitcockpit-node.x86_64                                                                                                   5.3.1_20260105082100generic-1.generic                                                                           openitcockpit
openitcockpit-nsta.x86_64                                                                                                   5.3.1_20260105094640RHEL9-1.RHEL9                                                                               openitcockpit
openitcockpit-selinux.x86_64                                                                                                5.3.1_20260105100010RHEL9-1.RHEL9                                                                               openitcockpit
openitcockpit-statusengine3-worker.x86_64                                                                                   3.8.2_20260105094922RHEL9-1.RHEL9                                                                               openitcockpit
openssh.x86_64                                                                                                              8.7p1-47.el9_7.alma.1                                                                                           baseos
openssh-clients.x86_64                                                                                                      8.7p1-47.el9_7.alma.1                                                                                           baseos
openssh-server.x86_64                                                                                                       8.7p1-47.el9_7.alma.1                                                                                           baseos
openssl.x86_64                                                                                                              1:3.5.1-4.el9_7                                                                                                 baseos
openssl-libs.x86_64                                                                                                         1:3.5.1-4.el9_7                                                                                                 baseos
osinfo-db.noarch                                                                                                            20250606-1.el9                                                                                                  appstream
ostree-libs.x86_64                                                                                                          2025.6-1.el9_7                                                                                                  appstream
pcp-conf.x86_64                                                                                                             6.3.7-5.el9                                                                                                     appstream
pcp-libs.x86_64                                                                                                             6.3.7-5.el9                                                                                                     appstream
perl-Net-SSLeay.x86_64                                                                                                      1.94-3.el9                                                                                                      appstream
php-bcmath.x86_64                                                                                                           8.1.34-1.module_php.8.1.el9.remi                                                                                remi-modular
php-cli.x86_64                                                                                                              8.1.34-1.module_php.8.1.el9.remi                                                                                remi-modular
php-common.x86_64                                                                                                           8.1.34-1.module_php.8.1.el9.remi                                                                                remi-modular
php-fpm.x86_64                                                                                                              8.1.34-1.module_php.8.1.el9.remi                                                                                remi-modular
php-gd.x86_64                                                                                                               8.1.34-1.module_php.8.1.el9.remi                                                                                remi-modular
php-imap.x86_64                                                                                                             8.1.34-1.module_php.8.1.el9.remi                                                                                remi-modular
php-intl.x86_64                                                                                                             8.1.34-1.module_php.8.1.el9.remi                                                                                remi-modular
php-ldap.x86_64                                                                                                             8.1.34-1.module_php.8.1.el9.remi                                                                                remi-modular
php-mbstring.x86_64                                                                                                         8.1.34-1.module_php.8.1.el9.remi                                                                                remi-modular
php-mysqlnd.x86_64                                                                                                          8.1.34-1.module_php.8.1.el9.remi                                                                                remi-modular
php-opcache.x86_64                                                                                                          8.1.34-1.module_php.8.1.el9.remi                                                                                remi-modular
php-pdo.x86_64                                                                                                              8.1.34-1.module_php.8.1.el9.remi                                                                                remi-modular
php-pecl-imagick-im7.x86_64                                                                                                 3.8.1-1.module_php.8.1.el9.remi                                                                                 remi-modular
php-pecl-redis6.x86_64                                                                                                      6.3.0-1.module_php.8.1.el9.remi                                                                                 remi-modular
php-process.x86_64                                                                                                          8.1.34-1.module_php.8.1.el9.remi                                                                                remi-modular
php-soap.x86_64                                                                                                             8.1.34-1.module_php.8.1.el9.remi                                                                                remi-modular
php-xml.x86_64                                                                                                              8.1.34-1.module_php.8.1.el9.remi                                                                                remi-modular
policycoreutils.x86_64                                                                                                      3.6-3.el9                                                                                                       baseos
policycoreutils-python-utils.noarch                                                                                         3.6-3.el9                                                                                                       appstream
polkit.x86_64                                                                                                               0.117-14.el9                                                                                                    baseos
polkit-libs.x86_64                                                                                                          0.117-14.el9                                                                                                    baseos
poppler.x86_64                                                                                                              21.01.0-23.el9_7                                                                                                appstream
poppler-glib.x86_64                                                                                                         21.01.0-23.el9_7                                                                                                appstream
python-unversioned-command.noarch                                                                                           3.9.25-2.el9_7                                                                                                  appstream
python3.x86_64                                                                                                              3.9.25-2.el9_7                                                                                                  baseos
python3-audit.x86_64                                                                                                        3.1.5-7.el9                                                                                                     appstream
python3-dasbus.noarch                                                                                                       1.5-1.el9                                                                                                       appstream
python3-dateutil.noarch                                                                                                     1:2.9.0.post0-1.el9_7                                                                                           baseos
python3-dnf.noarch                                                                                                          4.14.0-31.el9.alma.1                                                                                            baseos
python3-dnf-plugin-post-transaction-actions.noarch                                                                          4.3.0-24.el9_7                                                                                                  baseos
python3-dnf-plugins-core.noarch                                                                                             4.3.0-24.el9_7                                                                                                  baseos
python3-hawkey.x86_64                                                                                                       0.69.0-16.el9.alma.1                                                                                            baseos
python3-libdnf.x86_64                                                                                                       0.69.0-16.el9.alma.1                                                                                            baseos
python3-libs.x86_64                                                                                                         3.9.25-2.el9_7                                                                                                  baseos
python3-libxml2.x86_64                                                                                                      2.9.13-14.el9_7                                                                                                 baseos
python3-perf.x86_64                                                                                                         5.14.0-611.16.1.el9_7                                                                                           appstream
python3-policycoreutils.noarch                                                                                              3.6-3.el9                                                                                                       appstream
python3-rpm.x86_64                                                                                                          4.16.1.3-39.el9                                                                                                 baseos
python3-setuptools.noarch                                                                                                   53.0.0-15.el9                                                                                                   baseos
python3-setuptools-wheel.noarch                                                                                             53.0.0-15.el9                                                                                                   baseos
qemu-guest-agent.x86_64                                                                                                     17:9.1.0-29.el9_7.3.alma.1                                                                                      appstream
redhat-rpm-config.noarch                                                                                                    210-1.el9.alma.1                                                                                                appstream
redis.x86_64                                                                                                                6.2.20-2.el9_7                                                                                                  appstream
remi-release.noarch                                                                                                         9.7-4.el9.remi                                                                                                  remi-safe
rootfiles.noarch                                                                                                            8.1-35.el9                                                                                                      baseos
rpm.x86_64                                                                                                                  4.16.1.3-39.el9                                                                                                 baseos
rpm-build-libs.x86_64                                                                                                       4.16.1.3-39.el9                                                                                                 baseos
rpm-libs.x86_64                                                                                                             4.16.1.3-39.el9                                                                                                 baseos
rpm-plugin-audit.x86_64                                                                                                     4.16.1.3-39.el9                                                                                                 baseos
rpm-plugin-selinux.x86_64                                                                                                   4.16.1.3-39.el9                                                                                                 baseos
rpm-plugin-systemd-inhibit.x86_64                                                                                           4.16.1.3-39.el9                                                                                                 appstream
rpm-sign-libs.x86_64                                                                                                        4.16.1.3-39.el9                                                                                                 baseos
rsyslog.x86_64                                                                                                              8.2506.0-2.el9                                                                                                  appstream
rsyslog-logrotate.x86_64                                                                                                    8.2506.0-2.el9                                                                                                  appstream
selinux-policy.noarch                                                                                                       38.1.65-1.el9                                                                                                   baseos
selinux-policy-targeted.noarch                                                                                              38.1.65-1.el9                                                                                                   baseos
setroubleshoot-server.x86_64                                                                                                3.3.35-2.el9                                                                                                    appstream
shadow-utils.x86_64                                                                                                         2:4.9-15.el9                                                                                                    baseos
slirp4netns.x86_64                                                                                                          1.3.3-1.el9                                                                                                     appstream
sos.noarch                                                                                                                  4.10.1-2.el9                                                                                                    baseos
sqlite-libs.x86_64                                                                                                          3.34.1-9.el9_7                                                                                                  baseos
sscg.x86_64                                                                                                                 3.0.0-10.el9                                                                                                    appstream
sssd-client.x86_64                                                                                                          2.9.7-4.el9_7.1                                                                                                 baseos
sssd-common.x86_64                                                                                                          2.9.7-4.el9_7.1                                                                                                 baseos
sssd-kcm.x86_64                                                                                                             2.9.7-4.el9_7.1                                                                                                 baseos
sssd-nfs-idmap.x86_64                                                                                                       2.9.7-4.el9_7.1                                                                                                 baseos
sudo.x86_64                                                                                                                 1.9.5p2-13.el9                                                                                                  baseos
systemd.x86_64                                                                                                              252-55.el9_7.7.alma.1                                                                                           baseos
systemd-libs.x86_64                                                                                                         252-55.el9_7.7.alma.1                                                                                           baseos
systemd-pam.x86_64                                                                                                          252-55.el9_7.7.alma.1                                                                                           baseos
systemd-rpm-macros.noarch                                                                                                   252-55.el9_7.7.alma.1                                                                                           baseos
systemd-udev.x86_64                                                                                                         252-55.el9_7.7.alma.1                                                                                           baseos
systemtap-sdt-devel.x86_64                                                                                                  5.3-3.el9                                                                                                       appstream
tar.x86_64                                                                                                                  2:1.34-9.el9_7                                                                                                  baseos
tuned.noarch                                                                                                                2.26.0-1.el9                                                                                                    baseos
tzdata.noarch                                                                                                               2025c-1.el9                                                                                                     baseos
unzip.x86_64                                                                                                                6.0-59.el9                                                                                                      baseos
usermode.x86_64                                                                                                             1.114-7.el9                                                                                                     baseos
vim-common.x86_64                                                                                                           2:8.2.2637-23.el9_7                                                                                             appstream
vim-enhanced.x86_64                                                                                                         2:8.2.2637-23.el9_7                                                                                             appstream
vim-filesystem.noarch                                                                                                       2:8.2.2637-23.el9_7                                                                                             baseos
vim-minimal.x86_64                                                                                                          2:8.2.2637-23.el9_7                                                                                             baseos
virt-what.x86_64                                                                                                            1.27-2.el9                                                                                                      baseos
webkit2gtk3-jsc.x86_64                                                                                                      2.50.4-1.el9_7                                                                                                  appstream
xfsprogs.x86_64                                                                                                             6.4.0-7.el9                                                                                                     baseos
xorg-x11-server-Xvfb.x86_64                                                                                                 1.20.11-32.el9_7                                                                                                appstream
xorg-x11-server-common.x86_64                                                                                               1.20.11-32.el9_7                                                                                                appstream
yum.noarch                                                                                                                  4.14.0-31.el9.alma.1                                                                                            baseos
yum-utils.noarch                                                                                                            4.3.0-24.el9_7                                                                                                  baseos
Obsoleting Packages
grub2-tools.x86_64                                                                                                          1:2.06-114.el9_7.alma.1                                                                                         baseos
    grub2-tools.x86_64                                                                                                      1:2.06-104.el9_6.alma.1                                                                                         @baseos
grub2-tools-efi.x86_64                                                                                                      1:2.06-114.el9_7.alma.1                                                                                         baseos
    grub2-tools.x86_64                                                                                                      1:2.06-104.el9_6.alma.1                                                                                         @baseos
grub2-tools-extra.x86_64                                                                                                    1:2.06-114.el9_7.alma.1                                                                                         baseos
    grub2-tools.x86_64                                                                                                      1:2.06-104.el9_6.alma.1                                                                                         @baseos
grub2-tools-minimal.x86_64                                                                                                  1:2.06-114.el9_7.alma.1                                                                                         baseos
    grub2-tools.x86_64                                                                                                      1:2.06-104.el9_6.alma.1                                                                                         @baseos
mesa-dri-drivers.i686                                                                                                       25.0.7-3.el9_7.alma.1                                                                                           appstream
    mesa-libglapi.x86_64                                                                                                    24.2.8-3.el9_6.alma.1                                                                                           @appstream
mesa-dri-drivers.x86_64                                                                                                     25.0.7-3.el9_7.alma.1                                                                                           appstream
    mesa-libglapi.x86_64                                                                                                    24.2.8-3.el9_6.alma.1                                                                                           @appstream`

// dnf updateinfo list --available security
var availableSecurityUpdatesDnfOutput = `Last metadata expiration check: 0:52:58 ago on Wed Jan  7 18:13:57 2026.
ALSA-2025:23343 Moderate/Sec.  binutils-2.35.2-67.el9_7.1.x86_64
ALSA-2025:23343 Moderate/Sec.  binutils-gold-2.35.2-67.el9_7.1.x86_64
ALSA-2025:22175 Important/Sec. expat-2.5.0-5.el9_7.1.x86_64
ALSA-2025:22005 Moderate/Sec.  go-srpm-macros-3.6.0-12.el9_7.noarch
ALSA-2025:20532 Moderate/Sec.  grub2-common-1:2.06-114.el9_7.alma.1.noarch
ALSA-2025:20532 Moderate/Sec.  grub2-efi-x64-1:2.06-114.el9_7.alma.1.x86_64
ALSA-2025:20532 Moderate/Sec.  grub2-pc-1:2.06-114.el9_7.alma.1.x86_64
ALSA-2025:20532 Moderate/Sec.  grub2-pc-modules-1:2.06-114.el9_7.alma.1.noarch
ALSA-2025:20532 Moderate/Sec.  grub2-tools-1:2.06-114.el9_7.alma.1.x86_64
ALSA-2025:20532 Moderate/Sec.  grub2-tools-efi-1:2.06-114.el9_7.alma.1.x86_64
ALSA-2025:20532 Moderate/Sec.  grub2-tools-extra-1:2.06-114.el9_7.alma.1.x86_64
ALSA-2025:20532 Moderate/Sec.  grub2-tools-minimal-1:2.06-114.el9_7.alma.1.x86_64
ALSA-2025:23919 Important/Sec. httpd-filesystem-2.4.62-7.el9_7.3.noarch
ALSA-2025:19409 Moderate/Sec.  kernel-5.14.0-570.60.1.el9_6.x86_64
ALSA-2025:19930 Moderate/Sec.  kernel-5.14.0-570.62.1.el9_6.x86_64
ALSA-2025:22405 Moderate/Sec.  kernel-5.14.0-611.11.1.el9_7.x86_64
ALSA-2025:22865 Moderate/Sec.  kernel-5.14.0-611.13.1.el9_7.x86_64
ALSA-2025:23241 Important/Sec. kernel-5.14.0-611.16.1.el9_7.x86_64
ALSA-2025:20518 Moderate/Sec.  kernel-5.14.0-611.5.1.el9_7.x86_64
ALSA-2025:21926 Moderate/Sec.  kernel-5.14.0-611.9.1.el9_7.x86_64
ALSA-2025:19409 Moderate/Sec.  kernel-core-5.14.0-570.60.1.el9_6.x86_64
ALSA-2025:19930 Moderate/Sec.  kernel-core-5.14.0-570.62.1.el9_6.x86_64
ALSA-2025:22405 Moderate/Sec.  kernel-core-5.14.0-611.11.1.el9_7.x86_64
ALSA-2025:22865 Moderate/Sec.  kernel-core-5.14.0-611.13.1.el9_7.x86_64
ALSA-2025:23241 Important/Sec. kernel-core-5.14.0-611.16.1.el9_7.x86_64
ALSA-2025:20518 Moderate/Sec.  kernel-core-5.14.0-611.5.1.el9_7.x86_64
ALSA-2025:19409 Moderate/Sec.  kernel-headers-5.14.0-570.60.1.el9_6.x86_64
ALSA-2025:19930 Moderate/Sec.  kernel-headers-5.14.0-570.62.1.el9_6.x86_64
ALSA-2025:22405 Moderate/Sec.  kernel-headers-5.14.0-611.11.1.el9_7.x86_64
ALSA-2025:22865 Moderate/Sec.  kernel-headers-5.14.0-611.13.1.el9_7.x86_64
ALSA-2025:23241 Important/Sec. kernel-headers-5.14.0-611.16.1.el9_7.x86_64
ALSA-2025:20518 Moderate/Sec.  kernel-headers-5.14.0-611.5.1.el9_7.x86_64
ALSA-2025:21926 Moderate/Sec.  kernel-headers-5.14.0-611.9.1.el9_7.x86_64
ALSA-2025:19409 Moderate/Sec.  kernel-modules-5.14.0-570.60.1.el9_6.x86_64
ALSA-2025:19930 Moderate/Sec.  kernel-modules-5.14.0-570.62.1.el9_6.x86_64
ALSA-2025:22405 Moderate/Sec.  kernel-modules-5.14.0-611.11.1.el9_7.x86_64
ALSA-2025:22865 Moderate/Sec.  kernel-modules-5.14.0-611.13.1.el9_7.x86_64
ALSA-2025:23241 Important/Sec. kernel-modules-5.14.0-611.16.1.el9_7.x86_64
ALSA-2025:20518 Moderate/Sec.  kernel-modules-5.14.0-611.5.1.el9_7.x86_64
ALSA-2025:19409 Moderate/Sec.  kernel-modules-core-5.14.0-570.60.1.el9_6.x86_64
ALSA-2025:19930 Moderate/Sec.  kernel-modules-core-5.14.0-570.62.1.el9_6.x86_64
ALSA-2025:22405 Moderate/Sec.  kernel-modules-core-5.14.0-611.11.1.el9_7.x86_64
ALSA-2025:22865 Moderate/Sec.  kernel-modules-core-5.14.0-611.13.1.el9_7.x86_64
ALSA-2025:23241 Important/Sec. kernel-modules-core-5.14.0-611.16.1.el9_7.x86_64
ALSA-2025:20518 Moderate/Sec.  kernel-modules-core-5.14.0-611.5.1.el9_7.x86_64
ALSA-2025:19409 Moderate/Sec.  kernel-tools-5.14.0-570.60.1.el9_6.x86_64
ALSA-2025:19930 Moderate/Sec.  kernel-tools-5.14.0-570.62.1.el9_6.x86_64
ALSA-2025:22405 Moderate/Sec.  kernel-tools-5.14.0-611.11.1.el9_7.x86_64
ALSA-2025:22865 Moderate/Sec.  kernel-tools-5.14.0-611.13.1.el9_7.x86_64
ALSA-2025:23241 Important/Sec. kernel-tools-5.14.0-611.16.1.el9_7.x86_64
ALSA-2025:20518 Moderate/Sec.  kernel-tools-5.14.0-611.5.1.el9_7.x86_64
ALSA-2025:19409 Moderate/Sec.  kernel-tools-libs-5.14.0-570.60.1.el9_6.x86_64
ALSA-2025:19930 Moderate/Sec.  kernel-tools-libs-5.14.0-570.62.1.el9_6.x86_64
ALSA-2025:22405 Moderate/Sec.  kernel-tools-libs-5.14.0-611.11.1.el9_7.x86_64
ALSA-2025:22865 Moderate/Sec.  kernel-tools-libs-5.14.0-611.13.1.el9_7.x86_64
ALSA-2025:23241 Important/Sec. kernel-tools-libs-5.14.0-611.16.1.el9_7.x86_64
ALSA-2025:20518 Moderate/Sec.  kernel-tools-libs-5.14.0-611.5.1.el9_7.x86_64
ALSA-2025:20959 Important/Sec. libsoup-2.72.0-12.el9_7.1.x86_64
ALSA-2025:20943 Moderate/Sec.  libssh-0.10.4-15.el9_7.x86_64
ALSA-2025:23483 Moderate/Sec.  libssh-0.10.4-17.el9_7.x86_64
ALSA-2025:20943 Moderate/Sec.  libssh-config-0.10.4-15.el9_7.noarch
ALSA-2025:23483 Moderate/Sec.  libssh-config-0.10.4-17.el9_7.noarch
ALSA-2025:20954 Important/Sec. libsss_certmap-2.9.7-4.el9_7.1.x86_64
ALSA-2025:20954 Important/Sec. libsss_idmap-2.9.7-4.el9_7.1.x86_64
ALSA-2025:20954 Important/Sec. libsss_nss_idmap-2.9.7-4.el9_7.1.x86_64
ALSA-2025:20954 Important/Sec. libsss_sudo-2.9.7-4.el9_7.1.x86_64
ALSA-2025:20956 Important/Sec. libtiff-4.4.0-15.el9_7.2.x86_64
ALSA-2025:22376 Moderate/Sec.  libxml2-2.9.13-14.el9_7.x86_64
ALSA-2025:23109 Moderate/Sec.  mysql-8.0.44-1.el9_7.x86_64
ALSA-2025:23109 Moderate/Sec.  mysql-common-8.0.44-1.el9_7.x86_64
ALSA-2025:23109 Moderate/Sec.  mysql-errmsg-8.0.44-1.el9_7.x86_64
ALSA-2025:23109 Moderate/Sec.  mysql-libs-8.0.44-1.el9_7.x86_64
ALSA-2025:23109 Moderate/Sec.  mysql-server-8.0.44-1.el9_7.x86_64
ALSA-2025:23480 Moderate/Sec.  openssh-8.7p1-47.el9_7.alma.1.x86_64
ALSA-2025:23480 Moderate/Sec.  openssh-clients-8.7p1-47.el9_7.alma.1.x86_64
ALSA-2025:23480 Moderate/Sec.  openssh-server-8.7p1-47.el9_7.alma.1.x86_64
ALSA-2025:21255 Moderate/Sec.  openssl-1:3.5.1-4.el9_7.x86_64
ALSA-2025:21255 Moderate/Sec.  openssl-libs-1:3.5.1-4.el9_7.x86_64
ALSA-2025:23342 Moderate/Sec.  python-unversioned-command-3.9.25-2.el9_7.noarch
ALSA-2025:23342 Moderate/Sec.  python3-3.9.25-2.el9_7.x86_64
ALSA-2025:23342 Moderate/Sec.  python3-libs-3.9.25-2.el9_7.x86_64
ALSA-2025:22376 Moderate/Sec.  python3-libxml2-2.9.13-14.el9_7.x86_64
ALSA-2025:19409 Moderate/Sec.  python3-perf-5.14.0-570.60.1.el9_6.x86_64
ALSA-2025:19930 Moderate/Sec.  python3-perf-5.14.0-570.62.1.el9_6.x86_64
ALSA-2025:22405 Moderate/Sec.  python3-perf-5.14.0-611.11.1.el9_7.x86_64
ALSA-2025:22865 Moderate/Sec.  python3-perf-5.14.0-611.13.1.el9_7.x86_64
ALSA-2025:23241 Important/Sec. python3-perf-5.14.0-611.16.1.el9_7.x86_64
ALSA-2025:20518 Moderate/Sec.  python3-perf-5.14.0-611.5.1.el9_7.x86_64
ALSA-2025:21926 Moderate/Sec.  python3-perf-5.14.0-611.9.1.el9_7.x86_64
ALSA-2025:19237 Important/Sec. redis-6.2.20-1.el9_6.x86_64
ALSA-2025:20926 Important/Sec. redis-6.2.20-2.el9_7.x86_64
ALSA-2025:20559 Low/Sec.       shadow-utils-2:4.9-15.el9.x86_64
ALSA-2025:20936 Important/Sec. sqlite-libs-3.34.1-9.el9_7.x86_64
ALSA-2025:20954 Important/Sec. sssd-client-2.9.7-4.el9_7.1.x86_64
ALSA-2025:20954 Important/Sec. sssd-common-2.9.7-4.el9_7.1.x86_64
ALSA-2025:20954 Important/Sec. sssd-kcm-2.9.7-4.el9_7.1.x86_64
ALSA-2025:20954 Important/Sec. sssd-nfs-idmap-2.9.7-4.el9_7.1.x86_64
ALSA-2025:22660 Moderate/Sec.  systemd-252-55.el9_7.7.alma.1.x86_64
ALSA-2025:22660 Moderate/Sec.  systemd-libs-252-55.el9_7.7.alma.1.x86_64
ALSA-2025:22660 Moderate/Sec.  systemd-pam-252-55.el9_7.7.alma.1.x86_64
ALSA-2025:22660 Moderate/Sec.  systemd-rpm-macros-252-55.el9_7.7.alma.1.noarch
ALSA-2025:22660 Moderate/Sec.  systemd-udev-252-55.el9_7.7.alma.1.x86_64
ALSA-2025:20945 Moderate/Sec.  vim-common-2:8.2.2637-23.el9_7.x86_64
ALSA-2025:20945 Moderate/Sec.  vim-enhanced-2:8.2.2637-23.el9_7.x86_64
ALSA-2025:20945 Moderate/Sec.  vim-filesystem-2:8.2.2637-23.el9_7.noarch
ALSA-2025:20945 Moderate/Sec.  vim-minimal-2:8.2.2637-23.el9_7.x86_64
ALSA-2025:20922 Important/Sec. webkit2gtk3-jsc-2.50.1-1.el9_7.x86_64
ALSA-2025:22790 Important/Sec. webkit2gtk3-jsc-2.50.3-1.el9_7.x86_64
ALSA-2025:23700 Important/Sec. webkit2gtk3-jsc-2.50.4-1.el9_7.x86_64
ALSA-2025:20961 Moderate/Sec.  xorg-x11-server-Xvfb-1.20.11-32.el9_7.x86_64
ALSA-2025:20961 Moderate/Sec.  xorg-x11-server-common-1.20.11-32.el9_7.x86_64`

func TestDnfManager_getInstalledPackagesWithCancel(t *testing.T) {

	dnf := DnfManager{}

	if !dnf.IsAvailable() {
		t.Skip("dnf is not available, skipping test")
	}

	ctx := t.Context()
	output, err := dnf.getInstalledPackagesWithCancel(ctx)
	if err != nil {
		t.Fatalf("Error getting installed packages: %v", err)
	}

	if !strings.HasPrefix(output, "acl") {
		t.Errorf("Unexpected output, got: %s", output)
	}
}

func TestDnfManager_parseDnfListInstalledOutput(t *testing.T) {
	dnf := DnfManager{}

	packages, err := dnf.parseDnfListInstalledOutput(dnfListInstalledOutput)
	if err != nil {
		t.Fatalf("Error parsing dnf list installed output: %v", err)
	}

	// Table-driven test cases for a few known packages
	tests := []struct {
		name        string
		pkgName     string
		wantVersion string
		wantDesc    string
	}{
		{
			name:        "acl",
			pkgName:     "acl",
			wantVersion: "2.3.1-4.el9",
			wantDesc:    "",
		},
		{
			name:        "bash",
			pkgName:     "bash",
			wantVersion: "5.1.8-9.el9",
			wantDesc:    "",
		},
		{
			name:        "ca-certificates",
			pkgName:     "ca-certificates",
			wantVersion: "2025.2.80_v9.0.305-91.el9",
			wantDesc:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, pkg := range packages {
				if pkg.Name == tt.pkgName {
					found = true
					if pkg.Version != tt.wantVersion {
						t.Errorf("%s: got version %q, want %q", tt.pkgName, pkg.Version, tt.wantVersion)
					}
					if pkg.Description != tt.wantDesc {
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

func TestDnfManager_parseDnfCheckUpdateOutput(t *testing.T) {
	dnf := DnfManager{}

	// Parse installed packages from the updated test data
	installedPkgs, err := dnf.parseDnfListInstalledOutput(dnfListInstalledOutputForUpdateTest)
	if err != nil {
		t.Fatalf("Error parsing installed packages: %v", err)
	}

	var securityPackages map[string]bool

	// Parse available updates (assume dnfAvailableUpdatesOutput is defined elsewhere in the file)
	updates, err := dnf.parseDnfCheckUpdateOutput(dnfAvailableUpdatesOutput, installedPkgs, securityPackages)
	if err != nil {
		t.Fatalf("Error parsing dnf check-update output: %v", err)
	}

	// Table-driven test cases for a few known updates from the new sample
	tests := []struct {
		name             string
		pkgName          string
		wantCurrentVer   string
		wantAvailableVer string
	}{
		{
			name:             "ImageMagick7-djvu update",
			pkgName:          "ImageMagick7-djvu",
			wantCurrentVer:   "1:7.1.2.8-1.el9.remi",
			wantAvailableVer: "1:7.1.2.12-1.el9.remi",
		},
		{
			name:             "avahi-glib update",
			pkgName:          "avahi-glib",
			wantCurrentVer:   "0.8-22.el9_6.1",
			wantAvailableVer: "0.8-23.el9",
		},
		{
			name:             "binutils update",
			pkgName:          "binutils",
			wantCurrentVer:   "2.35.2-63.el9",
			wantAvailableVer: "2.35.2-67.el9_7.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, upd := range updates {
				if upd.Name == tt.pkgName {
					found = true
					if upd.CurrentVersion != tt.wantCurrentVer {
						t.Errorf("%s: got current version %q, want %q", tt.pkgName, upd.CurrentVersion, tt.wantCurrentVer)
					}
					if upd.AvailableVersion != tt.wantAvailableVer {
						t.Errorf("%s: got available version %q, want %q", tt.pkgName, upd.AvailableVersion, tt.wantAvailableVer)
					}
				}
			}
			if !found {
				t.Errorf("Update for package %q not found in parsed output", tt.pkgName)
			}
		})
	}
}
func TestDnfManager_removeArchFromPkgName(t *testing.T) {
	dnf := DnfManager{}

	tests := []struct {
		input    string
		expected string
	}{
		{"bash.x86_64", "bash"},
		{"coreutils-single.x86_64", "coreutils-single"},
		{"basesystem.noarch", "basesystem"},
		{"libstdc++.x86_64", "libstdc++"},
		{"python3-libs.x86_64", "python3-libs"},
		{"mesa-dri-drivers.i686", "mesa-dri-drivers"},
		{"openssl.arm64", "openssl"},
		{"foo.ppc64le", "foo"},
		{"bar.aarch64", "bar"},
		{"package.with.dots.x86_64", "package.with.dots"},
		{"noarchsuffix", "noarchsuffix"},
		{"another.package.i386", "another.package"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := dnf.removeArchFromPkgName(tt.input)
			if got != tt.expected {
				t.Errorf("removeArchFromPkgName(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
func TestDnfManager_removeVersionFromPkgName(t *testing.T) {
	dnf := DnfManager{}

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "grub2-tools-efi-1:2.06-114.el9_7.alma.1.x86_64",
			expected: "grub2-tools-efi",
		},
		{
			input:    "bash-5.1.8-9.el9.x86_64",
			expected: "bash",
		},
		{
			input:    "libstdc++-11.5.0-11.el9.alma.1.x86_64",
			expected: "libstdc++",
		},
		{
			input:    "python3-libs-3.9.25-2.el9_7.x86_64",
			expected: "python3-libs",
		},
		{
			input:    "mesa-dri-drivers-25.0.7-3.el9_7.alma.1.x86_64",
			expected: "mesa-dri-drivers",
		},
		{
			input:    "openssl-1:3.5.1-4.el9_7.arm64",
			expected: "openssl",
		},
		{
			input:    "foo-bar-2.0.1-1.el9.x86_64",
			expected: "foo-bar",
		},
		{
			input:    "nohyphenbeforeversion",
			expected: "nohyphenbeforeversion",
		},
		{
			input:    "package-without-version.x86_64",
			expected: "package-without-version.x86_64",
		},
		{
			input:    "another-package-1.2.3-4.el9.i386",
			expected: "another-package",
		},
		{
			input:    "another-package-1",
			expected: "another-package",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := dnf.removeVersionFromPkgName(tt.input)
			if got != tt.expected {
				t.Errorf("removeVersionFromPkgName(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
func TestDnfManager_parseDnfCheckUpdateOutput_basic(t *testing.T) {
	dnf := DnfManager{}

	installed := []Package{
		{Name: "bash", Version: "5.0.17-2.fc32"},
		{Name: "coreutils", Version: "8.32-4.fc32"},
	}
	securityPkgs := map[string]bool{
		"bash": true,
	}

	output := `
bash.x86_64           5.0.18-1.fc32           @fedora
coreutils.x86_64      8.32-5.fc32             @fedora
nano.x86_64           4.9-1.fc32              @fedora
`
	updates, err := dnf.parseDnfCheckUpdateOutput(output, installed, securityPkgs)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	tests := []struct {
		name          string
		wantCurrent   string
		wantAvailable string
		wantSecurity  bool
	}{
		{"bash", "5.0.17-2.fc32", "5.0.18-1.fc32", true},
		{"coreutils", "8.32-4.fc32", "8.32-5.fc32", false},
		{"nano", "", "4.9-1.fc32", false},
	}

	for _, tt := range tests {
		found := false
		for _, upd := range updates {
			if upd.Name == tt.name {
				found = true
				if upd.CurrentVersion != tt.wantCurrent {
					t.Errorf("%s: got current %q, want %q", tt.name, upd.CurrentVersion, tt.wantCurrent)
				}
				if upd.AvailableVersion != tt.wantAvailable {
					t.Errorf("%s: got available %q, want %q", tt.name, upd.AvailableVersion, tt.wantAvailable)
				}
				if upd.IsSecurityUpdate != tt.wantSecurity {
					t.Errorf("%s: got IsSecurityUpdate %v, want %v", tt.name, upd.IsSecurityUpdate, tt.wantSecurity)
				}
			}
		}
		if !found {
			t.Errorf("Update for %q not found", tt.name)
		}
	}
}

func TestDnfManager_parseDnfCheckUpdateOutput_realWorld(t *testing.T) {
	dnf := DnfManager{}

	installedPkgs, err := dnf.parseDnfListInstalledOutput(dnfListInstalledOutputForUpdateTest)
	if err != nil {
		t.Fatalf("Error parsing installed packages: %v", err)
	}

	securityPkgs, err := dnf.parseDnfSecurityUpdateOutput(availableSecurityUpdatesDnfOutput)
	if err != nil {
		t.Fatalf("Error parsing security updates: %v", err)
	}
	output := dnfAvailableUpdatesOutput
	updates, err := dnf.parseDnfCheckUpdateOutput(output, installedPkgs, securityPkgs)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(updates) == 0 {
		t.Errorf("Expected some updates, got 0")
	}

	if len(updates) != 282 {
		t.Errorf("Expected 282 updates, got %d", len(updates))
	}

	securityUpdateCount := 0
	for _, upd := range updates {
		if upd.IsSecurityUpdate {
			securityUpdateCount++
		}
	}

	if securityUpdateCount == 0 {
		t.Errorf("Expected at least one security update, got 0")
	}

	// Some packages may appear multiple times with different advisories.
	// We use a map to avoid duplicates.
	// This makes testing and determining the expected count a bit harder. Sorry.
	// ALSA-2025:19409 Moderate/Sec.  kernel-5.14.0-570.60.1.el9_6.x86_64
	// ALSA-2025:19930 Moderate/Sec.  kernel-5.14.0-570.62.1.el9_6.x86_64
	// ALSA-2025:22405 Moderate/Sec.  kernel-5.14.0-611.11.1.el9_7.x86_64
	// ALSA-2025:19409 Moderate/Sec.  kernel-core-5.14.0-570.60.1.el9_6.x86_64
	// ALSA-2025:19930 Moderate/Sec.  kernel-core-5.14.0-570.62.1.el9_6.x86_64
	// ALSA-2025:19409 Moderate/Sec.  kernel-headers-5.14.0-570.60.1.el9_6.x86_64
	// ALSA-2025:19930 Moderate/Sec.  kernel-headers-5.14.0-570.62.1.el9_6.x86_64
	if securityUpdateCount != 63 {
		t.Errorf("Expected 63 security updates, got %d", securityUpdateCount)
	}
}

func TestDnfManager_parseDnfCheckUpdateOutput_ignoresHeadersAndShortLines(t *testing.T) {
	dnf := DnfManager{}

	installed := []Package{
		{Name: "foo", Version: "1.0"},
	}
	securityPkgs := map[string]bool{}

	output := `
Last metadata expiration check: 0:00:01 ago on Thu 01 Jan 1970 00:00:01 UTC.
Failed to set locale, defaulting to C.UTF-8

foo.x86_64           1.1           @repo
bar.x86_64           2.0           @repo

shortline
`
	updates, err := dnf.parseDnfCheckUpdateOutput(output, installed, securityPkgs)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(updates) != 2 {
		t.Errorf("Expected 2 updates, got %d", len(updates))
	}
}

func TestDnfManager_parseDnfCheckUpdateOutput_emptyOutput(t *testing.T) {
	dnf := DnfManager{}
	updates, err := dnf.parseDnfCheckUpdateOutput("", nil, nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(updates) != 0 {
		t.Errorf("Expected 0 updates, got %d", len(updates))
	}
}

func TestDnfManager_parseDnfCheckUpdateOutput_securityMapNil(t *testing.T) {
	dnf := DnfManager{}
	installed := []Package{
		{Name: "foo", Version: "1.0"},
	}
	output := "foo.x86_64 1.1 @repo"
	updates, err := dnf.parseDnfCheckUpdateOutput(output, installed, nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(updates) != 1 {
		t.Fatalf("Expected 1 update, got %d", len(updates))
	}
	if updates[0].IsSecurityUpdate {
		t.Errorf("Expected IsSecurityUpdate to be false, got true")
	}
}
func TestDnfManager_parseDnfSecurityUpdateOutput(t *testing.T) {
	dnf := DnfManager{}

	tests := []struct {
		name     string
		input    string
		expected map[string]bool
	}{
		{
			name:  "single ALSA advisory",
			input: `ALSA-2025:20532 Moderate/Sec.  grub2-tools-efi-1:2.06-114.el9_7.alma.1.x86_64`,
			expected: map[string]bool{
				"grub2-tools-efi": true,
			},
		},
		{
			name: "multiple advisories",
			input: `ALSA-2025:20532 Moderate/Sec.  grub2-tools-efi-1:2.06-114.el9_7.alma.1.x86_64
ALSA-2025:20532 Moderate/Sec.  grub2-tools-1:2.06-114.el9_7.alma.1.x86_64
CESA-2025:12345 Important/Sec.  bash-5.1.8-9.el9.x86_64`,
			expected: map[string]bool{
				"grub2-tools-efi": true,
				"grub2-tools":     true,
				"bash":            true,
			},
		},
		{
			name: "ignores headers and short lines",
			input: `Last metadata expiration check: 0:52:58 ago on Wed Jan  7 18:13:57 2026.
Failed to set locale, defaulting to C.UTF-8
ALSA-2025:20532 Moderate/Sec.  grub2-tools-efi-1:2.06-114.el9_7.alma.1.x86_64
shortline`,
			expected: map[string]bool{
				"grub2-tools-efi": true,
			},
		},
		{
			name:     "empty input",
			input:    ``,
			expected: map[string]bool{},
		},
		{
			name:     "line with less than 3 fields",
			input:    `ALSA-2025:20532 Moderate/Sec.`,
			expected: map[string]bool{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dnf.parseDnfSecurityUpdateOutput(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if len(got) != len(tt.expected) {
				t.Errorf("Expected %d entries, got %d", len(tt.expected), len(got))
			}
			for k := range tt.expected {
				if !got[k] {
					t.Errorf("Expected package %q in result", k)
				}
			}
		})
	}
}
