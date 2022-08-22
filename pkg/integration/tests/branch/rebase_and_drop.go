package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RebaseAndDrop = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase onto another branch, deal with the conflicts. Also mark a commit to be dropped before continuing.",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		commonRebaseSetup(shell)
		// addin a couple additional commits so that we can drop one
		shell.EmptyCommit("to drop")
		shell.EmptyCommit("to keep")
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

		input.SwitchToCommitsWindow()
		input.NextItem()
		input.PressKeys(keys.Universal.Remove)
		assert.MatchSelectedLine(Contains("to drop"))
		assert.MatchSelectedLine(Contains("drop"))

		input.SwitchToFilesWindow()

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
		assert.MatchViewContent("commits", Contains("to keep"))
		assert.MatchViewContent("commits", NotContains("to drop"))
	},
})
