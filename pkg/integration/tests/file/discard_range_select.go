package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardRangeSelect = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discard a range of files using range select",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("dir2/file-2b", "old content")
		shell.CreateFileAndAdd("dir3/file-3b", "old content")
		shell.Commit("first commit")
		shell.UpdateFile("dir2/file-2b", "new content")
		shell.UpdateFile("dir3/file-3b", "new content")

		shell.CreateFile("dir1/file-1a", "")
		shell.CreateFile("dir1/file-1b", "")
		shell.CreateFile("dir2/file-2a", "")
		shell.CreateFile("dir3/file-3a", "")
		shell.CreateFile("file-a", "")
		shell.CreateFile("file-b", "")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ dir1"),
				Equals("    ?? file-1a"),
				Equals("    ?? file-1b"),
				Equals("  ▼ dir2"),
				Equals("    ?? file-2a"),
				Equals("     M file-2b"),
				Equals("  ▼ dir3"),
				Equals("    ?? file-3a"),
				Equals("     M file-3b"),
				Equals("  ?? file-a"),
				Equals("  ?? file-b"),
			).
			NavigateToLine(Contains("file-1b")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file-2a")).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir1"),
				Equals("    ?? file-1a"),
				Equals("    ?? file-1b").IsSelected(),
				Equals("  ▼ dir2").IsSelected(),
				Equals("    ?? file-2a").IsSelected(),
				Equals("     M file-2b"),
				Equals("  ▼ dir3"),
				Equals("    ?? file-3a"),
				Equals("     M file-3b"),
				Equals("  ?? file-a"),
				Equals("  ?? file-b"),
			).
			// Discard
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard all changes")).
					Confirm()
			}).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir1"),
				Equals("    ?? file-1a"),
				Equals("  ▼ dir3").IsSelected(),
				Equals("    ?? file-3a"),
				Equals("     M file-3b"),
				Equals("  ?? file-a"),
				Equals("  ?? file-b"),
			).
			// Verify you can discard collapsed directories in range select
			PressEnter().
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file-a")).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir1"),
				Equals("    ?? file-1a"),
				Equals("  ▶ dir3").IsSelected(),
				Equals("  ?? file-a").IsSelected(),
				Equals("  ?? file-b"),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard all changes")).
					Confirm()
			}).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir1"),
				Equals("    ?? file-1a"),
				Equals("  ?? file-b").IsSelected(),
			)
	},
})
