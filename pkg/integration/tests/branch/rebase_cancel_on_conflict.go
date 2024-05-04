package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var RebaseCancelOnConflict = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase onto another branch, cancel when there are conflicts.",
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

		t.ExpectPopup().Menu().
			Title(Equals("Conflicts!")).
			Select(Contains("Abort the rebase")).
			Cancel()

		t.Views().Branches().
			IsFocused()

		t.Views().Files().
			Lines(
				Contains("UU file"),
			)
	},
})
