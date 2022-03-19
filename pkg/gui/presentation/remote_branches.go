package presentation

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

func GetRemoteBranchListDisplayStrings(branches []*models.RemoteBranch, diffName string) [][]string {
	return slices.Map(branches, func(branch *models.RemoteBranch) []string {
		diffed := branch.FullName() == diffName
		return getRemoteBranchDisplayStrings(branch, diffed)
	})
}

// getRemoteBranchDisplayStrings returns the display string of branch
func getRemoteBranchDisplayStrings(b *models.RemoteBranch, diffed bool) []string {
	textStyle := GetBranchTextStyle(b.Name)
	if diffed {
		textStyle = theme.DiffTerminalColor
	}

	return []string{textStyle.Sprint(b.Name)}
}
