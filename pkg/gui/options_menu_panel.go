package gui

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) getBindings(context types.Context) []*types.Binding {
	var (
		bindingsGlobal, bindingsPanel []*types.Binding
	)

	bindings, _ := gui.GetInitialKeybindings()
	bindings = append(gui.GetCustomCommandKeybindings(), bindings...)

	for _, binding := range bindings {
		if GetKeyDisplay(binding.Key) != "" && binding.Description != "" {
			if len(binding.Contexts) == 0 {
				bindingsGlobal = append(bindingsGlobal, binding)
			} else if utils.IncludesString(binding.Contexts, string(context.GetKey())) {
				bindingsPanel = append(bindingsPanel, binding)
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
	context := gui.currentContext()
	bindings := gui.getBindings(context)

	menuItems := make([]*types.MenuItem, len(bindings))

	for i, binding := range bindings {
		binding := binding // note to self, never close over loop variables
		menuItems[i] = &types.MenuItem{
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

	return gui.c.Menu(types.CreateMenuOptions{
		Title:      strings.Title(gui.c.Tr.LcMenu),
		Items:      menuItems,
		HideCancel: true,
	})
}
