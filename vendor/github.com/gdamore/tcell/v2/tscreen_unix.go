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

// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris zos

package tcell

// initialize is used at application startup, and sets up the initial values
// including file descriptors used for terminals and saving the initial state
// so that it can be restored when the application terminates.
func (t *tScreen) initialize() error {
	var err error
	if t.tty == nil {
		t.tty, err = NewDevTty()
		if err != nil {
			return err
		}
	}
	return nil
}
