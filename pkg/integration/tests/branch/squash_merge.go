package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SquashMerge = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Squash merge a branch both with and without committing",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("original-branch").
			EmptyCommit("one").
			NewBranch("change-worktree-branch").
			CreateFileAndAdd("work", "content").
			Commit("work").
			Checkout("original-branch").
			NewBranch("change-commit-branch").
			CreateFileAndAdd("file", "content").
			Commit("file").
			Checkout("original-branch")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().TopLines(
			Contains("one"),
		)

		t.Views().Branches().
			Focus().
			Lines(
				Contains("original-branch").IsSelected(),
				Contains("change-commit-branch"),
				Contains("change-worktree-branch"),
			).
			SelectNextItem().
			Press(keys.Branches.MergeIntoCurrentBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Merge")).
			Select(Contains("Squash merge and commit")).
			Confirm()

		t.Views().Commits().TopLines(
			Contains("Squash merge change-commit-branch into original-branch"),
			Contains("one"),
		)

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("change-worktree-branch")).
			Press(keys.Branches.MergeIntoCurrentBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Merge")).
			Select(Contains("Squash merge and leave uncommitted")).
			Confirm()

		t.Views().Files().Focus().Lines(
			Contains("work"),
		)
	},
})
