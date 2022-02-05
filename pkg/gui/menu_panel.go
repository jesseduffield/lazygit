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

func (gui *Gui) handleMenuClose() error {
	return gui.returnFromContext()
}

// note: items option is mutated by this function
func (gui *Gui) createMenu(opts types.CreateMenuOptions) error {
	if !opts.HideCancel {
		// this is mutative but I'm okay with that for now
		opts.Items = append(opts.Items, &types.MenuItem{
			DisplayStrings: []string{gui.c.Tr.LcCancel},
			OnPress: func() error {
				return nil
			},
		})
	}

	for _, item := range opts.Items {
		if item.OpensMenu && item.DisplayStrings != nil {
			return errors.New("Message for the developer of this app: you've set opensMenu with displaystrings on the menu panel. Bad developer!. Apologies, user")
		}
	}

	x0, y0, x1, y1 := gui.getConfirmationPanelDimensionsForContentHeight(len(opts.Items))
	menuView, _ := gui.g.SetView("menu", x0, y0, x1, y1, 0)
	menuView.Title = opts.Title
	menuView.FgColor = theme.GocuiDefaultTextColor
	menuView.SetOnSelectItem(gui.onSelectItemWrapper(func(selectedLine int) error {
		return nil
	}))

	gui.State.Contexts.Menu.SetMenuItems(opts.Items)
	gui.State.Contexts.Menu.GetPanelState().SetSelectedLineIdx(0)
	_ = gui.c.PostRefreshUpdate(gui.State.Contexts.Menu)

	// TODO: ensure that if we're opened a menu from within a menu that it renders correctly
	return gui.c.PushContext(gui.State.Contexts.Menu)
}
