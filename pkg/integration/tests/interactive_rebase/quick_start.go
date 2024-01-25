package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var QuickStart = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Quick-starts an interactive rebase in several contexts",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// we're going to test the following:
		// * quick start from main fails
		// * quick start from feature branch starts from main
		// * quick start from branch with merge commit starts from merge commit

		shell.NewBranch("main")
		shell.EmptyCommit("initial commit")
		shell.EmptyCommit("last main commit")

		shell.NewBranch("feature-branch")
		shell.NewBranch("branch-to-merge")
		shell.NewBranch("branch-with-merge-commit")

		shell.Checkout("feature-branch")
		shell.EmptyCommit("feature-branch one")
		shell.EmptyCommit("feature-branch two")

		shell.Checkout("branch-to-merge")
		shell.EmptyCommit("branch-to-merge one")
		shell.EmptyCommit("branch-to-merge two")

		shell.Checkout("branch-with-merge-commit")
		shell.EmptyCommit("branch-with-merge one")
		shell.EmptyCommit("branch-with-merge two")

		shell.Merge("branch-to-merge")

		shell.EmptyCommit("branch-with-merge three")

		shell.Checkout("main")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("last main commit"),
				Contains("initial commit"),
			).
			// Verify we can't quick start from main
			Press(keys.Commits.StartInteractiveRebase)

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Equals("Cannot start interactive rebase: the HEAD commit is a merge commit or is present on the main branch, so there is no appropriate base commit to start the rebase from. You can start an interactive rebase from a specific commit by selecting the commit and pressing `e`.")).
			Confirm()

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("feature-branch")).
			Press(keys.Universal.Select)

		t.Views().Commits().
			Focus().
			Lines(
				Contains("feature-branch two").IsSelected(),
				Contains("feature-branch one"),
				Contains("last main commit"),
				Contains("initial commit"),
			).
			// Verify quick start picks the last commit on the main branch
			Press(keys.Commits.StartInteractiveRebase).
			Lines(
				Contains("feature-branch two").IsSelected(),
				Contains("feature-branch one"),
				Contains("last main commit").Contains("YOU ARE HERE"),
				Contains("initial commit"),
			).
			// Try again, verify we fail because we're already rebasing
			Press(keys.Commits.StartInteractiveRebase)

		t.ExpectToast(Equals("Disabled: Can't perform this action during a rebase"))
		t.Common().AbortRebase()

		// Verify if a merge commit is present on the branch we start from there
		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("branch-with-merge-commit")).
			Press(keys.Universal.Select)

		t.Views().Commits().
			Focus().
			Lines(
				Contains("branch-with-merge three").IsSelected(),
				Contains("Merge branch 'branch-to-merge'"),
				Contains("branch-to-merge two"),
				Contains("branch-to-merge one"),
				Contains("branch-with-merge two"),
				Contains("branch-with-merge one"),
				Contains("last main commit"),
				Contains("initial commit"),
			).
			Press(keys.Commits.StartInteractiveRebase).
			Lines(
				Contains("branch-with-merge three").IsSelected(),
				Contains("Merge branch 'branch-to-merge'").Contains("YOU ARE HERE"),
				Contains("branch-to-merge two"),
				Contains("branch-to-merge one"),
				Contains("branch-with-merge two"),
				Contains("branch-with-merge one"),
				Contains("last main commit"),
				Contains("initial commit"),
			)
	},
})
