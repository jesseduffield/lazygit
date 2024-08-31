package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var QuickStartKeepSelection = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Starts an interactive rebase and checks that the same commit stays selected",
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
			CreateNCommitsStartingAt(3, 2).
			NewBranch("branch2").
			CreateNCommitsStartingAt(3, 5)

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI commit 07").IsSelected(),
				Contains("CI commit 06"),
				Contains("CI commit 05"),
				Contains("CI * commit 04"),
				Contains("CI commit 03"),
				Contains("CI commit 02"),
				Contains("CI commit 01"),
			).
			NavigateToLine(Contains("commit 02")).
			Press(keys.Commits.StartInteractiveRebase).
			Lines(
				Contains("pick").Contains("CI commit 07"),
				Contains("pick").Contains("CI commit 06"),
				Contains("pick").Contains("CI commit 05"),
				Contains("update-ref").Contains("branch1"),
				Contains("pick").Contains("CI commit 04"),
				Contains("pick").Contains("CI commit 03"),
				Contains("CI commit 02").IsSelected(),
				Contains("CI <-- YOU ARE HERE --- commit 01"),
			)
	},
})
