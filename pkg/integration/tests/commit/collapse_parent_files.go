package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CollapseParentFiles = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing backspace in the commit files panel jumps to the parent directory and collapses it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("dir1/subd1/subfile0", "file0\n")
		shell.CreateFileAndAdd("dir2/d2_file1", "d2f1 content\n")
		shell.Commit("add files in two dirs")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("add files in two dirs").IsSelected(),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ dir1/subd1"),
				Equals("    A subfile0"),
				Equals("  ▼ dir2"),
				Equals("    A d2_file1"),
			).
			NavigateToLine(Contains("subfile0"))

		// backspace jumps to and collapses the immediate (compressed) parent "dir1/subd1"
		t.Views().CommitFiles().
			Press(keys.Files.CollapseParent).
			Lines(
				Equals("▼ /"),
				Equals("  ▶ dir1/subd1").IsSelected(),
				Equals("  ▼ dir2"),
				Equals("    A d2_file1"),
			)

		// backspace again jumps up to and collapses the root
		t.Views().CommitFiles().
			Press(keys.Files.CollapseParent).
			Lines(
				Equals("▶ /").IsSelected(),
			)
	},
})
