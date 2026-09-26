package gocui

import (
	"strings"
	"testing"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/stretchr/testify/assert"
)

func describeReplies(replies []terminalReply) []string {
	result := []string{}
	for _, reply := range replies {
		if reply.background != nil {
			result = append(result, reply.background.Hex())
		} else {
			result = append(result, reply.colorScheme.String())
		}
	}
	return result
}

func TestTerminalReplyScanner(t *testing.T) {
	scenarios := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "dark report",
			input:    "\x1b[?997;1n",
			expected: []string{"dark"},
		},
		{
			name:     "light report",
			input:    "\x1b[?997;2n",
			expected: []string{"light"},
		},
		{
			name:     "background ended by ST",
			input:    "\x1b]11;rgb:1e1e/1e1e/1e1e\x1b\\",
			expected: []string{"#1e1e1e"},
		},
		{
			name:     "background ended by BEL",
			input:    "\x1b]11;rgb:ffff/ffff/ffff\a",
			expected: []string{"#ffffff"},
		},
		{
			name:     "background with two hex digits per component",
			input:    "\x1b]11;rgb:fd/f6/e3\a",
			expected: []string{"#fdf6e3"},
		},
		{
			name:     "background with one hex digit per component",
			input:    "\x1b]11;rgb:f/0/f\a",
			expected: []string{"#ff00ff"},
		},
		{
			name:     "background with alpha",
			input:    "\x1b]11;rgba:0000/2b2b/3636/ffff\a",
			expected: []string{"#002b36"},
		},
		{
			name:     "replies among other input",
			input:    "j\x1b[?997;1n\x1b[A\x1b[<0;10;5M\x1b]11;rgb:0000/0000/0000\x1b\\\x1b[?62;22c",
			expected: []string{"dark", "#000000"},
		},
		{
			name:     "other reports",
			input:    "\x1b[?997;3n\x1b[0n\x1b[?996n\x1b]10;rgb:ffff/ffff/ffff\a\x1b]4;1;rgb:ffff/0000/0000\a",
			expected: []string{},
		},
		{
			name:     "malformed backgrounds",
			input:    "\x1b]11;rgb:zz/00/00\a\x1b]11;rgb:10/20\a\x1b]11;rgb:12345/0/0\a\x1b]11;rgb://\a\x1b]11;#ffffff\a",
			expected: []string{},
		},
		{
			name:     "reply after a sequence that was cut short",
			input:    "\x1b]11;rgb:\x1b[?997;2n\x1b[?99\x1b]11;rgb:00/00/00\a",
			expected: []string{"light", "#000000"},
		},
		{
			name:     "reply after an oversized sequence",
			input:    "\x1b]11;" + strings.Repeat("x", 100) + "\a\x1b[?" + strings.Repeat("9", 100) + "n\x1b[?997;1n",
			expected: []string{"dark"},
		},
		{
			name:     "oversized sequence that starts like a reply",
			input:    "\x1b]11;rgb:00/00/00" + strings.Repeat("x", 100) + "\a",
			expected: []string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			scanner := terminalReplyScanner{}
			assert.Equal(t, s.expected, describeReplies(scanner.scan([]byte(s.input))))
		})

		t.Run(s.name+", split in two", func(t *testing.T) {
			for i := range len(s.input) {
				scanner := terminalReplyScanner{}
				replies := scanner.scan([]byte(s.input[:i]))
				replies = append(replies, scanner.scan([]byte(s.input[i:]))...)
				assert.Equal(t, s.expected, describeReplies(replies), "split at %d", i)
			}
		})

		t.Run(s.name+", byte by byte", func(t *testing.T) {
			scanner := terminalReplyScanner{}
			replies := []terminalReply{}
			for i := range len(s.input) {
				replies = append(replies, scanner.scan([]byte{s.input[i]})...)
			}
			assert.Equal(t, s.expected, describeReplies(replies))
		})
	}
}

func TestColorSchemeQueries(t *testing.T) {
	allQueries := enableColorSchemeReports + requestColorScheme + requestBackgroundColor

	scenarios := []struct {
		name           string
		term           string
		termProgram    string
		tcellNegotiate string
		expected       string
	}{
		{name: "xterm", term: "xterm-256color", expected: allQueries},
		{name: "tmux", term: "tmux-256color", termProgram: "tmux", expected: allQueries},
		{name: "Terminal.app", term: "xterm-256color", termProgram: "Apple_Terminal", expected: requestBackgroundColor},
		{name: "WezTerm", term: "xterm-256color", termProgram: "WezTerm", expected: requestBackgroundColor},
		{name: "st", term: "st-256color", expected: ""},
		{name: "Linux console", term: "linux", expected: ""},
		{name: "VT100", term: "vt100", expected: ""},
		{name: "negotiation disabled", term: "xterm-256color", tcellNegotiate: "disable", expected: ""},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, colorSchemeQueries(s.term, s.termProgram, s.tcellNegotiate))
		})
	}
}

type fakeTty struct {
	input   string
	written strings.Builder
}

var _ tcell.Tty = &fakeTty{}

func (self *fakeTty) Start() error                          { return nil }
func (self *fakeTty) Stop() error                           { return nil }
func (self *fakeTty) Drain() error                          { return nil }
func (self *fakeTty) NotifyResize(chan<- bool)              {}
func (self *fakeTty) WindowSize() (tcell.WindowSize, error) { return tcell.WindowSize{}, nil }
func (self *fakeTty) Close() error                          { return nil }
func (self *fakeTty) Write(p []byte) (int, error)           { return self.written.Write(p) }
func (self *fakeTty) Read(p []byte) (int, error) {
	n := copy(p, self.input)
	self.input = self.input[n:]
	return n, nil
}

func (self *fakeTty) takeWritten() string {
	written := self.written.String()
	self.written.Reset()
	return written
}

type colorSchemeTtyTest struct {
	fake          *fakeTty
	tty           *colorSchemeTty
	notifications []string
}

func newColorSchemeTtyTest(termProgram string) *colorSchemeTtyTest {
	test := &colorSchemeTtyTest{fake: &fakeTty{}}
	test.tty = &colorSchemeTty{
		Tty:     test.fake,
		queries: colorSchemeQueries("xterm-256color", termProgram, ""),
	}
	test.tty.subscribe(func(colorScheme DetectedColorScheme) {
		test.notifications = append(test.notifications, colorScheme.String())
	})
	return test
}

// feed has the terminal send input, and has tcell read it
func (self *colorSchemeTtyTest) feed(input string) {
	self.fake.input = input
	_, _ = self.tty.Read(make([]byte, 128))
}

func (self *colorSchemeTtyTest) takeNotifications() []string {
	notifications := self.notifications
	self.notifications = nil
	return notifications
}

func TestColorSchemeTtyPrefersTheBackground(t *testing.T) {
	test := newColorSchemeTtyTest("")

	assert.NoError(t, test.tty.Start())
	assert.Equal(t, enableColorSchemeReports+requestColorScheme+requestBackgroundColor, test.fake.takeWritten())

	// The terminal says it's light, but its background is dark
	test.feed("\x1b[?997;2n")
	assert.Equal(t, []string{"light (reported by the terminal)"}, test.takeNotifications())
	assert.Equal(t, "", test.fake.takeWritten(), "the background was asked for already")

	test.feed("\x1b]11;rgb:1e1e/1e1e/1e1e\x1b\\\x1b[?62;22c")
	assert.Equal(t, []string{"dark (background #1e1e1e)"}, test.takeNotifications())
	assert.Equal(t, "", test.fake.takeWritten())

	// A later report only makes us ask for the background again
	test.feed("\x1b[?997;2n")
	assert.Nil(t, test.takeNotifications())
	assert.Equal(t, requestBackgroundColor, test.fake.takeWritten())

	test.feed("\x1b[?997;2n")
	assert.Equal(t, "", test.fake.takeWritten(), "the background was asked for already")

	test.feed("\x1b]11;rgb:ffff/ffff/ffff\x1b\\")
	assert.Equal(t, []string{"light (background #ffffff)"}, test.takeNotifications())

	// The same background again is no change
	test.feed("\x1b[?997;2n")
	assert.Equal(t, requestBackgroundColor, test.fake.takeWritten())
	test.feed("\x1b]11;rgb:ffff/ffff/ffff\x1b\\")
	assert.Nil(t, test.takeNotifications())
}

func TestColorSchemeTtyWithoutBackground(t *testing.T) {
	test := newColorSchemeTtyTest("")

	assert.NoError(t, test.tty.Start())
	test.fake.takeWritten()

	test.feed("\x1b[?997;1n\x1b[?62;22c")
	assert.Equal(t, []string{"dark (reported by the terminal)"}, test.takeNotifications())

	test.feed("\x1b[?997;2n")
	assert.Equal(t, []string{"light (reported by the terminal)"}, test.takeNotifications())
	assert.Equal(t, "", test.fake.takeWritten(), "a terminal that didn't answer before isn't asked again")

	test.tty.onFocusGained()
	assert.Equal(t, "", test.fake.takeWritten())
}

func TestColorSchemeTtyAsksAgainOnFocus(t *testing.T) {
	test := newColorSchemeTtyTest("Apple_Terminal")

	assert.NoError(t, test.tty.Start())
	assert.Equal(t, requestBackgroundColor, test.fake.takeWritten())

	test.tty.onFocusGained()
	assert.Equal(t, "", test.fake.takeWritten(), "the background was asked for already")

	test.feed("\x1b]11;rgb:0000/0000/0000\a")
	assert.Equal(t, []string{"dark (background #000000)"}, test.takeNotifications())

	test.tty.onFocusGained()
	assert.Equal(t, requestBackgroundColor, test.fake.takeWritten())
	test.tty.onFocusGained()
	assert.Equal(t, "", test.fake.takeWritten(), "the background was asked for already")

	test.feed("\x1b]11;rgb:ffff/ffff/ffff\a")
	assert.Equal(t, []string{"light (background #ffffff)"}, test.takeNotifications())
}

func TestColorSchemeTtyStopAndStart(t *testing.T) {
	test := newColorSchemeTtyTest("")

	assert.NoError(t, test.tty.Start())
	test.fake.takeWritten()
	test.feed("\x1b[?997;1n\x1b]11;rgb:0000/0000/0000\a")

	assert.NoError(t, test.tty.Stop())
	assert.Equal(t, disableColorSchemeReports, test.fake.takeWritten())

	test.tty.onFocusGained()
	assert.Equal(t, "", test.fake.takeWritten(), "nothing is written to a stopped tty")

	// Anything may have changed while we were stopped, so ask again
	assert.NoError(t, test.tty.Start())
	assert.Equal(t, enableColorSchemeReports+requestColorScheme+requestBackgroundColor, test.fake.takeWritten())
}

func TestColorSchemeTtyStopWithoutReports(t *testing.T) {
	test := newColorSchemeTtyTest("WezTerm")

	assert.NoError(t, test.tty.Start())
	test.fake.takeWritten()

	assert.NoError(t, test.tty.Stop())
	assert.Equal(t, "", test.fake.takeWritten(), "reports were never turned on")
}

func TestColorSchemeOfBackground(t *testing.T) {
	for background, expected := range map[string]ColorScheme{
		"rgb:0000/0000/0000": ColorSchemeDark,
		"rgb:1e/1e/1e":       ColorSchemeDark,
		"rgb:00/2b/36":       ColorSchemeDark,
		"rgb:70/70/70":       ColorSchemeDark,
		"rgb:80/80/80":       ColorSchemeLight,
		"rgb:fd/f6/e3":       ColorSchemeLight,
		"rgb:ff/ff/ff":       ColorSchemeLight,
	} {
		color, ok := parseBackgroundColorReport("11;" + background)
		assert.True(t, ok)
		assert.Equal(t, expected, colorSchemeOfBackground(color), background)
	}
}

// startWaiting starts waiting for replies, with a timeout too long to matter
func (self *colorSchemeTtyTest) startWaiting() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		self.tty.waitForReplies(time.Minute)
		close(done)
	}()
	return done
}

func isClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

func assertDoneWaiting(t *testing.T, done <-chan struct{}) {
	t.Helper()
	assert.Eventually(t, func() bool { return isClosed(done) }, time.Second, time.Millisecond)
}

func assertStillWaiting(t *testing.T, done <-chan struct{}) {
	t.Helper()
	assert.Never(t, func() bool { return isClosed(done) }, 50*time.Millisecond, time.Millisecond)
}

func TestColorSchemeTtyWaitsForReplies(t *testing.T) {
	test := newColorSchemeTtyTest("")

	assert.NoError(t, test.tty.Start())
	test.feed("\x1b[?997;1n\x1b]11;rgb:0000/0000/0000\a")
	assertDoneWaiting(t, test.startWaiting())

	assert.NoError(t, test.tty.Stop())
	assert.NoError(t, test.tty.Start())

	done := test.startWaiting()
	assertStillWaiting(t, done)
	test.feed("\x1b[?997;1n")
	assertStillWaiting(t, done)
	test.feed("\x1b]11;rgb:0000/0000/0000\a")
	assertDoneWaiting(t, done)

	test.tty.onFocusGained()
	done = test.startWaiting()
	assertStillWaiting(t, done)
	test.feed("\x1b]11;rgb:0000/0000/0000\a")
	assertDoneWaiting(t, done)
}

func TestColorSchemeTtyGivesUpWaiting(t *testing.T) {
	test := newColorSchemeTtyTest("")

	assert.NoError(t, test.tty.Start())
	test.feed("\x1b[?997;1n\x1b]11;rgb:0000/0000/0000\a")
	test.tty.onFocusGained()

	start := time.Now()
	test.tty.waitForReplies(20 * time.Millisecond)
	assert.GreaterOrEqual(t, time.Since(start), 20*time.Millisecond)
}

func TestColorSchemeTtyDoesntWaitForRepliesThatNeverCame(t *testing.T) {
	test := newColorSchemeTtyTest("")

	// The terminal answers nothing
	assert.NoError(t, test.tty.Start())
	assertDoneWaiting(t, test.startWaiting())

	// The terminal reports its color scheme, but not its background
	test.feed("\x1b[?997;1n\x1b[?62;22c")
	assert.NoError(t, test.tty.Stop())
	assert.NoError(t, test.tty.Start())

	done := test.startWaiting()
	assertStillWaiting(t, done)
	test.feed("\x1b[?997;1n")
	assertDoneWaiting(t, done)
}
