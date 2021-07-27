package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

func (gui *Gui) informationStr() string {
	for _, mode := range gui.modeStatuses() {
		if mode.isActive() {
			return mode.description()
		}
	}

	if gui.g.Mouse {
		donate := style.FgMagenta.SetUnderline(true).Sprint(gui.Tr.Donate)
		askQuestion := style.FgYellow.SetUnderline(true).Sprint(gui.Tr.AskQuestion)
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
		return gui.OSCommand.OpenLink(constants.Links.Donate)
	} else if cx <= len(gui.Tr.Donate)+1+len(gui.Tr.AskQuestion) {
		return gui.OSCommand.OpenLink(constants.Links.Discussions)
	}
	return nil
}
