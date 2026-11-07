package helpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/stefanhaller/git-todo-parser/todo"
)

type MergeAndRebaseHelper struct {
	c *HelperCommon

	// Whether the "continue the rebase/merge?" prompt is currently on screen.
	// We use this to auto-dismiss it if the operation stops being in the state
	// that the prompt is offering to act on (e.g. it was continued or aborted
	// externally), so the user isn't left with a stale prompt. Only accessed on
	// the UI thread.
	continueRebasePromptShowing bool
}

func NewMergeAndRebaseHelper(
	c *HelperCommon,
) *MergeAndRebaseHelper {
	return &MergeAndRebaseHelper{
		c: c,
	}
}

type RebaseOption string

const (
	REBASE_OPTION_CONTINUE string = "continue"
	REBASE_OPTION_ABORT    string = "abort"
	REBASE_OPTION_SKIP     string = "skip"
)

func (self *MergeAndRebaseHelper) CreateRebaseOptionsMenu() error {
	type optionAndKey struct {
		option string
		keys   []gocui.Key
	}

	options := []optionAndKey{
		{option: REBASE_OPTION_CONTINUE, keys: menuKey('c')},
		{option: REBASE_OPTION_ABORT, keys: menuKey('a')},
	}

	if self.c.Git().Status.WorkingTreeState().CanSkip() {
		options = append(options, optionAndKey{
			option: REBASE_OPTION_SKIP, keys: menuKey('s'),
		})
	}

	menuItems := lo.Map(options, func(row optionAndKey, _ int) *types.MenuItem {
		return &types.MenuItem{
			Label: row.option,
			OnPress: func() error {
				if row.option == REBASE_OPTION_ABORT {
					return self.abortMergeOrRebase(false)
				}
				return self.genericMergeCommand(row.option)
			},
			Keys: row.keys,
		}
	})

	title := self.c.Git().Status.WorkingTreeState().OptionsMenuTitle(self.c.Tr)
	return self.c.Menu(types.CreateMenuOptions{Title: title, Items: menuItems})
}

func (self *MergeAndRebaseHelper) ContinueRebase() error {
	return self.genericMergeCommand(REBASE_OPTION_CONTINUE)
}

func (self *MergeAndRebaseHelper) genericMergeCommand(command string) error {
	// The menu/prompt/confirm handlers that reach here run on the UI thread and
	// spin up a worker (via the waiting status below) to do the actual work.
	return self.genericMergeCommandImpl(command, true, false)
}

// genericMergeCommandImpl runs a merge/rebase continue/skip/abort and handles
// the result. Continuing can be slow (it may replay many commits), so the
// non-subprocess path runs on a worker with a waiting status.
//
// showWaitingStatus is false only for the recursive auto-skip in
// CheckMergeOrRebaseWithRefreshOptions and for aborting after discarding
// changes in abortMergeOrRebase; both already run on a worker, so they must not
// spin up a second waiting status. calledFromWorker is used only by the
// subprocess path below: it's true for that recursive worker skip and false for
// genericMergeCommand's UI-thread invocation, so the post-action refresh picks
// RefreshFromWorker vs Refresh correctly.
func (self *MergeAndRebaseHelper) genericMergeCommandImpl(command string, showWaitingStatus bool, calledFromWorker bool) error {
	status := self.c.Git().Status.WorkingTreeState()

	if status.None() {
		return errors.New(self.c.Tr.NotMergingOrRebasing)
	}

	self.c.LogAction(fmt.Sprintf("Merge/Rebase: %s", command))
	effectiveStatus := status.Effective()
	if effectiveStatus == models.WORKING_TREE_STATE_REBASING {
		todoFile, err := os.ReadFile(
			filepath.Join(self.c.Git().RepoPaths.WorktreeGitDirPath(), "rebase-merge/git-rebase-todo"),
		)

		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			self.c.LogCommand(string(todoFile), false)
		}
	}

	isMerging := effectiveStatus == models.WORKING_TREE_STATE_MERGING ||
		effectiveStatus == models.WORKING_TREE_STATE_SQUASH_MERGING
	selectHeadCommitOnSuccess := command == REBASE_OPTION_CONTINUE && isMerging

	// we should end up with a command like 'git merge --continue'

	// it's impossible for a rebase to require a commit so we'll use a subprocess only if it's a merge
	needsSubprocess := (isMerging && command != REBASE_OPTION_ABORT && self.c.UserConfig().Git.Merging.ManualCommit) ||
		// but we'll also use a subprocess if we have exec todos; those are likely to be lengthy build
		// tasks whose output the user will want to see in the terminal
		(effectiveStatus == models.WORKING_TREE_STATE_REBASING && command != REBASE_OPTION_ABORT && self.hasExecTodos(calledFromWorker))

	if needsSubprocess {
		// TODO: see if we should be calling more of the code from self.Git.Rebase.GenericMergeOrRebaseAction
		success, err := self.c.RunSubprocess(self.c.Git().Rebase.GenericMergeOrRebaseActionCmdObj(effectiveStatus, command))
		self.refreshAfterMergeOrRebase(types.RefreshOptions{
			CommitSelection: commitSelectionAfterMerge(success && selectHeadCommitOnSuccess),
		}, calledFromWorker)
		self.RecordWhetherMergeOrRebaseStartedInLazygit()
		return err
	}

	// runAction always ends up on a worker: either the waiting status below spins
	// one up, or we're the recursive auto-skip reached from
	// CheckMergeOrRebaseWithRefreshOptions or the abort in abortMergeOrRebase,
	// which already run on one.
	runAction := func() error {
		result := self.c.Git().Rebase.GenericMergeOrRebaseAction(effectiveStatus, command)
		return self.CheckMergeOrRebaseWithRefreshOptions(result,
			types.RefreshOptions{
				CommitSelection: commitSelectionAfterMerge(result == nil && selectHeadCommitOnSuccess),
			})
	}

	if showWaitingStatus {
		return self.c.WithWaitingStatus(status.Title(self.c.Tr), func(gocui.Task) error {
			return runAction()
		})
	}
	return runAction()
}

// commitSelectionAfterMerge maps whether a merge/rebase/pull created a new
// commit at HEAD to the corresponding commit-selection behavior: select that
// new commit, or otherwise keep the previous selection by hash.
func commitSelectionAfterMerge(createdNewCommit bool) types.CommitSelectionBehavior {
	if createdNewCommit {
		return types.SelectHeadCommit
	}
	return types.KeepCommitSelectionByHash
}

func (self *MergeAndRebaseHelper) hasExecTodos(calledFromWorker bool) bool {
	check := func() bool {
		for _, commit := range self.c.Model().Commits {
			if !commit.IsTODO() {
				break
			}
			if commit.Action == todo.Exec {
				return true
			}
		}
		return false
	}

	// This reads the model, which is only safe on the UI thread, so bounce there
	// when we're being called from a worker.
	if !calledFromWorker {
		return check()
	}

	result := false
	_ = self.c.GocuiGui().OnUIThreadAndWait(func() {
		result = check()
	})
	return result
}

var conflictStrings = []string{
	"Failed to merge in the changes",
	"When you have resolved this problem",
	"fix conflicts",
	"Resolve all conflicts manually",
	"Merge conflict in file",
	"hint: after resolving the conflicts",
	"CONFLICT (content):",
}

func isMergeConflictErr(errStr string) bool {
	for _, str := range conflictStrings {
		if strings.Contains(errStr, str) {
			return true
		}
	}

	return false
}

// RecordWhetherMergeOrRebaseStartedInLazygit is called right after we run a
// merge/rebase/cherry-pick/revert step. If it left an operation in progress,
// that operation is one we started, which is what later lets us auto-prompt to
// continue it once its conflicts are resolved. If nothing is in progress
// anymore (the step completed or aborted the operation), we clear the flag.
func (self *MergeAndRebaseHelper) RecordWhetherMergeOrRebaseStartedInLazygit() {
	self.c.State().GetRepoState().SetMergeOrRebaseStartedInLazygit(
		self.c.Git().Status.WorkingTreeState().Any())
}

// CheckMergeOrRebaseWithRefreshOptions handles the result of a merge/rebase
// step and refreshes. It always runs on a worker (the WithWaitingStatus /
// WithWaitingStatusBlockingInput / WithInlineStatus handlers).
func (self *MergeAndRebaseHelper) CheckMergeOrRebaseWithRefreshOptions(result error, refreshOptions types.RefreshOptions) error {
	self.refreshAfterMergeOrRebase(refreshOptions, true)

	self.RecordWhetherMergeOrRebaseStartedInLazygit()

	if result == nil {
		return nil
	} else if strings.Contains(result.Error(), "No changes - did you forget to use") {
		if self.c.Git().Status.WorkingTreeState().ApplyingPatches {
			// "git am" says this not only if resolving the conflicts of a patch
			// left nothing to commit, but also if the patch didn't apply at
			// all. Without --3way, it leaves the files alone then, and skipping
			// would lose the patch. Show git's explanation instead.
			return result
		}
		return self.genericMergeCommandImpl(REBASE_OPTION_SKIP, false, true)
	} else if strings.Contains(result.Error(), "The previous cherry-pick is now empty") {
		return self.genericMergeCommandImpl(REBASE_OPTION_SKIP, false, true)
	} else if strings.Contains(result.Error(), "No rebase in progress?") {
		// assume in this case that we're already done
		return nil
	}
	return self.CheckForConflicts(result)
}

// refreshAfterMergeOrRebase issues the post-action refresh on the entry point
// that matches the thread the merge/rebase ran on: RefreshFromWorker for the
// worker callers, Refresh for the merge/rebase-continue subprocess path that
// stays on the UI thread.
func (self *MergeAndRebaseHelper) refreshAfterMergeOrRebase(refreshOptions types.RefreshOptions, calledFromWorker bool) {
	if calledFromWorker {
		self.c.RefreshFromWorker(refreshOptions)
	} else {
		self.c.Refresh(refreshOptions)
	}
}

func (self *MergeAndRebaseHelper) CheckMergeOrRebase(result error) error {
	return self.CheckMergeOrRebaseWithRefreshOptions(result, types.RefreshOptions{})
}

// Like CheckMergeOrRebase, but for operations that create a new commit at HEAD
// (a merge, or a pull that merges): on success it selects that new commit,
// which the keep-selection-by-hash logic can't do since the commit didn't exist
// before the refresh.
func (self *MergeAndRebaseHelper) CheckMergeOrRebaseAndSelectHeadCommit(result error) error {
	return self.CheckMergeOrRebaseWithRefreshOptions(result,
		types.RefreshOptions{CommitSelection: commitSelectionAfterMerge(result == nil)})
}

func (self *MergeAndRebaseHelper) CheckForConflicts(result error) error {
	if result == nil {
		return nil
	}

	if isMergeConflictErr(result.Error()) {
		return self.PromptForConflictHandling()
	}

	return result
}

func (self *MergeAndRebaseHelper) PromptForConflictHandling() error {
	operation := self.c.Git().Status.WorkingTreeState().OperationName(self.c.Tr)
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.FoundConflictsTitle,
		Items: []*types.MenuItem{
			{
				Label: self.c.Tr.ViewConflictsMenuItem,
				OnPress: func() error {
					self.c.Context().Push(self.c.Contexts().Files, types.OnFocusOpts{})
					return nil
				},
			},
			{
				Label: utils.ResolvePlaceholderString(self.c.Tr.AbortOperationMenuItem,
					map[string]string{"operation": operation}),
				OnPress: func() error {
					return self.abortMergeOrRebase(false)
				},
				Keys: menuKey('a'),
			},
		},
		HideCancel: true,
	})
}

func (self *MergeAndRebaseHelper) AbortMergeOrRebaseWithConfirm() error {
	return self.abortMergeOrRebase(true)
}

// abortMergeOrRebase aborts the operation in progress. If alwaysConfirm is
// false, it asks the user to confirm only if it has to discard changes before
// git can abort.
func (self *MergeAndRebaseHelper) abortMergeOrRebase(alwaysConfirm bool) error {
	status := self.c.Git().Status.WorkingTreeState()
	placeholders := map[string]string{
		"operation": status.OperationName(self.c.Tr),
	}

	pathsToDiscard := pathsToDiscardBeforeAbort(status, self.c.Model().Files, self.c.Model().Submodules)
	if len(pathsToDiscard) > 0 {
		self.c.Confirm(types.ConfirmOpts{
			Title: utils.ResolvePlaceholderString(self.c.Tr.AbortOperationTitle, placeholders),
			Prompt: utils.ResolvePlaceholderString(self.c.Tr.AbortOperationDiscardsChangesPrompt, placeholders) +
				"\n\n- " + strings.Join(pathsToDiscard, "\n- "),
			HandleConfirm: func() error {
				return self.c.WithWaitingStatus(status.Title(self.c.Tr), func(gocui.Task) error {
					self.c.LogAction(self.c.Tr.Actions.DiscardAllUnstagedChangesInFile)
					if err := self.c.Git().WorkingTree.DiscardUnstagedChanges(pathsToDiscard); err != nil {
						return err
					}
					return self.genericMergeCommandImpl(REBASE_OPTION_ABORT, false, true)
				})
			},
		})
		return nil
	}

	if alwaysConfirm {
		self.c.Confirm(types.ConfirmOpts{
			Title:  utils.ResolvePlaceholderString(self.c.Tr.AbortOperationTitle, placeholders),
			Prompt: utils.ResolvePlaceholderString(self.c.Tr.AbortOperationPrompt, placeholders),
			HandleConfirm: func() error {
				return self.genericMergeCommand(REBASE_OPTION_ABORT)
			},
		})
		return nil
	}

	return self.genericMergeCommand(REBASE_OPTION_ABORT)
}

// pathsToDiscardBeforeAbort returns the files whose unstaged changes we have to
// discard before git can abort the operation in progress.
//
// Git refuses to start a merge, cherry-pick, revert or "git am" if there are
// staged changes, or if a file that it has to change has unstaged changes. It
// does start a squash merge that it can fast-forward while there are staged
// changes, but we don't offer squash merges then. Therefore, if a file has
// staged changes during the operation, git staged them (or the user did after
// git stopped), and any unstaged changes on top of them were made after git
// stopped, for example to fix the build. Git's abort command keeps unstaged
// changes, but it fails if they are on top of staged changes that it has to
// reset. A rebase is different. It refuses to start if there are any changes,
// and its abort command discards all changes.
func pathsToDiscardBeforeAbort(status models.WorkingTreeState, files []*models.File, submoduleConfigs []*models.SubmoduleConfig) []string {
	if status.Effective() == models.WORKING_TREE_STATE_REBASING {
		return nil
	}

	return lo.FilterMap(files, func(file *models.File, _ int) (string, bool) {
		return file.Path, file.HasStagedChanges && file.HasUnstagedChanges && !file.HasMergeConflicts &&
			!file.IsSubmodule(submoduleConfigs)
	})
}

// AbortSquashMergeIfNoChangesAreLeft aborts a squash merge in progress if there
// are no changes left, for example because the user discarded them. Otherwise,
// SQUASH_MSG would keep the squash merge in progress until the next commit,
// and that commit would get its message. Untracked files don't count, because
// a squash merge stages the files that it adds.
func (self *MergeAndRebaseHelper) AbortSquashMergeIfNoChangesAreLeft() error {
	if self.c.Git().Status.WorkingTreeState().Effective() != models.WORKING_TREE_STATE_SQUASH_MERGING {
		return nil
	}

	if lo.SomeBy(self.c.Model().Files, func(file *models.File) bool { return file.Tracked || file.HasStagedChanges }) {
		return nil
	}

	return self.genericMergeCommand(REBASE_OPTION_ABORT)
}

// PromptToContinueRebase asks the user if they want to continue the rebase/merge that's in progress
func (self *MergeAndRebaseHelper) PromptToContinueRebase() {
	self.continueRebasePromptShowing = true
	self.c.Confirm(types.ConfirmOpts{
		Title: self.c.Tr.Continue,
		Prompt: utils.ResolvePlaceholderString(self.c.Tr.ContinueAfterConflictsResolved, map[string]string{
			"operation": self.c.Git().Status.WorkingTreeState().OperationName(self.c.Tr),
		}),
		HandleClose: func() error {
			self.continueRebasePromptShowing = false
			return nil
		},
		HandleConfirm: func() error {
			self.continueRebasePromptShowing = false

			// By the time we get here, we might have unstaged changes again,
			// e.g. if the user had to fix build errors after resolving the
			// conflicts, but after lazygit opened the prompt already. Ask again
			// to auto-stage these.

			// Need to refresh the files to be really sure if this is the case.
			// We would otherwise be relying on lazygit's auto-refresh on focus,
			// but this is not supported by all terminals or on all platforms.
			//
			// The model.Files update is bounced onto the UI thread, so we have
			// to read it in Then; reading it inline here would see the previous
			// model.
			self.c.Refresh(types.RefreshOptions{
				Scope: []types.RefreshableView{types.FILES},
				Then: func() error {
					unstagedFiles := GetUnstagedFilesExceptSubmodules(self.c.Model().Files, self.c.Model().Submodules)
					if len(unstagedFiles) > 0 {
						self.c.Confirm(types.ConfirmOpts{
							Title:  self.c.Tr.Continue,
							Prompt: self.c.Tr.UnstagedFilesAfterConflictsResolved,
							HandleConfirm: func() error {
								self.c.LogAction(self.c.Tr.Actions.StageAllFiles)
								if err := self.c.Git().WorkingTree.StageFiles(unstagedFiles, []string{}); err != nil {
									return err
								}

								return self.genericMergeCommand(REBASE_OPTION_CONTINUE)
							},
						})

						return nil
					}

					return self.genericMergeCommand(REBASE_OPTION_CONTINUE)
				},
			})

			return nil
		},
	})
}

// DismissContinueRebasePromptIfShowing closes the "continue the rebase/merge?"
// prompt if it's currently on screen. It's called when the operation is no
// longer in the state the prompt is offering to act on (e.g. it was continued
// or aborted outside lazygit, or new conflicts have appeared), so that the
// user isn't left with a prompt whose "continue" would now be wrong or fail.
// Must be called on the UI thread.
func (self *MergeAndRebaseHelper) DismissContinueRebasePromptIfShowing() {
	if !self.continueRebasePromptShowing {
		return
	}

	self.continueRebasePromptShowing = false

	// Guard against popping something else: while our prompt is up no other
	// popup can open, and confirming or closing it would have cleared the flag,
	// so if it's set the confirmation context is ours.
	if self.c.Context().Current() == self.c.Contexts().Confirmation {
		self.c.Context().Pop()
	}
}

func (self *MergeAndRebaseHelper) RebaseOntoRef(ref string) error {
	checkedOutBranch := self.c.Model().Branches[0]
	checkedOutBranchName := checkedOutBranch.Name
	var disabledReason, baseBranchDisabledReason *types.DisabledReason
	if checkedOutBranchName == ref {
		disabledReason = &types.DisabledReason{Text: self.c.Tr.CantRebaseOntoSelf}
	}

	baseBranch, err := self.c.Git().Loaders.BranchLoader.GetBaseBranch(checkedOutBranch, self.c.Model().MainBranches)
	if err != nil {
		return err
	}
	if baseBranch == "" {
		baseBranch = self.c.Tr.CouldNotDetermineBaseBranch
		baseBranchDisabledReason = &types.DisabledReason{Text: self.c.Tr.CouldNotDetermineBaseBranch}
	}

	menuItems := []*types.MenuItem{
		{
			Label: utils.ResolvePlaceholderString(self.c.Tr.SimpleRebase,
				map[string]string{"ref": ref},
			),
			Keys:           menuKey('s'),
			DisabledReason: disabledReason,
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.RebaseBranch)
				baseCommit := self.c.Modes().MarkedBaseCommit.GetHash()
				return self.c.WithWaitingStatus(self.c.Tr.RebasingStatus, func(task gocui.Task) error {
					var err error
					if baseCommit != "" {
						err = self.c.Git().Rebase.RebaseBranchFromBaseCommit(ref, baseCommit)
					} else {
						err = self.c.Git().Rebase.RebaseBranch(ref)
					}
					err = self.CheckMergeOrRebase(err)
					if err == nil {
						self.c.OnUIThread(func() error {
							return self.ResetMarkedBaseCommit()
						})
					}
					return err
				})
			},
		},
		{
			Label: utils.ResolvePlaceholderString(self.c.Tr.InteractiveRebase,
				map[string]string{"ref": ref},
			),
			Keys:           menuKey('i'),
			DisabledReason: disabledReason,
			Tooltip:        self.c.Tr.InteractiveRebaseTooltip,
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.RebaseBranch)
				baseCommit := self.c.Modes().MarkedBaseCommit.GetHash()
				return self.c.WithWaitingStatus(self.c.Tr.RebasingStatus, func(task gocui.Task) error {
					var err error
					if baseCommit != "" {
						err = self.c.Git().Rebase.EditRebaseFromBaseCommit(ref, baseCommit)
					} else {
						err = self.c.Git().Rebase.EditRebase(ref)
					}
					if err = self.CheckMergeOrRebase(err); err != nil {
						return err
					}
					self.c.OnUIThread(func() error {
						if err := self.ResetMarkedBaseCommit(); err != nil {
							return err
						}
						self.c.Context().Push(self.c.Contexts().LocalCommits, types.OnFocusOpts{})
						return nil
					})
					return nil
				})
			},
		},
		{
			Label: utils.ResolvePlaceholderString(self.c.Tr.RebaseOntoBaseBranch,
				map[string]string{"baseBranch": ShortBranchName(baseBranch)},
			),
			Keys:           menuKey('b'),
			DisabledReason: baseBranchDisabledReason,
			Tooltip:        self.c.Tr.RebaseOntoBaseBranchTooltip,
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.RebaseBranch)
				baseCommit := self.c.Modes().MarkedBaseCommit.GetHash()
				return self.c.WithWaitingStatus(self.c.Tr.RebasingStatus, func(task gocui.Task) error {
					var err error
					if baseCommit != "" {
						err = self.c.Git().Rebase.RebaseBranchFromBaseCommit(baseBranch, baseCommit)
					} else {
						err = self.c.Git().Rebase.RebaseBranch(baseBranch)
					}
					err = self.CheckMergeOrRebase(err)
					if err == nil {
						self.c.OnUIThread(func() error {
							return self.ResetMarkedBaseCommit()
						})
					}
					return err
				})
			},
		},
	}

	title := utils.ResolvePlaceholderString(
		lo.Ternary(self.c.Modes().MarkedBaseCommit.GetHash() != "",
			self.c.Tr.RebasingFromBaseCommitTitle,
			self.c.Tr.RebasingTitle),
		map[string]string{
			"checkedOutBranch": checkedOutBranchName,
		},
	)

	return self.c.Menu(types.CreateMenuOptions{
		Title: title,
		Items: menuItems,
	})
}

func (self *MergeAndRebaseHelper) MergeRefIntoCheckedOutBranch(refName string) error {
	if self.c.Git().Branch.IsHeadDetached() {
		return errors.New("Cannot merge branch in detached head state. You might have checked out a commit directly or a remote branch, in which case you should checkout the local branch you want to be on")
	}
	checkedOutBranchName := self.c.Model().Branches[0].Name
	if checkedOutBranchName == refName {
		return errors.New(self.c.Tr.CantMergeBranchIntoItself)
	}

	wantFastForward, wantNonFastForward := self.fastForwardMergeUserPreference()
	canFastForward := self.c.Git().Branch.IsAncestor("HEAD", refName)

	var firstRegularMergeItem *types.MenuItem
	var secondRegularMergeItem *types.MenuItem
	var fastForwardMergeItem *types.MenuItem

	if !wantNonFastForward && (wantFastForward || canFastForward) {
		firstRegularMergeItem = &types.MenuItem{
			Label:   self.c.Tr.RegularMergeFastForward,
			OnPress: self.RegularMerge(refName, git_commands.MERGE_VARIANT_REGULAR),
			Keys:    menuKey('m'),
			Tooltip: utils.ResolvePlaceholderString(
				self.c.Tr.RegularMergeFastForwardTooltip,
				map[string]string{
					"checkedOutBranch": checkedOutBranchName,
					"selectedBranch":   refName,
				},
			),
		}
		fastForwardMergeItem = firstRegularMergeItem

		secondRegularMergeItem = &types.MenuItem{
			Label:   self.c.Tr.RegularMergeNonFastForward,
			OnPress: self.RegularMerge(refName, git_commands.MERGE_VARIANT_NON_FAST_FORWARD),
			Keys:    menuKey('n'),
			Tooltip: utils.ResolvePlaceholderString(
				self.c.Tr.RegularMergeNonFastForwardTooltip,
				map[string]string{
					"checkedOutBranch": checkedOutBranchName,
					"selectedBranch":   refName,
				},
			),
		}
	} else {
		firstRegularMergeItem = &types.MenuItem{
			Label:   self.c.Tr.RegularMergeNonFastForward,
			OnPress: self.RegularMerge(refName, git_commands.MERGE_VARIANT_REGULAR),
			Keys:    menuKey('m'),
			Tooltip: utils.ResolvePlaceholderString(
				self.c.Tr.RegularMergeNonFastForwardTooltip,
				map[string]string{
					"checkedOutBranch": checkedOutBranchName,
					"selectedBranch":   refName,
				},
			),
		}

		secondRegularMergeItem = &types.MenuItem{
			Label:   self.c.Tr.RegularMergeFastForward,
			OnPress: self.RegularMerge(refName, git_commands.MERGE_VARIANT_FAST_FORWARD),
			Keys:    menuKey('f'),
			Tooltip: utils.ResolvePlaceholderString(
				self.c.Tr.RegularMergeFastForwardTooltip,
				map[string]string{
					"checkedOutBranch": checkedOutBranchName,
					"selectedBranch":   refName,
				},
			),
		}
		fastForwardMergeItem = secondRegularMergeItem
	}

	if !canFastForward {
		fastForwardMergeItem.DisabledReason = &types.DisabledReason{
			Text: utils.ResolvePlaceholderString(
				self.c.Tr.CannotFastForwardMerge,
				map[string]string{
					"checkedOutBranch": checkedOutBranchName,
					"selectedBranch":   refName,
				},
			),
		}
	}

	// Git refuses to merge while there are staged changes, but not if it can
	// fast-forward the squash merge. It would add the changes of the branch to
	// the staged ones then, and we couldn't tell them apart anymore.
	var squashMergeDisabledReason *types.DisabledReason
	if AnyStagedFiles(self.c.Model().Files) {
		squashMergeDisabledReason = &types.DisabledReason{Text: self.c.Tr.CantSquashMergeWithStagedChanges}
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.Merge,
		Items: []*types.MenuItem{
			firstRegularMergeItem,
			secondRegularMergeItem,
			{
				Label:          self.c.Tr.SquashMergeUncommitted,
				OnPress:        self.SquashMergeUncommitted(refName, checkedOutBranchName),
				DisabledReason: squashMergeDisabledReason,
				Keys:           menuKey('s'),
				Tooltip: utils.ResolvePlaceholderString(
					self.c.Tr.SquashMergeUncommittedTooltip,
					map[string]string{
						"selectedBranch": refName,
					},
				),
			},
			{
				Label:          self.c.Tr.SquashMergeCommitted,
				OnPress:        self.SquashMergeCommitted(refName, checkedOutBranchName),
				DisabledReason: squashMergeDisabledReason,
				Keys:           menuKey('S'),
				Tooltip: utils.ResolvePlaceholderString(
					self.c.Tr.SquashMergeCommittedTooltip,
					map[string]string{
						"checkedOutBranch": checkedOutBranchName,
						"selectedBranch":   refName,
					},
				),
			},
		},
	})
}

func (self *MergeAndRebaseHelper) RegularMerge(refName string, variant git_commands.MergeVariant) func() error {
	return func() error {
		self.c.LogAction(self.c.Tr.Actions.Merge)
		return self.c.WithWaitingStatus(self.c.Tr.MergingStatus, func(gocui.Task) error {
			err := self.c.Git().Branch.Merge(refName, variant)
			return self.CheckMergeOrRebaseAndSelectHeadCommit(err)
		})
	}
}

func (self *MergeAndRebaseHelper) SquashMergeUncommitted(refName, checkedOutBranchName string) func() error {
	return func() error {
		self.c.LogAction(self.c.Tr.Actions.SquashMerge)
		return self.c.WithWaitingStatus(self.c.Tr.MergingStatus, func(gocui.Task) error {
			err := self.c.Git().Branch.SquashMerge(refName, self.squashMergeMessage(refName, checkedOutBranchName))
			err = self.CheckMergeOrRebase(err)

			// The user asked us not to commit the changes, so don't offer to
			// once they have resolved the conflicts
			self.c.State().GetRepoState().SetMergeOrRebaseStartedInLazygit(false)

			return err
		})
	}
}

func (self *MergeAndRebaseHelper) SquashMergeCommitted(refName, checkedOutBranchName string) func() error {
	return func() error {
		self.c.LogAction(self.c.Tr.Actions.SquashMerge)
		return self.c.WithWaitingStatus(self.c.Tr.MergingStatus, func(gocui.Task) error {
			err := self.c.Git().Branch.SquashMerge(refName, self.squashMergeMessage(refName, checkedOutBranchName))
			if err == nil {
				err = self.c.Git().Rebase.GenericMergeOrRebaseAction(models.WORKING_TREE_STATE_SQUASH_MERGING, REBASE_OPTION_CONTINUE)
			}
			return self.CheckMergeOrRebaseAndSelectHeadCommit(err)
		})
	}
}

func (self *MergeAndRebaseHelper) squashMergeMessage(refName, checkedOutBranchName string) string {
	return utils.ResolvePlaceholderString(self.c.UserConfig().Git.Merging.SquashMergeMessage, map[string]string{
		"selectedRef":   refName,
		"currentBranch": checkedOutBranchName,
	})
}

// Returns wantsFastForward, wantsNonFastForward. These will never both be true, but they can both be false.
func (self *MergeAndRebaseHelper) fastForwardMergeUserPreference() (bool, bool) {
	// Check user config first, because it takes precedence over git config
	mergingArgs := self.c.UserConfig().Git.Merging.Args
	if strings.Contains(mergingArgs, "--ff") { // also covers "--ff-only"
		return true, false
	}

	if strings.Contains(mergingArgs, "--no-ff") {
		return false, true
	}

	// Then check git config
	mergeFfConfig := self.c.Git().Config.GetMergeFF()
	if mergeFfConfig == "true" || mergeFfConfig == "only" {
		return true, false
	}

	if mergeFfConfig == "false" {
		return false, true
	}

	return false, false
}

func (self *MergeAndRebaseHelper) ResetMarkedBaseCommit() error {
	self.c.Modes().MarkedBaseCommit.Reset()
	self.c.PostRefreshUpdate(self.c.Contexts().LocalCommits)
	return nil
}
