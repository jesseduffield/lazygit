// Generated automatically.  DO NOT HAND-EDIT.

package wy50

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// Wyse 50
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:         "wy50",
		Aliases:      []string{"wyse50"},
		Columns:      80,
		Lines:        24,
		Bell:         "\a",
		Clear:        "\x1b+$<20>",
		ShowCursor:   "\x1b`1",
		HideCursor:   "\x1b`0",
		AttrOff:      "\x1b(\x1bH\x03",
		Dim:          "\x1b`7\x1b)",
		Reverse:      "\x1b`6\x1b)",
		PadChar:      "\x00",
		AltChars:     "a;j5k3l2m1n8q:t4u9v=w0x6",
		EnterAcs:     "\x1bH\x02",
		ExitAcs:      "\x1bH\x03",
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
