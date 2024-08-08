package helpers

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type WindowHelper struct {
	c          *HelperCommon
	viewHelper *ViewHelper
}

func NewWindowHelper(c *HelperCommon, viewHelper *ViewHelper) *WindowHelper {
	return &WindowHelper{
		c:          c,
		viewHelper: viewHelper,
	}
}

// A window refers to a place on the screen which can hold one or more views.
// A view is a box that renders content, and within a window only one view will
// appear at a time. When a view appears within a window, it occupies the whole
// space. Right now most windows are 1:1 with views, except for commitFiles which
// is a view that moves between windows

func (self *WindowHelper) GetViewNameForWindow(window string) string {
	viewName, ok := self.windowViewNameMap().Get(window)
	if !ok {
		panic(fmt.Sprintf("Viewname not found for window: %s", window))
	}

	return viewName
}

func (self *WindowHelper) GetContextForWindow(window string) types.Context {
	viewName := self.GetViewNameForWindow(window)

	context, ok := self.viewHelper.ContextForView(viewName)
	if !ok {
		panic("TODO: fix this")
	}

	return context
}

// for now all we actually care about is the context's view so we're storing that
func (self *WindowHelper) SetWindowContext(c types.Context) {
	if c.IsTransient() {
		self.resetWindowContext(c)
	}

	self.windowViewNameMap().Set(c.GetWindowName(), c.GetViewName())
}

func (self *WindowHelper) windowViewNameMap() *utils.ThreadSafeMap[string, string] {
	return self.c.State().GetRepoState().GetWindowViewNameMap()
}

func (self *WindowHelper) CurrentWindow() string {
	return self.c.Context().Current().GetWindowName()
}

// assumes the context's windowName has been set to the new window if necessary
func (self *WindowHelper) resetWindowContext(c types.Context) {
	for _, windowName := range self.windowViewNameMap().Keys() {
		viewName, ok := self.windowViewNameMap().Get(windowName)
		if !ok {
			continue
		}
		if viewName == c.GetViewName() && windowName != c.GetWindowName() {
			for _, context := range self.c.Contexts().Flatten() {
				if context.GetKey() != c.GetKey() && context.GetWindowName() == windowName {
					self.windowViewNameMap().Set(windowName, context.GetViewName())
				}
			}
		}
	}
}

// moves given context's view to the top of the window
func (self *WindowHelper) MoveToTopOfWindow(context types.Context) {
	view := context.GetView()
	if view == nil {
		return
	}

	window := context.GetWindowName()

	topView := self.TopViewInWindow(window, true)

	if topView != nil && view.Name() != topView.Name() {
		if err := self.c.GocuiGui().SetViewOnTopOf(view.Name(), topView.Name()); err != nil {
			self.c.Log.Error(err)
		}
	}
}

func (self *WindowHelper) TopViewInWindow(windowName string, includeInvisibleViews bool) *gocui.View {
	// now I need to find all views in that same window, via contexts. And I guess then I need to find the index of the highest view in that list.
	viewNamesInWindow := self.viewNamesInWindow(windowName)

	// The views list is ordered highest-last, so we're grabbing the last view of the window
	var topView *gocui.View
	for _, currentView := range self.c.GocuiGui().Views() {
		if lo.Contains(viewNamesInWindow, currentView.Name()) && (currentView.Visible || includeInvisibleViews) {
			topView = currentView
		}
	}

	return topView
}

func (self *WindowHelper) viewNamesInWindow(windowName string) []string {
	result := []string{}
	for _, context := range self.c.Contexts().Flatten() {
		if context.GetWindowName() == windowName {
			result = append(result, context.GetViewName())
		}
	}

	return result
}

func (self *WindowHelper) WindowForView(viewName string) string {
	context, ok := self.viewHelper.ContextForView(viewName)
	if !ok {
		panic("todo: deal with this")
	}

	return context.GetWindowName()
}

func (self *WindowHelper) SideWindows() []string {
	return []string{"status", "files", "branches", "commits", "stash"}
}
