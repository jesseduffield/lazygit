package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FastForwardBranchesInWorktreesWithOneFailing = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Fast-forward two branches checked out in other worktrees when the fast-forward fails in the first one",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "content")
		shell.Commit("one")
		shell.NewBranch("branch-a")
		shell.CreateFileAndAdd("a-file", "a")
		shell.Commit("a")
		shell.Checkout("master")
		shell.NewBranch("branch-b")
		shell.CreateFileAndAdd("b-file", "b")
		shell.Commit("b")
		shell.Checkout("master")

		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("branch-a", "origin/branch-a")
		shell.SetBranchUpstream("branch-b", "origin/branch-b")

		// remove a commit from each branch so that there is something to
		// fast-forward to
		shell.RunCommand([]string{"git", "branch", "-f", "branch-a", "branch-a~1"})
		shell.RunCommand([]string{"git", "branch", "-f", "branch-b", "branch-b~1"})

		shell.AddWorktreeCheckout("branch-a", "../worktree-a")
		shell.AddWorktreeCheckout("branch-b", "../worktree-b")

		// an untracked file that the fast-forward in worktree-a would overwrite
		shell.CreateFile("../worktree-a/a-file", "mine")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("branch-a (worktree worktree-a) ↓1"),
				Contains("branch-b (worktree worktree-b) ↓1"),
			).
			SelectNextItem().
			Press(keys.Universal.ToggleRangeSelect).
			SelectNextItem().
			Press(keys.Branches.FastForward)

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("would be overwritten by merge")).
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("master"),
				Contains("branch-a (worktree worktree-a) ↓1").IsSelected(),
				Contains("branch-b (worktree worktree-b) ✓").IsSelected(),
			)
	},
})
