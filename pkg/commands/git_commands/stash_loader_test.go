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
				ExpectGitArgs([]string{"stash", "list", "-z", "--pretty=%H|%ct|%gs"}, "", nil),
			[]*models.StashEntry{},
		},
		{
			"Several stash entries found",
			"",
			oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"stash", "list", "-z", "--pretty=%H|%ct|%gs"},
					fmt.Sprintf("fa1afe1|%d|WIP on add-pkg-commands-test: 55c6af2 increase parallel build\x00deadbeef|%d|WIP on master: bb86a3f update github template\x00",
						hoursAgo,
						daysAgo,
					), nil),
			[]*models.StashEntry{
				{
					Index:   0,
					Name:    "WIP on add-pkg-commands-test: 55c6af2 increase parallel build",
					Recency: "3h",
					Hash:    "fa1afe1",
				},
				{
					Index:   1,
					Name:    "WIP on master: bb86a3f update github template",
					Recency: "3d",
					Hash:    "deadbeef",
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
