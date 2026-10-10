package presentation

import (
	"strings"
	"testing"
	"time"

	"github.com/gookit/color"
	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/xo/terminfo"
)

func TestGetReflogCommitListDisplayStrings(t *testing.T) {
	scenarios := []struct {
		testName        string
		commitOpts      []models.NewCommitOpts
		fullDescription bool
		timeFormat      string
		shortTimeFormat string
		now             time.Time
		startIdx        int
		endIdx          int
		expected        string
	}{
		{
			testName:   "no commits",
			commitOpts: []models.NewCommitOpts{},
			startIdx:   0,
			endIdx:     1,
			now:        time.Date(2020, 1, 1, 5, 3, 4, 0, time.UTC),
			expected:   "",
		},
		{
			testName: "some commits",
			commitOpts: []models.NewCommitOpts{
				{Name: "checkout: moving from master to mybranch", Hash: "hash1"},
				{Name: "commit: make a change", Hash: "hash2"},
			},
			startIdx: 0,
			endIdx:   2,
			now:      time.Date(2020, 1, 1, 5, 3, 4, 0, time.UTC),
			expected: formatExpected(`
		hash1 checkout: moving from master to mybranch
		hash2 commit: make a change
				`),
		},
		{
			testName: "full description",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit: today", Hash: "hash1", UnixTimestamp: 1577844184},
				{Name: "commit: a while ago", Hash: "hash2", UnixTimestamp: 1576844184},
			},
			fullDescription: true,
			timeFormat:      "2006-01-02",
			shortTimeFormat: "3:04PM",
			startIdx:        0,
			endIdx:          2,
			now:             time.Date(2020, 1, 1, 5, 3, 4, 0, time.UTC),
			expected: formatExpected(`
		hash1 2:03AM     commit: today
		hash2 2019-12-20 commit: a while ago
				`),
		},
		{
			testName: "only showing commits from today",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit: today", Hash: "hash1", UnixTimestamp: 1577844184},
				{Name: "commit: a while ago", Hash: "hash2", UnixTimestamp: 1576844184},
			},
			fullDescription: true,
			timeFormat:      "2006-01-02",
			shortTimeFormat: "3:04PM",
			startIdx:        0,
			endIdx:          1,
			now:             time.Date(2020, 1, 1, 5, 3, 4, 0, time.UTC),
			expected: formatExpected(`
		hash1 2:03AM     commit: today
				`),
		},
	}

	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelNone)
	defer color.ForceSetColorLevel(oldColorLevel)

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			hashPool := &utils.StringPool{}

			commits := lo.Map(s.commitOpts,
				func(opts models.NewCommitOpts, _ int) *models.Commit { return models.NewCommit(hashPool, opts) })

			result := GetReflogCommitListDisplayStrings(
				commits,
				s.startIdx,
				s.endIdx,
				s.fullDescription,
				set.New[string](),
				"",
				s.now,
				s.timeFormat,
				s.shortTimeFormat,
				false,
			)

			renderedLines, _ := utils.RenderDisplayStrings(result, nil)
			renderedResult := strings.Join(renderedLines, "\n")
			t.Logf("\n%s", renderedResult)

			assert.EqualValues(t, s.expected, renderedResult)
		})
	}
}
