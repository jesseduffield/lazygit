package helpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type MergeAndRebaseHelper struct {
	c          *HelperCommon
	refsHelper *RefsHelper
}

func NewMergeAndRebaseHelper(
	c *HelperCommon,
	refsHelper *RefsHelper,
) *MergeAndRebaseHelper {
	return &MergeAndRebaseHelper{
		c:          c,
		refsHelper: refsHelper,
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

	if self.c.Git().Status.WorkingTreeState() == enums.REBASE_MODE_REBASING {
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

	var title string
	if self.c.Git().Status.WorkingTreeState() == enums.REBASE_MODE_MERGING {
		title = self.c.Tr.MergeOptionsTitle
	} else {
		title = self.c.Tr.RebaseOptionsTitle
	}

	return self.c.Menu(types.CreateMenuOptions{Title: title, Items: menuItems})
}

func (self *MergeAndRebaseHelper) genericMergeCommand(command string) error {
	status := self.c.Git().Status.WorkingTreeState()

	if status != enums.REBASE_MODE_MERGING && status != enums.REBASE_MODE_REBASING {
		return errors.New(self.c.Tr.NotMergingOrRebasing)
	}

	self.c.LogAction(fmt.Sprintf("Merge/Rebase: %s", command))
	if status == enums.REBASE_MODE_REBASING {
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

	commandType := ""
	switch status {
	case enums.REBASE_MODE_MERGING:
		commandType = "merge"
	case enums.REBASE_MODE_REBASING:
		commandType = "rebase"
	default:
		// shouldn't be possible to land here
	}

	// we should end up with a command like 'git merge --continue'

	// it's impossible for a rebase to require a commit so we'll use a subprocess only if it's a merge
	if status == enums.REBASE_MODE_MERGING && command != REBASE_OPTION_ABORT && self.c.UserConfig.Git.Merging.ManualCommit {
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

var conflictStrings = []string{
	"Failed to merge in the changes",
	"When you have resolved this problem",
	"fix conflicts",
	"Resolve all conflicts manually",
	"Merge conflict in file",
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
	if err := self.c.Refresh(refreshOptions); err != nil {
		return err
	}
	if result == nil {
		return nil
	} else if strings.Contains(result.Error(), "No changes - did you forget to use") {
		return self.genericMergeCommand(REBASE_OPTION_SKIP)
	} else if strings.Contains(result.Error(), "The previous cherry-pick is now empty") {
		return self.genericMergeCommand(REBASE_OPTION_CONTINUE)
	} else if strings.Contains(result.Error(), "No rebase in progress?") {
		// assume in this case that we're already done
		return nil
	} else {
		return self.CheckForConflicts(result)
	}
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
	mode := self.workingTreeStateNoun()
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.FoundConflictsTitle,
		Items: []*types.MenuItem{
			{
				Label: self.c.Tr.ViewConflictsMenuItem,
				OnPress: func() error {
					return self.c.PushContext(self.c.Contexts().Files)
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
	mode := self.workingTreeStateNoun()
	return self.c.Confirm(types.ConfirmOpts{
		Title:  fmt.Sprintf(self.c.Tr.AbortTitle, mode),
		Prompt: fmt.Sprintf(self.c.Tr.AbortPrompt, mode),
		HandleConfirm: func() error {
			return self.genericMergeCommand(REBASE_OPTION_ABORT)
		},
	})
}

func (self *MergeAndRebaseHelper) workingTreeStateNoun() string {
	workingTreeState := self.c.Git().Status.WorkingTreeState()
	switch workingTreeState {
	case enums.REBASE_MODE_NONE:
		return ""
	case enums.REBASE_MODE_MERGING:
		return "merge"
	default:
		return "rebase"
	}
}

// PromptToContinueRebase asks the user if they want to continue the rebase/merge that's in progress
func (self *MergeAndRebaseHelper) PromptToContinueRebase() error {
	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.Continue,
		Prompt: self.c.Tr.ConflictsResolved,
		HandleConfirm: func() error {
			return self.genericMergeCommand(REBASE_OPTION_CONTINUE)
		},
	})
}

func (self *MergeAndRebaseHelper) RebaseOntoRef(ref string) error {
	checkedOutBranch := self.refsHelper.GetCheckedOutRef().Name
	menuItems := []*types.MenuItem{
		{
			Label: self.c.Tr.SimpleRebase,
			Key:   's',
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
			Label:   self.c.Tr.InteractiveRebase,
			Key:     'i',
			Tooltip: self.c.Tr.InteractiveRebaseTooltip,
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
				return self.c.PushContext(self.c.Contexts().LocalCommits)
			},
		},
	}

	title := utils.ResolvePlaceholderString(
		lo.Ternary(self.c.Modes().MarkedBaseCommit.GetHash() != "",
			self.c.Tr.RebasingFromBaseCommitTitle,
			self.c.Tr.RebasingTitle),
		map[string]string{
			"checkedOutBranch": checkedOutBranch,
			"ref":              ref,
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
	checkedOutBranchName := self.refsHelper.GetCheckedOutRef().Name
	if checkedOutBranchName == refName {
		return errors.New(self.c.Tr.CantMergeBranchIntoItself)
	}
	prompt := utils.ResolvePlaceholderString(
		self.c.Tr.ConfirmMerge,
		map[string]string{
			"checkedOutBranch": checkedOutBranchName,
			"selectedBranch":   refName,
		},
	)

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.MergeConfirmTitle,
		Prompt: prompt,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.Merge)
			err := self.c.Git().Branch.Merge(refName, git_commands.MergeOpts{})
			return self.CheckMergeOrRebase(err)
		},
	})
}

func (self *MergeAndRebaseHelper) ResetMarkedBaseCommit() error {
	self.c.Modes().MarkedBaseCommit.Reset()
	return self.c.PostRefreshUpdate(self.c.Contexts().LocalCommits)
}
