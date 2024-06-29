package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageDeletedRangeSelect = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage a range of deleted files using range select",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("dir1/file-a", "")
		shell.CreateFileAndAdd("dir2/file-b", "")
		shell.CreateFileAndAdd("file-c", "")
		shell.CreateFileAndAdd("file-d", "")
		shell.Commit("first commit")

		shell.DeleteFile("dir1/file-a")
		shell.DeleteFile("dir2/file-b")
		shell.DeleteFile("file-c")
		shell.DeleteFile("file-d")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("▼ dir1").IsSelected(),
				Contains("   D").Contains("file-a"),
				Contains("▼ dir2"),
				Contains("   D").Contains("file-b"),
				Contains(" D").Contains("file-c"),
				Contains(" D").Contains("file-d"),
			).
			NavigateToLine(Contains("file-b")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file-c")).
			// Stage a deleted file and nested file
			PressPrimaryAction().
			Lines(
				Contains("▼ dir1"),
				Contains("   D").Contains("file-a"),
				Contains("▼ dir2"),
				Contains("  D ").Contains("file-b").IsSelected(),
				Contains("D ").Contains("file-c").IsSelected(),
				Contains(" D").Contains("file-d"),
			).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file-a")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file-d")).
			// Stage the entire selection of files while some deleted files are already staged
			PressPrimaryAction().
			Lines(
				Contains("▼ dir1"),
				Contains("  D ").Contains("file-a").IsSelected(),
				Contains("▼ dir2").IsSelected(),
				Contains("  D ").Contains("file-b").IsSelected(),
				Contains("D ").Contains("file-c").IsSelected(),
				Contains("D ").Contains("file-d").IsSelected(),
			).
			// Unstage; back to everything being unstaged
			PressPrimaryAction().
			Lines(
				Contains("▼ dir1"),
				Contains("   D").Contains("file-a").IsSelected(),
				Contains("▼ dir2").IsSelected(),
				Contains("   D").Contains("file-b").IsSelected(),
				Contains(" D").Contains("file-c").IsSelected(),
				Contains(" D").Contains("file-d").IsSelected(),
			)
	},
})
