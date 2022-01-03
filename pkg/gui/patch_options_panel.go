package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
)

func (gui *Gui) handleCreatePatchOptionsMenu() error {
	if !gui.GitCommand.PatchManager.Active() {
		return gui.createErrorPanel(gui.Tr.NoPatchError)
	}

	menuItems := []*menuItem{
		{
			displayString: "reset patch",
			onPress:       gui.handleResetPatch,
		},
		{
			displayString: "apply patch",
			onPress:       func() error { return gui.handleApplyPatch(false) },
		},
		{
			displayString: "apply patch in reverse",
			onPress:       func() error { return gui.handleApplyPatch(true) },
		},
	}

	if gui.GitCommand.PatchManager.CanRebase && gui.GitCommand.WorkingTreeState() == enums.REBASE_MODE_NONE {
		menuItems = append(menuItems, []*menuItem{
			{
				displayString: fmt.Sprintf("remove patch from original commit (%s)", gui.GitCommand.PatchManager.To),
				onPress:       gui.handleDeletePatchFromCommit,
			},
			{
				displayString: "move patch out into index",
				onPress:       gui.handleMovePatchIntoWorkingTree,
			},
			{
				displayString: "move patch into new commit",
				onPress:       gui.handlePullPatchIntoNewCommit,
			},
		}...)

		if gui.currentContext().GetKey() == gui.State.Contexts.BranchCommits.GetKey() {
			selectedCommit := gui.getSelectedLocalCommit()
			if selectedCommit != nil && gui.GitCommand.PatchManager.To != selectedCommit.Sha {
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
		}
	}

	return gui.createMenu(gui.Tr.PatchOptionsTitle, menuItems, createMenuOptions{showCancel: true})
}

func (gui *Gui) getPatchCommitIndex() int {
	for index, commit := range gui.State.Commits {
		if commit.Sha == gui.GitCommand.PatchManager.To {
			return index
		}
	}
	return -1
}

func (gui *Gui) validateNormalWorkingTreeState() (bool, error) {
	if gui.GitCommand.WorkingTreeState() != enums.REBASE_MODE_NONE {
		return false, gui.createErrorPanel(gui.Tr.CantPatchWhileRebasingError)
	}
	return true, nil
}

func (gui *Gui) returnFocusFromLineByLinePanelIfNecessary() error {
	if gui.State.MainContext == MAIN_PATCH_BUILDING_CONTEXT_KEY {
		return gui.handleEscapePatchBuildingPanel()
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

	return gui.WithWaitingStatus(gui.Tr.RebasingStatus, func() error {
		commitIndex := gui.getPatchCommitIndex()
		err := gui.GitCommand.WithSpan(gui.Tr.Spans.RemovePatchFromCommit).DeletePatchesFromCommit(gui.State.Commits, commitIndex, gui.GitCommand.PatchManager)
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

	return gui.WithWaitingStatus(gui.Tr.RebasingStatus, func() error {
		commitIndex := gui.getPatchCommitIndex()
		err := gui.GitCommand.WithSpan(gui.Tr.Spans.MovePatchToSelectedCommit).MovePatchToSelectedCommit(gui.State.Commits, commitIndex, gui.State.Panels.Commits.SelectedLineIdx, gui.GitCommand.PatchManager)
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleMovePatchIntoWorkingTree() error {
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

	if err := gui.returnFocusFromLineByLinePanelIfNecessary(); err != nil {
		return err
	}

	pull := func(stash bool) error {
		return gui.WithWaitingStatus(gui.Tr.RebasingStatus, func() error {
			commitIndex := gui.getPatchCommitIndex()
			err := gui.GitCommand.WithSpan(gui.Tr.Spans.MovePatchIntoIndex).MovePatchIntoIndex(gui.State.Commits, commitIndex, gui.GitCommand.PatchManager, stash)
			return gui.handleGenericMergeCommandResult(err)
		})
	}

	if len(gui.trackedFiles()) > 0 {
		return gui.ask(askOpts{
			title:  gui.Tr.MustStashTitle,
			prompt: gui.Tr.MustStashWarning,
			handleConfirm: func() error {
				return pull(true)
			},
		})
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

	return gui.WithWaitingStatus(gui.Tr.RebasingStatus, func() error {
		commitIndex := gui.getPatchCommitIndex()
		err := gui.GitCommand.WithSpan(gui.Tr.Spans.MovePatchIntoNewCommit).PullPatchIntoNewCommit(gui.State.Commits, commitIndex, gui.GitCommand.PatchManager)
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleApplyPatch(reverse bool) error {
	if err := gui.returnFocusFromLineByLinePanelIfNecessary(); err != nil {
		return err
	}

	span := gui.Tr.Spans.ApplyPatch
	if reverse {
		span = "Apply patch in reverse"
	}
	if err := gui.GitCommand.WithSpan(span).PatchManager.ApplyPatches(reverse); err != nil {
		return gui.surfaceError(err)
	}
	return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
}

func (gui *Gui) handleResetPatch() error {
	gui.GitCommand.PatchManager.Reset()
	if gui.currentContextKeyIgnoringPopups() == MAIN_PATCH_BUILDING_CONTEXT_KEY {
		if err := gui.pushContext(gui.State.Contexts.CommitFiles); err != nil {
			return err
		}
	}
	return gui.refreshCommitFilesView()
}
