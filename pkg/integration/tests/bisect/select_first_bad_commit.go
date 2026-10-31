package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectFirstBadCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Select the first bad commit when marking the last commit to test finishes the bisect",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(10).
			StartBisect("HEAD~2", "HEAD~5")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsFocused().
			TopLines(
				Contains("CI commit-10"),
				Contains("CI commit-09"),
				Contains("CI commit-08").Contains("<-- bad"),
				Contains("CI commit-07").Contains("<-- current").IsSelected(),
				Contains("CI commit-06"),
				Contains("CI commit-05").Contains("<-- good"),
			).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`Mark .* as good`)).Confirm()

				t.ExpectPopup().Alert().Title(Equals("Bisect complete")).Content(MatchesRegexp("(?s)commit-08.*Do you want to reset")).Cancel()
			}).
			SelectedLine(Contains("CI commit-08").Contains("<-- bad"))
	},
})
