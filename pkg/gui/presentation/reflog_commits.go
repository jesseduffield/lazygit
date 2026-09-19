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

func GetReflogCommitListDisplayStrings(commits []*models.Commit, startIdx int, endIdx int, fullDescription bool, cherryPickedCommitHashSet *set.Set[string], diffName string, now time.Time, timeFormat string, shortTimeFormat string, parseEmoji bool) [][]string {
	if startIdx >= len(commits) {
		return nil
	}

	var displayFunc func(*models.Commit, reflogCommitDisplayAttributes) []string
	reservedDateWidth := 0
	if fullDescription {
		displayFunc = getFullDescriptionDisplayStringsForReflogCommit
		// See getReservedColumnWidths for why the oldest entry alone decides
		// how much width the date column needs.
		reservedDateWidth = utils.StringWidth(utils.UnixToDateSmart(
			now, commits[len(commits)-1].UnixTimestamp, timeFormat, shortTimeFormat))
	} else {
		displayFunc = getDisplayStringsForReflogCommit
	}

	return lo.Map(commits[startIdx:endIdx], func(commit *models.Commit, _ int) []string {
		diffed := commit.Hash() == diffName
		cherryPicked := cherryPickedCommitHashSet.Includes(commit.Hash())
		return displayFunc(commit,
			reflogCommitDisplayAttributes{
				cherryPicked:      cherryPicked,
				diffed:            diffed,
				parseEmoji:        parseEmoji,
				timeFormat:        timeFormat,
				shortTimeFormat:   shortTimeFormat,
				now:               now,
				reservedDateWidth: reservedDateWidth,
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
	// The width the date column needs for the whole reflog, not just for the
	// lines that are on screen
	reservedDateWidth int
}

func getFullDescriptionDisplayStringsForReflogCommit(c *models.Commit, attrs reflogCommitDisplayAttributes) []string {
	name := c.Name
	if attrs.parseEmoji {
		name = emoji.Sprint(name)
	}

	date := style.FgMagenta.Sprint(
		utils.UnixToDateSmart(attrs.now, c.UnixTimestamp, attrs.timeFormat, attrs.shortTimeFormat))

	return []string{
		reflogHashColor(attrs.cherryPicked, attrs.diffed).Sprint(c.ShortHash()),
		utils.WithPadding(date, attrs.reservedDateWidth, utils.AlignLeft),
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
