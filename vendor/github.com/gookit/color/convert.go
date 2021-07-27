package color

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

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
	basic2hexMap = map[uint8]string{
		30: "000000", // black
		31: "c51e14", // red
		32: "1dc121", // green
		33: "c7c329", // yellow
		34: "0a2fc4", // blue
		35: "c839c5", // magenta
		36: "20c5c6", // cyan
		37: "c7c7c7", // white
		90: "686868", // lightBlack/darkGray
		91: "fd6f6b", // lightRed
		92: "67f86f", // lightGreen
		93: "fffa72", // lightYellow
		94: "6a76fb", // lightBlue
		95: "fd7cfc", // lightMagenta
		96: "68fdfe", // lightCyan
		97: "ffffff", // lightWhite
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
 * HEX code <=> RGB/True color code
 *************************************************************/

// Hex2rgb alias of the HexToRgb()
func Hex2rgb(hex string) []int { return HexToRgb(hex) }

// HexToRGB alias of the HexToRgb()
func HexToRGB(hex string) []int { return HexToRgb(hex) }

// HexToRgb convert hex color string to RGB numbers
//
// Usage:
// 	rgb := HexToRgb("ccc") // rgb: [204 204 204]
// 	rgb := HexToRgb("aabbcc") // rgb: [170 187 204]
// 	rgb := HexToRgb("#aabbcc") // rgb: [170 187 204]
// 	rgb := HexToRgb("0xad99c0") // rgb: [170 187 204]
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
//	hex := RgbToHex([]int{170, 187, 204}) // hex: "aabbcc"
func RgbToHex(rgb []int) string {
	hexNodes := make([]string, len(rgb))

	for _, v := range rgb {
		hexNodes = append(hexNodes, strconv.FormatInt(int64(v), 16))
	}
	return strings.Join(hexNodes, "")
}

/*************************************************************
 * 4bit(16) color <=> RGB/True color
 *************************************************************/

// Basic2hex convert basic color to hex string.
func Basic2hex(val uint8) string {
	return basic2hexMap[val]
}

// Hex2basic convert hex string to basic color code.
func Hex2basic(hex string) uint8 {
	return hex2basicMap[hex]
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

// Rgb2ansi alias of the RgbToAnsi()
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
