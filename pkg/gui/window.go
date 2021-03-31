package gui

// A window refers to a place on the screen which can hold one or more views.
// A view is a box that renders content, and within a window only one view will
// appear at a time. When a view appears within a window, it occupies the whole
// space. Right now most windows are 1:1 with views, except for commitFiles which
// is a view that moves between windows

func (gui *Gui) getViewNameForWindow(window string) string {
	viewName, ok := gui.State.WindowViewNameMap[window]
	if !ok {
		return window
	}

	return viewName
}

func (gui *Gui) getWindowForViewName(viewName string) string {
	if viewName == "commitFiles" {
		return gui.Contexts.CommitFiles.Context.GetWindowName()
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

func (gui *Gui) resetWindowForView(viewName string) {
	window := gui.getWindowForViewName(viewName)
	// we assume here that the window contains as its default view a view with the same name as the window
	gui.State.WindowViewNameMap[window] = window
}
