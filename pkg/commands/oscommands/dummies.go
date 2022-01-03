package oscommands

import (
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// NewDummyOSCommand creates a new dummy OSCommand for testing
func NewDummyOSCommand() *OSCommand {
	osCmd := NewOSCommand(utils.NewDummyCommon(), dummyPlatform)

	return osCmd
}

func NewDummyCmdObjBuilder(runner ICmdObjRunner) *CmdObjBuilder {
	return &CmdObjBuilder{
		runner:    runner,
		logCmdObj: func(ICmdObj) {},
		platform:  dummyPlatform,
	}
}

var dummyPlatform = &Platform{
	OS:              "darwin",
	Shell:           "bash",
	ShellArg:        "-c",
	OpenCommand:     "open {{filename}}",
	OpenLinkCommand: "open {{link}}",
}

func NewDummyOSCommandWithRunner(runner *FakeCmdObjRunner) *OSCommand {
	osCommand := NewOSCommand(utils.NewDummyCommon(), dummyPlatform)
	osCommand.Cmd = NewDummyCmdObjBuilder(runner)

	return osCommand
}
