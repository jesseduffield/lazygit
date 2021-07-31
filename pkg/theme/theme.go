package theme

import (
	"github.com/gookit/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

var (
	// GocuiDefaultTextColor does the same as DefaultTextColor but this one only colors gocui default text colors
	GocuiDefaultTextColor gocui.Attribute

	// ActiveBorderColor is the border color of the active frame
	ActiveBorderColor gocui.Attribute

	// InactiveBorderColor is the border color of the inactive active frames
	InactiveBorderColor gocui.Attribute

	// GocuiSelectedLineBgColor is the background color for the selected line in gocui
	GocuiSelectedLineBgColor gocui.Attribute

	OptionsColor gocui.Attribute

	// DefaultTextColor is the default text color
	DefaultTextColor = style.FgWhite

	// DefaultHiTextColor is the default highlighted text color
	DefaultHiTextColor = style.FgLightWhite

	// SelectedLineBgColor is the background color for the selected line
	SelectedLineBgColor = style.New()

	// SelectedRangeBgColor is the background color of the selected range of lines
	SelectedRangeBgColor = style.New()

	OptionsFgColor = style.New()

	DiffTerminalColor = style.FgMagenta
)

// UpdateTheme updates all theme variables
func UpdateTheme(themeConfig config.ThemeConfig) {
	ActiveBorderColor = GetGocuiStyle(themeConfig.ActiveBorderColor)
	InactiveBorderColor = GetGocuiStyle(themeConfig.InactiveBorderColor)
	SelectedLineBgColor = GetTextStyle(themeConfig.SelectedLineBgColor, true)
	SelectedRangeBgColor = GetTextStyle(themeConfig.SelectedRangeBgColor, true)
	GocuiSelectedLineBgColor = GetGocuiStyle(themeConfig.SelectedLineBgColor)
	OptionsColor = GetGocuiStyle(themeConfig.OptionsTextColor)
	OptionsFgColor = GetTextStyle(themeConfig.OptionsTextColor, false)

	isLightTheme := themeConfig.LightTheme
	if isLightTheme {
		DefaultTextColor = style.FgBlack
		DefaultHiTextColor = style.FgBlackLighter
		GocuiDefaultTextColor = gocui.ColorBlack
	} else {
		DefaultTextColor = style.FgWhite
		DefaultHiTextColor = style.FgLightWhite
		GocuiDefaultTextColor = gocui.ColorWhite
	}
}

// GetAttribute gets the gocui color attribute from the string
func GetGocuiAttribute(key string) gocui.Attribute {
	if utils.IsValidHexValue(key) {
		values := color.HEX(key).Values()
		return gocui.NewRGBColor(int32(values[0]), int32(values[1]), int32(values[2]))
	}

	colorMap := map[string]gocui.Attribute{
		"default":   gocui.ColorDefault,
		"black":     gocui.ColorBlack,
		"red":       gocui.ColorRed,
		"green":     gocui.ColorGreen,
		"yellow":    gocui.ColorYellow,
		"blue":      gocui.ColorBlue,
		"magenta":   gocui.ColorMagenta,
		"cyan":      gocui.ColorCyan,
		"white":     gocui.ColorWhite,
		"bold":      gocui.AttrBold,
		"reverse":   gocui.AttrReverse,
		"underline": gocui.AttrUnderline,
	}
	value, present := colorMap[key]
	if present {
		return value
	}
	return gocui.ColorWhite
}

// GetGocuiStyle bitwise OR's a list of attributes obtained via the given keys
func GetGocuiStyle(keys []string) gocui.Attribute {
	var attribute gocui.Attribute
	for _, key := range keys {
		attribute |= GetGocuiAttribute(key)
	}
	return attribute
}

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
