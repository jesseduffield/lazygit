// Generated automatically.  DO NOT HAND-EDIT.

package vt320

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// DEC VT320 7 bit terminal
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:              "vt320",
		Aliases:           []string{"vt300"},
		Columns:           80,
		Lines:             24,
		Clear:             "\x1b[H\x1b[2J",
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
		AutoMargin:        true,
	})
}
