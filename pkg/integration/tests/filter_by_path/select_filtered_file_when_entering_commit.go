package filter_by_path

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectFilteredFileWhenEnteringCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filter commits by file path, then enter a commit and ensure the file is selected",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "")
		shell.CreateFileAndAdd("dir/file2", "")
		shell.Commit("add files")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.GlobalPress(keys.Universal.FilteringMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Filtering")).
			Select(Contains("Enter path to filter by")).
			Confirm()

		t.ExpectPopup().Prompt().
			Title(Equals("Enter path:")).
			Type("dir/file2").
			Confirm()

		t.Views().Commits().
			Focus().
			Lines(
				Contains("add files").IsSelected(),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir"),
				Equals("    A file2").IsSelected(),
				Equals("  A file1"),
			)
	},
})
