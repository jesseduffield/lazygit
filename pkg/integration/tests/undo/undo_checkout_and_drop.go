package undo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var UndoCheckoutAndDrop = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Drop some commits and then undo/redo the actions",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")
		shell.EmptyCommit("three")
		shell.EmptyCommit("four")

		shell.NewBranch("other_branch")
		shell.Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// we're going to drop a commit, switch branch, drop a commit there, then undo everything, then redo everything.

		confirmCommitDrop := func() {
			t.ExpectPopup().Confirmation().
				Title(Equals("Delete commit")).
				Content(Equals("Are you sure you want to delete this commit?")).
				Confirm()
		}

		confirmUndoDrop := func() {
			t.ExpectPopup().Confirmation().
				Title(Equals("Undo")).
				Content(MatchesRegexp(`Are you sure you want to hard reset to '.*'\? An auto-stash will be performed if necessary\.`)).
				Confirm()
		}

		confirmRedoDrop := func() {
			t.ExpectPopup().Confirmation().
				Title(Equals("Redo")).
				Content(MatchesRegexp(`Are you sure you want to hard reset to '.*'\? An auto-stash will be performed if necessary\.`)).
				Confirm()
		}

		t.Views().Commits().Focus().
			Lines(
				Contains("four").IsSelected(),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			).
			Press(keys.Universal.Remove).
			Tap(confirmCommitDrop).
			Lines(
				Contains("three").IsSelected(),
				Contains("two"),
				Contains("one"),
			)

		t.Views().Branches().Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("other_branch"),
			).
			SelectNextItem().
			// checkout branch
			PressPrimaryAction().
			Lines(
				Contains("other_branch").IsSelected(),
				Contains("master"),
			)

		// drop the commit in the 'other_branch' branch too
		t.Views().Commits().Focus().
			Lines(
				Contains("four").IsSelected(),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			).
			Press(keys.Universal.Remove).
			Tap(confirmCommitDrop).
			Lines(
				Contains("three").IsSelected(),
				Contains("two"),
				Contains("one"),
			).
			Press(keys.Universal.Undo).
			Tap(confirmUndoDrop).
			Lines(
				Contains("four").IsSelected(),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			).
			Press(keys.Universal.Undo).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Undo")).
					Content(Contains("Are you sure you want to checkout 'master'?")).
					Confirm()

				t.Views().Branches().
					Lines(
						Contains("master").IsSelected(),
						Contains("other_branch"),
					)
			}).
			Lines(
				Contains("three").IsSelected(),
				Contains("two"),
				Contains("one"),
			).
			Press(keys.Universal.Undo).
			Tap(confirmUndoDrop).
			Lines(
				Contains("four").IsSelected(),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			).
			Press(keys.Universal.Redo).
			Tap(confirmRedoDrop).
			Lines(
				Contains("three").IsSelected(),
				Contains("two"),
				Contains("one"),
			).
			Press(keys.Universal.Redo).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Redo")).
					Content(Contains("Are you sure you want to checkout 'other_branch'?")).
					Confirm()

				t.Views().Branches().
					Lines(
						Contains("other_branch").IsSelected(),
						Contains("master"),
					)
			}).
			Press(keys.Universal.Redo).
			Tap(confirmRedoDrop).
			Lines(
				Contains("three").IsSelected(),
				Contains("two"),
				Contains("one"),
			)
	},
})
