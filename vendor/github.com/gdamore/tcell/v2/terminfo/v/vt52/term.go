// Generated automatically.  DO NOT HAND-EDIT.

package vt52

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// DEC VT52
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:         "vt52",
		Columns:      80,
		Lines:        24,
		Bell:         "\a",
		Clear:        "\x1bH\x1bJ",
		EnterKeypad:  "\x1b=",
		ExitKeypad:   "\x1b>",
		PadChar:      "\x00",
		AltChars:     "+h.k0affggolpnqprrss",
		EnterAcs:     "\x1bF",
		ExitAcs:      "\x1bG",
		SetCursor:    "\x1bY%p1%' '%+%c%p2%' '%+%c",
		CursorBack1:  "\x1bD",
		CursorUp1:    "\x1bA",
		KeyUp:        "\x1bA",
		KeyDown:      "\x1bB",
		KeyRight:     "\x1bC",
		KeyLeft:      "\x1bD",
		KeyBackspace: "\b",
		KeyF1:        "\x1bP",
		KeyF2:        "\x1bQ",
		KeyF3:        "\x1bR",
		KeyF5:        "\x1b?t",
		KeyF6:        "\x1b?u",
		KeyF7:        "\x1b?v",
		KeyF8:        "\x1b?w",
		KeyF9:        "\x1b?x",
	})
}
