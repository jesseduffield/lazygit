package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// splitting this action out into its own file because it's self-contained

type FilesRemoveController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &FilesRemoveController{}

func NewFilesRemoveController(
	common *controllerCommon,
) *FilesRemoveController {
	return &FilesRemoveController{
		baseController:   baseController{},
		controllerCommon: common,
	}
}

func (self *FilesRemoveController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.Remove),
			Handler:     self.checkSelectedFileNode(self.remove),
			Description: self.c.Tr.LcViewDiscardOptions,
			OpensMenu:   true,
		},
	}

	return bindings
}

func (self *FilesRemoveController) remove(node *filetree.FileNode) error {
	var menuItems []*types.MenuItem
	if node.File == nil {
		menuItems = []*types.MenuItem{
			{
				Label: self.c.Tr.LcDiscardAllChanges,
				OnPress: func() error {
					self.c.LogAction(self.c.Tr.Actions.DiscardAllChangesInDirectory)
					if err := self.git.WorkingTree.DiscardAllDirChanges(node); err != nil {
						return self.c.Error(err)
					}
					return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
				},
				Key: 'd',
			},
		}

		if node.GetHasStagedChanges() && node.GetHasUnstagedChanges() {
			menuItems = append(menuItems, &types.MenuItem{
				Label: self.c.Tr.LcDiscardUnstagedChanges,
				OnPress: func() error {
					self.c.LogAction(self.c.Tr.Actions.DiscardUnstagedChangesInDirectory)
					if err := self.git.WorkingTree.DiscardUnstagedDirChanges(node); err != nil {
						return self.c.Error(err)
					}

					return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
				},
				Key: 'u',
			})
		}
	} else {
		file := node.File

		submodules := self.model.Submodules
		if file.IsSubmodule(submodules) {
			submodule := file.SubmoduleConfig(submodules)

			menuItems = []*types.MenuItem{
				{
					Label: self.c.Tr.LcSubmoduleStashAndReset,
					OnPress: func() error {
						return self.ResetSubmodule(submodule)
					},
				},
			}
		} else {
			menuItems = []*types.MenuItem{
				{
					Label: self.c.Tr.LcDiscardAllChanges,
					OnPress: func() error {
						self.c.LogAction(self.c.Tr.Actions.DiscardAllChangesInFile)
						if err := self.git.WorkingTree.DiscardAllFileChanges(file); err != nil {
							return self.c.Error(err)
						}
						return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
					},
					Key: 'd',
				},
			}

			if file.HasStagedChanges && file.HasUnstagedChanges {
				menuItems = append(menuItems, &types.MenuItem{
					Label: self.c.Tr.LcDiscardUnstagedChanges,
					OnPress: func() error {
						self.c.LogAction(self.c.Tr.Actions.DiscardAllUnstagedChangesInFile)
						if err := self.git.WorkingTree.DiscardUnstagedFileChanges(file); err != nil {
							return self.c.Error(err)
						}

						return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
					},
					Key: 'u',
				})
			}
		}
	}

	return self.c.Menu(types.CreateMenuOptions{Title: node.GetPath(), Items: menuItems})
}

func (self *FilesRemoveController) ResetSubmodule(submodule *models.SubmoduleConfig) error {
	return self.c.WithWaitingStatus(self.c.Tr.LcResettingSubmoduleStatus, func() error {
		self.c.LogAction(self.c.Tr.Actions.ResetSubmodule)

		file := self.helpers.WorkingTree.FileForSubmodule(submodule)
		if file != nil {
			if err := self.git.WorkingTree.UnStageFile(file.Names(), file.Tracked); err != nil {
				return self.c.Error(err)
			}
		}

		if err := self.git.Submodule.Stash(submodule); err != nil {
			return self.c.Error(err)
		}
		if err := self.git.Submodule.Reset(submodule); err != nil {
			return self.c.Error(err)
		}

		return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES, types.SUBMODULES}})
	})
}

func (self *FilesRemoveController) checkSelectedFileNode(callback func(*filetree.FileNode) error) func() error {
	return func() error {
		node := self.context().GetSelected()
		if node == nil {
			return nil
		}

		return callback(node)
	}
}

func (self *FilesRemoveController) Context() types.Context {
	return self.context()
}

func (self *FilesRemoveController) context() *context.WorkingTreeContext {
	return self.contexts.Files
}
