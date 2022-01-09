package loaders

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetStashEntries(t *testing.T) {
	type scenario struct {
		testName             string
		filterPath           string
		runner               oscommands.ICmdObjRunner
		expectedStashEntries []*models.StashEntry
	}

	scenarios := []scenario{
		{
			"No stash entries found",
			"",
			oscommands.NewFakeRunner(t).
				Expect(`git stash list --pretty='%gs'`, "", nil),
			[]*models.StashEntry{},
		},
		{
			"Several stash entries found",
			"",
			oscommands.NewFakeRunner(t).
				Expect(
					`git stash list --pretty='%gs'`,
					"WIP on add-pkg-commands-test: 55c6af2 increase parallel build\nWIP on master: bb86a3f update github template",
					nil,
				),
			[]*models.StashEntry{
				{
					Index: 0,
					Name:  "WIP on add-pkg-commands-test: 55c6af2 increase parallel build",
				},
				{
					Index: 1,
					Name:  "WIP on master: bb86a3f update github template",
				},
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			cmd := oscommands.NewDummyCmdObjBuilder(s.runner)

			loader := NewStashLoader(utils.NewDummyCommon(), cmd)

			assert.EqualValues(t, s.expectedStashEntries, loader.GetStashEntries(s.filterPath))
		})
	}
}
