package gui

import (
	"errors"
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

func (gui *Gui) getMenuOptions() map[string]string {
	keybindingConfig := gui.c.UserConfig.Keybinding

	return map[string]string{
		gui.getKeyDisplay(keybindingConfig.Universal.Return): gui.c.Tr.LcClose,
		fmt.Sprintf("%s %s", gui.getKeyDisplay(keybindingConfig.Universal.PrevItem), gui.getKeyDisplay(keybindingConfig.Universal.NextItem)): gui.c.Tr.LcNavigate,
		gui.getKeyDisplay(keybindingConfig.Universal.Select): gui.c.Tr.LcExecute,
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

	for _, item := range opts.Items {
		if item.OpensMenu && item.LabelColumns != nil {
			return errors.New("Message for the developer of this app: you've set opensMenu with displaystrings on the menu panel. Bad developer!. Apologies, user")
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
