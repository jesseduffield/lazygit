// +build !windows

package pty

import (
	"os"
	"syscall"
	"unsafe"
)

// InheritSize applies the terminal size of master to slave. This should be run
// in a signal handler for syscall.SIGWINCH to automatically resize the slave when
// the master receives a window size change notification.
func InheritSize(master, slave *os.File) error {
	size, err := GetsizeFull(master)
	if err != nil {
		return err
	}
	err = Setsize(slave, size)
	if err != nil {
		return err
	}
	return nil
}

// Setsize resizes t to s.
func Setsize(t *os.File, ws *Winsize) error {
	return windowRectCall(ws, t.Fd(), syscall.TIOCSWINSZ)
}

// GetsizeFull returns the full terminal size description.
func GetsizeFull(t *os.File) (size *Winsize, err error) {
	var ws Winsize
	err = windowRectCall(&ws, t.Fd(), syscall.TIOCGWINSZ)
	return &ws, err
}

// Getsize returns the number of rows (lines) and cols (positions
// in each line) in terminal t.
func Getsize(t *os.File) (rows, cols int, err error) {
	ws, err := GetsizeFull(t)
	return int(ws.Rows), int(ws.Cols), err
}

// Winsize describes the terminal size.
type Winsize struct {
	Rows uint16 // ws_row: Number of rows (in cells)
	Cols uint16 // ws_col: Number of columns (in cells)
	X    uint16 // ws_xpixel: Width in pixels
	Y    uint16 // ws_ypixel: Height in pixels
}

func windowRectCall(ws *Winsize, fd, a2 uintptr) error {
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		a2,
		uintptr(unsafe.Pointer(ws)),
	)
	if errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}
