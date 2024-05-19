package git_commands

import (
	"errors"
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/assert"
)

var reflogOutput = strings.Replace(`c3c4b66b64c97ffeecde|1643150483|checkout: moving from A to B|51baa8c1
c3c4b66b64c97ffeecde|1643150483|checkout: moving from B to A|51baa8c1
c3c4b66b64c97ffeecde|1643150483|checkout: moving from A to B|51baa8c1
c3c4b66b64c97ffeecde|1643150483|checkout: moving from master to A|51baa8c1
f4ddf2f0d4be4ccc7efa|1643149435|checkout: moving from A to master|51baa8c1
`, "|", "\x00", -1)

func TestGetReflogCommits(t *testing.T) {
	type scenario struct {
		testName                string
		runner                  *oscommands.FakeCmdObjRunner
		lastReflogCommit        *models.Commit
		filterPath              string
		filterAuthor            string
		expectedCommits         []*models.Commit
		expectedOnlyObtainedNew bool
		expectedError           error
	}

	scenarios := []scenario{
		{
			testName: "no reflog entries",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"-c", "log.showSignature=false", "log", "-g", "--abbrev=40", "--format=%h%x00%ct%x00%gs%x00%p"}, "", nil),

			lastReflogCommit:        nil,
			expectedCommits:         []*models.Commit{},
			expectedOnlyObtainedNew: false,
			expectedError:           nil,
		},
		{
			testName: "some reflog entries",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"-c", "log.showSignature=false", "log", "-g", "--abbrev=40", "--format=%h%x00%ct%x00%gs%x00%p"}, reflogOutput, nil),

			lastReflogCommit: nil,
			expectedCommits: []*models.Commit{
				{
					Hash:          "c3c4b66b64c97ffeecde",
					Name:          "checkout: moving from A to B",
					Status:        models.StatusReflog,
					UnixTimestamp: 1643150483,
					Parents:       []string{"51baa8c1"},
				},
				{
					Hash:          "c3c4b66b64c97ffeecde",
					Name:          "checkout: moving from B to A",
					Status:        models.StatusReflog,
					UnixTimestamp: 1643150483,
					Parents:       []string{"51baa8c1"},
				},
				{
					Hash:          "c3c4b66b64c97ffeecde",
					Name:          "checkout: moving from A to B",
					Status:        models.StatusReflog,
					UnixTimestamp: 1643150483,
					Parents:       []string{"51baa8c1"},
				},
				{
					Hash:          "c3c4b66b64c97ffeecde",
					Name:          "checkout: moving from master to A",
					Status:        models.StatusReflog,
					UnixTimestamp: 1643150483,
					Parents:       []string{"51baa8c1"},
				},
				{
					Hash:          "f4ddf2f0d4be4ccc7efa",
					Name:          "checkout: moving from A to master",
					Status:        models.StatusReflog,
					UnixTimestamp: 1643149435,
					Parents:       []string{"51baa8c1"},
				},
			},
			expectedOnlyObtainedNew: false,
			expectedError:           nil,
		},
		{
			testName: "some reflog entries where last commit is given",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"-c", "log.showSignature=false", "log", "-g", "--abbrev=40", "--format=%h%x00%ct%x00%gs%x00%p"}, reflogOutput, nil),

			lastReflogCommit: &models.Commit{
				Hash:          "c3c4b66b64c97ffeecde",
				Name:          "checkout: moving from B to A",
				Status:        models.StatusReflog,
				UnixTimestamp: 1643150483,
				Parents:       []string{"51baa8c1"},
			},
			expectedCommits: []*models.Commit{
				{
					Hash:          "c3c4b66b64c97ffeecde",
					Name:          "checkout: moving from A to B",
					Status:        models.StatusReflog,
					UnixTimestamp: 1643150483,
					Parents:       []string{"51baa8c1"},
				},
			},
			expectedOnlyObtainedNew: true,
			expectedError:           nil,
		},
		{
			testName: "when passing filterPath",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"-c", "log.showSignature=false", "log", "-g", "--abbrev=40", "--format=%h%x00%ct%x00%gs%x00%p", "--follow", "--", "path"}, reflogOutput, nil),

			lastReflogCommit: &models.Commit{
				Hash:          "c3c4b66b64c97ffeecde",
				Name:          "checkout: moving from B to A",
				Status:        models.StatusReflog,
				UnixTimestamp: 1643150483,
				Parents:       []string{"51baa8c1"},
			},
			filterPath: "path",
			expectedCommits: []*models.Commit{
				{
					Hash:          "c3c4b66b64c97ffeecde",
					Name:          "checkout: moving from A to B",
					Status:        models.StatusReflog,
					UnixTimestamp: 1643150483,
					Parents:       []string{"51baa8c1"},
				},
			},
			expectedOnlyObtainedNew: true,
			expectedError:           nil,
		},
		{
			testName: "when passing filterAuthor",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"-c", "log.showSignature=false", "log", "-g", "--abbrev=40", "--format=%h%x00%ct%x00%gs%x00%p", "--author=John Doe <john@doe.com>"}, reflogOutput, nil),

			lastReflogCommit: &models.Commit{
				Hash:          "c3c4b66b64c97ffeecde",
				Name:          "checkout: moving from B to A",
				Status:        models.StatusReflog,
				UnixTimestamp: 1643150483,
				Parents:       []string{"51baa8c1"},
			},
			filterAuthor: "John Doe <john@doe.com>",
			expectedCommits: []*models.Commit{
				{
					Hash:          "c3c4b66b64c97ffeecde",
					Name:          "checkout: moving from A to B",
					Status:        models.StatusReflog,
					UnixTimestamp: 1643150483,
					Parents:       []string{"51baa8c1"},
				},
			},
			expectedOnlyObtainedNew: true,
			expectedError:           nil,
		},
		{
			testName: "when command returns error",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"-c", "log.showSignature=false", "log", "-g", "--abbrev=40", "--format=%h%x00%ct%x00%gs%x00%p"}, "", errors.New("haha")),

			lastReflogCommit:        nil,
			filterPath:              "",
			expectedCommits:         nil,
			expectedOnlyObtainedNew: false,
			expectedError:           errors.New("haha"),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.testName, func(t *testing.T) {
			builder := &ReflogCommitLoader{
				Common: utils.NewDummyCommon(),
				cmd:    oscommands.NewDummyCmdObjBuilder(scenario.runner),
			}

			commits, onlyObtainednew, err := builder.GetReflogCommits(scenario.lastReflogCommit, scenario.filterPath, scenario.filterAuthor)
			assert.Equal(t, scenario.expectedOnlyObtainedNew, onlyObtainednew)
			assert.Equal(t, scenario.expectedError, err)
			t.Logf("actual commits: \n%s", litter.Sdump(commits))
			assert.Equal(t, scenario.expectedCommits, commits)

			scenario.runner.CheckForMissingCalls()
		})
	}
}
