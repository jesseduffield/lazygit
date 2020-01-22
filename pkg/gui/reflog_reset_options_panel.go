package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

type reflogResetOption struct {
	handler     func() error
	description string
	command     string
}

// GetDisplayStrings is a function.
func (r *reflogResetOption) GetDisplayStrings(isFocused bool) []string {
	return []string{r.description, color.New(color.FgRed).Sprint(r.command)}
}

func (gui *Gui) handleCreateReflogResetMenu(g *gocui.Gui, v *gocui.View) error {
	commit := gui.getSelectedReflogCommit()

	resetFunction := func(reset func(string) error) func() error {
		return func() error {
			if err := reset(commit.Sha); err != nil {
				return gui.createErrorPanel(gui.g, err.Error())
			}

			gui.State.Panels.ReflogCommits.SelectedLine = 0

			return gui.refreshSidePanels(gui.g)
		}
	}

	options := []*reflogResetOption{
		{
			description: gui.Tr.SLocalize("hardReset"),
			command:     fmt.Sprintf("reset --hard %s", commit.Sha),
			handler:     resetFunction(gui.GitCommand.ResetHard),
		},
		{
			description: gui.Tr.SLocalize("softReset"),
			command:     fmt.Sprintf("reset --soft %s", commit.Sha),
			handler:     resetFunction(gui.GitCommand.ResetSoft),
		},
		{
			description: gui.Tr.SLocalize("cancel"),
			handler: func() error {
				return nil
			},
		},
	}

	handleMenuPress := func(index int) error {
		return options[index].handler()
	}

	return gui.createMenu("", options, len(options), handleMenuPress)
}
