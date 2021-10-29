package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/theme"
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
	textStyle := theme.DefaultTextColor
	if diffed {
		textStyle = theme.DiffTerminalColor
	}
	return []string{textStyle.Sprint(t.Name)}
}
