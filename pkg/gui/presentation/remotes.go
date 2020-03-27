package presentation

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetRemoteListDisplayStrings(remotes []*commands.Remote) [][]string {
	lines := make([][]string, len(remotes))

	for i := range remotes {
		lines[i] = getRemoteDisplayStrings(remotes[i])
	}

	return lines
}

// getRemoteDisplayStrings returns the display string of branch
func getRemoteDisplayStrings(r *commands.Remote) []string {
	branchCount := len(r.Branches)

	return []string{r.Name, utils.ColoredString(fmt.Sprintf("%d branches", branchCount), color.FgBlue)}
}
