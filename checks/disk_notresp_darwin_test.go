package checks

import (
	"testing"
	"unsafe"
)

func TestVfsQueryStructSizes(t *testing.T) {
	// sizeof(struct vfsidctl) and sizeof(struct vfsquery) on 64 bit macOS
	if size := unsafe.Sizeof(vfsIDCtl{}); size != 80 {
		t.Errorf("vfsidctl has %d bytes, expected 80", size)
	}
	if size := unsafe.Sizeof(vfsQuery{}); size != 128 {
		t.Errorf("vfsquery has %d bytes, expected 128", size)
	}
}

func TestSysctlNameToMib(t *testing.T) {
	mib, err := sysctlNameToMib("vfs.generic.ctlbyfsid")
	if err != nil {
		t.Fatal(err)
	}
	if len(mib) != 3 || mib[0] != 3 || mib[1] != 0 {
		t.Errorf("unexpected mib %v", mib)
	}
	// Does not panic without root, all file systems respond
	t.Logf("not responding: %v", notRespondingMounts())
}
