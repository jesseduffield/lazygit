package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StagingFolderStagesOnlyTrackedFilesInTrackedOnlyFilter = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging entire folder in tracked only view, should stage only tracked files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("test")
		shell.Chdir("test")
		shell.CreateFileAndAdd("file-tracked", "foo")

		shell.Commit("first commit")

		shell.CreateFile("file-untracked", "bar")
		shell.UpdateFile("file-tracked", "baz")

		shell.Chdir("..")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Lines(
				Contains(`test`).IsSelected(),
				Contains(`M file-tracked`),
				Contains(`?? file-untracked`),
			).
			Press(keys.Files.OpenStatusFilter).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Filtering")).
					Select(Contains("Show only tracked files")).
					Confirm()
			}).
			Lines(
				Contains(`test`).IsSelected(),
				Contains(`file-tracked`),
			).
			Press(keys.Universal.Select).
			Press(keys.Files.OpenStatusFilter).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Filtering")).
					Select(Contains("No filter")).
					Confirm()
			}).
			Lines(
				Contains(`test`).IsSelected(),
				Contains(`M  file-tracked`), // double space means it's staged
				Contains(`?? file-untracked`),
			)
	},
})
