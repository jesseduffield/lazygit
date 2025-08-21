package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var HideMergeCommits = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Test hiding merge commits functionality",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// Create a repo with merge commits to test the functionality
		shell.CreateFile("main.go", "package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}")
		shell.EmptyCommit("Initial commit")

		// Create a feature branch
		shell.NewBranch("feature-branch")
		shell.CreateFile("feature.go", "package main\n\nfunc feature() {\n\tprintln(\"Feature!\")\n}")
		shell.EmptyCommit("Add feature")

		// Switch back to master (default branch) and create another commit
		shell.Checkout("master")
		shell.CreateFile("utils.go", "package main\n\nfunc utils() {\n\tprintln(\"Utils!\")\n}")
		shell.EmptyCommit("Add utils")

		// Merge feature branch into master (this creates a merge commit)
		shell.Merge("feature-branch")

		// Create another commit after merge
		shell.CreateFile("final.go", "package main\n\nfunc final() {\n\tprintln(\"Final!\")\n}")
		shell.EmptyCommit("Final commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			// Verify that merge commit is initially visible
			Content(Contains(`Merge branch 'feature-branch'`)).
			Content(Contains(`Final commit`)).
			Content(Contains(`Add utils`)).
			Content(Contains(`Add feature`)).
			Content(Contains(`Initial commit`)).
			// Open filtering menu
			Press(keys.Universal.FilteringMenu).
			// Check that Hide merge commits option is available and toggle it
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Filtering")).
					Select(Contains("Hide merge commits")).
					Confirm()
			}).
			// Verify that merge commit is hidden
			Content(DoesNotContain(`Merge branch 'feature-branch'`)).
			Content(Contains(`Final commit`)).
			Content(Contains(`Add utils`)).
			Content(Contains(`Add feature`)).
			Content(Contains(`Initial commit`)).
			// Open filtering menu again to toggle off
			Press(keys.Universal.FilteringMenu).
			// Verify Hide merge commits is still selected and toggle it off
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Filtering")).
					Select(Contains("Hide merge commits")).
					Confirm()
			}).
			// Verify that merge commit is visible again
			Content(Contains(`Merge branch 'feature-branch'`)).
			Content(Contains(`Final commit`)).
			Content(Contains(`Add utils`)).
			Content(Contains(`Add feature`)).
			Content(Contains(`Initial commit`))
	},
})
