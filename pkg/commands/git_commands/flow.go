package git_commands

import (
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

type FlowCommands struct {
	*GitCommon
}

func NewFlowCommands(
	gitCommon *GitCommon,
) *FlowCommands {
	return &FlowCommands{
		GitCommon: gitCommon,
	}
}

func (self *FlowCommands) GitFlowEnabled() bool {
	return len(self.config.GetGitFlowPrefixMap()) > 0
}

func (self *FlowCommands) FinishCmdObj(branchName string) (*oscommands.CmdObj, error) {
	prefixMap := self.config.GetGitFlowPrefixMap()

	prefixPart, suffix, ok := strings.Cut(branchName, "/")
	if !ok || prefixPart == "" || suffix == "" {
		return nil, errors.New(self.Tr.NotAGitFlowBranch)
	}
	prefix := prefixPart + "/"

	branchType := prefixMap[prefix]
	if branchType == "" {
		return nil, errors.New(self.Tr.NotAGitFlowBranch)
	}

	cmdArgs := NewGitCmd("flow").Arg(branchType, "finish", suffix).ToArgv()

	return self.cmd.New(cmdArgs), nil
}

func (self *FlowCommands) StartCmdObj(branchType string, name string) *oscommands.CmdObj {
	cmdArgs := NewGitCmd("flow").Arg(branchType, "start", name).ToArgv()

	return self.cmd.New(cmdArgs)
}
