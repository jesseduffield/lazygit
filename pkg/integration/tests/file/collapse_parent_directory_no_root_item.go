package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CollapseParentDirectoryNoRootItem = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Collapsing the directory that the selected item is in, when the root item is disabled",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowRootItemInFileTree = false
	},
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
				Equals("▼ compressed/chain").IsSelected(),
				Equals("  ?? file-c"),
				Equals("▼ dir"),
				Equals("  ?? file-b"),
				Equals("  ▼ sub"),
				Equals("    ?? file-a"),
			)

		t.Views().Files().
			NavigateToLine(Contains("file-a")).
			Press(keys.Files.CollapseParentDirectory).
			Lines(
				Equals("▼ compressed/chain"),
				Equals("  ?? file-c"),
				Equals("▼ dir"),
				Equals("  ?? file-b"),
				Equals("  ▶ sub").IsSelected(),
			).
			Press(keys.Files.CollapseParentDirectory).
			Lines(
				Equals("▼ compressed/chain"),
				Equals("  ?? file-c"),
				Equals("▶ dir").IsSelected(),
			)

		// Without a root item, "dir" is already at the top level.
		t.Views().Files().
			Press(keys.Files.CollapseParentDirectory).
			Lines(
				Equals("▼ compressed/chain"),
				Equals("  ?? file-c"),
				Equals("▶ dir").IsSelected(),
			)

		t.Views().Files().
			NavigateToLine(Contains("file-c")).
			Press(keys.Files.CollapseParentDirectory).
			Lines(
				Equals("▶ compressed/chain").IsSelected(),
				Equals("▶ dir"),
			)
	},
})
