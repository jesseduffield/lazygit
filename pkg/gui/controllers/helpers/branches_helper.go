package helpers

import (
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type BranchesHelper struct {
	c              *HelperCommon
	worktreeHelper *WorktreeHelper
}

func NewBranchesHelper(c *HelperCommon, worktreeHelper *WorktreeHelper) *BranchesHelper {
	return &BranchesHelper{
		c:              c,
		worktreeHelper: worktreeHelper,
	}
}

func (self *BranchesHelper) ConfirmLocalDelete(branch *models.Branch) error {
	if self.checkedOutByOtherWorktree(branch) {
		return self.promptWorktreeBranchDelete(branch)
	}

	isMerged, err := self.c.Git().Branch.IsBranchMerged(branch, self.c.Model().MainBranches)
	if err != nil {
		return err
	}

	doDelete := func() error {
		return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(_ gocui.Task) error {
			self.c.LogAction(self.c.Tr.Actions.DeleteLocalBranch)
			if err := self.c.Git().Branch.LocalDelete(branch.Name, true); err != nil {
				return err
			}
			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES}})
		})
	}

	if isMerged {
		return doDelete()
	}

	title := self.c.Tr.ForceDeleteBranchTitle
	message := utils.ResolvePlaceholderString(
		self.c.Tr.ForceDeleteBranchMessage,
		map[string]string{
			"selectedBranchName": branch.Name,
		},
	)

	self.c.Confirm(types.ConfirmOpts{
		Title:  title,
		Prompt: message,
		HandleConfirm: func() error {
			return doDelete()
		},
	})

	return nil
}

func (self *BranchesHelper) ConfirmDeleteRemote(remoteName string, branchName string) error {
	title := utils.ResolvePlaceholderString(
		self.c.Tr.DeleteBranchTitle,
		map[string]string{
			"selectedBranchName": branchName,
		},
	)
	prompt := utils.ResolvePlaceholderString(
		self.c.Tr.DeleteRemoteBranchPrompt,
		map[string]string{
			"selectedBranchName": branchName,
			"upstream":           remoteName,
		},
	)
	self.c.Confirm(types.ConfirmOpts{
		Title:  title,
		Prompt: prompt,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(task gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.DeleteRemoteBranch)
				if err := self.c.Git().Remote.DeleteRemoteBranch(task, remoteName, branchName); err != nil {
					return err
				}
				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
			})
		},
	})

	return nil
}

func (self *BranchesHelper) ConfirmLocalAndRemoteDelete(branch *models.Branch) error {
	if self.checkedOutByOtherWorktree(branch) {
		return self.promptWorktreeBranchDelete(branch)
	}

	isMerged, err := self.c.Git().Branch.IsBranchMerged(branch, self.c.Model().MainBranches)
	if err != nil {
		return err
	}

	prompt := utils.ResolvePlaceholderString(
		self.c.Tr.DeleteLocalAndRemoteBranchPrompt,
		map[string]string{
			"localBranchName":  branch.Name,
			"remoteBranchName": branch.UpstreamBranch,
			"remoteName":       branch.UpstreamRemote,
		},
	)

	if !isMerged {
		prompt += "\n\n" + utils.ResolvePlaceholderString(
			self.c.Tr.ForceDeleteBranchMessage,
			map[string]string{
				"selectedBranchName": branch.Name,
			},
		)
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.DeleteLocalAndRemoteBranch,
		Prompt: prompt,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(task gocui.Task) error {
				// Delete the remote branch first so that we keep the local one
				// in case of failure
				self.c.LogAction(self.c.Tr.Actions.DeleteRemoteBranch)
				if err := self.c.Git().Remote.DeleteRemoteBranch(task, branch.UpstreamRemote, branch.Name); err != nil {
					return err
				}

				self.c.LogAction(self.c.Tr.Actions.DeleteLocalBranch)
				if err := self.c.Git().Branch.LocalDelete(branch.Name, true); err != nil {
					return err
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
			})
		},
	})

	return nil
}

func ShortBranchName(fullBranchName string) string {
	return strings.TrimPrefix(strings.TrimPrefix(fullBranchName, "refs/heads/"), "refs/remotes/")
}

func (self *BranchesHelper) checkedOutByOtherWorktree(branch *models.Branch) bool {
	return git_commands.CheckedOutByOtherWorktree(branch, self.c.Model().Worktrees)
}

func (self *BranchesHelper) worktreeForBranch(branch *models.Branch) (*models.Worktree, bool) {
	return git_commands.WorktreeForBranch(branch, self.c.Model().Worktrees)
}

func (self *BranchesHelper) promptWorktreeBranchDelete(selectedBranch *models.Branch) error {
	worktree, ok := self.worktreeForBranch(selectedBranch)
	if !ok {
		self.c.Log.Error("promptWorktreeBranchDelete out of sync with list of worktrees")
		return nil
	}

	title := utils.ResolvePlaceholderString(self.c.Tr.BranchCheckedOutByWorktree, map[string]string{
		"worktreeName": worktree.Name,
		"branchName":   selectedBranch.Name,
	})
	return self.c.Menu(types.CreateMenuOptions{
		Title: title,
		Items: []*types.MenuItem{
			{
				Label: self.c.Tr.SwitchToWorktree,
				OnPress: func() error {
					return self.worktreeHelper.Switch(worktree, context.LOCAL_BRANCHES_CONTEXT_KEY)
				},
			},
			{
				Label:   self.c.Tr.DetachWorktree,
				Tooltip: self.c.Tr.DetachWorktreeTooltip,
				OnPress: func() error {
					return self.worktreeHelper.Detach(worktree)
				},
			},
			{
				Label: self.c.Tr.RemoveWorktree,
				OnPress: func() error {
					return self.worktreeHelper.Remove(worktree, false)
				},
			},
		},
	})
}
