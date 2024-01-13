package termenv

import (
	"errors"

	"github.com/mattn/go-isatty"
)

var (
	// ErrStatusReport gets returned when the terminal can't be queried.
	ErrStatusReport = errors.New("unable to retrieve status report")
)

const (
	// Control Sequence Introducer
	CSI = "\x1b["
	// Operating System Command
	OSC = "\x1b]"
)

func (o *Output) isTTY() bool {
	if len(o.environ.Getenv("CI")) > 0 {
		return false
	}
	if o.TTY() == nil {
		return false
	}

	return isatty.IsTerminal(o.TTY().Fd())
}

// ColorProfile returns the supported color profile:
// Ascii, ANSI, ANSI256, or TrueColor.
func ColorProfile() Profile {
	return output.ColorProfile()
}

// ForegroundColor returns the terminal's default foreground color.
func ForegroundColor() Color {
	return output.ForegroundColor()
}

// BackgroundColor returns the terminal's default background color.
func BackgroundColor() Color {
	return output.BackgroundColor()
}

// HasDarkBackground returns whether terminal uses a dark-ish background.
func HasDarkBackground() bool {
	return output.HasDarkBackground()
}

// EnvNoColor returns true if the environment variables explicitly disable color output
// by setting NO_COLOR (https://no-color.org/)
// or CLICOLOR/CLICOLOR_FORCE (https://bixense.com/clicolors/)
// If NO_COLOR is set, this will return true, ignoring CLICOLOR/CLICOLOR_FORCE
// If CLICOLOR=="0", it will be true only if CLICOLOR_FORCE is also "0" or is unset.
func (o *Output) EnvNoColor() bool {
	return o.environ.Getenv("NO_COLOR") != "" || (o.environ.Getenv("CLICOLOR") == "0" && !o.cliColorForced())
}

// EnvNoColor returns true if the environment variables explicitly disable color output
// by setting NO_COLOR (https://no-color.org/)
// or CLICOLOR/CLICOLOR_FORCE (https://bixense.com/clicolors/)
// If NO_COLOR is set, this will return true, ignoring CLICOLOR/CLICOLOR_FORCE
// If CLICOLOR=="0", it will be true only if CLICOLOR_FORCE is also "0" or is unset.
func EnvNoColor() bool {
	return output.EnvNoColor()
}

// EnvColorProfile returns the color profile based on environment variables set
// Supports NO_COLOR (https://no-color.org/)
// and CLICOLOR/CLICOLOR_FORCE (https://bixense.com/clicolors/)
// If none of these environment variables are set, this behaves the same as ColorProfile()
// It will return the Ascii color profile if EnvNoColor() returns true
// If the terminal does not support any colors, but CLICOLOR_FORCE is set and not "0"
// then the ANSI color profile will be returned.
func EnvColorProfile() Profile {
	return output.EnvColorProfile()
}

// EnvNoColor returns true if the environment variables explicitly disable color output
// by setting NO_COLOR (https://no-color.org/)
// or CLICOLOR/CLICOLOR_FORCE (https://bixense.com/clicolors/)
// If NO_COLOR is set, this will return true, ignoring CLICOLOR/CLICOLOR_FORCE
// If CLICOLOR=="0", it will be true only if CLICOLOR_FORCE is also "0" or is unset.
func (o *Output) EnvColorProfile() Profile {
	if o.EnvNoColor() {
		return Ascii
	}
	p := o.ColorProfile()
	if o.cliColorForced() && p == Ascii {
		return ANSI
	}
	return p
}

func (o *Output) cliColorForced() bool {
	if forced := o.environ.Getenv("CLICOLOR_FORCE"); forced != "" {
		return forced != "0"
	}
	return false
}
