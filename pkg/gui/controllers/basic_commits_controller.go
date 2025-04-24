package controllers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context/traits"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
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
	GetSelectionRangeAndMode() (int, int, traits.RangeSelectMode)
	SetSelectionRangeAndMode(int, int, traits.RangeSelectMode)
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
		ListControllerTrait: NewListControllerTrait(
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
			// Putting this in BasicCommitsController even though we really only want it in the commits
			// panel. But I find it important that this ends up next to "New Branch", and I couldn't
			// find another way to achieve this. It's not such a big deal to have it in subcommits and
			// reflog too, I'd say.
			Key:               opts.GetKey(opts.Config.Branches.MoveCommitsToNewBranch),
			Handler:           self.c.Helpers().Refs.MoveCommitsToNewBranch,
			GetDisabledReason: self.c.Helpers().Refs.CanMoveCommitsToNewBranch,
			Description:       self.c.Tr.MoveCommitsToNewBranch,
			Tooltip:           self.c.Tr.MoveCommitsToNewBranchTooltip,
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
		{
			Key:               opts.GetKey(opts.Config.Commits.SelectCommitsOfCurrentBranch),
			Handler:           self.selectCommitsOfCurrentBranch,
			GetDisabledReason: self.require(self.canSelectCommitsOfCurrentBranch),
			Description:       self.c.Tr.SelectCommitsOfCurrentBranch,
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

func (self *BasicCommitsController) getCommitMessageBody(hash string) string {
	commitMessageBody, err := self.c.Git().Commit.GetCommitMessage(hash)
	if err != nil {
		return ""
	}
	_, body := self.c.Helpers().Commits.SplitCommitMessageAndDescription(commitMessageBody)
	return body
}

func (self *BasicCommitsController) copyCommitAttribute(commit *models.Commit) error {
	commitMessageBody := self.getCommitMessageBody(commit.Hash())
	var commitMessageBodyDisabled *types.DisabledReason
	if commitMessageBody == "" {
		commitMessageBodyDisabled = &types.DisabledReason{
			Text: self.c.Tr.CommitHasNoMessageBody,
		}
	}

	items := []*types.MenuItem{
		{
			Label: self.c.Tr.CommitHash,
			OnPress: func() error {
				return self.copyCommitHashToClipboard(commit)
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
			Label:          self.c.Tr.CommitMessageBody,
			DisabledReason: commitMessageBodyDisabled,
			OnPress: func() error {
				return self.copyCommitMessageBodyToClipboard(commitMessageBody)
			},
			Key: 'b',
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
	}

	commitTagsItem := types.MenuItem{
		Label: self.c.Tr.CommitTags,
		OnPress: func() error {
			return self.copyCommitTagsToClipboard(commit)
		},
		Key: 't',
	}

	if len(commit.Tags) == 0 {
		commitTagsItem.DisabledReason = &types.DisabledReason{Text: self.c.Tr.CommitHasNoTags}
	}

	items = append(items, &commitTagsItem)

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.Actions.CopyCommitAttributeToClipboard,
		Items: items,
	})
}

func (self *BasicCommitsController) copyCommitHashToClipboard(commit *models.Commit) error {
	self.c.LogAction(self.c.Tr.Actions.CopyCommitHashToClipboard)
	if err := self.c.OS().CopyToClipboard(commit.Hash()); err != nil {
		return err
	}

	self.c.Toast(fmt.Sprintf("'%s' %s", commit.Hash(), self.c.Tr.CopiedToClipboard))
	return nil
}

func (self *BasicCommitsController) copyCommitURLToClipboard(commit *models.Commit) error {
	url, err := self.c.Helpers().Host.GetCommitURL(commit.Hash())
	if err != nil {
		return err
	}

	self.c.LogAction(self.c.Tr.Actions.CopyCommitURLToClipboard)
	if err := self.c.OS().CopyToClipboard(url); err != nil {
		return err
	}

	self.c.Toast(self.c.Tr.CommitURLCopiedToClipboard)
	return nil
}

func (self *BasicCommitsController) copyCommitDiffToClipboard(commit *models.Commit) error {
	diff, err := self.c.Git().Commit.GetCommitDiff(commit.Hash())
	if err != nil {
		return err
	}

	self.c.LogAction(self.c.Tr.Actions.CopyCommitDiffToClipboard)
	if err := self.c.OS().CopyToClipboard(diff); err != nil {
		return err
	}

	self.c.Toast(self.c.Tr.CommitDiffCopiedToClipboard)
	return nil
}

func (self *BasicCommitsController) copyAuthorToClipboard(commit *models.Commit) error {
	author, err := self.c.Git().Commit.GetCommitAuthor(commit.Hash())
	if err != nil {
		return err
	}

	formattedAuthor := fmt.Sprintf("%s <%s>", author.Name, author.Email)

	self.c.LogAction(self.c.Tr.Actions.CopyCommitAuthorToClipboard)
	if err := self.c.OS().CopyToClipboard(formattedAuthor); err != nil {
		return err
	}

	self.c.Toast(self.c.Tr.CommitAuthorCopiedToClipboard)
	return nil
}

func (self *BasicCommitsController) copyCommitMessageToClipboard(commit *models.Commit) error {
	message, err := self.c.Git().Commit.GetCommitMessage(commit.Hash())
	if err != nil {
		return err
	}

	self.c.LogAction(self.c.Tr.Actions.CopyCommitMessageToClipboard)
	if err := self.c.OS().CopyToClipboard(message); err != nil {
		return err
	}

	self.c.Toast(self.c.Tr.CommitMessageCopiedToClipboard)
	return nil
}

func (self *BasicCommitsController) copyCommitMessageBodyToClipboard(commitMessageBody string) error {
	self.c.LogAction(self.c.Tr.Actions.CopyCommitMessageBodyToClipboard)
	if err := self.c.OS().CopyToClipboard(commitMessageBody); err != nil {
		return err
	}

	self.c.Toast(self.c.Tr.CommitMessageBodyCopiedToClipboard)
	return nil
}

func (self *BasicCommitsController) copyCommitSubjectToClipboard(commit *models.Commit) error {
	message, err := self.c.Git().Commit.GetCommitSubject(commit.Hash())
	if err != nil {
		return err
	}

	self.c.LogAction(self.c.Tr.Actions.CopyCommitSubjectToClipboard)
	if err := self.c.OS().CopyToClipboard(message); err != nil {
		return err
	}

	self.c.Toast(self.c.Tr.CommitSubjectCopiedToClipboard)
	return nil
}

func (self *BasicCommitsController) copyCommitTagsToClipboard(commit *models.Commit) error {
	message := strings.Join(commit.Tags, "\n")

	self.c.LogAction(self.c.Tr.Actions.CopyCommitTagsToClipboard)
	if err := self.c.OS().CopyToClipboard(message); err != nil {
		return err
	}

	self.c.Toast(self.c.Tr.CommitTagsCopiedToClipboard)
	return nil
}

func (self *BasicCommitsController) openInBrowser(commit *models.Commit) error {
	url, err := self.c.Helpers().Host.GetCommitURL(commit.Hash())
	if err != nil {
		return err
	}

	self.c.LogAction(self.c.Tr.Actions.OpenCommitInBrowser)
	if err := self.c.OS().OpenLink(url); err != nil {
		return err
	}

	return nil
}

func (self *BasicCommitsController) newBranch(commit *models.Commit) error {
	return self.c.Helpers().Refs.NewBranch(commit.RefName(), commit.Description(), "")
}

func (self *BasicCommitsController) createResetMenu(commit *models.Commit) error {
	return self.c.Helpers().Refs.CreateGitResetMenu(commit.Hash())
}

func (self *BasicCommitsController) checkout(commit *models.Commit) error {
	return self.c.Helpers().Refs.CreateCheckoutMenu(commit)
}

func (self *BasicCommitsController) copyRange(*models.Commit) error {
	return self.c.Helpers().CherryPick.CopyRange(self.context.GetCommits(), self.context)
}

func (self *BasicCommitsController) canCopyCommits(selectedCommits []*models.Commit, startIdx int, endIdx int) *types.DisabledReason {
	for _, commit := range selectedCommits {
		if commit.Hash() == "" {
			return &types.DisabledReason{Text: self.c.Tr.CannotCherryPickNonCommit, ShowErrorInPanel: true}
		}
	}

	return nil
}

func (self *BasicCommitsController) handleOldCherryPickKey() error {
	msg := utils.ResolvePlaceholderString(self.c.Tr.OldCherryPickKeyWarning,
		map[string]string{
			"copy":  keybindings.Label(self.c.UserConfig().Keybinding.Commits.CherryPickCopy),
			"paste": keybindings.Label(self.c.UserConfig().Keybinding.Commits.PasteCommits),
		})

	return errors.New(msg)
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

func (self *BasicCommitsController) canSelectCommitsOfCurrentBranch() *types.DisabledReason {
	if index := self.findFirstCommitAfterCurrentBranch(); index <= 0 {
		return &types.DisabledReason{Text: self.c.Tr.NoCommitsThisBranch}
	}

	return nil
}

func (self *BasicCommitsController) findFirstCommitAfterCurrentBranch() int {
	_, index, ok := lo.FindIndexOf(self.context.GetCommits(), func(c *models.Commit) bool {
		return c.IsMerge() || c.Status == models.StatusMerged
	})

	if !ok {
		return 0
	}

	return index
}

func (self *BasicCommitsController) selectCommitsOfCurrentBranch() error {
	index := self.findFirstCommitAfterCurrentBranch()
	if index <= 0 {
		return nil
	}

	_, _, mode := self.context.GetSelectionRangeAndMode()
	if mode != traits.RangeSelectModeSticky {
		// If we are in sticky range mode already, keep that; otherwise, open a non-sticky range
		mode = traits.RangeSelectModeNonSticky
	}
	// Create the range from bottom to top, so that when you cancel the range,
	// the head commit is selected
	self.context.SetSelectionRangeAndMode(0, index-1, mode)
	self.context.HandleFocus(types.OnFocusOpts{})
	return nil
}
