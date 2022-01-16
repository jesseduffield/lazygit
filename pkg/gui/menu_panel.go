package gui

import (
	"errors"
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
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
func (gui *Gui) createMenu(opts popup.CreateMenuOptions) error {
	if !opts.HideCancel {
		// this is mutative but I'm okay with that for now
		opts.Items = append(opts.Items, &popup.MenuItem{
			DisplayStrings: []string{gui.c.Tr.LcCancel},
			OnPress: func() error {
				return nil
			},
		})
	}

	gui.State.MenuItems = opts.Items

	stringArrays := make([][]string, len(opts.Items))
	for i, item := range opts.Items {
		if item.OpensMenu && item.DisplayStrings != nil {
			return errors.New("Message for the developer of this app: you've set opensMenu with displaystrings on the menu panel. Bad developer!. Apologies, user")
		}

		if item.DisplayStrings == nil {
			styledStr := item.DisplayString
			if item.OpensMenu {
				styledStr = opensMenuStyle(styledStr)
			}
			stringArrays[i] = []string{styledStr}
		} else {
			stringArrays[i] = item.DisplayStrings
		}
	}

	list := utils.RenderDisplayStrings(stringArrays)

	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(false, list)
	menuView, _ := gui.g.SetView("menu", x0, y0, x1, y1, 0)
	menuView.Title = opts.Title
	menuView.FgColor = theme.GocuiDefaultTextColor
	menuView.SetOnSelectItem(gui.onSelectItemWrapper(func(selectedLine int) error {
		return nil
	}))
	menuView.SetContent(list)
	gui.State.Panels.Menu.SelectedLineIdx = 0

	return gui.c.PushContext(gui.State.Contexts.Menu)
}

func (gui *Gui) getSelectedMenuItem() *popup.MenuItem {
	if len(gui.State.MenuItems) == 0 {
		return nil
	}

	return gui.State.MenuItems[gui.State.Panels.Menu.SelectedLineIdx]
}
