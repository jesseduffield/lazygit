package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterPreservesSelectionOnModelChange = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that when a filter is active and the model changes, the selection is preserved",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.NewBranch("branch-alpha")
		shell.NewBranch("branch-beta")
		shell.NewBranch("branch-gamma")
		shell.NewBranch("checked-out-branch")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("checked-out-branch").IsSelected(),
				Contains("branch-alpha"),
				Contains("branch-beta"),
				Contains("branch-gamma"),
				Contains("master"),
			).
			FilterOrSearch("branch-").
			Lines(
				Contains("branch-alpha").IsSelected(),
				Contains("branch-beta"),
				Contains("branch-gamma"),
			).
			// Move cursor to a non-zero position
			SelectNextItem().
			SelectNextItem().
			Lines(
				Contains("branch-alpha"),
				Contains("branch-beta"),
				Contains("branch-gamma").IsSelected(),
			)

		// Trigger a model update while staying on the Branches view.
		// Using a shell command that creates a new branch sorting after
		// branch-gamma, so the selection index still points to the same item.
		t.GlobalPress(keys.Universal.ExecuteShellCommand)
		t.ExpectPopup().Prompt().
			Title(Equals("Shell command:")).
			Type("git branch branch-zeta").
			Confirm()

		// Verify that the selection is still on branch-gamma (not reset to 0)
		t.Views().Branches().
			Lines(
				Contains("branch-alpha"),
				Contains("branch-beta"),
				Contains("branch-gamma").IsSelected(),
				Contains("branch-zeta"),
			)
	},
})
