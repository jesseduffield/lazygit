package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// note: this is required to simulate the clipboard during CI
func expectClipboard(t *TestDriver, matcher *TextMatcher) {
	defer t.Shell().DeleteFile("clipboard")

	t.FileSystem().FileContent("clipboard", matcher)
}

var CopyMenu = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The copy menu allows to copy name and diff of selected/all files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.CopyToClipboardCmd = "echo {{text}} > clipboard"
	},
	SetupRepo: func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Disabled item
		t.Views().Files().
			IsEmpty().
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("File name")).
					Tooltip(Equals("Disabled: No item selected")).
					Confirm().
					Tap(func() {
						t.ExpectToast(Equals("Disabled: No item selected"))
					}).
					Cancel()
			})

		t.Shell().
			CreateDir("dir").
			CreateFile("dir/1-unstaged_file", "unstaged content")

		// Empty content (new file)
		t.Views().Files().
			Press(keys.Universal.Refresh).
			Lines(
				Contains("dir").IsSelected(),
				Contains("unstaged_file"),
			).
			SelectNextItem().
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Diff of selected file")).
					Tooltip(Contains("Disabled: Nothing to copy")).
					Confirm().
					Tap(func() {
						t.ExpectToast(Equals("Disabled: Nothing to copy"))
					}).
					Cancel()
			}).
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Diff of all files")).
					Tooltip(Contains("Disabled: Nothing to copy")).
					Confirm().
					Tap(func() {
						t.ExpectToast(Equals("Disabled: Nothing to copy"))
					}).
					Cancel()
			})

		t.Shell().
			GitAdd("dir/1-unstaged_file").
			Commit("commit-unstaged").
			UpdateFile("dir/1-unstaged_file", "unstaged content (new)").
			CreateFileAndAdd("dir/2-staged_file", "staged content").
			Commit("commit-staged").
			UpdateFile("dir/2-staged_file", "staged content (new)").
			GitAdd("dir/2-staged_file")

		// Copy file name
		t.Views().Files().
			Press(keys.Universal.Refresh).
			Lines(
				Contains("dir"),
				Contains("unstaged_file").IsSelected(),
				Contains("staged_file"),
			).
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("File name")).
					Confirm()

				t.ExpectToast(Equals("File name copied to clipboard"))

				expectClipboard(t, Contains("unstaged_file"))
			})

		// Copy file path
		t.Views().Files().
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Path")).
					Confirm()

				t.ExpectToast(Equals("File path copied to clipboard"))

				expectClipboard(t, Contains("dir/1-unstaged_file"))
			})

		// Selected path diff on a single (unstaged) file
		t.Views().Files().
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Diff of selected file")).
					Tooltip(Equals("If there are staged items, this command considers only them. Otherwise, it considers all the unstaged ones.")).
					Confirm()

				t.ExpectToast(Equals("File diff copied to clipboard"))

				expectClipboard(t, Contains("+unstaged content (new)"))
			})

		// Selected path diff with staged and unstaged files
		t.Views().Files().
			SelectPreviousItem().
			Lines(
				Contains("dir").IsSelected(),
				Contains("unstaged_file"),
				Contains("staged_file"),
			).
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Diff of selected file")).
					Tooltip(Equals("If there are staged items, this command considers only them. Otherwise, it considers all the unstaged ones.")).
					Confirm()

				t.ExpectToast(Equals("File diff copied to clipboard"))

				expectClipboard(t, Contains("+staged content (new)"))
			})

		// All files diff with staged files
		t.Views().Files().
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Diff of all files")).
					Tooltip(Equals("If there are staged items, this command considers only them. Otherwise, it considers all the unstaged ones.")).
					Confirm()

				t.ExpectToast(Equals("All files diff copied to clipboard"))

				expectClipboard(t, Contains("+staged content (new)"))
			})

		// All files diff with no staged files
		t.Views().Files().
			SelectNextItem().
			SelectNextItem().
			Lines(
				Contains("dir"),
				Contains("unstaged_file"),
				Contains("staged_file").IsSelected(),
			).
			Press(keys.Universal.Select).
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Diff of all files")).
					Tooltip(Equals("If there are staged items, this command considers only them. Otherwise, it considers all the unstaged ones.")).
					Confirm()

				t.ExpectToast(Equals("All files diff copied to clipboard"))

				expectClipboard(t, Contains("+staged content (new)").Contains("+unstaged content (new)"))
			})
	},
})
