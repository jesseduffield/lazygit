package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SortLocalBranches = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Sort local branches by recency, date or alphabetically",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("commit").
			NewBranch("first").
			EmptyCommitWithDate("commit", "2023-04-07 10:00:00").
			NewBranch("second").
			EmptyCommitWithDate("commit", "2023-04-07 12:00:00").
			NewBranch("third").
			EmptyCommitWithDate("commit", "2023-04-07 11:00:00").
			Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// sorted by recency by default
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("third"),
				Contains("second"),
				Contains("first"),
			).
			SelectNextItem() // to test that the selection jumps back to the top when sorting

		t.Views().Branches().
			Press(keys.Branches.SortOrder)

		t.ExpectPopup().Menu().Title(Equals("Sort order")).
			Select(Contains("-committerdate")).
			Confirm()

		t.Views().Branches().
			IsFocused().
			Lines(
				Contains("master").IsSelected(),
				Contains("second"),
				Contains("third"),
				Contains("first"),
			)

		t.Views().Branches().
			Press(keys.Branches.SortOrder)

		t.ExpectPopup().Menu().Title(Equals("Sort order")).
			Select(Contains("refname")).
			Confirm()

		t.Views().Branches().
			IsFocused().
			Lines(
				Contains("master").IsSelected(),
				Contains("first"),
				Contains("second"),
				Contains("third"),
			)
	},
})
