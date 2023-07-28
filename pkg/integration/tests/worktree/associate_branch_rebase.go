package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// This is important because `git worktree list` will show a worktree being in a detached head state (which is true)
// when it's in the middle of a rebase, but it won't tell you about the branch it's on.
// Even so, if you attempt to check out that branch from another worktree git won't let you, so we need to
// keep track of the association ourselves.

// We need different logic for associated the branch depending on whether it's a main worktree or
// linked worktree, so this test handles both.

var AssociateBranchRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that when you start a rebase in a linked or main worktree, Lazygit still associates the worktree with the branch",
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
			)

		// start a rebase on the main worktree
		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("commit 2")).
			Press(keys.Universal.Edit)

		t.Views().Information().Content(Contains("Rebasing"))

		t.Views().Branches().
			Focus().
			// switch to linked worktree
			NavigateToLine(Contains("newbranch")).
			Press(keys.Universal.Select).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Switch to worktree")).
					Content(Equals("This branch is checked out by worktree linked-worktree. Do you want to switch to that worktree?")).
					Confirm()

				t.Views().Information().Content(DoesNotContain("Rebasing"))
			}).
			Lines(
				Contains("newbranch").IsSelected(),
				Contains("mybranch (worktree)"),
			)

		// start a rebase on the linked worktree
		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("commit 2")).
			Press(keys.Universal.Edit)

		t.Views().Information().Content(Contains("Rebasing"))

		// switch back to main worktree
		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("mybranch")).
			Press(keys.Universal.Select).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Switch to worktree")).
					Content(Equals("This branch is checked out by worktree repo. Do you want to switch to that worktree?")).
					Confirm()
			}).
			Lines(
				Contains("(no branch").IsSelected(),
				Contains("mybranch"),
				// even though the linked worktree is rebasing, we still associate it with the branch
				Contains("newbranch (worktree)"),
			)
	},
})
