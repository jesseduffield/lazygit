package git_commands

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

func TestGetCurrentSequencerTodoState(t *testing.T) {
	revParseArgs := []string{"rev-parse", "--verify", "--quiet", "HEAD"}

	scenarios := []struct {
		testName      string
		files         map[string]string
		runner        *oscommands.FakeCmdObjRunner
		expectedState sequencerTodoState
	}{
		{
			testName:      "CHERRY_PICK_HEAD exists",
			files:         map[string]string{"CHERRY_PICK_HEAD": "1234567\n"},
			runner:        oscommands.NewFakeRunner(t),
			expectedState: sequencerTodoStarted,
		},
		{
			testName:      "REVERT_HEAD exists",
			files:         map[string]string{"REVERT_HEAD": "1234567\n"},
			runner:        oscommands.NewFakeRunner(t),
			expectedState: sequencerTodoStarted,
		},
		{
			testName:      "HEAD hasn't moved since git tried to apply a commit",
			files:         map[string]string{"sequencer/abort-safety": "1234567\n"},
			runner:        oscommands.NewFakeRunner(t).ExpectGitArgs(revParseArgs, "1234567\n", nil),
			expectedState: sequencerTodoNotStarted,
		},
		{
			testName:      "HEAD has moved since git tried to apply a commit",
			files:         map[string]string{"sequencer/abort-safety": "1234567\n"},
			runner:        oscommands.NewFakeRunner(t).ExpectGitArgs(revParseArgs, "2345678\n", nil),
			expectedState: sequencerTodoCommitted,
		},
		{
			testName:      "HEAD is unborn, and there is no abort-safety file",
			files:         map[string]string{},
			runner:        oscommands.NewFakeRunner(t).ExpectGitArgs(revParseArgs, "", errors.New("error")),
			expectedState: sequencerTodoNotStarted,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitDir := t.TempDir()
			for path, content := range s.files {
				fullPath := filepath.Join(gitDir, path)
				assert.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0o755))
				assert.NoError(t, os.WriteFile(fullPath, []byte(content), 0o644))
			}

			state := getCurrentSequencerTodoState(oscommands.NewDummyCmdObjBuilder(s.runner), gitDir)

			assert.Equal(t, s.expectedState, state)
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestRepeatCurrentSequencerTodo(t *testing.T) {
	revParseArgs := []string{"rev-parse", "--verify", "--quiet", "HEAD"}
	todos := "pick 1234567 one\npick 2345678 two\n"
	repeatedTodos := "pick 1234567 one\npick 1234567 one\npick 2345678 two\n"

	setUp := func(t *testing.T, files map[string]string) string {
		t.Helper()
		gitDir := t.TempDir()
		for path, content := range files {
			fullPath := filepath.Join(gitDir, path)
			assert.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0o755))
			assert.NoError(t, os.WriteFile(fullPath, []byte(content), 0o644))
		}
		return gitDir
	}

	readTodos := func(t *testing.T, gitDir string) string {
		t.Helper()
		content, err := os.ReadFile(filepath.Join(gitDir, "sequencer", "todo"))
		assert.NoError(t, err)
		return string(content)
	}

	t.Run("no cherry-pick or revert of several commits in progress", func(t *testing.T) {
		gitDir := setUp(t, map[string]string{})
		runner := oscommands.NewFakeRunner(t)

		undo, err := repeatCurrentSequencerTodo(oscommands.NewDummyCmdObjBuilder(runner), gitDir, '#')

		assert.NoError(t, err)
		undo()
		assert.NoFileExists(t, filepath.Join(gitDir, "sequencer", "todo"))
		runner.CheckForMissingCalls()
	})

	t.Run("git started to apply the current todo", func(t *testing.T) {
		gitDir := setUp(t, map[string]string{"sequencer/todo": todos, "CHERRY_PICK_HEAD": "1234567\n"})
		runner := oscommands.NewFakeRunner(t)

		_, err := repeatCurrentSequencerTodo(oscommands.NewDummyCmdObjBuilder(runner), gitDir, '#')

		assert.NoError(t, err)
		assert.Equal(t, todos, readTodos(t, gitDir))
		runner.CheckForMissingCalls()
	})

	t.Run("git stopped before it started to apply the current todo", func(t *testing.T) {
		gitDir := setUp(t, map[string]string{"sequencer/todo": todos, "sequencer/abort-safety": "abcdef0\n"})
		runner := oscommands.NewFakeRunner(t).ExpectGitArgs(revParseArgs, "abcdef0\n", nil)

		undo, err := repeatCurrentSequencerTodo(oscommands.NewDummyCmdObjBuilder(runner), gitDir, '#')

		assert.NoError(t, err)
		assert.Equal(t, repeatedTodos, readTodos(t, gitDir))

		// If git didn't touch the todos, undo restores them
		undo()
		assert.Equal(t, todos, readTodos(t, gitDir))
		runner.CheckForMissingCalls()
	})

	t.Run("undo leaves the todos alone if git has changed them", func(t *testing.T) {
		gitDir := setUp(t, map[string]string{"sequencer/todo": todos, "sequencer/abort-safety": "abcdef0\n"})
		runner := oscommands.NewFakeRunner(t).ExpectGitArgs(revParseArgs, "abcdef0\n", nil)

		undo, err := repeatCurrentSequencerTodo(oscommands.NewDummyCmdObjBuilder(runner), gitDir, '#')
		assert.NoError(t, err)
		gitsTodos := "pick 2345678 two\n"
		assert.NoError(t, os.WriteFile(filepath.Join(gitDir, "sequencer", "todo"), []byte(gitsTodos), 0o644))

		undo()
		assert.Equal(t, gitsTodos, readTodos(t, gitDir))
		runner.CheckForMissingCalls()
	})
}
