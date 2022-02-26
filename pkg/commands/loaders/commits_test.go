package loaders

import (
	"path/filepath"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func NewDummyCommitLoader() *CommitLoader {
	cmn := utils.NewDummyCommon()

	return &CommitLoader{
		Common:               cmn,
		cmd:                  nil,
		getCurrentBranchName: func() (string, string, error) { return "master", "master", nil },
		getRebaseMode:        func() (enums.RebaseMode, error) { return enums.REBASE_MODE_NONE, nil },
		dotGitDir:            ".git",
		readFile: func(filename string) ([]byte, error) {
			return []byte(""), nil
		},
		walkFiles: func(root string, fn filepath.WalkFunc) error {
			return nil
		},
	}
}

const commitsOutput = `0eea75e8c631fba6b58135697835d58ba4c18dbc|1640826609|Jesse Duffield| (HEAD -> better-tests)|b21997d6b4cbdf84b149|better typing for rebase mode
b21997d6b4cbdf84b149d8e6a2c4d06a8e9ec164|1640824515|Jesse Duffield| (origin/better-tests)|e94e8fc5b6fab4cb755f|fix logging
e94e8fc5b6fab4cb755f29f1bdb3ee5e001df35c|1640823749|Jesse Duffield||d8084cd558925eb7c9c3|refactor
d8084cd558925eb7c9c38afeed5725c21653ab90|1640821426|Jesse Duffield||65f910ebd85283b5cce9|WIP
65f910ebd85283b5cce9bf67d03d3f1a9ea3813a|1640821275|Jesse Duffield||26c07b1ab33860a1a759|WIP
26c07b1ab33860a1a7591a0638f9925ccf497ffa|1640750752|Jesse Duffield||3d4470a6c072208722e5|WIP
3d4470a6c072208722e5ae9a54bcb9634959a1c5|1640748818|Jesse Duffield||053a66a7be3da43aacdc|WIP
053a66a7be3da43aacdc7aa78e1fe757b82c4dd2|1640739815|Jesse Duffield||985fe482e806b172aea4|refactoring the config struct`

func TestGetCommits(t *testing.T) {
	type scenario struct {
		testName          string
		runner            *oscommands.FakeCmdObjRunner
		expectedCommits   []*models.Commit
		expectedError     error
		rebaseMode        enums.RebaseMode
		currentBranchName string
		opts              GetCommitsOptions
	}

	scenarios := []scenario{
		{
			testName:          "should return no commits if there are none",
			rebaseMode:        enums.REBASE_MODE_NONE,
			currentBranchName: "master",
			opts:              GetCommitsOptions{RefName: "HEAD", IncludeRebaseCommits: false},
			runner: oscommands.NewFakeRunner(t).
				Expect(`git merge-base "HEAD" "HEAD"@{u}`, "b21997d6b4cbdf84b149d8e6a2c4d06a8e9ec164", nil).
				Expect(`git log "HEAD" --topo-order  --oneline --pretty=format:"%H|%at|%aN|%d|%p|%s" --abbrev=40`, "", nil),

			expectedCommits: []*models.Commit{},
			expectedError:   nil,
		},
		{
			testName:          "should return commits if they are present",
			rebaseMode:        enums.REBASE_MODE_NONE,
			currentBranchName: "master",
			opts:              GetCommitsOptions{RefName: "HEAD", IncludeRebaseCommits: false},
			runner: oscommands.NewFakeRunner(t).
				// here it's seeing which commits are yet to be pushed
				Expect(`git merge-base "HEAD" "HEAD"@{u}`, "b21997d6b4cbdf84b149d8e6a2c4d06a8e9ec164", nil).
				// here it's actually getting all the commits in a formatted form, one per line
				Expect(`git log "HEAD" --topo-order  --oneline --pretty=format:"%H|%at|%aN|%d|%p|%s" --abbrev=40`, commitsOutput, nil).
				// here it's seeing where our branch diverged from the master branch so that we can mark that commit and parent commits as 'merged'
				Expect(`git merge-base "HEAD" "master"`, "26c07b1ab33860a1a7591a0638f9925ccf497ffa", nil),

			expectedCommits: []*models.Commit{
				{
					Sha:           "0eea75e8c631fba6b58135697835d58ba4c18dbc",
					Name:          "better typing for rebase mode",
					Status:        "unpushed",
					Action:        "",
					Tags:          []string{},
					ExtraInfo:     "(HEAD -> better-tests)",
					Author:        "Jesse Duffield",
					UnixTimestamp: 1640826609,
					Parents: []string{
						"b21997d6b4cbdf84b149",
					},
				},
				{
					Sha:           "b21997d6b4cbdf84b149d8e6a2c4d06a8e9ec164",
					Name:          "fix logging",
					Status:        "pushed",
					Action:        "",
					Tags:          []string{},
					ExtraInfo:     "(origin/better-tests)",
					Author:        "Jesse Duffield",
					UnixTimestamp: 1640824515,
					Parents: []string{
						"e94e8fc5b6fab4cb755f",
					},
				},
				{
					Sha:           "e94e8fc5b6fab4cb755f29f1bdb3ee5e001df35c",
					Name:          "refactor",
					Status:        "pushed",
					Action:        "",
					Tags:          []string{},
					ExtraInfo:     "",
					Author:        "Jesse Duffield",
					UnixTimestamp: 1640823749,
					Parents: []string{
						"d8084cd558925eb7c9c3",
					},
				},
				{
					Sha:           "d8084cd558925eb7c9c38afeed5725c21653ab90",
					Name:          "WIP",
					Status:        "pushed",
					Action:        "",
					Tags:          []string{},
					ExtraInfo:     "",
					Author:        "Jesse Duffield",
					UnixTimestamp: 1640821426,
					Parents: []string{
						"65f910ebd85283b5cce9",
					},
				},
				{
					Sha:           "65f910ebd85283b5cce9bf67d03d3f1a9ea3813a",
					Name:          "WIP",
					Status:        "pushed",
					Action:        "",
					Tags:          []string{},
					ExtraInfo:     "",
					Author:        "Jesse Duffield",
					UnixTimestamp: 1640821275,
					Parents: []string{
						"26c07b1ab33860a1a759",
					},
				},
				{
					Sha:           "26c07b1ab33860a1a7591a0638f9925ccf497ffa",
					Name:          "WIP",
					Status:        "merged",
					Action:        "",
					Tags:          []string{},
					ExtraInfo:     "",
					Author:        "Jesse Duffield",
					UnixTimestamp: 1640750752,
					Parents: []string{
						"3d4470a6c072208722e5",
					},
				},
				{
					Sha:           "3d4470a6c072208722e5ae9a54bcb9634959a1c5",
					Name:          "WIP",
					Status:        "merged",
					Action:        "",
					Tags:          []string{},
					ExtraInfo:     "",
					Author:        "Jesse Duffield",
					UnixTimestamp: 1640748818,
					Parents: []string{
						"053a66a7be3da43aacdc",
					},
				},
				{
					Sha:           "053a66a7be3da43aacdc7aa78e1fe757b82c4dd2",
					Name:          "refactoring the config struct",
					Status:        "merged",
					Action:        "",
					Tags:          []string{},
					ExtraInfo:     "",
					Author:        "Jesse Duffield",
					UnixTimestamp: 1640739815,
					Parents: []string{
						"985fe482e806b172aea4",
					},
				},
			},
			expectedError: nil,
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.testName, func(t *testing.T) {
			builder := &CommitLoader{
				Common: utils.NewDummyCommon(),
				cmd:    oscommands.NewDummyCmdObjBuilder(scenario.runner),
				getCurrentBranchName: func() (string, string, error) {
					return scenario.currentBranchName, scenario.currentBranchName, nil
				},
				getRebaseMode: func() (enums.RebaseMode, error) { return scenario.rebaseMode, nil },
				dotGitDir:     ".git",
				readFile: func(filename string) ([]byte, error) {
					return []byte(""), nil
				},
				walkFiles: func(root string, fn filepath.WalkFunc) error {
					return nil
				},
			}

			commits, err := builder.GetCommits(scenario.opts)

			assert.Equal(t, scenario.expectedCommits, commits)
			assert.Equal(t, scenario.expectedError, err)

			scenario.runner.CheckForMissingCalls()
		})
	}
}
