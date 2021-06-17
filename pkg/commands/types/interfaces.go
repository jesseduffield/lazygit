package types

import "github.com/jesseduffield/lazygit/pkg/config"

//counterfeiter:generate . IGitConfigMgr
type IGitConfigMgr interface {
	GetPager(width int) string
	ColorArg() string
	GetConfigValue(key string) string
	UsingGpg() bool
	GetUserConfig() *config.UserConfig
	GetPushToCurrent() bool
	GetUserConfigDir() string
	FindRemoteForBranchInConfig(branchName string) (string, error)
	GetDebug() bool
}

//counterfeiter:generate . ICommander
type ICommander interface {
	Run(cmdObj ICmdObj) error
	RunWithOutput(cmdObj ICmdObj) (string, error)
	RunGitCmdFromStr(cmdStr string) error
	BuildGitCmdObjFromStr(cmdStr string) ICmdObj
	BuildShellCmdObj(command string) ICmdObj
	SkipEditor(cmdObj ICmdObj)
	Quote(string) string
}
