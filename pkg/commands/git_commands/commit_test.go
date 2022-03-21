package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestCommitRewordCommit(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"commit", "--allow-empty", "--amend", "--only", "-m", "test"}, "", nil)
	instance := buildCommitCommands(commonDeps{runner: runner})

	assert.NoError(t, instance.RewordLastCommit("test"))
	runner.CheckForMissingCalls()
}

func TestCommitResetToCommit(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"reset", "--hard", "78976bc"}, "", nil)

	instance := buildCommitCommands(commonDeps{runner: runner})

	assert.NoError(t, instance.ResetToCommit("78976bc", "hard", []string{}))
	runner.CheckForMissingCalls()
}

func TestCommitCommitObj(t *testing.T) {
	type scenario struct {
		testName             string
		message              string
		configSignoff        bool
		configSkipHookPrefix string
		expected             string
	}

	scenarios := []scenario{
		{
			testName:             "Commit",
			message:              "test",
			configSignoff:        false,
			configSkipHookPrefix: "",
			expected:             `git commit -m "test"`,
		},
		{
			testName:             "Commit with --no-verify flag",
			message:              "WIP: test",
			configSignoff:        false,
			configSkipHookPrefix: "WIP",
			expected:             `git commit --no-verify -m "WIP: test"`,
		},
		{
			testName:             "Commit with multiline message",
			message:              "line1\nline2",
			configSignoff:        false,
			configSkipHookPrefix: "",
			expected:             `git commit -m "line1" -m "line2"`,
		},
		{
			testName:             "Commit with signoff",
			message:              "test",
			configSignoff:        true,
			configSkipHookPrefix: "",
			expected:             `git commit --signoff -m "test"`,
		},
		{
			testName:             "Commit with signoff and no-verify",
			message:              "WIP: test",
			configSignoff:        true,
			configSkipHookPrefix: "WIP",
			expected:             `git commit --no-verify --signoff -m "WIP: test"`,
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			userConfig := config.GetDefaultConfig()
			userConfig.Git.Commit.SignOff = s.configSignoff
			userConfig.Git.SkipHookPrefix = s.configSkipHookPrefix

			instance := buildCommitCommands(commonDeps{userConfig: userConfig})

			cmdStr := instance.CommitCmdObj(s.message).ToString()
			assert.Equal(t, s.expected, cmdStr)
		})
	}
}

func TestCommitCreateFixupCommit(t *testing.T) {
	type scenario struct {
		testName string
		sha      string
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			testName: "valid case",
			sha:      "12345",
			runner: oscommands.NewFakeRunner(t).
				Expect(`git commit --fixup=12345`, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildCommitCommands(commonDeps{runner: s.runner})
			s.test(instance.CreateFixupCommit(s.sha))
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestCommitShowCmdObj(t *testing.T) {
	type scenario struct {
		testName    string
		filterPath  string
		contextSize int
		expected    string
	}

	scenarios := []scenario{
		{
			testName:    "Default case without filter path",
			filterPath:  "",
			contextSize: 3,
			expected:    "git show --submodule --color=always --unified=3 --no-renames --stat -p 1234567890 ",
		},
		{
			testName:    "Default case with filter path",
			filterPath:  "file.txt",
			contextSize: 3,
			expected:    `git show --submodule --color=always --unified=3 --no-renames --stat -p 1234567890  -- "file.txt"`,
		},
		{
			testName:    "Show diff with custom context size",
			filterPath:  "",
			contextSize: 77,
			expected:    "git show --submodule --color=always --unified=77 --no-renames --stat -p 1234567890 ",
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			userConfig := config.GetDefaultConfig()
			userConfig.Git.DiffContextSize = s.contextSize

			instance := buildCommitCommands(commonDeps{userConfig: userConfig})

			cmdStr := instance.ShowCmdObj("1234567890", s.filterPath).ToString()
			assert.Equal(t, s.expected, cmdStr)
		})
	}
}

func TestGetCommitMsg(t *testing.T) {
	type scenario struct {
		testName       string
		input          string
		expectedOutput string
	}
	scenarios := []scenario{
		{
			"empty",
			` commit deadbeef`,
			``,
		},
		{
			"no line breaks (single line)",
			`commit deadbeef
use generics to DRY up context code`,
			`use generics to DRY up context code`,
		},
		{
			"with line breaks",
			`commit deadbeef
Merge pull request #1750 from mark2185/fix-issue-template

'git-rev parse' should be 'git rev-parse'`,
			`Merge pull request #1750 from mark2185/fix-issue-template

'git-rev parse' should be 'git rev-parse'`,
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildCommitCommands(commonDeps{
				runner: oscommands.NewFakeRunner(t).Expect("git rev-list --format=%B --max-count=1 deadbeef", s.input, nil),
			})

			output, err := instance.GetCommitMessage("deadbeef")

			assert.NoError(t, err)

			assert.Equal(t, s.expectedOutput, output)
		})
	}
}
