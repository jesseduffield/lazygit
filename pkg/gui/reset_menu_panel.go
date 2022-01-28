package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) resetToRef(ref string, strength string, envVars []string) error {
	if err := gui.Git.Commit.ResetToCommit(ref, strength, envVars); err != nil {
		return gui.PopupHandler.Error(err)
	}

	gui.State.Panels.Commits.SelectedLineIdx = 0
	gui.State.Panels.ReflogCommits.SelectedLineIdx = 0
	// loading a heap of commits is slow so we limit them whenever doing a reset
	gui.State.Panels.Commits.LimitCommits = true

	if err := gui.pushContext(gui.State.Contexts.BranchCommits); err != nil {
		return err
	}

	if err := gui.refreshSidePanels(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES, types.BRANCHES, types.REFLOG, types.COMMITS}}); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) createResetMenu(ref string) error {
	strengths := []string{"soft", "mixed", "hard"}
	menuItems := make([]*popup.MenuItem, len(strengths))
	for i, strength := range strengths {
		strength := strength
		menuItems[i] = &popup.MenuItem{
			DisplayStrings: []string{
				fmt.Sprintf("%s reset", strength),
				style.FgRed.Sprintf("reset --%s %s", strength, ref),
			},
			OnPress: func() error {
				gui.logAction("Reset")
				return gui.resetToRef(ref, strength, []string{})
			},
		}
	}

	return gui.PopupHandler.Menu(popup.CreateMenuOptions{
		Title: fmt.Sprintf("%s %s", gui.Tr.LcResetTo, ref),
		Items: menuItems,
	})
}
