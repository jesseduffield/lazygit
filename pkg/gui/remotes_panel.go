package gui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) getSelectedRemote() *commands.Remote {
	selectedLine := gui.State.Panels.Remotes.SelectedLine
	if selectedLine == -1 {
		return nil
	}

	return gui.State.Remotes[selectedLine]
}

func (gui *Gui) handleRemotesClick(g *gocui.Gui, v *gocui.View) error {
	itemCount := len(gui.State.Remotes)
	handleSelect := gui.handleRemoteSelect
	selectedLine := &gui.State.Panels.Remotes.SelectedLine

	return gui.handleClick(v, itemCount, selectedLine, handleSelect)
}

func (gui *Gui) handleRemoteSelect(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	gui.State.SplitMainPanel = false

	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	gui.getMainView().Title = "Remote"

	remote := gui.getSelectedRemote()
	gui.focusPoint(0, gui.State.Panels.Menu.SelectedLine, gui.State.MenuItemCount, v)
	if err := gui.focusPoint(0, gui.State.Panels.Remotes.SelectedLine, len(gui.State.Remotes), v); err != nil {
		return err
	}

	return gui.renderString(g, "main", fmt.Sprintf("%s\nUrls:\n%s", utils.ColoredString(remote.Name, color.FgGreen), strings.Join(remote.Urls, "\n")))
}

// gui.refreshStatus is called at the end of this because that's when we can
// be sure there is a state.Remotes array to pick the current remote from
func (gui *Gui) refreshRemotes() error {
	remotes, err := gui.GitCommand.GetRemotes()
	if err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
	}

	gui.State.Remotes = remotes

	gui.g.Update(func(g *gocui.Gui) error {
		gui.refreshSelectedLine(&gui.State.Panels.Remotes.SelectedLine, len(gui.State.Remotes))
		return nil
	})

	return nil
}
