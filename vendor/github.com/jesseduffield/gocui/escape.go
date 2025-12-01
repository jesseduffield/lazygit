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
}

type (
	escapeState int
	fontEffect  int
)

type instruction interface{ isInstruction() }

type eraseInLineFromCursor struct{}

func (self eraseInLineFromCursor) isInstruction() {}

type noInstruction struct{}

func (self noInstruction) isInstruction() {}

const (
	stateNone escapeState = iota
	stateEscape
	stateCSI
	stateParams
	stateOSC
	stateOSCWaitForParams
	stateOSCParams
	stateOSCHyperlink
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
	errCSITooLong    = errors.New("CSI escape sequence is too long")
	errOSCParseError = errors.New("OSC escape sequence parsing error")
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
	}
	return ei
}

// reset sets the escapeInterpreter in initial state.
func (ei *escapeInterpreter) reset() {
	ei.state = stateNone
	ei.curFgColor = ColorDefault
	ei.curBgColor = ColorDefault
	ei.csiParam = nil
}

func (ei *escapeInterpreter) instructionRead() {
	ei.instruction = noInstruction{}
}

// parseOne parses a character (grapheme cluster). If isEscape is true, it means that the character
// is part of an escape sequence, and as such should not be printed verbatim. Otherwise, it's not an
// escape sequence.
func (ei *escapeInterpreter) parseOne(ch []byte) (isEscape bool, err error) {
	// Sanity checks
	if len(ei.csiParam) > 20 {
		return false, errCSITooLong
	}
	if len(ei.csiParam) > 0 && len(ei.csiParam[len(ei.csiParam)-1]) > 255 {
		return false, errCSITooLong
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
		default:
			return false, errNotCSI
		}
	case stateCSI:
		switch {
		case len(ch) == 1 && ch[0] >= '0' && ch[0] <= '9':
			ei.csiParam = append(ei.csiParam, "")
		case characterEquals(ch, 'm'):
			ei.csiParam = append(ei.csiParam, "0")
		case characterEquals(ch, 'K'):
			// fall through
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
			if err := ei.outputCSI(); err != nil {
				return false, errCSIParseError
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
		default:
			return false, errCSIParseError
		}
	case stateOSC:
		if characterEquals(ch, '8') {
			ei.state = stateOSCWaitForParams
			ei.hyperlink.Reset()
			return true, nil
		}

		ei.state = stateOSCSkipUnknown
		return true, nil
	case stateOSCWaitForParams:
		if !characterEquals(ch, ';') {
			return true, errOSCParseError
		}

		ei.state = stateOSCParams
		return true, nil
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
