package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AddForExistingBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a worktree that checks out an existing branch (no new branch), and confirm the option is disabled for an already-checked-out branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.NewBranchFrom("otherbranch", "mybranch")
		shell.Checkout("mybranch")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("mybranch").IsSelected(),
				Contains("otherbranch"),
			).
			// the current branch is checked out by this worktree, so "Worktree
			// for 'mybranch'" is disabled
			Press(keys.Universal.NewWorktree).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("New worktree")).
					Select(Contains("New worktree for 'mybranch'")).
					Tooltip(Contains("Branch mybranch is checked out by worktree repo")).
					Confirm().
					Tap(func() {
						t.ExpectToast(Contains("Branch mybranch is checked out by worktree repo"))
					}).
					Cancel()
			}).
			// otherbranch is not checked out anywhere, so we can make a worktree for it
			NavigateToLine(Contains("otherbranch")).
			Press(keys.Universal.NewWorktree).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("New worktree")).
					Select(Contains("New worktree for 'otherbranch'")).
					Confirm()

				t.ExpectPopup().Menu().
					Title(Equals("Worktree location")).
					Confirm()
			})

		// we've switched into the new worktree, which has otherbranch checked out
		t.Views().Branches().
			IsFocused().
			Lines(
				Contains("otherbranch").IsSelected(),
				Contains("mybranch (worktree repo)"),
			)

		t.Views().Status().
			Content(Contains("repo(otherbranch) → otherbranch"))
	},
})
