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

package tcell

import "io"

// Tty is an abstraction of a tty (traditionally "teletype").  This allows applications to
// provide for alternate backends, as there are situations where the traditional /dev/tty
// does not work, or where more flexible handling is required.  This interface is for use
// with the terminfo-style based API.  It extends the io.ReadWriter API.  It is reasonable
// that the implementation might choose to use different underlying files for the Reader
// and Writer sides of this API, as part of it's internal implementation.
type Tty interface {
	// Start is used to activate the Tty for use.  Upon return the terminal should be
	// in raw mode, non-blocking, etc.  The implementation should take care of saving
	// any state that is required so that it may be restored when Stop is called.
	Start() error

	// Stop is used to stop using this Tty instance.  This may be a suspend, so that other
	// terminal based applications can run in the foreground.  Implementations should
	// restore any state collected at Start(), and return to ordinary blocking mode, etc.
	// Drain is called first to drain the input.  Once this is called, no more Read
	// or Write calls will be made until Start is called again.
	Stop() error

	// Drain is called before Stop, and ensures that the reader will wake up appropriately
	// if it was blocked.  This workaround is required for /dev/tty on certain UNIX systems
	// to ensure that Read() does not block forever.  This typically arranges for the tty driver
	// to send data immediately (e.g. VMIN and VTIME both set zero) and sets a deadline on input.
	// Implementations may reasonably make this a no-op.  There will still be control sequences
	// emitted between the time this is called, and when Stop is called.
	Drain() error

	// NotifyResize is used register a callback when the tty thinks the dimensions have
	// changed.  The standard UNIX implementation links this to a handler for SIGWINCH.
	// If the supplied callback is nil, then any handler should be unregistered.
	NotifyResize(cb func())

	// WindowSize is called to determine the terminal dimensions.  This might be determined
	// by an ioctl or other means.
	WindowSize() (WindowSize, error)

	io.ReadWriteCloser
}
