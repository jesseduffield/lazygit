package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FastForwardRewrittenBranchWithLocalCommits = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Try to fast-forward a branch that has a commit of its own on top of a rewritten upstream branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		createBranchRewrittenOnTheRemote(shell)

		shell.Checkout("feature")
		shell.EmptyCommit("mine")
		shell.Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("feature ↓2↑3"),
			).
			SelectNextItem().
			Press(keys.Branches.FastForward)

		t.ExpectPopup().Alert().Title(Equals("Error")).
			Content(Contains("Cannot fast-forward 'feature' because it has commits")).
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("master"),
				Contains("feature ↓2↑3").IsSelected(),
			)
	},
})
