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
		shell.CreateFileAndAdd("test/file-tracked", "foo")

		shell.Commit("first commit")

		shell.CreateFile("test/file-untracked", "bar")
		shell.UpdateFile("test/file-tracked", "baz")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Lines(
				Equals("▼ test").IsSelected(),
				Equals("   M file-tracked"),
				Equals("  ?? file-untracked"),
			).
			Press(keys.Files.OpenStatusFilter).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Filtering")).
					Select(Contains("Show only tracked files")).
					Confirm()
			}).
			Lines(
				Equals("▼ test").IsSelected(),
				Equals("   M file-tracked"),
			).
			PressPrimaryAction().
			Press(keys.Files.OpenStatusFilter).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Filtering")).
					Select(Contains("No filter")).
					Confirm()
			}).
			Lines(
				Equals("▼ test").IsSelected(),
				Equals("  M  file-tracked"), // 'M' is now in the left column, so file is staged
				Equals("  ?? file-untracked"),
			)
	},
})
