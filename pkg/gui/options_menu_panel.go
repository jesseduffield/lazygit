package gui

import (
	"strings"

	"github.com/go-errors/errors"

	"github.com/jesseduffield/gocui"
)

func (gui *Gui) getBindings(v *gocui.View) []*Binding {
	var (
		bindingsGlobal, bindingsPanel []*Binding
	)

	bindings := gui.GetCurrentKeybindings()

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

func (gui *Gui) handleCreateOptionsMenu(g *gocui.Gui, v *gocui.View) error {
	bindings := gui.getBindings(v)

	handleMenuPress := func(index int) error {
		if bindings[index].Key == nil {
			return nil
		}
		if index >= len(bindings) {
			return errors.New("Index is greater than size of bindings")
		}
		err := gui.handleMenuClose(g, v)
		if err != nil {
			return err
		}
		return bindings[index].Handler(g, v)
	}

	return gui.createMenu(strings.Title(gui.Tr.SLocalize("menu")), bindings, handleMenuPress)
}
