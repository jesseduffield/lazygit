package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ReorderSidePanels = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reorder the side panels with gui.sidePanels, swapping the branches and commits panels",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.SidePanels = []config.SidePanel{
			{"status"},
			{"files", "worktrees", "submodules"},
			{"commits", "reflog"},
			{"branches", "remotes", "tags"},
			{"stash"},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(2)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// The third panel is now commits and the fourth is branches (the reverse
		// of the default order), so their jump keys are swapped.
		t.Views().Files().IsFocused().
			Press(keys.Universal.JumpToBlock[2])
		t.Views().Commits().IsFocused().
			Press(keys.Universal.JumpToBlock[3])
		t.Views().Branches().IsFocused()
	},
})
