package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) handleHelpPress(g *gocui.Gui, v *gocui.View) error {
	lineNumber := gui.getItemPosition(v)
	if len(gui.State.Keys) > lineNumber {
		err := gui.handleHelpClose(g, v)
		if err != nil {
			return err
		}
		return gui.State.Keys[lineNumber].Handler(g, v)
	}
	return nil
}

func (gui *Gui) handleHelpSelect(g *gocui.Gui, v *gocui.View) error {
	// doing nothing for now
	// but it is needed for switch in newLineFocused
	return nil
}

func (gui *Gui) renderHelpOptions(g *gocui.Gui) error {
	optionsMap := map[string]string{
		"esc/q": gui.Tr.SLocalize("close"),
		"↑ ↓":   gui.Tr.SLocalize("navigate"),
		"space": gui.Tr.SLocalize("execute"),
	}
	return gui.renderOptionsMap(g, optionsMap)
}

func (gui *Gui) handleHelpClose(g *gocui.Gui, v *gocui.View) error {
	// better to delete because for example after closing update confirmation panel,
	// the focus isn't set back to any of panels and one is unable to even quit
	//_, err := g.SetViewOnBottom(v.Name())
	err := g.DeleteView("help")
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

func (gui *Gui) handleHelp(g *gocui.Gui, v *gocui.View) error {
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
	helpView, _ := g.SetView("help", x0, y0, x1, y1-1, 0)
	helpView.Title = strings.Title(gui.Tr.SLocalize("help"))
	helpView.FgColor = gocui.ColorWhite

	if err := gui.renderHelpOptions(g); err != nil {
		return err
	}

	fmt.Fprint(helpView, content)

	g.Update(func(g *gocui.Gui) error {
		_, err := g.SetViewOnTop("help")
		if err != nil {
			return err
		}
		return gui.switchFocus(g, v, helpView)
	})
	return nil
}
