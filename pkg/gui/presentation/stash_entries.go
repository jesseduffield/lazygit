package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetStashEntryListDisplayStrings(stashEntries []*commands.StashEntry, diffName string) [][]string {
	lines := make([][]string, len(stashEntries))

	for i := range stashEntries {
		diffed := stashEntries[i].RefName() == diffName
		lines[i] = getStashEntryDisplayStrings(stashEntries[i], diffed)
	}

	return lines
}

// getStashEntryDisplayStrings returns the display string of branch
func getStashEntryDisplayStrings(s *commands.StashEntry, diffed bool) []string {
	attr := theme.DefaultTextColor
	if diffed {
		attr = theme.DiffTerminalColor
	}
	return []string{utils.ColoredString(s.Name, attr)}
}
