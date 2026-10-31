package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MarkCandidatesWhileTestingMergeBase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Mark the commits that the bisect hasn't ruled out yet while git tests the merge base of the good and the bad commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(2).
			NewBranch("fixed").
			CreateNCommitsStartingAt(1, 3).
			Checkout("master").
			CreateNCommitsStartingAt(2, 4).
			// The good commit-03 is on another branch, so git checks out the
			// merge base commit-02 to test first
			StartBisect("master", "fixed")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("CI commit-05").Contains("<-- bad"),
				MatchesRegexp(`\?.*commit-04`),
				Contains("CI commit-02").Contains("<-- current").IsSelected(),
				// This is an ancestor of the good commit-03, so the bisect has
				// ruled it out
				Contains("CI commit-01").DoesNotContain("?"),
			)
	},
})
