package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowResultAgain = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show the result of a finished bisect again instead of the bisect menu",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(10).
			StartBisect("HEAD~2", "HEAD~5").
			// Mark commit-07, the commit that git checked out to test. This
			// finishes the bisect, because commit-08 is bad.
			RunCommand([]string{"git", "bisect", "good"})
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsFocused().
			SelectedLine(Contains("CI commit-08").Contains("<-- bad")).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Alert().Title(Equals("Bisect complete")).Content(MatchesRegexp("(?s)commit-08.*Do you want to reset")).Cancel()
			})

		t.Views().Branches().
			Focus().
			Press(keys.Commits.ViewBisectOptions)

		t.ExpectToast(Equals("Disabled: Bisecting is in progress; use the commits panel"))
	},
})
