package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NewWorktreePicker = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "From the worktrees panel, the picker suggests only branches not already checked out, guards verbatim type-ins, and checks out an existing branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.NewBranchFrom("feature", "mybranch")
		shell.Checkout("mybranch")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Worktrees().
			Focus().
			Lines(
				Contains("(main worktree)"),
			).
			Press(keys.Universal.New).
			Tap(func() {
				// mybranch is checked out by the current worktree, so it's not
				// suggested; feature is
				t.ExpectPopup().Prompt().
					Title(Equals("New worktree for branch")).
					SuggestionLines(Contains("feature")).
					// typing a checked-out branch verbatim is still rejected
					Type("mybranch").
					Confirm()

				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Contains("Branch mybranch is checked out by worktree repo")).
					Confirm()
			}).
			Press(keys.Universal.New).
			Tap(func() {
				// picking an existing branch checks it out (no new branch)
				t.ExpectPopup().Prompt().
					Title(Equals("New worktree for branch")).
					Type("feature").
					Confirm()

				t.ExpectPopup().Menu().
					Title(Equals("Worktree location")).
					Confirm()
			}).
			// we stay in the worktrees panel, now switched into the new worktree
			IsFocused().
			Lines(
				Contains("feature").IsSelected(),
				Contains("(main worktree)"),
			)

		t.Views().Status().
			Content(Contains("repo(feature) → feature"))
	},
})
