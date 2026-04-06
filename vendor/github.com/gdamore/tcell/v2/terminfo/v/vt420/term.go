// This file was originally generated automatically,
// but it is edited to correct for errors in the VT420
// terminfo data.  Additionally we have added extended
// information for the extended F-keys.

package vt420

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// DEC VT420
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:              "vt420",
		Columns:           80,
		Lines:             24,
		Clear:             "\x1b[H\x1b[2J$<50>",
		ShowCursor:        "\x1b[?25h",
		HideCursor:        "\x1b[?25l",
		AttrOff:           "\x1b[m\x1b(B$<2>",
		Underline:         "\x1b[4m",
		Bold:              "\x1b[1m$<2>",
		Blink:             "\x1b[5m$<2>",
		Reverse:           "\x1b[7m$<2>",
		EnterKeypad:       "\x1b=",
		ExitKeypad:        "\x1b>",
		PadChar:           "\x00",
		AltChars:          "``aaffggjjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:          "\x1b(0$<2>",
		ExitAcs:           "\x1b(B$<4>",
		EnableAcs:         "\x1b)0",
		EnableAutoMargin:  "\x1b[?7h",
		DisableAutoMargin: "\x1b[?7l",
		SetCursor:         "\x1b[%i%p1%d;%p2%dH$<10>",
		AutoMargin:        true,
	})
}
