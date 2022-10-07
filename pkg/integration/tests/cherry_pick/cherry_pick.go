package cherry_pick

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CherryPick = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cherry pick commits from the subcommits view, without conflicts",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("base").
			NewBranch("first-branch").
			NewBranch("second-branch").
			Checkout("first-branch").
			EmptyCommit("one").
			EmptyCommit("two").
			Checkout("second-branch").
			EmptyCommit("three").
			EmptyCommit("four").
			Checkout("first-branch")
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		input.SwitchToBranchesWindow()
		assert.CurrentViewName("localBranches")

		assert.MatchSelectedLine(Contains("first-branch"))
		input.NextItem()
		assert.MatchSelectedLine(Contains("second-branch"))

		input.Enter()

		assert.CurrentViewName("subCommits")
		assert.MatchSelectedLine(Contains("four"))
		input.PressKeys(keys.Commits.CherryPickCopy)
		assert.MatchViewContent("information", Contains("1 commit copied"))

		input.NextItem()
		assert.MatchSelectedLine(Contains("three"))
		input.PressKeys(keys.Commits.CherryPickCopy)
		assert.MatchViewContent("information", Contains("2 commits copied"))

		input.SwitchToCommitsWindow()
		assert.CurrentViewName("commits")

		assert.MatchSelectedLine(Contains("two"))
		input.PressKeys(keys.Commits.PasteCommits)
		assert.InAlert()
		assert.MatchCurrentViewContent(Contains("Are you sure you want to cherry-pick the copied commits onto this branch?"))

		input.Confirm()
		assert.CurrentViewName("commits")
		assert.MatchSelectedLine(Contains("four"))
		input.NextItem()
		assert.MatchSelectedLine(Contains("three"))
		input.NextItem()
		assert.MatchSelectedLine(Contains("two"))

		assert.MatchViewContent("information", Contains("2 commits copied"))
		input.PressKeys(keys.Universal.Return)
		assert.MatchViewContent("information", NotContains("commits copied"))
	},
})
