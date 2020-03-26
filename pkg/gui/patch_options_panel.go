package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

func (gui *Gui) handleCreatePatchOptionsMenu(g *gocui.Gui, v *gocui.View) error {
	if !gui.GitCommand.PatchManager.CommitSelected() {
		return gui.createErrorPanel(gui.g, gui.Tr.SLocalize("NoPatchError"))
	}

	menuItems := []*menuItem{
		{
			displayString: fmt.Sprintf("remove patch from original commit (%s)", gui.GitCommand.PatchManager.CommitSha),
			onPress:       gui.handleDeletePatchFromCommit,
		},
		{
			displayString: "pull patch out into index",
			onPress:       gui.handlePullPatchIntoWorkingTree,
		},
		{
			displayString: "apply patch",
			onPress:       gui.handleApplyPatch,
		},
		{
			displayString: "reset patch",
			onPress:       gui.handleResetPatch,
		},
	}

	selectedCommit := gui.getSelectedCommit(gui.g)
	if selectedCommit != nil && gui.GitCommand.PatchManager.CommitSha != selectedCommit.Sha {
		// adding this option to index 1
		menuItems = append(
			menuItems[:1],
			append(
				[]*menuItem{
					{
						displayString: fmt.Sprintf("move patch to selected commit (%s)", selectedCommit.Sha),
						onPress:       gui.handleMovePatchToSelectedCommit,
					},
				}, menuItems[1:]...,
			)...,
		)
	}

	return gui.createMenu(gui.Tr.SLocalize("PatchOptionsTitle"), menuItems, createMenuOptions{showCancel: true})
}

func (gui *Gui) getPatchCommitIndex() int {
	for index, commit := range gui.State.Commits {
		if commit.Sha == gui.GitCommand.PatchManager.CommitSha {
			return index
		}
	}
	return -1
}

func (gui *Gui) validateNormalWorkingTreeState() (bool, error) {
	if gui.State.WorkingTreeState != "normal" {
		return false, gui.createErrorPanel(gui.g, gui.Tr.SLocalize("CantPatchWhileRebasingError"))
	}
	return true, nil
}

func (gui *Gui) returnFocusFromLineByLinePanelIfNecessary() error {
	if gui.State.MainContext == "patch-building" {
		return gui.handleEscapePatchBuildingPanel(gui.g, nil)
	}
	return nil
}

func (gui *Gui) handleDeletePatchFromCommit() error {
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

	if err := gui.returnFocusFromLineByLinePanelIfNecessary(); err != nil {
		return err
	}

	return gui.WithWaitingStatus(gui.Tr.SLocalize("RebasingStatus"), func() error {
		commitIndex := gui.getPatchCommitIndex()
		err := gui.GitCommand.DeletePatchesFromCommit(gui.State.Commits, commitIndex, gui.GitCommand.PatchManager)
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleMovePatchToSelectedCommit() error {
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

	if err := gui.returnFocusFromLineByLinePanelIfNecessary(); err != nil {
		return err
	}

	return gui.WithWaitingStatus(gui.Tr.SLocalize("RebasingStatus"), func() error {
		commitIndex := gui.getPatchCommitIndex()
		err := gui.GitCommand.MovePatchToSelectedCommit(gui.State.Commits, commitIndex, gui.State.Panels.Commits.SelectedLine, gui.GitCommand.PatchManager)
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handlePullPatchIntoWorkingTree() error {
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

	if err := gui.returnFocusFromLineByLinePanelIfNecessary(); err != nil {
		return err
	}

	return gui.WithWaitingStatus(gui.Tr.SLocalize("RebasingStatus"), func() error {
		commitIndex := gui.getPatchCommitIndex()
		err := gui.GitCommand.PullPatchIntoIndex(gui.State.Commits, commitIndex, gui.GitCommand.PatchManager)
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleApplyPatch() error {
	if err := gui.returnFocusFromLineByLinePanelIfNecessary(); err != nil {
		return err
	}

	if err := gui.GitCommand.PatchManager.ApplyPatches(false); err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
	}
	return gui.refreshSidePanels(gui.g)
}

func (gui *Gui) handleResetPatch() error {
	gui.GitCommand.PatchManager.Reset()
	return gui.refreshCommitFilesView()
}
