package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestEditFilesCmd(t *testing.T) {
	type scenario struct {
		filenames      []string
		osConfig       config.OSConfig
		expectedCmdStr string
		suspend        bool
	}

	scenarios := []scenario{
		{
			filenames:      []string{"test"},
			osConfig:       config.OSConfig{},
			expectedCmdStr: `vim -- "test"`,
			suspend:        true,
		},
		{
			filenames: []string{"test"},
			osConfig: config.OSConfig{
				Edit: "nano {{filename}}",
			},
			expectedCmdStr: `nano "test"`,
			suspend:        true,
		},
		{
			filenames: []string{"file/with space"},
			osConfig: config.OSConfig{
				EditPreset: "sublime",
			},
			expectedCmdStr: `subl -- "file/with space"`,
			suspend:        false,
		},
		{
			filenames: []string{"multiple", "files"},
			osConfig: config.OSConfig{
				EditPreset: "sublime",
			},
			expectedCmdStr: `subl -- "multiple" "files"`,
			suspend:        false,
		},
	}

	for _, s := range scenarios {
		userConfig := config.GetDefaultConfig()
		userConfig.OS = s.osConfig

		instance := buildFileCommands(commonDeps{
			userConfig: userConfig,
		})

		cmdStr, suspend := instance.GetEditCmdStr(s.filenames)
		assert.Equal(t, s.expectedCmdStr, cmdStr)
		assert.Equal(t, s.suspend, suspend)
	}
}

func TestEditFileAtLineCmd(t *testing.T) {
	type scenario struct {
		filename       string
		lineNumber     int
		osConfig       config.OSConfig
		expectedCmdStr string
		suspend        bool
	}

	scenarios := []scenario{
		{
			filename:       "test",
			lineNumber:     42,
			osConfig:       config.OSConfig{},
			expectedCmdStr: `vim +42 -- "test"`,
			suspend:        true,
		},
		{
			filename:   "test",
			lineNumber: 35,
			osConfig: config.OSConfig{
				EditAtLine: "nano +{{line}} {{filename}}",
			},
			expectedCmdStr: `nano +35 "test"`,
			suspend:        true,
		},
		{
			filename:   "file/with space",
			lineNumber: 12,
			osConfig: config.OSConfig{
				EditPreset: "sublime",
			},
			expectedCmdStr: `subl -- "file/with space":12`,
			suspend:        false,
		},
	}

	for _, s := range scenarios {
		userConfig := config.GetDefaultConfig()
		userConfig.OS = s.osConfig

		instance := buildFileCommands(commonDeps{
			userConfig: userConfig,
		})

		cmdStr, suspend := instance.GetEditAtLineCmdStr(s.filename, s.lineNumber)
		assert.Equal(t, s.expectedCmdStr, cmdStr)
		assert.Equal(t, s.suspend, suspend)
	}
}

func TestEditFileAtLineAndWaitCmd(t *testing.T) {
	type scenario struct {
		filename       string
		lineNumber     int
		osConfig       config.OSConfig
		expectedCmdStr string
	}

	scenarios := []scenario{
		{
			filename:       "test",
			lineNumber:     42,
			osConfig:       config.OSConfig{},
			expectedCmdStr: `vim +42 -- "test"`,
		},
		{
			filename:   "file/with space",
			lineNumber: 12,
			osConfig: config.OSConfig{
				EditPreset: "sublime",
			},
			expectedCmdStr: `subl --wait -- "file/with space":12`,
		},
	}

	for _, s := range scenarios {
		userConfig := config.GetDefaultConfig()
		userConfig.OS = s.osConfig

		instance := buildFileCommands(commonDeps{
			userConfig: userConfig,
		})

		cmdStr := instance.GetEditAtLineAndWaitCmdStr(s.filename, s.lineNumber)
		assert.Equal(t, s.expectedCmdStr, cmdStr)
	}
}

func TestGuessDefaultEditor(t *testing.T) {
	type scenario struct {
		gitConfigMockResponses map[string]string
		getenv                 func(string) string
		expectedResult         string
	}

	scenarios := []scenario{
		{
			gitConfigMockResponses: nil,
			getenv: func(env string) string {
				return ""
			},
			expectedResult: "",
		},
		{
			gitConfigMockResponses: map[string]string{"core.editor": "nano"},
			getenv: func(env string) string {
				return ""
			},
			expectedResult: "nano",
		},
		{
			gitConfigMockResponses: map[string]string{"core.editor": "code -w"},
			getenv: func(env string) string {
				return ""
			},
			expectedResult: "code",
		},
		{
			gitConfigMockResponses: nil,
			getenv: func(env string) string {
				if env == "VISUAL" {
					return "emacs"
				}

				return ""
			},
			expectedResult: "emacs",
		},
		{
			gitConfigMockResponses: nil,
			getenv: func(env string) string {
				if env == "EDITOR" {
					return "bbedit -w"
				}

				return ""
			},
			expectedResult: "bbedit",
		},
	}

	for _, s := range scenarios {
		instance := buildFileCommands(commonDeps{
			gitConfig: git_config.NewFakeGitConfig(s.gitConfigMockResponses),
			getenv:    s.getenv,
		})

		assert.Equal(t, s.expectedResult, instance.guessDefaultEditor())
	}
}
