package color

import (
	"fmt"
	"strconv"
	"strings"
)

// 24 bit RGB color
// RGB:
//
//	R 0-255 G 0-255 B 0-255
//	R 00-FF G 00-FF B 00-FF (16进制)
//
// Format:
//
//	ESC[ … 38;2;<r>;<g>;<b> … m // Select RGB foreground color
//	ESC[ … 48;2;<r>;<g>;<b> … m // Choose RGB background color
//
// links:
//
//	https://zh.wikipedia.org/wiki/ANSI%E8%BD%AC%E4%B9%89%E5%BA%8F%E5%88%97#24位
//
// example:
//
//	fg: \x1b[38;2;30;144;255mMESSAGE\x1b[0m
//	bg: \x1b[48;2;30;144;255mMESSAGE\x1b[0m
//	both: \x1b[38;2;233;90;203;48;2;30;144;255mMESSAGE\x1b[0m
const (
	TplFgRGB = "38;2;%d;%d;%d"
	TplBgRGB = "48;2;%d;%d;%d"
	FgRGBPfx = "38;2;"
	BgRGBPfx = "48;2;"
)

// mark color is fg or bg.
const (
	AsFg uint8 = iota
	AsBg
)

/*************************************************************
 * RGB Color(Bit24Color, TrueColor)
 *************************************************************/

// RGBColor definition.
// Support RGB color on Windows CMD, PowerShell
//
// The first to third digits represent the color value.
// The last digit represents the foreground(0), background(1), >1 is unset value
//
// Usage:
//
//	// 0, 1, 2 is R,G,B.
//	// 3rd: Fg=0, Bg=1, >1: unset value
//	RGBColor{30,144,255, 0}
//	RGBColor{30,144,255, 1}
type RGBColor [4]uint8

// create an empty RGBColor
var emptyRGBColor = RGBColor{3: 99}

// RGB color create.
//
// Usage:
//
//	c := RGB(30,144,255)
//	c := RGB(30,144,255, true)
//	c.Print("message")
func RGB(r, g, b uint8, isBg ...bool) RGBColor {
	rgb := RGBColor{r, g, b}
	if len(isBg) > 0 && isBg[0] {
		rgb[3] = AsBg
	}

	return rgb
}

// Rgb alias of the RGB()
func Rgb(r, g, b uint8, isBg ...bool) RGBColor { return RGB(r, g, b, isBg...) }

// Bit24 alias of the RGB()
func Bit24(r, g, b uint8, isBg ...bool) RGBColor { return RGB(r, g, b, isBg...) }

// RgbFromInt create instance from int r,g,b value
func RgbFromInt(r, g, b int, isBg ...bool) RGBColor { return RGB(uint8(r), uint8(g), uint8(b), isBg...) }

// RgbFromInts create instance from []int r,g,b value
func RgbFromInts(rgb []int, isBg ...bool) RGBColor {
	return RGB(uint8(rgb[0]), uint8(rgb[1]), uint8(rgb[2]), isBg...)
}

// HEX create RGB color from a HEX color string.
//
// Usage:
//
//	c := HEX("ccc") // rgb: [204 204 204]
//	c := HEX("aabbcc") // rgb: [170 187 204]
//	c := HEX("#aabbcc")
//	c := HEX("0xaabbcc")
//	c.Print("message")
func HEX(hex string, isBg ...bool) RGBColor {
	if rgb := HexToRgb(hex); len(rgb) > 0 {
		return RGB(uint8(rgb[0]), uint8(rgb[1]), uint8(rgb[2]), isBg...)
	}

	// mark is empty
	return emptyRGBColor
}

// Hex alias of the HEX()
func Hex(hex string, isBg ...bool) RGBColor { return HEX(hex, isBg...) }

// RGBFromHEX quick RGBColor from hex string, alias of HEX()
func RGBFromHEX(hex string, isBg ...bool) RGBColor { return HEX(hex, isBg...) }

// HSL create RGB color from a hsl value.
// more see HslToRgb()
func HSL(h, s, l float64, isBg ...bool) RGBColor {
	rgb := HslToRgb(h, s, l)
	return RGB(rgb[0], rgb[1], rgb[2], isBg...)
}

// Hsl alias of the HSL()
func Hsl(h, s, l float64, isBg ...bool) RGBColor { return HSL(h, s, l, isBg...) }

// HSLInt create RGB color from a hsl int value.
// more see HslIntToRgb()
func HSLInt(h, s, l int, isBg ...bool) RGBColor {
	rgb := HslIntToRgb(h, s, l)
	return RGB(rgb[0], rgb[1], rgb[2], isBg...)
}

// HslInt alias of the HSLInt()
func HslInt(h, s, l int, isBg ...bool) RGBColor { return HSLInt(h, s, l, isBg...) }

// RGBFromSlice quick RGBColor from slice[3]
func RGBFromSlice(rgb []uint8, isBg ...bool) RGBColor { return RGB(rgb[0], rgb[1], rgb[2], isBg...) }

// RGBFromString create RGB color from a string.
// Support use color name in the {namedRgbMap}
//
// Usage:
//
//	c := RGBFromString("170,187,204")
//	c.Print("message")
//
//	c := RGBFromString("brown")
//	c.Print("message with color brown")
func RGBFromString(rgb string, isBg ...bool) RGBColor {
	// use color name in the {namedRgbMap}
	if rgbVal, ok := namedRgbMap[rgb]; ok {
		rgb = rgbVal
	}

	// use rgb string.
	ss := stringToArr(rgb, ",")
	if len(ss) != 3 {
		return emptyRGBColor
	}

	var ar [3]uint8
	for i, val := range ss {
		iv, err := strconv.Atoi(val)
		if err != nil || !isValidUint8(iv) {
			return emptyRGBColor
		}

		ar[i] = uint8(iv)
	}

	return RGB(ar[0], ar[1], ar[2], isBg...)
}

// Set terminal by rgb/true color code
func (c RGBColor) Set() error { return SetTerminal(c.String()) }

// Reset terminal. alias of the ResetTerminal()
func (c RGBColor) Reset() error { return ResetTerminal() }

// Print print message
func (c RGBColor) Print(a ...any) {
	doPrintV2(c.String(), fmt.Sprint(a...))
}

// Printf format and print message
func (c RGBColor) Printf(format string, a ...any) {
	doPrintV2(c.String(), fmt.Sprintf(format, a...))
}

// Println print message with newline
func (c RGBColor) Println(a ...any) {
	doPrintlnV2(c.String(), a)
}

// Sprint returns rendered message
func (c RGBColor) Sprint(a ...any) string { return RenderCode(c.String(), a...) }

// Sprintf returns format and rendered message
func (c RGBColor) Sprintf(format string, a ...any) string {
	return RenderString(c.String(), fmt.Sprintf(format, a...))
}

// Values to RGB values
func (c RGBColor) Values() []int {
	return []int{int(c[0]), int(c[1]), int(c[2])}
}

// Code to color code string without prefix. eg: "204;123;56"
func (c RGBColor) Code() string { return fmt.Sprintf("%d;%d;%d", c[0], c[1], c[2]) }

// Hex color rgb to hex string. as in "ff0080".
func (c RGBColor) Hex() string { return fmt.Sprintf("%02x%02x%02x", c[0], c[1], c[2]) }

// RgbString to color code string without prefix. eg: "204,123,56"
func (c RGBColor) RgbString() string {
	return fmt.Sprintf("%d,%d,%d", c[0], c[1], c[2])
}

// FullCode to color code string with prefix
func (c RGBColor) FullCode() string { return c.String() }

// String to color code string with prefix. eg: "38;2;204;123;56"
func (c RGBColor) String() string {
	if c[3] == AsFg {
		return fmt.Sprintf(TplFgRGB, c[0], c[1], c[2])
	}

	if c[3] == AsBg {
		return fmt.Sprintf(TplBgRGB, c[0], c[1], c[2])
	}

	// c[3] > 1 is empty
	return ""
}

// ToBg convert to background color
func (c RGBColor) ToBg() RGBColor {
	c[3] = AsBg
	return c
}

// ToFg convert to foreground color
func (c RGBColor) ToFg() RGBColor {
	c[3] = AsFg
	return c
}

// IsEmpty value
func (c RGBColor) IsEmpty() bool { return c[3] > AsBg }

// IsValid value
// func (c RGBColor) IsValid() bool {
// 	return c[3] <= AsBg
// }

// C256 returns the closest approximate 256 (8 bit) color
func (c RGBColor) C256() Color256 {
	return C256(RgbTo256(c[0], c[1], c[2]), c[3] == AsBg)
}

// Basic returns the closest approximate 16 (4 bit) color
func (c RGBColor) Basic() Color {
	// return Color(RgbToAnsi(c[0], c[1], c[2], c[3] == AsBg))
	return Color(Rgb2basic(c[0], c[1], c[2], c[3] == AsBg))
}

// Color returns the closest approximate 16 (4 bit) color
func (c RGBColor) Color() Color { return c.Basic() }

// C16 returns the closest approximate 16 (4 bit) color
func (c RGBColor) C16() Color { return c.Basic() }

/*************************************************************
 * RGB Style
 *************************************************************/

// RGBStyle supports set foreground and background color
//
// All are composed of 4 digits uint8, the first three digits are the color value;
// The last bit is different from RGBColor, here it indicates whether the value is set.
//
//	1    Has been set
//	^1   Not set
type RGBStyle struct {
	// Name of the style
	Name string
	// color options of the style
	opts Opts
	// fg and bg color
	fg, bg RGBColor
}

// NewRGBStyle create a RGBStyle.
func NewRGBStyle(fg RGBColor, bg ...RGBColor) *RGBStyle {
	s := &RGBStyle{}
	if len(bg) > 0 {
		s.SetBg(bg[0])
	}

	return s.SetFg(fg)
}

// HEXStyle create a RGBStyle from HEX color string.
//
// Usage:
//
//	s := HEXStyle("aabbcc", "eee")
//	s.Print("message")
func HEXStyle(fg string, bg ...string) *RGBStyle {
	s := &RGBStyle{}
	if len(bg) > 0 {
		s.SetBg(HEX(bg[0]))
	}

	if len(fg) > 0 {
		s.SetFg(HEX(fg))
	}
	return s
}

// RGBStyleFromString create a RGBStyle from color value string.
//
// Usage:
//
//	s := RGBStyleFromString("170,187,204", "70,87,4")
//	s.Print("message")
func RGBStyleFromString(fg string, bg ...string) *RGBStyle {
	s := &RGBStyle{}
	if len(bg) > 0 {
		s.SetBg(RGBFromString(bg[0]))
	}

	return s.SetFg(RGBFromString(fg))
}

// Set fg and bg color, can also with color options
func (s *RGBStyle) Set(fg, bg RGBColor, opts ...Color) *RGBStyle {
	return s.SetFg(fg).SetBg(bg).SetOpts(opts)
}

// SetFg set fg color
func (s *RGBStyle) SetFg(fg RGBColor) *RGBStyle {
	fg[3] = 1 // add fixed value, mark is valid
	s.fg = fg
	return s
}

// SetBg set bg color
func (s *RGBStyle) SetBg(bg RGBColor) *RGBStyle {
	bg[3] = 1 // add fixed value, mark is valid
	s.bg = bg
	return s
}

// SetOpts set color options
func (s *RGBStyle) SetOpts(opts Opts) *RGBStyle {
	s.opts = opts
	return s
}

// AddOpts add options
func (s *RGBStyle) AddOpts(opts ...Color) *RGBStyle {
	s.opts.Add(opts...)
	return s
}

// Print print message
func (s *RGBStyle) Print(a ...any) {
	doPrintV2(s.String(), fmt.Sprint(a...))
}

// Printf format and print message
func (s *RGBStyle) Printf(format string, a ...any) {
	doPrintV2(s.String(), fmt.Sprintf(format, a...))
}

// Println print message with newline
func (s *RGBStyle) Println(a ...any) {
	doPrintlnV2(s.String(), a)
}

// Sprint returns rendered message
func (s *RGBStyle) Sprint(a ...any) string { return RenderCode(s.String(), a...) }

// Sprintf returns format and rendered message
func (s *RGBStyle) Sprintf(format string, a ...any) string {
	return RenderString(s.String(), fmt.Sprintf(format, a...))
}

// Code convert to color code string
func (s *RGBStyle) Code() string { return s.String() }

// FullCode convert to color code string
func (s *RGBStyle) FullCode() string { return s.String() }

// String convert to color code string
func (s *RGBStyle) String() string {
	var ss []string
	// last value ensure is enable.
	if s.fg[3] == 1 {
		ss = append(ss, fmt.Sprintf(TplFgRGB, s.fg[0], s.fg[1], s.fg[2]))
	}

	if s.bg[3] == 1 {
		ss = append(ss, fmt.Sprintf(TplBgRGB, s.bg[0], s.bg[1], s.bg[2]))
	}

	if s.opts.IsValid() {
		ss = append(ss, s.opts.String())
	}

	return strings.Join(ss, ";")
}

// IsEmpty style
func (s *RGBStyle) IsEmpty() bool { return s.fg[3] != 1 && s.bg[3] != 1 }
