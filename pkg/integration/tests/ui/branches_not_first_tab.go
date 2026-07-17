package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BranchesNotFirstTab = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "With gui.sidePanels grouping branches behind another tab, no ghost view must appear over the side panels",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.SidePanels = []config.SidePanel{
			{"worktrees", "branches", "remotes"},
			{"files"},
			{"commits", "tags"},
			{"stash"},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// The remote branches and sub-commits views are only shown after
		// drilling into a remote or a branch; at startup both must be hidden,
		// or they'd cover the side panels.
		t.Views().RemoteBranches().
			/* EXPECTED:
			IsInvisible()
			ACTUAL: */
			IsVisible()
		t.Views().SubCommits().IsInvisible()
	},
})
