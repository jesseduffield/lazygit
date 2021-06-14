package commands

import (
	"os"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

type getGitConfigValueFunc func(key string) (string, error)

type GitConfig struct {
	commander         *Commander
	pushToCurrent     bool
	userConfig        *config.UserConfig
	getGitConfigValue getGitConfigValueFunc
}

func NewGitConfig(commander *Commander, userConfig *config.UserConfig, getGitConfigValue getGitConfigValueFunc, log *logrus.Entry) *GitConfig {
	gitConfig := &GitConfig{
		commander:         commander,
		getGitConfigValue: getGitConfigValue,
		userConfig:        userConfig,
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

func (c *GitConfig) GetPager(width int) string {
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

func (c *GitConfig) colorArg() string {
	return c.userConfig.Git.Paging.ColorArg
}

func (c *GitConfig) GetConfigValue(key string) string {
	output, _ := c.getGitConfigValue(key)
	return output
}

func (c *GitConfig) configuredPager() string {
	if os.Getenv("GIT_PAGER") != "" {
		return os.Getenv("GIT_PAGER")
	}
	if os.Getenv("PAGER") != "" {
		return os.Getenv("PAGER")
	}
	output, err := c.commander.RunWithOutput(BuildGitCmdObjFromStr("config --get-all core.pager"))
	if err != nil {
		return ""
	}
	trimmedOutput := strings.TrimSpace(output)
	return strings.Split(trimmedOutput, "\n")[0]
}

// UsingGpg tells us whether the user has gpg enabled so that we can know
// whether we need to run a subprocess to allow them to enter their password
func (c *GitConfig) UsingGpg() bool {
	overrideGpg := c.userConfig.Git.OverrideGpg
	if overrideGpg {
		return false
	}

	gpgsign := c.GetConfigValue("commit.gpgsign")
	value := strings.ToLower(gpgsign)

	return value == "true" || value == "1" || value == "yes" || value == "on"
}

func (c *Git) FindRemoteForBranchInConfig(branchName string) (string, error) {
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
