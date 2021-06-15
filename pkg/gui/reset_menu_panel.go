package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) resetToRef(ref string, strength commands.ResetStrength, span string, options commands.ResetToRefOpts) error {
	if err := gui.Git.WithSpan(span).Branches().ResetToRef(ref, strength, options); err != nil {
		return gui.SurfaceError(err)
	}

	gui.State.Panels.Commits.SelectedLineIdx = 0
	gui.State.Panels.ReflogCommits.SelectedLineIdx = 0
	// loading a heap of commits is slow so we limit them whenever doing a reset
	gui.State.Panels.Commits.LimitCommits = true

	if err := gui.pushContext(gui.State.Contexts.BranchCommits); err != nil {
		return err
	}

	if err := gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{FILES, BRANCHES, REFLOG, COMMITS}}); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) createResetMenu(ref string) error {
	strengths := []commands.ResetStrength{commands.SOFT, commands.MIXED, commands.HARD}
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
				return gui.resetToRef(ref, strength, "Reset", commands.ResetToRefOpts{})
			},
		}
	}

	return gui.createMenu(fmt.Sprintf("%s %s", gui.Tr.LcResetTo, ref), menuItems, createMenuOptions{showCancel: true})
}
