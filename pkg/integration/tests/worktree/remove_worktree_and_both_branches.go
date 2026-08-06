package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RemoveWorktreeAndBothBranches = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "From the worktrees panel, remove a worktree and delete both its local and remote branch in one go",
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
		t.Views().Worktrees().
			Focus().
			Lines(
				Contains("(main worktree)").IsSelected(),
				Contains("linked-worktree"),
			).
			NavigateToLine(Contains("linked-worktree")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Remove worktree 'linked-worktree'?")).
					Select(Contains("Remove worktree and delete local and remote branch")).
					Confirm()

				// mybranch isn't fully merged, so we get the force-delete warning
				t.ExpectPopup().Confirmation().
					Title(Equals("Force delete branch")).
					Content(Equals("'mybranch' is not fully merged. Are you sure you want to delete it?")).
					Confirm()
			}).
			Lines(
				Contains("(main worktree)").IsSelected(),
			)

		// The remote branch is gone too
		t.Views().Remotes().
			Focus().
			Lines(Contains("origin")).
			PressEnter()

		t.Views().RemoteBranches().
			IsEmpty()

		// And so is the local branch
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
			)
	},
})
