package presentation

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

func GetWorktreeListDisplayStrings(worktrees []*models.Worktree) [][]string {
	return slices.Map(worktrees, func(worktree *models.Worktree) []string {
		return getWorktreeDisplayStrings(worktree)
	})
}

// getWorktreeDisplayStrings returns the display string of branch
func getWorktreeDisplayStrings(w *models.Worktree) []string {
	textStyle := theme.DefaultTextColor

	current := ""
	currentColor := style.FgCyan
	if w.Current() {
		current = "  *"
		currentColor = style.FgGreen
	}

	icon := icons.IconForWorktree(w, false)
	if w.Missing() {
		textStyle = style.FgRed
		icon = icons.IconForWorktree(w, true)
	}

	res := make([]string, 0, 3)
	res = append(res, currentColor.Sprint(current))
	if icons.IsIconEnabled() {
		res = append(res, textStyle.Sprint(icon))
	}
	res = append(res, textStyle.Sprint(w.Name()))
	return res
}
