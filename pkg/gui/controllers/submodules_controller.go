package controllers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SubmodulesController struct {
	baseController
	*controllerCommon

	enterSubmodule func(submodule *models.SubmoduleConfig) error
}

var _ types.IController = &SubmodulesController{}

func NewSubmodulesController(
	controllerCommon *controllerCommon,
	enterSubmodule func(submodule *models.SubmoduleConfig) error,
) *SubmodulesController {
	return &SubmodulesController{
		baseController:   baseController{},
		controllerCommon: controllerCommon,
		enterSubmodule:   enterSubmodule,
	}
}

func (self *SubmodulesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.GoInto),
			Handler:     self.checkSelected(self.enter),
			Description: self.c.Tr.LcEnterSubmodule,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Remove),
			Handler:     self.checkSelected(self.remove),
			Description: self.c.Tr.LcRemoveSubmodule,
		},
		{
			Key:         opts.GetKey(opts.Config.Submodules.Update),
			Handler:     self.checkSelected(self.update),
			Description: self.c.Tr.LcSubmoduleUpdate,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.New),
			Handler:     self.add,
			Description: self.c.Tr.LcAddSubmodule,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Edit),
			Handler:     self.checkSelected(self.editURL),
			Description: self.c.Tr.LcEditSubmoduleUrl,
		},
		{
			Key:         opts.GetKey(opts.Config.Submodules.Init),
			Handler:     self.checkSelected(self.init),
			Description: self.c.Tr.LcInitSubmodule,
		},
		{
			Key:         opts.GetKey(opts.Config.Submodules.BulkMenu),
			Handler:     self.openBulkActionsMenu,
			Description: self.c.Tr.LcViewBulkSubmoduleOptions,
			OpensMenu:   true,
		},
	}
}

func (self *SubmodulesController) GetOnClick() func() error {
	return self.checkSelected(self.enter)
}

func (self *SubmodulesController) enter(submodule *models.SubmoduleConfig) error {
	return self.enterSubmodule(submodule)
}

func (self *SubmodulesController) add() error {
	return self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.LcNewSubmoduleUrl,
		HandleConfirm: func(submoduleUrl string) error {
			nameSuggestion := filepath.Base(strings.TrimSuffix(submoduleUrl, filepath.Ext(submoduleUrl)))

			return self.c.Prompt(types.PromptOpts{
				Title:          self.c.Tr.LcNewSubmoduleName,
				InitialContent: nameSuggestion,
				HandleConfirm: func(submoduleName string) error {
					return self.c.Prompt(types.PromptOpts{
						Title:          self.c.Tr.LcNewSubmodulePath,
						InitialContent: submoduleName,
						HandleConfirm: func(submodulePath string) error {
							return self.c.WithWaitingStatus(self.c.Tr.LcAddingSubmoduleStatus, func() error {
								self.c.LogAction(self.c.Tr.Actions.AddSubmodule)
								err := self.git.Submodule.Add(submoduleName, submodulePath, submoduleUrl)
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
		Title:          fmt.Sprintf(self.c.Tr.LcUpdateSubmoduleUrl, submodule.Name),
		InitialContent: submodule.Url,
		HandleConfirm: func(newUrl string) error {
			return self.c.WithWaitingStatus(self.c.Tr.LcUpdatingSubmoduleUrlStatus, func() error {
				self.c.LogAction(self.c.Tr.Actions.UpdateSubmoduleUrl)
				err := self.git.Submodule.UpdateUrl(submodule.Name, submodule.Path, newUrl)
				if err != nil {
					_ = self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUBMODULES}})
			})
		},
	})
}

func (self *SubmodulesController) init(submodule *models.SubmoduleConfig) error {
	return self.c.WithWaitingStatus(self.c.Tr.LcInitializingSubmoduleStatus, func() error {
		self.c.LogAction(self.c.Tr.Actions.InitialiseSubmodule)
		err := self.git.Submodule.Init(submodule.Path)
		if err != nil {
			_ = self.c.Error(err)
		}

		return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUBMODULES}})
	})
}

func (self *SubmodulesController) openBulkActionsMenu() error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.LcBulkSubmoduleOptions,
		Items: []*types.MenuItem{
			{
				LabelColumns: []string{self.c.Tr.LcBulkInitSubmodules, style.FgGreen.Sprint(self.git.Submodule.BulkInitCmdObj().ToString())},
				OnPress: func() error {
					return self.c.WithWaitingStatus(self.c.Tr.LcRunningCommand, func() error {
						self.c.LogAction(self.c.Tr.Actions.BulkInitialiseSubmodules)
						err := self.git.Submodule.BulkInitCmdObj().Run()
						if err != nil {
							return self.c.Error(err)
						}

						return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUBMODULES}})
					})
				},
				Key: 'i',
			},
			{
				LabelColumns: []string{self.c.Tr.LcBulkUpdateSubmodules, style.FgYellow.Sprint(self.git.Submodule.BulkUpdateCmdObj().ToString())},
				OnPress: func() error {
					return self.c.WithWaitingStatus(self.c.Tr.LcRunningCommand, func() error {
						self.c.LogAction(self.c.Tr.Actions.BulkUpdateSubmodules)
						if err := self.git.Submodule.BulkUpdateCmdObj().Run(); err != nil {
							return self.c.Error(err)
						}

						return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUBMODULES}})
					})
				},
				Key: 'u',
			},
			{
				LabelColumns: []string{self.c.Tr.LcBulkDeinitSubmodules, style.FgRed.Sprint(self.git.Submodule.BulkDeinitCmdObj().ToString())},
				OnPress: func() error {
					return self.c.WithWaitingStatus(self.c.Tr.LcRunningCommand, func() error {
						self.c.LogAction(self.c.Tr.Actions.BulkDeinitialiseSubmodules)
						if err := self.git.Submodule.BulkDeinitCmdObj().Run(); err != nil {
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
	return self.c.WithWaitingStatus(self.c.Tr.LcUpdatingSubmoduleStatus, func() error {
		self.c.LogAction(self.c.Tr.Actions.UpdateSubmodule)
		err := self.git.Submodule.Update(submodule.Path)
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
			if err := self.git.Submodule.Delete(submodule); err != nil {
				return self.c.Error(err)
			}

			return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUBMODULES, types.FILES}})
		},
	})
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
	return self.contexts.Submodules
}
