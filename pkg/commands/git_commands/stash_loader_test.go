package git_commands

import (
	"fmt"
	"testing"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/stretchr/testify/assert"
)

func TestGetStashEntries(t *testing.T) {
	type scenario struct {
		testName             string
		filterPath           string
		runner               oscommands.ICmdObjRunner
		expectedStashEntries []*models.StashEntry
	}

	hoursAgo := time.Now().Unix() - 3*3600 - 1800
	daysAgo := time.Now().Unix() - 3*3600*24 - 3600*12

	scenarios := []scenario{
		{
			"No stash entries found",
			"",
			oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"stash", "list", "-z", "--pretty=%ct|%gs"}, "", nil),
			[]*models.StashEntry{},
		},
		{
			"Several stash entries found",
			"",
			oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"stash", "list", "-z", "--pretty=%ct|%gs"},
					fmt.Sprintf("%d|WIP on add-pkg-commands-test: 55c6af2 increase parallel build\x00%d|WIP on master: bb86a3f update github template\x00",
						hoursAgo,
						daysAgo,
					), nil),
			[]*models.StashEntry{
				{
					Index:   0,
					Name:    "WIP on add-pkg-commands-test: 55c6af2 increase parallel build",
					Recency: "3h",
				},
				{
					Index:   1,
					Name:    "WIP on master: bb86a3f update github template",
					Recency: "3d",
				},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			cmd := oscommands.NewDummyCmdObjBuilder(s.runner)

			loader := NewStashLoader(common.NewDummyCommon(), cmd)

			assert.EqualValues(t, s.expectedStashEntries, loader.GetStashEntries(s.filterPath))
		})
	}
}
