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

// This file describes a generic VT input processor.  It parses key sequences,
// (input bytes) and loads them into events.  It expects UTF-8 or UTF-16 as the input
// feed, along with ECMA-48 sequences.  The assumption here is that all potential
// key sequences are unambiguous between terminal variants (analysis of extant terminfo
// data appears to support this conjecture). This allows us to implement  this once,
// in the most efficient and terminal-agnostic way possible.
//
// There is unfortunately *one* conflict, with aixterm, for CSI-P - which is KeyDelete
// in aixterm, but F1 in others.

//go:build !js && !wasm
// +build !js,!wasm

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

	"github.com/gdamore/tcell/v3/vt"
)

type inputState int

const (
	istInit = inputState(iota)
	istUtf  // utf8 state
	istEsc  // escape
	istCsi  // control sequence introducer
	istOsc  // operating system command
	istDcs  // device control string
	istSos  // start of string (unused)
	istPm   // privacy message (unused)
	istApc  // application program command
	istSt   // string terminator
	istSs2  // single shift 2
	istSs3  // single shift 3
	istLnx  // linux F-key (not ECMA-48 compliant - bogus CSI)
	istXda  // extended device attributes (ESC P Ps ST)
)

func newInputParser(eq chan<- Event) *inputParser {
	return &inputParser{
		evch: eq,
		buf:  make([]rune, 0, 128),
	}
}

type inputParser struct {
	buf       []rune       // bytes to process (ingest data)
	utfBuf    []byte       // accrued UTF8 bytes
	strBuf    []byte       // accrued string data (for ST, OSC, etc.)
	csiParams []byte       // accrued parameter bytes for CSI (and SS3)
	csiInterm []byte       // accrued intermediate bytes for CSI
	escChar   byte         // last byte for escape
	escaped   bool         // true if next key should be modified by ESC
	btnsDown  ButtonMask   // mouse buttons down (excludes wheel buttons)
	state     inputState   // tracks processor state
	strState  inputState   // saved str state (needed for ST)
	l         sync.Mutex   // protects local state
	evch      chan<- Event // where events are routed
	rows      int          // used for clipping mouse coordinates
	cols      int          // used for clipping mouse coordinates
	keyTime   time.Time    // time of last key press / byte ingested
	nested    *inputParser // for buggy win32-input-mode implementations
	surrogate rune         // high surrogate pair seen (for Win32 input mode)
}

// Waiting returns true if the processor is waiting for
// some more input (i.e. we are not in in the initial state.)
// This can occur when we have ambiguous escape sequences, such
// as the lone escape.  If this is typed, we expect at least a minimal
// inter-key delay before the next stroke occurs, and the caller
// should check for waiting, and call Scan() or ScanUTF8() to
// finish the processing.  (Typically after a delay of around 100ms.)
func (ip *inputParser) Waiting() bool {
	ip.l.Lock()
	defer ip.l.Unlock()
	return ip.state != istInit
}

func (ip *inputParser) SetSize(w, h int) {
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
func (ip *inputParser) post(ev Event) {
	if ip.escaped {
		ip.escaped = false
		if ke, ok := ev.(*EventKey); ok {
			ev = NewEventKey(ke.Key(), ke.Str(), ke.Modifiers()|ModAlt)
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

type csiParamMode struct {
	M rune // Mode
	P int  // Parameter (first)
}

type keyMap struct {
	Key  Key
	Mod  ModMask
	Rune rune
}

var csiAllKeys = map[csiParamMode]keyMap{
	{M: 'A'}:         {Key: KeyUp},
	{M: 'B'}:         {Key: KeyDown},
	{M: 'C'}:         {Key: KeyRight},
	{M: 'D'}:         {Key: KeyLeft},
	{M: 'E'}:         {Key: KeyClear},
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
	{M: '^', P: 1}:   {Key: KeyHome, Mod: ModCtrl},
	{M: '^', P: 2}:   {Key: KeyInsert, Mod: ModCtrl},
	{M: '^', P: 3}:   {Key: KeyDelete, Mod: ModCtrl},
	{M: '^', P: 4}:   {Key: KeyEnd, Mod: ModCtrl},
	{M: '^', P: 5}:   {Key: KeyPgUp, Mod: ModCtrl},
	{M: '^', P: 6}:   {Key: KeyPgDn, Mod: ModCtrl},
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
	{M: '@', P: 1}:   {Key: KeyHome, Mod: ModShift | ModCtrl},
	{M: '@', P: 2}:   {Key: KeyInsert, Mod: ModShift | ModCtrl},
	{M: '@', P: 3}:   {Key: KeyDelete, Mod: ModShift | ModCtrl},
	{M: '@', P: 4}:   {Key: KeyEnd, Mod: ModShift | ModCtrl},
	{M: '@', P: 5}:   {Key: KeyPgUp, Mod: ModShift | ModCtrl},
	{M: '@', P: 6}:   {Key: KeyPgDn, Mod: ModShift | ModCtrl},
	{M: '@', P: 7}:   {Key: KeyHome, Mod: ModShift | ModCtrl},
	{M: '@', P: 8}:   {Key: KeyEnd, Mod: ModShift | ModCtrl},
	{M: '$', P: 1}:   {Key: KeyHome, Mod: ModShift},
	{M: '$', P: 2}:   {Key: KeyInsert, Mod: ModShift},
	{M: '$', P: 3}:   {Key: KeyDelete, Mod: ModShift},
	{M: '$', P: 5}:   {Key: KeyPgUp, Mod: ModShift},
	{M: '$', P: 6}:   {Key: KeyPgDn, Mod: ModShift},
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
var csiUKeys = map[int]keyMap{
	27:    {Key: KeyESC},
	9:     {Key: KeyTAB},
	13:    {Key: KeyEnter},
	127:   {Key: KeyBS},
	57358: {Key: KeyCapsLock},
	57359: {Key: KeyScrollLock},
	57360: {Key: KeyNumLock},
	57361: {Key: KeyPrint},
	57362: {Key: KeyPause},
	57363: {Key: KeyMenu},
	57376: {Key: KeyF13},
	57377: {Key: KeyF14},
	57378: {Key: KeyF15},
	57379: {Key: KeyF16},
	57380: {Key: KeyF17},
	57381: {Key: KeyF18},
	57382: {Key: KeyF19},
	57383: {Key: KeyF20},
	57384: {Key: KeyF21},
	57385: {Key: KeyF22},
	57386: {Key: KeyF23},
	57387: {Key: KeyF24},
	57388: {Key: KeyF25},
	57389: {Key: KeyF26},
	57390: {Key: KeyF27},
	57391: {Key: KeyF28},
	57392: {Key: KeyF29},
	57393: {Key: KeyF30},
	57394: {Key: KeyF31},
	57395: {Key: KeyF32},
	57396: {Key: KeyF33},
	57397: {Key: KeyF34},
	57398: {Key: KeyF35},
	57399: {Key: KeyRune, Rune: '0'}, // KP 0
	57400: {Key: KeyRune, Rune: '1'}, // KP 1
	57401: {Key: KeyRune, Rune: '2'}, // KP 2
	57402: {Key: KeyRune, Rune: '3'}, // KP 3
	57403: {Key: KeyRune, Rune: '4'}, // KP 4
	57404: {Key: KeyRune, Rune: '5'}, // KP 5
	57405: {Key: KeyRune, Rune: '6'}, // KP 6
	57406: {Key: KeyRune, Rune: '7'}, // KP 7
	57407: {Key: KeyRune, Rune: '8'}, // KP 8
	57408: {Key: KeyRune, Rune: '9'}, // KP 9
	57409: {Key: KeyRune, Rune: '.'}, // KP_DECIMAL
	57410: {Key: KeyRune, Rune: '/'}, // KP_DIVIDE
	57411: {Key: KeyRune, Rune: '*'}, // KP_MULTIPLY
	57412: {Key: KeyRune, Rune: '-'}, // KP_SUBTRACT
	57413: {Key: KeyRune, Rune: '+'}, // KP_ADD
	57414: {Key: KeyEnter},           // KP_ENTER
	57415: {Key: KeyRune, Rune: '='}, // KP_EQUAL
	57416: {Key: KeyClear},           // KP_SEPARATOR
	57417: {Key: KeyLeft},            // KP_LEFT
	57418: {Key: KeyRight},           // KP_RIGHT
	57419: {Key: KeyUp},              // KP_UP
	57420: {Key: KeyDown},            // KP_DOWN
	57421: {Key: KeyPgUp},            // KP_PG_UP
	57422: {Key: KeyPgDn},            // KP_PG_DN
	57423: {Key: KeyHome},            // KP_HOME
	57424: {Key: KeyEnd},             // KP_END
	57425: {Key: KeyInsert},          // KP_INSERT
	57426: {Key: KeyDelete},          // KP_DELETE
	// 57427: {Key: KeyBegin},          // KP_BEGIN

	// TODO: Media keys
}

// windows virtual key codes per microsoft
var winKeys = map[int]Key{
	0x03: KeyCancel,    // vkCancel
	0x08: KeyBackspace, // vkBackspace
	0x09: KeyTab,       // vkTab
	0x0d: KeyEnter,     // vkReturn
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
	'E': KeyClear,
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

func (ip *inputParser) scan() {
	for _, r := range ip.buf {
		ip.buf = ip.buf[1:]
		ip.escChar = 0
		ip.keyTime = time.Now()
		if r >= 0xA0 {
			// 8-bit extended Unicode we just treat as such - this will swallow anything else queued up
			ip.state = istInit
			ip.post(NewEventKey(KeyRune, string(r), ModNone))
			continue
		} else if r >= 0x80 {
			// ISO 2022 control chars
			ip.state = istEsc
			r -= 0x40
			// we fall through so it will be treated as the 7-bit equivalent
		}
		switch ip.state {
		case istInit:
			switch r {
			case '\x1b':
				// escape.. pending
				ip.state = istEsc
				ip.escChar = 0
			case '\t':
				ip.post(NewEventKey(KeyTab, "", ModNone))
			case '\b', '\x7F':
				ip.post(NewEventKey(KeyBackspace, "", ModNone))
			case '\r':
				ip.post(NewEventKey(KeyEnter, "", ModNone))
			default:
				// Control keys - legacy handling
				if r == 0 {
					ip.post(NewEventKey(KeyRune, " ", ModCtrl))
				} else if r < ' ' {
					ip.post(NewEventKey(KeyRune, string(r+0x40), ModCtrl))
				} else {
					ip.post(NewEventKey(KeyRune, string(r), ModNone))
				}
			}
		case istEsc:
			switch r {
			case '[':
				ip.state = istCsi
				ip.csiInterm = nil
				ip.csiParams = nil
				ip.escChar = byte(r)
			case ']':
				ip.state = istOsc
				ip.strBuf = nil
				ip.escChar = byte(r)
			case 'N':
				ip.state = istSs2 // no known uses
				ip.strBuf = nil
				ip.escChar = byte(r)
			case 'O':
				ip.state = istSs3
				ip.csiParams = nil
				ip.strBuf = nil
				ip.escChar = byte(r)
			case 'P':
				ip.state = istXda
				ip.csiParams = nil
				ip.strBuf = nil
				ip.escChar = byte(r)
			case 'X':
				ip.state = istSos
				ip.strBuf = nil
				ip.escChar = byte(r)
			case '^':
				ip.state = istPm
				ip.strBuf = nil
				ip.escChar = byte(r)
			case '_':
				ip.state = istApc
				ip.strBuf = nil
				ip.escChar = byte(r)
			case '\\':
				// string terminator reached, (orphaned?)
				ip.state = istInit
			case '\t':
				// Linux console only, does not conform to ECMA
				ip.state = istInit
				ip.post(NewEventKey(KeyBacktab, "", ModNone))
			default:
				if r == '\x1b' {
					// leading ESC to capture alt
					ip.escaped = true
					ip.escChar = byte(r)
				} else {
					// treat as alt-key ... legacy emulators only (no CSI-u or other)
					ip.state = istInit
					mod := ModAlt
					if r < ' ' {
						mod |= ModCtrl
						r += 0x60
					}
					ip.post(NewEventKey(KeyRune, string(r), mod))
				}
			}
		case istCsi:
			// usual case for incoming keys
			// NB: rxvt uses terminating '$' which is not a legal CSI terminator,
			// for certain shifted key sequences.  We special case this, and it's ok
			// because no other terminal seems to use this for CSI intermediates from
			// the terminal to the host (queries in the other direction can use it.)
			// However, this is only true if the first parameter does not have a "?",
			// because it *does* collide with DEC private mode queries otherwise.
			if r == '\x1b' {
				// Per ECMA-48 §5.3.1, ESC restarts the escape
				// sequence machine from any intermediate state.
				ip.state = istEsc
				ip.escChar = 0
			} else if r >= 0x30 && r <= 0x3F { // parameter bytes
				ip.csiParams = append(ip.csiParams, byte(r))
			} else if r == '$' && len(ip.csiParams) > 0 && ip.csiParams[0] != '?' { // rxvt non-standard
				ip.handleCsi(r, ip.csiParams, ip.csiInterm)
			} else if r >= 0x20 && r <= 0x2F { // intermediate bytes, rarely used
				ip.csiInterm = append(ip.csiInterm, byte(r))
			} else if r >= 0x40 && r <= 0x7F { // final byte
				ip.handleCsi(r, ip.csiParams, ip.csiInterm)
			} else {
				// bad parse, just swallow it all
				ip.state = istInit
			}
		case istSs2:
			// No known uses for SS2
			ip.state = istInit

		case istSs3: // typically application mode keys or older terminals
			ip.state = istInit
			// some SS3 sequences (old VTE) encode modifiers here just like CSI
			if r == '\x1b' {
				// Per ECMA-48 §5.3.1, ESC restarts the escape
				// sequence machine from any intermediate state.
				ip.state = istEsc
				ip.escChar = 0
			} else if r >= 0x30 && r <= 0x3F {
				ip.csiParams = append(ip.csiParams, byte(r))
				ip.state = istSs3
			} else if k, ok := ss3Keys[r]; ok {
				// If there are no parameters, then it's simple without modifiers.
				// The options for parameters are "1;<modifiers>" , or ";modifiers" (empty
				// first parameter defaults to 1), or just <modifiers>.  If a sequence has
				// parameters that do not match one of these forms, we just discard it.
				if len(ip.csiParams) == 0 {
					// simple SS3 case
					ip.post(NewEventKey(k, "", ModNone))
				} else if parts := strings.Split(string(ip.csiParams), ";"); len(parts) >= 1 {
					// SS3 with modifier (old style).  Note old terminfo would declare these as high
					// numbered function keys, but we encode as modified since that's how they are entered.
					if len(parts) >= 2 {
						if m, err := strconv.Atoi(parts[1]); err == nil && (parts[0] == "1" || parts[0] == "") {
							ip.post(NewEventKey(k, "", calcModifier(m)))
						}
					} else if m, err := strconv.Atoi(parts[0]); err == nil {
						ip.post(NewEventKey(k, "", calcModifier(m)))
					}
				}
			}

		case istPm, istApc, istSos, istDcs: // these we just eat
			switch r {
			case '\x1b':
				ip.strState = ip.state
				ip.state = istSt
			case '\x07': // bell - some send this instead of ST
				ip.state = istInit
			}

		case istXda:
			switch r {
			case '\x1b':
				ip.strState = ip.state
				ip.state = istSt
			case '\x07':
				ip.handleXda(string(ip.strBuf))
			default:
				ip.strBuf = append(ip.strBuf, byte(r&0x7f))
			}

		case istOsc: // not sure if used
			switch r {
			case '\x1b':
				ip.strState = ip.state
				ip.state = istSt
			case '\x07':
				ip.handleOsc(string(ip.strBuf))
			default:
				ip.strBuf = append(ip.strBuf, byte(r&0x7f))
			}
		case istSt:
			if r == '\\' || r == '\x07' {
				ip.state = istInit
				switch ip.strState {
				case istOsc:
					ip.handleOsc(string(ip.strBuf))
				case istXda:
					ip.handleXda(string(ip.strBuf))
				case istPm, istApc, istSos, istDcs:
					ip.state = istInit
				}
			} else {
				ip.strBuf = append(ip.strBuf, '\x1b', byte(r))
				ip.state = ip.strState
			}
		case istLnx:
			// linux console does not follow ECMA
			if k, ok := linuxFKeys[r]; ok {
				ip.post(NewEventKey(k, "", ModNone))
			}
			ip.state = istInit
		}
	}

	if ip.state != istInit && time.Since(ip.keyTime) > time.Millisecond*50 {
		if ip.state == istEsc {
			ip.post(NewEventKey(KeyEscape, "", ModNone))
		} else if ec := ip.escChar; ec != 0 {
			ip.post(NewEventKey(KeyRune, string(ec), ModAlt))
		}
		// if we take too long between bytes, reset the state machine.
		ip.state = istInit
	}
}

func (ip *inputParser) handleOsc(str string) {
	ip.state = istInit
	if content, ok := strings.CutPrefix(str, "52;c;"); ok {
		decoded := make([]byte, base64.StdEncoding.DecodedLen(len(content)))
		if count, err := base64.StdEncoding.Decode(decoded, []byte(content)); err == nil {
			ip.post(NewEventClipboard(decoded[:count]))
			return
		}
	}
}

func (ip *inputParser) handleXda(str string) {
	ip.state = istInit
	if content, ok := strings.CutPrefix(str, ">|"); ok {
		// two approaches, one with version like (1.23) another with just spaces
		if name, vers, ok := strings.Cut(content, "("); ok && strings.HasSuffix(vers, ")") {
			name = strings.TrimSpace(name)
			vers = strings.TrimSpace(strings.TrimSuffix(vers, ")"))
			ip.post(&eventTermName{Name: name, Version: vers})
		} else if name, vers, ok = strings.Cut(content, " "); ok {
			ip.post(&eventTermName{Name: name, Version: vers})
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

func (ip *inputParser) handleMouse(mode rune, params []int) {

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

	button := ButtonNone
	mod := ModNone

	// Mouse wheel has bit 6 set, no release events.  It should be noted
	// that wheel events are sometimes misdelivered as mouse button events
	// during a click-drag, so we debounce these, considering them to be
	// button press events unless we see an intervening release event.
	// This excludes motion (bit 5) and modifiers (bits 2, 3, 4) for now.
	switch btn & 0xC3 {
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
	case 0x80:
		button = Button4
	case 0x81:
		button = Button5
	case 0x82:
		button = Button6
	case 0x83:
		button = Button7
	}

	switch mode {
	case 'm':
		if (ip.btnsDown & button) == 0 {
			// a release without a corresponding press, so clear it
			button = ButtonNone
		} else {
			ip.btnsDown &^= button
			button = ip.btnsDown
		}

	case 'M':
		if btn&0x20 != 0 && button != ButtonNone && (ip.btnsDown&button) == 0 {
			// Ghostty may send out motion signals that indicate a button has
			// been pressed, even when the button is not actually pressed.
			// Do not create a synthetic button-down state from these packets.
			button = ip.btnsDown
			break
		}
		// record this press
		ip.btnsDown |= button
		// and use the full set so can see chords
		button = ip.btnsDown
		// mice wheel do not have release events
		ip.btnsDown &^= (WheelDown | WheelUp | WheelLeft | WheelRight)
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

func (ip *inputParser) handleWinKey(P []int) {
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

	// these terminals never send ambiguous escapes
	ip.escaped = false

	if P[0] == 0 && P[1] == 0 && P[2] > 0 && P[2] < 0x80 { // only ASCII in win32-input-mode
		if ip.nested == nil {
			ip.nested = &inputParser{
				evch: ip.evch,
				rows: ip.rows,
				cols: ip.cols,
			}
		}
		if P[2] > 0 {
			ip.nested.ScanUTF8([]byte{byte(P[2])})
		}
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
	} else if chr >= 0xD800 && chr <= 0xDBFF {
		// high surrogate pair
		ip.surrogate = chr
		return
	} else if chr >= 0xDC00 && chr <= 0xDFFF {
		// low surrogate pair
		chr = utf16.DecodeRune(ip.surrogate, chr)
	} else if P[0] == 0x10 || P[0] == 0x11 || P[0] == 0x12 || P[0] == 0x14 {
		// lone modifiers
		ip.surrogate = 0
		return
	}

	ip.surrogate = 0

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
		if key != KeyRune {
			ip.post(NewEventKey(key, "", mod))
		} else if chr != 0 {
			ip.post(NewEventKey(KeyRune, string(chr), mod))
		}
	}
}

func (ip *inputParser) handlePrimaryDA(params []int) {
	if len(params) < 1 {
		return
	}
	evDA := &eventPrimaryAttributes{Class: params[0]}
	params = params[1:]
	if evDA.Class >= 60 {
		for _, v := range params {
			switch v {
			case 3:
				evDA.ReGIS = true
			case 4:
				evDA.Sixel = true
			case 9:
				evDA.National = true
			case 12:
				evDA.SerboCroation = true
			case 22:
				evDA.Color = true
			case 23:
				evDA.Greek = true
			case 24:
				evDA.Turkish = true
			case 42:
				evDA.Latin2 = true
			case 52:
				evDA.Clipboard = true
			}
		}
	}
	ip.post(evDA)
}

func (ip *inputParser) handlePrivateModeResponse(params []int) {
	for len(params) < 2 {
		params = append(params, 0)
	}
	if params[1] >= 0 && params[1] <= 4 {
		ev := &eventPrivateMode{
			Mode:   vt.PrivateMode(params[0]),
			Status: vt.ModeStatus(params[1]),
		}
		ip.post(ev)
	}
}

func (ip *inputParser) handleKittyMode(params []int) {
	if len(params) == 1 && params[0] >= 0 && params[0] < 32 {
		ev := &eventKittyKbdMode{
			Mode: KittyKbdMode(params[0] & 0xffff),
		}
		ip.post(ev)
	}
}

func (ip *inputParser) handleXTermMode(params []int) {
	if len(params) >= 1 && params[0] == 4 {
		if len(params) == 1 {
			params = append(params, 0)
		}
		ev := &eventXTermKbdMode{
			Mode: XtermKbdMode(params[1] & 0x3),
		}
		ip.post(ev)
	}
}

func (ip *inputParser) handleCsi(mode rune, params []byte, intermediate []byte) {

	// reset state
	ip.state = istInit

	var parts []string
	var P []int
	hasLT := false
	hasQM := false
	hasGT := false
	pstr := string(params)
	// extract numeric parameters
	if strings.HasPrefix(pstr, "<") {
		hasLT = true
		pstr = pstr[1:]
	} else if strings.HasPrefix(pstr, "?") {
		hasQM = true
		pstr = pstr[1:]
	} else if strings.HasPrefix(pstr, ">") {
		hasGT = true
		pstr = pstr[1:]
	}

	if pstr != "" && pstr[0] >= '0' && pstr[0] <= '9' {
		parts = strings.Split(pstr, ";")
		for i := range parts {
			if parts[i] != "" {
				if n, e := strconv.ParseInt(parts[i], 10, 32); e == nil {
					P = append(P, int(n))
				}
			} else {
				P = append(P, 0)
			}
		}
	}
	var P0 int
	if len(P) > 0 {
		P0 = P[0]
	}

	if hasLT && len(intermediate) == 0 {
		switch mode {
		case 'm', 'M': // mouse event, we only do SGR tracking
			ip.handleMouse(mode, P)
		}
		return
	}
	if hasQM {
		switch mode {
		case 'c':
			if len(intermediate) == 0 {
				ip.handlePrimaryDA(P)
			}
		case 'y':
			if string(intermediate) == "$" {
				ip.handlePrivateModeResponse(P)
			}
		case 'u':
			if len(intermediate) == 0 {
				ip.handleKittyMode(P)
			}
		}
		return
	}
	if hasGT {
		switch mode {
		case 'm':
			if len(intermediate) == 0 {
				ip.handleXTermMode(P)
			}
		}
		return
	}

	if len(intermediate) != 0 {
		// we don't know what to do with these for now
		return
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
		ip.state = istLnx
		return
	case 'u':
		// CSI-u kitty keyboard protocol, is unambiguous
		if len(P) > 0 {
			mod := ModNone
			key := KeyRune
			chr := rune(0)
			if k1, ok := csiUKeys[P0]; ok {
				key = k1.Key
				chr = k1.Rune
			} else {
				chr = rune(P0)
			}
			if len(P) > 1 {
				mod = calcModifier(P[1])
			}
			if key != KeyRune {
				ip.post(NewEventKey(key, "", mod))
			} else if chr != 0 {
				ip.post(NewEventKey(KeyRune, string(chr), mod))
			}
			return
		}
	case '_':
		if len(P) > 0 {
			ip.handleWinKey(P)
			return
		}
	case 't':
		if len(P) < 1 {
			break
		}
		switch P[0] {
		case 8:
			if len(P) > 2 {
				// window size report
				h := P[1]
				w := P[2]
				if h != ip.rows || w != ip.cols {
					ip.SetSize(w, h)
				}
				return
			}
		case 48:
			if len(P) > 2 {
				// window resize report
				ip.post(NewEventResize(P[2], P[1]))
				return
			}
		}
	case '~':
		if len(P) >= 2 {
			mod := calcModifier(P[1])
			if ks, ok := csiAllKeys[csiParamMode{M: mode, P: P0}]; ok {
				ip.post(NewEventKey(ks.Key, "", mod))
				return
			}
			if P0 == 27 && len(P) > 2 && P[2] > 0 && P[2] <= 0xff {
				if P[2] < ' ' || P[2] == 0x7F {
					ip.post(NewEventKey(Key(P[2]), "", mod))
				} else {
					ip.post(NewEventKey(KeyRune, string(rune(P[2])), mod))
				}
				return
			}
		}
	}

	if ks, ok := csiAllKeys[csiParamMode{M: mode, P: P0}]; ok {
		if mode == '~' && len(P) > 1 && ks.Mod == ModNone {
			// apply modifiers if present
			ks.Mod = calcModifier(P[1])
		} else if mode == 'P' && os.Getenv("TERM") == "aixterm" {
			ks.Key = KeyDelete // aixterm hack - conflicts with kitty protocol
		}
		ip.post(NewEventKey(ks.Key, "", ks.Mod))
		return
	}

	// this might have been an SS3 style key with modifiers applied
	if k, ok := ss3Keys[mode]; ok && P0 == 1 && len(P) > 1 {
		ip.post(NewEventKey(k, "", calcModifier(P[1])))
		return
	}
	// if we got here we just swallow the unknown sequence
}

func (ip *inputParser) ScanUTF8(b []byte) {
	ip.l.Lock()
	defer ip.l.Unlock()

	ip.utfBuf = append(ip.utfBuf, b...)
	for len(ip.utfBuf) > 0 {
		// fast path, basic ascii, also includes ISO2022 8-bit controls
		if ip.utfBuf[0] < 0xA0 {
			ip.buf = append(ip.buf, rune(ip.utfBuf[0]))
			ip.utfBuf = ip.utfBuf[1:]
		} else {
			r, utfLen := utf8.DecodeRune(ip.utfBuf)
			if r == utf8.RuneError {
				// discard the leading byte as bad,
				// hopefully it will recover.
				utfLen = 1
			} else {
				ip.buf = append(ip.buf, r)
			}
			ip.utfBuf = ip.utfBuf[utfLen:]
		}
	}

	ip.scan()
}

// Scan scans the existing input, but does not take new content.
// This is typically called after a delay when Waiting() is true.
func (ip *inputParser) Scan() {
	ip.l.Lock()
	ip.scan()
	ip.l.Unlock()
}

// Private events between input and tscreen.

// eventPrimaryAttributes is for primary device attributes -- this should be
// the last event returned during initial handshaking
type eventPrimaryAttributes struct {
	EventTime
	Class         int  // Terminal class, 1 is vt100, vt101, 6 is vt102, > 60 for vt200 and up
	ReGIS         bool // Terminal supports ReGIS graphics (DA 3)
	Sixel         bool // Terminal supports Sixel graphics (DA 4)
	National      bool // Terminal supports national replacement character sets (DA 9)
	SerboCroation bool // Serbo-Croatian(DA 12)
	Color         bool // Terminal supports color (DA 22)
	Greek         bool // Greek (DA 23)
	Turkish       bool // Turkish (DA 24)
	Latin2        bool // ISO Latin-2 (DA 42)
	Clipboard     bool // OSC 52 support (DA 52)
}

// eventTermName is for extended attributes
type eventTermName struct {
	EventTime
	Name    string
	Version string
}

type eventPrivateMode struct {
	EventTime
	Mode   vt.PrivateMode // numeric mode e.g. 7 for auto-margin, 1006 for SGR mouse reports, etc
	Status vt.ModeStatus  // value of status
}

type KittyKbdMode uint16

const (
	KittyKbdModeOff       = KittyKbdMode(0)  // Disable Kitty keyboard mode
	KittyKbdModeBase      = KittyKbdMode(1)  // Enable disambiguated keys
	KittyKbdModeEvents    = KittyKbdMode(2)  // Report event types (e.g. key release)
	KittyKbdModeAlternate = KittyKbdMode(4)  // Report alternate keys
	KittyKbdModeAll       = KittyKbdMode(8)  // Report all keys using kitty keyboard protocol
	KittyKbdModeText      = KittyKbdMode(16) // Report associated text
)

type eventKittyKbdMode struct {
	EventTime
	Mode KittyKbdMode
}

type XtermKbdMode uint16

const (
	XtermKbdModeOff  = XtermKbdMode(0) // Disabled
	XtermKbdModeBase = XtermKbdMode(1) // Enabled except for ones with legacy behavior
	XtermKbdModeExt  = XtermKbdMode(2) // Enabled for all modified keys
	XtermKbdModeAll  = XtermKbdMode(3) // Send all keys (including unmodified)
)

type eventXTermKbdMode struct {
	EventTime
	Mode XtermKbdMode
}
