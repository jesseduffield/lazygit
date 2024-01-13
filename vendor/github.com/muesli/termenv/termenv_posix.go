//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd
// +build darwin dragonfly freebsd linux netbsd openbsd

package termenv

import (
	"golang.org/x/sys/unix"
)

func isForeground(fd int) bool {
	pgrp, err := unix.IoctlGetInt(fd, unix.TIOCGPGRP)
	if err != nil {
		return false
	}

	return pgrp == unix.Getpgrp()
}
