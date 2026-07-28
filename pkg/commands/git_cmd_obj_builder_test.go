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
	)

	assert.Equal(t, "/path/to/repo", builder.New([]string{"git", "status"}).GetCmd().Dir)
	assert.Equal(t, "/path/to/repo", builder.NewShell("git status", "").GetCmd().Dir)
}
