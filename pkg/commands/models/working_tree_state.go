package models

import "github.com/jesseduffield/lazygit/pkg/i18n"

// The state of the working tree. Several of these can be true at once.
// In particular, the concrete multi-state combinations that can occur in
// practice are Rebasing+CherryPicking, Rebasing+Reverting, and
// Rebasing+ApplyingPatches. Theoretically, I guess Rebasing+Merging could also
// happen, but it probably won't in practice.
type WorkingTreeState struct {
	Rebasing      bool
	Merging       bool
	CherryPicking bool
	Reverting     bool
	// git am
	ApplyingPatches bool
}

func (self WorkingTreeState) Any() bool {
	return self.Rebasing || self.Merging || self.CherryPicking || self.Reverting || self.ApplyingPatches
}

func (self WorkingTreeState) None() bool {
	return !self.Any()
}

type EffectiveWorkingTreeState int

const (
	// this means we're neither rebasing nor merging, cherry-picking, or reverting
	WORKING_TREE_STATE_NONE EffectiveWorkingTreeState = iota
	WORKING_TREE_STATE_REBASING
	WORKING_TREE_STATE_MERGING
	WORKING_TREE_STATE_CHERRY_PICKING
	WORKING_TREE_STATE_REVERTING
	WORKING_TREE_STATE_APPLYING_PATCHES
)

// Effective returns the "current" state; if several states are true at once,
// this is the one that should be displayed in status views, and it's the one
// that the user can continue or abort.
//
// As an example, if you are stopped in an interactive rebase, and then you
// perform a cherry-pick, and the cherry-pick conflicts, then both
// WorkingTreeState.Rebasing and WorkingTreeState.CherryPicking are true.
// The effective state is cherry-picking, because that's the one you can
// continue or abort. It is not possible to continue the rebase without first
// aborting the cherry-pick.
func (self WorkingTreeState) Effective() EffectiveWorkingTreeState {
	if self.Reverting {
		return WORKING_TREE_STATE_REVERTING
	}
	if self.CherryPicking {
		return WORKING_TREE_STATE_CHERRY_PICKING
	}
	if self.Merging {
		return WORKING_TREE_STATE_MERGING
	}
	if self.ApplyingPatches {
		return WORKING_TREE_STATE_APPLYING_PATCHES
	}
	if self.Rebasing {
		return WORKING_TREE_STATE_REBASING
	}
	return WORKING_TREE_STATE_NONE
}

func (self WorkingTreeState) Title(tr *i18n.TranslationSet) string {
	return map[EffectiveWorkingTreeState]string{
		WORKING_TREE_STATE_REBASING:         tr.RebasingStatus,
		WORKING_TREE_STATE_MERGING:          tr.MergingStatus,
		WORKING_TREE_STATE_CHERRY_PICKING:   tr.CherryPickingStatus,
		WORKING_TREE_STATE_REVERTING:        tr.RevertingStatus,
		WORKING_TREE_STATE_APPLYING_PATCHES: tr.ApplyingPatchesStatus,
	}[self.Effective()]
}

func (self WorkingTreeState) LowerCaseTitle(tr *i18n.TranslationSet) string {
	return map[EffectiveWorkingTreeState]string{
		WORKING_TREE_STATE_REBASING:         tr.LowercaseRebasingStatus,
		WORKING_TREE_STATE_MERGING:          tr.LowercaseMergingStatus,
		WORKING_TREE_STATE_CHERRY_PICKING:   tr.LowercaseCherryPickingStatus,
		WORKING_TREE_STATE_REVERTING:        tr.LowercaseRevertingStatus,
		WORKING_TREE_STATE_APPLYING_PATCHES: tr.LowercaseApplyingPatchesStatus,
	}[self.Effective()]
}

func (self WorkingTreeState) OptionsMenuTitle(tr *i18n.TranslationSet) string {
	return map[EffectiveWorkingTreeState]string{
		WORKING_TREE_STATE_REBASING:         tr.RebaseOptionsTitle,
		WORKING_TREE_STATE_MERGING:          tr.MergeOptionsTitle,
		WORKING_TREE_STATE_CHERRY_PICKING:   tr.CherryPickOptionsTitle,
		WORKING_TREE_STATE_REVERTING:        tr.RevertOptionsTitle,
		WORKING_TREE_STATE_APPLYING_PATCHES: tr.ApplyPatchesOptionsTitle,
	}[self.Effective()]
}

func (self WorkingTreeState) OptionsMapTitle(tr *i18n.TranslationSet) string {
	return map[EffectiveWorkingTreeState]string{
		WORKING_TREE_STATE_REBASING:         tr.ViewRebaseOptions,
		WORKING_TREE_STATE_MERGING:          tr.ViewMergeOptions,
		WORKING_TREE_STATE_CHERRY_PICKING:   tr.ViewCherryPickOptions,
		WORKING_TREE_STATE_REVERTING:        tr.ViewRevertOptions,
		WORKING_TREE_STATE_APPLYING_PATCHES: tr.ViewApplyPatchesOptions,
	}[self.Effective()]
}

// OperationName returns a phrase for the current operation that we can use in
// sentences like "Abort <operation>".
func (self WorkingTreeState) OperationName(tr *i18n.TranslationSet) string {
	return map[EffectiveWorkingTreeState]string{
		WORKING_TREE_STATE_REBASING:         tr.RebaseOperationName,
		WORKING_TREE_STATE_MERGING:          tr.MergeOperationName,
		WORKING_TREE_STATE_CHERRY_PICKING:   tr.CherryPickOperationName,
		WORKING_TREE_STATE_REVERTING:        tr.RevertOperationName,
		WORKING_TREE_STATE_APPLYING_PATCHES: tr.ApplyPatchesOperationName,
	}[self.Effective()]
}

func (self WorkingTreeState) CommandName() string {
	return map[EffectiveWorkingTreeState]string{
		WORKING_TREE_STATE_REBASING:         "rebase",
		WORKING_TREE_STATE_MERGING:          "merge",
		WORKING_TREE_STATE_CHERRY_PICKING:   "cherry-pick",
		WORKING_TREE_STATE_REVERTING:        "revert",
		WORKING_TREE_STATE_APPLYING_PATCHES: "am",
	}[self.Effective()]
}

func (self WorkingTreeState) CanShowTodos() bool {
	return self.Rebasing || self.CherryPicking || self.Reverting || self.ApplyingPatches
}

func (self WorkingTreeState) CanSkip() bool {
	return self.Rebasing || self.CherryPicking || self.Reverting || self.ApplyingPatches
}
