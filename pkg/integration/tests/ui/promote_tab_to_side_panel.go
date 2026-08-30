package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PromoteTabToSidePanel = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Promote the worktrees tab to its own top-level side panel via gui.sidePanels",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		// Worktrees is pulled out of the files panel into its own panel.
		cfg.GetUserConfig().Gui.SidePanels = []config.SidePanel{
			{"status"},
			{"files", "submodules"},
			{"worktrees"},
			{"branches", "remotes", "tags"},
			{"commits", "reflog"},
			{"stash"},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(2)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Worktrees is now its own panel in the third position, reachable by its
		// jump key rather than as a tab of the files panel.
		t.Views().Files().IsFocused().
			Press(keys.Universal.JumpToBlock[2])
		t.Views().Worktrees().IsFocused().
			Press(keys.Universal.JumpToBlock[1])

		// The files panel's tabs are now just files and submodules, so cycling
		// tabs from files goes straight to submodules.
		t.Views().Files().IsFocused().
			Press(keys.Universal.NextTab)
		t.Views().Submodules().IsFocused()
	},
})
