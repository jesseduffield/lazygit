package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

func (gui *Gui) resetToRef(ref string, strength string, envVars []string) error {
	if err := gui.Git.Commit.ResetToCommit(ref, strength, envVars); err != nil {
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
				style.FgRed.Sprintf("reset --%s %s", strength, ref),
			},
			onPress: func() error {
				gui.logAction("Reset")
				return gui.resetToRef(ref, strength, []string{})
			},
		}
	}

	return gui.createMenu(createMenuOptions{
		title: fmt.Sprintf("%s %s", gui.Tr.LcResetTo, ref),
		items: menuItems,
	})
}
