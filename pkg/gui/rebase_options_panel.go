package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
)

type option struct {
	value string
}

// GetDisplayStrings is a function.
func (r *option) GetDisplayStrings(isFocused bool) []string {
	return []string{r.value}
}

func (gui *Gui) handleCreateRebaseOptionsMenu(g *gocui.Gui, v *gocui.View) error {
	options := []*option{
		{value: "continue"},
		{value: "abort"},
	}

	if gui.State.WorkingTreeState == "rebasing" {
		options = append(options, &option{value: "skip"})
	}

	handleMenuPress := func(index int) error {
		command := options[index].value
		return gui.genericMergeCommand(command)
	}

	var title string
	if gui.State.WorkingTreeState == "merging" {
		title = gui.Tr.SLocalize("MergeOptionsTitle")
	} else {
		title = gui.Tr.SLocalize("RebaseOptionsTitle")
	}

	return gui.createMenu(title, options, handleMenuPress)
}

func (gui *Gui) genericMergeCommand(command string) error {
	status := gui.State.WorkingTreeState

	if status != "merging" && status != "rebasing" {
		return gui.createErrorPanel(gui.g, gui.Tr.SLocalize("NotMergingOrRebasing"))
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
	if err := gui.refreshSidePanels(gui.g); err != nil {
		return err
	}
	if result == nil {
		return nil
	} else if result == gui.Errors.ErrSubProcess {
		return result
	} else if strings.Contains(result.Error(), "No changes - did you forget to use") {
		return gui.genericMergeCommand("skip")
	} else if strings.Contains(result.Error(), "When you have resolved this problem") || strings.Contains(result.Error(), "fix conflicts") || strings.Contains(result.Error(), "Resolve all conflicts manually") {
		return gui.createConfirmationPanel(gui.g, gui.getFilesView(), gui.Tr.SLocalize("FoundConflictsTitle"), gui.Tr.SLocalize("FoundConflicts"),
			func(g *gocui.Gui, v *gocui.View) error {
				return nil
			}, func(g *gocui.Gui, v *gocui.View) error {
				return gui.genericMergeCommand("abort")
			},
		)
	} else {
		return gui.createErrorPanel(gui.g, result.Error())
	}
}
