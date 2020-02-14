package gui

import (
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) handleCreateDiscardMenu(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return err
		}
		return nil
	}

	menuItems := []*menuItem{
		{
			displayString: gui.Tr.SLocalize("discardAllChanges"),
			onPress: func() error {
				if err := gui.GitCommand.DiscardAllFileChanges(file); err != nil {
					return err
				}
				return gui.refreshFiles()
			},
		},
	}

	if file.HasStagedChanges && file.HasUnstagedChanges {
		menuItems = append(menuItems, &menuItem{
			displayString: gui.Tr.SLocalize("discardUnstagedChanges"),
			onPress: func() error {
				if err := gui.GitCommand.DiscardUnstagedFileChanges(file); err != nil {
					return err
				}

				return gui.refreshFiles()
			},
		})
	}

	return gui.createMenu(file.Name, menuItems, createMenuOptions{showCancel: true})
}
