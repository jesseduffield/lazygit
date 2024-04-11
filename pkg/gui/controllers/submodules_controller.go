package controllers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type SubmodulesController struct {
	baseController
	*ListControllerTrait[*models.SubmoduleConfig]
	c *ControllerCommon
}

var _ types.IController = &SubmodulesController{}

func NewSubmodulesController(
	c *ControllerCommon,
) *SubmodulesController {
	return &SubmodulesController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait[*models.SubmoduleConfig](
			c,
			c.Contexts().Submodules,
			c.Contexts().Submodules.GetSelected,
			c.Contexts().Submodules.GetSelectedItems,
		),
		c: c,
	}
}

func (self *SubmodulesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:               opts.GetKey(opts.Config.Universal.GoInto),
			Handler:           self.withItem(self.enter),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Enter,
			Tooltip: utils.ResolvePlaceholderString(self.c.Tr.EnterSubmoduleTooltip,
				map[string]string{"escape": keybindings.Label(opts.Config.Universal.Return)}),
			DisplayOnScreen: true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Select),
			Handler:           self.withItem(self.enter),
			GetDisabledReason: self.require(self.singleItemSelected()),
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Remove),
			Handler:           self.withItem(self.remove),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Remove,
			Tooltip:           self.c.Tr.RemoveSubmoduleTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Submodules.Update),
			Handler:           self.withItem(self.update),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Update,
			Tooltip:           self.c.Tr.SubmoduleUpdateTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:             opts.GetKey(opts.Config.Universal.New),
			Handler:         self.add,
			Description:     self.c.Tr.NewSubmodule,
			DisplayOnScreen: true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Edit),
			Handler:           self.withItem(self.editURL),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.EditSubmoduleUrl,
		},
		{
			Key:               opts.GetKey(opts.Config.Submodules.Init),
			Handler:           self.withItem(self.init),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Initialize,
			Tooltip:           self.c.Tr.InitSubmoduleTooltip,
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
	return self.withItemGraceful(self.enter)
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
					style.FgGreen.Sprint(submodule.FullName()),
					style.FgYellow.Sprint(submodule.FullPath()),
					style.FgCyan.Sprint(submodule.Url),
				)

				file := self.c.Helpers().WorkingTree.FileForSubmodule(submodule)
				if file == nil {
					task = types.NewRenderStringTask(prefix)
				} else {
					cmdObj := self.c.Git().WorkingTree.WorktreeFileDiffCmdObj(file, false, !file.HasUnstagedChanges && file.HasStagedChanges)
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
									return err
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
		Title:          fmt.Sprintf(self.c.Tr.UpdateSubmoduleUrl, submodule.FullName()),
		InitialContent: submodule.Url,
		HandleConfirm: func(newUrl string) error {
			return self.c.WithWaitingStatus(self.c.Tr.UpdatingSubmoduleUrlStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.UpdateSubmoduleUrl)
				err := self.c.Git().Submodule.UpdateUrl(submodule, newUrl)
				if err != nil {
					return err
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
			return err
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
							return err
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
							return err
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
							return err
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
			return err
		}

		return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUBMODULES}})
	})
}

func (self *SubmodulesController) remove(submodule *models.SubmoduleConfig) error {
	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.RemoveSubmodule,
		Prompt: fmt.Sprintf(self.c.Tr.RemoveSubmodulePrompt, submodule.FullName()),
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.RemoveSubmodule)
			if err := self.c.Git().Submodule.Delete(submodule); err != nil {
				return err
			}

			return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUBMODULES, types.FILES}})
		},
	})
}

func (self *SubmodulesController) easterEgg() error {
	return self.c.PushContext(self.c.Contexts().Snake)
}

func (self *SubmodulesController) context() *context.SubmodulesContext {
	return self.c.Contexts().Submodules
}
