package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var InteractiveRebaseOfCopiedBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Check that interactively rebasing a branch that is a copy of another branch doesn't affect the original branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig: func(config *config.AppConfig) {
		config.AppState.GitLogShowGraph = "never"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			NewBranch("branch1").
			CreateNCommits(3).
			NewBranch("branch2")

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI * commit 03"),
				Contains("CI commit 02"),
				Contains("CI commit 01"),
			).
			NavigateToLine(Contains("commit 01")).
			Press(keys.Universal.Edit).
			Lines(
				// No update-ref todo for branch1 here, even though command-line git would have added it
				Contains("pick").Contains("CI commit 03"),
				Contains("pick").Contains("CI commit 02"),
				Contains("CI <-- YOU ARE HERE --- commit 01"),
			)
	},
})
