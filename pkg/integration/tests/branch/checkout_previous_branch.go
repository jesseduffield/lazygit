package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CheckoutPreviousBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Checkout to the previous branch using the checkout previous branch functionality",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(3).
			NewBranch("previous-branch").
			EmptyCommit("previous commit").
			Checkout("master").
			EmptyCommit("master commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("previous-branch"),
			)

		// Press the checkout previous branch key (should checkout previous-branch)
		t.Views().Branches().
			Press(keys.Branches.CheckoutPreviousBranch).
			Lines(
				Contains("previous-branch").IsSelected(),
				Contains("master"),
			)

		// Verify we're on previous-branch
		t.Git().CurrentBranchName("previous-branch")

		// Press again to go back to master
		t.Views().Branches().
			Press(keys.Branches.CheckoutPreviousBranch).
			Lines(
				Contains("master").IsSelected(),
				Contains("previous-branch"),
			)

		// Verify we're back on master
		t.Git().CurrentBranchName("master")
	},
})
