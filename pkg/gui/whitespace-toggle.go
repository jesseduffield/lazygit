package gui

func (gui *Gui) toggleWhitespaceInDiffView() error {
	gui.IgnoreWhitespaceInDiffView = !gui.IgnoreWhitespaceInDiffView

	toastMessage := gui.Tr.ShowingWhitespaceInDiffView
	if gui.IgnoreWhitespaceInDiffView {
		toastMessage = gui.Tr.IgnoringWhitespaceInDiffView
	}
	gui.raiseToast(toastMessage)

	return gui.refreshFilesAndSubmodules()
}
