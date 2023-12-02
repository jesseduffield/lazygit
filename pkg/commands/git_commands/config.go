package git_commands

import (
	"os"
	"strconv"
	"strings"

	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/go-git/v5/config"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
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

func (self *ConfigCommands) ConfiguredPager() string {
	if os.Getenv("GIT_PAGER") != "" {
		return os.Getenv("GIT_PAGER")
	}
	if os.Getenv("PAGER") != "" {
		return os.Getenv("PAGER")
	}
	output := self.gitConfig.Get("core.pager")
	return strings.Split(output, "\n")[0]
}

func (self *ConfigCommands) GetPager(width int) string {
	useConfig := self.UserConfig.Git.Paging.UseConfig
	if useConfig {
		pager := self.ConfiguredPager()
		return strings.Split(pager, "| less")[0]
	}

	templateValues := map[string]string{
		"columnWidth": strconv.Itoa(width/2 - 6),
	}

	pagerTemplate := string(self.UserConfig.Git.Paging.Pager)
	return utils.ResolvePlaceholderString(pagerTemplate, templateValues)
}

// UsingGpg tells us whether the user has gpg enabled so that we can know
// whether we need to run a subprocess to allow them to enter their password
func (self *ConfigCommands) UsingGpg() bool {
	overrideGpg := self.UserConfig.Git.OverrideGpg
	if overrideGpg {
		return false
	}

	return self.gitConfig.GetBool("commit.gpgsign")
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
