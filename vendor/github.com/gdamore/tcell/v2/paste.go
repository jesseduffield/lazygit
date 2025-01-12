// Copyright 2024 The TCell Authors
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
//
// An event with .Start() true will be sent to mark the start of a bracketed paste,
// followed by a number of keys (string data) for the content, ending with the
// an event with .End() true.
type EventPaste struct {
	start bool
	t     time.Time
	data  []byte
}

// When returns the time when this EventPaste was created.
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

// NewEventClipboard returns a new NewEventClipboard with a data payload
func NewEventClipboard(data []byte) *EventClipboard {
	return &EventClipboard{t: time.Now(), data: data}
}

// EventClipboard represents data from the clipboard,
// in response to a GetClipboard request.
type EventClipboard struct {
	t    time.Time
	data []byte
}

// Data returns the attached binary data.
func (ev *EventClipboard) Data() []byte {
	return ev.data
}

// When returns the time when this event was created.
func (ev *EventClipboard) When() time.Time {
	return ev.t
}
