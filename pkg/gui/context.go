package gui

// changeContext is a helper function for when we want to change a 'main' context
// which currently just means a context that affects both the main and secondary views
// other views can have their context changed directly but this function helps
// keep the main and secondary views in sync
func (gui *Gui) changeMainViewsContext(context string) {
	if gui.State.MainContext == context {
		return
	}

	switch context {
	case "normal", "patch-building", "staging", "merging":
		gui.getMainView().Context = context
		gui.getSecondaryView().Context = context
	}

	gui.State.MainContext = context
}
