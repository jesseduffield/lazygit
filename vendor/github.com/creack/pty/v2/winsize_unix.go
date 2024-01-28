//go:build !windows
// +build !windows

package pty

import (
	"os"
	"syscall"
	"unsafe"
)

// Setsize resizes t to s.
func Setsize(t FdHolder, ws *Winsize) error {
	//nolint:gosec // Expected unsafe pointer for Syscall call.
	return ioctl(t.(*os.File), syscall.TIOCSWINSZ, uintptr(unsafe.Pointer(ws)))
}

// GetsizeFull returns the full terminal size description.
func GetsizeFull(t FdHolder) (size *Winsize, err error) {
	var ws Winsize

	//nolint:gosec // Expected unsafe pointer for Syscall call.
	if err := ioctl(t.(*os.File), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&ws))); err != nil {
		return nil, err
	}
	return &ws, nil
}
