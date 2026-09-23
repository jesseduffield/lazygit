package helpers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gocui"
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

			self.c.OnUIThread(func() error {
				self.c.Contexts().Branches.CollapseRangeSelectionToTop()
				return nil
			})
			self.c.RefreshFromWorker(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES}})
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
				self.c.RefreshFromWorker(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
				if resetRemoteBranchesSelection {
					self.c.OnUIThread(func() error {
						self.c.Contexts().RemoteBranches.CollapseRangeSelectionToTop()
						return nil
					})
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

				self.c.OnUIThread(func() error {
					self.c.Contexts().Branches.CollapseRangeSelectionToTop()
					return nil
				})
				self.c.RefreshFromWorker(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
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

		self.c.OnUIThread(func() error {
			self.c.Contexts().Branches.CollapseRangeSelectionToTop()
			return nil
		})
		self.c.RefreshFromWorker(types.RefreshOptions{
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

		self.c.OnUIThread(func() error {
			self.c.Contexts().Branches.CollapseRangeSelectionToTop()
			return nil
		})
		self.c.RefreshFromWorker(types.RefreshOptions{
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

// fetchGeneration must be the repo generation from when the fetch started,
// captured by the caller before running the fetch: the background fetch
// doesn't block repo switching and is a network call, so the window in which
// the user can switch repos spans the whole fetch, not just this refresh.
func (self *BranchesHelper) PostFetchRefresh(fetchErr error, background bool, fetchGeneration int) error {
	scope := []types.RefreshableView{
		types.BRANCHES, types.COMMITS, types.REMOTES, types.TAGS, types.PULL_REQUESTS,
	}
	// AutoForwardBranches needs a fresh worktree model to skip branches that are checked out elsewhere.
	if self.c.UserConfig().Git.AutoForwardBranches != "none" {
		scope = append(scope, types.WORKTREES)
	}
	// AutoForwardBranches reads Model.Branches, which the branches refresh writes
	// via a bounce, so it has to run in Then rather than right after Refresh
	// returns (where it would still see the previous branches).
	self.c.RefreshFromWorker(types.RefreshOptions{
		Scope:      scope,
		Background: background,
		Then: func() error {
			if fetchErr != nil {
				return nil
			}
			// Then callbacks are not generation-guarded, so check explicitly:
			// if the repo was switched since the fetch started, don't forward
			// this repo's branches on the strength of another repo's fetch.
			if self.c.State().GetRepoGeneration() != fetchGeneration {
				return nil
			}
			err := self.AutoForwardBranches(background)
			if background && err != nil {
				// The background poller discards this return value, so surface
				// the error in the log rather than as a popup for background work.
				self.c.Log.Error(err)
				return nil
			}
			return err
		},
	})
	return fetchErr
}

// One of the branches that a fast-forward is about to bring to its upstream
type branchToForward struct {
	branch *models.Branch
	// the worktree that the branch is checked out in, nil if there is none
	worktree *models.Worktree
	// whether the branch has to be reset to its upstream because it has
	// diverged from it, as opposed to being moved forward
	reset bool
}

// Updates the given branches to their upstream branches, fetching those first.
// A branch that is behind its upstream is moved forward to it; one that has
// diverged from it is reset to it, as long as it has no commits of its own. If
// any of the branches can't be updated, none of them is.
func (self *BranchesHelper) FastForwardBranches(branches []*models.Branch) error {
	// The worktrees come from the model, so they have to be looked up here,
	// before the work moves to a worker
	toForward := lo.Map(branches, func(branch *models.Branch, _ int) *branchToForward {
		worktree, _ := self.worktreeForBranch(branch)
		return &branchToForward{branch: branch, worktree: worktree}
	})
	anyCheckedOut := lo.SomeBy(toForward, func(f *branchToForward) bool { return f.worktree != nil })

	return self.withFastForwardingStatus(branches, func(task gocui.Task) error {
		defer func() {
			if anyCheckedOut {
				// The files of those worktrees have changed as well
				self.c.RefreshFromWorker(types.RefreshOptions{})
			} else {
				self.c.RefreshFromWorker(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES}})
			}
		}()

		self.c.LogAction(self.c.Tr.Actions.FastForwardBranch)

		if err := self.fetchUpstreamBranches(task, branches); err != nil {
			return err
		}

		// Look at all the branches before moving any of them, so that one we
		// have to refuse leaves the others alone too
		for _, f := range toForward {
			if err := self.planForwardingBranch(f); err != nil {
				return err
			}
		}

		return self.forwardBranches(toForward)
	})
}

// Runs f with all the given branches shown as being fast-forwarded
func (self *BranchesHelper) withFastForwardingStatus(branches []*models.Branch, f func(gocui.Task) error) error {
	return self.c.WithInlineStatus(branches[0], types.ItemOperationFastForwarding, context.LOCAL_BRANCHES_CONTEXT_KEY, func(task gocui.Task) error {
		for _, branch := range branches[1:] {
			self.c.State().SetItemOperation(branch, types.ItemOperationFastForwarding)
		}
		defer func() {
			for _, branch := range branches[1:] {
				self.c.State().ClearItemOperation(branch)
			}
		}()

		return f(task)
	})
}

func (self *BranchesHelper) fetchUpstreamBranches(task gocui.Task, branches []*models.Branch) error {
	remotes := lo.Uniq(lo.Map(branches, func(branch *models.Branch, _ int) string {
		return branch.UpstreamRemote
	}))

	for _, remote := range remotes {
		remoteBranches := []string{}
		for _, branch := range branches {
			if branch.UpstreamRemote == remote {
				remoteBranches = append(remoteBranches, branch.UpstreamBranch)
			}
		}

		if err := self.c.Git().Sync.FetchRemoteBranches(task, remote, remoteBranches); err != nil {
			return err
		}
	}

	return nil
}

// Works out whether the branch has to be reset to its upstream, and returns an
// error if it can't be brought there at all.
func (self *BranchesHelper) planForwardingBranch(f *branchToForward) error {
	f.reset = !self.c.Git().Branch.IsAncestor(
		f.branch.FullRefName(), f.branch.FullUpstreamRefName())
	if !f.reset {
		return nil
	}

	// Moving the branch to its upstream means giving up the commits it is
	// ahead by, so make sure that none of them is ours
	hasLocalOnlyCommits, err := self.c.Git().Branch.HasLocalOnlyCommits(f.branch)
	if err != nil {
		return err
	}
	if hasLocalOnlyCommits {
		return errors.New(utils.ResolvePlaceholderString(
			self.c.Tr.FwdLocalOnlyCommits,
			map[string]string{"branchName": f.branch.Name},
		))
	}

	if f.worktree != nil {
		// Resetting the branch changes the files of the worktree under the
		// user's feet, so only do it while they have no changes of their own
		// there
		worktreeGitDir, worktreePath := self.worktreeArgs(f.worktree)
		hasChanges, err := self.c.Git().WorkingTree.HasChangesToTrackedFiles(worktreeGitDir, worktreePath)
		if err != nil {
			return err
		}
		if hasChanges {
			return errors.New(utils.ResolvePlaceholderString(
				self.c.Tr.FwdUncommittedChanges,
				map[string]string{"branchName": f.branch.Name},
			))
		}
	}

	return nil
}

func (self *BranchesHelper) forwardBranches(toForward []*branchToForward) error {
	// The branches that aren't checked out anywhere are nothing but refs to
	// update, so they can all be done in one go
	updateCommands := ""
	for _, f := range toForward {
		if f.worktree == nil {
			updateCommands += fmt.Sprintf("update %s %s %s\n",
				f.branch.FullRefName(), f.branch.FullUpstreamRefName(), f.branch.CommitHash)
		}
	}

	if updateCommands != "" {
		self.c.LogCommand(strings.TrimRight(updateCommands, "\n"), false)
		if err := self.c.Git().Branch.UpdateBranchRefs(updateCommands); err != nil {
			return err
		}
	}

	// A branch that is checked out somewhere needs the files of that worktree
	// to be updated along with it
	for _, f := range toForward {
		if f.worktree == nil {
			continue
		}

		worktreeGitDir, worktreePath := self.worktreeArgs(f.worktree)

		var err error
		if f.reset {
			err = self.c.Git().WorkingTree.ResetKeep(
				f.branch.FullUpstreamRefName(), worktreeGitDir, worktreePath)
		} else {
			err = self.c.Git().Branch.FastForwardMerge(
				f.branch.FullUpstreamRefName(), worktreeGitDir, worktreePath)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// Returns the git dir and the path to pass for the given worktree; both are
// empty for the current one, which git commands use by default anyway.
func (self *BranchesHelper) worktreeArgs(worktree *models.Worktree) (string, string) {
	if worktree.IsCurrent {
		return "", ""
	}

	return worktree.GitDir, worktree.Path
}

func (self *BranchesHelper) AutoForwardBranches(background bool) error {
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

	self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES}, Background: background})

	return err
}
