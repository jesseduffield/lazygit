package gui

import (
	"time"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
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

	// A search counts the matches in everything the view holds, so a re-render of a
	// view that is being searched is read all the way to the end (as opening the
	// search prompt reads it, see MainViewController.openSearch). Lines left unread
	// hold matches the search doesn't know about, and would add themselves to the
	// "x of y" as the user scrolled far enough to load them.
	if v.IsSearching() {
		return tasks.LinesToRead{
			Total:               -1,
			InitialRefreshAfter: linesForFirstRefresh,
		}
	}

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

	if windowName == "files" {
		if err := gui.applyFilesTab(tabs[tabIndex].Tab); err != nil {
			return err
		}
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
	if view != nil && view.Name() == "files" {
		view.TabIndex = gui.filesTabIndex()
	}
	return view
}

func (gui *Gui) filesTabIndex() int {
	if gui.State == nil || gui.State.Contexts == nil || gui.State.Contexts.Files == nil {
		return 0
	}
	if !gui.c.UserConfig().Gui.ShowFileChangeTabs {
		return 0
	}

	filter := gui.State.Contexts.Files.GetStatusFilter()
	switch filter {
	case filetree.DisplayUnstaged:
		return 1
	case filetree.DisplayStaged:
		return 2
	default:
		return 0
	}
}

func (gui *Gui) applyFilesTab(tabName string) error {
	if gui.State == nil || gui.State.Contexts == nil || gui.State.Contexts.Files == nil {
		return nil
	}

	if !gui.c.UserConfig().Gui.ShowFileChangeTabs {
		return nil
	}

	filesContext := gui.State.Contexts.Files
	desiredFilter := filetree.DisplayAll

	switch tabName {
	case gui.c.Tr.UnstagedChanges:
		desiredFilter = filetree.DisplayUnstaged
	case gui.c.Tr.StagedChanges:
		desiredFilter = filetree.DisplayStaged
	case gui.c.Tr.FilesTitle:
		desiredFilter = filetree.DisplayAll
	default:
		return nil
	}

	filesContext.FileTreeViewModel.SetStatusFilter(desiredFilter)
	filesContext.GetView().Subtitle = ""
	filesContext.GetView().TabIndex = gui.filesTabIndex()
	gui.postRefreshUpdate(filesContext)
	return nil
}

func (gui *Gui) render() {
	gui.c.OnUIThread(func() error { return nil })
}

// renderContentOnly triggers a re-render that skips the layout pass and only
// redraws the views whose content changed (relying on tcell's cell-level dirty
// tracking to emit just the cells that actually differ). Use it when only a
// view's content changed, not the window layout.
func (gui *Gui) renderContentOnly() {
	gui.c.OnUIThreadContentOnly(func() error { return nil })
}

// postRefreshUpdate is to be called on a context after the state that it depends on has been refreshed
// if the context's view is set to another context we do nothing.
// if the context's view is the current view we trigger a focus; re-selecting the current item.
func (gui *Gui) postRefreshUpdate(c types.Context, opts types.OnFocusOpts) {
	t := time.Now()
	defer func() {
		gui.Log.Infof("postRefreshUpdate for %s took %s", c.GetKey(), time.Since(t))
	}()

	c.HandleRender()

	// The render may have given the context its first item, or taken its last one
	// away, which decides whether its view draws a selection at all.
	gui.State.ContextMgr.updateSelectionHighlights()

	if gui.currentViewName() == c.GetInputViewName() {
		c.HandleFocus(opts)
	} else {
		// The FocusLine call is included in the HandleFocus method which we
		// call for focused views above; but we need to call it here for
		// non-focused views to ensure that an inactive selection is painted
		// correctly, and that integration tests see the up to date selection
		// state.
		c.FocusLine(!opts.KeepScrollPosition)
		if opts.SkipMainViewUpdate {
			return
		}

		currentCtx := gui.State.ContextMgr.Current()
		if currentCtx.GetKey() == context.NORMAL_MAIN_CONTEXT_KEY || currentCtx.GetKey() == context.NORMAL_SECONDARY_CONTEXT_KEY {
			sidePanelContext := gui.State.ContextMgr.NextInStack(currentCtx)
			if sidePanelContext != nil && sidePanelContext.GetKey() == c.GetKey() {
				sidePanelContext.HandleRenderToMain()
			}
		} else if c.GetKey() == gui.State.ContextMgr.CurrentStatic().GetKey() {
			// If our view is not the current one, but it is the current static context, then this
			// can only mean that a popup is showing. In that case we want to refresh the main view
			// behind the popup.
			c.HandleRenderToMain()
		}
	}
}
