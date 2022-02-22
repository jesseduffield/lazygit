package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ReflogController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &ReflogController{}

func NewReflogController(
	common *controllerCommon,
) *ReflogController {
	return &ReflogController{
		baseController:   baseController{},
		controllerCommon: common,
	}
}

func (self *ReflogController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
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
			Key:         opts.GetKey(opts.Config.Commits.CherryPickCopy),
			Handler:     opts.Guards.OutsideFilterMode(self.checkSelected(self.copy)),
			Description: self.c.Tr.LcCherryPickCopy,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.CherryPickCopyRange),
			Handler:     opts.Guards.OutsideFilterMode(self.checkSelected(self.copyRange)),
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

func (self *ReflogController) checkSelected(callback func(*models.Commit) error) func() error {
	return func() error {
		commit := self.context().GetSelected()
		if commit == nil {
			return nil
		}

		return callback(commit)
	}
}

func (self *ReflogController) Context() types.Context {
	return self.context()
}

func (self *ReflogController) context() *context.ReflogCommitsContext {
	return self.contexts.ReflogCommits
}

func (self *ReflogController) checkout(commit *models.Commit) error {
	err := self.c.Ask(types.AskOpts{
		Title:  self.c.Tr.LcCheckoutCommit,
		Prompt: self.c.Tr.SureCheckoutThisCommit,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.CheckoutReflogCommit)
			return self.helpers.Refs.CheckoutRef(commit.Sha, types.CheckoutRefOptions{})
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (self *ReflogController) openResetMenu(commit *models.Commit) error {
	return self.helpers.Refs.CreateGitResetMenu(commit.Sha)
}

func (self *ReflogController) copy(commit *models.Commit) error {
	return self.helpers.CherryPick.Copy(commit, self.model.FilteredReflogCommits, self.context())
}

func (self *ReflogController) copyRange(commit *models.Commit) error {
	return self.helpers.CherryPick.CopyRange(self.context().GetSelectedLineIdx(), self.model.FilteredReflogCommits, self.context())
}
