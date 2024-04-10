package presentation

import (
	"fmt"

	"github.com/lobes/lazytask/pkg/commands/models"
	"github.com/lobes/lazytask/pkg/commands/types/enums"
	"github.com/lobes/lazytask/pkg/config"
	"github.com/lobes/lazytask/pkg/gui/presentation/icons"
	"github.com/lobes/lazytask/pkg/gui/style"
	"github.com/lobes/lazytask/pkg/gui/types"
	"github.com/lobes/lazytask/pkg/i18n"
)

func FormatStatus(
	repoName string,
	currentBranch *models.Branch,
	itemOperation types.ItemOperation,
	linkedWorktreeName string,
	workingTreeState enums.RebaseMode,
	tr *i18n.TranslationSet,
	userConfig *config.UserConfig,
) string {
	status := ""

	if currentBranch.IsRealBranch() {
		status += ColoredBranchStatus(currentBranch, itemOperation, tr, userConfig) + " "
	}

	if workingTreeState != enums.REBASE_MODE_NONE {
		status += style.FgYellow.Sprintf("(%s) ", FormatWorkingTreeStateLower(tr, workingTreeState))
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
	status += fmt.Sprintf("%s → %s ", repoName, name)

	return status
}
