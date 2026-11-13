//go:build aix && ppc64
// +build aix,ppc64

package checks

import (
	"fmt"
	"syscall"
	"unsafe"
)

// libperfstat is called without cgo, like golang.org/x/sys/unix calls the libc on AIX

//go:cgo_import_dynamic libperfstat_perfstat_cpu_total perfstat_cpu_total "libperfstat.a/shr_64.o"
//go:cgo_import_dynamic libperfstat_perfstat_cpu perfstat_cpu "libperfstat.a/shr_64.o"
//go:cgo_import_dynamic libperfstat_perfstat_memory_total perfstat_memory_total "libperfstat.a/shr_64.o"
//go:cgo_import_dynamic libperfstat_perfstat_partition_total perfstat_partition_total "libperfstat.a/shr_64.o"

//go:linkname libperfstat_perfstat_cpu_total libperfstat_perfstat_cpu_total
//go:linkname libperfstat_perfstat_cpu libperfstat_perfstat_cpu
//go:linkname libperfstat_perfstat_memory_total libperfstat_perfstat_memory_total
//go:linkname libperfstat_perfstat_partition_total libperfstat_perfstat_partition_total

type perfstatFunc uintptr

var (
	libperfstat_perfstat_cpu_total       perfstatFunc
	libperfstat_perfstat_cpu             perfstatFunc
	libperfstat_perfstat_memory_total    perfstatFunc
	libperfstat_perfstat_partition_total perfstatFunc
)

// Implemented in perfstat_aix_ppc64.s
func perfstatSyscall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

// perfstatCall calls a perfstat function with the signature
// int perfstat_xxx(perfstat_id_t *name, perfstat_xxx_t *userbuff, int sizeof_userbuff, int desired_number)
func perfstatCall(fn *perfstatFunc, name []byte, buf []byte, size, count int) (int, error) {
	var namePointer, bufPointer uintptr
	if len(name) > 0 {
		namePointer = uintptr(unsafe.Pointer(&name[0]))
	}
	if len(buf) > 0 {
		bufPointer = uintptr(unsafe.Pointer(&buf[0]))
	}
	r1, _, errno := perfstatSyscall6(uintptr(unsafe.Pointer(fn)), 4, namePointer, bufPointer, uintptr(size), uintptr(count), 0, 0)
	if int32(r1) == -1 {
		return -1, errno
	}
	return int(int32(r1)), nil
}

// perfstatTotal reads a perfstat_xxx_total_t structure. The sizes are tried in the order of the slice,
// as perfstat rejects sizes of structure versions it does not know.
func perfstatTotal(name string, fn *perfstatFunc, sizes ...int) ([]byte, error) {
	var lastErr error
	for _, size := range sizes {
		buf := make([]byte, size)
		_, err := perfstatCall(fn, nil, buf, size, 1)
		if err == nil {
			return buf, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("%s: %w", name, lastErr)
}

func perfstatCpuTotal() ([]byte, error) {
	return perfstatTotal("perfstat_cpu_total", &libperfstat_perfstat_cpu_total, perfstatCpuTotalSize, perfstatCpuTotalSize71)
}

func perfstatPartitionTotal() ([]byte, error) {
	return perfstatTotal("perfstat_partition_total", &libperfstat_perfstat_partition_total, perfstatPartitionSize, perfstatPartitionSize71)
}

func perfstatMemoryTotal() ([]byte, error) {
	return perfstatTotal("perfstat_memory_total", &libperfstat_perfstat_memory_total, perfstatMemoryTotalSize)
}

// perfstatCpus reads perfstat_cpu_t of all logical CPUs
func perfstatCpus() ([]perfstatTimes, error) {
	count, err := perfstatCall(&libperfstat_perfstat_cpu, nil, nil, perfstatCpuSize, 0)
	if err != nil {
		return nil, fmt.Errorf("perfstat_cpu: %w", err)
	}
	if count <= 0 || count > perfstatMaxLogicalCpus {
		return nil, fmt.Errorf("perfstat_cpu: unexpected number of CPUs %d", count)
	}

	buf := make([]byte, count*perfstatCpuSize)
	// An empty name (FIRST_CPU) starts with the first CPU
	first := make([]byte, perfstatIdentifierLength)
	read, err := perfstatCall(&libperfstat_perfstat_cpu, first, buf, perfstatCpuSize, count)
	if err != nil {
		return nil, fmt.Errorf("perfstat_cpu: %w", err)
	}
	if read > count {
		read = count
	}
	return decodeCpus(buf, read)
}
