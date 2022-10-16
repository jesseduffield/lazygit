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
		ExpectGitArgs([]string{"stash", "save", "A stash message"}, "", nil)
	instance := buildStashCommands(commonDeps{runner: runner})

	assert.NoError(t, instance.Save("A stash message"))
	runner.CheckForMissingCalls()
}

func TestStashStore(t *testing.T) {
	type scenario struct {
		testName string
		sha      string
		message  string
		expected []string
	}

	scenarios := []scenario{
		{
			testName: "Non-empty message",
			sha:      "0123456789abcdef",
			message:  "New stash name",
			expected: []string{"stash", "store", "0123456789abcdef", "-m", "New stash name"},
		},
		{
			testName: "Empty message",
			sha:      "0123456789abcdef",
			message:  "",
			expected: []string{"stash", "store", "0123456789abcdef"},
		},
		{
			testName: "Space message",
			sha:      "0123456789abcdef",
			message:  "  ",
			expected: []string{"stash", "store", "0123456789abcdef"},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			runner := oscommands.NewFakeRunner(t).
				ExpectGitArgs(s.expected, "", nil)
			instance := buildStashCommands(commonDeps{runner: runner})

			assert.NoError(t, instance.Store(s.sha, s.message))
			runner.CheckForMissingCalls()
		})
	}
}

func TestStashSha(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"rev-parse", "refs/stash@{5}"}, "14d94495194651adfd5f070590df566c11d28243\n", nil)
	instance := buildStashCommands(commonDeps{runner: runner})

	sha, err := instance.Sha(5)
	assert.NoError(t, err)
	assert.Equal(t, "14d94495194651adfd5f070590df566c11d28243", sha)
	runner.CheckForMissingCalls()
}

func TestStashStashEntryCmdObj(t *testing.T) {
	type scenario struct {
		testName    string
		index       int
		contextSize int
		expected    string
	}

	scenarios := []scenario{
		{
			testName:    "Default case",
			index:       5,
			contextSize: 3,
			expected:    "git stash show -p --stat --color=always --unified=3 stash@{5}",
		},
		{
			testName:    "Show diff with custom context size",
			index:       5,
			contextSize: 77,
			expected:    "git stash show -p --stat --color=always --unified=77 stash@{5}",
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			userConfig := config.GetDefaultConfig()
			userConfig.Git.DiffContextSize = s.contextSize
			instance := buildStashCommands(commonDeps{userConfig: userConfig})

			cmdStr := instance.ShowStashEntryCmdObj(s.index).ToString()
			assert.Equal(t, s.expected, cmdStr)
		})
	}
}

func TestStashRename(t *testing.T) {
	type scenario struct {
		testName         string
		index            int
		message          string
		expectedShaCmd   []string
		shaResult        string
		expectedDropCmd  []string
		expectedStoreCmd []string
	}

	scenarios := []scenario{
		{
			testName:         "Default case",
			index:            3,
			message:          "New message",
			expectedShaCmd:   []string{"rev-parse", "refs/stash@{3}"},
			shaResult:        "f0d0f20f2f61ffd6d6bfe0752deffa38845a3edd\n",
			expectedDropCmd:  []string{"stash", "drop", "stash@{3}"},
			expectedStoreCmd: []string{"stash", "store", "f0d0f20f2f61ffd6d6bfe0752deffa38845a3edd", "-m", "New message"},
		},
		{
			testName:         "Empty message",
			index:            4,
			message:          "",
			expectedShaCmd:   []string{"rev-parse", "refs/stash@{4}"},
			shaResult:        "f0d0f20f2f61ffd6d6bfe0752deffa38845a3edd\n",
			expectedDropCmd:  []string{"stash", "drop", "stash@{4}"},
			expectedStoreCmd: []string{"stash", "store", "f0d0f20f2f61ffd6d6bfe0752deffa38845a3edd"},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			runner := oscommands.NewFakeRunner(t).
				ExpectGitArgs(s.expectedShaCmd, s.shaResult, nil).
				ExpectGitArgs(s.expectedDropCmd, "", nil).
				ExpectGitArgs(s.expectedStoreCmd, "", nil)
			instance := buildStashCommands(commonDeps{runner: runner})

			err := instance.Rename(s.index, s.message)
			assert.NoError(t, err)
		})
	}
}
