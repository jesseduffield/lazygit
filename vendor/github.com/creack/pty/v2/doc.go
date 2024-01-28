// Package pty provides functions for working with Unix terminals.
package pty

import (
	"errors"
	"io"
	"time"
)

// ErrUnsupported is returned if a function is not
// available on the current platform.
var ErrUnsupported = errors.New("unsupported")

// Open a pty and its corresponding tty.
func Open() (Pty, Tty, error) {
	return open()
}

// FdHolder surfaces the Fd() method of the underlying handle.
type FdHolder interface {
	Fd() uintptr
}

// DeadlineHolder surfaces the SetDeadline() method to sets the read and write deadlines.
type DeadlineHolder interface {
	SetDeadline(t time.Time) error
}

// Pty for terminal control in current process.
//
//   - For Unix systems, the real type is *os.File.
//   - For Windows, the real type is a *WindowsPty for ConPTY handle.
type Pty interface {
	// FdHolder is intended to resize / control ioctls of the TTY of the child process in current process.
	FdHolder

	Name() string

	// WriteString is only used to identify Pty and Tty.
	WriteString(s string) (n int, err error)

	io.ReadWriteCloser
}

// Tty for data I/O in child process.
//
//   - For Unix systems, the real type is *os.File.
//   - For Windows, the real type is a *WindowsTty, which is a combination of two pipe file.
type Tty interface {
	// FdHolder Fd only intended for manual InheritSize from Pty.
	FdHolder

	Name() string

	io.ReadWriteCloser
}
