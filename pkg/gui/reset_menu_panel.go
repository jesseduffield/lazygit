package gui

import (
	"fmt"

	"github.com/fatih/color"
)

func (gui *Gui) createResetMenu(ref string) error {
	strengths := []string{"soft", "mixed", "hard"}
	menuItems := make([]*menuItem, len(strengths))
	for i, strength := range strengths {
		innerStrength := strength
		menuItems[i] = &menuItem{
			displayStrings: []string{
				fmt.Sprintf("%s reset", strength),
				color.New(color.FgRed).Sprint(
					fmt.Sprintf("reset --%s %s", strength, ref),
				),
			},
			onPress: func() error {
				if err := gui.GitCommand.ResetToCommit(ref, innerStrength); err != nil {
					return gui.createErrorPanel(gui.g, err.Error())
				}

				gui.switchCommitsPanelContext("branch-commits")
				gui.State.Panels.Commits.SelectedLine = 0
				gui.State.Panels.ReflogCommits.SelectedLine = 0

				if err := gui.refreshCommits(gui.g); err != nil {
					return err
				}
				if err := gui.refreshFiles(); err != nil {
					return err
				}
				if err := gui.resetOrigin(gui.getCommitsView()); err != nil {
					return err
				}

				return gui.handleCommitSelect(gui.g, gui.getCommitsView())
			},
		}
	}

	return gui.createMenu(fmt.Sprintf("%s %s", gui.Tr.SLocalize("resetTo"), ref), menuItems, createMenuOptions{showCancel: true})
}
