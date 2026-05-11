package style

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/gookit/color"
	"github.com/stretchr/testify/assert"
	"github.com/xo/terminfo"
)

func TestMerge(t *testing.T) {
	type scenario struct {
		name          string
		toMerge       []TextStyle
		expectedStyle TextStyle
		expectedStr   string
	}

	fgRed := color.FgRed
	bgRed := color.BgRed
	fgBlue := color.FgBlue

	rgbPinkLib := color.Rgb(0xFF, 0x00, 0xFF)
	rgbPink := NewRGBColor(rgbPinkLib)

	rgbYellowLib := color.Rgb(0xFF, 0xFF, 0x00)
	rgbYellow := NewRGBColor(rgbYellowLib)

	strToPrint := "foo"

	scenarios := []scenario{
		{
			"no color",
			nil,
			TextStyle{Style: color.Style{}},
			"foo",
		},
		{
			"only fg color",
			[]TextStyle{FgRed},
			TextStyle{fg: &Color{basic: &fgRed}, Style: color.Style{fgRed}},
			"\x1b[31mfoo\x1b[0m",
		},
		{
			"only bg color",
			[]TextStyle{BgRed},
			TextStyle{bg: &Color{basic: &bgRed}, Style: color.Style{bgRed}},
			"\x1b[41mfoo\x1b[0m",
		},
		{
			"fg and bg color",
			[]TextStyle{FgBlue, BgRed},
			TextStyle{
				fg:    &Color{basic: &fgBlue},
				bg:    &Color{basic: &bgRed},
				Style: color.Style{fgBlue, bgRed},
			},
			"\x1b[34;41mfoo\x1b[0m",
		},
		{
			"single attribute",
			[]TextStyle{AttrBold},
			TextStyle{
				decoration: Decoration{bold: true},
				Style:      color.Style{color.OpBold},
			},
			"\x1b[1mfoo\x1b[0m",
		},
		{
			"multiple attributes",
			[]TextStyle{AttrBold, AttrUnderline},
			TextStyle{
				decoration: Decoration{
					bold:      true,
					underline: true,
				},
				Style: color.Style{color.OpBold, color.OpUnderscore},
			},
			"\x1b[1;4mfoo\x1b[0m",
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
				Style: color.Style{fgBlue, bgRed, color.OpBold, color.OpUnderscore},
			},
			"\x1b[34;41;1;4mfoo\x1b[0m",
		},
		{
			"rgb fg color",
			[]TextStyle{New().SetFg(rgbPink)},
			TextStyle{
				fg:    &rgbPink,
				Style: color.NewRGBStyle(rgbPinkLib).SetOpts(color.Opts{}),
			},
			// '38;2' qualifies an RGB foreground color
			"\x1b[38;2;255;0;255mfoo\x1b[0m",
		},
		{
			"rgb fg and bg color",
			[]TextStyle{New().SetFg(rgbPink).SetBg(rgbYellow)},
			TextStyle{
				fg:    &rgbPink,
				bg:    &rgbYellow,
				Style: color.NewRGBStyle(rgbPinkLib, rgbYellowLib).SetOpts(color.Opts{}),
			},
			// '48;2' qualifies an RGB background color
			"\x1b[38;2;255;0;255;48;2;255;255;0mfoo\x1b[0m",
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
				Style: color.NewRGBStyle(rgbPinkLib, rgbYellowLib).SetOpts(color.Opts{color.OpBold, color.OpUnderscore}),
			},
			"\x1b[38;2;255;0;255;48;2;255;255;0;1;4mfoo\x1b[0m",
		},
		{
			"mix color-16 (background) with rgb (foreground)",
			[]TextStyle{New().SetFg(rgbYellow), BgRed},
			TextStyle{
				fg: &rgbYellow,
				bg: &Color{basic: &bgRed},
				Style: color.NewRGBStyle(
					rgbYellowLib,
					fgRed.RGB(), // We need to use FG here,  https://github.com/gookit/color/issues/39
				).SetOpts(color.Opts{}),
			},
			"\x1b[38;2;255;255;0;48;2;197;30;20mfoo\x1b[0m",
		},
		{
			"mix color-16 (foreground) with rgb (background)",
			[]TextStyle{FgRed, New().SetBg(rgbYellow)},
			TextStyle{
				fg: &Color{basic: &fgRed},
				bg: &rgbYellow,
				Style: color.NewRGBStyle(
					fgRed.RGB(),
					rgbYellowLib,
				).SetOpts(color.Opts{}),
			},
			"\x1b[38;2;197;30;20;48;2;255;255;0mfoo\x1b[0m",
		},
	}

	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelMillions)
	defer color.ForceSetColorLevel(oldColorLevel)

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			style := New()
			for _, other := range s.toMerge {
				style = style.MergeStyle(other)
			}
			assert.Equal(t, s.expectedStyle, style)
			assert.Equal(t, s.expectedStr, style.Sprint(strToPrint))
		})
	}
}

func TestTemplateFuncMapAddColors(t *testing.T) {
	type scenario struct {
		name   string
		tmpl   string
		expect string
	}

	scenarios := []scenario{
		{
			"normal template",
			"{{ .Foo }}",
			"bar",
		},
		{
			"colored string",
			"{{ .Foo | red }}",
			"\x1b[31mbar\x1b[0m",
		},
		{
			"string with decorator",
			"{{ .Foo | bold }}",
			"\x1b[1mbar\x1b[0m",
		},
		{
			"string with color and decorator",
			"{{ .Foo | bold | red }}",
			"\x1b[31m\x1b[1mbar\x1b[0m\x1b[0m",
		},
		{
			"multiple string with different colors",
			"{{ .Foo | red }} - {{ .Foo | blue }}",
			"\x1b[31mbar\x1b[0m - \x1b[34mbar\x1b[0m",
		},
	}

	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelMillions)
	defer color.ForceSetColorLevel(oldColorLevel)

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			tmpl, err := template.New("test template").Funcs(TemplateFuncMapAddColors(template.FuncMap{})).Parse(s.tmpl)
			assert.NoError(t, err)

			buff := bytes.NewBuffer(nil)
			err = tmpl.Execute(buff, struct{ Foo string }{"bar"})
			assert.NoError(t, err)

			assert.Equal(t, s.expect, buff.String())
		})
	}
}
