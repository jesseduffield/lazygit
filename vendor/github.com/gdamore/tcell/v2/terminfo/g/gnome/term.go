// Generated automatically.  DO NOT HAND-EDIT.

package gnome

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// GNOME Terminal
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:              "gnome",
		Columns:           80,
		Lines:             24,
		Colors:            8,
		Clear:             "\x1b[H\x1b[2J",
		EnterCA:           "\x1b7\x1b[?47h",
		ExitCA:            "\x1b[2J\x1b[?47l\x1b8",
		ShowCursor:        "\x1b[?25h",
		HideCursor:        "\x1b[?25l",
		AttrOff:           "\x1b[0m\x0f",
		Underline:         "\x1b[4m",
		Bold:              "\x1b[1m",
		Dim:               "\x1b[2m",
		Italic:            "\x1b[3m",
		Reverse:           "\x1b[7m",
		EnterKeypad:       "\x1b[?1h\x1b=",
		ExitKeypad:        "\x1b[?1l\x1b>",
		SetFg:             "\x1b[3%p1%dm",
		SetBg:             "\x1b[4%p1%dm",
		SetFgBg:           "\x1b[3%p1%d;4%p2%dm",
		ResetFgBg:         "\x1b[39;49m",
		PadChar:           "\x00",
		AltChars:          "``aaffggiijjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:          "\x0e",
		ExitAcs:           "\x0f",
		EnableAcs:         "\x1b)0",
		EnableAutoMargin:  "\x1b[?7h",
		DisableAutoMargin: "\x1b[?7l",
		Mouse:             "\x1b[M",
		SetCursor:         "\x1b[%i%p1%d;%p2%dH",
		AutoMargin:        true,
		XTermLike:         true,
	})

	// GNOME Terminal with xterm 256-colors
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:              "gnome-256color",
		Columns:           80,
		Lines:             24,
		Colors:            256,
		Clear:             "\x1b[H\x1b[2J",
		EnterCA:           "\x1b7\x1b[?47h",
		ExitCA:            "\x1b[2J\x1b[?47l\x1b8",
		ShowCursor:        "\x1b[?25h",
		HideCursor:        "\x1b[?25l",
		AttrOff:           "\x1b[0m\x0f",
		Underline:         "\x1b[4m",
		Bold:              "\x1b[1m",
		Dim:               "\x1b[2m",
		Italic:            "\x1b[3m",
		Reverse:           "\x1b[7m",
		EnterKeypad:       "\x1b[?1h\x1b=",
		ExitKeypad:        "\x1b[?1l\x1b>",
		SetFg:             "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;m",
		SetBg:             "\x1b[%?%p1%{8}%<%t4%p1%d%e%p1%{16}%<%t10%p1%{8}%-%d%e48;5;%p1%d%;m",
		SetFgBg:           "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;;%?%p2%{8}%<%t4%p2%d%e%p2%{16}%<%t10%p2%{8}%-%d%e48;5;%p2%d%;m",
		ResetFgBg:         "\x1b[39;49m",
		PadChar:           "\x00",
		AltChars:          "``aaffggiijjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:          "\x0e",
		ExitAcs:           "\x0f",
		EnableAcs:         "\x1b)0",
		EnableAutoMargin:  "\x1b[?7h",
		DisableAutoMargin: "\x1b[?7l",
		Mouse:             "\x1b[M",
		SetCursor:         "\x1b[%i%p1%d;%p2%dH",
		AutoMargin:        true,
		XTermLike:         true,
	})
}
