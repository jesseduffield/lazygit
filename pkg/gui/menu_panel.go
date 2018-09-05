package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) handleMenuPress(g *gocui.Gui, v *gocui.View) error {
	lineNumber := gui.getItemPosition(v)
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

func (gui *Gui) GetKey(binding Binding) string {
	r, ok := binding.Key.(rune)
	key := ""

	if ok {
		key = string(r)
	} else if binding.KeyReadable != "" {
		key = binding.KeyReadable
	}

	return key
}

func (gui *Gui) getMaxKeyLength(bindings []Binding) int {
	max := 0
	for _, binding := range bindings {
		keyLength := len(gui.GetKey(binding))
		if keyLength > max {
			max = keyLength
		}
	}
	return max
}

func (gui *Gui) handleMenu(g *gocui.Gui, v *gocui.View) error {
	// clear keys slice, so we don't have ghost elements
	gui.State.Keys = gui.State.Keys[:0]
	content := ""
	bindings := gui.GetKeybindings()
	padWidth := gui.getMaxKeyLength(bindings)

	for _, binding := range bindings {
		if key := gui.GetKey(binding); key != "" && binding.ViewName == v.Name() && binding.Description != "" {
			content += fmt.Sprintf("%s  %s\n", utils.WithPadding(key, padWidth), binding.Description)
			gui.State.Keys = append(gui.State.Keys, binding)
		}
	}

	// y1-1 so there will not be an extra space at the end of panel
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(g, content)
	menuView, _ := g.SetView("menu", x0, y0, x1, y1-1, 0)
	menuView.Title = strings.Title(gui.Tr.SLocalize("menu"))
	menuView.FgColor = gocui.ColorWhite

	if err := gui.renderMenuOptions(g); err != nil {
		return err
	}

	fmt.Fprint(menuView, content)

	g.Update(func(g *gocui.Gui) error {
		_, err := g.SetViewOnTop("menu")
		if err != nil {
			return err
		}
		return gui.switchFocus(g, v, menuView)
	})
	return nil
}
