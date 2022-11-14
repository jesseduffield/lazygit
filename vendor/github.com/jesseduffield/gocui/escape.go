// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"strconv"

	"github.com/go-errors/errors"
)

type escapeInterpreter struct {
	state                  escapeState
	curch                  rune
	csiParam               []string
	curFgColor, curBgColor Attribute
	mode                   OutputMode
	instruction            instruction
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
	stateOSCEscape

	bold               fontEffect = 1
	faint              fontEffect = 2
	italic             fontEffect = 3
	underline          fontEffect = 4
	blink              fontEffect = 5
	reverse            fontEffect = 7
	strike             fontEffect = 9
	setForegroundColor fontEffect = 38
	setBackgroundColor fontEffect = 48
)

var (
	errNotCSI        = errors.New("Not a CSI escape sequence")
	errCSIParseError = errors.New("CSI escape sequence parsing error")
	errCSITooLong    = errors.New("CSI escape sequence is too long")
)

// runes in case of error will output the non-parsed runes as a string.
func (ei *escapeInterpreter) runes() []rune {
	switch ei.state {
	case stateNone:
		return []rune{0x1b}
	case stateEscape:
		return []rune{0x1b, ei.curch}
	case stateCSI:
		return []rune{0x1b, '[', ei.curch}
	case stateParams:
		ret := []rune{0x1b, '['}
		for _, s := range ei.csiParam {
			ret = append(ret, []rune(s)...)
			ret = append(ret, ';')
		}
		return append(ret, ei.curch)
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

// parseOne parses a rune. If isEscape is true, it means that the rune is part
// of an escape sequence, and as such should not be printed verbatim. Otherwise,
// it's not an escape sequence.
func (ei *escapeInterpreter) parseOne(ch rune) (isEscape bool, err error) {
	// Sanity checks
	if len(ei.csiParam) > 20 {
		return false, errCSITooLong
	}
	if len(ei.csiParam) > 0 && len(ei.csiParam[len(ei.csiParam)-1]) > 255 {
		return false, errCSITooLong
	}

	ei.curch = ch

	switch ei.state {
	case stateNone:
		if ch == 0x1b {
			ei.state = stateEscape
			return true, nil
		}
		return false, nil
	case stateEscape:
		switch ch {
		case '[':
			ei.state = stateCSI
			return true, nil
		case ']':
			ei.state = stateOSC
			return true, nil
		default:
			return false, errNotCSI
		}
	case stateCSI:
		switch {
		case ch >= '0' && ch <= '9':
			ei.csiParam = append(ei.csiParam, "")
		case ch == 'm':
			ei.csiParam = append(ei.csiParam, "0")
		case ch == 'K':
			// fall through
		default:
			return false, errCSIParseError
		}
		ei.state = stateParams
		fallthrough
	case stateParams:
		switch {
		case ch >= '0' && ch <= '9':
			ei.csiParam[len(ei.csiParam)-1] += string(ch)
			return true, nil
		case ch == ';':
			ei.csiParam = append(ei.csiParam, "")
			return true, nil
		case ch == 'm':
			var err error
			switch ei.mode {
			case OutputNormal:
				err = ei.outputNormal()
			case Output256:
				err = ei.output256()
			case OutputTrue:
				err = ei.outputTrue()
			}
			if err != nil {
				return false, errCSIParseError
			}

			ei.state = stateNone
			ei.csiParam = nil
			return true, nil
		case ch == 'K':
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
		switch ch {
		case 0x1b:
			ei.state = stateOSCEscape
			return true, nil
		}
		return true, nil
	case stateOSCEscape:
		ei.state = stateNone
		return true, nil
	}
	return false, nil
}

// outputNormal provides 8 different colors:
//   black, red, green, yellow, blue, magenta, cyan, white
func (ei *escapeInterpreter) outputNormal() error {
	for _, param := range ei.csiParam {
		p, err := strconv.Atoi(param)
		if err != nil {
			return errCSIParseError
		}

		switch {
		case p >= 30 && p <= 37:
			ei.curFgColor = Get256Color(int32(p) - 30)
		case p == 39:
			ei.curFgColor = ColorDefault
		case p >= 40 && p <= 47:
			ei.curBgColor = Get256Color(int32(p) - 40)
		case p == 49:
			ei.curBgColor = ColorDefault
		case p == 0:
			ei.curFgColor = ColorDefault
			ei.curBgColor = ColorDefault
		case p >= 21 && p <= 29:
			ei.curFgColor &= ^getFontEffect(p - 20)
		default:
			ei.curFgColor |= getFontEffect(p)
		}
	}

	return nil
}

// output256 allows you to leverage the 256-colors terminal mode:
//   0x01 - 0x08: the 8 colors as in OutputNormal
//   0x09 - 0x10: Color* | AttrBold
//   0x11 - 0xe8: 216 different colors
//   0xe9 - 0x1ff: 24 different shades of grey
func (ei *escapeInterpreter) output256() error {
	if len(ei.csiParam) < 3 {
		return ei.outputNormal()
	}

	mode, err := strconv.Atoi(ei.csiParam[1])
	if err != nil {
		return errCSIParseError
	}
	if mode != 5 {
		return ei.outputNormal()
	}

	for _, param := range splitFgBg(ei.csiParam, 3) {
		fgbg, err := strconv.Atoi(param[0])
		if err != nil {
			return errCSIParseError
		}
		color, err := strconv.Atoi(param[2])
		if err != nil {
			return errCSIParseError
		}

		switch fontEffect(fgbg) {
		case setForegroundColor:
			ei.curFgColor = Get256Color(int32(color))

			for _, s := range param[3:] {
				p, err := strconv.Atoi(s)
				if err != nil {
					return errCSIParseError
				}

				ei.curFgColor |= getFontEffect(p)
			}
		case setBackgroundColor:
			ei.curBgColor = Get256Color(int32(color))
		default:
			return errCSIParseError
		}
	}
	return nil
}

// outputTrue allows you to leverage the true-color terminal mode.
//
// Works with rgb ANSI sequence: `\x1b[38;2;<r>;<g>;<b>m`, `\x1b[48;2;<r>;<g>;<b>m`
func (ei *escapeInterpreter) outputTrue() error {
	if len(ei.csiParam) < 5 {
		return ei.output256()
	}

	mode, err := strconv.Atoi(ei.csiParam[1])
	if err != nil {
		return errCSIParseError
	}
	if mode != 2 {
		return ei.output256()
	}

	for _, param := range splitFgBg(ei.csiParam, 5) {
		fgbg, err := strconv.Atoi(param[0])
		if err != nil {
			return errCSIParseError
		}
		colr, err := strconv.Atoi(param[2])
		if err != nil {
			return errCSIParseError
		}
		colg, err := strconv.Atoi(param[3])
		if err != nil {
			return errCSIParseError
		}
		colb, err := strconv.Atoi(param[4])
		if err != nil {
			return errCSIParseError
		}
		color := NewRGBColor(int32(colr), int32(colg), int32(colb))

		switch fontEffect(fgbg) {
		case setForegroundColor:
			ei.curFgColor = color

			for _, s := range param[5:] {
				p, err := strconv.Atoi(s)
				if err != nil {
					return errCSIParseError
				}

				ei.curFgColor |= getFontEffect(p)
			}
		case setBackgroundColor:
			ei.curBgColor = color
		default:
			return errCSIParseError
		}
	}
	return nil
}

// splitFgBg splits foreground and background color according to ANSI sequence.
//
// num (number of segments in ansi) is used to determine if it's 256 mode or rgb mode (3 - 256-color, 5 - rgb-color)
func splitFgBg(params []string, num int) [][]string {
	var out [][]string
	var current []string
	for _, p := range params {
		if len(current) == num && (p == "48" || p == "38") {
			out = append(out, current)
			current = []string{}
		}
		current = append(current, p)
	}

	if len(current) > 0 {
		out = append(out, current)
	}

	return out
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
