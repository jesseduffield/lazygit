package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetTagListDisplayStrings(tags []*models.Tag, diffName string) [][]string {
	lines := make([][]string, len(tags))

	for i := range tags {
		diffed := tags[i].Name == diffName
		lines[i] = getTagDisplayStrings(tags[i], diffed)
	}

	return lines
}

// getTagDisplayStrings returns the display string of branch
func getTagDisplayStrings(t *models.Tag, diffed bool) []string {
	attr := theme.DefaultTextColor
	if diffed {
		attr = theme.DiffTerminalColor
	}
	return []string{utils.ColoredString(t.Name, attr)}
}
