package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CheckOutBranchBeforeBisecting = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Check out a different branch after marking a commit as good, but before marking one as bad",
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
		t.Views().Commits().
			Focus().
			SelectedLine(Contains("CI commit-03")).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`Mark .* as good`)).Confirm()
			}).
			Lines(
				Contains("CI commit-03").Contains("<-- good"),
				Contains("CI commit-02"),
				Contains("CI commit-01"),
			)

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("feature")).
			PressPrimaryAction().
			Lines(
				Contains("* feature"),
				Contains("master"),
			)

		t.Views().Commits().
			Lines(
				Contains("CI commit-06"),
				Contains("CI commit-05"),
				Contains("CI commit-04"),
				Contains("CI commit-03").Contains("<-- good"),
				Contains("CI commit-02"),
				Contains("CI commit-01"),
			)
	},
})
