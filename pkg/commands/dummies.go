package commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// NewDummyGitCommand creates a new dummy GitCommand for testing
func NewDummyGitCommand() *GitCommand {
	return NewDummyGitCommandWithOSCommand(oscommands.NewDummyOSCommand())
}

// NewDummyGitCommandWithOSCommand creates a new dummy GitCommand for testing
func NewDummyGitCommandWithOSCommand(osCommand *oscommands.OSCommand) *GitCommand {
	return NewGitCommandAux(
		utils.NewDummyCommon(),
		osCommand,
		utils.NewDummyGitConfig(),
		".git",
		nil,
	)
}

func NewDummyGitCommandWithRunner(runner oscommands.ICmdObjRunner) *GitCommand {
	builder := oscommands.NewDummyCmdObjBuilder(runner)
	gitCommand := NewDummyGitCommand()
	gitCommand.Cmd = builder
	gitCommand.OSCommand.Cmd = builder

	return gitCommand
}
