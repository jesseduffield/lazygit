package gui

import (
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) handleCreateDiscardMenu(g *gocui.Gui, v *gocui.View) error {
	file, dir, err := gui.getSelectedDirOrFile()
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return err
		}
		return nil
	}

	if file != nil {
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
					if err := gui.GitCommand.DiscardUnstagedFileChanges(file.Name); err != nil {
						return gui.surfaceError(err)
					}

					return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{FILES}})
				},
			})
		}

		return gui.createMenu(file.Name, menuItems, createMenuOptions{showCancel: true})
	} else if dir != nil {
		menuItems := []*menuItem{
			// {
			// 	displayString: gui.Tr.SLocalize("discardAllChanges"),
			// 	onPress: func() error {
			// 		if err := gui.GitCommand.DiscardAllFileChanges(file); err != nil {
			// 			return gui.surfaceError(err)
			// 		}
			// 		return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{FILES}})
			// 	},
			// },
		}

		if dir.HasStagedChanges && dir.HasUnstagedChanges {
			menuItems = append(menuItems, &menuItem{
				displayString: gui.Tr.SLocalize("discardUnstagedChanges"),
				onPress: func() error {
					if err := gui.GitCommand.DiscardUnstagedFileChanges(dir.AbsolutePath()); err != nil {
						return gui.surfaceError(err)
					}

					return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{FILES}})
				},
			})
		}

		return gui.createMenu(dir.Name, menuItems, createMenuOptions{showCancel: true})
	}

	return nil
}
