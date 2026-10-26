package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EditFilesMissingFromWorkingTree = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Edit or open files of a commit that a later commit deleted",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
		config.GetUserConfig().OS.Edit = "echo {{filename}} > edit-command"
		config.GetUserConfig().OS.EditAtLine = "echo {{filename}}:{{line}} > edit-command"
		config.GetUserConfig().OS.Open = "echo {{filename}} > open-command"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\n")
		shell.CreateFileAndAdd("file2", "two\n")
		shell.CreateFileAndAdd("file3", "three\n")
		shell.Commit("add files")
		shell.DeleteFileAndAdd("file2")
		shell.DeleteFileAndAdd("file3")
		shell.Commit("delete files")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("delete files").IsSelected(),
				Contains("add files"),
			).
			NavigateToLine(Contains("add files")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Contains("A file1"),
				Contains("A file2"),
				Contains("A file3"),
			).
			NavigateToLine(Contains("file2")).
			Press(keys.Universal.Edit).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Contains("/repo/file2' doesn't exist in the working tree")).
					Confirm()
			}).
			Press(keys.Universal.OpenFile).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Contains("/repo/file2' doesn't exist in the working tree")).
					Confirm()
			}).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			NavigateToLine(Contains("+two")).
			Press(keys.Universal.Edit).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Contains("/repo/file2' doesn't exist in the working tree")).
					Confirm()
			}).
			PressEscape()

		t.Views().CommitFiles().
			IsFocused().
			Press(keys.Universal.ToggleRangeSelect).
			SelectNextItem().
			SelectedLines(
				Contains("A file2"),
				Contains("A file3"),
			).
			Press(keys.Universal.Edit).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Contains("None of the selected files exist in the working tree at '")).
					Confirm()
			}).
			// If only some of the files are missing, the others are opened without
			// complaint
			NavigateToLine(Contains("file1")).
			SelectedLines(
				Contains("A file1"),
				Contains("A file2"),
			).
			Press(keys.Universal.Edit).
			Tap(func() {
				t.FileSystem().FileContent("edit-command",
					Contains("/repo/file1\n").DoesNotContain("file2"))
			})
	},
})
