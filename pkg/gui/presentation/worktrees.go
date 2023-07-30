package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/samber/lo"
)

func GetWorktreeDisplayStrings(tr *i18n.TranslationSet, worktrees []*models.Worktree) [][]string {
	return lo.Map(worktrees, func(worktree *models.Worktree, _ int) []string {
		return GetWorktreeDisplayString(
			tr,
			worktree)
	})
}

func GetWorktreeDisplayString(tr *i18n.TranslationSet, worktree *models.Worktree) []string {
	textStyle := theme.DefaultTextColor

	current := ""
	currentColor := style.FgCyan
	if worktree.IsCurrent {
		current = "  *"
		currentColor = style.FgGreen
	}

	icon := icons.IconForWorktree(false)
	if worktree.IsPathMissing {
		textStyle = style.FgRed
		icon = icons.IconForWorktree(true)
	}

	res := []string{}
	res = append(res, currentColor.Sprint(current))
	if icons.IsIconEnabled() {
		res = append(res, textStyle.Sprint(icon))
	}

	name := worktree.Name
	if worktree.IsMain {
		name += " " + tr.MainWorktree
	}
	if worktree.IsPathMissing && !icons.IsIconEnabled() {
		name += " " + tr.MissingWorktree
	}
	res = append(res, textStyle.Sprint(name))
	return res
}
