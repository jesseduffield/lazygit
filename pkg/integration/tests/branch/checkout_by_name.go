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
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		input.SwitchToBranchesView()

		assert.CurrentView().Lines(
			Contains("master"),
			Contains("@"),
		)
		input.NextItem()

		input.Press(keys.Branches.CheckoutBranchByName)

		input.Prompt().Title(Equals("Branch name:")).Type("new-branch").Confirm()

		input.Alert(Equals("Branch not found"), Equals("Branch not found. Create a new branch named new-branch?"))

		assert.CurrentView().Name("localBranches").
			Lines(
				MatchesRegexp(`\*.*new-branch`),
				Contains("master"),
				Contains("@"),
			).
			SelectedLine(Contains("new-branch"))
	},
})
