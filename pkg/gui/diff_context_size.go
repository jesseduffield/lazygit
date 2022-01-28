package gui

import (
	"errors"
)

var CONTEXT_KEYS_SHOWING_DIFFS = []ContextKey{
	FILES_CONTEXT_KEY,
	COMMIT_FILES_CONTEXT_KEY,
	STASH_CONTEXT_KEY,
	BRANCH_COMMITS_CONTEXT_KEY,
	SUB_COMMITS_CONTEXT_KEY,
	MAIN_STAGING_CONTEXT_KEY,
	MAIN_PATCH_BUILDING_CONTEXT_KEY,
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
			return gui.PopupHandler.Error(err)
		}

		gui.UserConfig.Git.DiffContextSize = gui.UserConfig.Git.DiffContextSize + 1
		return gui.handleDiffContextSizeChange()
	}

	return nil
}

func (gui *Gui) DecreaseContextInDiffView() error {
	old_size := gui.UserConfig.Git.DiffContextSize

	if isShowingDiff(gui) && old_size > 1 {
		if err := gui.CheckCanChangeContext(); err != nil {
			return gui.PopupHandler.Error(err)
		}

		gui.UserConfig.Git.DiffContextSize = old_size - 1
		return gui.handleDiffContextSizeChange()
	}

	return nil
}

func (gui *Gui) handleDiffContextSizeChange() error {
	currentContext := gui.currentStaticContext()
	switch currentContext.GetKey() {
	// we make an exception for our staging and patch building contexts because they actually need to refresh their state afterwards.
	case MAIN_PATCH_BUILDING_CONTEXT_KEY:
		return gui.handleRefreshPatchBuildingPanel(-1)
	case MAIN_STAGING_CONTEXT_KEY:
		return gui.handleRefreshStagingPanel(false, -1)
	default:
		return currentContext.HandleRenderToMain()
	}
}

func (gui *Gui) CheckCanChangeContext() error {
	if gui.Git.Patch.PatchManager.Active() {
		return errors.New(gui.Tr.CantChangeContextSizeError)
	}

	return nil
}
