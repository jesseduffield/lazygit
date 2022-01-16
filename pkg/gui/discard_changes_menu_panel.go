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
				DisplayString: gui.c.Tr.LcDiscardAllChanges,
				OnPress: func() error {
					gui.c.LogAction(gui.c.Tr.Actions.DiscardAllChangesInDirectory)
					if err := gui.git.WorkingTree.DiscardAllDirChanges(node); err != nil {
						return gui.c.Error(err)
					}
					return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
				},
			},
		}

		if node.GetHasStagedChanges() && node.GetHasUnstagedChanges() {
			menuItems = append(menuItems, &popup.MenuItem{
				DisplayString: gui.c.Tr.LcDiscardUnstagedChanges,
				OnPress: func() error {
					gui.c.LogAction(gui.c.Tr.Actions.DiscardUnstagedChangesInDirectory)
					if err := gui.git.WorkingTree.DiscardUnstagedDirChanges(node); err != nil {
						return gui.c.Error(err)
					}

					return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
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
					DisplayString: gui.c.Tr.LcSubmoduleStashAndReset,
					OnPress: func() error {
						return gui.Controllers.Files.ResetSubmodule(submodule)
					},
				},
			}
		} else {
			menuItems = []*popup.MenuItem{
				{
					DisplayString: gui.c.Tr.LcDiscardAllChanges,
					OnPress: func() error {
						gui.c.LogAction(gui.c.Tr.Actions.DiscardAllChangesInFile)
						if err := gui.git.WorkingTree.DiscardAllFileChanges(file); err != nil {
							return gui.c.Error(err)
						}
						return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
					},
				},
			}

			if file.HasStagedChanges && file.HasUnstagedChanges {
				menuItems = append(menuItems, &popup.MenuItem{
					DisplayString: gui.c.Tr.LcDiscardUnstagedChanges,
					OnPress: func() error {
						gui.c.LogAction(gui.c.Tr.Actions.DiscardAllUnstagedChangesInFile)
						if err := gui.git.WorkingTree.DiscardUnstagedFileChanges(file); err != nil {
							return gui.c.Error(err)
						}

						return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
					},
				})
			}
		}
	}

	return gui.c.Menu(popup.CreateMenuOptions{Title: node.GetPath(), Items: menuItems})
}
