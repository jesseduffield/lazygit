package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RebaseFromMarkedBase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase onto another branch from a marked base commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			NewBranch("base-branch").
			EmptyCommit("one").
			EmptyCommit("two").
			EmptyCommit("three").
			NewBranch("active-branch").
			EmptyCommit("active one").
			EmptyCommit("active two").
			EmptyCommit("active three").
			Checkout("base-branch").
			NewBranch("target-branch").
			EmptyCommit("target one").
			EmptyCommit("target two").
			Checkout("active-branch")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("active three"),
				Contains("active two"),
				Contains("active one"),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			).
			NavigateToLine(Contains("active one")).
			Press(keys.Commits.MarkCommitAsBaseForRebase).
			Lines(
				Contains("active three").Contains("✓"),
				Contains("active two").Contains("✓"),
				Contains("↑↑↑ Will rebase from here ↑↑↑ active one"),
				Contains("three").DoesNotContain("✓"),
				Contains("two").DoesNotContain("✓"),
				Contains("one").DoesNotContain("✓"),
			)

		t.Views().Information().Content(Contains("Marked a base commit for rebase"))

		t.Views().Branches().
			Focus().
			Lines(
				Contains("active-branch"),
				Contains("target-branch"),
				Contains("base-branch"),
			).
			SelectNextItem().
			Press(keys.Branches.RebaseBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Rebase 'active-branch' from marked base onto 'target-branch'")).
			Select(Contains("Simple rebase")).
			Confirm()

		t.Views().Commits().Lines(
			Contains("active three").DoesNotContain("✓"),
			Contains("active two").DoesNotContain("✓"),
			Contains("target two").DoesNotContain("✓"),
			Contains("target one").DoesNotContain("✓"),
			Contains("three").DoesNotContain("✓"),
			Contains("two").DoesNotContain("✓"),
			Contains("one").DoesNotContain("✓"),
		)
	},
})
