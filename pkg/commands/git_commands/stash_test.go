package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestStashDrop(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"stash", "drop", "stash@{1}"}, "", nil)
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
