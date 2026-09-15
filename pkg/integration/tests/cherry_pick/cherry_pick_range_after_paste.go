package cherry_pick

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CherryPickRangeAfterPaste = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Regression test: range-copy multiple commits after a previous paste",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "recency"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("base").
			NewBranch("target").
			NewBranch("source").
			EmptyCommit("one").
			EmptyCommit("two").
			EmptyCommit("three").
			EmptyCommit("four").
			EmptyCommit("five").
			Checkout("target")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("target").IsSelected(),
				Contains("source"),
				Contains("master"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("five").IsSelected(),
				Contains("four"),
				Contains("three"),
				Contains("two"),
				Contains("one"),
				Contains("base"),
			).
			Press(keys.Commits.CherryPickCopy)

		t.Views().Commits().
			Focus().
			Lines(
				Contains("base").IsSelected(),
			).
			Press(keys.Commits.PasteCommits).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Cherry-pick")).
					Content(Equals("Are you sure you want to cherry-pick the 1 copied commit(s) onto this branch?")).
					Confirm()
			}).
			Lines(
				Contains("five"),
				Contains("base").IsSelected(),
			).
			Tap(func() {
				// After paste, CherryPicking.DidPaste is true, so it looks to the user as if no
				// commits are copied:
				t.Views().Information().Content(DoesNotContain("commits copied"))
			})

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("source")).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			NavigateToLine(Contains("four")).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Commits.CherryPickCopy).
			Tap(func() {
				t.Views().Information().Content(Contains("3 commits copied"))
			})

		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("base")).
			Press(keys.Commits.PasteCommits).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Cherry-pick")).
					Content(Equals("Are you sure you want to cherry-pick the 3 copied commit(s) onto this branch?")).
					Confirm()
			})

		t.Views().Commits().Lines(
			Contains("four"),
			Contains("three"),
			Contains("two"),
			Contains("five"),
			Contains("base").IsSelected(),
		)
	},
})
