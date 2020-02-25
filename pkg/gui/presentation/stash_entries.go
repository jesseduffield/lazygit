package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
)

func GetStashEntryListDisplayStrings(stashEntries []*commands.StashEntry) [][]string {
	lines := make([][]string, len(stashEntries))

	for i := range stashEntries {
		lines[i] = getStashEntryDisplayStrings(stashEntries[i])
	}

	return lines
}

// getStashEntryDisplayStrings returns the display string of branch
func getStashEntryDisplayStrings(s *commands.StashEntry) []string {
	return []string{s.DisplayString}
}
