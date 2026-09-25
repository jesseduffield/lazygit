package config

import (
	"fmt"
	"strings"
	"text/template"

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

// DiffRendererValues are what the command of a diff renderer can refer to.
type DiffRendererValues struct {
	// The width of the view that the diff is rendered into
	Width int
	// The number of lines of context around each hunk
	DiffContext uint64
}

func (self *DiffRendererConfigManager) GetStdinFilterCommand(values DiffRendererValues) (string, error) {
	currentDiffRendererConfig := self.currentDiffRendererConfig()
	if currentDiffRendererConfig == nil || currentDiffRendererConfig.getType() != DiffRendererType_StdinFilter {
		return "", nil
	}

	return currentDiffRendererConfig.resolveCommand(values)
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

func (self *DiffRendererConfigManager) GetExternalDiffCommand(values DiffRendererValues) (string, error) {
	currentDiffRendererConfig := self.currentDiffRendererConfig()
	if currentDiffRendererConfig == nil || currentDiffRendererConfig.getType() != DiffRendererType_ExtDiff {
		return "", nil
	}

	return currentDiffRendererConfig.resolveCommand(values)
}

// resolveCommand resolves the renderer's command, which is a Go template with
// the values it can refer to as its variables. A variable can be written with
// or without the leading dot, as in {{.width}} or {{width}}.
func (self *DiffRendererConfig) resolveCommand(values DiffRendererValues) (string, error) {
	variables := map[string]any{
		"width": values.Width,
	}
	switch self.getType() {
	case DiffRendererType_StdinFilter:
		variables["columnWidth"] = values.Width/2 - 6
	case DiffRendererType_ExtDiff:
		variables["diffContext"] = values.DiffContext
	case DiffRendererType_RawGit:
		// has no command
	}

	funcs := template.FuncMap{}
	for name, value := range variables {
		funcs[name] = func() any { return value }
	}

	command, err := utils.ResolveTemplate(string(self.Command), variables, funcs)
	if err != nil {
		return "", fmt.Errorf("git.diffRenderers: can't use the command '%s': %w", self.Command, err)
	}
	return command, nil
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
