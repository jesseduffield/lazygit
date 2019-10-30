package theme

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/spf13/viper"
)

var (
	// DefaultTextColor is the default text color
	DefaultTextColor = color.FgWhite
	// DefaultHiTextColor is the default highlighted text color
	DefaultHiTextColor = color.FgHiWhite

	// GocuiDefaultTextColor does the same as DefaultTextColor but this one only colors gocui default text colors
	GocuiDefaultTextColor gocui.Attribute

	// ActiveBorderColor is the border color of the active frame
	ActiveBorderColor gocui.Attribute

	// InactiveBorderColor is the border color of the inactive active frames
	InactiveBorderColor gocui.Attribute
)

// UpdateTheme updates all theme variables
func UpdateTheme(userConfig *viper.Viper) {
	ActiveBorderColor = getColor(userConfig.GetStringSlice("gui.theme.activeBorderColor"))
	InactiveBorderColor = getColor(userConfig.GetStringSlice("gui.theme.inactiveBorderColor"))

	isLightTheme := userConfig.GetBool("gui.theme.lightTheme")
	if isLightTheme {
		DefaultTextColor = color.FgBlack
		DefaultHiTextColor = color.FgHiBlack
		GocuiDefaultTextColor = gocui.ColorBlack
	} else {
		DefaultTextColor = color.FgWhite
		DefaultHiTextColor = color.FgHiWhite
		GocuiDefaultTextColor = gocui.ColorWhite
	}
}

// getAttribute gets the gocui color attribute from the string
func getAttribute(key string) gocui.Attribute {
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

// getColor bitwise OR's a list of attributes obtained via the given keys
func getColor(keys []string) gocui.Attribute {
	var attribute gocui.Attribute
	for _, key := range keys {
		attribute |= getAttribute(key)
	}
	return attribute
}

// GetAttribute gets the gocui color attribute from the string
func GetAttribute(key string) gocui.Attribute {
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

// GetColor bitwise OR's a list of attributes obtained via the given keys
func GetColor(keys []string) gocui.Attribute {
	var attribute gocui.Attribute
	for _, key := range keys {
		attribute |= GetAttribute(key)
	}
	return attribute
}
