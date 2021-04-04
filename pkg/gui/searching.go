package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) handleOpenSearch(viewName string) error {
	view, err := gui.g.View(viewName)
	if err != nil {
		return nil
	}

	gui.State.Searching.isSearching = true
	gui.State.Searching.view = view

	gui.renderString(gui.Views.Search, "")

	if err := gui.pushContext(gui.State.Contexts.Search); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleSearch() error {
	gui.State.Searching.searchString = gui.Views.Search.Buffer()
	gui.Log.Warn(gui.State.Searching.searchString)
	if err := gui.returnFromContext(); err != nil {
		return err
	}

	view := gui.State.Searching.view
	if view == nil {
		return nil
	}

	if err := view.Search(gui.State.Searching.searchString); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) onSelectItemWrapper(innerFunc func(int) error) func(int, int, int) error {
	keybindingConfig := gui.Config.GetUserConfig().Keybinding

	return func(y int, index int, total int) error {
		if total == 0 {
			gui.renderString(
				gui.Views.Search,
				fmt.Sprintf(
					"no matches for '%s' %s",
					gui.State.Searching.searchString,
					utils.ColoredString(
						fmt.Sprintf("%s: exit search mode", gui.getKeyDisplay(keybindingConfig.Universal.Return)),
						theme.OptionsFgColor,
					),
				),
			)
			return nil
		}
		gui.renderString(
			gui.Views.Search,
			fmt.Sprintf(
				"matches for '%s' (%d of %d) %s",
				gui.State.Searching.searchString,
				index+1,
				total,
				utils.ColoredString(
					fmt.Sprintf(
						"%s: next match, %s: previous match, %s: exit search mode",
						gui.getKeyDisplay(keybindingConfig.Universal.NextMatch),
						gui.getKeyDisplay(keybindingConfig.Universal.PrevMatch),
						gui.getKeyDisplay(keybindingConfig.Universal.Return),
					),
					theme.OptionsFgColor,
				),
			),
		)
		if err := innerFunc(y); err != nil {
			return err
		}
		return nil
	}
}

func (gui *Gui) onSearchEscape() error {
	gui.State.Searching.isSearching = false
	if gui.State.Searching.view != nil {
		gui.State.Searching.view.ClearSearch()
		gui.State.Searching.view = nil
	}

	return nil
}

func (gui *Gui) handleSearchEscape() error {
	if err := gui.onSearchEscape(); err != nil {
		return err
	}

	if err := gui.returnFromContext(); err != nil {
		return err
	}

	return nil
}
