// Copyright 2015 The TCell Authors
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

// EventResize is sent when the window size changes.
type EventResize struct {
	t  time.Time
	ws WindowSize
}

// NewEventResize creates an EventResize with the new updated window size,
// which is given in character cells.
func NewEventResize(width, height int) *EventResize {
	ws := WindowSize{
		Width:  width,
		Height: height,
	}
	return &EventResize{t: time.Now(), ws: ws}
}

// When returns the time when the Event was created.
func (ev *EventResize) When() time.Time {
	return ev.t
}

// Size returns the new window size as width, height in character cells.
func (ev *EventResize) Size() (int, int) {
	return ev.ws.Width, ev.ws.Height
}

// PixelSize returns the new window size as width, height in pixels. The size
// will be 0,0 if the screen doesn't support this feature
func (ev *EventResize) PixelSize() (int, int) {
	return ev.ws.PixelWidth, ev.ws.PixelHeight
}

type WindowSize struct {
	Width       int
	Height      int
	PixelWidth  int
	PixelHeight int
}

// CellDimensions returns the dimensions of a single cell, in pixels
func (ws WindowSize) CellDimensions() (int, int) {
	if ws.PixelWidth == 0 || ws.PixelHeight == 0 {
		return 0, 0
	}
	return (ws.PixelWidth / ws.Width), (ws.PixelHeight / ws.Height)
}
