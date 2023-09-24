package git_commands

import (
	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

type GitCommon struct {
	*common.Common
	version   *GitVersion
	cmd       oscommands.ICmdObjBuilder
	os        *oscommands.OSCommand
	repoPaths *RepoPaths
	repo      *gogit.Repository
	config    *ConfigCommands
}

func NewGitCommon(
	cmn *common.Common,
	version *GitVersion,
	cmd oscommands.ICmdObjBuilder,
	osCommand *oscommands.OSCommand,
	repoPaths *RepoPaths,
	repo *gogit.Repository,
	config *ConfigCommands,
) *GitCommon {
	return &GitCommon{
		Common:    cmn,
		version:   version,
		cmd:       cmd,
		os:        osCommand,
		repoPaths: repoPaths,
		repo:      repo,
		config:    config,
	}
}
