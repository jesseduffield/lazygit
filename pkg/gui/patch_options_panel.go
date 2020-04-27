package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

func (gui *Gui) handleCreatePatchOptionsMenu(g *gocui.Gui, v *gocui.View) error {
	if !gui.GitCommand.PatchManager.CommitSelected() {
		return gui.createErrorPanel(gui.Tr.SLocalize("NoPatchError"))
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
			displayString: "pull patch into new commit",
			onPress:       gui.handlePullPatchIntoNewCommit,
		},
		{
			displayString: "apply patch",
			onPress:       func() error { return gui.handleApplyPatch(false) },
		},
		{
			displayString: "apply patch in reverse",
			onPress:       func() error { return gui.handleApplyPatch(true) },
		},
		{
			displayString: "reset patch",
			onPress:       gui.handleResetPatch,
		},
	}

	selectedCommit := gui.getSelectedCommit()
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
	if gui.GitCommand.WorkingTreeState() != "normal" {
		return false, gui.createErrorPanel(gui.Tr.SLocalize("CantPatchWhileRebasingError"))
	}
	if gui.GitCommand.WorkingTreeState() != "normal" {
		return false, gui.createErrorPanel(gui.Tr.SLocalize("CantPatchWhileRebasingError"))
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

	pull := func(stash bool) error {
		return gui.WithWaitingStatus(gui.Tr.SLocalize("RebasingStatus"), func() error {
			commitIndex := gui.getPatchCommitIndex()
			err := gui.GitCommand.PullPatchIntoIndex(gui.State.Commits, commitIndex, gui.GitCommand.PatchManager, stash)
			return gui.handleGenericMergeCommandResult(err)
		})
	}

	if len(gui.trackedFiles()) > 0 {
		return gui.createConfirmationPanel(gui.g, gui.g.CurrentView(), true, gui.Tr.SLocalize("MustStashTitle"), gui.Tr.SLocalize("MustStashWarning"), func(*gocui.Gui, *gocui.View) error {
			return pull(true)
		}, nil)
	} else {
		return pull(false)
	}
}

func (gui *Gui) handlePullPatchIntoNewCommit() error {
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

	if err := gui.returnFocusFromLineByLinePanelIfNecessary(); err != nil {
		return err
	}

	return gui.WithWaitingStatus(gui.Tr.SLocalize("RebasingStatus"), func() error {
		commitIndex := gui.getPatchCommitIndex()
		err := gui.GitCommand.PullPatchIntoNewCommit(gui.State.Commits, commitIndex, gui.GitCommand.PatchManager)
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleApplyPatch(reverse bool) error {
	if err := gui.returnFocusFromLineByLinePanelIfNecessary(); err != nil {
		return err
	}

	if err := gui.GitCommand.PatchManager.ApplyPatches(reverse); err != nil {
		return gui.surfaceError(err)
	}
	return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
}

func (gui *Gui) handleResetPatch() error {
	gui.GitCommand.PatchManager.Reset()
	return gui.refreshCommitFilesView()
}
