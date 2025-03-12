package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardUnstagedRangeSelect = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discard unstaged changed in a range of files using range select",
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
		shell.CreateFileAndAdd("dir2/file-c", "")
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
				Equals("    A  file-c"),
				Equals("     M file-d"),
				Equals("  ?? file-e"),
				Equals("  ?? file-f"),
			).
			NavigateToLine(Contains("file-b")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file-c")).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir1"),
				Equals("    ?? file-a"),
				Equals("    ?? file-b").IsSelected(),
				Equals("  ▼ dir2").IsSelected(),
				Equals("    A  file-c").IsSelected(),
				Equals("     M file-d"),
				Equals("  ?? file-e"),
				Equals("  ?? file-f"),
			).
			// Discard
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard unstaged changes")).
					Confirm()
			}).
			// file-b is gone because it was selected and contained no staged changes.
			// file-c is still there because it contained no unstaged changes
			// file-d is gone because it was selected via dir2 and contained only unstaged changes
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir1"),
				Equals("    ?? file-a"),
				Equals("  ▼ dir2"),
				// Re-selecting file-c because it's where the selected line index
				// was before performing the action.
				Equals("    A  file-c").IsSelected(),
				Equals("  ?? file-e"),
				Equals("  ?? file-f"),
			)
	},
})
