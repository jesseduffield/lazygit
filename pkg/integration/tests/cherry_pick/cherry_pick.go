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
	Run: func(shell *Shell, input *Input, keys config.KeybindingConfig) {
		input.Views().Branches().
			Focus().
			Lines(
				Contains("first-branch"),
				Contains("second-branch"),
				Contains("master"),
			).
			SelectNextItem().
			PressEnter()

		input.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("four").IsSelected(),
				Contains("three"),
				Contains("base"),
			).
			// copy commits 'four' and 'three'
			Press(keys.Commits.CherryPickCopy)

		input.Views().Information().Content(Contains("1 commit copied"))

		input.Views().SubCommits().
			SelectNextItem().
			Press(keys.Commits.CherryPickCopy)

		input.Views().Information().Content(Contains("2 commits copied"))

		input.Views().Commits().
			Focus().
			Lines(
				Contains("two").IsSelected(),
				Contains("one"),
				Contains("base"),
			).
			Press(keys.Commits.PasteCommits)

		input.ExpectAlert().
			Title(Equals("Cherry-Pick")).
			Content(Contains("Are you sure you want to cherry-pick the copied commits onto this branch?")).
			Confirm()

		input.Views().Commits().
			IsFocused().
			Lines(
				Contains("four"),
				Contains("three"),
				Contains("two"),
				Contains("one"),
				Contains("base"),
			)

			// we need to manually exit out of cherrry pick mode
		input.Views().Information().Content(Contains("2 commits copied"))

		input.Views().Commits().
			PressEscape()

		input.Views().Information().Content(DoesNotContain("commits copied"))
	},
})
