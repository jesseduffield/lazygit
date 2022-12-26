package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var RebaseAndDrop = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase onto another branch, deal with the conflicts. Also mark a commit to be dropped before continuing.",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.MergeConflictsSetup(shell)
		// addin a couple additional commits so that we can drop one
		shell.EmptyCommit("to remove")
		shell.EmptyCommit("to keep")
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		input.SwitchToBranchesView()

		assert.CurrentView().Lines(
			Contains("first-change-branch"),
			Contains("second-change-branch"),
			Contains("original-branch"),
		)

		assert.View("commits").TopLines(
			Contains("to keep"),
			Contains("to remove"),
			Contains("first change"),
			Contains("original"),
		)

		input.NextItem()
		input.Press(keys.Branches.RebaseBranch)

		input.AcceptConfirmation(Equals("Rebasing"), Contains("Are you sure you want to rebase 'first-change-branch' on top of 'second-change-branch'?"))

		assert.View("information").Content(Contains("rebasing"))

		input.AcceptConfirmation(Equals("Auto-merge failed"), Contains("Conflicts!"))

		assert.CurrentView().
			Name("files").
			SelectedLine(Contains("file"))

		input.SwitchToCommitsView()
		assert.CurrentView().
			TopLines(
				MatchesRegexp(`pick.*to keep`),
				MatchesRegexp(`pick.*to remove`),
				MatchesRegexp("YOU ARE HERE.*second-change-branch unrelated change"),
				MatchesRegexp("second change"),
				MatchesRegexp("original"),
			).
			SelectedLineIdx(0)

		input.NextItem()
		input.Press(keys.Universal.Remove)
		assert.CurrentView().SelectedLine(MatchesRegexp(`drop.*to remove`))

		input.SwitchToFilesView()

		// not using Confirm() convenience method because I suspect we might change this
		// keybinding to something more bespoke
		input.Press(keys.Universal.Confirm)

		assert.CurrentView().Name("mergeConflicts")
		input.PrimaryAction()

		input.AcceptConfirmation(Equals("continue"), Contains("all merge conflicts resolved. Continue?"))

		assert.View("information").Content(NotContains("rebasing"))

		assert.View("commits").TopLines(
			Contains("to keep"),
			Contains("second-change-branch unrelated change"),
			Contains("second change"),
			Contains("original"),
		)
	},
})
