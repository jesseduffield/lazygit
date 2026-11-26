package git_commands

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

func TestResumeRebaseAfterLockError(t *testing.T) {
	newCmdObj := func(args ...string) *oscommands.CmdObj {
		return oscommands.NewDummyCmdObjBuilder(nil).New(args)
	}

	writeRebaseMergeFiles := func(gitDir string, done string, todo string) {
		dir := filepath.Join(gitDir, "rebase-merge")
		assert.NoError(t, os.MkdirAll(dir, 0o755))
		assert.NoError(t, os.WriteFile(filepath.Join(dir, "done"), []byte(done), 0o644))
		assert.NoError(t, os.WriteFile(filepath.Join(dir, "git-rebase-todo"), []byte(todo), 0o644))
	}

	// Sets up the recovery for cmdObj, then calls runUntilLockError to do what
	// the command did before the lock stopped it, and returns the command that
	// the recovery wants to run next.
	recoverFromLockError := func(gitDir string, cmdObj *oscommands.CmdObj, runUntilLockError func()) *oscommands.CmdObj {
		resumeRebaseAfterLockError(cmdObj, gitDir, '#', func() *oscommands.CmdObj {
			return newCmdObj("git", "rebase", "--continue")
		})
		runUntilLockError()
		return cmdObj.GetLockErrorRecovery()()
	}

	t.Run("runs the command again if it didn't start the rebase", func(t *testing.T) {
		cmdObj := newCmdObj("git", "rebase", "--interactive", "main")

		assert.Same(t, cmdObj, recoverFromLockError(t.TempDir(), cmdObj, func() {}))
	})

	t.Run("runs the command again if it didn't change the rebase", func(t *testing.T) {
		gitDir := t.TempDir()
		writeRebaseMergeFiles(gitDir,
			"pick 1234567 one\n",
			"pick 2345678 two\n")
		cmdObj := newCmdObj("git", "rebase", "--continue")

		assert.Same(t, cmdObj, recoverFromLockError(gitDir, cmdObj, func() {}))
	})

	t.Run("continues the rebase if git rescheduled the todo that failed", func(t *testing.T) {
		gitDir := t.TempDir()

		next := recoverFromLockError(gitDir, newCmdObj("git", "rebase", "--interactive", "main"), func() {
			writeRebaseMergeFiles(gitDir,
				"pick 1234567 one\npick 2345678 two\n",
				"pick 2345678 two\npick 3456789 three\n")
		})

		assert.Equal(t, "git rebase --continue", next.ToString())

		// If continuing fails before it changes the rebase, we run it again
		assert.Same(t, next, next.GetLockErrorRecovery()())

		// If continuing fails after it has changed the rebase, and git
		// rescheduled the todo that failed, we continue again
		writeRebaseMergeFiles(gitDir,
			"pick 1234567 one\npick 2345678 two\npick 2345678 two\npick 3456789 three\n",
			"pick 3456789 three\n")
		nextAfterContinuing := next.GetLockErrorRecovery()()
		assert.Equal(t, "git rebase --continue", nextAfterContinuing.ToString())
		assert.NotSame(t, next, nextAfterContinuing)

		// If the rebase is gone after continuing failed, we give up
		assert.NoError(t, os.RemoveAll(filepath.Join(gitDir, "rebase-merge")))
		assert.Nil(t, nextAfterContinuing.GetLockErrorRecovery()())
	})

	t.Run("continues the rebase if an older git version rescheduled the todo that failed", func(t *testing.T) {
		gitDir := t.TempDir()

		next := recoverFromLockError(gitDir, newCmdObj("git", "rebase", "--interactive", "main"), func() {
			writeRebaseMergeFiles(gitDir,
				"pick 1234567 one\npick 2345678 two\npick 1234567 one\n",
				"pick 2345678 two\npick 3456789 three\n")
		})

		assert.Equal(t, "git rebase --continue", next.ToString())
	})

	t.Run("gives up if git didn't reschedule the todo that failed", func(t *testing.T) {
		gitDir := t.TempDir()

		assert.Nil(t, recoverFromLockError(gitDir, newCmdObj("git", "rebase", "--interactive", "main"), func() {
			writeRebaseMergeFiles(gitDir,
				"pick 1234567 one\npick 2345678 two\n",
				"pick 3456789 three\n")
		}))
	})

	t.Run("gives up if a rebase without a rebase-merge directory is in progress", func(t *testing.T) {
		gitDir := t.TempDir()

		assert.Nil(t, recoverFromLockError(gitDir, newCmdObj("git", "pull"), func() {
			assert.NoError(t, os.MkdirAll(filepath.Join(gitDir, "rebase-apply"), 0o755))
		}))
	})
}

func TestAbortSequenceAfterLockError(t *testing.T) {
	scenarios := []struct {
		testName           string
		numCommits         int
		sequencerExists    bool
		runner             *oscommands.FakeCmdObjRunner
		expectRecovery     bool
		expectRunningAgain bool
	}{
		{
			testName:       "single commit",
			numCommits:     1,
			runner:         oscommands.NewFakeRunner(t),
			expectRecovery: false,
		},
		{
			testName:           "the lock stopped the command before it started",
			numCommits:         2,
			runner:             oscommands.NewFakeRunner(t),
			expectRecovery:     true,
			expectRunningAgain: true,
		},
		{
			testName:        "the lock stopped the command partway through",
			numCommits:      2,
			sequencerExists: true,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"cherry-pick", "--abort"}, "", nil),
			expectRecovery:     true,
			expectRunningAgain: true,
		},
		{
			testName:        "aborting fails",
			numCommits:      2,
			sequencerExists: true,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"cherry-pick", "--abort"}, "", errors.New("error")),
			expectRecovery:     true,
			expectRunningAgain: false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitDir := t.TempDir()
			if s.sequencerExists {
				assert.NoError(t, os.MkdirAll(filepath.Join(gitDir, "sequencer"), 0o755))
			}
			builder := oscommands.NewDummyCmdObjBuilder(s.runner)
			cmdObj := builder.New([]string{"git", "cherry-pick", "1234567", "2345678"})

			abortSequenceAfterLockError(cmdObj, gitDir, s.numCommits, func() *oscommands.CmdObj {
				return builder.New([]string{"git", "cherry-pick", "--abort"})
			})

			recovery := cmdObj.GetLockErrorRecovery()
			if !s.expectRecovery {
				assert.Nil(t, recovery)
			} else if s.expectRunningAgain {
				assert.Same(t, cmdObj, recovery())
			} else {
				assert.Nil(t, recovery())
			}
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestRetryContinueOrSkipAfterLockError(t *testing.T) {
	newCmdObj := func(args ...string) *oscommands.CmdObj {
		return oscommands.NewDummyCmdObjBuilder(nil).New(args)
	}

	writeTodos := func(gitDir string, todos string) {
		dir := filepath.Join(gitDir, "sequencer")
		assert.NoError(t, os.MkdirAll(dir, 0o755))
		assert.NoError(t, os.WriteFile(filepath.Join(dir, "todo"), []byte(todos), 0o644))
	}

	t.Run("runs the command again if it didn't change the todos", func(t *testing.T) {
		gitDir := t.TempDir()
		writeTodos(gitDir, "pick 1234567 one\npick 2345678 two\n")
		cmdObj := newCmdObj("git", "cherry-pick", "--continue")

		retryContinueOrSkipAfterLockError(cmdObj, filepath.Join(gitDir, "sequencer", "todo"))

		assert.Same(t, cmdObj, cmdObj.GetLockErrorRecovery()())
	})

	t.Run("runs the command again for a cherry-pick of a single commit", func(t *testing.T) {
		cmdObj := newCmdObj("git", "cherry-pick", "--continue")

		retryContinueOrSkipAfterLockError(cmdObj, filepath.Join(t.TempDir(), "sequencer", "todo"))

		assert.Same(t, cmdObj, cmdObj.GetLockErrorRecovery()())
	})

	t.Run("gives up if it changed the todos", func(t *testing.T) {
		gitDir := t.TempDir()
		writeTodos(gitDir, "pick 1234567 one\npick 2345678 two\n")
		cmdObj := newCmdObj("git", "cherry-pick", "--continue")

		retryContinueOrSkipAfterLockError(cmdObj, filepath.Join(gitDir, "sequencer", "todo"))
		writeTodos(gitDir, "pick 2345678 two\n")

		assert.Nil(t, cmdObj.GetLockErrorRecovery()())
	})

	t.Run("gives up if the cherry-pick is gone", func(t *testing.T) {
		gitDir := t.TempDir()
		writeTodos(gitDir, "pick 1234567 one\n")
		cmdObj := newCmdObj("git", "cherry-pick", "--continue")

		retryContinueOrSkipAfterLockError(cmdObj, filepath.Join(gitDir, "sequencer", "todo"))
		assert.NoError(t, os.RemoveAll(filepath.Join(gitDir, "sequencer")))

		assert.Nil(t, cmdObj.GetLockErrorRecovery()())
	})
}
