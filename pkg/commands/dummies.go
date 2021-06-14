package commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// NewDummyGit creates a new dummy Git for testing
func NewDummyGit() *Git {
	return NewDummyGitWithOS(oscommands.NewDummyOS())
}

func NewDummyGitConfig() *GitConfig {
	return &GitConfig{
		getGitConfigValue: func(string) (string, error) { return "", nil },
	}
}

// NewDummyGitWithOS creates a new dummy Git for testing
func NewDummyGitWithOS(oS *oscommands.OS) *Git {
	return &Git{
		GitConfig: &GitConfig{},
		log:       utils.NewDummyLog(),
		os:        oS,
		tr:        i18n.NewTranslationSet(utils.NewDummyLog()),
		config:    config.NewDummyAppConfig(),
	}
}
