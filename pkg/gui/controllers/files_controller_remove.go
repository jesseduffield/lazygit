package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// splitting this action out into its own file because it's self-contained

func (self *FilesController) remove(node *filetree.FileNode) error {
	var menuItems []*types.MenuItem
	if node.File == nil {
		menuItems = []*types.MenuItem{
			{
				DisplayString: self.c.Tr.LcDiscardAllChanges,
				OnPress: func() error {
					self.c.LogAction(self.c.Tr.Actions.DiscardAllChangesInDirectory)
					if err := self.git.WorkingTree.DiscardAllDirChanges(node); err != nil {
						return self.c.Error(err)
					}
					return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
				},
			},
		}

		if node.GetHasStagedChanges() && node.GetHasUnstagedChanges() {
			menuItems = append(menuItems, &types.MenuItem{
				DisplayString: self.c.Tr.LcDiscardUnstagedChanges,
				OnPress: func() error {
					self.c.LogAction(self.c.Tr.Actions.DiscardUnstagedChangesInDirectory)
					if err := self.git.WorkingTree.DiscardUnstagedDirChanges(node); err != nil {
						return self.c.Error(err)
					}

					return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
				},
			})
		}
	} else {
		file := node.File

		submodules := self.getSubmodules()
		if file.IsSubmodule(submodules) {
			submodule := file.SubmoduleConfig(submodules)

			menuItems = []*types.MenuItem{
				{
					DisplayString: self.c.Tr.LcSubmoduleStashAndReset,
					OnPress: func() error {
						return self.ResetSubmodule(submodule)
					},
				},
			}
		} else {
			menuItems = []*types.MenuItem{
				{
					DisplayString: self.c.Tr.LcDiscardAllChanges,
					OnPress: func() error {
						self.c.LogAction(self.c.Tr.Actions.DiscardAllChangesInFile)
						if err := self.git.WorkingTree.DiscardAllFileChanges(file); err != nil {
							return self.c.Error(err)
						}
						return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
					},
				},
			}

			if file.HasStagedChanges && file.HasUnstagedChanges {
				menuItems = append(menuItems, &types.MenuItem{
					DisplayString: self.c.Tr.LcDiscardUnstagedChanges,
					OnPress: func() error {
						self.c.LogAction(self.c.Tr.Actions.DiscardAllUnstagedChangesInFile)
						if err := self.git.WorkingTree.DiscardUnstagedFileChanges(file); err != nil {
							return self.c.Error(err)
						}

						return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
					},
				})
			}
		}
	}

	return self.c.Menu(types.CreateMenuOptions{Title: node.GetPath(), Items: menuItems})
}
