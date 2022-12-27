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
		input.SwitchToBranchesView()

		assert.Views().Current().Lines(
			Contains("first-branch"),
			Contains("second-branch"),
			Contains("master"),
		)

		input.NextItem()

		input.Enter()

		assert.Views().Current().Name("subCommits").Lines(
			Contains("four"),
			Contains("three"),
			Contains("base"),
		)

		// copy commits 'four' and 'three'
		input.Press(keys.Commits.CherryPickCopy)
		assert.Views().ByName("information").Content(Contains("1 commit copied"))
		input.NextItem()
		input.Press(keys.Commits.CherryPickCopy)
		assert.Views().ByName("information").Content(Contains("2 commits copied"))

		input.SwitchToCommitsView()

		assert.Views().Current().Lines(
			Contains("two"),
			Contains("one"),
			Contains("base"),
		)

		input.Press(keys.Commits.PasteCommits)
		input.Alert().
			Title(Equals("Cherry-Pick")).
			Content(Contains("Are you sure you want to cherry-pick the copied commits onto this branch?")).
			Confirm()

		assert.Views().Current().Name("commits").Lines(
			Contains("four"),
			Contains("three"),
			Contains("two"),
			Contains("one"),
			Contains("base"),
		)

		assert.Views().ByName("information").Content(Contains("2 commits copied"))
		input.Press(keys.Universal.Return)
		assert.Views().ByName("information").Content(DoesNotContain("commits copied"))
	},
})
