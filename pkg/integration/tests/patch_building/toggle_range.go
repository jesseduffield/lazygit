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
				Contains("▼ dir1").IsSelected(),
				Contains("  A").Contains("file1-a"),
				Contains("  A").Contains("file2-a"),
				Contains("  A").Contains("file3-a"),
				Contains("▼ dir2"),
				Contains("  A").Contains("file1-b"),
				Contains("  A").Contains("file2-b"),
				Contains("  A").Contains("file3-b"),
			).
			NavigateToLine(Contains("file1-a")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file3-a")).
			PressPrimaryAction().
			Lines(
				Contains("▼ dir1"),
				Contains("  ●").Contains("file1-a").IsSelected(),
				Contains("  ●").Contains("file2-a").IsSelected(),
				Contains("  ●").Contains("file3-a").IsSelected(),
				Contains("▼ dir2"),
				Contains("  A").Contains("file1-b"),
				Contains("  A").Contains("file2-b"),
				Contains("  A").Contains("file3-b"),
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
				Contains("▼ dir1"),
				Contains("  ●").Contains("file1-a"),
				Contains("  ●").Contains("file2-a"),
				Contains("  ●").Contains("file3-a"),
				Contains("▼ dir2"),
				Contains("  A").Contains("file1-b"),
				Contains("  A").Contains("file2-b"),
				Contains("  ◐").Contains("file3-b").IsSelected(),
			).
			NavigateToLine(Contains("dir1")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("dir2")).
			PressPrimaryAction().
			Lines(
				Contains("▼ dir1").IsSelected(),
				Contains("  ●").Contains("file1-a").IsSelected(),
				Contains("  ●").Contains("file2-a").IsSelected(),
				Contains("  ●").Contains("file3-a").IsSelected(),
				Contains("▼ dir2").IsSelected(),
				Contains("  ●").Contains("file1-b"),
				Contains("  ●").Contains("file2-b"),
				Contains("  ●").Contains("file3-b"),
			).
			PressPrimaryAction().
			Lines(
				Contains("▼ dir1").IsSelected(),
				Contains("  A").Contains("file1-a").IsSelected(),
				Contains("  A").Contains("file2-a").IsSelected(),
				Contains("  A").Contains("file3-a").IsSelected(),
				Contains("▼ dir2").IsSelected(),
				Contains("  A").Contains("file1-b"),
				Contains("  A").Contains("file2-b"),
				Contains("  A").Contains("file3-b"),
			)
	},
})
