package gui

func (gui *Gui) toggleWhitespaceInDiffView() error {
	gui.State.IgnoreWhitespaceInDiffView = !gui.State.IgnoreWhitespaceInDiffView

	toastMessage := gui.Tr.ShowingWhitespaceInDiffView
	if gui.State.IgnoreWhitespaceInDiffView {
		toastMessage = gui.Tr.IgnoringWhitespaceInDiffView
	}
	gui.raiseToast(toastMessage)

	return gui.refreshFilesAndSubmodules()
}
