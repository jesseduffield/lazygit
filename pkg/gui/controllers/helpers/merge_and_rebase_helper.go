package helpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/stefanhaller/git-todo-parser/todo"
)

type MergeAndRebaseHelper struct {
	c *HelperCommon
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
		key    types.Key
	}

	options := []optionAndKey{
		{option: REBASE_OPTION_CONTINUE, key: 'c'},
		{option: REBASE_OPTION_ABORT, key: 'a'},
	}

	if self.c.Git().Status.WorkingTreeState().CanSkip() {
		options = append(options, optionAndKey{
			option: REBASE_OPTION_SKIP, key: 's',
		})
	}

	menuItems := lo.Map(options, func(row optionAndKey, _ int) *types.MenuItem {
		return &types.MenuItem{
			Label: row.option,
			OnPress: func() error {
				return self.genericMergeCommand(row.option)
			},
			Key: row.key,
		}
	})

	title := self.c.Git().Status.WorkingTreeState().OptionsMenuTitle(self.c.Tr)
	return self.c.Menu(types.CreateMenuOptions{Title: title, Items: menuItems})
}

func (self *MergeAndRebaseHelper) ContinueRebase() error {
	return self.genericMergeCommand(REBASE_OPTION_CONTINUE)
}

func (self *MergeAndRebaseHelper) genericMergeCommand(command string) error {
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

	commandType := status.CommandName()

	// we should end up with a command like 'git merge --continue'

	// it's impossible for a rebase to require a commit so we'll use a subprocess only if it's a merge
	needsSubprocess := (effectiveStatus == models.WORKING_TREE_STATE_MERGING && command != REBASE_OPTION_ABORT && self.c.UserConfig().Git.Merging.ManualCommit) ||
		// but we'll also use a subprocess if we have exec todos; those are likely to be lengthy build
		// tasks whose output the user will want to see in the terminal
		(effectiveStatus == models.WORKING_TREE_STATE_REBASING && command != REBASE_OPTION_ABORT && self.hasExecTodos())

	if needsSubprocess {
		// TODO: see if we should be calling more of the code from self.Git.Rebase.GenericMergeOrRebaseAction
		return self.c.RunSubprocessAndRefresh(
			self.c.Git().Rebase.GenericMergeOrRebaseActionCmdObj(commandType, command),
		)
	}
	result := self.c.Git().Rebase.GenericMergeOrRebaseAction(commandType, command)
	if err := self.CheckMergeOrRebase(result); err != nil {
		return err
	}
	return nil
}

func (self *MergeAndRebaseHelper) hasExecTodos() bool {
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

func (self *MergeAndRebaseHelper) CheckMergeOrRebaseWithRefreshOptions(result error, refreshOptions types.RefreshOptions) error {
	self.c.Refresh(refreshOptions)

	if result == nil {
		return nil
	} else if strings.Contains(result.Error(), "No changes - did you forget to use") {
		return self.genericMergeCommand(REBASE_OPTION_SKIP)
	} else if strings.Contains(result.Error(), "The previous cherry-pick is now empty") {
		return self.genericMergeCommand(REBASE_OPTION_SKIP)
	} else if strings.Contains(result.Error(), "No rebase in progress?") {
		// assume in this case that we're already done
		return nil
	}
	return self.CheckForConflicts(result)
}

func (self *MergeAndRebaseHelper) CheckMergeOrRebase(result error) error {
	return self.CheckMergeOrRebaseWithRefreshOptions(result, types.RefreshOptions{Mode: types.ASYNC})
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
	mode := self.c.Git().Status.WorkingTreeState().CommandName()
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
				Label: fmt.Sprintf(self.c.Tr.AbortMenuItem, mode),
				OnPress: func() error {
					return self.genericMergeCommand(REBASE_OPTION_ABORT)
				},
				Key: 'a',
			},
		},
		HideCancel: true,
	})
}

func (self *MergeAndRebaseHelper) AbortMergeOrRebaseWithConfirm() error {
	// prompt user to confirm that they want to abort, then do it
	mode := self.c.Git().Status.WorkingTreeState().CommandName()
	self.c.Confirm(types.ConfirmOpts{
		Title:  fmt.Sprintf(self.c.Tr.AbortTitle, mode),
		Prompt: fmt.Sprintf(self.c.Tr.AbortPrompt, mode),
		HandleConfirm: func() error {
			return self.genericMergeCommand(REBASE_OPTION_ABORT)
		},
	})

	return nil
}

// PromptToContinueRebase asks the user if they want to continue the rebase/merge that's in progress
func (self *MergeAndRebaseHelper) PromptToContinueRebase() error {
	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.Continue,
		Prompt: fmt.Sprintf(self.c.Tr.ConflictsResolved, self.c.Git().Status.WorkingTreeState().CommandName()),
		HandleConfirm: func() error {
			// By the time we get here, we might have unstaged changes again,
			// e.g. if the user had to fix build errors after resolving the
			// conflicts, but after lazygit opened the prompt already. Ask again
			// to auto-stage these.

			// Need to refresh the files to be really sure if this is the case.
			// We would otherwise be relying on lazygit's auto-refresh on focus,
			// but this is not supported by all terminals or on all platforms.
			self.c.Refresh(types.RefreshOptions{
				Mode: types.SYNC, Scope: []types.RefreshableView{types.FILES},
			})

			root := self.c.Contexts().Files.FileTreeViewModel.GetRoot()
			if root.GetHasUnstagedChanges() {
				self.c.Confirm(types.ConfirmOpts{
					Title:  self.c.Tr.Continue,
					Prompt: self.c.Tr.UnstagedFilesAfterConflictsResolved,
					HandleConfirm: func() error {
						self.c.LogAction(self.c.Tr.Actions.StageAllFiles)
						if err := self.c.Git().WorkingTree.StageAll(true); err != nil {
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
			Key:            's',
			DisabledReason: disabledReason,
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.RebaseBranch)
				return self.c.WithWaitingStatus(self.c.Tr.RebasingStatus, func(task gocui.Task) error {
					baseCommit := self.c.Modes().MarkedBaseCommit.GetHash()
					var err error
					if baseCommit != "" {
						err = self.c.Git().Rebase.RebaseBranchFromBaseCommit(ref, baseCommit)
					} else {
						err = self.c.Git().Rebase.RebaseBranch(ref)
					}
					err = self.CheckMergeOrRebase(err)
					if err == nil {
						return self.ResetMarkedBaseCommit()
					}
					return err
				})
			},
		},
		{
			Label: utils.ResolvePlaceholderString(self.c.Tr.InteractiveRebase,
				map[string]string{"ref": ref},
			),
			Key:            'i',
			DisabledReason: disabledReason,
			Tooltip:        self.c.Tr.InteractiveRebaseTooltip,
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.RebaseBranch)
				baseCommit := self.c.Modes().MarkedBaseCommit.GetHash()
				var err error
				if baseCommit != "" {
					err = self.c.Git().Rebase.EditRebaseFromBaseCommit(ref, baseCommit)
				} else {
					err = self.c.Git().Rebase.EditRebase(ref)
				}
				if err = self.CheckMergeOrRebase(err); err != nil {
					return err
				}
				if err = self.ResetMarkedBaseCommit(); err != nil {
					return err
				}
				self.c.Context().Push(self.c.Contexts().LocalCommits, types.OnFocusOpts{})
				return nil
			},
		},
		{
			Label: utils.ResolvePlaceholderString(self.c.Tr.RebaseOntoBaseBranch,
				map[string]string{"baseBranch": ShortBranchName(baseBranch)},
			),
			Key:            'b',
			DisabledReason: baseBranchDisabledReason,
			Tooltip:        self.c.Tr.RebaseOntoBaseBranchTooltip,
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.RebaseBranch)
				return self.c.WithWaitingStatus(self.c.Tr.RebasingStatus, func(task gocui.Task) error {
					baseCommit := self.c.Modes().MarkedBaseCommit.GetHash()
					var err error
					if baseCommit != "" {
						err = self.c.Git().Rebase.RebaseBranchFromBaseCommit(baseBranch, baseCommit)
					} else {
						err = self.c.Git().Rebase.RebaseBranch(baseBranch)
					}
					err = self.CheckMergeOrRebase(err)
					if err == nil {
						return self.ResetMarkedBaseCommit()
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
	canFastForward := self.c.Git().Branch.CanDoFastForwardMerge(refName)

	var firstRegularMergeItem *types.MenuItem
	var secondRegularMergeItem *types.MenuItem
	var fastForwardMergeItem *types.MenuItem

	if !wantNonFastForward && (wantFastForward || canFastForward) {
		firstRegularMergeItem = &types.MenuItem{
			Label:   self.c.Tr.RegularMergeFastForward,
			OnPress: self.RegularMerge(refName, git_commands.MERGE_VARIANT_REGULAR),
			Key:     'm',
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
			Key:     'n',
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
			Key:     'm',
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
			Key:     'f',
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

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.Merge,
		Items: []*types.MenuItem{
			firstRegularMergeItem,
			secondRegularMergeItem,
			{
				Label:   self.c.Tr.SquashMergeUncommitted,
				OnPress: self.SquashMergeUncommitted(refName),
				Key:     's',
				Tooltip: utils.ResolvePlaceholderString(
					self.c.Tr.SquashMergeUncommittedTooltip,
					map[string]string{
						"selectedBranch": refName,
					},
				),
			},
			{
				Label:   self.c.Tr.SquashMergeCommitted,
				OnPress: self.SquashMergeCommitted(refName, checkedOutBranchName),
				Key:     'S',
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
		err := self.c.Git().Branch.Merge(refName, variant)
		return self.CheckMergeOrRebase(err)
	}
}

func (self *MergeAndRebaseHelper) SquashMergeUncommitted(refName string) func() error {
	return func() error {
		self.c.LogAction(self.c.Tr.Actions.SquashMerge)
		err := self.c.Git().Branch.Merge(refName, git_commands.MERGE_VARIANT_SQUASH)
		return self.CheckMergeOrRebase(err)
	}
}

func (self *MergeAndRebaseHelper) SquashMergeCommitted(refName, checkedOutBranchName string) func() error {
	return func() error {
		self.c.LogAction(self.c.Tr.Actions.SquashMerge)
		err := self.c.Git().Branch.Merge(refName, git_commands.MERGE_VARIANT_SQUASH)
		if err = self.CheckMergeOrRebase(err); err != nil {
			return err
		}
		message := utils.ResolvePlaceholderString(self.c.UserConfig().Git.Merging.SquashMergeMessage, map[string]string{
			"selectedRef":   refName,
			"currentBranch": checkedOutBranchName,
		})
		err = self.c.Git().Commit.CommitCmdObj(message, "", false).Run()
		if err != nil {
			return err
		}
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		return nil
	}
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
