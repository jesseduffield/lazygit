package style

import (
	"github.com/gookit/color"
)

type TextStyle struct {
	fg         *Color
	bg         *Color
	decoration Decoration
}

type Sprinter interface {
	Sprint(a ...interface{}) string
	Sprintf(format string, a ...interface{}) string
}

func (b TextStyle) Sprint(a ...interface{}) string {
	return b.deriveStyle().Sprint(a...)
}

func (b TextStyle) Sprintf(format string, a ...interface{}) string {
	return b.deriveStyle().Sprintf(format, a...)
}

func (b TextStyle) SetBold() TextStyle {
	b.decoration.SetBold()
	return b
}

func (b TextStyle) SetUnderline() TextStyle {
	b.decoration.SetUnderline()
	return b
}

func (b TextStyle) SetReverse() TextStyle {
	b.decoration.SetReverse()
	return b
}

func (b TextStyle) deriveStyle() Sprinter {
	// TODO: consider caching
	return deriveStyle(b.fg, b.bg, b.decoration)
}

func deriveStyle(fg *Color, bg *Color, decoration Decoration) Sprinter {
	if fg == nil && bg == nil {
		return color.Style(decoration.ToOpts())
	}

	isRgb := (fg != nil && fg.IsRGB()) || (bg != nil && bg.IsRGB())
	if isRgb {
		s := &color.RGBStyle{}
		if fg != nil {
			s.SetFg(*fg.ToRGB().rgb)
		}
		if bg != nil {
			s.SetBg(*bg.ToRGB().rgb)
		}
		s.SetOpts(decoration.ToOpts())
		return s
	}

	style := make([]color.Color, 0, 5)

	if fg != nil {
		style = append(style, *fg.basic)
	}

	if bg != nil {
		style = append(style, *bg.basic)
	}

	style = append(style, decoration.ToOpts()...)

	return color.Style(style)
}

// // Need to convert bg to fg otherwise .RGB wont work
// // for more info see https://github.com/gookit/color/issues/39
// rgbBg := (*b.bg - 10).RGB()
// rgbBg[3] = 1
// *res.bg = rgbBg
// res.style = *color.NewRGBStyle(*res.fg, rgbBg)

func (b TextStyle) MergeStyle(other TextStyle) TextStyle {
	b.decoration = b.decoration.Merge(other.decoration)

	if other.fg != nil {
		b.fg = other.fg
	}

	if other.bg != nil {
		b.bg = other.bg
	}

	return b
}
