// Copyright 2020 The TCell Authors
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

import (
	"time"
)

// EventPaste is used to mark the start and end of a bracketed paste.
// An event with .Start() true will be sent to mark the start.
// Then a number of keys will be sent to indicate that the content
// is pasted in.  At the end, an event with .Start() false will be sent.
type EventPaste struct {
	start bool
	t     time.Time
}

// When returns the time when this EventMouse was created.
func (ev *EventPaste) When() time.Time {
	return ev.t
}

// Start returns true if this is the start of a paste.
func (ev *EventPaste) Start() bool {
	return ev.start
}

// End returns true if this is the end of a paste.
func (ev *EventPaste) End() bool {
	return !ev.start
}

// NewEventPaste returns a new EventPaste.
func NewEventPaste(start bool) *EventPaste {
	return &EventPaste{t: time.Now(), start: start}
}
