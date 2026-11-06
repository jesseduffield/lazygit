package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var GlobMainBranches = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that glob patterns in mainBranches config match correctly",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.MainBranches = []string{"release/*"}
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("initial").
			NewBranch("release/1").
			EmptyCommit("release 1 commit").
			NewBranchFrom("feature", "release/1").
			EmptyCommit("feature 1").
			EmptyCommit("feature 2")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("feature 2").IsSelected(),
				Contains("feature 1"),
				Contains("release 1 commit"),
				Contains("initial"),
			).
			Press(keys.Commits.SelectCommitsOfCurrentBranch).
			Lines(
				Contains("feature 2").IsSelected(),
				Contains("feature 1").IsSelected(),
				Contains("release 1 commit"),
				Contains("initial"),
			)
	},
})
