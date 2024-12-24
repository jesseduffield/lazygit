//go:build solaris
//+build solaris

package pty

import (
	"syscall"
	"unsafe"
)

//go:cgo_import_dynamic libc_ioctl ioctl "libc.so"
//go:linkname procioctl libc_ioctl
var procioctl uintptr

const (
	// see /usr/include/sys/stropts.h
	I_PUSH = uintptr((int32('S')<<8 | 002))
	I_STR  = uintptr((int32('S')<<8 | 010))
	I_FIND = uintptr((int32('S')<<8 | 013))

	// see /usr/include/sys/ptms.h
	ISPTM   = (int32('P') << 8) | 1
	UNLKPT  = (int32('P') << 8) | 2
	PTSSTTY = (int32('P') << 8) | 3
	ZONEPT  = (int32('P') << 8) | 4
	OWNERPT = (int32('P') << 8) | 5

	// see /usr/include/sys/termios.h
	TIOCSWINSZ = (uint32('T') << 8) | 103
	TIOCGWINSZ = (uint32('T') << 8) | 104
)

type strioctl struct {
	icCmd     int32
	icTimeout int32
	icLen     int32
	icDP      unsafe.Pointer
}

// Defined in asm_solaris_amd64.s.
func sysvicall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

func ioctl(fd, cmd, ptr uintptr) error {
	if _, _, errno := sysvicall6(uintptr(unsafe.Pointer(&procioctl)), 3, fd, cmd, ptr, 0, 0, 0); errno != 0 {
		return errno
	}
	return nil
}
