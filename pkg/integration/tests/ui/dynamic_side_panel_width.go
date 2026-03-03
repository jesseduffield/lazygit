package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DynamicSidePanelWidth = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify dynamic side panel width resizes panels when switching focus",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.DynamicSidePanelWidth = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(5)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsFocused()

		t.Views().Commits().
			Focus().
			IsFocused().
			Lines(
				Contains("commit 05").IsSelected(),
				Contains("commit 04"),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
			)

		t.Views().Branches().
			Focus().
			IsFocused()

		t.Views().Commits().
			Focus().
			IsFocused().
			PressEnter()

		t.Views().CommitFiles().
			IsFocused()

		t.Views().Commits().
			Focus().
			IsFocused()
	},
})
