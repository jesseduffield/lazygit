// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"strconv"
	"strings"

	"github.com/go-errors/errors"
)

type escapeInterpreter struct {
	state                  escapeState
	curch                  string
	csiParam               []string
	curFgColor, curBgColor Attribute
	mode                   OutputMode
	instruction            instruction
	hyperlink              strings.Builder

	// the OSC number being accumulated while we don't yet know which OSC this is
	oscNumber strings.Builder
	// the payload of an OSC 1717 per-line diff-metadata sequence (see
	// diff-line-metadata-notes.md), accumulated like hyperlink
	metadata strings.Builder

	// ConPTY emits cursor-positioning escapes (CUP) to skip over blank
	// rows rather than emitting LFs for them. To convert those into row
	// advances the view can act on, we track where in the pseudo-terminal
	// screen the cursor currently is. 1-based to match the escape
	// sequences.
	//
	// We also have to track the column, but only well enough to count
	// soft-wraps when written content runs past the right edge: ConPTY's
	// CUPs are addressed against its post-wrap screen, so a logical line
	// long enough to wrap in ConPTY's screen counts for two rows from the
	// next CUP's perspective. Column accuracy past wrap-counting isn't
	// modelled — we don't track the col argument of CUPs, and most
	// pager-style emitters use col 1 anyway.
	screenRow, screenCol int

	// The screen width that soft-wraps are counted against (see
	// notifyCellsWritten). It's a snapshot of the view's InnerWidth taken on
	// the UI thread (in NewView, and refreshed per render via
	// View.SetContentWidth), rather than read live from the view's dimensions:
	// a view's output is written from a task goroutine, and reading the live
	// dimensions there would race the UI thread updating them during layout.
	screenColMax int
}

type (
	escapeState int
	fontEffect  int
)

type instruction interface{ isInstruction() }

type eraseInLineFromCursor struct{}

func (self eraseInLineFromCursor) isInstruction() {}

// cursorDown asks the view to advance N rows. Emitted when CUP / CUD /
// CNL / VPA targets a row past the current one; backward moves are
// ignored because the view's buffer is line-based and can't undo.
type cursorDown struct{ n int }

func (self cursorDown) isInstruction() {}

// cursorForward asks the view to materialize N space cells. Emitted
// when CUF advances the cursor right — ConPTY uses CUF (often paired
// with ECH) to encode runs of default-colored spaces compactly, so we
// have to render the gap, not just bump a counter.
type cursorForward struct{ n int }

func (self cursorForward) isInstruction() {}

type noInstruction struct{}

func (self noInstruction) isInstruction() {}

const (
	stateNone escapeState = iota
	stateEscape
	stateCharacterSetDesignation
	stateCSI
	stateParams
	stateCSIDiscard
	stateOSC
	stateOSCParams
	stateOSCHyperlink
	stateOSCMetadata
	stateOSCEndEscape
	stateOSCSkipUnknown

	bold      fontEffect = 1
	faint     fontEffect = 2
	italic    fontEffect = 3
	underline fontEffect = 4
	blink     fontEffect = 5
	reverse   fontEffect = 7
	strike    fontEffect = 9

	setForegroundColor     int = 38
	defaultForegroundColor int = 39
	setBackgroundColor     int = 48
	defaultBackgroundColor int = 49
)

var (
	errNotCSI        = errors.New("Not a CSI escape sequence")
	errCSIParseError = errors.New("CSI escape sequence parsing error")
)

// characters in case of error will output the non-parsed characters as a string.
func (ei *escapeInterpreter) characters() []string {
	switch ei.state {
	case stateNone:
		return []string{"\x1b"}
	case stateEscape:
		return []string{"\x1b", ei.curch}
	case stateCSI:
		return []string{"\x1b", "[", ei.curch}
	case stateParams:
		ret := []string{"\x1b", "["}
		for _, s := range ei.csiParam {
			ret = append(ret, s)
			ret = append(ret, ";")
		}
		return append(ret, ei.curch)
	default:
	}
	return nil
}

// newEscapeInterpreter returns an escapeInterpreter that will be able to parse
// terminal escape sequences.
func newEscapeInterpreter(mode OutputMode) *escapeInterpreter {
	ei := &escapeInterpreter{
		state:       stateNone,
		curFgColor:  ColorDefault,
		curBgColor:  ColorDefault,
		mode:        mode,
		instruction: noInstruction{},
		screenRow:   1,
		screenCol:   1,
	}
	return ei
}

// reset sets the escapeInterpreter in initial state. Note: this only resets
// escape-parsing state. Screen cursor state survives so that mid-stream
// malformed escapes don't desync the row tracking from the view.
func (ei *escapeInterpreter) reset() {
	ei.state = stateNone
	ei.curFgColor = ColorDefault
	ei.curBgColor = ColorDefault
	ei.csiParam = nil
}

// resetScreenCursor returns the screen-cursor tracking to the top of the
// pseudo-terminal screen. Called when the view is rewound before a fresh pty
// render, and on cursor-home (which ConPTY emits at the start of each screen)
// for views that aren't rewound in lockstep — see the CUP handling in parseOne.
func (ei *escapeInterpreter) resetScreenCursor() {
	ei.screenRow = 1
	ei.screenCol = 1
}

// notifyRowAdvance must be called by the view whenever it advances to the
// next row in response to an LF / CRLF outside of an escape sequence
// (i.e. the row transitions the parser doesn't see directly). Keeps the
// parser's notion of the current screen row in sync with the view.
func (ei *escapeInterpreter) notifyRowAdvance() {
	ei.screenRow++
	ei.screenCol = 1
}

// notifyColumnReset must be called when the view processes a bare CR
// (column reset without row advance). Keeps screenCol in sync so wrap
// counting starts over from col 1.
func (ei *escapeInterpreter) notifyColumnReset() {
	ei.screenCol = 1
}

// notifyCellsWritten must be called after the view writes visible cells
// to its buffer. Advances the parser's idea of the cursor by `width`
// columns; if that crosses the right edge of a `screenColMax`-wide pty
// screen, the corresponding number of soft-wraps are added to screenRow
// so subsequent CUPs land on the right line.
func (ei *escapeInterpreter) notifyCellsWritten(width int) {
	if ei.screenColMax <= 0 {
		return
	}
	// One column at a time: matches ConPTY's "pending wrap" semantics
	// where the cursor stays at col max+1 after writing the rightmost
	// cell and only wraps on the next cell. Loops over individual
	// columns rather than doing the math in one shot so wide cells on a
	// row boundary still wrap cleanly.
	for range width {
		if ei.screenCol > ei.screenColMax {
			ei.screenRow++
			ei.screenCol = 1
		}
		ei.screenCol++
	}
}

// emitCursorAdvance schedules a cursorDown instruction for the next time
// the view checks ei.instruction, advancing the parser's screen row by
// the same amount. n <= 0 is a no-op (backward / same-row CUPs are
// ignored — the view's buffer is line-based and can't undo).
func (ei *escapeInterpreter) emitCursorAdvance(n int) {
	if n <= 0 {
		return
	}
	ei.instruction = cursorDown{n: n}
	ei.screenRow += n
	ei.screenCol = 1
}

// firstParamOrDefault returns the first CSI parameter parsed as an int,
// or dflt if it's absent / empty / unparseable.
func (ei *escapeInterpreter) firstParamOrDefault(dflt int) int {
	if len(ei.csiParam) == 0 || ei.csiParam[0] == "" {
		return dflt
	}
	n, err := strconv.Atoi(ei.csiParam[0])
	if err != nil {
		return dflt
	}
	return n
}

func (ei *escapeInterpreter) instructionRead() {
	ei.instruction = noInstruction{}
}

// parseOne parses a character (grapheme cluster). If isEscape is true, it means that the character
// is part of an escape sequence, and as such should not be printed verbatim. Otherwise, it's not an
// escape sequence.
func (ei *escapeInterpreter) parseOne(ch []byte) (isEscape bool, err error) {
	// Sanity checks: if a sequence has grown absurdly long, stop
	// accumulating state and just swallow bytes until its final byte —
	// much better than leaking the accumulated garbage into the view.
	if len(ei.csiParam) > 20 || (len(ei.csiParam) > 0 && len(ei.csiParam[len(ei.csiParam)-1]) > 255) {
		ei.state = stateCSIDiscard
		ei.csiParam = nil
		return true, nil
	}

	ei.curch = string(ch)

	switch ei.state {
	case stateNone:
		if characterEquals(ch, 0x1b) {
			ei.state = stateEscape
			return true, nil
		}
		return false, nil
	case stateEscape:
		switch {
		case characterEquals(ch, '['):
			ei.state = stateCSI
			return true, nil
		case characterEquals(ch, ']'):
			ei.state = stateOSC
			return true, nil
		case characterEquals(ch, '('),
			characterEquals(ch, ')'),
			characterEquals(ch, '*'),
			characterEquals(ch, '+'):
			ei.state = stateCharacterSetDesignation
			return true, nil
		case len(ch) == 1 && ch[0] >= 0x30 && ch[0] <= 0x7E:
			// Single-byte ESC sequence (e.g. ESC c = RIS). We don't
			// interpret these, but we must consume them so they don't
			// leak into the view as literal text.
			ei.state = stateNone
			return true, nil
		default:
			return false, errNotCSI
		}
	case stateCharacterSetDesignation:
		// Not supported, so just skip it
		ei.state = stateNone
		return true, nil
	case stateCSI:
		switch {
		case len(ch) == 1 && ch[0] >= '0' && ch[0] <= '9':
			ei.csiParam = append(ei.csiParam, "")
		case characterEquals(ch, 'm'):
			ei.csiParam = append(ei.csiParam, "0")
		case characterEquals(ch, 'K'),
			characterEquals(ch, 'H'), characterEquals(ch, 'f'), characterEquals(ch, 'd'),
			characterEquals(ch, 'B'), characterEquals(ch, 'E'),
			characterEquals(ch, 'C'):
			// fall through — let stateParams handle these with default
			// params (CUP/VPA default to row 1, CUD/CNL/CUF default to
			// advance by 1).
		case characterEquals(ch, ';'):
			// Empty first param ([;Xm ≡ [0;Xm). Seed a slot for the
			// empty param; stateParams will append the next one when it
			// re-reads this ';' via the fallthrough.
			ei.csiParam = append(ei.csiParam, "")
		case len(ch) == 1 && ch[0] >= 0x3C && ch[0] <= 0x3F:
			// Private-mode prefix byte (<, =, >, ?). We don't interpret
			// DEC private-mode sequences, but must consume them so they
			// don't leak into the view as literal text. Seed an empty
			// param so the subsequent digits land on a valid slot.
			ei.csiParam = append(ei.csiParam, "")
			ei.state = stateParams
			return true, nil
		case len(ch) == 1 && ch[0] >= 0x20 && ch[0] <= 0x2F:
			// CSI intermediate byte. A sequence with intermediates is
			// one we don't implement; consume the rest until the final
			// byte.
			ei.state = stateCSIDiscard
			ei.csiParam = nil
			return true, nil
		case len(ch) == 1 && ch[0] >= 0x40 && ch[0] <= 0x7E:
			// Valid CSI final byte we don't implement — swallow.
			ei.state = stateNone
			ei.csiParam = nil
			return true, nil
		default:
			return false, errCSIParseError
		}
		ei.state = stateParams
		fallthrough
	case stateParams:
		switch {
		case len(ch) == 1 && ch[0] >= '0' && ch[0] <= '9':
			ei.csiParam[len(ei.csiParam)-1] += string(ch)
			return true, nil
		case characterEquals(ch, ';'):
			ei.csiParam = append(ei.csiParam, "")
			return true, nil
		case characterEquals(ch, 'm'):
			// outputCSI applies params left-to-right and mutates as it
			// goes, so on failure some leading params may already have
			// taken effect (e.g. `[1;;m` would leave AttrBold set before
			// hitting the empty param). Snapshot the colors beforehand
			// and restore them on error so a malformed SGR is truly a
			// no-op rather than a partial apply.
			savedFg, savedBg := ei.curFgColor, ei.curBgColor
			if err := ei.outputCSI(); err != nil {
				ei.curFgColor, ei.curBgColor = savedFg, savedBg
			}
			ei.state = stateNone
			ei.csiParam = nil
			return true, nil
		case characterEquals(ch, 'K'):
			p := 0
			if len(ei.csiParam) != 0 && ei.csiParam[0] != "" {
				p, err = strconv.Atoi(ei.csiParam[0])
				if err != nil {
					return false, errCSIParseError
				}
			}

			if p == 0 {
				ei.instruction = eraseInLineFromCursor{}
			} else {
				// non-zero values of P not supported
				ei.instruction = noInstruction{}
			}

			ei.state = stateNone
			ei.csiParam = nil
			return true, nil
		case characterEquals(ch, 'H'), characterEquals(ch, 'f'),
			characterEquals(ch, 'd'):
			// CUP / HVP (absolute (row, col), col ignored) or VPA (absolute row).
			targetRow := ei.firstParamOrDefault(1)
			if targetRow <= 1 {
				// Cursor home. ConPTY emits this (after [2J) at the start of
				// every screen, so it marks where ConPTY's coordinate origin
				// now sits. Re-anchor our row tracking to the current write
				// position rather than treating it as a backward move: a view
				// that isn't rewound in lockstep with ConPTY's screen (the
				// command log) would otherwise carry stale drift, making every
				// later absolute CUP compute a negative, dropped advance and
				// collapsing the blank rows ConPTY positioned with.
				ei.resetScreenCursor()
			} else {
				// Skip forward to the target row; ignore backward moves.
				ei.emitCursorAdvance(targetRow - ei.screenRow)
			}
			ei.state = stateNone
			ei.csiParam = nil
			return true, nil
		case characterEquals(ch, 'B'), characterEquals(ch, 'E'):
			// CUD / CNL — relative row advance by N. CNL also resets
			// the column, which we don't track, so the two are
			// equivalent for our purposes.
			ei.emitCursorAdvance(ei.firstParamOrDefault(1))
			ei.state = stateNone
			ei.csiParam = nil
			return true, nil
		case characterEquals(ch, 'C'):
			// CUF — cursor forward N. Emit space cells so the gap
			// renders. (screenCol is updated by the view via
			// notifyCellsWritten as those spaces are emitted.)
			ei.instruction = cursorForward{n: ei.firstParamOrDefault(1)}
			ei.state = stateNone
			ei.csiParam = nil
			return true, nil
		case len(ch) == 1 && ch[0] >= 0x20 && ch[0] <= 0x2F:
			// CSI intermediate byte after params. The final byte will
			// have a semantic we don't implement (e.g. `[0 q` =
			// DECSCUSR); consume everything until it arrives.
			ei.state = stateCSIDiscard
			ei.csiParam = nil
			return true, nil
		case len(ch) == 1 && ch[0] >= 0x40 && ch[0] <= 0x7E:
			// Valid CSI final byte we don't implement — swallow the
			// whole sequence rather than printing it as text.
			ei.state = stateNone
			ei.csiParam = nil
			return true, nil
		default:
			return false, errCSIParseError
		}
	case stateCSIDiscard:
		// Consume the rest of a CSI sequence whose semantic we don't
		// interpret (one with intermediate bytes, or one the sanity
		// checks at the top of parseOne bailed out of). Any byte in the
		// final-byte range ends it.
		if len(ch) == 1 && ch[0] >= 0x40 && ch[0] <= 0x7E {
			ei.state = stateNone
		}
		return true, nil
	case stateOSC:
		// Accumulate the OSC number until its terminating ';', then dispatch on
		// it. (The previous code only recognised the single-digit '8'; a number
		// like 1717 needs more than one character.)
		switch {
		case len(ch) == 1 && ch[0] >= '0' && ch[0] <= '9':
			ei.oscNumber.WriteByte(ch[0])
			return true, nil
		case characterEquals(ch, ';'):
			switch ei.oscNumber.String() {
			case "8":
				ei.hyperlink.Reset()
				ei.state = stateOSCParams
			case "1717":
				ei.metadata.Reset()
				ei.state = stateOSCMetadata
			default:
				ei.state = stateOSCSkipUnknown
			}
			ei.oscNumber.Reset()
			return true, nil
		default:
			// Not a recognized OSC; skip to its terminator (handling the case
			// where this character already is one).
			ei.oscNumber.Reset()
			switch {
			case characterEquals(ch, 0x07):
				ei.state = stateNone
			case characterEquals(ch, 0x1b):
				ei.state = stateOSCEndEscape
			default:
				ei.state = stateOSCSkipUnknown
			}
			return true, nil
		}
	case stateOSCParams:
		if characterEquals(ch, ';') {
			ei.state = stateOSCHyperlink
		}
		return true, nil
	case stateOSCHyperlink:
		switch {
		case characterEquals(ch, 0x07):
			ei.state = stateNone
		case characterEquals(ch, 0x1b):
			ei.state = stateOSCEndEscape
		default:
			ei.hyperlink.Write(ch)
		}
		return true, nil
	case stateOSCMetadata:
		switch {
		case characterEquals(ch, 0x07):
			ei.state = stateNone
		case characterEquals(ch, 0x1b):
			ei.state = stateOSCEndEscape
		default:
			ei.metadata.Write(ch)
		}
		return true, nil
	case stateOSCEndEscape:
		ei.state = stateNone
		return true, nil
	case stateOSCSkipUnknown:
		switch {
		case characterEquals(ch, 0x07):
			ei.state = stateNone
		case characterEquals(ch, 0x1b):
			ei.state = stateOSCEndEscape
		}
		return true, nil
	}
	return false, nil
}

func (ei *escapeInterpreter) outputCSI() error {
	n := len(ei.csiParam)
	for i := 0; i < n; {
		p, err := strconv.Atoi(ei.csiParam[i])
		if err != nil {
			return errCSIParseError
		}

		skip := 1
		switch {
		case p == 0: // reset style and color
			ei.curFgColor = ColorDefault
			ei.curBgColor = ColorDefault
		case p >= 1 && p <= 9: // set style
			ei.curFgColor |= getFontEffect(p)
		case p >= 21 && p <= 29: // reset style
			ei.curFgColor &= ^getFontEffect(p - 20)
		case p >= 30 && p <= 37: // set foreground color
			ei.curFgColor &= AttrStyleBits
			ei.curFgColor |= Get256Color(int32(p) - 30)
		case p == setForegroundColor: // set foreground color (256-color or true color)
			var color Attribute
			var err error
			color, skip, err = ei.csiColor(ei.csiParam[i:])
			if err != nil {
				return err
			}
			ei.curFgColor &= AttrStyleBits
			ei.curFgColor |= color
		case p == defaultForegroundColor: // reset foreground color
			ei.curFgColor &= AttrStyleBits
			ei.curFgColor |= ColorDefault
		case p >= 40 && p <= 47: // set background color
			ei.curBgColor &= AttrStyleBits
			ei.curBgColor |= Get256Color(int32(p) - 40)
		case p == setBackgroundColor: // set background color (256-color or true color)
			var color Attribute
			var err error
			color, skip, err = ei.csiColor(ei.csiParam[i:])
			if err != nil {
				return err
			}
			ei.curBgColor &= AttrStyleBits
			ei.curBgColor |= color
		case p == defaultBackgroundColor: // reset background color
			ei.curBgColor &= AttrStyleBits
			ei.curBgColor |= ColorDefault
		case p >= 90 && p <= 97: // set bright foreground color
			ei.curFgColor &= AttrStyleBits
			ei.curFgColor |= Get256Color(int32(p) - 90 + 8)
		case p >= 100 && p <= 107: // set bright background color
			ei.curBgColor &= AttrStyleBits
			ei.curBgColor |= Get256Color(int32(p) - 100 + 8)
		default:
		}
		i += skip
	}

	return nil
}

func (ei *escapeInterpreter) csiColor(param []string) (color Attribute, skip int, err error) {
	if len(param) < 2 {
		return 0, 0, errCSIParseError
	}

	switch param[1] {
	case "2":
		// 24-bit color
		if ei.mode < OutputTrue {
			return 0, 0, errCSIParseError
		}
		if len(param) < 5 {
			return 0, 0, errCSIParseError
		}
		var red, green, blue int
		red, err = strconv.Atoi(param[2])
		if err != nil {
			return 0, 0, errCSIParseError
		}
		green, err = strconv.Atoi(param[3])
		if err != nil {
			return 0, 0, errCSIParseError
		}
		blue, err = strconv.Atoi(param[4])
		if err != nil {
			return 0, 0, errCSIParseError
		}
		return NewRGBColor(int32(red), int32(green), int32(blue)), 5, nil
	case "5":
		// 8-bit color
		if ei.mode < Output256 {
			return 0, 0, errCSIParseError
		}
		if len(param) < 3 {
			return 0, 0, errCSIParseError
		}
		var hex int
		hex, err = strconv.Atoi(param[2])
		if err != nil {
			return 0, 0, errCSIParseError
		}
		return Get256Color(int32(hex)), 3, nil
	default:
		return 0, 0, errCSIParseError
	}
}

func getFontEffect(f int) Attribute {
	switch fontEffect(f) {
	case bold:
		return AttrBold
	case faint:
		return AttrDim
	case italic:
		return AttrItalic
	case underline:
		return AttrUnderline
	case blink:
		return AttrBlink
	case reverse:
		return AttrReverse
	case strike:
		return AttrStrikeThrough
	}
	return AttrNone
}
