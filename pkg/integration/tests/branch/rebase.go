package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var Rebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase onto another branch, deal with the conflicts.",
	ExtraCmdArgs: "",
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

		t.ExpectPopup().Confirmation().
			Title(Equals("Rebasing")).
			Content(Contains("Are you sure you want to rebase 'first-change-branch' on top of 'second-change-branch'?")).
			Confirm()

		t.ExpectPopup().Confirmation().
			Title(Equals("Auto-merge failed")).
			Content(Contains("Conflicts!")).
			Confirm()

		t.Views().Files().
			IsFocused().
			SelectedLine(Contains("file")).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			PressPrimaryAction()

		t.Views().Information().Content(Contains("rebasing"))

		t.ExpectPopup().Confirmation().
			Title(Equals("continue")).
			Content(Contains("all merge conflicts resolved. Continue?")).
			Confirm()

		t.Views().Information().Content(DoesNotContain("rebasing"))

		t.Views().Commits().TopLines(
			Contains("second-change-branch unrelated change"),
			Contains("second change"),
			Contains("original"),
		)
	},
})
