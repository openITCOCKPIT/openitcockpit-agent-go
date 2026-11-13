package checks

import (
	"unsafe"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// vfsIDCtl is struct vfsidctl of sys/mount.h (64 bit)
type vfsIDCtl struct {
	Version int32
	Fsid    [2]int32
	Ptr     uint64
	Len     uint64
	Spare   [12]uint32
}

// vfsQuery is struct vfsquery of sys/mount.h
type vfsQuery struct {
	Flags uint32
	Spare [31]uint32
}

const (
	vfsCtlVersion1  = 0x01
	vfsCtlQuery     = 0x00010003
	vqNotResponding = 0x0001
	sysctlName2Oid  = 3
	sysctlMaxName   = 24
)

// notRespondingMounts returns the mountpoints of the file systems which do not respond (e.g. NFS or SMB).
// After the first timeout macOS answers statfs of such a file system from the cache, so the disk check
// would report old values. The state is queried with sysctl vfs.generic.ctlbyfsid (VFS_CTL_QUERY),
// the kernel answers it without asking the server. It requires root, nil is returned otherwise.
func notRespondingMounts() map[string]bool {
	mib, err := sysctlNameToMib("vfs.generic.ctlbyfsid")
	if err != nil {
		log.Debugln("DiskCheck: vfs.generic.ctlbyfsid not available: ", err)
		return nil
	}
	mib = append(mib, vfsCtlQuery)

	stats, err := getfsstatNoWait()
	if err != nil {
		log.Debugln("DiskCheck: getfsstat: ", err)
		return nil
	}

	notResponding := map[string]bool{}
	for i := range stats {
		if isNotResponding(mib, &stats[i]) {
			notResponding[unix.ByteSliceToString(stats[i].Mntonname[:])] = true
		}
	}
	return notResponding
}

// isNotResponding returns true if the file system reports VQ_NOTRESP. Errors (e.g. no permission or
// file systems without VFS_CTL_QUERY like devfs) are ignored.
func isNotResponding(mib []int32, stat *unix.Statfs_t) bool {
	request := vfsIDCtl{
		Version: vfsCtlVersion1,
		Fsid:    stat.Fsid.Val,
	}
	var query vfsQuery
	queryLength := unsafe.Sizeof(query)
	err := rawSysctl(mib, unsafe.Pointer(&query), &queryLength, unsafe.Pointer(&request), unsafe.Sizeof(request))
	if err != nil {
		return false
	}
	return query.Flags&vqNotResponding != 0
}

// sysctlNameToMib resolves a sysctl name like sysctlnametomib(3)
func sysctlNameToMib(name string) ([]int32, error) {
	var mib [sysctlMaxName]int32
	length := unsafe.Sizeof(mib)
	nameBytes := []byte(name)
	err := rawSysctl([]int32{0, sysctlName2Oid}, unsafe.Pointer(&mib[0]), &length, unsafe.Pointer(&nameBytes[0]), uintptr(len(nameBytes)))
	if err != nil {
		return nil, err
	}
	return mib[:length/4], nil
}

// rawSysctl calls sysctl(2) with a new value, which is not possible with the functions of x/sys/unix
func rawSysctl(mib []int32, old unsafe.Pointer, oldLength *uintptr, newValue unsafe.Pointer, newLength uintptr) error {
	_, _, errno := unix.Syscall6(
		unix.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mib[0])),
		uintptr(len(mib)),
		uintptr(old),
		uintptr(unsafe.Pointer(oldLength)),
		uintptr(newValue),
		newLength,
	)
	if errno != 0 {
		return errno
	}
	return nil
}
