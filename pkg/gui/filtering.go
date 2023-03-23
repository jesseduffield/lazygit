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

func (gui *Gui) setFiltering(path string) error {
	gui.State.Modes.Filtering.SetPath(path)
	if gui.State.ScreenMode == types.SCREEN_NORMAL {
		gui.State.ScreenMode = types.SCREEN_HALF
	}

	if err := gui.c.PushContext(gui.State.Contexts.LocalCommits); err != nil {
		return err
	}

	return gui.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.COMMITS}, Then: func() {
		gui.State.Contexts.LocalCommits.SetSelectedLineIdx(0)
	}})
}
