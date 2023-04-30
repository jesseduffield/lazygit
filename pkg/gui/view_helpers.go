package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/tasks"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/spkg/bom"
)

func (gui *Gui) resetOrigin(v *gocui.View) error {
	_ = v.SetCursor(0, 0)
	return v.SetOrigin(0, 0)
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

// renderString resets the origin of a view and sets its content
func (gui *Gui) renderString(view *gocui.View, s string) error {
	if err := view.SetOrigin(0, 0); err != nil {
		return err
	}
	if err := view.SetCursor(0, 0); err != nil {
		return err
	}
	gui.setViewContent(view, s)
	return nil
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

	c := gui.c.CurrentContext()

	if c == gui.State.Contexts.Menu {
		gui.resizeMenu()
	} else if c == gui.State.Contexts.Confirmation || c == gui.State.Contexts.Suggestions {
		gui.resizeConfirmationPanel()
	} else if c == gui.State.Contexts.CommitMessage || c == gui.State.Contexts.CommitDescription {
		gui.resizeCommitMessagePanels()
	}

	return nil
}

func (gui *Gui) resizeMenu() {
	itemCount := gui.State.Contexts.Menu.GetList().Len()
	offset := 3
	panelWidth := gui.getConfirmationPanelWidth()
	x0, y0, x1, y1 := gui.getPopupPanelDimensionsForContentHeight(panelWidth, itemCount+offset)
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
	x0, y0, x1, y1 := gui.getPopupPanelDimensionsAux(panelWidth, panelHeight)
	confirmationViewBottom := y1 - suggestionsViewHeight
	_, _ = gui.g.SetView(gui.Views.Confirmation.Name(), x0, y0, x1, confirmationViewBottom, 0)

	suggestionsViewTop := confirmationViewBottom + 1
	_, _ = gui.g.SetView(gui.Views.Suggestions.Name(), x0, suggestionsViewTop, x1, suggestionsViewTop+suggestionsViewHeight, 0)
}

func (gui *Gui) resizeCommitMessagePanels() {
	panelWidth := gui.getConfirmationPanelWidth()
	content := gui.Views.CommitDescription.TextArea.GetContent()
	summaryViewHeight := 3
	panelHeight := gui.getMessageHeight(false, content, panelWidth)
	minHeight := 7
	if panelHeight < minHeight {
		panelHeight = minHeight
	}
	x0, y0, x1, y1 := gui.getPopupPanelDimensionsAux(panelWidth, panelHeight)

	_, _ = gui.g.SetView(gui.Views.CommitMessage.Name(), x0, y0, x1, y0+summaryViewHeight-1, 0)
	_, _ = gui.g.SetView(gui.Views.CommitDescription.Name(), x0, y0+summaryViewHeight, x1, y1+summaryViewHeight, 0)
}

func (gui *Gui) globalOptionsMap() map[string]string {
	keybindingConfig := gui.c.UserConfig.Keybinding

	return map[string]string{
		fmt.Sprintf("%s/%s", keybindings.Label(keybindingConfig.Universal.ScrollUpMain), keybindings.Label(keybindingConfig.Universal.ScrollDownMain)):                                                                                                               gui.c.Tr.LcScroll,
		fmt.Sprintf("%s %s %s %s", keybindings.Label(keybindingConfig.Universal.PrevBlock), keybindings.Label(keybindingConfig.Universal.NextBlock), keybindings.Label(keybindingConfig.Universal.PrevItem), keybindings.Label(keybindingConfig.Universal.NextItem)): gui.c.Tr.LcNavigate,
		keybindings.Label(keybindingConfig.Universal.Return):         gui.c.Tr.LcCancel,
		keybindings.Label(keybindingConfig.Universal.Quit):           gui.c.Tr.LcQuit,
		keybindings.Label(keybindingConfig.Universal.OptionMenuAlt1): gui.c.Tr.LcMenu,
		fmt.Sprintf("%s-%s", keybindings.Label(keybindingConfig.Universal.JumpToBlock[0]), keybindings.Label(keybindingConfig.Universal.JumpToBlock[len(keybindingConfig.Universal.JumpToBlock)-1])): gui.c.Tr.LcJump,
		fmt.Sprintf("%s/%s", keybindings.Label(keybindingConfig.Universal.ScrollLeft), keybindings.Label(keybindingConfig.Universal.ScrollRight)):                                                    gui.c.Tr.LcScrollLeftRight,
	}
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

	context, ok := gui.contextForView(viewName)
	if !ok {
		return nil
	}

	return gui.c.PushContext(context)
}

func (gui *Gui) contextForView(viewName string) (types.Context, bool) {
	view, err := gui.g.View(viewName)
	if err != nil {
		return nil, false
	}

	for _, context := range gui.State.Contexts.Flatten() {
		if context.GetViewName() == view.Name() {
			return context, true
		}
	}

	return nil, false
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
	context := gui.currentStaticContext()
	view, _ := gui.g.View(context.GetViewName())
	return view
}

func (gui *Gui) render() {
	gui.c.OnUIThread(func() error { return nil })
}
