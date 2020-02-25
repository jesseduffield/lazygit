package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
)

func GetTagListDisplayStrings(tags []*commands.Tag) [][]string {
	lines := make([][]string, len(tags))

	for i := range tags {
		lines[i] = getTagDisplayStrings(tags[i])
	}

	return lines
}

// getTagDisplayStrings returns the display string of branch
func getTagDisplayStrings(t *commands.Tag) []string {
	return []string{t.Name}
}
