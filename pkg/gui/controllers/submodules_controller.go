package controllers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SubmodulesController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &SubmodulesController{}

func NewSubmodulesController(
	controllerCommon *ControllerCommon,
) *SubmodulesController {
	return &SubmodulesController{
		baseController: baseController{},
		c:              controllerCommon,
	}
}

func (self *SubmodulesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.GoInto),
			Handler:     self.checkSelected(self.enter),
			Description: self.c.Tr.EnterSubmodule,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Select),
			Handler:     self.checkSelected(self.enter),
			Description: self.c.Tr.EnterSubmodule,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Remove),
			Handler:     self.checkSelected(self.remove),
			Description: self.c.Tr.RemoveSubmodule,
		},
		{
			Key:         opts.GetKey(opts.Config.Submodules.Update),
			Handler:     self.checkSelected(self.update),
			Description: self.c.Tr.SubmoduleUpdate,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.New),
			Handler:     self.add,
			Description: self.c.Tr.AddSubmodule,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Edit),
			Handler:     self.checkSelected(self.editURL),
			Description: self.c.Tr.EditSubmoduleUrl,
		},
		{
			Key:         opts.GetKey(opts.Config.Submodules.Init),
			Handler:     self.checkSelected(self.init),
			Description: self.c.Tr.InitSubmodule,
		},
		{
			Key:         opts.GetKey(opts.Config.Submodules.BulkMenu),
			Handler:     self.openBulkActionsMenu,
			Description: self.c.Tr.ViewBulkSubmoduleOptions,
			OpensMenu:   true,
		},
		{
			Key:         nil,
			Handler:     self.easterEgg,
			Description: self.c.Tr.EasterEgg,
		},
	}
}

func (self *SubmodulesController) GetOnClick() func() error {
	return self.checkSelected(self.enter)
}

func (self *SubmodulesController) GetOnRenderToMain() func() error {
	return func() error {
		return self.c.Helpers().Diff.WithDiffModeCheck(func() error {
			var task types.UpdateTask
			submodule := self.context().GetSelected()
			if submodule == nil {
				task = types.NewRenderStringTask("No submodules")
			} else {
				prefix := fmt.Sprintf(
					"Name: %s\nPath: %s\nUrl:  %s\n\n",
					style.FgGreen.Sprint(submodule.Name),
					style.FgYellow.Sprint(submodule.Path),
					style.FgCyan.Sprint(submodule.Url),
				)

				file := self.c.Helpers().WorkingTree.FileForSubmodule(submodule)
				if file == nil {
					task = types.NewRenderStringTask(prefix)
				} else {
					cmdObj := self.c.Git().WorkingTree.WorktreeFileDiffCmdObj(file, false, !file.HasUnstagedChanges && file.HasStagedChanges, self.c.GetAppState().IgnoreWhitespaceInDiffView)
					task = types.NewRunCommandTaskWithPrefix(cmdObj.GetCmd(), prefix)
				}
			}

			return self.c.RenderToMainViews(types.RefreshMainOpts{
				Pair: self.c.MainViewPairs().Normal,
				Main: &types.ViewUpdateOpts{
					Title: "Submodule",
					Task:  task,
				},
			})
		})
	}
}

func (self *SubmodulesController) enter(submodule *models.SubmoduleConfig) error {
	return self.c.Helpers().Repos.EnterSubmodule(submodule)
}

func (self *SubmodulesController) add() error {
	return self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.NewSubmoduleUrl,
		HandleConfirm: func(submoduleUrl string) error {
			nameSuggestion := filepath.Base(strings.TrimSuffix(submoduleUrl, filepath.Ext(submoduleUrl)))

			return self.c.Prompt(types.PromptOpts{
				Title:          self.c.Tr.NewSubmoduleName,
				InitialContent: nameSuggestion,
				HandleConfirm: func(submoduleName string) error {
					return self.c.Prompt(types.PromptOpts{
						Title:          self.c.Tr.NewSubmodulePath,
						InitialContent: submoduleName,
						HandleConfirm: func(submodulePath string) error {
							return self.c.WithWaitingStatus(self.c.Tr.AddingSubmoduleStatus, func(gocui.Task) error {
								self.c.LogAction(self.c.Tr.Actions.AddSubmodule)
								err := self.c.Git().Submodule.Add(submoduleName, submodulePath, submoduleUrl)
								if err != nil {
									_ = self.c.Error(err)
								}

								return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUBMODULES}})
							})
						},
					})
				},
			})
		},
	})
}

func (self *SubmodulesController) editURL(submodule *models.SubmoduleConfig) error {
	return self.c.Prompt(types.PromptOpts{
		Title:          fmt.Sprintf(self.c.Tr.UpdateSubmoduleUrl, submodule.Name),
		InitialContent: submodule.Url,
		HandleConfirm: func(newUrl string) error {
			return self.c.WithWaitingStatus(self.c.Tr.UpdatingSubmoduleUrlStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.UpdateSubmoduleUrl)
				err := self.c.Git().Submodule.UpdateUrl(submodule.Name, submodule.Path, newUrl)
				if err != nil {
					_ = self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUBMODULES}})
			})
		},
	})
}

func (self *SubmodulesController) init(submodule *models.SubmoduleConfig) error {
	return self.c.WithWaitingStatus(self.c.Tr.InitializingSubmoduleStatus, func(gocui.Task) error {
		self.c.LogAction(self.c.Tr.Actions.InitialiseSubmodule)
		err := self.c.Git().Submodule.Init(submodule.Path)
		if err != nil {
			_ = self.c.Error(err)
		}

		return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUBMODULES}})
	})
}

func (self *SubmodulesController) openBulkActionsMenu() error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.BulkSubmoduleOptions,
		Items: []*types.MenuItem{
			{
				LabelColumns: []string{self.c.Tr.BulkInitSubmodules, style.FgGreen.Sprint(self.c.Git().Submodule.BulkInitCmdObj().ToString())},
				OnPress: func() error {
					return self.c.WithWaitingStatus(self.c.Tr.RunningCommand, func(gocui.Task) error {
						self.c.LogAction(self.c.Tr.Actions.BulkInitialiseSubmodules)
						err := self.c.Git().Submodule.BulkInitCmdObj().Run()
						if err != nil {
							return self.c.Error(err)
						}

						return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUBMODULES}})
					})
				},
				Key: 'i',
			},
			{
				LabelColumns: []string{self.c.Tr.BulkUpdateSubmodules, style.FgYellow.Sprint(self.c.Git().Submodule.BulkUpdateCmdObj().ToString())},
				OnPress: func() error {
					return self.c.WithWaitingStatus(self.c.Tr.RunningCommand, func(gocui.Task) error {
						self.c.LogAction(self.c.Tr.Actions.BulkUpdateSubmodules)
						if err := self.c.Git().Submodule.BulkUpdateCmdObj().Run(); err != nil {
							return self.c.Error(err)
						}

						return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUBMODULES}})
					})
				},
				Key: 'u',
			},
			{
				LabelColumns: []string{self.c.Tr.BulkDeinitSubmodules, style.FgRed.Sprint(self.c.Git().Submodule.BulkDeinitCmdObj().ToString())},
				OnPress: func() error {
					return self.c.WithWaitingStatus(self.c.Tr.RunningCommand, func(gocui.Task) error {
						self.c.LogAction(self.c.Tr.Actions.BulkDeinitialiseSubmodules)
						if err := self.c.Git().Submodule.BulkDeinitCmdObj().Run(); err != nil {
							return self.c.Error(err)
						}

						return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUBMODULES}})
					})
				},
				Key: 'd',
			},
		},
	})
}

func (self *SubmodulesController) update(submodule *models.SubmoduleConfig) error {
	return self.c.WithWaitingStatus(self.c.Tr.UpdatingSubmoduleStatus, func(gocui.Task) error {
		self.c.LogAction(self.c.Tr.Actions.UpdateSubmodule)
		err := self.c.Git().Submodule.Update(submodule.Path)
		if err != nil {
			_ = self.c.Error(err)
		}

		return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUBMODULES}})
	})
}

func (self *SubmodulesController) remove(submodule *models.SubmoduleConfig) error {
	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.RemoveSubmodule,
		Prompt: fmt.Sprintf(self.c.Tr.RemoveSubmodulePrompt, submodule.Name),
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.RemoveSubmodule)
			if err := self.c.Git().Submodule.Delete(submodule); err != nil {
				return self.c.Error(err)
			}

			return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUBMODULES, types.FILES}})
		},
	})
}

func (self *SubmodulesController) easterEgg() error {
	return self.c.PushContext(self.c.Contexts().Snake)
}

func (self *SubmodulesController) checkSelected(callback func(*models.SubmoduleConfig) error) func() error {
	return func() error {
		submodule := self.context().GetSelected()
		if submodule == nil {
			return nil
		}

		return callback(submodule)
	}
}

func (self *SubmodulesController) Context() types.Context {
	return self.context()
}

func (self *SubmodulesController) context() *context.SubmodulesContext {
	return self.c.Contexts().Submodules
}
