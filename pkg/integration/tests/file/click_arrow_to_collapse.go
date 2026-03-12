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
		// Initial state: all expanded
		// Row 0: "▼ /"          arrow at column 0
		// Row 1: "  ▼ dir"      arrow at column 2
		// Row 2: "    ?? file-one"
		// Row 3: "  ▼ dir2"     arrow at column 2
		// Row 4: "    ?? file-two"
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

		// Click the arrow again to expand it
		t.Views().Files().
			Click(2, 1).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir").IsSelected(),
				Equals("    ?? file-one"),
				Equals("  ▼ dir2"),
				Equals("    ?? file-two"),
			)

		// Click the arrow on the root "/" (row 0, column 0) to collapse everything
		t.Views().Files().
			Click(0, 0).
			Lines(
				Equals("▶ /").IsSelected(),
			)

		// Click again to expand
		t.Views().Files().
			Click(0, 0).
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ dir"),
				Equals("    ?? file-one"),
				Equals("  ▼ dir2"),
				Equals("    ?? file-two"),
			)
	},
})
