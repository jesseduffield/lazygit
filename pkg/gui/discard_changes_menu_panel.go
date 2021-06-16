package gui

import . "github.com/jesseduffield/lazygit/pkg/gui/types"

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
					if err := gui.Git.WithSpan(gui.Tr.Spans.DiscardAllChangesInDirectory).Worktree().DiscardAllDirChanges(node); err != nil {
						return gui.SurfaceError(err)
					}
					return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC, Scope: []RefreshableView{FILES}})
				},
			},
		}

		if node.GetHasStagedChanges() && node.GetHasUnstagedChanges() {
			menuItems = append(menuItems, &menuItem{
				displayString: gui.Tr.LcDiscardUnstagedChanges,
				onPress: func() error {
					if err := gui.Git.WithSpan(gui.Tr.Spans.DiscardUnstagedChangesInDirectory).Worktree().DiscardUnstagedDirChanges(node); err != nil {
						return gui.SurfaceError(err)
					}

					return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC, Scope: []RefreshableView{FILES}})
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
						if err := gui.Git.WithSpan(gui.Tr.Spans.DiscardAllChangesInFile).Worktree().DiscardAllFileChanges(file); err != nil {
							return gui.SurfaceError(err)
						}
						return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC, Scope: []RefreshableView{FILES}})
					},
				},
			}

			if file.HasStagedChanges && file.HasUnstagedChanges {
				menuItems = append(menuItems, &menuItem{
					displayString: gui.Tr.LcDiscardUnstagedChanges,
					onPress: func() error {
						if err := gui.Git.WithSpan(gui.Tr.Spans.DiscardAllUnstagedChangesInFile).Worktree().DiscardUnstagedFileChanges(file); err != nil {
							return gui.SurfaceError(err)
						}

						return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC, Scope: []RefreshableView{FILES}})
					},
				})
			}
		}
	}

	return gui.createMenu(node.GetPath(), menuItems, createMenuOptions{showCancel: true})
}
