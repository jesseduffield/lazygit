package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveAcrossBranchBoundaryOutsideRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a commit across a branch boundary in a stack of branches",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.MainBranches = []string{"master"}
		config.GetAppState().GitLogShowGraph = "never"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(1).
			NewBranch("branch1").
			CreateNCommitsStartingAt(2, 2).
			NewBranch("branch2").
			CreateNCommitsStartingAt(2, 4)

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI commit 05").IsSelected(),
				Contains("CI commit 04"),
				Contains("CI * commit 03"),
				Contains("CI commit 02"),
				Contains("CI commit 01"),
			).
			NavigateToLine(Contains("commit 04")).
			Press(keys.Commits.MoveDownCommit).
			Lines(
				Contains("CI commit 05"),
				Contains("CI * commit 03"),
				Contains("CI commit 04").IsSelected(),
				Contains("CI commit 02"),
				Contains("CI commit 01"),
			)
	},
})
