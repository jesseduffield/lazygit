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
	Run: func(shell *Shell, input *Input, keys config.KeybindingConfig) {
		input.Views().Branches().
			Focus().
			Lines(
				Contains("first-change-branch"),
				Contains("second-change-branch"),
				Contains("original-branch"),
			)

		input.Views().Commits().
			TopLines(
				Contains("to keep").IsSelected(),
				Contains("to remove"),
				Contains("first change"),
				Contains("original"),
			).
			SelectNextItem().
			Press(keys.Branches.RebaseBranch)

		input.ExpectConfirmation().
			Title(Equals("Rebasing")).
			Content(Contains("Are you sure you want to rebase 'first-change-branch' on top of 'second-change-branch'?")).
			Confirm()

		input.Views().Information().Content(Contains("rebasing"))

		input.ExpectConfirmation().
			Title(Equals("Auto-merge failed")).
			Content(Contains("Conflicts!")).
			Confirm()

		input.Views().Files().IsFocused().
			SelectedLine(MatchesRegexp("UU.*file"))

		input.Views().Commits().
			Focus().
			TopLines(
				MatchesRegexp(`pick.*to keep`).IsSelected(),
				MatchesRegexp(`pick.*to remove`),
				MatchesRegexp("YOU ARE HERE.*second-change-branch unrelated change"),
				MatchesRegexp("second change"),
				MatchesRegexp("original"),
			).
			SelectNextItem().
			Press(keys.Universal.Remove).
			TopLines(
				MatchesRegexp(`pick.*to keep`),
				MatchesRegexp(`drop.*to remove`).IsSelected(),
				MatchesRegexp("YOU ARE HERE.*second-change-branch unrelated change"),
				MatchesRegexp("second change"),
				MatchesRegexp("original"),
			)

		input.Views().Files().
			Focus().
			PressEnter()

		input.Views().MergeConflicts().
			IsFocused().
			PressPrimaryAction()

		input.ExpectConfirmation().
			Title(Equals("continue")).
			Content(Contains("all merge conflicts resolved. Continue?")).
			Confirm()

		input.Views().Information().Content(DoesNotContain("rebasing"))

		input.Views().Commits().TopLines(
			Contains("to keep"),
			Contains("second-change-branch unrelated change").IsSelected(),
			Contains("second change"),
			Contains("original"),
		)
	},
})
