package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// This is verifying logic that is subject to change (we're just doing the easiest approach)
// There are two other UX flows we could have:
// 1) associate window tab states with the repo, so that when you switch back to a repo you get the same window tab states
// 2) retain the same window tab states when switching repos
// Option 1 is straightforward, but option 2 is harder because you'd need to deactivate any views containing dependent
// content e.g. the sub-commits view.

var ResetWindowTabs = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that window tabs are reset whenever switching repos",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.EmptyCommit("commit 2")
		shell.EmptyCommit("commit 3")
		shell.AddWorktree("mybranch", "../linked-worktree", "newbranch")
		shell.AddFileInWorktree("../linked-worktree")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// focus the remotes tab i.e. the second tab in the branches window
		t.Views().Remotes().
			Focus()

		t.Views().Worktrees().
			Focus().
			Lines(
				Contains("repo (main)").IsSelected(),
				Contains("linked-worktree"),
			).
			NavigateToLine(Contains("linked-worktree")).
			Press(keys.Universal.Select).
			Lines(
				Contains("linked-worktree").IsSelected(),
				Contains("repo (main)"),
			).
			// navigate back to the branches window
			Press(keys.Universal.NextBlock)

		t.Views().Branches().
			IsFocused()
	},
})
