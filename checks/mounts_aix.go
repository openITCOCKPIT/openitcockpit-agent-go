package checks

import (
	"context"
)

// fstabPath is the AIX file with the filesystems and their mount options (stanza format)
const fstabPath = "/etc/filesystems"

// parseMountTable returns the filesystems of /etc/filesystems which are mounted at boot
func parseMountTable(data []byte) []*resultFstabEntry {
	return aixFstabEntries(parseAixFilesystems(data))
}

// currentMountpoints returns the mountpoints of the mount command, it only reads the mount table of the kernel
func currentMountpoints(ctx context.Context) ([]string, error) {
	output, err := runAixCommand(ctx, "mount")
	if err != nil {
		return nil, err
	}
	mounts := parseAixMount(output)
	result := make([]string, 0, len(mounts))
	for _, mount := range mounts {
		result = append(result, mount.Mountpoint)
	}
	return result, nil
}
