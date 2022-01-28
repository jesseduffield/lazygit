package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) handleCreateDiscardMenu() error {
	node := gui.getSelectedFileNode()
	if node == nil {
		return nil
	}

	var menuItems []*popup.MenuItem
	if node.File == nil {
		menuItems = []*popup.MenuItem{
			{
				DisplayString: gui.Tr.LcDiscardAllChanges,
				OnPress: func() error {
					gui.logAction(gui.Tr.Actions.DiscardAllChangesInDirectory)
					if err := gui.Git.WorkingTree.DiscardAllDirChanges(node); err != nil {
						return gui.PopupHandler.Error(err)
					}
					return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
				},
			},
		}

		if node.GetHasStagedChanges() && node.GetHasUnstagedChanges() {
			menuItems = append(menuItems, &popup.MenuItem{
				DisplayString: gui.Tr.LcDiscardUnstagedChanges,
				OnPress: func() error {
					gui.logAction(gui.Tr.Actions.DiscardUnstagedChangesInDirectory)
					if err := gui.Git.WorkingTree.DiscardUnstagedDirChanges(node); err != nil {
						return gui.PopupHandler.Error(err)
					}

					return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
				},
			})
		}
	} else {
		file := node.File

		submodules := gui.State.Submodules
		if file.IsSubmodule(submodules) {
			submodule := file.SubmoduleConfig(submodules)

			menuItems = []*popup.MenuItem{
				{
					DisplayString: gui.Tr.LcSubmoduleStashAndReset,
					OnPress: func() error {
						return gui.resetSubmodule(submodule)
					},
				},
			}
		} else {
			menuItems = []*popup.MenuItem{
				{
					DisplayString: gui.Tr.LcDiscardAllChanges,
					OnPress: func() error {
						gui.logAction(gui.Tr.Actions.DiscardAllChangesInFile)
						if err := gui.Git.WorkingTree.DiscardAllFileChanges(file); err != nil {
							return gui.PopupHandler.Error(err)
						}
						return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
					},
				},
			}

			if file.HasStagedChanges && file.HasUnstagedChanges {
				menuItems = append(menuItems, &popup.MenuItem{
					DisplayString: gui.Tr.LcDiscardUnstagedChanges,
					OnPress: func() error {
						gui.logAction(gui.Tr.Actions.DiscardAllUnstagedChangesInFile)
						if err := gui.Git.WorkingTree.DiscardUnstagedFileChanges(file); err != nil {
							return gui.PopupHandler.Error(err)
						}

						return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
					},
				})
			}
		}
	}

	return gui.PopupHandler.Menu(popup.CreateMenuOptions{Title: node.GetPath(), Items: menuItems})
}
