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

// This terminal definition is hand-coded, as the default terminfo for
// this terminal is busted with respect to color.  Unlike pretty much every
// other ANSI compliant terminal, this terminal cannot combine foreground and
// background escapes.  The default terminfo also only provides escapes for
// 16-bit color.

package sun

import "github.com/stefanhaller/tcell/v2/terminfo"

func init() {

	// Sun Microsystems Inc. workstation console
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:         "sun",
		Aliases:      []string{"sun1", "sun2"},
		Columns:      80,
		Lines:        34,
		Bell:         "\a",
		Clear:        "\f",
		AttrOff:      "\x1b[m",
		Reverse:      "\x1b[7m",
		PadChar:      "\x00",
		SetCursor:    "\x1b[%i%p1%d;%p2%dH",
		CursorBack1:  "\b",
		CursorUp1:    "\x1b[A",
		KeyUp:        "\x1b[A",
		KeyDown:      "\x1b[B",
		KeyRight:     "\x1b[C",
		KeyLeft:      "\x1b[D",
		KeyInsert:    "\x1b[247z",
		KeyDelete:    "\u007f",
		KeyBackspace: "\b",
		KeyHome:      "\x1b[214z",
		KeyEnd:       "\x1b[220z",
		KeyPgUp:      "\x1b[216z",
		KeyPgDn:      "\x1b[222z",
		KeyF1:        "\x1b[224z",
		KeyF2:        "\x1b[225z",
		KeyF3:        "\x1b[226z",
		KeyF4:        "\x1b[227z",
		KeyF5:        "\x1b[228z",
		KeyF6:        "\x1b[229z",
		KeyF7:        "\x1b[230z",
		KeyF8:        "\x1b[231z",
		KeyF9:        "\x1b[232z",
		KeyF10:       "\x1b[233z",
		KeyF11:       "\x1b[234z",
		KeyF12:       "\x1b[235z",
		AutoMargin:   true,
		InsertChar:   "\x1b[@",
	})

	// Sun Microsystems Workstation console with color support (IA systems)
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:         "sun-color",
		Columns:      80,
		Lines:        34,
		Colors:       256,
		Bell:         "\a",
		Clear:        "\f",
		AttrOff:      "\x1b[m",
		Bold:         "\x1b[1m",
		Reverse:      "\x1b[7m",
		SetFg:        "\x1b[38;5;%p1%dm",
		SetBg:        "\x1b[48;5;%p1%dm",
		ResetFgBg:    "\x1b[0m",
		PadChar:      "\x00",
		SetCursor:    "\x1b[%i%p1%d;%p2%dH",
		CursorBack1:  "\b",
		CursorUp1:    "\x1b[A",
		KeyUp:        "\x1b[A",
		KeyDown:      "\x1b[B",
		KeyRight:     "\x1b[C",
		KeyLeft:      "\x1b[D",
		KeyInsert:    "\x1b[247z",
		KeyDelete:    "\u007f",
		KeyBackspace: "\b",
		KeyHome:      "\x1b[214z",
		KeyEnd:       "\x1b[220z",
		KeyPgUp:      "\x1b[216z",
		KeyPgDn:      "\x1b[222z",
		KeyF1:        "\x1b[224z",
		KeyF2:        "\x1b[225z",
		KeyF3:        "\x1b[226z",
		KeyF4:        "\x1b[227z",
		KeyF5:        "\x1b[228z",
		KeyF6:        "\x1b[229z",
		KeyF7:        "\x1b[230z",
		KeyF8:        "\x1b[231z",
		KeyF9:        "\x1b[232z",
		KeyF10:       "\x1b[233z",
		KeyF11:       "\x1b[234z",
		KeyF12:       "\x1b[235z",
		AutoMargin:   true,
		InsertChar:   "\x1b[@",
	})
}
