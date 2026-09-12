package git_commands

import (
	"os"
	"path/filepath"
	"sync"
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
	Svn         *SvnCommands
	// indexLock serializes index-writing git commands (git add, git reset, etc.)
	// with git svn operations (fetch, rebase, dcommit) that internally call
	// git update-index and take index.lock.
	// For non-SVN repos the lock is uncontended (nanosecond overhead).
	// For SVN repos: all git svn operations and staging commands take the lock
	// with Lock() (blocking wait). If a background fetch is running, staging
	// blocks on the UI thread until the fetch finishes (the fetch's "Fetching"
	// spinner is visible). The caller (press/toggleStagedAll) pauses background
	// refreshes for the duration to prevent a racing git status from
	// overwriting the optimistic render with the stale pre-add state.
	indexLock sync.Mutex
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
		if self.Common != nil {
			self.Common.Log.Info("Detected Git-SVN repository (found .git/svn)")
		}
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

// IndexLock returns the mutex used to serialize index-writing commands
// (git add/reset) with background git svn fetch.
func (self *GitCommon) IndexLock() *sync.Mutex {
	return &self.indexLock
}
