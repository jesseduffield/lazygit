package commands

import (
	"os"
	"strconv"
	"strings"

	gogit "github.com/jesseduffield/go-git/v5"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

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
	GetDotGitDir() string
}

type getGitConfigValueFunc func(key string) (string, error)

type GitConfigMgr struct {
	ICommander

	// Push to current determines whether the user has configured to push to the remote branch of the same name as the current or not
	pushToCurrent bool

	repo              *gogit.Repository
	userConfig        *config.UserConfig
	userConfigDir     string
	getGitConfigValue getGitConfigValueFunc
	debug             bool
	dotGitDir         string
}

func NewGitConfigMgr(commander ICommander, userConfig *config.UserConfig, userConfigDir string, getGitConfigValue getGitConfigValueFunc, log *logrus.Entry, repo *gogit.Repository, debug bool, dotGitDir string) *GitConfigMgr {
	gitConfig := &GitConfigMgr{
		ICommander:        commander,
		getGitConfigValue: getGitConfigValue,
		userConfig:        userConfig,
		userConfigDir:     userConfigDir,
		repo:              repo,
		debug:             debug,
		dotGitDir:         dotGitDir,
	}

	output, err := commander.RunWithOutput(
		BuildGitCmdObjFromStr("config --get push.default"),
	)
	pushToCurrent := false
	if err != nil {
		log.Errorf("error reading git config: %v", err)
	} else {
		pushToCurrent = strings.TrimSpace(output) == "current"
	}

	gitConfig.pushToCurrent = pushToCurrent

	return gitConfig
}

func (c *GitConfigMgr) GetDotGitDir() string {
	return c.dotGitDir
}

func (c *GitConfigMgr) GetUserConfig() *config.UserConfig {
	return c.userConfig
}

func (c *GitConfigMgr) GetUserConfigDir() string {
	return c.userConfigDir
}

func (c *GitConfigMgr) GetPager(width int) string {
	useConfig := c.userConfig.Git.Paging.UseConfig
	if useConfig {
		pager := c.configuredPager()
		return strings.Split(pager, "| less")[0]
	}

	templateValues := map[string]string{
		"columnWidth": strconv.Itoa(width/2 - 6),
	}

	pagerTemplate := c.userConfig.Git.Paging.Pager
	return utils.ResolvePlaceholderString(pagerTemplate, templateValues)
}

func (c *GitConfigMgr) ColorArg() string {
	return c.userConfig.Git.Paging.ColorArg
}

func (c *GitConfigMgr) GetConfigValue(key string) string {
	output, _ := c.getGitConfigValue(key)
	return output
}

func (c *GitConfigMgr) configuredPager() string {
	if os.Getenv("GIT_PAGER") != "" {
		return os.Getenv("GIT_PAGER")
	}
	if os.Getenv("PAGER") != "" {
		return os.Getenv("PAGER")
	}
	output, err := c.RunWithOutput(BuildGitCmdObjFromStr("config --get-all core.pager"))
	if err != nil {
		return ""
	}
	trimmedOutput := strings.TrimSpace(output)
	return strings.Split(trimmedOutput, "\n")[0]
}

// UsingGpg tells us whether the user has gpg enabled so that we can know
// whether we need to run a subprocess to allow them to enter their password
func (c *GitConfigMgr) UsingGpg() bool {
	overrideGpg := c.userConfig.Git.OverrideGpg
	if overrideGpg {
		return false
	}

	gpgsign := c.GetConfigValue("commit.gpgsign")
	value := strings.ToLower(gpgsign)

	return value == "true" || value == "1" || value == "yes" || value == "on"
}

func (c *GitConfigMgr) FindRemoteForBranchInConfig(branchName string) (string, error) {
	conf, err := c.repo.Config()
	if err != nil {
		return "", err
	}

	for configBranchName, configBranch := range conf.Branches {
		if configBranchName == branchName {
			return configBranch.Remote, nil
		}
	}

	return "", nil
}

func (c *GitConfigMgr) GetPushToCurrent() bool {
	return c.pushToCurrent
}

func (c *GitConfigMgr) GetDebug() bool {
	return c.debug
}
