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
	ctx := self.c.CurrentContext()
	// Don't show menu while displaying popup.
	if ctx.GetKind() == types.PERSISTENT_POPUP || ctx.GetKind() == types.TEMPORARY_POPUP {
		return nil
	}

	headingMap := map[string]string{
		"local":          self.c.Tr.KeybindingsMenuSectionLocal,
		"global":         self.c.Tr.KeybindingsMenuSectionGlobal,
		"navigation":     self.c.Tr.KeybindingsMenuSectionNavigation,
		"file":           "File",
		"workingTree":    "Working tree",
		"commit":         "Commit",
		"stash":          "Stash",
		"filterDisplay":  "Filter/Display",
		"sync":           "Sync",
		"customCommands": "Custom commands",
		"customPatch":    "Custom patch",
		"rebase":         "Rebase",
		"misc":           "Miscellaneous",
		"diff":           "Diff",
	}

	menuItems := []*types.MenuItem{}

	for groupKey, bindings := range self.getBindingGroups(ctx) {
		section := &types.MenuSection{
			Title:  headingMap[groupKey],
			Column: 1,
		}

		menuItems = append(menuItems,
			lo.Map(bindings, func(binding *types.Binding, _ int) *types.MenuItem {
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
					Section: section,
				}
			})...)
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title:           self.c.Tr.Keybindings,
		Items:           menuItems,
		HideCancel:      true,
		ColumnAlignment: []utils.Alignment{utils.AlignRight, utils.AlignLeft},
	})
}

func (self *OptionsMenuAction) getBindingGroups(context types.Context) map[string][]*types.Binding {
	bindings, _ := self.c.GetInitialKeybindingsWithCustomCommands()

	result := map[string][]*types.Binding{}
	appendToGroup := func(groupKey string, binding *types.Binding) {
		if _, ok := result[groupKey]; !ok {
			result[groupKey] = []*types.Binding{}
		}
		result[groupKey] = append(result[groupKey], binding)
	}

	for _, binding := range bindings {
		if keybindings.LabelFromKey(binding.Key) != "" && binding.Description != "" {
			if binding.ViewName == context.GetViewName() || binding.ViewName == "" {
				appendToGroup(binding.Tag, binding)
			}
		}
	}

	return result
}

// We shouldn't really need to do this. We should define alternative keys for the same
// handler in the keybinding struct.
func uniqueBindings(bindings []*types.Binding) []*types.Binding {
	return lo.UniqBy(bindings, func(binding *types.Binding) string {
		return binding.Description
	})
}
