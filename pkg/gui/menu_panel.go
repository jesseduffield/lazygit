package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) handleMenuSelect(g *gocui.Gui, v *gocui.View) error {
	// doing nothing for now
	// but it is needed for switch in newLineFocused
	return nil
}

func (gui *Gui) renderMenuOptions(g *gocui.Gui) error {
	optionsMap := map[string]string{
		"esc/q": gui.Tr.SLocalize("close"),
		"↑ ↓":   gui.Tr.SLocalize("navigate"),
		"space": gui.Tr.SLocalize("execute"),
	}
	return gui.renderOptionsMap(g, optionsMap)
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

	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(gui.g, list)
	menuView, _ := gui.g.SetView("menu", x0, y0, x1, y1, 0)
	menuView.Title = strings.Title(gui.Tr.SLocalize("menu"))
	menuView.FgColor = gocui.ColorWhite
	menuView.Clear()
	fmt.Fprint(menuView, list)

	if err := gui.renderMenuOptions(gui.g); err != nil {
		return err
	}

	wrappedHandlePress := func(g *gocui.Gui, v *gocui.View) error {
		lineNumber := gui.getItemPosition(v)
		return handlePress(lineNumber)
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
