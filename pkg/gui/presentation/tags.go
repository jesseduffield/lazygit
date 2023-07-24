package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/samber/lo"
)

func GetTagListDisplayStrings(tags []*models.Tag, diffName string) [][]string {
	return lo.Map(tags, func(tag *models.Tag, _ int) []string {
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
	descriptionColor := style.FgYellow
	res = append(res, textStyle.Sprint(t.Name), descriptionColor.Sprint(t.Description()))
	return res
}
