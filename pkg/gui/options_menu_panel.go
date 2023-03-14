package gui

import (
	"log"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

func (gui *Gui) getBindings(context types.Context) []*types.Binding {
	var bindingsGlobal, bindingsPanel, bindingsNavigation []*types.Binding

	bindings, _ := gui.GetInitialKeybindings()
	customBindings, err := gui.CustomCommandsClient.GetCustomCommandKeybindings()
	if err != nil {
		log.Fatal(err)
	}
	bindings = append(customBindings, bindings...)

	for _, binding := range bindings {
		if keybindings.LabelFromKey(binding.Key) != "" && binding.Description != "" {
			if binding.ViewName == "" {
				bindingsGlobal = append(bindingsGlobal, binding)
			} else if binding.Tag == "navigation" {
				bindingsNavigation = append(bindingsNavigation, binding)
			} else if binding.ViewName == context.GetViewName() {
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
	return lo.UniqBy(bindings, func(binding *types.Binding) string {
		return binding.Description
	})
}

func (gui *Gui) handleCreateOptionsMenu() error {
	ctx := gui.currentContext()
	// Don't show menu while displaying popup.
	if ctx.GetKind() == types.PERSISTENT_POPUP || ctx.GetKind() == types.TEMPORARY_POPUP {
		return nil
	}

	bindings := gui.getBindings(ctx)

	menuItems := slices.Map(bindings, func(binding *types.Binding) *types.MenuItem {
		return &types.MenuItem{
			OpensMenu: binding.OpensMenu,
			Label:     binding.Description,
			OnPress: func() error {
				if binding.Handler == nil {
					return nil
				}

				return binding.Handler()
			},
			Key:     binding.Key,
			Tooltip: binding.Tooltip,
		}
	})

	return gui.c.Menu(types.CreateMenuOptions{
		Title:      gui.c.Tr.MenuTitle,
		Items:      menuItems,
		HideCancel: true,
	})
}
