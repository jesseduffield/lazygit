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
	output := c.GitConfig.Get("core.pager")
	return strings.Split(output, "\n")[0]
}

func (c *GitCommand) GetPager(width int) string {
	useConfig := c.UserConfig.Git.Paging.UseConfig
	if useConfig {
		pager := c.ConfiguredPager()
		return strings.Split(pager, "| less")[0]
	}

	templateValues := map[string]string{
		"columnWidth": strconv.Itoa(width/2 - 6),
	}

	pagerTemplate := c.UserConfig.Git.Paging.Pager
	return utils.ResolvePlaceholderString(pagerTemplate, templateValues)
}

func (c *GitCommand) colorArg() string {
	return c.UserConfig.Git.Paging.ColorArg
}

// UsingGpg tells us whether the user has gpg enabled so that we can know
// whether we need to run a subprocess to allow them to enter their password
func (c *GitCommand) UsingGpg() bool {
	overrideGpg := c.UserConfig.Git.OverrideGpg
	if overrideGpg {
		return false
	}

	return c.GitConfig.GetBool("commit.gpgsign")
}
