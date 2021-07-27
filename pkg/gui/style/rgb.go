package style

import (
	"fmt"

	"github.com/gookit/color"
)

type RGBTextStyle struct {
	opts  color.Opts
	fgSet bool
	fg    color.RGBColor
	bg    *color.RGBColor
	style color.RGBStyle
}

func (b RGBTextStyle) Sprint(a ...interface{}) string {
	return b.style.Sprint(a...)
}

func (b RGBTextStyle) Sprintf(format string, a ...interface{}) string {
	return b.style.Sprintf(format, a...)
}

func (b RGBTextStyle) setOpt(opt color.Color, v bool) RGBTextStyle {
	if v {
		// Add value
		for _, listOpt := range b.opts {
			if listOpt == opt {
				return b
			}
		}
		b.opts = append(b.opts, opt)
	} else {
		// Remove value
		for idx, listOpt := range b.opts {
			if listOpt == opt {
				b.opts = append(b.opts[:idx], b.opts[idx+1:]...)
				return b
			}
		}
	}
	return b
}

func (b RGBTextStyle) SetBold(v bool) TextStyle {
	b = b.setOpt(color.OpBold, v)
	b.style.SetOpts(b.opts)
	return b
}

func (b RGBTextStyle) SetReverse(v bool) TextStyle {
	b = b.setOpt(color.OpReverse, v)
	b.style.SetOpts(b.opts)
	return b
}

func (b RGBTextStyle) SetUnderline(v bool) TextStyle {
	b = b.setOpt(color.OpUnderscore, v)
	b.style.SetOpts(b.opts)
	return b
}

func (b RGBTextStyle) SetColor(style TextStyle) TextStyle {
	var rgbStyle RGBTextStyle

	switch typedStyle := style.(type) {
	case BasicTextStyle:
		rgbStyle = typedStyle.convertToRGB()
	case RGBTextStyle:
		rgbStyle = typedStyle
	default:
		panic(fmt.Sprintf("got %T but expected BasicTextStyle or RGBTextStyle", typedStyle))
	}

	for _, opt := range rgbStyle.GetOpts() {
		b = b.setOpt(opt, true)
	}

	if rgbStyle.fgSet {
		b.fg = rgbStyle.fg
		b.style.SetFg(rgbStyle.fg)
		b.fgSet = true
	}

	if rgbStyle.bg != nil {
		// Making sure to copy value
		b.bg = &color.RGBColor{}
		*b.bg = *rgbStyle.bg
		b.style.SetBg(*rgbStyle.bg)
	}

	return b
}

func (b RGBTextStyle) SetRGBColor(red, green, blue uint8, background bool) TextStyle {
	parsedColor := color.Rgb(red, green, blue, background)
	if background {
		b.bg = &parsedColor
		b.style.SetBg(parsedColor)
	} else {
		b.fg = parsedColor
		b.style.SetFg(parsedColor)
		b.fgSet = true
	}
	return b
}

func (b RGBTextStyle) GetOpts() color.Opts {
	return b.opts
}
