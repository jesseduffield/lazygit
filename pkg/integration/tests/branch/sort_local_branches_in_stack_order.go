package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SortLocalBranchesInStackOrder = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Sort each group of branches that share a committer date so that a branch comes before the ones it is based on",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.41.0"),
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "date"
	},
	SetupRepo: func(shell *Shell) {
		// Rebasing a stack of branches gives every commit it creates the same
		// committer date; give each of these two stacks one date of its own for
		// the same effect.
		dateOfFirstStack := "2024-01-01 10:00:00"
		dateOfSecondStack := "2024-01-01 11:00:00"

		shell.
			EmptyCommitWithDate("base", dateOfFirstStack).
			NewBranch("one-bottom").
			EmptyCommitWithDate("one-bottom", dateOfFirstStack).
			NewBranch("one-middle").
			EmptyCommitWithDate("one-middle", dateOfFirstStack).
			NewBranch("one-top").
			EmptyCommitWithDate("one-top", dateOfFirstStack).
			NewBranchFrom("unrelated", "master").
			EmptyCommitWithDate("unrelated", dateOfFirstStack).
			NewBranchFrom("two-bottom", "master").
			EmptyCommitWithDate("two-bottom", dateOfSecondStack).
			NewBranch("two-top").
			EmptyCommitWithDate("two-top", dateOfSecondStack).
			Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("two-top"),
				Contains("two-bottom"),
				Contains("one-top"),
				Contains("one-middle"),
				Contains("one-bottom"),
				Contains("unrelated"),
			)
	},
})
