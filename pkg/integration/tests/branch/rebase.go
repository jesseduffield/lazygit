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
		input.SwitchToBranchesView()

		assert.View("localBranches").Lines(
			Contains("first-change-branch"),
			Contains("second-change-branch"),
			Contains("original-branch"),
		)

		assert.View("commits").TopLines(
			Contains("first change"),
			Contains("original"),
		)

		input.NextItem()
		input.Press(keys.Branches.RebaseBranch)

		input.Confirmation().
			Title(Equals("Rebasing")).
			Content(Contains("Are you sure you want to rebase 'first-change-branch' on top of 'second-change-branch'?")).
			Confirm()
		input.Confirmation().
			Title(Equals("Auto-merge failed")).
			Content(Contains("Conflicts!")).
			Confirm()

		assert.CurrentView().Name("files").SelectedLine(Contains("file"))

		// not using Confirm() convenience method because I suspect we might change this
		// keybinding to something more bespoke
		input.Press(keys.Universal.Confirm)

		assert.CurrentView().Name("mergeConflicts")
		input.PrimaryAction()

		assert.View("information").Content(Contains("rebasing"))

		input.Confirmation().
			Title(Equals("continue")).
			Content(Contains("all merge conflicts resolved. Continue?")).
			Confirm()

		assert.View("information").Content(DoesNotContain("rebasing"))

		assert.View("commits").TopLines(
			Contains("second-change-branch unrelated change"),
			Contains("second change"),
			Contains("original"),
		)
	},
})
