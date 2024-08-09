package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Basic = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Start a git bisect to find a bad commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			NewBranch("mybranch").
			CreateNCommits(10)
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetAppState().GitLogShowGraph = "never"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		markCommitAsBad := func() {
			t.Views().Commits().
				Press(keys.Commits.ViewBisectOptions)

			t.ExpectPopup().Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`Mark .* as bad`)).Confirm()
		}

		markCommitAsGood := func() {
			t.Views().Commits().
				Press(keys.Commits.ViewBisectOptions)

			t.ExpectPopup().Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`Mark .* as good`)).Confirm()
		}

		t.Views().Commits().
			Focus().
			SelectedLine(Contains("CI commit 10")).
			NavigateToLine(Contains("CI commit 09")).
			Tap(func() {
				markCommitAsBad()

				t.Views().Information().Content(Contains("Bisecting"))
			}).
			SelectedLine(Contains("<-- bad")).
			NavigateToLine(Contains("CI commit 02")).
			Tap(markCommitAsGood).
			TopLines(Contains("CI commit 10")).
			// lazygit will land us in the commit between our good and bad commits.
			SelectedLine(Contains("CI commit 05").Contains("<-- current")).
			Tap(markCommitAsBad).
			SelectedLine(Contains("CI commit 04").Contains("<-- current")).
			Tap(func() {
				markCommitAsGood()

				// commit 5 is the culprit because we marked 4 as good and 5 as bad.
				t.ExpectPopup().Alert().Title(Equals("Bisect complete")).Content(MatchesRegexp("(?s)commit 05.*Do you want to reset")).Confirm()
			}).
			IsFocused().
			Content(Contains("CI commit 04"))

		t.Views().Information().Content(DoesNotContain("Bisecting"))
	},
})
