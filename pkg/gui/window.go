package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

// A window refers to a place on the screen which can hold one or more views.
// A view is a box that renders content, and within a window only one view will
// appear at a time. When a view appears within a window, it occupies the whole
// space. Right now most windows are 1:1 with views, except for commitFiles which
// is a view that moves between windows

func (gui *Gui) initialWindowViewNameMap(contextTree *context.ContextTree) *utils.ThreadSafeMap[string, string] {
	result := utils.NewThreadSafeMap[string, string]()

	for _, context := range contextTree.Flatten() {
		result.Set(context.GetWindowName(), context.GetViewName())
	}

	return result
}

func (gui *Gui) getViewNameForWindow(window string) string {
	viewName, ok := gui.State.WindowViewNameMap.Get(window)
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

	gui.State.WindowViewNameMap.Set(c.GetWindowName(), c.GetViewName())
}

func (gui *Gui) currentWindow() string {
	return gui.currentContext().GetWindowName()
}

// assumes the context's windowName has been set to the new window if necessary
func (gui *Gui) resetWindowContext(c types.Context) {
	for _, windowName := range gui.State.WindowViewNameMap.Keys() {
		viewName, ok := gui.State.WindowViewNameMap.Get(windowName)
		if !ok {
			continue
		}
		if viewName == c.GetViewName() && windowName != c.GetWindowName() {
			for _, context := range gui.State.Contexts.Flatten() {
				if context.GetKey() != c.GetKey() && context.GetWindowName() == windowName {
					gui.State.WindowViewNameMap.Set(windowName, context.GetViewName())
				}
			}
		}
	}
}

// moves given context's view to the top of the window
func (gui *Gui) moveToTopOfWindow(context types.Context) {
	view := context.GetView()
	if view == nil {
		return
	}

	window := context.GetWindowName()

	topView := gui.topViewInWindow(window)

	if view.Name() != topView.Name() {
		if err := gui.g.SetViewOnTopOf(view.Name(), topView.Name()); err != nil {
			gui.Log.Error(err)
		}
	}
}

func (gui *Gui) topViewInWindow(windowName string) *gocui.View {
	// now I need to find all views in that same window, via contexts. And I guess then I need to find the index of the highest view in that list.
	viewNamesInWindow := gui.viewNamesInWindow(windowName)

	// The views list is ordered highest-last, so we're grabbing the last view of the window
	var topView *gocui.View
	for _, currentView := range gui.g.Views() {
		if lo.Contains(viewNamesInWindow, currentView.Name()) {
			topView = currentView
		}
	}

	return topView
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
