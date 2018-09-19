// +build !windows

package termbox

// Eterm
var eterm_keys = []string{
	"\x1b[11~", "\x1b[12~", "\x1b[13~", "\x1b[14~", "\x1b[15~", "\x1b[17~", "\x1b[18~", "\x1b[19~", "\x1b[20~", "\x1b[21~", "\x1b[23~", "\x1b[24~", "\x1b[2~", "\x1b[3~", "\x1b[7~", "\x1b[8~", "\x1b[5~", "\x1b[6~", "\x1b[A", "\x1b[B", "\x1b[D", "\x1b[C",
}
var eterm_funcs = []string{
	"\x1b7\x1b[?47h", "\x1b[2J\x1b[?47l\x1b8", "\x1b[?25h", "\x1b[?25l", "\x1b[H\x1b[2J", "\x1b[m\x0f", "\x1b[4m", "\x1b[1m", "\x1b[5m", "\x1b[7m", "", "", "", "",
}

// screen
var screen_keys = []string{
	"\x1bOP", "\x1bOQ", "\x1bOR", "\x1bOS", "\x1b[15~", "\x1b[17~", "\x1b[18~", "\x1b[19~", "\x1b[20~", "\x1b[21~", "\x1b[23~", "\x1b[24~", "\x1b[2~", "\x1b[3~", "\x1b[1~", "\x1b[4~", "\x1b[5~", "\x1b[6~", "\x1bOA", "\x1bOB", "\x1bOD", "\x1bOC",
}
var screen_funcs = []string{
	"\x1b[?1049h", "\x1b[?1049l", "\x1b[34h\x1b[?25h", "\x1b[?25l", "\x1b[H\x1b[J", "\x1b[m\x0f", "\x1b[4m", "\x1b[1m", "\x1b[5m", "\x1b[7m", "\x1b[?1h\x1b=", "\x1b[?1l\x1b>", ti_mouse_enter, ti_mouse_leave,
}

// xterm
var xterm_keys = []string{
	"\x1bOP", "\x1bOQ", "\x1bOR", "\x1bOS", "\x1b[15~", "\x1b[17~", "\x1b[18~", "\x1b[19~", "\x1b[20~", "\x1b[21~", "\x1b[23~", "\x1b[24~", "\x1b[2~", "\x1b[3~", "\x1bOH", "\x1bOF", "\x1b[5~", "\x1b[6~", "\x1bOA", "\x1bOB", "\x1bOD", "\x1bOC",
}
var xterm_funcs = []string{
	"\x1b[?1049h", "\x1b[?1049l", "\x1b[?12l\x1b[?25h", "\x1b[?25l", "\x1b[H\x1b[2J", "\x1b(B\x1b[m", "\x1b[4m", "\x1b[1m", "\x1b[5m", "\x1b[7m", "\x1b[?1h\x1b=", "\x1b[?1l\x1b>", ti_mouse_enter, ti_mouse_leave,
}

// rxvt-unicode
var rxvt_unicode_keys = []string{
	"\x1b[11~", "\x1b[12~", "\x1b[13~", "\x1b[14~", "\x1b[15~", "\x1b[17~", "\x1b[18~", "\x1b[19~", "\x1b[20~", "\x1b[21~", "\x1b[23~", "\x1b[24~", "\x1b[2~", "\x1b[3~", "\x1b[7~", "\x1b[8~", "\x1b[5~", "\x1b[6~", "\x1b[A", "\x1b[B", "\x1b[D", "\x1b[C",
}
var rxvt_unicode_funcs = []string{
	"\x1b[?1049h", "\x1b[r\x1b[?1049l", "\x1b[?25h", "\x1b[?25l", "\x1b[H\x1b[2J", "\x1b[m\x1b(B", "\x1b[4m", "\x1b[1m", "\x1b[5m", "\x1b[7m", "\x1b=", "\x1b>", ti_mouse_enter, ti_mouse_leave,
}

// linux
var linux_keys = []string{
	"\x1b[[A", "\x1b[[B", "\x1b[[C", "\x1b[[D", "\x1b[[E", "\x1b[17~", "\x1b[18~", "\x1b[19~", "\x1b[20~", "\x1b[21~", "\x1b[23~", "\x1b[24~", "\x1b[2~", "\x1b[3~", "\x1b[1~", "\x1b[4~", "\x1b[5~", "\x1b[6~", "\x1b[A", "\x1b[B", "\x1b[D", "\x1b[C",
}
var linux_funcs = []string{
	"", "", "\x1b[?25h\x1b[?0c", "\x1b[?25l\x1b[?1c", "\x1b[H\x1b[J", "\x1b[0;10m", "\x1b[4m", "\x1b[1m", "\x1b[5m", "\x1b[7m", "", "", "", "",
}

// rxvt-256color
var rxvt_256color_keys = []string{
	"\x1b[11~", "\x1b[12~", "\x1b[13~", "\x1b[14~", "\x1b[15~", "\x1b[17~", "\x1b[18~", "\x1b[19~", "\x1b[20~", "\x1b[21~", "\x1b[23~", "\x1b[24~", "\x1b[2~", "\x1b[3~", "\x1b[7~", "\x1b[8~", "\x1b[5~", "\x1b[6~", "\x1b[A", "\x1b[B", "\x1b[D", "\x1b[C",
}
var rxvt_256color_funcs = []string{
	"\x1b7\x1b[?47h", "\x1b[2J\x1b[?47l\x1b8", "\x1b[?25h", "\x1b[?25l", "\x1b[H\x1b[2J", "\x1b[m\x0f", "\x1b[4m", "\x1b[1m", "\x1b[5m", "\x1b[7m", "\x1b=", "\x1b>", ti_mouse_enter, ti_mouse_leave,
}

var terms = []struct {
	name  string
	keys  []string
	funcs []string
}{
	{"Eterm", eterm_keys, eterm_funcs},
	{"screen", screen_keys, screen_funcs},
	{"xterm", xterm_keys, xterm_funcs},
	{"rxvt-unicode", rxvt_unicode_keys, rxvt_unicode_funcs},
	{"linux", linux_keys, linux_funcs},
	{"rxvt-256color", rxvt_256color_keys, rxvt_256color_funcs},
}
