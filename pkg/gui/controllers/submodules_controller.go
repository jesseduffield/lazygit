package controllers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SubmodulesController struct {
	c       *types.ControllerCommon
	context types.IListContext
	git     *commands.GitCommand

	enterSubmodule       func(submodule *models.SubmoduleConfig) error
	getSelectedSubmodule func() *models.SubmoduleConfig
}

var _ types.IController = &SubmodulesController{}

func NewSubmodulesController(
	c *types.ControllerCommon,
	context types.IListContext,
	git *commands.GitCommand,
	enterSubmodule func(submodule *models.SubmoduleConfig) error,
	getSelectedSubmodule func() *models.SubmoduleConfig,
) *SubmodulesController {
	return &SubmodulesController{
		c:                    c,
		context:              context,
		git:                  git,
		enterSubmodule:       enterSubmodule,
		getSelectedSubmodule: getSelectedSubmodule,
	}
}

func (self *SubmodulesController) Keybindings(getKey func(key string) interface{}, config config.KeybindingConfig, guards types.KeybindingGuards) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         getKey(config.Universal.GoInto),
			Handler:     self.checkSelected(self.enter),
			Description: self.c.Tr.LcEnterSubmodule,
		},
		{
			Key:         getKey(config.Universal.Remove),
			Handler:     self.checkSelected(self.remove),
			Description: self.c.Tr.LcRemoveSubmodule,
		},
		{
			Key:         getKey(config.Submodules.Update),
			Handler:     self.checkSelected(self.update),
			Description: self.c.Tr.LcSubmoduleUpdate,
		},
		{
			Key:         getKey(config.Universal.New),
			Handler:     self.add,
			Description: self.c.Tr.LcAddSubmodule,
		},
		{
			Key:         getKey(config.Universal.Edit),
			Handler:     self.checkSelected(self.editURL),
			Description: self.c.Tr.LcEditSubmoduleUrl,
		},
		{
			Key:         getKey(config.Submodules.Init),
			Handler:     self.checkSelected(self.init),
			Description: self.c.Tr.LcInitSubmodule,
		},
		{
			Key:         getKey(config.Submodules.BulkMenu),
			Handler:     self.openBulkActionsMenu,
			Description: self.c.Tr.LcViewBulkSubmoduleOptions,
			OpensMenu:   true,
		},
		{
			Key:     gocui.MouseLeft,
			Handler: func() error { return self.context.HandleClick(self.checkSelected(self.enter)) },
		},
	}

	return append(bindings, self.context.Keybindings(getKey, config, guards)...)
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
				DisplayStrings: []string{self.c.Tr.LcBulkInitSubmodules, style.FgGreen.Sprint(self.git.Submodule.BulkInitCmdObj().ToString())},
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
			},
			{
				DisplayStrings: []string{self.c.Tr.LcBulkUpdateSubmodules, style.FgYellow.Sprint(self.git.Submodule.BulkUpdateCmdObj().ToString())},
				OnPress: func() error {
					return self.c.WithWaitingStatus(self.c.Tr.LcRunningCommand, func() error {
						self.c.LogAction(self.c.Tr.Actions.BulkUpdateSubmodules)
						if err := self.git.Submodule.BulkUpdateCmdObj().Run(); err != nil {
							return self.c.Error(err)
						}

						return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUBMODULES}})
					})
				},
			},
			{
				DisplayStrings: []string{self.c.Tr.LcBulkDeinitSubmodules, style.FgRed.Sprint(self.git.Submodule.BulkDeinitCmdObj().ToString())},
				OnPress: func() error {
					return self.c.WithWaitingStatus(self.c.Tr.LcRunningCommand, func() error {
						self.c.LogAction(self.c.Tr.Actions.BulkDeinitialiseSubmodules)
						if err := self.git.Submodule.BulkDeinitCmdObj().Run(); err != nil {
							return self.c.Error(err)
						}

						return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUBMODULES}})
					})
				},
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
	return self.c.Ask(types.AskOpts{
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
		submodule := self.getSelectedSubmodule()
		if submodule == nil {
			return nil
		}

		return callback(submodule)
	}
}

func (self *SubmodulesController) Context() types.Context {
	return self.context
}
