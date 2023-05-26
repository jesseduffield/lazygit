package cherry_pick

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var CherryPickConflicts = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cherry pick commits from the subcommits view, with conflicts",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.MergeConflictsSetup(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("first-change-branch"),
				Contains("second-change-branch"),
				Contains("original-branch"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			TopLines(
				Contains("second-change-branch unrelated change"),
				Contains("second change"),
			).
			Press(keys.Commits.CherryPickCopy).
			Tap(func() {
				t.Views().Information().Content(Contains("1 commit copied"))
			}).
			SelectNextItem().
			Press(keys.Commits.CherryPickCopy)

		t.Views().Information().Content(Contains("2 commits copied"))

		t.Views().Commits().
			Focus().
			TopLines(
				Contains("first change"),
			).
			Press(keys.Commits.PasteCommits)

		t.ExpectPopup().Alert().
			Title(Equals("Cherry-pick")).
			Content(Contains("Are you sure you want to cherry-pick the copied commits onto this branch?")).
			Confirm()

		t.Common().AcknowledgeConflicts()

		t.Views().Files().
			IsFocused().
			SelectedLine(Contains("file")).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			// picking 'Second change'
			SelectNextItem().
			PressPrimaryAction()

		t.Common().ContinueOnConflictsResolved()

		t.Views().Files().IsEmpty()

		t.Views().Commits().
			Focus().
			TopLines(
				Contains("second-change-branch unrelated change"),
				Contains("second change"),
				Contains("first change"),
			).
			SelectNextItem().
			Tap(func() {
				// because we picked 'Second change' when resolving the conflict,
				// we now see this commit as having replaced First Change with Second Change,
				// as opposed to replacing 'Original' with 'Second change'
				t.Views().Main().
					Content(Contains("-First Change")).
					Content(Contains("+Second Change"))

				t.Views().Information().Content(Contains("2 commits copied"))
			}).
			PressEscape().
			Tap(func() {
				t.Views().Information().Content(DoesNotContain("commits copied"))
			})
	},
})
