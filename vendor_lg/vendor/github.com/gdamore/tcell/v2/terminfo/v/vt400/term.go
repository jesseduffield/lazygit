// Generated automatically.  DO NOT HAND-EDIT.

package vt400

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// DEC VT400 24x80 column autowrap
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:              "vt400",
		Aliases:           []string{"vt400-24", "dec-vt400"},
		Columns:           80,
		Lines:             24,
		Clear:             "\x1b[H\x1b[J$<10/>",
		ShowCursor:        "\x1b[?25h",
		HideCursor:        "\x1b[?25l",
		AttrOff:           "\x1b[m\x1b(B",
		Underline:         "\x1b[4m",
		Bold:              "\x1b[1m",
		Blink:             "\x1b[5m",
		Reverse:           "\x1b[7m",
		EnterKeypad:       "\x1b[?1h\x1b=",
		ExitKeypad:        "\x1b[?1l\x1b>",
		PadChar:           "\x00",
		AltChars:          "``aaffggjjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:          "\x1b(0",
		ExitAcs:           "\x1b(B",
		EnableAutoMargin:  "\x1b[?7h",
		DisableAutoMargin: "\x1b[?7l",
		SetCursor:         "\x1b[%i%p1%d;%p2%dH",
		CursorBack1:       "\b",
		CursorUp1:         "\x1b[A",
		KeyUp:             "\x1bOA",
		KeyDown:           "\x1bOB",
		KeyRight:          "\x1bOC",
		KeyLeft:           "\x1bOD",
		KeyBackspace:      "\b",
		KeyF1:             "\x1bOP",
		KeyF2:             "\x1bOQ",
		KeyF3:             "\x1bOR",
		KeyF4:             "\x1bOS",
		KeyF6:             "\x1b[17~",
		KeyF7:             "\x1b[18~",
		KeyF8:             "\x1b[19~",
		KeyF9:             "\x1b[20~",
		AutoMargin:        true,
		InsertChar:        "\x1b[@",
	})
}
