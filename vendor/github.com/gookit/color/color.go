/*
Package color is command line color library.
Support rich color rendering output, universal API method, compatible with Windows system

Source code and other details for the project are available at GitHub:

	https://github.com/gookit/color

For more usage, please see README and tests.
*/
package color

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// color render templates
//
// ESC 操作的表示:
//
//	"\033"(Octal 8进制) = "\x1b"(Hexadecimal 16进制) = 27 (10进制)
const (
	// StartSet chars
	StartSet = "\x1b["
	// ResetSet close all properties.
	ResetSet = "\x1b[0m"
	// SettingTpl string.
	SettingTpl = "\x1b[%sm"
	// FullColorTpl for build color code
	FullColorTpl = "\x1b[%sm%s\x1b[0m"
	// CodeSuffix string for color code.
	CodeSuffix = "[0m"
)

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
)

var (
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
	// the color support level for current terminal
	// needVTP - need enable VTP, only for Windows OS
	colorLevel, needVTP = detectTermColorLevel()
	// match color codes
	codeRegex = regexp.MustCompile(CodeExpr)
)

// TermColorLevel Get the currently supported color level
func TermColorLevel() Level { return colorLevel }

// SupportColor Whether the current environment supports color output
func SupportColor() bool { return colorLevel > LevelNo }

// Support256Color Whether the current environment supports 256-color output
func Support256Color() bool { return colorLevel > Level16 }

// SupportTrueColor Whether the current environment supports (RGB)True-color output
func SupportTrueColor() bool { return colorLevel > Level256 }

/*************************************************************
 * global settings
 *************************************************************/

// Set console color attributes
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
func NotRenderTag() { RenderTag = false }

// SetOutput set default colored text output
func SetOutput(w io.Writer) { output = w }

// ResetOutput reset output
func ResetOutput() { output = os.Stdout }

// ResetOptions reset all package option setting
func ResetOptions() {
	RenderTag = true
	Enable = true
	output = os.Stdout
}

// ForceSetColorLevel force open color render
func ForceSetColorLevel(level Level) Level {
	oldLevelVal := colorLevel
	colorLevel = level
	return oldLevelVal
}

// ForceColor force open color render
func ForceColor() Level { return ForceOpenColor() }

// ForceOpenColor force open color render
func ForceOpenColor() Level {
	// TODO should set level to ?
	return ForceSetColorLevel(LevelRgb)
}

// EnableDebug enable debug mode
func EnableDebug() { debugMode = true }

// ResetDebug reset debug mode
func ResetDebug() { debugMode = false }

// InnerErrs info
func InnerErrs() []error { return innerErrs }

/*************************************************************
 * render color code
 *************************************************************/

// RenderCode render message by color code.
//
// Usage:
//
//	msg := RenderCode("3;32;45", "some", "message")
func RenderCode(code string, args ...any) string {
	var message string

	// Fast path optimizations
	if ln := len(args); ln == 1 {
		// Single argument - avoid fmt.Sprint overhead
		if str, ok := args[0].(string); ok {
			message = str
		} else {
			message = fmt.Sprint(args[0])
		}
	} else if ln == 2 {
		// Two arguments - common case, try to optimize if both are strings
		if str1, ok1 := args[0].(string); ok1 {
			if str2, ok2 := args[1].(string); ok2 {
				message = str1 + str2
			} else {
				message = fmt.Sprint(args...)
			}
		} else {
			message = fmt.Sprint(args...)
		}
	} else if ln == 0 {
		return ""
	} else {
		// Multiple arguments - use fmt.Sprint for safety
		message = fmt.Sprint(args...)
	}

	if len(code) == 0 {
		return message
	}

	// disabled OR not support color
	if !Enable || !SupportColor() {
		return ClearCode(message)
	}

	// return fmt.Sprintf(FullColorTpl, code, message)
	return StartSet + code + "m" + message + ResetSet
}

// RenderWithSpaces Render code with spaces.
// If the number of args is > 1, a space will be added between the args
func RenderWithSpaces(code string, args ...any) string {
	msg := formatLikePrintln(args)
	if len(code) == 0 {
		return msg
	}

	// disabled OR not support color
	if !Enable || !SupportColor() {
		return ClearCode(msg)
	}

	return StartSet + code + "m" + msg + ResetSet
}

// RenderString render a string with color code.
//
// Usage:
//
//	msg := RenderString("3;32;45", "a message")
func RenderString(code string, str string) string {
	if len(code) == 0 || str == "" {
		return str
	}

	// disabled OR not support color
	if !Enable || !SupportColor() {
		return ClearCode(str)
	}

	// return fmt.Sprintf(FullColorTpl, code, str)
	return StartSet + code + "m" + str + ResetSet
}

// ClearCode clear color codes.
//
// eg:
//
//	"\033[36;1mText\x1b[0m" -> "Text"
func ClearCode(str string) string {
	if !strings.Contains(str, CodeSuffix) {
		return str
	}
	return codeRegex.ReplaceAllString(str, "")
}
