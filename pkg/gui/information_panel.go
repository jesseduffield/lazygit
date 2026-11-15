package gui

import (
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) informationStr() string {
	if activeMode, ok := gui.helpers.Mode.GetActiveMode(); ok {
		return activeMode.InfoLabel()
	}

	return gui.Config.GetVersion()
}

func (gui *Gui) handleInfoClick() error {
	if !gui.g.Mouse {
		return nil
	}

	view := gui.Views.Information

	cx, _ := view.Cursor()
	width := view.Width()

	if activeMode, ok := gui.helpers.Mode.GetActiveMode(); ok {
		if width-cx > utils.StringWidth(gui.c.Tr.ResetInParentheses) {
			return nil
		}
		return activeMode.Reset()
	}

	return nil
}
