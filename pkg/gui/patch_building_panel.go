package gui

import (
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// getFromAndReverseArgsForDiff tells us the from and reverse args to be used in a diff command. If we're not in diff mode we'll end up with the equivalent of a `git show` i.e `git diff blah^..blah`.
func (gui *Gui) getFromAndReverseArgsForDiff(to string) (string, bool) {
	from := to + "^"
	reverse := false

	if gui.State.Modes.Diffing.Active() {
		reverse = gui.State.Modes.Diffing.Reverse
		from = gui.State.Modes.Diffing.Ref
	}

	return from, reverse
}

func (gui *Gui) refreshPatchBuildingPanel(selectedLineIdx int) error {
	if !gui.Git.Patch.PatchManager.Active() {
		return gui.handleEscapePatchBuildingPanel()
	}

	gui.Views.Main.Title = "Patch"
	gui.Views.Secondary.Title = "Custom Patch"

	// get diff from commit file that's currently selected
	node := gui.getSelectedCommitFileNode()
	if node == nil {
		return nil
	}

	to := gui.State.CommitFileTreeViewModel.GetParent()
	from, reverse := gui.getFromAndReverseArgsForDiff(to)
	diff, err := gui.Git.WorkingTree.ShowFileDiff(from, to, reverse, node.GetPath(), true)
	if err != nil {
		return err
	}

	secondaryDiff := gui.Git.Patch.PatchManager.RenderPatchForFile(node.GetPath(), true, false, true)
	if err != nil {
		return err
	}

	empty, err := gui.refreshLineByLinePanel(diff, secondaryDiff, false, selectedLineIdx)
	if err != nil {
		return err
	}

	if empty {
		return gui.handleEscapePatchBuildingPanel()
	}

	return nil
}

func (gui *Gui) handleRefreshPatchBuildingPanel(selectedLineIdx int) error {
	gui.Mutexes.LineByLinePanelMutex.Lock()
	defer gui.Mutexes.LineByLinePanelMutex.Unlock()

	return gui.refreshPatchBuildingPanel(selectedLineIdx)
}

func (gui *Gui) onPatchBuildingFocus(selectedLineIdx int) error {
	gui.Mutexes.LineByLinePanelMutex.Lock()
	defer gui.Mutexes.LineByLinePanelMutex.Unlock()

	if gui.State.Panels.LineByLine == nil || selectedLineIdx != -1 {
		return gui.refreshPatchBuildingPanel(selectedLineIdx)
	}

	return nil
}

func (gui *Gui) handleToggleSelectionForPatch() error {
	err := gui.withLBLActiveCheck(func(state *LblPanelState) error {
		toggleFunc := gui.Git.Patch.PatchManager.AddFileLineRange
		filename := gui.getSelectedCommitFileName()
		includedLineIndices, err := gui.Git.Patch.PatchManager.GetFileIncLineIndices(filename)
		if err != nil {
			return err
		}
		currentLineIsStaged := utils.IncludesInt(includedLineIndices, state.GetSelectedLineIdx())
		if currentLineIsStaged {
			toggleFunc = gui.Git.Patch.PatchManager.RemoveFileLineRange
		}

		// add range of lines to those set for the file
		node := gui.getSelectedCommitFileNode()
		if node == nil {
			return nil
		}

		firstLineIdx, lastLineIdx := state.SelectedRange()

		if err := toggleFunc(node.GetPath(), firstLineIdx, lastLineIdx); err != nil {
			// might actually want to return an error here
			gui.Log.Error(err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err := gui.refreshCommitFilesView(); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleEscapePatchBuildingPanel() error {
	gui.escapeLineByLinePanel()

	if gui.Git.Patch.PatchManager.IsEmpty() {
		gui.Git.Patch.PatchManager.Reset()
	}

	if gui.currentContext().GetKey() == gui.State.Contexts.PatchBuilding.GetKey() {
		return gui.pushContext(gui.State.Contexts.CommitFiles)
	} else {
		// need to re-focus in case the secondary view should now be hidden
		return gui.currentContext().HandleFocus()
	}
}

func (gui *Gui) secondaryPatchPanelUpdateOpts() *viewUpdateOpts {
	if gui.Git.Patch.PatchManager.Active() {
		patch := gui.Git.Patch.PatchManager.RenderAggregatedPatchColored(false)

		return &viewUpdateOpts{
			title:     "Custom Patch",
			noWrap:    true,
			highlight: true,
			task:      NewRenderStringWithoutScrollTask(patch),
		}
	}

	return nil
}
