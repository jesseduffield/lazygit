package commands

import (
	"fmt"

	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

func (c *GitCommand) FlowStart(branchType string, name string) ICmdObj {
	return BuildGitCmdObjFromStr(fmt.Sprintf("flow %s start %s", branchType, name))
}

func (c *GitCommand) FlowFinish(branchType string, name string) ICmdObj {
	return BuildGitCmdObjFromStr(fmt.Sprintf("flow %s finish %s", branchType, name))
}

func (c *GitCommand) GetGitFlowRegexpConfig() (string, error) {
	return c.RunCommandWithOutput(BuildGitCmdObjFromStr("config --local --get-regexp gitflow"))
}
