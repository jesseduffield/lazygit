package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var InitialSidePanelNotFirstTab = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "gui.initialSidePanel names a tab that isn't the first one of its panel, so it must be brought to the front",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.SidePanels = []config.SidePanel{
			{"status"},
			{"files"},
			{"branches"},
			{"reflog", "commits"},
			{"stash"},
		}
		cfg.GetUserConfig().Gui.InitialSidePanel = "commits"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Commits is shown in front of the reflog tab it shares a panel with,
		// rather than being focused behind it.
		t.Views().Commits().
			IsActiveTab().
			IsFocused().
			Lines(
				Contains("one"),
			)
	},
})
