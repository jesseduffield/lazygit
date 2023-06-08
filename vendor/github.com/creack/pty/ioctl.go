//go:build !windows && !solaris && !aix
// +build !windows,!solaris,!aix

package pty

import "syscall"

const (
	TIOCGWINSZ = syscall.TIOCGWINSZ
	TIOCSWINSZ = syscall.TIOCSWINSZ
)

func ioctl(fd, cmd, ptr uintptr) error {
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, cmd, ptr)
	if e != 0 {
		return e
	}
	return nil
}
