// Generated automatically.  DO NOT HAND-EDIT.

package emacs

import "github.com/stefanhaller/tcell/v2/terminfo"

func init() {

	// gnu emacs term.el terminal emulation
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:        "eterm",
		Columns:     80,
		Lines:       24,
		Bell:        "\a",
		Clear:       "\x1b[H\x1b[J",
		EnterCA:     "\x1b7\x1b[?47h",
		ExitCA:      "\x1b[2J\x1b[?47l\x1b8",
		AttrOff:     "\x1b[m",
		Underline:   "\x1b[4m",
		Bold:        "\x1b[1m",
		Reverse:     "\x1b[7m",
		PadChar:     "\x00",
		SetCursor:   "\x1b[%i%p1%d;%p2%dH",
		CursorBack1: "\b",
		CursorUp1:   "\x1b[A",
		AutoMargin:  true,
	})

	// Emacs term.el terminal emulator term-protocol-version 0.96
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:         "eterm-color",
		Columns:      80,
		Lines:        24,
		Colors:       8,
		Bell:         "\a",
		Clear:        "\x1b[H\x1b[J",
		AttrOff:      "\x1b[m",
		Underline:    "\x1b[4m",
		Bold:         "\x1b[1m",
		Blink:        "\x1b[5m",
		Reverse:      "\x1b[7m",
		SetFg:        "\x1b[%p1%{30}%+%dm",
		SetBg:        "\x1b[%p1%'('%+%dm",
		SetFgBg:      "\x1b[%p1%{30}%+%d;%p2%'('%+%dm",
		ResetFgBg:    "\x1b[39;49m",
		PadChar:      "\x00",
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
		KeyEnd:       "\x1b[4~",
		KeyPgUp:      "\x1b[5~",
		KeyPgDn:      "\x1b[6~",
		AutoMargin:   true,
	})
}
