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
	Run: func(shell *Shell, input *Input, keys config.KeybindingConfig) {
		input.Views().Branches().
			Focus().
			Lines(
				Contains("first-change-branch"),
				Contains("second-change-branch"),
				Contains("original-branch"),
			).
			SelectNextItem().
			PressEnter()

		input.Views().SubCommits().
			IsFocused().
			TopLines(
				Contains("second-change-branch unrelated change"),
				Contains("second change"),
			).
			Press(keys.Commits.CherryPickCopy)

		input.Views().Information().Content(Contains("1 commit copied"))

		input.Views().SubCommits().
			SelectNextItem().
			Press(keys.Commits.CherryPickCopy)

		input.Views().Information().Content(Contains("2 commits copied"))

		input.Views().Commits().
			Focus().
			TopLines(
				Contains("first change"),
			).
			Press(keys.Commits.PasteCommits)

		input.ExpectAlert().Title(Equals("Cherry-Pick")).Content(Contains("Are you sure you want to cherry-pick the copied commits onto this branch?")).Confirm()

		input.ExpectConfirmation().
			Title(Equals("Auto-merge failed")).
			Content(Contains("Conflicts!")).
			Confirm()

		input.Views().Files().
			IsFocused().
			SelectedLine(Contains("file")).
			PressEnter()

		input.Views().MergeConflicts().
			IsFocused().
			// picking 'Second change'
			SelectNextItem().
			PressPrimaryAction()

		input.ExpectConfirmation().
			Title(Equals("continue")).
			Content(Contains("all merge conflicts resolved. Continue?")).
			Confirm()

		input.Model().WorkingTreeFileCount(0)

		input.Views().Commits().
			Focus().
			TopLines(
				Contains("second-change-branch unrelated change"),
				Contains("second change"),
				Contains("first change"),
			).
			SelectNextItem()

		// because we picked 'Second change' when resolving the conflict,
		// we now see this commit as having replaced First Change with Second Change,
		// as opposed to replacing 'Original' with 'Second change'
		input.Views().Main().
			Content(Contains("-First Change")).
			Content(Contains("+Second Change"))

		input.Views().Information().Content(Contains("2 commits copied"))

		input.Views().Commits().
			PressEscape()

		input.Views().Information().Content(DoesNotContain("commits copied"))
	},
})
