package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Delete = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Try to delete the checked out branch first (to no avail), and then delete another branch.",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("blah").
			NewBranch("branch-one").
			NewBranch("branch-two")
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		input.SwitchToBranchesWindow()
		assert.CurrentViewName("localBranches")

		assert.MatchSelectedLine(Contains("branch-two"))
		input.PressKeys(keys.Universal.Remove)
		assert.InAlert()
		assert.MatchCurrentViewContent(Contains("You cannot delete the checked out branch!"))

		input.Confirm()

		input.NextItem()
		assert.MatchSelectedLine(Contains("branch-one"))
		input.PressKeys(keys.Universal.Remove)
		assert.InConfirm()
		assert.MatchCurrentViewContent(Contains("Are you sure you want to delete the branch 'branch-one'?"))
		input.Confirm()
		assert.CurrentViewName("localBranches")
		assert.MatchSelectedLine(Contains("master"))
		assert.MatchCurrentViewContent(NotContains("branch-one"))
	},
})
