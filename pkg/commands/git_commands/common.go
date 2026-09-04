package git_commands

import (
	"os"
	"path/filepath"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
)

type GitCommon struct {
	*common.Common
	version     *GitVersion
	cmd         oscommands.ICmdObjBuilder
	os          *oscommands.OSCommand
	repoPaths   *RepoPaths
	config      *ConfigCommands
	diffRendererConfigManager *config.DiffRendererConfigManager
	IsGitSvnRepo bool
}

func (self *GitCommon) detectGitSvnRepo() {
	if self.Common != nil && !self.Common.UserConfig().Git.EnableGitSvnCompat {
		self.IsGitSvnRepo = false
		return
	}

	if self.repoPaths == nil {
		self.IsGitSvnRepo = false
		return
	}

	svnDir := filepath.Join(self.repoPaths.RepoGitDirPath(), "svn")
	if info, err := os.Stat(svnDir); err == nil && info.IsDir() {
		self.IsGitSvnRepo = true
		self.Common.Log.Info("Detected Git-SVN repository (found .git/svn)")
	} else {
		self.IsGitSvnRepo = false
	}
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
	gitCommon := &GitCommon{
		Common:      cmn,
		version:     version,
		cmd:         cmd,
		os:          osCommand,
		repoPaths:   repoPaths,
		config:      config,
		diffRendererConfigManager: diffRendererConfigManager,
	}
	gitCommon.detectGitSvnRepo()
	return gitCommon
}

func (self *GitCommon) IsSvnRepo() bool {
	return self.IsGitSvnRepo
}
