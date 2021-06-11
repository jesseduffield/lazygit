package gui

func (gui *Gui) toggleWhitespaceInDiffView() error {
	return gui.setIgnoreWhitespaceFlag(!gui.State.IgnoreWhitespaceInDiffView)
}

func (gui *Gui) setIgnoreWhitespaceFlag(shouldIgnoreWhitespace bool) error {
	if gui.State.IgnoreWhitespaceInDiffView == shouldIgnoreWhitespace {
		return nil
	}

	gui.State.IgnoreWhitespaceInDiffView = shouldIgnoreWhitespace

	toastMessage := gui.Tr.ShowingWhitespaceInDiffView
	if gui.State.IgnoreWhitespaceInDiffView {
		toastMessage = gui.Tr.IgnoringWhitespaceInDiffView
	}
	gui.raiseToast(toastMessage)

	return gui.refreshFilesAndSubmodules()
}
