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
		AutoMargin:        true,
		InsertChar:        "\x1b[@",
	})
}
