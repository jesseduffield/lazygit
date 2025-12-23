// Generated automatically.  DO NOT HAND-EDIT.

package simpleterm

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// aka simpleterm
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:          "st",
		Aliases:       []string{"stterm"},
		Columns:       80,
		Lines:         24,
		Colors:        8,
		Clear:         "\x1b[H\x1b[2J",
		EnterCA:       "\x1b[?1049h",
		ExitCA:        "\x1b[?1049l",
		ShowCursor:    "\x1b[?25h",
		HideCursor:    "\x1b[?25l",
		AttrOff:       "\x1b[0m",
		Underline:     "\x1b[4m",
		Bold:          "\x1b[1m",
		Dim:           "\x1b[2m",
		Italic:        "\x1b[3m",
		Blink:         "\x1b[5m",
		Reverse:       "\x1b[7m",
		EnterKeypad:   "\x1b[?1h\x1b=",
		ExitKeypad:    "\x1b[?1l\x1b>",
		SetFg:         "\x1b[3%p1%dm",
		SetBg:         "\x1b[4%p1%dm",
		SetFgBg:       "\x1b[3%p1%d;4%p2%dm",
		ResetFgBg:     "\x1b[39;49m",
		AltChars:      "+C,D-A.B0E``aaffgghFiGjjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:      "\x1b(0",
		ExitAcs:       "\x1b(B",
		EnableAcs:     "\x1b)0",
		StrikeThrough: "\x1b[9m",
		Mouse:         "\x1b[M",
		SetCursor:     "\x1b[%i%p1%d;%p2%dH",
		AutoMargin:    true,
		XTermLike:     true,
	})

	// simpleterm with 256 colors
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:          "st-256color",
		Aliases:       []string{"stterm-256color"},
		Columns:       80,
		Lines:         24,
		Colors:        256,
		Clear:         "\x1b[H\x1b[2J",
		EnterCA:       "\x1b[?1049h",
		ExitCA:        "\x1b[?1049l",
		ShowCursor:    "\x1b[?25h",
		HideCursor:    "\x1b[?25l",
		AttrOff:       "\x1b[0m",
		Underline:     "\x1b[4m",
		Bold:          "\x1b[1m",
		Dim:           "\x1b[2m",
		Italic:        "\x1b[3m",
		Blink:         "\x1b[5m",
		Reverse:       "\x1b[7m",
		EnterKeypad:   "\x1b[?1h\x1b=",
		ExitKeypad:    "\x1b[?1l\x1b>",
		SetFg:         "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;m",
		SetBg:         "\x1b[%?%p1%{8}%<%t4%p1%d%e%p1%{16}%<%t10%p1%{8}%-%d%e48;5;%p1%d%;m",
		SetFgBg:       "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;;%?%p2%{8}%<%t4%p2%d%e%p2%{16}%<%t10%p2%{8}%-%d%e48;5;%p2%d%;m",
		ResetFgBg:     "\x1b[39;49m",
		AltChars:      "+C,D-A.B0E``aaffgghFiGjjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:      "\x1b(0",
		ExitAcs:       "\x1b(B",
		EnableAcs:     "\x1b)0",
		StrikeThrough: "\x1b[9m",
		Mouse:         "\x1b[M",
		SetCursor:     "\x1b[%i%p1%d;%p2%dH",
		AutoMargin:    true,
		XTermLike:     true,
	})
}
