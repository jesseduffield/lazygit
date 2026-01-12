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

//go:build !(js && wasm)
// +build !js !wasm

package tcell

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"maps"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"golang.org/x/term"
	"golang.org/x/text/transform"

	"github.com/gdamore/tcell/v2/terminfo"
)

// NewTerminfoScreen returns a Screen that uses the stock TTY interface
// and POSIX terminal control, combined with a terminfo description taken from
// the $TERM environment variable.  It returns an error if the terminal
// is not supported for any reason.
//
// For terminals that do not support dynamic resize events, the $LINES
// $COLUMNS environment variables can be set to the actual window size,
// otherwise defaults taken from the terminal database are used.
func NewTerminfoScreen() (Screen, error) {
	return NewTerminfoScreenFromTty(nil)
}

// LookupTerminfo attempts to find a definition for the named $TERM falling
// back to attempting to parse the output from infocmp.
func LookupTerminfo(name string) (ti *terminfo.Terminfo, e error) {
	ti, e = terminfo.LookupTerminfo(name)
	if e != nil {
		ti, e = loadDynamicTerminfo(name)
		if e != nil {
			return nil, e
		}
		terminfo.AddTerminfo(ti)
	}

	return
}

var defaultTerm string

// NewTerminfoScreenFromTtyTerminfo returns a Screen using a custom Tty
// implementation  and custom terminfo specification.
// If the passed in tty is nil, then a reasonable default (typically /dev/tty)
// is presumed, at least on UNIX hosts. (Windows hosts will typically fail this
// call altogether.)
// If passed terminfo is nil, then TERM environment variable is queried for
// terminal specification.
func NewTerminfoScreenFromTtyTerminfo(tty Tty, ti *terminfo.Terminfo) (s Screen, e error) {
	term := defaultTerm
	if term == "" {
		term = os.Getenv("TERM")
	}
	if ti == nil {
		ti, e = LookupTerminfo(term)
		if e != nil {
			return nil, e
		}
	}

	t := &tScreen{ti: ti, tty: tty}

	if len(ti.Mouse) > 0 {
		t.mouse = []byte(ti.Mouse)
	}
	t.prepareKeys()
	t.buildAcsMap()
	t.resizeQ = make(chan bool, 1)
	t.fallback = make(map[rune]string)
	maps.Copy(t.fallback, RuneFallbacks)

	return &baseScreen{screenImpl: t}, nil
}

// NewTerminfoScreenFromTty returns a Screen using a custom Tty implementation.
// If the passed in tty is nil, then a reasonable default (typically /dev/tty)
// is presumed, at least on UNIX hosts. (Windows hosts will typically fail this
// call altogether.)
func NewTerminfoScreenFromTty(tty Tty) (Screen, error) {
	return NewTerminfoScreenFromTtyTerminfo(tty, nil)
}

// tKeyCode represents a combination of a key code and modifiers.
type tKeyCode struct {
	key Key
	mod ModMask
}

// tScreen represents a screen backed by a terminfo implementation.
type tScreen struct {
	ti             *terminfo.Terminfo
	tty            Tty
	h              int
	w              int
	fini           bool
	cells          CellBuffer
	buffering      bool // true if we are collecting writes to buf instead of sending directly to out
	buf            bytes.Buffer
	curstyle       Style
	style          Style
	resizeQ        chan bool
	quit           chan struct{}
	keychan        chan []byte
	cx             int
	cy             int
	mouse          []byte
	clear          bool
	cursorx        int
	cursory        int
	acs            map[rune]string
	charset        string
	encoder        transform.Transformer
	decoder        transform.Transformer
	fallback       map[rune]string
	colors         map[Color]Color
	palette        []Color
	truecolor      bool
	escaped        bool
	buttondn       bool
	finiOnce       sync.Once
	enablePaste    string
	disablePaste   string
	enterUrl       string
	exitUrl        string
	setWinSize     string
	enableFocus    string
	disableFocus   string
	doubleUnder    string
	curlyUnder     string
	dottedUnder    string
	dashedUnder    string
	underColor     string
	underRGB       string
	underFg        string // reset underline color to foreground
	cursorStyles   map[CursorStyle]string
	cursorStyle    CursorStyle
	cursorColor    Color
	cursorRGB      string
	cursorFg       string
	saved          *term.State
	stopQ          chan struct{}
	eventQ         chan Event
	running        bool
	wg             sync.WaitGroup
	mouseFlags     MouseFlags
	pasteEnabled   bool
	focusEnabled   bool
	setTitle       string
	saveTitle      string
	restoreTitle   string
	title          string
	setClipboard   string
	startSyncOut   string
	endSyncOut     string
	enableCsiU     string
	disableCsiU    string
	disableEmojiWA bool // if true don't try to workaround emoji bugs
	input          InputProcessor

	sync.Mutex
}

func (t *tScreen) Init() error {
	if e := t.initialize(); e != nil {
		return e
	}

	t.keychan = make(chan []byte, 10)

	t.charset = getCharset()
	if enc := GetEncoding(t.charset); enc != nil {
		t.encoder = enc.NewEncoder()
		t.decoder = enc.NewDecoder()
	} else {
		return ErrNoCharset
	}
	ti := t.ti

	// environment overrides
	w := ti.Columns
	h := ti.Lines
	if i, _ := strconv.Atoi(os.Getenv("LINES")); i != 0 {
		h = i
	}
	if i, _ := strconv.Atoi(os.Getenv("COLUMNS")); i != 0 {
		w = i
	}
	if t.ti.SetFgBgRGB != "" || t.ti.SetFgRGB != "" || t.ti.SetBgRGB != "" {
		t.truecolor = true
	}
	// A user who wants to have his themes honored can
	// set this environment variable.
	if os.Getenv("TCELL_TRUECOLOR") == "disable" {
		t.truecolor = false
	}
	// clip to reasonable limits
	nColors := min(t.nColors(), 256)
	t.colors = make(map[Color]Color, nColors)
	t.palette = make([]Color, nColors)
	for i := range nColors {
		t.palette[i] = Color(i) | ColorValid
		// identity map for our builtin colors
		t.colors[Color(i)|ColorValid] = Color(i) | ColorValid
	}

	t.quit = make(chan struct{})
	t.eventQ = make(chan Event, 256)
	t.input = NewInputProcessor(t.eventQ)

	t.Lock()
	t.cx = -1
	t.cy = -1
	t.style = StyleDefault
	t.cells.Resize(w, h)
	t.cursorx = -1
	t.cursory = -1
	t.resize()
	t.Unlock()

	if err := t.engage(); err != nil {
		return err
	}

	return nil
}

func (t *tScreen) prepareBracketedPaste() {
	// Another workaround for lack of reporting in terminfo.
	// We assume if the terminal has a mouse entry, that it
	// offers bracketed paste.  But we allow specific overrides
	// via our terminal database.
	if t.ti.Mouse != "" || t.ti.XTermLike {
		t.enablePaste = "\x1b[?2004h"
		t.disablePaste = "\x1b[?2004l"
	}
}

func (t *tScreen) prepareUnderlines() {
	if t.ti.XTermLike {
		t.doubleUnder = "\x1b[4:2m"
		t.curlyUnder = "\x1b[4:3m"
		t.dottedUnder = "\x1b[4:4m"
		t.dashedUnder = "\x1b[4:5m"
		t.underColor = "\x1b[58:5:%p1%dm"
		t.underRGB = "\x1b[58:2::%p1%d:%p2%d:%p3%dm"
		t.underFg = "\x1b[59m"
	}
}

func (t *tScreen) prepareExtendedOSC() {
	// Linux is a special beast - because it has a mouse entry, but does
	// not swallow these OSC commands properly.
	if strings.Contains(t.ti.Name, "linux") {
		return
	}
	// More stuff for limits in terminfo.  This time we are applying
	// the most common OSC (operating system commands).  Generally
	// terminals that don't understand these will ignore them.
	// Again, we condition this based on mouse capabilities.
	if t.ti.Mouse != "" || t.ti.XTermLike {
		t.enterUrl = "\x1b]8;%p2%s;%p1%s\x1b\\"
		t.exitUrl = "\x1b]8;;\x1b\\"
	}

	if t.ti.Mouse != "" || t.ti.XTermLike {
		t.setWinSize = "\x1b[8;%p1%p2%d;%dt"
	}

	if t.ti.Mouse != "" || t.ti.XTermLike {
		t.enableFocus = "\x1b[?1004h"
		t.disableFocus = "\x1b[?1004l"
	}

	if t.ti.XTermLike {
		t.saveTitle = "\x1b[22;2t"
		t.restoreTitle = "\x1b[23;2t"
		// this also tries to request that UTF-8 is allowed in the title
		t.setTitle = "\x1b[>2t\x1b]2;%p1%s\x1b\\"
	}

	if t.setClipboard == "" && t.ti.XTermLike {
		// this string takes a base64 string and sends it to the clipboard.
		// it will also be able to retrieve the clipboard using "?" as the
		// sent string, when we support that.
		t.setClipboard = "\x1b]52;c;%p1%s\x1b\\"
	}

	if t.startSyncOut == "" && t.ti.XTermLike {
		// this is in theory a queryable private mode, but we just assume it will be ok
		// The terminals we have been able to test it all either just swallow it, or
		// handle it.
		t.startSyncOut = "\x1b[?2026h"
		t.endSyncOut = "\x1b[?2026l"
	}

	if t.enableCsiU == "" && t.ti.XTermLike {
		if runtime.GOOS == "windows" && os.Getenv("TERM") == "" {
			// on Windows, if we don't have a TERM, use only win32-input-mode
			t.enableCsiU = "\x1b[?9001h"
			t.disableCsiU = "\x1b[?9001l"
		} else {
			// three advanced keyboard protocols:
			// - xterm modifyOtherKeys (uses CSI 27 ~ )
			// - kitty csi-u (uses CSI u)
			// - win32-input-mode (uses CSI _)
			t.enableCsiU = "\x1b[>4;2m" + "\x1b[>1u" + "\x1b[?9001h"
			t.disableCsiU = "\x1b[?9001l" + "\x1b[<u" + "\x1b[>4;0m"
		}
	}
}

func (t *tScreen) prepareCursorStyles() {
	if t.ti.Mouse != "" || t.ti.XTermLike {
		t.cursorStyles = map[CursorStyle]string{
			CursorStyleDefault:           "\x1b[0 q",
			CursorStyleBlinkingBlock:     "\x1b[1 q",
			CursorStyleSteadyBlock:       "\x1b[2 q",
			CursorStyleBlinkingUnderline: "\x1b[3 q",
			CursorStyleSteadyUnderline:   "\x1b[4 q",
			CursorStyleBlinkingBar:       "\x1b[5 q",
			CursorStyleSteadyBar:         "\x1b[6 q",
		}
		if t.cursorRGB == "" {
			t.cursorRGB = "\x1b]12;#%p1%02x%p2%02x%p3%02x\007"
			t.cursorFg = "\x1b]112\007"
		}
	}
}

func (t *tScreen) prepareKeys() {
	ti := t.ti
	if strings.HasPrefix(ti.Name, "xterm") {
		// assume its some form of XTerm clone
		t.ti.XTermLike = true
		ti.XTermLike = true
	}
	t.prepareBracketedPaste()
	t.prepareCursorStyles()
	t.prepareUnderlines()
	t.prepareExtendedOSC()
}

func (t *tScreen) Fini() {
	t.finiOnce.Do(t.finish)
}

func (t *tScreen) finish() {
	close(t.quit)
	t.finalize()
}

func (t *tScreen) SetStyle(style Style) {
	t.Lock()
	if !t.fini {
		t.style = style
	}
	t.Unlock()
}

func (t *tScreen) encodeStr(s string) []byte {

	var dstBuf [128]byte
	var buf []byte
	nb := dstBuf[:]
	dst := 0
	var err error
	if enc := t.encoder; enc != nil {
		enc.Reset()
		dst, _, err = enc.Transform(nb, []byte(s), true)
	}
	if err != nil || dst == 0 || nb[0] == '\x1a' {
		// Combining characters are elided
		r, _ := utf8.DecodeRuneInString(s)
		if len(buf) == 0 {
			if acs, ok := t.acs[r]; ok {
				buf = append(buf, []byte(acs)...)
			} else if fb, ok := t.fallback[r]; ok {
				buf = append(buf, []byte(fb)...)
			} else {
				buf = append(buf, '?')
			}
		}
	} else {
		buf = append(buf, nb[:dst]...)
	}

	return buf
}

func (t *tScreen) sendFgBg(fg Color, bg Color, attr AttrMask) AttrMask {
	ti := t.ti
	if t.Colors() == 0 {
		// foreground vs background, we calculate luminance
		// and possibly do a reverse video
		if !fg.Valid() {
			return attr
		}
		v, ok := t.colors[fg]
		if !ok {
			v = FindColor(fg, []Color{ColorBlack, ColorWhite})
			t.colors[fg] = v
		}
		switch v {
		case ColorWhite:
			return attr
		case ColorBlack:
			return attr ^ AttrReverse
		}
	}

	if fg == ColorReset || bg == ColorReset {
		t.TPuts(ti.ResetFgBg)
	}
	if t.truecolor {
		if ti.SetFgBgRGB != "" && fg.IsRGB() && bg.IsRGB() {
			r1, g1, b1 := fg.RGB()
			r2, g2, b2 := bg.RGB()
			t.TPuts(ti.TParm(ti.SetFgBgRGB,
				int(r1), int(g1), int(b1),
				int(r2), int(g2), int(b2)))
			return attr
		}

		if fg.IsRGB() && ti.SetFgRGB != "" {
			r, g, b := fg.RGB()
			t.TPuts(ti.TParm(ti.SetFgRGB, int(r), int(g), int(b)))
			fg = ColorDefault
		}

		if bg.IsRGB() && ti.SetBgRGB != "" {
			r, g, b := bg.RGB()
			t.TPuts(ti.TParm(ti.SetBgRGB,
				int(r), int(g), int(b)))
			bg = ColorDefault
		}
	}

	if fg.Valid() {
		if v, ok := t.colors[fg]; ok {
			fg = v
		} else {
			v = FindColor(fg, t.palette)
			t.colors[fg] = v
			fg = v
		}
	}

	if bg.Valid() {
		if v, ok := t.colors[bg]; ok {
			bg = v
		} else {
			v = FindColor(bg, t.palette)
			t.colors[bg] = v
			bg = v
		}
	}

	if fg.Valid() && bg.Valid() && ti.SetFgBg != "" {
		t.TPuts(ti.TParm(ti.SetFgBg, int(fg&0xff), int(bg&0xff)))
	} else {
		if fg.Valid() && ti.SetFg != "" {
			t.TPuts(ti.TParm(ti.SetFg, int(fg&0xff)))
		}
		if bg.Valid() && ti.SetBg != "" {
			t.TPuts(ti.TParm(ti.SetBg, int(bg&0xff)))
		}
	}
	return attr
}

func (t *tScreen) drawCell(x, y int) int {

	ti := t.ti

	str, style, width := t.cells.Get(x, y)
	if !t.cells.Dirty(x, y) {
		return width
	}

	if y == t.h-1 && x == t.w-1 && t.ti.AutoMargin && ti.DisableAutoMargin == "" && ti.InsertChar != "" {
		// our solution is somewhat goofy.
		// we write to the second to the last cell what we want in the last cell, then we
		// insert a character at that 2nd to last position to shift the last column into
		// place, then we rewrite that 2nd to last cell.  Old terminals suck.
		t.TPuts(ti.TGoto(x-1, y))
		defer func() {
			t.TPuts(ti.TGoto(x-1, y))
			t.TPuts(ti.InsertChar)
			t.cy = y
			t.cx = x - 1
			t.cells.SetDirty(x-1, y, true)
			_ = t.drawCell(x-1, y)
			t.TPuts(t.ti.TGoto(0, 0))
			t.cy = 0
			t.cx = 0
		}()
	} else if t.cy != y || t.cx != x {
		t.TPuts(ti.TGoto(x, y))
		t.cx = x
		t.cy = y
	}

	if style == StyleDefault {
		style = t.style
	}
	if style != t.curstyle {
		fg, bg, attrs := style.fg, style.bg, style.attrs

		t.TPuts(ti.AttrOff)

		attrs = t.sendFgBg(fg, bg, attrs)
		if attrs&AttrBold != 0 {
			t.TPuts(ti.Bold)
		}
		if us, uc := style.ulStyle, style.ulColor; us != UnderlineStyleNone {
			if t.underColor != "" || t.underRGB != "" {
				if uc == ColorReset {
					t.TPuts(t.underFg)
				} else if uc.IsRGB() {
					if t.underRGB != "" {
						r, g, b := uc.RGB()
						t.TPuts(ti.TParm(t.underRGB, int(r), int(g), int(b)))
					} else {
						if v, ok := t.colors[uc]; ok {
							uc = v
						} else {
							v = FindColor(uc, t.palette)
							t.colors[uc] = v
							uc = v
						}
						t.TPuts(ti.TParm(t.underColor, int(uc&0xff)))
					}
				} else if uc.Valid() {
					t.TPuts(ti.TParm(t.underColor, int(uc&0xff)))
				}
			}
			t.TPuts(ti.Underline) // to ensure everyone gets at least a basic underline
			switch us {
			case UnderlineStyleDouble:
				t.TPuts(t.doubleUnder)
			case UnderlineStyleCurly:
				t.TPuts(t.curlyUnder)
			case UnderlineStyleDotted:
				t.TPuts(t.dottedUnder)
			case UnderlineStyleDashed:
				t.TPuts(t.dashedUnder)
			}
		}
		if attrs&AttrReverse != 0 {
			t.TPuts(ti.Reverse)
		}
		if attrs&AttrBlink != 0 {
			t.TPuts(ti.Blink)
		}
		if attrs&AttrDim != 0 {
			t.TPuts(ti.Dim)
		}
		if attrs&AttrItalic != 0 {
			t.TPuts(ti.Italic)
		}
		if attrs&AttrStrikeThrough != 0 {
			t.TPuts(ti.StrikeThrough)
		}

		// URL string can be long, so don't send it unless we really need to
		if t.enterUrl != "" && t.curstyle.url != style.url {
			if style.url != "" {
				t.TPuts(ti.TParm(t.enterUrl, style.url, style.urlId))
			} else {
				t.TPuts(t.exitUrl)
			}
		}

		t.curstyle = style
	}

	// now emit runes - taking care to not overrun width with a
	// wide character, and to ensure that we emit exactly one regular
	// character followed up by any residual combing characters

	if width < 1 {
		width = 1
	}

	buf := t.encodeStr(str)
	str = string(buf)

	if width > 1 && str == "?" {
		// No FullWidth character support
		str = "? "
		t.cx = -1
	}

	if x > t.w-width {
		// too wide to fit; emit a single space instead
		width = 1
		str = " "
	}
	if width > 1 {
		// Clobber over any content in the next cell.
		// This fixes a problem with some terminals where overwriting two
		// adjacent single cells with a wide rune would leave an image
		// of the second cell.  This is a workaround for buggy terminals.
		t.writeString("  \b\b")
	}

	t.writeString(str)
	t.cx += width
	t.cells.SetDirty(x, y, false)
	if width > 1 {
		t.cx = -1
	}

	return width
}

func (t *tScreen) ShowCursor(x, y int) {
	t.Lock()
	t.cursorx = x
	t.cursory = y
	t.Unlock()
}

func (t *tScreen) SetCursor(cs CursorStyle, cc Color) {
	t.Lock()
	t.cursorStyle = cs
	t.cursorColor = cc
	t.Unlock()
}

func (t *tScreen) HideCursor() {
	t.ShowCursor(-1, -1)
}

func (t *tScreen) showCursor() {

	x, y := t.cursorx, t.cursory
	w, h := t.cells.Size()
	if x < 0 || y < 0 || x >= w || y >= h {
		t.hideCursor()
		return
	}
	t.TPuts(t.ti.TGoto(x, y))
	t.TPuts(t.ti.ShowCursor)
	if t.cursorStyles != nil {
		if esc, ok := t.cursorStyles[t.cursorStyle]; ok {
			t.TPuts(esc)
		}
	}
	if t.cursorRGB != "" {
		if t.cursorColor == ColorReset {
			t.TPuts(t.cursorFg)
		} else if t.cursorColor.Valid() {
			r, g, b := t.cursorColor.RGB()
			t.TPuts(t.ti.TParm(t.cursorRGB, int(r), int(g), int(b)))
		}
	}
	t.cx = x
	t.cy = y
}

// writeString sends a string to the terminal. The string is sent as-is and
// this function does not expand inline padding indications (of the form
// $<[delay]> where [delay] is msec). In order to have these expanded, use
// TPuts. If the screen is "buffering", the string is collected in a buffer,
// with the intention that the entire buffer be sent to the terminal in one
// write operation at some point later.
func (t *tScreen) writeString(s string) {
	if t.buffering {
		_, _ = io.WriteString(&t.buf, s)
	} else {
		_, _ = io.WriteString(t.tty, s)
	}
}

func (t *tScreen) TPuts(s string) {
	if t.buffering {
		t.ti.TPuts(&t.buf, s)
	} else {
		t.ti.TPuts(t.tty, s)
	}
}

func (t *tScreen) Show() {
	t.Lock()
	if !t.fini {
		t.resize()
		t.draw()
	}
	t.Unlock()
}

func (t *tScreen) clearScreen() {
	t.TPuts(t.ti.AttrOff)
	t.TPuts(t.exitUrl)
	_ = t.sendFgBg(t.style.fg, t.style.bg, AttrNone)
	t.TPuts(t.ti.Clear)
	t.clear = false
}

func (t *tScreen) startBuffering() {
	t.TPuts(t.startSyncOut)
}

func (t *tScreen) endBuffering() {
	t.TPuts(t.endSyncOut)
}

func (t *tScreen) hideCursor() {
	// does not update cursor position
	if t.ti.HideCursor != "" {
		t.TPuts(t.ti.HideCursor)
	} else {
		// No way to hide cursor, stick it
		// at bottom right of screen
		t.cx, t.cy = t.cells.Size()
		t.TPuts(t.ti.TGoto(t.cx, t.cy))
	}
}

func (t *tScreen) draw() {
	// clobber cursor position, because we're going to change it all
	t.cx = -1
	t.cy = -1
	// make no style assumptions
	t.curstyle = styleInvalid

	t.buf.Reset()
	t.buffering = true
	t.startBuffering()
	defer func() {
		t.buffering = false
		t.endBuffering()
	}()

	// hide the cursor while we move stuff around
	t.hideCursor()

	if t.clear {
		t.clearScreen()
	}

	for y := 0; y < t.h; y++ {
		for x := 0; x < t.w; x++ {
			width := t.drawCell(x, y)
			if width > 1 {
				if x+1 < t.w {
					// this is necessary so that if we ever
					// go back to drawing that cell, we
					// actually will *draw* it.
					t.cells.SetDirty(x+1, y, true)
				}
			}
			x += width - 1
		}
	}

	// restore the cursor
	t.showCursor()

	_, _ = t.buf.WriteTo(t.tty)
}

func (t *tScreen) EnableMouse(flags ...MouseFlags) {
	var f MouseFlags
	flagsPresent := false
	for _, flag := range flags {
		f |= flag
		flagsPresent = true
	}
	if !flagsPresent {
		f = MouseMotionEvents | MouseDragEvents | MouseButtonEvents
	}

	t.Lock()
	t.mouseFlags = f
	t.enableMouse(f)
	t.Unlock()
}

func (t *tScreen) enableMouse(f MouseFlags) {
	// Rather than using terminfo to find mouse escape sequences, we rely on the fact that
	// pretty much *every* terminal that supports mouse tracking follows the
	// XTerm standards (the modern ones).
	if len(t.mouse) != 0 {
		// start by disabling all tracking.
		t.TPuts("\x1b[?1000l\x1b[?1002l\x1b[?1003l\x1b[?1006l")
		if f&MouseButtonEvents != 0 {
			t.TPuts("\x1b[?1000h")
		}
		if f&MouseDragEvents != 0 {
			t.TPuts("\x1b[?1002h")
		}
		if f&MouseMotionEvents != 0 {
			t.TPuts("\x1b[?1003h")
		}
		if f&(MouseButtonEvents|MouseDragEvents|MouseMotionEvents) != 0 {
			t.TPuts("\x1b[?1006h")
		}
	}

}

func (t *tScreen) DisableMouse() {
	t.Lock()
	t.mouseFlags = 0
	t.enableMouse(0)
	t.Unlock()
}

func (t *tScreen) EnablePaste() {
	t.Lock()
	t.pasteEnabled = true
	t.enablePasting(true)
	t.Unlock()
}

func (t *tScreen) DisablePaste() {
	t.Lock()
	t.pasteEnabled = false
	t.enablePasting(false)
	t.Unlock()
}

func (t *tScreen) enablePasting(on bool) {
	var s string
	if on {
		s = t.enablePaste
	} else {
		s = t.disablePaste
	}
	if s != "" {
		t.TPuts(s)
	}
}

func (t *tScreen) EnableFocus() {
	t.Lock()
	t.focusEnabled = true
	t.enableFocusReporting()
	t.Unlock()
}

func (t *tScreen) DisableFocus() {
	t.Lock()
	t.focusEnabled = false
	t.disableFocusReporting()
	t.Unlock()
}

func (t *tScreen) enableFocusReporting() {
	if t.enableFocus != "" {
		t.TPuts(t.enableFocus)
	}
}

func (t *tScreen) disableFocusReporting() {
	if t.disableFocus != "" {
		t.TPuts(t.disableFocus)
	}
}

func (t *tScreen) Size() (int, int) {
	t.Lock()
	w, h := t.w, t.h
	t.Unlock()
	return w, h
}

func (t *tScreen) resize() {
	ws, err := t.tty.WindowSize()
	if err != nil {
		return
	}
	if ws.Width == t.w && ws.Height == t.h {
		return
	}
	t.cx = -1
	t.cy = -1

	t.cells.Resize(ws.Width, ws.Height)
	t.cells.Invalidate()
	t.h = ws.Height
	t.w = ws.Width
	t.input.SetSize(ws.Width, ws.Height)
}

func (t *tScreen) Colors() int {
	if os.Getenv("NO_COLOR") != "" {
		return 0
	}
	// this doesn't change, no need for lock
	if t.truecolor {
		return 1 << 24
	}
	return t.ti.Colors
}

// nColors returns the size of the built-in palette.
// This is distinct from Colors(), as it will generally
// always be a small number. (<= 256)
func (t *tScreen) nColors() int {
	if os.Getenv("NO_COLOR") != "" {
		return 0
	}
	return t.ti.Colors
}

// vtACSNames is a map of bytes defined by terminfo that are used in
// the terminals Alternate Character Set to represent other glyphs.
// For example, the upper left corner of the box drawing set can be
// displayed by printing "l" while in the alternate character set.
// It's not quite that simple, since the "l" is the terminfo name,
// and it may be necessary to use a different character based on
// the terminal implementation (or the terminal may lack support for
// this altogether).  See buildAcsMap below for detail.
var vtACSNames = map[byte]rune{
	'+': RuneRArrow,
	',': RuneLArrow,
	'-': RuneUArrow,
	'.': RuneDArrow,
	'0': RuneBlock,
	'`': RuneDiamond,
	'a': RuneCkBoard,
	'b': '␉', // VT100, Not defined by terminfo
	'c': '␌', // VT100, Not defined by terminfo
	'd': '␋', // VT100, Not defined by terminfo
	'e': '␊', // VT100, Not defined by terminfo
	'f': RuneDegree,
	'g': RunePlMinus,
	'h': RuneBoard,
	'i': RuneLantern,
	'j': RuneLRCorner,
	'k': RuneURCorner,
	'l': RuneULCorner,
	'm': RuneLLCorner,
	'n': RunePlus,
	'o': RuneS1,
	'p': RuneS3,
	'q': RuneHLine,
	'r': RuneS7,
	's': RuneS9,
	't': RuneLTee,
	'u': RuneRTee,
	'v': RuneBTee,
	'w': RuneTTee,
	'x': RuneVLine,
	'y': RuneLEqual,
	'z': RuneGEqual,
	'{': RunePi,
	'|': RuneNEqual,
	'}': RuneSterling,
	'~': RuneBullet,
}

// buildAcsMap builds a map of characters that we translate from Unicode to
// alternate character encodings.  To do this, we use the standard VT100 ACS
// maps.  This is only done if the terminal lacks support for Unicode; we
// always prefer to emit Unicode glyphs when we are able.
func (t *tScreen) buildAcsMap() {
	acsstr := t.ti.AltChars
	t.acs = make(map[rune]string)
	for len(acsstr) > 2 {
		srcv := acsstr[0]
		dstv := string(acsstr[1])
		if r, ok := vtACSNames[srcv]; ok {
			t.acs[r] = t.ti.EnterAcs + dstv + t.ti.ExitAcs
		}
		acsstr = acsstr[2:]
	}
}

func (t *tScreen) scanInput(buf *bytes.Buffer) {
	for buf.Len() > 0 {
		utf := make([]byte, min(8, max(buf.Len()*2, 128)))
		nOut, nIn, e := t.decoder.Transform(utf, buf.Bytes(), true)
		_ = buf.Next(nIn)
		t.input.ScanUTF8(utf[:nOut])
		if e == transform.ErrShortSrc {
			return
		}
	}
}

func (t *tScreen) mainLoop(stopQ chan struct{}) {
	defer t.wg.Done()
	buf := &bytes.Buffer{}
	for {
		select {
		case <-stopQ:
			return
		case <-t.quit:
			return
		case <-t.resizeQ:
			t.Lock()
			t.cx = -1
			t.cy = -1
			t.resize()
			t.cells.Invalidate()
			t.draw()
			t.Unlock()
			continue
		case chunk := <-t.keychan:
			buf.Write(chunk)
			t.scanInput(buf)
		}
	}
}

func (t *tScreen) inputLoop(stopQ chan struct{}) {

	defer t.wg.Done()
	for {
		select {
		case <-stopQ:
			return
		default:
		}
		chunk := make([]byte, 128)
		n, e := t.tty.Read(chunk)
		switch e {
		case nil:
		default:
			t.Lock()
			running := t.running
			t.Unlock()
			if running {
				select {
				case t.eventQ <- NewEventError(e):
				case <-t.quit:
				}
			}
			return
		}
		if n > 0 {
			t.keychan <- chunk[:n]
		}
	}
}

func (t *tScreen) Sync() {
	t.Lock()
	t.cx = -1
	t.cy = -1
	if !t.fini {
		t.resize()
		t.clear = true
		t.cells.Invalidate()
		t.draw()
	}
	t.Unlock()
}

func (t *tScreen) CharacterSet() string {
	return t.charset
}

func (t *tScreen) RegisterRuneFallback(orig rune, fallback string) {
	t.Lock()
	t.fallback[orig] = fallback
	t.Unlock()
}

func (t *tScreen) UnregisterRuneFallback(orig rune) {
	t.Lock()
	delete(t.fallback, orig)
	t.Unlock()
}

func (t *tScreen) CanDisplay(r rune, checkFallbacks bool) bool {

	if enc := t.encoder; enc != nil {
		nb := make([]byte, 6)

		enc.Reset()
		dst, _, err := enc.Transform(nb, []byte(string(r)), true)
		if dst != 0 && err == nil && nb[0] != '\x1A' {
			return true
		}
	}
	// Terminal fallbacks always permitted, since we assume they are
	// basically nearly perfect renditions.
	if _, ok := t.acs[r]; ok {
		return true
	}
	if !checkFallbacks {
		return false
	}
	if _, ok := t.fallback[r]; ok {
		return true
	}
	return false
}

func (t *tScreen) HasMouse() bool {
	return len(t.mouse) != 0
}

func (t *tScreen) HasKey(_ Key) bool {
	// We always return true
	return true
}

func (t *tScreen) SetSize(w, h int) {
	if t.setWinSize != "" {
		t.TPuts(t.ti.TParm(t.setWinSize, w, h))
	}
	t.cells.Invalidate()
	t.resize()
}

func (t *tScreen) Resize(int, int, int, int) {}

func (t *tScreen) Suspend() error {
	t.disengage()
	return nil
}

func (t *tScreen) Resume() error {
	return t.engage()
}

func (t *tScreen) Tty() (Tty, bool) {
	return t.tty, true
}

// engage is used to place the terminal in raw mode and establish screen size, etc.
// Think of this is as tcell "engaging" the clutch, as it's going to be driving the
// terminal interface.
func (t *tScreen) engage() error {
	t.Lock()
	defer t.Unlock()
	if t.tty == nil {
		return ErrNoScreen
	}
	t.tty.NotifyResize(func() {
		select {
		case t.resizeQ <- true:
		default:
		}
	})
	if t.running {
		return errors.New("already engaged")
	}
	if err := t.tty.Start(); err != nil {
		return err
	}
	t.running = true
	if ws, err := t.tty.WindowSize(); err == nil && ws.Width != 0 && ws.Height != 0 {
		t.cells.Resize(ws.Width, ws.Height)
	}
	stopQ := make(chan struct{})
	t.stopQ = stopQ
	t.enableMouse(t.mouseFlags)
	t.enablePasting(t.pasteEnabled)
	if t.focusEnabled {
		t.enableFocusReporting()
	}
	ti := t.ti
	if os.Getenv("TCELL_ALTSCREEN") != "disable" {
		// Technically this may not be right, but every terminal we know about
		// (even Wyse 60) uses this to enter the alternate screen buffer, and
		// possibly save and restore the window title and/or icon.
		// (In theory there could be terminals that don't support X,Y cursor
		// positions without a setup command, but we don't support them.)
		t.TPuts(ti.EnterCA)
		t.TPuts(t.saveTitle)
	}
	t.TPuts(ti.EnterKeypad)
	t.TPuts(ti.HideCursor)
	t.TPuts(ti.EnableAcs)
	t.TPuts(ti.DisableAutoMargin)
	t.TPuts(ti.Clear)
	if t.title != "" && t.setTitle != "" {
		t.TPuts(t.ti.TParm(t.setTitle, t.title))
	}
	t.TPuts(t.enableCsiU)

	t.wg.Add(2)
	go t.inputLoop(stopQ)
	go t.mainLoop(stopQ)
	return nil
}

// disengage is used to release the terminal back to support from the caller.
// Think of this as tcell disengaging the clutch, so that another application
// can take over the terminal interface.  This restores the TTY mode that was
// present when the application was first started.
func (t *tScreen) disengage() {

	t.Lock()
	if !t.running {
		t.Unlock()
		return
	}

	t.running = false
	stopQ := t.stopQ
	close(stopQ)
	_ = t.tty.Drain()
	t.Unlock()

	t.tty.NotifyResize(nil)
	// wait for everything to shut down
	t.wg.Wait()

	// shutdown the screen and disable special modes (e.g. mouse and bracketed paste)
	ti := t.ti
	t.cells.Resize(0, 0)
	t.TPuts(ti.ShowCursor)
	if t.cursorStyles != nil && t.cursorStyle != CursorStyleDefault {
		t.TPuts(t.cursorStyles[CursorStyleDefault])
	}
	if t.cursorFg != "" && t.cursorColor.Valid() {
		t.TPuts(t.cursorFg)
	}
	t.TPuts(ti.ResetFgBg)
	t.TPuts(ti.AttrOff)
	t.TPuts(ti.ExitKeypad)
	t.TPuts(ti.EnableAutoMargin)
	t.TPuts(t.disableCsiU)
	if os.Getenv("TCELL_ALTSCREEN") != "disable" {
		if t.restoreTitle != "" {
			t.TPuts(t.restoreTitle)
		}
		t.TPuts(ti.Clear) // only needed if ExitCA is empty
		t.TPuts(ti.ExitCA)
	}
	t.enableMouse(0)
	t.enablePasting(false)
	t.disableFocusReporting()

	_ = t.tty.Stop()
}

// Beep emits a beep to the terminal.
func (t *tScreen) Beep() error {
	t.writeString(string(byte(7)))
	return nil
}

// finalize is used to at application shutdown, and restores the terminal
// to it's initial state.  It should not be called more than once.
func (t *tScreen) finalize() {
	t.disengage()
	_ = t.tty.Close()
}

func (t *tScreen) StopQ() <-chan struct{} {
	return t.quit
}

func (t *tScreen) EventQ() chan Event {
	return t.eventQ
}

func (t *tScreen) GetCells() *CellBuffer {
	return &t.cells
}

func (t *tScreen) SetTitle(title string) {
	t.Lock()
	t.title = title
	if t.setTitle != "" && t.running {
		t.TPuts(t.ti.TParm(t.setTitle, title))
	}
	t.Unlock()
}

func (t *tScreen) SetClipboard(data []byte) {
	// Post binary data to the system clipboard.  It might be UTF-8, it might not be.
	t.Lock()
	if t.setClipboard != "" {
		encoded := base64.StdEncoding.EncodeToString(data)
		t.TPuts(t.ti.TParm(t.setClipboard, encoded))
	}
	t.Unlock()
}

func (t *tScreen) GetClipboard() {
	t.Lock()
	if t.setClipboard != "" {
		t.TPuts(t.ti.TParm(t.setClipboard, "?"))
	}
	t.Unlock()
}
