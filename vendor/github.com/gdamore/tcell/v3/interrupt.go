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

// EventInterrupt is a generic wakeup event.  Its can be used to
// to request a redraw.  It can carry an arbitrary payload, as well.
type EventInterrupt struct {
	EventTime
	v any
}

// Data is used to obtain the opaque event payload.
func (ev *EventInterrupt) Data() any {
	return ev.v
}

// NewEventInterrupt creates an EventInterrupt with the given payload.
func NewEventInterrupt(data any) *EventInterrupt {
	ev := &EventInterrupt{v: data}
	ev.SetEventNow()
	return ev
}
