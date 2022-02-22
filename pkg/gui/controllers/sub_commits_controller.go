package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SubCommitsController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &SubCommitsController{}

func NewSubCommitsController(
	common *controllerCommon,
) *SubCommitsController {
	return &SubCommitsController{
		baseController:   baseController{},
		controllerCommon: common,
	}
}

func (self *SubCommitsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.Select),
			Handler:     self.checkSelected(self.checkout),
			Description: self.c.Tr.LcCheckoutCommit,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.ViewResetOptions),
			Handler:     self.checkSelected(self.openResetMenu),
			Description: self.c.Tr.LcViewResetOptions,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.New),
			Handler:     self.checkSelected(self.newBranch),
			Description: self.c.Tr.LcNewBranch,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.CherryPickCopy),
			Handler:     self.checkSelected(self.copy),
			Description: self.c.Tr.LcCherryPickCopy,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.CherryPickCopyRange),
			Handler:     self.checkSelected(self.copyRange),
			Description: self.c.Tr.LcCherryPickCopyRange,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.ResetCherryPick),
			Handler:     self.helpers.CherryPick.Reset,
			Description: self.c.Tr.LcResetCherryPick,
		},
	}

	return bindings
}

func (self *SubCommitsController) checkSelected(callback func(*models.Commit) error) func() error {
	return func() error {
		commit := self.context().GetSelected()
		if commit == nil {
			return nil
		}

		return callback(commit)
	}
}

func (self *SubCommitsController) Context() types.Context {
	return self.context()
}

func (self *SubCommitsController) context() *context.ReflogCommitsContext {
	return self.contexts.ReflogCommits
}

func (self *SubCommitsController) checkout(commit *models.Commit) error {
	err := self.c.Ask(types.AskOpts{
		Title:  self.c.Tr.LcCheckoutCommit,
		Prompt: self.c.Tr.SureCheckoutThisCommit,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.CheckoutCommit)
			return self.helpers.Refs.CheckoutRef(commit.Sha, types.CheckoutRefOptions{})
		},
	})
	if err != nil {
		return err
	}

	self.context().SetSelectedLineIdx(0)

	return nil
}

func (self *SubCommitsController) openResetMenu(commit *models.Commit) error {
	return self.helpers.Refs.CreateGitResetMenu(commit.Sha)
}

func (self *SubCommitsController) newBranch(commit *models.Commit) error {
	return self.helpers.Refs.NewBranch(commit.RefName(), commit.Description(), "")
}

func (self *SubCommitsController) copy(commit *models.Commit) error {
	return self.helpers.CherryPick.Copy(commit, self.model.SubCommits, self.context())
}

func (self *SubCommitsController) copyRange(commit *models.Commit) error {
	return self.helpers.CherryPick.CopyRange(self.context().GetSelectedLineIdx(), self.model.SubCommits, self.context())
}
