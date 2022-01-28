package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type RebaseOption string

const (
	REBASE_OPTION_CONTINUE = "continue"
	REBASE_OPTION_ABORT    = "abort"
	REBASE_OPTION_SKIP     = "skip"
)

func (gui *Gui) handleCreateRebaseOptionsMenu() error {
	options := []string{REBASE_OPTION_CONTINUE, REBASE_OPTION_ABORT}

	if gui.Git.Status.WorkingTreeState() == enums.REBASE_MODE_REBASING {
		options = append(options, REBASE_OPTION_SKIP)
	}

	menuItems := make([]*popup.MenuItem, len(options))
	for i, option := range options {
		// note to self. Never, EVER, close over loop variables in a function
		option := option
		menuItems[i] = &popup.MenuItem{
			DisplayString: option,
			OnPress: func() error {
				return gui.genericMergeCommand(option)
			},
		}
	}

	var title string
	if gui.Git.Status.WorkingTreeState() == enums.REBASE_MODE_MERGING {
		title = gui.Tr.MergeOptionsTitle
	} else {
		title = gui.Tr.RebaseOptionsTitle
	}

	return gui.PopupHandler.Menu(popup.CreateMenuOptions{Title: title, Items: menuItems})
}

func (gui *Gui) genericMergeCommand(command string) error {
	status := gui.Git.Status.WorkingTreeState()

	if status != enums.REBASE_MODE_MERGING && status != enums.REBASE_MODE_REBASING {
		return gui.PopupHandler.ErrorMsg(gui.Tr.NotMergingOrRebasing)
	}

	gui.logAction(fmt.Sprintf("Merge/Rebase: %s", command))

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
	if status == enums.REBASE_MODE_MERGING && command != REBASE_OPTION_ABORT && gui.UserConfig.Git.Merging.ManualCommit {
		// TODO: see if we should be calling more of the code from gui.Git.Rebase.GenericMergeOrRebaseAction
		return gui.runSubprocessWithSuspenseAndRefresh(
			gui.Git.Rebase.GenericMergeOrRebaseActionCmdObj(commandType, command),
		)
	}
	result := gui.Git.Rebase.GenericMergeOrRebaseAction(commandType, command)
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
	if err := gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC}); err != nil {
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
		return gui.PopupHandler.Ask(popup.AskOpts{
			Title:               gui.Tr.FoundConflictsTitle,
			Prompt:              gui.Tr.FoundConflicts,
			HandlersManageFocus: true,
			HandleConfirm: func() error {
				return gui.pushContext(gui.State.Contexts.Files)
			},
			HandleClose: func() error {
				if err := gui.returnFromContext(); err != nil {
					return err
				}

				return gui.genericMergeCommand(REBASE_OPTION_ABORT)
			},
		})
	} else {
		return gui.PopupHandler.ErrorMsg(result.Error())
	}
}

func (gui *Gui) abortMergeOrRebaseWithConfirm() error {
	// prompt user to confirm that they want to abort, then do it
	mode := gui.workingTreeStateNoun()
	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  fmt.Sprintf(gui.Tr.AbortTitle, mode),
		Prompt: fmt.Sprintf(gui.Tr.AbortPrompt, mode),
		HandleConfirm: func() error {
			return gui.genericMergeCommand(REBASE_OPTION_ABORT)
		},
	})
}

func (gui *Gui) workingTreeStateNoun() string {
	workingTreeState := gui.Git.Status.WorkingTreeState()
	switch workingTreeState {
	case enums.REBASE_MODE_NONE:
		return ""
	case enums.REBASE_MODE_MERGING:
		return "merge"
	default:
		return "rebase"
	}
}
