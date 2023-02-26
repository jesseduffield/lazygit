package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Rebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Begins an interactive rebase, then fixups, drops, and squashes some commits",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.EmptyCommit("first commit to edit")
		shell.EmptyCommit("commit to squash")
		shell.EmptyCommit("second commit to edit")
		shell.EmptyCommit("commit to drop")

		shell.CreateFileAndAdd("fixup-commit-file", "fixup-commit-file")
		shell.Commit("commit to fixup")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit to fixup"),
				Contains("commit to drop"),
				Contains("second commit to edit"),
				Contains("commit to squash"),
				Contains("first commit to edit"),
				Contains("initial commit"),
			).
			NavigateToLine(Contains("first commit to edit")).
			Press(keys.Universal.Edit).
			Lines(
				MatchesRegexp("pick.*commit to fixup"),
				MatchesRegexp("pick.*commit to drop"),
				MatchesRegexp("pick.*second commit to edit"),
				MatchesRegexp("pick.*commit to squash"),
				MatchesRegexp("YOU ARE HERE.*first commit to edit").IsSelected(),
				Contains("initial commit"),
			).
			SelectPreviousItem().
			Press(keys.Commits.SquashDown).
			Lines(
				MatchesRegexp("pick.*commit to fixup"),
				MatchesRegexp("pick.*commit to drop"),
				MatchesRegexp("pick.*second commit to edit"),
				MatchesRegexp("squash.*commit to squash").IsSelected(),
				MatchesRegexp("YOU ARE HERE.*first commit to edit"),
				Contains("initial commit"),
			).
			SelectPreviousItem().
			Press(keys.Universal.Edit).
			Lines(
				MatchesRegexp("pick.*commit to fixup"),
				MatchesRegexp("pick.*commit to drop"),
				MatchesRegexp("edit.*second commit to edit").IsSelected(),
				MatchesRegexp("squash.*commit to squash"),
				MatchesRegexp("YOU ARE HERE.*first commit to edit"),
				Contains("initial commit"),
			).
			SelectPreviousItem().
			Press(keys.Universal.Remove).
			Lines(
				MatchesRegexp("pick.*commit to fixup"),
				MatchesRegexp("drop.*commit to drop").IsSelected(),
				MatchesRegexp("edit.*second commit to edit"),
				MatchesRegexp("squash.*commit to squash"),
				MatchesRegexp("YOU ARE HERE.*first commit to edit"),
				Contains("initial commit"),
			).
			SelectPreviousItem().
			Press(keys.Commits.MarkCommitAsFixup).
			Lines(
				MatchesRegexp("fixup.*commit to fixup").IsSelected(),
				MatchesRegexp("drop.*commit to drop"),
				MatchesRegexp("edit.*second commit to edit"),
				MatchesRegexp("squash.*commit to squash"),
				MatchesRegexp("YOU ARE HERE.*first commit to edit"),
				Contains("initial commit"),
			).
			Tap(func() {
				t.Common().ContinueRebase()
			}).
			Lines(
				MatchesRegexp("fixup.*commit to fixup").IsSelected(),
				MatchesRegexp("drop.*commit to drop"),
				MatchesRegexp("YOU ARE HERE.*second commit to edit"),
				MatchesRegexp("first commit to edit"),
				Contains("initial commit"),
			).
			Tap(func() {
				t.Common().ContinueRebase()
			}).
			Lines(
				Contains("second commit to edit").IsSelected(),
				Contains("first commit to edit"),
				Contains("initial commit"),
			).
			Tap(func() {
				// commit 4 was squashed into 6 so we assert that their messages have been concatenated
				t.Views().Main().Content(
					Contains("second commit to edit").
						// file from fixup commit is present
						Contains("fixup-commit-file").
						// but message is not (because it's a fixup, not a squash)
						DoesNotContain("commit to fixup"),
				)
			}).
			SelectNextItem().
			Tap(func() {
				// commit 4 was squashed into 6 so we assert that their messages have been concatenated
				t.Views().Main().Content(
					Contains("first commit to edit").
						// message from squashed commit has been concatenated with message other commit
						Contains("commit to squash"),
				)
			})
	},
})
