// Copyright 2025 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vt

// KeyEvent is a key event.
type KeyEvent struct {
	Down   bool     // true if event is for key down event
	Repeat int      // if > 1, a repeat count
	Key    Key      // Key symbol.
	Base   BaseKey  // base key code (physical key, e.g 'a'), may be zero if same as code
	VK     WinVK    // Windows virtual key. 0 for none, or if not known.
	SC     ScanCode // Windows scan code. 0 for none, or if not known.
	Mod    Modifier // modifiers
	Utf    string   // if non-empty, the unicode content for this
}

type Modifier int

const (
	ModNone   = Modifier(0)
	ModLShift = Modifier(1 << iota)
	ModRShift
	ModLCtrl
	ModRCtrl
	ModLAlt
	ModRAlt
	ModLMeta
	ModRMeta
	ModLHyper
	ModRHyper
	ModCapsLock
	ModNumLock
)

func (m Modifier) IsShift() bool    { return (m & (ModLShift | ModRShift)) != 0 }
func (m Modifier) IsCtrl() bool     { return (m & (ModLCtrl | ModRCtrl)) != 0 }
func (m Modifier) IsAlt() bool      { return (m & (ModLAlt | ModRAlt)) != 0 }
func (m Modifier) IsMeta() bool     { return (m & (ModLMeta | ModRMeta)) != 0 }
func (m Modifier) IsHyper() bool    { return (m & (ModLHyper | ModRHyper)) != 0 }
func (m Modifier) IsNumLock() bool  { return m&ModNumLock != 0 }
func (m Modifier) IsCapsLock() bool { return m&ModCapsLock != 0 }
func (m Modifier) IsCapitals() bool { return m.IsCapsLock() != m.IsShift() }
func (m Modifier) IsAltGr() bool    { return m.IsCtrl() && m.IsAlt() }

// Button is the mouse button pressed or released.
type Button int

const (
	NoButton = Button(0)         // No buttons are pressed.
	Button1  = Button(1 << iota) // Usually left most button.
	Button2                      // Usually right most button.
	Button3                      // Usually middle button.
	Button4
	Button5
	Button6
	Button7
	Button8
	WheelUp    // Wheel motion up/away from user.
	WheelDown  // Wheel motion down/towards user.
	WheelLeft  // Wheel motion to left.
	WheelRight // Wheel motion to right.
)

// MouseEvent reports a single mouse event.  Only a single button
// may be reported for a given event.  The application will have
// to keep state.  As buttons are never pressed exactly simultaneously,
// the backend will send chords as a series of presses followed by a series
// of releases.
type MouseEvent struct {
	Position Coord    // Location of pointer.
	Button   Button   // Buttons pressed.
	Down     bool     // True on press, false on release.
	Motion   bool     // True if mouse moved at least once cell.
	Mod      Modifier // Modifiers (for modified click).
}

// encodeButton just encodes the XTerm style button details into a byte
func (ev MouseEvent) encodeButton() byte {
	var btn byte
	switch ev.Button {
	case NoButton:
		btn = 3
	case Button1:
		btn = 0
	case Button2: // intentionally reversed with button 3
		btn = 2
	case Button3:
		btn = 1
	case WheelUp:
		btn = 0x40
	case WheelDown:
		btn = 0x41
	case WheelLeft:
		btn = 0x42
	case WheelRight:
		btn = 0x43
	case Button4:
		btn = 0x80
	case Button5:
		btn = 0x81
	case Button6:
		btn = 0x82
	case Button7:
		btn = 0x83
	default:
		btn = 3
	}
	if ev.Motion {
		btn += 0x20
	}
	if ev.Mod.IsShift() {
		btn += 4
	}
	if ev.Mod.IsAlt() || ev.Mod.IsMeta() {
		btn += 8
	}
	if ev.Mod.IsCtrl() {
		btn += 16
	}
	return btn
}
