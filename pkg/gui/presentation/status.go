package presentation

import (
	"fmt"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
)

func FormatStatus(
	repoName string,
	currentBranch *models.Branch,
	itemOperation types.ItemOperation,
	linkedWorktreeName string,
	workingTreeState models.WorkingTreeState,
	tr *i18n.TranslationSet,
	userConfig *config.UserConfig,
) string {
	status := ""

	if currentBranch.IsRealBranch() {
		status += BranchStatus(currentBranch, itemOperation, tr, time.Now(), userConfig)
		if status != "" {
			status += " "
		}
	}

	if workingTreeState.Any() {
		status += style.FgYellow.Sprintf("(%s) ", workingTreeState.LowerCaseTitle(tr))
	}

	name := GetBranchTextStyle(currentBranch.Name).Sprint(currentBranch.Name)
	// If the user is in a linked worktree (i.e. not the main worktree) we'll display that
	if linkedWorktreeName != "" {
		icon := ""
		if icons.IsIconEnabled() {
			icon = icons.LINKED_WORKTREE_ICON + " "
		}
		repoName = fmt.Sprintf("%s(%s%s)", repoName, icon, style.FgCyan.Sprint(linkedWorktreeName))
	}
	status += fmt.Sprintf("%s → %s", repoName, name)

	return status
}
