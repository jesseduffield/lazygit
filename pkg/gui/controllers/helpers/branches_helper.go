package helpers

import (
	"errors"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
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

func (self *BranchesHelper) ConfirmLocalDelete(branches []*models.Branch) error {
	if len(branches) > 1 {
		if lo.SomeBy(branches, func(branch *models.Branch) bool { return self.checkedOutByOtherWorktree(branch) }) {
			return errors.New(self.c.Tr.SomeBranchesCheckedOutByWorktreeError)
		}
	} else if self.checkedOutByOtherWorktree(branches[0]) {
		return self.promptWorktreeBranchDelete(branches[0])
	}

	allBranchesMerged, err := self.allBranchesMerged(branches)
	if err != nil {
		return err
	}

	doDelete := func() error {
		return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(_ gocui.Task) error {
			self.c.LogAction(self.c.Tr.Actions.DeleteLocalBranch)
			branchNames := lo.Map(branches, func(branch *models.Branch, _ int) string { return branch.Name })
			if err := self.c.Git().Branch.LocalDelete(branchNames, true); err != nil {
				return err
			}
			selectionStart, _ := self.c.Contexts().Branches.GetSelectionRange()
			self.c.Contexts().Branches.SetSelectedLineIdx(selectionStart)
			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES}})
		})
	}

	if allBranchesMerged {
		return doDelete()
	}

	title := self.c.Tr.ForceDeleteBranchTitle
	var message string
	if len(branches) == 1 {
		message = utils.ResolvePlaceholderString(
			self.c.Tr.ForceDeleteBranchMessage,
			map[string]string{
				"selectedBranchName": branches[0].Name,
			},
		)
	} else {
		message = self.c.Tr.ForceDeleteBranchesMessage
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:  title,
		Prompt: message,
		HandleConfirm: func() error {
			return doDelete()
		},
	})

	return nil
}

func (self *BranchesHelper) ConfirmDeleteRemote(remoteBranches []*models.RemoteBranch) error {
	var title string
	if len(remoteBranches) == 1 {
		title = utils.ResolvePlaceholderString(
			self.c.Tr.DeleteBranchTitle,
			map[string]string{
				"selectedBranchName": remoteBranches[0].Name,
			},
		)
	} else {
		title = self.c.Tr.DeleteBranchesTitle
	}
	var prompt string
	if len(remoteBranches) == 1 {
		prompt = utils.ResolvePlaceholderString(
			self.c.Tr.DeleteRemoteBranchPrompt,
			map[string]string{
				"selectedBranchName": remoteBranches[0].Name,
				"upstream":           remoteBranches[0].RemoteName,
			},
		)
	} else {
		prompt = self.c.Tr.DeleteRemoteBranchesPrompt
	}
	self.c.Confirm(types.ConfirmOpts{
		Title:  title,
		Prompt: prompt,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(task gocui.Task) error {
				if err := self.deleteRemoteBranches(remoteBranches, task); err != nil {
					return err
				}
				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
			})
		},
	})

	return nil
}

func (self *BranchesHelper) ConfirmLocalAndRemoteDelete(branches []*models.Branch) error {
	if lo.SomeBy(branches, func(branch *models.Branch) bool { return self.checkedOutByOtherWorktree(branch) }) {
		return errors.New(self.c.Tr.SomeBranchesCheckedOutByWorktreeError)
	}

	allBranchesMerged, err := self.allBranchesMerged(branches)
	if err != nil {
		return err
	}

	var prompt string
	if len(branches) == 1 {
		prompt = utils.ResolvePlaceholderString(
			self.c.Tr.DeleteLocalAndRemoteBranchPrompt,
			map[string]string{
				"localBranchName":  branches[0].Name,
				"remoteBranchName": branches[0].UpstreamBranch,
				"remoteName":       branches[0].UpstreamRemote,
			},
		)
	} else {
		prompt = self.c.Tr.DeleteLocalAndRemoteBranchesPrompt
	}

	if !allBranchesMerged {
		if len(branches) == 1 {
			prompt += "\n\n" + utils.ResolvePlaceholderString(
				self.c.Tr.ForceDeleteBranchMessage,
				map[string]string{
					"selectedBranchName": branches[0].Name,
				},
			)
		} else {
			prompt += "\n\n" + self.c.Tr.ForceDeleteBranchesMessage
		}
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.DeleteLocalAndRemoteBranch,
		Prompt: prompt,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(task gocui.Task) error {
				// Delete the remote branches first so that we keep the local ones
				// in case of failure
				remoteBranches := lo.Map(branches, func(branch *models.Branch, _ int) *models.RemoteBranch {
					return &models.RemoteBranch{Name: branch.UpstreamBranch, RemoteName: branch.UpstreamRemote}
				})
				if err := self.deleteRemoteBranches(remoteBranches, task); err != nil {
					return err
				}

				self.c.LogAction(self.c.Tr.Actions.DeleteLocalBranch)
				branchNames := lo.Map(branches, func(branch *models.Branch, _ int) string { return branch.Name })
				if err := self.c.Git().Branch.LocalDelete(branchNames, true); err != nil {
					return err
				}

				selectionStart, _ := self.c.Contexts().Branches.GetSelectionRange()
				self.c.Contexts().Branches.SetSelectedLineIdx(selectionStart)

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

func (self *BranchesHelper) allBranchesMerged(branches []*models.Branch) (bool, error) {
	allBranchesMerged := true
	for _, branch := range branches {
		isMerged, err := self.c.Git().Branch.IsBranchMerged(branch, self.c.Model().MainBranches)
		if err != nil {
			return false, err
		}
		if !isMerged {
			allBranchesMerged = false
			break
		}
	}
	return allBranchesMerged, nil
}

func (self *BranchesHelper) deleteRemoteBranches(remoteBranches []*models.RemoteBranch, task gocui.Task) error {
	remotes := lo.GroupBy(remoteBranches, func(branch *models.RemoteBranch) string { return branch.RemoteName })
	for remote, branches := range remotes {
		self.c.LogAction(self.c.Tr.Actions.DeleteRemoteBranch)
		branchNames := lo.Map(branches, func(branch *models.RemoteBranch, _ int) string { return branch.Name })
		if err := self.c.Git().Remote.DeleteRemoteBranch(task, remote, branchNames); err != nil {
			return err
		}
	}
	return nil
}
