package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SwapInRebaseWithConflict = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Via an edit-triggered rebase, swap two commits, causing a conflict. Then resolve the conflict and continue",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("myfile", "one")
		shell.Commit("commit one")
		shell.UpdateFileAndAdd("myfile", "two")
		shell.Commit("commit two")
		shell.UpdateFileAndAdd("myfile", "three")
		shell.Commit("commit three")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit three").IsSelected(),
				Contains("commit two"),
				Contains("commit one"),
			).
			NavigateToLine(Contains("commit one")).
			Press(keys.Universal.Edit).
			Lines(
				Contains("commit three"),
				Contains("commit two"),
				Contains("YOU ARE HERE").Contains("commit one").IsSelected(),
			).
			SelectPreviousItem().
			Press(keys.Commits.MoveUpCommit).
			Lines(
				Contains("commit two").IsSelected(),
				Contains("commit three"),
				Contains("YOU ARE HERE").Contains("commit one"),
			).
			Tap(func() {
				t.Common().ContinueRebase()
			})

		handleConflictsFromSwap(t)
	},
})
