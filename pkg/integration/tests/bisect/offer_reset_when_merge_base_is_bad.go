package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OfferResetWhenMergeBaseIsBad = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Offer to reset instead of showing the bisect menu when git stopped because the merge base is bad",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(2).
			NewBranch("fixed").
			CreateNCommitsStartingAt(1, 3).
			Checkout("master").
			CreateNCommitsStartingAt(2, 4).
			// Git checks out the merge base commit-02 to test first. If it is
			// bad, git can't find the first bad commit, and stops.
			StartBisect("master", "fixed").
			RunCommandExpectError([]string{"git", "bisect", "bad"})
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("CI commit-05"),
				Contains("CI commit-04"),
				Contains("CI commit-02").Contains("<-- bad").IsSelected(),
				Contains("CI commit-01"),
			).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Confirmation().Title(Equals("Reset 'git bisect'")).Content(Equals("Are you sure you want to reset 'git bisect'?")).Cancel()
			})
	},
})
