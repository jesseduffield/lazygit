package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CheckoutDuringBisect = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Checkout a different commit during a bisect and verify the current marker follows HEAD",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			NewBranch("mybranch").
			CreateNCommits(10)
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			SelectedLine(Contains("CI commit 10")).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`Mark .* as bad`)).Confirm()
			}).
			NavigateToLine(Contains("CI commit 01")).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`Mark .* as good`)).Confirm()
			}).
			// bisect has auto-selected commit 05 as current
			Lines(
				Contains("CI commit 10").Contains("<-- bad"),
				Contains("CI commit 09").DoesNotContain("<--"),
				Contains("CI commit 08").DoesNotContain("<--"),
				Contains("CI commit 07").DoesNotContain("<--"),
				Contains("CI commit 06").DoesNotContain("<--"),
				Contains("CI commit 05").Contains("<-- current").IsSelected(),
				Contains("CI commit 04").DoesNotContain("<--"),
				Contains("CI commit 03").DoesNotContain("<--"),
				Contains("CI commit 02").DoesNotContain("<--"),
				Contains("CI commit 01").Contains("<-- good"),
			).
			// now checkout commit 08 manually
			NavigateToLine(Contains("CI commit 08")).
			Press(keys.Commits.CheckoutCommit).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Contains("Checkout branch or commit")).
					Select(MatchesRegexp("Checkout commit .* as detached head")).
					Confirm()
			})

		// after checkout, go back to commits panel and verify the current
		// marker moved to commit 08
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI commit 10").Contains("<-- bad"),
				Contains("CI commit 09").DoesNotContain("<--"),
				Contains("CI commit 08").Contains("<-- current"),
				Contains("CI commit 07").DoesNotContain("<--"),
				Contains("CI commit 06").DoesNotContain("<--"),
				Contains("CI commit 05").DoesNotContain("<--"),
				Contains("CI commit 04").DoesNotContain("<--"),
				Contains("CI commit 03").DoesNotContain("<--"),
				Contains("CI commit 02").DoesNotContain("<--"),
				Contains("CI commit 01").Contains("<-- good"),
			).
			// mark the manually checked-out commit as good via bisect menu
			NavigateToLine(Contains("CI commit 08")).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				// the bisect menu should show the HEAD commit (08) as current
				t.ExpectPopup().Menu().Title(Equals("Bisect")).
					Select(MatchesRegexp(`Mark .* as good`)).Confirm()
			}).
			// after marking 08 as good, bisect narrows the range
			SelectedLine(Contains("<-- current"))
	},
})
