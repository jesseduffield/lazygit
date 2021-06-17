package commands

import (
	"fmt"

	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

//counterfeiter:generate . IFlowMgr
type IFlowMgr interface {
	Start(branchType string, name string) ICmdObj
	Finish(branchType string, name string) ICmdObj
	GetGitFlowRegexpConfig() (string, error)
}

type FlowMgr struct {
	ICommander

	config IGitConfigMgr
}

func NewFlowMgr(
	commander ICommander,
	config IGitConfigMgr,
) *FlowMgr {
	return &FlowMgr{
		ICommander: commander,
		config:     config,
	}
}

func (c *FlowMgr) Start(branchType string, name string) ICmdObj {
	return BuildGitCmdObjFromStr(fmt.Sprintf("flow %s start %s", branchType, name))
}

func (c *FlowMgr) Finish(branchType string, name string) ICmdObj {
	return BuildGitCmdObjFromStr(fmt.Sprintf("flow %s finish %s", branchType, name))
}

func (c *FlowMgr) GetGitFlowRegexpConfig() (string, error) {
	return c.RunWithOutput(BuildGitCmdObjFromStr("config --local --get-regexp gitflow"))
}
