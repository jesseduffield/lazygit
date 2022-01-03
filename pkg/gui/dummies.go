package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/updates"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// NewDummyGui creates a new dummy GUI for testing
func NewDummyUpdater() *updates.Updater {
	newAppConfig := config.NewDummyAppConfig()
	dummyUpdater, _ := updates.NewUpdater(utils.NewDummyCommon(), newAppConfig, oscommands.NewDummyOSCommand())
	return dummyUpdater
}

func NewDummyGui() *Gui {
	newAppConfig := config.NewDummyAppConfig()
	dummyGui, _ := NewGui(utils.NewDummyCommon(), commands.NewDummyGitCommand(), oscommands.NewDummyOSCommand(), newAppConfig, NewDummyUpdater(), "", false)
	return dummyGui
}
