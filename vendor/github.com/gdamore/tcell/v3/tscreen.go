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

//go:build (!js && !wasm) || (js && wasm)
// +build !js,!wasm js,wasm

package tcell

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gdamore/tcell/v3/color"
	"github.com/gdamore/tcell/v3/vt"
	"golang.org/x/text/transform"
)

// NewTerminfoScreen returns a Screen that uses the stock TTY interface
// and POSIX terminal control, combined with a terminfo description taken from
// the $TERM environment variable.  It returns an error if the terminal
// is not supported for any reason.
//
// For terminals that do not support dynamic resize events, the $LINES
// $COLUMNS environment variables can be set to the actual window size,
// otherwise defaults taken from the terminal database are used.
func NewTerminfoScreen(opts ...TerminfoScreenOption) (Screen, error) {
	return NewTerminfoScreenFromTty(nil, opts...)
}

type TerminfoScreenOption interface {
	apply(*tScreen)
}

// OptColors forces the number of colors, overriding the value
// of the color count that would be detected by the environment.
// If the value is 0, then color is forced off.  Other reasonable values
// are 8, 16, 88, 256, or 1<<24.  The latter case intrinsically enables
// 24-bit color as well.
type OptColors int

func (o OptColors) apply(t *tScreen) {
	t.ncolor = min(int(o), 256)
	t.truecolor = o > 256
	t.noColor = o == 0
}

// OptTerm overrides the detection of $TERM.
type OptTerm string

func (o OptTerm) apply(t *tScreen) {
	t.term = string(o)
}

// OptAltScreen controls whether the alternate screen buffer is used.
// The default is true. The TCELL_ALTSCREEN=disable environment override
// is still honored.
type OptAltScreen bool

func (o OptAltScreen) apply(t *tScreen) {
	t.altScreen = bool(o)
}

// OptSanitizeContent enables stripping control characters from content passed
// to Put and PutStr. This is safer, but a little slower than leaving content
// unsanitized.
type OptSanitizeContent bool

func (o OptSanitizeContent) apply(t *tScreen) {
	t.cells.sanitizeContent = bool(o)
}

// OptAdvancedKeys enables richer key reporting where supported.  In this mode
// key events may include release state, repeat counts, and physical keys, and
// ASCII control letters are reported as KeyRune with ModCtrl instead of
// KeyCtrlA through KeyCtrlZ.  Shift-Tab is reported as KeyTab with ModShift,
// rather than KeyBacktab.
type OptAdvancedKeys bool

func (o OptAdvancedKeys) apply(t *tScreen) {
	t.advancedKeys = bool(o)
}

// OptKeyboardProtocol forces the keyboard reporting protocol instead of using
// startup negotiation. The zero value forces legacy keyboard reporting.
type OptKeyboardProtocol KeyProtocol

func (o OptKeyboardProtocol) apply(t *tScreen) {
	t.forceKeyboardProtocol(KeyProtocol(o))
}

// OptNegotiation controls whether terminal capabilities are negotiated during
// startup. The default is true.
type OptNegotiation bool

func (o OptNegotiation) apply(t *tScreen) {
	t.negotiate = bool(o)
}

// OptControlStringLimit sets the maximum inbound control-string payload size
// accepted from the terminal before the parser drops the sequence. This limits
// OSC and XDA strings, including OSC 52 clipboard strings; OSC 52 is the
// protocol used for writing clipboard data through the terminal. The default is
// 64 KiB; a value of 0 disables the limit.
type OptControlStringLimit int

func (o OptControlStringLimit) apply(t *tScreen) {
	t.controlStringLimit = max(int(o), 0)
}

// Some terminal escapes that are basically universal.
// We would really like to be able to use private mode queries for some of
// these but generally we've found that support for queries is not always present,
// even when the private modes can be controlled. It appears that *all* terminals
// will happily swallow the escapes that they do not recognize, with the small annoyance
// in "st" where it prints error messages to its stderr (which is usually not visible
// to the user unless they started it from another terminal session).  But apart from
// the complaint to stderr from "st", everything else is fine.
const (
	enableAutoMargin  = "\x1b[?7h" // dec private mode 7
	setCursorPosition = "\x1b[%[1]d;%[2]dH"
	sgr0              = "\x1b[m" // attrOff
	bold              = "\x1b[1m"
	dim               = "\x1b[2m"
	italic            = "\x1b[3m"
	underline         = "\x1b[4m"
	blink             = "\x1b[5m"
	reverse           = "\x1b[7m"
	strikeThrough     = "\x1b[9m"
	clear             = "\x1b[H\x1b[J"
	doubleUnder       = "\x1b[4:2m"
	curlyUnder        = "\x1b[4:3m"
	dottedUnder       = "\x1b[4:4m"
	dashedUnder       = "\x1b[4:5m"
	underColor        = "\x1b[58:5:%dm"
	underRGB          = "\x1b[58:2::%d:%d:%dm"
	underFg           = "\x1b[59m"
	enableAltChars    = "\x1b(B\x1b)0"                      // set G0 as US-ASCII, G1 as DEC line drawing
	startAltChars     = "\x0e"                              // aka Shift-Out
	endAltChars       = "\x0f"                              // aka Shift-In
	setFg8            = "\x1b[3%dm"                         // for colors less than 8
	setFg256          = "\x1b[38;5;%dm"                     // for colors less than 256
	setFgRgb          = "\x1b[38;2;%d;%d;%dm"               // for RGB
	setBg8            = "\x1b[4%dm"                         // color colors less than 8
	setBg256          = "\x1b[48;5;%dm"                     // for colors less than 256
	setBgRgb          = "\x1b[48;2;%d;%d;%dm"               // for RGB
	setFgBgRgb        = "\x1b[38;2;%d;%d;%d;48;2;%d;%d;%dm" // for RGB, in one shot
	enterCA           = "\x1b[?1049h"                       // alternate screen
	exitCA            = "\x1b[?1049l"                       // alternate screen
	enterKeypad       = "\x1b[?1h\x1b="                     // Note mode 1 might not be supported everywhere
	exitKeypad        = "\x1b[?1l\x1b>"                     // Also mode 1
	requestWindowSize = "\x1b[18t"                          // For modern terminals
	requestPrimaryDA  = "\x1b[c"                            // Request primary device attributes
	requestExtAttr    = "\x1b[>q"                           // Request extended attribute (emulator name and version)
	setClipboard      = "\x1b]52;c;%s\x1b\\"                // Clipboard content is base64
	notifyDesktop9    = "\x1b]9;%[2]s\x1b\\"                // Args are title, body (but OSC 9 only has body)
	notifyDesktop777  = "\x1b]777;notify;%s;%s\x1b\\"       // Most commonly supported
	queryKittyKbd     = "\x1b[?u"                           // Query for Kitty keyboard support
	enableKittyKbd    = "\x1b[=1u"                          // Technically this pushes
	enableKittyKbdAdv = "\x1b[=15u"                         // disambiguation, events, alternate keys, all keys
	disableKittyKbd   = "\x1b[=0u"                          // Technically this means pop previous mode
	queryXTermKbd     = "\x1b[?4m"                          // Query for XTerm modify other keys support
	enableXTermKbd    = "\x1b[>4;2m"                        // Enable modify other keys protocol
	disableXTermKbd   = "\x1b[>4;0m"                        // Disable modify other keys protocol
)

// NewTerminfoScreenFromTty returns a Screen using a custom Tty implementation.
// If the passed in tty is nil, then a reasonable default (typically /dev/tty)
// is presumed, at least on UNIX hosts. (Windows hosts will typically fail this
// call altogether.)
func NewTerminfoScreenFromTty(tty Tty, opts ...TerminfoScreenOption) (Screen, error) {
	t := &tScreen{
		tty:                tty,
		altScreen:          true,
		negotiate:          true,
		controlStringLimit: defaultControlStringLimit,
	}

	t.prepareCursorStyles()
	t.prepareExtendedOSC()
	t.buildAcsMap()
	t.resizeQ = make(chan bool, 1)
	t.fallback = make(map[rune]string)
	maps.Copy(t.fallback, RuneFallbacks)
	for _, o := range opts {
		o.apply(t)
	}

	return &baseScreen{screenImpl: t}, nil
}

// tScreen represents a screen backed by a terminfo implementation.
type tScreen struct {
	tty                Tty
	h                  int
	w                  int
	fini               bool
	cells              CellBuffer
	buffering          bool // true if we are collecting writes to buf instead of sending directly to out
	buf                bytes.Buffer
	curstyle           Style
	style              Style
	resizeQ            chan bool
	quit               chan struct{}
	keyQ               chan []byte
	cx                 int
	cy                 int
	cls                bool // clear screen
	cursorx            int
	cursory            int
	acs                map[rune]string
	charset            string
	encoder            transform.Transformer
	decoder            transform.Transformer
	fallback           map[rune]string
	ncolor             int
	colors             map[color.Color]color.Color
	palette            []color.Color
	truecolor          bool
	noColor            bool
	legacy             bool
	hasClipboard       bool // true if OSC 52 reported via DA1
	finiOnce           sync.Once
	enterUrl           string
	exitUrl            string
	setWinSize         string
	cursorStyles       map[CursorStyle]string
	cursorStyle        CursorStyle
	cursorColor        color.Color
	cursorRGB          string
	cursorFg           string
	stopQ              chan struct{}
	eventQ             chan Event
	initQ              chan Event
	initted            bool
	running            bool
	startTime          time.Time
	wg                 sync.WaitGroup
	mouseFlags         MouseFlags
	pasteEnabled       bool
	focusEnabled       bool
	setTitle           string
	saveTitle          string
	restoreTitle       string
	title              string
	setClipboard       string
	notifyDesktop      string
	termName           string
	termVers           string
	term               string // value from $TERM
	altScreen          bool
	inlineResize       bool
	haveMouse          bool
	haveMouseSgr       bool
	haveKittyKbd       bool
	haveWin32Kbd       bool
	haveXTermKbd       bool
	forcedKbd          KeyProtocol
	forceKbd           bool
	negotiate          bool
	mouseDisabled      bool
	advancedKeys       bool
	controlStringLimit int
	input              *inputParser
	sync.Mutex
}

func (t *tScreen) useAltScreen() bool {
	return t.altScreen && os.Getenv("TCELL_ALTSCREEN") != "disable"
}

func validKeyboardProtocol(p KeyProtocol) bool {
	switch p {
	case LegacyKeyboard, KittyKeyboard, Win32Keyboard, XTermKeyboard:
		return true
	default:
		return false
	}
}

func parseKeyboardProtocol(s string) (KeyProtocol, bool) {
	switch s {
	case "legacy":
		return LegacyKeyboard, true
	case "kitty":
		return KittyKeyboard, true
	case "win32":
		return Win32Keyboard, true
	case "xterm":
		return XTermKeyboard, true
	default:
		return LegacyKeyboard, false
	}
}

func (t *tScreen) forceKeyboardProtocol(p KeyProtocol) bool {
	if !validKeyboardProtocol(p) {
		return false
	}
	t.forcedKbd = p
	t.forceKbd = true
	return true
}

func (t *tScreen) applyKeyboardProtocolOverride() {
	if !t.forceKbd {
		return
	}
	t.haveKittyKbd = t.forcedKbd == KittyKeyboard
	t.haveWin32Kbd = t.forcedKbd == Win32Keyboard
	t.haveXTermKbd = t.forcedKbd == XTermKeyboard
}

func (t *tScreen) applyEnvironmentOverrides() {
	switch os.Getenv("TCELL_KEYBOARD_PROTOCOL") {
	case "auto":
		t.forceKbd = false
	case "":
	default:
		if p, ok := parseKeyboardProtocol(os.Getenv("TCELL_KEYBOARD_PROTOCOL")); ok {
			t.forceKeyboardProtocol(p)
		}
	}

	switch os.Getenv("TCELL_NEGOTIATE") {
	case "auto":
		t.negotiate = true
	case "disable":
		t.negotiate = false
	}

	t.mouseDisabled = os.Getenv("TCELL_MOUSE") == "disable"
}

func (t *tScreen) Init() error {
	if e := t.initialize(); e != nil {
		return e
	}

	t.startTime = time.Now()
	t.keyQ = make(chan []byte, 10)

	t.charset = getCharset()
	if enc := GetEncoding(t.charset); enc != nil {
		t.encoder = enc.NewEncoder()
		t.decoder = enc.NewDecoder()
	} else {
		return ErrNoCharset
	}

	// environment overrides
	w := 80
	h := 24
	if i, _ := strconv.Atoi(os.Getenv("LINES")); i != 0 {
		h = i
	}
	if i, _ := strconv.Atoi(os.Getenv("COLUMNS")); i != 0 {
		w = i
	}
	if t.term == "" {
		t.term = os.Getenv("TERM")
	}
	nterm := t.term

	if t.ncolor == 0 && !t.noColor {
		cterm := os.Getenv("COLORTERM")

		// On Windows, enable 24-bit color by default (all terminals there are 24-bit capable)
		if runtime.GOOS == "windows" {
			t.truecolor = true
			t.ncolor = 256
		} else if slices.Contains([]string{"truecolor", "direct", "24bit"}, cterm) || strings.HasSuffix(nterm, "-direct") || strings.HasSuffix(nterm, "-truecolor") {
			t.truecolor = true
			t.ncolor = 256 // base 8-bit palette
		} else if strings.HasSuffix(nterm, "-256color") || strings.Contains(cterm, "256") {
			t.ncolor = 256
		} else if strings.HasSuffix(nterm, "-88color") {
			t.ncolor = 88
		} else if strings.HasSuffix(nterm, "-16color") {
			t.ncolor = 16
		} else if strings.Contains(nterm, "color") || cterm != "" {
			t.ncolor = 8
		} else if strings.Contains(nterm, "mono") || strings.HasSuffix(nterm, "-m") { // monochrome variants
			t.ncolor = 0
		} else if strings.Contains(nterm, "ansi") || slices.Contains([]string{"dtterm", "xterm", "aixterm", "linux"}, nterm) {
			t.ncolor = 8
		} else if strings.HasPrefix(nterm, "vt") || nterm == "sun" {
			// legacy DEC VT 100/220 etc. family.  (technically the VT525 can do ANSI, but they should set to ansi)
			t.ncolor = 0
		} else {
			// best guess - this covers all the modern variants like ghostty,
			t.ncolor = 256
		}
		if os.Getenv("NO_COLOR") != "" {
			t.truecolor = false
			t.ncolor = 0
			t.noColor = true
		}
		// A user who wants to have his themes honored can set this environment variable.
		if os.Getenv("TCELL_TRUECOLOR") == "disable" {
			t.truecolor = false
		}
	}

	if strings.HasPrefix(nterm, "vt") || strings.Contains(nterm, "ansi") || nterm == "linux" || nterm == "sun" || nterm == "sun-color" {
		// these terminals are "legacy" and not expected to support most OSC functions
		t.legacy = true
	}

	t.applyEnvironmentOverrides()

	t.initted = false
	t.quit = make(chan struct{})
	t.initQ = make(chan Event, 32)
	t.eventQ = make(chan Event, 128)
	t.input = newInputParser(t.filterEvents())
	t.input.advanced = t.advancedKeys
	t.input.controlStringMax = t.controlStringLimit

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

	// clip to reasonable limits
	nColors := min(t.ncolor, 256)
	t.colors = make(map[color.Color]color.Color, nColors)
	t.palette = make([]color.Color, nColors)
	for i := range nColors {
		t.palette[i] = color.PaletteColor(i)
		// identity map for our builtin colors
		t.colors[color.PaletteColor(i)] = color.PaletteColor(i)
	}

	return nil
}

func (t *tScreen) processInitQ() {
	// NB: called with lock held
	if t.initted {
		return
	}

	expire := time.After(time.Second)

	for {
		select {
		case <-expire:
			t.initted = true
			return
		case ev := <-t.initQ:
			switch ev := ev.(type) {
			case *eventPrimaryAttributes:
				if ev.Color && t.ncolor == 0 && !t.noColor {
					t.ncolor = 8
				}
				if ev.Clipboard && t.setClipboard == "" {
					t.setClipboard = setClipboard
				}
				t.hasClipboard = ev.Clipboard
				t.initted = true
				return
			case *eventTermName:
				// terminal specific overrides
				t.termName = ev.Name
				t.termVers = ev.Version
				switch ev.Name {
				case "iTerm2":
					// Some terminals can use OSC 9.  Unfortunately we can only discover
					// them using this means.  It appears that pretty much all of them
					// except iTerm2 also support more standard OSC 777, and it seems like
					// only Kitty has its OSC 99 thing, but it also does OSC 777 well.
					t.notifyDesktop = notifyDesktop9
				}
			case *eventPrivateMode:
				switch ev.Mode {
				case vt.PmResizeReports:
					t.inlineResize = ev.Status.Changeable()
				case vt.PmMouseSgr:
					t.haveMouseSgr = ev.Status.Changeable()
				case vt.PmMouseButton:
					t.haveMouse = ev.Status.Changeable()
				case vt.PmWin32Input:
					t.haveWin32Kbd = ev.Status.Changeable()
				}
			case *eventKittyKbdMode:
				t.haveKittyKbd = true
			case *eventXTermKbdMode:
				t.haveXTermKbd = true
			}
		}
	}
}

func (t *tScreen) filterEvents() chan Event {
	inQ := make(chan Event, 128)
	go func() {
		for {
			var ev Event
			select {
			case ev = <-inQ:
			case <-t.quit:
				return
			}
			switch ev.(type) {
			case *eventTermName, *eventPrimaryAttributes, *eventPrivateMode, *eventKittyKbdMode, *eventXTermKbdMode:
				select {
				case t.initQ <- ev:
				default:
				}

			default:
				t.eventQ <- ev
			}
		}
	}()
	return inQ
}

func (t *tScreen) prepareExtendedOSC() {
	if t.legacy {
		return
	}

	// OSC 8 is for enter/exit URL.
	t.enterUrl = "\x1b]8;%[2]s;%[1]s\x1b\\"
	t.exitUrl = "\x1b]8;;\x1b\\"

	// CSI .. t is for window operations.
	t.setWinSize = "\x1b[8;%[2]d;%[1]dt"
	t.saveTitle = "\x1b[22;2t"
	t.restoreTitle = "\x1b[23;2t"
	// this also tries to request that UTF-8 is allowed in the title
	t.setTitle = "\x1b[>2t\x1b]2;%s\x1b\\"

	// OSC 52 is for saving to the clipboard.
	// this string takes a base64 string and sends it to the clipboard.
	// it will also be able to retrieve the clipboard using "?" as the
	// sent string, when we support that.
	t.setClipboard = setClipboard

	// OSC 777 is the desktop notification supported by a variety of
	// newer terminals.  (There was also OSC 9 and OSC 99, but they
	// are not as widely deployed, and OSC 9 is not unique.)
	t.notifyDesktop = notifyDesktop777
}

func (t *tScreen) prepareCursorStyles() {
	t.cursorStyles = map[CursorStyle]string{
		CursorStyleDefault:           "\x1b[0 q",
		CursorStyleBlinkingBlock:     "\x1b[1 q",
		CursorStyleSteadyBlock:       "\x1b[2 q",
		CursorStyleBlinkingUnderline: "\x1b[3 q",
		CursorStyleSteadyUnderline:   "\x1b[4 q",
		CursorStyleBlinkingBar:       "\x1b[5 q",
		CursorStyleSteadyBar:         "\x1b[6 q",
	}
	if t.legacy {
		return
	}
	if t.cursorRGB == "" {
		t.cursorRGB = "\x1b]12;#%02x%02x%02x\007"
		t.cursorFg = "\x1b]112\007"
	}
}

func (t *tScreen) Fini() {
	// Ensure that enough time passes for terminals to  finish sending
	// their initial response (gnome-terminal sends terminal dimensions
	// asynchronously later than the response to primary DA for some reason.)
	if time.Since(t.startTime) < 50*time.Millisecond {
		time.Sleep(time.Millisecond * 50)
	}
	t.finiOnce.Do(t.finish)
}

func (t *tScreen) finish() {
	t.Lock()
	t.fini = true
	t.Unlock()
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

// resolvePalette looks up a color to obtain the palette entry for it.
func (t *tScreen) resolvePalette(c Color) Color {
	if v, ok := t.colors[c]; ok {
		return v
	}
	v := color.Find(c, t.palette)
	t.colors[c] = v
	return v
}

// sendFgBg sends the foreground and background.  It is assumed that sgr0
// was already emitted prior to calling this (so colors are already in default).
func (t *tScreen) sendFgBg(fg Color, bg Color, attr AttrMask) AttrMask {
	if t.Colors() == 0 {
		// foreground vs background, we calculate luminance
		// and possibly do a reverse video
		if !fg.Valid() {
			return attr
		}
		v, ok := t.colors[fg]
		if !ok {
			v = color.Find(fg, []Color{ColorBlack, ColorWhite})
			t.colors[fg] = v
		}
		switch v {
		case ColorWhite:
			return attr
		case ColorBlack:
			return attr ^ AttrReverse
		}
	}

	if t.truecolor {
		if fg.IsRGB() && bg.IsRGB() {
			r1, g1, b1 := fg.RGB()
			r2, g2, b2 := bg.RGB()
			t.Printf(setFgBgRgb, r1, g1, b1, r2, g2, b2)
			return attr
		}

		if fg.IsRGB() {
			r, g, b := fg.RGB()
			t.Printf(setFgRgb, r, g, b)
			fg = ColorDefault
		}

		if bg.IsRGB() {
			r, g, b := bg.RGB()
			t.Printf(setBgRgb, r, g, b)
			bg = ColorDefault
		}
	}

	if fg.Valid() {
		fg = t.resolvePalette(fg)
		fgc := fg & 0xffffff
		if fgc < 8 {
			t.Printf(setFg8, fgc)
		} else if fgc < 256 {
			t.Printf(setFg256, fgc)
		}
	}

	if bg.Valid() {
		bg = t.resolvePalette(bg)
		bgc := bg & 0xffffff
		if bgc < 8 {
			t.Printf(setBg8, bgc)
		} else if bgc < 256 {
			t.Printf(setBg256, bgc)
		}
	}

	return attr
}

// emitAttrs dumps prints the attributes, aside from underline that is special
// The assumption is that sgr0 was already printed ahead of this.
func (t *tScreen) emitAttrs(attrs AttrMask) {

	if attrs&AttrBold != 0 {
		t.Print(bold)
	}
	if attrs&AttrReverse != 0 {
		t.Print(reverse)
	}
	if attrs&AttrBlink != 0 {
		t.Print(blink)
	}
	if attrs&AttrDim != 0 {
		t.Print(dim)
	}
	if attrs&AttrItalic != 0 {
		t.Print(italic)
	}
	if attrs&AttrStrikeThrough != 0 {
		t.Print(strikeThrough)
	}
}

// emitUl dumps prints the underline, which may be colored.
// The assumption is that sgr0 was already printed ahead of this.
func (t *tScreen) emitUnderline(us UnderlineStyle, uc Color) {
	if us != UnderlineStyleNone {
		// NB: under color should have been reset by sgr0
		if uc.IsRGB() {
			r, g, b := uc.RGB()
			uc = t.resolvePalette(uc)
			t.Printf(underColor, uc&0xff)
			t.Printf(underRGB, r, g, b)
		} else if uc.Valid() {
			t.Printf(underColor, uc&0xff)
		}

		t.Print(underline) // to ensure everyone gets at least a basic underline
		switch us {
		case UnderlineStyleDouble:
			t.Print(doubleUnder)
		case UnderlineStyleCurly:
			t.Print(curlyUnder)
		case UnderlineStyleDotted:
			t.Print(dottedUnder)
		case UnderlineStyleDashed:
			t.Print(dashedUnder)
		}
	}
}

// emitUrl either emits a url (OSC 8), or if the string is empty
// then the OSC 8 to exit the URL.  It should only be called if we
// either have a new URL, or need to exit an old one, as it always emits
// the OSC 8 sequence (if OSC 8 is supported).
func (t *tScreen) emitUrl(u urlInfo) {
	if t.enterUrl != "" {
		if u.url != "" {
			t.Printf(t.enterUrl, u.url, u.id)
		} else {
			t.Print(t.exitUrl)
		}
	}
}

// urlNeedsEmission reports whether a hyperlink transition has any wire effect.
// Url ids can be staged before the Url itself, and id-only transitions have no
// OSC 8 representation of their own.
func urlNeedsEmission(oldUrl, newUrl urlInfo) bool {
	return oldUrl != newUrl && (oldUrl.url != "" || newUrl.url != "")
}

func (t *tScreen) drawCell(x, y int) int {

	str, style, width := t.cells.Get(x, y)
	if !t.cells.Dirty(x, y) {
		return width
	}

	if t.cy != y || t.cx != x {
		t.Printf(setCursorPosition, y+1, x+1)
		t.cx = x
		t.cy = y
	}

	if style == StyleDefault {
		style = t.style
	}
	if style != t.curstyle {
		fg, bg, attrs := style.fg, style.bg, style.attrs

		t.Print(sgr0)

		attrs = t.sendFgBg(fg, bg, attrs)
		t.emitAttrs(attrs)
		t.emitUnderline(style.ulStyle, style.ulColor)

		var newUrl urlInfo
		var oldUrl urlInfo
		if t.curstyle.url != nil {
			oldUrl = *t.curstyle.url
		}
		if style.url != nil {
			newUrl = *style.url
		}
		// URL string can be long, so don't send it unless we really need to.
		if urlNeedsEmission(oldUrl, newUrl) {
			t.emitUrl(newUrl)
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
	if width > 1 && x+width < t.w {
		// Clobber over any content in the next cell.
		// This fixes a problem with some terminals where overwriting two
		// adjacent single cells with a wide rune would leave an image
		// of the second cell.  This is a workaround for buggy terminals.
		t.Print("  \b\b")
	}
	t.Print(str)
	t.cx += width
	t.cells.SetDirty(x, y, false)
	if width > 1 && len([]rune(str)) > 1 {
		t.cx = -1
	}

	return width
}

func (t *tScreen) ShowCursor(x, y int) {
	t.Lock()
	if t.fini {
		t.Unlock()
		return
	}
	t.cursorx = x
	t.cursory = y
	t.Unlock()
}

func (t *tScreen) SetCursor(cs CursorStyle, cc Color) {
	t.Lock()
	if t.fini {
		t.Unlock()
		return
	}
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
	t.Printf(setCursorPosition, y+1, x+1)
	t.Print(vt.PmShowCursor.Enable())
	if t.cursorStyles != nil {
		if esc, ok := t.cursorStyles[t.cursorStyle]; ok {
			t.Print(esc)
		}
	}
	if t.cursorRGB != "" {
		if t.cursorColor == ColorReset {
			t.Print(t.cursorFg)
		} else if t.cursorColor.Valid() {
			r, g, b := t.cursorColor.RGB()
			t.Printf(t.cursorRGB, r, g, b)
		}
	}
	t.cx = x
	t.cy = y
}

func (t *tScreen) Write(b []byte) (int, error) {
	if t.buffering {
		return t.buf.Write(b)
	}
	return t.tty.Write(b)
}

func (t *tScreen) Print(s string) {
	_, _ = io.WriteString(t, s)
}

func (t *tScreen) Printf(f string, args ...any) {
	_, _ = fmt.Fprintf(t, f, args...)
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
	t.Print(sgr0)
	t.Print(t.exitUrl)
	_ = t.sendFgBg(t.style.fg, t.style.bg, AttrNone)

	t.Print(clear)

	t.cls = false
}

func (t *tScreen) startBuffering() {
	t.Print(vt.PmSyncOutput.Enable())
}

func (t *tScreen) endBuffering() {
	t.Print(vt.PmSyncOutput.Disable())
}

func (t *tScreen) hideCursor() {
	// just in case we cannot hide it, move it to the end
	t.cx, t.cy = t.cells.Size()
	t.Printf(setCursorPosition, t.cy+1, t.cx+1)
	// then hide it
	t.Print(vt.PmShowCursor.Disable())
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

	if t.cls {
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

	if t.curstyle.url != nil && t.curstyle.url.url != "" {
		t.emitUrl(urlInfo{})
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
	if t.fini {
		t.Unlock()
		return
	}
	t.mouseFlags = f
	t.enableMouse(f)
	t.Unlock()
}

func (t *tScreen) enableMouse(f MouseFlags) {
	// Rather than using terminfo to find mouse escape sequences, we rely on the fact that
	// pretty much *every* terminal that supports mouse tracking follows the
	// XTerm standards (the modern ones).  It is expected that all terminals understand
	// the same DEC private modes.  Note that the SGR mode is required for the mouse sequences
	// to be understood.

	// We rely on dec private mode queries for this.
	// If your terminal doesn't support these, then ask them to fix it.
	// Note that as of macOS 26, macOS Terminal does not support them,
	// so we enable the mouse unconditionally unless we get a report
	// that says we have mouse, but not SGR mouse.  This is suboptimal, but
	// a concession forced by the sorry state of terminal emulators.
	if t.mouseDisabled {
		f = 0
	}
	if f != 0 && t.haveMouse && !t.haveMouseSgr {
		return
	}

	// start by disabling all tracking.
	t.Print(vt.PmMouseButton.Disable())
	t.Print(vt.PmMouseDrag.Disable())
	t.Print(vt.PmMouseMotion.Disable())
	t.Print(vt.PmMouseSgr.Disable())
	t.Print(vt.PmMouseSgrPixel.Disable())

	pixel := f&MousePixelEvents != 0
	t.input.SetPixelMouse(pixel)

	if f&(MouseButtonEvents|MouseDragEvents|MouseMotionEvents) != 0 {
		t.Print(vt.PmMouseButton.Enable())
	}
	if f&MouseDragEvents != 0 {
		t.Print(vt.PmMouseDrag.Enable())
	}
	if f&MouseMotionEvents != 0 {
		t.Print(vt.PmMouseMotion.Enable())
	}
	if f&(MouseButtonEvents|MouseDragEvents|MouseMotionEvents) != 0 {
		if pixel {
			t.Print(vt.PmMouseSgrPixel.Enable())
		} else {
			t.Print(vt.PmMouseSgr.Enable())
		}
	}
}

func (t *tScreen) DisableMouse() {
	t.Lock()
	if t.fini {
		t.Unlock()
		return
	}
	t.mouseFlags = 0
	t.enableMouse(0)
	t.Unlock()
}

func (t *tScreen) EnablePaste() {
	t.Lock()
	if t.fini {
		t.Unlock()
		return
	}
	t.pasteEnabled = true
	t.enablePasting(true)
	t.Unlock()
}

func (t *tScreen) DisablePaste() {
	t.Lock()
	if t.fini {
		t.Unlock()
		return
	}
	t.pasteEnabled = false
	t.enablePasting(false)
	t.Unlock()
}

func (t *tScreen) enablePasting(on bool) {
	var s string
	if on {
		s = vt.PmBracketedPaste.Enable()
	} else {
		s = vt.PmBracketedPaste.Disable()
	}
	if s != "" {
		t.Print(s)
	}
}

func (t *tScreen) EnableFocus() {
	t.Lock()
	if t.fini {
		t.Unlock()
		return
	}
	t.focusEnabled = true
	t.enableFocusReporting()
	t.Unlock()
}

func (t *tScreen) DisableFocus() {
	t.Lock()
	if t.fini {
		t.Unlock()
		return
	}
	t.focusEnabled = false
	t.disableFocusReporting()
	t.Unlock()
}

func (t *tScreen) enableFocusReporting() {
	t.Print(vt.PmFocusReports.Enable())
}

func (t *tScreen) disableFocusReporting() {
	t.Print(vt.PmFocusReports.Disable())
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
	// this doesn't change, no need for lock
	if t.truecolor {
		return 1 << 24
	}
	return t.ncolor
}

// vtACSNames is a map of bytes defined by terminfo that are used in
// the terminals Alternate Character Set to represent other glyphs.
// For example, the upper left corner of the box drawing set can be
// displayed by printing "l" while in the alternate character set.
// It's not quite that simple, since the "l" is the terminfo name,
// and it may be necessary to use a different character based on
// the terminal implementation (or the terminal may lack support for
// this altogether).  These values are from the DEC VT100, and all
// modern terminal emulators support this as charset 0.
var vtACSNames = map[byte]rune{
	'`': RuneDiamond,
	'a': RuneCkBoard,
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
	const acsstr = "``aaffggjjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~"

	t.acs = make(map[rune]string)
	for b, r := range vtACSNames {
		t.acs[r] = startAltChars + string(b) + endAltChars
	}
}

func (t *tScreen) scanInput(buf *bytes.Buffer) {
	// The end of the buffer isn't necessarily the end of the input, because
	// large inputs are chunked. Set atEOF to false so the UTF-8 validating decoder
	// returns ErrShortSrc instead of ErrInvalidUTF8 for incomplete multi-byte codepoints.
	const atEOF = false

	for buf.Len() > 0 {
		utf := make([]byte, min(8, max(buf.Len()*2, 128)))
		nOut, nIn, e := t.decoder.Transform(utf, buf.Bytes(), atEOF)
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
	var ta <-chan time.Time
	for {
		select {
		case <-stopQ:
			return
		case <-t.quit:
			return
		case <-t.resizeQ:
			go func() {
				t.Lock()
				t.cx = -1
				t.cy = -1
				t.resize()
				t.cells.Invalidate()
				t.draw()
				t.Unlock()
			}()
			continue
		case chunk := <-t.keyQ:
			buf.Write(chunk)
			t.scanInput(buf)
			if t.input.Waiting() {
				ta = time.After(time.Millisecond * 100)
			} else {
				ta = nil
			}
		case <-ta:
			t.input.Scan()
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
			t.keyQ <- chunk[:n]
		}
	}
}

func (t *tScreen) Sync() {
	t.Lock()
	t.cx = -1
	t.cy = -1
	if !t.fini {
		t.resize()
		t.cls = true
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

func (t *tScreen) SetSize(w, h int) {
	t.Lock()
	defer t.Unlock()
	if t.fini {
		return
	}
	if t.setWinSize != "" {
		t.Printf(t.setWinSize, w, h)
	}
	t.cells.Invalidate()
	t.resize()
}

func (t *tScreen) Resize(int, int, int, int) {}

func (t *tScreen) Suspend() error {
	t.Lock()
	if t.fini {
		t.Unlock()
		return nil
	}
	finish := t.disengageStart()
	t.Unlock()
	if finish {
		t.disengageFinish()
	}
	return nil
}

func (t *tScreen) Resume() error {
	t.Lock()
	defer t.Unlock()
	if t.fini {
		return nil
	}
	return t.engageLocked()
}

func (t *tScreen) Tty() (Tty, bool) {
	return t.tty, true
}

func (t *tScreen) applyKnownTerminalProfile(goos, termProgram string) bool {
	switch termProgram {
	case "Apple_Terminal":
		// macOS Terminal.app cannot handle the startup queries, but it does
		// support modern mouse reporting.
		t.haveMouse = true
		t.haveMouseSgr = true
		t.termName = "Terminal.app"
		t.termVers = os.Getenv("TERM_PROGRAM_VERSION")
		return true
	case "WezTerm":
		// The WezTerm keyboard protocol to use is in theory driven by its
		// own configuration, but we have found this unreliable because it
		// does not mask unsupported capabilities.  Furthermore, on Windows
		// builds the kitty protocol implementation is broken, while on other
		// builds win32-input-mode is broken.  This is a best effort to make
		// WezTerm work reasonably; our stronger advice is to choose another
		// terminal program altogether.  This workaround will probably not
		// apply to ssh sessions, as TERM_PROGRAM is not normally propagated.
		if goos == "windows" {
			t.haveWin32Kbd = true
		} else {
			t.haveKittyKbd = true
			t.haveWin32Kbd = false
		}
		t.haveMouse = true
		t.haveMouseSgr = true
		t.initted = true
		t.termName = "WezTerm"
		t.termVers = os.Getenv("TERM_PROGRAM_VERSION")
		return true
	}
	return false
}

func useVTWindowSizeQuery(goos string) bool {
	return goos != "windows"
}

func useXTermKeyboardQuery(goos string) bool {
	return goos != "windows"
}

// engage is used to place the terminal in raw mode and establish screen size, etc.
// Think of this is as tcell "engaging" the clutch, as it's going to be driving the
// terminal interface.
func (t *tScreen) engage() error {
	t.Lock()
	defer t.Unlock()
	return t.engageLocked()
}

// engageLocked is engage's implementation when t's lock is already held.
func (t *tScreen) engageLocked() error {
	if t.tty == nil {
		return ErrNoScreen
	}
	if t.running {
		return errors.New("already engaged")
	}
	if err := t.tty.Start(); err != nil {
		return err
	}

	stopQ := make(chan struct{})
	t.stopQ = stopQ
	t.wg.Add(2)
	go t.inputLoop(stopQ)
	go t.mainLoop(stopQ)

	if !t.initted {
		// macOS Terminal.app is brain damaged
		// https://garrett.damore.org/2025/12/macos-terminal-still-missing-mark-apple.html
		// Eventually they'll hopefully fix this.  As the environment variable
		// does not convey by default via ssh, remote sessions might see spurious characters
		// emitted during startup.  See the blog post for alternatives.
		if !t.applyKnownTerminalProfile(runtime.GOOS, os.Getenv("TERM_PROGRAM")) && t.negotiate {
			if useVTWindowSizeQuery(runtime.GOOS) {
				t.Print(requestWindowSize)
			}
			t.Print(vt.PmResizeReports.Query())
			t.Print(vt.PmMouseButton.Query())
			t.Print(vt.PmMouseSgr.Query())
			if !t.forceKbd {
				t.Print(vt.PmWin32Input.Query())
				t.Print(queryKittyKbd)
				if useXTermKeyboardQuery(runtime.GOOS) {
					// XTerm's modifyOtherKeys mode is mainly useful for XTerm
					// itself, and we do not use it on Windows.
					t.Print(queryXTermKbd)
				}
			}
			t.Print(requestExtAttr)
		}
		if !t.negotiate {
			t.initted = true
		} else if !t.initted {
			t.Print(requestPrimaryDA) // NB: MUST BE LAST
		}
	}
	t.processInitQ()
	t.applyKeyboardProtocolOverride()
	if t.useAltScreen() {
		// Technically this may not be right, but every terminal we know about
		// (even Wyse 60) uses this to enter the alternate screen buffer, and
		// possibly save and restore the window title and/or icon.
		// (In theory there could be terminals that don't support X,Y cursor
		// positions without a setup command, but we don't support them.)
		t.Print(enterCA)
		t.Print(t.saveTitle)
	}
	if t.haveWin32Kbd {
		t.Print(vt.PmWin32Input.Enable())
	} else if t.haveKittyKbd {
		if t.advancedKeys {
			t.Print(enableKittyKbdAdv)
		} else {
			t.Print(enableKittyKbd)
		}
	} else if t.haveXTermKbd {
		t.Print(enableXTermKbd)
	}

	t.running = true
	if ws, err := t.tty.WindowSize(); err == nil && ws.Width != 0 && ws.Height != 0 {
		t.cells.Resize(ws.Width, ws.Height)
	}
	t.enableMouse(t.mouseFlags)
	t.enablePasting(t.pasteEnabled)
	if t.focusEnabled {
		t.enableFocusReporting()
	}
	t.Print(enterKeypad)
	t.Print(enableAltChars)
	t.Print(vt.PmShowCursor.Disable())
	t.Print(vt.PmAutoMargin.Disable())
	t.Print(clear)
	if t.title != "" && t.setTitle != "" {
		t.Printf(t.setTitle, t.title)
	}
	if t.negotiate && useVTWindowSizeQuery(runtime.GOOS) {
		t.Print(requestWindowSize)
	}

	if t.inlineResize {
		t.Print(vt.PmResizeReports.Enable())
	} else {
		t.tty.NotifyResize(t.resizeQ)
	}
	return nil
}

// disengage is used to release the terminal back to support from the caller.
// Think of this as tcell disengaging the clutch, so that another application
// can take over the terminal interface.  This restores the TTY mode that was
// present when the application was first started.
func (t *tScreen) disengage() {
	t.Lock()
	finish := t.disengageStart()
	t.Unlock()
	if finish {
		t.disengageFinish()
	}
}

// disengageStart begins a disengage operation while t's lock is already held.
// It returns true when disengageFinish must be called after releasing the lock.
func (t *tScreen) disengageStart() bool {
	if !t.running {
		return false
	}

	t.running = false
	if t.inlineResize {
		t.Print(vt.PmResizeReports.Disable())
	} else {
		t.tty.NotifyResize(nil)
	}
	stopQ := t.stopQ
	close(stopQ)
	_ = t.tty.Drain()
	return true
}

// disengageFinish completes a disengage operation after disengageStart has
// released the running loops.
func (t *tScreen) disengageFinish() {
	// wait for everything to shut down
	t.wg.Wait()

	// shutdown the screen and disable special modes (e.g. mouse and bracketed paste)
	t.cells.Resize(0, 0)
	t.Print(vt.PmShowCursor.Enable())
	if t.cursorStyles != nil && t.cursorStyle != CursorStyleDefault {
		t.Print(t.cursorStyles[CursorStyleDefault])
	}
	if t.cursorFg != "" && t.cursorColor.Valid() {
		t.Print(t.cursorFg)
	}
	t.Print(exitKeypad)
	t.Print(sgr0)
	t.Print(vt.PmAutoMargin.Enable())
	if t.haveWin32Kbd {
		t.Print(vt.PmWin32Input.Disable())
	}
	if t.haveKittyKbd {
		t.Print(disableKittyKbd)
	}
	if t.haveXTermKbd {
		t.Print(disableXTermKbd)
	}

	// Hack for Windows.
	if runtime.GOOS == "windows" {
		t.Print(vt.PmWin32Input.Disable())
	}

	// t.Print(t.disableCsiU)
	if t.useAltScreen() {
		t.Print(t.restoreTitle)
		t.Print(clear)
		t.Print(exitCA)
	}
	t.enableMouse(0)
	t.enablePasting(false)
	t.disableFocusReporting()

	_ = t.tty.Stop()
}

// Beep emits a beep to the terminal.
func (t *tScreen) Beep() error {
	t.Lock()
	defer t.Unlock()
	if t.fini {
		return nil
	}
	t.Print(string(byte(7)))
	return nil
}

// finalize is used to at application shutdown, and restores the terminal
// to it's initial state.  It should not be called more than once.
func (t *tScreen) finalize() {
	t.disengage()
	_ = t.tty.Close()
	close(t.eventQ)
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
	if t.fini {
		t.Unlock()
		return
	}
	t.title = stripOSCControlsIfNeeded(title)
	if t.setTitle != "" && t.running {
		t.Printf(t.setTitle, t.title)
	}
	t.Unlock()
}

func (t *tScreen) SetClipboard(data []byte) {
	// Post binary data to the system clipboard.  It might be UTF-8, it might not be.
	t.Lock()
	if t.fini {
		t.Unlock()
		return
	}
	if t.setClipboard != "" {
		encoded := base64.StdEncoding.EncodeToString(data)
		t.Printf(t.setClipboard, encoded)
	}
	t.Unlock()
}

func (t *tScreen) GetClipboard() {
	t.Lock()
	if t.fini {
		t.Unlock()
		return
	}
	if t.setClipboard != "" {
		t.Printf(t.setClipboard, "?")
	}
	t.Unlock()
}

func (t *tScreen) HasClipboard() bool {
	return t.hasClipboard
}

func (t *tScreen) ShowNotification(title string, body string) {
	t.Lock()
	if t.fini {
		t.Unlock()
		return
	}
	t.Printf(t.notifyDesktop, stripOSCControlsIfNeeded(title), stripOSCControlsIfNeeded(body))
	t.Unlock()
}

func (t *tScreen) Terminal() (string, string) {
	t.Lock()
	defer t.Unlock()
	return t.termName, t.termVers
}

func (t *tScreen) KeyboardProtocol() KeyProtocol {
	t.Lock()
	defer t.Unlock()
	if t.haveWin32Kbd {
		return Win32Keyboard
	}
	if t.haveKittyKbd {
		return KittyKeyboard
	}
	if t.haveXTermKbd {
		return XTermKeyboard
	}
	return LegacyKeyboard
}
