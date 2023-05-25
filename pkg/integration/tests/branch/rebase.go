package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var Rebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase onto another branch, deal with the conflicts.",
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
			Title(Equals("Rebase 'first-change-branch' onto 'second-change-branch'")).
			Select(Contains("Simple rebase")).
			Confirm()

		t.Common().AcknowledgeConflicts()

		t.Views().Files().
			IsFocused().
			SelectedLine(Contains("file")).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Rebasing"))

		t.Common().ContinueOnConflictsResolved()

		t.Views().Information().Content(DoesNotContain("Rebasing"))

		t.Views().Commits().TopLines(
			Contains("second-change-branch unrelated change"),
			Contains("second change"),
			Contains("original"),
		)
	},
})
