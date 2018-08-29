package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
)

var keys []Binding

func (gui *Gui) handleHelpPress(g *gocui.Gui, v *gocui.View) error {
	lineNumber := gui.getItemPosition(v)
	err := gui.handleHelpClose(g, v)
	if err != nil {
		return err
	}
	return keys[lineNumber].Handler(g, v)
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
	err := g.DeleteView(v.Name())
	if err != nil {
		return err
	}
	return gui.returnFocus(g, v)
}

func (gui *Gui) handleHelp(g *gocui.Gui, v *gocui.View) error {
	// clear keys slice, so we don't have ghost elements
	keys = keys[:0]
	content := ""
	bindings := gui.getKeybindings()
	maxX, maxY := g.Size()
	x := maxX * 3 / 4
	y := 5
	helpView, _ := g.SetView("help", maxX-x, y, x, maxY-y, 0)
	helpView.Title = strings.Title(gui.Tr.SLocalize("help"))

	gui.renderHelpOptions(g)

	for _, binding := range bindings {
		if binding.ViewName == v.Name() && binding.Description != "" && binding.KeyReadable != "" {
			content += fmt.Sprintf(" %s - %s\n", binding.KeyReadable, binding.Description)
			keys = append(keys, binding)
		}
	}

	// for testing
	/*content += "first\n"
	content += "second\n"
	content += "third\n"
	*/
	gui.renderString(g, "help", content)

	g.Update(func(g *gocui.Gui) error {
		g.SetViewOnTop("help")
		gui.switchFocus(g, v, helpView)
		return nil
	})
	return nil
}
