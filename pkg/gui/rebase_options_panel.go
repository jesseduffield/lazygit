package gui

import (
	"fmt"
	"strings"

	. "github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) handleCreateRebaseOptionsMenu() error {
	options := []string{"continue", "abort"}
	title := gui.Tr.MergeOptionsTitle

	if gui.Git.Status().IsRebasing() {
		options = append(options, "skip")
		title = gui.Tr.RebaseOptionsTitle
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

	return gui.createMenu(title, menuItems, createMenuOptions{showCancel: true})
}

func (gui *Gui) genericMergeCommand(action string) error {
	if gui.Git.Status().InNormalWorkingTreeState() {
		return gui.CreateErrorPanel(gui.Tr.NotMergingOrRebasing)
	}

	gitCommand := gui.Git.WithSpan(fmt.Sprintf("Merge/Rebase: %s", action))

	// it's impossible for a rebase to require a commit so we'll use a subprocess only if it's a merge
	if gui.Git.Status().IsMerging() && action != "abort" && gui.Config.GetUserConfig().Git.Merging.ManualCommit {
		return gui.runSubprocessWithSuspenseAndRefresh(
			gitCommand.GenericMergeOrRebaseCmdObj(action),
		)
	}

	command := gui.Git.MergeOrRebase()
	actionErr := gui.Git.GenericMergeOrRebaseAction(command, action)
	if err := gui.handleGenericMergeCommandResult(actionErr); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) handleGenericMergeCommandResult(result error) error {
	if err := gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC}); err != nil {
		return err
	}
	if result == nil {
		return nil
	} else if strings.Contains(result.Error(), "No changes - did you forget to use") {
		return gui.genericMergeCommand("skip")
	} else if strings.Contains(result.Error(), "The previous cherry-pick is now empty") {
		return gui.genericMergeCommand("continue")
	} else if strings.Contains(result.Error(), "No rebase in progress?") {
		// assume in this case that we're already done
		return nil
	} else if strings.Contains(result.Error(), "When you have resolved this problem") || strings.Contains(result.Error(), "fix conflicts") || strings.Contains(result.Error(), "Resolve all conflicts manually") {
		return gui.Ask(AskOpts{
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

				return gui.genericMergeCommand("abort")
			},
		})
	} else {
		return gui.CreateErrorPanel(result.Error())
	}
}
