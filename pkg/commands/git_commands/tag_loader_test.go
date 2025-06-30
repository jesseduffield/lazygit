package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/stretchr/testify/assert"
)

const tagsOutput = "refs/tags/tag1\x00tag\x00this is my message\n" +
	"refs/tags/tag2\x00commit\x00\n" +
	"refs/tags/tag3\x00tag\x00this is my other message\n"

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
				ExpectGitArgs([]string{"for-each-ref", "--sort=-creatordate", "--format=%(refname)%00%(objecttype)%00%(contents:subject)", "refs/tags"}, "", nil),
			expectedTags:  []*models.Tag{},
			expectedError: nil,
		},
		{
			testName: "should return tags if present",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"for-each-ref", "--sort=-creatordate", "--format=%(refname)%00%(objecttype)%00%(contents:subject)", "refs/tags"}, tagsOutput, nil),
			expectedTags: []*models.Tag{
				{Name: "tag1", Message: "this is my message", IsAnnotated: true},
				{Name: "tag2", Message: "", IsAnnotated: false},
				{Name: "tag3", Message: "this is my other message", IsAnnotated: true},
			},
			expectedError: nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.testName, func(t *testing.T) {
			loader := &TagLoader{
				Common: common.NewDummyCommon(),
				cmd:    oscommands.NewDummyCmdObjBuilder(scenario.runner),
			}

			tags, err := loader.GetTags()

			assert.Equal(t, scenario.expectedTags, tags)
			assert.Equal(t, scenario.expectedError, err)

			scenario.runner.CheckForMissingCalls()
		})
	}
}
