package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StartFromSubCommitsAndReflog = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Start a bisect by marking a commit of a branch as bad and a reflog entry as good",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(3).
			NewBranch("feature").
			CreateNCommitsStartingAt(3, 4).
			Checkout("master")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("feature")).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			SelectedLine(Contains("CI commit-06")).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect")).
					Lines(
						MatchesRegexp(`b Mark commit [0-9a-f]+ as bad \(start bisect\)`),
						MatchesRegexp(`g Mark commit [0-9a-f]+ as good \(start bisect\)`),
						Contains("t Choose bisect terms"),
						Contains("Cancel"),
					).
					Select(Contains("as bad")).Confirm()
			}).
			PressEscape()

		t.Views().ReflogCommits().
			Focus().
			NavigateToLine(Contains("commit-01")).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect")).
					Lines(
						MatchesRegexp(`b Mark commit [0-9a-f]+ as bad`),
						MatchesRegexp(`g Mark commit [0-9a-f]+ as good`),
						Contains("r Reset bisect"),
						Contains("Cancel"),
					).
					Select(Contains("as good")).Confirm()
			})

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("CI commit-06").Contains("<-- bad"),
				Contains("CI commit-05"),
				Contains("CI commit-04"),
				Contains("CI commit-03").Contains("<-- current").IsSelected(),
				Contains("CI commit-02"),
				Contains("CI commit-01").Contains("<-- good"),
			)
	},
})
