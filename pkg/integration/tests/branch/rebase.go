package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var originalFileContent = `
This
Is
The
Original
File
`

var firstChangeFileContent = `
This
Is
The
First Change
File
`

var secondChangeFileContent = `
This
Is
The
Second Change
File
`

// prepares us for a rebase that has conflicts
var commonRebaseSetup = func(shell *Shell) {
	shell.
		NewBranch("original-branch").
		EmptyCommit("one").
		EmptyCommit("two").
		EmptyCommit("three").
		CreateFileAndAdd("file", originalFileContent).
		Commit("original").
		NewBranch("first-change-branch").
		UpdateFileAndAdd("file", firstChangeFileContent).
		Commit("first change").
		Checkout("original-branch").
		NewBranch("second-change-branch").
		UpdateFileAndAdd("file", secondChangeFileContent).
		Commit("second change").
		EmptyCommit("second-change-branch unrelated change").
		Checkout("first-change-branch")
}

var Rebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase onto another branch, deal with the conflicts.",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		commonRebaseSetup(shell)
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

		assert.InConfirm()
		assert.MatchCurrentViewContent(Contains("all merge conflicts resolved. Continue?"))
		input.Confirm()

		// this proves we actually have integrated the changes from second-change-branch
		assert.MatchViewContent("commits", Contains("second-change-branch unrelated change"))
	},
})
