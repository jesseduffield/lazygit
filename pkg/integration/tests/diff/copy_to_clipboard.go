package diff

import (
	"os"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// note: this is required to simulate the clipboard during CI
func expectClipboard(t *TestDriver, matcher *TextMatcher) {
	defer t.Shell().DeleteFile("clipboard")

	t.FileSystem().FileContent("clipboard", matcher)
}

var CopyToClipboard = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The copy menu allows to copy name and diff of selected/all files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > clipboard"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("dir")
		shell.CreateFileAndAdd("dir/file1", "1st line\n")
		shell.Commit("1")
		shell.UpdateFileAndAdd("dir/file1", "1st line\n2nd line\n")
		shell.CreateFileAndAdd("dir/file2", "file2\n")
		shell.Commit("2")
		shell.UpdateFileAndAdd("dir/file1", "1st line\n2nd line\n3rd line\n")
		shell.Commit("3")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("3").IsSelected(),
				Contains("2"),
				Contains("1"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("dir").IsSelected(),
				Contains("file1"),
				Contains("file2"),
			).
			NavigateToLine(Contains("file1")).
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("File name")).
					Confirm().
					Tap(func() {
						t.ExpectToast(Equals("File name copied to clipboard"))
						expectClipboard(t, Equals("file1"))
					})
			}).
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Relative path")).
					Confirm().
					Tap(func() {
						t.ExpectToast(Equals("File path copied to clipboard"))
						expectClipboard(t, Equals("dir/file1"))
					})
			}).
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Absolute path")).
					Confirm().
					Tap(func() {
						t.ExpectToast(Equals("File path copied to clipboard"))
						repoDir, _ := os.Getwd()
						// On windows the following path would have backslashes, but we don't run integration tests on windows yet.
						expectClipboard(t, Equals(repoDir+"/dir/file1"))
					})
			}).
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Diff of selected file")).
					Confirm().
					Tap(func() {
						t.ExpectToast(Equals("File diff copied to clipboard"))
						expectClipboard(t,
							Contains("diff --git a/dir/file1 b/dir/file1").Contains("+2nd line").DoesNotContain("+1st line").
								DoesNotContain("diff --git a/dir/file2 b/dir/file2").DoesNotContain("+file2"))
					})
			}).
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Diff of all files")).
					Confirm().
					Tap(func() {
						t.ExpectToast(Equals("All files diff copied to clipboard"))
						expectClipboard(t,
							Contains("diff --git a/dir/file1 b/dir/file1").Contains("+2nd line").DoesNotContain("+1st line").
								Contains("diff --git a/dir/file2 b/dir/file2").Contains("+file2"))
					})
			}).
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Content of selected file")).
					Confirm().
					Tap(func() {
						t.ExpectToast(Equals("File content copied to clipboard"))
						expectClipboard(t, Equals("1st line\n2nd line\n"))
					})
			})

		t.Views().Commits().
			Focus().
			// Select commits 1 and 2
			Press(keys.Universal.RangeSelectDown).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("dir").IsSelected(),
				Contains("file1"),
				Contains("file2"),
			).
			NavigateToLine(Contains("file1")).
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Diff of selected file")).
					Confirm().
					Tap(func() {
						t.ExpectToast(Equals("File diff copied to clipboard"))
						expectClipboard(t,
							Contains("diff --git a/dir/file1 b/dir/file1").Contains("+1st line").Contains("+2nd line"))
					})
			}).
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Content of selected file")).
					Confirm().
					Tap(func() {
						t.ExpectToast(Equals("File content copied to clipboard"))
						expectClipboard(t, Equals("1st line\n2nd line\n"))
					})
			})
	},
})
