package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MergeNonFastForward = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Merge a branch into another using non-fast-forward merge",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("original-branch").
			EmptyCommit("one").
			NewBranch("branch1").
			EmptyCommit("branch1").
			Checkout("original-branch").
			NewBranchFrom("branch2", "original-branch").
			EmptyCommit("branch2").
			Checkout("original-branch")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("original-branch").IsSelected(),
				Contains("branch1"),
				Contains("branch2"),
			).
			SelectNextItem().
			Press(keys.Branches.MergeIntoCurrentBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Merge")).
			TopLines(
				Contains("Regular merge (fast-forward)"),
				Contains("Regular merge (with merge commit)"),
			).
			Select(Contains("Regular merge (with merge commit)")).
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("⏣─╮ Merge branch 'branch1' into original-branch").IsSelected(),
				Contains("│ ◯ * branch1"),
				Contains("◯─╯ one"),
			)

		// Check that branch2 shows the non-fast-forward option first
		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("branch2")).
			Press(keys.Branches.MergeIntoCurrentBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Merge")).
			TopLines(
				Contains("Regular merge (with merge commit)"),
				Contains("Regular merge (fast-forward)"),
			).
			Select(Contains("Regular merge (with merge commit)")).
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("⏣─╮ Merge branch 'branch2' into original-branch").IsSelected(),
				Contains("│ ◯ * branch2"),
				Contains("⏣─│─╮ Merge branch 'branch1' into original-branch"),
				Contains("│ │ ◯ * branch1"),
				Contains("◯─┴─╯ one"),
			)
	},
})
