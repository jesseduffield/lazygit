package controllers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type RebaseHelper struct {
	c                              *types.ControllerCommon
	contexts                       *context.ContextTree
	git                            *commands.GitCommand
	takeOverMergeConflictScrolling func()
}

func NewRebaseHelper(
	c *types.ControllerCommon,
	contexts *context.ContextTree,
	git *commands.GitCommand,
	takeOverMergeConflictScrolling func(),
) *RebaseHelper {
	return &RebaseHelper{
		c:                              c,
		contexts:                       contexts,
		git:                            git,
		takeOverMergeConflictScrolling: takeOverMergeConflictScrolling,
	}
}

type RebaseOption string

const (
	REBASE_OPTION_CONTINUE string = "continue"
	REBASE_OPTION_ABORT    string = "abort"
	REBASE_OPTION_SKIP     string = "skip"
)

func (self *RebaseHelper) CreateRebaseOptionsMenu() error {
	options := []string{REBASE_OPTION_CONTINUE, REBASE_OPTION_ABORT}

	if self.git.Status.WorkingTreeState() == enums.REBASE_MODE_REBASING {
		options = append(options, REBASE_OPTION_SKIP)
	}

	menuItems := make([]*types.MenuItem, len(options))
	for i, option := range options {
		// note to self. Never, EVER, close over loop variables in a function
		option := option
		menuItems[i] = &types.MenuItem{
			DisplayString: option,
			OnPress: func() error {
				return self.genericMergeCommand(option)
			},
		}
	}

	var title string
	if self.git.Status.WorkingTreeState() == enums.REBASE_MODE_MERGING {
		title = self.c.Tr.MergeOptionsTitle
	} else {
		title = self.c.Tr.RebaseOptionsTitle
	}

	return self.c.Menu(types.CreateMenuOptions{Title: title, Items: menuItems})
}

func (self *RebaseHelper) genericMergeCommand(command string) error {
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
}

func isMergeConflictErr(errStr string) bool {
	for _, str := range conflictStrings {
		if strings.Contains(errStr, str) {
			return true
		}
	}

	return false
}

func (self *RebaseHelper) CheckMergeOrRebase(result error) error {
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
		return self.c.Ask(types.AskOpts{
			Title:               self.c.Tr.FoundConflictsTitle,
			Prompt:              self.c.Tr.FoundConflicts,
			HandlersManageFocus: true,
			HandleConfirm: func() error {
				return self.c.PushContext(self.contexts.Files)
			},
			HandleClose: func() error {
				if err := self.c.PopContext(); err != nil {
					return err
				}

				return self.genericMergeCommand(REBASE_OPTION_ABORT)
			},
		})
	} else {
		return self.c.ErrorMsg(result.Error())
	}
}

func (self *RebaseHelper) AbortMergeOrRebaseWithConfirm() error {
	// prompt user to confirm that they want to abort, then do it
	mode := self.workingTreeStateNoun()
	return self.c.Ask(types.AskOpts{
		Title:  fmt.Sprintf(self.c.Tr.AbortTitle, mode),
		Prompt: fmt.Sprintf(self.c.Tr.AbortPrompt, mode),
		HandleConfirm: func() error {
			return self.genericMergeCommand(REBASE_OPTION_ABORT)
		},
	})
}

func (self *RebaseHelper) workingTreeStateNoun() string {
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
func (self *RebaseHelper) PromptToContinueRebase() error {
	self.takeOverMergeConflictScrolling()

	return self.c.Ask(types.AskOpts{
		Title:  "continue",
		Prompt: self.c.Tr.ConflictsResolved,
		HandleConfirm: func() error {
			return self.genericMergeCommand(REBASE_OPTION_CONTINUE)
		},
	})
}
