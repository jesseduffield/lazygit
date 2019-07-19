package gui

import (
	"github.com/jesseduffield/gocui"
)

// GetAttribute gets the gocui color attribute from the string
func (gui *Gui) GetAttribute(key string) gocui.Attribute {
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
func (gui *Gui) GetColor(keys []string) gocui.Attribute {
	var attribute gocui.Attribute
	for _, key := range keys {
		attribute |= gui.GetAttribute(key)
	}
	return attribute
}

// GetOptionsPanelTextColor gets the color of the options panel text
func (gui *Gui) GetOptionsPanelTextColor() (gocui.Attribute, error) {
	userConfig := gui.Config.GetUserConfig()
	optionsColor := userConfig.GetStringSlice("gui.theme.optionsTextColor")
	return gui.GetColor(optionsColor), nil
}

// SetColorScheme sets the color scheme for the app based on the user config
func (gui *Gui) SetColorScheme() error {
	userConfig := gui.Config.GetUserConfig()
	activeBorderColor := userConfig.GetStringSlice("gui.theme.activeBorderColor")
	inactiveBorderColor := userConfig.GetStringSlice("gui.theme.inactiveBorderColor")
	gui.g.FgColor = gui.GetColor(inactiveBorderColor)
	gui.g.SelFgColor = gui.GetColor(activeBorderColor)
	return nil
}
