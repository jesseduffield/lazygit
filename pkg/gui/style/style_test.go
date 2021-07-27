package style

import (
	"testing"

	"github.com/gookit/color"
	"github.com/stretchr/testify/assert"
)

func TestNewStyle(t *testing.T) {
	type scenario struct {
		name          string
		fg, bg        color.Color
		expectedStyle color.Style
	}

	scenarios := []scenario{
		{
			"no color",
			0, 0,
			color.Style{},
		},
		{
			"only fg color",
			color.FgRed, 0,
			color.Style{color.FgRed},
		},
		{
			"only bg color",
			0, color.BgRed,
			color.Style{color.BgRed},
		},
		{
			"fg and bg color",
			color.FgBlue, color.BgRed,
			color.Style{color.FgBlue, color.BgRed},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			style := New(s.fg, s.bg)
			basicStyle, ok := style.(BasicTextStyle)
			assert.True(t, ok, "New(..) should return a interface of type BasicTextStyle")
			assert.Equal(t, s.fg, basicStyle.fg)
			assert.Equal(t, s.bg, basicStyle.bg)
			assert.Equal(t, []color.Color(nil), basicStyle.opts)
			assert.Equal(t, s.expectedStyle, basicStyle.style)
		})
	}
}

func TestBasicSetColor(t *testing.T) {
	type scenario struct {
		name       string
		colorToSet BasicTextStyle
		expect     BasicTextStyle
	}

	scenarios := []scenario{
		{
			"empty color",
			BasicTextStyle{},
			BasicTextStyle{fg: color.FgRed, bg: color.BgBlue, opts: []color.Color{color.OpBold}}},
		{
			"set new fg color",
			BasicTextStyle{fg: color.FgCyan},
			BasicTextStyle{fg: color.FgCyan, bg: color.BgBlue, opts: []color.Color{color.OpBold}},
		},
		{
			"set new bg color",
			BasicTextStyle{bg: color.BgGray},
			BasicTextStyle{fg: color.FgRed, bg: color.BgGray, opts: []color.Color{color.OpBold}},
		},
		{
			"set new fg and bg color",
			BasicTextStyle{fg: color.FgCyan, bg: color.BgGray},
			BasicTextStyle{fg: color.FgCyan, bg: color.BgGray, opts: []color.Color{color.OpBold}},
		},
		{
			"add options",
			BasicTextStyle{opts: []color.Color{color.OpUnderscore}},
			BasicTextStyle{fg: color.FgRed, bg: color.BgBlue, opts: []color.Color{color.OpBold, color.OpUnderscore}},
		},
		{
			"add options that already exists",
			BasicTextStyle{opts: []color.Color{color.OpBold}},
			BasicTextStyle{fg: color.FgRed, bg: color.BgBlue, opts: []color.Color{color.OpBold}},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			style, ok := New(color.FgRed, color.BgBlue).
				SetBold(true).
				SetColor(s.colorToSet).(BasicTextStyle)
			assert.True(t, ok, "SetColor should return a interface of type BasicTextStyle if the input was also BasicTextStyle")

			style.style = nil
			assert.Equal(t, s.expect, style)
		})
	}
}

func TestRGBSetColor(t *testing.T) {
	type scenario struct {
		name       string
		colorToSet TextStyle
		expect     RGBTextStyle
	}

	red := color.FgRed.RGB()
	cyan := color.FgCyan.RGB()
	blue := color.FgBlue.RGB()
	gray := color.FgGray.RGB()

	toBg := func(c color.RGBColor) *color.RGBColor {
		c[3] = 1
		return &c
	}

	scenarios := []scenario{
		{
			"empty RGBTextStyle input",
			RGBTextStyle{},
			RGBTextStyle{fgSet: true, fg: red, bg: toBg(blue), opts: []color.Color{color.OpBold}},
		},
		{
			"empty BasicTextStyle input",
			BasicTextStyle{},
			RGBTextStyle{fgSet: true, fg: red, bg: toBg(blue), opts: []color.Color{color.OpBold}},
		},
		{
			"set fg and bg color using BasicTextStyle",
			BasicTextStyle{fg: color.FgCyan, bg: color.BgGray},
			RGBTextStyle{fgSet: true, fg: cyan, bg: toBg(gray), opts: []color.Color{color.OpBold}},
		},
		{
			"set fg and bg color using RGBTextStyle",
			RGBTextStyle{fgSet: true, fg: cyan, bg: toBg(gray)},
			RGBTextStyle{fgSet: true, fg: cyan, bg: toBg(gray), opts: []color.Color{color.OpBold}},
		},
		{
			"add options",
			RGBTextStyle{opts: []color.Color{color.OpUnderscore}},
			RGBTextStyle{fgSet: true, fg: red, bg: toBg(blue), opts: []color.Color{color.OpBold, color.OpUnderscore}},
		},
		{
			"add options using BasicTextStyle",
			BasicTextStyle{opts: []color.Color{color.OpUnderscore}},
			RGBTextStyle{fgSet: true, fg: red, bg: toBg(blue), opts: []color.Color{color.OpBold, color.OpUnderscore}},
		},
		{
			"add options that already exists",
			RGBTextStyle{opts: []color.Color{color.OpBold}},
			RGBTextStyle{fgSet: true, fg: red, bg: toBg(blue), opts: []color.Color{color.OpBold}},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			style, ok := New(color.FgRed, color.BgBlue).SetBold(true).(BasicTextStyle)
			assert.True(t, ok, "SetBold should return a interface of type BasicTextStyle")

			rgbStyle, ok := style.convertToRGB().SetColor(s.colorToSet).(RGBTextStyle)
			assert.True(t, ok, "SetColor should return a interface of type RGBTextColor")

			rgbStyle.style = color.RGBStyle{}
			assert.Equal(t, s.expect, rgbStyle)
		})
	}
}

func TestConvertBasicToRGB(t *testing.T) {
	type scenario struct {
		name string
		test func(*testing.T)
	}

	scenarios := []scenario{
		{
			"convert to rgb with fg",
			func(t *testing.T) {
				basicStyle, ok := New(color.FgRed, 0).(BasicTextStyle)
				assert.True(t, ok, "New(..) should return a interface of type BasicTextStyle")

				rgbStyle := basicStyle.convertToRGB()
				assert.True(t, rgbStyle.fgSet)
				assert.Equal(t, color.RGB(197, 30, 20), rgbStyle.fg)
				assert.Nil(t, rgbStyle.bg)
			},
		},
		{
			"convert to rgb with fg and bg",
			func(t *testing.T) {
				basicStyle, ok := New(color.FgRed, color.BgRed).(BasicTextStyle)
				assert.True(t, ok, "New(..) should return a interface of type BasicTextStyle")

				rgbStyle := basicStyle.convertToRGB()
				assert.True(t, rgbStyle.fgSet)
				assert.Equal(t, color.RGB(197, 30, 20), rgbStyle.fg)
				assert.Equal(t, color.RGB(197, 30, 20, true), *rgbStyle.bg)
			},
		},
		{
			"convert to rgb using SetRGBColor",
			func(t *testing.T) {
				style := New(color.FgRed, 0)
				rgbStyle, ok := style.SetRGBColor(255, 00, 255, true).(RGBTextStyle)
				assert.True(t, ok, "SetRGBColor should return a interface of type RGBTextStyle")

				assert.True(t, rgbStyle.fgSet)
				assert.Equal(t, color.RGB(197, 30, 20), rgbStyle.fg)
				assert.Equal(t, color.RGB(255, 0, 255, true), *rgbStyle.bg)
			},
		},
		{
			"convert to rgb using SetRGBColor multiple times",
			func(t *testing.T) {
				style := New(color.FgRed, 0)
				rgbStyle, ok := style.SetRGBColor(00, 255, 255, false).SetRGBColor(255, 00, 255, true).(RGBTextStyle)
				assert.True(t, ok, "SetRGBColor should return a interface of type RGBTextStyle")

				assert.True(t, rgbStyle.fgSet)
				assert.Equal(t, color.RGB(0, 255, 255), rgbStyle.fg)
				assert.Equal(t, color.RGB(255, 0, 255, true), *rgbStyle.bg)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, s.test)
	}
}

func TestSettingAtributes(t *testing.T) {
	type scenario struct {
		name         string
		test         func(s TextStyle) TextStyle
		expectedOpts []color.Color
	}

	scenarios := []scenario{
		{
			"no attributes",
			func(s TextStyle) TextStyle {
				return s
			},
			[]color.Color{},
		},
		{
			"set single attribute",
			func(s TextStyle) TextStyle {
				return s.SetBold(true)
			},
			[]color.Color{color.OpBold},
		},
		{
			"set multiple attributes",
			func(s TextStyle) TextStyle {
				return s.SetBold(true).SetUnderline(true)
			},
			[]color.Color{color.OpBold, color.OpUnderscore},
		},
		{
			"unset a attributes",
			func(s TextStyle) TextStyle {
				return s.SetBold(true).SetBold(false)
			},
			[]color.Color{},
		},
		{
			"unset a attributes with multiple attributes",
			func(s TextStyle) TextStyle {
				return s.SetBold(true).SetUnderline(true).SetBold(false)
			},
			[]color.Color{color.OpUnderscore},
		},
		{
			"unset all attributes with multiple attributes",
			func(s TextStyle) TextStyle {
				return s.SetBold(true).SetUnderline(true).SetBold(false).SetUnderline(false)
			},
			[]color.Color{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			// Test basic style
			style := New(color.FgRed, 0)
			basicStyle, ok := style.(BasicTextStyle)
			assert.True(t, ok, "New(..) should return a interface of type BasicTextStyle")
			basicStyle, ok = s.test(basicStyle).(BasicTextStyle)
			assert.True(t, ok, "underlaying type should not be changed after test")
			assert.Len(t, basicStyle.opts, len(s.expectedOpts))
			for _, opt := range basicStyle.opts {
				assert.Contains(t, s.expectedOpts, opt)
			}
			for _, opt := range s.expectedOpts {
				assert.Contains(t, basicStyle.style, opt)
			}

			// Test RGB style
			rgbStyle := New(color.FgRed, 0).(BasicTextStyle).convertToRGB()
			rgbStyle, ok = s.test(rgbStyle).(RGBTextStyle)
			assert.True(t, ok, "underlaying type should not be changed after test")
			assert.Len(t, rgbStyle.opts, len(s.expectedOpts))
			for _, opt := range rgbStyle.opts {
				assert.Contains(t, s.expectedOpts, opt)
			}
		})
	}
}
