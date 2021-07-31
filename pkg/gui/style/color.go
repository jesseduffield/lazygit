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

func (c *Color) IsRGB() bool {
	return c.rgb != nil
}

func (c *Color) ToRGB() Color {
	if c.IsRGB() {
		return *c
	}

	rgb := c.basic.RGB()
	c.rgb = &rgb

	return NewRGBColor(rgb)
}
