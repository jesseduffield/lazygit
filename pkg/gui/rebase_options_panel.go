package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands"
)

type RebaseOption string

const (
	REBASE_OPTION_CONTINUE = "continue"
	REBASE_OPTION_ABORT    = "abort"
	REBASE_OPTION_SKIP     = "skip"
)

func (gui *Gui) handleCreateRebaseOptionsMenu() error {
	options := []string{REBASE_OPTION_CONTINUE, REBASE_OPTION_ABORT}

	if gui.GitCommand.WorkingTreeState() == commands.REBASE_MODE_REBASING {
		options = append(options, REBASE_OPTION_SKIP)
	}

	menuItems := make([]*menuItem, len(options))
	for i, option := range options {
		// note to self. Never, EVER, close over loop variables in a function
		option := option
		menuItems[i] = &menuItem{
			displayString: option,
			onPress: func() error {
				return gui.genericMergeCommand(option)
			},
		}
	}

	var title string
	if gui.GitCommand.WorkingTreeState() == commands.REBASE_MODE_MERGING {
		title = gui.Tr.MergeOptionsTitle
	} else {
		title = gui.Tr.RebaseOptionsTitle
	}

	return gui.createMenu(title, menuItems, createMenuOptions{showCancel: true})
}

func (gui *Gui) genericMergeCommand(command string) error {
	status := gui.GitCommand.WorkingTreeState()

	if status != commands.REBASE_MODE_MERGING && status != commands.REBASE_MODE_REBASING {
		return gui.createErrorPanel(gui.Tr.NotMergingOrRebasing)
	}

	gitCommand := gui.GitCommand.WithSpan(fmt.Sprintf("Merge/Rebase: %s", command))

	commandType := strings.Replace(status, "ing", "e", 1)
	// we should end up with a command like 'git merge --continue'

	// it's impossible for a rebase to require a commit so we'll use a subprocess only if it's a merge
	if status == commands.REBASE_MODE_MERGING && command != REBASE_OPTION_ABORT && gui.Config.GetUserConfig().Git.Merging.ManualCommit {
		sub := gitCommand.NewCmdObj("git " + commandType + " --" + command)
		if sub != nil {
			return gui.runSubprocessWithSuspenseAndRefresh(sub)
		}
		return nil
	}
	result := gitCommand.GenericMergeOrRebaseAction(commandType, command)
	if err := gui.handleGenericMergeCommandResult(result); err != nil {
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

func (gui *Gui) handleGenericMergeCommandResult(result error) error {
	if err := gui.refreshSidePanels(refreshOptions{mode: ASYNC}); err != nil {
		return err
	}
	if result == nil {
		return nil
	} else if strings.Contains(result.Error(), "No changes - did you forget to use") {
		return gui.genericMergeCommand(REBASE_OPTION_SKIP)
	} else if strings.Contains(result.Error(), "The previous cherry-pick is now empty") {
		return gui.genericMergeCommand(REBASE_OPTION_CONTINUE)
	} else if strings.Contains(result.Error(), "No rebase in progress?") {
		// assume in this case that we're already done
		return nil
	} else if isMergeConflictErr(result.Error()) {
		return gui.ask(askOpts{
			title:               gui.Tr.FoundConflictsTitle,
			prompt:              gui.Tr.FoundConflicts,
			handlersManageFocus: true,
			handleConfirm: func() error {
				return gui.pushContext(gui.State.Contexts.Files)
			},
			handleClose: func() error {
				if err := gui.returnFromContext(); err != nil {
					return err
				}

				return gui.genericMergeCommand(REBASE_OPTION_ABORT)
			},
		})
	} else {
		return gui.createErrorPanel(result.Error())
	}
}

func (gui *Gui) abortMergeOrRebaseWithConfirm() error {
	// prompt user to confirm that they want to abort, then do it
	mode := gui.workingTreeStateNoun()
	return gui.ask(askOpts{
		title:  fmt.Sprintf(gui.Tr.AbortTitle, mode),
		prompt: fmt.Sprintf(gui.Tr.AbortPrompt, mode),
		handleConfirm: func() error {
			return gui.genericMergeCommand(REBASE_OPTION_ABORT)
		},
	})
}

func (gui *Gui) workingTreeStateNoun() string {
	workingTreeState := gui.GitCommand.WorkingTreeState()
	switch workingTreeState {
	case commands.REBASE_MODE_NORMAL:
		return ""
	case commands.REBASE_MODE_MERGING:
		return "merge"
	default:
		return "rebase"
	}
}
