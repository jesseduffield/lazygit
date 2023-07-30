package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Crud = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "From the worktrees view, add a work tree, switch to it, switch back, and remove it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("mybranch"),
			)

		t.Views().Status().
			Lines(
				Contains("repo → mybranch"),
			)

		t.Views().Worktrees().
			Focus().
			Lines(
				Contains("repo (main)"),
			).
			Press(keys.Universal.New).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Worktree")).
					Select(Contains(`Create worktree from ref`).DoesNotContain(("detached"))).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New worktree base ref")).
					InitialText(Equals("mybranch")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New worktree path")).
					Type("../linked-worktree").
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New branch name (leave blank to checkout mybranch)")).
					Type("newbranch").
					Confirm()
			}).
			Lines(
				Contains("linked-worktree").IsSelected(),
				Contains("repo (main)"),
			).
			// confirm we're still in the same view
			IsFocused()

		// status panel includes the worktree if it's a linked worktree
		t.Views().Status().
			Lines(
				Contains("repo(linked-worktree) → newbranch"),
			)

		t.Views().Branches().
			Lines(
				Contains("newbranch"),
				Contains("mybranch"),
			)

		t.Views().Worktrees().
			// confirm we can't remove the current worktree
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Equals("You cannot remove the current worktree!")).
					Confirm()
			}).
			// confirm we cannot remove the main worktree
			NavigateToLine(Contains("repo (main)")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Equals("You cannot remove the main worktree!")).
					Confirm()
			}).
			// switch back to main worktree
			Press(keys.Universal.Select).
			Lines(
				Contains("repo (main)").IsSelected(),
				Contains("linked-worktree"),
			)

		t.Views().Branches().
			Lines(
				Contains("mybranch"),
				Contains("newbranch"),
			)

		t.Views().Worktrees().
			// remove linked worktree
			NavigateToLine(Contains("linked-worktree")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Remove worktree")).
					Content(Contains("Are you sure you want to remove worktree 'linked-worktree'?")).
					Confirm()
			}).
			Lines(
				Contains("repo (main)").IsSelected(),
			)
	},
})
