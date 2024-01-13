//go:build solaris || illumos
// +build solaris illumos

package termenv

import (
	"golang.org/x/sys/unix"
)

func isForeground(fd int) bool {
	pgrp, err := unix.IoctlGetInt(fd, unix.TIOCGPGRP)
	if err != nil {
		return false
	}

	g, err := unix.Getpgrp()
	if err != nil {
		return false
	}

	return pgrp == g
}
