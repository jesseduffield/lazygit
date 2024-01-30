package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/updates"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func NewDummyUpdater() *updates.Updater {
	newAppConfig := config.NewDummyAppConfig()
	dummyUpdater, _ := updates.NewUpdater(utils.NewDummyCommon(), newAppConfig, oscommands.NewDummyOSCommand())
	return dummyUpdater
}

// NewDummyGui creates a new dummy GUI for testing
func NewDummyGui() *Gui {
	newAppConfig := config.NewDummyAppConfig()
	dummyGui, _ := NewGui(utils.NewDummyCommon(), newAppConfig, &git_commands.GitVersion{Major: 2, Minor: 0, Patch: 0}, NewDummyUpdater(), false, "", nil)
	return dummyGui
}
