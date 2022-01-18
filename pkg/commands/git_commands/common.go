package git_commands

import (
	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

type GitCommon struct {
	*common.Common
	cmd       oscommands.ICmdObjBuilder
	os        *oscommands.OSCommand
	dotGitDir string
	repo      *gogit.Repository
	config    *ConfigCommands
}

func NewGitCommon(
	cmn *common.Common,
	cmd oscommands.ICmdObjBuilder,
	osCommand *oscommands.OSCommand,
	dotGitDir string,
	repo *gogit.Repository,
	config *ConfigCommands,
) *GitCommon {
	return &GitCommon{
		Common:    cmn,
		cmd:       cmd,
		os:        osCommand,
		dotGitDir: dotGitDir,
		repo:      repo,
		config:    config,
	}
}
