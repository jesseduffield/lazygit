package gui

import (
	"errors"
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type menuItem struct {
	displayString  string
	displayStrings []string
	onPress        func() error
	// only applies when displayString is used
	opensMenu bool
}

func (gui *Gui) getMenuOptions() map[string]string {
	keybindingConfig := gui.UserConfig.Keybinding

	return map[string]string{
		gui.getKeyDisplay(keybindingConfig.Universal.Return): gui.Tr.LcClose,
		fmt.Sprintf("%s %s", gui.getKeyDisplay(keybindingConfig.Universal.PrevItem), gui.getKeyDisplay(keybindingConfig.Universal.NextItem)): gui.Tr.LcNavigate,
		gui.getKeyDisplay(keybindingConfig.Universal.Select): gui.Tr.LcExecute,
	}
}

func (gui *Gui) handleMenuClose() error {
	return gui.returnFromContext()
}

type createMenuOptions struct {
	hideCancel bool
}

func (gui *Gui) createMenu(title string, items []*menuItem, createMenuOptions createMenuOptions) error {
	if !createMenuOptions.hideCancel {
		// this is mutative but I'm okay with that for now
		items = append(items, &menuItem{
			displayStrings: []string{gui.Tr.LcCancel},
			onPress: func() error {
				return nil
			},
		})
	}

	gui.State.MenuItems = items

	stringArrays := make([][]string, len(items))
	for i, item := range items {
		if item.opensMenu && item.displayStrings != nil {
			return errors.New("Message for the developer of this app: you've set opensMenu with displaystrings on the menu panel. Bad developer!. Apologies, user")
		}

		if item.displayStrings == nil {
			styledStr := item.displayString
			if item.opensMenu {
				styledStr = opensMenuStyle(styledStr)
			}
			stringArrays[i] = []string{styledStr}
		} else {
			stringArrays[i] = item.displayStrings
		}
	}

	list := utils.RenderDisplayStrings(stringArrays)

	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(false, list)
	menuView, _ := gui.g.SetView("menu", x0, y0, x1, y1, 0)
	menuView.Title = title
	menuView.FgColor = theme.GocuiDefaultTextColor
	menuView.SetOnSelectItem(gui.onSelectItemWrapper(func(selectedLine int) error {
		return nil
	}))
	menuView.SetContent(list)
	gui.State.Panels.Menu.SelectedLineIdx = 0

	return gui.pushContext(gui.State.Contexts.Menu)
}

func (gui *Gui) onMenuPress() error {
	selectedLine := gui.State.Panels.Menu.SelectedLineIdx
	if err := gui.returnFromContext(); err != nil {
		return err
	}

	if err := gui.State.MenuItems[selectedLine].onPress(); err != nil {
		return err
	}

	return nil
}
