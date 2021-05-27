package gui

import "github.com/jesseduffield/gocui"

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

func (gui *Gui) getWindowForView(view *gocui.View) string {
	if view == gui.Views.CommitFiles {
		return gui.State.Contexts.CommitFiles.GetWindowName()
	}

	return view.Name()
}

func (gui *Gui) setViewAsActiveForWindow(view *gocui.View) {
	if gui.State.WindowViewNameMap == nil {
		gui.State.WindowViewNameMap = map[string]string{}
	}

	gui.State.WindowViewNameMap[gui.getWindowForView(view)] = view.Name()
}

func (gui *Gui) currentWindow() string {
	return gui.getWindowForView(gui.g.CurrentView())
}

func (gui *Gui) resetWindowForView(view *gocui.View) {
	window := gui.getWindowForView(view)
	// we assume here that the window contains as its default view a view with the same name as the window
	gui.State.WindowViewNameMap[window] = window
}
