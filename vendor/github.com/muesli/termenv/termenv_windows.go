//go:build windows
// +build windows

package termenv

import (
	"fmt"
	"strconv"

	"golang.org/x/sys/windows"
)

func (o *Output) ColorProfile() Profile {
	if !o.isTTY() {
		return Ascii
	}

	if o.environ.Getenv("ConEmuANSI") == "ON" {
		return TrueColor
	}

	winVersion, _, buildNumber := windows.RtlGetNtVersionNumbers()
	if buildNumber < 10586 || winVersion < 10 {
		// No ANSI support before Windows 10 build 10586.
		if o.environ.Getenv("ANSICON") != "" {
			conVersion := o.environ.Getenv("ANSICON_VER")
			cv, err := strconv.ParseInt(conVersion, 10, 64)
			if err != nil || cv < 181 {
				// No 8 bit color support before v1.81 release.
				return ANSI
			}

			return ANSI256
		}

		return Ascii
	}
	if buildNumber < 14931 {
		// No true color support before build 14931.
		return ANSI256
	}

	return TrueColor
}

func (o Output) foregroundColor() Color {
	// default gray
	return ANSIColor(7)
}

func (o Output) backgroundColor() Color {
	// default black
	return ANSIColor(0)
}

// EnableWindowsANSIConsole enables virtual terminal processing on Windows
// platforms. This allows the use of ANSI escape sequences in Windows console
// applications. Ensure this gets called before anything gets rendered with
// termenv.
//
// Returns the original console mode and an error if one occurred.
func EnableWindowsANSIConsole() (uint32, error) {
	handle, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if err != nil {
		return 0, err
	}

	var mode uint32
	err = windows.GetConsoleMode(handle, &mode)
	if err != nil {
		return 0, err
	}

	// See https://docs.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences
	if mode&windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING != windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING {
		vtpmode := mode | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		if err := windows.SetConsoleMode(handle, vtpmode); err != nil {
			return 0, err
		}
	}

	return mode, nil
}

// RestoreWindowsConsole restores the console mode to a previous state.
func RestoreWindowsConsole(mode uint32) error {
	handle, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if err != nil {
		return err
	}

	return windows.SetConsoleMode(handle, mode)
}

// EnableVirtualTerminalProcessing enables virtual terminal processing on
// Windows for o and returns a function that restores o to its previous state.
// On non-Windows platforms, or if o does not refer to a terminal, then it
// returns a non-nil no-op function and no error.
func EnableVirtualTerminalProcessing(o *Output) (restoreFunc func() error, err error) {
	// There is nothing to restore until we set the console mode.
	restoreFunc = func() error {
		return nil
	}

	// If o is not a tty, then there is nothing to do.
	tty := o.TTY()
	if tty == nil {
		return
	}

	// Get the current console mode. If there is an error, assume that o is not
	// a terminal, discard the error, and return.
	var mode uint32
	if err2 := windows.GetConsoleMode(windows.Handle(tty.Fd()), &mode); err2 != nil {
		return
	}

	// If virtual terminal processing is already set, then there is nothing to
	// do and nothing to restore.
	if mode&windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING == windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING {
		return
	}

	// Enable virtual terminal processing. See
	// https://docs.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences
	if err2 := windows.SetConsoleMode(windows.Handle(tty.Fd()), mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING); err2 != nil {
		err = fmt.Errorf("windows.SetConsoleMode: %w", err2)
		return
	}

	// Set the restore function. We maintain a reference to the tty in the
	// closure (rather than just its handle) to ensure that the tty is not
	// closed by a finalizer.
	restoreFunc = func() error {
		return windows.SetConsoleMode(windows.Handle(tty.Fd()), mode)
	}

	return
}
