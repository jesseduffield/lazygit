package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) handleCreatePatchOptionsMenu() error {
	if !gui.Git.Patch.PatchManager.Active() {
		return gui.PopupHandler.ErrorMsg(gui.Tr.NoPatchError)
	}

	menuItems := []*popup.MenuItem{
		{
			DisplayString: "reset patch",
			OnPress:       gui.handleResetPatch,
		},
		{
			DisplayString: "apply patch",
			OnPress:       func() error { return gui.handleApplyPatch(false) },
		},
		{
			DisplayString: "apply patch in reverse",
			OnPress:       func() error { return gui.handleApplyPatch(true) },
		},
	}

	if gui.Git.Patch.PatchManager.CanRebase && gui.Git.Status.WorkingTreeState() == enums.REBASE_MODE_NONE {
		menuItems = append(menuItems, []*popup.MenuItem{
			{
				DisplayString: fmt.Sprintf("remove patch from original commit (%s)", gui.Git.Patch.PatchManager.To),
				OnPress:       gui.handleDeletePatchFromCommit,
			},
			{
				DisplayString: "move patch out into index",
				OnPress:       gui.handleMovePatchIntoWorkingTree,
			},
			{
				DisplayString: "move patch into new commit",
				OnPress:       gui.handlePullPatchIntoNewCommit,
			},
		}...)

		if gui.currentContext().GetKey() == gui.State.Contexts.BranchCommits.GetKey() {
			selectedCommit := gui.getSelectedLocalCommit()
			if selectedCommit != nil && gui.Git.Patch.PatchManager.To != selectedCommit.Sha {
				// adding this option to index 1
				menuItems = append(
					menuItems[:1],
					append(
						[]*popup.MenuItem{
							{
								DisplayString: fmt.Sprintf("move patch to selected commit (%s)", selectedCommit.Sha),
								OnPress:       gui.handleMovePatchToSelectedCommit,
							},
						}, menuItems[1:]...,
					)...,
				)
			}
		}
	}

	return gui.PopupHandler.Menu(popup.CreateMenuOptions{Title: gui.Tr.PatchOptionsTitle, Items: menuItems})
}

func (gui *Gui) getPatchCommitIndex() int {
	for index, commit := range gui.State.Commits {
		if commit.Sha == gui.Git.Patch.PatchManager.To {
			return index
		}
	}
	return -1
}

func (gui *Gui) validateNormalWorkingTreeState() (bool, error) {
	if gui.Git.Status.WorkingTreeState() != enums.REBASE_MODE_NONE {
		return false, gui.PopupHandler.ErrorMsg(gui.Tr.CantPatchWhileRebasingError)
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

	return gui.PopupHandler.WithWaitingStatus(gui.Tr.RebasingStatus, func() error {
		commitIndex := gui.getPatchCommitIndex()
		gui.logAction(gui.Tr.Actions.RemovePatchFromCommit)
		err := gui.Git.Patch.DeletePatchesFromCommit(gui.State.Commits, commitIndex)
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

	return gui.PopupHandler.WithWaitingStatus(gui.Tr.RebasingStatus, func() error {
		commitIndex := gui.getPatchCommitIndex()
		gui.logAction(gui.Tr.Actions.MovePatchToSelectedCommit)
		err := gui.Git.Patch.MovePatchToSelectedCommit(gui.State.Commits, commitIndex, gui.State.Panels.Commits.SelectedLineIdx)
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
		return gui.PopupHandler.WithWaitingStatus(gui.Tr.RebasingStatus, func() error {
			commitIndex := gui.getPatchCommitIndex()
			gui.logAction(gui.Tr.Actions.MovePatchIntoIndex)
			err := gui.Git.Patch.MovePatchIntoIndex(gui.State.Commits, commitIndex, stash)
			return gui.handleGenericMergeCommandResult(err)
		})
	}

	if len(gui.trackedFiles()) > 0 {
		return gui.PopupHandler.Ask(popup.AskOpts{
			Title:  gui.Tr.MustStashTitle,
			Prompt: gui.Tr.MustStashWarning,
			HandleConfirm: func() error {
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

	return gui.PopupHandler.WithWaitingStatus(gui.Tr.RebasingStatus, func() error {
		commitIndex := gui.getPatchCommitIndex()
		gui.logAction(gui.Tr.Actions.MovePatchIntoNewCommit)
		err := gui.Git.Patch.PullPatchIntoNewCommit(gui.State.Commits, commitIndex)
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleApplyPatch(reverse bool) error {
	if err := gui.returnFocusFromLineByLinePanelIfNecessary(); err != nil {
		return err
	}

	action := gui.Tr.Actions.ApplyPatch
	if reverse {
		action = "Apply patch in reverse"
	}
	gui.logAction(action)
	if err := gui.Git.Patch.PatchManager.ApplyPatches(reverse); err != nil {
		return gui.PopupHandler.Error(err)
	}
	return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC})
}

func (gui *Gui) handleResetPatch() error {
	gui.Git.Patch.PatchManager.Reset()
	if gui.currentContextKeyIgnoringPopups() == MAIN_PATCH_BUILDING_CONTEXT_KEY {
		if err := gui.pushContext(gui.State.Contexts.CommitFiles); err != nil {
			return err
		}
	}
	return gui.refreshCommitFilesView()
}
