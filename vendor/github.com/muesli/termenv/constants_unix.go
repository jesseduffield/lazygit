//go:build (darwin || dragonfly || freebsd || netbsd || openbsd) && !solaris && !illumos
// +build darwin dragonfly freebsd netbsd openbsd
// +build !solaris
// +build !illumos

package termenv

import "golang.org/x/sys/unix"

const (
	tcgetattr = unix.TIOCGETA
	tcsetattr = unix.TIOCSETA
)
