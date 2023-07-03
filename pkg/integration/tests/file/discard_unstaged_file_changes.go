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
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("MM").Contains("file-one").IsSelected(),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("file-one")).
					Select(Contains("Discard unstaged changes")).
					Confirm()
			}).
			Lines(
				Contains("M ").Contains("file-one").IsSelected(),
			)

		t.FileSystem().FileContent("file-one", Equals("original content\nnew content\n"))
	},
})
