package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ClickArrowToCollapse = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Click the arrow on a directory to collapse/expand it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("dir")
		shell.CreateFile("dir/file-one", "original content\n")
		shell.CreateDir("dir2")
		shell.CreateFile("dir2/file-two", "original content\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ dir"),
				Equals("    ?? file-one"),
				Equals("  ▼ dir2"),
				Equals("    ?? file-two"),
			)

		// Click the arrow on "dir" (row 1, column 2) to collapse it
		t.Views().Files().
			Click(2, 1).
			Lines(
				Equals("▼ /"),
				Equals("  ▶ dir").IsSelected(),
				Equals("  ▼ dir2"),
				Equals("    ?? file-two"),
			)

		// Click one to the right of the arrow on "dir2" (row 2, column 3) to collapse it
		// Arrow + space after should register a collapse toggle
		t.Views().Files().
			Click(3, 2).
			Lines(
				Equals("▼ /"),
				Equals("  ▶ dir"),
				Equals("  ▶ dir2").IsSelected(),
			)

		// Click one to the left of the arrow on "dir2" (row 2, column 1)
		// Space before arrow should not register a collapse toggle
		t.Views().Files().
			Click(1, 2).
			Lines(
				Equals("▼ /"),
				Equals("  ▶ dir"),
				Equals("  ▶ dir2").IsSelected(),
			)

		// Clicking on the file/directory name "dir" should change selected but not toggle collapse
		t.Views().Files().
			Click(5, 1).
			Lines(
				Equals("▼ /"),
				Equals("  ▶ dir").IsSelected(),
				Equals("  ▶ dir2"),
			)

		// Click the arrow again to expand it
		t.Views().Files().
			Click(2, 1).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir").IsSelected(),
				Equals("    ?? file-one"),
				Equals("  ▶ dir2"),
			)

		// Click the arrow on the root "/" (row 0, column 0) to collapse everything
		t.Views().Files().
			Click(0, 0).
			Lines(
				Equals("▶ /").IsSelected(),
			)
	},
})
