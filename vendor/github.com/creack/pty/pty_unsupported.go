//go:build !linux && !darwin && !freebsd && !dragonfly && !netbsd && !openbsd && !solaris
// +build !linux,!darwin,!freebsd,!dragonfly,!netbsd,!openbsd,!solaris

package pty

import (
	"os"
)

func open() (pty, tty *os.File, err error) {
	return nil, nil, ErrUnsupported
}
