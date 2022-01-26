//go:build dragonfly
//+build dragonfly

package pty

import (
	"errors"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

// same code as pty_darwin.go
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

	if err := unlockpt(p); err != nil {
		return nil, nil, err
	}

	t, err := os.OpenFile(sname, os.O_RDWR, 0)
	if err != nil {
		return nil, nil, err
	}
	return p, t, nil
}

func grantpt(f *os.File) error {
	_, err := isptmaster(f.Fd())
	return err
}

func unlockpt(f *os.File) error {
	_, err := isptmaster(f.Fd())
	return err
}

func isptmaster(fd uintptr) (bool, error) {
	err := ioctl(fd, syscall.TIOCISPTMASTER, 0)
	return err == nil, err
}

var (
	emptyFiodgnameArg fiodgnameArg
	ioctl_FIODNAME    = _IOW('f', 120, unsafe.Sizeof(emptyFiodgnameArg))
)

func ptsname(f *os.File) (string, error) {
	name := make([]byte, _C_SPECNAMELEN)
	fa := fiodgnameArg{Name: (*byte)(unsafe.Pointer(&name[0])), Len: _C_SPECNAMELEN, Pad_cgo_0: [4]byte{0, 0, 0, 0}}

	err := ioctl(f.Fd(), ioctl_FIODNAME, uintptr(unsafe.Pointer(&fa)))
	if err != nil {
		return "", err
	}

	for i, c := range name {
		if c == 0 {
			s := "/dev/" + string(name[:i])
			return strings.Replace(s, "ptm", "pts", -1), nil
		}
	}
	return "", errors.New("TIOCPTYGNAME string not NUL-terminated")
}
