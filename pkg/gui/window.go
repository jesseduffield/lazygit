package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

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

// for now all we actually care about is the context's view so we're storing that
func (gui *Gui) setWindowContext(c types.Context) {
	if gui.State.WindowViewNameMap == nil {
		gui.State.WindowViewNameMap = map[string]string{}
	}

	if c.IsTransient() {
		gui.resetWindowContext(c)
	}

	gui.State.WindowViewNameMap[c.GetWindowName()] = c.GetViewName()
}

func (gui *Gui) currentWindow() string {
	return gui.currentContext().GetWindowName()
}

// assumes the context's windowName has been set to the new window if necessary
func (gui *Gui) resetWindowContext(c types.Context) {
	for windowName, viewName := range gui.State.WindowViewNameMap {
		if viewName == c.GetViewName() && windowName != c.GetWindowName() {
			// we assume here that the window contains as its default view a view with the same name as the window
			gui.State.WindowViewNameMap[windowName] = windowName
		}
	}
}
