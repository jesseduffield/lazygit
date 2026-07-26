package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var HideSidePanel = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Hide a side panel by omitting it from gui.sidePanels",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		// No stash panel.
		cfg.GetUserConfig().Gui.SidePanels = []config.SidePanel{
			{"status"},
			{"files", "worktrees", "submodules"},
			{"branches", "remotes", "tags"},
			{"commits", "reflog"},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(2)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Commits is now the last panel; cycling forward from it wraps around to
		// the status panel, skipping the hidden stash panel entirely.
		t.Views().Files().IsFocused().
			Press(keys.Universal.JumpToBlock[3])
		t.Views().Commits().IsFocused().
			Press(keys.Universal.NextBlock)
		t.Views().Status().IsFocused()
	},
})
