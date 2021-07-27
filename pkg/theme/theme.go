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
	DefaultTextColor = style.New(color.FgWhite, 0)

	// DefaultHiTextColor is the default highlighted text color
	DefaultHiTextColor = style.New(color.FgLightWhite, 0)

	// SelectedLineBgColor is the background color for the selected line
	SelectedLineBgColor = style.New(0, 0)

	// SelectedRangeBgColor is the background color of the selected range of lines
	SelectedRangeBgColor = style.New(0, 0)

	OptionsFgColor = style.New(0, 0)

	DiffTerminalColor = style.New(color.FgMagenta, 0)
)

// UpdateTheme updates all theme variables
func UpdateTheme(themeConfig config.ThemeConfig) {
	ActiveBorderColor = GetGocuiColor(themeConfig.ActiveBorderColor)
	InactiveBorderColor = GetGocuiColor(themeConfig.InactiveBorderColor)
	SelectedLineBgColor = style.SetConfigStyles(SelectedLineBgColor, themeConfig.SelectedLineBgColor, true)
	SelectedRangeBgColor = style.SetConfigStyles(SelectedRangeBgColor, themeConfig.SelectedRangeBgColor, true)
	GocuiSelectedLineBgColor = GetGocuiColor(themeConfig.SelectedLineBgColor)
	OptionsColor = GetGocuiColor(themeConfig.OptionsTextColor)
	OptionsFgColor = style.SetConfigStyles(OptionsFgColor, themeConfig.OptionsTextColor, false)

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
	r, g, b, validHexColor := utils.GetHexColorValues(key)
	if validHexColor {
		return gocui.NewRGBColor(int32(r), int32(g), int32(b))
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

// GetGocuiColor bitwise OR's a list of attributes obtained via the given keys
func GetGocuiColor(keys []string) gocui.Attribute {
	var attribute gocui.Attribute
	for _, key := range keys {
		attribute |= GetGocuiAttribute(key)
	}
	return attribute
}
