package color

import (
	"fmt"
	"regexp"
	"strings"
)

// output colored text like use html tag. (not support windows cmd)
const (
	// MatchExpr regex to match color tags
	//
	// Notice: golang 不支持反向引用. 即不支持使用 \1 引用第一个匹配 ([a-z=;]+)
	// MatchExpr = `<([a-z=;]+)>(.*?)<\/\1>`
	// 所以调整一下 统一使用 `</>` 来结束标签，例如 "<info>some text</>"
	//
	// allow custom attrs, eg: "<fg=white;bg=blue;op=bold>content</>"
	// (?s:...) s - 让 "." 匹配换行
	MatchExpr = `<([0-9a-zA-Z_=,;]+)>(?s:(.*?))<\/>`

	// AttrExpr regex to match custom color attributes
	// eg: "<fg=white;bg=blue;op=bold>content</>"
	AttrExpr = `(fg|bg|op)[\s]*=[\s]*([0-9a-zA-Z,]+);?`

	// StripExpr regex used for removing color tags
	// StripExpr = `<[\/]?[a-zA-Z=;]+>`
	// 随着上面的做一些调整
	StripExpr = `<[\/]?[0-9a-zA-Z_=,;]*>`
)

var (
	attrRegex  = regexp.MustCompile(AttrExpr)
	matchRegex = regexp.MustCompile(MatchExpr)
	stripRegex = regexp.MustCompile(StripExpr)
)

/*************************************************************
 * internal defined color tags
 *************************************************************/

// There are internal defined fg color tags
//
// Usage:
//
//	<tag>content text</>
//
// @notice 加 0 在前面是为了防止之前的影响到现在的设置
var colorTags = map[string]string{
	// basic tags
	"red":      "0;31",
	"red1":     "1;31", // with bold
	"redB":     "1;31",
	"red_b":    "1;31",
	"blue":     "0;34",
	"blue1":    "1;34", // with bold
	"blueB":    "1;34",
	"blue_b":   "1;34",
	"cyan":     "0;36",
	"cyan1":    "1;36", // with bold
	"cyanB":    "1;36",
	"cyan_b":   "1;36",
	"green":    "0;32",
	"green1":   "1;32", // with bold
	"greenB":   "1;32",
	"green_b":  "1;32",
	"black":    "0;30",
	"white":    "1;37",
	"default":  "0;39", // no color
	"normal":   "0;39", // no color
	"brown":    "0;33", // #A52A2A
	"yellow":   "0;33",
	"ylw0":     "0;33",
	"yellowB":  "1;33", // with bold
	"ylw1":     "1;33",
	"ylwB":     "1;33",
	"magenta":  "0;35",
	"mga":      "0;35", // short name
	"magentaB": "1;35", // with bold
	"magenta1": "1;35",
	"mgb":      "1;35",
	"mga1":     "1;35",
	"mgaB":     "1;35",

	// light/hi tags

	"gray":          "0;90",
	"darkGray":      "0;90",
	"dark_gray":     "0;90",
	"lightYellow":   "0;93",
	"light_yellow":  "0;93",
	"hiYellow":      "0;93",
	"hi_yellow":     "0;93",
	"hiYellowB":     "1;93", // with bold
	"hi_yellow_b":   "1;93",
	"lightMagenta":  "0;95",
	"light_magenta": "0;95",
	"hiMagenta":     "0;95",
	"hi_magenta":    "0;95",
	"lightMagenta1": "1;95", // with bold
	"hiMagentaB":    "1;95", // with bold
	"hi_magenta_b":  "1;95",
	"lightRed":      "0;91",
	"light_red":     "0;91",
	"hiRed":         "0;91",
	"hi_red":        "0;91",
	"lightRedB":     "1;91", // with bold
	"light_red_b":   "1;91",
	"hi_red_b":      "1;91",
	"lightGreen":    "0;92",
	"light_green":   "0;92",
	"hiGreen":       "0;92",
	"hi_green":      "0;92",
	"lightGreenB":   "1;92",
	"light_green_b": "1;92",
	"hi_green_b":    "1;92",
	"lightBlue":     "0;94",
	"light_blue":    "0;94",
	"hiBlue":        "0;94",
	"hi_blue":       "0;94",
	"lightBlueB":    "1;94",
	"light_blue_b":  "1;94",
	"hi_blue_b":     "1;94",
	"lightCyan":     "0;96",
	"light_cyan":    "0;96",
	"hiCyan":        "0;96",
	"hi_cyan":       "0;96",
	"lightCyanB":    "1;96",
	"light_cyan_b":  "1;96",
	"hi_cyan_b":     "1;96",
	"lightWhite":    "0;97;40",
	"light_white":   "0;97;40",

	// option
	"bold":       "1",
	"b":          "1",
	"italic":     "3",
	"i":          "3", // italic
	"underscore": "4",
	"us":         "4", // short name for 'underscore'
	"blink":      "5",
	"fb":         "6", // fast blink
	"reverse":    "7",
	"st":         "9", // strikethrough

	// alert tags, like bootstrap's alert
	"suc":     "1;32", // same "green" and "bold"
	"success": "1;32",
	"info":    "0;32", // same "green",
	"comment": "0;33", // same "brown"
	"note":    "36;1",
	"notice":  "36;4",
	"warn":    "0;1;33",
	"warning": "0;30;43",
	"primary": "0;34",
	"danger":  "1;31", // same "red" but add bold
	"err":     "97;41",
	"error":   "97;41", // fg light white; bg red
}

/*************************************************************
 * internal defined tag attributes
 *************************************************************/

// built-in attributes for fg,bg 16-colors and op codes.
var (
	attrFgs = map[string]string{
		// basic colors

		"black":   FgBlack.Code(),
		"red":     "31",
		"green":   "32",
		"brown":   "33", // #A52A2A
		"yellow":  "33",
		"ylw":     "33",
		"blue":    "34",
		"cyan":    "36",
		"magenta": "35",
		"mga":     "35",
		"white":   FgWhite.Code(),
		"default": "39", // no color
		"normal":  "39", // no color

		// light/hi colors

		"darkGray":      FgDarkGray.Code(),
		"dark_gray":     "90",
		"gray":          "90",
		"lightYellow":   "93",
		"light_yellow":  "93",
		"hiYellow":      "93",
		"hi_yellow":     "93",
		"lightMagenta":  "95",
		"light_magenta": "95",
		"hiMagenta":     "95",
		"hi_magenta":    "95",
		"hi_mga":        "95",
		"lightRed":      "91",
		"light_red":     "91",
		"hiRed":         "91",
		"hi_red":        "91",
		"lightGreen":    "92",
		"light_green":   "92",
		"hiGreen":       "92",
		"hi_green":      "92",
		"lightBlue":     "94",
		"light_blue":    "94",
		"hiBlue":        "94",
		"hi_blue":       "94",
		"lightCyan":     "96",
		"light_cyan":    "96",
		"hiCyan":        "96",
		"hi_cyan":       "96",
		"lightWhite":    "97",
		"light_white":   "97",
	}

	attrBgs = map[string]string{
		// basic colors

		"black":   BgBlack.Code(),
		"red":     "41",
		"green":   "42",
		"brown":   "43", // #A52A2A
		"yellow":  "43",
		"ylw":     "43",
		"blue":    "44",
		"cyan":    "46",
		"magenta": "45",
		"mga":     "45",
		"white":   FgWhite.Code(),
		"default": "49", // no color
		"normal":  "49", // no color

		// light/hi colors

		"darkGray":      BgDarkGray.Code(),
		"dark_gray":     "100",
		"gray":          "100",
		"lightYellow":   "103",
		"light_yellow":  "103",
		"hiYellow":      "103",
		"hi_yellow":     "103",
		"lightMagenta":  "105",
		"light_magenta": "105",
		"hiMagenta":     "105",
		"hi_magenta":    "105",
		"hi_mga":        "105",
		"lightRed":      "101",
		"light_red":     "101",
		"hiRed":         "101",
		"hi_red":        "101",
		"lightGreen":    "102",
		"light_green":   "102",
		"hiGreen":       "102",
		"hi_green":      "102",
		"lightBlue":     "104",
		"light_blue":    "104",
		"hiBlue":        "104",
		"hi_blue":       "104",
		"lightCyan":     "106",
		"light_cyan":    "106",
		"hiCyan":        "106",
		"hi_cyan":       "106",
		"lightWhite":    BgLightWhite.Code(),
		"light_white":   "107",
	}

	attrOpts = map[string]string{
		"reset":         OpReset.Code(),
		"bold":          OpBold.Code(),
		"b":             OpBold.Code(),
		"fuzzy":         OpFuzzy.Code(),
		"italic":        OpItalic.Code(),
		"i":             OpItalic.Code(),
		"underscore":    OpUnderscore.Code(),
		"us":            OpUnderscore.Code(),
		"u":             OpUnderscore.Code(),
		"blink":         OpBlink.Code(),
		"fastblink":     OpFastBlink.Code(),
		"fb":            OpFastBlink.Code(),
		"reverse":       OpReverse.Code(),
		"concealed":     OpConcealed.Code(),
		"strikethrough": OpStrikethrough.Code(),
		"st":            OpStrikethrough.Code(),
	}
)

/*************************************************************
 * parse color tags
 *************************************************************/

var (
	tagParser = TagParser{}
	// regex for match color 256 code
	rxNumStr  = regexp.MustCompile("^[0-9]{1,3}$")
	rxHexCode = regexp.MustCompile("^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$")
)

// TagParser struct
type TagParser struct {
	disable bool
}

// NewTagParser create
func NewTagParser() *TagParser {
	return &TagParser{}
}

// func (tp *TagParser) Disable() *TagParser {
// 	tp.disable = true
// 	return tp
// }

// ParseByEnv parse given string. will check package setting.
func (tp *TagParser) ParseByEnv(str string) string {
	// disable handler TAG
	if !RenderTag {
		return str
	}

	// disable OR not support color
	if !Enable || !SupportColor() {
		return ClearTag(str)
	}
	return tp.Parse(str)
}

// Parse given string, replace color tag and return rendered string
//
// Use built in tags:
//
//	<TAG_NAME>CONTENT</>
//	// e.g: `<info>message</>`
//
// Custom tag attributes:
//
//	`<fg=VALUE;bg=VALUE;op=VALUES>CONTENT</>`
//	// e.g: `<fg=167;bg=232>wel</>`
func (tp *TagParser) Parse(str string) string {
	// not contains color tag
	if !strings.Contains(str, "</>") {
		return str
	}

	// find color tags by regex. str eg: "<fg=white;bg=blue;op=bold>content</>"
	matched := matchRegex.FindAllStringSubmatch(str, -1)

	// item: 0 full text 1 tag name 2 tag content
	for _, item := range matched {
		full, tag, body := repairMatchedTag(item[0], item[1], item[2])

		// use defined color tag name: "<info>content</>" -> tag: "info"
		if !strings.ContainsRune(tag, '=') {
			if code := colorTags[tag]; len(code) > 0 {
				str = strings.Replace(str, full, RenderString(code, body), 1)
			} else if code, ok := namedRgbMap[tag]; ok {
				code = strings.Replace(code, ",", ";", -1)
				now := RenderString(FgRGBPfx+code, body)
				str = strings.Replace(str, full, now, 1)
			}
			continue
		}

		// custom color in tag
		// - basic: "fg=white;bg=blue;op=bold"
		if code := ParseCodeFromAttr(tag); len(code) > 0 {
			str = strings.Replace(str, full, RenderString(code, body), 1)
		}
	}

	return str
}

func repairMatchedTag(full, tag, body string) (string, string, string) {
	if len(colorTags[tag]) == 0 && len(namedRgbMap[tag]) == 0 && len(ParseCodeFromAttr(tag)) == 0 {
		if matched := matchRegex.FindAllStringSubmatch(strings.TrimPrefix(full, "<"+tag+">"), -1); len(matched) > 0 {
			full = matched[0][0]
			tag = matched[0][1]
			body = matched[0][2]
			return repairMatchedTag(full, tag, body)
		}
	}

	return full, tag, body
}

// ReplaceTag parse string, replace color tag and return rendered string
func ReplaceTag(str string) string {
	return tagParser.ParseByEnv(str)
}

// ParseCodeFromAttr parse color attributes.
//
// attr format:
//
//	// VALUE please see var: FgColors, BgColors, AllOptions
//	"fg=VALUE;bg=VALUE;op=VALUE"
//
// 16 color:
//
//	"fg=yellow"
//	"bg=red"
//	"op=bold,underscore" // option is allow multi value
//	"fg=white;bg=blue;op=bold"
//	"fg=white;op=bold,underscore"
//
// 256 color:
//
//	"fg=167"
//	"fg=167;bg=23"
//	"fg=167;bg=23;op=bold"
//
// True color:
//
//	// hex
//	"fg=fc1cac"
//	"fg=fc1cac;bg=c2c3c4"
//	// r,g,b
//	"fg=23,45,214"
//	"fg=23,45,214;bg=109,99,88"
func ParseCodeFromAttr(attr string) (code string) {
	if !strings.ContainsRune(attr, '=') {
		return
	}

	attr = strings.Trim(attr, ";=,")
	if len(attr) == 0 {
		return
	}

	var codes []string
	matched := attrRegex.FindAllStringSubmatch(attr, -1)

	for _, item := range matched {
		pos, val := item[1], item[2]
		switch pos {
		case "fg":
			if code, ok := attrFgs[val]; ok { // attr fg
				codes = append(codes, code)
			} else if code := rgbHex256toCode(val, false); code != "" {
				codes = append(codes, code)
			}
		case "bg":
			if code, ok := attrBgs[val]; ok { // attr bg
				codes = append(codes, code)
			} else if code := rgbHex256toCode(val, true); code != "" {
				codes = append(codes, code)
			}
		case "op": // options allow multi value
			if strings.Contains(val, ",") {
				ns := strings.Split(val, ",")
				for _, n := range ns {
					if code, ok := attrOpts[n]; ok { // attr ops
						codes = append(codes, code)
					}
				}
			} else if code, ok := attrOpts[val]; ok {
				codes = append(codes, code)
			}
		}
	}

	return strings.Join(codes, ";")
}

func rgbHex256toCode(val string, isBg bool) (code string) {
	if len(val) == 6 && rxHexCode.MatchString(val) { // hex: "fc1cac"
		code = HEX(val, isBg).String()
	} else if strings.ContainsRune(val, ',') { // rgb: "231,178,161"
		code = strings.Replace(val, ",", ";", -1)
		if isBg {
			code = BgRGBPfx + code
		} else {
			code = FgRGBPfx + code
		}
	} else if len(val) < 4 && rxNumStr.MatchString(val) { // 256 code
		if isBg {
			code = Bg256Pfx + val
		} else {
			code = Fg256Pfx + val
		}
	}
	return
}

// ClearTag clear all tag for a string
func ClearTag(s string) string {
	if !strings.Contains(s, "</>") {
		return s
	}
	return stripRegex.ReplaceAllString(s, "")
}

/*************************************************************
 * helper methods
 *************************************************************/

// GetTagCode get color code by tag name
func GetTagCode(name string) string { return colorTags[name] }

// ApplyTag for messages
func ApplyTag(tag string, a ...any) string {
	return RenderCode(GetTagCode(tag), a...)
}

// WrapTag wrap a tag for a string "<tag>content</>"
func WrapTag(s string, tag string) string {
	if s == "" || tag == "" {
		return s
	}
	return fmt.Sprintf("<%s>%s</>", tag, s)
}

// GetColorTags get all internal color tags
func GetColorTags() map[string]string {
	return colorTags
}

// IsDefinedTag is defined tag name
func IsDefinedTag(name string) bool {
	_, ok := colorTags[name]
	return ok
}

/*************************************************************
 * Tag extra
 *************************************************************/

// Tag value is a defined style name
// Usage:
//
//	Tag("info").Println("message")
type Tag string

// Print messages
func (tg Tag) Print(a ...any) {
	name := string(tg)
	str := fmt.Sprint(a...)

	if stl := GetStyle(name); !stl.IsEmpty() {
		stl.Print(str)
	} else {
		doPrintV2(GetTagCode(name), str)
	}
}

// Printf format and print messages
func (tg Tag) Printf(format string, a ...any) {
	name := string(tg)
	str := fmt.Sprintf(format, a...)

	if stl := GetStyle(name); !stl.IsEmpty() {
		stl.Print(str)
	} else {
		doPrintV2(GetTagCode(name), str)
	}
}

// Println messages line
func (tg Tag) Println(a ...any) {
	name := string(tg)
	if stl := GetStyle(name); !stl.IsEmpty() {
		stl.Println(a...)
	} else {
		doPrintlnV2(GetTagCode(name), a)
	}
}

// Sprint render messages
func (tg Tag) Sprint(a ...any) string {
	return RenderCode(GetTagCode(string(tg)), a...)
}

// Sprintf format and render messages
func (tg Tag) Sprintf(format string, a ...any) string {
	tag := string(tg)
	str := fmt.Sprintf(format, a...)

	return RenderString(GetTagCode(tag), str)
}
