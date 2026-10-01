package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CollapseParentDirectory = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Collapsing the directory that the selected item is in",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("dir")
		shell.CreateDir("dir/sub")
		shell.CreateFile("dir/sub/file-a", "original content\n")
		shell.CreateFile("dir/file-b", "original content\n")
		shell.CreateDir("compressed")
		shell.CreateDir("compressed/chain")
		shell.CreateFile("compressed/chain/file-c", "original content\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ compressed/chain"),
				Equals("    ?? file-c"),
				Equals("  ▼ dir"),
				Equals("    ?? file-b"),
				Equals("    ▼ sub"),
				Equals("      ?? file-a"),
			)

		// Each press folds up one more level, leaving the cursor on the
		// directory that was just collapsed.
		t.Views().Files().
			NavigateToLine(Contains("file-a")).
			Press(keys.Files.CollapseParentDirectory).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ compressed/chain"),
				Equals("    ?? file-c"),
				Equals("  ▼ dir"),
				Equals("    ?? file-b"),
				Equals("    ▶ sub").IsSelected(),
			).
			Press(keys.Files.CollapseParentDirectory).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ compressed/chain"),
				Equals("    ?? file-c"),
				Equals("  ▶ dir").IsSelected(),
			).
			Press(keys.Files.CollapseParentDirectory).
			Lines(
				Equals("▶ /").IsSelected(),
			)

		// The root is at the top level, so there's nothing left to collapse.
		t.Views().Files().
			Press(keys.Files.CollapseParentDirectory).
			Lines(
				Equals("▶ /").IsSelected(),
			)

		// A chain of single-child directories is one node, and collapses as one.
		t.Views().Files().
			Press(keys.Files.ExpandAll).
			NavigateToLine(Contains("file-c")).
			Press(keys.Files.CollapseParentDirectory).
			Lines(
				Equals("▼ /"),
				Equals("  ▶ compressed/chain").IsSelected(),
				Equals("  ▼ dir"),
				Equals("    ?? file-b"),
				Equals("    ▼ sub"),
				Equals("      ?? file-a"),
			)
	},
})
