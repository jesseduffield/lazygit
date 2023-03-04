package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EditNonTodoCommitDuringRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Tries to edit a non-todo commit while already rebasing, resulting in an error message",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(2)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 02").IsSelected(),
				Contains("commit 01"),
			).
			Press(keys.Universal.Edit).
			Lines(
				Contains("<-- YOU ARE HERE --- commit 02"),
				Contains("commit 01"),
			).
			NavigateToLine(Contains("commit 01")).
			Press(keys.Universal.Edit)

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("Can't perform this action during a rebase")).
			Confirm()
	},
})
