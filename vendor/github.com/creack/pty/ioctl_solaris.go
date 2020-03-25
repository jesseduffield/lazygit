package pty

import (
	"golang.org/x/sys/unix"
	"unsafe"
)

const (
	// see /usr/include/sys/stropts.h
	I_PUSH  = uintptr((int32('S')<<8 | 002))
	I_STR   = uintptr((int32('S')<<8 | 010))
	I_FIND  = uintptr((int32('S')<<8 | 013))
	// see /usr/include/sys/ptms.h
	ISPTM   = (int32('P') << 8) | 1
	UNLKPT  = (int32('P') << 8) | 2
	PTSSTTY = (int32('P') << 8) | 3
	ZONEPT  = (int32('P') << 8) | 4
	OWNERPT = (int32('P') << 8) | 5
)

type strioctl struct {
	ic_cmd    int32
	ic_timout int32
	ic_len    int32
	ic_dp     unsafe.Pointer
}

func ioctl(fd, cmd, ptr uintptr) error {
	return unix.IoctlSetInt(int(fd), uint(cmd), int(ptr))
}
