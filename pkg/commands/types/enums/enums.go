package enums

import "github.com/jesseduffield/lazygit/pkg/i18n"

type WorkingTreeState int

// Note: this is a bitfield, several of these can be set at once
const (
	// this means we're neither rebasing nor merging, cherry-picking, or reverting
	WORKING_TREE_STATE_NONE     WorkingTreeState = 0
	WORKING_TREE_STATE_REBASING                  = 1 << iota
	WORKING_TREE_STATE_MERGING
	WORKING_TREE_STATE_CHERRY_PICKING
	WORKING_TREE_STATE_REVERTING
)

// Effective returns the "current" state; if several states are true at once,
// this is the one that should be displayed in status views, and it's the one
// that the user can continue or abort.
//
// As an example, if you are stopped in an interactive rebase, and then you
// perform a cherry-pick, and the cherry-pick conflicts, then both
// WORKING_TREE_STATE_REBASING and WORKING_TREE_STATE_CHERRY_PICKING are true.
// The effective state is cherry-picking, because that's the one you can
// continue or abort. It is not possible to continue the rebase without first
// aborting the cherry-pick.
func (self WorkingTreeState) Effective() WorkingTreeState {
	if self.IsReverting() {
		return WORKING_TREE_STATE_REVERTING
	}
	if self.IsCherryPicking() {
		return WORKING_TREE_STATE_CHERRY_PICKING
	}
	if self.IsMerging() {
		return WORKING_TREE_STATE_MERGING
	}
	if self.IsRebasing() {
		return WORKING_TREE_STATE_REBASING
	}
	return WORKING_TREE_STATE_NONE
}

func (self WorkingTreeState) IsMerging() bool {
	return (self & WORKING_TREE_STATE_MERGING) != 0
}

func (self WorkingTreeState) IsRebasing() bool {
	return (self & WORKING_TREE_STATE_REBASING) != 0
}

func (self WorkingTreeState) IsCherryPicking() bool {
	return (self & WORKING_TREE_STATE_CHERRY_PICKING) != 0
}

func (self WorkingTreeState) IsReverting() bool {
	return (self & WORKING_TREE_STATE_REVERTING) != 0
}

func (self WorkingTreeState) Title(tr *i18n.TranslationSet) string {
	return map[WorkingTreeState]string{
		WORKING_TREE_STATE_REBASING:       tr.RebasingStatus,
		WORKING_TREE_STATE_MERGING:        tr.MergingStatus,
		WORKING_TREE_STATE_CHERRY_PICKING: tr.CherryPickingStatus,
		WORKING_TREE_STATE_REVERTING:      tr.RevertingStatus,
	}[self.Effective()]
}

func (self WorkingTreeState) LowerCaseTitle(tr *i18n.TranslationSet) string {
	return map[WorkingTreeState]string{
		WORKING_TREE_STATE_REBASING:       tr.LowercaseRebasingStatus,
		WORKING_TREE_STATE_MERGING:        tr.LowercaseMergingStatus,
		WORKING_TREE_STATE_CHERRY_PICKING: tr.LowercaseCherryPickingStatus,
		WORKING_TREE_STATE_REVERTING:      tr.LowercaseRevertingStatus,
	}[self.Effective()]
}

func (self WorkingTreeState) OptionsMenuTitle(tr *i18n.TranslationSet) string {
	return map[WorkingTreeState]string{
		WORKING_TREE_STATE_REBASING:       tr.RebaseOptionsTitle,
		WORKING_TREE_STATE_MERGING:        tr.MergeOptionsTitle,
		WORKING_TREE_STATE_CHERRY_PICKING: tr.CherryPickOptionsTitle,
		WORKING_TREE_STATE_REVERTING:      tr.RevertOptionsTitle,
	}[self.Effective()]
}

func (self WorkingTreeState) OptionsMapTitle(tr *i18n.TranslationSet) string {
	return map[WorkingTreeState]string{
		WORKING_TREE_STATE_REBASING:       tr.ViewRebaseOptions,
		WORKING_TREE_STATE_MERGING:        tr.ViewMergeOptions,
		WORKING_TREE_STATE_CHERRY_PICKING: tr.ViewCherryPickOptions,
		WORKING_TREE_STATE_REVERTING:      tr.ViewRevertOptions,
	}[self.Effective()]
}

func (self WorkingTreeState) CommandName() string {
	return map[WorkingTreeState]string{
		WORKING_TREE_STATE_REBASING:       "rebase",
		WORKING_TREE_STATE_MERGING:        "merge",
		WORKING_TREE_STATE_CHERRY_PICKING: "cherry-pick",
		WORKING_TREE_STATE_REVERTING:      "revert",
	}[self.Effective()]
}

func (self WorkingTreeState) CanShowTodos() bool {
	return self.IsRebasing() || self.IsCherryPicking() || self.IsReverting()
}
