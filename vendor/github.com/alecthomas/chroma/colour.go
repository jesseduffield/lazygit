package chroma

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ANSI2RGB maps ANSI colour names, as supported by Chroma, to hex RGB values.
var ANSI2RGB = map[string]string{
	"#ansiblack":     "000000",
	"#ansidarkred":   "7f0000",
	"#ansidarkgreen": "007f00",
	"#ansibrown":     "7f7fe0",
	"#ansidarkblue":  "00007f",
	"#ansipurple":    "7f007f",
	"#ansiteal":      "007f7f",
	"#ansilightgray": "e5e5e5",
	// Normal
	"#ansidarkgray":  "555555",
	"#ansired":       "ff0000",
	"#ansigreen":     "00ff00",
	"#ansiyellow":    "ffff00",
	"#ansiblue":      "0000ff",
	"#ansifuchsia":   "ff00ff",
	"#ansiturquoise": "00ffff",
	"#ansiwhite":     "ffffff",

	// Aliases without the "ansi" prefix, because...why?
	"#black":     "000000",
	"#darkred":   "7f0000",
	"#darkgreen": "007f00",
	"#brown":     "7f7fe0",
	"#darkblue":  "00007f",
	"#purple":    "7f007f",
	"#teal":      "007f7f",
	"#lightgray": "e5e5e5",
	// Normal
	"#darkgray":  "555555",
	"#red":       "ff0000",
	"#green":     "00ff00",
	"#yellow":    "ffff00",
	"#blue":      "0000ff",
	"#fuchsia":   "ff00ff",
	"#turquoise": "00ffff",
	"#white":     "ffffff",
}

// Colour represents an RGB colour.
type Colour int32

// NewColour creates a Colour directly from RGB values.
func NewColour(r, g, b uint8) Colour {
	return ParseColour(fmt.Sprintf("%02x%02x%02x", r, g, b))
}

// Distance between this colour and another.
//
// This uses the approach described here (https://www.compuphase.com/cmetric.htm).
// This is not as accurate as LAB, et. al. but is *vastly* simpler and sufficient for our needs.
func (c Colour) Distance(e2 Colour) float64 {
	ar, ag, ab := int64(c.Red()), int64(c.Green()), int64(c.Blue())
	br, bg, bb := int64(e2.Red()), int64(e2.Green()), int64(e2.Blue())
	rmean := (ar + br) / 2
	r := ar - br
	g := ag - bg
	b := ab - bb
	return math.Sqrt(float64((((512 + rmean) * r * r) >> 8) + 4*g*g + (((767 - rmean) * b * b) >> 8)))
}

// Brighten returns a copy of this colour with its brightness adjusted.
//
// If factor is negative, the colour is darkened.
//
// Uses approach described here (http://www.pvladov.com/2012/09/make-color-lighter-or-darker.html).
func (c Colour) Brighten(factor float64) Colour {
	r := float64(c.Red())
	g := float64(c.Green())
	b := float64(c.Blue())

	if factor < 0 {
		factor++
		r *= factor
		g *= factor
		b *= factor
	} else {
		r = (255-r)*factor + r
		g = (255-g)*factor + g
		b = (255-b)*factor + b
	}
	return NewColour(uint8(r), uint8(g), uint8(b))
}

// BrightenOrDarken brightens a colour if it is < 0.5 brighteness or darkens if > 0.5 brightness.
func (c Colour) BrightenOrDarken(factor float64) Colour {
	if c.Brightness() < 0.5 {
		return c.Brighten(factor)
	}
	return c.Brighten(-factor)
}

// Brightness of the colour (roughly) in the range 0.0 to 1.0
func (c Colour) Brightness() float64 {
	return (float64(c.Red()) + float64(c.Green()) + float64(c.Blue())) / 255.0 / 3.0
}

// ParseColour in the forms #rgb, #rrggbb, #ansi<colour>, or #<colour>.
// Will return an "unset" colour if invalid.
func ParseColour(colour string) Colour {
	colour = normaliseColour(colour)
	n, err := strconv.ParseUint(colour, 16, 32)
	if err != nil {
		return 0
	}
	return Colour(n + 1)
}

// MustParseColour is like ParseColour except it panics if the colour is invalid.
//
// Will panic if colour is in an invalid format.
func MustParseColour(colour string) Colour {
	parsed := ParseColour(colour)
	if !parsed.IsSet() {
		panic(fmt.Errorf("invalid colour %q", colour))
	}
	return parsed
}

// IsSet returns true if the colour is set.
func (c Colour) IsSet() bool { return c != 0 }

func (c Colour) String() string   { return fmt.Sprintf("#%06x", int(c-1)) }
func (c Colour) GoString() string { return fmt.Sprintf("Colour(0x%06x)", int(c-1)) }

// Red component of colour.
func (c Colour) Red() uint8 { return uint8(((c - 1) >> 16) & 0xff) }

// Green component of colour.
func (c Colour) Green() uint8 { return uint8(((c - 1) >> 8) & 0xff) }

// Blue component of colour.
func (c Colour) Blue() uint8 { return uint8((c - 1) & 0xff) }

// Colours is an orderable set of colours.
type Colours []Colour

func (c Colours) Len() int           { return len(c) }
func (c Colours) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Colours) Less(i, j int) bool { return c[i] < c[j] }

// Convert colours to #rrggbb.
func normaliseColour(colour string) string {
	if ansi, ok := ANSI2RGB[colour]; ok {
		return ansi
	}
	if strings.HasPrefix(colour, "#") {
		colour = colour[1:]
		if len(colour) == 3 {
			return colour[0:1] + colour[0:1] + colour[1:2] + colour[1:2] + colour[2:3] + colour[2:3]
		}
	}
	return colour
}
