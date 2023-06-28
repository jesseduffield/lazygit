package controllers

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type OptionsMenuAction struct {
	c *ControllerCommon
}

func (self *OptionsMenuAction) Call() error {
	ctx := self.c.CurrentContext()
	// Don't show menu while displaying popup.
	if ctx.GetKind() == types.PERSISTENT_POPUP || ctx.GetKind() == types.TEMPORARY_POPUP {
		return nil
	}

	bindings := self.getBindings(ctx)

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

	return self.c.Menu(types.CreateMenuOptions{
		Title:      self.c.Tr.Keybindings,
		Items:      menuItems,
		HideCancel: true,
	})
}

func (self *OptionsMenuAction) getBindings(context types.Context) []*types.Binding {
	var bindingsGlobal, bindingsPanel, bindingsNavigation []*types.Binding

	bindings, _ := self.c.GetInitialKeybindingsWithCustomCommands()

	for _, binding := range bindings {
		if keybindings.LabelFromKey(binding.Key) != "" && binding.Description != "" && binding.Key != nil {
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
