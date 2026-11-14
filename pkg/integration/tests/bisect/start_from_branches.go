package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StartFromBranches = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Start a bisect by marking the heads of two branches as good and bad in the branches panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(3).
			NewBranch("feature").
			CreateNCommitsStartingAt(6, 4).
			Checkout("master")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("feature"),
			).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect")).
					Lines(
						Contains("b Mark branch 'master' as bad (start bisect)"),
						Contains("g Mark branch 'master' as good (start bisect)"),
						Contains("t Choose bisect terms"),
						Contains("Cancel"),
					).
					Select(Contains("as good")).Confirm()

				t.Views().Information().Content(Contains("Bisecting"))
			}).
			NavigateToLine(Contains("feature")).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect")).
					Lines(
						Contains("b Mark branch 'feature' as bad"),
						Contains("g Mark branch 'feature' as good"),
						Contains("r Reset bisect"),
						Contains("Cancel"),
					).
					Select(Contains("as bad")).Confirm()
			})

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("CI commit-09").Contains("<-- bad"),
				Contains("CI commit-08"),
				Contains("CI commit-07"),
				Contains("CI commit-06").Contains("<-- current").IsSelected(),
				Contains("CI commit-05"),
				Contains("CI commit-04"),
				Contains("CI commit-03").Contains("<-- good"),
				Contains("CI commit-02"),
				Contains("CI commit-01"),
			)

		// Once git knows both a good and a bad commit, the bisect continues in
		// the commits panel
		t.Views().Branches().
			Focus().
			Press(keys.Commits.ViewBisectOptions)

		t.ExpectToast(Equals("Disabled: Bisecting is in progress; use the commits panel"))
	},
})
