package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

func GetRemoteBranchListDisplayStrings(branches []*models.RemoteBranch, diffName string) [][]string {
	lines := make([][]string, len(branches))

	for i := range branches {
		diffed := branches[i].FullName() == diffName
		lines[i] = getRemoteBranchDisplayStrings(branches[i], diffed)
	}

	return lines
}

// getRemoteBranchDisplayStrings returns the display string of branch
func getRemoteBranchDisplayStrings(b *models.RemoteBranch, diffed bool) []string {
	nameColorAttr := GetBranchColor(b.Name)
	if diffed {
		nameColorAttr = theme.DiffTerminalColor
	}

	return []string{nameColorAttr.Sprint(b.Name)}
}
