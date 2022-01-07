package commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

type CustomCommands struct {
	*common.Common

	cmd oscommands.ICmdObjBuilder
}

func NewCustomCommands(
	common *common.Common,
	cmd oscommands.ICmdObjBuilder,
) *CustomCommands {
	return &CustomCommands{
		Common: common,
		cmd:    cmd,
	}
}

// Only to be used for the sake of running custom commands specified by the user.
// If you want to run a new command, try finding a place for it in one of the neighbouring
// files, or creating a new BlahCommands struct to hold it.
func (self *CustomCommands) RunWithOutput(cmdStr string) (string, error) {
	return self.cmd.New(cmdStr).RunWithOutput()
}
