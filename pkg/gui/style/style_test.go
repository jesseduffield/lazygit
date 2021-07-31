package style

import (
	"testing"

	"github.com/gookit/color"
	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	type scenario struct {
		name          string
		toMerge       []TextStyle
		expectedStyle TextStyle
	}

	fgRed := color.FgRed
	bgRed := color.BgRed
	fgBlue := color.FgBlue

	rgbPinkLib := color.Rgb(0xFF, 0x00, 0xFF)
	rgbPink := NewRGBColor(rgbPinkLib)

	rgbYellowLib := color.Rgb(0xFF, 0xFF, 0x00)
	rgbYellow := NewRGBColor(rgbYellowLib)

	scenarios := []scenario{
		{
			"no color",
			nil,
			TextStyle{style: color.Style{}},
		},
		{
			"only fg color",
			[]TextStyle{FgRed},
			TextStyle{fg: &Color{basic: &fgRed}, style: color.Style{fgRed}},
		},
		{
			"only bg color",
			[]TextStyle{BgRed},
			TextStyle{bg: &Color{basic: &bgRed}, style: color.Style{bgRed}},
		},
		{
			"fg and bg color",
			[]TextStyle{FgBlue, BgRed},
			TextStyle{
				fg:    &Color{basic: &fgBlue},
				bg:    &Color{basic: &bgRed},
				style: color.Style{fgBlue, bgRed},
			},
		},
		{
			"single attribute",
			[]TextStyle{AttrBold},
			TextStyle{
				decoration: Decoration{bold: true},
				style:      color.Style{color.OpBold},
			},
		},
		{
			"multiple attributes",
			[]TextStyle{AttrBold, AttrUnderline},
			TextStyle{
				decoration: Decoration{
					bold:      true,
					underline: true,
				},
				style: color.Style{color.OpBold, color.OpUnderscore},
			},
		},
		{
			"multiple attributes and colors",
			[]TextStyle{AttrBold, FgBlue, AttrUnderline, BgRed},
			TextStyle{
				fg: &Color{basic: &fgBlue},
				bg: &Color{basic: &bgRed},
				decoration: Decoration{
					bold:      true,
					underline: true,
				},
				style: color.Style{fgBlue, bgRed, color.OpBold, color.OpUnderscore},
			},
		},
		{
			"rgb fg color",
			[]TextStyle{New().SetFg(rgbPink)},
			TextStyle{
				fg:    &rgbPink,
				style: color.NewRGBStyle(rgbPinkLib).SetOpts(color.Opts{}),
			},
		},
		{
			"rgb fg and bg color",
			[]TextStyle{New().SetFg(rgbPink).SetBg(rgbYellow)},
			TextStyle{
				fg:    &rgbPink,
				bg:    &rgbYellow,
				style: color.NewRGBStyle(rgbPinkLib, rgbYellowLib).SetOpts(color.Opts{}),
			},
		},
		{
			"rgb fg and bg color with opts",
			[]TextStyle{AttrBold, New().SetFg(rgbPink).SetBg(rgbYellow), AttrUnderline},
			TextStyle{
				fg: &rgbPink,
				bg: &rgbYellow,
				decoration: Decoration{
					bold:      true,
					underline: true,
				},
				style: color.NewRGBStyle(rgbPinkLib, rgbYellowLib).SetOpts(color.Opts{color.OpBold, color.OpUnderscore}),
			},
		},
		{
			"mix color-16 with rgb colors",
			[]TextStyle{New().SetFg(rgbYellow), BgRed},
			TextStyle{
				fg: &rgbYellow,
				bg: &Color{basic: &bgRed},
				style: color.NewRGBStyle(
					rgbYellowLib,
					fgRed.RGB(), // We need to use FG here,  https://github.com/gookit/color/issues/39
				).SetOpts(color.Opts{}),
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			style := New()
			for _, other := range s.toMerge {
				style = style.MergeStyle(other)
			}
			assert.Equal(t, s.expectedStyle, style)
		})
	}
}
