package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) handleMenuPress(g *gocui.Gui, v *gocui.View) error {
	lineNumber := gui.getItemPosition(v)
	if gui.State.Keys[lineNumber].Key == nil {
		return nil
	}
	if len(gui.State.Keys) > lineNumber {
		err := gui.handleMenuClose(g, v)
		if err != nil {
			return err
		}
		return gui.State.Keys[lineNumber].Handler(g, v)
	}
	return nil
}

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
	// better to delete because for example after closing update confirmation panel,
	// the focus isn't set back to any of panels and one is unable to even quit
	//_, err := g.SetViewOnBottom(v.Name())
	err := g.DeleteView("menu")
	if err != nil {
		return err
	}
	return gui.returnFocus(g, v)
}

func (gui *Gui) getBindings(v *gocui.View) []*Binding {
	var (
		bindingsGlobal, bindingsPanel []*Binding
	)

	bindings := gui.GetKeybindings()

	for _, binding := range bindings {
		if binding.GetKey() != "" && binding.Description != "" {
			switch binding.ViewName {
			case "":
				bindingsGlobal = append(bindingsGlobal, binding)
			case v.Name():
				bindingsPanel = append(bindingsPanel, binding)
			}
		}
	}

	// append dummy element to have a separator between
	// panel and global keybindings
	bindingsPanel = append(bindingsPanel, &Binding{})
	return append(bindingsPanel, bindingsGlobal...)
}

func (gui *Gui) handleMenu(g *gocui.Gui, v *gocui.View) error {
	gui.State.Keys = gui.getBindings(v)

	list, err := utils.RenderList(gui.State.Keys)
	if err != nil {
		return err
	}

	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(g, list)
	menuView, _ := g.SetView("menu", x0, y0, x1, y1, 0)
	menuView.Title = strings.Title(gui.Tr.SLocalize("menu"))
	menuView.FgColor = gocui.ColorWhite
	menuView.Clear()
	fmt.Fprint(menuView, list)

	if err := gui.renderMenuOptions(g); err != nil {
		return err
	}

	g.Update(func(g *gocui.Gui) error {
		_, err := g.SetViewOnTop("menu")
		if err != nil {
			return err
		}
		return gui.switchFocus(g, v, menuView)
	})
	return nil
}
