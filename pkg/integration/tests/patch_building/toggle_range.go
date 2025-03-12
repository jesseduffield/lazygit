package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ToggleRange = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Check multi select toggle logic",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("dir1")
		shell.CreateFileAndAdd("dir1/file1-a", "d2f1 first line\nsecond line\nthird line\n")
		shell.CreateFileAndAdd("dir1/file2-a", "d1f2 first line\n")
		shell.CreateFileAndAdd("dir1/file3-a", "d1f3 first line\n")

		shell.CreateDir("dir2")
		shell.CreateFileAndAdd("dir2/file1-b", "d2f1 first line\nsecond line\nthird line\n")
		shell.CreateFileAndAdd("dir2/file2-b", "d2f2 first line\n")
		shell.CreateFileAndAdd("dir2/file3-b", "d2f3 first line\nsecond line\n")

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
				Equals("  ▼ dir1"),
				Equals("    A file1-a"),
				Equals("    A file2-a"),
				Equals("    A file3-a"),
				Equals("  ▼ dir2"),
				Equals("    A file1-b"),
				Equals("    A file2-b"),
				Equals("    A file3-b"),
			).
			NavigateToLine(Contains("file1-a")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file3-a")).
			PressPrimaryAction().
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir1"),
				Equals("    ● file1-a").IsSelected(),
				Equals("    ● file2-a").IsSelected(),
				Equals("    ● file3-a").IsSelected(),
				Equals("  ▼ dir2"),
				Equals("    A file1-b"),
				Equals("    A file2-b"),
				Equals("    A file3-b"),
			).
			PressEscape().
			NavigateToLine(Contains("file3-b")).
			PressEnter()

		t.Views().Main().IsFocused().
			NavigateToLine(Contains("second line")).
			PressPrimaryAction().
			PressEscape()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir1"),
				Equals("    ● file1-a"),
				Equals("    ● file2-a"),
				Equals("    ● file3-a"),
				Equals("  ▼ dir2"),
				Equals("    A file1-b"),
				Equals("    A file2-b"),
				Equals("    ◐ file3-b").IsSelected(),
			).
			NavigateToLine(Contains("dir1")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("dir2")).
			PressPrimaryAction().
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir1").IsSelected(),
				Equals("    ● file1-a").IsSelected(),
				Equals("    ● file2-a").IsSelected(),
				Equals("    ● file3-a").IsSelected(),
				Equals("  ▼ dir2").IsSelected(),
				Equals("    ● file1-b"),
				Equals("    ● file2-b"),
				Equals("    ● file3-b"),
			).
			PressPrimaryAction().
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir1").IsSelected(),
				Equals("    A file1-a").IsSelected(),
				Equals("    A file2-a").IsSelected(),
				Equals("    A file3-a").IsSelected(),
				Equals("  ▼ dir2").IsSelected(),
				Equals("    A file1-b"),
				Equals("    A file2-b"),
				Equals("    A file3-b"),
			)
	},
})
