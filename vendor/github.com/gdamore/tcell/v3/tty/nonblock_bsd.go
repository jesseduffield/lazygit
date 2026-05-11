// Copyright 2021 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build darwin || dragonfly || freebsd || netbsd || openbsd
// +build darwin dragonfly freebsd netbsd openbsd

package tty

import (
	"syscall"

	"golang.org/x/sys/unix"
)

// BSD systems use TIOC style ioctls.

// tcSetBufParams is used by the tty driver on UNIX systems to configure the
// buffering parameters (minimum character count and minimum wait time in msec.)
// This also waits for output to drain first.
func tcSetBufParams(fd int, vMin uint8, vTime uint8) error {
	_ = syscall.SetNonblock(fd, true)
	tio, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)
	if err != nil {
		return err
	}
	tio.Cc[unix.VMIN] = vMin
	tio.Cc[unix.VTIME] = vTime
	if err = unix.IoctlSetTermios(fd, unix.TIOCSETAW, tio); err != nil {
		return err
	}
	return nil
}

// tcFlushInput discards any queued input before the caller starts reading from
// the tty. This avoids stale bytes, such as delayed mouse reports, from being
// delivered to the next foreground application.
func tcFlushInput(fd int) error {
	return unix.IoctlSetPointerInt(fd, unix.TIOCFLUSH, unix.TCIFLUSH)
}
