package oscommands

import (
	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// NewDummyOSCommand creates a new dummy OSCommand for testing
func NewDummyOSCommand() *OSCommand {
	return NewOSCommand(utils.NewDummyCommon())
}

func NewCmdObjBuilderDummy(runner ICmdObjRunner) ICmdObjBuilder {
	return &CmdObjBuilder{
		runner:    runner,
		logCmdObj: func(ICmdObj) {},
		command:   secureexec.Command,
		platform: &Platform{
			OS:              "darwin",
			Shell:           "bash",
			ShellArg:        "-c",
			OpenCommand:     "open {{filename}}",
			OpenLinkCommand: "open {{link}}",
		},
	}
}
