package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RetainedWindowFocus = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that the focused context in each window is retained when switching worktrees",
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

		t.Views().Remotes().
			IsFocused()
	},
})
