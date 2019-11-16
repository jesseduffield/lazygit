package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

type patchMenuOption struct {
	displayName string
	function    func() error
}

// GetDisplayStrings is a function.
func (o *patchMenuOption) GetDisplayStrings(isFocused bool) []string {
	return []string{o.displayName}
}

func (gui *Gui) handleCreatePatchOptionsMenu(g *gocui.Gui, v *gocui.View) error {
	if !gui.GitCommand.PatchManager.CommitSelected() {
		return gui.createErrorPanel(gui.g, gui.Tr.SLocalize("NoPatchError"))
	}

	options := []*patchMenuOption{
		{displayName: fmt.Sprintf("remove patch from original commit (%s)", gui.GitCommand.PatchManager.CommitSha), function: gui.handleDeletePatchFromCommit},
		{displayName: "pull patch out into index", function: gui.handlePullPatchIntoWorkingTree},
		{displayName: "reset patch", function: gui.handleResetPatch},
	}

	selectedCommit := gui.getSelectedCommit(gui.g)
	if selectedCommit != nil && gui.GitCommand.PatchManager.CommitSha != selectedCommit.Sha {
		// adding this option to index 1
		options = append(
			options[:1],
			append(
				[]*patchMenuOption{
					{
						displayName: fmt.Sprintf("move patch to selected commit (%s)", selectedCommit.Sha),
						function:    gui.handleMovePatchToSelectedCommit,
					},
				}, options[1:]...,
			)...,
		)
	}

	handleMenuPress := func(index int) error {
		return options[index].function()
	}

	return gui.createMenu(gui.Tr.SLocalize("PatchOptionsTitle"), options, len(options), handleMenuPress)
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

func (gui *Gui) handleResetPatch() error {
	gui.GitCommand.PatchManager.Reset()
	return gui.refreshCommitFilesView()
}
