package theme

import (
	"github.com/gookit/color"
	"github.com/lobes/lazytask/pkg/gui/style"
	"github.com/lobes/lazytask/pkg/utils"
)

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
		case "strikethrough":
			s = s.SetStrikethrough()
		default:
			value, present := style.ColorMap[key]
			if present {
				var c style.TextStyle
				if background {
					c = value.Background
				} else {
					c = value.Foreground
				}
				s = s.MergeStyle(c)
			} else if utils.IsValidHexValue(key) {
				c := style.NewRGBColor(color.HEX(key, background))
				if background {
					s = s.SetBg(c)
				} else {
					s = s.SetFg(c)
				}
			}
		}
	}

	return s
}
