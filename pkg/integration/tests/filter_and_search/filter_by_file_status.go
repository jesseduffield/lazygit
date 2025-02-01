package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterByFileStatus = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filtering to show untracked files in repo that hides them by default",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		// need to set untracked files to not be displayed in git config
		shell.SetConfig("status.showUntrackedFiles", "no")

		shell.CreateFileAndAdd("file-tracked", "foo")

		shell.Commit("first commit")

		shell.CreateFile("file-untracked", "bar")
		shell.UpdateFile("file-tracked", "baz")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Lines(
				Contains(`file-tracked`).IsSelected(),
			).
			Press(keys.Files.OpenStatusFilter).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Filtering")).
					Select(Contains("Show only untracked files")).
					Confirm()
			}).
			Lines(
				Contains(`file-untracked`).IsSelected(),
			).
			Press(keys.Files.OpenStatusFilter).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Filtering")).
					Select(Contains("Show only tracked files")).
					Confirm()
			}).
			Lines(
				Contains(`file-tracked`).IsSelected(),
			).
			Press(keys.Files.OpenStatusFilter).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Filtering")).
					Select(Contains("No filter")).
					Confirm()
			}).
			Lines(
				Contains(`file-tracked`).IsSelected(),
			)
	},
})
