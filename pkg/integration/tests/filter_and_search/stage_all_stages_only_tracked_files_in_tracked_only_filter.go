package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageAllStagesOnlyTrackedFilesInTrackedOnlyFilter = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging all files in tracked only view should stage only tracked files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file-tracked", "foo")

		shell.Commit("first commit")

		shell.CreateFile("file-untracked", "bar")
		shell.UpdateFile("file-tracked", "baz")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Lines(
				Equals("▼ /").IsSelected(),
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
				Equals(" M file-tracked"),
			).
			Press(keys.Files.ToggleStagedAll).
			Press(keys.Files.OpenStatusFilter).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Filtering")).
					Select(Contains("No filter")).
					Confirm()
			}).
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  M  file-tracked"), // 'M' is now in the left column, so file is staged
				Equals("  ?? file-untracked"),
			)
	},
})
