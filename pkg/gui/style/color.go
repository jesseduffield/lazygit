package style

import "github.com/gookit/color"

type Color struct {
	rgb   *color.RGBColor
	basic *color.Color
}

func NewRGBColor(cl color.RGBColor) Color {
	c := Color{}
	c.rgb = &cl
	return c
}

func NewBasicColor(cl color.Color) Color {
	c := Color{}
	c.basic = &cl
	return c
}

func (c Color) IsRGB() bool {
	return c.rgb != nil
}

func (c Color) ToRGB(isBg bool) Color {
	if c.IsRGB() {
		return c
	}

	if isBg {
		// We need to convert bg color to fg color
		// This is a gookit/color bug,
		// https://github.com/gookit/color/issues/39
		return NewRGBColor((*c.basic - 10).RGB())
	}

	return NewRGBColor(c.basic.RGB())
}
