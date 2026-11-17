//go:build linux || freebsd || aix
// +build linux freebsd aix

package checks

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
)

// CheckMounts reports filesystems of the /etc/fstab (AIX: /etc/filesystems) which are not mounted,
// e.g. a NFS share after a reboot
type CheckMounts struct {
	// fstab and mountpoints can be changed by tests
	fstab       string
	mountpoints func(ctx context.Context) ([]string, error)
}

// Name will be used in the response as check name
func (c *CheckMounts) Name() string {
	return "mounts"
}

type resultFstabEntry struct {
	Device     string `json:"device"`     // e.g.: UUID=..., /dev/sdb1 or nas:/export
	Mountpoint string `json:"mountpoint"` // e.g.: /mnt/backup
	Fstype     string `json:"fstype"`     // e.g.: ext4, nfs
	Options    string `json:"options"`    // e.g.: defaults,nofail
	Nofail     bool   `json:"nofail"`     // The boot does not fail, if the filesystem is missing (nofail, FreeBSD: failok)
	Network    bool   `json:"network"`    // Network filesystem (nfs, cifs, ... or _netdev)
}

type resultMounts struct {
	Fstab    string              `json:"fstab"`
	Expected int                 `json:"expected"` // Number of filesystems of the fstab, which should be mounted
	Missing  []*resultFstabEntry `json:"missing"`
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckMounts) Run(ctx context.Context) (interface{}, error) {
	path := c.fstab
	if path == "" {
		path = fstabPath
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	mountpointsFunc := c.mountpoints
	if mountpointsFunc == nil {
		mountpointsFunc = currentMountpoints
	}
	mounted, err := mountpointsFunc(ctx)
	if err != nil {
		return nil, err
	}
	isMounted := make(map[string]bool, len(mounted))
	for _, mountpoint := range mounted {
		isMounted[filepath.Clean(mountpoint)] = true
	}

	result := &resultMounts{Fstab: path, Missing: []*resultFstabEntry{}}
	for _, entry := range parseMountTable(data) {
		result.Expected++
		if !isMounted[entry.Mountpoint] {
			result.Missing = append(result.Missing, entry)
		}
	}
	return result, nil
}

// Configure the command or return false if the command was disabled
func (c *CheckMounts) Configure(cfg *config.Configuration) (bool, error) {
	if !cfg.Mounts {
		return false, nil
	}
	_, err := os.Stat(fstabPath)
	return err == nil, nil
}

// networkFilesystems are filesystems which require the network
var networkFilesystems = map[string]bool{
	"nfs": true, "nfs4": true, "cifs": true, "smbfs": true, "smb3": true, "sshfs": true, "fuse.sshfs": true,
	"glusterfs": true, "ceph": true, "fuse.ceph": true, "davfs": true, "9p": true,
}

// parseFstab returns the filesystems of a fstab which should be mounted
func parseFstab(data []byte) []*resultFstabEntry {
	var entries []*resultFstabEntry
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		entry := &resultFstabEntry{
			Device:     unescapeMountField(fields[0]),
			Mountpoint: unescapeMountField(fields[1]),
			Fstype:     strings.ToLower(fields[2]),
		}
		if len(fields) > 3 {
			entry.Options = fields[3]
		}
		options := map[string]bool{}
		for _, option := range strings.Split(entry.Options, ",") {
			options[strings.ToLower(option)] = true
		}

		// Swap, entries without mountpoint and filesystems which are not mounted at boot
		if entry.Fstype == "swap" || entry.Mountpoint == "none" || !strings.HasPrefix(entry.Mountpoint, "/") || options["noauto"] {
			continue
		}
		entry.Mountpoint = filepath.Clean(entry.Mountpoint)
		// FreeBSD calls the option failok
		entry.Nofail = options["nofail"] || options["failok"]
		entry.Network = options["_netdev"] || networkFilesystems[entry.Fstype]
		entries = append(entries, entry)
	}
	return entries
}

// unescapeMountField decodes the octal escapes of /proc/mounts and /etc/fstab, e.g. \040 for a space
func unescapeMountField(value string) string {
	if !strings.Contains(value, `\`) {
		return value
	}
	var builder strings.Builder
	for i := 0; i < len(value); i++ {
		if value[i] == '\\' && i+4 <= len(value) {
			if code, err := strconv.ParseUint(value[i+1:i+4], 8, 8); err == nil {
				builder.WriteByte(byte(code))
				i += 3
				continue
			}
		}
		builder.WriteByte(value[i])
	}
	return builder.String()
}
