package controllers

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CustomPatchOptionsMenuAction struct {
	c *ControllerCommon
}

func (self *CustomPatchOptionsMenuAction) Call() error {
	if !self.c.Git().Patch.PatchBuilder.Active() {
		return self.c.ErrorMsg(self.c.Tr.NoPatchError)
	}

	if self.c.Git().Patch.PatchBuilder.IsEmpty() {
		return self.c.ErrorMsg(self.c.Tr.EmptyPatchError)
	}

	menuItems := []*types.MenuItem{
		{
			Label:   self.c.Tr.ResetPatch,
			Tooltip: self.c.Tr.ResetPatchTooltip,
			OnPress: self.c.Helpers().PatchBuilding.Reset,
			Key:     'c',
		},
		{
			Label:   self.c.Tr.ApplyPatch,
			Tooltip: self.c.Tr.ApplyPatchTooltip,
			OnPress: func() error { return self.handleApplyPatch(false) },
			Key:     'a',
		},
		{
			Label:   self.c.Tr.ApplyPatchInReverse,
			Tooltip: self.c.Tr.ApplyPatchInReverseTooltip,
			OnPress: func() error { return self.handleApplyPatch(true) },
			Key:     'r',
		},
	}

	if self.c.Git().Patch.PatchBuilder.CanRebase && self.c.Git().Status.WorkingTreeState() == enums.REBASE_MODE_NONE {
		menuItems = append(menuItems, []*types.MenuItem{
			{
				Label:   fmt.Sprintf(self.c.Tr.RemovePatchFromOriginalCommit, self.c.Git().Patch.PatchBuilder.To),
				Tooltip: self.c.Tr.RemovePatchFromOriginalCommitTooltip,
				OnPress: self.handleDeletePatchFromCommit,
				Key:     'd',
			},
			{
				Label:   self.c.Tr.MovePatchOutIntoIndex,
				Tooltip: self.c.Tr.MovePatchOutIntoIndexTooltip,
				OnPress: self.handleMovePatchIntoWorkingTree,
				Key:     'i',
			},
			{
				Label:   self.c.Tr.MovePatchIntoNewCommit,
				Tooltip: self.c.Tr.MovePatchIntoNewCommitTooltip,
				OnPress: self.handlePullPatchIntoNewCommit,
				Key:     'n',
			},
		}...)

		if self.c.CurrentContext().GetKey() == self.c.Contexts().LocalCommits.GetKey() {
			selectedCommit := self.c.Contexts().LocalCommits.GetSelected()
			if selectedCommit != nil && self.c.Git().Patch.PatchBuilder.To != selectedCommit.Sha {

				var disabledReason *types.DisabledReason
				if self.c.Contexts().LocalCommits.AreMultipleItemsSelected() {
					disabledReason = &types.DisabledReason{Text: self.c.Tr.RangeSelectNotSupported}
				}

				// adding this option to index 1
				menuItems = append(
					menuItems[:1],
					append(
						[]*types.MenuItem{
							{
								Label:          fmt.Sprintf(self.c.Tr.MovePatchToSelectedCommit, selectedCommit.Sha),
								Tooltip:        self.c.Tr.MovePatchToSelectedCommitTooltip,
								OnPress:        self.handleMovePatchToSelectedCommit,
								Key:            'm',
								DisabledReason: disabledReason,
							},
						}, menuItems[1:]...,
					)...,
				)
			}
		}
	}

	menuItems = append(menuItems, []*types.MenuItem{
		{
			Label:   self.c.Tr.CopyPatchToClipboard,
			OnPress: func() error { return self.copyPatchToClipboard() },
			Key:     'y',
		},
	}...)

	return self.c.Menu(types.CreateMenuOptions{Title: self.c.Tr.PatchOptionsTitle, Items: menuItems})
}

func (self *CustomPatchOptionsMenuAction) getPatchCommitIndex() int {
	for index, commit := range self.c.Model().Commits {
		if commit.Sha == self.c.Git().Patch.PatchBuilder.To {
			return index
		}
	}
	return -1
}

func (self *CustomPatchOptionsMenuAction) validateNormalWorkingTreeState() (bool, error) {
	if self.c.Git().Status.WorkingTreeState() != enums.REBASE_MODE_NONE {
		return false, self.c.ErrorMsg(self.c.Tr.CantPatchWhileRebasingError)
	}
	return true, nil
}

func (self *CustomPatchOptionsMenuAction) returnFocusFromPatchExplorerIfNecessary() error {
	if self.c.CurrentContext().GetKey() == self.c.Contexts().CustomPatchBuilder.GetKey() {
		return self.c.Helpers().PatchBuilding.Escape()
	}
	return nil
}

func (self *CustomPatchOptionsMenuAction) handleDeletePatchFromCommit() error {
	if ok, err := self.validateNormalWorkingTreeState(); !ok {
		return err
	}

	if err := self.returnFocusFromPatchExplorerIfNecessary(); err != nil {
		return err
	}

	return self.c.WithWaitingStatus(self.c.Tr.RebasingStatus, func(gocui.Task) error {
		commitIndex := self.getPatchCommitIndex()
		self.c.LogAction(self.c.Tr.Actions.RemovePatchFromCommit)
		err := self.c.Git().Patch.DeletePatchesFromCommit(self.c.Model().Commits, commitIndex)
		return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
	})
}

func (self *CustomPatchOptionsMenuAction) handleMovePatchToSelectedCommit() error {
	if ok, err := self.validateNormalWorkingTreeState(); !ok {
		return err
	}

	if err := self.returnFocusFromPatchExplorerIfNecessary(); err != nil {
		return err
	}

	return self.c.WithWaitingStatus(self.c.Tr.RebasingStatus, func(gocui.Task) error {
		commitIndex := self.getPatchCommitIndex()
		self.c.LogAction(self.c.Tr.Actions.MovePatchToSelectedCommit)
		err := self.c.Git().Patch.MovePatchToSelectedCommit(self.c.Model().Commits, commitIndex, self.c.Contexts().LocalCommits.GetSelectedLineIdx())
		return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
	})
}

func (self *CustomPatchOptionsMenuAction) handleMovePatchIntoWorkingTree() error {
	if ok, err := self.validateNormalWorkingTreeState(); !ok {
		return err
	}

	if err := self.returnFocusFromPatchExplorerIfNecessary(); err != nil {
		return err
	}

	pull := func(stash bool) error {
		return self.c.WithWaitingStatus(self.c.Tr.RebasingStatus, func(gocui.Task) error {
			commitIndex := self.getPatchCommitIndex()
			self.c.LogAction(self.c.Tr.Actions.MovePatchIntoIndex)
			err := self.c.Git().Patch.MovePatchIntoIndex(self.c.Model().Commits, commitIndex, stash)
			return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
		})
	}

	if self.c.Helpers().WorkingTree.IsWorkingTreeDirty() {
		return self.c.Confirm(types.ConfirmOpts{
			Title:  self.c.Tr.MustStashTitle,
			Prompt: self.c.Tr.MustStashWarning,
			HandleConfirm: func() error {
				return pull(true)
			},
		})
	} else {
		return pull(false)
	}
}

func (self *CustomPatchOptionsMenuAction) handlePullPatchIntoNewCommit() error {
	if ok, err := self.validateNormalWorkingTreeState(); !ok {
		return err
	}

	if err := self.returnFocusFromPatchExplorerIfNecessary(); err != nil {
		return err
	}

	commitIndex := self.getPatchCommitIndex()
	return self.c.Helpers().Commits.OpenCommitMessagePanel(
		&helpers.OpenCommitMessagePanelOpts{
			// Pass a commit index of one less than the moved-from commit, so that
			// you can press up arrow once to recall the original commit message:
			CommitIndex:      commitIndex - 1,
			InitialMessage:   "",
			SummaryTitle:     self.c.Tr.CommitSummaryTitle,
			DescriptionTitle: self.c.Tr.CommitDescriptionTitle,
			PreserveMessage:  false,
			OnConfirm: func(summary string, description string) error {
				return self.c.WithWaitingStatus(self.c.Tr.RebasingStatus, func(gocui.Task) error {
					_ = self.c.Helpers().Commits.PopCommitMessageContexts()
					self.c.LogAction(self.c.Tr.Actions.MovePatchIntoNewCommit)
					err := self.c.Git().Patch.PullPatchIntoNewCommit(self.c.Model().Commits, commitIndex, summary, description)
					return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
				})
			},
		},
	)
}

func (self *CustomPatchOptionsMenuAction) handleApplyPatch(reverse bool) error {
	if err := self.returnFocusFromPatchExplorerIfNecessary(); err != nil {
		return err
	}

	action := self.c.Tr.Actions.ApplyPatch
	if reverse {
		action = "Apply patch in reverse"
	}
	self.c.LogAction(action)
	if err := self.c.Git().Patch.ApplyCustomPatch(reverse); err != nil {
		return self.c.Error(err)
	}
	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (self *CustomPatchOptionsMenuAction) copyPatchToClipboard() error {
	patch := self.c.Git().Patch.PatchBuilder.RenderAggregatedPatch(true)

	self.c.LogAction(self.c.Tr.Actions.CopyPatchToClipboard)
	if err := self.c.OS().CopyToClipboard(patch); err != nil {
		return self.c.Error(err)
	}

	self.c.Toast(self.c.Tr.PatchCopiedToClipboard)

	return nil
}
