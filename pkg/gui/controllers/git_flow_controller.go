package controllers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type GitFlowController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &GitFlowController{}

func NewGitFlowController(
	common *controllerCommon,
) *GitFlowController {
	return &GitFlowController{
		baseController:   baseController{},
		controllerCommon: common,
	}
}

func (self *GitFlowController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Branches.ViewGitFlowOptions),
			Handler:     self.checkSelected(self.handleCreateGitFlowMenu),
			Description: self.c.Tr.LcGitFlowOptions,
			OpensMenu:   true,
		},
	}

	return bindings
}

func (self *GitFlowController) handleCreateGitFlowMenu(branch *models.Branch) error {
	if !self.git.Flow.GitFlowEnabled() {
		return self.c.ErrorMsg("You need to install git-flow and enable it in this repo to use git-flow features")
	}

	startHandler := func(branchType string) func() error {
		return func() error {
			title := utils.ResolvePlaceholderString(self.c.Tr.NewGitFlowBranchPrompt, map[string]string{"branchType": branchType})

			return self.c.Prompt(types.PromptOpts{
				Title: title,
				HandleConfirm: func(name string) error {
					self.c.LogAction(self.c.Tr.Actions.GitFlowStart)
					return self.c.RunSubprocessAndRefresh(
						self.git.Flow.StartCmdObj(branchType, name),
					)
				},
			})
		}
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: "git flow",
		Items: []*types.MenuItem{
			{
				// not localising here because it's one to one with the actual git flow commands
				Label: fmt.Sprintf("finish branch '%s'", branch.Name),
				OnPress: func() error {
					return self.gitFlowFinishBranch(branch.Name)
				},
			},
			{
				Label:   "start feature",
				OnPress: startHandler("feature"),
				Key:     'f',
			},
			{
				Label:   "start hotfix",
				OnPress: startHandler("hotfix"),
				Key:     'h',
			},
			{
				Label:   "start bugfix",
				OnPress: startHandler("bugfix"),
				Key:     'b',
			},
			{
				Label:   "start release",
				OnPress: startHandler("release"),
				Key:     'r',
			},
		},
	})
}

func (self *GitFlowController) gitFlowFinishBranch(branchName string) error {
	cmdObj, err := self.git.Flow.FinishCmdObj(branchName)
	if err != nil {
		return self.c.Error(err)
	}

	self.c.LogAction(self.c.Tr.Actions.GitFlowFinish)
	return self.c.RunSubprocessAndRefresh(cmdObj)
}

func (self *GitFlowController) checkSelected(callback func(*models.Branch) error) func() error {
	return func() error {
		node := self.context().GetSelected()
		if node == nil {
			return nil
		}

		return callback(node)
	}
}

func (self *GitFlowController) Context() types.Context {
	return self.context()
}

func (self *GitFlowController) context() *context.BranchesContext {
	return self.contexts.Branches
}
