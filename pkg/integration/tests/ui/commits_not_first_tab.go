package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitsNotFirstTab = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "With gui.sidePanels grouping commits behind another tab, no ghost view must appear over the side panels",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.SidePanels = []config.SidePanel{
			{"branches", "worktrees", "remotes"},
			{"files"},
			{"tags", "commits"},
			{"stash"},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// The commit files view is only shown after drilling into a commit; at
		// startup it must be hidden, or it'd cover the side panels.
		t.Views().CommitFiles().
			IsInvisible()
	},
})
