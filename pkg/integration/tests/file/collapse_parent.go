package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CollapseParent = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing backspace on a selected file jumps to its parent directory and collapses it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("dir")
		shell.CreateDir("dir/subdir")
		shell.CreateFile("dir/subdir/file-one", "original content\n")
		shell.CreateDir("dir2")
		shell.CreateFile("dir2/file-two", "original content\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ dir/subdir"),
				Equals("    ?? file-one"),
				Equals("  ▼ dir2"),
				Equals("    ?? file-two"),
			).
			// select "file-one" nested two levels deep
			NavigateToLine(Contains("file-one"))

		// backspace jumps to and collapses the immediate (compressed) parent "dir/subdir"
		t.Views().Files().
			Press(keys.Files.CollapseParent).
			Lines(
				Equals("▼ /"),
				Equals("  ▶ dir/subdir").IsSelected(),
				Equals("  ▼ dir2"),
				Equals("    ?? file-two"),
			)

		// backspace again jumps up to and collapses the root
		t.Views().Files().
			Press(keys.Files.CollapseParent).
			Lines(
				Equals("▶ /").IsSelected(),
			)

		// backspace on an already-collapsed root is a no-op
		t.Views().Files().
			Press(keys.Files.CollapseParent).
			Lines(
				Equals("▶ /").IsSelected(),
			)
	},
})
