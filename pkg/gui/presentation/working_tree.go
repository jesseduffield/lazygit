package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/i18n"
)

func FormatWorkingTreeStateTitle(tr *i18n.TranslationSet, rebaseMode enums.RebaseMode) string {
	switch rebaseMode {
	case enums.REBASE_MODE_REBASING:
		return tr.RebasingStatus
	case enums.REBASE_MODE_MERGING:
		return tr.MergingStatus
	default:
		// should never actually display this
		return "none"
	}
}

func FormatWorkingTreeStateLower(tr *i18n.TranslationSet, rebaseMode enums.RebaseMode) string {
	switch rebaseMode {
	case enums.REBASE_MODE_REBASING:
		return tr.LowercaseRebasingStatus
	case enums.REBASE_MODE_MERGING:
		return tr.LowercaseMergingStatus
	default:
		// should never actually display this
		return "none"
	}
}
