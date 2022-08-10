package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/spkg/bom"
)

func (gui *Gui) resetOrigin(v *gocui.View) error {
	_ = v.SetCursor(0, 0)
	return v.SetOrigin(0, 0)
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
	panelHeight := gui.getMessageHeight(true, prompt, panelWidth) + suggestionsViewHeight
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensionsAux(panelWidth, panelHeight)
	confirmationViewBottom := y1 - suggestionsViewHeight
	_, _ = gui.g.SetView(gui.Views.Confirmation.Name(), x0, y0, x1, confirmationViewBottom, 0)

	suggestionsViewTop := confirmationViewBottom + 1
	_, _ = gui.g.SetView(gui.Views.Suggestions.Name(), x0, suggestionsViewTop, x1, suggestionsViewTop+suggestionsViewHeight, 0)
}

func (gui *Gui) globalOptionsMap() map[string]string {
	keybindingConfig := gui.c.UserConfig.Keybinding

	return map[string]string{
		fmt.Sprintf("%s/%s", keybindings.Label(keybindingConfig.Universal.ScrollUpMain), keybindings.Label(keybindingConfig.Universal.ScrollDownMain)):                                                                                                               gui.c.Tr.LcScroll,
		fmt.Sprintf("%s %s %s %s", keybindings.Label(keybindingConfig.Universal.PrevBlock), keybindings.Label(keybindingConfig.Universal.NextBlock), keybindings.Label(keybindingConfig.Universal.PrevItem), keybindings.Label(keybindingConfig.Universal.NextItem)): gui.c.Tr.LcNavigate,
		keybindings.Label(keybindingConfig.Universal.Return):     gui.c.Tr.LcCancel,
		keybindings.Label(keybindingConfig.Universal.Quit):       gui.c.Tr.LcQuit,
		keybindings.Label(keybindingConfig.Universal.OptionMenu): gui.c.Tr.LcMenu,
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
