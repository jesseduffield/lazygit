package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) validateNotInFilterMode() bool {
	if gui.State.Modes.Filtering.Active() {
		_ = gui.c.Confirm(types.ConfirmOpts{
			Title:         gui.c.Tr.MustExitFilterModeTitle,
			Prompt:        gui.c.Tr.MustExitFilterModePrompt,
			HandleConfirm: gui.helpers.Mode.ExitFilterMode,
		})

		return false
	}
	return true
}
