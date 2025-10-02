//go:build plan9
// +build plan9

// Copyright 2025 The TCell Authors
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

import "os"

// initialize on Plan 9: if no TTY was provided, use the Plan 9 TTY.
func (t *tScreen) initialize() error {
    if os.Getenv("TERM") == "" {
        // TERM should be "vt100" in a vt(1) window; color/mouse support will be limited.
        _ = os.Setenv("TERM", "vt100")
    }
	if t.tty == nil {
		tty, err := NewDevTty()
		if err != nil {
			return err
		}
		t.tty = tty
	}
	return nil
}
