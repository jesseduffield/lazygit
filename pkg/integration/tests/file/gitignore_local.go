package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var GitignoreLocal = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that ignoring a file uses the local .gitignore if one exists in the same directory",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		// Create root .gitignore
		shell.CreateFile(".gitignore", "")
		// Create subdirectory with its own .gitignore
		shell.CreateDir("subdir_with_gitignore")
		shell.CreateFile("subdir_with_gitignore/.gitignore", "")
		shell.CreateFile("subdir_with_gitignore/file_to_ignore", "")
		// Create subdirectory without .gitignore
		shell.CreateDir("subdir_without_gitignore")
		shell.CreateFile("subdir_without_gitignore/another_file", "")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ subdir_with_gitignore"),
				Equals("    ?? .gitignore"),
				Equals("    ?? file_to_ignore"),
				Equals("  ▼ subdir_without_gitignore"),
				Equals("    ?? another_file"),
				Equals("  ?? .gitignore"),
			).
			// Navigate to subdir_with_gitignore/file_to_ignore
			NavigateToLine(Contains("file_to_ignore")).
			Press(keys.Files.IgnoreFile).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Ignore or exclude file")).Select(Contains("Add to .gitignore")).Confirm()
				// Should be added to the local .gitignore with just the basename
				t.FileSystem().FileContent("subdir_with_gitignore/.gitignore", Equals("file_to_ignore\n"))
				// Root .gitignore should remain empty
				t.FileSystem().FileContent(".gitignore", Equals(""))
			}).
			// Navigate to subdir_without_gitignore/another_file
			NavigateToLine(Contains("another_file")).
			Press(keys.Files.IgnoreFile).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Ignore or exclude file")).Select(Contains("Add to .gitignore")).Confirm()
				// Should be added to root .gitignore with full path since no local .gitignore exists
				t.FileSystem().FileContent(".gitignore", Equals("subdir_without_gitignore/another_file\n"))
				// Local .gitignore should not have been created
				t.FileSystem().PathNotPresent("subdir_without_gitignore/.gitignore")
			})
	},
})
