//go:build !linux && !darwin && !freebsd && !dragonfly && !netbsd && !openbsd && !solaris && !windows
// +build !linux,!darwin,!freebsd,!dragonfly,!netbsd,!openbsd,!solaris,!windows

package pty

import (
	"os"
)

func open() (pty, tty *os.File, err error) {
	return nil, nil, ErrUnsupported
}
