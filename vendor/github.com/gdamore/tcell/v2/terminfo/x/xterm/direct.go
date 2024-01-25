// Copyright 2021 The TCell Authors
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

// This terminal definition is derived from the xterm-256color definition, but
// makes use of the RGB property these terminals have to support direct color.
// The terminfo entry for this uses a new format for the color handling introduced
// by ncurses 6.1 (and used by nobody else), so this override ensures we get
// good handling even in the face of this.

package xterm

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// derived from xterm-256color, but adds full RGB support
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:          "xterm-direct",
		Aliases:       []string{"xterm-truecolor"},
		Columns:       80,
		Lines:         24,
		Colors:        256,
		Bell:          "\a",
		Clear:         "\x1b[H\x1b[2J",
		EnterCA:       "\x1b[?1049h\x1b[22;0;0t",
		ExitCA:        "\x1b[?1049l\x1b[23;0;0t",
		ShowCursor:    "\x1b[?12l\x1b[?25h",
		HideCursor:    "\x1b[?25l",
		AttrOff:       "\x1b(B\x1b[m",
		Underline:     "\x1b[4m",
		Bold:          "\x1b[1m",
		Dim:           "\x1b[2m",
		Italic:        "\x1b[3m",
		Blink:         "\x1b[5m",
		Reverse:       "\x1b[7m",
		EnterKeypad:   "\x1b[?1h\x1b=",
		ExitKeypad:    "\x1b[?1l\x1b>",
		SetFg:         "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;m",
		SetBg:         "\x1b[%?%p1%{8}%<%t4%p1%d%e%p1%{16}%<%t10%p1%{8}%-%d%e48;5;%p1%d%;m",
		SetFgBg:       "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;;%?%p2%{8}%<%t4%p2%d%e%p2%{16}%<%t10%p2%{8}%-%d%e48;5;%p2%d%;m",
		SetFgRGB:      "\x1b[38;2;%p1%d;%p2%d;%p3%dm",
		SetBgRGB:      "\x1b[48;2;%p1%d;%p2%d;%p3%dm",
		SetFgBgRGB:    "\x1b[38;2;%p1%d;%p2%d;%p3%d;48;2;%p4%d;%p5%d;%p6%dm",
		ResetFgBg:     "\x1b[39;49m",
		AltChars:      "``aaffggiijjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:      "\x1b(0",
		ExitAcs:       "\x1b(B",
		StrikeThrough: "\x1b[9m",
		Mouse:         "\x1b[M",
		SetCursor:     "\x1b[%i%p1%d;%p2%dH",
		CursorBack1:   "\b",
		CursorUp1:     "\x1b[A",
		KeyUp:         "\x1bOA",
		KeyDown:       "\x1bOB",
		KeyRight:      "\x1bOC",
		KeyLeft:       "\x1bOD",
		KeyInsert:     "\x1b[2~",
		KeyDelete:     "\x1b[3~",
		KeyBackspace:  "\u007f",
		KeyHome:       "\x1bOH",
		KeyEnd:        "\x1bOF",
		KeyPgUp:       "\x1b[5~",
		KeyPgDn:       "\x1b[6~",
		KeyF1:         "\x1bOP",
		KeyF2:         "\x1bOQ",
		KeyF3:         "\x1bOR",
		KeyF4:         "\x1bOS",
		KeyF5:         "\x1b[15~",
		KeyF6:         "\x1b[17~",
		KeyF7:         "\x1b[18~",
		KeyF8:         "\x1b[19~",
		KeyF9:         "\x1b[20~",
		KeyF10:        "\x1b[21~",
		KeyF11:        "\x1b[23~",
		KeyF12:        "\x1b[24~",
		KeyBacktab:    "\x1b[Z",
		Modifiers:     1,
		AutoMargin:    true,
		TrueColor:     true,
	})
}
