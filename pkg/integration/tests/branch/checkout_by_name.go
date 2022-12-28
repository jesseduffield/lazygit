package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CheckoutByName = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Try to checkout branch by name. Verify that it also works on the branch with the special name @.",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(3).
			NewBranch("@").
			Checkout("master").
			EmptyCommit("blah")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("@"),
			).
			SelectNextItem().
			Press(keys.Branches.CheckoutBranchByName).
			Tap(func() {
				t.ExpectPopup().Prompt().Title(Equals("Branch name:")).Type("new-branch").Confirm()

				t.ExpectPopup().Alert().Title(Equals("Branch not found")).Content(Equals("Branch not found. Create a new branch named new-branch?")).Confirm()
			}).
			Lines(
				MatchesRegexp(`\*.*new-branch`).IsSelected(),
				Contains("master"),
				Contains("@"),
			)
	},
})
