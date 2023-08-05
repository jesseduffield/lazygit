// Copyright 2022 The TCell Authors
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

package terminfo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// ErrTermNotFound indicates that a suitable terminal entry could
	// not be found.  This can result from either not having TERM set,
	// or from the TERM failing to support certain minimal functionality,
	// in particular absolute cursor addressability (the cup capability)
	// is required.  For example, legacy "adm3" lacks this capability,
	// whereas the slightly newer "adm3a" supports it.  This failure
	// occurs most often with "dumb".
	ErrTermNotFound = errors.New("terminal entry not found")
)

// Terminfo represents a terminfo entry.  Note that we use friendly names
// in Go, but when we write out JSON, we use the same names as terminfo.
// The name, aliases and smous, rmous fields do not come from terminfo directly.
type Terminfo struct {
	Name         string
	Aliases      []string
	Columns      int    // cols
	Lines        int    // lines
	Colors       int    // colors
	Bell         string // bell
	Clear        string // clear
	EnterCA      string // smcup
	ExitCA       string // rmcup
	ShowCursor   string // cnorm
	HideCursor   string // civis
	AttrOff      string // sgr0
	Underline    string // smul
	Bold         string // bold
	Blink        string // blink
	Reverse      string // rev
	Dim          string // dim
	Italic       string // sitm
	EnterKeypad  string // smkx
	ExitKeypad   string // rmkx
	SetFg        string // setaf
	SetBg        string // setab
	ResetFgBg    string // op
	SetCursor    string // cup
	CursorBack1  string // cub1
	CursorUp1    string // cuu1
	PadChar      string // pad
	KeyBackspace string // kbs
	KeyF1        string // kf1
	KeyF2        string // kf2
	KeyF3        string // kf3
	KeyF4        string // kf4
	KeyF5        string // kf5
	KeyF6        string // kf6
	KeyF7        string // kf7
	KeyF8        string // kf8
	KeyF9        string // kf9
	KeyF10       string // kf10
	KeyF11       string // kf11
	KeyF12       string // kf12
	KeyF13       string // kf13
	KeyF14       string // kf14
	KeyF15       string // kf15
	KeyF16       string // kf16
	KeyF17       string // kf17
	KeyF18       string // kf18
	KeyF19       string // kf19
	KeyF20       string // kf20
	KeyF21       string // kf21
	KeyF22       string // kf22
	KeyF23       string // kf23
	KeyF24       string // kf24
	KeyF25       string // kf25
	KeyF26       string // kf26
	KeyF27       string // kf27
	KeyF28       string // kf28
	KeyF29       string // kf29
	KeyF30       string // kf30
	KeyF31       string // kf31
	KeyF32       string // kf32
	KeyF33       string // kf33
	KeyF34       string // kf34
	KeyF35       string // kf35
	KeyF36       string // kf36
	KeyF37       string // kf37
	KeyF38       string // kf38
	KeyF39       string // kf39
	KeyF40       string // kf40
	KeyF41       string // kf41
	KeyF42       string // kf42
	KeyF43       string // kf43
	KeyF44       string // kf44
	KeyF45       string // kf45
	KeyF46       string // kf46
	KeyF47       string // kf47
	KeyF48       string // kf48
	KeyF49       string // kf49
	KeyF50       string // kf50
	KeyF51       string // kf51
	KeyF52       string // kf52
	KeyF53       string // kf53
	KeyF54       string // kf54
	KeyF55       string // kf55
	KeyF56       string // kf56
	KeyF57       string // kf57
	KeyF58       string // kf58
	KeyF59       string // kf59
	KeyF60       string // kf60
	KeyF61       string // kf61
	KeyF62       string // kf62
	KeyF63       string // kf63
	KeyF64       string // kf64
	KeyInsert    string // kich1
	KeyDelete    string // kdch1
	KeyHome      string // khome
	KeyEnd       string // kend
	KeyHelp      string // khlp
	KeyPgUp      string // kpp
	KeyPgDn      string // knp
	KeyUp        string // kcuu1
	KeyDown      string // kcud1
	KeyLeft      string // kcub1
	KeyRight     string // kcuf1
	KeyBacktab   string // kcbt
	KeyExit      string // kext
	KeyClear     string // kclr
	KeyPrint     string // kprt
	KeyCancel    string // kcan
	Mouse        string // kmous
	AltChars     string // acsc
	EnterAcs     string // smacs
	ExitAcs      string // rmacs
	EnableAcs    string // enacs
	KeyShfRight  string // kRIT
	KeyShfLeft   string // kLFT
	KeyShfHome   string // kHOM
	KeyShfEnd    string // kEND
	KeyShfInsert string // kIC
	KeyShfDelete string // kDC

	// These are non-standard extensions to terminfo.  This includes
	// true color support, and some additional keys.  Its kind of bizarre
	// that shifted variants of left and right exist, but not up and down.
	// Terminal support for these are going to vary amongst XTerm
	// emulations, so don't depend too much on them in your application.

	StrikeThrough           string // smxx
	SetFgBg                 string // setfgbg
	SetFgBgRGB              string // setfgbgrgb
	SetFgRGB                string // setfrgb
	SetBgRGB                string // setbrgb
	KeyShfUp                string // shift-up
	KeyShfDown              string // shift-down
	KeyShfPgUp              string // shift-kpp
	KeyShfPgDn              string // shift-knp
	KeyCtrlUp               string // ctrl-up
	KeyCtrlDown             string // ctrl-left
	KeyCtrlRight            string // ctrl-right
	KeyCtrlLeft             string // ctrl-left
	KeyMetaUp               string // meta-up
	KeyMetaDown             string // meta-left
	KeyMetaRight            string // meta-right
	KeyMetaLeft             string // meta-left
	KeyAltUp                string // alt-up
	KeyAltDown              string // alt-left
	KeyAltRight             string // alt-right
	KeyAltLeft              string // alt-left
	KeyCtrlHome             string
	KeyCtrlEnd              string
	KeyMetaHome             string
	KeyMetaEnd              string
	KeyAltHome              string
	KeyAltEnd               string
	KeyAltShfUp             string
	KeyAltShfDown           string
	KeyAltShfLeft           string
	KeyAltShfRight          string
	KeyMetaShfUp            string
	KeyMetaShfDown          string
	KeyMetaShfLeft          string
	KeyMetaShfRight         string
	KeyCtrlShfUp            string
	KeyCtrlShfDown          string
	KeyCtrlShfLeft          string
	KeyCtrlShfRight         string
	KeyCtrlShfHome          string
	KeyCtrlShfEnd           string
	KeyAltShfHome           string
	KeyAltShfEnd            string
	KeyMetaShfHome          string
	KeyMetaShfEnd           string
	EnablePaste             string // bracketed paste mode
	DisablePaste            string
	PasteStart              string
	PasteEnd                string
	Modifiers               int
	InsertChar              string // string to insert a character (ich1)
	AutoMargin              bool   // true if writing to last cell in line advances
	TrueColor               bool   // true if the terminal supports direct color
	CursorDefault           string
	CursorBlinkingBlock     string
	CursorSteadyBlock       string
	CursorBlinkingUnderline string
	CursorSteadyUnderline   string
	CursorBlinkingBar       string
	CursorSteadyBar         string
	EnterUrl                string
	ExitUrl                 string
	SetWindowSize           string
	EnableFocusReporting    string
	DisableFocusReporting   string
}

const (
	ModifiersNone  = 0
	ModifiersXTerm = 1
)

type stack []interface{}

func (st stack) Push(v interface{}) stack {
	if b, ok := v.(bool); ok {
		if b {
			return append(st, 1)
		} else {
			return append(st, 0)
		}
	}
	return append(st, v)
}

func (st stack) PopString() (string, stack) {
	if len(st) > 0 {
		e := st[len(st)-1]
		var s string
		switch v := e.(type) {
		case int:
			s = strconv.Itoa(v)
		case string:
			s = v
		}
		return s, st[:len(st)-1]
	}
	return "", st

}
func (st stack) PopInt() (int, stack) {
	if len(st) > 0 {
		e := st[len(st)-1]
		var i int
		switch v := e.(type) {
		case int:
			i = v
		case string:
			i, _ = strconv.Atoi(v)
		}
		return i, st[:len(st)-1]
	}
	return 0, st
}

// static vars
var svars [26]string

type paramsBuffer struct {
	out bytes.Buffer
	buf bytes.Buffer
}

// Start initializes the params buffer with the initial string data.
// It also locks the paramsBuffer.  The caller must call End() when
// finished.
func (pb *paramsBuffer) Start(s string) {
	pb.out.Reset()
	pb.buf.Reset()
	pb.buf.WriteString(s)
}

// End returns the final output from TParam, but it also releases the lock.
func (pb *paramsBuffer) End() string {
	s := pb.out.String()
	return s
}

// NextCh returns the next input character to the expander.
func (pb *paramsBuffer) NextCh() (byte, error) {
	return pb.buf.ReadByte()
}

// PutCh "emits" (rather schedules for output) a single byte character.
func (pb *paramsBuffer) PutCh(ch byte) {
	pb.out.WriteByte(ch)
}

// PutString schedules a string for output.
func (pb *paramsBuffer) PutString(s string) {
	pb.out.WriteString(s)
}

// TParm takes a terminfo parameterized string, such as setaf or cup, and
// evaluates the string, and returns the result with the parameter
// applied.
func (t *Terminfo) TParm(s string, p ...interface{}) string {
	var stk stack
	var a string
	var ai, bi int
	var dvars [26]string
	var params [9]interface{}
	var pb = &paramsBuffer{}

	pb.Start(s)

	// make sure we always have 9 parameters -- makes it easier
	// later to skip checks
	for i := 0; i < len(params) && i < len(p); i++ {
		params[i] = p[i]
	}

	const (
		emit = iota
		toEnd
		toElse
	)

	skip := emit

	for {

		ch, err := pb.NextCh()
		if err != nil {
			break
		}

		if ch != '%' {
			if skip == emit {
				pb.PutCh(ch)
			}
			continue
		}

		ch, err = pb.NextCh()
		if err != nil {
			// XXX Error
			break
		}
		if skip == toEnd {
			if ch == ';' {
				skip = emit
			}
			continue
		} else if skip == toElse {
			if ch == 'e' || ch == ';' {
				skip = emit
			}
			continue
		}

		switch ch {
		case '%': // quoted %
			pb.PutCh(ch)

		case 'i': // increment both parameters (ANSI cup support)
			if i, ok := params[0].(int); ok {
				params[0] = i + 1
			}
			if i, ok := params[1].(int); ok {
				params[1] = i + 1
			}

		case 's':
			// NB: 's', 'c', and 'd' below are special cased for
			// efficiency.  They could be handled by the richer
			// format support below, less efficiently.
			a, stk = stk.PopString()
			pb.PutString(a)

		case 'c':
			// Integer as special character.
			ai, stk = stk.PopInt()
			pb.PutCh(byte(ai))

		case 'd':
			ai, stk = stk.PopInt()
			pb.PutString(strconv.Itoa(ai))

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'x', 'X', 'o', ':':
			// This is pretty suboptimal, but this is rarely used.
			// None of the mainstream terminals use any of this,
			// and it would surprise me if this code is ever
			// executed outside test cases.
			f := "%"
			if ch == ':' {
				ch, _ = pb.NextCh()
			}
			f += string(ch)
			for ch == '+' || ch == '-' || ch == '#' || ch == ' ' {
				ch, _ = pb.NextCh()
				f += string(ch)
			}
			for (ch >= '0' && ch <= '9') || ch == '.' {
				ch, _ = pb.NextCh()
				f += string(ch)
			}
			switch ch {
			case 'd', 'x', 'X', 'o':
				ai, stk = stk.PopInt()
				pb.PutString(fmt.Sprintf(f, ai))
			case 's':
				a, stk = stk.PopString()
				pb.PutString(fmt.Sprintf(f, a))
			case 'c':
				ai, stk = stk.PopInt()
				pb.PutString(fmt.Sprintf(f, ai))
			}

		case 'p': // push parameter
			ch, _ = pb.NextCh()
			ai = int(ch - '1')
			if ai >= 0 && ai < len(params) {
				stk = stk.Push(params[ai])
			} else {
				stk = stk.Push(0)
			}

		case 'P': // pop & store variable
			ch, _ = pb.NextCh()
			if ch >= 'A' && ch <= 'Z' {
				svars[int(ch-'A')], stk = stk.PopString()
			} else if ch >= 'a' && ch <= 'z' {
				dvars[int(ch-'a')], stk = stk.PopString()
			}

		case 'g': // recall & push variable
			ch, _ = pb.NextCh()
			if ch >= 'A' && ch <= 'Z' {
				stk = stk.Push(svars[int(ch-'A')])
			} else if ch >= 'a' && ch <= 'z' {
				stk = stk.Push(dvars[int(ch-'a')])
			}

		case '\'': // push(char) - the integer value of it
			ch, _ = pb.NextCh()
			_, _ = pb.NextCh() // must be ' but we don't check
			stk = stk.Push(int(ch))

		case '{': // push(int)
			ai = 0
			ch, _ = pb.NextCh()
			for ch >= '0' && ch <= '9' {
				ai *= 10
				ai += int(ch - '0')
				ch, _ = pb.NextCh()
			}
			// ch must be '}' but no verification
			stk = stk.Push(ai)

		case 'l': // push(strlen(pop))
			a, stk = stk.PopString()
			stk = stk.Push(len(a))

		case '+':
			bi, stk = stk.PopInt()
			ai, stk = stk.PopInt()
			stk = stk.Push(ai + bi)

		case '-':
			bi, stk = stk.PopInt()
			ai, stk = stk.PopInt()
			stk = stk.Push(ai - bi)

		case '*':
			bi, stk = stk.PopInt()
			ai, stk = stk.PopInt()
			stk = stk.Push(ai * bi)

		case '/':
			bi, stk = stk.PopInt()
			ai, stk = stk.PopInt()
			if bi != 0 {
				stk = stk.Push(ai / bi)
			} else {
				stk = stk.Push(0)
			}

		case 'm': // push(pop mod pop)
			bi, stk = stk.PopInt()
			ai, stk = stk.PopInt()
			if bi != 0 {
				stk = stk.Push(ai % bi)
			} else {
				stk = stk.Push(0)
			}

		case '&': // AND
			bi, stk = stk.PopInt()
			ai, stk = stk.PopInt()
			stk = stk.Push(ai & bi)

		case '|': // OR
			bi, stk = stk.PopInt()
			ai, stk = stk.PopInt()
			stk = stk.Push(ai | bi)

		case '^': // XOR
			bi, stk = stk.PopInt()
			ai, stk = stk.PopInt()
			stk = stk.Push(ai ^ bi)

		case '~': // bit complement
			ai, stk = stk.PopInt()
			stk = stk.Push(ai ^ -1)

		case '!': // logical NOT
			ai, stk = stk.PopInt()
			stk = stk.Push(ai == 0)

		case '=': // numeric compare
			bi, stk = stk.PopInt()
			ai, stk = stk.PopInt()
			stk = stk.Push(ai == bi)

		case '>': // greater than, numeric
			bi, stk = stk.PopInt()
			ai, stk = stk.PopInt()
			stk = stk.Push(ai > bi)

		case '<': // less than, numeric
			bi, stk = stk.PopInt()
			ai, stk = stk.PopInt()
			stk = stk.Push(ai < bi)

		case '?': // start conditional

		case ';':
			skip = emit

		case 't':
			ai, stk = stk.PopInt()
			if ai == 0 {
				skip = toElse
			}

		case 'e':
			skip = toEnd

		default:
			pb.PutString("%" + string(ch))
		}
	}

	return pb.End()
}

// TPuts emits the string to the writer, but expands inline padding
// indications (of the form $<[delay]> where [delay] is msec) to
// a suitable time (unless the terminfo string indicates this isn't needed
// by specifying npc - no padding).  All Terminfo based strings should be
// emitted using this function.
func (t *Terminfo) TPuts(w io.Writer, s string) {
	for {
		beg := strings.Index(s, "$<")
		if beg < 0 {
			// Most strings don't need padding, which is good news!
			_, _ = io.WriteString(w, s)
			return
		}
		_, _ = io.WriteString(w, s[:beg])
		s = s[beg+2:]
		end := strings.Index(s, ">")
		if end < 0 {
			// unterminated.. just emit bytes unadulterated
			_, _ = io.WriteString(w, "$<"+s)
			return
		}
		val := s[:end]
		s = s[end+1:]
		padus := 0
		unit := time.Millisecond
		dot := false
	loop:
		for i := range val {
			switch val[i] {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				padus *= 10
				padus += int(val[i] - '0')
				if dot {
					unit /= 10
				}
			case '.':
				if !dot {
					dot = true
				} else {
					break loop
				}
			default:
				break loop
			}
		}

		// Curses historically uses padding to achieve "fine grained"
		// delays. We have much better clocks these days, and so we
		// do not rely on padding but simply sleep a bit.
		if len(t.PadChar) > 0 {
			time.Sleep(unit * time.Duration(padus))
		}
	}
}

// TGoto returns a string suitable for addressing the cursor at the given
// row and column.  The origin 0, 0 is in the upper left corner of the screen.
func (t *Terminfo) TGoto(col, row int) string {
	return t.TParm(t.SetCursor, row, col)
}

// TColor returns a string corresponding to the given foreground and background
// colors.  Either fg or bg can be set to -1 to elide.
func (t *Terminfo) TColor(fi, bi int) string {
	rv := ""
	// As a special case, we map bright colors to lower versions if the
	// color table only holds 8.  For the remaining 240 colors, the user
	// is out of luck.  Someday we could create a mapping table, but its
	// not worth it.
	if t.Colors == 8 {
		if fi > 7 && fi < 16 {
			fi -= 8
		}
		if bi > 7 && bi < 16 {
			bi -= 8
		}
	}
	if t.Colors > fi && fi >= 0 {
		rv += t.TParm(t.SetFg, fi)
	}
	if t.Colors > bi && bi >= 0 {
		rv += t.TParm(t.SetBg, bi)
	}
	return rv
}

var (
	dblock    sync.Mutex
	terminfos = make(map[string]*Terminfo)
)

// AddTerminfo can be called to register a new Terminfo entry.
func AddTerminfo(t *Terminfo) {
	dblock.Lock()
	terminfos[t.Name] = t
	for _, x := range t.Aliases {
		terminfos[x] = t
	}
	dblock.Unlock()
}

// LookupTerminfo attempts to find a definition for the named $TERM.
func LookupTerminfo(name string) (*Terminfo, error) {
	if name == "" {
		// else on windows: index out of bounds
		// on the name[0] reference below
		return nil, ErrTermNotFound
	}

	addtruecolor := false
	add256color := false
	switch os.Getenv("COLORTERM") {
	case "truecolor", "24bit", "24-bit":
		addtruecolor = true
	}
	dblock.Lock()
	t := terminfos[name]
	dblock.Unlock()

	// If the name ends in -truecolor, then fabricate an entry
	// from the corresponding -256color, -color, or bare terminal.
	if t != nil && t.TrueColor {
		addtruecolor = true
	} else if t == nil && strings.HasSuffix(name, "-truecolor") {

		suffixes := []string{
			"-256color",
			"-88color",
			"-color",
			"",
		}
		base := name[:len(name)-len("-truecolor")]
		for _, s := range suffixes {
			if t, _ = LookupTerminfo(base + s); t != nil {
				addtruecolor = true
				break
			}
		}
	}

	// If the name ends in -256color, maybe fabricate using the xterm 256 color sequences
	if t == nil && strings.HasSuffix(name, "-256color") {
		suffixes := []string{
			"-88color",
			"-color",
		}
		base := name[:len(name)-len("-256color")]
		for _, s := range suffixes {
			if t, _ = LookupTerminfo(base + s); t != nil {
				add256color = true
				break
			}
		}
	}

	if t == nil {
		return nil, ErrTermNotFound
	}

	switch os.Getenv("TCELL_TRUECOLOR") {
	case "":
	case "disable":
		addtruecolor = false
	default:
		addtruecolor = true
	}

	// If the user has requested 24-bit color with $COLORTERM, then
	// amend the value (unless already present).  This means we don't
	// need to have a value present.
	if addtruecolor &&
		t.SetFgBgRGB == "" &&
		t.SetFgRGB == "" &&
		t.SetBgRGB == "" {

		// Supply vanilla ISO 8613-6:1994 24-bit color sequences.
		t.SetFgRGB = "\x1b[38;2;%p1%d;%p2%d;%p3%dm"
		t.SetBgRGB = "\x1b[48;2;%p1%d;%p2%d;%p3%dm"
		t.SetFgBgRGB = "\x1b[38;2;%p1%d;%p2%d;%p3%d;" +
			"48;2;%p4%d;%p5%d;%p6%dm"
	}

	if add256color {
		t.Colors = 256
		t.SetFg = "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;m"
		t.SetBg = "\x1b[%?%p1%{8}%<%t4%p1%d%e%p1%{16}%<%t10%p1%{8}%-%d%e48;5;%p1%d%;m"
		t.SetFgBg = "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;;%?%p2%{8}%<%t4%p2%d%e%p2%{16}%<%t10%p2%{8}%-%d%e48;5;%p2%d%;m"
		t.ResetFgBg = "\x1b[39;49m"
	}
	return t, nil
}
