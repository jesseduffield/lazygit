package gui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/updates"
)

// NewDummyGui creates a new dummy GUI for testing
func NewDummyUpdater() *updates.Updater {
    DummyUpdater, _ := updates.NewUpdater(utils.NewDummyLog(), config.NewDummyAppConfig(), oscommands.NewDummyOSCommand(), i18n.NewTranslationSet(utils.NewDummyLog()))
    return DummyUpdater
}

func NewDummyGui() *Gui {
        DummyGui, _ := NewGui(utils.NewDummyLog(), commands.NewDummyGitCommand(), oscommands.NewDummyOSCommand(), i18n.NewTranslationSet(utils.NewDummyLog()), config.NewDummyAppConfig(), NewDummyUpdater(), "", false)
	return DummyGui
}
