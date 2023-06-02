package git_commands

import (
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCommitStoreLoaderLoad(t *testing.T) {
	commitsOutputStore := strings.Replace(`a|b
b|c d
c|e
d|e
e|`, "|", "\x00", -1)

	args := []string{"log", "--all", "--pretty=format:%H%x00%P"}

	type scenario struct {
		testName        string
		runner          *oscommands.FakeCmdObjRunner
		startingCommits []models.ImmutableCommit
		expectedCommits []models.ImmutableCommit
		expectedError   error
	}

	scenarios := []scenario{
		{
			testName: "should return no commits if there are none",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(args, "", nil),

			startingCommits: []models.ImmutableCommit{},
			expectedCommits: []models.ImmutableCommit{},
			expectedError:   nil,
		},
		{
			testName: "properly processes commits",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(args, commitsOutputStore, nil),

			startingCommits: []models.ImmutableCommit{},
			expectedCommits: []models.ImmutableCommit{
				models.NewImmutableCommit("a", []string{"b"}),
				models.NewImmutableCommit("b", []string{"c", "d"}),
				models.NewImmutableCommit("c", []string{"e"}),
				models.NewImmutableCommit("d", []string{"e"}),
				models.NewImmutableCommit("e", []string{}),
			},
			expectedError: nil,
		},
		{
			testName: "merges into exising commits",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(args, commitsOutputStore, nil),

			startingCommits: []models.ImmutableCommit{
				// this one is present in the command output
				models.NewImmutableCommit("a", []string{"b"}),
				// this one isn't
				models.NewImmutableCommit("f", []string{"g"}),
			},
			expectedCommits: []models.ImmutableCommit{
				models.NewImmutableCommit("a", []string{"b"}),
				models.NewImmutableCommit("b", []string{"c", "d"}),
				models.NewImmutableCommit("c", []string{"e"}),
				models.NewImmutableCommit("d", []string{"e"}),
				models.NewImmutableCommit("e", []string{}),
				models.NewImmutableCommit("f", []string{"g"}),
			},
			expectedError: nil,
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.testName, func(t *testing.T) {
			common := utils.NewDummyCommon()

			builder := &CommitStoreLoader{
				Common: common,
				cmd:    oscommands.NewDummyCmdObjBuilder(scenario.runner),
			}

			commitStore := models.NewCommitStore()
			commitStore.AddSlice(scenario.startingCommits)

			err := builder.Load(commitStore)

			assert.EqualValues(t, scenario.expectedCommits, commitStore.Slice())
			assert.Equal(t, scenario.expectedError, err)

			scenario.runner.CheckForMissingCalls()
		})
	}
}
