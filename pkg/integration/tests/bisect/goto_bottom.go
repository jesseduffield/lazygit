package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var GotoBottom = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Go to the bottom of the commits list in steps, stopping only at the bad commit, the current commit, and the first good commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(3).
			NewBranch("mybranch").
			CreateNCommitsStartingAt(7, 4).
			StartBisect("HEAD~1", "HEAD~8")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI commit-10"),
				Contains("CI commit-09").Contains("<-- bad"),
				Contains("CI commit-08"),
				Contains("CI commit-07"),
				Contains("CI commit-06"),
				Contains("CI commit-05").Contains("<-- current").IsSelected(),
				Contains("CI commit-04"),
				Contains("CI commit-03"),
				Contains("CI commit-02").Contains("<-- good"),
				Contains("CI commit-01"),
			).
			Press(keys.Universal.GotoTop).
			Press(keys.Universal.GotoBottom).
			SelectedLine(Contains("CI commit-09")).
			Press(keys.Universal.GotoBottom).
			SelectedLine(Contains("CI commit-05")).
			Press(keys.Universal.GotoBottom).
			SelectedLine(Contains("CI commit-02")).
			Press(keys.Universal.GotoBottom).
			SelectedLine(Contains("CI commit-01"))
	},
})
