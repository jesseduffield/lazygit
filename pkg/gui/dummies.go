package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/updates"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// NewDummyGui creates a new dummy GUI for testing
func NewDummyUpdater() *updates.Updater {
	newAppConfig := config.NewDummyAppConfig()
	DummyUpdater, _ := updates.NewUpdater(utils.NewDummyLog(), newAppConfig, oscommands.NewDummyOSCommand(), i18n.NewTranslationSet(utils.NewDummyLog(), newAppConfig.GetUserConfig().Gui.Language))
	return DummyUpdater
}

func NewDummyGui() *Gui {
	newAppConfig := config.NewDummyAppConfig()
	DummyGui, _ := NewGui(utils.NewDummyLog(), commands.NewDummyGitCommand(), oscommands.NewDummyOSCommand(), i18n.NewTranslationSet(utils.NewDummyLog(), newAppConfig.GetUserConfig().Gui.Language), newAppConfig, NewDummyUpdater(), "", false)
	return DummyGui
}
