package gui

func (gui *Gui) handleCreateDiscardMenu() error {
	node := gui.getSelectedFileNode()
	if node == nil {
		return nil
	}

	var menuItems []*menuItem
	if node.File == nil {
		menuItems = []*menuItem{
			{
				displayString: gui.Tr.LcDiscardAllChanges,
				onPress: func() error {
					if err := gui.GitCommand.DiscardAllDirChanges(node); err != nil {
						return gui.surfaceError(err)
					}
					return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []RefreshableView{FILES}})
				},
			},
		}

		if node.GetHasStagedChanges() && node.GetHasUnstagedChanges() {
			menuItems = append(menuItems, &menuItem{
				displayString: gui.Tr.LcDiscardUnstagedChanges,
				onPress: func() error {
					if err := gui.GitCommand.DiscardUnstagedDirChanges(node); err != nil {
						return gui.surfaceError(err)
					}

					return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []RefreshableView{FILES}})
				},
			})
		}
	} else {
		file := node.File

		submodules := gui.State.Submodules
		if file.IsSubmodule(submodules) {
			submodule := file.SubmoduleConfig(submodules)

			menuItems = []*menuItem{
				{
					displayString: gui.Tr.LcSubmoduleStashAndReset,
					onPress: func() error {
						return gui.handleResetSubmodule(submodule)
					},
				},
			}
		} else {
			menuItems = []*menuItem{
				{
					displayString: gui.Tr.LcDiscardAllChanges,
					onPress: func() error {
						if err := gui.GitCommand.DiscardAllFileChanges(file); err != nil {
							return gui.surfaceError(err)
						}
						return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []RefreshableView{FILES}})
					},
				},
			}

			if file.HasStagedChanges && file.HasUnstagedChanges {
				menuItems = append(menuItems, &menuItem{
					displayString: gui.Tr.LcDiscardUnstagedChanges,
					onPress: func() error {
						if err := gui.GitCommand.DiscardUnstagedFileChanges(file); err != nil {
							return gui.surfaceError(err)
						}

						return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []RefreshableView{FILES}})
					},
				})
			}
		}
	}

	return gui.createMenu(node.GetPath(), menuItems, createMenuOptions{showCancel: true})
}
