package controllers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// This controller is for all contexts that contain a list of commits.

var _ types.IController = &BasicCommitsController{}

type ContainsCommits interface {
	types.Context
	GetSelected() *models.Commit
	GetCommits() []*models.Commit
	GetSelectedLineIdx() int
}

type BasicCommitsController struct {
	baseController
	*controllerCommon
	context ContainsCommits
}

func NewBasicCommitsController(controllerCommon *controllerCommon, context ContainsCommits) *BasicCommitsController {
	return &BasicCommitsController{
		baseController:   baseController{},
		controllerCommon: controllerCommon,
		context:          context,
	}
}

func (self *BasicCommitsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Commits.CheckoutCommit),
			Handler:     self.checkSelected(self.checkout),
			Description: self.c.Tr.LcCheckoutCommit,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.CopyCommitAttributeToClipboard),
			Handler:     self.checkSelected(self.copyCommitAttribute),
			Description: self.c.Tr.LcCopyCommitAttributeToClipboard,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.OpenInBrowser),
			Handler:     self.checkSelected(self.openInBrowser),
			Description: self.c.Tr.LcOpenCommitInBrowser,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.New),
			Handler:     self.checkSelected(self.newBranch),
			Description: self.c.Tr.LcCreateNewBranchFromCommit,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.ViewResetOptions),
			Handler:     self.checkSelected(self.createResetMenu),
			Description: self.c.Tr.LcViewResetOptions,
			OpensMenu:   true,
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

func (self *BasicCommitsController) checkSelected(callback func(*models.Commit) error) func() error {
	return func() error {
		commit := self.context.GetSelected()
		if commit == nil {
			return nil
		}

		return callback(commit)
	}
}

func (self *BasicCommitsController) Context() types.Context {
	return self.context
}

func (self *BasicCommitsController) copyCommitAttribute(commit *models.Commit) error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.Actions.CopyCommitAttributeToClipboard,
		Items: []*types.MenuItem{
			{
				Label: self.c.Tr.LcCommitSha,
				OnPress: func() error {
					return self.copyCommitSHAToClipboard(commit)
				},
				Key: 's',
			},
			{
				Label: self.c.Tr.LcCommitURL,
				OnPress: func() error {
					return self.copyCommitURLToClipboard(commit)
				},
				Key: 'u',
			},
			{
				Label: self.c.Tr.LcCommitDiff,
				OnPress: func() error {
					return self.copyCommitDiffToClipboard(commit)
				},
				Key: 'd',
			},
			{
				Label: self.c.Tr.LcCommitMessage,
				OnPress: func() error {
					return self.copyCommitMessageToClipboard(commit)
				},
				Key: 'm',
			},
			{
				Label: self.c.Tr.LcCommitAuthor,
				OnPress: func() error {
					return self.copyAuthorToClipboard(commit)
				},
				Key: 'a',
			},
		},
	})
}

func (self *BasicCommitsController) copyCommitSHAToClipboard(commit *models.Commit) error {
	self.c.LogAction(self.c.Tr.Actions.CopyCommitSHAToClipboard)
	if err := self.os.CopyToClipboard(commit.Sha); err != nil {
		return self.c.Error(err)
	}

	self.c.Toast(self.c.Tr.CommitSHACopiedToClipboard)
	return nil
}

func (self *BasicCommitsController) copyCommitURLToClipboard(commit *models.Commit) error {
	url, err := self.helpers.Host.GetCommitURL(commit.Sha)
	if err != nil {
		return err
	}

	self.c.LogAction(self.c.Tr.Actions.CopyCommitURLToClipboard)
	if err := self.os.CopyToClipboard(url); err != nil {
		return self.c.Error(err)
	}

	self.c.Toast(self.c.Tr.CommitURLCopiedToClipboard)
	return nil
}

func (self *BasicCommitsController) copyCommitDiffToClipboard(commit *models.Commit) error {
	diff, err := self.git.Commit.GetCommitDiff(commit.Sha)
	if err != nil {
		return self.c.Error(err)
	}

	self.c.LogAction(self.c.Tr.Actions.CopyCommitDiffToClipboard)
	if err := self.os.CopyToClipboard(diff); err != nil {
		return self.c.Error(err)
	}

	self.c.Toast(self.c.Tr.CommitDiffCopiedToClipboard)
	return nil
}

func (self *BasicCommitsController) copyAuthorToClipboard(commit *models.Commit) error {
	author, err := self.git.Commit.GetCommitAuthor(commit.Sha)
	if err != nil {
		return self.c.Error(err)
	}

	formattedAuthor := fmt.Sprintf("%s <%s>", author.Name, author.Email)

	self.c.LogAction(self.c.Tr.Actions.CopyCommitAuthorToClipboard)
	if err := self.os.CopyToClipboard(formattedAuthor); err != nil {
		return self.c.Error(err)
	}

	self.c.Toast(self.c.Tr.CommitAuthorCopiedToClipboard)
	return nil
}

func (self *BasicCommitsController) copyCommitMessageToClipboard(commit *models.Commit) error {
	message, err := self.git.Commit.GetCommitMessage(commit.Sha)
	if err != nil {
		return self.c.Error(err)
	}

	self.c.LogAction(self.c.Tr.Actions.CopyCommitMessageToClipboard)
	if err := self.os.CopyToClipboard(message); err != nil {
		return self.c.Error(err)
	}

	self.c.Toast(self.c.Tr.CommitMessageCopiedToClipboard)
	return nil
}

func (self *BasicCommitsController) openInBrowser(commit *models.Commit) error {
	url, err := self.helpers.Host.GetCommitURL(commit.Sha)
	if err != nil {
		return self.c.Error(err)
	}

	self.c.LogAction(self.c.Tr.Actions.OpenCommitInBrowser)
	if err := self.os.OpenLink(url); err != nil {
		return self.c.Error(err)
	}

	return nil
}

func (self *BasicCommitsController) newBranch(commit *models.Commit) error {
	return self.helpers.Refs.NewBranch(commit.RefName(), commit.Description(), "")
}

func (self *BasicCommitsController) createResetMenu(commit *models.Commit) error {
	return self.helpers.Refs.CreateGitResetMenu(commit.Sha)
}

func (self *BasicCommitsController) checkout(commit *models.Commit) error {
	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.LcCheckoutCommit,
		Prompt: self.c.Tr.SureCheckoutThisCommit,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.CheckoutCommit)
			return self.helpers.Refs.CheckoutRef(commit.Sha, types.CheckoutRefOptions{})
		},
	})
}

func (self *BasicCommitsController) copy(commit *models.Commit) error {
	return self.helpers.CherryPick.Copy(commit, self.context.GetCommits(), self.context)
}

func (self *BasicCommitsController) copyRange(*models.Commit) error {
	return self.helpers.CherryPick.CopyRange(self.context.GetSelectedLineIdx(), self.context.GetCommits(), self.context)
}
