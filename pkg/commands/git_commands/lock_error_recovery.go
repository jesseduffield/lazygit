package git_commands

import (
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
