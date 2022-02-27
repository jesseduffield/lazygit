package gui

import (
	"log"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) getBindings(context types.Context) []*types.Binding {
	var (
		bindingsGlobal, bindingsPanel, bindingsNavigation []*types.Binding
	)

	bindings, _ := gui.GetInitialKeybindings()
	customBindings, err := gui.CustomCommandsClient.GetCustomCommandKeybindings()
	if err != nil {
		log.Fatal(err)
	}
	bindings = append(customBindings, bindings...)

	for _, binding := range bindings {
		if GetKeyDisplay(binding.Key) != "" && binding.Description != "" {
			if len(binding.Contexts) == 0 && binding.ViewName == "" {
				bindingsGlobal = append(bindingsGlobal, binding)
			} else if binding.Tag == "navigation" {
				bindingsNavigation = append(bindingsNavigation, binding)
			} else if utils.IncludesString(binding.Contexts, string(context.GetKey())) {
				bindingsPanel = append(bindingsPanel, binding)
			}
		}
	}

	resultBindings := []*types.Binding{}
	resultBindings = append(resultBindings, uniqueBindings(bindingsPanel)...)
	// adding a separator between the panel-specific bindings and the other bindings
	resultBindings = append(resultBindings, &types.Binding{})
	resultBindings = append(resultBindings, uniqueBindings(bindingsGlobal)...)
	resultBindings = append(resultBindings, uniqueBindings(bindingsNavigation)...)

	return resultBindings
}

// We shouldn't really need to do this. We should define alternative keys for the same
// handler in the keybinding struct.
func uniqueBindings(bindings []*types.Binding) []*types.Binding {
	keys := make(map[string]bool)
	result := make([]*types.Binding, 0)

	for _, binding := range bindings {
		if _, ok := keys[binding.Description]; !ok {
			keys[binding.Description] = true
			result = append(result, binding)
		}
	}

	return result
}

func (gui *Gui) displayDescription(binding *types.Binding) string {
	if binding.OpensMenu {
		return presentation.OpensMenuStyle(binding.Description)
	}

	return style.FgCyan.Sprint(binding.Description)
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
