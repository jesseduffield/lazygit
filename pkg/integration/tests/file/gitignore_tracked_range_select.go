package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var GitignoreTrackedRangeSelect = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Range-select across a tracked directory and its children (parent+child case)",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("dir-tracked/file-a", "x")
		shell.CreateFileAndAdd("dir-tracked/file-b", "x")
		shell.CreateFileAndAdd("tracked1", "x")
		shell.Commit("initial")
		shell.UpdateFile("dir-tracked/file-a", "y")
		shell.UpdateFile("dir-tracked/file-b", "y")
		shell.UpdateFile("tracked1", "y")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ dir-tracked"),
				Equals("     M file-a"),
				Equals("     M file-b"),
				Equals("   M tracked1"),
			).
			NavigateToLine(Contains("dir-tracked")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("tracked1")).
			Press(keys.Files.IgnoreFile).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Ignore or exclude file")).
					Select(Contains("Add to .gitignore")).
					Confirm()

				t.ExpectPopup().Confirmation().
					Title(Equals("Ignore tracked file")).
					Content(Contains("tracked file")).
					Confirm()

				t.FileSystem().FileContent(".gitignore", Equals("/dir-tracked\n/tracked1\n"))
			})
	},
})
