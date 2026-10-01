package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var InitialSidePanel = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The side panel named by gui.initialSidePanel is focused at startup",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.InitialSidePanel = "commits"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("one"),
			)
	},
})
