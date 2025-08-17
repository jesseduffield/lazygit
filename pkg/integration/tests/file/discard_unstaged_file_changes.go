package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardUnstagedFileChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding unstaged changes in a file",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file-one", "original content\n")

		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file-one", "original content\nnew content\n")
		shell.UpdateFile("file-one", "original content\nnew content\neven newer content\n")

		shell.CreateFileAndAdd("file-two", "original content\n")
		shell.UpdateFile("file-two", "original content\nnew content\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  MM file-one"),
				Equals("  AM file-two"),
			).
			SelectNextItem().
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard unstaged changes")).
					Confirm()
			}).
			Lines(
				Equals("▼ /"),
				Equals("  M  file-one").IsSelected(),
				Equals("  AM file-two"),
			).
			SelectNextItem().
			Lines(
				Equals("▼ /"),
				Equals("  M  file-one"),
				Equals("  AM file-two").IsSelected(),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard unstaged changes")).
					Confirm()
			}).
			Lines(
				Equals("▼ /"),
				Equals("  M  file-one"),
				Equals("  A  file-two").IsSelected(),
			)

		t.FileSystem().FileContent("file-one", Equals("original content\nnew content\n"))
		t.FileSystem().FileContent("file-two", Equals("original content\n"))
	},
})
