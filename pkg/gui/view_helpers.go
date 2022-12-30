package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/tasks"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/spkg/bom"
)

func (gui *Gui) resetViewOrigin(v *gocui.View) {
	if err := v.SetCursor(0, 0); err != nil {
		gui.Log.Error(err)
	}

	if err := v.SetOrigin(0, 0); err != nil {
		gui.Log.Error(err)
	}
}

// Returns the number of lines that we should read initially from a cmd task so
// that the scrollbar has the correct size, along with the number of lines after
// which the view is filled and we can do a first refresh.
func (gui *Gui) linesToReadFromCmdTask(v *gocui.View) tasks.LinesToRead {
	_, height := v.Size()
	_, oy := v.Origin()

	linesForFirstRefresh := height + oy + 10

	// We want to read as many lines initially as necessary to let the
	// scrollbar go to its minimum height, so that the scrollbar thumb doesn't
	// change size as you scroll down.
	minScrollbarHeight := 2
	linesToReadForAccurateScrollbar := height*(height-1)/minScrollbarHeight + oy

	// However, cap it at some arbitrary max limit, so that we don't get
	// performance problems for huge monitors or tiny font sizes
	if linesToReadForAccurateScrollbar > 5000 {
		linesToReadForAccurateScrollbar = 5000
	}

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

func (gui *Gui) resizeCurrentPopupPanel() error {
	v := gui.g.CurrentView()
	if v == nil {
		return nil
	}

	if v == gui.Views.Menu {
		gui.resizeMenu()
	} else if v == gui.Views.Confirmation || v == gui.Views.Suggestions {
		gui.resizeConfirmationPanel()
	} else if gui.isPopupPanel(v.Name()) {
		return gui.resizePopupPanel(v, v.Buffer())
	}

	return nil
}

func (gui *Gui) resizePopupPanel(v *gocui.View, content string) error {
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(v.Wrap, content)
	_, err := gui.g.SetView(v.Name(), x0, y0, x1, y1, 0)
	return err
}

func (gui *Gui) resizeMenu() {
	itemCount := gui.State.Contexts.Menu.GetList().Len()
	offset := 3
	panelWidth := gui.getConfirmationPanelWidth()
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensionsForContentHeight(panelWidth, itemCount+offset)
	menuBottom := y1 - offset
	_, _ = gui.g.SetView(gui.Views.Menu.Name(), x0, y0, x1, menuBottom, 0)

	tooltipTop := menuBottom + 1
	tooltipHeight := gui.getMessageHeight(true, gui.State.Contexts.Menu.GetSelected().Tooltip, panelWidth) + 2 // plus 2 for the frame
	_, _ = gui.g.SetView(gui.Views.Tooltip.Name(), x0, tooltipTop, x1, tooltipTop+tooltipHeight-1, 0)
}

func (gui *Gui) resizeConfirmationPanel() {
	suggestionsViewHeight := 0
	if gui.Views.Suggestions.Visible {
		suggestionsViewHeight = 11
	}
	panelWidth := gui.getConfirmationPanelWidth()
	prompt := gui.Views.Confirmation.Buffer()
	wrap := true
	if gui.Views.Confirmation.Editable {
		prompt = gui.Views.Confirmation.TextArea.GetContent()
		wrap = false
	}
	panelHeight := gui.getMessageHeight(wrap, prompt, panelWidth) + suggestionsViewHeight
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensionsAux(panelWidth, panelHeight)
	confirmationViewBottom := y1 - suggestionsViewHeight
	_, _ = gui.g.SetView(gui.Views.Confirmation.Name(), x0, y0, x1, confirmationViewBottom, 0)

	suggestionsViewTop := confirmationViewBottom + 1
	_, _ = gui.g.SetView(gui.Views.Suggestions.Name(), x0, suggestionsViewTop, x1, suggestionsViewTop+suggestionsViewHeight, 0)
}

func (gui *Gui) isPopupPanel(viewName string) bool {
	return viewName == "commitMessage" || viewName == "confirmation" || viewName == "menu"
}

func (gui *Gui) popupPanelFocused() bool {
	return gui.isPopupPanel(gui.currentViewName())
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

	return gui.c.PushContext(context)
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
	context := gui.c.CurrentStaticContext()
	view, _ := gui.g.View(context.GetViewName())
	return view
}

func (gui *Gui) render() {
	gui.c.OnUIThread(func() error { return nil })
}

// postRefreshUpdate is to be called on a context after the state that it depends on has been refreshed
// if the context's view is set to another context we do nothing.
// if the context's view is the current view we trigger a focus; re-selecting the current item.
func (gui *Gui) postRefreshUpdate(c types.Context) error {
	if err := c.HandleRender(); err != nil {
		return err
	}

	if gui.currentViewName() == c.GetViewName() {
		if err := c.HandleFocus(types.OnFocusOpts{}); err != nil {
			return err
		}
	}

	return nil
}
