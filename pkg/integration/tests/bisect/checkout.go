package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Checkout = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Start a git bisect and checkout a different commit within the bisect range",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(10)
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Log.ShowGraph = "never"
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
			// Verify we're at the bisect-selected commit (commit 05)
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
			// Navigate to a different commit and checkout
			NavigateToLine(Contains("commit 07")).
			PressPrimaryAction()

		// Confirm the checkout menu and select detached head checkout
		t.ExpectPopup().Menu().
			Title(Contains("Checkout branch or commit")).
			Select(MatchesRegexp("Checkout commit .* as detached head")).
			Confirm()

		// After checkout, focus switches to branches panel - just verify it's focused
		t.Views().Branches().
			IsFocused()

		// Switch back to commits panel and verify the current marker moved
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI commit 10").Contains("<-- bad"),
				Contains("CI commit 09").DoesNotContain("<--"),
				Contains("CI commit 08").DoesNotContain("<--"),
				// The current marker should now be on commit 07 (where HEAD is),
				// not on commit 05 (where bisect expected us to be)
				Contains("CI commit 07").Contains("<-- current"),
				Contains("CI commit 06").DoesNotContain("<--"),
				Contains("CI commit 05").DoesNotContain("<--"),
				Contains("CI commit 04").DoesNotContain("<--"),
				Contains("CI commit 03").DoesNotContain("<--"),
				Contains("CI commit 02").DoesNotContain("<--"),
				Contains("CI commit 01").Contains("<-- good"),
			)
	},
})
