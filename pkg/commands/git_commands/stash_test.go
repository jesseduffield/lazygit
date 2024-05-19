package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestStashDrop(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"stash", "drop", "stash@{1}"}, "Dropped refs/stash@{1} (98e9cca532c37c766107093010c72e26f2c24c04)\n", nil)
	instance := buildStashCommands(commonDeps{runner: runner})

	assert.NoError(t, instance.Drop(1))
	runner.CheckForMissingCalls()
}

func TestStashApply(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"stash", "apply", "stash@{1}"}, "", nil)
	instance := buildStashCommands(commonDeps{runner: runner})

	assert.NoError(t, instance.Apply(1))
	runner.CheckForMissingCalls()
}

func TestStashPop(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"stash", "pop", "stash@{1}"}, "", nil)
	instance := buildStashCommands(commonDeps{runner: runner})

	assert.NoError(t, instance.Pop(1))
	runner.CheckForMissingCalls()
}

func TestStashSave(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"stash", "push", "-m", "A stash message"}, "", nil)
	instance := buildStashCommands(commonDeps{runner: runner})

	assert.NoError(t, instance.Push("A stash message"))
	runner.CheckForMissingCalls()
}

func TestStashStore(t *testing.T) {
	type scenario struct {
		testName string
		hash     string
		message  string
		expected []string
	}

	scenarios := []scenario{
		{
			testName: "Non-empty message",
			hash:     "0123456789abcdef",
			message:  "New stash name",
			expected: []string{"stash", "store", "-m", "New stash name", "0123456789abcdef"},
		},
		{
			testName: "Empty message",
			hash:     "0123456789abcdef",
			message:  "",
			expected: []string{"stash", "store", "0123456789abcdef"},
		},
		{
			testName: "Space message",
			hash:     "0123456789abcdef",
			message:  "  ",
			expected: []string{"stash", "store", "0123456789abcdef"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			runner := oscommands.NewFakeRunner(t).
				ExpectGitArgs(s.expected, "", nil)
			instance := buildStashCommands(commonDeps{runner: runner})

			assert.NoError(t, instance.Store(s.hash, s.message))
			runner.CheckForMissingCalls()
		})
	}
}

func TestStashHash(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"rev-parse", "refs/stash@{5}"}, "14d94495194651adfd5f070590df566c11d28243\n", nil)
	instance := buildStashCommands(commonDeps{runner: runner})

	hash, err := instance.Hash(5)
	assert.NoError(t, err)
	assert.Equal(t, "14d94495194651adfd5f070590df566c11d28243", hash)
	runner.CheckForMissingCalls()
}

func TestStashStashEntryCmdObj(t *testing.T) {
	type scenario struct {
		testName         string
		index            int
		contextSize      int
		ignoreWhitespace bool
		expected         []string
	}

	scenarios := []scenario{
		{
			testName:         "Default case",
			index:            5,
			contextSize:      3,
			ignoreWhitespace: false,
			expected:         []string{"git", "-C", "/path/to/worktree", "stash", "show", "-p", "--stat", "--color=always", "--unified=3", "stash@{5}"},
		},
		{
			testName:         "Show diff with custom context size",
			index:            5,
			contextSize:      77,
			ignoreWhitespace: false,
			expected:         []string{"git", "-C", "/path/to/worktree", "stash", "show", "-p", "--stat", "--color=always", "--unified=77", "stash@{5}"},
		},
		{
			testName:         "Default case",
			index:            5,
			contextSize:      3,
			ignoreWhitespace: true,
			expected:         []string{"git", "-C", "/path/to/worktree", "stash", "show", "-p", "--stat", "--color=always", "--unified=3", "--ignore-all-space", "stash@{5}"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			userConfig := config.GetDefaultConfig()
			appState := &config.AppState{}
			appState.IgnoreWhitespaceInDiffView = s.ignoreWhitespace
			appState.DiffContextSize = s.contextSize
			repoPaths := RepoPaths{
				worktreePath: "/path/to/worktree",
			}
			instance := buildStashCommands(commonDeps{userConfig: userConfig, appState: appState, repoPaths: &repoPaths})

			cmdStr := instance.ShowStashEntryCmdObj(s.index).Args()
			assert.Equal(t, s.expected, cmdStr)
		})
	}
}

func TestStashRename(t *testing.T) {
	type scenario struct {
		testName         string
		index            int
		message          string
		expectedHashCmd  []string
		hashResult       string
		expectedDropCmd  []string
		expectedStoreCmd []string
	}

	scenarios := []scenario{
		{
			testName:         "Default case",
			index:            3,
			message:          "New message",
			expectedHashCmd:  []string{"rev-parse", "refs/stash@{3}"},
			hashResult:       "f0d0f20f2f61ffd6d6bfe0752deffa38845a3edd\n",
			expectedDropCmd:  []string{"stash", "drop", "stash@{3}"},
			expectedStoreCmd: []string{"stash", "store", "-m", "New message", "f0d0f20f2f61ffd6d6bfe0752deffa38845a3edd"},
		},
		{
			testName:         "Empty message",
			index:            4,
			message:          "",
			expectedHashCmd:  []string{"rev-parse", "refs/stash@{4}"},
			hashResult:       "f0d0f20f2f61ffd6d6bfe0752deffa38845a3edd\n",
			expectedDropCmd:  []string{"stash", "drop", "stash@{4}"},
			expectedStoreCmd: []string{"stash", "store", "f0d0f20f2f61ffd6d6bfe0752deffa38845a3edd"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			runner := oscommands.NewFakeRunner(t).
				ExpectGitArgs(s.expectedHashCmd, s.hashResult, nil).
				ExpectGitArgs(s.expectedDropCmd, "", nil).
				ExpectGitArgs(s.expectedStoreCmd, "", nil)
			instance := buildStashCommands(commonDeps{runner: runner})

			err := instance.Rename(s.index, s.message)
			assert.NoError(t, err)
		})
	}
}
