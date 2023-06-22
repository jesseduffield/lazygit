package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PickRescheduled = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Makes a pick during a rebase fail because it would overwrite an untracked file",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "1\n").Commit("one")
		shell.UpdateFileAndAdd("file2", "2\n").Commit("two")
		shell.UpdateFileAndAdd("file3", "3\n").Commit("three")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("three").IsSelected(),
				Contains("two"),
				Contains("one"),
			).
			NavigateToLine(Contains("one")).
			Press(keys.Universal.Edit).
			Lines(
				Contains("pick").Contains("three"),
				Contains("pick").Contains("two"),
				Contains("<-- YOU ARE HERE --- one").IsSelected(),
			).
			Tap(func() {
				t.Shell().CreateFile("file3", "other content\n")
				t.Common().ContinueRebase()
				t.ExpectPopup().Alert().Title(Equals("Error")).
					Content(Contains("The following untracked working tree files would be overwritten by merge").
						Contains("Please move or remove them before you merge.")).
					Confirm()
			}).
			Lines(
				Contains("pick").Contains("three"),
				Contains("<-- YOU ARE HERE --- two"),
				Contains("one"),
			)
	},
})
