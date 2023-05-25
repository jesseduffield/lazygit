package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

func (gui *Gui) handleOpenSearch(viewName string) error {
	view, err := gui.g.View(viewName)
	if err != nil {
		return nil
	}

	gui.State.Searching.isSearching = true
	gui.State.Searching.view = view

	gui.Views.Search.ClearTextArea()

	if err := gui.c.PushContext(gui.State.Contexts.Search); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleSearch() error {
	gui.State.Searching.searchString = gui.Views.Search.TextArea.GetContent()
	if err := gui.c.PopContext(); err != nil {
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
	keybindingConfig := gui.c.UserConfig.Keybinding

	return func(y int, index int, total int) error {
		if total == 0 {
			gui.c.SetViewContent(
				gui.Views.Search,
				fmt.Sprintf(
					gui.Tr.NoMatchesFor,
					gui.State.Searching.searchString,
					theme.OptionsFgColor.Sprintf(gui.Tr.ExitSearchMode, keybindings.Label(keybindingConfig.Universal.Return)),
				),
			)
			return nil
		}
		gui.c.SetViewContent(
			gui.Views.Search,
			fmt.Sprintf(
				gui.Tr.MatchesFor,
				gui.State.Searching.searchString,
				index+1,
				total,
				theme.OptionsFgColor.Sprintf(
					gui.Tr.SearchKeybindings,
					keybindings.Label(keybindingConfig.Universal.NextMatch),
					keybindings.Label(keybindingConfig.Universal.PrevMatch),
					keybindings.Label(keybindingConfig.Universal.Return),
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

	if err := gui.c.PopContext(); err != nil {
		return err
	}

	return nil
}
