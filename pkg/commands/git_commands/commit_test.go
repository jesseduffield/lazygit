package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestCommitRewordCommit(t *testing.T) {
	type scenario struct {
		testName string
		runner   *oscommands.FakeCmdObjRunner
		input    string
	}
	scenarios := []scenario{
		{
			"Single line reword",
			oscommands.NewFakeRunner(t).ExpectGitArgs([]string{"commit", "--allow-empty", "--amend", "--only", "-m", "test"}, "", nil),
			"test",
		},
		{
			"Multi line reword",
			oscommands.NewFakeRunner(t).ExpectGitArgs([]string{"commit", "--allow-empty", "--amend", "--only", "-m", "test", "-m", "line 2\nline 3"}, "", nil),
			"test\nline 2\nline 3",
		},
	}
	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildCommitCommands(commonDeps{runner: s.runner})

			assert.NoError(t, instance.RewordLastCommit(s.input))
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestCommitResetToCommit(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"reset", "--hard", "78976bc"}, "", nil)

	instance := buildCommitCommands(commonDeps{runner: runner})

	assert.NoError(t, instance.ResetToCommit("78976bc", "hard", []string{}))
	runner.CheckForMissingCalls()
}

func TestCommitCommitCmdObj(t *testing.T) {
	type scenario struct {
		testName             string
		message              string
		configSignoff        bool
		configSkipHookPrefix string
		expectedArgs         []string
	}

	scenarios := []scenario{
		{
			testName:             "Commit",
			message:              "test",
			configSignoff:        false,
			configSkipHookPrefix: "",
			expectedArgs:         []string{"commit", "-m", "test"},
		},
		{
			testName:             "Commit with --no-verify flag",
			message:              "WIP: test",
			configSignoff:        false,
			configSkipHookPrefix: "WIP",
			expectedArgs:         []string{"commit", "--no-verify", "-m", "WIP: test"},
		},
		{
			testName:             "Commit with multiline message",
			message:              "line1\nline2",
			configSignoff:        false,
			configSkipHookPrefix: "",
			expectedArgs:         []string{"commit", "-m", "line1", "-m", "line2"},
		},
		{
			testName:             "Commit with signoff",
			message:              "test",
			configSignoff:        true,
			configSkipHookPrefix: "",
			expectedArgs:         []string{"commit", "--signoff", "-m", "test"},
		},
		{
			testName:             "Commit with signoff and no-verify",
			message:              "WIP: test",
			configSignoff:        true,
			configSkipHookPrefix: "WIP",
			expectedArgs:         []string{"commit", "--no-verify", "--signoff", "-m", "WIP: test"},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			userConfig := config.GetDefaultConfig()
			userConfig.Git.Commit.SignOff = s.configSignoff
			userConfig.Git.SkipHookPrefix = s.configSkipHookPrefix

			runner := oscommands.NewFakeRunner(t).ExpectGitArgs(s.expectedArgs, "", nil)
			instance := buildCommitCommands(commonDeps{userConfig: userConfig, runner: runner})

			assert.NoError(t, instance.CommitCmdObj(s.message).Run())
			runner.CheckForMissingCalls()
		})
	}
}

func TestCommitCommitEditorCmdObj(t *testing.T) {
	type scenario struct {
		testName      string
		configSignoff bool
		expected      []string
	}

	scenarios := []scenario{
		{
			testName:      "Commit using editor",
			configSignoff: false,
			expected:      []string{"commit"},
		},
		{
			testName:      "Commit with --signoff",
			configSignoff: true,
			expected:      []string{"commit", "--signoff"},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			userConfig := config.GetDefaultConfig()
			userConfig.Git.Commit.SignOff = s.configSignoff

			runner := oscommands.NewFakeRunner(t).ExpectGitArgs(s.expected, "", nil)
			instance := buildCommitCommands(commonDeps{userConfig: userConfig, runner: runner})

			assert.NoError(t, instance.CommitEditorCmdObj().Run())
			runner.CheckForMissingCalls()
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
				ExpectGitArgs([]string{"commit", "--fixup=12345"}, "", nil),
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
		testName         string
		filterPath       string
		contextSize      int
		ignoreWhitespace bool
		expected         []string
	}

	scenarios := []scenario{
		{
			testName:         "Default case without filter path",
			filterPath:       "",
			contextSize:      3,
			ignoreWhitespace: false,
			expected:         []string{"show", "--submodule", "--color=always", "--unified=3", "--stat", "--decorate", "-p", "1234567890"},
		},
		{
			testName:         "Default case with filter path",
			filterPath:       "file.txt",
			contextSize:      3,
			ignoreWhitespace: false,
			expected:         []string{"show", "--submodule", "--color=always", "--unified=3", "--stat", "--decorate", "-p", "1234567890", "--", "file.txt"},
		},
		{
			testName:         "Show diff with custom context size",
			filterPath:       "",
			contextSize:      77,
			ignoreWhitespace: false,
			expected:         []string{"show", "--submodule", "--color=always", "--unified=77", "--stat", "--decorate", "-p", "1234567890"},
		},
		{
			testName:         "Show diff, ignoring whitespace",
			filterPath:       "",
			contextSize:      77,
			ignoreWhitespace: true,
			expected:         []string{"show", "--submodule", "--color=always", "--unified=77", "--stat", "--decorate", "-p", "1234567890", "--ignore-all-space"},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			userConfig := config.GetDefaultConfig()
			userConfig.Git.DiffContextSize = s.contextSize

			runner := oscommands.NewFakeRunner(t).ExpectGitArgs(s.expected, "", nil)
			instance := buildCommitCommands(commonDeps{userConfig: userConfig, runner: runner})

			assert.NoError(t, instance.ShowCmdObj("1234567890", s.filterPath, s.ignoreWhitespace).Run())
			runner.CheckForMissingCalls()
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
				runner: oscommands.NewFakeRunner(t).ExpectGitArgs([]string{"rev-list", "--format=%B", "--max-count=1", "deadbeef"}, s.input, nil),
			})

			output, err := instance.GetCommitMessage("deadbeef")

			assert.NoError(t, err)

			assert.Equal(t, s.expectedOutput, output)
		})
	}
}

func TestGetCommitMessageFromHistory(t *testing.T) {
	type scenario struct {
		testName string
		runner   *oscommands.FakeCmdObjRunner
		test     func(string, error)
	}
	scenarios := []scenario{
		{
			"Empty message",
			oscommands.NewFakeRunner(t).ExpectGitArgs([]string{"log", "-1", "--skip=2", "--pretty=%H"}, "", nil).ExpectGitArgs([]string{"rev-list", "--format=%B", "--max-count=1"}, "", nil),
			func(output string, err error) {
				assert.Error(t, err)
			},
		},
		{
			"Default case to retrieve a commit in history",
			oscommands.NewFakeRunner(t).ExpectGitArgs([]string{"log", "-1", "--skip=2", "--pretty=%H"}, "sha3 \n", nil).ExpectGitArgs([]string{"rev-list", "--format=%B", "--max-count=1", "sha3"}, `commit sha3
				use generics to DRY up context code`, nil),
			func(output string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "use generics to DRY up context code", output)
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildCommitCommands(commonDeps{runner: s.runner})

			output, err := instance.GetCommitMessageFromHistory(2)

			s.test(output, err)
		})
	}
}
