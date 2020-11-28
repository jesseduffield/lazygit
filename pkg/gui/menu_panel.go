package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type menuItem struct {
	displayString  string
	displayStrings []string
	onPress        func() error
}

// every item in a list context needs an ID
func (i *menuItem) ID() string {
	if i.displayString != "" {
		return i.displayString
	}

	return strings.Join(i.displayStrings, "-")
}

// list panel functions

func (gui *Gui) handleMenuSelect() error {
	return nil
}

// specific functions

func (gui *Gui) getMenuOptions() map[string]string {
	keybindingConfig := gui.Config.GetUserConfig().Keybinding

	return map[string]string{
		gui.getKeyDisplay(keybindingConfig.Universal.Return): gui.Tr.LcClose,
		fmt.Sprintf("%s %s", gui.getKeyDisplay(keybindingConfig.Universal.PrevItem), gui.getKeyDisplay(keybindingConfig.Universal.NextItem)): gui.Tr.LcNavigate,
		gui.getKeyDisplay(keybindingConfig.Universal.Select): gui.Tr.LcExecute,
	}
}

func (gui *Gui) handleMenuClose(g *gocui.Gui, v *gocui.View) error {
	_ = g.DeleteView("menu")
	return gui.returnFromContext()
}

type createMenuOptions struct {
	showCancel bool
}

func (gui *Gui) createMenu(title string, items []*menuItem, createMenuOptions createMenuOptions) error {
	if createMenuOptions.showCancel {
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
		if item.displayStrings == nil {
			stringArrays[i] = []string{item.displayString}
		} else {
			stringArrays[i] = item.displayStrings
		}
	}

	list := utils.RenderDisplayStrings(stringArrays)

	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(false, list)
	menuView, _ := gui.g.SetView("menu", x0, y0, x1, y1, 0)
	menuView.Title = title
	menuView.FgColor = theme.GocuiDefaultTextColor
	menuView.ContainsList = true
	menuView.Clear()
	menuView.SetOnSelectItem(gui.onSelectItemWrapper(func(selectedLine int) error {
		return nil
	}))
	fmt.Fprint(menuView, list)
	gui.State.Panels.Menu.SelectedLineIdx = 0

	gui.g.Update(func(g *gocui.Gui) error {
		return gui.pushContext(gui.Contexts.Menu.Context)
	})
	return nil
}

func (gui *Gui) onMenuPress() error {
	selectedLine := gui.State.Panels.Menu.SelectedLineIdx
	if err := gui.State.MenuItems[selectedLine].onPress(); err != nil {
		return err
	}

	return gui.returnFromContext()
}
