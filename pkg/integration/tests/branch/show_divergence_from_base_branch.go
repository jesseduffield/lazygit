package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowDivergenceFromBaseBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show divergence from base branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowDivergenceFromBaseBranch = "arrowAndNumber"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("master 1").
			EmptyCommit("master 2").
			EmptyCommit("master 3").
			NewBranchFrom("feature", "master^").
			EmptyCommit("feature 1").
			EmptyCommit("feature 2")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("feature ↓1").IsSelected(),
				Contains("master"),
			).
			Press(keys.Branches.SetUpstream)

		t.ExpectPopup().Menu().Title(Contains("Upstream")).
			Select(Contains("View divergence from base branch (master)")).Confirm()

		t.Views().SubCommits().
			IsFocused().
			Title(Contains("Commits (feature <-> master)")).
			Lines(
				DoesNotContainAnyOf("↓", "↑").Contains("--- Remote ---"),
				Contains("↓").Contains("master 3"),
				DoesNotContainAnyOf("↓", "↑").Contains("--- Local ---"),
				Contains("↑").Contains("feature 2"),
				Contains("↑").Contains("feature 1"),
			)
	},
})
