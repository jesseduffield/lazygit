package gui

import (
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/tasks"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/spkg/bom"
)

func (gui *Gui) resetViewOrigin(v *gocui.View) {
	v.SetCursor(0, 0)
	v.SetOrigin(0, 0)
}

// Returns the number of lines that we should read initially from a cmd task so
// that the scrollbar has the correct size, along with the number of lines after
// which the view is filled and we can do a first refresh.
func (gui *Gui) linesToReadFromCmdTask(v *gocui.View) tasks.LinesToRead {
	height := v.InnerHeight()
	oy := v.OriginY()

	linesForFirstRefresh := height + oy + 10

	// We want to read as many lines initially as necessary to let the
	// scrollbar go to its minimum height, so that the scrollbar thumb doesn't
	// change size as you scroll down.
	minScrollbarHeight := 1
	linesToReadForAccurateScrollbar := min(
		// However, cap it at some arbitrary max limit, so that we don't get
		// performance problems for huge monitors or tiny font sizes
		height*(height-1)/minScrollbarHeight+oy, 5000)

	return tasks.LinesToRead{
		Total:               linesToReadForAccurateScrollbar,
		InitialRefreshAfter: linesForFirstRefresh,
	}
}

func (gui *Gui) cleanString(s string) string {
	output := string(bom.Clean([]byte(s)))
	return utils.NormalizeLinefeeds(output)
}

func (gui *Gui) setViewContent(v *gocui.View, s string) {
	v.SetContent(gui.cleanString(s))
}

func (gui *Gui) currentViewName() string {
	currentView := gui.g.CurrentView()
	if currentView == nil {
		return ""
	}
	return currentView.Name()
}

func (gui *Gui) onViewTabClick(windowName string, tabIndex int) error {
	tabs := gui.viewTabMap()[windowName]
	if len(tabs) == 0 {
		return nil
	}

	viewName := tabs[tabIndex].ViewName

	context, ok := gui.helpers.View.ContextForView(viewName)
	if !ok {
		return nil
	}

	gui.c.Context().Push(context, types.OnFocusOpts{})
	return nil
}

func (gui *Gui) handleNextTab() error {
	view := getTabbedView(gui)
	if view == nil {
		return nil
	}

	for _, context := range gui.State.Contexts.Flatten() {
		if context.GetViewName() == view.Name() {
			return gui.onViewTabClick(
				context.GetWindowName(),
				utils.ModuloWithWrap(view.TabIndex+1, len(view.Tabs)),
			)
		}
	}

	return nil
}

func (gui *Gui) handlePrevTab() error {
	view := getTabbedView(gui)
	if view == nil {
		return nil
	}

	for _, context := range gui.State.Contexts.Flatten() {
		if context.GetViewName() == view.Name() {
			return gui.onViewTabClick(
				context.GetWindowName(),
				utils.ModuloWithWrap(view.TabIndex-1, len(view.Tabs)),
			)
		}
	}

	return nil
}

func getTabbedView(gui *Gui) *gocui.View {
	// It safe assumption that only static contexts have tabs
	context := gui.c.Context().CurrentStatic()
	view, _ := gui.g.View(context.GetViewName())
	return view
}

func (gui *Gui) render() {
	gui.c.OnUIThread(func() error { return nil })
}

// postRefreshUpdate is to be called on a context after the state that it depends on has been refreshed
// if the context's view is set to another context we do nothing.
// if the context's view is the current view we trigger a focus; re-selecting the current item.
func (gui *Gui) postRefreshUpdate(c types.Context) {
	t := time.Now()
	defer func() {
		gui.Log.Infof("postRefreshUpdate for %s took %s", c.GetKey(), time.Since(t))
	}()

	c.HandleRender()

	if gui.currentViewName() == c.GetViewName() {
		c.HandleFocus(types.OnFocusOpts{})
	} else {
		// The FocusLine call is included in the HandleFocus method which we
		// call for focused views above; but we need to call it here for
		// non-focused views to ensure that an inactive selection is painted
		// correctly, and that integration tests see the up to date selection
		// state.
		c.FocusLine()

		currentCtx := gui.State.ContextMgr.Current()
		if currentCtx.GetKey() == context.NORMAL_MAIN_CONTEXT_KEY || currentCtx.GetKey() == context.NORMAL_SECONDARY_CONTEXT_KEY {
			// Searching can't cope well with the view being updated while it is being searched.
			// We might be able to fix the problems with this, but it doesn't seem easy, so for now
			// just don't rerender the view while searching, on the assumption that users will probably
			// either search or change their data, but not both at the same time.
			if !currentCtx.GetView().IsSearching() {
				sidePanelContext := gui.State.ContextMgr.NextInStack(currentCtx)
				if sidePanelContext != nil && sidePanelContext.GetKey() == c.GetKey() {
					sidePanelContext.HandleRenderToMain()
				}
			}
		} else if c.GetKey() == gui.State.ContextMgr.CurrentStatic().GetKey() {
			// If our view is not the current one, but it is the current static context, then this
			// can only mean that a popup is showing. In that case we want to refresh the main view
			// behind the popup.
			c.HandleRenderToMain()
		}
	}
}
