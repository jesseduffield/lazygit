package config

import (
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

type PagerConfig struct {
	getUserConfig func() *UserConfig
	pagerIndex    int
}

func NewPagerConfig(getUserConfig func() *UserConfig) *PagerConfig {
	return &PagerConfig{getUserConfig: getUserConfig}
}

func (self *PagerConfig) currentPagerConfig() *PagingConfig {
	pagers := self.getUserConfig().Git.Pagers
	if len(pagers) == 0 {
		return nil
	}

	// Guard against the pager index being out of range, which can happen if the user
	// has removed pagers from their config file while lazygit is running.
	if self.pagerIndex >= len(pagers) {
		self.pagerIndex = 0
	}

	return &pagers[self.pagerIndex]
}

func (self *PagerConfig) GetPagerCommand(width int) string {
	currentPagerConfig := self.currentPagerConfig()
	if currentPagerConfig == nil {
		return ""
	}

	templateValues := map[string]string{
		"columnWidth": strconv.Itoa(width/2 - 6),
	}

	pagerTemplate := string(currentPagerConfig.Pager)
	return utils.ResolvePlaceholderString(pagerTemplate, templateValues)
}

func (self *PagerConfig) GetColorArg() string {
	currentPagerConfig := self.currentPagerConfig()
	if currentPagerConfig == nil {
		return "always"
	}

	colorArg := currentPagerConfig.ColorArg
	if colorArg == "" {
		return "always"
	}
	return colorArg
}

func (self *PagerConfig) GetExternalDiffCommand() string {
	currentPagerConfig := self.currentPagerConfig()
	if currentPagerConfig == nil {
		return ""
	}
	return currentPagerConfig.ExternalDiffCommand
}

func (self *PagerConfig) GetUseExternalDiffGitConfig() bool {
	currentPagerConfig := self.currentPagerConfig()
	if currentPagerConfig == nil {
		return false
	}
	return currentPagerConfig.UseExternalDiffGitConfig
}

func (self *PagerConfig) CyclePagers() {
	self.pagerIndex = (self.pagerIndex + 1) % len(self.getUserConfig().Git.Pagers)
}

func (self *PagerConfig) CurrentPagerIndex() (int, int) {
	return self.pagerIndex, len(self.getUserConfig().Git.Pagers)
}

// CurrentPagerName returns a name for the current pager, suitable for showing
// to the user. It returns an empty string if no name can be derived; callers
// should substitute a localized fallback in that case.
func (self *PagerConfig) CurrentPagerName() string {
	currentPagerConfig := self.currentPagerConfig()
	if currentPagerConfig == nil {
		return ""
	}
	return currentPagerConfig.displayName()
}

// CurrentPagerUsesGitConfigDiff reports whether the current pager defers to
// git's own external diff config. Such an entry has no name we can derive (the
// actual command may even vary per file via .gitattributes), so callers show a
// generic label rather than treating it like the default no-pager entry.
func (self *PagerConfig) CurrentPagerUsesGitConfigDiff() bool {
	currentPagerConfig := self.currentPagerConfig()
	return currentPagerConfig != nil && currentPagerConfig.UseExternalDiffGitConfig
}

func (self *PagingConfig) displayName() string {
	if self.Name != "" {
		return self.Name
	}
	if word := firstWord(string(self.Pager)); word != "" {
		return word
	}
	return firstWord(self.ExternalDiffCommand)
}

func firstWord(command string) string {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}
