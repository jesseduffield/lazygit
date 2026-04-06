// Generated automatically.  DO NOT HAND-EDIT.

package dtterm

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// CDE desktop terminal
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:              "dtterm",
		Columns:           80,
		Lines:             24,
		Colors:            8,
		Clear:             "\x1b[H\x1b[J",
		ShowCursor:        "\x1b[?25h",
		HideCursor:        "\x1b[?25l",
		AttrOff:           "\x1b[m\x0f",
		Underline:         "\x1b[4m",
		Bold:              "\x1b[1m",
		Dim:               "\x1b[2m",
		Blink:             "\x1b[5m",
		Reverse:           "\x1b[7m",
		SetFg:             "\x1b[3%p1%dm",
		SetBg:             "\x1b[4%p1%dm",
		SetFgBg:           "\x1b[3%p1%d;4%p2%dm",
		ResetFgBg:         "\x1b[39;49m",
		PadChar:           "\x00",
		AltChars:          "``aaffggjjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:          "\x0e",
		ExitAcs:           "\x0f",
		EnableAcs:         "\x1b(B\x1b)0",
		EnableAutoMargin:  "\x1b[?7h",
		DisableAutoMargin: "\x1b[?7l",
		SetCursor:         "\x1b[%i%p1%d;%p2%dH",
		AutoMargin:        true,
	})
}
