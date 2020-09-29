package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

func (gui *Gui) submoduleFromFile(file *models.File) *models.SubmoduleConfig {
	for _, config := range gui.State.Submodules {
		if config.Name == file.Name {
			return config
		}
	}

	return nil
}

func (gui *Gui) handleCreateDiscardMenu(g *gocui.Gui, v *gocui.View) error {
	file := gui.getSelectedFile()
	if file == nil {
		return nil
	}

	var menuItems []*menuItem

	submoduleConfigs := gui.State.Submodules
	if file.IsSubmodule(submoduleConfigs) {
		submoduleConfig := file.SubmoduleConfig(submoduleConfigs)

		menuItems = []*menuItem{
			{
				displayString: gui.Tr.SLocalize("submoduleStashAndReset"),
				onPress: func() error {
					if err := gui.GitCommand.UnStageFile(file.Name, file.Tracked); err != nil {
						return gui.surfaceError(err)
					}
					if err := gui.GitCommand.SubmoduleStash(submoduleConfig); err != nil {
						return gui.surfaceError(err)
					}
					if err := gui.GitCommand.SubmoduleReset(submoduleConfig); err != nil {
						return gui.surfaceError(err)
					}
					return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{FILES}})
				},
			},
		}
	} else {
		menuItems = []*menuItem{
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

	}

	return gui.createMenu(file.Name, menuItems, createMenuOptions{showCancel: true})
}
