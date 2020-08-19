package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) handleOpenSearch(g *gocui.Gui, v *gocui.View) error {
	gui.State.Searching.isSearching = true
	gui.State.Searching.view = v

	gui.renderString("search", "")

	if err := gui.switchContext(gui.Contexts.Search.Context); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleSearch(g *gocui.Gui, v *gocui.View) error {
	gui.State.Searching.searchString = gui.getSearchView().Buffer()
	if err := gui.switchContextToView(gui.State.Searching.view.Name()); err != nil {
		return err
	}

	if err := gui.State.Searching.view.Search(gui.State.Searching.searchString); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) onSelectItemWrapper(innerFunc func(int) error) func(int, int, int) error {
	return func(y int, index int, total int) error {
		if total == 0 {
			gui.renderString(
				"search",
				fmt.Sprintf(
					"no matches for '%s' %s",
					gui.State.Searching.searchString,
					utils.ColoredString(
						fmt.Sprintf("%s: exit search mode", gui.getKeyDisplay("universal.return")),
						theme.OptionsFgColor,
					),
				),
			)
			return nil
		}
		gui.renderString(
			"search",
			fmt.Sprintf(
				"matches for '%s' (%d of %d) %s",
				gui.State.Searching.searchString,
				index+1,
				total,
				utils.ColoredString(
					fmt.Sprintf(
						"%s: next match, %s: previous match, %s: exit search mode",
						gui.getKeyDisplay("universal.nextMatch"),
						gui.getKeyDisplay("universal.prevMatch"),
						gui.getKeyDisplay("universal.return"),
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

func (gui *Gui) handleSearchEscape(g *gocui.Gui, v *gocui.View) error {
	if err := gui.switchContextToView(gui.State.Searching.view.Name()); err != nil {
		return err
	}

	return gui.onSearchEscape()
}
