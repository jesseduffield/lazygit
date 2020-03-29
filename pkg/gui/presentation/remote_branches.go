package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetRemoteBranchListDisplayStrings(branches []*commands.RemoteBranch, diffName string) [][]string {
	lines := make([][]string, len(branches))

	for i := range branches {
		diffed := branches[i].FullName() == diffName
		lines[i] = getRemoteBranchDisplayStrings(branches[i], diffed)
	}

	return lines
}

// getRemoteBranchDisplayStrings returns the display string of branch
func getRemoteBranchDisplayStrings(b *commands.RemoteBranch, diffed bool) []string {
	nameColorAttr := GetBranchColor(b.Name)
	if diffed {
		nameColorAttr = theme.DiffTerminalColor
	}

	displayName := utils.ColoredString(b.Name, nameColorAttr)

	return []string{displayName}
}
