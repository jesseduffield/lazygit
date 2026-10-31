package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StartFromTagsAndRemoteBranches = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Start a bisect by marking a tag as good and a remote branch as bad",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(2).
			CreateLightweightTag("v1.0", "HEAD").
			CreateNCommitsStartingAt(4, 3).
			CloneIntoRemote("origin").
			HardReset("HEAD~2")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Tags().
			Focus().
			Lines(
				Contains("v1.0").IsSelected(),
			).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect")).
					Lines(
						Contains("b Mark tag 'v1.0' as bad (start bisect)"),
						Contains("g Mark tag 'v1.0' as good (start bisect)"),
						Contains("t Choose bisect terms"),
						Contains("Cancel"),
					).
					Select(Contains("as good")).Confirm()
			})

		t.Views().Remotes().
			Focus().
			Lines(
				Contains("origin").IsSelected(),
			).
			PressEnter()

		t.Views().RemoteBranches().
			IsFocused().
			Lines(
				Contains("master").IsSelected(),
			).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect")).
					Lines(
						Contains("b Mark branch 'origin/master' as bad"),
						Contains("g Mark branch 'origin/master' as good"),
						Contains("r Reset bisect"),
						Contains("Cancel"),
					).
					Select(Contains("as bad")).Confirm()
			})

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("CI commit-06").Contains("<-- bad"),
				Contains("CI commit-05"),
				Contains("CI commit-04").Contains("<-- current").IsSelected(),
				Contains("CI commit-03"),
				Contains("CI v1.0 commit-02").Contains("<-- good"),
				Contains("CI commit-01"),
			)
	},
})
