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

// Package vt provides common definitions for VT derived terminals and applications.
// This includes the venerable VT100, XTerm, and newer emulators such as Kitty and
// the Windows Terminal.
package vt

import "fmt"

// PrivateMode describes a DEC Private Mode.
type PrivateMode int

const (
	PmAppCursor        PrivateMode = 1    // Application cursor keys.
	PmVT52             PrivateMode = 2    // Clear to enable VT52 compatibility (not supported).
	PmColumns          PrivateMode = 3    // Set to enable 132 columns, reset for 80 columns.
	PmScrolling        PrivateMode = 4    // Smooth scrolling (jump by default).
	PmScreen           PrivateMode = 5    // Set to reverse dark and light on screen.
	PmOrigin           PrivateMode = 6    // Coordinates are relative to margins.
	PmAutoMargin       PrivateMode = 7    // Automatically wrap at margin.
	PmAutoRepeat       PrivateMode = 8    // Enable automatic key repeat.
	PmMouseX10         PrivateMode = 9    // Legacy (X10) mouse reporting.
	PmBlinkCursor      PrivateMode = 12   // Blinking (on) or steady (off) cursor.
	PmPrintFF          PrivateMode = 18   // Print form feed after printing screen.
	PmPrintExtent      PrivateMode = 19   // Print full screen (on) or scrolling region (off).
	PmShowCursor       PrivateMode = 25   // Show the cursor (default on).
	PmCharSet          PrivateMode = 42   // Enable national (on) or multinational (off) character sets.
	PmLeftRightMargin  PrivateMode = 69   // Enable left and right margins
	PmMouseButton      PrivateMode = 1000 // Report mouse button events.
	PmMouseDrag        PrivateMode = 1002 // Report mouse motion events when button depressed, requires PmMouseButton.
	PmMouseMotion      PrivateMode = 1003 // Report mouse motion events, requires PmMouseButton.
	PmFocusReports     PrivateMode = 1004 // Send focus gained or lost reports.
	PmMouseSgr         PrivateMode = 1006 // Use SGR sequences for mouse reports.
	PmMouseSgrPixel    PrivateMode = 1016 // Use SGR sequences for mouse reports, using pixel-level coordinates.
	PmAltScreen        PrivateMode = 1049 // 47 and 1047 are alternates, but we use 1049
	PmBracketedPaste   PrivateMode = 2004 // Bracket pasted text with bracketed paste escape sequences.
	PmSyncOutput       PrivateMode = 2026 // Buffer output when enabled, updating screen when reset.
	PmGraphemeClusters PrivateMode = 2027 // Support for grapheme cluster handling.
	PmResizeReports    PrivateMode = 2048 // Send in-band resize reports.
	PmWin32Input       PrivateMode = 9001 // Use Win32-Input-Mode for keyboard reports
)

// Enable returns the string used to enable this private mode.
func (pm PrivateMode) Enable() string {
	return fmt.Sprintf("\x1b[?%dh", pm)
}

// Disable returns the string used to disable this private mode.
func (pm PrivateMode) Disable() string {
	return fmt.Sprintf("\x1b[?%dl", pm)
}

// Query returns the string used to query the state of this private mode.
func (pm PrivateMode) Query() string {
	return fmt.Sprintf("\x1b[?%d$p", pm)
}

// Reply returns a string representing a query reply for the given mode and status.
func (pm PrivateMode) Reply(status ModeStatus) string {
	return fmt.Sprintf("\x1b[?%d;%d$y", pm, status)
}

// ModeStatus represents the status of the mode.
type ModeStatus int

// AnsiMode are modes standardized in ECMA-48.
// They use CSI-h and CSI-l (no question mark).
type AnsiMode int

// Enable returns the string used to enable this ANSI mode.
func (pm AnsiMode) Enable() string {
	return fmt.Sprintf("\x1b[%dh", pm)
}

// Disable returns the string used to disable this ANSI mode.
func (pm AnsiMode) Disable() string {
	return fmt.Sprintf("\x1b[%dl", pm)
}

// Query returns the string used to query the state of this ANSI mode.
func (pm AnsiMode) Query() string {
	return fmt.Sprintf("\x1b[%d$p", pm)
}

// Reply returns a string representing a query reply for the given mode and status.
func (pm AnsiMode) Reply(status ModeStatus) string {
	return fmt.Sprintf("\x1b[%d;%d$y", pm, status)
}

const (
	AmKeyboardAction AnsiMode = 2  // Lock the keyboard.
	AmInsertReplace  AnsiMode = 4  // Insert or replace characters when a new character is added.
	AmSendReceive    AnsiMode = 12 // XON or XOFF.
	AmNewLineMode    AnsiMode = 20 // If true, LF emits CR as well, and Return sends both CR and LF.
)

const (
	ModeNA        ModeStatus = 0 // Mode is not supported (or unknown)
	ModeOn        ModeStatus = 1 // Mode is on (e.g. via CSI-h)
	ModeOff       ModeStatus = 2 // Mode is off (e.g. via CSI-l)
	ModeOnLocked  ModeStatus = 3 // Mode is hardwired on
	ModeOffLocked ModeStatus = 4 // Mode is hardwired off
)

// Changeable indicates that the mode may be changed.
func (ms ModeStatus) Changeable() bool {
	return ms == ModeOn || ms == ModeOff
}
