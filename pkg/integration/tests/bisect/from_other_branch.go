package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FromOtherBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Opening lazygit when bisect has been started from another branch. There's an issue where we don't reselect the current branch if we mark the current branch as bad so this test side-steps that problem",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("only commit on master"). // this'll ensure we have a master branch
			NewBranch("other").
			CreateNCommits(10).
			Checkout("master").
			StartBisect("other~2", "other~5")
	},
	SetupConfig: func(cfg *config.AppConfig) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Information().Content(Contains("Bisecting"))

		t.Views().Commits().
			Focus().
			TopLines(
				MatchesRegexp(`<-- bad.*commit 08`),
				MatchesRegexp(`<-- current.*commit 07`),
				MatchesRegexp(`\?.*commit 06`),
				MatchesRegexp(`<-- good.*commit 05`),
			).
			SelectNextItem().
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`Mark .* as good`)).Confirm()

				t.ExpectPopup().Alert().Title(Equals("Bisect complete")).Content(MatchesRegexp("(?s)commit 08.*Do you want to reset")).Confirm()

				t.Views().Information().Content(DoesNotContain("Bisecting"))
			}).
			// back in master branch which just had the one commit
			Lines(
				Contains("only commit on master"),
			)
	},
})
