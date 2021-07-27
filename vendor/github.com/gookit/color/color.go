/*
Package color is Command line color library.
Support rich color rendering output, universal API method, compatible with Windows system

Source code and other details for the project are available at GitHub:

	https://github.com/gookit/color

More usage please see README and tests.
*/
package color

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/xo/terminfo"
)

// terminal color available level alias of the terminfo.ColorLevel*
const (
	LevelNo  = terminfo.ColorLevelNone     // not support color.
	Level16  = terminfo.ColorLevelBasic    // 3/4 bit color supported
	Level256 = terminfo.ColorLevelHundreds // 8 bit color supported
	LevelRgb = terminfo.ColorLevelMillions // (24 bit)true color supported
)

// color render templates
// ESC 操作的表示:
// 	"\033"(Octal 8进制) = "\x1b"(Hexadecimal 16进制) = 27 (10进制)
const (
	SettingTpl   = "\x1b[%sm"
	FullColorTpl = "\x1b[%sm%s\x1b[0m"
)

// ResetSet Close all properties.
const ResetSet = "\x1b[0m"

// CodeExpr regex to clear color codes eg "\033[1;36mText\x1b[0m"
const CodeExpr = `\033\[[\d;?]+m`

var (
	// Enable switch color render and display
	//
	// NOTICE:
	// if ENV: NO_COLOR is not empty, will disable color render.
	Enable = os.Getenv("NO_COLOR") == ""
	// RenderTag render HTML tag on call color.Xprint, color.PrintX
	RenderTag = true
	// debug mode for development.
	//
	// set env:
	// 	COLOR_DEBUG_MODE=on
	// or:
	// 	COLOR_DEBUG_MODE=on go run ./_examples/envcheck.go
	debugMode = os.Getenv("COLOR_DEBUG_MODE") == "on"
	// inner errors record on detect color level
	innerErrs []error
	// output the default io.Writer message print
	output io.Writer = os.Stdout
	// mark current env, It's like in `cmd.exe`
	// if not in windows, it's always is False.
	isLikeInCmd bool
	// the color support level for current terminal
	// needVTP - need enable VTP, only for windows OS
	colorLevel, needVTP = detectTermColorLevel()
	// match color codes
	codeRegex = regexp.MustCompile(CodeExpr)
	// mark current env is support color.
	// Always: isLikeInCmd != supportColor
	// supportColor = IsSupportColor()
)

// TermColorLevel value on current ENV
func TermColorLevel() terminfo.ColorLevel {
	return colorLevel
}

// SupportColor on the current ENV
func SupportColor() bool {
	return colorLevel > terminfo.ColorLevelNone
}

// Support16Color on the current ENV
// func Support16Color() bool {
// 	return colorLevel > terminfo.ColorLevelNone
// }

// Support256Color on the current ENV
func Support256Color() bool {
	return colorLevel > terminfo.ColorLevelBasic
}

// SupportTrueColor on the current ENV
func SupportTrueColor() bool {
	return colorLevel > terminfo.ColorLevelHundreds
}

/*************************************************************
 * global settings
 *************************************************************/

// Set set console color attributes
func Set(colors ...Color) (int, error) {
	code := Colors2code(colors...)
	err := SetTerminal(code)
	return 0, err
}

// Reset reset console color attributes
func Reset() (int, error) {
	err := ResetTerminal()
	return 0, err
}

// Disable disable color output
func Disable() bool {
	oldVal := Enable
	Enable = false
	return oldVal
}

// NotRenderTag on call color.Xprint, color.PrintX
func NotRenderTag() {
	RenderTag = false
}

// SetOutput set default colored text output
func SetOutput(w io.Writer) {
	output = w
}

// ResetOutput reset output
func ResetOutput() {
	output = os.Stdout
}

// ResetOptions reset all package option setting
func ResetOptions() {
	RenderTag = true
	Enable = true
	output = os.Stdout
}

// ForceColor force open color render
func ForceSetColorLevel(level terminfo.ColorLevel) terminfo.ColorLevel {
	oldLevelVal := colorLevel
	colorLevel = level
	return oldLevelVal
}

// ForceColor force open color render
func ForceColor() terminfo.ColorLevel {
	return ForceOpenColor()
}

// ForceOpenColor force open color render
func ForceOpenColor() terminfo.ColorLevel {
	// TODO should set level to ?
	return ForceSetColorLevel(terminfo.ColorLevelMillions)
}

// IsLikeInCmd check result
// Deprecated
func IsLikeInCmd() bool {
	return isLikeInCmd
}

// InnerErrs info
func InnerErrs() []error {
	return innerErrs
}

/*************************************************************
 * render color code
 *************************************************************/

// RenderCode render message by color code.
// Usage:
// 	msg := RenderCode("3;32;45", "some", "message")
func RenderCode(code string, args ...interface{}) string {
	var message string
	if ln := len(args); ln == 0 {
		return ""
	}

	message = fmt.Sprint(args...)
	if len(code) == 0 {
		return message
	}

	// disabled OR not support color
	if !Enable || !SupportColor() {
		return ClearCode(message)
	}

	return fmt.Sprintf(FullColorTpl, code, message)
}

// RenderWithSpaces Render code with spaces.
// If the number of args is > 1, a space will be added between the args
func RenderWithSpaces(code string, args ...interface{}) string {
	message := formatArgsForPrintln(args)
	if len(code) == 0 {
		return message
	}

	// disabled OR not support color
	if !Enable || !SupportColor() {
		return ClearCode(message)
	}

	return fmt.Sprintf(FullColorTpl, code, message)
}

// RenderString render a string with color code.
// Usage:
// 	msg := RenderString("3;32;45", "a message")
func RenderString(code string, str string) string {
	if len(code) == 0 || str == "" {
		return str
	}

	// disabled OR not support color
	if !Enable || !SupportColor() {
		return ClearCode(str)
	}

	return fmt.Sprintf(FullColorTpl, code, str)
}

// ClearCode clear color codes.
// eg: "\033[36;1mText\x1b[0m" -> "Text"
func ClearCode(str string) string {
	return codeRegex.ReplaceAllString(str, "")
}
