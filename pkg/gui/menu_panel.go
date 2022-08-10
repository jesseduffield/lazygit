package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) getMenuOptions() map[string]string {
	keybindingConfig := gui.c.UserConfig.Keybinding

	return map[string]string{
		keybindings.Label(keybindingConfig.Universal.Return): gui.c.Tr.LcClose,
		fmt.Sprintf("%s %s", keybindings.Label(keybindingConfig.Universal.PrevItem), keybindings.Label(keybindingConfig.Universal.NextItem)): gui.c.Tr.LcNavigate,
		keybindings.Label(keybindingConfig.Universal.Select): gui.c.Tr.LcExecute,
	}
}

// note: items option is mutated by this function
func (gui *Gui) createMenu(opts types.CreateMenuOptions) error {
	if !opts.HideCancel {
		// this is mutative but I'm okay with that for now
		opts.Items = append(opts.Items, &types.MenuItem{
			LabelColumns: []string{gui.c.Tr.LcCancel},
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
			item.LabelColumns[0] = presentation.OpensMenuStyle(item.LabelColumns[0])
		}

		maxColumnSize = utils.Max(maxColumnSize, len(item.LabelColumns))
	}

	for _, item := range opts.Items {
		if len(item.LabelColumns) < maxColumnSize {
			// we require that each item has the same number of columns so we're padding out with blank strings
			// if this item has too few
			item.LabelColumns = append(item.LabelColumns, make([]string, maxColumnSize-len(item.LabelColumns))...)
		}
	}

	gui.State.Contexts.Menu.SetMenuItems(opts.Items)
	gui.State.Contexts.Menu.SetSelectedLineIdx(0)

	gui.Views.Menu.Title = opts.Title
	gui.Views.Menu.FgColor = theme.GocuiDefaultTextColor
	gui.Views.Menu.SetOnSelectItem(gui.onSelectItemWrapper(func(selectedLine int) error {
		return nil
	}))

	gui.Views.Tooltip.Wrap = true
	gui.Views.Tooltip.FgColor = theme.GocuiDefaultTextColor
	gui.Views.Tooltip.Visible = true

	// resetting keybindings so that the menu-specific keybindings are registered
	if err := gui.resetKeybindings(); err != nil {
		return err
	}

	_ = gui.c.PostRefreshUpdate(gui.State.Contexts.Menu)

	// TODO: ensure that if we're opened a menu from within a menu that it renders correctly
	return gui.c.PushContext(gui.State.Contexts.Menu)
}
