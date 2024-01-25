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
				Contains("▼ dir1").IsSelected(),
				Contains("  ??").Contains("file-a"),
				Contains("  ??").Contains("file-b"),
				Contains("▼ dir2"),
				Contains("  ??").Contains("file-c"),
				Contains("   M").Contains("file-d"),
				Contains("??").Contains("file-e"),
				Contains("??").Contains("file-f"),
			).
			NavigateToLine(Contains("file-b")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file-c")).
			// Stage
			PressPrimaryAction().
			Lines(
				Contains("▼ dir1"),
				Contains("  ??").Contains("file-a"),
				Contains("  A ").Contains("file-b").IsSelected(),
				Contains("▼ dir2").IsSelected(),
				Contains("  A ").Contains("file-c").IsSelected(),
				// Staged because dir2 was part of the selection when he hit space
				Contains("  M ").Contains("file-d"),
				Contains("??").Contains("file-e"),
				Contains("??").Contains("file-f"),
			).
			// Unstage; back to everything being unstaged
			PressPrimaryAction().
			Lines(
				Contains("▼ dir1"),
				Contains("  ??").Contains("file-a"),
				Contains("  ??").Contains("file-b").IsSelected(),
				Contains("▼ dir2").IsSelected(),
				Contains("  ??").Contains("file-c").IsSelected(),
				Contains("   M").Contains("file-d"),
				Contains("??").Contains("file-e"),
				Contains("??").Contains("file-f"),
			).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("dir2")).
			// Verify that collapsed directories can be included in the range.
			// Collapse the directory
			PressEnter().
			Lines(
				Contains("▼ dir1"),
				Contains("  ??").Contains("file-a"),
				Contains("  ??").Contains("file-b"),
				Contains("▶ dir2").IsSelected(),
				Contains("??").Contains("file-e"),
				Contains("??").Contains("file-f"),
			).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file-e")).
			// Stage
			PressPrimaryAction().
			Lines(
				Contains("▼ dir1"),
				Contains("  ??").Contains("file-a"),
				Contains("  ??").Contains("file-b"),
				Contains("▶ dir2").IsSelected(),
				Contains("A ").Contains("file-e").IsSelected(),
				Contains("??").Contains("file-f"),
			).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("dir2")).
			// Expand the directory again to verify it's been staged
			PressEnter().
			Lines(
				Contains("▼ dir1"),
				Contains("  ??").Contains("file-a"),
				Contains("  ??").Contains("file-b"),
				Contains("▼ dir2").IsSelected(),
				Contains("  A ").Contains("file-c"),
				Contains("  M ").Contains("file-d"),
				Contains("A ").Contains("file-e"),
				Contains("??").Contains("file-f"),
			)
	},
})
