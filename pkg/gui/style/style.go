package style

import (
	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type TextStyle interface {
	Sprint(a ...interface{}) string
	Sprintf(format string, a ...interface{}) string
	SetBold(v bool) TextStyle
	SetReverse(v bool) TextStyle
	SetUnderline(v bool) TextStyle
	SetColor(style TextStyle) TextStyle
	SetRGBColor(r, g, b uint8, background bool) TextStyle
}

var (
	FgWhite        = New(color.FgWhite, 0)
	FgLightWhite   = New(color.FgLightWhite, 0)
	FgBlack        = New(color.FgBlack, 0)
	FgBlackLighter = New(color.FgBlack.Light(), 0)
	FgCyan         = New(color.FgCyan, 0)
	FgRed          = New(color.FgRed, 0)
	FgGreen        = New(color.FgGreen, 0)
	FgBlue         = New(color.FgBlue, 0)
	FgYellow       = New(color.FgYellow, 0)
	FgMagenta      = New(color.FgMagenta, 0)

	BgWhite   = New(0, color.BgWhite)
	BgBlack   = New(0, color.BgBlack)
	BgRed     = New(0, color.BgRed)
	BgGreen   = New(0, color.BgGreen)
	BgYellow  = New(0, color.BgYellow)
	BgBlue    = New(0, color.BgBlue)
	BgMagenta = New(0, color.BgMagenta)
	BgCyan    = New(0, color.BgCyan)

	AttrUnderline = New(0, 0).SetUnderline(true)
	AttrBold      = New(0, 0).SetUnderline(true)
)

func New(fg color.Color, bg color.Color, opts ...color.Color) TextStyle {
	return BasicTextStyle{
		fg:    fg,
		bg:    bg,
		opts:  opts,
		style: color.Style{},
	}.deriveStyle()
}

func SetConfigStyles(s TextStyle, keys []string, background bool) TextStyle {
	for _, key := range keys {
		colorMap := map[string]struct {
			forground  TextStyle
			background TextStyle
		}{
			"default": {FgWhite, BgBlack},
			"black":   {FgBlack, BgBlack},
			"red":     {FgRed, BgRed},
			"green":   {FgGreen, BgGreen},
			"yellow":  {FgYellow, BgYellow},
			"blue":    {FgBlue, BgBlue},
			"magenta": {FgMagenta, BgMagenta},
			"cyan":    {FgCyan, BgCyan},
			"white":   {FgWhite, BgWhite},
		}
		value, present := colorMap[key]
		if present {
			if background {
				s = s.SetColor(value.background)
			} else {
				s = s.SetColor(value.forground)
			}
			continue
		}

		if key == "bold" {
			s = s.SetBold(true)
			continue
		} else if key == "reverse" {
			s = s.SetReverse(true)
			continue
		} else if key == "underline" {
			s = s.SetUnderline(true)
			continue
		}

		r, g, b, validHexColor := utils.GetHexColorValues(key)
		if validHexColor {
			s = s.SetRGBColor(r, g, b, background)
			continue
		}
	}

	return s
}
