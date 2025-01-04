//go:build ignore
// +build ignore

// Copyright 2020 The TCell Authors
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

// This command is used to generate suitable configuration files in either
// go syntax or in JSON.  It defaults to JSON output on stdout.  If no
// term values are specified on the command line, then $TERM is used.
//
// Usage is like this:
//
// mkinfo [-go file.go] [-quiet] [-nofatal] [-I <import>] [-P <pkg}] [<term>...]
//
// -go       specifies Go output into the named file.  Use - for stdout.
// -nofatal  indicates that errors loading definitions should not be fatal
// -P pkg    use the supplied package name
// -I import use the named import instead of github.com/gdamore/tcell/v2/terminfo
//

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2/terminfo"
)

type termcap struct {
	name    string
	desc    string
	aliases []string
	bools   map[string]bool
	nums    map[string]int
	strs    map[string]string
}

func (tc *termcap) getnum(s string) int {
	return (tc.nums[s])
}

func (tc *termcap) getflag(s string) bool {
	return (tc.bools[s])
}

func (tc *termcap) getstr(s string) string {
	return (tc.strs[s])
}

const (
	NONE = iota
	CTRL
	ESC
)

var notaddressable = errors.New("terminal not cursor addressable")

func unescape(s string) string {
	// Various escapes are in \x format.  Control codes are
	// encoded as ^M (carat followed by ASCII equivalent).
	// Escapes are: \e, \E - escape
	//  \0 NULL, \n \l \r \t \b \f \s for equivalent C escape.
	buf := &bytes.Buffer{}
	esc := NONE

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch esc {
		case NONE:
			switch c {
			case '\\':
				esc = ESC
			case '^':
				esc = CTRL
			default:
				buf.WriteByte(c)
			}
		case CTRL:
			buf.WriteByte(c ^ 1<<6)
			esc = NONE
		case ESC:
			switch c {
			case 'E', 'e':
				buf.WriteByte(0x1b)
			case '0', '1', '2', '3', '4', '5', '6', '7':
				if i+2 < len(s) && s[i+1] >= '0' && s[i+1] <= '7' && s[i+2] >= '0' && s[i+2] <= '7' {
					buf.WriteByte(((c - '0') * 64) + ((s[i+1] - '0') * 8) + (s[i+2] - '0'))
					i = i + 2
				} else if c == '0' {
					buf.WriteByte(0)
				}
			case 'n':
				buf.WriteByte('\n')
			case 'r':
				buf.WriteByte('\r')
			case 't':
				buf.WriteByte('\t')
			case 'b':
				buf.WriteByte('\b')
			case 'f':
				buf.WriteByte('\f')
			case 's':
				buf.WriteByte(' ')
			case 'l':
				panic("WTF: weird format: " + s)
			default:
				buf.WriteByte(c)
			}
			esc = NONE
		}
	}
	return (buf.String())
}

func (tc *termcap) setupterm(name string) error {
	cmd := exec.Command("infocmp", "-x", "-1", name)
	output := &bytes.Buffer{}
	cmd.Stdout = output

	tc.strs = make(map[string]string)
	tc.bools = make(map[string]bool)
	tc.nums = make(map[string]int)

	err := cmd.Run()
	if err != nil {
		return err
	}

	// Now parse the output.
	// We get comment lines (starting with "#"), followed by
	// a header line that looks like "<name>|<alias>|...|<desc>"
	// then capabilities, one per line, starting with a tab and ending
	// with a comma and newline.
	lines := strings.Split(output.String(), "\n")
	for len(lines) > 0 && strings.HasPrefix(lines[0], "#") {
		lines = lines[1:]
	}

	// Ditch trailing empty last line
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	header := lines[0]
	if strings.HasSuffix(header, ",") {
		header = header[:len(header)-1]
	}
	names := strings.Split(header, "|")
	tc.name = names[0]
	names = names[1:]
	if len(names) > 0 {
		tc.desc = names[len(names)-1]
		names = names[:len(names)-1]
	}
	tc.aliases = names
	for _, val := range lines[1:] {
		if (!strings.HasPrefix(val, "\t")) ||
			(!strings.HasSuffix(val, ",")) {
			return (errors.New("malformed infocmp: " + val))
		}

		val = val[1:]
		val = val[:len(val)-1]

		if k := strings.SplitN(val, "=", 2); len(k) == 2 {
			tc.strs[k[0]] = unescape(k[1])
		} else if k := strings.SplitN(val, "#", 2); len(k) == 2 {
			if u, err := strconv.ParseUint(k[1], 0, 0); err != nil {
				return (err)
			} else {
				tc.nums[k[0]] = int(u)
			}
		} else {
			tc.bools[val] = true
		}
	}
	return nil
}

// This program is used to collect data from the system's terminfo library,
// and write it into Go source code.  That is, we maintain our terminfo
// capabilities encoded in the program.  It should never need to be run by
// an end user, but developers can use this to add codes for additional
// terminal types.
func getinfo(name string) (*terminfo.Terminfo, string, error) {
	var tc termcap
	if err := tc.setupterm(name); err != nil {
		if err != nil {
			return nil, "", err
		}
	}
	t := &terminfo.Terminfo{}
	// If this is an alias record, then just emit the alias
	t.Name = tc.name
	if t.Name != name {
		return t, "", nil
	}
	t.Aliases = tc.aliases
	t.Colors = tc.getnum("colors")
	t.Columns = tc.getnum("cols")
	t.Lines = tc.getnum("lines")
	t.Bell = tc.getstr("bel")
	t.Clear = tc.getstr("clear")
	t.EnterCA = tc.getstr("smcup")
	t.ExitCA = tc.getstr("rmcup")
	t.ShowCursor = tc.getstr("cnorm")
	t.HideCursor = tc.getstr("civis")
	t.AttrOff = tc.getstr("sgr0")
	t.Underline = tc.getstr("smul")
	t.Bold = tc.getstr("bold")
	t.Blink = tc.getstr("blink")
	t.Dim = tc.getstr("dim")
	t.Italic = tc.getstr("sitm")
	t.Reverse = tc.getstr("rev")
	t.EnterKeypad = tc.getstr("smkx")
	t.ExitKeypad = tc.getstr("rmkx")
	t.SetFg = tc.getstr("setaf")
	t.SetBg = tc.getstr("setab")
	t.ResetFgBg = tc.getstr("op")
	t.SetCursor = tc.getstr("cup")
	t.CursorBack1 = tc.getstr("cub1")
	t.CursorUp1 = tc.getstr("cuu1")
	t.KeyF1 = tc.getstr("kf1")
	t.KeyF2 = tc.getstr("kf2")
	t.KeyF3 = tc.getstr("kf3")
	t.KeyF4 = tc.getstr("kf4")
	t.KeyF5 = tc.getstr("kf5")
	t.KeyF6 = tc.getstr("kf6")
	t.KeyF7 = tc.getstr("kf7")
	t.KeyF8 = tc.getstr("kf8")
	t.KeyF9 = tc.getstr("kf9")
	t.KeyF10 = tc.getstr("kf10")
	t.KeyF11 = tc.getstr("kf11")
	t.KeyF12 = tc.getstr("kf12")
	t.KeyInsert = tc.getstr("kich1")
	t.KeyDelete = tc.getstr("kdch1")
	t.KeyBackspace = tc.getstr("kbs")
	t.KeyHome = tc.getstr("khome")
	t.KeyEnd = tc.getstr("kend")
	t.KeyUp = tc.getstr("kcuu1")
	t.KeyDown = tc.getstr("kcud1")
	t.KeyRight = tc.getstr("kcuf1")
	t.KeyLeft = tc.getstr("kcub1")
	t.KeyPgDn = tc.getstr("knp")
	t.KeyPgUp = tc.getstr("kpp")
	t.KeyBacktab = tc.getstr("kcbt")
	t.KeyExit = tc.getstr("kext")
	t.KeyCancel = tc.getstr("kcan")
	t.KeyPrint = tc.getstr("kprt")
	t.KeyHelp = tc.getstr("khlp")
	t.KeyClear = tc.getstr("kclr")
	t.AltChars = tc.getstr("acsc")
	t.EnterAcs = tc.getstr("smacs")
	t.ExitAcs = tc.getstr("rmacs")
	t.EnableAcs = tc.getstr("enacs")
	t.StrikeThrough = tc.getstr("smxx")
	t.Mouse = tc.getstr("kmous")

	t.Modifiers = terminfo.ModifiersNone

	// Terminfo lacks descriptions for a bunch of modified keys,
	// but modern XTerm and emulators often have them. We detect
	// this based on compatible definitions for shifted right.
	// We also choose to use our modifiers for function keys --
	// the terminfo entries list these all as higher coded escape
	// keys, but it's nicer to match them to modifiers.
	if tc.getstr("kRIT") == "\x1b[1;2C" {
		t.Modifiers = terminfo.ModifiersXTerm
	} else {
		// Lookup high level function keys.
		t.KeyShfInsert = tc.getstr("kIC")
		t.KeyShfDelete = tc.getstr("kDC")
		t.KeyShfRight = tc.getstr("kRIT")
		t.KeyShfLeft = tc.getstr("kLFT")
		t.KeyShfHome = tc.getstr("kHOM")
		t.KeyShfEnd = tc.getstr("kEND")
		t.KeyF13 = tc.getstr("kf13")
		t.KeyF14 = tc.getstr("kf14")
		t.KeyF15 = tc.getstr("kf15")
		t.KeyF16 = tc.getstr("kf16")
		t.KeyF17 = tc.getstr("kf17")
		t.KeyF18 = tc.getstr("kf18")
		t.KeyF19 = tc.getstr("kf19")
		t.KeyF20 = tc.getstr("kf20")
		t.KeyF21 = tc.getstr("kf21")
		t.KeyF22 = tc.getstr("kf22")
		t.KeyF23 = tc.getstr("kf23")
		t.KeyF24 = tc.getstr("kf24")
		t.KeyF25 = tc.getstr("kf25")
		t.KeyF26 = tc.getstr("kf26")
		t.KeyF27 = tc.getstr("kf27")
		t.KeyF28 = tc.getstr("kf28")
		t.KeyF29 = tc.getstr("kf29")
		t.KeyF30 = tc.getstr("kf30")
		t.KeyF31 = tc.getstr("kf31")
		t.KeyF32 = tc.getstr("kf32")
		t.KeyF33 = tc.getstr("kf33")
		t.KeyF34 = tc.getstr("kf34")
		t.KeyF35 = tc.getstr("kf35")
		t.KeyF36 = tc.getstr("kf36")
		t.KeyF37 = tc.getstr("kf37")
		t.KeyF38 = tc.getstr("kf38")
		t.KeyF39 = tc.getstr("kf39")
		t.KeyF40 = tc.getstr("kf40")
		t.KeyF41 = tc.getstr("kf41")
		t.KeyF42 = tc.getstr("kf42")
		t.KeyF43 = tc.getstr("kf43")
		t.KeyF44 = tc.getstr("kf44")
		t.KeyF45 = tc.getstr("kf45")
		t.KeyF46 = tc.getstr("kf46")
		t.KeyF47 = tc.getstr("kf47")
		t.KeyF48 = tc.getstr("kf48")
		t.KeyF49 = tc.getstr("kf49")
		t.KeyF50 = tc.getstr("kf50")
		t.KeyF51 = tc.getstr("kf51")
		t.KeyF52 = tc.getstr("kf52")
		t.KeyF53 = tc.getstr("kf53")
		t.KeyF54 = tc.getstr("kf54")
		t.KeyF55 = tc.getstr("kf55")
		t.KeyF56 = tc.getstr("kf56")
		t.KeyF57 = tc.getstr("kf57")
		t.KeyF58 = tc.getstr("kf58")
		t.KeyF59 = tc.getstr("kf59")
		t.KeyF60 = tc.getstr("kf60")
		t.KeyF61 = tc.getstr("kf61")
		t.KeyF62 = tc.getstr("kf62")
		t.KeyF63 = tc.getstr("kf63")
		t.KeyF64 = tc.getstr("kf64")
	}

	// And the same thing for rxvt.
	// It seems that urxvt at least send ESC as ALT prefix for these,
	// although some places seem to indicate a separate ALT key sequence.
	// Users are encouraged to update to an emulator that more closely
	// matches xterm for better functionality.
	if t.KeyShfRight == "\x1b[c" && t.KeyShfLeft == "\x1b[d" {
		t.KeyShfUp = "\x1b[a"
		t.KeyShfDown = "\x1b[b"
		t.KeyCtrlUp = "\x1b[Oa"
		t.KeyCtrlDown = "\x1b[Ob"
		t.KeyCtrlRight = "\x1b[Oc"
		t.KeyCtrlLeft = "\x1b[Od"
	}
	if t.KeyShfHome == "\x1b[7$" && t.KeyShfEnd == "\x1b[8$" {
		t.KeyCtrlHome = "\x1b[7^"
		t.KeyCtrlEnd = "\x1b[8^"
	}

	// Technically the RGB flag that is provided for xterm-direct is not
	// quite right.  The problem is that the -direct flag that was introduced
	// with ncurses 6.1 requires a parsing for the parameters that we lack.
	// For this case we'll just assume it's XTerm compatible.  Someday this
	// may be incorrect, but right now it is correct, and nobody uses it
	// anyway.
	if tc.getflag("Tc") {
		// This presumes XTerm 24-bit true color.
		t.TrueColor = true
	} else if tc.getflag("RGB") {
		// This is for xterm-direct, which uses a different scheme entirely.
		// (ncurses went a very different direction from everyone else, and
		// so it's unlikely anything is using this definition.)
		t.TrueColor = true
		t.SetBg = "\x1b[%?%p1%{8}%<%t4%p1%d%e%p1%{16}%<%t10%p1%{8}%-%d%e48;5;%p1%d%;m"
		t.SetFg = "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;m"
	}

	// If the kmous entry is present, then we need to record the
	// the codes to enter and exit mouse mode.  Sadly, this is not
	// part of the terminfo databases anywhere that I've found, but
	// is an extension.  The escape codes are documented in the XTerm
	// manual, and all terminals that have kmous are expected to
	// use these same codes, unless explicitly configured otherwise
	// vi XM.  Note that in any event, we only known how to parse either
	// x11 or SGR mouse events -- if your terminal doesn't support one
	// of these two forms, you maybe out of luck.
	t.MouseMode = tc.getstr("XM")
	if t.Mouse != "" && t.MouseMode == "" {
		// we anticipate that all xterm mouse tracking compatible
		// terminals understand mouse tracking (1000), but we hope
		// that those that don't understand any-event tracking (1003)
		// will at least ignore it.  Likewise we hope that terminals
		// that don't understand SGR reporting (1006) just ignore it.
		t.MouseMode = "%?%p1%{1}%=%t%'h'%Pa%e%'l'%Pa%;" +
			"\x1b[?1000%ga%c\x1b[?1002%ga%c\x1b[?1003%ga%c\x1b[?1006%ga%c"
	}

	// We only support colors in ANSI 8 or 256 color mode.
	if t.Colors < 8 || t.SetFg == "" {
		t.Colors = 0
	}
	if t.SetCursor == "" {
		return nil, "", notaddressable
	}

	// For padding, we lookup the pad char.  If that isn't present,
	// and npc is *not* set, then we assume a null byte.
	t.PadChar = tc.getstr("pad")
	if t.PadChar == "" {
		if !tc.getflag("npc") {
			t.PadChar = "\u0000"
		}
	}

	// For terminals that use "standard" SGR sequences, lets combine the
	// foreground and background together.
	if strings.HasPrefix(t.SetFg, "\x1b[") &&
		strings.HasPrefix(t.SetBg, "\x1b[") &&
		strings.HasSuffix(t.SetFg, "m") &&
		strings.HasSuffix(t.SetBg, "m") {
		fg := t.SetFg[:len(t.SetFg)-1]
		r := regexp.MustCompile("%p1")
		bg := r.ReplaceAllString(t.SetBg[2:], "%p2")
		t.SetFgBg = fg + ";" + bg
	}

	return t, tc.desc, nil
}

func dotGoAddInt(w io.Writer, n string, i int) {
	if i == 0 {
		// initialized to 0, ignore
		return
	}
	fmt.Fprintf(w, "\t\t%-13s %d,\n", n+":", i)
}
func dotGoAddStr(w io.Writer, n string, s string) {
	if s == "" {
		return
	}
	fmt.Fprintf(w, "\t\t%-13s %q,\n", n+":", s)
}
func dotGoAddFlag(w io.Writer, n string, b bool) {
	if !b {
		// initialized to 0, ignore
		return
	}
	fmt.Fprintf(w, "\t\t%-13s true,\n", n+":")
}

func dotGoAddArr(w io.Writer, n string, a []string) {
	if len(a) == 0 {
		return
	}
	fmt.Fprintf(w, "\t\t%-13s []string{", n+":")
	did := false
	for _, b := range a {
		if did {
			fmt.Fprint(w, ", ")
		}
		did = true
		fmt.Fprintf(w, "%q", b)
	}
	fmt.Fprintln(w, "},")
}

func dotGoHeader(w io.Writer, packname, tipackname string) {
	fmt.Fprintln(w, "// Generated automatically.  DO NOT HAND-EDIT.")
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "package %s\n", packname)
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "import \"%s\"\n", tipackname)
	fmt.Fprintln(w, "")
}

func dotGoTrailer(w io.Writer) {
}

func dotGoInfo(w io.Writer, terms []*TData) {

	fmt.Fprintln(w, "func init() {")
	for _, t := range terms {
		fmt.Fprintf(w, "\n\t// %s\n", t.Desc)
		fmt.Fprintln(w, "\tterminfo.AddTerminfo(&terminfo.Terminfo{")
		dotGoAddStr(w, "Name", t.Name)
		dotGoAddArr(w, "Aliases", t.Aliases)
		dotGoAddInt(w, "Columns", t.Columns)
		dotGoAddInt(w, "Lines", t.Lines)
		dotGoAddInt(w, "Colors", t.Colors)
		dotGoAddStr(w, "Bell", t.Bell)
		dotGoAddStr(w, "Clear", t.Clear)
		dotGoAddStr(w, "EnterCA", t.EnterCA)
		dotGoAddStr(w, "ExitCA", t.ExitCA)
		dotGoAddStr(w, "ShowCursor", t.ShowCursor)
		dotGoAddStr(w, "HideCursor", t.HideCursor)
		dotGoAddStr(w, "AttrOff", t.AttrOff)
		dotGoAddStr(w, "Underline", t.Underline)
		dotGoAddStr(w, "Bold", t.Bold)
		dotGoAddStr(w, "Dim", t.Dim)
		dotGoAddStr(w, "Italic", t.Italic)
		dotGoAddStr(w, "Blink", t.Blink)
		dotGoAddStr(w, "Reverse", t.Reverse)
		dotGoAddStr(w, "EnterKeypad", t.EnterKeypad)
		dotGoAddStr(w, "ExitKeypad", t.ExitKeypad)
		dotGoAddStr(w, "SetFg", t.SetFg)
		dotGoAddStr(w, "SetBg", t.SetBg)
		dotGoAddStr(w, "SetFgBg", t.SetFgBg)
		dotGoAddStr(w, "ResetFgBg", t.ResetFgBg)
		dotGoAddStr(w, "PadChar", t.PadChar)
		dotGoAddStr(w, "AltChars", t.AltChars)
		dotGoAddStr(w, "EnterAcs", t.EnterAcs)
		dotGoAddStr(w, "ExitAcs", t.ExitAcs)
		dotGoAddStr(w, "EnableAcs", t.EnableAcs)
		dotGoAddStr(w, "SetFgRGB", t.SetFgRGB)
		dotGoAddStr(w, "SetBgRGB", t.SetBgRGB)
		dotGoAddStr(w, "SetFgBgRGB", t.SetFgBgRGB)
		dotGoAddStr(w, "StrikeThrough", t.StrikeThrough)
		dotGoAddStr(w, "Mouse", t.Mouse)
		dotGoAddStr(w, "MouseMode", t.MouseMode)
		dotGoAddStr(w, "SetCursor", t.SetCursor)
		dotGoAddStr(w, "CursorBack1", t.CursorBack1)
		dotGoAddStr(w, "CursorUp1", t.CursorUp1)
		dotGoAddStr(w, "KeyUp", t.KeyUp)
		dotGoAddStr(w, "KeyDown", t.KeyDown)
		dotGoAddStr(w, "KeyRight", t.KeyRight)
		dotGoAddStr(w, "KeyLeft", t.KeyLeft)
		dotGoAddStr(w, "KeyInsert", t.KeyInsert)
		dotGoAddStr(w, "KeyDelete", t.KeyDelete)
		dotGoAddStr(w, "KeyBackspace", t.KeyBackspace)
		dotGoAddStr(w, "KeyHome", t.KeyHome)
		dotGoAddStr(w, "KeyEnd", t.KeyEnd)
		dotGoAddStr(w, "KeyPgUp", t.KeyPgUp)
		dotGoAddStr(w, "KeyPgDn", t.KeyPgDn)
		dotGoAddStr(w, "KeyF1", t.KeyF1)
		dotGoAddStr(w, "KeyF2", t.KeyF2)
		dotGoAddStr(w, "KeyF3", t.KeyF3)
		dotGoAddStr(w, "KeyF4", t.KeyF4)
		dotGoAddStr(w, "KeyF5", t.KeyF5)
		dotGoAddStr(w, "KeyF6", t.KeyF6)
		dotGoAddStr(w, "KeyF7", t.KeyF7)
		dotGoAddStr(w, "KeyF8", t.KeyF8)
		dotGoAddStr(w, "KeyF9", t.KeyF9)
		dotGoAddStr(w, "KeyF10", t.KeyF10)
		dotGoAddStr(w, "KeyF11", t.KeyF11)
		dotGoAddStr(w, "KeyF12", t.KeyF12)
		// Extended keys.  We don't report these if they are going to be
		// handled as if they were XTerm sequences.
		dotGoAddStr(w, "KeyF13", t.KeyF13)
		dotGoAddStr(w, "KeyF14", t.KeyF14)
		dotGoAddStr(w, "KeyF15", t.KeyF15)
		dotGoAddStr(w, "KeyF16", t.KeyF16)
		dotGoAddStr(w, "KeyF17", t.KeyF17)
		dotGoAddStr(w, "KeyF18", t.KeyF18)
		dotGoAddStr(w, "KeyF19", t.KeyF19)
		dotGoAddStr(w, "KeyF20", t.KeyF20)
		dotGoAddStr(w, "KeyF21", t.KeyF21)
		dotGoAddStr(w, "KeyF22", t.KeyF22)
		dotGoAddStr(w, "KeyF23", t.KeyF23)
		dotGoAddStr(w, "KeyF24", t.KeyF24)
		dotGoAddStr(w, "KeyF25", t.KeyF25)
		dotGoAddStr(w, "KeyF26", t.KeyF26)
		dotGoAddStr(w, "KeyF27", t.KeyF27)
		dotGoAddStr(w, "KeyF28", t.KeyF28)
		dotGoAddStr(w, "KeyF29", t.KeyF29)
		dotGoAddStr(w, "KeyF30", t.KeyF30)
		dotGoAddStr(w, "KeyF31", t.KeyF31)
		dotGoAddStr(w, "KeyF32", t.KeyF32)
		dotGoAddStr(w, "KeyF33", t.KeyF33)
		dotGoAddStr(w, "KeyF34", t.KeyF34)
		dotGoAddStr(w, "KeyF35", t.KeyF35)
		dotGoAddStr(w, "KeyF36", t.KeyF36)
		dotGoAddStr(w, "KeyF37", t.KeyF37)
		dotGoAddStr(w, "KeyF38", t.KeyF38)
		dotGoAddStr(w, "KeyF39", t.KeyF39)
		dotGoAddStr(w, "KeyF40", t.KeyF40)
		dotGoAddStr(w, "KeyF41", t.KeyF41)
		dotGoAddStr(w, "KeyF42", t.KeyF42)
		dotGoAddStr(w, "KeyF43", t.KeyF43)
		dotGoAddStr(w, "KeyF44", t.KeyF44)
		dotGoAddStr(w, "KeyF45", t.KeyF45)
		dotGoAddStr(w, "KeyF46", t.KeyF46)
		dotGoAddStr(w, "KeyF47", t.KeyF47)
		dotGoAddStr(w, "KeyF48", t.KeyF48)
		dotGoAddStr(w, "KeyF49", t.KeyF49)
		dotGoAddStr(w, "KeyF50", t.KeyF50)
		dotGoAddStr(w, "KeyF51", t.KeyF51)
		dotGoAddStr(w, "KeyF52", t.KeyF52)
		dotGoAddStr(w, "KeyF53", t.KeyF53)
		dotGoAddStr(w, "KeyF54", t.KeyF54)
		dotGoAddStr(w, "KeyF55", t.KeyF55)
		dotGoAddStr(w, "KeyF56", t.KeyF56)
		dotGoAddStr(w, "KeyF57", t.KeyF57)
		dotGoAddStr(w, "KeyF58", t.KeyF58)
		dotGoAddStr(w, "KeyF59", t.KeyF59)
		dotGoAddStr(w, "KeyF60", t.KeyF60)
		dotGoAddStr(w, "KeyF61", t.KeyF61)
		dotGoAddStr(w, "KeyF62", t.KeyF62)
		dotGoAddStr(w, "KeyF63", t.KeyF63)
		dotGoAddStr(w, "KeyF64", t.KeyF64)
		dotGoAddStr(w, "KeyCancel", t.KeyCancel)
		dotGoAddStr(w, "KeyPrint", t.KeyPrint)
		dotGoAddStr(w, "KeyExit", t.KeyExit)
		dotGoAddStr(w, "KeyHelp", t.KeyHelp)
		dotGoAddStr(w, "KeyClear", t.KeyClear)
		dotGoAddStr(w, "KeyBacktab", t.KeyBacktab)
		dotGoAddStr(w, "KeyShfLeft", t.KeyShfLeft)
		dotGoAddStr(w, "KeyShfRight", t.KeyShfRight)
		dotGoAddStr(w, "KeyShfUp", t.KeyShfUp)
		dotGoAddStr(w, "KeyShfDown", t.KeyShfDown)
		dotGoAddStr(w, "KeyShfHome", t.KeyShfHome)
		dotGoAddStr(w, "KeyShfEnd", t.KeyShfEnd)
		dotGoAddStr(w, "KeyShfInsert", t.KeyShfInsert)
		dotGoAddStr(w, "KeyShfDelete", t.KeyShfDelete)
		dotGoAddStr(w, "KeyCtrlUp", t.KeyCtrlUp)
		dotGoAddStr(w, "KeyCtrlDown", t.KeyCtrlDown)
		dotGoAddStr(w, "KeyCtrlRight", t.KeyCtrlRight)
		dotGoAddStr(w, "KeyCtrlLeft", t.KeyCtrlLeft)
		dotGoAddStr(w, "KeyCtrlHome", t.KeyCtrlHome)
		dotGoAddStr(w, "KeyCtrlEnd", t.KeyCtrlEnd)
		dotGoAddInt(w, "Modifiers", t.Modifiers)
		dotGoAddFlag(w, "TrueColor", t.TrueColor)
		fmt.Fprintln(w, "\t})")
	}
	fmt.Fprintln(w, "}")
}

var packname = ""
var tipackname = "github.com/gdamore/tcell/v2/terminfo"

func dotGoFile(fname string, terms []*TData) error {
	w := os.Stdout
	var e error
	if fname != "-" && fname != "" {
		if w, e = os.Create(fname); e != nil {
			return e
		}
	}
	if packname == "" {
		packname = strings.Replace(terms[0].Name, "-", "_", -1)
	}
	dotGoHeader(w, packname, tipackname)
	dotGoInfo(w, terms)
	dotGoTrailer(w)
	if w != os.Stdout {
		w.Close()
	}
	cmd := exec.Command("go", "fmt", fname)
	cmd.Run()
	return nil
}

type TData struct {
	Desc string

	terminfo.Terminfo
}

func main() {
	gofile := ""
	nofatal := false
	quiet := false
	all := false

	flag.StringVar(&gofile, "go", "", "generate go source in named file")
	flag.StringVar(&tipackname, "I", tipackname, "import package path")
	flag.StringVar(&packname, "P", packname, "package name (go source)")
	flag.BoolVar(&nofatal, "nofatal", false, "errors are not fatal")
	flag.BoolVar(&quiet, "quiet", false, "suppress error messages")
	flag.BoolVar(&all, "all", false, "load all terminals from terminfo")
	flag.Parse()
	var e error

	args := flag.Args()
	if len(args) == 0 {
		args = []string{os.Getenv("TERM")}
	}

	tdata := make([]*TData, 0)

	for _, term := range args {
		if t, desc, e := getinfo(term); e != nil {
			if all && e == notaddressable {
				continue
			}
			if !quiet {
				fmt.Fprintf(os.Stderr,
					"Failed loading %s: %v\n", term, e)
			}
			if !nofatal {
				os.Exit(1)
			}
		} else {
			tdata = append(tdata, &TData{
				Desc:     desc,
				Terminfo: *t,
			})
		}
	}

	if len(tdata) == 0 {
		// No data.
		os.Exit(0)
	}

	e = dotGoFile(gofile, tdata)
	if e != nil {
		fmt.Fprintf(os.Stderr, "Failed %s: %v", gofile, e)
		os.Exit(1)
	}
}
