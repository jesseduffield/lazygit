// Generated automatically.  DO NOT HAND-EDIT.

package vt220

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// DEC VT220
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:              "vt220",
		Aliases:           []string{"vt200"},
		Columns:           80,
		Lines:             24,
		Clear:             "\x1b[H\x1b[J",
		ShowCursor:        "\x1b[?25h",
		HideCursor:        "\x1b[?25l",
		AttrOff:           "\x1b[m\x1b(B",
		Underline:         "\x1b[4m",
		Bold:              "\x1b[1m",
		Blink:             "\x1b[5m",
		Reverse:           "\x1b[7m",
		PadChar:           "\x00",
		AltChars:          "``aaffggjjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:          "\x1b(0$<2>",
		ExitAcs:           "\x1b(B$<4>",
		EnableAcs:         "\x1b)0",
		EnableAutoMargin:  "\x1b[?7h",
		DisableAutoMargin: "\x1b[?7l",
		SetCursor:         "\x1b[%i%p1%d;%p2%dH",
		AutoMargin:        true,
	})
}
