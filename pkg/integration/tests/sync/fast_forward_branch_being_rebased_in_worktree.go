package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FastForwardBranchBeingRebasedInWorktree = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Try to fast-forward a branch that is being rebased in another worktree",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "content")
		shell.Commit("one")
		shell.EmptyCommit("two")
		shell.EmptyCommit("three")

		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("master", "origin/master")

		// remove a commit so that there is something to fast-forward to
		shell.HardReset("HEAD~1")

		shell.NewBranch("other")
		shell.AddWorktreeCheckout("master", "../linked-worktree")

		// the failing exec stops the rebase after picking "two", with HEAD
		// detached from master
		shell.RunCommandExpectError([]string{"git", "-C", "../linked-worktree", "rebase", "--exec", "false", "HEAD~1"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("other").IsSelected(),
				Contains("master (worktree linked-worktree) ↓1"),
			).
			NavigateToLine(Contains("master")).
			Press(keys.Branches.FastForward)

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Equals("Cannot fast-forward 'master' because it is being rebased or bisected in worktree linked-worktree")).
			Confirm()

		t.Views().Worktrees().
			Focus().
			NavigateToLine(Contains("linked-worktree")).
			Press(keys.Universal.Select)

		t.Views().Information().Content(Contains("Rebasing"))

		t.Views().Commits().
			Lines(
				Contains("─── Pending rebase todos"),
				Contains("─── Commits"),
				Contains("two"),
				Contains("one"),
			)
	},
})
