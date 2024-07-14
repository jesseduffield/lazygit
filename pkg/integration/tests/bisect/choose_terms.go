package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChooseTerms = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Start a git bisect by choosing 'broken/fixed' as bisect terms",
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
		markCommitAsFixed := func() {
			t.Views().Commits().
				Press(keys.Commits.ViewBisectOptions)

			t.ExpectPopup().Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`Mark .* as fixed`)).Confirm()
		}

		markCommitAsBroken := func() {
			t.Views().Commits().
				Press(keys.Commits.ViewBisectOptions)

			t.ExpectPopup().Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`Mark .* as broken`)).Confirm()
		}

		t.Views().Commits().
			Focus().
			SelectedLine(Contains("CI commit 10")).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect")).Select(Contains("Choose bisect terms")).Confirm()
				t.ExpectPopup().Prompt().Title(Equals("Term for old/good commit:")).Type("broken").Confirm()
				t.ExpectPopup().Prompt().Title(Equals("Term for new/bad commit:")).Type("fixed").Confirm()
			}).
			NavigateToLine(Contains("CI commit 09")).
			Tap(markCommitAsFixed).
			SelectedLine(Contains("<-- fixed")).
			NavigateToLine(Contains("CI commit 02")).
			Tap(markCommitAsBroken).
			Lines(
				Contains("CI commit 10").DoesNotContain("<--"),
				Contains("CI commit 09").Contains("<-- fixed"),
				Contains("CI commit 08").DoesNotContain("<--"),
				Contains("CI commit 07").DoesNotContain("<--"),
				Contains("CI commit 06").DoesNotContain("<--"),
				Contains("CI commit 05").Contains("<-- current").IsSelected(),
				Contains("CI commit 04").DoesNotContain("<--"),
				Contains("CI commit 03").DoesNotContain("<--"),
				Contains("CI commit 02").Contains("<-- broken"),
				Contains("CI commit 01").DoesNotContain("<--"),
			).
			Tap(markCommitAsFixed).
			SelectedLine(Contains("CI commit 04").Contains("<-- current")).
			Tap(func() {
				markCommitAsBroken()

				// commit 5 is the culprit because we marked 4 as broken and 5 as fixed.
				t.ExpectPopup().Alert().Title(Equals("Bisect complete")).Content(MatchesRegexp("(?s)commit 05.*Do you want to reset")).Confirm()
			}).
			IsFocused().
			Content(Contains("CI commit 04"))

		t.Views().Information().Content(DoesNotContain("Bisecting"))
	},
})
