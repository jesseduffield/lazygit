package gui

import (
	"github.com/lobes/lazytask/pkg/commands/git_commands"
	"github.com/lobes/lazytask/pkg/commands/oscommands"
	"github.com/lobes/lazytask/pkg/config"
	"github.com/lobes/lazytask/pkg/updates"
	"github.com/lobes/lazytask/pkg/utils"
)

func NewDummyUpdater() *updates.Updater {
	newAppConfig := config.NewDummyAppConfig()
	dummyUpdater, _ := updates.NewUpdater(utils.NewDummyCommon(), newAppConfig, oscommands.NewDummyOSCommand())
	return dummyUpdater
}

// NewDummyGui creates a new dummy GUI for testing
func NewDummyGui() *Gui {
	newAppConfig := config.NewDummyAppConfig()
	dummyGui, _ := NewGui(utils.NewDummyCommon(), newAppConfig, &git_commands.GitVersion{}, NewDummyUpdater(), false, "", nil)
	return dummyGui
}
