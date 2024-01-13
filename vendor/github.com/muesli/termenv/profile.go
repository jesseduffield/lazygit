package termenv

import (
	"image/color"
	"strconv"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

// Profile is a color profile: Ascii, ANSI, ANSI256, or TrueColor.
type Profile int

const (
	// TrueColor, 24-bit color profile
	TrueColor = Profile(iota)
	// ANSI256, 8-bit color profile
	ANSI256
	// ANSI, 4-bit color profile
	ANSI
	// Ascii, uncolored profile
	Ascii //nolint:revive
)

// String returns a new Style.
func (p Profile) String(s ...string) Style {
	return Style{
		profile: p,
		string:  strings.Join(s, " "),
	}
}

// Convert transforms a given Color to a Color supported within the Profile.
func (p Profile) Convert(c Color) Color {
	if p == Ascii {
		return NoColor{}
	}

	switch v := c.(type) {
	case ANSIColor:
		return v

	case ANSI256Color:
		if p == ANSI {
			return ansi256ToANSIColor(v)
		}
		return v

	case RGBColor:
		h, err := colorful.Hex(string(v))
		if err != nil {
			return nil
		}
		if p != TrueColor {
			ac := hexToANSI256Color(h)
			if p == ANSI {
				return ansi256ToANSIColor(ac)
			}
			return ac
		}
		return v
	}

	return c
}

// Color creates a Color from a string. Valid inputs are hex colors, as well as
// ANSI color codes (0-15, 16-255).
func (p Profile) Color(s string) Color {
	if len(s) == 0 {
		return nil
	}

	var c Color
	if strings.HasPrefix(s, "#") {
		c = RGBColor(s)
	} else {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil
		}

		if i < 16 {
			c = ANSIColor(i)
		} else {
			c = ANSI256Color(i)
		}
	}

	return p.Convert(c)
}

// FromColor creates a Color from a color.Color.
func (p Profile) FromColor(c color.Color) Color {
	col, _ := colorful.MakeColor(c)
	return p.Color(col.Hex())
}
