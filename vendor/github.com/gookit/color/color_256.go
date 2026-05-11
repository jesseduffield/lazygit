package color

import (
	"fmt"
	"strconv"
	"strings"
)

/*
from wikipedia, 256 color:
   ESC[ … 38;5;<n> … m选择前景色
   ESC[ … 48;5;<n> … m选择背景色
     0-  7：标准颜色（同 ESC[30–37m）
     8- 15：高强度颜色（同 ESC[90–97m）
    16-231：6 × 6 × 6 立方（216色）: 16 + 36 × r + 6 × g + b (0 ≤ r, g, b ≤ 5)
   232-255：从黑到白的24阶灰度色
*/

// tpl for 8 bit 256 color(`2^8`)
//
// format:
//
//		ESC[ … 38;5;<n> … m // 选择前景色
//	 ESC[ … 48;5;<n> … m // 选择背景色
//
// example:
//
//	fg "\x1b[38;5;242m"
//	bg "\x1b[48;5;208m"
//	both "\x1b[38;5;242;48;5;208m"
//
// links:
//
//	https://zh.wikipedia.org/wiki/ANSI%E8%BD%AC%E4%B9%89%E5%BA%8F%E5%88%97#8位
const (
	TplFg256 = "38;5;%d"
	TplBg256 = "48;5;%d"
	Fg256Pfx = "38;5;"
	Bg256Pfx = "48;5;"
)

/*************************************************************
 * 8bit(256) Color: Bit8Color Color256
 *************************************************************/

// Color256 256 color (8 bit), uint8 range at 0 - 255.
// Support 256 color on windows CMD, PowerShell
//
// 颜色值使用10进制和16进制都可 0x98 = 152
//
// The color consists of two uint8:
//
//	0: color value
//	1: color type; Fg=0, Bg=1, >1: unset value
//
// example:
//
//	fg color: [152, 0]
//	bg color: [152, 1]
//
// lint warn - Name starts with package name
type Color256 [2]uint8
type Bit8Color = Color256 // alias

var emptyC256 = Color256{1: 99}

// Bit8 create a color256
func Bit8(val uint8, isBg ...bool) Color256 { return C256(val, isBg...) }

// C256 create a color256
func C256(val uint8, isBg ...bool) Color256 {
	bc := Color256{val}

	// mark is bg color
	if len(isBg) > 0 && isBg[0] {
		bc[1] = AsBg
	}

	return bc
}

// Set terminal by 256 color code
func (c Color256) Set() error { return SetTerminal(c.String()) }

// Reset terminal. alias of the ResetTerminal()
func (c Color256) Reset() error { return ResetTerminal() }

// Print print message
func (c Color256) Print(a ...any) {
	doPrintV2(c.String(), fmt.Sprint(a...))
}

// Printf format and print message
func (c Color256) Printf(format string, a ...any) {
	doPrintV2(c.String(), fmt.Sprintf(format, a...))
}

// Println print message with newline
func (c Color256) Println(a ...any) {
	doPrintlnV2(c.String(), a)
}

// Sprint returns rendered message
func (c Color256) Sprint(a ...any) string { return RenderCode(c.String(), a...) }

// Sprintf returns format and rendered message
func (c Color256) Sprintf(format string, a ...any) string {
	return RenderString(c.String(), fmt.Sprintf(format, a...))
}

// C16 convert color-256 to 16 color.
func (c Color256) C16() Color { return c.Basic() }

// Basic convert color-256 to basic 16 color.
func (c Color256) Basic() Color { return Color(c[0]) /* TODO */ }

// RGB convert color-256 to RGB color.
func (c Color256) RGB() RGBColor {
	return RGBFromSlice(C256ToRgb(c[0]), c[1] == AsBg)
}

// RGBColor convert color-256 to RGB color.
func (c Color256) RGBColor() RGBColor { return c.RGB() }

// Value return color value
func (c Color256) Value() uint8 { return c[0] }

// Code convert to color code string. eg: "12"
func (c Color256) Code() string { return strconv.Itoa(int(c[0])) }

// FullCode convert to color code string with prefix. eg: "38;5;12"
func (c Color256) FullCode() string { return c.String() }

// String convert to color code string with prefix. eg: "38;5;12"
func (c Color256) String() string {
	if c[1] == AsFg { // 0 is Fg
		return Fg256Pfx + strconv.Itoa(int(c[0]))
	}

	if c[1] == AsBg { // 1 is Bg
		return Bg256Pfx + strconv.Itoa(int(c[0]))
	}
	return "" // empty
}

// IsFg color
func (c Color256) IsFg() bool { return c[1] == AsFg }

// ToFg 256 color
func (c Color256) ToFg() Color256 {
	c[1] = AsFg
	return c
}

// IsBg color
func (c Color256) IsBg() bool { return c[1] == AsBg }

// ToBg 256 color
func (c Color256) ToBg() Color256 {
	c[1] = AsBg
	return c
}

// IsEmpty value
func (c Color256) IsEmpty() bool { return c[1] > 1 }

/*************************************************************
 * 8bit(256) Style
 *************************************************************/

// Style256 definition
//
// 前/背景色
// 都是由两位uint8组成, 第一位是色彩值；
// 第二位与 Bit8Color 不一样的是，在这里表示是否设置了值 0 未设置 !=0 已设置
type Style256 struct {
	// Name of the style
	Name string
	// color options of the style
	opts Opts
	// fg and bg color
	fg, bg Color256
}

// S256 create a color256 style
//
// Usage:
//
//	s := color.S256()
//	s := color.S256(132) // fg
//	s := color.S256(132, 203) // fg and bg
func S256(fgAndBg ...uint8) *Style256 {
	s := &Style256{}
	vl := len(fgAndBg)
	if vl > 0 { // with fg
		s.fg = Color256{fgAndBg[0], 1}

		if vl > 1 { // and with bg
			s.bg = Color256{fgAndBg[1], 1}
		}
	}

	return s
}

// Set fg and bg color value, can also with color options
func (s *Style256) Set(fgVal, bgVal uint8, opts ...Color) *Style256 {
	s.fg = Color256{fgVal, 1}
	s.bg = Color256{bgVal, 1}
	s.opts.Add(opts...)
	return s
}

// SetBg set bg color value
func (s *Style256) SetBg(bgVal uint8) *Style256 {
	s.bg = Color256{bgVal, 1}
	return s
}

// SetFg set fg color value
func (s *Style256) SetFg(fgVal uint8) *Style256 {
	s.fg = Color256{fgVal, 1}
	return s
}

// SetOpts set options
func (s *Style256) SetOpts(opts Opts) *Style256 {
	s.opts = opts
	return s
}

// AddOpts add options
func (s *Style256) AddOpts(opts ...Color) *Style256 {
	s.opts.Add(opts...)
	return s
}

// Print message
func (s *Style256) Print(a ...any) {
	doPrintV2(s.String(), fmt.Sprint(a...))
}

// Printf format and print message
func (s *Style256) Printf(format string, a ...any) {
	doPrintV2(s.String(), fmt.Sprintf(format, a...))
}

// Println print message with newline
func (s *Style256) Println(a ...any) {
	doPrintlnV2(s.String(), a)
}

// Sprint returns rendered message
func (s *Style256) Sprint(a ...any) string { return RenderCode(s.Code(), a...) }

// Sprintf returns format and rendered message
func (s *Style256) Sprintf(format string, a ...any) string {
	return RenderString(s.Code(), fmt.Sprintf(format, a...))
}

// Code convert to color code string
func (s *Style256) Code() string { return s.String() }

// String convert to color code string
func (s *Style256) String() string {
	var ss []string
	if s.fg[1] > 0 {
		ss = append(ss, Fg256Pfx+strconv.FormatInt(int64(s.fg[0]), 10))
	}

	if s.bg[1] > 0 {
		ss = append(ss, Bg256Pfx+strconv.FormatInt(int64(s.bg[0]), 10))
	}

	if s.opts.IsValid() {
		ss = append(ss, s.opts.String())
	}
	return strings.Join(ss, ";")
}
