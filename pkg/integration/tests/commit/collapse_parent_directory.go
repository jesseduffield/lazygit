package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CollapseParentDirectory = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Collapsing the directory that the selected file of a commit is in",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("dir/sub/file-a", "original content\n")
		shell.CreateFileAndAdd("dir/file-b", "original content\n")
		shell.CreateFileAndAdd("compressed/chain/file-c", "original content\n")
		shell.Commit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ compressed/chain"),
				Equals("    A file-c"),
				Equals("  ▼ dir"),
				Equals("    A file-b"),
				Equals("    ▼ sub"),
				Equals("      A file-a"),
			)

		// Each press folds up one more level, leaving the cursor on the
		// directory that was just collapsed.
		t.Views().CommitFiles().
			NavigateToLine(Contains("file-a")).
			Press(keys.Files.CollapseParentDirectory).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ compressed/chain"),
				Equals("    A file-c"),
				Equals("  ▼ dir"),
				Equals("    A file-b"),
				Equals("    ▶ sub").IsSelected(),
			).
			Press(keys.Files.CollapseParentDirectory).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ compressed/chain"),
				Equals("    A file-c"),
				Equals("  ▶ dir").IsSelected(),
			).
			Press(keys.Files.CollapseParentDirectory).
			Lines(
				Equals("▶ /").IsSelected(),
			)

		// The root is at the top level, so there's nothing left to collapse.
		t.Views().CommitFiles().
			Press(keys.Files.CollapseParentDirectory).
			Lines(
				Equals("▶ /").IsSelected(),
			)

		// A chain of single-child directories is one node, and collapses as one.
		t.Views().CommitFiles().
			Press(keys.Files.ExpandAll).
			NavigateToLine(Contains("file-c")).
			Press(keys.Files.CollapseParentDirectory).
			Lines(
				Equals("▼ /"),
				Equals("  ▶ compressed/chain").IsSelected(),
				Equals("  ▼ dir"),
				Equals("    A file-b"),
				Equals("    ▼ sub"),
				Equals("      A file-a"),
			)
	},
})
