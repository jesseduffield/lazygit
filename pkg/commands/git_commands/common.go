package git_commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
)

type GitCommon struct {
	*common.Common
	version                   *GitVersion
	cmd                       oscommands.ICmdObjBuilder
	os                        *oscommands.OSCommand
	repoPaths                 *RepoPaths
	config                    *ConfigCommands
	diffRendererConfigManager *config.DiffRendererConfigManager
}

func NewGitCommon(
	cmn *common.Common,
	version *GitVersion,
	cmd oscommands.ICmdObjBuilder,
	osCommand *oscommands.OSCommand,
	repoPaths *RepoPaths,
	config *ConfigCommands,
	diffRendererConfigManager *config.DiffRendererConfigManager,
) *GitCommon {
	return &GitCommon{
		Common:                    cmn,
		version:                   version,
		cmd:                       cmd,
		os:                        osCommand,
		repoPaths:                 repoPaths,
		config:                    config,
		diffRendererConfigManager: diffRendererConfigManager,
	}
}
