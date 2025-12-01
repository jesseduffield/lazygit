// Generated automatically.  DO NOT HAND-EDIT.

package cygwin

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// ANSI emulation for Cygwin
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:       "cygwin",
		Colors:     8,
		Clear:      "\x1b[H\x1b[J",
		EnterCA:    "\x1b7\x1b[?47h",
		ExitCA:     "\x1b[2J\x1b[?47l\x1b8",
		AttrOff:    "\x1b[0;10m",
		Underline:  "\x1b[4m",
		Bold:       "\x1b[1m",
		Reverse:    "\x1b[7m",
		SetFg:      "\x1b[3%p1%dm",
		SetBg:      "\x1b[4%p1%dm",
		SetFgBg:    "\x1b[3%p1%d;4%p2%dm",
		ResetFgBg:  "\x1b[39;49m",
		PadChar:    "\x00",
		AltChars:   "+\x10,\x11-\x18.\x190\xdb`\x04a\xb1f\xf8g\xf1h\xb0j\xd9k\xbfl\xdam\xc0n\xc5o~p\xc4q\xc4r\xc4s_t\xc3u\xb4v\xc1w\xc2x\xb3y\xf3z\xf2{\xe3|\xd8}\x9c~\xfe",
		EnterAcs:   "\x1b[11m",
		ExitAcs:    "\x1b[10m",
		SetCursor:  "\x1b[%i%p1%d;%p2%dH",
		AutoMargin: true,
		InsertChar: "\x1b[@",
	})
}
