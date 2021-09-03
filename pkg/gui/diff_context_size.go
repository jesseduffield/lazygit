package gui

func isShowingDiff(gui *Gui) bool {
	key := gui.currentStaticContext().GetKey()

	return key == FILES_CONTEXT_KEY || key == COMMIT_FILES_CONTEXT_KEY || key == STASH_CONTEXT_KEY || key == BRANCH_COMMITS_CONTEXT_KEY || key == SUB_COMMITS_CONTEXT_KEY || key == MAIN_STAGING_CONTEXT_KEY || key == MAIN_PATCH_BUILDING_CONTEXT_KEY
}

func (gui *Gui) IncreaseContextInDiffView() error {
	if isShowingDiff(gui) {
		gui.Config.GetUserConfig().Git.DiffContextSize = gui.Config.GetUserConfig().Git.DiffContextSize + 1
		return gui.postRefreshUpdate(gui.currentStaticContext())
	}

	return nil
}
