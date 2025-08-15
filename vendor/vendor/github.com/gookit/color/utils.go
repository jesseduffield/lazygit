package color

import (
	"fmt"
	"io"
	"log"
	"strings"
)

// SetTerminal by given code.
func SetTerminal(code string) error {
	if !Enable || !SupportColor() {
		return nil
	}

	_, err := fmt.Fprintf(output, SettingTpl, code)
	return err
}

// ResetTerminal terminal setting.
func ResetTerminal() error {
	if !Enable || !SupportColor() {
		return nil
	}

	_, err := fmt.Fprint(output, ResetSet)
	return err
}

/*************************************************************
 * print methods(will auto parse color tags)
 *************************************************************/

// Print render color tag and print messages
func Print(a ...interface{}) {
	Fprint(output, a...)
}

// Printf format and print messages
func Printf(format string, a ...interface{}) {
	Fprintf(output, format, a...)
}

// Println messages with new line
func Println(a ...interface{}) {
	Fprintln(output, a...)
}

// Fprint print rendered messages to writer
// Notice: will ignore print error
func Fprint(w io.Writer, a ...interface{}) {
	_, err := fmt.Fprint(w, Render(a...))
	saveInternalError(err)

	// if isLikeInCmd {
	// 	renderColorCodeOnCmd(func() {
	// 		_, _ = fmt.Fprint(w, Render(a...))
	// 	})
	// } else {
	// 	_, _ = fmt.Fprint(w, Render(a...))
	// }
}

// Fprintf print format and rendered messages to writer.
// Notice: will ignore print error
func Fprintf(w io.Writer, format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	_, err := fmt.Fprint(w, ReplaceTag(str))
	saveInternalError(err)
}

// Fprintln print rendered messages line to writer
// Notice: will ignore print error
func Fprintln(w io.Writer, a ...interface{}) {
	str := formatArgsForPrintln(a)
	_, err := fmt.Fprintln(w, ReplaceTag(str))
	saveInternalError(err)
}

// Lprint passes colored messages to a log.Logger for printing.
// Notice: should be goroutine safe
func Lprint(l *log.Logger, a ...interface{}) {
	l.Print(Render(a...))
}

// Render parse color tags, return rendered string.
// Usage:
//	text := Render("<info>hello</> <cyan>world</>!")
//	fmt.Println(text)
func Render(a ...interface{}) string {
	if len(a) == 0 {
		return ""
	}

	return ReplaceTag(fmt.Sprint(a...))
}

// Sprint parse color tags, return rendered string
func Sprint(a ...interface{}) string {
	if len(a) == 0 {
		return ""
	}

	return ReplaceTag(fmt.Sprint(a...))
}

// Sprintf format and return rendered string
func Sprintf(format string, a ...interface{}) string {
	return ReplaceTag(fmt.Sprintf(format, a...))
}

// String alias of the ReplaceTag
func String(s string) string {
	return ReplaceTag(s)
}

// Text alias of the ReplaceTag
func Text(s string) string {
	return ReplaceTag(s)
}

/*************************************************************
 * helper methods for print
 *************************************************************/

// new implementation, support render full color code on pwsh.exe, cmd.exe
func doPrintV2(code, str string) {
	_, err := fmt.Fprint(output, RenderString(code, str))
	saveInternalError(err)

	// if isLikeInCmd {
	// 	renderColorCodeOnCmd(func() {
	// 		_, _ = fmt.Fprint(output, RenderString(code, str))
	// 	})
	// } else {
	// 	_, _ = fmt.Fprint(output, RenderString(code, str))
	// }
}

// new implementation, support render full color code on pwsh.exe, cmd.exe
func doPrintlnV2(code string, args []interface{}) {
	str := formatArgsForPrintln(args)
	_, err := fmt.Fprintln(output, RenderString(code, str))
	saveInternalError(err)
}

// if use Println, will add spaces for each arg
func formatArgsForPrintln(args []interface{}) (message string) {
	if ln := len(args); ln == 0 {
		message = ""
	} else if ln == 1 {
		message = fmt.Sprint(args[0])
	} else {
		message = fmt.Sprintln(args...)
		// clear last "\n"
		message = message[:len(message)-1]
	}
	return
}

/*************************************************************
 * helper methods
 *************************************************************/

// is on debug mode
// func isDebugMode() bool {
// 	return debugMode == "on"
// }

func debugf(f string, v ...interface{}) {
	if debugMode {
		fmt.Print("COLOR_DEBUG: ")
		fmt.Printf(f, v...)
		fmt.Println()
	}
}

// equals: return ok ? val1 : val2
func compareVal(ok bool, val1, val2 uint8) uint8 {
	if ok {
		return val1
	}
	return val2
}

func saveInternalError(err error) {
	if err != nil {
		debugf("inner error: %s", err.Error())
		innerErrs = append(innerErrs, err)
	}
}

func stringToArr(str, sep string) (arr []string) {
	str = strings.TrimSpace(str)
	if str == "" {
		return
	}

	ss := strings.Split(str, sep)
	for _, val := range ss {
		if val = strings.TrimSpace(val); val != "" {
			arr = append(arr, val)
		}
	}
	return
}
