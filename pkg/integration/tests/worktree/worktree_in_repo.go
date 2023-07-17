package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var WorktreeInRepo = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Add a worktree inside the repo, then remove the directory and confirm the worktree is removed",
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
					Type("linked-worktree").
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
			// switch back to main worktree
			NavigateToLine(Contains("repo (main)")).
			Press(keys.Universal.Select).
			Lines(
				Contains("repo (main)").IsSelected(),
				Contains("linked-worktree"),
			)

		t.Views().Files().
			Focus().
			Lines(
				Contains("linked-worktree"),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("linked-worktree")).
					Select(Contains("Discard all changes")).
					Confirm()
			}).
			IsEmpty()

		// confirm worktree appears as missing
		t.Views().Worktrees().
			Focus().
			Lines(
				Contains("repo (main)").IsSelected(),
				Contains("linked-worktree (missing)"),
			)
	},
})
