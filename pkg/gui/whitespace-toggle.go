package gui

func (gui *Gui) toggleWhitespaceInDiffView() error {
	gui.IgnoreWhitespaceInDiffView = !gui.IgnoreWhitespaceInDiffView

	toastMessage := gui.c.Tr.ShowingWhitespaceInDiffView
	if gui.IgnoreWhitespaceInDiffView {
		toastMessage = gui.c.Tr.IgnoringWhitespaceInDiffView
	}
	gui.c.Toast(toastMessage)

	return gui.refreshFilesAndSubmodules()
}
