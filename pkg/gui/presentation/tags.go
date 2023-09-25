package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

func GetTagListDisplayStrings(
	tags []*models.Tag,
	getRefOperation func(branch *models.Tag) types.RefOperation,
	diffName string,
	tr *i18n.TranslationSet,
) [][]string {
	return lo.Map(tags, func(tag *models.Tag, _ int) []string {
		diffed := tag.Name == diffName
		return getTagDisplayStrings(tag, getRefOperation(tag), diffed, tr)
	})
}

// getTagDisplayStrings returns the display string of branch
func getTagDisplayStrings(t *models.Tag, refOperation types.RefOperation, diffed bool, tr *i18n.TranslationSet) []string {
	textStyle := theme.DefaultTextColor
	if diffed {
		textStyle = theme.DiffTerminalColor
	}
	res := make([]string, 0, 2)
	if icons.IsIconEnabled() {
		res = append(res, textStyle.Sprint(icons.IconForTag(t)))
	}
	descriptionColor := style.FgYellow
	descriptionStr := descriptionColor.Sprint(t.Description())
	refOperationStr := refOperationToString(refOperation, tr)
	if refOperationStr != "" {
		descriptionStr = style.FgCyan.Sprint(refOperationStr+" "+utils.Loader()) + " " + descriptionStr
	}
	res = append(res, textStyle.Sprint(t.Name), descriptionStr)
	return res
}
