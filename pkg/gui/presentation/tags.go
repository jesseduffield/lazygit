package presentation

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

func GetTagListDisplayStrings(tags []*models.Tag, diffName string) [][]string {
	return slices.Map(tags, func(tag *models.Tag) []string {
		diffed := tag.Name == diffName
		return getTagDisplayStrings(tag, diffed)
	})
}

// getTagDisplayStrings returns the display string of branch
func getTagDisplayStrings(t *models.Tag, diffed bool) []string {
	textStyle := theme.DefaultTextColor
	if diffed {
		textStyle = theme.DiffTerminalColor
	}
	res := make([]string, 0, 2)
	if icons.IsIconEnabled() {
		res = append(res, textStyle.Sprint(icons.IconForTag(t)))
	}
	res = append(res, textStyle.Sprint(t.Name))
	return res
}
