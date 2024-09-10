package controllers

import (
	"errors"
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type GitFlowController struct {
	baseController
	*ListControllerTrait[*models.Branch]
	c *ControllerCommon
}

var _ types.IController = &GitFlowController{}

func NewGitFlowController(
	c *ControllerCommon,
) *GitFlowController {
	return &GitFlowController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait[*models.Branch](
			c,
			c.Contexts().Branches,
			c.Contexts().Branches.GetSelected,
			c.Contexts().Branches.GetSelectedItems,
		),
		c: c,
	}
}

func (self *GitFlowController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Branches.ViewGitFlowOptions),
			Handler:     self.withItem(self.handleCreateGitFlowMenu),
			Description: self.c.Tr.GitFlowOptions,
			OpensMenu:   true,
		},
	}

	return bindings
}

func (self *GitFlowController) handleCreateGitFlowMenu(branch *models.Branch) error {
	if !self.c.Git().Flow.GitFlowEnabled() {
		return errors.New("You need to install git-flow and enable it in this repo to use git-flow features")
	}

	startHandler := func(branchType string) func() error {
		return func() error {
			title := utils.ResolvePlaceholderString(self.c.Tr.NewGitFlowBranchPrompt, map[string]string{"branchType": branchType})

			self.c.Prompt(types.PromptOpts{
				Title: title,
				HandleConfirm: func(name string) error {
					self.c.LogAction(self.c.Tr.Actions.GitFlowStart)
					return self.c.RunSubprocessAndRefresh(
						self.c.Git().Flow.StartCmdObj(branchType, name),
					)
				},
			})

			return nil
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
				DisabledReason: self.require(self.singleItemSelected())(),
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
	cmdObj, err := self.c.Git().Flow.FinishCmdObj(branchName)
	if err != nil {
		return err
	}

	self.c.LogAction(self.c.Tr.Actions.GitFlowFinish)
	return self.c.RunSubprocessAndRefresh(cmdObj)
}
