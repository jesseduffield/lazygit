// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux,!appengine netbsd openbsd

// Package term provides the IsTerminal function.
package term

import (
	"syscall"
	"unsafe"
)

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd uintptr) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd,
		ioctlGetTermios, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}
