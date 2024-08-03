package oscommands

import (
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// NewDummyOSCommand creates a new dummy OSCommand for testing
func NewDummyOSCommand() *OSCommand {
	osCmd := NewOSCommand(utils.NewDummyCommon(), config.NewDummyAppConfig(), dummyPlatform, NewNullGuiIO(utils.NewDummyLog()))

	return osCmd
}

type OSCommandDeps struct {
	Common       *common.Common
	Platform     *Platform
	GetenvFn     func(string) string
	RemoveFileFn func(string) error
	Cmd          *CmdObjBuilder
	TempDir      string
}

func NewDummyOSCommandWithDeps(deps OSCommandDeps) *OSCommand {
	common := deps.Common
	if common == nil {
		common = utils.NewDummyCommon()
	}

	platform := deps.Platform
	if platform == nil {
		platform = dummyPlatform
	}

	return &OSCommand{
		Common:       common,
		Platform:     platform,
		getenvFn:     deps.GetenvFn,
		removeFileFn: deps.RemoveFileFn,
		guiIO:        NewNullGuiIO(utils.NewDummyLog()),
		tempDir:      deps.TempDir,
	}
}

func NewDummyCmdObjBuilder(runner ICmdObjRunner) *CmdObjBuilder {
	return &CmdObjBuilder{
		runner:   runner,
		platform: dummyPlatform,
	}
}

var dummyPlatform = &Platform{
	OS:                  "darwin",
	Shell:               "bash",
	InteractiveShell:    "bash",
	ShellArg:            "-c",
	InteractiveShellArg: "-i",
	OpenCommand:         "open {{filename}}",
	OpenLinkCommand:     "open {{link}}",
}

func NewDummyOSCommandWithRunner(runner *FakeCmdObjRunner) *OSCommand {
	osCommand := NewOSCommand(utils.NewDummyCommon(), config.NewDummyAppConfig(), dummyPlatform, NewNullGuiIO(utils.NewDummyLog()))
	osCommand.Cmd = NewDummyCmdObjBuilder(runner)

	return osCommand
}
