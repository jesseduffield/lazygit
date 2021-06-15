package commands_test

import (
	"fmt"

	. "github.com/jesseduffield/lazygit/pkg/commands/commandsfakes"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

var mockQuote = func(str string) string { return fmt.Sprintf("\"%s\"", str) }

func NewFakeCommander() *FakeICommander {
	commander := &FakeICommander{}
	commander.BuildGitCmdObjFromStrCalls(func(command string) ICmdObj {
		return oscommands.NewCmdObjFromStr("git " + command)
	})
	commander.RunGitCmdFromStrCalls(func(command string) error {
		return commander.Run(commander.BuildGitCmdObjFromStr((command)))
	})
	return commander
}
