package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var RebaseConflictsFixBuildErrors = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase onto another branch, deal with the conflicts. While continue prompt is showing, fix build errors; get another prompt when continuing.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.MergeConflictsSetup(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().TopLines(
			Contains("first change"),
			Contains("original"),
		)

		t.Views().Branches().
			Focus().
			Lines(
				Contains("first-change-branch"),
				Contains("second-change-branch"),
				Contains("original-branch"),
			).
			SelectNextItem().
			Press(keys.Branches.RebaseBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Rebase 'first-change-branch'")).
			Select(Contains("Simple rebase")).
			Confirm()

		t.Common().AcknowledgeConflicts()

		t.Views().Files().
			IsFocused().
			SelectedLine(Contains("file")).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			SelectNextItem().
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Rebasing"))

		popup := t.ExpectPopup().Confirmation().
			Title(Equals("Continue")).
			Content(Contains("All merge conflicts resolved. Continue?"))

		// While the popup is showing, fix some build errors
		t.Shell().UpdateFile("file", "make it compile again")

		// Continue
		popup.Confirm()

		t.ExpectPopup().Confirmation().
			Title(Equals("Continue")).
			Content(Contains("Files have been modified since conflicts were resolved. Auto-stage them and continue?")).
			Confirm()

		t.Views().Information().Content(DoesNotContain("Rebasing"))

		t.Views().Commits().TopLines(
			Contains("first change"),
			Contains("second-change-branch unrelated change"),
			Contains("second change"),
			Contains("original"),
		)
	},
})
