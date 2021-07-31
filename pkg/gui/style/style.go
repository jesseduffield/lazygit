package style

import (
	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

var (
	// FgWhite        = New(pointerTo(color.FgWhite), nil)
	FgWhite        = FromBasicFg(color.FgWhite)
	FgLightWhite   = FromBasicFg(color.FgLightWhite)
	FgBlack        = FromBasicFg(color.FgBlack)
	FgBlackLighter = FromBasicFg(color.FgBlack.Light())
	FgCyan         = FromBasicFg(color.FgCyan)
	FgRed          = FromBasicFg(color.FgRed)
	FgGreen        = FromBasicFg(color.FgGreen)
	FgBlue         = FromBasicFg(color.FgBlue)
	FgYellow       = FromBasicFg(color.FgYellow)
	FgMagenta      = FromBasicFg(color.FgMagenta)

	BgWhite   = FromBasicBg(color.BgWhite)
	BgBlack   = FromBasicBg(color.BgBlack)
	BgRed     = FromBasicBg(color.BgRed)
	BgGreen   = FromBasicBg(color.BgGreen)
	BgYellow  = FromBasicBg(color.BgYellow)
	BgBlue    = FromBasicBg(color.BgBlue)
	BgMagenta = FromBasicBg(color.BgMagenta)
	BgCyan    = FromBasicBg(color.BgCyan)

	AttrUnderline = New().SetUnderline()
	AttrBold      = New().SetBold()
)

func New() TextStyle {
	return TextStyle{}
}

func FromBasicFg(fg color.Color) TextStyle {
	s := New()
	c := NewBasicColor(fg)
	s.fg = &c
	return s
}

func FromBasicBg(bg color.Color) TextStyle {
	s := New()
	c := NewBasicColor(bg)
	s.bg = &c
	return s
}

var colorMap = map[string]struct {
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

func SetConfigStyles(keys []string, background bool) TextStyle {
	s := New()

	for _, key := range keys {
		switch key {
		case "bold":
			s = s.SetBold()
		case "reverse":
			s = s.SetReverse()
		case "underline":
			s = s.SetUnderline()
		default:
			value, present := colorMap[key]
			if present {
				var c TextStyle
				if background {
					c = value.background
				} else {
					c = value.forground
				}
				s = s.MergeStyle(c)
			} else if utils.IsValidHexValue(key) {
				c := NewRGBColor(color.HEX(key, background))
				s.bg = &c
			}
		}
	}

	return s
}
