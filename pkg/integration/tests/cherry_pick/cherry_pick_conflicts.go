package cherry_pick

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var CherryPickConflicts = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cherry pick commits from the subcommits view, with conflicts",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.MergeConflictsSetup(shell)
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		input.SwitchToBranchesView()
		assert.CurrentView().Lines(
			Contains("first-change-branch"),
			Contains("second-change-branch"),
			Contains("original-branch"),
		)

		input.NextItem()

		input.Enter()

		assert.CurrentView().Name("subCommits").TopLines(
			Contains("second-change-branch unrelated change"),
			Contains("second change"),
		)

		input.Press(keys.Commits.CherryPickCopy)
		assert.View("information").Content(Contains("1 commit copied"))

		input.NextItem()
		input.Press(keys.Commits.CherryPickCopy)
		assert.View("information").Content(Contains("2 commits copied"))

		input.SwitchToCommitsView()

		assert.CurrentView().TopLines(
			Contains("first change"),
		)

		input.Press(keys.Commits.PasteCommits)
		input.Alert(Equals("Cherry-Pick"), Contains("Are you sure you want to cherry-pick the copied commits onto this branch?"))

		input.AcceptConfirmation(Equals("Auto-merge failed"), Contains("Conflicts!"))

		assert.CurrentView().Name("files")
		assert.CurrentView().SelectedLine(Contains("file"))

		// not using Confirm() convenience method because I suspect we might change this
		// keybinding to something more bespoke
		input.Press(keys.Universal.Confirm)

		assert.CurrentView().Name("mergeConflicts")
		// picking 'Second change'
		input.NextItem()
		input.PrimaryAction()

		input.AcceptConfirmation(Equals("continue"), Contains("all merge conflicts resolved. Continue?"))

		assert.CurrentView().Name("files")
		assert.WorkingTreeFileCount(0)

		input.SwitchToCommitsView()

		assert.CurrentView().TopLines(
			Contains("second-change-branch unrelated change"),
			Contains("second change"),
			Contains("first change"),
		)
		input.NextItem()
		// because we picked 'Second change' when resolving the conflict,
		// we now see this commit as having replaced First Change with Second Change,
		// as opposed to replacing 'Original' with 'Second change'
		assert.MainView().
			Content(Contains("-First Change")).
			Content(Contains("+Second Change"))

		assert.View("information").Content(Contains("2 commits copied"))
		input.Press(keys.Universal.Return)
		assert.View("information").Content(NotContains("commits copied"))
	},
})
