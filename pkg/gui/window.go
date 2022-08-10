package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// A window refers to a place on the screen which can hold one or more views.
// A view is a box that renders content, and within a window only one view will
// appear at a time. When a view appears within a window, it occupies the whole
// space. Right now most windows are 1:1 with views, except for commitFiles which
// is a view that moves between windows

func (gui *Gui) initialWindowViewNameMap(contextTree *context.ContextTree) map[string]string {
	result := map[string]string{}

	for _, context := range contextTree.Flatten() {
		result[context.GetWindowName()] = context.GetViewName()
	}

	return result
}

func (gui *Gui) getViewNameForWindow(window string) string {
	viewName, ok := gui.State.WindowViewNameMap[window]
	if !ok {
		panic(fmt.Sprintf("Viewname not found for window: %s", window))
	}

	return viewName
}

func (gui *Gui) getContextForWindow(window string) types.Context {
	viewName := gui.getViewNameForWindow(window)

	context, ok := gui.contextForView(viewName)
	if !ok {
		panic("TODO: fix this")
	}

	return context
}

// for now all we actually care about is the context's view so we're storing that
func (gui *Gui) setWindowContext(c types.Context) {
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
			for _, context := range gui.State.Contexts.Flatten() {
				if context.GetKey() != c.GetKey() && context.GetWindowName() == windowName {
					gui.State.WindowViewNameMap[windowName] = context.GetViewName()
				}
			}
		}
	}
}

func (gui *Gui) moveToTopOfWindow(context types.Context) {
	view := context.GetView()
	if view == nil {
		return
	}

	window := context.GetWindowName()

	// now I need to find all views in that same window, via contexts. And I guess then I need to find the index of the highest view in that list.
	viewNamesInWindow := gui.viewNamesInWindow(window)

	// The views list is ordered highest-last, so we're grabbing the last view of the window
	topView := view
	for _, currentView := range gui.g.Views() {
		if lo.Contains(viewNamesInWindow, currentView.Name()) {
			topView = currentView
		}
	}

	if err := gui.g.SetViewOnTopOf(view.Name(), topView.Name()); err != nil {
		gui.Log.Error(err)
	}
}

func (gui *Gui) viewNamesInWindow(windowName string) []string {
	result := []string{}
	for _, context := range gui.State.Contexts.Flatten() {
		if context.GetWindowName() == windowName {
			result = append(result, context.GetViewName())
		}
	}

	return result
}
