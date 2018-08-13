// +build darwin freebsd openbsd netbsd dragonfly
// +build !appengine,!gopherjs

package logrus

import "golang.org/x/sys/unix"

const ioctlReadTermios = unix.TIOCGETA

type Termios unix.Termios
