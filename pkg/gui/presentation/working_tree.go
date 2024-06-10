package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/i18n"
)

func FormatWorkingTreeStateTitle(tr *i18n.TranslationSet, workingTreeState enums.WorkingTreeState) string {
	switch workingTreeState {
	case enums.WORKING_TREE_STATE_REBASING:
		return tr.RebasingStatus
	case enums.WORKING_TREE_STATE_MERGING:
		return tr.MergingStatus
	default:
		// should never actually display this
		return "none"
	}
}

func FormatWorkingTreeStateLower(tr *i18n.TranslationSet, workingTreeState enums.WorkingTreeState) string {
	switch workingTreeState {
	case enums.WORKING_TREE_STATE_REBASING:
		return tr.LowercaseRebasingStatus
	case enums.WORKING_TREE_STATE_MERGING:
		return tr.LowercaseMergingStatus
	default:
		// should never actually display this
		return "none"
	}
}
