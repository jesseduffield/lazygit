package gui

import (
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) getBindings(v *gocui.View) []*types.Binding {
	var (
		bindingsGlobal, bindingsPanel []*types.Binding
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
	bindingsPanel = append(bindingsPanel, &types.Binding{})
	return append(bindingsPanel, bindingsGlobal...)
}

func (gui *Gui) displayDescription(binding *types.Binding) string {
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

	menuItems := make([]*popup.MenuItem, len(bindings))

	for i, binding := range bindings {
		binding := binding // note to self, never close over loop variables
		menuItems[i] = &popup.MenuItem{
			DisplayStrings: []string{GetKeyDisplay(binding.Key), gui.displayDescription(binding)},
			OnPress: func() error {
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

	return gui.c.Menu(popup.CreateMenuOptions{
		Title:      strings.Title(gui.c.Tr.LcMenu),
		Items:      menuItems,
		HideCancel: true,
	})
}
