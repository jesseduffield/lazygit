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

		assert.MatchSelectedLine(Contains("first-change-branch"))
		input.NextItem()
		assert.MatchSelectedLine(Contains("second-change-branch"))
		input.PressKeys(keys.Branches.RebaseBranch)

		assert.InConfirm()
		assert.MatchCurrentViewContent(Contains("Are you sure you want to rebase 'first-change-branch' onto 'second-change-branch'?"))
		input.Confirm()

		assert.InConfirm()
		assert.MatchCurrentViewContent(Contains("Conflicts!"))
		input.Confirm()

		assert.CurrentViewName("files")
		assert.MatchSelectedLine(Contains("file"))

		// not using Confirm() convenience method because I suspect we might change this
		// keybinding to something more bespoke
		input.PressKeys(keys.Universal.Confirm)

		assert.CurrentViewName("mergeConflicts")
		input.PrimaryAction()

		assert.MatchViewContent("information", Contains("rebasing"))
		assert.InConfirm()
		assert.MatchCurrentViewContent(Contains("all merge conflicts resolved. Continue?"))
		input.Confirm()
		assert.MatchViewContent("information", NotContains("rebasing"))

		// this proves we actually have integrated the changes from second-change-branch
		assert.MatchViewContent("commits", Contains("second-change-branch unrelated change"))
	},
})
