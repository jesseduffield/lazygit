package gui

// A window refers to a place on the screen which can hold one or more views.
// A view is a box that renders content, and within a window only one view will
// appear at a time. When a view appears within a window, it occupies the whole
// space. Right now most windows are 1:1 with views, except for commitFiles which
// is a view belonging to the 'commits' window, alongside the 'commits' view.

func (gui *Gui) getViewNameForWindow(window string) string {
	viewName, ok := gui.State.WindowViewNameMap[window]
	if !ok {
		return window
	}

	return viewName
}

func (gui *Gui) getWindowForViewName(viewName string) string {
	// should soft-code this
	if viewName == "commitFiles" {
		return "commits"
	}

	return viewName
}

func (gui *Gui) setViewAsActiveForWindow(viewName string) {
	if gui.State.WindowViewNameMap == nil {
		gui.State.WindowViewNameMap = map[string]string{}
	}

	gui.State.WindowViewNameMap[gui.getWindowForViewName(viewName)] = viewName
}

func (gui *Gui) currentWindow() string {
	return gui.getWindowForViewName(gui.currentViewName())
}
