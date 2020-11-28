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

func (gui *Gui) refreshPatchBuildingPanel(selectedLineIdx int, state *lBlPanelState) error {
	if !gui.GitCommand.PatchManager.Active() {
		return gui.handleEscapePatchBuildingPanel()
	}

	gui.splitMainPanel(true)

	gui.getMainView().Title = "Patch"
	gui.getSecondaryView().Title = "Custom Patch"

	// get diff from commit file that's currently selected
	commitFile := gui.getSelectedCommitFile()
	if commitFile == nil {
		return nil
	}

	to := commitFile.Parent
	from, reverse := gui.getFromAndReverseArgsForDiff(to)
	diff, err := gui.GitCommand.ShowFileDiff(from, to, reverse, commitFile.Name, true)
	if err != nil {
		return err
	}

	secondaryDiff := gui.GitCommand.PatchManager.RenderPatchForFile(commitFile.Name, true, false, true)
	if err != nil {
		return err
	}

	empty, err := gui.refreshLineByLinePanel(diff, secondaryDiff, false, selectedLineIdx, state)
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

	return gui.refreshPatchBuildingPanel(selectedLineIdx, gui.State.Panels.LineByLine)
}

func (gui *Gui) handleToggleSelectionForPatch() error {
	err := gui.withLBLActiveCheck(func(state *lBlPanelState) error {
		toggleFunc := gui.GitCommand.PatchManager.AddFileLineRange
		filename := gui.getSelectedCommitFileName()
		includedLineIndices, err := gui.GitCommand.PatchManager.GetFileIncLineIndices(filename)
		if err != nil {
			return err
		}
		currentLineIsStaged := utils.IncludesInt(includedLineIndices, state.SelectedLineIdx)
		if currentLineIsStaged {
			toggleFunc = gui.GitCommand.PatchManager.RemoveFileLineRange
		}

		// add range of lines to those set for the file
		commitFile := gui.getSelectedCommitFile()
		if commitFile == nil {
			return nil
		}

		if err := toggleFunc(commitFile.Name, state.FirstLineIdx, state.LastLineIdx); err != nil {
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

	if gui.GitCommand.PatchManager.IsEmpty() {
		gui.GitCommand.PatchManager.Reset()
	}

	if gui.currentContext().GetKey() == gui.Contexts.PatchBuilding.Context.GetKey() {
		return gui.pushContext(gui.Contexts.CommitFiles.Context)
	} else {
		// need to re-focus in case the secondary view should now be hidden
		return gui.currentContext().HandleFocus()
	}
}

func (gui *Gui) secondaryPatchPanelUpdateOpts() *viewUpdateOpts {
	if gui.GitCommand.PatchManager.Active() {
		patch := gui.GitCommand.PatchManager.RenderAggregatedPatchColored(false)

		return &viewUpdateOpts{
			title:     "Custom Patch",
			noWrap:    true,
			highlight: true,
			task:      gui.createRenderStringWithoutScrollTask(patch),
		}
	}

	return nil
}
