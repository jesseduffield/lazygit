package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageRangeSelect = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage/unstage a range of files using range select",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("dir2/file-d", "old content")
		shell.Commit("first commit")
		shell.UpdateFile("dir2/file-d", "new content")

		shell.CreateFile("dir1/file-a", "")
		shell.CreateFile("dir1/file-b", "")
		shell.CreateFile("dir2/file-c", "")
		shell.CreateFile("file-e", "")
		shell.CreateFile("file-f", "")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ dir1"),
				Equals("    ?? file-a"),
				Equals("    ?? file-b"),
				Equals("  ▼ dir2"),
				Equals("    ?? file-c"),
				Equals("     M file-d"),
				Equals("  ?? file-e"),
				Equals("  ?? file-f"),
			).
			NavigateToLine(Contains("file-b")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file-c")).
			// Stage
			PressPrimaryAction().
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir1"),
				Equals("    ?? file-a"),
				Equals("    A  file-b").IsSelected(),
				Equals("  ▼ dir2").IsSelected(),
				Equals("    A  file-c").IsSelected(),
				// Staged because dir2 was part of the selection when he hit space
				Equals("    M  file-d"),
				Equals("  ?? file-e"),
				Equals("  ?? file-f"),
			).
			// Unstage; back to everything being unstaged
			PressPrimaryAction().
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir1"),
				Equals("    ?? file-a"),
				Equals("    ?? file-b").IsSelected(),
				Equals("  ▼ dir2").IsSelected(),
				Equals("    ?? file-c").IsSelected(),
				Equals("     M file-d"),
				Equals("  ?? file-e"),
				Equals("  ?? file-f"),
			).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("dir2")).
			// Verify that collapsed directories can be included in the range.
			// Collapse the directory
			PressEnter().
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir1"),
				Equals("    ?? file-a"),
				Equals("    ?? file-b"),
				Equals("  ▶ dir2").IsSelected(),
				Equals("  ?? file-e"),
				Equals("  ?? file-f"),
			).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file-e")).
			// Stage
			PressPrimaryAction().
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir1"),
				Equals("    ?? file-a"),
				Equals("    ?? file-b"),
				Equals("  ▶ dir2").IsSelected(),
				Equals("  A  file-e").IsSelected(),
				Equals("  ?? file-f"),
			).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("dir2")).
			// Expand the directory again to verify it's been staged
			PressEnter().
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir1"),
				Equals("    ?? file-a"),
				Equals("    ?? file-b"),
				Equals("  ▼ dir2").IsSelected(),
				Equals("    A  file-c"),
				Equals("    M  file-d"),
				Equals("  A  file-e"),
				Equals("  ?? file-f"),
			)
	},
})
