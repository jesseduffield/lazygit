package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RebaseOntoBaseBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase the current branch onto its base branch",
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
		t.Views().Commits().Lines(
			Contains("feature 2"),
			Contains("feature 1"),
			Contains("master 2"),
			Contains("master 1"),
		)

		t.Views().Branches().
			Focus().
			Lines(
				Contains("feature ↓1").IsSelected(),
				Contains("master"),
			).
			Press(keys.Branches.RebaseBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Rebase 'feature'")).
			Select(Contains("Rebase onto base branch (master)")).
			Confirm()

		t.Views().Commits().Lines(
			Contains("feature 2"),
			Contains("feature 1"),
			Contains("master 3"),
			Contains("master 2"),
			Contains("master 1"),
		)
	},
})
