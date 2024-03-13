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
		Bell:              "\a",
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
		CursorBack1:       "\b",
		CursorUp1:         "\x1b[A",
		KeyUp:             "\x1b[A",
		KeyDown:           "\x1b[B",
		KeyRight:          "\x1b[C",
		KeyLeft:           "\x1b[D",
		KeyInsert:         "\x1b[2~",
		KeyDelete:         "\x1b[3~",
		KeyBackspace:      "\b",
		KeyPgUp:           "\x1b[5~",
		KeyPgDn:           "\x1b[6~",
		KeyF1:             "\x1bOP",
		KeyF2:             "\x1bOQ",
		KeyF3:             "\x1bOR",
		KeyF4:             "\x1bOS",
		KeyF6:             "\x1b[17~",
		KeyF7:             "\x1b[18~",
		KeyF8:             "\x1b[19~",
		KeyF9:             "\x1b[20~",
		KeyF10:            "\x1b[21~",
		KeyF11:            "\x1b[23~",
		KeyF12:            "\x1b[24~",
		KeyF13:            "\x1b[25~",
		KeyF14:            "\x1b[26~",
		KeyF17:            "\x1b[31~",
		KeyF18:            "\x1b[32~",
		KeyF19:            "\x1b[33~",
		KeyF20:            "\x1b[34~",
		KeyHelp:           "\x1b[28~",
		AutoMargin:        true,
	})
}
