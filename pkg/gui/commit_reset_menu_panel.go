package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) handleCreateCommitResetMenu(g *gocui.Gui, v *gocui.View) error {
	commit := gui.getSelectedCommit(g)
	if commit == nil {
		return gui.createErrorPanel(gui.g, gui.Tr.SLocalize("NoCommitsThisBranch"))
	}

	strengths := []string{"soft", "mixed", "hard"}
	menuItems := make([]*menuItem, len(strengths))
	for i, strength := range strengths {
		innerStrength := strength
		menuItems[i] = &menuItem{
			displayStrings: []string{
				fmt.Sprintf("%s reset", strength),
				color.New(color.FgRed).Sprint(
					fmt.Sprintf("reset --%s %s", strength, commit.Sha),
				),
			},
			onPress: func() error {
				if err := gui.GitCommand.ResetToCommit(commit.Sha, innerStrength); err != nil {
					return err
				}

				if err := gui.refreshCommits(g); err != nil {
					return err
				}
				if err := gui.refreshFiles(); err != nil {
					return err
				}
				if err := gui.resetOrigin(gui.getCommitsView()); err != nil {
					return err
				}

				gui.State.Panels.Commits.SelectedLine = 0
				return gui.handleCommitSelect(g, gui.getCommitsView())
			},
		}
	}

	return gui.createMenu(fmt.Sprintf("%s %s", gui.Tr.SLocalize("resetTo"), commit.Sha), menuItems, createMenuOptions{showCancel: true})
}
