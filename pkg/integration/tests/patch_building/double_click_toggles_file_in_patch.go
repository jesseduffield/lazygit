package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DoubleClickTogglesFileInPatch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Double-clicking a file or directory of a commit toggles whether it is in the custom patch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("dir/file1", "file1 content\n")
		shell.CreateFileAndAdd("dir/file2", "file2 content\n")
		shell.CreateFileAndAdd("file3", "file3 content\n")
		shell.Commit("first commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("first commit").IsSelected(),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ dir"),
				Equals("    A file1"),
				Equals("    A file2"),
				Equals("  A file3"),
			).
			DoubleClick(5, 4).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir"),
				Equals("    A file1"),
				Equals("    A file2"),
				Equals("  ● file3").IsSelected(),
			)

		t.Views().Information().Content(Contains("Building patch"))

		t.Views().CommitFiles().
			DoubleClick(5, 1).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir").IsSelected(),
				Equals("    ● file1"),
				Equals("    ● file2"),
				Equals("  ● file3"),
			).
			DoubleClick(5, 4).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir"),
				Equals("    ● file1"),
				Equals("    ● file2"),
				Equals("  A file3").IsSelected(),
			).
			DoubleClick(5, 1).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir").IsSelected(),
				Equals("    A file1"),
				Equals("    A file2"),
				Equals("  A file3"),
			)

		t.Views().Information().Content(DoesNotContain("Building patch"))
	},
})
