package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DoubleClickTogglesStaging = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Double-clicking a file or directory toggles whether it is staged",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("dir/file1", "")
		shell.CreateFile("dir/file2", "")
		shell.CreateFile("file3", "")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ dir"),
				Equals("    ?? file1"),
				Equals("    ?? file2"),
				Equals("  ?? file3"),
			).
			DoubleClick(6, 4).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir"),
				Equals("    ?? file1"),
				Equals("    ?? file2"),
				Equals("  A  file3").IsSelected(),
			).
			DoubleClick(5, 1).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir").IsSelected(),
				Equals("    A  file1"),
				Equals("    A  file2"),
				Equals("  A  file3"),
			).
			DoubleClick(6, 4).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir"),
				Equals("    A  file1"),
				Equals("    A  file2"),
				Equals("  ?? file3").IsSelected(),
			).
			DoubleClick(5, 1).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir").IsSelected(),
				Equals("    ?? file1"),
				Equals("    ?? file2"),
				Equals("  ?? file3"),
			)
	},
})
