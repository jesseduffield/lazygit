package config

import (
	"strconv"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

type PagerConfig struct {
	getUserConfig func() *UserConfig
}

func NewPagerConfig(getUserConfig func() *UserConfig) *PagerConfig {
	return &PagerConfig{getUserConfig: getUserConfig}
}

func (self *PagerConfig) GetPagerCommand(width int) string {
	templateValues := map[string]string{
		"columnWidth": strconv.Itoa(width/2 - 6),
	}

	pagerTemplate := string(self.getUserConfig().Git.Paging.Pager)
	return utils.ResolvePlaceholderString(pagerTemplate, templateValues)
}

func (self *PagerConfig) GetColorArg() string {
	return self.getUserConfig().Git.Paging.ColorArg
}

func (self *PagerConfig) GetExternalDiffCommand() string {
	return self.getUserConfig().Git.Paging.ExternalDiffCommand
}

func (self *PagerConfig) GetUseExternalDiffGitConfig() bool {
	return self.getUserConfig().Git.Paging.UseExternalDiffGitConfig
}
