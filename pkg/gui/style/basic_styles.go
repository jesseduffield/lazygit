package style

import (
	"github.com/gookit/color"
)

var (
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

	ColorMap = map[string]struct {
		Foreground TextStyle
		Background TextStyle
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
)

func FromBasicFg(fg color.Color) TextStyle {
	return New().SetFg(NewBasicColor(fg))
}

func FromBasicBg(bg color.Color) TextStyle {
	return New().SetBg(NewBasicColor(bg))
}
