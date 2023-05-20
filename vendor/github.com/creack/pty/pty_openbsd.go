//go:build openbsd
// +build openbsd

package pty

import (
	"os"
	"syscall"
	"unsafe"
)

func cInt8ToString(in []int8) string {
	var s []byte
	for _, v := range in {
		if v == 0 {
			break
		}
		s = append(s, byte(v))
	}
	return string(s)
}

func open() (pty, tty *os.File, err error) {
	/*
	 * from ptm(4):
	 * The PTMGET command allocates a free pseudo terminal, changes its
	 * ownership to the caller, revokes the access privileges for all previous
	 * users, opens the file descriptors for the pty and tty devices and
	 * returns them to the caller in struct ptmget.
	 */

	p, err := os.OpenFile("/dev/ptm", os.O_RDWR|syscall.O_CLOEXEC, 0)
	if err != nil {
		return nil, nil, err
	}
	defer p.Close()

	var ptm ptmget
	if err := ioctl(p.Fd(), uintptr(ioctl_PTMGET), uintptr(unsafe.Pointer(&ptm))); err != nil {
		return nil, nil, err
	}

	pty = os.NewFile(uintptr(ptm.Cfd), cInt8ToString(ptm.Cn[:]))
	tty = os.NewFile(uintptr(ptm.Sfd), cInt8ToString(ptm.Sn[:]))

	return pty, tty, nil
}
