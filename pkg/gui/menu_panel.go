package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/samber/lo"
)

// note: items option is mutated by this function
func (gui *Gui) createMenu(opts types.CreateMenuOptions) error {
	if !opts.HideCancel {
		// this is mutative but I'm okay with that for now
		opts.Items = append(opts.Items, &types.MenuItem{
			LabelColumns: []string{gui.c.Tr.Cancel},
			OnPress: func() error {
				if opts.OnCancel != nil {
					return opts.OnCancel()
				}
				return nil
			},
		})
	}

	maxColumnSize := 1

	// Only the primary key of each navigation binding is reserved as
	// essential; alternates (e.g. the historical j/k that lived under
	// `*Alt` fields) stay available to be reused by menu items, which
	// take precedence over the inherited list bindings.
	essentialKeys := []gocui.Key{
		config.GetValidatedKeyBindingKeys(gui.c.UserConfig().Keybinding.Universal.ConfirmMenu)[0],
		config.GetValidatedKeyBindingKeys(gui.c.UserConfig().Keybinding.Universal.Return)[0],
		config.GetValidatedKeyBindingKeys(gui.c.UserConfig().Keybinding.Universal.PrevItem)[0],
		config.GetValidatedKeyBindingKeys(gui.c.UserConfig().Keybinding.Universal.NextItem)[0],
	}

	for _, item := range opts.Items {
		if item.LabelColumns == nil {
			item.LabelColumns = []string{item.Label}
		}

		if item.OpensMenu {
			item.LabelColumns[0] = fmt.Sprintf("%s...", item.LabelColumns[0])
		}

		maxColumnSize = max(maxColumnSize, len(item.LabelColumns))

		// Remove all item keybindings that are the same as one of the essential bindings
		if !opts.KeepConflictingKeybindings {
			item.Keys = lo.Filter(item.Keys, func(k gocui.Key, _ int) bool {
				return !lo.Contains(essentialKeys, k)
			})
		}
	}

	for _, item := range opts.Items {
		if len(item.LabelColumns) < maxColumnSize {
			// we require that each item has the same number of columns so we're padding out with blank strings
			// if this item has too few
			item.LabelColumns = append(item.LabelColumns, make([]string, maxColumnSize-len(item.LabelColumns))...)
		}
	}

	gui.State.Contexts.Menu.SetMenuItems(opts.Items, opts.ColumnAlignment)
	gui.State.Contexts.Menu.SetPrompt(opts.Prompt)
	gui.State.Contexts.Menu.SetAllowFilteringKeybindings(opts.AllowFilteringKeybindings)
	gui.State.Contexts.Menu.SetKeybindingsTakePrecedence(!opts.KeepConflictingKeybindings)
	gui.State.Contexts.Menu.SetOnCancel(opts.OnCancel)
	gui.State.Contexts.Menu.SetSelection(0)

	gui.Views.Menu.SetOriginY(0)

	gui.Views.Menu.Title = opts.Title
	gui.Views.Menu.FgColor = theme.GocuiDefaultTextColor

	gui.Views.Tooltip.Wrap = true
	gui.Views.Tooltip.FgColor = theme.GocuiDefaultTextColor
	gui.Views.Tooltip.Visible = true

	// resetting keybindings so that the menu-specific keybindings are registered
	if err := gui.resetKeybindings(); err != nil {
		return err
	}

	gui.c.PostRefreshUpdate(gui.State.Contexts.Menu)

	// TODO: ensure that if we're opened a menu from within a menu that it renders correctly
	gui.c.Context().Push(gui.State.Contexts.Menu, types.OnFocusOpts{})
	return nil
}
