package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type OptionsMenuAction struct {
	c *ControllerCommon
}

func (self *OptionsMenuAction) Call() error {
	ctx := self.c.Context().Current()
	local, global, navigation := self.getBindings(ctx)

	menuItems := []*types.MenuItem{}

	appendBindings := func(bindings []*types.Binding, section *types.MenuSection) {
		menuItems = append(menuItems,
			lo.Map(bindings, func(binding *types.Binding, _ int) *types.MenuItem {
				var disabledReason *types.DisabledReason
				if binding.GetDisabledReason != nil {
					disabledReason = binding.GetDisabledReason()
				}
				return &types.MenuItem{
					OpensMenu: binding.OpensMenu,
					Label:     binding.Description,
					OnPress: func() error {
						if binding.Handler == nil {
							return nil
						}

						return self.c.IGuiCommon.CallKeybindingHandler(binding)
					},
					Key:            binding.Key,
					Tooltip:        binding.Tooltip,
					DisabledReason: disabledReason,
					Section:        section,
				}
			})...)
	}

	appendBindings(local, &types.MenuSection{Title: self.c.Tr.KeybindingsMenuSectionLocal, Column: 1})
	appendBindings(global, &types.MenuSection{Title: self.c.Tr.KeybindingsMenuSectionGlobal, Column: 1})
	appendBindings(navigation, &types.MenuSection{Title: self.c.Tr.KeybindingsMenuSectionNavigation, Column: 1})

	return self.c.Menu(types.CreateMenuOptions{
		Title:           self.c.Tr.Keybindings,
		Items:           menuItems,
		HideCancel:      true,
		ColumnAlignment: []utils.Alignment{utils.AlignRight, utils.AlignLeft},
	})
}

// Returns three slices of bindings: local, global, and navigation
func (self *OptionsMenuAction) getBindings(context types.Context) ([]*types.Binding, []*types.Binding, []*types.Binding) {
	var bindingsGlobal, bindingsPanel, bindingsNavigation []*types.Binding

	bindings, _ := self.c.GetInitialKeybindingsWithCustomCommands()

	for _, binding := range bindings {
		if keybindings.LabelFromKey(binding.Key) != "" && binding.Description != "" {
			if binding.ViewName == "" {
				bindingsGlobal = append(bindingsGlobal, binding)
			} else if binding.ViewName == context.GetViewName() {
				if binding.Tag == "navigation" {
					bindingsNavigation = append(bindingsNavigation, binding)
				} else {
					bindingsPanel = append(bindingsPanel, binding)
				}
			}
		}
	}

	return uniqueBindings(bindingsPanel), uniqueBindings(bindingsGlobal), uniqueBindings(bindingsNavigation)
}

// We shouldn't really need to do this. We should define alternative keys for the same
// handler in the keybinding struct.
func uniqueBindings(bindings []*types.Binding) []*types.Binding {
	return lo.UniqBy(bindings, func(binding *types.Binding) string {
		return binding.Description
	})
}
