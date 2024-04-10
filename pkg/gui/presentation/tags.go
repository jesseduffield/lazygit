package presentation

import (
	"time"

	"github.com/lobes/lazytask/pkg/commands/models"
	"github.com/lobes/lazytask/pkg/config"
	"github.com/lobes/lazytask/pkg/gui/presentation/icons"
	"github.com/lobes/lazytask/pkg/gui/style"
	"github.com/lobes/lazytask/pkg/gui/types"
	"github.com/lobes/lazytask/pkg/i18n"
	"github.com/lobes/lazytask/pkg/theme"
	"github.com/lobes/lazytask/pkg/utils"
	"github.com/samber/lo"
)

func GetTagListDisplayStrings(
	tags []*models.Tag,
	getItemOperation func(item types.HasUrn) types.ItemOperation,
	diffName string,
	tr *i18n.TranslationSet,
	userConfig *config.UserConfig,
) [][]string {
	return lo.Map(tags, func(tag *models.Tag, _ int) []string {
		diffed := tag.Name == diffName
		return getTagDisplayStrings(tag, getItemOperation(tag), diffed, tr, userConfig)
	})
}

// getTagDisplayStrings returns the display string of branch
func getTagDisplayStrings(
	t *models.Tag,
	itemOperation types.ItemOperation,
	diffed bool,
	tr *i18n.TranslationSet,
	userConfig *config.UserConfig,
) []string {
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
	itemOperationStr := ItemOperationToString(itemOperation, tr)
	if itemOperationStr != "" {
		descriptionStr = style.FgCyan.Sprint(itemOperationStr+" "+utils.Loader(time.Now(), userConfig.Gui.Spinner)) + " " + descriptionStr
	}
	res = append(res, textStyle.Sprint(t.Name), descriptionStr)
	return res
}
