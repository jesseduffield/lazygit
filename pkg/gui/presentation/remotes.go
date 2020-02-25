package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
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
	return []string{r.Name}
}
