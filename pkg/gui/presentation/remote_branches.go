package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetRemoteBranchListDisplayStrings(branches []*commands.RemoteBranch) [][]string {
	lines := make([][]string, len(branches))

	for i := range branches {
		lines[i] = getRemoteBranchDisplayStrings(branches[i])
	}

	return lines
}

// getRemoteBranchDisplayStrings returns the display string of branch
func getRemoteBranchDisplayStrings(b *commands.RemoteBranch) []string {
	displayName := utils.ColoredString(b.Name, GetBranchColor(b.Name))

	return []string{displayName}
}
