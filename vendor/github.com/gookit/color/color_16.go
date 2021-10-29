package color

import (
	"fmt"
	"strconv"
)

// Color Color16, 16 color value type
// 3(2^3=8) OR 4(2^4=16) bite color.
type Color uint8
type Basic = Color // alias of Color

// Opts basic color options. code: 0 - 9
type Opts []Color

// Add option value
func (o *Opts) Add(ops ...Color) {
	for _, op := range ops {
		if uint8(op) < 10 {
			*o = append(*o, op)
		}
	}
}

// IsValid options
func (o Opts) IsValid() bool {
	return len(o) > 0
}

// IsEmpty options
func (o Opts) IsEmpty() bool {
	return len(o) == 0
}

// String options to string. eg: "1;3"
func (o Opts) String() string {
	return Colors2code(o...)
}

/*************************************************************
 * Basic 16 color definition
 *************************************************************/

// Base value for foreground/background color
const (
	FgBase uint8 = 30
	BgBase uint8 = 40
	// hi color base code
	HiFgBase uint8 = 90
	HiBgBase uint8 = 100
)

// Foreground colors. basic foreground colors 30 - 37
const (
	FgBlack Color = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta // 品红
	FgCyan    // 青色
	FgWhite
	// FgDefault revert default FG
	FgDefault Color = 39
)

// Extra foreground color 90 - 97(非标准)
const (
	FgDarkGray Color = iota + 90 // 亮黑（灰）
	FgLightRed
	FgLightGreen
	FgLightYellow
	FgLightBlue
	FgLightMagenta
	FgLightCyan
	FgLightWhite
	// FgGray is alias of FgDarkGray
	FgGray Color = 90 // 亮黑（灰）
)

// Background colors. basic background colors 40 - 47
const (
	BgBlack Color = iota + 40
	BgRed
	BgGreen
	BgYellow // BgBrown like yellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
	// BgDefault revert default BG
	BgDefault Color = 49
)

// Extra background color 100 - 107(非标准)
const (
	BgDarkGray Color = iota + 100
	BgLightRed
	BgLightGreen
	BgLightYellow
	BgLightBlue
	BgLightMagenta
	BgLightCyan
	BgLightWhite
	// BgGray is alias of BgDarkGray
	BgGray Color = 100
)

// Option settings
const (
	OpReset         Color = iota // 0 重置所有设置
	OpBold                       // 1 加粗
	OpFuzzy                      // 2 模糊(不是所有的终端仿真器都支持)
	OpItalic                     // 3 斜体(不是所有的终端仿真器都支持)
	OpUnderscore                 // 4 下划线
	OpBlink                      // 5 闪烁
	OpFastBlink                  // 5 快速闪烁(未广泛支持)
	OpReverse                    // 7 颠倒的 交换背景色与前景色
	OpConcealed                  // 8 隐匿的
	OpStrikethrough              // 9 删除的，删除线(未广泛支持)
)

// There are basic and light foreground color aliases
const (
	Red     = FgRed
	Cyan    = FgCyan
	Gray    = FgDarkGray // is light Black
	Blue    = FgBlue
	Black   = FgBlack
	Green   = FgGreen
	White   = FgWhite
	Yellow  = FgYellow
	Magenta = FgMagenta

	// special

	Bold   = OpBold
	Normal = FgDefault

	// extra light

	LightRed     = FgLightRed
	LightCyan    = FgLightCyan
	LightBlue    = FgLightBlue
	LightGreen   = FgLightGreen
	LightWhite   = FgLightWhite
	LightYellow  = FgLightYellow
	LightMagenta = FgLightMagenta

	HiRed     = FgLightRed
	HiCyan    = FgLightCyan
	HiBlue    = FgLightBlue
	HiGreen   = FgLightGreen
	HiWhite   = FgLightWhite
	HiYellow  = FgLightYellow
	HiMagenta = FgLightMagenta

	BgHiRed     = BgLightRed
	BgHiCyan    = BgLightCyan
	BgHiBlue    = BgLightBlue
	BgHiGreen   = BgLightGreen
	BgHiWhite   = BgLightWhite
	BgHiYellow  = BgLightYellow
	BgHiMagenta = BgLightMagenta
)

// Bit4 an method for create Color
func Bit4(code uint8) Color {
	return Color(code)
}

/*************************************************************
 * Color render methods
 *************************************************************/

// Name get color code name.
func (c Color) Name() string {
	name, ok := basic2nameMap[uint8(c)]
	if ok {
		return name
	}
	return "unknown"
}

// Text render a text message
func (c Color) Text(message string) string {
	return RenderString(c.String(), message)
}

// Render messages by color setting
// Usage:
// 		green := color.FgGreen.Render
// 		fmt.Println(green("message"))
func (c Color) Render(a ...interface{}) string {
	return RenderCode(c.String(), a...)
}

// Renderln messages by color setting.
// like Println, will add spaces for each argument
// Usage:
// 		green := color.FgGreen.Renderln
// 		fmt.Println(green("message"))
func (c Color) Renderln(a ...interface{}) string {
	return RenderWithSpaces(c.String(), a...)
}

// Sprint render messages by color setting. is alias of the Render()
func (c Color) Sprint(a ...interface{}) string {
	return RenderCode(c.String(), a...)
}

// Sprintf format and render message.
// Usage:
// 	green := color.Green.Sprintf
//  colored := green("message")
func (c Color) Sprintf(format string, args ...interface{}) string {
	return RenderString(c.String(), fmt.Sprintf(format, args...))
}

// Print messages.
// Usage:
// 		color.Green.Print("message")
// OR:
// 		green := color.FgGreen.Print
// 		green("message")
func (c Color) Print(args ...interface{}) {
	doPrintV2(c.Code(), fmt.Sprint(args...))
}

// Printf format and print messages.
// Usage:
// 		color.Cyan.Printf("string %s", "arg0")
func (c Color) Printf(format string, a ...interface{}) {
	doPrintV2(c.Code(), fmt.Sprintf(format, a...))
}

// Println messages with new line
func (c Color) Println(a ...interface{}) {
	doPrintlnV2(c.String(), a)
}

// Light current color. eg: 36(FgCyan) -> 96(FgLightCyan).
// Usage:
// 	lightCyan := Cyan.Light()
// 	lightCyan.Print("message")
func (c Color) Light() Color {
	val := int(c)
	if val >= 30 && val <= 47 {
		return Color(uint8(c) + 60)
	}

	// don't change
	return c
}

// Darken current color. eg. 96(FgLightCyan) -> 36(FgCyan)
// Usage:
// 	cyan := LightCyan.Darken()
// 	cyan.Print("message")
func (c Color) Darken() Color {
	val := int(c)
	if val >= 90 && val <= 107 {
		return Color(uint8(c) - 60)
	}

	// don't change
	return c
}

// C256 convert 16 color to 256-color code.
func (c Color) C256() Color256 {
	val := uint8(c)
	if val < 10 { // is option code
		return emptyC256 // empty
	}

	var isBg uint8
	if val >= BgBase && val <= 47 { // is bg
		isBg = AsBg
		val = val - 10 // to fg code
	} else if val >= HiBgBase && val <= 107 { // is hi bg
		isBg = AsBg
		val = val - 10 // to fg code
	}

	if c256, ok := basicTo256Map[val]; ok {
		return Color256{c256, isBg}
	}

	// use raw value direct convert
	return Color256{val}
}

// RGB convert 16 color to 256-color code.
func (c Color) RGB() RGBColor {
	val := uint8(c)
	if val < 10 { // is option code
		return emptyRGBColor
	}

	return HEX(Basic2hex(val))
}

// Code convert to code string. eg "35"
func (c Color) Code() string {
	// return fmt.Sprintf("%d", c)
	return strconv.Itoa(int(c))
}

// String convert to code string. eg "35"
func (c Color) String() string {
	// return fmt.Sprintf("%d", c)
	return strconv.Itoa(int(c))
}

// IsValid color value
func (c Color) IsValid() bool {
	return c < 107
}

/*************************************************************
 * basic color maps
 *************************************************************/

// FgColors foreground colors map
var FgColors = map[string]Color{
	"black":   FgBlack,
	"red":     FgRed,
	"green":   FgGreen,
	"yellow":  FgYellow,
	"blue":    FgBlue,
	"magenta": FgMagenta,
	"cyan":    FgCyan,
	"white":   FgWhite,
	"default": FgDefault,
}

// BgColors background colors map
var BgColors = map[string]Color{
	"black":   BgBlack,
	"red":     BgRed,
	"green":   BgGreen,
	"yellow":  BgYellow,
	"blue":    BgBlue,
	"magenta": BgMagenta,
	"cyan":    BgCyan,
	"white":   BgWhite,
	"default": BgDefault,
}

// ExFgColors extra foreground colors map
var ExFgColors = map[string]Color{
	"darkGray":     FgDarkGray,
	"lightRed":     FgLightRed,
	"lightGreen":   FgLightGreen,
	"lightYellow":  FgLightYellow,
	"lightBlue":    FgLightBlue,
	"lightMagenta": FgLightMagenta,
	"lightCyan":    FgLightCyan,
	"lightWhite":   FgLightWhite,
}

// ExBgColors extra background colors map
var ExBgColors = map[string]Color{
	"darkGray":     BgDarkGray,
	"lightRed":     BgLightRed,
	"lightGreen":   BgLightGreen,
	"lightYellow":  BgLightYellow,
	"lightBlue":    BgLightBlue,
	"lightMagenta": BgLightMagenta,
	"lightCyan":    BgLightCyan,
	"lightWhite":   BgLightWhite,
}

// Options color options map
// Deprecated
// NOTICE: please use AllOptions instead.
var Options = AllOptions

// AllOptions color options map
var AllOptions = map[string]Color{
	"reset":      OpReset,
	"bold":       OpBold,
	"fuzzy":      OpFuzzy,
	"italic":     OpItalic,
	"underscore": OpUnderscore,
	"blink":      OpBlink,
	"reverse":    OpReverse,
	"concealed":  OpConcealed,
}

var (
	// TODO basic name alias
	// basicNameAlias = map[string]string{}

	// basic color name to code
	name2basicMap = initName2basicMap()
	// basic2nameMap basic color code to name
	basic2nameMap = map[uint8]string{
		30: "black",
		31: "red",
		32: "green",
		33: "yellow",
		34: "blue",
		35: "magenta",
		36: "cyan",
		37: "white",
		// hi color code
		90: "lightBlack",
		91: "lightRed",
		92: "lightGreen",
		93: "lightYellow",
		94: "lightBlue",
		95: "lightMagenta",
		96: "lightCyan",
		97: "lightWhite",
		// options
		0: "reset",
		1: "bold",
		2: "fuzzy",
		3: "italic",
		4: "underscore",
		5: "blink",
		7: "reverse",
		8: "concealed",
	}
)

// Basic2nameMap data
func Basic2nameMap() map[uint8]string {
	return basic2nameMap
}

func initName2basicMap() map[string]uint8 {
	n2b := make(map[string]uint8, len(basic2nameMap))
	for u, s := range basic2nameMap {
		n2b[s] = u
	}
	return n2b
}
