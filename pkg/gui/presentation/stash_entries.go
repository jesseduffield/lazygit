package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

func GetStashEntryListDisplayStrings(stashEntries []*models.StashEntry, diffName string) [][]string {
	lines := make([][]string, len(stashEntries))

	for i := range stashEntries {
		diffed := stashEntries[i].RefName() == diffName
		lines[i] = getStashEntryDisplayStrings(stashEntries[i], diffed)
	}

	return lines
}

// getStashEntryDisplayStrings returns the display string of branch
func getStashEntryDisplayStrings(s *models.StashEntry, diffed bool) []string {
	textStyle := theme.DefaultTextColor
	if diffed {
		textStyle = theme.DiffTerminalColor
	}
	return []string{textStyle.Sprint(s.Name)}
}
