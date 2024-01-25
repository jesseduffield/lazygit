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
				Contains("▼ dir1").IsSelected(),
				Contains("  ??").Contains("file-1a"),
				Contains("  ??").Contains("file-1b"),
				Contains("▼ dir2"),
				Contains("  ??").Contains("file-2a"),
				Contains("   M").Contains("file-2b"),
				Contains("▼ dir3"),
				Contains("  ??").Contains("file-3a"),
				Contains("   M").Contains("file-3b"),
				Contains("??").Contains("file-a"),
				Contains("??").Contains("file-b"),
			).
			NavigateToLine(Contains("file-1b")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file-2a")).
			Lines(
				Contains("▼ dir1"),
				Contains("  ??").Contains("file-1a"),
				Contains("  ??").Contains("file-1b").IsSelected(),
				Contains("▼ dir2").IsSelected(),
				Contains("  ??").Contains("file-2a").IsSelected(),
				Contains("   M").Contains("file-2b"),
				Contains("▼ dir3"),
				Contains("  ??").Contains("file-3a"),
				Contains("   M").Contains("file-3b"),
				Contains("??").Contains("file-a"),
				Contains("??").Contains("file-b"),
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
				Contains("▼ dir1"),
				Contains("  ??").Contains("file-1a"),
				Contains("▼ dir3").IsSelected(),
				Contains("  ??").Contains("file-3a"),
				Contains("   M").Contains("file-3b"),
				Contains("??").Contains("file-a"),
				Contains("??").Contains("file-b"),
			).
			// Verify you can discard collapsed directories in range select
			PressEnter().
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file-a")).
			Lines(
				Contains("▼ dir1"),
				Contains("  ??").Contains("file-1a"),
				Contains("▶ dir3").IsSelected(),
				Contains("??").Contains("file-a").IsSelected(),
				Contains("??").Contains("file-b"),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard all changes")).
					Confirm()
			}).
			Lines(
				Contains("▼ dir1"),
				Contains("  ??").Contains("file-1a"),
				Contains("??").Contains("file-b").IsSelected(),
			)
	},
})
