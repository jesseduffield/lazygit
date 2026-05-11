//go:build !windows
// +build !windows

// The method in the file has no effect
// Only for compatibility with non-Windows systems

package color

import (
	"strings"
	"syscall"
)

// detect special term color support
func detectSpecialTermColor(termVal string) (Level, bool) {
	if termVal == "" {
		return LevelNo, false
	}

	debugf("terminfo check fail - fallback detect color by check TERM value")

	// on TERM=screen:
	// - support 256, not support true-color. test on macOS
	if termVal == "screen" {
		return Level256, false
	}

	if strings.Contains(termVal, "256color") {
		return Level256, false
	}

	if strings.Contains(termVal, "xterm") {
		return Level256, false
		// return terminfo.ColorLevelBasic, false
	}

	// return LevelNo, nil
	return Level16, false
}

// IsTerminal returns true if the given file descriptor is a terminal.
//
// Usage:
//
//	IsTerminal(os.Stdout.Fd())
func IsTerminal(fd uintptr) bool {
	return fd == uintptr(syscall.Stdout) || fd == uintptr(syscall.Stdin) || fd == uintptr(syscall.Stderr)
}
