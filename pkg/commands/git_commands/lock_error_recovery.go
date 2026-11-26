package git_commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// resumeRebaseAfterLockError sets the lock-error recovery (see
// CmdObj.RecoverFromLockErrorWith) of cmdObj, a command that starts or
// continues a rebase in the worktree whose git dir is gitDir. Call it right
// before running the command, since it records the state of the rebase at
// that time.
//
// If a lock stops the command before it changes the state of the rebase, we
// run the command again. This includes a command that hasn't started the
// rebase yet, since git cleans up after itself in that case. If a lock stops
// the rebase partway through, git normally reschedules the todo that failed,
// and the command made by newContinueCmdObj ("git rebase --continue") runs it
// again. In some cases git doesn't reschedule it, for example if the lock
// stops a pick after it has applied the changes but before it has committed
// them. Continuing could then lose the commit, so we give up and leave the
// rebase to the user.
func resumeRebaseAfterLockError(
	cmdObj *oscommands.CmdObj,
	gitDir string,
	commentChar byte,
	newContinueCmdObj func() *oscommands.CmdObj,
) *oscommands.CmdObj {
	var recoverFrom func(cmdObj *oscommands.CmdObj, stateBefore string) *oscommands.CmdObj
	recoverFrom = func(cmdObj *oscommands.CmdObj, stateBefore string) *oscommands.CmdObj {
		return cmdObj.RecoverFromLockErrorWith(func() *oscommands.CmdObj {
			state := rebaseState(gitDir)
			if state == stateBefore {
				return cmdObj
			}

			if !wasLastRebaseTodoRescheduled(gitDir, commentChar) {
				return nil
			}

			return recoverFrom(newContinueCmdObj(), state)
		})
	}

	return recoverFrom(cmdObj, rebaseState(gitDir))
}

// abortSequenceAfterLockError sets the lock-error recovery (see
// CmdObj.RecoverFromLockErrorWith) of cmdObj, a cherry-pick or revert of
// numCommits commits in the worktree whose git dir is gitDir.
//
// If a lock stops a cherry-pick or revert of several commits after it has
// applied some of them, git leaves the operation in progress in
// .git/sequencer. Unlike a rebase, it doesn't reschedule the commit that
// failed, and continuing the operation would skip that commit. So we abort the
// operation with the command made by newAbortCmdObj, and run the command
// again. The runner calls the recovery after its last attempt too, so we
// also abort when it gives up, and the user isn't left with an operation that
// stopped halfway.
func abortSequenceAfterLockError(
	cmdObj *oscommands.CmdObj,
	gitDir string,
	numCommits int,
	newAbortCmdObj func() *oscommands.CmdObj,
) *oscommands.CmdObj {
	// For a single commit, git doesn't use .git/sequencer. If a lock stops it,
	// it has either not changed anything yet, or it has staged the changes, and
	// running it again fails with an error. Git also doesn't check whether
	// another operation is in progress in .git/sequencer, so we must not abort
	// that one.
	if numCommits < 2 {
		return cmdObj
	}

	return cmdObj.RecoverFromLockErrorWith(func() *oscommands.CmdObj {
		// Git refuses to start a cherry-pick or revert of several commits while
		// one is in progress, so if there is one now, our command started it.
		if dirExists(filepath.Join(gitDir, "sequencer")) {
			if err := newAbortCmdObj().Run(); err != nil {
				return nil
			}
		}
		return cmdObj
	})
}

// retryContinueOrSkipAfterLockError sets the lock-error recovery (see
// CmdObj.RecoverFromLockErrorWith) of cmdObj, a command that continues or skips
// a step of an operation in progress, such as a cherry-pick or revert.
// progressFile is the file in which git records how far the operation has got.
// Call this function right before running the command, since it records the
// contents of progressFile at that time.
//
// If a lock stops the command before it changes progressFile, we run the
// command again. Otherwise it has got past the step that it continued or
// skipped, and the lock stopped it at the next one. Unlike a rebase, git doesn't
// reschedule that step (see abortSequenceAfterLockError), and running the
// command again could skip it. So we give up and leave the operation to the
// user.
func retryContinueOrSkipAfterLockError(cmdObj *oscommands.CmdObj, progressFile string) *oscommands.CmdObj {
	progressBefore, _ := os.ReadFile(progressFile)
	return cmdObj.RecoverFromLockErrorWith(func() *oscommands.CmdObj {
		progress, _ := os.ReadFile(progressFile)
		if bytes.Equal(progress, progressBefore) {
			return cmdObj
		}
		return nil
	})
}

// rebaseState returns a snapshot of the rebase that is in progress in gitDir,
// made of the files in which git records how far it has got. It's empty if no
// rebase is in progress.
func rebaseState(gitDir string) string {
	var state strings.Builder
	for _, dir := range []string{"rebase-merge", "rebase-apply"} {
		if dirExists(filepath.Join(gitDir, dir)) {
			state.WriteString(dir + "\n")
		}
	}
	for _, file := range []string{"rebase-merge/done", "rebase-merge/git-rebase-todo", "rebase-apply/next"} {
		if content, err := os.ReadFile(filepath.Join(gitDir, file)); err == nil {
			state.WriteString(file + "\n")
			state.Write(content)
		}
	}
	return state.String()
}

// wasLastRebaseTodoRescheduled tells whether git put the todo that it executed
// last back into the todo list of the rebase that is in progress in gitDir.
// It's false if no rebase is in progress, or if the rebase doesn't use the
// rebase-merge directory; that kind of rebase never reschedules anything.
func wasLastRebaseTodoRescheduled(gitDir string, commentChar byte) bool {
	doneTodos, err := utils.ReadRebaseTodoFile(filepath.Join(gitDir, "rebase-merge", "done"), commentChar)
	if err != nil {
		return false
	}

	todos, err := utils.ReadRebaseTodoFile(filepath.Join(gitDir, "rebase-merge", "git-rebase-todo"), commentChar)
	if err != nil {
		return false
	}

	return utils.LastDoneTodoWasRescheduled(doneTodos, todos)
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
