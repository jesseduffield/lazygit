package gui

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) refreshStagingPanel(forceSecondaryFocused bool, selectedLineIdx int) error {
	gui.splitMainPanel(true)

	file := gui.getSelectedFile()
	if file == nil || (!file.HasUnstagedChanges && !file.HasStagedChanges) {
		return gui.handleStagingEscape()
	}

	secondaryFocused := false
	if forceSecondaryFocused {
		secondaryFocused = true
	} else if gui.State.Panels.LineByLine != nil {
		secondaryFocused = gui.State.Panels.LineByLine.SecondaryFocused
	}

	if (secondaryFocused && !file.HasStagedChanges) || (!secondaryFocused && !file.HasUnstagedChanges) {
		secondaryFocused = !secondaryFocused
	}

	if secondaryFocused {
		gui.Views.Main.Title = gui.Tr.StagedChanges
		gui.Views.Secondary.Title = gui.Tr.UnstagedChanges
	} else {
		gui.Views.Main.Title = gui.Tr.UnstagedChanges
		gui.Views.Secondary.Title = gui.Tr.StagedChanges
	}

	// note for custom diffs, we'll need to send a flag here saying not to use the custom diff
	diff := gui.Git.WorktreeFileDiff(file, true, secondaryFocused)
	secondaryDiff := gui.Git.WorktreeFileDiff(file, true, !secondaryFocused)

	// if we have e.g. a deleted file with nothing else to the diff will have only
	// 4-5 lines in which case we'll swap panels
	if len(strings.Split(diff, "\n")) < 5 {
		if len(strings.Split(secondaryDiff, "\n")) < 5 {
			return gui.handleStagingEscape()
		}
		secondaryFocused = !secondaryFocused
		diff, secondaryDiff = secondaryDiff, diff
	}

	empty, err := gui.refreshLineByLinePanel(diff, secondaryDiff, secondaryFocused, selectedLineIdx)
	if err != nil {
		return err
	}

	if empty {
		return gui.handleStagingEscape()
	}

	return nil
}

func (gui *Gui) handleTogglePanelClick() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		state.SecondaryFocused = !state.SecondaryFocused

		return gui.refreshStagingPanel(false, gui.Views.Secondary.SelectedLineIdx())
	})
}

func (gui *Gui) handleRefreshStagingPanel(forceSecondaryFocused bool, selectedLineIdx int) error {
	gui.Mutexes.LineByLinePanelMutex.Lock()
	defer gui.Mutexes.LineByLinePanelMutex.Unlock()

	return gui.refreshStagingPanel(forceSecondaryFocused, selectedLineIdx)
}

func (gui *Gui) handleTogglePanel() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		state.SecondaryFocused = !state.SecondaryFocused
		return gui.refreshStagingPanel(false, -1)
	})
}

func (gui *Gui) handleStagingEscape() error {
	gui.escapeLineByLinePanel()

	return gui.pushContext(gui.State.Contexts.Files)
}

func (gui *Gui) handleToggleStagedSelection() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		return gui.applySelection(state.SecondaryFocused, state)
	})
}

func (gui *Gui) handleResetSelection() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		if state.SecondaryFocused {
			// for backwards compatibility
			return gui.applySelection(true, state)
		}

		if !gui.Config.GetUserConfig().Gui.SkipUnstageLineWarning {
			return gui.Ask(AskOpts{
				Title:               gui.Tr.UnstageLinesTitle,
				Prompt:              gui.Tr.UnstageLinesPrompt,
				HandlersManageFocus: true,
				HandleConfirm: func() error {
					return gui.withLBLActiveCheck(func(state *LblPanelState) error {
						if err := gui.pushContext(gui.State.Contexts.Staging); err != nil {
							return err
						}

						return gui.applySelection(true, state)
					})
				},
				HandleClose: func() error {
					return gui.pushContext(gui.State.Contexts.Staging)
				},
			})
		} else {
			return gui.applySelection(true, state)
		}
	})
}

func (gui *Gui) applySelection(reverse bool, state *LblPanelState) error {
	file := gui.getSelectedFile()
	if file == nil {
		return nil
	}

	firstLineIdx, lastLineIdx := state.SelectedRange()
	patch := patch.ModifiedPatchForRange(gui.Log, file.Name, state.GetDiff(), firstLineIdx, lastLineIdx, reverse, false)

	if patch == "" {
		return nil
	}

	// apply the patch then refresh this panel
	// create a new temp file with the patch, then call git apply with that patch
	applyFlags := []string{}
	if !reverse || state.SecondaryFocused {
		applyFlags = append(applyFlags, "cached")
	}
	err := gui.Git.WithSpan(gui.Tr.Spans.ApplyPatch).Worktree().ApplyPatch(patch, applyFlags...)
	if err != nil {
		return gui.SurfaceError(err)
	}

	if state.SelectingRange() {
		state.SetLineSelectMode()
	}

	if err := gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{FILES}}); err != nil {
		return err
	}
	if err := gui.refreshStagingPanel(false, -1); err != nil {
		return err
	}
	return nil
}
