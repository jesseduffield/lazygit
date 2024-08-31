package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var QuickStartKeepSelectionRange = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Starts an interactive rebase and checks that the same commit range stays selected",
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
			CreateNCommitsStartingAt(2, 4).
			NewBranch("branch3").
			CreateNCommitsStartingAt(2, 6)

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("commit 04")).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.RangeSelectDown).
			Lines(
				Contains("CI commit 07"),
				Contains("CI commit 06"),
				Contains("CI * commit 05"),
				Contains("CI commit 04").IsSelected(),
				Contains("CI * commit 03").IsSelected(),
				Contains("CI commit 02").IsSelected(),
				Contains("CI commit 01"),
			).
			Press(keys.Commits.StartInteractiveRebase).
			Lines(
				Contains("CI commit 07"),
				Contains("CI commit 06"),
				Contains("update-ref").Contains("branch2"),
				Contains("CI commit 05"),
				Contains("CI commit 04").IsSelected(),
				Contains("update-ref").Contains("branch1").IsSelected(),
				Contains("CI commit 03").IsSelected(),
				Contains("CI commit 02").IsSelected(),
				Contains("CI <-- YOU ARE HERE --- commit 01"),
			)
	},
})
