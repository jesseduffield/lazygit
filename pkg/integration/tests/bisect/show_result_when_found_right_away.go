package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowResultWhenFoundRightAway = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show the result instead of the bisect menu when git found the first bad commit without testing any",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(3).
			// The good commit is the parent of the bad one, so git doesn't
			// check out anything to test
			StartBisect("HEAD", "HEAD~1")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("CI commit-03").Contains("<-- bad").IsSelected(),
				Contains("CI commit-02").Contains("<-- good"),
				Contains("CI commit-01"),
			).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				/* EXPECTED:
				t.ExpectPopup().Alert().Title(Equals("Bisect complete")).Content(MatchesRegexp("(?s)commit-03.*Do you want to reset")).Cancel()
				ACTUAL: */
				t.ExpectPopup().Menu().Title(Equals("Bisect")).TopLines(MatchesRegexp(`b Mark commit [0-9a-f]+ as bad`)).Cancel()
			})

		t.Views().Branches().
			Focus().
			Press(keys.Commits.ViewBisectOptions)

		t.ExpectToast(Equals("Disabled: Bisecting is in progress; use the commits panel"))
	},
})
