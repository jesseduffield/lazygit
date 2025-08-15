// Copyright 2020 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import "github.com/gdamore/tcell/v2"

// Attribute affects the presentation of characters, such as color, boldness, etc.
type Attribute uint64

const (
	// ColorDefault is used to leave the Color unchanged from whatever system or teminal default may exist.
	ColorDefault = Attribute(tcell.ColorDefault)

	// AttrIsValidColor is used to indicate the color value is actually
	// valid (initialized).  This is useful to permit the zero value
	// to be treated as the default.
	AttrIsValidColor = Attribute(tcell.ColorValid)

	// AttrIsRGBColor is used to indicate that the Attribute value is RGB value of color.
	// The lower order 3 bytes are RGB.
	// (It's not a color in basic ANSI range 256).
	AttrIsRGBColor = Attribute(tcell.ColorIsRGB)

	// AttrColorBits is a mask where color is located in Attribute
	AttrColorBits = 0xffffffffff // roughly 5 bytes, tcell uses 4 bytes and half-byte as a special flags for color (rest is reserved for future)

	// AttrStyleBits is a mask where character attributes (e.g.: bold, italic, underline) are located in Attribute
	AttrStyleBits = 0xffffff0000000000 // remaining 3 bytes in the 8 bytes Attribute (tcell is not using it, so we should be fine)
)

// Color attributes. These colors are compatible with tcell.Color type and can be expanded like:
//
//	g.FgColor := gocui.Attribute(tcell.ColorLime)
const (
	ColorBlack Attribute = AttrIsValidColor + iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

// grayscale indexes (for backward compatibility with termbox-go original grayscale)
var grayscale = []tcell.Color{
	16, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244,
	245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255, 231,
}

// Attributes are not colors, but effects (e.g.: bold, dim) which affect the display of text.
// They can be combined.
const (
	AttrBold Attribute = 1 << (40 + iota)
	AttrBlink
	AttrReverse
	AttrUnderline
	AttrDim
	AttrItalic
	AttrStrikeThrough
	AttrNone Attribute = 0 // Just normal text.
)

// AttrAll represents all the text effect attributes turned on
const AttrAll = AttrBold | AttrBlink | AttrReverse | AttrUnderline | AttrDim | AttrItalic

// IsValidColor indicates if the Attribute is a valid color value (has been set).
func (a Attribute) IsValidColor() bool {
	return a&AttrIsValidColor != 0
}

// Hex returns the color's hexadecimal RGB 24-bit value with each component
// consisting of a single byte, ala R << 16 | G << 8 | B.  If the color
// is unknown or unset, -1 is returned.
//
// This function produce the same output as `tcell.Hex()` with additional
// support for `termbox-go` colors (to 256).
func (a Attribute) Hex() int32 {
	if !a.IsValidColor() {
		return -1
	}
	tc := getTcellColor(a, OutputTrue)
	return tc.Hex()
}

// RGB returns the red, green, and blue components of the color, with
// each component represented as a value 0-255.  If the color
// is unknown or unset, -1 is returned for each component.
//
// This function produce the same output as `tcell.RGB()` with additional
// support for `termbox-go` colors (to 256).
func (a Attribute) RGB() (int32, int32, int32) {
	v := a.Hex()
	if v < 0 {
		return -1, -1, -1
	}
	return (v >> 16) & 0xff, (v >> 8) & 0xff, v & 0xff
}

// GetColor creates a Color from a color name (W3C name). A hex value may
// be supplied as a string in the format "#ffffff".
func GetColor(color string) Attribute {
	return Attribute(tcell.GetColor(color))
}

// Get256Color creates Attribute which stores ANSI color (0-255)
func Get256Color(color int32) Attribute {
	return Attribute(color) | AttrIsValidColor
}

// GetRGBColor creates Attribute which stores RGB color.
// Color is passed as 24bit RGB value, where R << 16 | G << 8 | B
func GetRGBColor(color int32) Attribute {
	return Attribute(color) | AttrIsValidColor | AttrIsRGBColor
}

// NewRGBColor creates Attribute which stores RGB color.
func NewRGBColor(r, g, b int32) Attribute {
	return Attribute(tcell.NewRGBColor(r, g, b))
}

// getTcellColor transform  Attribute into tcell.Color
func getTcellColor(c Attribute, omode OutputMode) tcell.Color {
	c = c & AttrColorBits
	// Default color is 0 in tcell/v2 and was 0 in termbox-go, so we are good here
	if c == ColorDefault {
		return tcell.ColorDefault
	}

	tc := tcell.ColorDefault
	// Check if we have valid color
	if c.IsValidColor() {
		tc = tcell.Color(c)
	} else if c > 0 && c <= 256 {
		// It's not valid color, but it has value in range 1-256
		// This is old Attribute style of color from termbox-go (black=1, etc.)
		// convert to tcell color (black=0|ColorValid)
		tc = tcell.Color(c-1) | tcell.ColorValid
	}

	switch omode {
	case OutputTrue:
		return tc
	case OutputNormal:
		tc &= tcell.Color(0xf) | tcell.ColorValid
	case Output256:
		tc &= tcell.Color(0xff) | tcell.ColorValid
	case Output216:
		tc &= tcell.Color(0xff)
		if tc > 215 {
			return tcell.ColorDefault
		}
		tc += tcell.Color(16) | tcell.ColorValid
	case OutputGrayscale:
		tc &= tcell.Color(0x1f)
		if tc > 26 {
			return tcell.ColorDefault
		}
		tc = grayscale[tc] | tcell.ColorValid
	default:
		return tcell.ColorDefault
	}
	return tc
}
