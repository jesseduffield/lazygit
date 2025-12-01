// Copyright 2025 The TCell Authors
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

// This file describes a generic VT input processor.  It parses key sequences,
// (input bytes) and loads them into events.  It expects UTF-8 or UTF-16 as the input
// feed, along with ECMA-48 sequences.  The assumption here is that all potential
// key sequences are unambiguous between terminal variants (analysis of extant terminfo
// data appears to support this conjecture). This allows us to implement  this once,
// in the most efficient and terminal-agnostic way possible.
//
// There is unfortunately *one* conflict, with aixterm, for CSI-P - which is KeyDelete
// in aixterm, but F1 in others.

package tcell

import (
	"encoding/base64"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf16"
	"unicode/utf8"
)

type inpState int

const (
	inpStateInit = inpState(iota)
	inpStateUtf
	inpStateEsc
	inpStateCsi // control sequence introducer
	inpStateOsc // operating system command
	inpStateDcs // device control string
	inpStateSos // start of string (unused)
	inpStatePm  // privacy message (unused)
	inpStateApc // application program command
	inpStateSt  // string terminator
	inpStateSs2 // single shift 2
	inpStateSs3 // single shift 3
	inpStateLFK // linux F-key (not ECMA-48 compliant - bogus CSI)
)

type InputProcessor interface {
	ScanUTF8([]byte)
	ScanUTF16([]uint16)
	SetSize(rows, cols int)
}

func NewInputProcessor(eq chan<- Event) InputProcessor {
	return &inputProcessor{
		evch: eq,
		buf:  make([]rune, 0, 128),
	}
}

type inputProcessor struct {
	ut8       []byte
	ut16      []uint16
	buf       []rune
	scratch   []byte
	csiParams []byte
	csiInterm []byte
	escaped   bool
	btnDown   bool // mouse button tracking for broken terms
	state     inpState
	strState  inpState // saved str state (needed for ST)
	timer     *time.Timer
	expire    time.Time
	l         sync.Mutex
	encBuf    []rune
	evch      chan<- Event
	rows      int // used for clipping mouse coordinates
	cols      int // used for clipping mouse coordinates
	nested    *inputProcessor
}

func (ip *inputProcessor) SetSize(w, h int) {
	if ip.nested != nil {
		ip.nested.SetSize(w, h)
		return
	}
	go func() {
		ip.l.Lock()
		ip.rows = h
		ip.cols = w
		ip.post(NewEventResize(w, h))
		ip.l.Unlock()
	}()
}
func (ip *inputProcessor) post(ev Event) {
	if ip.escaped {
		ip.escaped = false
		if ke, ok := ev.(*EventKey); ok {
			ev = NewEventKey(ke.Key(), ke.Rune(), ke.Modifiers()|ModAlt)
		}
	} else if ke, ok := ev.(*EventKey); ok {
		switch ke.Key() {
		case keyPasteStart:
			ev = NewEventPaste(true)
		case keyPasteEnd:
			ev = NewEventPaste(false)
		}
	}

	ip.evch <- ev
}

func (ip *inputProcessor) escTimeout() {
	ip.l.Lock()
	defer ip.l.Unlock()
	if ip.state == inpStateEsc && ip.expire.Before(time.Now()) {
		// post it
		ip.state = inpStateInit
		ip.escaped = false
		ip.post(NewEventKey(KeyEsc, 0, ModNone))
	}
}

type csiParamMode struct {
	M rune // Mode
	P int  // Parameter (first)
}

type keyMap struct {
	Key Key
	Mod ModMask
}

var csiAllKeys = map[csiParamMode]keyMap{
	{M: 'A'}:         {Key: KeyUp},
	{M: 'B'}:         {Key: KeyDown},
	{M: 'C'}:         {Key: KeyRight},
	{M: 'D'}:         {Key: KeyLeft},
	{M: 'F'}:         {Key: KeyEnd},
	{M: 'H'}:         {Key: KeyHome},
	{M: 'L'}:         {Key: KeyInsert},
	{M: 'P'}:         {Key: KeyF1}, // except for aixterm, where this is Delete
	{M: 'Q'}:         {Key: KeyF2},
	{M: 'S'}:         {Key: KeyF4},
	{M: 'Z'}:         {Key: KeyBacktab},
	{M: 'a'}:         {Key: KeyUp, Mod: ModShift},
	{M: 'b'}:         {Key: KeyDown, Mod: ModShift},
	{M: 'c'}:         {Key: KeyRight, Mod: ModShift},
	{M: 'd'}:         {Key: KeyLeft, Mod: ModShift},
	{M: 'q', P: 1}:   {Key: KeyF1}, // all these 'q' are for aixterm
	{M: 'q', P: 2}:   {Key: KeyF2},
	{M: 'q', P: 3}:   {Key: KeyF3},
	{M: 'q', P: 4}:   {Key: KeyF4},
	{M: 'q', P: 5}:   {Key: KeyF5},
	{M: 'q', P: 6}:   {Key: KeyF6},
	{M: 'q', P: 7}:   {Key: KeyF7},
	{M: 'q', P: 8}:   {Key: KeyF8},
	{M: 'q', P: 9}:   {Key: KeyF9},
	{M: 'q', P: 10}:  {Key: KeyF10},
	{M: 'q', P: 11}:  {Key: KeyF11},
	{M: 'q', P: 12}:  {Key: KeyF12},
	{M: 'q', P: 13}:  {Key: KeyF13},
	{M: 'q', P: 14}:  {Key: KeyF14},
	{M: 'q', P: 15}:  {Key: KeyF15},
	{M: 'q', P: 16}:  {Key: KeyF16},
	{M: 'q', P: 17}:  {Key: KeyF17},
	{M: 'q', P: 18}:  {Key: KeyF18},
	{M: 'q', P: 19}:  {Key: KeyF19},
	{M: 'q', P: 20}:  {Key: KeyF20},
	{M: 'q', P: 21}:  {Key: KeyF21},
	{M: 'q', P: 22}:  {Key: KeyF22},
	{M: 'q', P: 23}:  {Key: KeyF23},
	{M: 'q', P: 24}:  {Key: KeyF24},
	{M: 'q', P: 25}:  {Key: KeyF25},
	{M: 'q', P: 26}:  {Key: KeyF26},
	{M: 'q', P: 27}:  {Key: KeyF27},
	{M: 'q', P: 28}:  {Key: KeyF28},
	{M: 'q', P: 29}:  {Key: KeyF29},
	{M: 'q', P: 30}:  {Key: KeyF30},
	{M: 'q', P: 31}:  {Key: KeyF31},
	{M: 'q', P: 32}:  {Key: KeyF32},
	{M: 'q', P: 33}:  {Key: KeyF33},
	{M: 'q', P: 34}:  {Key: KeyF34},
	{M: 'q', P: 35}:  {Key: KeyF35},
	{M: 'q', P: 36}:  {Key: KeyF36},
	{M: 'q', P: 144}: {Key: KeyClear},
	{M: 'q', P: 146}: {Key: KeyEnd},
	{M: 'q', P: 150}: {Key: KeyPgUp},
	{M: 'q', P: 154}: {Key: KeyPgDn},
	{M: 'z', P: 214}: {Key: KeyHome},
	{M: 'z', P: 216}: {Key: KeyPgUp},
	{M: 'z', P: 220}: {Key: KeyEnd},
	{M: 'z', P: 222}: {Key: KeyPgDn},
	{M: 'z', P: 224}: {Key: KeyF1},
	{M: 'z', P: 225}: {Key: KeyF2},
	{M: 'z', P: 226}: {Key: KeyF3},
	{M: 'z', P: 227}: {Key: KeyF4},
	{M: 'z', P: 228}: {Key: KeyF5},
	{M: 'z', P: 229}: {Key: KeyF6},
	{M: 'z', P: 230}: {Key: KeyF7},
	{M: 'z', P: 231}: {Key: KeyF8},
	{M: 'z', P: 232}: {Key: KeyF9},
	{M: 'z', P: 233}: {Key: KeyF10},
	{M: 'z', P: 234}: {Key: KeyF11},
	{M: 'z', P: 235}: {Key: KeyF12},
	{M: 'z', P: 247}: {Key: KeyInsert},
	{M: '^', P: 7}:   {Key: KeyHome, Mod: ModCtrl},
	{M: '^', P: 8}:   {Key: KeyEnd, Mod: ModCtrl},
	{M: '^', P: 11}:  {Key: KeyF23},
	{M: '^', P: 12}:  {Key: KeyF24},
	{M: '^', P: 13}:  {Key: KeyF25},
	{M: '^', P: 14}:  {Key: KeyF26},
	{M: '^', P: 15}:  {Key: KeyF27},
	{M: '^', P: 17}:  {Key: KeyF28}, // 16 is a gap
	{M: '^', P: 18}:  {Key: KeyF29},
	{M: '^', P: 19}:  {Key: KeyF30},
	{M: '^', P: 20}:  {Key: KeyF31},
	{M: '^', P: 21}:  {Key: KeyF32},
	{M: '^', P: 23}:  {Key: KeyF33}, // 22 is a gap
	{M: '^', P: 24}:  {Key: KeyF34},
	{M: '^', P: 25}:  {Key: KeyF35},
	{M: '^', P: 26}:  {Key: KeyF36}, // 27 is a gap
	{M: '^', P: 28}:  {Key: KeyF37},
	{M: '^', P: 29}:  {Key: KeyF38}, // 30 is a gap
	{M: '^', P: 31}:  {Key: KeyF39},
	{M: '^', P: 32}:  {Key: KeyF40},
	{M: '^', P: 33}:  {Key: KeyF41},
	{M: '^', P: 34}:  {Key: KeyF42},
	{M: '@', P: 23}:  {Key: KeyF43},
	{M: '@', P: 24}:  {Key: KeyF44},
	{M: '$', P: 2}:   {Key: KeyInsert, Mod: ModShift},
	{M: '$', P: 3}:   {Key: KeyDelete, Mod: ModShift},
	{M: '$', P: 7}:   {Key: KeyHome, Mod: ModShift},
	{M: '$', P: 8}:   {Key: KeyEnd, Mod: ModShift},
	{M: '$', P: 23}:  {Key: KeyF21},
	{M: '$', P: 24}:  {Key: KeyF22},
	{M: '~', P: 1}:   {Key: KeyHome},
	{M: '~', P: 2}:   {Key: KeyInsert},
	{M: '~', P: 3}:   {Key: KeyDelete},
	{M: '~', P: 4}:   {Key: KeyEnd},
	{M: '~', P: 5}:   {Key: KeyPgUp},
	{M: '~', P: 6}:   {Key: KeyPgDn},
	{M: '~', P: 7}:   {Key: KeyHome},
	{M: '~', P: 8}:   {Key: KeyEnd},
	{M: '~', P: 11}:  {Key: KeyF1},
	{M: '~', P: 12}:  {Key: KeyF2},
	{M: '~', P: 13}:  {Key: KeyF3},
	{M: '~', P: 14}:  {Key: KeyF4},
	{M: '~', P: 15}:  {Key: KeyF5},
	{M: '~', P: 17}:  {Key: KeyF6},
	{M: '~', P: 18}:  {Key: KeyF7},
	{M: '~', P: 19}:  {Key: KeyF8},
	{M: '~', P: 20}:  {Key: KeyF9},
	{M: '~', P: 21}:  {Key: KeyF10},
	{M: '~', P: 23}:  {Key: KeyF11},
	{M: '~', P: 24}:  {Key: KeyF12},
	{M: '~', P: 25}:  {Key: KeyF13},
	{M: '~', P: 26}:  {Key: KeyF14},
	{M: '~', P: 28}:  {Key: KeyF15}, // aka KeyHelp
	{M: '~', P: 29}:  {Key: KeyF16},
	{M: '~', P: 31}:  {Key: KeyF17},
	{M: '~', P: 32}:  {Key: KeyF18},
	{M: '~', P: 33}:  {Key: KeyF19},
	{M: '~', P: 34}:  {Key: KeyF20},
	{M: '~', P: 200}: {Key: keyPasteStart},
	{M: '~', P: 201}: {Key: keyPasteEnd},
}

// keys reported using Kitty csi-u protocol
var csiUKeys = map[int]Key{
	27:    KeyESC,
	9:     KeyTAB,
	13:    KeyEnter,
	127:   KeyBS,
	57358: KeyCapsLock,
	57359: KeyScrollLock,
	57360: KeyNumLock,
	57361: KeyPrint,
	57362: KeyPause,
	57363: KeyMenu,
	57376: KeyF13,
	57377: KeyF14,
	57378: KeyF15,
	57379: KeyF16,
	57380: KeyF17,
	57381: KeyF18,
	57382: KeyF19,
	57383: KeyF20,
	57384: KeyF21,
	57385: KeyF22,
	57386: KeyF23,
	57387: KeyF24,
	57388: KeyF25,
	57389: KeyF26,
	57390: KeyF27,
	57391: KeyF28,
	57392: KeyF29,
	57393: KeyF30,
	57394: KeyF31,
	57395: KeyF32,
	57396: KeyF33,
	57397: KeyF34,
	57398: KeyF35,
	// TODO: KP keys
	// TODO: Media keys
}

// windows virtual key codes per microsoft
var winKeys = map[int]Key{
	0x03: KeyCancel,    // vkCancel
	0x08: KeyBackspace, // vkBackspace
	0x09: KeyTab,       // vkTab
	0x0d: KeyEnter,     // vkReturn
	0x12: KeyClear,     // vClear
	0x13: KeyPause,     // vkPause
	0x1b: KeyEscape,    // vkEscape
	0x21: KeyPgUp,      // vkPrior
	0x22: KeyPgDn,      // vkNext
	0x23: KeyEnd,       // vkEnd
	0x24: KeyHome,      // vkHome
	0x25: KeyLeft,      // vkLeft
	0x26: KeyUp,        // vkUp
	0x27: KeyRight,     // vkRight
	0x28: KeyDown,      // vkDown
	0x2a: KeyPrint,     // vkPrint
	0x2c: KeyPrint,     // vkPrtScr
	0x2d: KeyInsert,    // vkInsert
	0x2e: KeyDelete,    // vkDelete
	0x2f: KeyHelp,      // vkHelp
	0x70: KeyF1,        // vkF1
	0x71: KeyF2,        // vkF2
	0x72: KeyF3,        // vkF3
	0x73: KeyF4,        // vkF4
	0x74: KeyF5,        // vkF5
	0x75: KeyF6,        // vkF6
	0x76: KeyF7,        // vkF7
	0x77: KeyF8,        // vkF8
	0x78: KeyF9,        // vkF9
	0x79: KeyF10,       // vkF10
	0x7a: KeyF11,       // vkF11
	0x7b: KeyF12,       // vkF12
	0x7c: KeyF13,       // vkF13
	0x7d: KeyF14,       // vkF14
	0x7e: KeyF15,       // vkF15
	0x7f: KeyF16,       // vkF16
	0x80: KeyF17,       // vkF17
	0x81: KeyF18,       // vkF18
	0x82: KeyF19,       // vkF19
	0x83: KeyF20,       // vkF20
	0x84: KeyF21,       // vkF21
	0x85: KeyF22,       // vkF22
	0x86: KeyF23,       // vkF23
	0x87: KeyF24,       // vkF24
}

// keys by their SS3 - used in application mode usually (legacy VT-style)
var ss3Keys = map[rune]Key{
	'A': KeyUp,
	'B': KeyDown,
	'C': KeyRight,
	'D': KeyLeft,
	'F': KeyEnd,
	'H': KeyHome,
	'P': KeyF1,
	'Q': KeyF2,
	'R': KeyF3,
	'S': KeyF4,
	't': KeyF5,
	'u': KeyF6,
	'v': KeyF7,
	'l': KeyF8,
	'w': KeyF9,
	'x': KeyF10,
}

// linux terminal uses these non ECMA keys prefixed by CSI-[
var linuxFKeys = map[rune]Key{
	'A': KeyF1,
	'B': KeyF2,
	'C': KeyF3,
	'D': KeyF4,
	'E': KeyF5,
}

func (ip *inputProcessor) scan() {
	for _, r := range ip.buf {
		ip.buf = ip.buf[1:]
		if r > 0x7F {
			// 8-bit extended Unicode we just treat as such - this will swallow anything else queued up
			ip.state = inpStateInit
			ip.post(NewEventKey(KeyRune, r, ModNone))
			continue
		}
		switch ip.state {
		case inpStateInit:
			switch r {
			case '\x1b':
				// escape.. pending
				ip.state = inpStateEsc
				if len(ip.buf) == 0 && ip.nested == nil {
					ip.expire = time.Now().Add(time.Millisecond * 50)
					ip.timer = time.AfterFunc(time.Millisecond*60, ip.escTimeout)
				}
			case '\t':
				ip.post(NewEventKey(KeyTab, 0, ModNone))
			case '\b', '\x7F':
				ip.post(NewEventKey(KeyBackspace, 0, ModNone))
			case '\r':
				ip.post(NewEventKey(KeyEnter, 0, ModNone))
			default:
				// Control keys - legacy handling
				if r < ' ' {
					ip.post(NewEventKey(KeyCtrlSpace+Key(r), 0, ModCtrl))
				} else {
					ip.post(NewEventKey(KeyRune, r, ModNone))
				}
			}
		case inpStateEsc:
			switch r {
			case '[':
				ip.state = inpStateCsi
				ip.csiInterm = nil
				ip.csiParams = nil
			case ']':
				ip.state = inpStateOsc
				ip.scratch = nil
			case 'N':
				ip.state = inpStateSs2 // no known uses
				ip.scratch = nil
			case 'O':
				ip.state = inpStateSs3
				ip.scratch = nil
			case 'X':
				ip.state = inpStateSos
				ip.scratch = nil
			case '^':
				ip.state = inpStatePm
				ip.scratch = nil
			case '_':
				ip.state = inpStateApc
				ip.scratch = nil
			case '\\':
				// string terminator reached, (orphaned?)
				ip.state = inpStateInit
			case '\t':
				// Linux console only, does not conform to ECMA
				ip.state = inpStateInit
				ip.post(NewEventKey(KeyBacktab, 0, ModNone))
			default:
				if r == '\x1b' {
					// leading ESC to capture alt
					ip.escaped = true
				} else {
					// treat as alt-key ... legacy emulators only (no CSI-u or other)
					ip.state = inpStateInit
					mod := ModAlt
					if r < ' ' {
						mod |= ModCtrl
						r += 0x60
					}
					ip.post(NewEventKey(KeyRune, r, mod))
				}
			}
		case inpStateCsi:
			// usual case for incoming keys
			if r >= 0x30 && r <= 0x3F { // parameter bytes
				ip.csiParams = append(ip.csiParams, byte(r))
			} else if r >= 0x20 && r <= 0x2F { // intermediate bytes, rarely used
				ip.csiInterm = append(ip.csiInterm, byte(r))
			} else if r >= 0x40 && r <= 0x7F { // final byte
				ip.handleCsi(r, ip.csiParams, ip.csiInterm)
			} else {
				// bad parse, just swallow it all
				ip.state = inpStateInit
			}
		case inpStateSs2:
			// No known uses for SS2
			ip.state = inpStateInit

		case inpStateSs3: // typically application mode keys or older terminals
			ip.state = inpStateInit
			if k, ok := ss3Keys[r]; ok {
				ip.post(NewEventKey(k, 0, ModNone))
			}

		case inpStatePm, inpStateApc, inpStateSos, inpStateDcs: // these we just eat
			switch r {
			case '\x1b':
				ip.strState = ip.state
				ip.state = inpStateSt
			case '\x07': // bell - some send this instead of ST
				ip.state = inpStateInit
			}

		case inpStateOsc: // not sure if used
			switch r {
			case '\x1b':
				ip.strState = ip.state
				ip.state = inpStateSt
			case '\x07':
				ip.handleOsc(string(ip.scratch))
			default:
				ip.scratch = append(ip.scratch, byte(r&0x7f))
			}
		case inpStateSt:
			if r == '\\' || r == '\x07' {
				ip.state = inpStateInit
				switch ip.strState {
				case inpStateOsc:
					ip.handleOsc(string(ip.scratch))
				case inpStatePm, inpStateApc, inpStateSos, inpStateDcs:
					ip.state = inpStateInit
				}
			} else {
				ip.scratch = append(ip.scratch, '\x1b', byte(r))
				ip.state = ip.strState
			}
		case inpStateLFK:
			// linux console does not follow ECMA
			if k, ok := linuxFKeys[r]; ok {
				ip.post(NewEventKey(k, 0, ModNone))
			}
			ip.state = inpStateInit
		}
	}
}

func (ip *inputProcessor) handleOsc(str string) {
	ip.state = inpStateInit
	if content, ok := strings.CutPrefix(str, "52;c;"); ok {
		decoded := make([]byte, base64.StdEncoding.DecodedLen(len(content)))
		if count, err := base64.StdEncoding.Decode(decoded, []byte(content)); err == nil {
			ip.post(NewEventClipboard(decoded[:count]))
			return
		}
	}
}

func calcModifier(n int) ModMask {
	n--
	m := ModNone
	if n&1 != 0 {
		m |= ModShift
	}
	if n&2 != 0 {
		m |= ModAlt
	}
	if n&4 != 0 {
		m |= ModCtrl
	}
	if n&8 != 0 {
		m |= ModMeta // kitty calls this Super
	}
	if n&16 != 0 {
		m |= ModHyper
	}
	if n&32 != 0 {
		m |= ModMeta // for now not separating from Super
	}
	// Not doing (kitty only):
	// caps_lock 0b1000000   (64)
	// num_lock  0b10000000  (128)

	return m
}

// func (ip *inputProcessor) handleMouse(x, y, btn int, down bool) *EventMouse {
func (ip *inputProcessor) handleMouse(mode rune, params []int) {

	// XTerm mouse events only report at most one button at a time,
	// which may include a wheel button.  Wheel motion events are
	// reported as single impulses, while other button events are reported
	// as separate press & release events.
	if len(params) < 3 {
		return
	}
	btn := params[0]
	// Some terminals will report mouse coordinates outside the
	// screen, especially with click-drag events.  Clip the coordinates
	// to the screen in that case.
	x := max(min(params[1]-1, ip.cols-1), 0)
	y := max(min(params[2]-1, ip.rows-1), 0)
	motion := (btn & 0x20) != 0
	scroll := (btn & 0x42) == 0x40
	btn &^= 0x20
	if mode == 'm' {
		// mouse release, clear all buttons
		btn |= 3
		btn &^= 0x40
		ip.btnDown = false
	} else if motion {
		/*
		 * Some broken terminals appear to send
		 * mouse button one motion events, instead of
		 * encoding 35 (no buttons) into these events.
		 * We resolve these by looking for a non-motion
		 * event first.
		 */
		if !ip.btnDown {
			btn |= 3
			btn &^= 0x40
		}
	} else if !scroll {
		ip.btnDown = true
	}

	button := ButtonNone
	mod := ModNone

	// Mouse wheel has bit 6 set, no release events.  It should be noted
	// that wheel events are sometimes misdelivered as mouse button events
	// during a click-drag, so we debounce these, considering them to be
	// button press events unless we see an intervening release event.
	switch btn & 0x43 {
	case 0:
		button = Button1
	case 1:
		button = Button3 // Note we prefer to treat right as button 2
	case 2:
		button = Button2 // And the middle button as button 3
	case 3:
		button = ButtonNone
	case 0x40:
		button = WheelUp
	case 0x41:
		button = WheelDown
	case 0x42:
		button = WheelLeft
	case 0x43:
		button = WheelRight
	}

	if btn&0x4 != 0 {
		mod |= ModShift
	}
	if btn&0x8 != 0 {
		mod |= ModAlt
	}
	if btn&0x10 != 0 {
		mod |= ModCtrl
	}

	ip.post(NewEventMouse(x, y, button, mod))
}

func (ip *inputProcessor) handleWinKey(P []int) {
	// win32-input-mode
	//  ^[ [ Vk ; Sc ; Uc ; Kd ; Cs ; Rc _
	// Vk: the value of wVirtualKeyCode - any number. If omitted, defaults to '0'.
	// Sc: the value of wVirtualScanCode - any number. If omitted, defaults to '0'.
	// Uc: the decimal value of UnicodeChar - for example, NUL is "0", LF is
	//     "10", the character 'A' is "65". If omitted, defaults to '0'.
	// Kd: the value of bKeyDown - either a '0' or '1'. If omitted, defaults to '0'.
	// Cs: the value of dwControlKeyState - any number. If omitted, defaults to '0'.
	// Rc: the value of wRepeatCount - any number. If omitted, defaults to '1'.
	//
	// Note that some 3rd party terminal emulators (not Terminal) suffer from a bug
	// where other events, such as mouse events, are doubly encoded, using Vk 0
	// for each character.  (So a CSI-M sequence is encoded as a series of CSI-_
	// sequences.)  We consider this a bug in those terminal emulators -- Windows 11
	// Terminal does not suffer this brain damage. (We've observed this with both Alacritty
	// and WezTerm.)
	for len(P) < 6 {
		P = append(P, 0) // ensure sufficient length
	}
	if P[3] == 0 {
		// key up event ignore ignore
		return
	}

	if P[0] == 0 && P[1] == 0 && P[2] > 0 && P[2] < 0x80 { // only ASCII in win32-input-mode
		if ip.nested == nil {
			ip.nested = &inputProcessor{
				evch: ip.evch,
				rows: ip.rows,
				cols: ip.cols,
			}
		}

		ip.nested.ScanUTF8([]byte{byte(P[2])})
		return
	}

	key := KeyRune
	chr := rune(P[2])
	mod := ModNone
	rpt := max(1, P[5])
	if k1, ok := winKeys[P[0]]; ok {
		chr = 0
		key = k1
	} else if chr == 0 && P[0] >= 0x30 && P[0] <= 0x39 {
		chr = rune(P[0])
	} else if chr < ' ' && P[0] >= 0x41 && P[0] <= 0x5a {
		key = Key(P[0])
		chr = 0
	} else if key == 0x11 || key == 0x13 || key == 0x14 {
		// lone modifiers
		return
	}

	// Modifiers
	if P[4]&0x010 != 0 {
		mod |= ModShift
	}
	if P[4]&0x000c != 0 {
		mod |= ModCtrl
	}
	if P[4]&0x0003 != 0 {
		mod |= ModAlt
	}
	if key == KeyRune && chr > ' ' && mod == ModShift {
		// filter out lone shift for printable chars
		mod = ModNone
	}
	if chr != 0 && mod&(ModCtrl|ModAlt) == ModCtrl|ModAlt {
		// Filter out ctrl+alt (it means AltGr)
		mod = ModNone
	}

	for range rpt {
		if key != KeyRune || chr != 0 {
			ip.post(NewEventKey(key, chr, mod))
		}
	}
}

func (ip *inputProcessor) handleCsi(mode rune, params []byte, intermediate []byte) {

	// reset state
	ip.state = inpStateInit

	if len(intermediate) != 0 {
		// we don't know what to do with these for now
		return
	}

	var parts []string
	var P []int
	hasLT := false
	pstr := string(params)
	// extract numeric parameters
	if strings.HasPrefix(pstr, "<") {
		hasLT = true
		pstr = pstr[1:]
	}
	if pstr != "" && pstr[0] >= '0' && pstr[0] <= '9' {
		parts = strings.Split(pstr, ";")
		for i := range parts {
			if parts[i] != "" {
				if n, e := strconv.ParseInt(parts[i], 10, 32); e == nil {
					P = append(P, int(n))
				}
			}
		}
	}
	var P0 int
	if len(P) > 0 {
		P0 = P[0]
	}

	if hasLT {
		switch mode {
		case 'm', 'M': // mouse event, we only do SGR tracking
			ip.handleMouse(mode, P)
		}
	}

	switch mode {
	case 'I': // focus in
		ip.post(NewEventFocus(true))
		return
	case 'O': // focus out
		ip.post(NewEventFocus(false))
		return
	case '[':
		// linux console F-key - CSI-[ modifies next key
		ip.state = inpStateLFK
		return
	case 'u':
		// CSI-u kitty keyboard protocol
		if len(P) > 0 && !hasLT {
			mod := ModNone
			key := KeyRune
			chr := rune(0)
			if k1, ok := csiUKeys[P0]; ok {
				key = k1
				chr = 0
			} else {
				chr = rune(P0)
			}
			if len(P) > 1 {
				mod = calcModifier(P[1])
			}
			ip.post(NewEventKey(key, chr, mod))
		}
		return
	case '_':
		if len(intermediate) == 0 && len(P) > 0 {
			ip.handleWinKey(P)
			return
		}
	case '~':
		if len(intermediate) == 0 && len(P) >= 2 {
			mod := calcModifier(P[1])
			if ks, ok := csiAllKeys[csiParamMode{M: mode, P: P0}]; ok {
				ip.post(NewEventKey(ks.Key, 0, mod))
				return
			}
			if P0 == 27 && len(P) > 2 && P[2] > 0 && P[2] <= 0xff {
				if P[2] < ' ' || P[2] == 0x7F {
					ip.post(NewEventKey(Key(P[2]), 0, mod))
				} else {
					ip.post(NewEventKey(KeyRune, rune(P[2]), mod))
				}
				return
			}
		}
	}

	if ks, ok := csiAllKeys[csiParamMode{M: mode, P: P0}]; ok && !hasLT {
		if mode == '~' && len(P) > 1 && ks.Mod == ModNone {
			// apply modifiers if present
			ks.Mod = calcModifier(P[1])
		} else if mode == 'P' && os.Getenv("TERM") == "aixterm" {
			ks.Key = KeyDelete // aixterm hack - conflicts with kitty protocol
		}
		ip.post(NewEventKey(ks.Key, 0, ks.Mod))
		return
	}

	// this might have been an SS3 style key with modifiers applied
	if k, ok := ss3Keys[mode]; ok && P0 == 1 && len(P) > 1 {
		ip.post(NewEventKey(k, 0, calcModifier(P[1])))
		return
	}
	// if we got here we just swallow the unknown sequence
}

func (ip *inputProcessor) ScanUTF8(b []byte) {
	ip.l.Lock()
	defer ip.l.Unlock()

	ip.ut8 = append(ip.ut8, b...)
	for len(ip.ut8) > 0 {
		// fast path, basic ascii
		if ip.ut8[0] < 0x7F {
			ip.buf = append(ip.buf, rune(ip.ut8[0]))
			ip.ut8 = ip.ut8[1:]
		} else {
			r, len := utf8.DecodeRune(ip.ut8)
			if r == utf8.RuneError {
				r = rune(ip.ut8[0])
				len = 1
			}
			ip.buf = append(ip.buf, r)
			ip.ut8 = ip.ut8[len:]
		}
	}

	ip.scan()
}

func (ip *inputProcessor) ScanUTF16(u []uint16) {
	ip.l.Lock()
	defer ip.l.Unlock()
	ip.ut16 = append(ip.ut16, u...)
	for len(ip.ut16) > 0 {
		if !utf16.IsSurrogate(rune(ip.ut16[0])) {
			ip.buf = append(ip.buf, rune(ip.ut16[0]))
			ip.ut16 = ip.ut16[1:]
		} else if len(ip.ut16) > 1 {
			ip.buf = append(ip.buf, utf16.DecodeRune(rune(ip.ut16[0]), rune(ip.ut16[1])))
			ip.ut16 = ip.ut16[2:]
		} else {
			break
		}
	}
}
