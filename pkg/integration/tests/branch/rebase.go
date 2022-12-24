package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var Rebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase onto another branch, deal with the conflicts.",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.MergeConflictsSetup(shell)
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		input.SwitchToBranchesWindow()
		assert.CurrentViewName("localBranches")

		assert.SelectedLine(Contains("first-change-branch"))
		input.NextItem()
		assert.SelectedLine(Contains("second-change-branch"))
		input.Press(keys.Branches.RebaseBranch)

		input.AcceptConfirmation(Equals("Rebasing"), Contains("Are you sure you want to rebase 'first-change-branch' on top of 'second-change-branch'?"))

		input.AcceptConfirmation(Equals("Auto-merge failed"), Contains("Conflicts!"))

		assert.CurrentViewName("files")
		assert.SelectedLine(Contains("file"))

		// not using Confirm() convenience method because I suspect we might change this
		// keybinding to something more bespoke
		input.Press(keys.Universal.Confirm)

		assert.CurrentViewName("mergeConflicts")
		input.PrimaryAction()

		assert.ViewContent("information", Contains("rebasing"))

		input.AcceptConfirmation(Equals("continue"), Contains("all merge conflicts resolved. Continue?"))

		assert.ViewContent("information", NotContains("rebasing"))

		// this proves we actually have integrated the changes from second-change-branch
		assert.ViewContent("commits", Contains("second-change-branch unrelated change"))
	},
})
