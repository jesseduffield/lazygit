// +build !windows

// The method in the file has no effect
// Only for compatibility with non-Windows systems

package color

import (
	"strings"
	"syscall"

	"github.com/xo/terminfo"
)

// detect special term color support
func detectSpecialTermColor(termVal string) (terminfo.ColorLevel, bool) {
	if termVal == "" {
		return terminfo.ColorLevelNone, false
	}

	debugf("terminfo check fail - fallback detect color by check TERM value")

	// on TERM=screen:
	// - support 256, not support true-color. test on macOS
	if termVal == "screen" {
		return terminfo.ColorLevelHundreds, false
	}

	if strings.Contains(termVal, "256color") {
		return terminfo.ColorLevelHundreds, false
	}

	if strings.Contains(termVal, "xterm") {
		return terminfo.ColorLevelHundreds, false
		// return terminfo.ColorLevelBasic, false
	}

	// return terminfo.ColorLevelNone, nil
	return terminfo.ColorLevelBasic, false
}

// IsTerminal returns true if the given file descriptor is a terminal.
//
// Usage:
// 	IsTerminal(os.Stdout.Fd())
func IsTerminal(fd uintptr) bool {
	return fd == uintptr(syscall.Stdout) || fd == uintptr(syscall.Stdin) || fd == uintptr(syscall.Stderr)
}
