// Copyright 2026 The TCell Authors
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

package vt

// Backend describes the backend of a terminal.
// This can be used to create a real emulator, while allowing the processor
// front end to handle the common details of parsing escape sequences, the state
// machine, and so forth. Backends support a limited set of common functionality,
// including a cursor. They only need to support writing at the cursor.
type Backend interface {

	// GetPrivateMode returns the status of a given private mode.
	GetPrivateMode(PrivateMode) ModeStatus

	// SetPrivateMode sets a private mode to the given status.
	// If either value is invalid, this should simply ignore the operation.
	SetPrivateMode(PrivateMode, ModeStatus) error

	// GetSize returns the size of the terminal in characters.
	// The X and Y are counts, so the bottom right cell should be at coordinate (X-1, Y-1).
	GetSize() Coord

	// Colors returns the number of colors this terminal can support.  For direct color,
	// return 1<<24. The XTerm palette is assumed. Monochrome terminals should return 0.
	Colors() int

	// Put content at the given location. The string in the cell might be a grapheme cluster, and the width can be
	// 0, 1, or 2.  The backend should try to optimize by keeping style and last coordinates if needed.
	// A graphical backend could use this and be completely stateless.
	Put(Coord, Cell)

	// GetPosition returns the cursor position.
	GetPosition() Coord

	// SetPosition sets the cursor position. If the position is out of bounds,
	// it should be clipped to the window size.
	SetPosition(Coord)

	// Reset resets the terminal to default state.
	Reset()

	// RaiseResize is called by the emulation layer when it has completed its own internal resizing.
	// The backend is responsible for sending a signal (if needed) to child processes as part of this
	// function.  (The emulation layer knows nothing of child processes.)
	RaiseResize()

	// Buffering is called by the emulator to indicate that the backend should buffer contents because
	// multiple updates are taking place.  This should be treated in addition to mode 2026, if the backend
	// supports it.  (Mode 2026 should only be supported by the backend if it actually supports true
	// double buffering.)
	Buffering(bool)

	// SetCursor is used to set the current cursor style.  If the backend does not support changing
	// the cursor shape, it should implement at least hidden, steady, and blinking (typically as a block).
	SetCursor(CursorStyle)
}

// Beeper can be implemented by a backend to indicate it can ring the bell or beep.
// This is typically done in response to a 0x07 bell.
type Beeper interface {
	Beep()
}

// Resizer adds notifications when the window size changes.
type Resizer interface {
	// NotifyResize registers a channel to be posted to if the window size changes.
	NotifyResize(chan<- bool)
}

// Titler adds support for setting the window title. (Typically this is OSC 2.)
// Note that for security reasons we only support setting this.
// We don't bother with icon titles, since few terminal emulators support it, and it
// would be hard for us to do this in any portable fashion.
type Titler interface {
	// SetWindowTitle only changes the window title.
	SetWindowTitle(string)
}

// MouseReporting determines what mouse events the backend reports.
type MouseReporting int

const (
	MouseDisabled = MouseReporting(iota) // No mouse reports at all.
	MouseButtons                         // Report button events only.
	MouseDrag                            // Report drag events.
	MouseMotion                          // Report motion events (movement).
)

// Mouser adds support configuring mouse reporting.
// We also assume that a mouse reporter can report focus events.
type Mouser interface {
	SetMouse(MouseReporting)
}

// Blitter implements a cell-level blit, where a rectangular range of cells is copied from one
// location to another.  The source and destination may overlap.  The old locations will remain
// unchanged except of course or cells overwritten by the blit. The content will also be clipped
// to the visible dimensions.
type Blitter interface {
	Blit(src, dst, dim Coord)
}

// Clipboard implements a clipboard or copy buffer for copy/paste activity.
// The backend may prevent sending clipboard data by returning an empty string
// for the clipboard.  Frequently this is done for security reasons.
type Clipboard interface {

	// SetClipboard sets the contents of the clipboard.
	SetClipboard([]byte)

	// GetClipboard gets the contents of the clipboard.
	// It will return nil if the operation is not supported.
	// An empty clipboard will be []byte{}
	GetClipboard() []byte
}

// AdvancedKeyboard provides raw keyboard events, which gives
// access to key presses, and releases, physical keys, and so forth.
// These keyboards should provide a mapping facility to obtain associated
// Unicode text via a layout, as well.
type AdvancedKeyboard interface {
	// IsAdvancedKeyboard returns true if the emulator supports full keyboard
	// reporting.  This must include key press and release events, and mapping
	// of physical keys (thus permitting key disambiguation).
	IsAdvancedKeyboard() bool
}
