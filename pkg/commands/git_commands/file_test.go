package git_commands

import (
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestEditFileCmdStr(t *testing.T) {
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

		s.test(instance.GetEditCmdStr(s.filename, 1))
		s.runner.CheckForMissingCalls()
	}
}
