package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Checkout = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Test checkout during bisect shows HEAD position correctly",
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
		// Start bisect by marking commit 09 as bad and commit 02 as good
		t.Views().Commits().
			Focus().
			SelectedLine(Contains("CI commit 10")).
			NavigateToLine(Contains("CI commit 09")).
			Press(keys.Commits.ViewBisectOptions)

		t.ExpectPopup().Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`Mark .* as bad`)).Confirm()

		t.Views().Information().Content(Contains("Bisecting"))

		t.Views().Commits().
			IsFocused().
			NavigateToLine(Contains("CI commit 02")).
			Press(keys.Commits.ViewBisectOptions)

		t.ExpectPopup().Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`Mark .* as good`)).Confirm()

		// Now we're bisecting, we should be at commit 05 (the midpoint)
		t.Views().Commits().
			IsFocused().
			SelectedLine(Contains("CI commit 05").Contains("<-- current"))

		// Try to checkout a different commit (commit 03)
		t.Views().Commits().
			NavigateToLine(Contains("CI commit 03")).
			PressPrimaryAction()

		// The checkout menu should appear
		t.ExpectPopup().Menu().
			Title(Contains("Checkout branch or commit")).
			Select(MatchesRegexp("Checkout commit .* as detached head")).
			Confirm()

		// After checkout, focus moves to branches panel
		t.Views().Branches().
			IsFocused()

		// Now go back to commits panel and verify HEAD position is shown
		t.Views().Commits().
			Focus().
			// Commit 05 should still show "<-- current" (bisect expected)
			// Commit 03 should show "<-- YOU ARE HERE" (actual HEAD position)
			Content(
				Contains("<-- bad").
					Contains("<-- current").
					Contains("<-- YOU ARE HERE").
					Contains("<-- good"),
			)

		// Bisect should still be active
		t.Views().Information().Content(Contains("Bisecting"))
	},
})
