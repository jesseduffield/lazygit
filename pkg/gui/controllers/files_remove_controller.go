package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// splitting this action out into its own file because it's self-contained

type FilesRemoveController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &FilesRemoveController{}

func NewFilesRemoveController(
	common *ControllerCommon,
) *FilesRemoveController {
	return &FilesRemoveController{
		baseController: baseController{},
		c:              common,
	}
}

func (self *FilesRemoveController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.Remove),
			Handler:     self.checkSelectedFileNode(self.remove),
			Description: self.c.Tr.ViewDiscardOptions,
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
				Label: self.c.Tr.DiscardAllChanges,
				OnPress: func() error {
					self.c.LogAction(self.c.Tr.Actions.DiscardAllChangesInDirectory)
					if err := self.c.Git().WorkingTree.DiscardAllDirChanges(node); err != nil {
						return self.c.Error(err)
					}
					return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
				},
				Key: 'x',
				Tooltip: utils.ResolvePlaceholderString(
					self.c.Tr.DiscardAllTooltip,
					map[string]string{
						"path": node.GetPath(),
					},
				),
			},
		}

		if node.GetHasStagedChanges() && node.GetHasUnstagedChanges() {
			menuItems = append(menuItems, &types.MenuItem{
				Label: self.c.Tr.DiscardUnstagedChanges,
				OnPress: func() error {
					self.c.LogAction(self.c.Tr.Actions.DiscardUnstagedChangesInDirectory)
					if err := self.c.Git().WorkingTree.DiscardUnstagedDirChanges(node); err != nil {
						return self.c.Error(err)
					}

					return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
				},
				Key: 'u',
				Tooltip: utils.ResolvePlaceholderString(
					self.c.Tr.DiscardUnstagedTooltip,
					map[string]string{
						"path": node.GetPath(),
					},
				),
			})
		}
	} else {
		file := node.File

		submodules := self.c.Model().Submodules
		if file.IsSubmodule(submodules) {
			submodule := file.SubmoduleConfig(submodules)

			menuItems = []*types.MenuItem{
				{
					Label: self.c.Tr.SubmoduleStashAndReset,
					OnPress: func() error {
						return self.ResetSubmodule(submodule)
					},
				},
			}
		} else {
			menuItems = []*types.MenuItem{
				{
					Label: self.c.Tr.DiscardAllChanges,
					OnPress: func() error {
						self.c.LogAction(self.c.Tr.Actions.DiscardAllChangesInFile)
						if err := self.c.Git().WorkingTree.DiscardAllFileChanges(file); err != nil {
							return self.c.Error(err)
						}
						return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
					},
					Key: 'x',
					Tooltip: utils.ResolvePlaceholderString(
						self.c.Tr.DiscardAllTooltip,
						map[string]string{
							"path": node.GetPath(),
						},
					),
				},
			}

			if file.HasStagedChanges && file.HasUnstagedChanges {
				menuItems = append(menuItems, &types.MenuItem{
					Label: self.c.Tr.DiscardUnstagedChanges,
					OnPress: func() error {
						self.c.LogAction(self.c.Tr.Actions.DiscardAllUnstagedChangesInFile)
						if err := self.c.Git().WorkingTree.DiscardUnstagedFileChanges(file); err != nil {
							return self.c.Error(err)
						}

						return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
					},
					Key: 'u',
					Tooltip: utils.ResolvePlaceholderString(
						self.c.Tr.DiscardUnstagedTooltip,
						map[string]string{
							"path": node.GetPath(),
						},
					),
				})
			}
		}
	}

	return self.c.Menu(types.CreateMenuOptions{Title: node.GetPath(), Items: menuItems})
}

func (self *FilesRemoveController) ResetSubmodule(submodule *models.SubmoduleConfig) error {
	return self.c.WithWaitingStatus(self.c.Tr.ResettingSubmoduleStatus, func(gocui.Task) error {
		self.c.LogAction(self.c.Tr.Actions.ResetSubmodule)

		file := self.c.Helpers().WorkingTree.FileForSubmodule(submodule)
		if file != nil {
			if err := self.c.Git().WorkingTree.UnStageFile(file.Names(), file.Tracked); err != nil {
				return self.c.Error(err)
			}
		}

		if err := self.c.Git().Submodule.Stash(submodule); err != nil {
			return self.c.Error(err)
		}
		if err := self.c.Git().Submodule.Reset(submodule); err != nil {
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
	return self.c.Contexts().Files
}
