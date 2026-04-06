// Generated automatically.  DO NOT HAND-EDIT.

package aixterm

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// IBM Aixterm Terminal Emulator
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:       "aixterm",
		Columns:    80,
		Lines:      25,
		Colors:     8,
		Clear:      "\x1b[H\x1b[J",
		AttrOff:    "\x1b[0;10m\x1b(B",
		Underline:  "\x1b[4m",
		Bold:       "\x1b[1m",
		Reverse:    "\x1b[7m",
		SetFg:      "\x1b[3%p1%dm",
		SetBg:      "\x1b[4%p1%dm",
		SetFgBg:    "\x1b[3%p1%d;4%p2%dm",
		ResetFgBg:  "\x1b[32m\x1b[40m",
		PadChar:    "\x00",
		AltChars:   "jjkkllmmnnqqttuuvvwwxx",
		EnterAcs:   "\x1b(0",
		ExitAcs:    "\x1b(B",
		SetCursor:  "\x1b[%i%p1%d;%p2%dH",
		AutoMargin: true,
	})
}
