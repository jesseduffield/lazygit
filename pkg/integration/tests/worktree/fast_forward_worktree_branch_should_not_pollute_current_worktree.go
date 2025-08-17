package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FastForwardWorktreeBranchShouldNotPolluteCurrentWorktree = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Fast-forward a linked worktree branch from another worktree",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// both main and linked worktree will have changed to fast-forward
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.EmptyCommit("two")
		shell.EmptyCommit("three")
		shell.NewBranch("newbranch")

		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("mybranch", "origin/mybranch")
		shell.SetBranchUpstream("newbranch", "origin/newbranch")

		// remove the 'three' commit so that we have something to pull from the remote
		shell.HardReset("HEAD^")
		shell.Checkout("mybranch")
		shell.HardReset("HEAD^")

		shell.AddWorktreeCheckout("newbranch", "../linked-worktree")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("mybranch").Contains("↓1").IsSelected(),
				Contains("newbranch (worktree)").Contains("↓1"),
			).
			Press(keys.Branches.FastForward).
			Lines(
				Contains("mybranch").Contains("✓").IsSelected(),
				Contains("newbranch (worktree)").Contains("↓1"),
			).
			NavigateToLine(Contains("newbranch (worktree)")).
			Press(keys.Branches.FastForward).
			Lines(
				Contains("mybranch").Contains("✓"),
				Contains("newbranch (worktree)").Contains("✓").IsSelected(),
			).
			NavigateToLine(Contains("mybranch"))

		// check the current worktree that it has no lines in the File changes pane
		t.Views().Files().
			Focus().
			Press(keys.Files.RefreshFiles).
			LineCount(EqualsInt(0))
	},
})
