package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardAllChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discard all changes of a file in the staging panel, then assert we land in the staging panel of the next file",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\n")
		shell.CreateFileAndAdd("file2", "1\n2\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nthree\nfour\n")
		shell.UpdateFile("file2", "1\n2\n3\n4\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("   M file1"),
				Equals("   M file2"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().Staging().
			IsFocused().
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(Contains("+three")).
			// discard the line
			Press(keys.Universal.Remove).
			Tap(func() {
				t.Common().ConfirmDiscardLines()
			}).
			SelectedLines(Contains("+four")).
			// discard the other line
			Press(keys.Universal.Remove).
			Tap(func() {
				t.Common().ConfirmDiscardLines()

				// because there are no more changes in file1 we switch to file2
				t.Views().Files().
					Lines(
						Equals(" M file2"),
					)
			}).
			// assert we are still in the staging panel, but now looking at the changes of the other file
			IsFocused().
			SelectedLines(Contains("+3"))
	},
})
