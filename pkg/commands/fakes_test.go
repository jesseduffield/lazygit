package commands_test

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands"
	. "github.com/jesseduffield/lazygit/pkg/commands/commandsfakes"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/utils"
	. "github.com/onsi/gomega"
)

func NewFakeCommander() *FakeICommander {
	commander := &FakeICommander{}

	commander.BuildGitCmdObjFromStrCalls(func(command string) ICmdObj {
		return oscommands.NewCmdObjFromStr("git " + command)
	})

	commander.RunGitCmdFromStrCalls(func(command string) error {
		return commander.Run(commander.BuildGitCmdObjFromStr((command)))
	})

	commander.RunCalls(func(cmdObj ICmdObj) error {
		_, err := commander.RunWithOutput(cmdObj)
		return err
	})

	commander.QuoteCalls(func(str string) string {
		return fmt.Sprintf("\"%s\"", str)
	})

	return commander
}

type ExpectedRunWithOutputCall struct {
	cmdStr    string
	outputStr string
	outputErr error
}

func SetExpectedRunWithOutputCalls(commander *FakeICommander, expectedCalls []ExpectedRunWithOutputCall) {
	i := 0
	commander.RunWithOutputCalls(func(cmdObj ICmdObj) (string, error) {
		// we shouldn't be calling this function any more times than we expect
		Expect(i).To(BeNumerically("<", len(expectedCalls)))

		call := expectedCalls[i]
		i += 1

		Expect(cmdObj.ToString()).To(Equal(call.cmdStr))
		return call.outputStr, call.outputErr
	})
}

func NewFakeMgrCtx(commander *FakeICommander, config *FakeIGitConfigMgr) *commands.MgrCtx {
	return commands.NewMgrCtx(commander, config, nil, utils.NewDummyLog(), nil, i18n.EnglishTranslationSet())
}
