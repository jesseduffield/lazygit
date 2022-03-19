package presentation

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

func GetStashEntryListDisplayStrings(stashEntries []*models.StashEntry, diffName string) [][]string {
	return slices.Map(stashEntries, func(stashEntry *models.StashEntry) []string {
		diffed := stashEntry.RefName() == diffName
		return getStashEntryDisplayStrings(stashEntry, diffed)
	})
}

// getStashEntryDisplayStrings returns the display string of branch
func getStashEntryDisplayStrings(s *models.StashEntry, diffed bool) []string {
	textStyle := theme.DefaultTextColor
	if diffed {
		textStyle = theme.DiffTerminalColor
	}
	return []string{textStyle.Sprint(s.Name)}
}
