package commands

import (
	"fmt"

	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

func (c *Git) FlowStart(branchType string, name string) ICmdObj {
	return BuildGitCmdObjFromStr(fmt.Sprintf("flow %s start %s", branchType, name))
}

func (c *Git) FlowFinish(branchType string, name string) ICmdObj {
	return BuildGitCmdObjFromStr(fmt.Sprintf("flow %s finish %s", branchType, name))
}

func (c *Git) GetGitFlowRegexpConfig() (string, error) {
	return c.RunWithOutput(BuildGitCmdObjFromStr("config --local --get-regexp gitflow"))
}
