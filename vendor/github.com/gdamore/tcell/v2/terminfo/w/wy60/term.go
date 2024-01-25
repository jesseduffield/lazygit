// Generated automatically.  DO NOT HAND-EDIT.

package wy60

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// Wyse 60
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:         "wy60",
		Aliases:      []string{"wyse60"},
		Columns:      80,
		Lines:        24,
		Bell:         "\a",
		Clear:        "\x1b+$<100>",
		EnterCA:      "\x1bw0",
		ExitCA:       "\x1bw1",
		ShowCursor:   "\x1b`1",
		HideCursor:   "\x1b`0",
		AttrOff:      "\x1b(\x1bH\x03\x1bG0\x1bcD",
		Underline:    "\x1bG8",
		Dim:          "\x1bGp",
		Blink:        "\x1bG2",
		Reverse:      "\x1bG4",
		PadChar:      "\x00",
		AltChars:     "+/,.0[a2fxgqh1ihjYk?lZm@nEqDtCu4vAwBx3yszr{c~~",
		EnterAcs:     "\x1bcE",
		ExitAcs:      "\x1bcD",
		SetCursor:    "\x1b=%p1%' '%+%c%p2%' '%+%c",
		CursorBack1:  "\b",
		CursorUp1:    "\v",
		KeyUp:        "\v",
		KeyDown:      "\n",
		KeyRight:     "\f",
		KeyLeft:      "\b",
		KeyInsert:    "\x1bQ",
		KeyDelete:    "\x1bW",
		KeyBackspace: "\b",
		KeyHome:      "\x1e",
		KeyPgUp:      "\x1bJ",
		KeyPgDn:      "\x1bK",
		KeyF1:        "\x01@\r",
		KeyF2:        "\x01A\r",
		KeyF3:        "\x01B\r",
		KeyF4:        "\x01C\r",
		KeyF5:        "\x01D\r",
		KeyF6:        "\x01E\r",
		KeyF7:        "\x01F\r",
		KeyF8:        "\x01G\r",
		KeyF9:        "\x01H\r",
		KeyF10:       "\x01I\r",
		KeyF11:       "\x01J\r",
		KeyF12:       "\x01K\r",
		KeyF13:       "\x01L\r",
		KeyF14:       "\x01M\r",
		KeyF15:       "\x01N\r",
		KeyF16:       "\x01O\r",
		KeyPrint:     "\x1bP",
		KeyBacktab:   "\x1bI",
		KeyShfHome:   "\x1b{",
		AutoMargin:   true,
	})
}
