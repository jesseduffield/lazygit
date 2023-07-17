package presentation

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func FormatStatus(currentBranch *models.Branch, linkedWorktreeName string, workingTreeState enums.RebaseMode, tr *i18n.TranslationSet) string {
	status := ""

	if currentBranch.IsRealBranch() {
		status += ColoredBranchStatus(currentBranch, tr) + " "
	}

	if workingTreeState != enums.REBASE_MODE_NONE {
		status += style.FgYellow.Sprintf("(%s) ", FormatWorkingTreeStateLower(tr, workingTreeState))
	}

	name := GetBranchTextStyle(currentBranch.Name).Sprint(currentBranch.Name)
	repoName := utils.GetCurrentRepoName()
	// If the user is in a linked worktree (i.e. not the main worktree) we'll display that
	if linkedWorktreeName != "" {
		repoName = fmt.Sprintf("%s(%s)", repoName, style.FgCyan.Sprint(linkedWorktreeName))
	}
	status += fmt.Sprintf("%s → %s ", repoName, name)

	return status
}
