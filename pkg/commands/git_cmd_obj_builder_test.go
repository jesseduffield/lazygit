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
	)

	assert.Contains(t, builder.New([]string{"git", "status"}).GetEnvVars(), git_commands.OptionalLocksEnvVar+"=0")
	assert.Contains(t, builder.NewShell("git status", "").GetEnvVars(), git_commands.OptionalLocksEnvVar+"=0")
}
