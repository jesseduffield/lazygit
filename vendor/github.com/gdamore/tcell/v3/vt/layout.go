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

import (
	"sort"
	"sync"
	"time"
)

// This file defines the layout mechanism for keyboards.
// All keyboards are assumed to pass "codes" that are mapped to one of
// the Key values (this is easy for USB because the codes are mostly the
// same as USB already!), and then we look up values.  This can be done
// using a simple map, and we allow inheritance so we don't have to redefine
// too many things.  For complex scenarios you'll probably want to make use
// of operating system facilities or supply your own tables (for example for
// Japanese keyboards.)

// ModifierMap is a map that represents keys that generate runes in various
// shift states (modifier states). A layout may have many of these, and they
// are searched until a match is fine.
type ModifierMap struct {
	When func(Modifier) bool // Map only applies if this returns true
	Map  map[Key]rune        // Map is the mapping from Key to specific rune when this map matches.
}

// KeyboardState represents the current state of the keyboard.
// A zero initialized value is ready for use.  Note that emulators that
// get the associated events from their operating system do not need to
// make use of this, but this structure makes it possible to build an emulator
// with a keyboard layout that is not known to the operating system.
//
// The keyboard state is assumed to be "single threaded", meaning only a single
// caller will operate on it at any given time.  Typically there is just a single
// keyboard polling thread or goroutine.
type KeyboardState struct {
	deadKey        map[rune]DeadKey // current dead key state
	mod            Modifier
	layout         *Layout
	lastKey        Key           // last key pressed
	lastRune       rune          // last rune for last key
	repeating      bool          // true if we are repeating
	repeatStart    time.Time     // when we started repeating
	repeatTime     time.Time     // last time we checked
	repeatDelay    time.Duration // delay before starting repeat
	repeatInterval time.Duration // duration between repeats
	pressed        map[Key]bool
	initialized    bool
}

// initialize the keyboard, lazily.
func (ks *KeyboardState) initialize() {
	if !ks.initialized {
		ks.pressed = make(map[Key]bool)
		if ks.layout == nil {
			ks.layout = KeyboardANSI
		}
		ks.repeatDelay = time.Millisecond * 250
		ks.repeatInterval = time.Millisecond * 30
		ks.initialized = true
	}
}

// reset the keyboard state.
func (ks *KeyboardState) reset() {
	ks.clearRepeat()
	ks.mod = 0
	ks.deadKey = nil
	ks.pressed = make(map[Key]bool)
}

// clear repeat clears any repeating key.
func (ks *KeyboardState) clearRepeat() {
	ks.repeating = false
	ks.repeatStart = time.Time{}
	ks.repeatTime = time.Time{}
	ks.lastRune = 0
	ks.lastKey = 0
}

// SetRepeat sets the repeat parameters. Note that this will only have any meaningful
// impact if the caller calls the Pressed function repeatedly (periodically) while
// a key is depressed.
//
// The repeat starts after a key has been held for for delay, with a new repeat
// added every interval.
//
// The caller should usually call this before processing keyboard events.  It must
// not be called concurrently with either of the Pressed or Release functions.
func (ks *KeyboardState) SetRepeat(delay time.Duration, interval time.Duration) {
	ks.initialize()
	ks.repeatDelay = delay
	ks.repeatInterval = interval
	ks.clearRepeat()
}

// SetLayout sets the layout this keyboard should use.
// This also resets the keyboard state.
func (ks *KeyboardState) SetLayout(km *Layout) {
	ks.initialize()
	ks.reset()
	ks.layout = km
}

// Pressed should be called when a given key is depressed.
func (ks *KeyboardState) Pressed(k Key) *KeyEvent {
	ks.initialize()
	event := &KeyEvent{
		Down: true,
		Key:  k,
		SC:   k.ScanCode(),
		Base: k.KittyBase(),
		Mod:  ks.mod,
	}
	// if another key was pressed, then clear the repeat state
	lastKey := ks.lastKey
	if lastKey != k && ks.repeatInterval != 0 {
		ks.clearRepeat()
		ks.repeatStart = time.Now().Add(ks.repeatDelay)
		ks.repeatTime = ks.repeatStart
	}
	wasPressed := ks.pressed[k]
	ks.pressed[k] = true
	ks.lastKey = k

	l := ks.layout
	event.VK = l.Virtual[k]
	if mod, ok := l.Locking[k]; ok {
		// locking modifiers never repeat
		if wasPressed {
			return nil
		}
		if ks.mod&mod == 0 {
			ks.mod |= mod
		} else {
			ks.mod &^= mod
		}
		ks.pressed[k] = true
		event.Mod = ks.mod
		return event
	}
	if mod, ok := l.Modifiers[k]; ok {
		if wasPressed {
			return nil
		}
		ks.mod |= mod
		event.Mod = ks.mod
		return event
	}

	// attempt to look up the rune for this
	r := l.KeyToUTF(k, ks.mod)
	if ks.deadKey == nil && l.DeadKeys != nil && r != 0 {
		if dk, ok := l.DeadKeys[r]; ok {
			ks.deadKey = dk.Next
			return event
		}
	}

	if dk := ks.deadKey; dk != nil {
		if n, ok := dk[r]; ok {
			if n.U != 0 {
				event.Utf = string(n.U)
				ks.lastRune = n.U
				ks.deadKey = nil
			} else {
				ks.deadKey = n.Next
			}
			return event
		}
		// failed lookup - ignore it
		ks.deadKey = nil
	}

	if r != 0 {
		event.Utf = string(r)
	}

	if lastKey == k && wasPressed && ks.repeatInterval > 0 {
		ks.repeating = true
		if time.Now().After(ks.repeatStart) {
			deltaT := time.Since(ks.repeatTime).Truncate(ks.repeatInterval)
			event.Repeat = int(deltaT / ks.repeatInterval)
			if ks.repeatTime == ks.repeatStart {
				// fence post - count the first one!
				event.Repeat++
			}
			ks.repeatTime = ks.repeatTime.Add(deltaT)
			// if we polled before the repeat interval then report nothing
			if event.Repeat == 0 {
				return nil
			}
		}
	}

	return event
}

// Released should be called when a given key is Released.
func (ks *KeyboardState) Released(k Key) *KeyEvent {
	ks.initialize()
	event := &KeyEvent{
		Down: false,
		Key:  k,
		SC:   k.ScanCode(),
		Base: k.KittyBase(),
		Mod:  ks.mod,
	}
	if ks.lastKey == k {
		ks.clearRepeat()
	}
	wasPressed := ks.pressed[k]
	delete(ks.pressed, k)

	l := ks.layout
	event.VK = l.Virtual[k]

	if mod, ok := ks.layout.Modifiers[k]; ok {
		if !wasPressed {
			// something weird
			return nil
		}
		ks.mod &^= mod

		event.Mod = ks.mod
		return event
	}

	if _, ok := ks.layout.Locking[k]; !ok {
		if r := l.KeyToUTF(k, ks.mod); r != 0 {
			event.Utf = string(r)
		}
	}

	// no real point in looking up UTF for key release, so we don't

	return event
}

// DeadKey is what happens when a dead key is pressed.  Either it starts an unresolved sequence,
// (in which case Next will be non-nil), or it resolves to a final rune (in which case U will be non-zero)
// This is also sometimes called a composed key sequence.
type DeadKey struct {
	Next map[rune]DeadKey // Next corresponds to the next key in the sequence for dead keys
	U    rune             // U is the rune that should be emitted on completion of the sequence.
}

// DeadRune is used internally to create a rune to represent a dead key.
// While any numbers can be used from this point on, the expectation is
// that this will be added to a smaller rune (such as an ASCII character)
// This rune value a bit inside the supplementary private use area B,
// in an attempt to minimize the chance of any conflicting use (for
// example nerd fonts.)
const DeadRune = 0x101000

// Layout is a structure that represents a keyboard layout.
// Applications or users may implement their own layouts by
// registering an instance of this.
type Layout struct {
	// Name is the name of the keyboard layout.
	// We prefer to use the same names that Microsoft uses for keyboard layouts.
	Name string

	// Base is a base keyboard layout, so that we can simplify by only
	// overriding keys we are are handling differently.
	Base *Layout

	// DeadKeys maps specific starting keys, to an emitted rune.
	// The keys should a rune value starting with DeadRune.
	DeadKeys map[rune]DeadKey

	// Locking are modifiers that toggle a locked state like NumLock or CapsLock.
	Locking map[Key]Modifier

	// Modifiers toggle a state, but only while the key is depressed.
	Modifiers map[Key]Modifier

	// Virtual maps keys to virtual keys.
	Virtual map[Key]WinVK

	// Maps is the list of key maps by modifier and mask.
	// Use more specific masks first - for example NumLock
	// mask might contain only the masks for the number keypad,
	// and that should be listed before the maps for the rest
	// of the keyboard.  The algorithm searches for the first
	// match, not the best match.
	Maps []ModifierMap
}

func (km *Layout) KeyToUTF(k Key, m Modifier) rune {
	for _, mm := range km.Maps {
		if mm.When != nil && !mm.When(m) {
			continue
		}
		if u, ok := mm.Map[k]; ok {
			return u
		}
	}
	if km.Base != nil {
		return km.Base.KeyToUTF(k, m)
	}
	return 0
}

// KeysUsLower is a list of lower case key maps.
var KeysUsLower = map[Key]rune{
	KeyA: 'a',
	KeyB: 'b',
	KeyC: 'c',
	KeyD: 'd',
	KeyE: 'e',
	KeyF: 'f',
	KeyG: 'g',
	KeyH: 'h',
	KeyI: 'i',
	KeyJ: 'j',
	KeyK: 'k',
	KeyL: 'l',
	KeyM: 'm',
	KeyN: 'n',
	KeyO: 'o',
	KeyP: 'p',
	KeyQ: 'q',
	KeyR: 'r',
	KeyS: 's',
	KeyT: 't',
	KeyU: 'u',
	KeyV: 'v',
	KeyW: 'w',
	KeyX: 'x',
	KeyY: 'y',
	KeyZ: 'z',
}

// KeysUsUpper is a list of upper case key maps.
var KeysUsUpper = map[Key]rune{
	KeyA: 'A',
	KeyB: 'B',
	KeyC: 'C',
	KeyD: 'D',
	KeyE: 'E',
	KeyF: 'F',
	KeyG: 'G',
	KeyH: 'H',
	KeyI: 'I',
	KeyJ: 'J',
	KeyK: 'K',
	KeyL: 'L',
	KeyM: 'M',
	KeyN: 'N',
	KeyO: 'O',
	KeyP: 'P',
	KeyQ: 'Q',
	KeyR: 'R',
	KeyS: 'S',
	KeyT: 'T',
	KeyU: 'U',
	KeyV: 'V',
	KeyW: 'W',
	KeyX: 'X',
	KeyY: 'Y',
	KeyZ: 'Z',
}

// KeysDigits is a map of number keys to corresponding digits.
var KeysDigits = map[Key]rune{
	Key1: '1',
	Key2: '2',
	Key3: '3',
	Key4: '4',
	Key5: '5',
	Key6: '6',
	Key7: '7',
	Key8: '8',
	Key9: '9',
	Key0: '0',
}

// KeysPadDigits is a map of the digits on the numeric keypad,
// when num lock is engaged.
var KeysPadDigits = map[Key]rune{
	KeyPad0:   '0',
	KeyPad1:   '1',
	KeyPad2:   '2',
	KeyPad3:   '3',
	KeyPad4:   '4',
	KeyPad5:   '5',
	KeyPad6:   '6',
	KeyPad7:   '7',
	KeyPad8:   '8',
	KeyPad9:   '9',
	KeyPadDec: '.',
}

var KeysPadOps = map[Key]rune{
	KeyPadMul: '*',
	KeyPadAdd: '+',
	KeyPadSub: '-',
	KeyPadDiv: '/',
}

// KeyboardANSI is the base PC keyboard used in ANSI (US) systems.
var KeyboardANSI = &Layout{
	Name: "US",

	Base: nil,

	// Virtual maps physical keys to virtual keys.
	Virtual: map[Key]WinVK{
		KeyEsc:          VkEscape,
		KeyF1:           VkF1,
		KeyF2:           VkF2,
		KeyF3:           VkF3,
		KeyF4:           VkF4,
		KeyF5:           VkF5,
		KeyF6:           VkF6,
		KeyF7:           VkF7,
		KeyF8:           VkF8,
		KeyF9:           VkF9,
		KeyF10:          VkF10,
		KeyF11:          VkF11,
		KeyF12:          VkF12,
		KeyF13:          VkF13,
		KeyF14:          VkF14,
		KeyF15:          VkF15,
		KeyF16:          VkF16,
		KeyF17:          VkF17,
		KeyF18:          VkF18,
		KeyF19:          VkF19,
		KeyF20:          VkF20,
		KeyF21:          VkF21,
		KeyF22:          VkF22,
		KeyF23:          VkF23,
		KeyF24:          VkF24,
		KeyPrtScr:       VkSnapshot,
		KeyScrLock:      VkScroll,
		KeyPause:        VkPause,
		KeyInsert:       VkInsert,
		KeyDelete:       VkDelete,
		KeyHome:         VkHome,
		KeyEnd:          VkEnd,
		KeyPgUp:         VkPrior,
		KeyPgDn:         VkNext,
		KeyLeft:         VkLeft,
		KeyRight:        VkRight,
		KeyUp:           VkUp,
		KeyDown:         VkDown,
		KeyNumLock:      VkNumLock,
		KeyCapsLock:     VkCapital,
		KeyLShift:       VkLShift,
		KeyLCtrl:        VkLControl,
		KeyLMeta:        VkLWin,
		KeyLAlt:         VkLMenu,
		KeyRShift:       VkRShift,
		KeyRCtrl:        VkRControl,
		KeyRMeta:        VkRWin,
		KeyRAlt:         VkRMenu,
		KeyMenu:         VkApps,
		KeyConvert:      VkConvert,
		KeyNonConvert:   VkNonConvert,
		KeyBackspace:    VkBack,
		KeyEnter:        VkReturn,
		KeySpace:        VkSpace,
		KeyTab:          VkTab,
		Key1:            Vk1,
		Key2:            Vk2,
		Key3:            Vk3,
		Key4:            Vk4,
		Key5:            Vk5,
		Key6:            Vk6,
		Key7:            Vk7,
		Key8:            Vk8,
		Key9:            Vk9,
		Key0:            Vk0,
		KeyA:            VkA,
		KeyB:            VkB,
		KeyC:            VkC,
		KeyD:            VkD,
		KeyE:            VkE,
		KeyF:            VkF,
		KeyG:            VkG,
		KeyH:            VkH,
		KeyI:            VkI,
		KeyJ:            VkJ,
		KeyK:            VkK,
		KeyL:            VkL,
		KeyM:            VkM,
		KeyN:            VkN,
		KeyO:            VkO,
		KeyP:            VkP,
		KeyQ:            VkQ,
		KeyR:            VkR,
		KeyS:            VkS,
		KeyT:            VkT,
		KeyU:            VkU,
		KeyV:            VkV,
		KeyW:            VkW,
		KeyX:            VkX,
		KeyY:            VkY,
		KeyZ:            VkZ,
		KeyPadMul:       VkMultiply,
		KeyPadAdd:       VkAdd,
		KeyPadSub:       VkSubtract,
		KeyPadDiv:       VkDivide,
		KeyEqual:        VkOemPlus,
		KeyComma:        VkOemComma,
		KeyMinus:        VkOemMinus,
		KeyPeriod:       VkOemPeriod,
		KeySlash:        VkOem2,
		KeyGrave:        VkOem3,
		KeyLBrace:       VkOem4,
		KeyBackslash:    VkOem5,
		KeyRBrace:       VkOem6,
		KeyQuote:        VkOem7,
		KeyIsoBackSlash: VkOem102,
		KeyPad0:         VkInsert,
		KeyPad1:         VkEnd,
		KeyPad2:         VkDown,
		KeyPad3:         VkNext,
		KeyPad4:         VkLeft,
		KeyPad5:         VkClear,
		KeyPad6:         VkRight,
		KeyPad7:         VkHome,
		KeyPad8:         VkUp,
		KeyPad9:         VkPrior,
		KeyPadDec:       VkDelete,
	},
	Maps: []ModifierMap{
		// Specials - without control
		{
			When: func(m Modifier) bool { return !m.IsCtrl() },
			Map: map[Key]rune{
				KeyTab:       '\t',
				KeyEnter:     '\r',
				KeyBackspace: '\b',
				KeyEsc:       '\x1b',
				KeySpace:     ' ',
				KeyPadEnter:  '\r',
			},
		},
		// Specials - with control (but without shift)
		{
			When: func(m Modifier) bool { return m.IsCtrl() && !m.IsShift() },
			Map: map[Key]rune{
				KeyTab:       '\t',
				KeyEnter:     '\n',
				KeyBackspace: '\x7f',
				KeyEsc:       '\x1b',
				KeySpace:     ' ',
				KeyPadEnter:  '\n',
			},
		},
		// Key pad operators
		{
			When: func(m Modifier) bool { return !m.IsAlt() },
			Map:  KeysPadOps,
		},
		// Numeric keypad when num lock is engaged
		{
			When: func(m Modifier) bool { return m.IsNumLock() },
			Map:  KeysPadDigits,
		},
		// Numbers - without shift
		{
			When: func(m Modifier) bool { return !m.IsShift() && !m.IsCtrl() },
			Map:  KeysDigits,
		},
		// Numbers - with shift - this is locale sensitive usually
		{
			When: func(m Modifier) bool { return m.IsShift() && !m.IsCtrl() },
			Map: map[Key]rune{
				Key1: '!',
				Key2: '@',
				Key3: '#',
				Key4: '$',
				Key5: '%',
				Key6: '^',
				Key7: '&',
				Key8: '*',
				Key9: '(',
				Key0: ')',
			},
		},
		// Special shift-control cases
		{
			When: func(m Modifier) bool { return m.IsCtrl() && m.IsShift() },
			Map: map[Key]rune{
				Key2:     0,
				Key6:     '\x1e',
				KeyMinus: '\x1f',
			},
		},
		// Letters - base (lower case)
		{
			When: func(m Modifier) bool { return !m.IsCtrl() && !m.IsCapitals() },
			Map:  KeysUsLower,
		},
		// Letters - capitals (either caps lock or shift, but not both)
		{
			When: func(m Modifier) bool { return !m.IsCtrl() && m.IsCapitals() },
			Map:  KeysUsUpper,
		},
		// OEM keys - base
		{
			When: func(m Modifier) bool { return !m.IsShift() && !m.IsCtrl() },
			Map: map[Key]rune{
				KeySemi:         ';',
				KeyEqual:        '=',
				KeyComma:        ',',
				KeyMinus:        '-',
				KeyPeriod:       '.',
				KeySlash:        '/',
				KeyGrave:        '`',
				KeyLBrace:       '[',
				KeyBackslash:    '\\',
				KeyRBrace:       ']',
				KeyQuote:        '\'',
				KeyIsoBackSlash: '\\',
			},
		},
		// OEM keys - shift
		{
			When: func(m Modifier) bool { return m.IsShift() && !m.IsCtrl() },
			Map: map[Key]rune{
				KeySemi:         ':',
				KeyEqual:        '+',
				KeyComma:        '<',
				KeyMinus:        '_',
				KeyPeriod:       '>',
				KeySlash:        '?',
				KeyGrave:        '~',
				KeyLBrace:       '{',
				KeyBackslash:    '|',
				KeyRBrace:       '}',
				KeyQuote:        '"',
				KeyIsoBackSlash: '|',
			},
		},
		// OEM keys - control (odd balls)
		{
			When: func(m Modifier) bool { return m.IsCtrl() && !m.IsShift() },
			Map: map[Key]rune{
				KeyLBrace:       '\x1b',
				KeyBackslash:    '\x1c',
				KeyRBrace:       '\x1d',
				KeyIsoBackSlash: '\x1c',
			},
		},
	},
	Modifiers: map[Key]Modifier{
		KeyLShift: ModLShift,
		KeyRShift: ModRShift,
		KeyLCtrl:  ModLCtrl,
		KeyRCtrl:  ModRCtrl,
		KeyLAlt:   ModLAlt,
		KeyRAlt:   ModRAlt,
		KeyRMeta:  ModRMeta,
		KeyLMeta:  ModLMeta,
		KeyRHyper: ModRHyper,
		KeyLHyper: ModLHyper,
	},
	Locking: map[Key]Modifier{
		KeyNumLock:  ModNumLock,
		KeyCapsLock: ModCapsLock,
	},
}

var allLayouts = map[string]*Layout{
	KeyboardANSI.Name: KeyboardANSI,
}

var layoutsLock sync.Mutex

// RegisterLayout registers the given layout.
func RegisterLayout(km *Layout) {
	layoutsLock.Lock()
	allLayouts[km.Name] = km
	layoutsLock.Unlock()
}

// GetLayout returns a keyboard layout for the given name.
// The layout must have been previously registered with RegisterLayout.
// (Builtin layouts do this as a consequence of importing the layout.)
func GetLayout(name string) *Layout {
	layoutsLock.Lock()
	defer layoutsLock.Unlock()
	return allLayouts[name]
}

// Layouts returns a list of all known layout names.
func Layouts() []string {
	layoutsLock.Lock()
	defer layoutsLock.Unlock()
	res := make([]string, 0, len(allLayouts))
	for k := range allLayouts {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}
