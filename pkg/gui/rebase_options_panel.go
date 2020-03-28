package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
)

func (gui *Gui) handleCreateRebaseOptionsMenu(g *gocui.Gui, v *gocui.View) error {
	options := []string{"continue", "abort"}

	if gui.GitCommand.WorkingTreeState() == "rebasing" {
		options = append(options, "skip")
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
	if gui.GitCommand.WorkingTreeState() == "merging" {
		title = gui.Tr.SLocalize("MergeOptionsTitle")
	} else {
		title = gui.Tr.SLocalize("RebaseOptionsTitle")
	}

	return gui.createMenu(title, menuItems, createMenuOptions{showCancel: true})
}

func (gui *Gui) genericMergeCommand(command string) error {
	status := gui.GitCommand.WorkingTreeState()

	if status != "merging" && status != "rebasing" {
		return gui.createErrorPanel(gui.Tr.SLocalize("NotMergingOrRebasing"))
	}

	commandType := strings.Replace(status, "ing", "e", 1)
	// we should end up with a command like 'git merge --continue'

	// it's impossible for a rebase to require a commit so we'll use a subprocess only if it's a merge
	if status == "merging" && command != "abort" && gui.Config.GetUserConfig().GetBool("git.merging.manualCommit") {
		sub := gui.OSCommand.PrepareSubProcess("git", commandType, fmt.Sprintf("--%s", command))
		if sub != nil {
			gui.SubProcess = sub
			return gui.Errors.ErrSubProcess
		}
		return nil
	}
	result := gui.GitCommand.GenericMerge(commandType, command)
	if err := gui.handleGenericMergeCommandResult(result); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) handleGenericMergeCommandResult(result error) error {
	if err := gui.refreshSidePanels(refreshOptions{mode: ASYNC}); err != nil {
		return err
	}
	if result == nil {
		return nil
	} else if result == gui.Errors.ErrSubProcess {
		return result
	} else if strings.Contains(result.Error(), "No changes - did you forget to use") {
		return gui.genericMergeCommand("skip")
	} else if strings.Contains(result.Error(), "The previous cherry-pick is now empty") {
		return gui.genericMergeCommand("continue")
	} else if strings.Contains(result.Error(), "No rebase in progress?") {
		// assume in this case that we're already done
		return nil
	} else if strings.Contains(result.Error(), "When you have resolved this problem") || strings.Contains(result.Error(), "fix conflicts") || strings.Contains(result.Error(), "Resolve all conflicts manually") {
		return gui.createConfirmationPanel(gui.g, gui.getFilesView(), true, gui.Tr.SLocalize("FoundConflictsTitle"), gui.Tr.SLocalize("FoundConflicts"),
			func(g *gocui.Gui, v *gocui.View) error {
				return nil
			}, func(g *gocui.Gui, v *gocui.View) error {
				return gui.genericMergeCommand("abort")
			},
		)
	} else {
		return gui.createErrorPanel(result.Error())
	}
}
