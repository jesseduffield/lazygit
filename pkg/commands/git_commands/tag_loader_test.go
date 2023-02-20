package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

const tagsOutput = `tag1 this is my message
tag2
tag3 this is my other message
`

func TestGetTags(t *testing.T) {
	type scenario struct {
		testName      string
		runner        *oscommands.FakeCmdObjRunner
		expectedTags  []*models.Tag
		expectedError error
	}

	scenarios := []scenario{
		{
			testName: "should return no tags if there are none",
			runner: oscommands.NewFakeRunner(t).
				Expect(`git tag --list -n --sort=-creatordate`, "", nil),
			expectedTags:  []*models.Tag{},
			expectedError: nil,
		},
		{
			testName: "should return tags if present",
			runner: oscommands.NewFakeRunner(t).
				Expect(`git tag --list -n --sort=-creatordate`, tagsOutput, nil),
			expectedTags: []*models.Tag{
				{Name: "tag1", Message: "this is my message"},
				{Name: "tag2", Message: ""},
				{Name: "tag3", Message: "this is my other message"},
			},
			expectedError: nil,
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.testName, func(t *testing.T) {
			loader := &TagLoader{
				Common: utils.NewDummyCommon(),
				cmd:    oscommands.NewDummyCmdObjBuilder(scenario.runner),
			}

			tags, err := loader.GetTags()

			assert.Equal(t, scenario.expectedTags, tags)
			assert.Equal(t, scenario.expectedError, err)

			scenario.runner.CheckForMissingCalls()
		})
	}
}
