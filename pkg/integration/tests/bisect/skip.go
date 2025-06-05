package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Skip = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Start a git bisect and skip a few commits (selected or current)",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(10)
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetAppState().GitLogShowGraph = "never"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			SelectedLine(Contains("commit 10")).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`Mark .* as bad`)).Confirm()
			}).
			NavigateToLine(Contains("commit 01")).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`Mark .* as good`)).Confirm()
				t.Views().Information().Content(Contains("Bisecting"))
			}).
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
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect")).
					// Does not show a "Skip selected commit" entry:
					Lines(
						Contains("b Mark current commit").Contains("as bad"),
						Contains("g Mark current commit").Contains("as good"),
						Contains("s Skip current commit"),
						Contains("r Reset bisect"),
						Contains("Cancel"),
					).
					Select(Contains("Skip current commit")).Confirm()
			}).
			// Skipping the current commit selects the new current commit:
			Lines(
				Contains("CI commit 10").Contains("<-- bad"),
				Contains("CI commit 09").DoesNotContain("<--"),
				Contains("CI commit 08").DoesNotContain("<--"),
				Contains("CI commit 07").DoesNotContain("<--"),
				Contains("CI commit 06").Contains("<-- current").IsSelected(),
				Contains("CI commit 05").Contains("<-- skipped"),
				Contains("CI commit 04").DoesNotContain("<--"),
				Contains("CI commit 03").DoesNotContain("<--"),
				Contains("CI commit 02").DoesNotContain("<--"),
				Contains("CI commit 01").Contains("<-- good"),
			).
			NavigateToLine(Contains("commit 07")).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect")).
					// Does show a "Skip selected commit" entry:
					Lines(
						Contains("b Mark current commit").Contains("as bad"),
						Contains("g Mark current commit").Contains("as good"),
						Contains("s Skip current commit"),
						Contains("S Skip selected commit"),
						Contains("r Reset bisect"),
						Contains("Cancel"),
					).
					Select(Contains("Skip selected commit")).Confirm()
			}).
			// Skipping a selected, non-current commit keeps the selection
			// there:
			SelectedLine(Contains("CI commit 07").Contains("<-- skipped"))
	},
})
