// Generated automatically.  DO NOT HAND-EDIT.

package beterm

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// BeOS Terminal
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:         "beterm",
		Columns:      80,
		Lines:        25,
		Colors:       8,
		Bell:         "\a",
		Clear:        "\x1b[H\x1b[J",
		AttrOff:      "\x1b[0;10m",
		Underline:    "\x1b[4m",
		Bold:         "\x1b[1m",
		Reverse:      "\x1b[7m",
		EnterKeypad:  "\x1b[?4h",
		ExitKeypad:   "\x1b[?4l",
		SetFg:        "\x1b[3%p1%dm",
		SetBg:        "\x1b[4%p1%dm",
		SetFgBg:      "\x1b[3%p1%d;4%p2%dm",
		ResetFgBg:    "\x1b[m",
		PadChar:      "\x00",
		SetCursor:    "\x1b[%i%p1%d;%p2%dH",
		CursorBack1:  "\b",
		CursorUp1:    "\x1b[A",
		KeyUp:        "\x1b[A",
		KeyDown:      "\x1b[B",
		KeyRight:     "\x1b[C",
		KeyLeft:      "\x1b[D",
		KeyInsert:    "\x1b[2~",
		KeyDelete:    "\x1b[3~",
		KeyBackspace: "\b",
		KeyHome:      "\x1b[1~",
		KeyEnd:       "\x1b[4~",
		KeyPgUp:      "\x1b[5~",
		KeyPgDn:      "\x1b[6~",
		KeyF1:        "\x1b[11~",
		KeyF2:        "\x1b[12~",
		KeyF3:        "\x1b[13~",
		KeyF4:        "\x1b[14~",
		KeyF5:        "\x1b[15~",
		KeyF6:        "\x1b[16~",
		KeyF7:        "\x1b[17~",
		KeyF8:        "\x1b[18~",
		KeyF9:        "\x1b[19~",
		KeyF10:       "\x1b[20~",
		KeyF11:       "\x1b[21~",
		KeyF12:       "\x1b[22~",
		AutoMargin:   true,
		InsertChar:   "\x1b[@",
	})
}
