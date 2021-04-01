// Copyright 2021 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build linux aix zos solaris

package tcell

import (
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

// NB: We might someday wish to move Windows to this model.   However,
// that would probably mean sacrificing some of the richer key reporting
// that we can obtain with the console API present on Windows.

// nonBlocking changes VMIN to 0, and VTIME to 1.  This basically ensures that
// we can wake up the input loop.  We only want to do this if we are going to interrupt
// that loop.  Normally we use VMIN 1 and VTIME 0, which ensures we pick up bytes when
// they come but don't spin burning cycles.
func (t *tScreen) nonBlocking(on bool) {
	fd := int(os.Stdin.Fd())
	tio, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return
	}
	if on {
		tio.Cc[unix.VMIN] = 0
		tio.Cc[unix.VTIME] = 0
	} else {
		// block for any output
		tio.Cc[unix.VTIME] = 0
		tio.Cc[unix.VMIN] = 1
	}

	_ = syscall.SetNonblock(fd, on)
	// We want to set this *right now*.
	_ = unix.IoctlSetTermios(fd, unix.TCSETS, tio)
}
