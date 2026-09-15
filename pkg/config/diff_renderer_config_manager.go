package config

import (
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type DiffRendererConfigManager struct {
	getUserConfig     func() *UserConfig
	diffRendererIndex int
}

type DiffRendererType int

const (
	DiffRendererType_StdinFilter DiffRendererType = iota
	DiffRendererType_ExtDiff
	DiffRendererType_RawGit
)

func NewDiffRendererConfigManager(getUserConfig func() *UserConfig) *DiffRendererConfigManager {
	return &DiffRendererConfigManager{getUserConfig: getUserConfig}
}

func (self *DiffRendererConfigManager) currentDiffRendererConfig() *DiffRendererConfig {
	diffRenderers := self.getUserConfig().Git.DiffRenderers
	if len(diffRenderers) == 0 {
		return nil
	}

	// Guard against the diff renderer index being out of range, which can happen if the user
	// has removed diff renderers from their config file while lazygit is running.
	if self.diffRendererIndex >= len(diffRenderers) {
		self.diffRendererIndex = 0
	}

	return &diffRenderers[self.diffRendererIndex]
}

func (self *DiffRendererConfig) getType() DiffRendererType {
	switch self.Type {
	case "stdinFilter", "":
		return DiffRendererType_StdinFilter
	case "extDiff":
		return DiffRendererType_ExtDiff
	case "rawGit":
		return DiffRendererType_RawGit
	}
	panic("invalid diff renderer type: " + self.Type)
}

func (self *DiffRendererConfigManager) GetDiffRendererType() DiffRendererType {
	currentDiffRendererConfig := self.currentDiffRendererConfig()
	if currentDiffRendererConfig == nil {
		return DiffRendererType_RawGit
	}
	return currentDiffRendererConfig.getType()
}

func (self *DiffRendererConfigManager) GetStdinFilterCommand(width int) string {
	currentDiffRendererConfig := self.currentDiffRendererConfig()
	if currentDiffRendererConfig == nil || currentDiffRendererConfig.getType() != DiffRendererType_StdinFilter {
		return ""
	}

	templateValues := map[string]string{
		"columnWidth": strconv.Itoa(width/2 - 6),
	}

	commandTemplate := string(currentDiffRendererConfig.Command)
	return utils.ResolvePlaceholderString(commandTemplate, templateValues)
}

func (self *DiffRendererConfigManager) GetColorArg() string {
	currentDiffRendererConfig := self.currentDiffRendererConfig()
	if currentDiffRendererConfig == nil || currentDiffRendererConfig.getType() != DiffRendererType_StdinFilter {
		return "always"
	}

	colorArg := currentDiffRendererConfig.ColorArg
	if colorArg == "" {
		return "always"
	}
	return colorArg
}

func (self *DiffRendererConfigManager) GetExternalDiffCommand(diffContext uint64) string {
	currentDiffRendererConfig := self.currentDiffRendererConfig()
	if currentDiffRendererConfig == nil || currentDiffRendererConfig.getType() != DiffRendererType_ExtDiff {
		return ""
	}

	templateValues := map[string]string{
		"diffContext": strconv.Itoa(int(diffContext)),
	}

	return utils.ResolvePlaceholderString(string(currentDiffRendererConfig.Command), templateValues)
}

func (self *DiffRendererConfigManager) GetRawGitArgs() []string {
	currentDiffRendererConfig := self.currentDiffRendererConfig()
	if currentDiffRendererConfig == nil || currentDiffRendererConfig.getType() != DiffRendererType_RawGit {
		return nil
	}
	return currentDiffRendererConfig.Args
}

func (self *DiffRendererConfigManager) CycleDiffRenderers() {
	self.diffRendererIndex = (self.diffRendererIndex + 1) % len(self.getUserConfig().Git.DiffRenderers)
}

func (self *DiffRendererConfigManager) CycleDiffRenderersBackward() {
	n := len(self.getUserConfig().Git.DiffRenderers)
	self.diffRendererIndex = (self.diffRendererIndex - 1 + n) % n
}

func (self *DiffRendererConfigManager) CurrentDiffRendererIndex() (int, int) {
	return self.diffRendererIndex, len(self.getUserConfig().Git.DiffRenderers)
}

// CurrentDiffRendererName returns a name for the current diff renderer, suitable for showing
// to the user.
func (self *DiffRendererConfigManager) CurrentDiffRendererName(tr *i18n.TranslationSet) string {
	currentDiffRendererConfig := self.currentDiffRendererConfig()
	if currentDiffRendererConfig == nil {
		return tr.DefaultDiffRendererName
	}

	if name := currentDiffRendererConfig.displayName(); name != "" {
		return name
	}

	if currentDiffRendererConfig.getType() == DiffRendererType_ExtDiff && currentDiffRendererConfig.Command == "" {
		return tr.ExternalDiffDiffRendererName
	}

	return tr.DefaultDiffRendererName
}

func (self *DiffRendererConfig) displayName() string {
	if self.Name != "" {
		return self.Name
	}
	if self.getType() == DiffRendererType_RawGit && len(self.Args) > 0 {
		return self.Args[0]
	}
	return firstWord(string(self.Command))
}

func firstWord(command string) string {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}
