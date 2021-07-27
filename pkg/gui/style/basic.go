package style

import (
	"fmt"

	"github.com/gookit/color"
)

type BasicTextStyle struct {
	fg   color.Color
	bg   color.Color
	opts []color.Color

	style color.Style
}

func (b BasicTextStyle) Sprint(a ...interface{}) string {
	return b.style.Sprint(a...)
}

func (b BasicTextStyle) Sprintf(format string, a ...interface{}) string {
	return b.style.Sprintf(format, a...)
}

func (b BasicTextStyle) deriveStyle() BasicTextStyle {
	// b.style[:0] makes sure to use the same slice memory
	if b.fg == 0 {
		// Fg is most of the time defined so we reverse the check
		b.style = b.style[:0]
	} else {
		b.style = append(b.style[:0], b.fg)
	}

	if b.bg != 0 {
		b.style = append(b.style, b.bg)
	}

	b.style = append(b.style, b.opts...)
	return b
}

func (b BasicTextStyle) setOpt(opt color.Color, v bool, deriveIfChanged bool) BasicTextStyle {
	if v {
		// Add value
		for _, listOpt := range b.opts {
			if listOpt == opt {
				// Option already added
				return b
			}
		}

		b.opts = append(b.opts, opt)
	} else {
		// Remove value
		for idx, listOpt := range b.opts {
			if listOpt == opt {
				b.opts = append(b.opts[:idx], b.opts[idx+1:]...)

				if deriveIfChanged {
					return b.deriveStyle()
				}
				return b
			}
		}
	}

	if deriveIfChanged {
		return b.deriveStyle()
	}
	return b
}

func (b BasicTextStyle) SetBold(v bool) TextStyle {
	return b.setOpt(color.OpBold, v, true)
}

func (b BasicTextStyle) SetReverse(v bool) TextStyle {
	return b.setOpt(color.OpReverse, v, true)
}

func (b BasicTextStyle) SetUnderline(v bool) TextStyle {
	return b.setOpt(color.OpUnderscore, v, true)
}

func (b BasicTextStyle) SetRGBColor(red, green, blue uint8, background bool) TextStyle {
	return b.convertToRGB().SetRGBColor(red, green, blue, background)
}

func (b BasicTextStyle) convertToRGB() RGBTextStyle {
	res := RGBTextStyle{
		fg:    b.fg.RGB(),
		fgSet: b.fg != 0,
		opts:  b.opts,
	}

	if b.bg != 0 {
		// Need to convert bg to fg otherwise .RGB wont work
		// for more info see https://github.com/gookit/color/issues/39
		rgbBg := (b.bg - 10).RGB()
		rgbBg[3] = 1
		res.bg = &rgbBg
		res.style = *color.NewRGBStyle(res.fg, rgbBg)
	} else {
		res.style = *color.NewRGBStyle(res.fg)
	}
	res.style.SetOpts(b.opts)

	return res
}

func (b BasicTextStyle) SetColor(other TextStyle) TextStyle {
	switch typedOther := other.(type) {
	case BasicTextStyle:
		if typedOther.fg != 0 {
			b.fg = typedOther.fg
		}
		if typedOther.bg != 0 {
			b.bg = typedOther.bg
		}
		for _, opt := range typedOther.opts {
			b = b.setOpt(opt, true, false)
		}
		return b.deriveStyle()
	case RGBTextStyle:
		bAsRGB := b.convertToRGB()

		for _, opt := range typedOther.opts {
			bAsRGB.setOpt(opt, true)
		}

		if typedOther.fgSet {
			bAsRGB.fg = typedOther.fg
			bAsRGB.style.SetFg(typedOther.fg)
		}

		if typedOther.bg != nil {
			// Making sure to copy the value
			bAsRGB.bg = &color.RGBColor{}
			*bAsRGB.bg = *typedOther.bg
			bAsRGB.style.SetBg(*typedOther.bg)
		}

		return bAsRGB
	default:
		panic(fmt.Sprintf("got %T but expected BasicTextStyle or RGBTextStyle", typedOther))
	}
}
