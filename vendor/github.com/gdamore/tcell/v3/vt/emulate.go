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
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/clipperhouse/uax29/v2/graphemes"
	"github.com/gdamore/tcell/v3/color"
)

// Emulator is a terminal emulator API. It implements the state machinery
// (escape parsing and so forth) associated with being a terminal emulator.
// The backend handles rendering the content, and some low level details.
//
// NOTE: This is not a committed interface yet, its entirely a work in progress.
type Emulator interface {
	// SetId sets our identity.
	SetId(name string, version string)

	// SendRaw sends raw data to the consumer.  This bypasses the normal encoding,
	// so it should be used with caution.
	SendRaw([]byte)

	// KeyEvent injects a keyboard event into the emulator, which will ultimately
	// result in data being sent via SendRaw.
	KeyEvent(ev KeyEvent)

	// ResizeEvent is called by a backend when the terminal has resized
	// This will send in-band resize notifications if the client has requested them.
	ResizeEvent(Coord)

	// MouseEvent is called by a backend to report mouse activity.
	MouseEvent(ev MouseEvent)

	// FocusEvent is called by a backend to report that focus is gained (true) or lost (false).
	FocusEvent(bool)

	// Drain waits until any queued but not processed input has finished processing.
	// It also wakes the reader.
	Drain() error

	// Start starts processing.
	Start() error

	// Stop stops processing.
	Stop() error

	// Reader reads data from the emulator.  These are bytes that would be transmitted
	// to a remote party.
	io.Reader

	// Writer writes data to the emulator.  These are commands that the emulator should process.
	io.Writer
}

// Style represents the styling of a cell.
// This is an interface to prevent direct modification.
type Style interface {
	Fg() color.Color              // Fg returns the foreground color.
	Bg() color.Color              // Bg returns the background color.
	Uc() color.Color              // Uc returns the underline color.k
	Attr() Attr                   // Attr returns the associated attributes.
	Url() (string, string)        // Url returns the URL and associated id if one was set.
	WithFg(color.Color) Style     // WithFg creates a new style with the foreground
	WithBg(color.Color) Style     // WithBg creates a new style with the background.
	WithUc(color.Color) Style     // WithUc creates a new style with the underline color
	WithAttr(Attr) Style          // WithAttr creates a new style with the attributes.
	WithUrl(string, string) Style // WithLink creates a new style with the URL and id.
	Equal(Style) bool             // Equal returns true if the styles are the same.
}

// styleStruct implements Style.  Note that it is possible to make this even more
// compact, but we don't think further optimization here on size will justify the
// complexity and runtime performance hit to do so.  We're also already only storing
// a class reference to this per cell.
type styleStruct struct {
	fg   color.Color
	bg   color.Color
	uc   color.Color // underline color
	attr Attr
	url  string // URL
	id   string // Id for link
}

var BaseStyle = &styleStruct{}

var asciiRuneStrings = func() [utf8.RuneSelf]string {
	var table [utf8.RuneSelf]string
	for i := 0; i < utf8.RuneSelf; i++ {
		table[i] = string(rune(i))
	}
	return table
}()

const (
	runeStringCacheSize      = 32
	clusterStringCacheSize   = 32
	clusterStringCacheMaxLen = 128
)

type runeStringCache struct {
	entries [runeStringCacheSize]runeStringCacheEntry
	n       int
}

type runeStringCacheEntry struct {
	r rune
	s string
}

func (c *runeStringCache) stringFor(r rune) string {
	for i := 0; i < c.n; i++ {
		if c.entries[i].r == r {
			return c.entries[i].s
		}
	}

	s := string(r)
	n := c.n
	if n < len(c.entries) {
		n++
	}
	copy(c.entries[1:n], c.entries[:n-1])
	c.entries[0] = runeStringCacheEntry{r: r, s: s}
	c.n = n
	return s
}

type clusterStringCache struct {
	entries [clusterStringCacheSize]clusterStringCacheEntry
	n       int
}

type clusterStringCacheEntry struct {
	n int
	b [clusterStringCacheMaxLen]byte
	s string
}

func (c *clusterStringCache) stringFor(cluster []byte) string {
	if len(cluster) == 0 {
		return ""
	}
	if len(cluster) > clusterStringCacheMaxLen {
		return string(cluster)
	}
	for i := 0; i < c.n; i++ {
		e := &c.entries[i]
		if e.n == len(cluster) && bytes.Equal(e.b[:e.n], cluster) {
			return e.s
		}
	}

	s := string(cluster)
	n := c.n
	if n < len(c.entries) {
		n++
	}
	copy(c.entries[1:n], c.entries[:n-1])
	e := &c.entries[0]
	e.n = len(cluster)
	copy(e.b[:], cluster)
	e.s = s
	c.n = n
	return s
}

func (ss *styleStruct) Fg() color.Color              { return ss.fg }
func (ss *styleStruct) Bg() color.Color              { return ss.bg }
func (ss *styleStruct) Uc() color.Color              { return ss.uc }
func (ss *styleStruct) Attr() Attr                   { return ss.attr }
func (ss *styleStruct) Url() (string, string)        { return ss.url, ss.id }
func (ss *styleStruct) WithFg(fg color.Color) Style  { ns := *ss; ns.fg = fg; return &ns }
func (ss *styleStruct) WithBg(bg color.Color) Style  { ns := *ss; ns.bg = bg; return &ns }
func (ss *styleStruct) WithUc(uc color.Color) Style  { ns := *ss; ns.uc = uc; return &ns }
func (ss *styleStruct) WithAttr(a Attr) Style        { ns := *ss; ns.attr = a; return &ns }
func (ss *styleStruct) WithUrl(url, id string) Style { ns := *ss; ns.url = url; ns.id = id; return &ns }
func (ss *styleStruct) Equal(other Style) bool {
	if s2, ok := other.(*styleStruct); ok {
		return *ss == *s2
	}
	// We have chosen not to support alternative implementations for this compare.
	// We could delegate to the other style, but that could lead to a loop if they
	// do the same.
	return (false)
}

// Cell is a representation of a display cell. Most consumers will not need this.
// Note, this is not the simplest possible representation, and a 256x256 cell
// display is going to need about 3MB to store it all, but it's simple, and adequate
// to retain pretty much all of what we need for Unicode.  We could save some memory
// by using explicit struct pointers and by eliminating grapheme cluster support,
// but modern users expect these features.
type Cell struct {
	C string // Content, it will be a grapheme cluster
	S Style  // Style, a pointer is used efficiency
	W int    // Display width (0, 1, or 2)
}

// NewEmulator creates an emulator instance on top of the given backend.
// The input is relative to the emulator, so it receives data from the host,
// whereas the emulator sends data to the application through the output.
func NewEmulator(be Backend) Emulator {
	stopQ := make(chan bool)
	defStyle := BaseStyle.WithFg(color.Silver).WithBg(color.Black)
	em := &emulator{
		be:           be,
		inBuf:        &bytes.Buffer{},
		writeQ:       make(chan any),
		readQ:        make(chan any, 1024),
		stopQ:        stopQ,
		style:        defStyle,
		defaultStyle: defStyle,
		ansiModes: map[AnsiMode]ModeStatus{
			AmNewLineMode: ModeOff,
		},
		localModes: map[PrivateMode]ModeStatus{
			PmAppCursor:       ModeOff,
			PmAutoMargin:      ModeOn,
			PmAutoRepeat:      ModeOn,
			PmVT52:            ModeOnLocked, // we never support VT52 mode (note ON means ANSI mode)
			PmLeftRightMargin: ModeOff,
			PmShowCursor:      ModeOn,
			PmBlinkCursor:     ModeOn,
			PmWin32Input:      ModeOff,
		},
		mouseReports: MouseDisabled,
	}
	if _, ok := be.(Resizer); ok {
		em.localModes[PmResizeReports] = ModeOff
	}

	// add mouse modes - we also add focus reporting mode
	if _, ok := be.(Mouser); ok {
		em.localModes[PmMouseX10] = ModeOff
		em.localModes[PmMouseButton] = ModeOff
		em.localModes[PmMouseDrag] = ModeOff
		em.localModes[PmMouseMotion] = ModeOff
		em.localModes[PmMouseSgr] = ModeOff
		em.localModes[PmFocusReports] = ModeOff
	}

	if ak, ok := be.(AdvancedKeyboard); ok && ak.IsAdvancedKeyboard() {
		em.localModes[PmWin32Input] = ModeOff
	}

	em.size = em.be.GetSize()
	em.topMargin = 0
	em.botMargin = em.size.Y - 1
	em.ltMargin = 0
	em.rtMargin = em.size.X - 1
	em.cells = make([]Cell, int(em.size.X)*int(em.size.Y))
	em.graphemeIter = *graphemes.FromBytes(nil)
	close(stopQ)
	em.inb = em.inbInit
	em.cursor = BlinkingBlock
	em.be.SetCursor(em.cursor)
	return em
}

// emulator is an implementation of a terminal emulator built on top of
// a Backend.  It implements the common escape sequence handling and high
// level functionality that a real terminal emulator, or a mock, would need.
type emulator struct {
	stopQ          chan bool
	writeQ         chan any // queues data from application to emulator
	readQ          chan any // queues data from emulator to application
	be             Backend
	inBuf          *bytes.Buffer // buffer queued for input
	inb            func(byte)    // input byte function (faster than state switch)
	style          Style
	defaultStyle   Style
	utfLen         int
	pos            Coord
	buffering      uint         // reference count - number of (re-entrant) buffering calls
	autoWrap       bool         // next character will wrap (auto margin, deferred until char emitted)
	sevenOnly      bool         // only allow 7-bit escapes (needed for KOI8, ShiftJIS, etc.)
	appKeyPad      bool         // use application key pad keys?
	name           string       // name of this emulator (used for extended attributes)
	vers           string       // version string of this emulator (used for extended attributes)
	saved          savedCursor  // data saved by save cursor (DECSC)
	sendLock       sync.Mutex   // ensures that send data cannot be intermixed
	modeLock       sync.RWMutex // protects localModes/ansiModes and related derived state
	tabStops       []Col        // tab stops, ordered. if nil every 8th position is used
	lastIndex      int          // index of last cell written + 1 (for grapheme clustering) (zero means none)
	graphemeBuf    []byte       // scratch buffer for grapheme clustering checks
	graphemeIter   graphemes.Iterator[[]byte]
	runeStrings    runeStringCache
	clusterStrings clusterStringCache
	cells          []Cell         // content of cells, we have to maintain our own copy (backend might or might not)
	mouseReports   MouseReporting // whether we have enabled mouse reports
	size           Coord          // physical window size
	topMargin      Row            // top margin, scrollable region includes this row
	botMargin      Row            // bottom margin, scrollable region includes this row
	ltMargin       Col            // left margin, scrollable region to the right
	rtMargin       Col            // right margin, scrollable region to the left
	cursor         CursorStyle    // current cursor style (visibility, blink, shape)

	localModes map[PrivateMode]ModeStatus // some modes we handle locally
	ansiModes  map[AnsiMode]ModeStatus    // some modes we handle locally
}

// savedCursor is the content we save when saving the cursor,
// which is more than just the cursor location itself.
type savedCursor struct {
	pos      Coord
	style    Style
	autoWrap bool
	// We should probably store OSC 8 data here, eventually.
	// TODO: Character sets
	// TODO: Origin mode (DEC Mode 6)
}

func (em *emulator) saveCursor() {
	em.saved.pos = em.getPosition()
	em.saved.style = em.style
	em.saved.autoWrap = em.autoWrap
}

func (em *emulator) restoreCursor() {
	em.setPosition(em.saved.pos)
	em.autoWrap = em.saved.autoWrap
	em.style = em.saved.style
}

func (em *emulator) bufferingStart() {
	em.buffering++
	if em.buffering == 1 {
		em.be.Buffering(true)
	}
}

func (em *emulator) bufferingEnd() {
	em.buffering--
	if em.buffering == 0 {
		em.be.Buffering(false)
	}
}

// inbInit processes bytes received in the "default" state. Most often these are just
// text characters to display on screen, but if ESC is seen then additional processing will result.
func (em *emulator) inbInit(b byte) {
	em.inBuf.Reset()

	// hot path - just doing ASCII directly.
	if b >= ' ' && b < 0x7f {
		// plain ascii
		em.putRune(rune(b))
		return
	}

	// For 8-bit encodings, we treat these as Fe sequences.
	// Basically the same as ESC followed by (b - 0x40).
	// TODO: condition this so that we do not do this if
	// the encoding cannot support it (UTF, 8859, and EUC encodings
	// are all fine here, but others like ShiftJIS or KOI8 might not be).
	if b >= 0x80 && b <= 0x9F && !em.sevenOnly {
		em.inbEsc(b - 0x40)
		return
	}

	// TODO: To support non-UTF-8 locales, include a check here for > 0x7F.  Those locales
	// might preclude 8-bit control sequences - 8859 character sets are fine, but e.g. KOI8,
	// and ShiftJIS use values in those ranges.

	switch b {
	case 0x1b: // ESC (escape)
		em.inb = em.inbEsc
	case 0x07: // BEL (bell)
		em.beep()
	case 0x08: // BS (backspace)
		em.moveLeft()
	case 0x09: // horizontal tab
		em.nextTab()
	case 0x0a, 0x0b, 0x0c: // LF (line feed), VF, FF
		em.processLineFeed()
	case 0x0d: // CR (carriage return)
		em.processCarriageReturn()
	case 0x0e: // TODO: SO
		em.lastIndex = 0
	case 0x0f: // TODO: SI
		em.lastIndex = 0
	case 0x18: //TODO Cancel (reset parser)
		em.lastIndex = 0
	default:
		// TODO: consider separating Unicode from other 8-bit character sets
		if b&0xE0 == 0xC0 {
			em.utfLen = 2
			em.inb = em.inbUTF
			em.inBuf.WriteByte(b)
		} else if b&0xF0 == 0xE0 {
			em.utfLen = 3
			em.inb = em.inbUTF
			em.inBuf.WriteByte(b)
		} else if b&0xF8 == 0xF0 {
			em.utfLen = 4
			em.inb = em.inbUTF
			em.inBuf.WriteByte(b)
		} else {
			em.lastIndex = 0
			em.beep()
		}
	}
}

// inbEsc processes the next byte after an escape character is seen.
func (em *emulator) inbEsc(b byte) {

	// By default, reset to init state. Other states will be set explicitly as needed.
	em.inb = em.inbInit
	em.lastIndex = 0

	switch b {
	case '[':
		em.inb = em.inbCSI
	case ']':
		em.inb = em.inbOSC
	case ' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/':
		// 0x20 - 0x2F -- usually followed by just one terminating character, but could include others
		em.inb = em.inbNF
		em.inBuf.WriteByte(b)
	case '^': // privacy message (PM)
		em.inb = em.inbStr
	case '_': // application program command (APC)
		em.inb = em.inbStr
	case '=':
		em.appKeyPad = true
	case '>':
		em.appKeyPad = false
	case 'D': // down one line (IND)
		em.processIndex()
	case 'E': // next line (NEL)
		em.nextLine()
	case 'H': // set tab stop (HTS) - VT52 is go home, but we do not support VT52
		em.setTabStop(em.getPosition().X)
	case 'M': // up one line (RI)
		em.processReverseIndex()
	case 'N': // single shift two (SS2) (TODO)
	case 'O': // single shift three (SS3) (TODO)
	case 'P':
		em.inb = em.inbStr // device control string (DCS) (TODO)
	case 'X': // start of string (SOS)
		em.inb = em.inbStr
	case 'Z': // DECID, obsolete form to get primary DA
		em.sendDA()
	case 'c': // RIS, soft reset
		em.softReset()
	case '6': // back index (DECBI, VT420, not widely supported)
		em.moveLeft()
	case '7': // save cursor (DECSC, VT100)
		em.saveCursor()
	case '8': // restore cursor (DECRC, VT100)
		em.restoreCursor()
	case '9': // forward index (DECFI, VT420, not widely supported)
		em.moveRight()
	default:
		// ESC-V and ESC-W are for guarded area (TODO)
		em.inb = em.inbInit
	}
}

// inbNF processes bytes that are part of an "nF" sequence (see ECMA-48).
func (em *emulator) inbNF(b byte) {
	if b >= 0x20 && b <= 0x2F {
		em.inBuf.WriteByte(b)
		return
	}
	if b < 0x20 || b > 0x7E { // not a valid sequence
		em.beep()
		em.inb = em.inbInit
		return
	}
	em.inBuf.WriteByte(b)
	em.inb = em.inbInit
	switch em.inBuf.String() {
	case "#8": // DECALN - fill screen with 'E'
		size := em.size
		em.autoWrap = false
		em.topMargin = 0
		em.botMargin = em.size.Y - 1
		em.ltMargin = 0
		em.rtMargin = em.size.X - 1
		em.setPosition(Coord{0, 0})
		em.style = em.style.WithAttr(Plain)
		// TODO: Reset DECOM (when we implement origin mode)
		for row := range size.Y {
			ix := em.index(Coord{X: 0, Y: row})
			for col := range size.X {
				em.cells[ix].S = em.style
				em.cells[ix].C = "E"
				em.cells[ix].W = 1
				em.be.Put(Coord{X: col, Y: row}, em.cells[ix])
				ix++
			}
		}
		// most implementations leave the cursor at home for this
		em.setPosition(Coord{0, 0})

		// case "%@": // TODO: select 8859-1
		// case "%G": // TODO: select UTF-8
		// case "(A": // TODO: select G0 as UK
		// case "(B": // TODO: select G0 as US
		// case "(C", "(5": // TODO: select G0 as Finnish
		// case "(H", "(7": // TODO: select G0 as Swedish
		// case "(K": // TODO: select G0 as German
		// case "(Q", "(9": // TODO: select G0 as French Canadian
		// case "(R", "(f": // TODO: select G0 as French
		// case "(Y": // TODO: select G0 as Italian
	}
}

// inbCSI handles bytes that are part of a CSI based sequence.
func (em *emulator) inbCSI(b byte) {
	if (b >= 0x30) && (b <= 0x3F) {
		em.inBuf.WriteByte(b) // parameter bytes
	} else if (b >= 0x20) && (b <= 0x2F) {
		em.inBuf.WriteByte(b) // intermediate bytes
	} else if b >= 0x40 && (b <= 0x7F) {
		em.inb = em.inbInit
		em.processCsi(b)
	} else {
		// error state
		em.beep()
		em.inb = em.inbInit
	}
}

// inbOSC handles bytes that are part of on OSC sequences (operating system command).
func (em *emulator) inbOSC(b byte) {
	switch b {
	case 0x9c, 0x07:
		em.inb = em.inbInit
		em.processOSC()
	case '\\':
		if buf := em.inBuf.Bytes(); len(buf) > 0 && buf[len(buf)-1] == 0x1b {
			em.inb = em.inbInit
			em.inBuf.Truncate(em.inBuf.Len() - 1)
			em.processOSC()
		} else {
			em.inBuf.WriteByte(b)
		}
	default:
		em.inBuf.WriteByte(b)
	}
}

// inbStr handles PM, SOS, and any other string we want to consume and discard.
func (em *emulator) inbStr(b byte) {
	switch b {
	case 0x9c, 0x07:
		em.inb = em.inbInit
	case '\\':
		if buf := em.inBuf.Bytes(); len(buf) > 0 && buf[len(buf)-1] == 0x1b {
			em.inb = em.inbInit
			em.inBuf.Truncate(em.inBuf.Len() - 1)
		} else {
			em.inBuf.WriteByte(b)
		}
	default:
		em.inBuf.WriteByte(b)
	}
}

// inbUTF handles continuation bytes for UTF-8 sequences.
func (em *emulator) inbUTF(b byte) {
	if b&0xC0 == 0x80 {
		// good continuation byte
		em.inBuf.WriteByte(b)
		if em.inBuf.Len() == em.utfLen {
			em.inb = em.inbInit
			r, _, err := em.inBuf.ReadRune()
			if err != nil {
				em.beep()
			} else {
				em.putRune(r)
			}
		}
	} else {
		em.beep()
		em.inb = em.inbInit
	}
}

func (em *emulator) beep() {
	if beeper, ok := em.be.(Beeper); ok {
		beeper.Beep()
	}
}

// numericParams splits the string consisting of numeric parameters into integers.
// It ensures a minimum number are present (needed for some safety cases).
// Empty strings default to zero.
func numericParams(str string, minimumLen int) ([]int, error) {
	ps := strings.Split(str, ";")
	pi := make([]int, max(len(ps), minimumLen))
	for i, str := range ps {
		if str != "" {
			if iv, err := strconv.Atoi(str); err != nil {
				return nil, err
			} else {
				pi[i] = iv
			}
		}
	}
	return pi, nil
}

// parseSgrColor grabs either 2 arguments, or 4 arguments for palette or rgb values
// used with SGR 38 and 48. The arguments must be numbers, and are returned as such.
func (em *emulator) parseSgrColor(args []string, words []string) (color.Color, []string, error) {
	if len(args) == 0 {
		if len(words) == 0 {
			return color.None, nil, errors.New("invalid color specification")
		}

		switch words[0] {
		case "2": // RGB (direct)
			if len(words) < 4 {
				return color.None, nil, errors.New("invalid color specification")
			}
			args = words[:4]
			words = words[4:]
		case "5": // palette index
			if len(words) < 2 {
				return color.None, nil, errors.New("invalid color specification")
			}
			args = words[:2]
			words = words[2:]
		default:
			return color.None, nil, errors.New("invalid color specification")
		}
	}

	switch args[0] {
	case "2": // RGB color
		if len(args) < 4 || em.be.Colors() <= 256 {
			return color.None, nil, errors.New("invalid color specification")
		}
		r, re := strconv.Atoi(args[1])
		g, ge := strconv.Atoi(args[2])
		b, be := strconv.Atoi(args[3])
		if re != nil || ge != nil || be != nil || r > 255 || g > 255 || b > 255 || r < 0 || g < 0 || b < 0 {
			return color.None, nil, errors.New("invalid color specification")
		}
		return color.NewRGBColor(int32(r), int32(g), int32(b)), words, nil
	case "5": // palette index
		if len(args) < 2 {
			return color.None, nil, errors.New("invalid color specification")
		}
		p, e := strconv.Atoi(args[1])
		if e != nil || p < 0 || p >= em.be.Colors() {
			return color.None, nil, errors.New("invalid color specification")
		}
		return color.IsValid | color.Color(p&0xff), words, nil
	}

	return color.None, nil, errors.New("invalid color specification")
}

func (em *emulator) pickColor(c color.Color, def color.Color) (color.Color, bool) {
	numColors := em.be.Colors()
	if numColors == 0 {
		return color.None, false
	}
	if c.Valid() {
		if c.IsRGB() {
			if numColors > 256 {
				return c, true
			}
		}
		if (int(c) & 255) < numColors {
			return c, true
		}
	}
	if c == color.Reset {
		return def, true
	}
	return color.None, false
}

// processSgr processes SGR commands (things that change how characters are displayed).
func (em *emulator) processSgr(str string) {
	words := strings.Split(str, ";")

	// technically parameters for 38 or 48 should be separated by colons, but due to historical
	// accident it is more common to see semicolon separation.  Underline styles are also separated
	// by a colon, if present.
	if len(words) == 0 {
		words = []string{"0"}
	}
	for len(words) > 0 {
		// we do this instead of a range so we can lop off
		// multiple words for SGR38 and 48.
		word := words[0]
		words = words[1:]

		if word == "" {
			word = "0"
		}
		args := []string(nil)
		if strings.Contains(word, ":") {
			args = strings.Split(word, ":")
			word = args[0]
			args = args[1:]
		}

		v, err := strconv.Atoi(word)
		if err != nil {
			// just swallow it for now
			return
		}
		switch v {
		case 0:
			em.style = em.defaultStyle
		case 1:
			em.style = em.style.WithAttr((em.style.Attr() &^ Dim) | Bold)
		case 2:
			em.style = em.style.WithAttr((em.style.Attr() &^ Bold) | Dim)
		case 3:
			em.style = em.style.WithAttr(em.style.Attr() | Italic)
		case 4:
			em.style = em.style.WithAttr((em.style.Attr() &^ UnderlineMask) | Underline)

			if len(args) > 0 {
				switch args[0] {
				case "2":
					em.style = em.style.WithAttr(em.style.Attr() | DoubleUnderline)
				case "3":
					em.style = em.style.WithAttr(em.style.Attr() | CurlyUnderline)
				case "4":
					em.style = em.style.WithAttr(em.style.Attr() | DottedUnderline)
				case "5":
					em.style = em.style.WithAttr(em.style.Attr() | DashedUnderline)
				}
			}
		case 5, 6:
			em.style = em.style.WithAttr(em.style.Attr() | Blink)
		case 7:
			em.style = em.style.WithAttr(em.style.Attr() | Reverse)
		case 8: // ignore, its for invisible
		case 9:
			em.style = em.style.WithAttr(em.style.Attr() | StrikeThrough)
		case 21: // Doubly underlined, per ECMA
			em.style = em.style.WithAttr((em.style.Attr() &^ UnderlineMask) | DoubleUnderline)
		case 22:
			em.style = em.style.WithAttr(em.style.Attr() &^ (Bold | Dim))
		case 23:
			em.style = em.style.WithAttr(em.style.Attr() &^ Italic)
		case 24:
			em.style = em.style.WithAttr(em.style.Attr() &^ UnderlineMask)
		case 25:
			em.style = em.style.WithAttr(em.style.Attr() &^ Blink)
		case 27:
			em.style = em.style.WithAttr(em.style.Attr() &^ Reverse)
		case 29:
			em.style = em.style.WithAttr(em.style.Attr() &^ StrikeThrough)

		case 30, 31, 32, 33, 34, 35, 36, 37: // simple foreground colors
			if c, ok := em.pickColor(color.Black+color.Color(v-30), em.defaultStyle.Fg()); ok {
				em.style = em.style.WithFg(c)
			}
		case 38:
			if c, rest, err := em.parseSgrColor(args, words); err == nil {
				words = rest
				em.style = em.style.WithFg(c)
			}
		case 39:
			if c, ok := em.pickColor(color.Reset, em.defaultStyle.Fg()); ok {
				em.style = em.style.WithFg(c)
			}
		case 40, 41, 42, 43, 44, 45, 46, 47: // simple background colors
			if c, ok := em.pickColor(color.Black+color.Color(v-40), em.defaultStyle.Bg()); ok {
				em.style = em.style.WithBg(c)
			}
		case 48:
			if c, rest, err := em.parseSgrColor(args, words); err == nil {
				words = rest
				em.style = em.style.WithBg(c)
			}
		case 49:
			if c, ok := em.pickColor(color.Reset, em.defaultStyle.Bg()); ok {
				em.style = em.style.WithBg(c)
			}
		case 53:
			em.style = em.style.WithAttr(em.style.Attr() | Overline)
		case 55:
			em.style = em.style.WithAttr(em.style.Attr() &^ Overline)
		case 58:
			if c, rest, err := em.parseSgrColor(args, words); err == nil {
				words = rest
				em.style = em.style.WithUc(c)
			}
		}
	}
}

// processCursorUp implements CUU.
func (em *emulator) processCursorUp(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		em.moveUpN(Row(max(1, pi[0])))
	}
}

// processCursorDown implements CUD.
func (em *emulator) processCursorDown(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		em.moveDownN(Row(max(1, pi[0])))
	}
}

// processCursorForward implements CUF.
func (em *emulator) processCursorForward(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		em.moveRightN(Col(max(1, pi[0])))
	}
}

// processCursorBackward implements CUB.
func (em *emulator) processCursorBackward(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		em.moveLeftN(Col(max(1, pi[0])))
	}
}

// processCursorNextLine implements CNL.
func (em *emulator) processCursorNextLine(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		em.moveDownN(Row(max(1, pi[0])))
		pos := em.getPosition()
		pos.X = 0
		em.setPosition(pos)
	}
}

// processCursorPreviousLine implements CPL.
func (em *emulator) processCursorPreviousLine(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		em.moveUpN(Row(max(1, pi[0])))
		pos := em.getPosition()
		pos.X = 0
		em.setPosition(pos)
	}
}

// processCursorColumn implements CHA.
func (em *emulator) processCursorColumn(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		pos := em.getPosition()
		// TODO: possibly clip to margins (origin mode)
		pos.X = min(Col(max(1, pi[0])), em.size.X) - 1
		em.setPosition(pos)
	}
}

// processCursorPosition implements CUP, and also HVP.
func (em *emulator) processCursorPosition(str string) {
	if pi, err := numericParams(str, 2); err == nil {
		em.autoWrap = false
		pos := em.getPosition()
		wsz := em.size
		row := Row(max(1, pi[0]))
		col := Col(max(1, pi[1]))
		row = max(1, min(row, wsz.Y))
		col = max(1, min(col, wsz.X))
		pos.X = col - 1
		pos.Y = row - 1
		em.setPosition(pos)
	}
}

// processCursorTab implements CHT.
func (em *emulator) processCursorTab(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		// Note: tab does not clear this field.
		for range max(1, pi[0]) {
			em.nextTab()
		}
	}
}

// processCursorBackTab implements CBT.
func (em *emulator) processCursorBackTab(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		for range max(1, pi[0]) {
			em.prevTab()
		}
	}
}

// processEraseDisplay implements ED.
func (em *emulator) processEraseDisplay(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.bufferingStart()
		defer em.bufferingEnd()

		switch pi[0] {
		case 0: // erase below
			em.eraseBelow()
		case 1: // erase above
			em.eraseAbove()
		case 2: // erase all
			em.eraseAll()
			// others not supported (3 is erase saved lines)
		}
	}
}

// processEraseLine implements EL.
func (em *emulator) processEraseLine(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.bufferingStart()
		defer em.bufferingEnd()

		switch pi[0] {
		case 0:
			em.eraseToLineEnd()
		case 1:
			em.eraseToLineStart()
		case 2:
			em.eraseLine()
		}
	}
}

// processEraseCharacter implements ECH.
// This ignores the margin.
func (em *emulator) processEraseCharacter(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		pos := em.pos
		// TODO: delete wide character if we are splitting it at the start
		em.bufferingStart()
		defer em.bufferingEnd()
		for range max(1, pi[0]) {

			em.eraseCell(pos)
			pos.X++
			if pos.X >= em.size.X {
				break
			}
		}
	}
}

// processScrollUp implements SU (VT420.)
func (em *emulator) processScrollUp(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.bufferingStart()
		defer em.bufferingEnd()
		for range max(pi[0], 1) {
			// TODO: consider faster jump scroll.
			// This should be something tunable as well.
			em.scrollUp()
		}
	}
}

// processScrollDown implements SD (VT420.)
func (em *emulator) processScrollDown(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.bufferingStart()
		defer em.bufferingEnd()
		for range max(pi[0], 1) {
			// TODO: consider faster jump scroll.
			// This should be something tunable as well.
			em.scrollDown()
		}
	}
}

// processWindowOps handles CSI ... t window operations.
func (em *emulator) processWindowOps(str string) {
	if pi, err := numericParams(str, 3); err == nil {
		switch pi[0] {
		case 8: // Resize window: CSI 8 ; rows ; cols t
			rows := pi[1]
			cols := pi[2]
			if rows < 1 || cols < 1 {
				return
			}
			size := Coord{X: Col(cols), Y: Row(rows)}
			if ws, ok := em.be.(interface{ SetSize(Coord) }); ok {
				ws.SetSize(size)
				em.applyResize(size)
			}

		case 18: // Report text area size: CSI 8 ; rows ; cols t
			em.SendRaw(fmt.Appendf(nil, "\x1b[8;%d;%dt", em.size.Y, em.size.X))
		}
	}
}

// processVerticalMargins implements DECSTBM (set top and bottom margins, VT220.)
func (em *emulator) processVerticalMargins(str string) {
	if pi, err := numericParams(str, 2); err == nil {
		pi[0] = max(1, pi[0])
		if pi[1] == 0 {
			pi[1] = int(em.size.Y)
		}
		pi[0]--
		pi[1]--

		top := min(Row(pi[0]), em.size.Y-1)
		bot := min(Row(pi[1]), em.size.Y-1)

		// no change if values are out of range
		if bot > top {
			em.topMargin = top
			em.botMargin = bot
			em.setPosition(Coord{X: 0, Y: 0})
		}
	}
}

// processHorizontalMargins implements DECSLRM (set left and right margins, VT400.)
// It only works if Private Mode 69 (Left and Right margins)
func (em *emulator) processHorizontalMargins(str string) {
	if em.getPrivateMode(PmLeftRightMargin) != ModeOn {
		// For compat with SCO and ANSI.SYS.
		if str == "" {
			em.saveCursor()
		}
		return
	}
	if pi, err := numericParams(str, 2); err == nil {
		pi[0] = max(1, pi[0])
		if pi[1] == 0 {
			pi[1] = int(em.size.X)
		}
		pi[0]--
		pi[1]--

		lm := min(Col(pi[0]), em.size.X-1)
		rm := min(Col(pi[1]), em.size.X-1)

		// no change if values are out of range
		if rm > lm {
			em.rtMargin = rm
			em.ltMargin = lm
			em.setPosition(Coord{X: 0, Y: 0})
		}
	}
}

// processIndex moves down, unless already on the bottom margin, in which case it scrolls Up.
func (em *emulator) processIndex() {
	pos := em.getPosition()
	em.autoWrap = false
	if pos.Y == em.botMargin && em.ltMargin <= pos.X && pos.X <= em.rtMargin {
		em.scrollUp()
	} else {
		em.processCursorDown("")
	}
}

// processCarriageReturn handle CR.
func (em *emulator) processCarriageReturn() {
	em.lastIndex = 0
	em.setPosition(Coord{0, em.getPosition().Y})
}

// processLineFeed is like IND, but if ANSI mode 20 is set, then a CR is appended as well.
func (em *emulator) processLineFeed() {
	em.lastIndex = 0
	em.processIndex()
	if em.getAnsiMode(AmNewLineMode) == ModeOn {
		em.processCarriageReturn()
	}
}

// processReverseIndex moves up, unless already on the top margin, in which case it scrolls down.
func (em *emulator) processReverseIndex() {
	pos := em.getPosition()
	if pos.Y == em.topMargin {
		em.scrollDown()
	} else {
		em.processCursorUp("")
	}
}

// processCursorRow implements VPA (set vertical position absolute)
func (em *emulator) processCursorRow(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		pos := em.getPosition()
		pos.Y = min(Row(max(1, pi[0])), em.size.Y) - 1
		em.setPosition(pos)
	}
}

// processCursorRowAdvance implements VPR (vertical position relative).
func (em *emulator) processCursorRowAdvance(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		pos := em.getPosition()
		pos.Y = min(pos.Y+Row(max(1, pi[0])), em.size.Y-1)
		em.setPosition(pos)
	}
}

// processInsertLine implements IL.
func (em *emulator) processInsertLine(str string) {
	// insert line only takes effect within the scrolling region
	if em.pos.X < em.ltMargin || em.pos.X > em.rtMargin ||
		em.pos.Y < em.topMargin || em.pos.Y > em.botMargin {
		return
	}

	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		num := Row(max(1, pi[0]))
		num = min(num, em.botMargin-em.pos.Y+1)

		// process as a scroll down at our position.
		top := em.topMargin
		em.topMargin = em.pos.Y
		for range num {
			em.scrollDown()
		}
		em.topMargin = top
		em.pos.X = 0
		em.setPosition(em.pos)
	}
}

// processDeleteLine implements DL.
func (em *emulator) processDeleteLine(str string) {
	// delete line only takes effect within the scrolling region
	if em.pos.X < em.ltMargin || em.pos.X > em.rtMargin ||
		em.pos.Y < em.topMargin || em.pos.Y > em.botMargin {
		return
	}
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		num := Row(max(1, pi[0]))
		num = min(num, em.botMargin-em.pos.Y+1)
		if num < 1 {
			return
		}
		// process as a scroll up at our position.
		top := em.topMargin
		em.topMargin = em.pos.Y
		for range num {
			em.scrollUp()
		}
		em.topMargin = top
		em.pos.X = 0
		em.setPosition(em.pos)
	}
}

// processDeleteCharacter implements DCH.
func (em *emulator) processDeleteCharacter(str string) {
	// only takes effect within the scrolling region
	if em.pos.X < em.ltMargin || em.pos.X > em.rtMargin ||
		em.pos.Y < em.topMargin || em.pos.Y > em.botMargin {
		return
	}
	em.bufferingStart()
	defer em.bufferingEnd()

	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		num := Col(max(1, pi[0]))
		num = min(num, em.rtMargin-em.pos.X+1)
		if num < 1 {
			return
		}
		// this is essentially a one line scroll left
		pos := em.pos

		// if we are breaking a wide rune, delete it (but preserve style)
		if em.pos.X > 0 {
			if ix := em.index(em.pos); em.cells[ix-1].W > 1 {
				em.cells[ix-1].C = ""
				em.cells[ix-1].W = 0
				em.be.Put(Coord{X: em.pos.X - 1, Y: em.pos.Y}, em.cells[ix-1])
			}
		}

		for range num {
			src := Coord{X: em.pos.X + 1, Y: em.pos.Y}
			dst := Coord{X: em.pos.X, Y: em.pos.Y}
			dim := Coord{X: em.rtMargin - em.pos.X, Y: 1}
			em.blit(src, dst, dim)
			em.eraseCell(Coord{X: em.rtMargin, Y: em.pos.Y})
		}
		em.pos = pos
		em.setPosition(em.pos)
	}
}

// processInsertCharacter implements ICH.
func (em *emulator) processInsertCharacter(str string) {
	em.autoWrap = false
	// only takes effect within the scrolling region -- HOWEVER,
	// ICH still resets auto-wrap in this case, unlike DCH.
	if em.pos.X < em.ltMargin || em.pos.X > em.rtMargin ||
		em.pos.Y < em.topMargin || em.pos.Y > em.botMargin {
		return
	}
	em.bufferingStart()
	defer em.bufferingEnd()
	if pi, err := numericParams(str, 1); err == nil {
		num := Col(max(1, pi[0]))
		num = min(num, em.rtMargin-em.pos.X+1)
		if num < 1 {
			return
		}

		// this is essentially a one line scroll right
		pos := em.pos

		// if we are breaking a wide rune, delete it
		if em.pos.X > 0 {
			if ix := em.index(em.pos); em.cells[ix-1].W > 1 {
				em.cells[ix-1].C = ""
				em.cells[ix-1].W = 0
				em.be.Put(Coord{X: em.pos.X - 1, Y: em.pos.Y}, em.cells[ix-1])
			}
		}

		for range num {
			src := Coord{X: em.pos.X, Y: em.pos.Y}
			dst := Coord{X: em.pos.X + 1, Y: em.pos.Y}
			dim := Coord{X: em.rtMargin - em.pos.X, Y: 1}
			em.blit(src, dst, dim)
			ix := em.index(Coord{X: em.pos.X, Y: em.pos.Y})
			em.cells[ix].C = ""
			em.cells[ix].W = 0
			em.cells[ix].S = em.style
			// NB: We don't use eraseCell, because we need to preserve attributes.
			em.be.Put(Coord{X: em.pos.X, Y: em.pos.Y}, em.cells[ix])
		}

		// if we clipped off the end of a wide character, then delete it.
		if ix := em.index(Coord{X: em.rtMargin, Y: em.pos.Y}); em.cells[ix].W > 1 {
			em.cells[ix].C = ""
			em.cells[ix].W = 0
			em.be.Put(Coord{X: em.rtMargin, Y: em.pos.Y}, em.cells[ix])
		}

		em.pos = pos
		em.setPosition(em.pos)
	}
}

// processSetMode implements SM (set ANSI mode).
func (em *emulator) processSetMode(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		for _, pm := range pi {
			em.setAnsiMode(AnsiMode(pm), ModeOn)
		}
	}
}

// processResetMode implements RM (reset ANSI mode).
func (em *emulator) processResetMode(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		for _, pm := range pi {
			em.setAnsiMode(AnsiMode(pm), ModeOff)
		}
	}
}

// processRequestMode implements DECRQM for ANSI modes.
// Only a single numeric parameter (mode number) can be supplied (VT300+)
func (em *emulator) processRequestMode(str string) {
	if m, err := strconv.Atoi(str); err == nil {
		am := AnsiMode(m)
		status := em.getAnsiMode(am)
		em.SendRaw([]byte(am.Reply(status)))
	}
}

// processSetPrivateMode implements DECSET (set private mode).
func (em *emulator) processSetPrivateMode(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		for _, pm := range pi {
			em.setPrivateMode(PrivateMode(pm), ModeOn)
		}
	}
}

// processResetPrivateMode implements DECRST (reset private mode).
func (em *emulator) processResetPrivateMode(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		for _, pm := range pi {
			em.setPrivateMode(PrivateMode(pm), ModeOff)
		}
	}
}

// processRequestPrivateMode implements DECRQM for private modes.
// Only a single numeric parameter (mode number) can be supplied (VT300+)
func (em *emulator) processRequestPrivateMode(str string) {
	if m, err := strconv.Atoi(str); err == nil {
		pm := PrivateMode(m)
		status := em.getPrivateMode(pm)
		em.SendRaw([]byte(pm.Reply(status)))
	}
}

// processTabReset implements DECST8C (set tab stops to every 8 chars)
func (em *emulator) processTabReset(str string) {
	if pi, err := numericParams(str, 1); err == nil && pi[0] == 5 {
		em.tabStops = nil
	}
}

// processTabClear implements TBC (clear horizontal tab).
func (em *emulator) processTabClear(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		switch pi[0] {
		case 0: // clear stop at current column
			em.clrTabStop(em.getPosition().X)
		case 3: // clear all columns
			em.tabStops = []Col{} // this is distinct from nil
		}
	}
}

// processPrimaryAttributes implements send DA.
func (em *emulator) processPrimaryAttributes(str string) {
	if pi, err := numericParams(str, 1); err == nil && pi[0] == 0 {
		em.sendDA()
	}
}

// processExtendedAttributes implements XTVERSION (send terminal name and version).
func (em *emulator) processExtendedAttributes(str string) {
	if pi, err := numericParams(str, 1); err == nil && pi[0] == 0 && em.name != "" {
		em.SendRaw(fmt.Appendf(nil, "\x1bP>|%s %s\x1b\\", em.name, em.vers))
	}
}

// processCursorStyle implements DECSCUSR (set cursor style).
func (em *emulator) processCursorStyle(str string) {
	// get previous visibility state, as we don't change it with this call.
	visible := em.cursor.IsVisible()
	if pi, err := numericParams(str, 1); err == nil {
		switch pi[0] {
		case 0, 1:
			em.cursor = BlinkingBlock
		case 2:
			em.cursor = SteadyBlock
		case 3:
			em.cursor = BlinkingUnderline
		case 4:
			em.cursor = SteadyUnderline
		case 5:
			em.cursor = BlinkingBar
		case 6:
			em.cursor = SteadyBar
		}
		if !visible {
			em.cursor = em.cursor.Hide()
		}
		em.be.SetCursor(em.cursor)
	}
}

// processCsi processes CSI sequences.
func (em *emulator) processCsi(final byte) {

	// CSI sequences are supported in several different possible ways:
	// parameters may have a prefix character that is not numeric, typically
	// indicating a whole different mode of operation than the final byte.
	// There may also be intermediate bytes, but we only look for one, because
	// the use cases we have this are that only a single intermediate byte is
	// sometimes used to affect function.  (E.g. $ in some cases.)
	cmd := ""
	if em.inBuf.Len() > 0 {
		if b := em.inBuf.Bytes()[0]; b > '9' && b <= '?' {
			cmd += string(b)
			em.inBuf.ReadByte()
		}
	}
	if l := em.inBuf.Len(); l > 0 {
		if b := em.inBuf.Bytes()[l-1]; b >= 0x20 && b <= 0x2F {
			cmd += string(b)
			em.inBuf.Truncate(l - 1)
		}
	}
	cmd += string(final)

	str := em.inBuf.String()
	switch cmd {

	case "@":
		em.processInsertCharacter(str)
	case "A":
		em.processCursorUp(str)
	case "B":
		em.processCursorDown(str)
	case "C":
		em.processCursorForward(str)
	case "D":
		em.processCursorBackward(str)
	case "E":
		em.processCursorNextLine(str)
	case "F":
		em.processCursorPreviousLine(str)
	case "G":
		em.processCursorColumn(str)
	case "H", "f":
		em.processCursorPosition(str)
	case "I":
		em.processCursorTab(str)
	case "J":
		em.processEraseDisplay(str)
	case "K":
		em.processEraseLine(str)
	case "L":
		em.processInsertLine(str)
	case "M":
		em.processDeleteLine(str)
	case "P":
		em.processDeleteCharacter(str)
	case "S":
		em.processScrollUp(str)
	case "T":
		em.processScrollDown(str)
	case "X":
		em.processEraseCharacter(str)
	case "Z":
		em.processCursorBackTab(str)
	case "c":
		em.processPrimaryAttributes(str)
	case "d":
		em.processCursorRow(str)
	case "e":
		em.processCursorRowAdvance(str)
	case "g":
		em.processTabClear(str)
	case "h":
		em.processSetMode(str)
	case "l":
		em.processResetMode(str)
	case "m":
		em.processSgr(str)
	case "n":
		em.deviceReport(str)
	case "r":
		em.processVerticalMargins(str)
	case "s":
		em.processHorizontalMargins(str)
	case "t":
		em.processWindowOps(str)
	case " q":
		em.processCursorStyle(str)
	case "?W":
		em.processTabReset(str)
	case "?h":
		em.processSetPrivateMode(str)
	case "?l":
		em.processResetPrivateMode(str)
	case "$p":
		em.processRequestMode(str)
	case "?$p":
		em.processRequestPrivateMode(str)
	case ">q":
		em.processExtendedAttributes(str)
	}
}

// processClipboard handles OSC 52 commands.
func (em *emulator) processClipboard(str string) {
	clipper, ok := em.be.(Clipboard)
	if !ok {
		return
	}

	// first parameter is the target.  We only have a single
	// target, and alias all possibilities to the same.
	parts := strings.SplitN(str, ";", 2)
	if len(parts) != 2 {
		return
	}

	if parts[1] == "?" {
		// request for clipboard content
		data := clipper.GetClipboard()
		if data != nil {
			em.SendRaw(fmt.Appendf(nil, "\x1b]52;c;%s\x1b\\", base64.StdEncoding.EncodeToString(data)))
		}
		return
	}

	buf := make([]byte, base64.StdEncoding.DecodedLen(len(parts[1])))
	if n, err := base64.StdEncoding.Decode(buf, []byte(parts[1])); err == nil {
		clipper.SetClipboard(buf[:n])
		return
	}

	clipper.SetClipboard([]byte{})
}

// processHyperLink handles OSC 8 commands.
func (em *emulator) processHyperLink(str string) {
	// format is params;URI params are colon separated key value pairs.
	// if the URI is absent, then the link is terminated.
	parts := strings.SplitN(str, ";", 2)
	if len(parts) == 2 {
		if parts[1] == "" {
			// No URI
			em.style = em.style.WithUrl("", "")
			return
		}
		url := parts[1]
		id := ""
		for pair := range strings.SplitSeq(parts[0], ":") {
			if val, ok := strings.CutPrefix(pair, "id="); ok {
				id = val
			}
		}
		em.style = em.style.WithUrl(url, id)
	}
}

// processOSC processes an operating system command.
func (em *emulator) processOSC() {

	// Every OSC we support has a number, semicolon, then string.
	ns, str, ok := strings.Cut(em.inBuf.String(), ";")
	if !ok {
		return
	}
	if num, err := strconv.Atoi(ns); err != nil {
		return
	} else {
		switch num {
		case 2: // Set window title
			if t, ok := em.be.(Titler); ok {
				// TODO: possibly validate the UTF-8 content?
				em.bufferingStart()
				defer em.bufferingEnd()
				t.SetWindowTitle(str)
			}
		case 8:
			em.processHyperLink(str)
		case 52:
			em.processClipboard(str)
		}
	}
}

func (em *emulator) getPosition() Coord {
	pos := em.be.GetPosition()
	em.pos = pos
	return em.pos
}

func (em *emulator) setPosition(pos Coord) {
	em.pos = pos
	em.be.SetPosition(pos)
}

func (em *emulator) deviceReport(s string) {
	switch s {
	case "5":
		em.SendRaw([]byte("\x1b[0n"))
	case "6":
		pos := em.getPosition()
		em.SendRaw(fmt.Appendf(nil, "\x1b[%d;%dR", pos.Y+1, pos.X+1))
	default: // ignore
	}
}

func (em *emulator) moveUpN(count Row) {
	for range count {
		em.moveUp()
	}
}

func (em *emulator) moveDownN(count Row) {
	for range count {
		em.moveDown()
	}
}

func (em *emulator) moveLeftN(count Col) {
	for range count {
		em.moveLeft()
	}
}

func (em *emulator) moveRightN(count Col) {
	for range count {
		em.moveRight()
	}
}

// moveDown moves down, to the limit of either bottom margin, or the bottom of the screen if outside the margin.
func (em *emulator) moveDown() {
	pos := em.getPosition()
	win := em.size
	if pos.Y == em.botMargin || pos.Y == win.Y-1 {
		return
	}
	pos.Y++
	em.setPosition(pos)
}

// moveUp moves up, to the limit of either top margin, or zero if outside the margin.
func (em *emulator) moveUp() {
	pos := em.getPosition()
	if pos.Y == 0 || pos.Y == em.topMargin {
		return
	}
	pos.Y--
	em.setPosition(pos)
}

func (em *emulator) moveLeft() {
	em.autoWrap = false
	em.lastIndex = 0
	pos := em.getPosition()
	if pos.X > 0 {
		pos.X--
		em.setPosition(pos)
	}
}

func (em *emulator) moveRight() {
	pos := em.getPosition()
	win := em.size
	if pos.X < win.X-1 {
		pos.X++
		em.setPosition(pos)
	}
}

// nextLine is like CNL with 1, but it optionally also scrolls.
func (em *emulator) nextLine() {
	em.autoWrap = false
	if em.pos.Y == em.botMargin {
		em.scrollUp()
	}
	em.moveDown()
	em.pos.X = 0
	em.setPosition(em.pos)
}

// blit performs a data move operation.  It does ignores margins.
func (em *emulator) blit(src, dst, dim Coord) {

	em.bufferingStart()
	defer em.bufferingEnd()

	// save the source and destination for the backend blit
	bsrc := src
	bdst := dst

	lim := em.size

	// clip to visible source
	if dim.X+src.X > lim.X {
		dim.X = lim.X - src.X
	}
	if dim.Y+src.Y > lim.Y {
		dim.Y = lim.Y - src.Y
	}
	// and clip to final destination
	if dim.X+dst.X > lim.X {
		dim.X = lim.X - dst.X
	}
	if dim.Y+dst.Y > lim.Y {
		dim.Y = lim.Y - dst.Y
	}

	// gap represents decrement when shifting to the next row --
	// skipping over the irrelevant cells. (The increment in the
	// index when going from last cell of row to first cell of next row,
	// or vice versa.)
	gap := int(lim.X - dim.X)

	// the following logic is carefully constructed to avoid expensive
	// operations in the loops (only addition or subtraction)
	if em.index(src) > em.index(dst) { // source appears later, so we can forward copy
		si := em.index(src)
		di := em.index(dst)
		for range dim.Y {
			for range dim.X {
				em.cells[di] = em.cells[si]
				di++
				si++
			}
			// advance to next row
			si += gap
			di += gap
		}
	} else { // source appears earlier, so we have to reverse copy
		src.Y += dim.Y - 1
		dst.Y += dim.Y - 1
		src.X += dim.X - 1
		dst.X += dim.X - 1
		si := em.index(src)
		di := em.index(dst)

		for range dim.Y {
			for range dim.X {
				em.cells[di] = em.cells[si]
				si--
				di--
			}
			si -= gap
			di -= gap
		}
	}

	// Now we possibly blit underneath.  We'll use the underlying
	// implementation's blit operation if it has one, else we'll
	// just rewrite the cells in linear order.
	if b, ok := em.be.(Blitter); ok {
		// The backend implements what should be a fast blit.
		b.Blit(bsrc, bdst, dim)
	} else {
		// This does math for each cell, so you're looking at a lot of multiplications
		// for a bit display -- 100x100 means 10,000x2 multiplications.  This could be
		// optimized, but this is really a fallback as every backend really *should* have
		// an efficient blit operation. (The ones that don't probably don't keep their own
		// state, such as a wrappers on top of TTYs.)  In these cases the cost of writing
		// the content is probably substantially dominant anyway.
		for row := range dim.Y {
			for col := range dim.X {
				pos := bdst
				pos.X += col
				pos.Y += row
				cell := em.cells[em.index(pos)]
				em.be.Put(pos, cell)
			}
		}
	}
}

func (em *emulator) scrollUp() {
	dim := Coord{X: em.rtMargin - em.ltMargin + 1, Y: em.botMargin - em.topMargin}
	src := Coord{X: em.ltMargin, Y: em.topMargin + 1}
	dst := Coord{X: em.ltMargin, Y: em.topMargin}

	// TODO: deal with wide characters broken across the margin
	pos := em.pos
	em.blit(src, dst, dim)
	bot := Coord{X: em.ltMargin, Y: em.botMargin}
	for bot.X <= em.rtMargin {
		em.eraseCell(bot)
		bot.X++
	}
	em.setPosition(pos)
}

func (em *emulator) scrollDown() {
	dim := Coord{X: em.rtMargin - em.ltMargin + 1, Y: em.botMargin - em.topMargin}
	src := Coord{X: em.ltMargin, Y: em.topMargin}
	dst := Coord{X: em.ltMargin, Y: em.topMargin + 1}

	pos := em.pos
	em.blit(src, dst, dim)
	em.setPosition(Coord{X: em.ltMargin, Y: em.topMargin})
	top := Coord{X: em.ltMargin, Y: em.topMargin}
	for top.X <= em.rtMargin {
		em.eraseCell(top)
		top.X++
	}
	em.setPosition(pos)
}

// nextTab advances to the next tab stop, or the end of
// the line if there is no further tab.
func (em *emulator) nextTab() {
	em.lastIndex = 0
	maxX := em.size.X - 1
	curX := em.getPosition().X
	if curX == maxX { // already at end
		return
	}
	nextX := maxX
	if em.tabStops == nil {
		// just advance to the next one
		nextX = min((curX+8)&^7, maxX)
	} else {
		for _, p := range em.tabStops {
			if p > curX {
				nextX = p
				break
			}
		}
	}
	em.setPosition(Coord{X: nextX, Y: em.pos.Y})
}

func (em *emulator) prevTab() {
	curX := em.getPosition().X
	if curX == 0 {
		return
	}
	nextX := curX - 1
	if em.tabStops == nil {
		nextX &^= 7
	} else if i, exist := slices.BinarySearch(em.tabStops, nextX); exist {
		nextX = em.tabStops[i]
	} else if i > 0 {
		nextX = em.tabStops[i-1]
	} else {
		nextX = 0
	}
	em.setPosition(Coord{X: nextX, Y: em.pos.Y})
}

// initTabStops initializes the tab stops assuming every 8th column
// is a tab stop.  This should only be called if the user is intentionally
// changing the tab stops, because it will no longer support expanding
// tab stops on resizing.
func (em *emulator) initTabStops() {
	if em.tabStops == nil {
		// no tab stop at offset 0 since that would be pointless
		em.tabStops = make([]Col, 0, int(em.size.X/8)+1)
		for col := Col(8); col < em.size.X; col += 8 {
			em.tabStops = append(em.tabStops, col)
		}
	}
}

// setTabStop sets a tab stop at the given location.
// This calls  initTabStops - please see the description of that function for ramifications.
func (em *emulator) setTabStop(ts Col) {
	em.initTabStops()
	if index, exist := slices.BinarySearch(em.tabStops, ts); !exist {
		em.tabStops = slices.Insert(em.tabStops, index, ts)
	}
}

// clrTabStop clears the tab stop at the given column.  This calls
// initTabStops - please see the description of that function for ramifications.
func (em *emulator) clrTabStop(ts Col) {
	em.initTabStops()
	em.tabStops = slices.DeleteFunc(em.tabStops, func(x Col) bool { return x == ts })
}

// index obtains the index in the cells slice for the given coordinates,
// which must be within the bounds of the display size.
func (em *emulator) index(c Coord) int {
	return int(c.Y)*int(em.size.X) + int(c.X)
}

// putRune puts out a single rune.  This might be a subsequent part of a grapheme cluster, in
// which case it will be emitted together with the preceding base character.
func (em *emulator) putRune(r rune) {
	dim := em.size

	if lastIdx := em.lastIndex; lastIdx != 0 {
		lastIdx--
		if pm := em.getPrivateMode(PmGraphemeClusters); pm == ModeOn || pm == ModeOnLocked {
			// ASCII-to-ASCII pairs cannot extend a grapheme cluster, except CRLF.
			prev := em.cells[lastIdx].C
			if len(prev) == 1 && prev[0] < utf8.RuneSelf && !shouldCheckGrapheme(prev[0], r) {
				// fall through to the normal single-rune path
			} else {
				// maybe we need to update the last index
				buf := em.graphemeBuf[:0]
				need := len(prev) + utf8.UTFMax
				if cap(buf) < need {
					buf = make([]byte, 0, need)
				}
				buf = append(buf, prev...)
				buf = utf8.AppendRune(buf, r)
				em.graphemeIter.SetText(buf)
				if em.graphemeIter.Next() && len(em.graphemeIter.Value()) == len(buf) {
					// we are adding to a cluster
					cluster := em.graphemeIter.Value()
					width := em.cells[lastIdx].W
					if w := textWidthOptions.Rune(r); w > width {
						width = w
					}
					if isRegionalIndicator(r) && width < 2 {
						width = 2
					}
					if r == '\uFE0F' && width < 2 {
						width = 2
					}
					em.cells[lastIdx].C = em.clusterString(cluster)
					em.cells[lastIdx].W = width
					col := Col(lastIdx) % dim.X
					row := Row(lastIdx / int(dim.X))
					// we may have to move position if this switches to wide, so recalculate expected end
					next := col + Col(width)
					if em.getPrivateMode(PmAutoMargin) == ModeOn && next >= dim.X {
						em.autoWrap = true
					}
					end := next
					if end >= dim.X {
						end = dim.X - 1
					}
					if width == 2 && col < dim.X-1 && em.cells[lastIdx+1].W != 0 {
						// erase the next cell before putting down a character
						em.cells[lastIdx+1].C = ""
						em.cells[lastIdx+1].S = em.cells[lastIdx].S
						em.cells[lastIdx+1].W = 0
						em.be.Put(Coord{X: col + 1, Y: row}, em.cells[lastIdx+1])
					}
					// we leave the em.lastIndex for now, we might keep extending this cluster
					em.be.Put(Coord{X: col, Y: row}, em.cells[lastIdx])
					em.setPosition(Coord{X: end, Y: row})
					em.graphemeBuf = buf[:0]
					return
				}
				em.graphemeBuf = buf[:0]
			}
		}
	}

	if em.autoWrap {
		em.nextLine()
	}

	autoMargin := em.getPrivateMode(PmAutoMargin) == ModeOn

	pos := em.getPosition()
	w := textWidthOptions.Rune(r)
	if autoMargin && pos.X+Col(w) >= dim.X {
		em.autoWrap = true
	}
	index := em.index(pos)
	em.cells[index].C = em.runeString(r)
	em.cells[index].S = em.style
	em.cells[index].W = w
	em.be.Put(em.pos, em.cells[index])
	em.lastIndex = index + 1

	if w == 2 && pos.X < dim.X-1 {
		index++
		em.cells[index].C = ""
		em.cells[index].S = em.style
		em.cells[index].W = 0
	}
	// Advance the cursor. This will stop at the margin.
	// Note that if auto margin is enabled, we will have set
	// autoWrap above if we were at the margin already.
	em.moveRightN(Col(w))
}

func (em *emulator) runeString(r rune) string {
	if r < utf8.RuneSelf {
		return asciiRuneStrings[r]
	}
	return em.runeStrings.stringFor(r)
}

func (em *emulator) clusterString(cluster []byte) string {
	return em.clusterStrings.stringFor(cluster)
}

func shouldCheckGrapheme(prev byte, r rune) bool {
	if r < utf8.RuneSelf {
		return prev == '\r' && r == '\n'
	}

	if unicode.Is(unicode.M, r) {
		return true
	}
	if r == '\u200d' {
		return true
	}
	if r >= 0xFE00 && r <= 0xFE0F {
		return true
	}
	if r >= 0xE0100 && r <= 0xE01EF {
		return true
	}

	return false
}

func isRegionalIndicator(r rune) bool {
	return r >= 0x1F1E6 && r <= 0x1F1FF
}

// eraseCell erases a single cell at the given offset.
// It clears attributes, but leaves the colors intact.
func (em *emulator) eraseCell(c Coord) {
	s := em.style.WithAttr(Plain)
	index := em.index(c)
	em.cells[index].C = ""
	em.cells[index].S = s
	em.cells[index].W = 0
	em.be.Put(c, em.cells[index])
}

// eraseBelow erases from (and including) the current cursor position to the end of the window.
func (em *emulator) eraseBelow() {
	size := em.size
	pos := em.getPosition()
	for x := pos.X; x < size.X; x++ {
		em.eraseCell(Coord{X: x, Y: pos.Y})
	}
	for y := pos.Y + 1; y < size.Y; y++ {
		for x := Col(0); x < size.X; x++ {
			em.eraseCell(Coord{X: x, Y: y})
		}
	}
	em.setPosition(pos)
}

// eraseAbove erases from the origin to (and including) the current cursor position.
func (em *emulator) eraseAbove() {
	size := em.size
	pos := em.getPosition()
	for y := Row(0); y < pos.Y; y++ {
		for x := Col(0); x < size.X; x++ {
			em.eraseCell(Coord{X: x, Y: y})
		}
	}
	for x := Col(0); x <= pos.X; x++ {
		em.eraseCell(Coord{X: x, Y: pos.Y})
	}
	em.setPosition(pos)
}

// eraseAll erases the entire screen. It uses the color, but resets all other attributes.
func (em *emulator) eraseAll() {
	size := em.size
	pos := em.getPosition()
	for y := Row(0); y < size.Y; y++ {
		for x := Col(0); x < size.X; x++ {
			em.eraseCell(Coord{X: x, Y: y})
		}
	}
	em.setPosition(pos)
}

// eraseToLineEnd erases to the end of the line, including the cursor position.
func (em *emulator) eraseToLineEnd() {
	size := em.size
	pos := em.getPosition()
	for x := pos.X; x < size.X; x++ {
		em.eraseCell(Coord{x, pos.Y})
	}
	em.setPosition(pos)
}

// eraseToLineStart erases to the start of the line, including the cursor position.
func (em *emulator) eraseToLineStart() {
	pos := em.getPosition()
	for x := Col(0); x <= pos.X; x++ {
		em.eraseCell(Coord{x, pos.Y})
	}
	em.setPosition(pos)
}

// eraseLine erases the entire line.
func (em *emulator) eraseLine() {
	size := em.size
	pos := em.getPosition()
	for x := range size.X {
		em.eraseCell(Coord{x, pos.Y})
	}
	em.setPosition(pos)
}

// softReset performs a soft reset.
func (em *emulator) softReset() {
	// TODO:
	// Select default character sets
	em.tabStops = nil
	em.autoWrap = false
	em.style = em.defaultStyle
	em.saved = savedCursor{style: em.defaultStyle}
	em.topMargin = 0
	em.botMargin = em.size.Y - 1
	em.ltMargin = 0
	em.rtMargin = em.size.X - 1
	em.appKeyPad = false
	em.be.Reset()
	// start by resetting all modes
	for _, am := range em.ansiModeKeys() {
		em.setAnsiMode(am, ModeOff) // NB: No effect for non-changeable modes
	}
	for _, pm := range em.privateModeKeys() {
		em.setPrivateMode(pm, ModeOff) // NB: No effect for non-changeable modes
	}
	// and set any that should reset on (auto-margin)
	em.setPrivateMode(PmAutoMargin, ModeOn)
	em.setPrivateMode(PmAutoRepeat, ModeOn)
	em.setPrivateMode(PmShowCursor, ModeOn)
	em.setPrivateMode(PmBlinkCursor, ModeOn)
	// set default cursor - matches VT defaults
	em.cursor = BlinkingBlock
	em.be.SetCursor(em.cursor)
	em.setPosition(Coord{0, 0})
	em.eraseAll()
}

func (em *emulator) ansiModeKeys() []AnsiMode {
	em.modeLock.RLock()
	defer em.modeLock.RUnlock()

	keys := make([]AnsiMode, 0, len(em.ansiModes))
	for am := range em.ansiModes {
		keys = append(keys, am)
	}
	return keys
}

func (em *emulator) privateModeKeys() []PrivateMode {
	em.modeLock.RLock()
	defer em.modeLock.RUnlock()

	keys := make([]PrivateMode, 0, len(em.localModes))
	for pm := range em.localModes {
		keys = append(keys, pm)
	}
	return keys
}

// sendDA ends the primary device attributes.
func (em *emulator) sendDA() {
	buf := &bytes.Buffer{}
	_, _ = fmt.Fprintf(buf, "\x1b[?63")
	if em.be.Colors() > 0 {
		_, _ = fmt.Fprintf(buf, ";22")
	}
	if _, ok := em.be.(Clipboard); ok {
		_, _ = fmt.Fprintf(buf, ";52")
	}
	// 9 for NRC?
	// 15 for graphics?
	buf.WriteRune('c')
	em.SendRaw(buf.Bytes())
}

// setAnsiMode sets the ANSI mode.
func (em *emulator) setAnsiMode(mode AnsiMode, ms ModeStatus) {
	if !ms.Changeable() {
		return
	}
	em.modeLock.Lock()
	defer em.modeLock.Unlock()
	if old, ok := em.ansiModes[mode]; ok && old.Changeable() {
		em.ansiModes[mode] = ms
	}
}

func (em *emulator) getAnsiMode(mode AnsiMode) ModeStatus {
	em.modeLock.RLock()
	defer em.modeLock.RUnlock()
	return em.ansiModes[mode]
}

// getPrivateMode returns the value of a DEC private mode.
func (em *emulator) getPrivateMode(pm PrivateMode) ModeStatus {
	em.modeLock.RLock()
	if ms, ok := em.localModes[pm]; ok {
		em.modeLock.RUnlock()
		return ms
	}
	em.modeLock.RUnlock()
	return em.be.GetPrivateMode(pm)
}

func (em *emulator) updateMouseReporting() {
	mi, ok := em.be.(Mouser)
	if !ok {
		return
	}
	em.modeLock.RLock()
	report := em.mouseReportingLocked()
	em.modeLock.RUnlock()
	mi.SetMouse(report)
}

// setPrivateMode sets the DEC private mode.
func (em *emulator) setPrivateMode(pm PrivateMode, ms ModeStatus) {
	if !ms.Changeable() {
		return
	}
	em.modeLock.Lock()
	old, ok := em.localModes[pm]
	if ok && old.Changeable() {
		em.localModes[pm] = ms

		var (
			setMouse  bool
			report    MouseReporting
			setCursor bool
			cursor    CursorStyle
		)

		switch pm {
		case PmMouseButton, PmMouseDrag, PmMouseMotion, PmMouseSgr, PmMouseSgrPixel, PmMouseX10:
			report = em.mouseReportingLocked()
			setMouse = true
		case PmShowCursor:
			if ms == ModeOn {
				em.cursor = em.cursor.Show()
			} else {
				em.cursor = em.cursor.Hide()
			}
			cursor = em.cursor
			setCursor = true
		case PmBlinkCursor:
			if ms == ModeOn {
				em.cursor = em.cursor.Blink()
			} else {
				em.cursor = em.cursor.Steady()

			}
			cursor = em.cursor
			setCursor = true
		}
		em.modeLock.Unlock()

		if setMouse {
			if mi, ok := em.be.(Mouser); ok {
				mi.SetMouse(report)
			}
		}
		if setCursor {
			em.be.SetCursor(cursor)
		}
		return
	}
	em.modeLock.Unlock()

	if em.be.GetPrivateMode(pm).Changeable() {
		_ = em.be.SetPrivateMode(pm, ms)
	}
}

func (em *emulator) mouseReportingLocked() MouseReporting {
	switch {
	case em.localModes[PmMouseButton] == ModeOn:
		em.mouseReports = MouseButtons
		if em.localModes[PmMouseMotion] == ModeOn {
			em.mouseReports = MouseMotion
		} else if em.localModes[PmMouseDrag] == ModeOn {
			em.mouseReports = MouseDrag
		}
	case em.localModes[PmMouseX10] == ModeOn:
		em.mouseReports = MouseButtons
	default:
		em.mouseReports = MouseDisabled
	}
	return em.mouseReports
}

// SendRaw allows raw data to be sent to the application.
// This is done in a thread-safe way, so that content is not intermingled.
func (em *emulator) SendRaw(b []byte) {
	em.sendLock.Lock()
	defer em.sendLock.Unlock()

	// Do not attempt to send *anything* if we are stopped.
	select {
	case <-em.stopQ:
		return
	default:
	}

	// Try to write to the readQ, but if we cannot, then wait until
	// either we can, or the stopQ is fired.  This ensures that we avoid
	// breaking up content if at all possible.
	for _, ch := range b {
		select {
		case em.readQ <- ch:
		default:
			select {
			case em.readQ <- ch:
			case <-em.stopQ:
				return
			}
		}
	}
}

// KeyEvent injects a keyboard event into the emulator
func (em *emulator) KeyEvent(ev KeyEvent) {

	if em.getPrivateMode(PmWin32Input) == ModeOn {
		em.keyWin32IM(ev)
	} else {
		// eliminate "control" keys (which keyboard maps provide) from consideration.
		// (We handle control keys explicitly.)
		if ev.Utf != "" && ev.Utf[0] < ' ' {
			ev.Utf = ""
		}
		// TODO: more add support for kitty, and maybe modify other keys
		em.keyLegacy(ev)
	}
}

// ResizeEvent is called by the backend when a resize occurs.  A real backend with a child
// process (essentially a "real emulator") should probably also fire SIGWINCH if appropriate.
// That would be the job of something other than this code.
func (em *emulator) ResizeEvent(size Coord) {
	select {
	case em.writeQ <- size:
	case <-em.stopQ:
	}
}

func (em *emulator) applyResize(size Coord) {
	// resize clobbers our content, until it is redrawn
	em.size = size
	// resizing resets the margins
	em.topMargin = 0
	em.botMargin = em.size.Y - 1
	em.ltMargin = 0
	em.rtMargin = em.size.X - 1
	em.cells = make([]Cell, int(em.size.X)*int(em.size.Y))
	for i := range em.cells {
		em.cells[i].S = em.defaultStyle
	}

	em.pos = em.getPosition()
	if em.getPrivateMode(PmResizeReports) == ModeOn { // NB: we never support "ModeOnLocked"
		// NB: for now we do not support pixel sizes
		em.SendRaw(fmt.Appendf(nil, "\x1b[48;%d;%d;0;0t", em.size.Y, em.size.X))
	}

	// Send a SIGWINCH or similar.
	em.be.RaiseResize()
}

var legacyKeys = map[Key]struct {
	K  string // unmodified key
	A  string // unmodified in application cursor mode (smkx)
	S  string // with shift (if empty use regular modifier)
	C  string // with control (if empty use regular modifier)
	CS string // with ctrl-shift
}{
	KeyF1:        {K: "\x1bOP"}, // SS3 P
	KeyF2:        {K: "\x1bOQ"}, // SS3 Q
	KeyF3:        {K: "\x1bOR"}, // SS3 R
	KeyF4:        {K: "\x1bOS"}, // SS3 S
	KeyF5:        {K: "\x1b[15~"},
	KeyF6:        {K: "\x1b[17~"},
	KeyF7:        {K: "\x1b[18~"},
	KeyF8:        {K: "\x1b[19~"},
	KeyF9:        {K: "\x1b[20~"},
	KeyF10:       {K: "\x1b[21~"},
	KeyF11:       {K: "\x1b[23~"},
	KeyF12:       {K: "\x1b[24~"},
	KeyF13:       {K: "\x1b[25~"},
	KeyF14:       {K: "\x1b[26~"},
	KeyF15:       {K: "\x1b[28~"},
	KeyF16:       {K: "\x1b[29~"},
	KeyF17:       {K: "\x1b[31~"},
	KeyF18:       {K: "\x1b[32~"},
	KeyF19:       {K: "\x1b[33~"},
	KeyF20:       {K: "\x1b[34~"},
	KeyUp:        {K: "\x1b[A", A: "\x1bOA"},
	KeyDown:      {K: "\x1b[B", A: "\x1bOB"},
	KeyRight:     {K: "\x1b[C", A: "\x1bOC"},
	KeyLeft:      {K: "\x1b[D", A: "\x1bOD"},
	KeyHome:      {K: "\x1b[H", A: "\x1bOH"},
	KeyEnd:       {K: "\x1b[F", A: "\x1bOF"},
	KeyPgUp:      {K: "\x1b[5~"},
	KeyPgDn:      {K: "\x1b[6~"},
	KeyDelete:    {K: "\x1b[3~"},
	KeyInsert:    {K: "\x1b[2~"},
	KeyMenu:      {K: "\x1b[29~"}, // also F16
	KeyTab:       {K: "\t", S: "\x1b[Z", CS: "\x1b[Z"},
	KeyBackspace: {K: "\x7f", S: "\x7f", C: "\x08", CS: "\x08"},
	KeySpace:     {K: " ", S: " ", C: "\x00", CS: "\x00"},
	KeyEnter:     {K: "\r", S: "\r", CS: "\r"}, // NB: consider using kitty encoding here
	KeyPadEnter:  {K: "\r", S: "\r", CS: "\r"}, // NB: consider using kitty encoding here
	KeyEsc:       {K: "\x1b", S: "\x1b", C: "\x1b"},
}

var legacyControls = map[Key]string{
	// These ones are weird legacy control sequences that we mostly
	// do not care about.  We don't include shifted variants.
	Key2:      "\x00",
	Key3:      "\x1b",
	Key4:      "\x1c",
	Key5:      "\x1d",
	Key6:      "\x1e",
	Key7:      "\x1f",
	Key8:      "\x7f",
	KeyLBrace: "\x1b",
	KeySlash:  "\x1c",
	KeyRBrace: "\x1d",
}

// legacyPadKeys are keys that are on the keypad, when not in numeric keypad mode.
// Note that num lock overrides this.
var legacyPadKeys = map[Key]struct {
	app string
	num string
}{
	KeyPadEnter: {"\x1bOM", "\r"},
	KeyPadMul:   {"\x1bOj", "*"},
	KeyPadAdd:   {"\x1bOk", "+"},
	KeyPadSub:   {"\x1bOm", "-"},
	KeyPadDiv:   {"\x1bOo", "/"},
	KeyPadDec:   {"\x1b[3~", "."}, // Del
	KeyPad0:     {"\x1b[2~", "0"}, // Ins
	KeyPad1:     {"\x1bOF", "1"},  // End
	KeyPad2:     {"\x1b[B", "2"},  // Down
	KeyPad3:     {"\x1b[6~", "3"}, // PgDn
	KeyPad4:     {"\x1b[D", "4"},  // Left
	KeyPad5:     {"\x1b[E", "5"},  // Clear/Begin
	KeyPad6:     {"\x1b[C", "6"},  // Right
	KeyPad7:     {"\x1bOH", "7"},  // Home
	KeyPad8:     {"\x1b[A", "8"},  // Up
	KeyPad9:     {"\x1b[5~", "9"}, // PgUp
	KeyPadEqual: {"\x1bOX", "="},
}

// repeatRaw is called to provide key repeat.  We limit key repeating to just 40,
// and we ensure that at least one is included.  We only repeat if key repeat is enabled.
func (em *emulator) repeatRaw(ev KeyEvent, data []byte) {
	if pm := em.getPrivateMode(PmAutoRepeat); pm == ModeOn || pm == ModeOnLocked {
		for range min(max(1, ev.Repeat), 40) {
			em.SendRaw(data)
		}
	} else {
		if ev.Repeat == 0 {
			em.SendRaw(data)
		}
	}
}

// noRepeatRaw is used to send a key that should never repeat.
// It will only send if the repeat count is zero.
func (em *emulator) noRepeatRaw(ev KeyEvent, data []byte) {
	if ev.Repeat == 0 {
		em.SendRaw(data)
	}
}

// keyLegacy handles a keyboard event when in legacy vt220 style mode.
func (em *emulator) keyLegacy(ev KeyEvent) {
	if !ev.Down { // legacy protocol does not support key release
		return
	}
	if ev.Mod.IsMeta() || ev.Mod.IsHyper() { // legacy protocol does not support these
		return
	}

	// Shift-Ctrl keys are never sent in the legacy protocol.  We do have to ensure
	// that if we are sending other Utf (for example with AltGr), then we still might
	// send it, but this is only an issue for non-ASCII runes. Also, this filter only
	// applies for "regular" keys (i.e. not function keys, cursor keys, etc.)
	if ev.Mod.IsShift() && ev.Mod.IsCtrl() && (ev.Utf == "" || ev.Utf[0] < 0x80) {
		if base := ev.Key.KittyBase(); base >= ' ' && base < 0x80 {
			return
		}
	}

	// keypad sequences
	if v, ok := legacyPadKeys[ev.Key]; ok {
		if ev.Mod&ModNumLock == 0 {
			if em.appKeyPad {
				em.repeatRaw(ev, []byte(v.app))
			} else {
				em.repeatRaw(ev, []byte(v.num))
			}
			return
		} else {
			ev.Utf = v.num
		}
	}

	// For control keys (e.g. control-J) we never emit a rune directly -- but we might later
	// add after decoding the key accordingly.
	if ev.Utf != "" && (ev.Mod == ModLCtrl || ev.Mod == ModRCtrl || ev.Utf[0] < ' ') {
		ev.Utf = ""
	}

	if ev.Utf != "" {
		if ev.Utf[0] < 0x80 && ev.Mod.IsAlt() { // ASCII might get alt
			em.noRepeatRaw(ev, fmt.Appendf(nil, "\x1b%s", ev.Utf))
		} else { // otherwise send the UTF as-is
			em.repeatRaw(ev, []byte(ev.Utf))
		}
		return
	}

	// some weird number control sequences - legacy compatibility
	// We do not repeat these.
	if v, ok := legacyControls[ev.Key]; ok && (ev.Mod == ModLCtrl || ev.Mod == ModRCtrl) {
		em.noRepeatRaw(ev, []byte(v))
		return
	}

	if v, ok := legacyKeys[ev.Key]; ok {
		str := ""
		match := false
		if !ev.Mod.IsShift() && !ev.Mod.IsCtrl() {
			if em.getPrivateMode(PmAppCursor) == ModeOn && v.A != "" {
				str = v.A
			} else {
				str = v.K
			}
			// AnsiMode 20 sends newline, but only in legacy mode.
			if str == "\r" && em.getAnsiMode(AmNewLineMode) == ModeOn {
				str = "\r\n"
			}
			match = true
		} else if ev.Mod.IsShift() && !ev.Mod.IsCtrl() {
			if str = v.S; str != "" {
				match = true
			}
		} else if ev.Mod.IsCtrl() && !ev.Mod.IsShift() {
			if str = v.C; str != "" {
				match = true
			}
		} else { // IsCtrl & IsShift
			if str = v.CS; str != "" {
				match = true
			}
		}
		if !match {
			// No specific modifiers present, lets add them. There are two cases,
			// one for SS3 based keys and another for CSI based keys.  SS3 based
			// keys are converted to CSI - 1 ; mod ; final
			// Note: legacy encoding does not use modifiers for alt or super - alt will be
			// determined by sending an escape prefix.
			mod := 0
			if ev.Mod.IsShift() {
				mod |= 1
			}
			if ev.Mod.IsCtrl() {
				mod |= 4
			}
			if strings.HasPrefix(v.K, "\x1bO") {
				str = fmt.Sprintf("\x1b[1;%d%c", mod+1, v.K[len(v.K)-1])
			} else {
				str = fmt.Sprintf("%s;%d%c", v.K[:len(v.K)-1], mod+1, v.K[len(v.K)-1])
			}
		}
		if ev.Mod.IsAlt() {
			// no repeating ALT sequences
			em.noRepeatRaw(ev, append([]byte{'\x1b'}, []byte(str)...)) // alt sends leading escape
		} else if ev.Mod.IsCtrl() {
			// no repeating CTRL sequences
			em.noRepeatRaw(ev, []byte(str))
		} else {
			// but other sequences (should just be shifted or unmodified)
			// are fine.  (E.g. we want to allow repeats of cursor keys)
			em.repeatRaw(ev, []byte(str))
		}
		return
	}

	// fallback control key handling
	if ev.Key >= KeyA && ev.Key <= KeyZ && ev.Mod.IsCtrl() {
		b := byte(ev.Key-KeyA) + 1 /* ctrl-A */
		if ev.Mod.IsAlt() {
			em.noRepeatRaw(ev, []byte{'\x1b', b})
		} else {
			em.noRepeatRaw(ev, []byte{b})
		}
		return
	}
}

var win32NoRepeat = map[Key]bool{
	KeyLShift:   true,
	KeyRShift:   true,
	KeyLCtrl:    true,
	KeyRCtrl:    true,
	KeyLAlt:     true,
	KeyRAlt:     true,
	KeyLMeta:    true,
	KeyRMeta:    true,
	KeyCapsLock: true,
	KeyNumLock:  true,
	KeyEnter:    true,
	KeyScrLock:  true,
	KeyPause:    true,
	KeyPrtScr:   true,
}

// keyWin32IM generates the sequence for a key event when in Win32 input mode.
// Win32 input mode is ESC [ Vk ; Sc ; Uc ; Kd ; Cs ; Rc _
// Note that we specifically do NOT doubly encode non-keyboard events -- those
// are already unambiguously handled within the protocol.  (Windows Terminal behaves
// the same way, but most 3rd party terminals do doubly encode.)
func (em *emulator) keyWin32IM(ev KeyEvent) {
	// Some keys that never repeat
	if pm := em.getPrivateMode(PmAutoRepeat); pm == ModeOff || pm == ModeOffLocked {
		if ev.Repeat != 0 {
			return
		}
	}
	r := rune(0)
	if ev.Utf != "" {
		runes := []rune(ev.Utf)
		if len(runes) == 1 {
			r = runes[0]
		}
	}
	kd := 0
	if ev.Down {
		kd = 1
	}
	cs := 0
	// Modifiers
	if ev.Mod&ModRAlt != 0 {
		cs |= 0x01
	}
	if ev.Mod&ModLAlt != 0 {
		cs |= 0x02
	}
	if ev.Mod&ModRCtrl != 0 {
		cs |= 0x04
	}
	if ev.Mod&ModLCtrl != 0 {
		cs |= 0x08
	}
	if ev.Mod.IsShift() {
		cs |= 0x10
	}
	if ev.Mod.IsNumLock() {
		cs |= 0x20
	}
	// NB: 0x40 is for scroll lock, we don't support it for now
	if ev.Mod.IsCapsLock() {
		cs |= 0x80
	}
	switch ev.Key {
	case KeyPadEnter:
	case KeyPadDiv:
	case KeyInsert:
	case KeyDelete:
	case KeyHome:
	case KeyEnd:
	case KeyPgUp:
	case KeyPgDn:
		cs |= 0x100 // enhanced
	}
	if win32NoRepeat[ev.Key] {
		if ev.Repeat > 0 {
			return
		}
		ev.Repeat = 1
	}
	em.SendRaw(fmt.Appendf(nil, "\x1b[%d;%d;%d;%d;%d;%d_", ev.VK, ev.SC, r, kd, cs, max(1, ev.Repeat)))
}

func (em *emulator) MouseEvent(ev MouseEvent) {
	if pm := em.getPrivateMode(PmMouseButton); pm == ModeOn {

		if em.getPrivateMode(PmMouseDrag) != ModeOn && em.getPrivateMode(PmMouseMotion) != ModeOn {
			// suppress motion events if the user didn't request
			if ev.Button == NoButton {
				return // if entire event was just motion, bail
			}
			ev.Motion = false
		}

		if pm = em.getPrivateMode(PmMouseSgr); pm == ModeOn {
			btn := ev.encodeButton()
			if ev.Down {
				em.SendRaw(fmt.Appendf(nil, "\x1b[<%d;%d;%dM", btn, ev.Position.X+1, ev.Position.Y+1))
			} else {
				em.SendRaw(fmt.Appendf(nil, "\x1b[<%d;%d;%dm", btn, ev.Position.X+1, ev.Position.Y+1))
			}
		} else {
			// Old style reporting (via 1000h).
			// Limitations of legacy VT200 reporting are that the coordinates must be between
			// 1 and 223 inclusive, and that once any release occurs all buttons are assumed
			// to be released.  (Please use SGR mode if at all possible.)
			// Further, this mode is not CSI compliant as the encoded values that arrive ahead of
			// the final character may be within the range of technically legal CSI final bytes.
			if !ev.Down {
				ev.Button = NoButton
			}
			btn := ev.encodeButton()
			data := append([]byte{'\x1b', '[', 'M', btn + 32},
				byte(min(ev.Position.X+1, 223)+32),
				byte(min(ev.Position.Y+1, 223)+32))
			em.SendRaw(data)
		}

	} else if pm := em.getPrivateMode(PmMouseX10); pm == ModeOn && ev.Down {
		// legacy X10 reporting only
		x := byte(min(ev.Position.X+1, 223)) + 32
		y := byte(min(ev.Position.Y+1, 223)) + 32
		// NB: we intentionally reverse buttons 2 & 3 (for xterm compatibility)
		switch ev.Button {
		case Button1:
			em.SendRaw([]byte{'\x1b', '[', 'M', ' ', x, y})
		case Button2:
			em.SendRaw([]byte{'\x1b', '[', 'M', '"', x, y})
		case Button3:
			em.SendRaw([]byte{'\x1b', '[', 'M', '!', x, y})
		}
	}
}

func (em *emulator) FocusEvent(focused bool) {
	if pm := em.getPrivateMode(PmFocusReports); pm == ModeOn {
		if focused {
			em.SendRaw([]byte{'\x1b', '[', 'I'})
		} else {
			em.SendRaw([]byte{'\x1b', '[', 'O'})
		}
	}
}

// SetId sets the terminal name and version.
func (em *emulator) SetId(name string, version string) {
	em.name = name
	em.vers = version
}

// Start the terminal emulator.
func (em *emulator) Start() error {
	select {
	case <-em.stopQ:
	default:
		// already running
		return errors.New("terminal already started")
	}
	stopQ := make(chan bool)
	em.stopQ = stopQ
	go em.run(stopQ)
	return nil
}

// Stop the terminal emulator.  This also wakes any blocked
// Read or Write calls, which will return an error.
func (em *emulator) Stop() error {
	select {
	case <-em.stopQ:
	default:
		close(em.stopQ)
	}
	return nil
}

// Drain pending output to the terminal emulator.
func (em *emulator) Drain() error {
	q := make(chan bool)
	select {
	case em.writeQ <- q:
	case <-em.stopQ:
	}
	select {
	case <-q:
	case <-em.stopQ:
	}
	// make sure to wake the reader
	select {
	case em.readQ <- true:
	default:
	}
	return nil
}

// Write data to the emulator (commands).
func (em *emulator) Write(data []byte) (n int, err error) {
	stopQ := em.stopQ
	writeQ := em.writeQ
	drainQ := make(chan bool)
	select {
	case writeQ <- data:
		// we add the drainQ for synchronization, so that we only
		// return after the the emulator has processed this.
		select {
		case <-stopQ:
			return 0, errors.New("terminal emulator stopped")
		case writeQ <- drainQ:
		}
		select {
		case <-stopQ:
			return 0, errors.New("terminal emulator stopped")
		case <-drainQ:
			return len(data), nil
		}
	case <-stopQ:
		return 0, errors.New("terminal emulator stopped")
	}
}

// Read data (key events, etc.) from the emulator.
func (em *emulator) Read(data []byte) (n int, err error) {
	stopQ := em.stopQ
	readQ := em.readQ

	n = 0
	if len(data) < 1 {
		return 0, nil
	}
	select {
	case <-stopQ:
		return 0, errors.New("terminal emulator stopped")
	case v := <-readQ:
		// The data arriving in the channel may be a byte, or it might be a bool
		// trying to force a wake up.  Note that the bool may be intermingled with other
		// bytes, so we check it. Also data may have arrived since the bool was posted,
		// so make sure we don't terminate until we have collected all the relevant data
		// that we can (up to the limit of what was requested.)
		if ch, ok := v.(byte); ok {
			data[n] = ch
			n++
		}
		for n < len(data) {
			select {
			case v = <-readQ:
				if ch, ok := v.(byte); ok {
					data[n] = ch
					n++
				}
			default:
				return n, nil
			}
		}
		return n, nil
	}
}

func (em *emulator) run(stopQ <-chan bool) {
	for {
		select {
		case item := <-em.writeQ:
			switch d := item.(type) {
			case byte:
				em.inb(d)
			case []byte:
				for _, ch := range d {
					em.inb(ch)
				}
			case chan bool:
				close(d)

			case Coord: // resize notification
				em.applyResize(d)
			}
		case <-stopQ:
			return
		}
	}
}
