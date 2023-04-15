package helpers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type MergeAndRebaseHelper struct {
	c          *types.HelperCommon
	contexts   *context.ContextTree
	git        *commands.GitCommand
	refsHelper *RefsHelper
}

func NewMergeAndRebaseHelper(
	c *types.HelperCommon,
	contexts *context.ContextTree,
	git *commands.GitCommand,
	refsHelper *RefsHelper,
) *MergeAndRebaseHelper {
	return &MergeAndRebaseHelper{
		c:          c,
		contexts:   contexts,
		git:        git,
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

	if self.git.Status.WorkingTreeState() == enums.REBASE_MODE_REBASING {
		options = append(options, optionAndKey{
			option: REBASE_OPTION_SKIP, key: 's',
		})
	}

	menuItems := slices.Map(options, func(row optionAndKey) *types.MenuItem {
		return &types.MenuItem{
			Label: row.option,
			OnPress: func() error {
				return self.genericMergeCommand(row.option)
			},
			Key: row.key,
		}
	})

	var title string
	if self.git.Status.WorkingTreeState() == enums.REBASE_MODE_MERGING {
		title = self.c.Tr.MergeOptionsTitle
	} else {
		title = self.c.Tr.RebaseOptionsTitle
	}

	return self.c.Menu(types.CreateMenuOptions{Title: title, Items: menuItems})
}

func (self *MergeAndRebaseHelper) genericMergeCommand(command string) error {
	status := self.git.Status.WorkingTreeState()

	if status != enums.REBASE_MODE_MERGING && status != enums.REBASE_MODE_REBASING {
		return self.c.ErrorMsg(self.c.Tr.NotMergingOrRebasing)
	}

	self.c.LogAction(fmt.Sprintf("Merge/Rebase: %s", command))

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
			self.git.Rebase.GenericMergeOrRebaseActionCmdObj(commandType, command),
		)
	}
	result := self.git.Rebase.GenericMergeOrRebaseAction(commandType, command)
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

func (self *MergeAndRebaseHelper) CheckMergeOrRebase(result error) error {
	if err := self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC}); err != nil {
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
	} else if isMergeConflictErr(result.Error()) {
		return self.c.Confirm(types.ConfirmOpts{
			Title:  self.c.Tr.FoundConflictsTitle,
			Prompt: self.c.Tr.FoundConflicts,
			HandleConfirm: func() error {
				return self.c.PushContext(self.contexts.Files)
			},
			HandleClose: func() error {
				return self.genericMergeCommand(REBASE_OPTION_ABORT)
			},
		})
	} else {
		return self.c.ErrorMsg(result.Error())
	}
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
	workingTreeState := self.git.Status.WorkingTreeState()
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
		Title:  "continue",
		Prompt: self.c.Tr.ConflictsResolved,
		HandleConfirm: func() error {
			return self.genericMergeCommand(REBASE_OPTION_CONTINUE)
		},
	})
}

func (self *MergeAndRebaseHelper) RebaseOntoRef(ref string) error {
	checkedOutBranch := self.refsHelper.GetCheckedOutRef().Name
	if ref == checkedOutBranch {
		return self.c.ErrorMsg(self.c.Tr.CantRebaseOntoSelf)
	}
	menuItems := []*types.MenuItem{
		{
			Label: self.c.Tr.SimpleRebase,
			Key:   's',
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.RebaseBranch)
				err := self.git.Rebase.RebaseBranch(ref)
				return self.CheckMergeOrRebase(err)
			},
		},
		{
			Label:   self.c.Tr.InteractiveRebase,
			Key:     'i',
			Tooltip: self.c.Tr.InteractiveRebaseTooltip,
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.RebaseBranch)
				err := self.git.Rebase.EditRebase(ref)
				if err = self.CheckMergeOrRebase(err); err != nil {
					return err
				}
				return self.c.PushContext(self.contexts.LocalCommits)
			},
		},
	}

	title := utils.ResolvePlaceholderString(
		self.c.Tr.RebasingTitle,
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
	if self.git.Branch.IsHeadDetached() {
		return self.c.ErrorMsg("Cannot merge branch in detached head state. You might have checked out a commit directly or a remote branch, in which case you should checkout the local branch you want to be on")
	}
	checkedOutBranchName := self.refsHelper.GetCheckedOutRef().Name
	if checkedOutBranchName == refName {
		return self.c.ErrorMsg(self.c.Tr.CantMergeBranchIntoItself)
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
			err := self.git.Branch.Merge(refName, git_commands.MergeOpts{})
			return self.CheckMergeOrRebase(err)
		},
	})
}
