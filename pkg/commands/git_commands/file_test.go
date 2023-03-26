package git_commands

import (
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestEditFileCmdStrLegacy(t *testing.T) {
	type scenario struct {
		filename                  string
		configEditCommand         string
		configEditCommandTemplate string
		runner                    *oscommands.FakeCmdObjRunner
		getenv                    func(string) string
		gitConfigMockResponses    map[string]string
		test                      func(string, error)
	}

	scenarios := []scenario{
		{
			filename:                  "test",
			configEditCommand:         "",
			configEditCommandTemplate: "{{editor}} {{filename}}",
			runner: oscommands.NewFakeRunner(t).
				Expect(`which vi`, "", errors.New("error")),
			getenv: func(env string) string {
				return ""
			},
			gitConfigMockResponses: nil,
			test: func(cmdStr string, err error) {
				assert.EqualError(t, err, "No editor defined in config file, $GIT_EDITOR, $VISUAL, $EDITOR, or git config")
			},
		},
		{
			filename:                  "test",
			configEditCommand:         "nano",
			configEditCommandTemplate: "{{editor}} {{filename}}",
			runner:                    oscommands.NewFakeRunner(t),
			getenv: func(env string) string {
				return ""
			},
			gitConfigMockResponses: nil,
			test: func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, `nano "test"`, cmdStr)
			},
		},
		{
			filename:                  "test",
			configEditCommand:         "",
			configEditCommandTemplate: "{{editor}} {{filename}}",
			runner:                    oscommands.NewFakeRunner(t),
			getenv: func(env string) string {
				return ""
			},
			gitConfigMockResponses: map[string]string{"core.editor": "nano"},
			test: func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, `nano "test"`, cmdStr)
			},
		},
		{
			filename:                  "test",
			configEditCommand:         "",
			configEditCommandTemplate: "{{editor}} {{filename}}",
			runner:                    oscommands.NewFakeRunner(t),
			getenv: func(env string) string {
				if env == "VISUAL" {
					return "nano"
				}

				return ""
			},
			gitConfigMockResponses: nil,
			test: func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, `nano "test"`, cmdStr)
			},
		},
		{
			filename:                  "test",
			configEditCommand:         "",
			configEditCommandTemplate: "{{editor}} {{filename}}",
			runner:                    oscommands.NewFakeRunner(t),
			getenv: func(env string) string {
				if env == "EDITOR" {
					return "emacs"
				}

				return ""
			},
			gitConfigMockResponses: nil,
			test: func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, `emacs "test"`, cmdStr)
			},
		},
		{
			filename:                  "test",
			configEditCommand:         "",
			configEditCommandTemplate: "{{editor}} {{filename}}",
			runner: oscommands.NewFakeRunner(t).
				Expect(`which vi`, "/usr/bin/vi", nil),
			getenv: func(env string) string {
				return ""
			},
			gitConfigMockResponses: nil,
			test: func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, `vi "test"`, cmdStr)
			},
		},
		{
			filename:                  "file/with space",
			configEditCommand:         "",
			configEditCommandTemplate: "{{editor}} {{filename}}",
			runner: oscommands.NewFakeRunner(t).
				Expect(`which vi`, "/usr/bin/vi", nil),
			getenv: func(env string) string {
				return ""
			},
			gitConfigMockResponses: nil,
			test: func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, `vi "file/with space"`, cmdStr)
			},
		},
		{
			filename:                  "open file/at line",
			configEditCommand:         "vim",
			configEditCommandTemplate: "{{editor}} +{{line}} {{filename}}",
			runner:                    oscommands.NewFakeRunner(t),
			getenv: func(env string) string {
				return ""
			},
			gitConfigMockResponses: nil,
			test: func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, `vim +1 "open file/at line"`, cmdStr)
			},
		},
		{
			filename:                  "default edit command template",
			configEditCommand:         "vim",
			configEditCommandTemplate: "",
			runner:                    oscommands.NewFakeRunner(t),
			getenv: func(env string) string {
				return ""
			},
			gitConfigMockResponses: nil,
			test: func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, `vim +1 -- "default edit command template"`, cmdStr)
			},
		},
	}

	for _, s := range scenarios {
		userConfig := config.GetDefaultConfig()
		userConfig.OS.EditCommand = s.configEditCommand
		userConfig.OS.EditCommandTemplate = s.configEditCommandTemplate

		instance := buildFileCommands(commonDeps{
			runner:     s.runner,
			userConfig: userConfig,
			gitConfig:  git_config.NewFakeGitConfig(s.gitConfigMockResponses),
			getenv:     s.getenv,
		})

		s.test(instance.GetEditCmdStrLegacy(s.filename, 1))
		s.runner.CheckForMissingCalls()
	}
}

func TestEditFileCmd(t *testing.T) {
	type scenario struct {
		filename               string
		osConfig               config.OSConfig
		expectedCmdStr         string
		expectedEditInTerminal bool
	}

	scenarios := []scenario{
		{
			filename:               "test",
			osConfig:               config.OSConfig{},
			expectedCmdStr:         `vim -- "test"`,
			expectedEditInTerminal: true,
		},
		{
			filename: "test",
			osConfig: config.OSConfig{
				Edit: "nano {{filename}}",
			},
			expectedCmdStr:         `nano "test"`,
			expectedEditInTerminal: true,
		},
		{
			filename: "file/with space",
			osConfig: config.OSConfig{
				EditPreset: "sublime",
			},
			expectedCmdStr:         `subl -- "file/with space"`,
			expectedEditInTerminal: false,
		},
	}

	for _, s := range scenarios {
		userConfig := config.GetDefaultConfig()
		userConfig.OS = s.osConfig

		instance := buildFileCommands(commonDeps{
			userConfig: userConfig,
		})

		cmdStr, editInTerminal := instance.GetEditCmdStr(s.filename)
		assert.Equal(t, s.expectedCmdStr, cmdStr)
		assert.Equal(t, s.expectedEditInTerminal, editInTerminal)
	}
}

func TestEditFileAtLineCmd(t *testing.T) {
	type scenario struct {
		filename               string
		lineNumber             int
		osConfig               config.OSConfig
		expectedCmdStr         string
		expectedEditInTerminal bool
	}

	scenarios := []scenario{
		{
			filename:               "test",
			lineNumber:             42,
			osConfig:               config.OSConfig{},
			expectedCmdStr:         `vim +42 -- "test"`,
			expectedEditInTerminal: true,
		},
		{
			filename:   "test",
			lineNumber: 35,
			osConfig: config.OSConfig{
				EditAtLine: "nano +{{line}} {{filename}}",
			},
			expectedCmdStr:         `nano +35 "test"`,
			expectedEditInTerminal: true,
		},
		{
			filename:   "file/with space",
			lineNumber: 12,
			osConfig: config.OSConfig{
				EditPreset: "sublime",
			},
			expectedCmdStr:         `subl -- "file/with space":12`,
			expectedEditInTerminal: false,
		},
	}

	for _, s := range scenarios {
		userConfig := config.GetDefaultConfig()
		userConfig.OS = s.osConfig

		instance := buildFileCommands(commonDeps{
			userConfig: userConfig,
		})

		cmdStr, editInTerminal := instance.GetEditAtLineCmdStr(s.filename, s.lineNumber)
		assert.Equal(t, s.expectedCmdStr, cmdStr)
		assert.Equal(t, s.expectedEditInTerminal, editInTerminal)
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
