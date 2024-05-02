package presentation

import (
	"time"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/kyokomi/emoji/v2"
	"github.com/samber/lo"
)

func GetReflogCommitListDisplayStrings(commits []*models.Commit, fullDescription bool, cherryPickedCommitHashSet *set.Set[string], diffName string, now time.Time, timeFormat string, shortTimeFormat string, parseEmoji bool) [][]string {
	var displayFunc func(*models.Commit, reflogCommitDisplayAttributes) []string
	if fullDescription {
		displayFunc = getFullDescriptionDisplayStringsForReflogCommit
	} else {
		displayFunc = getDisplayStringsForReflogCommit
	}

	return lo.Map(commits, func(commit *models.Commit, _ int) []string {
		diffed := commit.Hash == diffName
		cherryPicked := cherryPickedCommitHashSet.Includes(commit.Hash)
		return displayFunc(commit,
			reflogCommitDisplayAttributes{
				cherryPicked:    cherryPicked,
				diffed:          diffed,
				parseEmoji:      parseEmoji,
				timeFormat:      timeFormat,
				shortTimeFormat: shortTimeFormat,
				now:             now,
			})
	})
}

func reflogHashColor(cherryPicked, diffed bool) style.TextStyle {
	if diffed {
		return theme.DiffTerminalColor
	}

	hashColor := style.FgBlue
	if cherryPicked {
		hashColor = theme.CherryPickedCommitTextStyle
	}

	return hashColor
}

type reflogCommitDisplayAttributes struct {
	cherryPicked    bool
	diffed          bool
	parseEmoji      bool
	timeFormat      string
	shortTimeFormat string
	now             time.Time
}

func getFullDescriptionDisplayStringsForReflogCommit(c *models.Commit, attrs reflogCommitDisplayAttributes) []string {
	name := c.Name
	if attrs.parseEmoji {
		name = emoji.Sprint(name)
	}

	return []string{
		reflogHashColor(attrs.cherryPicked, attrs.diffed).Sprint(c.ShortHash()),
		style.FgMagenta.Sprint(utils.UnixToDateSmart(attrs.now, c.UnixTimestamp, attrs.timeFormat, attrs.shortTimeFormat)),
		theme.DefaultTextColor.Sprint(name),
	}
}

func getDisplayStringsForReflogCommit(c *models.Commit, attrs reflogCommitDisplayAttributes) []string {
	name := c.Name
	if attrs.parseEmoji {
		name = emoji.Sprint(name)
	}

	return []string{
		reflogHashColor(attrs.cherryPicked, attrs.diffed).Sprint(c.ShortHash()),
		theme.DefaultTextColor.Sprint(name),
	}
}
