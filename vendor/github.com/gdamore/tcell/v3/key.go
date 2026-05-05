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

import (
	"fmt"
	"strings"
)

// EventKey represents a key press.  Usually this is a key press followed
// by a key release, but since terminal programs don't have a way to report
// key release events, we usually get just one event.  If a key is held down
// then the terminal may synthesize repeated key presses at some predefined
// rate.  We have no control over that, nor visibility into it.
//
// In some cases, we can have a modifier key, such as ModAlt, that can be
// generated with a key press.  (This usually is represented by having the
// high bit set, or in some cases, by sending an ESC prior to the rune.)
//
// If the value of Key() is KeyRune, then the actual key value will be
// available (as a grapheme cluster) with the Str() method.
// This will be the case for most keys.
//
// In most situations, the modifiers will not be set.  For example, if the
// rune is 'A', this will be reported without the ModShift bit set, since
// really can't tell if the Shift key was pressed (it might have been CAPSLOCK,
// or a terminal that only can send capitals, or keyboard with separate
// capital letters from lower case letters).
//
// Generally, terminal applications have far less visibility into keyboard
// activity than graphical applications.  Hence, they should avoid depending
// overly much on availability of modifiers, or the availability of any
// specific keys.
type EventKey struct {
	EventTime
	mod ModMask
	key Key
	str string // string for key, usually just one character, but may be composed sequence
}

// Str returns the string corresponding to the key press, if it makes sense.
// The result is only defined if the value of Key() is KeyRune.  It will be
// either one key (e.g. 'A'), or could be a composed sequence.
func (ev *EventKey) Str() string {
	return ev.str
}

// Key returns a virtual key code.  We use this to identify specific key
// codes, such as KeyEnter, etc.  Most control and function keys are reported
// with unique Key values.  Normal alphanumeric and punctuation keys will
// generally return KeyRune here; the specific key can be further decoded
// using the Str() function.
func (ev *EventKey) Key() Key {
	return ev.key
}

// Modifiers returns the modifiers that were present with the key press.  Note
// that not all platforms and terminals support this equally well, and some
// cases we will not not know for sure.  Hence, applications should avoid
// using this in most circumstances.
func (ev *EventKey) Modifiers() ModMask {
	return ev.mod
}

// KeyNames holds the written names of special keys. Useful to echo back a key
// name, or to look up a key from a string value.
var KeyNames = map[Key]string{
	KeyEnter:      "Enter",
	KeyBackspace:  "Backspace",
	KeyTab:        "Tab",
	KeyBacktab:    "Backtab",
	KeyEsc:        "Esc",
	KeyBackspace2: "Backspace2",
	KeyDelete:     "Delete",
	KeyInsert:     "Insert",
	KeyUp:         "Up",
	KeyDown:       "Down",
	KeyLeft:       "Left",
	KeyRight:      "Right",
	KeyHome:       "Home",
	KeyEnd:        "End",
	KeyUpLeft:     "UpLeft",
	KeyUpRight:    "UpRight",
	KeyDownLeft:   "DownLeft",
	KeyDownRight:  "DownRight",
	KeyCenter:     "Center",
	KeyPgDn:       "PgDn",
	KeyPgUp:       "PgUp",
	KeyClear:      "Clear",
	KeyExit:       "Exit",
	KeyCancel:     "Cancel",
	KeyPause:      "Pause",
	KeyPrint:      "Print",
	KeyF1:         "F1",
	KeyF2:         "F2",
	KeyF3:         "F3",
	KeyF4:         "F4",
	KeyF5:         "F5",
	KeyF6:         "F6",
	KeyF7:         "F7",
	KeyF8:         "F8",
	KeyF9:         "F9",
	KeyF10:        "F10",
	KeyF11:        "F11",
	KeyF12:        "F12",
	KeyF13:        "F13",
	KeyF14:        "F14",
	KeyF15:        "F15",
	KeyF16:        "F16",
	KeyF17:        "F17",
	KeyF18:        "F18",
	KeyF19:        "F19",
	KeyF20:        "F20",
	KeyF21:        "F21",
	KeyF22:        "F22",
	KeyF23:        "F23",
	KeyF24:        "F24",
	KeyF25:        "F25",
	KeyF26:        "F26",
	KeyF27:        "F27",
	KeyF28:        "F28",
	KeyF29:        "F29",
	KeyF30:        "F30",
	KeyF31:        "F31",
	KeyF32:        "F32",
	KeyF33:        "F33",
	KeyF34:        "F34",
	KeyF35:        "F35",
	KeyF36:        "F36",
	KeyF37:        "F37",
	KeyF38:        "F38",
	KeyF39:        "F39",
	KeyF40:        "F40",
	KeyF41:        "F41",
	KeyF42:        "F42",
	KeyF43:        "F43",
	KeyF44:        "F44",
	KeyF45:        "F45",
	KeyF46:        "F46",
	KeyF47:        "F47",
	KeyF48:        "F48",
	KeyF49:        "F49",
	KeyF50:        "F50",
	KeyF51:        "F51",
	KeyF52:        "F52",
	KeyF53:        "F53",
	KeyF54:        "F54",
	KeyF55:        "F55",
	KeyF56:        "F56",
	KeyF57:        "F57",
	KeyF58:        "F58",
	KeyF59:        "F59",
	KeyF60:        "F60",
	KeyF61:        "F61",
	KeyF62:        "F62",
	KeyF63:        "F63",
	KeyF64:        "F64",
	KeyMenu:       "Menu",
	KeyCapsLock:   "CapsLock",
	KeyScrollLock: "ScrollLock",
	KeyNumLock:    "NumLock",
	KeyCtrlA:      "Ctrl-A",
	KeyCtrlB:      "Ctrl-B",
	KeyCtrlC:      "Ctrl-C",
	KeyCtrlD:      "Ctrl-D",
	KeyCtrlE:      "Ctrl-E",
	KeyCtrlF:      "Ctrl-F",
	KeyCtrlG:      "Ctrl-G",
	KeyCtrlH:      "Ctrl-H",
	KeyCtrlI:      "Ctrl-I",
	KeyCtrlJ:      "Ctrl-J",
	KeyCtrlK:      "Ctrl-K",
	KeyCtrlL:      "Ctrl-L",
	KeyCtrlM:      "Ctrl-M",
	KeyCtrlN:      "Ctrl-N",
	KeyCtrlO:      "Ctrl-O",
	KeyCtrlP:      "Ctrl-P",
	KeyCtrlQ:      "Ctrl-Q",
	KeyCtrlR:      "Ctrl-R",
	KeyCtrlS:      "Ctrl-S",
	KeyCtrlT:      "Ctrl-T",
	KeyCtrlU:      "Ctrl-U",
	KeyCtrlV:      "Ctrl-V",
	KeyCtrlW:      "Ctrl-W",
	KeyCtrlX:      "Ctrl-X",
	KeyCtrlY:      "Ctrl-Y",
	KeyCtrlZ:      "Ctrl-Z",
}

// Name returns a printable value or the key stroke.  This can be used
// when printing the event, for example.
func (ev *EventKey) Name() string {
	s := ""
	m := []string{}
	if ev.mod&ModShift != 0 {
		m = append(m, "Shift")
	}
	if ev.mod&ModAlt != 0 {
		m = append(m, "Alt")
	}
	if ev.mod&ModMeta != 0 {
		m = append(m, "Meta")
	}
	if ev.mod&ModCtrl != 0 {
		m = append(m, "Ctrl")
	}
	if ev.mod&ModHyper != 0 {
		m = append(m, "Hyper")
	}

	ok := false
	if s, ok = KeyNames[ev.key]; !ok {
		if ev.key == KeyRune {
			s = "Rune[" + ev.str + "]"
		} else {
			s = fmt.Sprintf("Key[%d,%s]", ev.key, ev.str)
		}
	}
	if len(m) != 0 {
		if ev.mod&ModCtrl != 0 && strings.HasPrefix(s, "Ctrl-") {
			s = s[5:]
		}
		return fmt.Sprintf("%s+%s", strings.Join(m, "+"), s)
	}
	return s
}

// NewEventKey attempts to create a suitable event.  It parses the various
// ASCII control sequences if KeyRune is passed for Key, but if the caller
// has more precise information it should set that specifically.  Callers
// that aren't sure about modifier state (most) should just pass ModNone.
func NewEventKey(k Key, str string, mod ModMask) *EventKey {
	ch := rune(0)
	if len(str) == 1 {
		ch = []rune(str)[0]
	}

	if k == KeyRune {
		if ch != 0 && (ch < ' ' || ch == 0x7f) {
			// Turn specials into proper key codes.  This is for
			// control characters and the DEL.
			k = Key(ch)
			if mod == ModNone && ch < ' ' {
				switch k {
				case KeyBackspace, KeyTab, KeyEsc, KeyEnter:
					// these keys are directly typeable without CTRL
					str = ""
				default:
					// most likely entered with a CTRL keypress
					mod = ModCtrl
				}
				ch = ch + '\x60'
			}
		}

		// For legacy reasons, if Ctrl is pressed with an ASCII alphabetic, then we
		// emit it as a KeyCtrlXX symbol.
		if mod == ModCtrl {
			// We don't do Ctrl-[ or backslash or those specially.
			if ch >= 'A' && ch <= 'Z' { // upper case
				k = KeyCtrlA + Key(ch-'A')
				str = ""
			} else if ch >= 'a' && ch <= 'z' { // lower case
				k = KeyCtrlA + Key(ch-'a')
				str = ""
			}
		}

		// Windows reports ModShift for shifted keys.  This is inconsistent
		// with UNIX, lets harmonize this.
		if mod == ModShift && str != "" {
			mod = ModNone
		}
	}

	// Backspace2 is just another name for backspace.
	if k == KeyBackspace2 {
		k = KeyBackspace
	}

	// Shift-Tab should be Backtab.
	if k == KeyTab && (mod&ModShift) != 0 {
		k = KeyBacktab
		mod &^= ModShift
	}
	ev := &EventKey{key: k, str: str, mod: mod}
	ev.SetEventNow()
	return ev
}

// ModMask is a mask of modifier keys.  Note that it will not always be
// possible to report modifier keys.
type ModMask int16

// These are the modifiers keys that can be sent either with a key press,
// or a mouse event.  Note that as of now, due to the confusion associated
// with Meta, and the lack of support for it on many/most platforms, the
// current implementations never use it.  Instead, they use ModAlt, even for
// events that could possibly have been distinguished from ModAlt.
const (
	ModShift ModMask = 1 << iota
	ModCtrl
	ModAlt
	ModMeta
	ModHyper
	ModNone ModMask = 0
)

// Key is a generic value for representing keys, and especially special
// keys (function keys, cursor movement keys, etc.)  For normal keys, like
// ASCII letters, we use KeyRune, and then expect the application to
// inspect the Str() member of the EventKey.
type Key int16

// This is the list of named keys.  KeyRune is special however, in that it is
// a place holder key indicating that a printable character was sent.  The
// actual value of the rune will be transported in the Rune of the associated
// EventKey.
const (
	KeyRune Key = iota + 256
	KeyUp
	KeyDown
	KeyRight
	KeyLeft
	KeyUpLeft
	KeyUpRight
	KeyDownLeft
	KeyDownRight
	KeyCenter
	KeyPgUp
	KeyPgDn
	KeyHome
	KeyEnd
	KeyInsert
	KeyDelete
	KeyHelp
	KeyExit
	KeyClear
	KeyCancel
	KeyPrint
	KeyPause
	KeyBacktab
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
	KeyF21
	KeyF22
	KeyF23
	KeyF24
	KeyF25
	KeyF26
	KeyF27
	KeyF28
	KeyF29
	KeyF30
	KeyF31
	KeyF32
	KeyF33
	KeyF34
	KeyF35
	KeyF36
	KeyF37
	KeyF38
	KeyF39
	KeyF40
	KeyF41
	KeyF42
	KeyF43
	KeyF44
	KeyF45
	KeyF46
	KeyF47
	KeyF48
	KeyF49
	KeyF50
	KeyF51
	KeyF52
	KeyF53
	KeyF54
	KeyF55
	KeyF56
	KeyF57
	KeyF58
	KeyF59
	KeyF60
	KeyF61
	KeyF62
	KeyF63
	KeyF64
	KeyMenu
	KeyCapsLock
	KeyScrollLock
	KeyNumLock
)

const (
	// These key codes are used internally, and will never appear to applications.
	keyPasteStart Key = iota + 16384
	keyPasteEnd
)

// These are the control keys, they will also be reported with the
// rune (lower case) and control modifier.  If the shift key
// or other modifiers are present then these will *NOT* be reported,
// but reported instead as KeyRune.
const (
	KeyCtrlA Key = iota + 65
	KeyCtrlB
	KeyCtrlC
	KeyCtrlD
	KeyCtrlE
	KeyCtrlF
	KeyCtrlG
	KeyCtrlH
	KeyCtrlI
	KeyCtrlJ
	KeyCtrlK
	KeyCtrlL
	KeyCtrlM
	KeyCtrlN
	KeyCtrlO
	KeyCtrlP
	KeyCtrlQ
	KeyCtrlR
	KeyCtrlS
	KeyCtrlT
	KeyCtrlU
	KeyCtrlV
	KeyCtrlW
	KeyCtrlX
	KeyCtrlY
	KeyCtrlZ
)

// Special values - these are fixed in an attempt to make it more likely
// that aliases will encode the same way.

// These are the defined ASCII values for key codes.  They generally match
// with KeyCtrl values.
const (
	KeyNUL Key = iota
	KeySOH
	KeySTX
	KeyETX
	KeyEOT
	KeyENQ
	KeyACK
	KeyBEL
	KeyBS
	KeyTAB
	KeyLF
	KeyVT
	KeyFF
	KeyCR
	KeySO
	KeySI
	KeyDLE
	KeyDC1
	KeyDC2
	KeyDC3
	KeyDC4
	KeyNAK
	KeySYN
	KeyETB
	KeyCAN
	KeyEM
	KeySUB
	KeyESC
	KeyFS
	KeyGS
	KeyRS
	KeyUS
	KeyDEL Key = 0x7F
)

// These keys are aliases for other names.
const (
	KeyBackspace = KeyBS
	KeyTab       = KeyTAB
	KeyEsc       = KeyESC
	KeyEscape    = KeyESC
	KeyEnter     = KeyCR

	// NB: This key will be translated to KeyBackspace
	KeyBackspace2 = KeyDEL
)
