package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Suggestions = components.NewIntegrationTest(components.NewIntegrationTestArgs{
	Description:  "Checking out a branch with name suggestions",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *components.Shell) {
		shell.
			EmptyCommit("my commit message").
			NewBranch("new-branch").
			NewBranch("new-branch-2").
			NewBranch("new-branch-3").
			NewBranch("branch-to-checkout").
			NewBranch("other-new-branch-2").
			NewBranch("other-new-branch-3")
	},
	Run: func(shell *components.Shell, input *components.Input, assert *components.Assert, keys config.KeybindingConfig) {
		input.SwitchToBranchesWindow()
		assert.CurrentViewName("localBranches")

		input.PressKeys(keys.Branches.CheckoutBranchByName)
		assert.CurrentViewName("confirmation")

		input.Type("branch-to")

		input.PressKeys(keys.Universal.TogglePanel)
		assert.CurrentViewName("suggestions")

		// we expect the first suggestion to be the branch we want because it most
		// closely matches what we typed in
		input.Confirm()

		assert.CurrentBranchName("branch-to-checkout")
	},
})
