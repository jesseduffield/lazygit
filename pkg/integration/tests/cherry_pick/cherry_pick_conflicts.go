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

		assert.MatchSelectedLine(Contains("first-change-branch"))
		input.NextItem()
		assert.MatchSelectedLine(Contains("second-change-branch"))

		input.Enter()

		assert.CurrentViewName("subCommits")
		assert.MatchSelectedLine(Contains("second-change-branch unrelated change"))
		input.PressKeys(keys.Commits.CherryPickCopy)
		assert.MatchViewContent("information", Contains("1 commit copied"))

		input.NextItem()
		assert.MatchSelectedLine(Contains("second change"))
		input.PressKeys(keys.Commits.CherryPickCopy)
		assert.MatchViewContent("information", Contains("2 commits copied"))

		input.SwitchToCommitsWindow()
		assert.CurrentViewName("commits")

		assert.MatchSelectedLine(Contains("first change"))
		input.PressKeys(keys.Commits.PasteCommits)
		assert.InAlert()
		assert.MatchCurrentViewContent(Contains("Are you sure you want to cherry-pick the copied commits onto this branch?"))

		input.Confirm()

		assert.MatchCurrentViewContent(Contains("Conflicts!"))
		input.Confirm()

		assert.CurrentViewName("files")
		assert.MatchSelectedLine(Contains("file"))

		// not using Confirm() convenience method because I suspect we might change this
		// keybinding to something more bespoke
		input.PressKeys(keys.Universal.Confirm)

		assert.CurrentViewName("mergeConflicts")
		// picking 'Second change'
		input.NextItem()
		input.PrimaryAction()

		assert.InConfirm()
		assert.MatchCurrentViewContent(Contains("all merge conflicts resolved. Continue?"))
		input.Confirm()

		assert.CurrentViewName("files")
		assert.WorkingTreeFileCount(0)

		input.SwitchToCommitsWindow()
		assert.CurrentViewName("commits")

		assert.MatchSelectedLine(Contains("second-change-branch unrelated change"))
		input.NextItem()
		assert.MatchSelectedLine(Contains("second change"))
		// because we picked 'Second change' when resolving the conflict,
		// we now see this commit as having replaced First Change with Second Change,
		// as opposed to replacing 'Original' with 'Second change'
		assert.MatchMainViewContent(Contains("-First Change"))
		assert.MatchMainViewContent(Contains("+Second Change"))
		input.NextItem()
		assert.MatchSelectedLine(Contains("first change"))

		assert.MatchViewContent("information", Contains("2 commits copied"))
		input.PressKeys(keys.Universal.Return)
		assert.MatchViewContent("information", NotContains("commits copied"))
	},
})
