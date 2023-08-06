// Generated automatically.  DO NOT HAND-EDIT.

package vt320

import "github.com/stefanhaller/tcell/v2/terminfo"

func init() {

	// dec vt320 7 bit terminal
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:         "vt320",
		Aliases:      []string{"vt300"},
		Columns:      80,
		Lines:        24,
		Bell:         "\a",
		Clear:        "\x1b[H\x1b[2J",
		ShowCursor:   "\x1b[?25h",
		HideCursor:   "\x1b[?25l",
		AttrOff:      "\x1b[m\x1b(B",
		Underline:    "\x1b[4m",
		Bold:         "\x1b[1m",
		Blink:        "\x1b[5m",
		Reverse:      "\x1b[7m",
		EnterKeypad:  "\x1b[?1h\x1b=",
		ExitKeypad:   "\x1b[?1l\x1b>",
		PadChar:      "\x00",
		AltChars:     "``aaffggjjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:     "\x1b(0",
		ExitAcs:      "\x1b(B",
		SetCursor:    "\x1b[%i%p1%d;%p2%dH",
		CursorBack1:  "\b",
		CursorUp1:    "\x1b[A",
		KeyUp:        "\x1bOA",
		KeyDown:      "\x1bOB",
		KeyRight:     "\x1bOC",
		KeyLeft:      "\x1bOD",
		KeyInsert:    "\x1b[2~",
		KeyDelete:    "\x1b[3~",
		KeyBackspace: "\u007f",
		KeyHome:      "\x1b[1~",
		KeyPgUp:      "\x1b[5~",
		KeyPgDn:      "\x1b[6~",
		KeyF1:        "\x1bOP",
		KeyF2:        "\x1bOQ",
		KeyF3:        "\x1bOR",
		KeyF4:        "\x1bOS",
		KeyF6:        "\x1b[17~",
		KeyF7:        "\x1b[18~",
		KeyF8:        "\x1b[19~",
		KeyF9:        "\x1b[20~",
		KeyF10:       "\x1b[21~",
		KeyF11:       "\x1b[23~",
		KeyF12:       "\x1b[24~",
		KeyF13:       "\x1b[25~",
		KeyF14:       "\x1b[26~",
		KeyF15:       "\x1b[28~",
		KeyF16:       "\x1b[29~",
		KeyF17:       "\x1b[31~",
		KeyF18:       "\x1b[32~",
		KeyF19:       "\x1b[33~",
		KeyF20:       "\x1b[34~",
		AutoMargin:   true,
	})
}
