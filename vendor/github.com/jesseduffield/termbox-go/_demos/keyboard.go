package main

import "github.com/nsf/termbox-go"
import "fmt"

type key struct {
	x  int
	y  int
	ch rune
}

var K_ESC = []key{{1, 1, 'E'}, {2, 1, 'S'}, {3, 1, 'C'}}
var K_F1 = []key{{6, 1, 'F'}, {7, 1, '1'}}
var K_F2 = []key{{9, 1, 'F'}, {10, 1, '2'}}
var K_F3 = []key{{12, 1, 'F'}, {13, 1, '3'}}
var K_F4 = []key{{15, 1, 'F'}, {16, 1, '4'}}
var K_F5 = []key{{19, 1, 'F'}, {20, 1, '5'}}
var K_F6 = []key{{22, 1, 'F'}, {23, 1, '6'}}
var K_F7 = []key{{25, 1, 'F'}, {26, 1, '7'}}
var K_F8 = []key{{28, 1, 'F'}, {29, 1, '8'}}
var K_F9 = []key{{33, 1, 'F'}, {34, 1, '9'}}
var K_F10 = []key{{36, 1, 'F'}, {37, 1, '1'}, {38, 1, '0'}}
var K_F11 = []key{{40, 1, 'F'}, {41, 1, '1'}, {42, 1, '1'}}
var K_F12 = []key{{44, 1, 'F'}, {45, 1, '1'}, {46, 1, '2'}}
var K_PRN = []key{{50, 1, 'P'}, {51, 1, 'R'}, {52, 1, 'N'}}
var K_SCR = []key{{54, 1, 'S'}, {55, 1, 'C'}, {56, 1, 'R'}}
var K_BRK = []key{{58, 1, 'B'}, {59, 1, 'R'}, {60, 1, 'K'}}
var K_LED1 = []key{{66, 1, '-'}}
var K_LED2 = []key{{70, 1, '-'}}
var K_LED3 = []key{{74, 1, '-'}}
var K_TILDE = []key{{1, 4, '`'}}
var K_TILDE_SHIFT = []key{{1, 4, '~'}}
var K_1 = []key{{4, 4, '1'}}
var K_1_SHIFT = []key{{4, 4, '!'}}
var K_2 = []key{{7, 4, '2'}}
var K_2_SHIFT = []key{{7, 4, '@'}}
var K_3 = []key{{10, 4, '3'}}
var K_3_SHIFT = []key{{10, 4, '#'}}
var K_4 = []key{{13, 4, '4'}}
var K_4_SHIFT = []key{{13, 4, '$'}}
var K_5 = []key{{16, 4, '5'}}
var K_5_SHIFT = []key{{16, 4, '%'}}
var K_6 = []key{{19, 4, '6'}}
var K_6_SHIFT = []key{{19, 4, '^'}}
var K_7 = []key{{22, 4, '7'}}
var K_7_SHIFT = []key{{22, 4, '&'}}
var K_8 = []key{{25, 4, '8'}}
var K_8_SHIFT = []key{{25, 4, '*'}}
var K_9 = []key{{28, 4, '9'}}
var K_9_SHIFT = []key{{28, 4, '('}}
var K_0 = []key{{31, 4, '0'}}
var K_0_SHIFT = []key{{31, 4, ')'}}
var K_MINUS = []key{{34, 4, '-'}}
var K_MINUS_SHIFT = []key{{34, 4, '_'}}
var K_EQUALS = []key{{37, 4, '='}}
var K_EQUALS_SHIFT = []key{{37, 4, '+'}}
var K_BACKSLASH = []key{{40, 4, '\\'}}
var K_BACKSLASH_SHIFT = []key{{40, 4, '|'}}
var K_BACKSPACE = []key{{44, 4, 0x2190}, {45, 4, 0x2500}, {46, 4, 0x2500}}
var K_INS = []key{{50, 4, 'I'}, {51, 4, 'N'}, {52, 4, 'S'}}
var K_HOM = []key{{54, 4, 'H'}, {55, 4, 'O'}, {56, 4, 'M'}}
var K_PGU = []key{{58, 4, 'P'}, {59, 4, 'G'}, {60, 4, 'U'}}
var K_K_NUMLOCK = []key{{65, 4, 'N'}}
var K_K_SLASH = []key{{68, 4, '/'}}
var K_K_STAR = []key{{71, 4, '*'}}
var K_K_MINUS = []key{{74, 4, '-'}}
var K_TAB = []key{{1, 6, 'T'}, {2, 6, 'A'}, {3, 6, 'B'}}
var K_q = []key{{6, 6, 'q'}}
var K_Q = []key{{6, 6, 'Q'}}
var K_w = []key{{9, 6, 'w'}}
var K_W = []key{{9, 6, 'W'}}
var K_e = []key{{12, 6, 'e'}}
var K_E = []key{{12, 6, 'E'}}
var K_r = []key{{15, 6, 'r'}}
var K_R = []key{{15, 6, 'R'}}
var K_t = []key{{18, 6, 't'}}
var K_T = []key{{18, 6, 'T'}}
var K_y = []key{{21, 6, 'y'}}
var K_Y = []key{{21, 6, 'Y'}}
var K_u = []key{{24, 6, 'u'}}
var K_U = []key{{24, 6, 'U'}}
var K_i = []key{{27, 6, 'i'}}
var K_I = []key{{27, 6, 'I'}}
var K_o = []key{{30, 6, 'o'}}
var K_O = []key{{30, 6, 'O'}}
var K_p = []key{{33, 6, 'p'}}
var K_P = []key{{33, 6, 'P'}}
var K_LSQB = []key{{36, 6, '['}}
var K_LCUB = []key{{36, 6, '{'}}
var K_RSQB = []key{{39, 6, ']'}}
var K_RCUB = []key{{39, 6, '}'}}
var K_ENTER = []key{
	{43, 6, 0x2591}, {44, 6, 0x2591}, {45, 6, 0x2591}, {46, 6, 0x2591},
	{43, 7, 0x2591}, {44, 7, 0x2591}, {45, 7, 0x21B5}, {46, 7, 0x2591},
	{41, 8, 0x2591}, {42, 8, 0x2591}, {43, 8, 0x2591}, {44, 8, 0x2591},
	{45, 8, 0x2591}, {46, 8, 0x2591},
}
var K_DEL = []key{{50, 6, 'D'}, {51, 6, 'E'}, {52, 6, 'L'}}
var K_END = []key{{54, 6, 'E'}, {55, 6, 'N'}, {56, 6, 'D'}}
var K_PGD = []key{{58, 6, 'P'}, {59, 6, 'G'}, {60, 6, 'D'}}
var K_K_7 = []key{{65, 6, '7'}}
var K_K_8 = []key{{68, 6, '8'}}
var K_K_9 = []key{{71, 6, '9'}}
var K_K_PLUS = []key{{74, 6, ' '}, {74, 7, '+'}, {74, 8, ' '}}
var K_CAPS = []key{{1, 8, 'C'}, {2, 8, 'A'}, {3, 8, 'P'}, {4, 8, 'S'}}
var K_a = []key{{7, 8, 'a'}}
var K_A = []key{{7, 8, 'A'}}
var K_s = []key{{10, 8, 's'}}
var K_S = []key{{10, 8, 'S'}}
var K_d = []key{{13, 8, 'd'}}
var K_D = []key{{13, 8, 'D'}}
var K_f = []key{{16, 8, 'f'}}
var K_F = []key{{16, 8, 'F'}}
var K_g = []key{{19, 8, 'g'}}
var K_G = []key{{19, 8, 'G'}}
var K_h = []key{{22, 8, 'h'}}
var K_H = []key{{22, 8, 'H'}}
var K_j = []key{{25, 8, 'j'}}
var K_J = []key{{25, 8, 'J'}}
var K_k = []key{{28, 8, 'k'}}
var K_K = []key{{28, 8, 'K'}}
var K_l = []key{{31, 8, 'l'}}
var K_L = []key{{31, 8, 'L'}}
var K_SEMICOLON = []key{{34, 8, ';'}}
var K_PARENTHESIS = []key{{34, 8, ':'}}
var K_QUOTE = []key{{37, 8, '\''}}
var K_DOUBLEQUOTE = []key{{37, 8, '"'}}
var K_K_4 = []key{{65, 8, '4'}}
var K_K_5 = []key{{68, 8, '5'}}
var K_K_6 = []key{{71, 8, '6'}}
var K_LSHIFT = []key{{1, 10, 'S'}, {2, 10, 'H'}, {3, 10, 'I'}, {4, 10, 'F'}, {5, 10, 'T'}}
var K_z = []key{{9, 10, 'z'}}
var K_Z = []key{{9, 10, 'Z'}}
var K_x = []key{{12, 10, 'x'}}
var K_X = []key{{12, 10, 'X'}}
var K_c = []key{{15, 10, 'c'}}
var K_C = []key{{15, 10, 'C'}}
var K_v = []key{{18, 10, 'v'}}
var K_V = []key{{18, 10, 'V'}}
var K_b = []key{{21, 10, 'b'}}
var K_B = []key{{21, 10, 'B'}}
var K_n = []key{{24, 10, 'n'}}
var K_N = []key{{24, 10, 'N'}}
var K_m = []key{{27, 10, 'm'}}
var K_M = []key{{27, 10, 'M'}}
var K_COMMA = []key{{30, 10, ','}}
var K_LANB = []key{{30, 10, '<'}}
var K_PERIOD = []key{{33, 10, '.'}}
var K_RANB = []key{{33, 10, '>'}}
var K_SLASH = []key{{36, 10, '/'}}
var K_QUESTION = []key{{36, 10, '?'}}
var K_RSHIFT = []key{{42, 10, 'S'}, {43, 10, 'H'}, {44, 10, 'I'}, {45, 10, 'F'}, {46, 10, 'T'}}
var K_ARROW_UP = []key{{54, 10, '('}, {55, 10, 0x2191}, {56, 10, ')'}}
var K_K_1 = []key{{65, 10, '1'}}
var K_K_2 = []key{{68, 10, '2'}}
var K_K_3 = []key{{71, 10, '3'}}
var K_K_ENTER = []key{{74, 10, 0x2591}, {74, 11, 0x2591}, {74, 12, 0x2591}}
var K_LCTRL = []key{{1, 12, 'C'}, {2, 12, 'T'}, {3, 12, 'R'}, {4, 12, 'L'}}
var K_LWIN = []key{{6, 12, 'W'}, {7, 12, 'I'}, {8, 12, 'N'}}
var K_LALT = []key{{10, 12, 'A'}, {11, 12, 'L'}, {12, 12, 'T'}}
var K_SPACE = []key{
	{14, 12, ' '}, {15, 12, ' '}, {16, 12, ' '}, {17, 12, ' '}, {18, 12, ' '},
	{19, 12, 'S'}, {20, 12, 'P'}, {21, 12, 'A'}, {22, 12, 'C'}, {23, 12, 'E'},
	{24, 12, ' '}, {25, 12, ' '}, {26, 12, ' '}, {27, 12, ' '}, {28, 12, ' '},
}
var K_RALT = []key{{30, 12, 'A'}, {31, 12, 'L'}, {32, 12, 'T'}}
var K_RWIN = []key{{34, 12, 'W'}, {35, 12, 'I'}, {36, 12, 'N'}}
var K_RPROP = []key{{38, 12, 'P'}, {39, 12, 'R'}, {40, 12, 'O'}, {41, 12, 'P'}}
var K_RCTRL = []key{{43, 12, 'C'}, {44, 12, 'T'}, {45, 12, 'R'}, {46, 12, 'L'}}
var K_ARROW_LEFT = []key{{50, 12, '('}, {51, 12, 0x2190}, {52, 12, ')'}}
var K_ARROW_DOWN = []key{{54, 12, '('}, {55, 12, 0x2193}, {56, 12, ')'}}
var K_ARROW_RIGHT = []key{{58, 12, '('}, {59, 12, 0x2192}, {60, 12, ')'}}
var K_K_0 = []key{{65, 12, ' '}, {66, 12, '0'}, {67, 12, ' '}, {68, 12, ' '}}
var K_K_PERIOD = []key{{71, 12, '.'}}

type combo struct {
	keys [][]key
}

var combos = []combo{
	{[][]key{K_TILDE, K_2, K_SPACE, K_LCTRL, K_RCTRL}},
	{[][]key{K_A, K_LCTRL, K_RCTRL}},
	{[][]key{K_B, K_LCTRL, K_RCTRL}},
	{[][]key{K_C, K_LCTRL, K_RCTRL}},
	{[][]key{K_D, K_LCTRL, K_RCTRL}},
	{[][]key{K_E, K_LCTRL, K_RCTRL}},
	{[][]key{K_F, K_LCTRL, K_RCTRL}},
	{[][]key{K_G, K_LCTRL, K_RCTRL}},
	{[][]key{K_H, K_BACKSPACE, K_LCTRL, K_RCTRL}},
	{[][]key{K_I, K_TAB, K_LCTRL, K_RCTRL}},
	{[][]key{K_J, K_LCTRL, K_RCTRL}},
	{[][]key{K_K, K_LCTRL, K_RCTRL}},
	{[][]key{K_L, K_LCTRL, K_RCTRL}},
	{[][]key{K_M, K_ENTER, K_K_ENTER, K_LCTRL, K_RCTRL}},
	{[][]key{K_N, K_LCTRL, K_RCTRL}},
	{[][]key{K_O, K_LCTRL, K_RCTRL}},
	{[][]key{K_P, K_LCTRL, K_RCTRL}},
	{[][]key{K_Q, K_LCTRL, K_RCTRL}},
	{[][]key{K_R, K_LCTRL, K_RCTRL}},
	{[][]key{K_S, K_LCTRL, K_RCTRL}},
	{[][]key{K_T, K_LCTRL, K_RCTRL}},
	{[][]key{K_U, K_LCTRL, K_RCTRL}},
	{[][]key{K_V, K_LCTRL, K_RCTRL}},
	{[][]key{K_W, K_LCTRL, K_RCTRL}},
	{[][]key{K_X, K_LCTRL, K_RCTRL}},
	{[][]key{K_Y, K_LCTRL, K_RCTRL}},
	{[][]key{K_Z, K_LCTRL, K_RCTRL}},
	{[][]key{K_LSQB, K_ESC, K_3, K_LCTRL, K_RCTRL}},
	{[][]key{K_4, K_BACKSLASH, K_LCTRL, K_RCTRL}},
	{[][]key{K_RSQB, K_5, K_LCTRL, K_RCTRL}},
	{[][]key{K_6, K_LCTRL, K_RCTRL}},
	{[][]key{K_7, K_SLASH, K_MINUS_SHIFT, K_LCTRL, K_RCTRL}},
	{[][]key{K_SPACE}},
	{[][]key{K_1_SHIFT, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_DOUBLEQUOTE, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_3_SHIFT, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_4_SHIFT, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_5_SHIFT, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_7_SHIFT, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_QUOTE}},
	{[][]key{K_9_SHIFT, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_0_SHIFT, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_8_SHIFT, K_K_STAR, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_EQUALS_SHIFT, K_K_PLUS, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_COMMA}},
	{[][]key{K_MINUS, K_K_MINUS}},
	{[][]key{K_PERIOD, K_K_PERIOD}},
	{[][]key{K_SLASH, K_K_SLASH}},
	{[][]key{K_0, K_K_0}},
	{[][]key{K_1, K_K_1}},
	{[][]key{K_2, K_K_2}},
	{[][]key{K_3, K_K_3}},
	{[][]key{K_4, K_K_4}},
	{[][]key{K_5, K_K_5}},
	{[][]key{K_6, K_K_6}},
	{[][]key{K_7, K_K_7}},
	{[][]key{K_8, K_K_8}},
	{[][]key{K_9, K_K_9}},
	{[][]key{K_PARENTHESIS, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_SEMICOLON}},
	{[][]key{K_LANB, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_EQUALS}},
	{[][]key{K_RANB, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_QUESTION, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_2_SHIFT, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_A, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_B, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_C, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_D, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_E, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_F, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_G, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_H, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_I, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_J, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_K, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_L, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_M, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_N, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_O, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_P, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_Q, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_R, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_S, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_T, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_U, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_V, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_W, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_X, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_Y, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_Z, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_LSQB}},
	{[][]key{K_BACKSLASH}},
	{[][]key{K_RSQB}},
	{[][]key{K_6_SHIFT, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_MINUS_SHIFT, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_TILDE}},
	{[][]key{K_a}},
	{[][]key{K_b}},
	{[][]key{K_c}},
	{[][]key{K_d}},
	{[][]key{K_e}},
	{[][]key{K_f}},
	{[][]key{K_g}},
	{[][]key{K_h}},
	{[][]key{K_i}},
	{[][]key{K_j}},
	{[][]key{K_k}},
	{[][]key{K_l}},
	{[][]key{K_m}},
	{[][]key{K_n}},
	{[][]key{K_o}},
	{[][]key{K_p}},
	{[][]key{K_q}},
	{[][]key{K_r}},
	{[][]key{K_s}},
	{[][]key{K_t}},
	{[][]key{K_u}},
	{[][]key{K_v}},
	{[][]key{K_w}},
	{[][]key{K_x}},
	{[][]key{K_y}},
	{[][]key{K_z}},
	{[][]key{K_LCUB, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_BACKSLASH_SHIFT, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_RCUB, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_TILDE_SHIFT, K_LSHIFT, K_RSHIFT}},
	{[][]key{K_8, K_BACKSPACE, K_LCTRL, K_RCTRL}},
}

var func_combos = []combo{
	{[][]key{K_F1}},
	{[][]key{K_F2}},
	{[][]key{K_F3}},
	{[][]key{K_F4}},
	{[][]key{K_F5}},
	{[][]key{K_F6}},
	{[][]key{K_F7}},
	{[][]key{K_F8}},
	{[][]key{K_F9}},
	{[][]key{K_F10}},
	{[][]key{K_F11}},
	{[][]key{K_F12}},
	{[][]key{K_INS}},
	{[][]key{K_DEL}},
	{[][]key{K_HOM}},
	{[][]key{K_END}},
	{[][]key{K_PGU}},
	{[][]key{K_PGD}},
	{[][]key{K_ARROW_UP}},
	{[][]key{K_ARROW_DOWN}},
	{[][]key{K_ARROW_LEFT}},
	{[][]key{K_ARROW_RIGHT}},
}

func print_tb(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func printf_tb(x, y int, fg, bg termbox.Attribute, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	print_tb(x, y, fg, bg, s)
}

func draw_key(k []key, fg, bg termbox.Attribute) {
	for _, k := range k {
		termbox.SetCell(k.x+2, k.y+4, k.ch, fg, bg)
	}
}

func draw_keyboard() {
	termbox.SetCell(0, 0, 0x250C, termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(79, 0, 0x2510, termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(0, 23, 0x2514, termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(79, 23, 0x2518, termbox.ColorWhite, termbox.ColorBlack)

	for i := 1; i < 79; i++ {
		termbox.SetCell(i, 0, 0x2500, termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(i, 23, 0x2500, termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(i, 17, 0x2500, termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(i, 4, 0x2500, termbox.ColorWhite, termbox.ColorBlack)
	}
	for i := 1; i < 23; i++ {
		termbox.SetCell(0, i, 0x2502, termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(79, i, 0x2502, termbox.ColorWhite, termbox.ColorBlack)
	}
	termbox.SetCell(0, 17, 0x251C, termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(79, 17, 0x2524, termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(0, 4, 0x251C, termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(79, 4, 0x2524, termbox.ColorWhite, termbox.ColorBlack)
	for i := 5; i < 17; i++ {
		termbox.SetCell(1, i, 0x2588, termbox.ColorYellow, termbox.ColorYellow)
		termbox.SetCell(78, i, 0x2588, termbox.ColorYellow, termbox.ColorYellow)
	}

	draw_key(K_ESC, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_F1, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_F2, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_F3, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_F4, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_F5, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_F6, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_F7, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_F8, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_F9, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_F10, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_F11, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_F12, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_PRN, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_SCR, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_BRK, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_LED1, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_LED2, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_LED3, termbox.ColorWhite, termbox.ColorBlue)

	draw_key(K_TILDE, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_1, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_2, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_3, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_4, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_5, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_6, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_7, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_8, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_9, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_0, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_MINUS, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_EQUALS, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_BACKSLASH, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_BACKSPACE, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_INS, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_HOM, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_PGU, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_K_NUMLOCK, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_K_SLASH, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_K_STAR, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_K_MINUS, termbox.ColorWhite, termbox.ColorBlue)

	draw_key(K_TAB, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_q, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_w, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_e, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_r, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_t, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_y, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_u, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_i, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_o, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_p, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_LSQB, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_RSQB, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_ENTER, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_DEL, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_END, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_PGD, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_K_7, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_K_8, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_K_9, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_K_PLUS, termbox.ColorWhite, termbox.ColorBlue)

	draw_key(K_CAPS, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_a, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_s, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_d, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_f, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_g, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_h, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_j, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_k, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_l, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_SEMICOLON, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_QUOTE, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_K_4, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_K_5, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_K_6, termbox.ColorWhite, termbox.ColorBlue)

	draw_key(K_LSHIFT, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_z, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_x, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_c, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_v, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_b, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_n, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_m, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_COMMA, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_PERIOD, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_SLASH, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_RSHIFT, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_ARROW_UP, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_K_1, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_K_2, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_K_3, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_K_ENTER, termbox.ColorWhite, termbox.ColorBlue)

	draw_key(K_LCTRL, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_LWIN, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_LALT, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_SPACE, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_RCTRL, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_RPROP, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_RWIN, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_RALT, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_ARROW_LEFT, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_ARROW_DOWN, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_ARROW_RIGHT, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_K_0, termbox.ColorWhite, termbox.ColorBlue)
	draw_key(K_K_PERIOD, termbox.ColorWhite, termbox.ColorBlue)

	printf_tb(33, 1, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorBlack, "Keyboard demo!")
	printf_tb(21, 2, termbox.ColorMagenta, termbox.ColorBlack, "(press CTRL+X and then CTRL+Q to exit)")
	printf_tb(15, 3, termbox.ColorMagenta, termbox.ColorBlack, "(press CTRL+X and then CTRL+C to change input mode)")

	inputmode := termbox.SetInputMode(termbox.InputCurrent)
	inputmode_str := ""
	switch {
	case inputmode&termbox.InputEsc != 0:
		inputmode_str = "termbox.InputEsc"
	case inputmode&termbox.InputAlt != 0:
		inputmode_str = "termbox.InputAlt"
	}

	if inputmode&termbox.InputMouse != 0 {
		inputmode_str += " | termbox.InputMouse"
	}
	printf_tb(3, 18, termbox.ColorWhite, termbox.ColorBlack, "Input mode: %s", inputmode_str)
}

var fcmap = []string{
	"CTRL+2, CTRL+~",
	"CTRL+A",
	"CTRL+B",
	"CTRL+C",
	"CTRL+D",
	"CTRL+E",
	"CTRL+F",
	"CTRL+G",
	"CTRL+H, BACKSPACE",
	"CTRL+I, TAB",
	"CTRL+J",
	"CTRL+K",
	"CTRL+L",
	"CTRL+M, ENTER",
	"CTRL+N",
	"CTRL+O",
	"CTRL+P",
	"CTRL+Q",
	"CTRL+R",
	"CTRL+S",
	"CTRL+T",
	"CTRL+U",
	"CTRL+V",
	"CTRL+W",
	"CTRL+X",
	"CTRL+Y",
	"CTRL+Z",
	"CTRL+3, ESC, CTRL+[",
	"CTRL+4, CTRL+\\",
	"CTRL+5, CTRL+]",
	"CTRL+6",
	"CTRL+7, CTRL+/, CTRL+_",
	"SPACE",
}

var fkmap = []string{
	"F1",
	"F2",
	"F3",
	"F4",
	"F5",
	"F6",
	"F7",
	"F8",
	"F9",
	"F10",
	"F11",
	"F12",
	"INSERT",
	"DELETE",
	"HOME",
	"END",
	"PGUP",
	"PGDN",
	"ARROW UP",
	"ARROW DOWN",
	"ARROW LEFT",
	"ARROW RIGHT",
}

func funckeymap(k termbox.Key) string {
	if k == termbox.KeyCtrl8 {
		return "CTRL+8, BACKSPACE 2" /* 0x7F */
	} else if k >= termbox.KeyArrowRight && k <= 0xFFFF {
		return fkmap[0xFFFF-k]
	} else if k <= termbox.KeySpace {
		return fcmap[k]
	}
	return "UNKNOWN"
}

func pretty_print_press(ev *termbox.Event) {
	printf_tb(3, 19, termbox.ColorWhite, termbox.ColorBlack, "Key: ")
	printf_tb(8, 19, termbox.ColorYellow, termbox.ColorBlack, "decimal: %d", ev.Key)
	printf_tb(8, 20, termbox.ColorGreen, termbox.ColorBlack, "hex:     0x%X", ev.Key)
	printf_tb(8, 21, termbox.ColorCyan, termbox.ColorBlack, "octal:   0%o", ev.Key)
	printf_tb(8, 22, termbox.ColorRed, termbox.ColorBlack, "string:  %s", funckeymap(ev.Key))

	printf_tb(54, 19, termbox.ColorWhite, termbox.ColorBlack, "Char: ")
	printf_tb(60, 19, termbox.ColorYellow, termbox.ColorBlack, "decimal: %d", ev.Ch)
	printf_tb(60, 20, termbox.ColorGreen, termbox.ColorBlack, "hex:     0x%X", ev.Ch)
	printf_tb(60, 21, termbox.ColorCyan, termbox.ColorBlack, "octal:   0%o", ev.Ch)
	printf_tb(60, 22, termbox.ColorRed, termbox.ColorBlack, "string:  %s", string(ev.Ch))

	modifier := "none"
	if ev.Mod != 0 {
		modifier = "termbox.ModAlt"
	}
	printf_tb(54, 18, termbox.ColorWhite, termbox.ColorBlack, "Modifier: %s", modifier)
}

func pretty_print_resize(ev *termbox.Event) {
	printf_tb(3, 19, termbox.ColorWhite, termbox.ColorBlack, "Resize event: %d x %d", ev.Width, ev.Height)
}

var counter = 0

func pretty_print_mouse(ev *termbox.Event) {
	printf_tb(3, 19, termbox.ColorWhite, termbox.ColorBlack, "Mouse event: %d x %d", ev.MouseX, ev.MouseY)
	button := ""
	switch ev.Key {
	case termbox.MouseLeft:
		button = "MouseLeft: %d"
	case termbox.MouseMiddle:
		button = "MouseMiddle: %d"
	case termbox.MouseRight:
		button = "MouseRight: %d"
	case termbox.MouseWheelUp:
		button = "MouseWheelUp: %d"
	case termbox.MouseWheelDown:
		button = "MouseWheelDown: %d"
	case termbox.MouseRelease:
		button = "MouseRelease: %d"
	}
	if ev.Mod&termbox.ModMotion != 0 {
		button += "*"
	}
	counter++
	printf_tb(43, 19, termbox.ColorWhite, termbox.ColorBlack, "Key: ")
	printf_tb(48, 19, termbox.ColorYellow, termbox.ColorBlack, button, counter)
}

func dispatch_press(ev *termbox.Event) {
	if ev.Mod&termbox.ModAlt != 0 {
		draw_key(K_LALT, termbox.ColorWhite, termbox.ColorRed)
		draw_key(K_RALT, termbox.ColorWhite, termbox.ColorRed)
	}

	var k *combo
	if ev.Key >= termbox.KeyArrowRight {
		k = &func_combos[0xFFFF-ev.Key]
	} else if ev.Ch < 128 {
		if ev.Ch == 0 && ev.Key < 128 {
			k = &combos[ev.Key]
		} else {
			k = &combos[ev.Ch]
		}
	}
	if k == nil {
		return
	}

	keys := k.keys
	for _, k := range keys {
		draw_key(k, termbox.ColorWhite, termbox.ColorRed)
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	draw_keyboard()
	termbox.Flush()
	inputmode := 0
	ctrlxpressed := false
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlS && ctrlxpressed {
				termbox.Sync()
			}
			if ev.Key == termbox.KeyCtrlQ && ctrlxpressed {
				break loop
			}
			if ev.Key == termbox.KeyCtrlC && ctrlxpressed {
				chmap := []termbox.InputMode{
					termbox.InputEsc | termbox.InputMouse,
					termbox.InputAlt | termbox.InputMouse,
					termbox.InputEsc,
					termbox.InputAlt,
				}
				inputmode++
				if inputmode >= len(chmap) {
					inputmode = 0
				}
				termbox.SetInputMode(chmap[inputmode])
			}
			if ev.Key == termbox.KeyCtrlX {
				ctrlxpressed = true
			} else {
				ctrlxpressed = false
			}

			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			draw_keyboard()
			dispatch_press(&ev)
			pretty_print_press(&ev)
			termbox.Flush()
		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			draw_keyboard()
			pretty_print_resize(&ev)
			termbox.Flush()
		case termbox.EventMouse:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			draw_keyboard()
			pretty_print_mouse(&ev)
			termbox.Flush()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
