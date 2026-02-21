package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterFilesStageAll = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Toggle all staging with a filter only stages visible files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("dir1")
		shell.CreateFile("dir1/apple-grape", "apple-grape content\n")
		shell.CreateFile("dir1/apple-orange", "apple-orange content\n")
		shell.CreateFile("dir1/grape-orange", "grape-orange content\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Lines(
				Equals("▼ dir1").IsSelected(),
				Equals("  ?? apple-grape"),
				Equals("  ?? apple-orange"),
				Equals("  ?? grape-orange"),
			).
			// Filter to show only "apple" files
			FilterOrSearch("apple").
			Lines(
				// first item is always selected after filtering
				Equals("▼ dir1").IsSelected(),
				Equals("  ?? apple-grape"),
				Equals("  ?? apple-orange"),
			).
			// Stage all visible files
			Press(keys.Files.ToggleStagedAll).
			// Clear the filter and verify only apple files are staged
			PressEscape()

		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ dir1").IsSelected(),
				Equals("  A  apple-grape"),
				Equals("  A  apple-orange"),
				Equals("  ?? grape-orange"),
			)
	},
})
