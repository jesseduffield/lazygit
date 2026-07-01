package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RemoveWorktreeAndBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "From the worktrees panel, remove a worktree and delete its branch in one go",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.NewBranch("newbranch")
		shell.EmptyCommit("commit on newbranch")
		shell.Checkout("mybranch")
		shell.AddWorktreeCheckout("newbranch", "../linked-worktree")
		shell.RunCommand([]string{"git", "worktree", "add", "--detach", "../detached-worktree"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Worktrees().
			Focus().
			Lines(
				Contains("(main worktree)").IsSelected(),
				Contains("detached-worktree"),
				Contains("linked-worktree"),
			).
			// A detached worktree has no branch, so neither delete action is offered
			NavigateToLine(Contains("detached-worktree")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Remove worktree 'detached-worktree'?")).
					Select(Contains("Remove worktree and delete branch")).
					Tooltip(Contains("This worktree is not checked out on a branch")).
					Select(Contains("Remove worktree and delete local and remote branch")).
					Tooltip(Contains("This worktree is not checked out on a branch")).
					Cancel()
			}).
			// Remove a worktree and delete its branch at once. newbranch has no
			// upstream, so deleting the remote branch too isn't offered.
			NavigateToLine(Contains("linked-worktree")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Remove worktree 'linked-worktree'?")).
					Select(Contains("Remove worktree and delete local and remote branch")).
					Tooltip(Contains("The selected branch has no upstream")).
					Select(Contains("Remove worktree and delete branch")).
					Confirm()

				// newbranch isn't fully merged, so we get the force-delete warning
				t.ExpectPopup().Confirmation().
					Title(Equals("Force delete branch")).
					Content(Equals("'newbranch' is not fully merged. Are you sure you want to delete it?")).
					Confirm()
			}).
			Lines(
				Contains("(main worktree)"),
				Contains("detached-worktree"),
			)

		// The branch is gone too
		t.Views().Branches().
			Focus().
			Lines(
				Contains("mybranch").IsSelected(),
			)
	},
})
