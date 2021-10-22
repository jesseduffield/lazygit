package commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// NewDummyGitCommand creates a new dummy GitCommand for testing
func NewDummyGitCommand() *GitCommand {
	return NewDummyGitCommandWithOSCommand(oscommands.NewDummyOSCommand())
}

// NewDummyGitCommandWithOSCommand creates a new dummy GitCommand for testing
func NewDummyGitCommandWithOSCommand(osCommand *oscommands.OSCommand) *GitCommand {
	newAppConfig := config.NewDummyAppConfig()
	return &GitCommand{
		Log:       utils.NewDummyLog(),
		OSCommand: osCommand,
		Tr:        i18n.NewTranslationSet(utils.NewDummyLog(), newAppConfig.GetUserConfig().Gui.Language),
		Config:    newAppConfig,
		GitConfig: git_config.NewFakeGitConfig(map[string]string{}),
	}
}
