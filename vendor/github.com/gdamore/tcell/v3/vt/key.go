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

// BaseKey is the Kitty protocol base key. These are kitty's representation of a scan code.
// As the Kitty protocol is likely the extended keyboard protocol we care about, we use
// this as the primary reporting mechanism. (It also helps that this may provide an easier
// fallback for implementations that don't have raw scan codes and are willing to assume an
// ANSI layout.)
type BaseKey rune

var shiftedBaseKeys map[BaseKey]rune

func (bk BaseKey) Shifted() rune {

	if s, ok := shiftedBaseKeys[bk]; ok {
		return s
	}
	if bk >= 'a' && bk <= 'z' {
		return rune(bk - 32)
	}
	return rune(bk)
}

// ScanCode is the scan code used by Windows for a key. These are physical key locations,
// and every physical should have a exactly one mapping here.
type ScanCode uint16

// Key represents the key code for a key. These are largely taken from the HID specification page 0x07,
// although there are gaps and some inconsistencies.  We chose this specification because it covers all
// keyboards we are likely to see in practice. Note that the zero value is reserved and will never be
// assigned a valid key.  A number of keys are from UNIX keyboards (e.g. Sun type 6 keyboards), and we do not
// explicitly support these because none of the common reporting protocols have a way to report them.
// Instead these should probably be mapped higher functions (F13 and up), if you're facing this.
// These are key locations on a US keyboard.  For example on AZERTY layout KeyQ corresponds to a key
// that has a printed "A", but is in the upper left position (below the digits.)
type Key uint16

const (
	KeyUTF          = Key(0x01) // virtual UTF-8 content
	KeyA            = Key(0x04)
	KeyB            = Key(0x05)
	KeyC            = Key(0x06)
	KeyD            = Key(0x07)
	KeyE            = Key(0x08)
	KeyF            = Key(0x09)
	KeyG            = Key(0x0A)
	KeyH            = Key(0x0B)
	KeyI            = Key(0x0C)
	KeyJ            = Key(0x0D)
	KeyK            = Key(0x0E)
	KeyL            = Key(0x0F)
	KeyM            = Key(0x10)
	KeyN            = Key(0x11)
	KeyO            = Key(0x12)
	KeyP            = Key(0x13)
	KeyQ            = Key(0x14)
	KeyR            = Key(0x15)
	KeyS            = Key(0x16)
	KeyT            = Key(0x17)
	KeyU            = Key(0x18)
	KeyV            = Key(0x19)
	KeyW            = Key(0x1A)
	KeyX            = Key(0x1B)
	KeyY            = Key(0x1C)
	KeyZ            = Key(0x1D)
	Key1            = Key(0x1E)
	Key2            = Key(0x1F)
	Key3            = Key(0x20)
	Key4            = Key(0x21)
	Key5            = Key(0x22)
	Key6            = Key(0x23)
	Key7            = Key(0x24)
	Key8            = Key(0x25)
	Key9            = Key(0x26)
	Key0            = Key(0x27)
	KeyEnter        = Key(0x28)
	KeyEsc          = Key(0x29)
	KeyBackspace    = Key(0x2a) // sometimes called delete
	KeyTab          = Key(0x2b)
	KeySpace        = Key(0x2c)
	KeyMinus        = Key(0x2d) // - and _
	KeyEqual        = Key(0x2e) // = and +
	KeyLBrace       = Key(0x2f) // [ and {
	KeyRBrace       = Key(0x30) // ] and }
	KeyBackslash    = Key(0x31) // \
	KeyIsoHash      = Key(0x32) // international only
	KeySemi         = Key(0x33) // ;
	KeyQuote        = Key(0x34) // ' and " aka apostrophe
	KeyGrave        = Key(0x35) // ` and ~
	KeyComma        = Key(0x36) // , and <
	KeyPeriod       = Key(0x37) // . and >
	KeySlash        = Key(0x38) // / and ?
	KeyCapsLock     = Key(0x39)
	KeyF1           = Key(0x3a)
	KeyF2           = Key(0x3b)
	KeyF3           = Key(0x3c)
	KeyF4           = Key(0x3d)
	KeyF5           = Key(0x3e)
	KeyF6           = Key(0x3f)
	KeyF7           = Key(0x40)
	KeyF8           = Key(0x41)
	KeyF9           = Key(0x42)
	KeyF10          = Key(0x43)
	KeyF11          = Key(0x44)
	KeyF12          = Key(0x45)
	KeyPrtScr       = Key(0x46)
	KeyScrLock      = Key(0x47)
	KeyPause        = Key(0x48)
	KeyInsert       = Key(0x49)
	KeyHome         = Key(0x4a)
	KeyPgUp         = Key(0x4b)
	KeyDelete       = Key(0x4c) // forward delete (DEL)
	KeyEnd          = Key(0x4d)
	KeyPgDn         = Key(0x4e)
	KeyRight        = Key(0x4f)
	KeyLeft         = Key(0x50)
	KeyDown         = Key(0x51)
	KeyUp           = Key(0x52)
	KeyNumLock      = Key(0x53) // also clear
	KeyPadDiv       = Key(0x54)
	KeyPadMul       = Key(0x55)
	KeyPadSub       = Key(0x56)
	KeyPadAdd       = Key(0x57)
	KeyPadEnter     = Key(0x58)
	KeyPad1         = Key(0x59) // also pad end
	KeyPad2         = Key(0x5a) // also pad down
	KeyPad3         = Key(0x5b) // also pad page down
	KeyPad4         = Key(0x5c) // also pad left
	KeyPad5         = Key(0x5d)
	KeyPad6         = Key(0x5e) // also pad right
	KeyPad7         = Key(0x5f) // also pad home
	KeyPad8         = Key(0x60) // also pad up
	KeyPad9         = Key(0x61) // also pad page up
	KeyPad0         = Key(0x62) // also pad insert
	KeyPadDec       = Key(0x63) // also pad delete
	KeyIsoBackSlash = Key(0x64) // international keyboards only
	KeyPadEqual     = Key(0x67)
	KeyF13          = Key(0x68)
	KeyF14          = Key(0x69)
	KeyF15          = Key(0x6a)
	KeyF16          = Key(0x6b)
	KeyF17          = Key(0x6c)
	KeyF18          = Key(0x6d)
	KeyF19          = Key(0x6e)
	KeyF20          = Key(0x6f)
	KeyF21          = Key(0x70)
	KeyF22          = Key(0x71)
	KeyF23          = Key(0x72)
	KeyF24          = Key(0x73)
	KeyMenu         = Key(0x76) // might also be 0x65, KeyApplication
	KeyPadComma     = Key(0x85)
	KeyPadEqSign    = Key(0x86)
	KeyIsoSlash     = Key(0x87) // also International 1
	KeyHiragana     = Key(0x88) // also International 2
	KeyYen          = Key(0x89) // found on JIS keyboards
	KeyConvert      = Key(0x8a)
	KeyNonConvert   = Key(0x8b)
	KeyAltErase     = Key(0x99) // e.g split space-erase bar
	KeySysReq       = Key(0x9a)
	KeyCancel       = Key(0x9b)
	KeyLCtrl        = Key(0xe0)
	KeyLShift       = Key(0xe1)
	KeyLAlt         = Key(0xe2)
	KeyLMeta        = Key(0xe3)
	KeyRCtrl        = Key(0xe4)
	KeyRShift       = Key(0xe5)
	KeyRAlt         = Key(0xe6)
	KeyRMeta        = Key(0xe7)
	KeyLHyper       = Key(0xff01) // software only, in reserved space
	KeyRHyper       = Key(0xff02) // software only, in reserved space
)

// scanCodes is a list of Windows scan codes for physical keys.
var scanCodes map[Key]ScanCode

// ScanCode returns the corresponding Windows Scan Code (not a VK!)
// for the given key. (Virtual keys should be determined by the
// host OS using the current keyboard layout.)
func (k Key) ScanCode() ScanCode {

	if w, ok := scanCodes[k]; ok {
		return w
	}
	return 0
}

// WinVK represents a windows virtual key code.
// These are similar to base keys, but a multiple scanned key codes
// may result in the same virtual key.  This can also be sensitive to
// the keyboard layout.
type WinVK rune

const (
	VkBack       = WinVK(0x08) // backspace
	VkTab        = WinVK(0x09)
	VkClear      = WinVK(0x0c)
	VkReturn     = WinVK(0x0d)
	VkShift      = WinVK(0x10)
	VkControl    = WinVK(0x11)
	VkMenu       = WinVK(0x12)
	VkPause      = WinVK(0x13)
	VkCapital    = WinVK(0x14) // caps lock
	VkKana       = WinVK(0x15)
	VkHangul     = WinVK(0x15)
	VkImeOn      = WinVK(0x16)
	VkJunja      = WinVK(0x17)
	VkFinal      = WinVK(0x18)
	VkKanji      = WinVK(0x19)
	VkImeOff     = WinVK(0x1a)
	VkEscape     = WinVK(0x1b)
	VkConvert    = WinVK(0x1c)
	VkNonConvert = WinVK(0x1d)
	VkAccept     = WinVK(0x1e)
	VkModeChange = WinVK(0x1f)
	VkSpace      = WinVK(0x20)
	VkPrior      = WinVK(0x21) // page up
	VkNext       = WinVK(0x22) // page down
	VkEnd        = WinVK(0x23)
	VkHome       = WinVK(0x24)
	VkLeft       = WinVK(0x25)
	VkUp         = WinVK(0x26)
	VkRight      = WinVK(0x27)
	VkDown       = WinVK(0x28)
	VkSelect     = WinVK(0x29)
	VkPrint      = WinVK(0x2a)
	VkExecute    = WinVK(0x2b)
	VkSnapshot   = WinVK(0x2c) // print screen
	VkInsert     = WinVK(0x2D)
	VkDelete     = WinVK(0x2E)
	VkHelp       = WinVK(0x2F)
	Vk0          = WinVK(0x30)
	Vk1          = WinVK(0x31)
	Vk2          = WinVK(0x32)
	Vk3          = WinVK(0x33)
	Vk4          = WinVK(0x34)
	Vk5          = WinVK(0x35)
	Vk6          = WinVK(0x36)
	Vk7          = WinVK(0x37)
	Vk8          = WinVK(0x38)
	Vk9          = WinVK(0x39)
	VkA          = WinVK(0x41)
	VkB          = WinVK(0x42)
	VkC          = WinVK(0x43)
	VkD          = WinVK(0x44)
	VkE          = WinVK(0x45)
	VkF          = WinVK(0x46)
	VkG          = WinVK(0x47)
	VkH          = WinVK(0x48)
	VkI          = WinVK(0x49)
	VkJ          = WinVK(0x4a)
	VkK          = WinVK(0x4b)
	VkL          = WinVK(0x4c)
	VkM          = WinVK(0x4d)
	VkN          = WinVK(0x4e)
	VkO          = WinVK(0x4f)
	VkP          = WinVK(0x50)
	VkQ          = WinVK(0x51)
	VkR          = WinVK(0x52)
	VkS          = WinVK(0x53)
	VkT          = WinVK(0x54)
	VkU          = WinVK(0x55)
	VkV          = WinVK(0x56)
	VkW          = WinVK(0x57)
	VkX          = WinVK(0x58)
	VkY          = WinVK(0x59)
	VkZ          = WinVK(0x5a)
	VkLWin       = WinVK(0x5b) // left meta
	VkRWin       = WinVK(0x5c) // right meta
	VkApps       = WinVK(0x5d) // menu key
	VkNumPad0    = WinVK(0x60)
	VkNumPad1    = WinVK(0x61)
	VkNumPad2    = WinVK(0x62)
	VkNumPad3    = WinVK(0x63)
	VkNumPad4    = WinVK(0x64)
	VkNumPad5    = WinVK(0x65)
	VkNumPad6    = WinVK(0x66)
	VkNumPad7    = WinVK(0x67)
	VkNumPad8    = WinVK(0x68)
	VkNumPad9    = WinVK(0x69)
	VkMultiply   = WinVK(0x6a) // pad multiply
	VkAdd        = WinVK(0x6b) // pad add
	VkSeparator  = WinVK(0x6c) // separator
	VkSubtract   = WinVK(0x6d) // pad subtract
	VkDecimal    = WinVK(0x6e) // pad decimal point
	VkDivide     = WinVK(0x6f) // pad divide
	VkF1         = WinVK(0x70)
	VkF2         = WinVK(0x71)
	VkF3         = WinVK(0x72)
	VkF4         = WinVK(0x73)
	VkF5         = WinVK(0x74)
	VkF6         = WinVK(0x75)
	VkF7         = WinVK(0x76)
	VkF8         = WinVK(0x77)
	VkF9         = WinVK(0x78)
	VkF10        = WinVK(0x79)
	VkF11        = WinVK(0x7a)
	VkF12        = WinVK(0x7b)
	VkF13        = WinVK(0x7c)
	VkF14        = WinVK(0x7d)
	VkF15        = WinVK(0x7e)
	VkF16        = WinVK(0x7f)
	VkF17        = WinVK(0x80)
	VkF18        = WinVK(0x81)
	VkF19        = WinVK(0x82)
	VkF20        = WinVK(0x83)
	VkF21        = WinVK(0x84)
	VkF22        = WinVK(0x85)
	VkF23        = WinVK(0x86)
	VkF24        = WinVK(0x87)
	VkNumLock    = WinVK(0x90)
	VkScroll     = WinVK(0x91) // scroll lock
	VkLShift     = WinVK(0xa0)
	VkRShift     = WinVK(0xa1)
	VkLControl   = WinVK(0xa2)
	VkRControl   = WinVK(0xa3)
	VkLMenu      = WinVK(0xa4) // left alt
	VkRMenu      = WinVK(0xa5) // right alt
	VkOem1       = WinVK(0xba) // ; and :
	VkOemPlus    = WinVK(0xbb) // = and +
	VkOemComma   = WinVK(0xbc) // , and <
	VkOemMinus   = WinVK(0xbd) // - and _
	VkOemPeriod  = WinVK(0xbe) // . and >
	VkOem2       = WinVK(0xbf) // / and ?
	VkOem3       = WinVK(0xc0) // ` and ~
	VkOem4       = WinVK(0xdb) // [ and {
	VkOem5       = WinVK(0xdc) // \ and |
	VkOem6       = WinVK(0xdd) // ] and }
	VkOem7       = WinVK(0xde) // ' and "
	VkOem8       = WinVK(0xdf) // right control for Canadian CSA
	VkOem102     = WinVK(0xe2) // ISO backslash
	VkPacket     = WinVK(0xe7)
)

var baseKeys map[Key]BaseKey

// KittyBase returns the corresponding Kitty "base" key for the given USB cod.
// If no corresponding value can be found, then zero is returned.  Note that
// some keys (such as F1) are valid, and recognized by Kitty, but do not use the
// base key encoding because they use another reporting format.
func (k Key) KittyBase() BaseKey {
	if w, ok := baseKeys[k]; ok {
		return w
	}
	return 0
}

func init() {
	// we place them in init to avoid incorrect processing in coverage checks.

	baseKeys = map[Key]BaseKey{
		KeyA:            'a',
		KeyB:            'b',
		KeyC:            'c',
		KeyD:            'd',
		KeyE:            'e',
		KeyF:            'f',
		KeyG:            'g',
		KeyH:            'h',
		KeyI:            'i',
		KeyJ:            'j',
		KeyK:            'k',
		KeyL:            'l',
		KeyM:            'm',
		KeyN:            'n',
		KeyO:            'o',
		KeyP:            'p',
		KeyQ:            'q',
		KeyR:            'r',
		KeyS:            's',
		KeyT:            't',
		KeyU:            'u',
		KeyV:            'v',
		KeyW:            'w',
		KeyX:            'x',
		KeyY:            'y',
		KeyZ:            'z',
		Key1:            '1',
		Key2:            '2',
		Key3:            '3',
		Key4:            '4',
		Key5:            '5',
		Key6:            '6',
		Key7:            '7',
		Key8:            '8',
		Key9:            '9',
		Key0:            '0',
		KeyEnter:        '\r',
		KeyEsc:          '\x1b',
		KeyBackspace:    '\x7f',
		KeyTab:          '\t',
		KeySpace:        ' ',
		KeyMinus:        '-',
		KeyEqual:        '=',
		KeyLBrace:       '[',
		KeyRBrace:       ']',
		KeyBackslash:    '\\',
		KeyIsoHash:      '#',
		KeySemi:         ';',
		KeyQuote:        '\'',
		KeyGrave:        '`',
		KeyComma:        ',',
		KeyPeriod:       '.',
		KeySlash:        '/',
		KeyCapsLock:     57357,
		KeyPrtScr:       57361,
		KeyScrLock:      57359,
		KeyPause:        57362,
		KeyNumLock:      57360,
		KeyPadDiv:       57410,
		KeyPadMul:       57411,
		KeyPadSub:       57412,
		KeyPadAdd:       57413,
		KeyPadEnter:     57414,
		KeyPad1:         57400,
		KeyPad2:         57401,
		KeyPad3:         57402,
		KeyPad4:         57403,
		KeyPad5:         57404,
		KeyPad6:         57405,
		KeyPad7:         57406,
		KeyPad8:         57407,
		KeyPad9:         57408,
		KeyPad0:         57399,
		KeyPadDec:       57409,
		KeyIsoBackSlash: '\\', // TODO: does kitty have another mapping for this?
		KeyMenu:         57363,
		KeyPadEqual:     57415,
		KeyF13:          57376,
		KeyF14:          57377,
		KeyF15:          57378,
		KeyF16:          57379,
		KeyF17:          57380,
		KeyF18:          57381,
		KeyF19:          57382,
		KeyF20:          57383,
		KeyF21:          57384,
		KeyF22:          57385,
		KeyF23:          57386,
		KeyF24:          57387, // NB: F25 up through 35 are notionally supported by Kitty, but not by us
		KeyPadComma:     57416,
		KeyIsoSlash:     '/', // Kitty cannot discriminate?
		KeyYen:          '¥',
		KeyLCtrl:        57442,
		KeyLShift:       57441,
		KeyLAlt:         57443,
		KeyLMeta:        57444,
		KeyRCtrl:        57448,
		KeyRShift:       57447,
		KeyRAlt:         57449,
		KeyRMeta:        57450,

		// KeyHiragana:   0, // TBD
		// KeyConvert:    0, // TBD
		// KeyNonConvert: 0, // TD

		// Windows uses a bunch of HID usages from
		// the consumer page (0x0c) for media playback, and
		// other applications. We just ignore them.
	}

	scanCodes = map[Key]ScanCode{
		KeyA:            0x1e,
		KeyB:            0x30,
		KeyC:            0x2e,
		KeyD:            0x20,
		KeyE:            0x12,
		KeyF:            0x21,
		KeyG:            0x22,
		KeyH:            0x23,
		KeyI:            0x17,
		KeyJ:            0x24,
		KeyK:            0x25,
		KeyL:            0x26,
		KeyM:            0x32,
		KeyN:            0x31,
		KeyO:            0x18,
		KeyP:            0x19,
		KeyQ:            0x10,
		KeyR:            0x13,
		KeyS:            0x1f,
		KeyT:            0x14,
		KeyU:            0x16,
		KeyV:            0x2f,
		KeyW:            0x11,
		KeyX:            0x2d,
		KeyY:            0x15,
		KeyZ:            0x2c,
		Key1:            0x02,
		Key2:            0x03,
		Key3:            0x04,
		Key4:            0x05,
		Key5:            0x06,
		Key6:            0x07,
		Key7:            0x08,
		Key8:            0x09,
		Key9:            0x0a,
		Key0:            0x0b,
		KeyEnter:        0x1c,
		KeyEsc:          0x01,
		KeyBackspace:    0x0e,
		KeyTab:          0x0f,
		KeySpace:        0x39,
		KeyMinus:        0x0c,
		KeyEqual:        0x0d,
		KeyLBrace:       0x1a,
		KeyRBrace:       0x1b,
		KeyBackslash:    0x2b,
		KeySemi:         0x27,
		KeyQuote:        0x28,
		KeyGrave:        0x29,
		KeyComma:        0x33,
		KeyPeriod:       0x34,
		KeySlash:        0x35,
		KeyCapsLock:     0x3a,
		KeyF1:           0x3b,
		KeyF2:           0x3c,
		KeyF3:           0x3d,
		KeyF4:           0x3e,
		KeyF5:           0x3f,
		KeyF6:           0x40,
		KeyF7:           0x41,
		KeyF8:           0x42,
		KeyF9:           0x43,
		KeyF10:          0x44,
		KeyF11:          0x57,
		KeyF12:          0x58,
		KeyPrtScr:       0x54,
		KeyScrLock:      0x46,
		KeyPause:        0xe046,
		KeyInsert:       0xe052,
		KeyHome:         0xe047,
		KeyPgUp:         0xe049,
		KeyDelete:       0xe053,
		KeyEnd:          0xe04f,
		KeyPgDn:         0xe051,
		KeyRight:        0xe04d,
		KeyLeft:         0xe04b,
		KeyDown:         0xe050,
		KeyUp:           0xe048,
		KeyNumLock:      0x45,
		KeyPadDiv:       0xe035,
		KeyPadMul:       0x37,
		KeyPadSub:       0x4a,
		KeyPadAdd:       0x4e,
		KeyPadEnter:     0xe01c,
		KeyPad1:         0x4f,
		KeyPad2:         0x50,
		KeyPad3:         0x51,
		KeyPad4:         0x4b,
		KeyPad5:         0x4c,
		KeyPad6:         0x4d,
		KeyPad7:         0x47,
		KeyPad8:         0x48,
		KeyPad9:         0x49,
		KeyPad0:         0x52,
		KeyPadDec:       0x53,
		KeyIsoBackSlash: 0x56,
		KeyPadEqual:     0x59,
		KeyF13:          0x64,
		KeyF14:          0x65,
		KeyF15:          0x66,
		KeyF16:          0x67,
		KeyF17:          0x68,
		KeyF18:          0x69,
		KeyF19:          0x6a,
		KeyF20:          0x6b,
		KeyF21:          0x6c,
		KeyF22:          0x6d,
		KeyF23:          0x6e,
		KeyF24:          0x76,
		KeyPadComma:     0x7e,
		KeyIsoSlash:     0x73,
		KeyHiragana:     0x70,
		KeyYen:          0x7d,
		KeyConvert:      0x79,
		KeyNonConvert:   0x7b,
		KeyLCtrl:        0x1d,
		KeyLShift:       0x2a,
		KeyLAlt:         0x38,
		KeyLMeta:        0xe05b,
		KeyRCtrl:        0xe01d,
		KeyRShift:       0x36,
		KeyRAlt:         0xe038,
		KeyRMeta:        0xe05c,
		KeyMenu:         0xe05d,

		// Windows uses a bunch of HID usages from
		// the consumer page (0x0c) for media playback, and
		// other applications. We just ignore them.
	}

	shiftedBaseKeys = map[BaseKey]rune{
		'`':  '~',
		'1':  '!',
		'2':  '@',
		'3':  '#',
		'4':  '$',
		'5':  '%',
		'6':  '^',
		'7':  '&',
		'8':  '*',
		'9':  '(',
		'0':  ')',
		'-':  '_',
		'=':  '+',
		'¥':  '|',
		'[':  '{',
		']':  '}',
		'\\': '|',
		';':  ':',
		'\'': '"',
		',':  '<',
		'.':  '>',
		'/':  '?',
	}
}
