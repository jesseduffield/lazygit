package models

import "github.com/jesseduffield/lazygit/pkg/i18n"

type WorkingTreeState int

const (
	// this means we're neither rebasing nor merging
	WORKING_TREE_STATE_NONE WorkingTreeState = iota
	WORKING_TREE_STATE_REBASING
	WORKING_TREE_STATE_MERGING
)

func (self WorkingTreeState) IsMerging() bool {
	return self == WORKING_TREE_STATE_MERGING
}

func (self WorkingTreeState) IsRebasing() bool {
	return self == WORKING_TREE_STATE_REBASING
}

func (self WorkingTreeState) Title(tr *i18n.TranslationSet) string {
	switch self {
	case WORKING_TREE_STATE_REBASING:
		return tr.RebasingStatus
	case WORKING_TREE_STATE_MERGING:
		return tr.MergingStatus
	default:
		// should never actually display this
		return "none"
	}
}

func (self WorkingTreeState) LowerCaseTitle(tr *i18n.TranslationSet) string {
	switch self {
	case WORKING_TREE_STATE_REBASING:
		return tr.LowercaseRebasingStatus
	case WORKING_TREE_STATE_MERGING:
		return tr.LowercaseMergingStatus
	default:
		// should never actually display this
		return "none"
	}
}

func (self WorkingTreeState) OptionsMenuTitle(tr *i18n.TranslationSet) string {
	if self == WORKING_TREE_STATE_MERGING {
		return tr.MergeOptionsTitle
	}
	return tr.RebaseOptionsTitle
}

func (self WorkingTreeState) CommandName() string {
	switch self {
	case WORKING_TREE_STATE_MERGING:
		return "merge"
	case WORKING_TREE_STATE_REBASING:
		return "rebase"
	default:
		// shouldn't be possible to land here
		return ""
	}
}
