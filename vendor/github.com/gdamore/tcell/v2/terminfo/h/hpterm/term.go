// Generated automatically.  DO NOT HAND-EDIT.

package hpterm

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// HP X11 terminal emulator (old)
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:         "hpterm",
		Aliases:      []string{"X-hpterm"},
		Columns:      80,
		Lines:        24,
		Bell:         "\a",
		Clear:        "\x1b&a0y0C\x1bJ",
		AttrOff:      "\x1b&d@\x0f",
		Underline:    "\x1b&dD",
		Bold:         "\x1b&dB",
		Dim:          "\x1b&dH",
		Reverse:      "\x1b&dB",
		EnterKeypad:  "\x1b&s1A",
		ExitKeypad:   "\x1b&s0A",
		PadChar:      "\x00",
		EnterAcs:     "\x0e",
		ExitAcs:      "\x0f",
		SetCursor:    "\x1b&a%p1%dy%p2%dC",
		CursorBack1:  "\b",
		CursorUp1:    "\x1bA",
		KeyUp:        "\x1bA",
		KeyDown:      "\x1bB",
		KeyRight:     "\x1bC",
		KeyLeft:      "\x1bD",
		KeyInsert:    "\x1bQ",
		KeyDelete:    "\x1bP",
		KeyBackspace: "\b",
		KeyHome:      "\x1bh",
		KeyPgUp:      "\x1bV",
		KeyPgDn:      "\x1bU",
		KeyF1:        "\x1bp",
		KeyF2:        "\x1bq",
		KeyF3:        "\x1br",
		KeyF4:        "\x1bs",
		KeyF5:        "\x1bt",
		KeyF6:        "\x1bu",
		KeyF7:        "\x1bv",
		KeyF8:        "\x1bw",
		KeyClear:     "\x1bJ",
		AutoMargin:   true,
	})
}
