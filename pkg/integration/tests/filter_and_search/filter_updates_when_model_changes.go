package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterUpdatesWhenModelChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that after deleting a branch the filter is reapplied to show only the remaining branches",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.NewBranch("branch-to-delete")
		shell.NewBranch("other")
		shell.NewBranch("checked-out-branch")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("checked-out-branch").IsSelected(),
				Contains("other"),
				Contains("branch-to-delete"),
				Contains("master"),
			).
			FilterOrSearch("branch").
			Lines(
				Contains("checked-out-branch").IsSelected(),
				Contains("branch-to-delete"),
			).
			SelectNextItem().
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Delete branch 'branch-to-delete'?")).
					Select(Contains("Delete local branch")).
					Confirm()
			}).
			Lines(
				Contains("checked-out-branch").IsSelected(),
			)

		// Verify that updating the filter works even if the view is not the active one
		t.Views().Files().Focus()

		// To do that, we use a custom command to create a new branch that matches the filter
		t.GlobalPress(keys.Universal.ExecuteShellCommand)
		t.ExpectPopup().Prompt().
			Title(Equals("Shell command:")).
			Type("git branch new-branch").
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("checked-out-branch").IsSelected(),
				Contains("new-branch"),
			)

		t.Views().Branches().
			Focus().
			// cancel the filter
			PressEscape().
			Lines(
				Contains("checked-out-branch").IsSelected(),
				Contains("other"),
				Contains("master"),
				Contains("new-branch"),
			)
	},
})
