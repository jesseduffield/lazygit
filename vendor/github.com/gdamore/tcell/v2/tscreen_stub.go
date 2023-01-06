//go:build js || plan9 || windows
// +build js plan9 windows

// Copyright 2022 The TCell Authors
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

// NB: We might someday wish to move Windows to this model.   However,
// that would probably mean sacrificing some of the richer key reporting
// that we can obtain with the console API present on Windows.

func (t *tScreen) initialize() error {
	if t.tty == nil {
		return ErrNoScreen
	}
	// If a tty was supplied (custom), it should work.
	// Custom screen implementations will need to provide a TTY
	// implementation that we can use.
	return nil
}
