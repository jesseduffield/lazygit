package commands_test

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands"
	. "github.com/jesseduffield/lazygit/pkg/commands/commandsfakes"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands/oscommandsfakes"
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

type ExpectedRunCall struct {
	cmdStr    string
	outputStr string
	outputErr error
}

func WithRunCalls(
	commander *FakeICommander, expectedCalls []ExpectedRunCall, f func(),
) {
	i := 0
	commander.RunWithOutputCalls(func(cmdObj ICmdObj) (string, error) {
		// we shouldn't be calling this function any more times than we expect
		Expect(i).To(
			BeNumerically("<", len(expectedCalls)),
			"Unexpected call of RunWithOutput:\n\t%s\n\nThis means the function you're testing has attempted to run more commands than expected. If the function is doing its job properly, you'll need to append the command to the `WithRunCalls` call", cmdObj.ToString(),
		)

		call := expectedCalls[i]
		i += 1

		Expect(cmdObj.ToString()).To(Equal(call.cmdStr))
		return call.outputStr, call.outputErr
	})

	f()

	Expect(i).To(
		BeNumerically("==", len(expectedCalls)),
		func() string {
			missingCalls := make([]string, len(expectedCalls)-i)
			for j, expectedCall := range expectedCalls[i:] {
				missingCalls[j] = expectedCall.cmdStr
			}
			return fmt.Sprintf("Received fewer calls to RunWithOutput than expected. Missing calls:\n%s", strings.Join(missingCalls, "\n"))
		},
	)
}

func NewFakeMgrCtx(commander *FakeICommander, config *FakeIGitConfigMgr, os *oscommandsfakes.FakeIOS) *commands.MgrCtx {
	if config == nil {
		config = &FakeIGitConfigMgr{}
	}

	if commander == nil {
		commander = &FakeICommander{}
	}

	if os == nil {
		os = &oscommandsfakes.FakeIOS{}
	}

	return commands.NewMgrCtx(commander, config, nil, utils.NewDummyLog(), os, i18n.EnglishTranslationSet())
}

func SuccessCall(cmdStr string) ExpectedRunCall {
	return ExpectedRunCall{cmdStr: cmdStr}
}
