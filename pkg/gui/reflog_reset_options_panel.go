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
	red := color.New(color.FgRed)

	resetFunction := func(reset func(string) error) func() error {
		return func() error {
			if err := reset(commit.Sha); err != nil {
				return gui.createErrorPanel(gui.g, err.Error())
			}

			gui.State.Panels.ReflogCommits.SelectedLine = 0

			return gui.refreshSidePanels(gui.g)
		}
	}

	menuItems := []*menuItem{
		{
			displayStrings: []string{
				gui.Tr.SLocalize("hardReset"),
				red.Sprint(fmt.Sprintf("reset --hard %s", commit.Sha)),
			},
			onPress: resetFunction(gui.GitCommand.ResetHard),
		},
		{
			displayStrings: []string{
				gui.Tr.SLocalize("softReset"),
				red.Sprint(fmt.Sprintf("reset --soft %s", commit.Sha)),
			},
			onPress: resetFunction(gui.GitCommand.ResetSoft),
		},
	}

	return gui.createMenu("", menuItems, createMenuOptions{showCancel: true})
}
