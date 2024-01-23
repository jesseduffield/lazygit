package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EditTheConflCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Swap two commits, causing a conflict; then try to interact with the 'confl' commit, which results in an error.",
	ExtraCmdArgs: []string{},
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
			Press(keys.Commits.MoveDownCommit).
			Tap(func() {
				t.Common().AcknowledgeConflicts()
			}).
			Focus().
			Lines(
				Contains("pick").Contains("commit two"),
				Contains("conflict").Contains("<-- YOU ARE HERE --- commit three"),
				Contains("commit one"),
			).
			NavigateToLine(Contains("<-- YOU ARE HERE --- commit three")).
			Press(keys.Commits.RenameCommit)

		t.ExpectToast(Contains("Disabled: Rewording commits while interactively rebasing is not currently supported"))
	},
})
