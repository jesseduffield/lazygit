package gui

import (
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) getBindings(v *gocui.View) []*Binding {
	var (
		bindingsGlobal, bindingsPanel []*Binding
	)

	bindings := append(gui.GetCustomCommandKeybindings(), gui.GetInitialKeybindings()...)

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

func (gui *Gui) displayDescription(binding *Binding) string {
	if binding.OpensMenu {
		return opensMenuStyle(binding.Description)
	}

	return style.FgCyan.Sprint(binding.Description)
}

func opensMenuStyle(str string) string {
	return style.FgMagenta.Sprintf("%s...", str)
}

func (gui *Gui) handleCreateOptionsMenu() error {
	view := gui.g.CurrentView()
	if view == nil {
		return nil
	}

	bindings := gui.getBindings(view)

	menuItems := make([]*menuItem, len(bindings))

	for i, binding := range bindings {
		binding := binding // note to self, never close over loop variables
		menuItems[i] = &menuItem{
			displayStrings: []string{GetKeyDisplay(binding.Key), gui.displayDescription(binding)},
			onPress: func() error {
				if binding.Key == nil {
					return nil
				}
				if err := gui.handleMenuClose(); err != nil {
					return err
				}
				return binding.Handler()
			},
		}
	}

	return gui.createMenu(strings.Title(gui.Tr.LcMenu), menuItems, createMenuOptions{})
}
