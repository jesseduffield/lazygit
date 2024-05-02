package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/samber/lo"
)

func GetStashEntryListDisplayStrings(stashEntries []*models.StashEntry, diffName string) [][]string {
	return lo.Map(stashEntries, func(stashEntry *models.StashEntry, _ int) []string {
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

	res := make([]string, 0, 3)
	res = append(res, style.FgCyan.Sprint(s.Recency))

	if icons.IsIconEnabled() {
		res = append(res, textStyle.Sprint(icons.IconForStash(s)))
	}

	res = append(res, textStyle.Sprint(s.Name))
	return res
}
