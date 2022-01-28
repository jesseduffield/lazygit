package controllers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// if Go let me do private struct embedding of structs with public fields (which it should)
// I would just do that. But alas.
type ControllerCommon struct {
	*common.Common
	IGuiCommon
}

type SubmodulesController struct {
	// I've said publicly that I'm against single-letter variable names but in this
	// case I would actually prefer a _zero_ letter variable name in the form of
	// struct embedding, but Go does not allow hiding public fields in an embedded struct
	// to the client
	c                    *ControllerCommon
	enterSubmoduleFn     func(submodule *models.SubmoduleConfig) error
	getSelectedSubmodule func() *models.SubmoduleConfig
	git                  *commands.GitCommand
	submodules           []*models.SubmoduleConfig
}

func NewSubmodulesController(
	c *ControllerCommon,
	enterSubmoduleFn func(submodule *models.SubmoduleConfig) error,
	git *commands.GitCommand,
	submodules []*models.SubmoduleConfig,
	getSelectedSubmodule func() *models.SubmoduleConfig,
) *SubmodulesController {
	return &SubmodulesController{
		c:                    c,
		enterSubmoduleFn:     enterSubmoduleFn,
		git:                  git,
		submodules:           submodules,
		getSelectedSubmodule: getSelectedSubmodule,
	}
}

func (self *SubmodulesController) Keybindings(getKey func(key string) interface{}, config config.KeybindingConfig) []*types.Binding {
	return []*types.Binding{
		{
			Key:         getKey(config.Universal.GoInto),
			Handler:     self.forSubmodule(self.enter),
			Description: self.c.Tr.LcEnterSubmodule,
		},
		{
			Key:         getKey(config.Universal.Remove),
			Handler:     self.forSubmodule(self.remove),
			Description: self.c.Tr.LcRemoveSubmodule,
		},
		{
			Key:         getKey(config.Submodules.Update),
			Handler:     self.forSubmodule(self.update),
			Description: self.c.Tr.LcSubmoduleUpdate,
		},
		{
			Key:         getKey(config.Universal.New),
			Handler:     self.add,
			Description: self.c.Tr.LcAddSubmodule,
		},
		{
			Key:         getKey(config.Universal.Edit),
			Handler:     self.forSubmodule(self.editURL),
			Description: self.c.Tr.LcEditSubmoduleUrl,
		},
		{
			Key:         getKey(config.Submodules.Init),
			Handler:     self.forSubmodule(self.init),
			Description: self.c.Tr.LcInitSubmodule,
		},
		{
			Key:         getKey(config.Submodules.BulkMenu),
			Handler:     self.openBulkActionsMenu,
			Description: self.c.Tr.LcViewBulkSubmoduleOptions,
			OpensMenu:   true,
		},
	}
}

func (self *SubmodulesController) enter(submodule *models.SubmoduleConfig) error {
	return self.enterSubmoduleFn(submodule)
}

func (self *SubmodulesController) add() error {
	return self.c.Prompt(popup.PromptOpts{
		Title: self.c.Tr.LcNewSubmoduleUrl,
		HandleConfirm: func(submoduleUrl string) error {
			nameSuggestion := filepath.Base(strings.TrimSuffix(submoduleUrl, filepath.Ext(submoduleUrl)))

			return self.c.Prompt(popup.PromptOpts{
				Title:          self.c.Tr.LcNewSubmoduleName,
				InitialContent: nameSuggestion,
				HandleConfirm: func(submoduleName string) error {

					return self.c.Prompt(popup.PromptOpts{
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
	return self.c.Prompt(popup.PromptOpts{
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
	return self.c.Menu(popup.CreateMenuOptions{
		Title: self.c.Tr.LcBulkSubmoduleOptions,
		Items: []*popup.MenuItem{
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
	return self.c.Ask(popup.AskOpts{
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

func (self *SubmodulesController) forSubmodule(callback func(*models.SubmoduleConfig) error) func() error {
	return func() error {
		submodule := self.getSelectedSubmodule()
		if submodule == nil {
			return nil
		}

		return callback(submodule)
	}
}
