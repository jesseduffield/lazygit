package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type menuItem struct {
	displayString  string
	displayStrings []string
	onPress        func() error
}

// list panel functions

func (gui *Gui) handleMenuSelect() error {
	return nil
}

// specific functions

func (gui *Gui) renderMenuOptions() error {
	optionsMap := map[string]string{
		gui.getKeyDisplay("universal.return"): gui.Tr.SLocalize("close"),
		fmt.Sprintf("%s %s", gui.getKeyDisplay("universal.prevItem"), gui.getKeyDisplay("universal.nextItem")): gui.Tr.SLocalize("navigate"),
		gui.getKeyDisplay("universal.select"): gui.Tr.SLocalize("execute"),
	}
	return gui.renderOptionsMap(optionsMap)
}

func (gui *Gui) menuConfirmationKeys() []interface{} {
	return []interface{}{gui.getKey("universal.select"), gui.getKey("universal.confirm"), gui.getKey("universal.confirm-alt1")}
}

func (gui *Gui) handleMenuClose(g *gocui.Gui, v *gocui.View) error {
	for _, key := range gui.menuConfirmationKeys() {
		if err := g.DeleteKeybinding("menu", key, gocui.ModNone); err != nil {
			return err
		}
	}
	err := g.DeleteView("menu")
	if err != nil {
		return err
	}
	return gui.returnFromContext()
}

type createMenuOptions struct {
	showCancel bool
}

func (gui *Gui) createMenu(title string, items []*menuItem, createMenuOptions createMenuOptions) error {
	if createMenuOptions.showCancel {
		// this is mutative but I'm okay with that for now
		items = append(items, &menuItem{
			displayStrings: []string{gui.Tr.SLocalize("cancel")},
			onPress: func() error {
				return nil
			},
		})
	}

	gui.State.MenuItemCount = len(items)

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
	gui.State.Panels.Menu.SelectedLine = 0

	wrappedHandlePress := func(g *gocui.Gui, v *gocui.View) error {
		selectedLine := gui.State.Panels.Menu.SelectedLine
		if err := items[selectedLine].onPress(); err != nil {
			return err
		}

		_, _ = gui.g.SetViewOnBottom("menu")

		return gui.returnFromContext()
	}

	gui.State.Panels.Menu.OnPress = wrappedHandlePress

	for _, key := range gui.menuConfirmationKeys() {
		_ = gui.g.DeleteKeybinding("menu", key, gocui.ModNone)

		if err := gui.g.SetKeybinding("menu", nil, key, gocui.ModNone, wrappedHandlePress); err != nil {
			return err
		}
	}

	gui.g.Update(func(g *gocui.Gui) error {
		return gui.switchContext(gui.Contexts.Menu.Context)
	})
	return nil
}
