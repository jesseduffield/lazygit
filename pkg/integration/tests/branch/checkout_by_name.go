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
		input.SwitchToBranchesWindow()
		assert.CurrentViewName("localBranches")

		assert.MatchSelectedLine(Contains("master"))
		input.NextItem()
		assert.MatchSelectedLine(Contains("@"))
		input.PressKeys(keys.Branches.CheckoutBranchByName)
		assert.InPrompt()
		assert.MatchCurrentViewTitle(Equals("Branch name:"))
		input.Type("new-branch")
		input.Confirm()
		assert.InAlert()
		assert.MatchCurrentViewContent(Equals("Branch not found. Create a new branch named new-branch?"))
		input.Confirm()

		assert.CurrentViewName("localBranches")
		assert.MatchSelectedLine(Contains("new-branch"))
	},
})
