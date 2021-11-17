package gui

import (
	"errors"
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
	// only applies when displayString is used
	opensMenu bool
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

func (gui *Gui) handleMenuClose() error {
	return gui.returnFromContextSync()
}

type createMenuOptions struct {
	showCancel     bool
	allowFiltering bool
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

	menuView, renderError := gui.buildMenuView(title, items)
	if renderError != nil {
		return renderError
	}

	if createMenuOptions.allowFiltering {
		filter := MenuPanelFilter{
			menuView:  menuView,
			menuItems: items,
			updateMenu: func(needle string, filteredItems []*menuItem) {
				menuView, _ := gui.buildMenuView(title, filteredItems)
				menuView.Search(needle)
			},
		}

		menuView.Editable = true

		menuView.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
			switch key {
			case gocui.KeyBackspace:
			case gocui.KeyBackspace2:
				filter.HandleSearchBackspace()
			case gocui.KeyDelete:
				filter.HandleResetSearch()
			default:
				char := string(ch)
				filter.HandleSearchKeystroke(char)
			}
			return false
		})
	}

	gui.g.Update(func(g *gocui.Gui) error {
		return gui.pushContext(gui.State.Contexts.Menu)
	})
	return nil
}

// buildMenuView formats a given set of menuItems and creates a menuView
func (gui *Gui) buildMenuView(title string, items []*menuItem) (*gocui.View, error) {
	list, err := formatListItems(items)
	if err != nil {
		return nil, err
	}

	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(false, list)
	menuView, _ := gui.g.SetView("menu", x0, y0, x1, y1, 0)
	menuView.Title = title
	menuView.FgColor = theme.GocuiDefaultTextColor
	menuView.SetOnSelectItem(gui.onSelectItemWrapper(func(selectedLine int) error {
		return nil
	}))
	menuView.SetContent(list)
	gui.State.MenuItems = items
	gui.State.Panels.Menu.SelectedLineIdx = 0
	return menuView, nil
}

// formatListItems formats a set of menuItems to a format string suitable for a menu
func formatListItems(items []*menuItem) (string, error) {
	stringArrays := make([][]string, len(items))
	for i, item := range items {
		if item.opensMenu && item.displayStrings != nil {
			return "", errors.New("Message for the developer of this app: you've set opensMenu with displaystrings on the menu panel. Bad developer!. Apologies, user")
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
	return list, nil
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
