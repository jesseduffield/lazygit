package theme

import (
	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

var colorMap = map[string]struct {
	foreground style.TextStyle
	background style.TextStyle
}{
	"default": {style.FgWhite, style.BgBlack},
	"black":   {style.FgBlack, style.BgBlack},
	"red":     {style.FgRed, style.BgRed},
	"green":   {style.FgGreen, style.BgGreen},
	"yellow":  {style.FgYellow, style.BgYellow},
	"blue":    {style.FgBlue, style.BgBlue},
	"magenta": {style.FgMagenta, style.BgMagenta},
	"cyan":    {style.FgCyan, style.BgCyan},
	"white":   {style.FgWhite, style.BgWhite},
}

func GetTextStyle(keys []string, background bool) style.TextStyle {
	s := style.New()

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
				var c style.TextStyle
				if background {
					c = value.background
				} else {
					c = value.foreground
				}
				s = s.MergeStyle(c)
			} else if utils.IsValidHexValue(key) {
				c := style.NewRGBColor(color.HEX(key, background))
				if background {
					s.SetBg(c)
				} else {
					s.SetFg(c)
				}
			}
		}
	}

	return s
}
