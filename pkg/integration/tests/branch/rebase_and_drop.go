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
			Contains("to keep").IsSelected(),
			Contains("to remove"),
			Contains("first change"),
			Contains("original"),
		)

		input.NextItem()
		input.Press(keys.Branches.RebaseBranch)

		input.InConfirm().
			Title(Equals("Rebasing")).
			Content(Contains("Are you sure you want to rebase 'first-change-branch' on top of 'second-change-branch'?")).
			Confirm()

		assert.View("information").Content(Contains("rebasing"))

		input.InConfirm().
			Title(Equals("Auto-merge failed")).
			Content(Contains("Conflicts!")).
			Confirm()

		assert.CurrentView().
			Name("files").
			SelectedLine(MatchesRegexp("UU.*file"))

		input.SwitchToCommitsView()
		assert.CurrentView().
			TopLines(
				MatchesRegexp(`pick.*to keep`).IsSelected(),
				MatchesRegexp(`pick.*to remove`),
				MatchesRegexp("YOU ARE HERE.*second-change-branch unrelated change"),
				MatchesRegexp("second change"),
				MatchesRegexp("original"),
			)

		input.NextItem()
		input.Press(keys.Universal.Remove)

		assert.CurrentView().
			TopLines(
				MatchesRegexp(`pick.*to keep`),
				MatchesRegexp(`drop.*to remove`).IsSelected(),
				MatchesRegexp("YOU ARE HERE.*second-change-branch unrelated change"),
				MatchesRegexp("second change"),
				MatchesRegexp("original"),
			)

		input.SwitchToFilesView()

		// not using Confirm() convenience method because I suspect we might change this
		// keybinding to something more bespoke
		input.Press(keys.Universal.Confirm)

		assert.CurrentView().Name("mergeConflicts")
		input.PrimaryAction()

		input.InConfirm().
			Title(Equals("continue")).
			Content(Contains("all merge conflicts resolved. Continue?")).
			Confirm()

		assert.View("information").Content(DoesNotContain("rebasing"))

		assert.View("commits").TopLines(
			Contains("to keep"),
			Contains("second-change-branch unrelated change").IsSelected(),
			Contains("second change"),
			Contains("original"),
		)
	},
})
