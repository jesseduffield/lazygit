package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) validateNotInFilterMode() bool {
	if gui.State.Modes.Filtering.Active() {
		_ = gui.c.Confirm(types.ConfirmOpts{
			Title:         gui.c.Tr.MustExitFilterModeTitle,
			Prompt:        gui.c.Tr.MustExitFilterModePrompt,
			HandleConfirm: gui.exitFilterMode,
		})

		return false
	}
	return true
}

func (gui *Gui) outsideFilterMode(f func() error) func() error {
	return func() error {
		if !gui.validateNotInFilterMode() {
			return nil
		}

		return f()
	}
}

func (gui *Gui) exitFilterMode() error {
	return gui.clearFiltering()
}

func (gui *Gui) clearFiltering() error {
	gui.State.Modes.Filtering.Reset()
	if gui.State.ScreenMode == SCREEN_HALF {
		gui.State.ScreenMode = SCREEN_NORMAL
	}

	return gui.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.COMMITS}})
}

func (gui *Gui) setFiltering(path string) error {
	gui.State.Modes.Filtering.SetPath(path)
	if gui.State.ScreenMode == SCREEN_NORMAL {
		gui.State.ScreenMode = SCREEN_HALF
	}

	if err := gui.c.PushContext(gui.State.Contexts.LocalCommits); err != nil {
		return err
	}

	return gui.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.COMMITS}, Then: func() {
		gui.State.Contexts.LocalCommits.SetSelectedLineIdx(0)
	}})
}
