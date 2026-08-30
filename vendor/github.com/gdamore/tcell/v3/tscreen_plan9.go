//go:build plan9
// +build plan9

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

// initialize on Plan 9: if no TTY was provided, use the Plan 9 TTY.
func (t *tScreen) initialize() error {
	if t.tty == nil {
		tty, err := NewDevTty()
		if err != nil {
			return err
		}
		t.tty = tty
	}
	return nil
}
