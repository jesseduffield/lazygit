package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/constants"
)

func (gui *Gui) informationStr() string {
	for _, mode := range gui.modeStatuses() {
		if mode.isActive() {
			return mode.description()
		}
	}

	if gui.g.Mouse {
		donate := color.New(color.FgMagenta, color.Underline).Sprint(gui.Tr.Donate)
		askQuestion := color.New(color.FgYellow, color.Underline).Sprint(gui.Tr.AskQuestion)
		return fmt.Sprintf("%s %s %s", donate, askQuestion, gui.Config.GetVersion())
	} else {
		return gui.Config.GetVersion()
	}
}

func (gui *Gui) handleInfoClick() error {
	if !gui.g.Mouse {
		return nil
	}

	view := gui.Views.Information

	cx, _ := view.Cursor()
	width, _ := view.Size()

	for _, mode := range gui.modeStatuses() {
		if mode.isActive() {
			if width-cx > len(gui.Tr.ResetInParentheses) {
				return nil
			}
			return mode.reset()
		}
	}

	// if we're not in an active mode we show the donate button
	if cx <= len(gui.Tr.Donate) {
		return gui.OS.OpenLink(constants.Links.Donate)
	} else if cx <= len(gui.Tr.Donate)+1+len(gui.Tr.AskQuestion) {
		return gui.OS.OpenLink(constants.Links.Discussions)
	}
	return nil
}
