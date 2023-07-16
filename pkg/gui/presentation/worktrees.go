package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/samber/lo"
)

func GetWorktreeDisplayStrings(worktrees []*models.Worktree, isCurrent func(*models.Worktree) bool, isMissing func(*models.Worktree) bool) [][]string {
	return lo.Map(worktrees, func(worktree *models.Worktree, _ int) []string {
		return GetWorktreeDisplayString(
			isCurrent(worktree),
			isMissing(worktree),
			worktree)
	})
}

func GetWorktreeDisplayString(isCurrent bool, isPathMissing bool, worktree *models.Worktree) []string {
	textStyle := theme.DefaultTextColor

	current := ""
	currentColor := style.FgCyan
	if isCurrent {
		current = "  *"
		currentColor = style.FgGreen
	}

	icon := icons.IconForWorktree(false)
	if isPathMissing {
		textStyle = style.FgRed
		icon = icons.IconForWorktree(true)
	}

	res := []string{}
	res = append(res, currentColor.Sprint(current))
	if icons.IsIconEnabled() {
		res = append(res, textStyle.Sprint(icon))
	}
	res = append(res, textStyle.Sprint(worktree.Name()))
	return res
}
