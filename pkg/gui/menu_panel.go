package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

// note: items option is mutated by this function
func (gui *Gui) createMenu(opts types.CreateMenuOptions) error {
	if !opts.HideCancel {
		// this is mutative but I'm okay with that for now
		opts.Items = append(opts.Items, &types.MenuItem{
			LabelColumns: []string{gui.c.Tr.Cancel},
			OnPress: func() error {
				return nil
			},
		})
	}

	maxColumnSize := 1

	for _, item := range opts.Items {
		if item.LabelColumns == nil {
			item.LabelColumns = []string{item.Label}
		}

		if item.OpensMenu {
			item.LabelColumns[0] = fmt.Sprintf("%s...", item.LabelColumns[0])
		}

		maxColumnSize = max(maxColumnSize, len(item.LabelColumns))
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
	gui.State.Contexts.Menu.SetSelection(0)

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
	gui.c.Context().Push(gui.State.Contexts.Menu)
	return nil
}
