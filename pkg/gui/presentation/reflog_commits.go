package presentation

import (
	"time"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/kyokomi/emoji/v2"
)

func GetReflogCommitListDisplayStrings(commits []*models.Commit, fullDescription bool, cherryPickedCommitShaSet *set.Set[string], diffName string, now time.Time, timeFormat string, shortTimeFormat string, parseEmoji bool) [][]string {
	var displayFunc func(*models.Commit, reflogCommitDisplayAttributes) []string
	if fullDescription {
		displayFunc = getFullDescriptionDisplayStringsForReflogCommit
	} else {
		displayFunc = getDisplayStringsForReflogCommit
	}

	return slices.Map(commits, func(commit *models.Commit) []string {
		diffed := commit.Sha == diffName
		cherryPicked := cherryPickedCommitShaSet.Includes(commit.Sha)
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

func reflogShaColor(cherryPicked, diffed bool) style.TextStyle {
	if diffed {
		return theme.DiffTerminalColor
	}

	shaColor := style.FgBlue
	if cherryPicked {
		shaColor = theme.CherryPickedCommitTextStyle
	}

	return shaColor
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
		reflogShaColor(attrs.cherryPicked, attrs.diffed).Sprint(c.ShortSha()),
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
		reflogShaColor(attrs.cherryPicked, attrs.diffed).Sprint(c.ShortSha()),
		theme.DefaultTextColor.Sprint(name),
	}
}
