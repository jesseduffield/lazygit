package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveRangeSpanningBranchBoundaryOutsideRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a range of commits that spans a branch boundary in a stack of branches",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.MainBranches = []string{"master"}
		config.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(1).
			NewBranch("branch1").
			CreateNCommitsStartingAt(2, 2).
			NewBranch("branch2").
			CreateNCommitsStartingAt(2, 4).
			NewBranch("branch3").
			CreateNCommitsStartingAt(2, 6)

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI commit-07").IsSelected(),
				Contains("CI commit-06"),
				Contains("CI * commit-05"),
				Contains("CI commit-04"),
				Contains("CI * commit-03"),
				Contains("CI commit-02"),
				Contains("CI commit-01"),
			).
			NavigateToLine(Contains("commit-06")).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Commits.MoveUpCommit).
			Lines(
				Contains("CI commit-06").IsSelected(),
				/* EXPECTED:
				Contains("CI * commit-05").IsSelected(),
				ACTUAL: */
				Contains("CI commit-05").IsSelected(),
				Contains("CI commit-04").IsSelected(),
				Contains("CI commit-07"),
				Contains("CI * commit-03"),
				Contains("CI commit-02"),
				Contains("CI commit-01"),
			)
	},
})
