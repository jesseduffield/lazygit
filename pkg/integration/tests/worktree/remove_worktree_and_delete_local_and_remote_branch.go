package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RemoveWorktreeAndDeleteLocalAndRemoteBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Delete the local branch, the remote branch, and the worktree of a single branch checked out in another worktree, all at once",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CloneIntoRemote("origin")
		shell.EmptyCommit("initial commit")
		shell.NewBranch("mybranch")
		shell.EmptyCommit("commit on mybranch")
		shell.PushBranchAndSetUpstream("origin", "mybranch")
		shell.EmptyCommit("commit not pushed to the remote") // so mybranch isn't fully merged
		shell.Checkout("master")
		shell.AddWorktreeCheckout("mybranch", "../linked-worktree")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("mybranch (worktree linked-worktree)"),
			).
			NavigateToLine(Contains("mybranch")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Delete branch 'mybranch'?")).
					Select(Contains("Delete local and remote branch")).
					Confirm()
			}).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Branch mybranch is checked out by worktree linked-worktree")).
					Select(Contains("Remove worktree and delete local and remote branch")).
					Confirm()

				// mybranch is not contained in master, so we get the force-delete warning
				t.ExpectPopup().Confirmation().
					Title(Equals("Force delete branch")).
					Content(Equals("'mybranch' is not fully merged. Are you sure you want to delete it?")).
					Confirm()
			}).
			// The local branch is gone
			Lines(
				Contains("master").IsSelected(),
			)

		// The remote branch is gone too
		t.Views().Remotes().
			Focus().
			Lines(Contains("origin")).
			PressEnter()

		t.Views().RemoteBranches().
			IsEmpty()

		// And so is the worktree
		t.Views().Worktrees().
			Focus().
			Lines(
				Contains("(main worktree)").IsSelected(),
			)
	},
})
