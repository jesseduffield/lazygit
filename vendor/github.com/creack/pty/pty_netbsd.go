//go:build netbsd
//+build netbsd

package pty

import (
	"errors"
	"os"
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

	if err := grantpt(p); err != nil {
		return nil, nil, err
	}

	// In NetBSD unlockpt() does nothing, so it isn't called here.

	t, err := os.OpenFile(sname, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}
	return p, t, nil
}

func ptsname(f *os.File) (string, error) {
	/*
	 * from ptsname(3): The ptsname() function is equivalent to:
	 * struct ptmget pm;
	 * ioctl(fd, TIOCPTSNAME, &pm) == -1 ? NULL : pm.sn;
	 */
	var ptm ptmget
	if err := ioctl(f.Fd(), uintptr(ioctl_TIOCPTSNAME), uintptr(unsafe.Pointer(&ptm))); err != nil {
		return "", err
	}
	name := make([]byte, len(ptm.Sn))
	for i, c := range ptm.Sn {
		name[i] = byte(c)
		if c == 0 {
			return string(name[:i]), nil
		}
	}
	return "", errors.New("TIOCPTSNAME string not NUL-terminated")
}

func grantpt(f *os.File) error {
	/*
	 * from grantpt(3): Calling grantpt() is equivalent to:
	 * ioctl(fd, TIOCGRANTPT, 0);
	 */
	return ioctl(f.Fd(), uintptr(ioctl_TIOCGRANTPT), 0)
}
