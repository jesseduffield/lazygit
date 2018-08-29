package gui

import (
	"fmt"
	"github.com/jesseduffield/gocui"
	"strings"
)

func (gui *Gui) renderHelpOptions(g *gocui.Gui) error {
	optionsMap := map[string]string{
		"esc/q":     gui.Tr.SLocalize("close"),
		"PgUp/PgDn": gui.Tr.SLocalize("scroll"),
	}
	return gui.renderOptionsMap(g, optionsMap)
}

func (gui *Gui) scrollUpHelp(g *gocui.Gui, v *gocui.View) error {
	mainView, _ := g.View("help")
	ox, oy := mainView.Origin()
	if oy >= 1 {
		return mainView.SetOrigin(ox, oy-gui.Config.GetUserConfig().GetInt("gui.scrollHeight"))
	}
	return nil
}

func (gui *Gui) scrollDownHelp(g *gocui.Gui, v *gocui.View) error {
	mainView, _ := g.View("help")
	ox, oy := mainView.Origin()
	if oy < len(mainView.BufferLines()) {
		return mainView.SetOrigin(ox, oy+gui.Config.GetUserConfig().GetInt("gui.scrollHeight"))
	}
	return nil
}

func (gui *Gui) handleHelpClose(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnBottom(v.Name())
	return gui.switchFocus(g, v, gui.getFilesView(g))
}

func (gui *Gui) handleHelp(g *gocui.Gui, v *gocui.View) error {
	content := ""
	curr := ""
	bindings := gui.getKeybindings()
	maxX, maxY := g.Size()
	helpView, _ := g.SetView("help", 0, 0, maxX-1, maxY-2, 0)
	helpView.Title = strings.Title(gui.Tr.SLocalize("help"))

	gui.renderHelpOptions(g)

	for _, binding := range bindings {
		if binding.Description != "" {
			if curr != binding.ViewName {
				curr = binding.ViewName
				content += fmt.Sprintf("\n%s:\n", strings.Title(curr))
			}
			content += fmt.Sprintf("  %s - %s\n", binding.KeyReadable, binding.Description)
		}
	}

	helpView.Write([]byte(content))

	g.Update(func(g *gocui.Gui) error {
		g.SetViewOnTop("help")
		gui.switchFocus(g, v, helpView)
		return nil
	})
	return nil
}
