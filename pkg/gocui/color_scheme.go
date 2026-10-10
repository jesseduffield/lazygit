package gocui

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/lucasb-eyer/go-colorful"
)

// ColorScheme says whether the terminal shows light text on a dark background,
// or dark text on a light one.
type ColorScheme int

const (
	ColorSchemeUnknown ColorScheme = iota
	ColorSchemeDark
	ColorSchemeLight
)

func (self ColorScheme) String() string {
	switch self {
	case ColorSchemeDark:
		return "dark"
	case ColorSchemeLight:
		return "light"
	default:
		return "unknown"
	}
}

// DetectedColorScheme is what the terminal has told us about its colors.
type DetectedColorScheme struct {
	ColorScheme ColorScheme
	// The background color that ColorScheme was derived from, as #rrggbb. It is
	// empty if the terminal didn't tell us its background color; ColorScheme is
	// then whatever the terminal said about itself, if anything.
	Background string
}

func (self DetectedColorScheme) String() string {
	if self.Background != "" {
		return fmt.Sprintf("%s (background %s)", self.ColorScheme, self.Background)
	}
	if self.ColorScheme != ColorSchemeUnknown {
		return fmt.Sprintf("%s (reported by the terminal)", self.ColorScheme)
	}
	return self.ColorScheme.String()
}

const (
	// Makes the terminal send a color scheme report whenever its colors change
	enableColorSchemeReports  = "\x1b[?2031h"
	disableColorSchemeReports = "\x1b[?2031l"
	// Asks for a color scheme report: CSI ? 997 ; 1 n for dark, 2 for light
	requestColorScheme = "\x1b[?996n"
	// Asks for the background color: OSC 11 ; rgb:RRRR/GGGG/BBBB, ended by BEL
	// or ST. Ending the request with BEL makes more terminals reply.
	requestBackgroundColor = "\x1b]11;?\a"
)

// colorSchemeTty is a tcell.Tty that finds out the terminal's color scheme.
//
// Some terminals say whether they are dark or light, and can report it again
// whenever that changes. What they base this on varies, though: some go by the
// terminal's own colors, others by the dark or light mode of the operating
// system, whether or not the terminal follows it. So we also ask for the
// background color, and when the terminal tells us that, it decides; a color
// scheme report is then only a sign that the background may have changed.
//
// The queries go out when tcell starts the tty, which is before tcell sends its
// own queries during Screen.Init and waits for their answer. Terminals answer in
// order, so the answers to ours have arrived by the time Init returns, and we
// know the color scheme before drawing anything without waiting for it.
//
// tcell doesn't understand the answers and drops them, so all we need to do is
// watch for them in the input as it goes by.
type colorSchemeTty struct {
	tcell.Tty

	// What we send when the tty starts
	queries string

	// Guards all of the fields below, and serializes our writes to the terminal
	// with tcell's
	mutex sync.Mutex

	started bool

	background colorful.Color
	// Whether the terminal has told us its background color
	haveBackground bool
	// Whether we have asked for the background color and are waiting for the
	// answer
	backgroundRequested bool
	// The color scheme that the terminal last reported for itself
	reported ColorScheme
	// Whether we have asked for a color scheme report and are waiting for it
	colorSchemeRequested bool
	// Closed once no answer that waitForReplies waits for is outstanding
	repliesArrived chan struct{}

	notified DetectedColorScheme
	onChange func(DetectedColorScheme)

	// Only used by Read, which tcell never calls concurrently with itself
	scanner terminalReplyScanner
}

var _ tcell.Tty = &colorSchemeTty{}

func newColorSchemeTty(tty tcell.Tty) *colorSchemeTty {
	return &colorSchemeTty{
		Tty: tty,
		queries: colorSchemeQueries(
			os.Getenv("TERM"),
			os.Getenv("TERM_PROGRAM"),
			os.Getenv("TCELL_NEGOTIATE"),
		),
	}
}

// colorSchemeQueries returns what to ask the terminal. It leaves out the
// terminals that tcell doesn't send its own queries to (see
// applyKnownTerminalProfile and the legacy terminals in tScreen.Init), except
// for those that are known to answer a request for the background color.
func colorSchemeQueries(term string, termProgram string, tcellNegotiate string) string {
	if tcellNegotiate == "disable" {
		return ""
	}

	if term == "st" || strings.HasPrefix(term, "st-") ||
		strings.HasPrefix(term, "vt") || strings.Contains(term, "ansi") ||
		term == "linux" || term == "sun" || term == "sun-color" {
		return ""
	}

	if termProgram == "Apple_Terminal" || termProgram == "WezTerm" {
		return requestBackgroundColor
	}

	return enableColorSchemeReports + requestColorScheme + requestBackgroundColor
}

// subscribe sets a function to call whenever the detected color scheme changes,
// and returns the one detected so far. The function is called on the goroutine
// that reads from the tty, so it must not block.
func (self *colorSchemeTty) subscribe(onChange func(DetectedColorScheme)) DetectedColorScheme {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.onChange = onChange
	self.notified = self.detected()
	return self.notified
}

func (self *colorSchemeTty) Start() error {
	if err := self.Tty.Start(); err != nil {
		return err
	}

	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.started = true
	self.writeLocked(self.queries)
	if strings.Contains(self.queries, requestBackgroundColor) {
		self.backgroundRequested = true
	}
	if strings.Contains(self.queries, requestColorScheme) {
		self.colorSchemeRequested = true
	}

	return nil
}

// waitForReplies waits until the terminal has answered what we asked it, but
// no longer than the timeout. Call it before handing the terminal to another
// program, or the answers would reach that program as if they were typed. It
// only waits for the answers that the terminal has given before, so a terminal
// that doesn't answer at all costs no time.
func (self *colorSchemeTty) waitForReplies(timeout time.Duration) {
	self.mutex.Lock()
	if !self.awaitingRepliesLocked() {
		self.mutex.Unlock()
		return
	}
	if self.repliesArrived == nil {
		self.repliesArrived = make(chan struct{})
	}
	repliesArrived := self.repliesArrived
	self.mutex.Unlock()

	select {
	case <-repliesArrived:
	case <-time.After(timeout):
	}
}

func (self *colorSchemeTty) awaitingRepliesLocked() bool {
	return (self.backgroundRequested && self.haveBackground) ||
		(self.colorSchemeRequested && self.reported != ColorSchemeUnknown)
}

func (self *colorSchemeTty) Stop() error {
	self.mutex.Lock()
	// Otherwise, the program we hand the terminal to would receive the reports
	// as if they were typed
	if strings.Contains(self.queries, enableColorSchemeReports) {
		self.writeLocked(disableColorSchemeReports)
	}
	self.started = false
	self.mutex.Unlock()

	return self.Tty.Stop()
}

func (self *colorSchemeTty) Write(p []byte) (int, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	return self.Tty.Write(p)
}

func (self *colorSchemeTty) Read(p []byte) (int, error) {
	n, err := self.Tty.Read(p)

	for _, reply := range self.scanner.scan(p[:n]) {
		self.handleReply(reply)
	}

	return n, err
}

// onFocusGained asks for the background color again, for the terminals that
// tell us their background color but not when it changes. When the operating
// system switches between dark and light mode, the user is usually busy
// elsewhere, so coming back to the terminal is a good time to check.
func (self *colorSchemeTty) onFocusGained() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.requestBackgroundColorLocked()
}

func (self *colorSchemeTty) handleReply(reply terminalReply) {
	self.mutex.Lock()

	if reply.background != nil {
		self.background = *reply.background
		self.haveBackground = true
		self.backgroundRequested = false
	} else {
		self.reported = reply.colorScheme
		self.colorSchemeRequested = false
		self.requestBackgroundColorLocked()
	}

	if self.repliesArrived != nil && !self.awaitingRepliesLocked() {
		close(self.repliesArrived)
		self.repliesArrived = nil
	}

	detected := self.detected()
	onChange := self.onChange
	changed := detected != self.notified
	if changed {
		self.notified = detected
	}

	self.mutex.Unlock()

	if changed && onChange != nil {
		onChange(detected)
	}
}

// requestBackgroundColorLocked asks for the background color, if the terminal
// has told us its background color before and isn't about to do it anyway.
func (self *colorSchemeTty) requestBackgroundColorLocked() {
	if !self.started || !self.haveBackground || self.backgroundRequested {
		return
	}

	self.writeLocked(requestBackgroundColor)
	self.backgroundRequested = true
}

func (self *colorSchemeTty) detected() DetectedColorScheme {
	if !self.haveBackground {
		return DetectedColorScheme{ColorScheme: self.reported}
	}

	return DetectedColorScheme{
		ColorScheme: colorSchemeOfBackground(self.background),
		Background:  self.background.Hex(),
	}
}

func (self *colorSchemeTty) writeLocked(s string) {
	if s == "" {
		return
	}

	// If the terminal is gone, tcell finds out when it next writes or reads
	_, _ = self.Tty.Write([]byte(s))
}

func colorSchemeOfBackground(background colorful.Color) ColorScheme {
	lightness, _, _ := background.Lab()
	if lightness < 0.5 {
		return ColorSchemeDark
	}
	return ColorSchemeLight
}

// terminalReply is one of the answers that colorSchemeTty asks the terminal
// for. Exactly one of its fields is set.
type terminalReply struct {
	colorScheme ColorScheme
	background  *colorful.Color
}

type terminalReplyScannerState int

const (
	scanningText terminalReplyScannerState = iota
	scanningEscape
	scanningCsi
	scanningOsc
	scanningOscEscape
)

// No reply we look for is longer than this; we don't look at longer sequences
const maxTerminalReplyLength = 64

// terminalReplyScanner finds the terminal's answers to our color scheme queries
// in its input. The input arrives in chunks, which may split an answer anywhere,
// so the scanner keeps its state from one chunk to the next.
type terminalReplyScanner struct {
	state     terminalReplyScannerState
	sequence  []byte
	oversized bool
}

func (self *terminalReplyScanner) scan(input []byte) []terminalReply {
	var replies []terminalReply

	for _, b := range input {
		switch self.state {
		case scanningText:
			if b == '\x1b' {
				self.state = scanningEscape
			}

		case scanningEscape:
			self.scanEscape(b)

		case scanningCsi:
			switch {
			case b == '\x1b':
				self.state = scanningEscape
			case b >= 0x20 && b <= 0x3f: // parameter and intermediate bytes
				self.appendToSequence(b)
			case b >= 0x40 && b <= 0x7e: // final byte
				self.state = scanningText
				if b == 'n' && !self.oversized {
					if colorScheme, ok := parseColorSchemeReport(string(self.sequence)); ok {
						replies = append(replies, terminalReply{colorScheme: colorScheme})
					}
				}
			default:
				self.state = scanningText
			}

		case scanningOsc:
			switch b {
			case '\a':
				self.state = scanningText
				replies = self.appendOscReply(replies)
			case '\x1b':
				self.state = scanningOscEscape
			default:
				self.appendToSequence(b)
			}

		case scanningOscEscape:
			if b == '\\' {
				self.state = scanningText
				replies = self.appendOscReply(replies)
			} else {
				// Not a string terminator, so a new escape sequence has begun and
				// the one before it was cut short
				self.scanEscape(b)
			}
		}
	}

	return replies
}

// scanEscape handles the byte after an ESC.
func (self *terminalReplyScanner) scanEscape(b byte) {
	switch b {
	case '[':
		self.startSequence(scanningCsi)
	case ']':
		self.startSequence(scanningOsc)
	case '\x1b':
		self.state = scanningEscape
	default:
		self.state = scanningText
	}
}

func (self *terminalReplyScanner) startSequence(state terminalReplyScannerState) {
	self.state = state
	self.sequence = self.sequence[:0]
	self.oversized = false
}

func (self *terminalReplyScanner) appendToSequence(b byte) {
	if len(self.sequence) >= maxTerminalReplyLength {
		self.oversized = true
		return
	}
	self.sequence = append(self.sequence, b)
}

func (self *terminalReplyScanner) appendOscReply(replies []terminalReply) []terminalReply {
	if self.oversized {
		return replies
	}

	if background, ok := parseBackgroundColorReport(string(self.sequence)); ok {
		replies = append(replies, terminalReply{background: &background})
	}

	return replies
}

// parseColorSchemeReport parses the parameters of CSI ? 997 ; <n> n
func parseColorSchemeReport(params string) (ColorScheme, bool) {
	switch params {
	case "?997;1":
		return ColorSchemeDark, true
	case "?997;2":
		return ColorSchemeLight, true
	default:
		return ColorSchemeUnknown, false
	}
}

// parseBackgroundColorReport parses the content of OSC 11 ; rgb:R/G/B, where
// each of R, G and B has from one to four hex digits. Some terminals send
// rgba:R/G/B/A instead.
func parseBackgroundColorReport(content string) (colorful.Color, bool) {
	spec, ok := strings.CutPrefix(content, "11;")
	if !ok {
		return colorful.Color{}, false
	}

	var components []string
	if rgb, ok := strings.CutPrefix(spec, "rgb:"); ok {
		components = strings.Split(rgb, "/")
		if len(components) != 3 {
			return colorful.Color{}, false
		}
	} else if rgba, ok := strings.CutPrefix(spec, "rgba:"); ok {
		components = strings.Split(rgba, "/")
		if len(components) != 4 {
			return colorful.Color{}, false
		}
	} else {
		return colorful.Color{}, false
	}

	values := [3]float64{}
	for i := range values {
		component := components[i]
		if len(component) < 1 || len(component) > 4 {
			return colorful.Color{}, false
		}
		value, err := strconv.ParseUint(component, 16, 16)
		if err != nil {
			return colorful.Color{}, false
		}
		maxValue := uint64(1)<<(4*len(component)) - 1
		values[i] = float64(value) / float64(maxValue)
	}

	return colorful.Color{R: values[0], G: values[1], B: values[2]}, true
}
