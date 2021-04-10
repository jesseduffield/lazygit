package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

func (gui *Gui) resetToRef(ref string, strength string, span string, options oscommands.RunCommandOptions) error {
	if err := gui.GitCommand.WithSpan(span).ResetToCommit(ref, strength, options); err != nil {
		return gui.surfaceError(err)
	}

	gui.State.Panels.Commits.SelectedLineIdx = 0
	gui.State.Panels.ReflogCommits.SelectedLineIdx = 0
	// loading a heap of commits is slow so we limit them whenever doing a reset
	gui.State.Panels.Commits.LimitCommits = true

	if err := gui.pushContext(gui.State.Contexts.BranchCommits); err != nil {
		return err
	}

	if err := gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{FILES, BRANCHES, REFLOG, COMMITS}}); err != nil {
		return err
	}

	return nil
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
				return gui.resetToRef(ref, strength, "Reset", oscommands.RunCommandOptions{})
			},
		}
	}

	return gui.createMenu(fmt.Sprintf("%s %s", gui.Tr.LcResetTo, ref), menuItems, createMenuOptions{showCancel: true})
}
