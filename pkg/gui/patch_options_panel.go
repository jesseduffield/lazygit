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
	m := gui.State.PatchManager
	if m == nil {
		return gui.createErrorPanel(gui.g, gui.Tr.SLocalize("NoPatchError"))
	}

	options := []*patchMenuOption{
		{displayName: "discard patch", function: gui.handleDeletePatchFromCommit},
		{displayName: "pull patch out into index", function: gui.handlePullPatchIntoWorkingTree},
		{displayName: "save patch to file"},
		{displayName: "clear patch", function: gui.handleClearPatch},
	}

	selectedCommit := gui.getSelectedCommit(gui.g)
	if selectedCommit != nil && gui.State.PatchManager.CommitSha != selectedCommit.Sha {
		options = append(options, &patchMenuOption{
			displayName: fmt.Sprintf("move patch to selected commit (%s)", selectedCommit.Sha),
			function:    gui.handleMovePatchToSelectedCommit,
		})
	}

	handleMenuPress := func(index int) error {
		return options[index].function()
	}

	return gui.createMenu(gui.Tr.SLocalize("PatchOptionsTitle"), options, len(options), handleMenuPress)
}

func (gui *Gui) getPatchCommitIndex() int {
	for index, commit := range gui.State.Commits {
		if commit.Sha == gui.State.PatchManager.CommitSha {
			return index
		}
	}
	return -1
}

func (gui *Gui) handleDeletePatchFromCommit() error {
	// TODO: deal with when we're already rebasing

	return gui.WithWaitingStatus(gui.Tr.SLocalize("RebasingStatus"), func() error {
		commitIndex := gui.getPatchCommitIndex()
		err := gui.GitCommand.DeletePatchesFromCommit(gui.State.Commits, commitIndex, gui.State.PatchManager)
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleMovePatchToSelectedCommit() error {
	// TODO: deal with when we're already rebasing

	return gui.WithWaitingStatus(gui.Tr.SLocalize("RebasingStatus"), func() error {
		commitIndex := gui.getPatchCommitIndex()
		err := gui.GitCommand.MovePatchToSelectedCommit(gui.State.Commits, commitIndex, gui.State.Panels.Commits.SelectedLine, gui.State.PatchManager)
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handlePullPatchIntoWorkingTree() error {
	// TODO: deal with when we're already rebasing

	return gui.WithWaitingStatus(gui.Tr.SLocalize("RebasingStatus"), func() error {
		commitIndex := gui.getPatchCommitIndex()
		err := gui.GitCommand.PullPatchIntoIndex(gui.State.Commits, commitIndex, gui.State.PatchManager)
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleClearPatch() error {
	gui.State.PatchManager = nil
	return gui.refreshCommitFilesView()
}
