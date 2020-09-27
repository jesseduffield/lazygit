package gui

import (
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) handleCreateDiscardMenu(g *gocui.Gui, v *gocui.View) error {
	file := gui.getSelectedFile()
	if file == nil {
		return nil
	}

	if file.IsSubmodule {
		// git submodule foreach '[[ "$name" == "renderers/chartify" ]] && git stash --include-untracked'
		// git submodule update --force renderers/chartify
	}

	menuItems := []*menuItem{
		{
			displayString: gui.Tr.SLocalize("discardAllChanges"),
			onPress: func() error {
				if err := gui.GitCommand.DiscardAllFileChanges(file); err != nil {
					return gui.surfaceError(err)
				}
				return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{FILES}})
			},
		},
	}

	if file.HasStagedChanges && file.HasUnstagedChanges {
		menuItems = append(menuItems, &menuItem{
			displayString: gui.Tr.SLocalize("discardUnstagedChanges"),
			onPress: func() error {
				if err := gui.GitCommand.DiscardUnstagedFileChanges(file); err != nil {
					return gui.surfaceError(err)
				}

				return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{FILES}})
			},
		})
	}

	return gui.createMenu(file.Name, menuItems, createMenuOptions{showCancel: true})
}
