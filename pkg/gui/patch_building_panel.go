package gui

import "github.com/samber/lo"

func (gui *Gui) refreshPatchBuildingPanel(selectedLineIdx int) error {
	if !gui.git.Patch.PatchManager.Active() {
		return gui.handleEscapePatchBuildingPanel()
	}

	gui.Views.Main.Title = "Patch"
	gui.Views.Secondary.Title = "Custom Patch"

	// get diff from commit file that's currently selected
	node := gui.State.Contexts.CommitFiles.GetSelected()
	if node == nil {
		return nil
	}

	ref := gui.State.Contexts.CommitFiles.CommitFileTreeViewModel.GetRef()
	to := ref.RefName()
	from, reverse := gui.State.Modes.Diffing.GetFromAndReverseArgsForDiff(ref.ParentRefName())
	diff, err := gui.git.WorkingTree.ShowFileDiff(from, to, reverse, node.GetPath(), true)
	if err != nil {
		return err
	}

	secondaryDiff := gui.git.Patch.PatchManager.RenderPatchForFile(node.GetPath(), true, false, true)
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
		toggleFunc := gui.git.Patch.PatchManager.AddFileLineRange
		filename := gui.getSelectedCommitFileName()
		includedLineIndices, err := gui.git.Patch.PatchManager.GetFileIncLineIndices(filename)
		if err != nil {
			return err
		}
		currentLineIsStaged := lo.Contains(includedLineIndices, state.GetSelectedLineIdx())
		if currentLineIsStaged {
			toggleFunc = gui.git.Patch.PatchManager.RemoveFileLineRange
		}

		// add range of lines to those set for the file
		node := gui.State.Contexts.CommitFiles.GetSelected()
		if node == nil {
			return nil
		}

		firstLineIdx, lastLineIdx := state.SelectedRange()

		if err := toggleFunc(node.GetPath(), firstLineIdx, lastLineIdx); err != nil {
			// might actually want to return an error here
			gui.c.Log.Error(err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	if err := gui.handleRefreshPatchBuildingPanel(-1); err != nil {
		return err
	}

	if err := gui.refreshCommitFilesContext(); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleEscapePatchBuildingPanel() error {
	gui.escapeLineByLinePanel()

	if gui.git.Patch.PatchManager.IsEmpty() {
		gui.git.Patch.PatchManager.Reset()
	}

	if gui.currentContext().GetKey() == gui.State.Contexts.PatchBuilding.GetKey() {
		return gui.c.PushContext(gui.State.Contexts.CommitFiles)
	} else {
		// need to re-focus in case the secondary view should now be hidden
		return gui.currentContext().HandleFocus()
	}
}

func (gui *Gui) secondaryPatchPanelUpdateOpts() *viewUpdateOpts {
	if gui.git.Patch.PatchManager.Active() {
		patch := gui.git.Patch.PatchManager.RenderAggregatedPatchColored(false)

		return &viewUpdateOpts{
			title:     "Custom Patch",
			noWrap:    true,
			highlight: true,
			context:   gui.State.Contexts.PatchBuilding,
			task:      NewRenderStringWithoutScrollTask(patch),
		}
	}

	return nil
}
