package presentation

import "github.com/jesseduffield/lazygit/pkg/commands/types/enums"

func FormatWorkingTreeState(rebaseMode enums.RebaseMode) string {
	switch rebaseMode {
	case enums.REBASE_MODE_REBASING:
		return "rebasing"
	case enums.REBASE_MODE_MERGING:
		return "merging"
	default:
		return "none"
	}
}
