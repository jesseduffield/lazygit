package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

func (gui *Gui) handleOpenSearch(viewName string) error {
	view, err := gui.g.View(viewName)
	if err != nil {
		return nil
	}

	gui.State.Modes.Searching.OnSearchPrompt(view, gui.c.CurrentContext().GetKey())

	gui.Views.Search.ClearTextArea()

	if err := gui.c.PushContext(gui.State.Contexts.Search); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleSearch() error {
	needle := gui.Views.Search.TextArea.GetContent()
	if err := gui.State.Modes.Searching.OnSearch(needle); err != nil {
		return err
	}

	if err := gui.c.PopContext(); err != nil {
		return err
	}

	return gui.c.PostRefreshUpdate(gui.currentContext())
}

func (gui *Gui) onSelectItemWrapper(innerFunc func(int) error) func(int, int, int) error {
	keybindingConfig := gui.c.UserConfig.Keybinding

	return func(y int, index int, total int) error {
		if total == 0 {
			return gui.renderString(
				gui.Views.Search,
				fmt.Sprintf(
					"no matches for '%s' %s",
					gui.State.Modes.Searching.GetSearchString(),
					theme.OptionsFgColor.Sprintf("%s: exit search mode", gui.getKeyDisplay(keybindingConfig.Universal.Return)),
				),
			)
		}
		_ = gui.renderString(
			gui.Views.Search,
			fmt.Sprintf(
				"matches for '%s' (%d of %d) %s",
				gui.State.Modes.Searching.GetSearchString(),
				index+1,
				total,
				theme.OptionsFgColor.Sprintf(
					"%s: next match, %s: previous match, %s: exit search mode",
					gui.getKeyDisplay(keybindingConfig.Universal.NextMatch),
					gui.getKeyDisplay(keybindingConfig.Universal.PrevMatch),
					gui.getKeyDisplay(keybindingConfig.Universal.Return),
				),
			),
		)
		if err := innerFunc(y); err != nil {
			return err
		}
		return nil
	}
}

func (gui *Gui) onSearchEscape() {
	gui.State.Modes.Searching.Escape()
}

func (gui *Gui) handleSearchPromptEscape() error {
	if err := gui.c.PopContext(); err != nil {
		return err
	}

	return gui.exitSearch()
}

func (gui *Gui) exitSearch() error {
	gui.onSearchEscape()

	return gui.c.PostRefreshUpdate(gui.currentContext())
}

func (gui *Gui) searchEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	matched := gui.handleEditorKeypress(v.TextArea, key, ch, mod, false)

	v.RenderTextArea()

	if err := gui.State.Modes.Searching.OnSearch(gui.Views.Search.TextArea.GetContent()); err != nil {
		gui.Log.Error(err)
	}

	if parentContext, ok := gui.parentContext(); ok {
		if err := gui.c.PostRefreshUpdate(parentContext); err != nil {
			gui.Log.Error(err)
		}
	}

	return matched
}
