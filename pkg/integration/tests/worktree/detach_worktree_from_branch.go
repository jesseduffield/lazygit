package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DetachWorktreeFromBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Detach a worktree from the branches view",
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
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("mybranch").IsSelected(),
				Contains("newbranch (worktree)"),
			).
			NavigateToLine(Contains("newbranch")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Branch newbranch is checked out by worktree linked-worktree")).
					Select(Equals("Detach worktree")).
					Confirm()
			}).
			Lines(
				Contains("mybranch"),
				Contains("newbranch").DoesNotContain("(worktree)").IsSelected(),
			)

		t.Views().Worktrees().
			Focus().
			Lines(
				Contains("repo (main)").IsSelected(),
				Contains("linked-worktree"),
			)
	},
})
