// Copyright 2025 The TCell Authors
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

package tcell

import (
	"errors"
)

var (
	// ErrNoScreen indicates that no suitable screen could be found.
	// This may result from attempting to run on a platform where there
	// is no support for either termios or console I/O (such as js),
	// or from running in an environment where there is no access to
	// a suitable console/terminal device.  (For example, running on
	// without a controlling TTY or with no /dev/tty on POSIX platforms.)
	ErrNoScreen = errors.New("no suitable screen available")

	// ErrNoCharset indicates that the locale environment the
	// program is not supported by the program, because no suitable
	// encoding was found for it.  This problem never occurs if
	// the environment is UTF-8 or UTF-16.
	ErrNoCharset = errors.New("character set not supported")

	// ErrEventQFull indicates that the event queue is full, and
	// cannot accept more events.
	ErrEventQFull = errors.New("event queue full")
)

// An EventError is an event representing some sort of error, and carries
// an error payload.
type EventError struct {
	EventTime
	err error
}

// Error implements the error.
func (ev *EventError) Error() string {
	return ev.err.Error()
}

// Unwrap exposes the underlying error payload so callers can use
// errors.Is / errors.As to match against sentinel values such as
// io.EOF.
func (ev *EventError) Unwrap() error {
	return ev.err
}

// NewEventError creates an ErrorEvent with the given error payload.
func NewEventError(err error) *EventError {
	ev := &EventError{err: err}
	ev.SetEventNow()
	return ev
}
