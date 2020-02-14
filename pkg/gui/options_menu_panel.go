package gui

import (
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) getBindings(v *gocui.View) []*Binding {
	var (
		bindingsGlobal, bindingsPanel []*Binding
	)

	bindings := gui.GetInitialKeybindings()

	for _, binding := range bindings {
		if GetKeyDisplay(binding.Key) != "" && binding.Description != "" {
			switch binding.ViewName {
			case "":
				bindingsGlobal = append(bindingsGlobal, binding)
			case v.Name():
				if len(binding.Contexts) == 0 || utils.IncludesString(binding.Contexts, v.Context) {
					bindingsPanel = append(bindingsPanel, binding)
				}
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

	menuItems := make([]*menuItem, len(bindings))

	for i, binding := range bindings {
		innerBinding := binding // note to self, never close over loop variables
		menuItems[i] = &menuItem{
			displayStrings: []string{GetKeyDisplay(innerBinding.Key), innerBinding.Description},
			onPress: func() error {
				if innerBinding.Key == nil {
					return nil
				}
				if err := gui.handleMenuClose(g, v); err != nil {
					return err
				}
				return innerBinding.Handler(g, v)
			},
		}
	}

	return gui.createMenu(strings.Title(gui.Tr.SLocalize("menu")), menuItems, createMenuOptions{})
}
