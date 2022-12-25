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
		input.SwitchToBranchesWindow()
		assert.CurrentViewName("localBranches")

		assert.CurrentViewLines(
			Contains("first-change-branch"),
			Contains("second-change-branch"),
			Contains("original-branch"),
		)

		input.NextItem()

		input.Enter()

		assert.CurrentViewName("subCommits")

		assert.CurrentViewTopLines(
			Contains("second-change-branch unrelated change"),
			Contains("second change"),
		)

		input.Press(keys.Commits.CherryPickCopy)
		assert.ViewContent("information", Contains("1 commit copied"))

		input.NextItem()
		input.Press(keys.Commits.CherryPickCopy)
		assert.ViewContent("information", Contains("2 commits copied"))

		input.SwitchToCommitsWindow()
		assert.CurrentViewName("commits")

		assert.CurrentViewTopLines(
			Contains("first change"),
		)

		input.Press(keys.Commits.PasteCommits)
		input.Alert(Equals("Cherry-Pick"), Contains("Are you sure you want to cherry-pick the copied commits onto this branch?"))

		input.AcceptConfirmation(Equals("Auto-merge failed"), Contains("Conflicts!"))

		assert.CurrentViewName("files")
		assert.SelectedLine(Contains("file"))

		// not using Confirm() convenience method because I suspect we might change this
		// keybinding to something more bespoke
		input.Press(keys.Universal.Confirm)

		assert.CurrentViewName("mergeConflicts")
		// picking 'Second change'
		input.NextItem()
		input.PrimaryAction()

		input.AcceptConfirmation(Equals("continue"), Contains("all merge conflicts resolved. Continue?"))

		assert.CurrentViewName("files")
		assert.WorkingTreeFileCount(0)

		input.SwitchToCommitsWindow()
		assert.CurrentViewName("commits")

		assert.CurrentViewTopLines(
			Contains("second-change-branch unrelated change"),
			Contains("second change"),
			Contains("first change"),
		)
		input.NextItem()
		// because we picked 'Second change' when resolving the conflict,
		// we now see this commit as having replaced First Change with Second Change,
		// as opposed to replacing 'Original' with 'Second change'
		assert.MainViewContent(Contains("-First Change"))
		assert.MainViewContent(Contains("+Second Change"))

		assert.ViewContent("information", Contains("2 commits copied"))
		input.Press(keys.Universal.Return)
		assert.ViewContent("information", NotContains("commits copied"))
	},
})
