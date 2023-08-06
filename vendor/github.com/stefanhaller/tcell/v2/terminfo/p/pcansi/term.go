// Generated automatically.  DO NOT HAND-EDIT.

package pcansi

import "github.com/stefanhaller/tcell/v2/terminfo"

func init() {

	// ibm-pc terminal programs claiming to be ansi
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:         "pcansi",
		Columns:      80,
		Lines:        24,
		Colors:       8,
		Bell:         "\a",
		Clear:        "\x1b[H\x1b[J",
		AttrOff:      "\x1b[0;10m",
		Underline:    "\x1b[4m",
		Bold:         "\x1b[1m",
		Blink:        "\x1b[5m",
		Reverse:      "\x1b[7m",
		SetFg:        "\x1b[3%p1%dm",
		SetBg:        "\x1b[4%p1%dm",
		SetFgBg:      "\x1b[3%p1%d;4%p2%dm",
		ResetFgBg:    "\x1b[37;40m",
		PadChar:      "\x00",
		AltChars:     "+\x10,\x11-\x18.\x190\xdb`\x04a\xb1f\xf8g\xf1h\xb0j\xd9k\xbfl\xdam\xc0n\xc5o~p\xc4q\xc4r\xc4s_t\xc3u\xb4v\xc1w\xc2x\xb3y\xf3z\xf2{\xe3|\xd8}\x9c~\xfe",
		EnterAcs:     "\x1b[12m",
		ExitAcs:      "\x1b[10m",
		SetCursor:    "\x1b[%i%p1%d;%p2%dH",
		CursorBack1:  "\x1b[D",
		CursorUp1:    "\x1b[A",
		KeyUp:        "\x1b[A",
		KeyDown:      "\x1b[B",
		KeyRight:     "\x1b[C",
		KeyLeft:      "\x1b[D",
		KeyBackspace: "\b",
		KeyHome:      "\x1b[H",
		AutoMargin:   true,
	})
}
