package commands

import (
	"os"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (c *GitCommand) ConfiguredPager() string {
	if os.Getenv("GIT_PAGER") != "" {
		return os.Getenv("GIT_PAGER")
	}
	if os.Getenv("PAGER") != "" {
		return os.Getenv("PAGER")
	}
	output, err := c.RunCommandWithOutput(BuildGitCmdObjFromStr("config --get-all core.pager"))
	if err != nil {
		return ""
	}
	trimmedOutput := strings.TrimSpace(output)
	return strings.Split(trimmedOutput, "\n")[0]
}

func (c *GitCommand) GetPager(width int) string {
	useConfig := c.config.GetUserConfig().Git.Paging.UseConfig
	if useConfig {
		pager := c.ConfiguredPager()
		return strings.Split(pager, "| less")[0]
	}

	templateValues := map[string]string{
		"columnWidth": strconv.Itoa(width/2 - 6),
	}

	pagerTemplate := c.config.GetUserConfig().Git.Paging.Pager
	return utils.ResolvePlaceholderString(pagerTemplate, templateValues)
}

func (c *GitCommand) colorArg() string {
	return c.config.GetUserConfig().Git.Paging.ColorArg
}

func (c *GitCommand) GetConfigValue(key string) string {
	output, _ := c.getGitConfigValue(key)
	return output
}

// UsingGpg tells us whether the user has gpg enabled so that we can know
// whether we need to run a subprocess to allow them to enter their password
func (c *GitCommand) UsingGpg() bool {
	overrideGpg := c.config.GetUserConfig().Git.OverrideGpg
	if overrideGpg {
		return false
	}

	gpgsign := c.GetConfigValue("commit.gpgsign")
	value := strings.ToLower(gpgsign)

	return value == "true" || value == "1" || value == "yes" || value == "on"
}

func (c *GitCommand) FindRemoteForBranchInConfig(branchName string) (string, error) {
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
