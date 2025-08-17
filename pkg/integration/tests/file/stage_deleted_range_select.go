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
		shell.CreateFileAndAdd("file-a", "")
		shell.CreateFileAndAdd("file-b", "")
		shell.Commit("first commit")

		shell.DeleteFile("file-a")
		shell.DeleteFile("file-b")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("   D file-a"),
				Equals("   D file-b"),
			).
			SelectNextItem().
			// Stage a single deleted file
			PressPrimaryAction().
			Lines(
				Equals("▼ /"),
				Equals("  D  file-a").IsSelected(),
				Equals("   D file-b"),
			).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file-b")).
			// Stage both files while a deleted file is already staged
			PressPrimaryAction().
			Lines(
				Equals("▼ /"),
				Equals("  D  file-a").IsSelected(),
				Equals("  D  file-b").IsSelected(),
			).
			// Unstage; back to everything being unstaged
			PressPrimaryAction().
			Lines(
				Equals("▼ /"),
				Equals("   D file-a").IsSelected(),
				Equals("   D file-b").IsSelected(),
			)
	},
})
