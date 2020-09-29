package presentation

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
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

	nameColorAttr := theme.DefaultTextColor
	if diffed {
		nameColorAttr = theme.DiffTerminalColor
	}

	return []string{utils.ColoredString(r.Name, nameColorAttr), utils.ColoredString(fmt.Sprintf("%d branches", branchCount), color.FgBlue)}
}
