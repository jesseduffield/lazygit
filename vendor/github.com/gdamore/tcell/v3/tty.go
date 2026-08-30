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

import "github.com/gdamore/tcell/v3/tty"

type Tty = tty.Tty

// NewDevTty obtains a default tty from the console or TTY (e.g. /dev/tty) for the process.
func NewDevTty() (Tty, error) {
	return tty.NewDevTty()
}

// NewDevTtyFromDev obtains a tty from the given device path. Not supported on Windows.
func NewDevTtyFromDev(dev string) (Tty, error) {
	return tty.NewDevTtyFromDev(dev)
}

// NewStdIoTty obtains a tty from stdin and stdout.
func NewStdIoTty() (Tty, error) {
	return tty.NewStdIoTty()
}
