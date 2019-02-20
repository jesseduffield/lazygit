package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) handleMenuSelect(g *gocui.Gui, v *gocui.View) error {
	return gui.focusPoint(0, gui.State.Panels.Menu.SelectedLine, v)
}

func (gui *Gui) handleMenuNextLine(g *gocui.Gui, v *gocui.View) error {
	panelState := gui.State.Panels.Menu
	gui.changeSelectedLine(&panelState.SelectedLine, v.LinesHeight(), false)

	return gui.handleMenuSelect(g, v)
}

func (gui *Gui) handleMenuPrevLine(g *gocui.Gui, v *gocui.View) error {
	panelState := gui.State.Panels.Menu
	gui.changeSelectedLine(&panelState.SelectedLine, v.LinesHeight(), true)

	return gui.handleMenuSelect(g, v)
}

// specific functions

func (gui *Gui) renderMenuOptions() error {
	optionsMap := map[string]string{
		"esc/q": gui.Tr.SLocalize("close"),
		"↑ ↓":   gui.Tr.SLocalize("navigate"),
		"space": gui.Tr.SLocalize("execute"),
	}
	return gui.renderOptionsMap(optionsMap)
}

func (gui *Gui) handleMenuClose(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteKeybinding("menu", gocui.KeySpace, gocui.ModNone); err != nil {
		return err
	}
	err := g.DeleteView("menu")
	if err != nil {
		return err
	}
	return gui.returnFocus(g, v)
}

func (gui *Gui) createMenu(items interface{}, handlePress func(int) error) error {
	list, err := utils.RenderList(items)
	if err != nil {
		return err
	}

	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(gui.g, false, list)
	menuView, _ := gui.g.SetView("menu", x0, y0, x1, y1, 0)
	menuView.Title = strings.Title(gui.Tr.SLocalize("menu"))
	menuView.FgColor = gocui.ColorWhite
	menuView.Clear()
	fmt.Fprint(menuView, list)
	gui.State.Panels.Menu.SelectedLine = 0

	wrappedHandlePress := func(g *gocui.Gui, v *gocui.View) error {
		selectedLine := gui.State.Panels.Menu.SelectedLine
		return handlePress(selectedLine)
	}

	if err := gui.g.SetKeybinding("menu", gocui.KeySpace, gocui.ModNone, wrappedHandlePress); err != nil {
		return err
	}

	gui.g.Update(func(g *gocui.Gui) error {
		if _, err := g.SetViewOnTop("menu"); err != nil {
			return err
		}
		currentView := gui.g.CurrentView()
		return gui.switchFocus(gui.g, currentView, menuView)
	})
	return nil
}
