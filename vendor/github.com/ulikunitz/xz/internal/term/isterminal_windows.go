// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package term

import (
	"syscall"
	"unsafe"
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var getConsoleMode = kernel32.NewProc("GetConsoleMode")

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd uintptr) bool {
	var st uint32
	r, _, e := syscall.Syscall(getConsoleMode.Addr(),
		2, fd, uintptr(unsafe.Pointer(&st)), 0)
	return r != 0 && e == 0
}
