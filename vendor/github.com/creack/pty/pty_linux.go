//go:build linux
//+build linux

package pty

import (
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

func open() (pty, tty *os.File, err error) {
	p, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		return nil, nil, err
	}
	// In case of error after this point, make sure we close the ptmx fd.
	defer func() {
		if err != nil {
			_ = p.Close() // Best effort.
		}
	}()

	sname, err := ptsname(p)
	if err != nil {
		return nil, nil, err
	}

	if err := unlockpt(p); err != nil {
		return nil, nil, err
	}

	t, err := os.OpenFile(sname, os.O_RDWR|syscall.O_NOCTTY, 0) //nolint:gosec // Expected Open from a variable.
	if err != nil {
		return nil, nil, err
	}
	return p, t, nil
}

func ptsname(f *os.File) (string, error) {
	var n _C_uint
	err := ioctl(f.Fd(), syscall.TIOCGPTN, uintptr(unsafe.Pointer(&n))) //nolint:gosec // Expected unsafe pointer for Syscall call.
	if err != nil {
		return "", err
	}
	return "/dev/pts/" + strconv.Itoa(int(n)), nil
}

func unlockpt(f *os.File) error {
	var u _C_int
	// use TIOCSPTLCK with a pointer to zero to clear the lock
	return ioctl(f.Fd(), syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&u))) //nolint:gosec // Expected unsafe pointer for Syscall call.
}
