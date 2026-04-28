package color

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// values from https://github.com/go-terminfo/terminfo
// var (
// RgbaBlack    = image_color.RGBA{0, 0, 0, 255}
// Red       = color.RGBA{205, 0, 0, 255}
// Green     = color.RGBA{0, 205, 0, 255}
// Orange    = color.RGBA{205, 205, 0, 255}
// Blue      = color.RGBA{0, 0, 238, 255}
// Magenta   = color.RGBA{205, 0, 205, 255}
// Cyan      = color.RGBA{0, 205, 205, 255}
// LightGrey = color.RGBA{229, 229, 229, 255}
//
// DarkGrey     = color.RGBA{127, 127, 127, 255}
// LightRed     = color.RGBA{255, 0, 0, 255}
// LightGreen   = color.RGBA{0, 255, 0, 255}
// Yellow       = color.RGBA{255, 255, 0, 255}
// LightBlue    = color.RGBA{92, 92, 255, 255}
// LightMagenta = color.RGBA{255, 0, 255, 255}
// LightCyan    = color.RGBA{0, 255, 255, 255}
// White        = color.RGBA{255, 255, 255, 255}
// )

var (
	// ---------- basic(16) <=> 256 color convert ----------
	basicTo256Map = map[uint8]uint8{
		30: 0,   // black 	000000
		31: 160, // red 	c51e14
		32: 34,  // green 	1dc121
		33: 184, // yellow 	c7c329
		34: 20,  // blue 	0a2fc4
		35: 170, // magenta c839c5
		36: 44,  // cyan 	20c5c6
		37: 188, // white 	c7c7c7
		90: 59,  // lightBlack 		686868
		91: 203, // lightRed 		fd6f6b
		92: 83,  // lightGreen 		67f86f
		93: 227, // lightYellow 	fffa72
		94: 69,  // lightBlue 		6a76fb
		95: 213, // lightMagenta 	fd7cfc
		96: 87,  // lightCyan 		68fdfe
		97: 15,  // lightWhite 		ffffff
	}

	// ---------- basic(16) <=> RGB color convert ----------
	// refer from Hyper app
	// Tip: only keep foreground color, background color need convert to foreground color for convert to RGB
	basic2hexMap = map[uint8]string{
		30: "000000", // black
		31: "c51e14", // red
		32: "1dc121", // green
		33: "c7c329", // yellow
		34: "0a2fc4", // blue
		35: "c839c5", // magenta
		36: "20c5c6", // cyan
		37: "c7c7c7", // white
		// - don't add bg color, convert to fg color for convert to RGB
		// 40:  "000000", // black
		// 41:  "c51e14", // red
		// 42:  "1dc121", // green
		// 43:  "c7c329", // yellow
		// 44:  "0a2fc4", // blue
		// 45:  "c839c5", // magenta
		// 46:  "20c5c6", // cyan
		// 47:  "c7c7c7", // white
		90: "686868", // lightBlack/darkGray
		91: "fd6f6b", // lightRed
		92: "67f86f", // lightGreen
		93: "fffa72", // lightYellow
		94: "6a76fb", // lightBlue
		95: "fd7cfc", // lightMagenta
		96: "68fdfe", // lightCyan
		97: "ffffff", // lightWhite
		// - don't add bg color
		// 100: "686868", // lightBlack/darkGray
		// 101: "fd6f6b", // lightRed
		// 102: "67f86f", // lightGreen
		// 103: "fffa72", // lightYellow
		// 104: "6a76fb", // lightBlue
		// 105: "fd7cfc", // lightMagenta
		// 106: "68fdfe", // lightCyan
		// 107: "ffffff", // lightWhite
	}
	// will convert data from basic2hexMap
	hex2basicMap = initHex2basicMap()

	// ---------- 256 <=> RGB color convert ----------
	// adapted from https://gist.github.com/MicahElliott/719710

	c256ToHexMap = init256ToHexMap()

	// rgb to 256 color look-up table
	// RGB hex => 256 code
	hexTo256Table = map[string]uint8{
		// Primary 3-bit (8 colors). Unique representation!
		"000000": 0,
		"800000": 1,
		"008000": 2,
		"808000": 3,
		"000080": 4,
		"800080": 5,
		"008080": 6,
		"c0c0c0": 7,

		// Equivalent "bright" versions of original 8 colors.
		"808080": 8,
		"ff0000": 9,
		"00ff00": 10,
		"ffff00": 11,
		"0000ff": 12,
		"ff00ff": 13,
		"00ffff": 14,
		"ffffff": 15,

		// values commented out below are duplicates from the prior sections

		// Strictly ascending.
		// "000000": 16,
		"000001": 16, // up: avoid key conflicts, value + 1
		"00005f": 17,
		"000087": 18,
		"0000af": 19,
		"0000d7": 20,
		// "0000ff": 21,
		"0000fe": 21, // up: avoid key conflicts, value - 1
		"005f00": 22,
		"005f5f": 23,
		"005f87": 24,
		"005faf": 25,
		"005fd7": 26,
		"005fff": 27,
		"008700": 28,
		"00875f": 29,
		"008787": 30,
		"0087af": 31,
		"0087d7": 32,
		"0087ff": 33,
		"00af00": 34,
		"00af5f": 35,
		"00af87": 36,
		"00afaf": 37,
		"00afd7": 38,
		"00afff": 39,
		"00d700": 40,
		"00d75f": 41,
		"00d787": 42,
		"00d7af": 43,
		"00d7d7": 44,
		"00d7ff": 45,
		// "00ff00": 46,
		"00ff01": 46, // up: avoid key conflicts, value + 1
		"00ff5f": 47,
		"00ff87": 48,
		"00ffaf": 49,
		"00ffd7": 50,
		// "00ffff": 51,
		"00fffe": 51, // up: avoid key conflicts, value - 1
		"5f0000": 52,
		"5f005f": 53,
		"5f0087": 54,
		"5f00af": 55,
		"5f00d7": 56,
		"5f00ff": 57,
		"5f5f00": 58,
		"5f5f5f": 59,
		"5f5f87": 60,
		"5f5faf": 61,
		"5f5fd7": 62,
		"5f5fff": 63,
		"5f8700": 64,
		"5f875f": 65,
		"5f8787": 66,
		"5f87af": 67,
		"5f87d7": 68,
		"5f87ff": 69,
		"5faf00": 70,
		"5faf5f": 71,
		"5faf87": 72,
		"5fafaf": 73,
		"5fafd7": 74,
		"5fafff": 75,
		"5fd700": 76,
		"5fd75f": 77,
		"5fd787": 78,
		"5fd7af": 79,
		"5fd7d7": 80,
		"5fd7ff": 81,
		"5fff00": 82,
		"5fff5f": 83,
		"5fff87": 84,
		"5fffaf": 85,
		"5fffd7": 86,
		"5fffff": 87,
		"870000": 88,
		"87005f": 89,
		"870087": 90,
		"8700af": 91,
		"8700d7": 92,
		"8700ff": 93,
		"875f00": 94,
		"875f5f": 95,
		"875f87": 96,
		"875faf": 97,
		"875fd7": 98,
		"875fff": 99,
		"878700": 100,
		"87875f": 101,
		"878787": 102,
		"8787af": 103,
		"8787d7": 104,
		"8787ff": 105,
		"87af00": 106,
		"87af5f": 107,
		"87af87": 108,
		"87afaf": 109,
		"87afd7": 110,
		"87afff": 111,
		"87d700": 112,
		"87d75f": 113,
		"87d787": 114,
		"87d7af": 115,
		"87d7d7": 116,
		"87d7ff": 117,
		"87ff00": 118,
		"87ff5f": 119,
		"87ff87": 120,
		"87ffaf": 121,
		"87ffd7": 122,
		"87ffff": 123,
		"af0000": 124,
		"af005f": 125,
		"af0087": 126,
		"af00af": 127,
		"af00d7": 128,
		"af00ff": 129,
		"af5f00": 130,
		"af5f5f": 131,
		"af5f87": 132,
		"af5faf": 133,
		"af5fd7": 134,
		"af5fff": 135,
		"af8700": 136,
		"af875f": 137,
		"af8787": 138,
		"af87af": 139,
		"af87d7": 140,
		"af87ff": 141,
		"afaf00": 142,
		"afaf5f": 143,
		"afaf87": 144,
		"afafaf": 145,
		"afafd7": 146,
		"afafff": 147,
		"afd700": 148,
		"afd75f": 149,
		"afd787": 150,
		"afd7af": 151,
		"afd7d7": 152,
		"afd7ff": 153,
		"afff00": 154,
		"afff5f": 155,
		"afff87": 156,
		"afffaf": 157,
		"afffd7": 158,
		"afffff": 159,
		"d70000": 160,
		"d7005f": 161,
		"d70087": 162,
		"d700af": 163,
		"d700d7": 164,
		"d700ff": 165,
		"d75f00": 166,
		"d75f5f": 167,
		"d75f87": 168,
		"d75faf": 169,
		"d75fd7": 170,
		"d75fff": 171,
		"d78700": 172,
		"d7875f": 173,
		"d78787": 174,
		"d787af": 175,
		"d787d7": 176,
		"d787ff": 177,
		"d7af00": 178,
		"d7af5f": 179,
		"d7af87": 180,
		"d7afaf": 181,
		"d7afd7": 182,
		"d7afff": 183,
		"d7d700": 184,
		"d7d75f": 185,
		"d7d787": 186,
		"d7d7af": 187,
		"d7d7d7": 188,
		"d7d7ff": 189,
		"d7ff00": 190,
		"d7ff5f": 191,
		"d7ff87": 192,
		"d7ffaf": 193,
		"d7ffd7": 194,
		"d7ffff": 195,
		// "ff0000": 196,
		"ff0001": 196, // up: avoid key conflicts, value + 1
		"ff005f": 197,
		"ff0087": 198,
		"ff00af": 199,
		"ff00d7": 200,
		// "ff00ff": 201,
		"ff00fe": 201, // up: avoid key conflicts, value - 1
		"ff5f00": 202,
		"ff5f5f": 203,
		"ff5f87": 204,
		"ff5faf": 205,
		"ff5fd7": 206,
		"ff5fff": 207,
		"ff8700": 208,
		"ff875f": 209,
		"ff8787": 210,
		"ff87af": 211,
		"ff87d7": 212,
		"ff87ff": 213,
		"ffaf00": 214,
		"ffaf5f": 215,
		"ffaf87": 216,
		"ffafaf": 217,
		"ffafd7": 218,
		"ffafff": 219,
		"ffd700": 220,
		"ffd75f": 221,
		"ffd787": 222,
		"ffd7af": 223,
		"ffd7d7": 224,
		"ffd7ff": 225,
		// "ffff00": 226,
		"ffff01": 226, // up: avoid key conflicts, value + 1
		"ffff5f": 227,
		"ffff87": 228,
		"ffffaf": 229,
		"ffffd7": 230,
		// "ffffff": 231,
		"fffffe": 231, // up: avoid key conflicts, value - 1

		// Gray-scale range.
		"080808": 232,
		"121212": 233,
		"1c1c1c": 234,
		"262626": 235,
		"303030": 236,
		"3a3a3a": 237,
		"444444": 238,
		"4e4e4e": 239,
		"585858": 240,
		"626262": 241,
		"6c6c6c": 242,
		"767676": 243,
		// "808080": 244,
		"808081": 244, // up: avoid key conflicts, value + 1
		"8a8a8a": 245,
		"949494": 246,
		"9e9e9e": 247,
		"a8a8a8": 248,
		"b2b2b2": 249,
		"bcbcbc": 250,
		"c6c6c6": 251,
		"d0d0d0": 252,
		"dadada": 253,
		"e4e4e4": 254,
		"eeeeee": 255,
	}

	incs = []uint8{0x00, 0x5f, 0x87, 0xaf, 0xd7, 0xff}
)

func initHex2basicMap() map[string]uint8 {
	h2b := make(map[string]uint8, len(basic2hexMap))
	// ini data map
	for u, s := range basic2hexMap {
		h2b[s] = u
	}
	return h2b
}

func init256ToHexMap() map[uint8]string {
	c256toh := make(map[uint8]string, len(hexTo256Table))
	// ini data map
	for hex, c256 := range hexTo256Table {
		c256toh[c256] = hex
	}
	return c256toh
}

// RgbTo256Table mapping data
func RgbTo256Table() map[string]uint8 {
	return hexTo256Table
}

// Colors2code convert colors to code. return like "32;45;3"
func Colors2code(colors ...Color) string {
	if len(colors) == 0 {
		return ""
	}

	var codes []string
	for _, color := range colors {
		codes = append(codes, color.String())
	}

	return strings.Join(codes, ";")
}

/*************************************************************
 * region HEX <=> RGB
 * HEX code <=> RGB/True color code
 *************************************************************/

// Hex2rgb alias of the HexToRgb()
func Hex2rgb(hex string) []int { return HexToRgb(hex) }

// HexToRGB alias of the HexToRgb()
func HexToRGB(hex string) []int { return HexToRgb(hex) }

// HexToRgb convert hex color string to RGB numbers
//
// Usage:
//
//	rgb := HexToRgb("ccc") // rgb: [204 204 204]
//	rgb := HexToRgb("aabbcc") // rgb: [170 187 204]
//	rgb := HexToRgb("#aabbcc") // rgb: [170 187 204]
//	rgb := HexToRgb("0xad99c0") // rgb: [170 187 204]
func HexToRgb(hex string) (rgb []int) {
	hex = strings.TrimSpace(hex)
	if hex == "" {
		return
	}

	// like from css. eg "#ccc" "#ad99c0"
	if hex[0] == '#' {
		hex = hex[1:]
	}

	hex = strings.ToLower(hex)
	switch len(hex) {
	case 3: // "ccc"
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	case 8: // "0xad99c0"
		hex = strings.TrimPrefix(hex, "0x")
	}

	// recheck
	if len(hex) != 6 {
		return
	}

	// convert string to int64
	if i64, err := strconv.ParseInt(hex, 16, 32); err == nil {
		color := int(i64)
		// parse int
		rgb = make([]int, 3)
		rgb[0] = color >> 16
		rgb[1] = (color & 0x00FF00) >> 8
		rgb[2] = color & 0x0000FF
	}
	return
}

// Rgb2hex alias of the RgbToHex()
func Rgb2hex(rgb []int) string { return RgbToHex(rgb) }

// RgbToHex convert RGB-code to hex-code
//
// Usage:
//
//	hex := RgbToHex([]int{170, 187, 204}) // hex: "aabbcc"
func RgbToHex(rgb []int) string {
	hexNodes := make([]string, 0, len(rgb))

	for _, v := range rgb {
		hexNodes = append(hexNodes, strconv.FormatInt(int64(v), 16))
	}
	return strings.Join(hexNodes, "")
}

/*************************************************************
 * region 4bit(16) <=> RGB
 * 4bit(16) color <=> RGB/True color
 *************************************************************/

// BasicToHex convert basic color to hex string.
func BasicToHex(val uint8) string {
	val = Bg2Fg(val)
	return basic2hexMap[val]
}

// Basic2hex convert basic color to hex string.
func Basic2hex(val uint8) string {
	return BasicToHex(val)
}

// Hex2basic convert hex string to basic color code.
func Hex2basic(hex string, asBg ...bool) uint8 {
	val := hex2basicMap[hex]

	if len(asBg) > 0 && asBg[0] {
		return Fg2Bg(val)
	}
	return val
}

// Rgb2basic alias of the RgbToAnsi()
func Rgb2basic(r, g, b uint8, isBg bool) uint8 {
	// is basic color, direct use static map data.
	hex := RgbToHex([]int{int(r), int(g), int(b)})
	if val, ok := hex2basicMap[hex]; ok {
		if isBg {
			return val + 10
		}
		return val
	}

	return RgbToAnsi(r, g, b, isBg)
}

// Rgb2ansi convert RGB-code to 16-code, alias of the RgbToAnsi()
func Rgb2ansi(r, g, b uint8, isBg bool) uint8 {
	return RgbToAnsi(r, g, b, isBg)
}

// RgbToAnsi convert RGB-code to 16-code
// refer https://github.com/radareorg/radare2/blob/master/libr/cons/rgb.c#L249-L271
func RgbToAnsi(r, g, b uint8, isBg bool) uint8 {
	var bright, c, k uint8
	base := compareVal(isBg, BgBase, FgBase)

	// eco bright-specific
	if r == 0x80 && g == 0x80 && b == 0x80 { // 0x80=128
		bright = 53
	} else if r == 0xff || g == 0xff || b == 0xff { // 0xff=255
		bright = 60
	} // else bright = 0

	if r == g && g == b {
		// 0x7f=127
		// r = (r > 0x7f) ? 1 : 0;
		r = compareVal(r > 0x7f, 1, 0)
		g = compareVal(g > 0x7f, 1, 0)
		b = compareVal(b > 0x7f, 1, 0)
	} else {
		k = (r + g + b) / 3

		// r = (r >= k) ? 1 : 0;
		r = compareVal(r >= k, 1, 0)
		g = compareVal(g >= k, 1, 0)
		b = compareVal(b >= k, 1, 0)
	}

	// c = (r ? 1 : 0) + (g ? (b ? 6 : 2) : (b ? 4 : 0))
	c = compareVal(r > 0, 1, 0)

	if g > 0 {
		c += compareVal(b > 0, 6, 2)
	} else {
		c += compareVal(b > 0, 4, 0)
	}
	return base + bright + c
}

/*************************************************************
 * 8bit(256) color <=> RGB/True color
 *************************************************************/

// Rgb2short convert RGB-code to 256-code
func Rgb2short(r, g, b uint8) uint8 {
	return RgbTo256(r, g, b)
}

// RgbTo256 convert RGB-code to 256-code
func RgbTo256(r, g, b uint8) uint8 {
	res := make([]uint8, 3)
	for partI, part := range [3]uint8{r, g, b} {
		i := 0
		for i < len(incs)-1 {
			s, b := incs[i], incs[i+1] // smaller, bigger
			if s <= part && part <= b {
				s1 := math.Abs(float64(s) - float64(part))
				b1 := math.Abs(float64(b) - float64(part))
				var closest uint8
				if s1 < b1 {
					closest = s
				} else {
					closest = b
				}
				res[partI] = closest
				break
			}
			i++
		}
	}
	hex := fmt.Sprintf("%02x%02x%02x", res[0], res[1], res[2])
	equiv := hexTo256Table[hex]
	return equiv
}

// C256ToRgb convert an 256 color code to RGB numbers
func C256ToRgb(val uint8) (rgb []uint8) {
	hex := c256ToHexMap[val]
	// convert to rgb code
	rgbInts := Hex2rgb(hex)

	return []uint8{
		uint8(rgbInts[0]),
		uint8(rgbInts[1]),
		uint8(rgbInts[2]),
	}
}

// C256ToRgbV1 convert an 256 color code to RGB numbers
// refer https://github.com/torvalds/linux/commit/cec5b2a97a11ade56a701e83044d0a2a984c67b4
func C256ToRgbV1(val uint8) (rgb []uint8) {
	var r, g, b uint8
	if val < 8 { // Standard colours.
		// r = val&1 ? 0xaa : 0x00;
		r = compareVal(val&1 == 1, 0xaa, 0x00)
		g = compareVal(val&2 == 2, 0xaa, 0x00)
		b = compareVal(val&4 == 4, 0xaa, 0x00)
	} else if val < 16 {
		// r = val & 1 ? 0xff : 0x55;
		r = compareVal(val&1 == 1, 0xff, 0x55)
		g = compareVal(val&2 == 2, 0xff, 0x55)
		b = compareVal(val&4 == 4, 0xff, 0x55)
	} else if val < 232 { /* 6x6x6 colour cube. */
		r = (val - 16) / 36 * 85 / 2
		g = (val - 16) / 6 % 6 * 85 / 2
		b = (val - 16) % 6 * 85 / 2
	} else { /* Grayscale ramp. */
		nv := uint8(int(val)*10 - 2312)
		// set value
		r, g, b = nv, nv, nv
	}

	return []uint8{r, g, b}
}

/**************************************************************
 * region HSL <=> RGB color
 ************************************************************
 * h,s,l = Hue, Saturation, Lightness 色相、饱和度、亮度
 *
 * refers
 *  http://en.wikipedia.org/wiki/HSL_color_space
 *  https://www.w3.org/TR/css-color-3/#hsl-color
 *  https://stackoverflow.com/questions/2353211/hsl-to-rgb-color-conversion
 *	https://github.com/less/less.js/blob/master/packages/less/src/less/functions/color.js
 *  https://github.com/d3/d3-color/blob/v3.0.1/README.md#hsl
 *
 * examples:
 *  color: hsl(0, 100%, 50%)   // red
 *  color: hsl(120, 100%, 50%) // lime
 *  color: hsl(120, 100%, 25%) // dark green
 *  color: hsl(120, 100%, 75%) // light green
 *  color: hsl(120, 75%, 75%)  // pastel green, and so on
 */

// HslIntToRgb Converts an HSL color value to RGB
// Assumes h: 0-360, s: 0-100%, l: 0-100%
// returns r, g, and b in the set [0, 255].
//
// Usage:
//
//	HslIntToRgb(0, 100, 50) // red
//	HslIntToRgb(120, 100, 50) // lime
//	HslIntToRgb(120, 100, 25) // dark green
//	HslIntToRgb(120, 100, 75) // light green
func HslIntToRgb(h, s, l int) (rgb []uint8) {
	return HslToRgb(float64(h)/360, float64(s)/100, float64(l)/100)
}

// HslToRgb Converts an HSL color value to RGB. Conversion formula
// adapted from http://en.wikipedia.org/wiki/HSL_color_space.
// Assumes h, s, and l are contained in the set [0, 1]
// returns r, g, and b in the set [0, 255].
//
// Usage:
//
//	rgbVals := HslToRgb(0, 1, 0.5) // red
func HslToRgb(h, s, l float64) (rgb []uint8) {
	var r, g, b float64

	if s == 0 { // achromatic
		r, g, b = l, l, l
	} else {
		// q = l < 0.5 ? l * (1 + s) : l + s - l*s
		var q float64
		if l < 0.5 {
			q = l * (1.0 + s)
		} else {
			q = l + s - l*s
		}

		var p = 2.0*l - q

		r = hue2rgb(p, q, h+1.0/3.0)
		g = hue2rgb(p, q, h)
		b = hue2rgb(p, q, h-1.0/3.0)
	}

	// return []uint8{uint8(r * 255), uint8(g * 255), uint8(b * 255)}
	return []uint8{
		uint8(math.Round(r * 255)),
		uint8(math.Round(g * 255)),
		uint8(math.Round(b * 255)),
	}
}

var hue2rgb = func(p, q, t float64) float64 {
	if t < 0.0 {
		t += 1
	}
	if t > 1.0 {
		t -= 1
	}

	if t < 1.0/6.0 {
		return p + (q-p)*6.0*t
	}

	if t < 1.0/2.0 {
		return q
	}

	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6.0
	}
	return p
}

// RgbStrToHslInts convert rgb(r,g,b) string to HSL int values.
func RgbStrToHslInts(rgbStr string) []int {
	f64s := RgbStrToHsl(rgbStr)
	if f64s == nil {
		return nil
	}

	return []int{
		int(math.Round(f64s[0] * 360)),
		int(math.Round(f64s[1] * 100)),
		int(math.Round(f64s[2] * 100)),
	}
}

// RgbStrToHsl convert rgb(r,g,b) string to HSL
func RgbStrToHsl(rgbStr string) []float64 {
	if pos := strings.IndexByte(rgbStr, '('); pos > 0 {
		rgbStr = strings.TrimRight(rgbStr[pos+1:], "()")
	}

	rgbVals := strings.Split(rgbStr, ",")
	if len(rgbVals) != 3 {
		return nil
	}

	r, e1 := strconv.ParseInt(strings.TrimSpace(rgbVals[0]), 10, 0)
	if e1 != nil {
		return nil
	}
	g, e2 := strconv.ParseInt(strings.TrimSpace(rgbVals[1]), 10, 0)
	if e2 != nil {
		return nil
	}
	b, e3 := strconv.ParseInt(strings.TrimSpace(rgbVals[2]), 10, 0)
	if e3 != nil {
		return nil
	}
	return RgbToHsl(uint8(r), uint8(g), uint8(b))
}

// RgbToHslInt Converts an RGB color value to HSL. Conversion formula
// Assumes r, g, and b are contained in the set [0, 255] and
// returns [h,s,l] h: 0-360, s: 0-100%, l: 0-100%.
func RgbToHslInt(r, g, b uint8) []int {
	f64s := RgbToHsl(r, g, b)
	return []int{
		int(math.Round(f64s[0] * 360)),
		int(math.Round(f64s[1] * 100)),
		int(math.Round(f64s[2] * 100)),
	}
}

// RgbToHsl Converts an RGB color value to HSL. Conversion formula
//
// adapted from http://en.wikipedia.org/wiki/HSL_color_space.
//
// e.g: rgb(59, 130, 246) = hsl(217, 91%, 60%)
//
// Assumes r, g, and b are contained in the set [0, 255] and
// returns h, s, and l in the set [0, 1].
func RgbToHsl(r, g, b uint8) []float64 {
	// to float64
	fr, fg, fb := float64(r), float64(g), float64(b)
	// percentage
	pr, pg, pb := fr/255.0, fg/255.0, fb/255.0

	ps := []float64{pr, pg, pb}
	sort.Float64s(ps)

	min1, max1 := ps[0], ps[2]
	// max := math.Max(math.Max(pr, pg), pb)
	// min := math.Min(math.Min(pr, pg), pb)
	mid := (max1 + min1) / 2.0 // call Lightness

	h, s, l := mid, mid, mid
	// calc Saturation
	if max1 == min1 {
		h, s = 0, 0 // achromatic
	} else {
		var d = max1 - min1 // 计算色差
		// s = l > 0.5 ? d / (2 - max1 - min1) : d / (max1 + min1)
		s = compareF64(l > 0.5, d/(2.0-max1-min1), d/(max1+min1))

		// calc Hue
		switch max1 {
		case pr:
			// h = (g - b) / d + (g < b ? 6 : 0)
			h = (pg - pb) / d
			h += compareF64(g < b, 6, 0)
		case pg:
			h = (pb-pr)/d + 2
		case pb:
			h = (pr-pg)/d + 4
		}

		h /= 6
	}

	return []float64{h, s, l}
}

/**************************************************************
 * region HSV/HSB <=> RGB color
 ************************************************************
 * h,s,v/b = Hue, Saturation, Value(Brightness) 色相、饱和度、值（亮度）
 *
 * refers
 *  https://stackoverflow.com/questions/2353211/hsl-to-rgb-color-conversion
 *	https://github.com/less/less.js/blob/master/packages/less/src/less/functions/color.js
 *  https://github.com/d3/d3-color/blob/v3.0.1/README.md#hsl
 */

// function aliases
var (
	HsvToRgb = HSVToRGB
	// HsvIntToRgbInts alias for HSVIntToRGBInts
	HsvIntToRgbInts = HSVIntToRGBInts
)

// HSVIntToRGBInts Converts an HSL color value to RGB slice. Conversion formula
// adapted from https://en.wikipedia.org/wiki/HSL_and_HSV#HSV_to_RGB
//
//  Assumes h: 0-360, s: 0-100, l: 0-100
//  returns r, g, and b in the set [0, 255].
func HSVIntToRGBInts(h, s, v int) []uint8 {
	r, g, b := HSVToRGB(float64(h), float64(s)/100, float64(v)/100)
	return []uint8{r, g, b}
}

// HSVToRGB Convert HSV values to RGB values
//   - inputs: h (0-360), s (0-1.0), v (0-1.0)
//   - returns: r, g, b (0-255)
func HSVToRGB(h, s, v float64) (r, g, b uint8) {
	// 1. 处理特殊情况：饱和度为0（灰色）
	if s == 0 {
		gray := uint8(v * 255)
		return gray, gray, gray
	}

	// 2. 确保h在0-360范围内
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}

	// 3. 将h转换为0-6范围
	hSector := float64(h) / 60.0
	fraction := hSector - math.Floor(hSector)
	sector := int(math.Floor(hSector))

	// 4. 计算中间值
	p := v * (1 - s)
	q := v * (1 - fraction*s)
	t := v * (1 - (1-fraction)*s)

	// 5. 根据色相扇区计算RGB
	switch sector {
	case 0:
		r = uint8(v * 255)
		g = uint8(t * 255)
		b = uint8(p * 255)
	case 1:
		r = uint8(q * 255)
		g = uint8(v * 255)
		b = uint8(p * 255)
	case 2:
		r = uint8(p * 255)
		g = uint8(v * 255)
		b = uint8(t * 255)
	case 3:
		r = uint8(p * 255)
		g = uint8(q * 255)
		b = uint8(v * 255)
	case 4:
		r = uint8(t * 255)
		g = uint8(p * 255)
		b = uint8(v * 255)
	case 5:
		r = uint8(v * 255)
		g = uint8(p * 255)
		b = uint8(q * 255)
	}

	return r, g, b
}

// function aliases
var (
	RgbToHsv      = RGBToHSV
	RgbToHsvInts  = RGBToHSVInts
	RgbToHsvSlice = RGBToHSVSlice
)

// RGBToHSVInts convert RGB to HSV int values.
//
//	r, g, b: [0, 255]  => h (0-360), s (0-100), v (0-100)
func RGBToHSVInts(r, g, b uint8) []int {
	h, s, v := RGBToHSV(r, g, b)
	return []int{int(h), int(math.Round(s * 100)), int(math.Round(v * 100))}
}

// RGBToHSVSlice Convert RGB values to HSV values slice.
//
//	r, g, b (0-255)  => h (0-360), s (0-1.0), v (0-1.0)
func RGBToHSVSlice(r, g, b uint8) []float64 {
	h, s, v := RGBToHSV(r, g, b)
	return []float64{h, s, v}
}

// RGBToHSV Convert RGB values to HSV values
//   - inputs: r, g, b (0-255)
//   - returns: h (0-360), s (0-1.0), v (0-1.0)
func RGBToHSV(r, g, b uint8) (h, s, v float64) {
	// 1. 将RGB值归一化到 [0, 1] 范围
	rNorm := float64(r) / 255.0
	gNorm := float64(g) / 255.0
	bNorm := float64(b) / 255.0

	// 2. 找出最大值和最小值
	max1 := math.Max(math.Max(rNorm, gNorm), bNorm)
	min1 := math.Min(math.Min(rNorm, gNorm), bNorm)
	delta := max1 - min1

	// 3. 计算明度 (Value)
	v = max1

	// 4. 计算饱和度 (Saturation)
	if max1 == 0 {
		// 黑色情况，饱和度为0
		s = 0
	} else {
		s = delta / max1
	}

	// 5. 计算色相 (Hue)
	if delta == 0 {
		// 灰色情况，色相为0
		h = 0
	} else {
		switch max1 {
		case rNorm:
			h = (gNorm - bNorm) / delta
			if gNorm < bNorm {
				h += 6
			}
		case gNorm:
			h = (bNorm-rNorm)/delta + 2
		case bNorm:
			h = (rNorm-gNorm)/delta + 4
		}
		h *= 60 // 转换为角度 (0-360度)
	}

	return h, s, v
}

//
// region Named RGB color
//

// Named rgb colors
// https://www.w3.org/TR/css-color-3/#svg-color
var namedRgbMap = map[string]string{
	"aliceblue":            "240,248,255", // #F0F8FF
	"antiquewhite":         "250,235,215", // #FAEBD7
	"aqua":                 "0,255,255",   // #00FFFF
	"aquamarine":           "127,255,212", // #7FFFD4
	"azure":                "240,255,255", // #F0FFFF
	"beige":                "245,245,220", // #F5F5DC
	"bisque":               "255,228,196", // #FFE4C4
	"black":                "0,0,0",       // #000000
	"blanchedalmond":       "255,235,205", // #FFEBCD
	"blue":                 "0,0,255",     // #0000FF
	"blueviolet":           "138,43,226",  // #8A2BE2
	"brown":                "165,42,42",   // #A52A2A
	"burlywood":            "222,184,135", // #DEB887
	"cadetblue":            "95,158,160",  // #5F9EA0
	"chartreuse":           "127,255,0",   // #7FFF00
	"chocolate":            "210,105,30",  // #D2691E
	"coral":                "255,127,80",  // #FF7F50
	"cornflowerblue":       "100,149,237", // #6495ED
	"cornsilk":             "255,248,220", // #FFF8DC
	"crimson":              "220,20,60",   // #DC143C
	"cyan":                 "0,255,255",   // #00FFFF
	"darkblue":             "0,0,139",     // #00008B
	"darkcyan":             "0,139,139",   // #008B8B
	"darkgoldenrod":        "184,134,11",  // #B8860B
	"darkgray":             "169,169,169", // #A9A9A9
	"darkgreen":            "0,100,0",     // #006400
	"darkgrey":             "169,169,169", // #A9A9A9
	"darkkhaki":            "189,183,107", // #BDB76B
	"darkmagenta":          "139,0,139",   // #8B008B
	"darkolivegreen":       "85,107,47",   // #556B2F
	"darkorange":           "255,140,0",   // #FF8C00
	"darkorchid":           "153,50,204",  // #9932CC
	"darkred":              "139,0,0",     // #8B0000
	"darksalmon":           "233,150,122", // #E9967A
	"darkseagreen":         "143,188,143", // #8FBC8F
	"darkslateblue":        "72,61,139",   // #483D8B
	"darkslategray":        "47,79,79",    // #2F4F4F
	"darkslategrey":        "47,79,79",    // #2F4F4F
	"darkturquoise":        "0,206,209",   // #00CED1
	"darkviolet":           "148,0,211",   // #9400D3
	"deeppink":             "255,20,147",  // #FF1493
	"deepskyblue":          "0,191,255",   // #00BFFF
	"dimgray":              "105,105,105", // #696969
	"dimgrey":              "105,105,105", // #696969
	"dodgerblue":           "30,144,255",  // #1E90FF
	"firebrick":            "178,34,34",   // #B22222
	"floralwhite":          "255,250,240", // #FFFAF0
	"forestgreen":          "34,139,34",   // #228B22
	"fuchsia":              "255,0,255",   // #FF00FF
	"gainsboro":            "220,220,220", // #DCDCDC
	"ghostwhite":           "248,248,255", // #F8F8FF
	"gold":                 "255,215,0",   // #FFD700
	"goldenrod":            "218,165,32",  // #DAA520
	"gray":                 "128,128,128", // #808080
	"green":                "0,128,0",     // #008000
	"greenyellow":          "173,255,47",  // #ADFF2F
	"grey":                 "128,128,128", // #808080
	"honeydew":             "240,255,240", // #F0FFF0
	"hotpink":              "255,105,180", // #FF69B4
	"indianred":            "205,92,92",   // #CD5C5C
	"indigo":               "75,0,130",    // #4B0082
	"ivory":                "255,255,240", // #FFFFF0
	"khaki":                "240,230,140", // #F0E68C
	"lavender":             "230,230,250", // #E6E6FA
	"lavenderblush":        "255,240,245", // #FFF0F5
	"lawngreen":            "124,252,0",   // #7CFC00
	"lemonchiffon":         "255,250,205", // #FFFACD
	"lightblue":            "173,216,230", // #ADD8E6
	"lightcoral":           "240,128,128", // #F08080
	"lightcyan":            "224,255,255", // #E0FFFF
	"lightgoldenrodyellow": "250,250,210", // #FAFAD2
	"lightgray":            "211,211,211", // #D3D3D3
	"lightgreen":           "144,238,144", // #90EE90
	"lightgrey":            "211,211,211", // #D3D3D3
	"lightpink":            "255,182,193", // #FFB6C1
	"lightsalmon":          "255,160,122", // #FFA07A
	"lightseagreen":        "32,178,170",  // #20B2AA
	"lightskyblue":         "135,206,250", // #87CEFA
	"lightslategray":       "119,136,153", // #778899
	"lightslategrey":       "119,136,153", // #778899
	"lightsteelblue":       "176,196,222", // #B0C4DE
	"lightyellow":          "255,255,224", // #FFFFE0
	"lime":                 "0,255,0",     // #00FF00
	"limegreen":            "50,205,50",   // #32CD32
	"linen":                "250,240,230", // #FAF0E6
	"magenta":              "255,0,255",   // #FF00FF
	"maroon":               "128,0,0",     // #800000
	"mediumaquamarine":     "102,205,170", // #66CDAA
	"mediumblue":           "0,0,205",     // #0000CD
	"mediumorchid":         "186,85,211",  // #BA55D3
	"mediumpurple":         "147,112,219", // #9370DB
	"mediumseagreen":       "60,179,113",  // #3CB371
	"mediumslateblue":      "123,104,238", // #7B68EE
	"mediumspringgreen":    "0,250,154",   // #00FA9A
	"mediumturquoise":      "72,209,204",  // #48D1CC
	"mediumvioletred":      "199,21,133",  // #C71585
	"midnightblue":         "25,25,112",   // #191970
	"mintcream":            "245,255,250", // #F5FFFA
	"mistyrose":            "255,228,225", // #FFE4E1
	"moccasin":             "255,228,181", // #FFE4B5
	"navajowhite":          "255,222,173", // #FFDEAD
	"navy":                 "0,0,128",     // #000080
	"oldlace":              "253,245,230", // #FDF5E6
	"olive":                "128,128,0",   // #808000
	"olivedrab":            "107,142,35",  // #6B8E23
	"orange":               "255,165,0",   // #FFA500
	"orangered":            "255,69,0",    // #FF4500
	"orchid":               "218,112,214", // #DA70D6
	"palegoldenrod":        "238,232,170", // #EEE8AA
	"palegreen":            "152,251,152", // #98FB98
	"paleturquoise":        "175,238,238", // #AFEEEE
	"palevioletred":        "219,112,147", // #DB7093
	"papayawhip":           "255,239,213", // #FFEFD5
	"peachpuff":            "255,218,185", // #FFDAB9
	"peru":                 "205,133,63",  // #CD853F
	"pink":                 "255,192,203", // #FFC0CB
	"plum":                 "221,160,221", // #DDA0DD
	"powderblue":           "176,224,230", // #B0E0E6
	"purple":               "128,0,128",   // #800080
	"red":                  "255,0,0",     // #FF0000
	"rosybrown":            "188,143,143", // #BC8F8F
	"royalblue":            "65,105,225",  // #4169E1
	"saddlebrown":          "139,69,19",   // #8B4513
	"salmon":               "250,128,114", // #FA8072
	"sandybrown":           "244,164,96",  // #F4A460
	"seagreen":             "46,139,87",   // #2E8B57
	"seashell":             "255,245,238", // #FFF5EE
	"sienna":               "160,82,45",   // #A0522D
	"silver":               "192,192,192", // #C0C0C0
	"skyblue":              "135,206,235", // #87CEEB
	"slateblue":            "106,90,205",  // #6A5ACD
	"slategray":            "112,128,144", // #708090
	"slategrey":            "112,128,144", // #708090
	"snow":                 "255,250,250", // #FFFAFA
	"springgreen":          "0,255,127",   // #00FF7F
	"steelblue":            "70,130,180",  // #4682B4
	"tan":                  "210,180,140", // #D2B48C
	"teal":                 "0,128,128",   // #008080
	"thistle":              "216,191,216", // #D8BFD8
	"tomato":               "255,99,71",   // #FF6347
	"turquoise":            "64,224,208",  // #40E0D0
	"violet":               "238,130,238", // #EE82EE
	"wheat":                "245,222,179", // #F5DEB3
	"white":                "255,255,255", // #FFFFFF
	"whitesmoke":           "245,245,245", // #F5F5F5
	"yellow":               "255,255,0",   // #FFFF00
	"yellowgreen":          "154,205,50",  // #9ACD32
}
