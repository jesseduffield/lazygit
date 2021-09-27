// +build windows

// Copyright 2021 The TCell Authors
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
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

type cScreen struct {
	in         syscall.Handle
	out        syscall.Handle
	cancelflag syscall.Handle
	scandone   chan struct{}
	evch       chan Event
	quit       chan struct{}
	curx       int
	cury       int
	style      Style
	clear      bool
	fini       bool
	vten       bool
	truecolor  bool
	running    bool

	w int
	h int

	oscreen consoleInfo
	ocursor cursorInfo
	oimode  uint32
	oomode  uint32
	cells   CellBuffer

	finiOnce sync.Once

	mouseEnabled bool
	wg           sync.WaitGroup
	stopQ        chan struct{}

	sync.Mutex
}

var winLock sync.Mutex

var winPalette = []Color{
	ColorBlack,
	ColorMaroon,
	ColorGreen,
	ColorNavy,
	ColorOlive,
	ColorPurple,
	ColorTeal,
	ColorSilver,
	ColorGray,
	ColorRed,
	ColorLime,
	ColorBlue,
	ColorYellow,
	ColorFuchsia,
	ColorAqua,
	ColorWhite,
}

var winColors = map[Color]Color{
	ColorBlack:   ColorBlack,
	ColorMaroon:  ColorMaroon,
	ColorGreen:   ColorGreen,
	ColorNavy:    ColorNavy,
	ColorOlive:   ColorOlive,
	ColorPurple:  ColorPurple,
	ColorTeal:    ColorTeal,
	ColorSilver:  ColorSilver,
	ColorGray:    ColorGray,
	ColorRed:     ColorRed,
	ColorLime:    ColorLime,
	ColorBlue:    ColorBlue,
	ColorYellow:  ColorYellow,
	ColorFuchsia: ColorFuchsia,
	ColorAqua:    ColorAqua,
	ColorWhite:   ColorWhite,
}

var (
	k32 = syscall.NewLazyDLL("kernel32.dll")
	u32 = syscall.NewLazyDLL("user32.dll")
)

// We have to bring in the kernel32 and user32 DLLs directly, so we can get
// access to some system calls that the core Go API lacks.
//
// Note that Windows appends some functions with W to indicate that wide
// characters (Unicode) are in use.  The documentation refers to them
// without this suffix, as the resolution is made via preprocessor.
var (
	procReadConsoleInput           = k32.NewProc("ReadConsoleInputW")
	procWaitForMultipleObjects     = k32.NewProc("WaitForMultipleObjects")
	procCreateEvent                = k32.NewProc("CreateEventW")
	procSetEvent                   = k32.NewProc("SetEvent")
	procGetConsoleCursorInfo       = k32.NewProc("GetConsoleCursorInfo")
	procSetConsoleCursorInfo       = k32.NewProc("SetConsoleCursorInfo")
	procSetConsoleCursorPosition   = k32.NewProc("SetConsoleCursorPosition")
	procSetConsoleMode             = k32.NewProc("SetConsoleMode")
	procGetConsoleMode             = k32.NewProc("GetConsoleMode")
	procGetConsoleScreenBufferInfo = k32.NewProc("GetConsoleScreenBufferInfo")
	procFillConsoleOutputAttribute = k32.NewProc("FillConsoleOutputAttribute")
	procFillConsoleOutputCharacter = k32.NewProc("FillConsoleOutputCharacterW")
	procSetConsoleWindowInfo       = k32.NewProc("SetConsoleWindowInfo")
	procSetConsoleScreenBufferSize = k32.NewProc("SetConsoleScreenBufferSize")
	procSetConsoleTextAttribute    = k32.NewProc("SetConsoleTextAttribute")
	procMessageBeep                = u32.NewProc("MessageBeep")
)

const (
	w32Infinite    = ^uintptr(0)
	w32WaitObject0 = uintptr(0)
)

const (
	// VT100/XTerm escapes understood by the console
	vtShowCursor = "\x1b[?25h"
	vtHideCursor = "\x1b[?25l"
	vtCursorPos  = "\x1b[%d;%dH" // Note that it is Y then X
	vtSgr0       = "\x1b[0m"
	vtBold       = "\x1b[1m"
	vtUnderline  = "\x1b[4m"
	vtBlink      = "\x1b[5m" // Not sure this is processed
	vtReverse    = "\x1b[7m"
	vtSetFg      = "\x1b[38;5;%dm"
	vtSetBg      = "\x1b[48;5;%dm"
	vtSetFgRGB   = "\x1b[38;2;%d;%d;%dm" // RGB
	vtSetBgRGB   = "\x1b[48;2;%d;%d;%dm" // RGB
)

// NewConsoleScreen returns a Screen for the Windows console associated
// with the current process.  The Screen makes use of the Windows Console
// API to display content and read events.
func NewConsoleScreen() (Screen, error) {
	return &cScreen{}, nil
}

func (s *cScreen) Init() error {
	s.evch = make(chan Event, 10)
	s.quit = make(chan struct{})
	s.scandone = make(chan struct{})

	in, e := syscall.Open("CONIN$", syscall.O_RDWR, 0)
	if e != nil {
		return e
	}
	s.in = in
	out, e := syscall.Open("CONOUT$", syscall.O_RDWR, 0)
	if e != nil {
		syscall.Close(s.in)
		return e
	}
	s.out = out

	s.truecolor = true

	// ConEmu handling of colors and scrolling when in terminal
	// mode is extremely problematic at the best.  The color
	// palette will scroll even though characters do not, when
	// emitting stuff for the last character.  In the future we
	// might change this to look at specific versions of ConEmu
	// if they fix the bug.
	if os.Getenv("ConEmuPID") != "" {
		s.truecolor = false
	}
	switch os.Getenv("TCELL_TRUECOLOR") {
	case "disable":
		s.truecolor = false
	case "enable":
		s.truecolor = true
	}

	s.Lock()

	s.curx = -1
	s.cury = -1
	s.style = StyleDefault
	s.getCursorInfo(&s.ocursor)
	s.getConsoleInfo(&s.oscreen)
	s.getOutMode(&s.oomode)
	s.getInMode(&s.oimode)
	s.resize()

	s.fini = false
	s.setInMode(modeResizeEn | modeExtndFlg)

	// 24-bit color is opt-in for now, because we can't figure out
	// to make it work consistently.
	if s.truecolor {
		s.setOutMode(modeVtOutput | modeNoAutoNL | modeCookedOut)
		var omode uint32
		s.getOutMode(&omode)
		if omode&modeVtOutput == modeVtOutput {
			s.vten = true
		} else {
			s.truecolor = false
			s.setOutMode(0)
		}
	} else {
		s.setOutMode(0)
	}

	s.Unlock()

	return s.engage()
}

func (s *cScreen) CharacterSet() string {
	// We are always UTF-16LE on Windows
	return "UTF-16LE"
}

func (s *cScreen) EnableMouse(...MouseFlags) {
	s.Lock()
	s.mouseEnabled = true
	s.enableMouse(true)
	s.Unlock()
}

func (s *cScreen) DisableMouse() {
	s.Lock()
	s.mouseEnabled = false
	s.enableMouse(false)
	s.Unlock()
}

func (s *cScreen) enableMouse(on bool) {
	if on {
		s.setInMode(modeResizeEn | modeMouseEn | modeExtndFlg)
	} else {
		s.setInMode(modeResizeEn | modeExtndFlg)
	}
}

// Windows lacks bracketed paste (for now)

func (s *cScreen) EnablePaste() {}

func (s *cScreen) DisablePaste() {}

func (s *cScreen) Fini() {
	s.disengage()
}

func (s *cScreen) disengage() {
	s.Lock()
	if !s.running {
		s.Unlock()
		return
	}
	s.running = false
	stopQ := s.stopQ
	procSetEvent.Call(uintptr(s.cancelflag))
	close(stopQ)
	s.Unlock()

	s.wg.Wait()

	s.setInMode(s.oimode)
	s.setOutMode(s.oomode)
	s.setBufferSize(int(s.oscreen.size.x), int(s.oscreen.size.y))
	s.clearScreen(StyleDefault, false)
	s.setCursorPos(0, 0, false)
	s.setCursorInfo(&s.ocursor)
	procSetConsoleTextAttribute.Call(
		uintptr(s.out),
		uintptr(s.mapStyle(StyleDefault)))
}

func (s *cScreen) engage() error {
	s.Lock()
	defer s.Unlock()
	if s.running {
		return errors.New("already engaged")
	}
	s.stopQ = make(chan struct{})
	cf, _, e := procCreateEvent.Call(
		uintptr(0),
		uintptr(1),
		uintptr(0),
		uintptr(0))
	if cf == uintptr(0) {
		return e
	}
	s.running = true
	s.cancelflag = syscall.Handle(cf)
	s.enableMouse(s.mouseEnabled)

	if s.vten {
		s.setOutMode(modeVtOutput | modeNoAutoNL | modeCookedOut)
	} else {
		s.setOutMode(0)
	}

	s.clearScreen(s.style, s.vten)
	s.hideCursor()

	s.cells.Invalidate()
	s.hideCursor()
	s.resize()
	s.draw()
	s.doCursor()

	s.wg.Add(1)
	go s.scanInput(s.stopQ)
	return nil
}

func (s *cScreen) PostEventWait(ev Event) {
	s.evch <- ev
}

func (s *cScreen) PostEvent(ev Event) error {
	select {
	case s.evch <- ev:
		return nil
	default:
		return ErrEventQFull
	}
}

func (s *cScreen) ChannelEvents(ch chan<- Event, quit <-chan struct{}) {
	defer close(ch)
	for {
		select {
		case <-quit:
			return
		case <-s.stopQ:
			return
		case ev := <-s.evch:
			select {
			case <-quit:
				return
			case <-s.stopQ:
				return
			case ch <- ev:
			}
		}
	}
}

func (s *cScreen) PollEvent() Event {
	select {
	case <-s.stopQ:
		return nil
	case ev := <-s.evch:
		return ev
	}
}

func (s *cScreen) HasPendingEvent() bool {
	return len(s.evch) > 0
}

type cursorInfo struct {
	size    uint32
	visible uint32
}

type coord struct {
	x int16
	y int16
}

func (c coord) uintptr() uintptr {
	// little endian, put x first
	return uintptr(c.x) | (uintptr(c.y) << 16)
}

type rect struct {
	left   int16
	top    int16
	right  int16
	bottom int16
}

func (s *cScreen) emitVtString(vs string) {
	esc := utf16.Encode([]rune(vs))
	syscall.WriteConsole(s.out, &esc[0], uint32(len(esc)), nil, nil)
}

func (s *cScreen) showCursor() {
	if s.vten {
		s.emitVtString(vtShowCursor)
	} else {
		s.setCursorInfo(&cursorInfo{size: 100, visible: 1})
	}
}

func (s *cScreen) hideCursor() {
	if s.vten {
		s.emitVtString(vtHideCursor)
	} else {
		s.setCursorInfo(&cursorInfo{size: 1, visible: 0})
	}
}

func (s *cScreen) ShowCursor(x, y int) {
	s.Lock()
	if !s.fini {
		s.curx = x
		s.cury = y
	}
	s.doCursor()
	s.Unlock()
}

func (s *cScreen) doCursor() {
	x, y := s.curx, s.cury

	if x < 0 || y < 0 || x >= s.w || y >= s.h {
		s.hideCursor()
	} else {
		s.setCursorPos(x, y, s.vten)
		s.showCursor()
	}
}

func (s *cScreen) HideCursor() {
	s.ShowCursor(-1, -1)
}

type inputRecord struct {
	typ  uint16
	_    uint16
	data [16]byte
}

const (
	keyEvent    uint16 = 1
	mouseEvent  uint16 = 2
	resizeEvent uint16 = 4
	menuEvent   uint16 = 8  // don't use
	focusEvent  uint16 = 16 // don't use
)

type mouseRecord struct {
	x     int16
	y     int16
	btns  uint32
	mod   uint32
	flags uint32
}

const (
	mouseDoubleClick uint32 = 0x2
	mouseHWheeled    uint32 = 0x8
	mouseVWheeled    uint32 = 0x4
	mouseMoved       uint32 = 0x1
)

type resizeRecord struct {
	x int16
	y int16
}

type keyRecord struct {
	isdown int32
	repeat uint16
	kcode  uint16
	scode  uint16
	ch     uint16
	mod    uint32
}

const (
	// Constants per Microsoft.  We don't put the modifiers
	// here.
	vkCancel = 0x03
	vkBack   = 0x08 // Backspace
	vkTab    = 0x09
	vkClear  = 0x0c
	vkReturn = 0x0d
	vkPause  = 0x13
	vkEscape = 0x1b
	vkSpace  = 0x20
	vkPrior  = 0x21 // PgUp
	vkNext   = 0x22 // PgDn
	vkEnd    = 0x23
	vkHome   = 0x24
	vkLeft   = 0x25
	vkUp     = 0x26
	vkRight  = 0x27
	vkDown   = 0x28
	vkPrint  = 0x2a
	vkPrtScr = 0x2c
	vkInsert = 0x2d
	vkDelete = 0x2e
	vkHelp   = 0x2f
	vkF1     = 0x70
	vkF2     = 0x71
	vkF3     = 0x72
	vkF4     = 0x73
	vkF5     = 0x74
	vkF6     = 0x75
	vkF7     = 0x76
	vkF8     = 0x77
	vkF9     = 0x78
	vkF10    = 0x79
	vkF11    = 0x7a
	vkF12    = 0x7b
	vkF13    = 0x7c
	vkF14    = 0x7d
	vkF15    = 0x7e
	vkF16    = 0x7f
	vkF17    = 0x80
	vkF18    = 0x81
	vkF19    = 0x82
	vkF20    = 0x83
	vkF21    = 0x84
	vkF22    = 0x85
	vkF23    = 0x86
	vkF24    = 0x87
)

var vkKeys = map[uint16]Key{
	vkCancel: KeyCancel,
	vkBack:   KeyBackspace,
	vkTab:    KeyTab,
	vkClear:  KeyClear,
	vkPause:  KeyPause,
	vkPrint:  KeyPrint,
	vkPrtScr: KeyPrint,
	vkPrior:  KeyPgUp,
	vkNext:   KeyPgDn,
	vkReturn: KeyEnter,
	vkEnd:    KeyEnd,
	vkHome:   KeyHome,
	vkLeft:   KeyLeft,
	vkUp:     KeyUp,
	vkRight:  KeyRight,
	vkDown:   KeyDown,
	vkInsert: KeyInsert,
	vkDelete: KeyDelete,
	vkHelp:   KeyHelp,
	vkF1:     KeyF1,
	vkF2:     KeyF2,
	vkF3:     KeyF3,
	vkF4:     KeyF4,
	vkF5:     KeyF5,
	vkF6:     KeyF6,
	vkF7:     KeyF7,
	vkF8:     KeyF8,
	vkF9:     KeyF9,
	vkF10:    KeyF10,
	vkF11:    KeyF11,
	vkF12:    KeyF12,
	vkF13:    KeyF13,
	vkF14:    KeyF14,
	vkF15:    KeyF15,
	vkF16:    KeyF16,
	vkF17:    KeyF17,
	vkF18:    KeyF18,
	vkF19:    KeyF19,
	vkF20:    KeyF20,
	vkF21:    KeyF21,
	vkF22:    KeyF22,
	vkF23:    KeyF23,
	vkF24:    KeyF24,
}

// NB: All Windows platforms are little endian.  We assume this
// never, ever change.  The following code is endian safe. and does
// not use unsafe pointers.
func getu32(v []byte) uint32 {
	return uint32(v[0]) + (uint32(v[1]) << 8) + (uint32(v[2]) << 16) + (uint32(v[3]) << 24)
}
func geti32(v []byte) int32 {
	return int32(getu32(v))
}
func getu16(v []byte) uint16 {
	return uint16(v[0]) + (uint16(v[1]) << 8)
}
func geti16(v []byte) int16 {
	return int16(getu16(v))
}

// Convert windows dwControlKeyState to modifier mask
func mod2mask(cks uint32) ModMask {
	mm := ModNone
	// Left or right control
	if (cks & (0x0008 | 0x0004)) != 0 {
		mm |= ModCtrl
	}
	// Left or right alt
	if (cks & (0x0002 | 0x0001)) != 0 {
		mm |= ModAlt
	}
	// Any shift
	if (cks & 0x0010) != 0 {
		mm |= ModShift
	}
	return mm
}

func mrec2btns(mbtns, flags uint32) ButtonMask {
	btns := ButtonNone
	if mbtns&0x1 != 0 {
		btns |= Button1
	}
	if mbtns&0x2 != 0 {
		btns |= Button2
	}
	if mbtns&0x4 != 0 {
		btns |= Button3
	}
	if mbtns&0x8 != 0 {
		btns |= Button4
	}
	if mbtns&0x10 != 0 {
		btns |= Button5
	}
	if mbtns&0x20 != 0 {
		btns |= Button6
	}
	if mbtns&0x40 != 0 {
		btns |= Button7
	}
	if mbtns&0x80 != 0 {
		btns |= Button8
	}

	if flags&mouseVWheeled != 0 {
		if mbtns&0x80000000 == 0 {
			btns |= WheelUp
		} else {
			btns |= WheelDown
		}
	}
	if flags&mouseHWheeled != 0 {
		if mbtns&0x80000000 == 0 {
			btns |= WheelRight
		} else {
			btns |= WheelLeft
		}
	}
	return btns
}

func (s *cScreen) getConsoleInput() error {
	// cancelFlag comes first as WaitForMultipleObjects returns the lowest index
	// in the event that both events are signalled.
	waitObjects := []syscall.Handle{s.cancelflag, s.in}
	// As arrays are contiguous in memory, a pointer to the first object is the
	// same as a pointer to the array itself.
	pWaitObjects := unsafe.Pointer(&waitObjects[0])

	rv, _, er := procWaitForMultipleObjects.Call(
		uintptr(len(waitObjects)),
		uintptr(pWaitObjects),
		uintptr(0),
		w32Infinite)
	// WaitForMultipleObjects returns WAIT_OBJECT_0 + the index.
	switch rv {
	case w32WaitObject0: // s.cancelFlag
		return errors.New("cancelled")
	case w32WaitObject0 + 1: // s.in
		rec := &inputRecord{}
		var nrec int32
		rv, _, er := procReadConsoleInput.Call(
			uintptr(s.in),
			uintptr(unsafe.Pointer(rec)),
			uintptr(1),
			uintptr(unsafe.Pointer(&nrec)))
		if rv == 0 {
			return er
		}
		if nrec != 1 {
			return nil
		}
		switch rec.typ {
		case keyEvent:
			krec := &keyRecord{}
			krec.isdown = geti32(rec.data[0:])
			krec.repeat = getu16(rec.data[4:])
			krec.kcode = getu16(rec.data[6:])
			krec.scode = getu16(rec.data[8:])
			krec.ch = getu16(rec.data[10:])
			krec.mod = getu32(rec.data[12:])

			if krec.isdown == 0 || krec.repeat < 1 {
				// its a key release event, ignore it
				return nil
			}
			if krec.ch != 0 {
				// synthesized key code
				for krec.repeat > 0 {
					// convert shift+tab to backtab
					if mod2mask(krec.mod) == ModShift && krec.ch == vkTab {
						s.PostEventWait(NewEventKey(KeyBacktab, 0,
							ModNone))
					} else {
						s.PostEventWait(NewEventKey(KeyRune, rune(krec.ch),
							mod2mask(krec.mod)))
					}
					krec.repeat--
				}
				return nil
			}
			key := KeyNUL // impossible on Windows
			ok := false
			if key, ok = vkKeys[krec.kcode]; !ok {
				return nil
			}
			for krec.repeat > 0 {
				s.PostEventWait(NewEventKey(key, rune(krec.ch),
					mod2mask(krec.mod)))
				krec.repeat--
			}

		case mouseEvent:
			var mrec mouseRecord
			mrec.x = geti16(rec.data[0:])
			mrec.y = geti16(rec.data[2:])
			mrec.btns = getu32(rec.data[4:])
			mrec.mod = getu32(rec.data[8:])
			mrec.flags = getu32(rec.data[12:])
			btns := mrec2btns(mrec.btns, mrec.flags)
			// we ignore double click, events are delivered normally
			s.PostEventWait(NewEventMouse(int(mrec.x), int(mrec.y), btns,
				mod2mask(mrec.mod)))

		case resizeEvent:
			var rrec resizeRecord
			rrec.x = geti16(rec.data[0:])
			rrec.y = geti16(rec.data[2:])
			s.PostEventWait(NewEventResize(int(rrec.x), int(rrec.y)))

		default:
		}
	default:
		return er
	}

	return nil
}

func (s *cScreen) scanInput(stopQ chan struct{}) {
	defer s.wg.Done()
	for {
		select {
		case <-stopQ:
			return
		default:
		}
		if e := s.getConsoleInput(); e != nil {
			return
		}
	}
}

// Windows console can display 8 characters, in either low or high intensity
func (s *cScreen) Colors() int {
	if s.vten {
		return 1 << 24
	}
	return 16
}

var vgaColors = map[Color]uint16{
	ColorBlack:   0,
	ColorMaroon:  0x4,
	ColorGreen:   0x2,
	ColorNavy:    0x1,
	ColorOlive:   0x6,
	ColorPurple:  0x5,
	ColorTeal:    0x3,
	ColorSilver:  0x7,
	ColorGrey:    0x8,
	ColorRed:     0xc,
	ColorLime:    0xa,
	ColorBlue:    0x9,
	ColorYellow:  0xe,
	ColorFuchsia: 0xd,
	ColorAqua:    0xb,
	ColorWhite:   0xf,
}

// Windows uses RGB signals
func mapColor2RGB(c Color) uint16 {
	winLock.Lock()
	if v, ok := winColors[c]; ok {
		c = v
	} else {
		v = FindColor(c, winPalette)
		winColors[c] = v
		c = v
	}
	winLock.Unlock()

	if vc, ok := vgaColors[c]; ok {
		return vc
	}
	return 0
}

// Map a tcell style to Windows attributes
func (s *cScreen) mapStyle(style Style) uint16 {
	f, b, a := style.Decompose()
	fa := s.oscreen.attrs & 0xf
	ba := (s.oscreen.attrs) >> 4 & 0xf
	if f != ColorDefault && f != ColorReset {
		fa = mapColor2RGB(f)
	}
	if b != ColorDefault && b != ColorReset {
		ba = mapColor2RGB(b)
	}
	var attr uint16
	// We simulate reverse by doing the color swap ourselves.
	// Apparently windows cannot really do this except in DBCS
	// views.
	if a&AttrReverse != 0 {
		attr = ba
		attr |= (fa << 4)
	} else {
		attr = fa
		attr |= (ba << 4)
	}
	if a&AttrBold != 0 {
		attr |= 0x8
	}
	if a&AttrDim != 0 {
		attr &^= 0x8
	}
	if a&AttrUnderline != 0 {
		// Best effort -- doesn't seem to work though.
		attr |= 0x8000
	}
	// Blink is unsupported
	return attr
}

func (s *cScreen) SetCell(x, y int, style Style, ch ...rune) {
	if len(ch) > 0 {
		s.SetContent(x, y, ch[0], ch[1:], style)
	} else {
		s.SetContent(x, y, ' ', nil, style)
	}
}

func (s *cScreen) SetContent(x, y int, mainc rune, combc []rune, style Style) {
	s.Lock()
	if !s.fini {
		s.cells.SetContent(x, y, mainc, combc, style)
	}
	s.Unlock()
}

func (s *cScreen) GetContent(x, y int) (rune, []rune, Style, int) {
	s.Lock()
	mainc, combc, style, width := s.cells.GetContent(x, y)
	s.Unlock()
	return mainc, combc, style, width
}

func (s *cScreen) sendVtStyle(style Style) {
	esc := &strings.Builder{}

	fg, bg, attrs := style.Decompose()

	esc.WriteString(vtSgr0)

	if attrs&(AttrBold|AttrDim) == AttrBold {
		esc.WriteString(vtBold)
	}
	if attrs&AttrBlink != 0 {
		esc.WriteString(vtBlink)
	}
	if attrs&AttrUnderline != 0 {
		esc.WriteString(vtUnderline)
	}
	if attrs&AttrReverse != 0 {
		esc.WriteString(vtReverse)
	}
	if fg.IsRGB() {
		r, g, b := fg.RGB()
		fmt.Fprintf(esc, vtSetFgRGB, r, g, b)
	} else if fg.Valid() {
		fmt.Fprintf(esc, vtSetFg, fg&0xff)
	}
	if bg.IsRGB() {
		r, g, b := bg.RGB()
		fmt.Fprintf(esc, vtSetBgRGB, r, g, b)
	} else if bg.Valid() {
		fmt.Fprintf(esc, vtSetBg, bg&0xff)
	}
	s.emitVtString(esc.String())
}

func (s *cScreen) writeString(x, y int, style Style, ch []uint16) {
	// we assume the caller has hidden the cursor
	if len(ch) == 0 {
		return
	}
	s.setCursorPos(x, y, s.vten)

	if s.vten {
		s.sendVtStyle(style)
	} else {
		procSetConsoleTextAttribute.Call(
			uintptr(s.out),
			uintptr(s.mapStyle(style)))
	}
	syscall.WriteConsole(s.out, &ch[0], uint32(len(ch)), nil, nil)
}

func (s *cScreen) draw() {
	// allocate a scratch line bit enough for no combining chars.
	// if you have combining characters, you may pay for extra allocs.
	if s.clear {
		s.clearScreen(s.style, s.vten)
		s.clear = false
		s.cells.Invalidate()
	}
	buf := make([]uint16, 0, s.w)
	wcs := buf[:]
	lstyle := styleInvalid

	lx, ly := -1, -1
	ra := make([]rune, 1)

	for y := 0; y < s.h; y++ {
		for x := 0; x < s.w; x++ {
			mainc, combc, style, width := s.cells.GetContent(x, y)
			dirty := s.cells.Dirty(x, y)
			if style == StyleDefault {
				style = s.style
			}

			if !dirty || style != lstyle {
				// write out any data queued thus far
				// because we are going to skip over some
				// cells, or because we need to change styles
				s.writeString(lx, ly, lstyle, wcs)
				wcs = buf[0:0]
				lstyle = StyleDefault
				if !dirty {
					continue
				}
			}
			if x > s.w-width {
				mainc = ' '
				combc = nil
				width = 1
			}
			if len(wcs) == 0 {
				lstyle = style
				lx = x
				ly = y
			}
			ra[0] = mainc
			wcs = append(wcs, utf16.Encode(ra)...)
			if len(combc) != 0 {
				wcs = append(wcs, utf16.Encode(combc)...)
			}
			for dx := 0; dx < width; dx++ {
				s.cells.SetDirty(x+dx, y, false)
			}
			x += width - 1
		}
		s.writeString(lx, ly, lstyle, wcs)
		wcs = buf[0:0]
		lstyle = styleInvalid
	}
}

func (s *cScreen) Show() {
	s.Lock()
	if !s.fini {
		s.hideCursor()
		s.resize()
		s.draw()
		s.doCursor()
	}
	s.Unlock()
}

func (s *cScreen) Sync() {
	s.Lock()
	if !s.fini {
		s.cells.Invalidate()
		s.hideCursor()
		s.resize()
		s.draw()
		s.doCursor()
	}
	s.Unlock()
}

type consoleInfo struct {
	size  coord
	pos   coord
	attrs uint16
	win   rect
	maxsz coord
}

func (s *cScreen) getConsoleInfo(info *consoleInfo) {
	procGetConsoleScreenBufferInfo.Call(
		uintptr(s.out),
		uintptr(unsafe.Pointer(info)))
}

func (s *cScreen) getCursorInfo(info *cursorInfo) {
	procGetConsoleCursorInfo.Call(
		uintptr(s.out),
		uintptr(unsafe.Pointer(info)))
}

func (s *cScreen) setCursorInfo(info *cursorInfo) {
	procSetConsoleCursorInfo.Call(
		uintptr(s.out),
		uintptr(unsafe.Pointer(info)))

}

func (s *cScreen) setCursorPos(x, y int, vtEnable bool) {
	if vtEnable {
		// Note that the string is Y first.  Origin is 1,1.
		s.emitVtString(fmt.Sprintf(vtCursorPos, y+1, x+1))
	} else {
		procSetConsoleCursorPosition.Call(
			uintptr(s.out),
			coord{int16(x), int16(y)}.uintptr())
	}
}

func (s *cScreen) setBufferSize(x, y int) {
	procSetConsoleScreenBufferSize.Call(
		uintptr(s.out),
		coord{int16(x), int16(y)}.uintptr())
}

func (s *cScreen) Size() (int, int) {
	s.Lock()
	w, h := s.w, s.h
	s.Unlock()

	return w, h
}

func (s *cScreen) resize() {
	info := consoleInfo{}
	s.getConsoleInfo(&info)

	w := int((info.win.right - info.win.left) + 1)
	h := int((info.win.bottom - info.win.top) + 1)

	if s.w == w && s.h == h {
		return
	}

	s.cells.Resize(w, h)
	s.w = w
	s.h = h

	s.setBufferSize(w, h)

	r := rect{0, 0, int16(w - 1), int16(h - 1)}
	procSetConsoleWindowInfo.Call(
		uintptr(s.out),
		uintptr(1),
		uintptr(unsafe.Pointer(&r)))
	s.PostEvent(NewEventResize(w, h))
}

func (s *cScreen) Clear() {
	s.Fill(' ', s.style)
}

func (s *cScreen) Fill(r rune, style Style) {
	s.Lock()
	if !s.fini {
		s.cells.Fill(r, style)
		s.clear = true
	}
	s.Unlock()
}

func (s *cScreen) clearScreen(style Style, vtEnable bool) {
	if vtEnable {
		s.sendVtStyle(style)
		row := strings.Repeat(" ", s.w)
		for y := 0; y < s.h; y++ {
			s.setCursorPos(0, y, vtEnable)
			s.emitVtString(row)
		}
		s.setCursorPos(0, 0, vtEnable)

	} else {
		pos := coord{0, 0}
		attr := s.mapStyle(style)
		x, y := s.w, s.h
		scratch := uint32(0)
		count := uint32(x * y)

		procFillConsoleOutputAttribute.Call(
			uintptr(s.out),
			uintptr(attr),
			uintptr(count),
			pos.uintptr(),
			uintptr(unsafe.Pointer(&scratch)))
		procFillConsoleOutputCharacter.Call(
			uintptr(s.out),
			uintptr(' '),
			uintptr(count),
			pos.uintptr(),
			uintptr(unsafe.Pointer(&scratch)))
	}
}

const (
	// Input modes
	modeExtndFlg uint32 = 0x0080
	modeMouseEn         = 0x0010
	modeResizeEn        = 0x0008
	modeCooked          = 0x0001
	modeVtInput         = 0x0200

	// Output modes
	modeCookedOut uint32 = 0x0001
	modeWrapEOL          = 0x0002
	modeVtOutput         = 0x0004
	modeNoAutoNL         = 0x0008
)

func (s *cScreen) setInMode(mode uint32) error {
	rv, _, err := procSetConsoleMode.Call(
		uintptr(s.in),
		uintptr(mode))
	if rv == 0 {
		return err
	}
	return nil
}

func (s *cScreen) setOutMode(mode uint32) error {
	rv, _, err := procSetConsoleMode.Call(
		uintptr(s.out),
		uintptr(mode))
	if rv == 0 {
		return err
	}
	return nil
}

func (s *cScreen) getInMode(v *uint32) {
	procGetConsoleMode.Call(
		uintptr(s.in),
		uintptr(unsafe.Pointer(v)))
}

func (s *cScreen) getOutMode(v *uint32) {
	procGetConsoleMode.Call(
		uintptr(s.out),
		uintptr(unsafe.Pointer(v)))
}

func (s *cScreen) SetStyle(style Style) {
	s.Lock()
	s.style = style
	s.Unlock()
}

// No fallback rune support, since we have Unicode.  Yay!

func (s *cScreen) RegisterRuneFallback(r rune, subst string) {
}

func (s *cScreen) UnregisterRuneFallback(r rune) {
}

func (s *cScreen) CanDisplay(r rune, checkFallbacks bool) bool {
	// We presume we can display anything -- we're Unicode.
	// (Sadly this not precisely true.  Combinings are especially
	// poorly supported under Windows.)
	return true
}

func (s *cScreen) HasMouse() bool {
	return true
}

func (s *cScreen) Resize(int, int, int, int) {}

func (s *cScreen) HasKey(k Key) bool {
	// Microsoft has codes for some keys, but they are unusual,
	// so we don't include them.  We include all the typical
	// 101, 105 key layout keys.
	valid := map[Key]bool{
		KeyBackspace: true,
		KeyTab:       true,
		KeyEscape:    true,
		KeyPause:     true,
		KeyPrint:     true,
		KeyPgUp:      true,
		KeyPgDn:      true,
		KeyEnter:     true,
		KeyEnd:       true,
		KeyHome:      true,
		KeyLeft:      true,
		KeyUp:        true,
		KeyRight:     true,
		KeyDown:      true,
		KeyInsert:    true,
		KeyDelete:    true,
		KeyF1:        true,
		KeyF2:        true,
		KeyF3:        true,
		KeyF4:        true,
		KeyF5:        true,
		KeyF6:        true,
		KeyF7:        true,
		KeyF8:        true,
		KeyF9:        true,
		KeyF10:       true,
		KeyF11:       true,
		KeyF12:       true,
		KeyRune:      true,
	}

	return valid[k]
}

func (s *cScreen) Beep() error {
	// A simple beep. If the sound card is not available, the sound is generated
	// using the speaker.
	//
	// Reference:
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messagebeep
	const simpleBeep = 0xffffffff
	if rv, _, err := procMessageBeep.Call(simpleBeep); rv == 0 {
		return err
	}
	return nil
}

func (s *cScreen) Suspend() error {
	s.disengage()
	return nil
}

func (s *cScreen) Resume() error {
	return s.engage()
}
