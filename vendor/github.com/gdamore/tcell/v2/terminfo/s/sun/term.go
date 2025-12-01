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

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// Sun Microsystems Inc. workstation console
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:       "sun",
		Aliases:    []string{"sun1", "sun2"},
		Columns:    80,
		Lines:      34,
		Clear:      "\f",
		AttrOff:    "\x1b[m",
		Reverse:    "\x1b[7m",
		PadChar:    "\x00",
		SetCursor:  "\x1b[%i%p1%d;%p2%dH",
		AutoMargin: true,
		InsertChar: "\x1b[@",
	})

	// Sun Microsystems Workstation console with color support (IA systems)
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:       "sun-color",
		Columns:    80,
		Lines:      34,
		Colors:     256,
		Clear:      "\f",
		AttrOff:    "\x1b[m",
		Bold:       "\x1b[1m",
		Reverse:    "\x1b[7m",
		SetFg:      "\x1b[38;5;%p1%dm",
		SetBg:      "\x1b[48;5;%p1%dm",
		ResetFgBg:  "\x1b[0m",
		PadChar:    "\x00",
		SetCursor:  "\x1b[%i%p1%d;%p2%dH",
		AutoMargin: true,
		InsertChar: "\x1b[@",
	})
}
