package git_commands

import (
	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/go-git/v5/config"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/common"
)

type ConfigCommands struct {
	*common.Common

	gitConfig git_config.IGitConfig
	repo      *gogit.Repository
}

func NewConfigCommands(
	common *common.Common,
	gitConfig git_config.IGitConfig,
	repo *gogit.Repository,
) *ConfigCommands {
	return &ConfigCommands{
		Common:    common,
		gitConfig: gitConfig,
		repo:      repo,
	}
}

type GpgConfigKey string

const (
	CommitGpgSign GpgConfigKey = "commit.gpgSign"
	TagGpgSign    GpgConfigKey = "tag.gpgSign"
)

// NeedsGpgSubprocess tells us whether the user has gpg enabled for the specified action type
// and needs a subprocess because they have a process where they manually
// enter their password every time a GPG action is taken
func (self *ConfigCommands) NeedsGpgSubprocess(key GpgConfigKey) bool {
	overrideGpg := self.UserConfig().Git.OverrideGpg
	if overrideGpg {
		return false
	}

	return self.gitConfig.GetBool(string(key))
}

func (self *ConfigCommands) NeedsGpgSubprocessForCommit() bool {
	return self.NeedsGpgSubprocess(CommitGpgSign)
}

func (self *ConfigCommands) GetGpgTagSign() bool {
	return self.gitConfig.GetBool(string(TagGpgSign))
}

func (self *ConfigCommands) GetCoreEditor() string {
	return self.gitConfig.Get("core.editor")
}

// GetRemoteURL returns current repo remote url
func (self *ConfigCommands) GetRemoteURL() string {
	return self.gitConfig.Get("remote.origin.url")
}

func (self *ConfigCommands) GetShowUntrackedFiles() string {
	return self.gitConfig.Get("status.showUntrackedFiles")
}

// this determines whether the user has configured to push to the remote branch of the same name as the current or not
func (self *ConfigCommands) GetPushToCurrent() bool {
	return self.gitConfig.Get("push.default") == "current"
}

// returns the repo's branches as specified in the git config
func (self *ConfigCommands) Branches() (map[string]*config.Branch, error) {
	conf, err := self.repo.Config()
	if err != nil {
		return nil, err
	}

	return conf.Branches, nil
}

func (self *ConfigCommands) GetGitFlowPrefixes() string {
	return self.gitConfig.GetGeneral("--local --get-regexp gitflow.prefix")
}

func (self *ConfigCommands) GetCoreCommentChar() byte {
	if commentCharStr := self.gitConfig.Get("core.commentChar"); len(commentCharStr) == 1 {
		return commentCharStr[0]
	}

	return '#'
}

func (self *ConfigCommands) GetRebaseUpdateRefs() bool {
	return self.gitConfig.GetBool("rebase.updateRefs")
}

func (self *ConfigCommands) GetMergeFF() string {
	return self.gitConfig.Get("merge.ff")
}

func (self *ConfigCommands) DropConfigCache() {
	self.gitConfig.DropCache()
}
