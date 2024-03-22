package controllers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// This controller is for all contexts that contain a list of commits.

var _ types.IController = &BasicCommitsController{}

type ContainsCommits interface {
	types.Context
	types.IListContext
	GetSelected() *models.Commit
	GetSelectedItems() ([]*models.Commit, int, int)
	GetCommits() []*models.Commit
	GetSelectedLineIdx() int
}

type BasicCommitsController struct {
	baseController
	*ListControllerTrait[*models.Commit]
	c       *ControllerCommon
	context ContainsCommits
}

func NewBasicCommitsController(c *ControllerCommon, context ContainsCommits) *BasicCommitsController {
	return &BasicCommitsController{
		baseController: baseController{},
		c:              c,
		context:        context,
		ListControllerTrait: NewListControllerTrait[*models.Commit](
			c,
			context,
			context.GetSelected,
			context.GetSelectedItems,
		),
	}
}

func (self *BasicCommitsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:               opts.GetKey(opts.Config.Commits.CheckoutCommit),
			Handler:           self.withItem(self.checkout),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Checkout,
			Tooltip:           self.c.Tr.CheckoutCommitTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.CopyCommitAttributeToClipboard),
			Handler:           self.withItem(self.copyCommitAttribute),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.CopyCommitAttributeToClipboard,
			Tooltip:           self.c.Tr.CopyCommitAttributeToClipboardTooltip,
			OpensMenu:         true,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.OpenInBrowser),
			Handler:           self.withItem(self.openInBrowser),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.OpenCommitInBrowser,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.New),
			Handler:           self.withItem(self.newBranch),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.CreateNewBranchFromCommit,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.ViewResetOptions),
			Handler:           self.withItem(self.createResetMenu),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.ViewResetOptions,
			Tooltip:           self.c.Tr.ResetTooltip,
			OpensMenu:         true,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.CherryPickCopy),
			Handler:           self.withItem(self.copyRange),
			GetDisabledReason: self.require(self.itemRangeSelected(self.canCopyCommits)),
			Description:       self.c.Tr.CherryPickCopy,
			Tooltip: utils.ResolvePlaceholderString(self.c.Tr.CherryPickCopyTooltip,
				map[string]string{
					"paste":  keybindings.Label(opts.Config.Commits.PasteCommits),
					"escape": keybindings.Label(opts.Config.Universal.Return),
				},
			),
			DisplayOnScreen: true,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.ResetCherryPick),
			Handler:     self.c.Helpers().CherryPick.Reset,
			Description: self.c.Tr.ResetCherryPick,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.OpenDiffTool),
			Handler:           self.withItem(self.openDiffTool),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.OpenDiffTool,
		},
		// Putting this at the bottom of the list so that it has the lowest priority,
		// meaning that if the user has configured another keybinding to the same key
		// then that will take precedence.
		{
			// Hardcoding this key because it's not configurable
			Key:     opts.GetKey("c"),
			Handler: self.handleOldCherryPickKey,
		},
	}

	return bindings
}

func (self *BasicCommitsController) copyCommitAttribute(commit *models.Commit) error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.Actions.CopyCommitAttributeToClipboard,
		Items: []*types.MenuItem{
			{
				Label: self.c.Tr.CommitSha,
				OnPress: func() error {
					return self.copyCommitSHAToClipboard(commit)
				},
			},
			{
				Label: self.c.Tr.CommitSubject,
				OnPress: func() error {
					return self.copyCommitSubjectToClipboard(commit)
				},
				Key: 's',
			},
			{
				Label: self.c.Tr.CommitMessage,
				OnPress: func() error {
					return self.copyCommitMessageToClipboard(commit)
				},
				Key: 'm',
			},
			{
				Label: self.c.Tr.CommitURL,
				OnPress: func() error {
					return self.copyCommitURLToClipboard(commit)
				},
				Key: 'u',
			},
			{
				Label: self.c.Tr.CommitDiff,
				OnPress: func() error {
					return self.copyCommitDiffToClipboard(commit)
				},
				Key: 'd',
			},
			{
				Label: self.c.Tr.CommitAuthor,
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
	if err := self.c.OS().CopyToClipboard(commit.Sha); err != nil {
		return self.c.Error(err)
	}

	self.c.Toast(fmt.Sprintf("'%s' %s", commit.Sha, self.c.Tr.CopiedToClipboard))
	return nil
}

func (self *BasicCommitsController) copyCommitURLToClipboard(commit *models.Commit) error {
	url, err := self.c.Helpers().Host.GetCommitURL(commit.Sha)
	if err != nil {
		return self.c.Error(err)
	}

	self.c.LogAction(self.c.Tr.Actions.CopyCommitURLToClipboard)
	if err := self.c.OS().CopyToClipboard(url); err != nil {
		return self.c.Error(err)
	}

	self.c.Toast(self.c.Tr.CommitURLCopiedToClipboard)
	return nil
}

func (self *BasicCommitsController) copyCommitDiffToClipboard(commit *models.Commit) error {
	diff, err := self.c.Git().Commit.GetCommitDiff(commit.Sha)
	if err != nil {
		return self.c.Error(err)
	}

	self.c.LogAction(self.c.Tr.Actions.CopyCommitDiffToClipboard)
	if err := self.c.OS().CopyToClipboard(diff); err != nil {
		return self.c.Error(err)
	}

	self.c.Toast(self.c.Tr.CommitDiffCopiedToClipboard)
	return nil
}

func (self *BasicCommitsController) copyAuthorToClipboard(commit *models.Commit) error {
	author, err := self.c.Git().Commit.GetCommitAuthor(commit.Sha)
	if err != nil {
		return self.c.Error(err)
	}

	formattedAuthor := fmt.Sprintf("%s <%s>", author.Name, author.Email)

	self.c.LogAction(self.c.Tr.Actions.CopyCommitAuthorToClipboard)
	if err := self.c.OS().CopyToClipboard(formattedAuthor); err != nil {
		return self.c.Error(err)
	}

	self.c.Toast(self.c.Tr.CommitAuthorCopiedToClipboard)
	return nil
}

func (self *BasicCommitsController) copyCommitMessageToClipboard(commit *models.Commit) error {
	message, err := self.c.Git().Commit.GetCommitMessage(commit.Sha)
	if err != nil {
		return self.c.Error(err)
	}

	self.c.LogAction(self.c.Tr.Actions.CopyCommitMessageToClipboard)
	if err := self.c.OS().CopyToClipboard(message); err != nil {
		return self.c.Error(err)
	}

	self.c.Toast(self.c.Tr.CommitMessageCopiedToClipboard)
	return nil
}

func (self *BasicCommitsController) copyCommitSubjectToClipboard(commit *models.Commit) error {
	message, err := self.c.Git().Commit.GetCommitSubject(commit.Sha)
	if err != nil {
		return self.c.Error(err)
	}

	self.c.LogAction(self.c.Tr.Actions.CopyCommitSubjectToClipboard)
	if err := self.c.OS().CopyToClipboard(message); err != nil {
		return self.c.Error(err)
	}

	self.c.Toast(self.c.Tr.CommitSubjectCopiedToClipboard)
	return nil
}

func (self *BasicCommitsController) openInBrowser(commit *models.Commit) error {
	url, err := self.c.Helpers().Host.GetCommitURL(commit.Sha)
	if err != nil {
		return self.c.Error(err)
	}

	self.c.LogAction(self.c.Tr.Actions.OpenCommitInBrowser)
	if err := self.c.OS().OpenLink(url); err != nil {
		return self.c.Error(err)
	}

	return nil
}

func (self *BasicCommitsController) newBranch(commit *models.Commit) error {
	return self.c.Helpers().Refs.NewBranch(commit.RefName(), commit.Description(), "")
}

func (self *BasicCommitsController) createResetMenu(commit *models.Commit) error {
	return self.c.Helpers().Refs.CreateGitResetMenu(commit.Sha)
}

func (self *BasicCommitsController) checkout(commit *models.Commit) error {
	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.CheckoutCommit,
		Prompt: self.c.Tr.SureCheckoutThisCommit,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.CheckoutCommit)
			return self.c.Helpers().Refs.CheckoutRef(commit.Sha, types.CheckoutRefOptions{})
		},
	})
}

func (self *BasicCommitsController) copyRange(*models.Commit) error {
	return self.c.Helpers().CherryPick.CopyRange(self.context.GetCommits(), self.context)
}

func (self *BasicCommitsController) canCopyCommits(selectedCommits []*models.Commit, startIdx int, endIdx int) *types.DisabledReason {
	for _, commit := range selectedCommits {
		if commit.Sha == "" {
			return &types.DisabledReason{Text: self.c.Tr.CannotCherryPickNonCommit, ShowErrorInPanel: true}
		}

		if commit.IsMerge() {
			return &types.DisabledReason{Text: self.c.Tr.CannotCherryPickMergeCommit, ShowErrorInPanel: true}
		}
	}

	return nil
}

func (self *BasicCommitsController) handleOldCherryPickKey() error {
	msg := utils.ResolvePlaceholderString(self.c.Tr.OldCherryPickKeyWarning,
		map[string]string{
			"copy":  keybindings.Label(self.c.UserConfig.Keybinding.Commits.CherryPickCopy),
			"paste": keybindings.Label(self.c.UserConfig.Keybinding.Commits.PasteCommits),
		})

	return self.c.ErrorMsg(msg)
}

func (self *BasicCommitsController) openDiffTool(commit *models.Commit) error {
	to := commit.RefName()
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(commit.ParentRefName())
	_, err := self.c.RunSubprocess(self.c.Git().Diff.OpenDiffToolCmdObj(
		git_commands.DiffToolCmdOptions{
			Filepath:    ".",
			FromCommit:  from,
			ToCommit:    to,
			Reverse:     reverse,
			IsDirectory: true,
			Staged:      false,
		}))
	return err
}
