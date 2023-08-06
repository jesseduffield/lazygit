// Generated automatically.  DO NOT HAND-EDIT.

package vt420

import "github.com/stefanhaller/tcell/v2/terminfo"

func init() {

	// DEC VT420
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:         "vt420",
		Columns:      80,
		Lines:        24,
		Bell:         "\a",
		Clear:        "\x1b[H\x1b[2J$<50>",
		ShowCursor:   "\x1b[?25h",
		HideCursor:   "\x1b[?25l",
		AttrOff:      "\x1b[m\x1b(B$<2>",
		Underline:    "\x1b[4m",
		Bold:         "\x1b[1m$<2>",
		Blink:        "\x1b[5m$<2>",
		Reverse:      "\x1b[7m$<2>",
		EnterKeypad:  "\x1b=",
		ExitKeypad:   "\x1b>",
		PadChar:      "\x00",
		AltChars:     "``aaffggjjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:     "\x1b(0$<2>",
		ExitAcs:      "\x1b(B$<4>",
		EnableAcs:    "\x1b)0",
		SetCursor:    "\x1b[%i%p1%d;%p2%dH$<10>",
		CursorBack1:  "\b",
		CursorUp1:    "\x1b[A",
		KeyUp:        "\x1b[A",
		KeyDown:      "\x1b[B",
		KeyRight:     "\x1b[C",
		KeyLeft:      "\x1b[D",
		KeyInsert:    "\x1b[2~",
		KeyDelete:    "\x1b[3~",
		KeyBackspace: "\b",
		KeyPgUp:      "\x1b[5~",
		KeyPgDn:      "\x1b[6~",
		KeyF1:        "\x1bOP",
		KeyF2:        "\x1bOQ",
		KeyF3:        "\x1bOR",
		KeyF4:        "\x1bOS",
		KeyF5:        "\x1b[17~",
		KeyF6:        "\x1b[18~",
		KeyF7:        "\x1b[19~",
		KeyF8:        "\x1b[20~",
		KeyF9:        "\x1b[21~",
		KeyF10:       "\x1b[29~",
		AutoMargin:   true,
	})
}
