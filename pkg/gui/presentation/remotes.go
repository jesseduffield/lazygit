package presentation

import (
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

func GetRemoteListDisplayStrings(
	remotes []*models.Remote,
	diffName string,
	getItemOperation func(item types.HasUrn) types.ItemOperation,
	tr *i18n.TranslationSet,
) [][]string {
	return lo.Map(remotes, func(remote *models.Remote, _ int) []string {
		diffed := remote.Name == diffName
		return getRemoteDisplayStrings(remote, diffed, getItemOperation(remote), tr)
	})
}

// getRemoteDisplayStrings returns the display string of branch
func getRemoteDisplayStrings(
	r *models.Remote,
	diffed bool,
	itemOperation types.ItemOperation,
	tr *i18n.TranslationSet,
) []string {
	branchCount := len(r.Branches)

	textStyle := theme.DefaultTextColor
	if diffed {
		textStyle = theme.DiffTerminalColor
	}

	res := make([]string, 0, 3)
	if icons.IsIconEnabled() {
		res = append(res, textStyle.Sprint(icons.IconForRemote(r)))
	}
	descriptionStr := style.FgBlue.Sprintf("%d branches", branchCount)
	itemOperationStr := ItemOperationToString(itemOperation, tr)
	if itemOperationStr != "" {
		descriptionStr += " " + style.FgCyan.Sprint(itemOperationStr+" "+utils.Loader(time.Now()))
	}
	res = append(res, textStyle.Sprint(r.Name), descriptionStr)
	return res
}
