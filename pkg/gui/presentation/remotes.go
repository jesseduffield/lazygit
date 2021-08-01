package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

func GetRemoteListDisplayStrings(remotes []*models.Remote, diffName string) [][]string {
	lines := make([][]string, len(remotes))

	for i := range remotes {
		diffed := remotes[i].Name == diffName
		lines[i] = getRemoteDisplayStrings(remotes[i], diffed)
	}

	return lines
}

// getRemoteDisplayStrings returns the display string of branch
func getRemoteDisplayStrings(r *models.Remote, diffed bool) []string {
	branchCount := len(r.Branches)

	textStyle := theme.DefaultTextColor
	if diffed {
		textStyle = theme.DiffTerminalColor
	}

	return []string{textStyle.Sprint(r.Name), style.FgBlue.Sprintf("%d branches", branchCount)}
}
