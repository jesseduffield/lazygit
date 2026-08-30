package commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// Every git command we build disables optional locks by default, so that our
// invocations never contend for index.lock (see git_commands.OptionalLocksEnvVar
// for the rationale). Commands that want the lock opt back in with
// CmdObj.RemoveEnvVar.
func TestGitCmdObjBuilderDisablesOptionalLocksByDefault(t *testing.T) {
	builder := NewGitCmdObjBuilder(
		utils.NewDummyLog(),
		oscommands.NewDummyCmdObjBuilder(oscommands.NewFakeRunner(t)),
		"/path/to/repo",
		nil,
	)

	assert.Contains(t, builder.New([]string{"git", "status"}).GetEnvVars(), git_commands.OptionalLocksEnvVar+"=0")
	assert.Contains(t, builder.NewShell("git status", "").GetEnvVars(), git_commands.OptionalLocksEnvVar+"=0")
}

// Every command the builder produces runs in the directory of the repo the
// builder was created for, not in the process's current directory: lazygit
// chdirs when switching repos, and commands built for the previous repo after
// that (e.g. by a background refresh still in flight) must keep addressing the
// repo they were built for.
func TestGitCmdObjBuilderPinsCommandsToRepoDir(t *testing.T) {
	builder := NewGitCmdObjBuilder(
		utils.NewDummyLog(),
		oscommands.NewDummyCmdObjBuilder(oscommands.NewFakeRunner(t)),
		"/path/to/repo",
		nil,
	)

	assert.Equal(t, "/path/to/repo", builder.New([]string{"git", "status"}).GetCmd().Dir)
	assert.Equal(t, "/path/to/repo", builder.NewShell("git status", "").GetCmd().Dir)
}

// A repo whose git dir isn't in its worktree can't be found by running a
// command there, so the builder has to tell every command where it is; see
// RepoPaths.GitLocationEnvVars. The process env says the same thing, but only
// for the repo lazygit is in right now, which isn't necessarily this one.
func TestGitCmdObjBuilderPinsCommandsToGitLocation(t *testing.T) {
	builder := NewGitCmdObjBuilder(
		utils.NewDummyLog(),
		oscommands.NewDummyCmdObjBuilder(oscommands.NewFakeRunner(t)),
		"/path/to/worktree",
		[]string{"GIT_DIR=/path/to/repo/.git", "GIT_WORK_TREE=/path/to/worktree"},
	)

	assert.Subset(t, builder.New([]string{"git", "status"}).GetEnvVars(),
		[]string{"GIT_DIR=/path/to/repo/.git", "GIT_WORK_TREE=/path/to/worktree"})
	assert.Subset(t, builder.NewShell("git status", "").GetEnvVars(),
		[]string{"GIT_DIR=/path/to/repo/.git", "GIT_WORK_TREE=/path/to/worktree"})
}
