package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var GitignoreRangeSelect = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Ignore and exclude multiple files at once via range select",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile(".gitignore", "")
		shell.CreateFile("toIgnore1", "")
		shell.CreateFile("toIgnore2", "")
		shell.CreateFile("toIgnore3", "")
		shell.CreateFile("toExclude1", "")
		shell.CreateFile("toExclude2", "")
		shell.CreateFile("toExclude3", "")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ?? .gitignore"),
				Equals("  ?? toExclude1"),
				Equals("  ?? toExclude2"),
				Equals("  ?? toExclude3"),
				Equals("  ?? toIgnore1"),
				Equals("  ?? toIgnore2"),
				Equals("  ?? toIgnore3"),
			).
			// Select range from toIgnore1 to toIgnore3
			NavigateToLine(Contains("toIgnore1")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("toIgnore3")).
			Press(keys.Files.IgnoreFile).
			// Ignore all selected files
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Ignore or exclude file")).
					Select(Contains("Add to .gitignore")).
					Confirm()

				t.FileSystem().FileContent(".gitignore", Equals("/toIgnore1\n/toIgnore2\n/toIgnore3\n"))
			}).
			// Dismiss the range select mode for the next set of steps
			Press(keys.Universal.ToggleRangeSelect).
			// Select range from toExclude1 to toExclude3
			NavigateToLine(Contains("toExclude1")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("toExclude3")).
			Press(keys.Files.IgnoreFile).
			// Exclude all selected files
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Ignore or exclude file")).
					Select(Contains("Add to .git/info/exclude")).
					Confirm()

				t.FileSystem().FileContent(".git/info/exclude", Contains("/toExclude1\n/toExclude2\n/toExclude3\n"))
			})
	},
})
