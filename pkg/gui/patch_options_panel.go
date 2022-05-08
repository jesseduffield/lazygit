package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) handleCreatePatchOptionsMenu() error {
	if !gui.git.Patch.PatchManager.Active() {
		return gui.c.ErrorMsg(gui.c.Tr.NoPatchError)
	}

	menuItems := []*types.MenuItem{
		{
			Label:   "reset patch",
			OnPress: gui.handleResetPatch,
			Key:     'c',
		},
		{
			Label:   "apply patch",
			OnPress: func() error { return gui.handleApplyPatch(false) },
			Key:     'a',
		},
		{
			Label:   "apply patch in reverse",
			OnPress: func() error { return gui.handleApplyPatch(true) },
			Key:     'r',
		},
	}

	if gui.git.Patch.PatchManager.CanRebase && gui.git.Status.WorkingTreeState() == enums.REBASE_MODE_NONE {
		menuItems = append(menuItems, []*types.MenuItem{
			{
				Label:   fmt.Sprintf("remove patch from original commit (%s)", gui.git.Patch.PatchManager.To),
				OnPress: gui.handleDeletePatchFromCommit,
				Key:     'd',
			},
			{
				Label:   "move patch out into index",
				OnPress: gui.handleMovePatchIntoWorkingTree,
				Key:     'i',
			},
			{
				Label:   "move patch into new commit",
				OnPress: gui.handlePullPatchIntoNewCommit,
				Key:     'n',
			},
		}...)

		if gui.currentContext().GetKey() == gui.State.Contexts.LocalCommits.GetKey() {
			selectedCommit := gui.getSelectedLocalCommit()
			if selectedCommit != nil && gui.git.Patch.PatchManager.To != selectedCommit.Sha {
				// adding this option to index 1
				menuItems = append(
					menuItems[:1],
					append(
						[]*types.MenuItem{
							{
								Label:   fmt.Sprintf("move patch to selected commit (%s)", selectedCommit.Sha),
								OnPress: gui.handleMovePatchToSelectedCommit,
								Key:     'm',
							},
						}, menuItems[1:]...,
					)...,
				)
			}
		}
	}

	return gui.c.Menu(types.CreateMenuOptions{Title: gui.c.Tr.PatchOptionsTitle, Items: menuItems})
}

func (gui *Gui) getPatchCommitIndex() int {
	for index, commit := range gui.State.Model.Commits {
		if commit.Sha == gui.git.Patch.PatchManager.To {
			return index
		}
	}
	return -1
}

func (gui *Gui) validateNormalWorkingTreeState() (bool, error) {
	if gui.git.Status.WorkingTreeState() != enums.REBASE_MODE_NONE {
		return false, gui.c.ErrorMsg(gui.c.Tr.CantPatchWhileRebasingError)
	}
	return true, nil
}

func (gui *Gui) returnFocusFromLineByLinePanelIfNecessary() error {
	if gui.State.MainContext == context.MAIN_PATCH_BUILDING_CONTEXT_KEY {
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

	return gui.c.WithWaitingStatus(gui.c.Tr.RebasingStatus, func() error {
		commitIndex := gui.getPatchCommitIndex()
		gui.c.LogAction(gui.c.Tr.Actions.RemovePatchFromCommit)
		err := gui.git.Patch.DeletePatchesFromCommit(gui.State.Model.Commits, commitIndex)
		return gui.helpers.MergeAndRebase.CheckMergeOrRebase(err)
	})
}

func (gui *Gui) handleMovePatchToSelectedCommit() error {
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

	if err := gui.returnFocusFromLineByLinePanelIfNecessary(); err != nil {
		return err
	}

	return gui.c.WithWaitingStatus(gui.c.Tr.RebasingStatus, func() error {
		commitIndex := gui.getPatchCommitIndex()
		gui.c.LogAction(gui.c.Tr.Actions.MovePatchToSelectedCommit)
		err := gui.git.Patch.MovePatchToSelectedCommit(gui.State.Model.Commits, commitIndex, gui.State.Contexts.LocalCommits.GetSelectedLineIdx())
		return gui.helpers.MergeAndRebase.CheckMergeOrRebase(err)
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
		return gui.c.WithWaitingStatus(gui.c.Tr.RebasingStatus, func() error {
			commitIndex := gui.getPatchCommitIndex()
			gui.c.LogAction(gui.c.Tr.Actions.MovePatchIntoIndex)
			err := gui.git.Patch.MovePatchIntoIndex(gui.State.Model.Commits, commitIndex, stash)
			return gui.helpers.MergeAndRebase.CheckMergeOrRebase(err)
		})
	}

	if gui.helpers.WorkingTree.IsWorkingTreeDirty() {
		return gui.c.Confirm(types.ConfirmOpts{
			Title:  gui.c.Tr.MustStashTitle,
			Prompt: gui.c.Tr.MustStashWarning,
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

	return gui.c.WithWaitingStatus(gui.c.Tr.RebasingStatus, func() error {
		commitIndex := gui.getPatchCommitIndex()
		gui.c.LogAction(gui.c.Tr.Actions.MovePatchIntoNewCommit)
		err := gui.git.Patch.PullPatchIntoNewCommit(gui.State.Model.Commits, commitIndex)
		return gui.helpers.MergeAndRebase.CheckMergeOrRebase(err)
	})
}

func (gui *Gui) handleApplyPatch(reverse bool) error {
	if err := gui.returnFocusFromLineByLinePanelIfNecessary(); err != nil {
		return err
	}

	action := gui.c.Tr.Actions.ApplyPatch
	if reverse {
		action = "Apply patch in reverse"
	}
	gui.c.LogAction(action)
	if err := gui.git.Patch.PatchManager.ApplyPatches(reverse); err != nil {
		return gui.c.Error(err)
	}
	return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (gui *Gui) handleResetPatch() error {
	gui.git.Patch.PatchManager.Reset()
	if gui.currentContextKeyIgnoringPopups() == context.MAIN_PATCH_BUILDING_CONTEXT_KEY {
		if err := gui.c.PushContext(gui.State.Contexts.CommitFiles); err != nil {
			return err
		}
	}

	if err := gui.handleEscapePatchBuildingPanel(); err != nil {
		return err
	}

	return gui.refreshCommitFilesContext()
}
