//go:build aix && ppc64

#include "textflag.h"

// The functions of libperfstat are called like the libc functions in golang.org/x/sys/unix (asm_aix_ppc64.s),
// the call itself is implemented in runtime/syscall_aix.go
TEXT ·perfstatSyscall6(SB),NOSPLIT,$0-88
	JMP	syscall·syscall6(SB)
