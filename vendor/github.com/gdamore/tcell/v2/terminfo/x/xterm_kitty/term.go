// Generated automatically.  DO NOT HAND-EDIT.

package xterm_kitty

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// KovIdTTY
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:              "xterm-kitty",
		Columns:           80,
		Lines:             24,
		Colors:            256,
		Clear:             "\x1b[H\x1b[2J",
		EnterCA:           "\x1b[?1049h",
		ExitCA:            "\x1b[?1049l",
		ShowCursor:        "\x1b[?12h\x1b[?25h",
		HideCursor:        "\x1b[?25l",
		AttrOff:           "\x1b(B\x1b[m",
		Underline:         "\x1b[4m",
		Bold:              "\x1b[1m",
		Dim:               "\x1b[2m",
		Italic:            "\x1b[3m",
		Reverse:           "\x1b[7m",
		EnterKeypad:       "\x1b[?1h",
		ExitKeypad:        "\x1b[?1l",
		SetFg:             "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;m",
		SetBg:             "\x1b[%?%p1%{8}%<%t4%p1%d%e%p1%{16}%<%t10%p1%{8}%-%d%e48;5;%p1%d%;m",
		SetFgBg:           "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;;%?%p2%{8}%<%t4%p2%d%e%p2%{16}%<%t10%p2%{8}%-%d%e48;5;%p2%d%;m",
		ResetFgBg:         "\x1b[39;49m",
		AltChars:          "++,,--..00``aaffgghhiijjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:          "\x1b(0",
		ExitAcs:           "\x1b(B",
		EnableAutoMargin:  "\x1b[?7h",
		DisableAutoMargin: "\x1b[?7l",
		StrikeThrough:     "\x1b[9m",
		Mouse:             "\x1b[M",
		SetCursor:         "\x1b[%i%p1%d;%p2%dH",
		TrueColor:         true,
		AutoMargin:        true,
		XTermLike:         true,
	})
}
