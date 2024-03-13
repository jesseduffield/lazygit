// Generated automatically.  DO NOT HAND-EDIT.

package vt100

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// DEC VT100 (w/advanced video)
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:              "vt100",
		Aliases:           []string{"vt100-am"},
		Columns:           80,
		Lines:             24,
		Bell:              "\a",
		Clear:             "\x1b[H\x1b[J$<50>",
		AttrOff:           "\x1b[m\x0f$<2>",
		Underline:         "\x1b[4m$<2>",
		Bold:              "\x1b[1m$<2>",
		Blink:             "\x1b[5m$<2>",
		Reverse:           "\x1b[7m$<2>",
		EnterKeypad:       "\x1b[?1h\x1b=",
		ExitKeypad:        "\x1b[?1l\x1b>",
		PadChar:           "\x00",
		AltChars:          "``aaffggjjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:          "\x0e",
		ExitAcs:           "\x0f",
		EnableAcs:         "\x1b(B\x1b)0",
		EnableAutoMargin:  "\x1b[?7h",
		DisableAutoMargin: "\x1b[?7l",
		SetCursor:         "\x1b[%i%p1%d;%p2%dH$<5>",
		CursorBack1:       "\b",
		CursorUp1:         "\x1b[A$<2>",
		KeyUp:             "\x1bOA",
		KeyDown:           "\x1bOB",
		KeyRight:          "\x1bOC",
		KeyLeft:           "\x1bOD",
		KeyBackspace:      "\b",
		KeyF1:             "\x1bOP",
		KeyF2:             "\x1bOQ",
		KeyF3:             "\x1bOR",
		KeyF4:             "\x1bOS",
		KeyF5:             "\x1bOt",
		KeyF6:             "\x1bOu",
		KeyF7:             "\x1bOv",
		KeyF8:             "\x1bOl",
		KeyF9:             "\x1bOw",
		KeyF10:            "\x1bOx",
		AutoMargin:        true,
	})
}
