package helpers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gocui"
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
		return self.promptWorktreeBranchDelete(
			branches[0],
			self.c.Tr.RemoveWorktreeAndDeleteBranch,
			self.c.Tr.DetachWorktreeAndDeleteBranch,
			self.deleteLocalBranchesContinuation(branches),
		)
	}

	return self.confirmForceIfUnmerged(branches, func() error {
		return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(_ gocui.Task) error {
			if err := self.doDeleteLocalBranches(branches); err != nil {
				return err
			}

			self.c.Contexts().Branches.CollapseRangeSelectionToTop()
			self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES}})
			return nil
		})
	})
}

func (self *BranchesHelper) ConfirmDeleteRemote(remoteBranches []*models.RemoteBranch, resetRemoteBranchesSelection bool) error {
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
				self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
				if resetRemoteBranchesSelection {
					self.c.Contexts().RemoteBranches.CollapseRangeSelectionToTop()
				}
				return nil
			})
		},
	})

	return nil
}

func (self *BranchesHelper) ConfirmLocalAndRemoteDelete(branches []*models.Branch) error {
	if len(branches) > 1 {
		if lo.SomeBy(branches, func(branch *models.Branch) bool { return self.checkedOutByOtherWorktree(branch) }) {
			return errors.New(self.c.Tr.SomeBranchesCheckedOutByWorktreeError)
		}
	} else if self.checkedOutByOtherWorktree(branches[0]) {
		return self.promptWorktreeBranchDelete(
			branches[0],
			self.c.Tr.RemoveWorktreeAndDeleteBothBranches,
			self.c.Tr.DetachWorktreeAndDeleteBothBranches,
			self.deleteLocalAndRemoteBranchesContinuation(branches),
		)
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
				if err := self.doDeleteLocalAndRemoteBranches(task, branches); err != nil {
					return err
				}

				self.c.Contexts().Branches.CollapseRangeSelectionToTop()
				self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
				return nil
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

// promptWorktreeBranchDelete handles deleting a branch that's checked out by
// another worktree: the worktree has to be removed or detached first to free the
// branch, so we offer both as menu items. Either way the branch is deleted
// afterwards (that's what the user asked for), via deleteBranches, which knows
// whether to delete just the local branch or the remote one too.
func (self *BranchesHelper) promptWorktreeBranchDelete(
	branch *models.Branch,
	removeLabel string,
	detachLabel string,
	deleteBranches func(gocui.Task) error,
) error {
	worktree, ok := self.worktreeForBranch(branch)
	if !ok {
		self.c.Log.Error("promptWorktreeBranchDelete out of sync with list of worktrees")
		return nil
	}

	title := utils.ResolvePlaceholderString(self.c.Tr.BranchCheckedOutByWorktree, map[string]string{
		"worktreeName": worktree.Name,
		"branchName":   branch.Name,
	})
	return self.c.Menu(types.CreateMenuOptions{
		Title: title,
		Items: []*types.MenuItem{
			{
				Label: removeLabel,
				Keys:  menuKey('r'),
				OnPress: func() error {
					return self.confirmForceIfUnmerged([]*models.Branch{branch}, func() error {
						return self.worktreeHelper.Remove(worktree, deleteBranches)
					})
				},
			},
			{
				Label:   detachLabel,
				Keys:    menuKey('d'),
				Tooltip: self.c.Tr.DetachWorktreeTooltip,
				OnPress: func() error {
					return self.confirmForceIfUnmerged([]*models.Branch{branch}, func() error {
						return self.worktreeHelper.Detach(worktree, deleteBranches)
					})
				},
			},
		},
	})
}

// RemoveWorktreeAndDeleteBranch removes the worktree and deletes the local branch
// it has checked out, force-warning first if the branch isn't fully merged. It's
// the worktrees-panel counterpart to deleting a worktree-checked-out branch from
// the branches panel.
func (self *BranchesHelper) RemoveWorktreeAndDeleteBranch(
	worktree *models.Worktree, branch *models.Branch,
) error {
	branches := []*models.Branch{branch}
	return self.removeWorktreeAndDelete(worktree, branches,
		self.deleteLocalBranchesContinuation(branches))
}

// RemoveWorktreeAndDeleteBothBranches is like RemoveWorktreeAndDeleteBranch but
// also deletes the branch's upstream.
func (self *BranchesHelper) RemoveWorktreeAndDeleteBothBranches(
	worktree *models.Worktree, branch *models.Branch,
) error {
	branches := []*models.Branch{branch}
	return self.removeWorktreeAndDelete(worktree, branches,
		self.deleteLocalAndRemoteBranchesContinuation(branches))
}

func (self *BranchesHelper) removeWorktreeAndDelete(
	worktree *models.Worktree, branches []*models.Branch, deleteBranches func(gocui.Task) error,
) error {
	return self.confirmForceIfUnmerged(branches, func() error {
		return self.worktreeHelper.Remove(worktree, deleteBranches)
	})
}

// confirmForceIfUnmerged runs onConfirm directly if all the branches are fully
// merged, and otherwise shows the force-delete warning first and runs onConfirm
// when the user confirms it.
func (self *BranchesHelper) confirmForceIfUnmerged(branches []*models.Branch, onConfirm func() error) error {
	allBranchesMerged, err := self.allBranchesMerged(branches)
	if err != nil {
		return err
	}
	if allBranchesMerged {
		return onConfirm()
	}

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
		Title:         self.c.Tr.ForceDeleteBranchTitle,
		Prompt:        message,
		HandleConfirm: onConfirm,
	})

	return nil
}

func (self *BranchesHelper) doDeleteLocalBranches(branches []*models.Branch) error {
	self.c.LogAction(self.c.Tr.Actions.DeleteLocalBranch)
	self.logBranchHashes(branches)
	branchNames := lo.Map(branches, func(branch *models.Branch, _ int) string { return branch.Name })
	return self.c.Git().Branch.LocalDelete(branchNames, true)
}

func (self *BranchesHelper) doDeleteLocalAndRemoteBranches(task gocui.Task, branches []*models.Branch) error {
	// Delete the remote branches first so that we keep the local ones
	// in case of failure
	remoteBranches := lo.Map(branches, func(branch *models.Branch, _ int) *models.RemoteBranch {
		return &models.RemoteBranch{Name: branch.UpstreamBranch, RemoteName: branch.UpstreamRemote}
	})
	if err := self.deleteRemoteBranches(remoteBranches, task); err != nil {
		return err
	}

	return self.doDeleteLocalBranches(branches)
}

// deleteLocalBranchesContinuation returns a worktree-removal continuation that
// deletes the local branches and refreshes once the worktree is out of the way.
func (self *BranchesHelper) deleteLocalBranchesContinuation(branches []*models.Branch) func(gocui.Task) error {
	return func(gocui.Task) error {
		if err := self.doDeleteLocalBranches(branches); err != nil {
			return err
		}

		self.c.Contexts().Branches.CollapseRangeSelectionToTop()
		self.c.Refresh(types.RefreshOptions{
			Mode:  types.ASYNC,
			Scope: []types.RefreshableView{types.WORKTREES, types.BRANCHES, types.FILES},
		})
		return nil
	}
}

// deleteLocalAndRemoteBranchesContinuation returns a worktree-removal
// continuation that deletes the local and remote branches and refreshes once the
// worktree is out of the way.
func (self *BranchesHelper) deleteLocalAndRemoteBranchesContinuation(branches []*models.Branch) func(gocui.Task) error {
	return func(task gocui.Task) error {
		if err := self.doDeleteLocalAndRemoteBranches(task, branches); err != nil {
			return err
		}

		self.c.Contexts().Branches.CollapseRangeSelectionToTop()
		self.c.Refresh(types.RefreshOptions{
			Mode:  types.ASYNC,
			Scope: []types.RefreshableView{types.WORKTREES, types.BRANCHES, types.REMOTES, types.FILES},
		})
		return nil
	}
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

func (self *BranchesHelper) logBranchHashes(branches []*models.Branch) {
	for _, branch := range branches {
		msg := utils.ResolvePlaceholderString(
			self.c.Tr.Log.DeletingBranch,
			map[string]string{
				"branchName": branch.Name,
				"hash":       branch.CommitHash,
			},
		)

		self.c.LogCommand(msg, false)
	}
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

func (self *BranchesHelper) PostFetchRefresh(fetchErr error, background bool) error {
	scope := []types.RefreshableView{
		types.BRANCHES, types.COMMITS, types.REMOTES, types.TAGS, types.PULL_REQUESTS,
	}
	// AutoForwardBranches needs a fresh worktree model to skip branches that are checked out elsewhere.
	if self.c.UserConfig().Git.AutoForwardBranches != "none" {
		scope = append(scope, types.WORKTREES)
	}
	self.c.Refresh(types.RefreshOptions{Scope: scope, Mode: types.SYNC, Background: background})
	if fetchErr != nil {
		return fetchErr
	}
	return self.AutoForwardBranches()
}

func (self *BranchesHelper) AutoForwardBranches() error {
	if self.c.UserConfig().Git.AutoForwardBranches == "none" {
		return nil
	}

	branches := self.c.Model().Branches
	if len(branches) == 0 {
		return nil
	}

	allBranches := self.c.UserConfig().Git.AutoForwardBranches == "allBranches"
	updateCommands := ""
	// The first branch is the currently checked out branch; skip it
	for _, branch := range branches[1:] {
		if branch.RemoteBranchStoredLocally() &&
			!self.checkedOutByOtherWorktree(branch) &&
			(allBranches || lo.Contains(self.c.UserConfig().Git.MainBranches, branch.Name)) {
			isStrictlyBehind := branch.IsBehindForPull() && !branch.IsAheadForPull()
			if isStrictlyBehind {
				updateCommands += fmt.Sprintf("update %s %s %s\n", branch.FullRefName(), branch.FullUpstreamRefName(), branch.CommitHash)
			}
		}
	}

	if updateCommands == "" {
		return nil
	}

	self.c.LogAction(self.c.Tr.Actions.AutoForwardBranches)
	self.c.LogCommand(strings.TrimRight(updateCommands, "\n"), false)
	err := self.c.Git().Branch.UpdateBranchRefs(updateCommands)

	self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES}, Mode: types.SYNC})

	return err
}
