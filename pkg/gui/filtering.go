package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) validateNotInFilterMode() (bool, error) {
	if gui.State.Modes.Filtering.Active() {
		err := gui.PopupHandler.Ask(popup.AskOpts{
			Title:         gui.Tr.MustExitFilterModeTitle,
			Prompt:        gui.Tr.MustExitFilterModePrompt,
			HandleConfirm: gui.exitFilterMode,
		})

		return false, err
	}
	return true, nil
}

func (gui *Gui) exitFilterMode() error {
	return gui.clearFiltering()
}

func (gui *Gui) clearFiltering() error {
	gui.State.Modes.Filtering.Reset()
	if gui.State.ScreenMode == SCREEN_HALF {
		gui.State.ScreenMode = SCREEN_NORMAL
	}

	return gui.refreshSidePanels(types.RefreshOptions{Scope: []types.RefreshableView{types.COMMITS}})
}

func (gui *Gui) setFiltering(path string) error {
	gui.State.Modes.Filtering.SetPath(path)
	if gui.State.ScreenMode == SCREEN_NORMAL {
		gui.State.ScreenMode = SCREEN_HALF
	}

	if err := gui.pushContext(gui.State.Contexts.BranchCommits); err != nil {
		return err
	}

	return gui.refreshSidePanels(types.RefreshOptions{Scope: []types.RefreshableView{types.COMMITS}, Then: func() {
		gui.State.Contexts.BranchCommits.GetPanelState().SetSelectedLineIdx(0)
	}})
}
