package config

import (
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestUserConfigValidate_enums(t *testing.T) {
	type testCase struct {
		value string
		valid bool
	}

	scenarios := []struct {
		name      string
		setup     func(config *UserConfig, value string)
		testCases []testCase
	}{
		{
			name: "Gui.StatusPanelView",
			setup: func(config *UserConfig, value string) {
				config.Gui.StatusPanelView = value
			},
			testCases: []testCase{
				{value: "dashboard", valid: true},
				{value: "allBranchesLog", valid: true},
				{value: "", valid: false},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Gui.ShowDivergenceFromBaseBranch",
			setup: func(config *UserConfig, value string) {
				config.Gui.ShowDivergenceFromBaseBranch = value
			},
			testCases: []testCase{
				{value: "none", valid: true},
				{value: "onlyArrow", valid: true},
				{value: "arrowAndNumber", valid: true},
				{value: "", valid: false},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Git.AutoForwardBranches",
			setup: func(config *UserConfig, value string) {
				config.Git.AutoForwardBranches = value
			},
			testCases: []testCase{
				{value: "none", valid: true},
				{value: "onlyMainBranches", valid: true},
				{value: "allBranches", valid: true},
				{value: "", valid: false},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Git.LocalBranchSortOrder",
			setup: func(config *UserConfig, value string) {
				config.Git.LocalBranchSortOrder = value
			},
			testCases: []testCase{
				{value: "date", valid: true},
				{value: "recency", valid: true},
				{value: "alphabetical", valid: true},
				{value: "", valid: false},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Git.RemoteBranchSortOrder",
			setup: func(config *UserConfig, value string) {
				config.Git.RemoteBranchSortOrder = value
			},
			testCases: []testCase{
				{value: "date", valid: true},
				{value: "recency", valid: false},
				{value: "alphabetical", valid: true},
				{value: "", valid: false},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Git.Log.Order",
			setup: func(config *UserConfig, value string) {
				config.Git.Log.Order = value
			},
			testCases: []testCase{
				{value: "date-order", valid: true},
				{value: "author-date-order", valid: true},
				{value: "topo-order", valid: true},
				{value: "default", valid: true},

				{value: "", valid: false},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Git.Log.ShowGraph",
			setup: func(config *UserConfig, value string) {
				config.Git.Log.ShowGraph = value
			},
			testCases: []testCase{
				{value: "always", valid: true},
				{value: "never", valid: true},
				{value: "when-maximised", valid: true},

				{value: "", valid: false},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Keybindings",
			setup: func(config *UserConfig, value string) {
				config.Keybinding.Universal.Quit = Keybinding{value}
			},
			testCases: []testCase{
				{value: "", valid: true},
				{value: "<disabled>", valid: true},
				{value: "q", valid: true},
				{value: "<c-c>", valid: true},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "JumpToBlock keybinding",
			setup: func(config *UserConfig, value string) {
				labels := strings.Split(value, ",")
				config.Keybinding.Universal.JumpToBlock = lo.Map(labels, func(label string, _ int) Keybinding {
					return Keybinding{label}
				})
			},
			testCases: []testCase{
				{value: "", valid: false},
				{value: "1,2,3", valid: false},
				{value: "1,2,3,4,5", valid: true},
				{value: "1,2,3,4,invalid", valid: false},
				{value: "1,2,3,4,5,6", valid: false},
			},
		},
		{
			name: "Custom command keybinding",
			setup: func(config *UserConfig, value string) {
				config.CustomCommands = []CustomCommand{
					{
						Key:     Keybinding{value},
						Command: "echo 'hello'",
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: true},
				{value: "<disabled>", valid: true},
				{value: "q", valid: true},
				{value: "<c-c>", valid: true},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Custom command keybinding in sub menu",
			setup: func(config *UserConfig, value string) {
				config.CustomCommands = []CustomCommand{
					{
						Key:         Keybinding{"X"},
						Description: "My Custom Commands",
						CommandMenu: []CustomCommand{
							{Key: Keybinding{value}, Command: "echo 'hello'", Context: "global"},
						},
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: true},
				{value: "<disabled>", valid: true},
				{value: "q", valid: true},
				{value: "<c-c>", valid: true},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Custom command keybinding in prompt menu",
			setup: func(config *UserConfig, value string) {
				config.CustomCommands = []CustomCommand{
					{
						Key:         Keybinding{"X"},
						Description: "My Custom Commands",
						Prompts: []CustomCommandPrompt{
							{
								Options: []CustomCommandMenuOption{
									{Key: Keybinding{value}},
								},
							},
						},
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: true},
				{value: "<disabled>", valid: true},
				{value: "q", valid: true},
				{value: "<c-c>", valid: true},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Custom command output",
			setup: func(config *UserConfig, value string) {
				config.CustomCommands = []CustomCommand{
					{
						Output: value,
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: true},
				{value: "none", valid: true},
				{value: "terminal", valid: true},
				{value: "log", valid: true},
				{value: "logWithPty", valid: true},
				{value: "popup", valid: true},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Custom command sub menu",
			setup: func(config *UserConfig, _ string) {
				config.CustomCommands = []CustomCommand{
					{
						Key:         Keybinding{"X"},
						Description: "My Custom Commands",
						CommandMenu: []CustomCommand{
							{Key: Keybinding{"1"}, Command: "echo 'hello'", Context: "global"},
						},
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: true},
			},
		},
		{
			name: "Custom command sub menu",
			setup: func(config *UserConfig, _ string) {
				config.CustomCommands = []CustomCommand{
					{
						Key:     Keybinding{"X"},
						Context: "global", // context is not allowed for submenus
						CommandMenu: []CustomCommand{
							{Key: Keybinding{"1"}, Command: "echo 'hello'", Context: "global"},
						},
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: false},
			},
		},
		{
			name: "Custom command sub menu",
			setup: func(config *UserConfig, _ string) {
				config.CustomCommands = []CustomCommand{
					{
						Key:         Keybinding{"X"},
						LoadingText: "loading", // other properties are not allowed for submenus (using loadingText as an example)
						CommandMenu: []CustomCommand{
							{Key: Keybinding{"1"}, Command: "echo 'hello'", Context: "global"},
						},
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: false},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			for _, testCase := range s.testCases {
				config := GetDefaultConfig()
				s.setup(config, testCase.value)
				err := config.Validate()

				if testCase.valid {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			}
		})
	}
}

func TestUserConfigValidate_spinnerFrames(t *testing.T) {
	scenarios := []struct {
		name   string
		frames []string
		valid  bool
	}{
		{name: "empty", frames: []string{}, valid: false},
		{name: "single frame", frames: []string{"|"}, valid: true},
		{name: "all same width", frames: []string{"|", "/", "-", "\\"}, valid: true},
		{name: "all same width, multi-char", frames: []string{".  ", ".. ", "..."}, valid: true},
		{name: "all same width, wide runes", frames: []string{"⠋", "⠙", "⠹"}, valid: true},
		{name: "differing widths", frames: []string{"|", "//"}, valid: false},
		{name: "first differs from rest", frames: []string{"||", "/", "-"}, valid: false},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			config := GetDefaultConfig()
			config.Gui.Spinner.Frames = s.frames
			err := config.Validate()

			if s.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestUserConfigValidate_pagers(t *testing.T) {
	scenarios := []struct {
		name  string
		pager PagingConfig
		valid bool
	}{
		{name: "empty", pager: PagingConfig{}, valid: true},
		{name: "pager only", pager: PagingConfig{Pager: "delta"}, valid: true},
		{name: "external diff command only", pager: PagingConfig{ExternalDiffCommand: "difft"}, valid: true},
		{name: "git config external diff only", pager: PagingConfig{UseExternalDiffGitConfig: true}, valid: true},
		{name: "pager and external diff command", pager: PagingConfig{Pager: "delta", ExternalDiffCommand: "difft"}, valid: false},
		{name: "pager and git config external diff", pager: PagingConfig{Pager: "delta", UseExternalDiffGitConfig: true}, valid: false},
		{name: "both external diff mechanisms", pager: PagingConfig{ExternalDiffCommand: "difft", UseExternalDiffGitConfig: true}, valid: false},
		{name: "all three", pager: PagingConfig{Pager: "delta", ExternalDiffCommand: "difft", UseExternalDiffGitConfig: true}, valid: false},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			config := GetDefaultConfig()
			config.Git.Pagers = []PagingConfig{s.pager}
			err := config.Validate()

			if s.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
