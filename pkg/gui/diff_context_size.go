package gui

import (
	"errors"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

var CONTEXT_KEYS_SHOWING_DIFFS = []types.ContextKey{
	context.FILES_CONTEXT_KEY,
	context.COMMIT_FILES_CONTEXT_KEY,
	context.STASH_CONTEXT_KEY,
	context.LOCAL_COMMITS_CONTEXT_KEY,
	context.SUB_COMMITS_CONTEXT_KEY,
	context.MAIN_STAGING_CONTEXT_KEY,
	context.MAIN_PATCH_BUILDING_CONTEXT_KEY,
}

func isShowingDiff(gui *Gui) bool {
	key := gui.currentStaticContext().GetKey()

	for _, contextKey := range CONTEXT_KEYS_SHOWING_DIFFS {
		if key == contextKey {
			return true
		}
	}
	return false
}

func (gui *Gui) IncreaseContextInDiffView() error {
	if isShowingDiff(gui) {
		if err := gui.CheckCanChangeContext(); err != nil {
			return gui.c.Error(err)
		}

		gui.c.UserConfig.Git.DiffContextSize = gui.c.UserConfig.Git.DiffContextSize + 1
		return gui.handleDiffContextSizeChange()
	}

	return nil
}

func (gui *Gui) DecreaseContextInDiffView() error {
	old_size := gui.c.UserConfig.Git.DiffContextSize

	if isShowingDiff(gui) && old_size > 1 {
		if err := gui.CheckCanChangeContext(); err != nil {
			return gui.c.Error(err)
		}

		gui.c.UserConfig.Git.DiffContextSize = old_size - 1
		return gui.handleDiffContextSizeChange()
	}

	return nil
}

func (gui *Gui) handleDiffContextSizeChange() error {
	currentContext := gui.currentStaticContext()
	switch currentContext.GetKey() {
	// we make an exception for our staging and patch building contexts because they actually need to refresh their state afterwards.
	case context.MAIN_PATCH_BUILDING_CONTEXT_KEY:
		return gui.handleRefreshPatchBuildingPanel(-1)
	case context.MAIN_STAGING_CONTEXT_KEY:
		return gui.handleRefreshStagingPanel(false, -1)
	default:
		return currentContext.HandleRenderToMain()
	}
}

func (gui *Gui) CheckCanChangeContext() error {
	if gui.git.Patch.PatchManager.Active() {
		return errors.New(gui.c.Tr.CantChangeContextSizeError)
	}

	return nil
}
