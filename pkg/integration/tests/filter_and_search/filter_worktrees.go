package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterWorktrees = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filtering worktrees by branch name",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.NewBranch("branch-aaa")
		shell.NewBranch("branch-xxx")
		shell.Checkout("master")
		shell.AddWorktreeCheckout("branch-aaa", "../worktree-xxx")
		shell.AddWorktreeCheckout("branch-xxx", "../worktree-1")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Worktrees().
			Focus().
			Lines(
				Contains("(main worktree)").IsSelected(),
				Contains("worktree-1   branch-xxx"),
				Contains("worktree-xxx branch-aaa"),
			).
			FilterOrSearch("xxx").
			Lines(
				Contains("worktree-1   branch-xxx"),
				Contains("worktree-xxx branch-aaa"),
			)
	},
})
