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

func (gui *Gui) handleMenuSelect(g *gocui.Gui, v *gocui.View) error {
	v.FocusPoint(0, gui.State.Panels.Menu.SelectedLine)
	return nil
}

// specific functions

func (gui *Gui) renderMenuOptions() error {
	optionsMap := map[string]string{
		fmt.Sprintf("%s/%s", gui.getKeyDisplay("universal.return"), gui.getKeyDisplay("universal.quit")):       gui.Tr.SLocalize("close"),
		fmt.Sprintf("%s %s", gui.getKeyDisplay("universal.prevItem"), gui.getKeyDisplay("universal.nextItem")): gui.Tr.SLocalize("navigate"),
		gui.getKeyDisplay("universal.select"): gui.Tr.SLocalize("execute"),
	}
	return gui.renderOptionsMap(optionsMap)
}

func (gui *Gui) handleMenuClose(g *gocui.Gui, v *gocui.View) error {
	for _, key := range []gocui.Key{gocui.KeySpace, gocui.KeyEnter} {
		if err := g.DeleteKeybinding("menu", key, gocui.ModNone); err != nil {
			return err
		}
	}
	err := g.DeleteView("menu")
	if err != nil {
		return err
	}
	return gui.returnFocus(g, v)
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

	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(gui.g, false, list)
	menuView, _ := gui.g.SetView("menu", x0, y0, x1, y1, 0)
	menuView.Title = title
	menuView.FgColor = theme.GocuiDefaultTextColor
	menuView.ContainsList = true
	menuView.Clear()
	fmt.Fprint(menuView, list)
	gui.State.Panels.Menu.SelectedLine = 0

	wrappedHandlePress := func(g *gocui.Gui, v *gocui.View) error {
		selectedLine := gui.State.Panels.Menu.SelectedLine
		if err := items[selectedLine].onPress(); err != nil {
			return err
		}

		if _, err := gui.g.View("menu"); err == nil {
			if _, err := gui.g.SetViewOnBottom("menu"); err != nil {
				return err
			}
		}

		return gui.returnFocus(gui.g, menuView)
	}

	gui.State.Panels.Menu.OnPress = wrappedHandlePress

	for _, key := range []gocui.Key{gocui.KeySpace, gocui.KeyEnter, 'y'} {
		_ = gui.g.DeleteKeybinding("menu", key, gocui.ModNone)

		if err := gui.g.SetKeybinding("menu", nil, key, gocui.ModNone, wrappedHandlePress); err != nil {
			return err
		}
	}

	gui.g.Update(func(g *gocui.Gui) error {
		if _, err := gui.g.View("menu"); err == nil {
			if _, err := g.SetViewOnTop("menu"); err != nil {
				return err
			}
		}
		currentView := gui.g.CurrentView()
		return gui.switchFocus(gui.g, currentView, menuView)
	})
	return nil
}
