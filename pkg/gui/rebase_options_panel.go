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
func (r *option) GetDisplayStrings() []string {
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
		return gui.genericRebaseCommand(command)
	}

	var title string
	if gui.State.WorkingTreeState == "merging" {
		title = gui.Tr.SLocalize("MergeOptionsTitle")
	} else {
		title = gui.Tr.SLocalize("RebaseOptionsTitle")
	}

	return gui.createMenu(title, options, handleMenuPress)
}

func (gui *Gui) genericRebaseCommand(command string) error {
	status := gui.State.WorkingTreeState

	if status != "merging" && status != "rebasing" {
		return gui.createErrorPanel(gui.g, gui.Tr.SLocalize("NotMergingOrRebasing"))
	}

	commandType := strings.Replace(status, "ing", "e", 1)
	// we should end up with a command like 'git merge --continue'
	fullCommand := fmt.Sprintf("git %s --%s", commandType, command)

	if err := gui.OSCommand.RunCommand(fullCommand); err != nil { // this guy freezes because it's opening your editor to enter the merge message
		return err
	}
	return nil
}
