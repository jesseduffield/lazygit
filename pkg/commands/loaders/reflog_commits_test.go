package loaders

import (
	"errors"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/assert"
)

const reflogOutput = `c3c4b66b64c97ffeecde 1643150483 checkout: moving from A to B
c3c4b66b64c97ffeecde 1643150483 checkout: moving from B to A
c3c4b66b64c97ffeecde 1643150483 checkout: moving from A to B
c3c4b66b64c97ffeecde 1643150483 checkout: moving from master to A
f4ddf2f0d4be4ccc7efa 1643149435 checkout: moving from A to master
`

func TestGetReflogCommits(t *testing.T) {
	type scenario struct {
		testName                string
		runner                  *oscommands.FakeCmdObjRunner
		lastReflogCommit        *models.Commit
		filterPath              string
		expectedCommits         []*models.Commit
		expectedOnlyObtainedNew bool
		expectedError           error
	}

	scenarios := []scenario{
		{
			testName: "no reflog entries",
			runner: oscommands.NewFakeRunner(t).
				Expect(`git log -g --abbrev=40 --format="%h %ct %gs"`, "", nil),

			lastReflogCommit:        nil,
			expectedCommits:         []*models.Commit{},
			expectedOnlyObtainedNew: false,
			expectedError:           nil,
		},
		{
			testName: "some reflog entries",
			runner: oscommands.NewFakeRunner(t).
				Expect(`git log -g --abbrev=40 --format="%h %ct %gs"`, reflogOutput, nil),

			lastReflogCommit: nil,
			expectedCommits: []*models.Commit{
				{
					Sha:           "c3c4b66b64c97ffeecde",
					Name:          "checkout: moving from A to B",
					Status:        "reflog",
					UnixTimestamp: 1643150483,
				},
				{
					Sha:           "c3c4b66b64c97ffeecde",
					Name:          "checkout: moving from B to A",
					Status:        "reflog",
					UnixTimestamp: 1643150483,
				},
				{
					Sha:           "c3c4b66b64c97ffeecde",
					Name:          "checkout: moving from A to B",
					Status:        "reflog",
					UnixTimestamp: 1643150483,
				},
				{
					Sha:           "c3c4b66b64c97ffeecde",
					Name:          "checkout: moving from master to A",
					Status:        "reflog",
					UnixTimestamp: 1643150483,
				},
				{
					Sha:           "f4ddf2f0d4be4ccc7efa",
					Name:          "checkout: moving from A to master",
					Status:        "reflog",
					UnixTimestamp: 1643149435,
				},
			},
			expectedOnlyObtainedNew: false,
			expectedError:           nil,
		},
		{
			testName: "some reflog entries where last commit is given",
			runner: oscommands.NewFakeRunner(t).
				Expect(`git log -g --abbrev=40 --format="%h %ct %gs"`, reflogOutput, nil),

			lastReflogCommit: &models.Commit{
				Sha:           "c3c4b66b64c97ffeecde",
				Name:          "checkout: moving from B to A",
				Status:        "reflog",
				UnixTimestamp: 1643150483,
			},
			expectedCommits: []*models.Commit{
				{
					Sha:           "c3c4b66b64c97ffeecde",
					Name:          "checkout: moving from A to B",
					Status:        "reflog",
					UnixTimestamp: 1643150483,
				},
			},
			expectedOnlyObtainedNew: true,
			expectedError:           nil,
		},
		{
			testName: "when passing filterPath",
			runner: oscommands.NewFakeRunner(t).
				Expect(`git log -g --abbrev=40 --format="%h %ct %gs" --follow -- "path"`, reflogOutput, nil),

			lastReflogCommit: &models.Commit{
				Sha:           "c3c4b66b64c97ffeecde",
				Name:          "checkout: moving from B to A",
				Status:        "reflog",
				UnixTimestamp: 1643150483,
			},
			filterPath: "path",
			expectedCommits: []*models.Commit{
				{
					Sha:           "c3c4b66b64c97ffeecde",
					Name:          "checkout: moving from A to B",
					Status:        "reflog",
					UnixTimestamp: 1643150483,
				},
			},
			expectedOnlyObtainedNew: true,
			expectedError:           nil,
		},
		{
			testName: "when command returns error",
			runner: oscommands.NewFakeRunner(t).
				Expect(`git log -g --abbrev=40 --format="%h %ct %gs"`, "", errors.New("haha")),

			lastReflogCommit:        nil,
			filterPath:              "",
			expectedCommits:         nil,
			expectedOnlyObtainedNew: false,
			expectedError:           errors.New("haha"),
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.testName, func(t *testing.T) {
			builder := &ReflogCommitLoader{
				Common: utils.NewDummyCommon(),
				cmd:    oscommands.NewDummyCmdObjBuilder(scenario.runner),
			}

			commits, onlyObtainednew, err := builder.GetReflogCommits(scenario.lastReflogCommit, scenario.filterPath)
			assert.Equal(t, scenario.expectedOnlyObtainedNew, onlyObtainednew)
			assert.Equal(t, scenario.expectedError, err)
			t.Logf("actual commits: \n%s", litter.Sdump(commits))
			assert.Equal(t, scenario.expectedCommits, commits)

			scenario.runner.CheckForMissingCalls()
		})
	}
}
