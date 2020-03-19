package gui

import (
	"fmt"

	"github.com/fatih/color"
)

func (gui *Gui) resetToRef(ref string, strength string) error {
	if err := gui.GitCommand.ResetToCommit(ref, strength); err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
	}

	if err := gui.switchCommitsPanelContext("branch-commits"); err != nil {
		return err
	}

	gui.State.Panels.Commits.SelectedLine = 0
	gui.State.Panels.ReflogCommits.SelectedLine = 0

	if err := gui.refreshCommits(gui.g); err != nil {
		return err
	}
	if err := gui.refreshFiles(); err != nil {
		return err
	}
	if err := gui.refreshBranches(gui.g); err != nil {
		return err
	}
	if err := gui.resetOrigin(gui.getCommitsView()); err != nil {
		return err
	}

	return gui.handleCommitSelect(gui.g, gui.getCommitsView())
}

func (gui *Gui) createResetMenu(ref string) error {
	strengths := []string{"soft", "mixed", "hard"}
	menuItems := make([]*menuItem, len(strengths))
	for i, strength := range strengths {
		strength := strength
		menuItems[i] = &menuItem{
			displayStrings: []string{
				fmt.Sprintf("%s reset", strength),
				color.New(color.FgRed).Sprint(
					fmt.Sprintf("reset --%s %s", strength, ref),
				),
			},
			onPress: func() error {
				return gui.resetToRef(ref, strength)
			},
		}
	}

	return gui.createMenu(fmt.Sprintf("%s %s", gui.Tr.SLocalize("resetTo"), ref), menuItems, createMenuOptions{showCancel: true})
}
