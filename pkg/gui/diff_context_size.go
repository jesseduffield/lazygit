package gui

import (
	"errors"
)

func isShowingDiff(gui *Gui) bool {
	key := gui.currentStaticContext().GetKey()

	return key == FILES_CONTEXT_KEY || key == COMMIT_FILES_CONTEXT_KEY || key == STASH_CONTEXT_KEY || key == BRANCH_COMMITS_CONTEXT_KEY || key == SUB_COMMITS_CONTEXT_KEY || key == MAIN_STAGING_CONTEXT_KEY || key == MAIN_PATCH_BUILDING_CONTEXT_KEY
}

func (gui *Gui) IncreaseContextInDiffView() error {
	if isShowingDiff(gui) {
		if err := gui.CheckCanChangeContext(); err != nil {
			return gui.surfaceError(err)
		}

		gui.Config.GetUserConfig().Git.DiffContextSize = gui.Config.GetUserConfig().Git.DiffContextSize + 1
		return gui.currentStaticContext().HandleRenderToMain()
	}

	return nil
}

func (gui *Gui) DecreaseContextInDiffView() error {
	old_size := gui.Config.GetUserConfig().Git.DiffContextSize

	if isShowingDiff(gui) && old_size > 1 {
		if err := gui.CheckCanChangeContext(); err != nil {
			return gui.surfaceError(err)
		}

		gui.Config.GetUserConfig().Git.DiffContextSize = old_size - 1
		return gui.currentStaticContext().HandleRenderToMain()
	}

	return nil
}

func (gui *Gui) CheckCanChangeContext() error {
	if gui.GitCommand.PatchManager.Active() {
		return errors.New(gui.Tr.CantChangeContextSizeError)
	}

	return nil
}
