package commands

import (
	"io"
	"io/ioutil"

	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// NewDummyGitCommand creates a new dummy GitCommand for testing
func NewDummyGitCommand() *GitCommand {
	return NewDummyGitCommandWithOSCommand(oscommands.NewDummyOSCommand())
}

// NewDummyGitCommandWithOSCommand creates a new dummy GitCommand for testing
func NewDummyGitCommandWithOSCommand(osCommand *oscommands.OSCommand) *GitCommand {
	runner := &oscommands.FakeCmdObjRunner{}
	builder := oscommands.NewDummyCmdObjBuilder(runner)

	return &GitCommand{
		Common:       utils.NewDummyCommon(),
		Cmd:          builder,
		OSCommand:    osCommand,
		GitConfig:    git_config.NewFakeGitConfig(map[string]string{}),
		GetCmdWriter: func() io.Writer { return ioutil.Discard },
	}
}

func NewDummyGitCommandWithRunner(runner oscommands.ICmdObjRunner) *GitCommand {
	builder := oscommands.NewDummyCmdObjBuilder(runner)
	gitCommand := NewDummyGitCommand()
	gitCommand.Cmd = builder
	gitCommand.OSCommand.Cmd = builder

	return gitCommand
}
