package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardRangeSelect = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding a range of files from an old commit.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("dir1/file0", "file0\n")
		shell.CreateFileAndAdd("dir1/dir2/file1", "file1\n")
		shell.CreateFileAndAdd("dir3/file1", "d3f1 content\n")
		shell.CreateFileAndAdd("dir3/file4", "d3f4 content\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("dir3/file1", "d3f1 content\nsecond line\n")
		shell.CreateFileAndAdd("dir3/file2", "d3f2 content\n")
		shell.CreateFileAndAdd("dir3/file3", "d3f3 content\n")
		shell.DeleteFileAndAdd("dir3/file4")
		shell.Commit("first commit to change")

		shell.CreateFileAndAdd("dir1/fileToRemove", "file to remove content\n")
		shell.CreateFileAndAdd("dir1/multiLineFile", "this file has\ncontent on\nthree lines\n")
		shell.CreateFileAndAdd("dir1/dir2/file2ToRemove", "file2 to remove content\n")
		shell.Commit("second commit to change")

		shell.CreateFileAndAdd("file3", "file3")
		shell.Commit("third commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("third commit").IsSelected(),
				Contains("second commit to change"),
				Contains("first commit to change"),
				Contains("first commit"),
			).
			NavigateToLine(Contains("first commit to change")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("dir3").IsSelected(),
				Contains("file1"),
				Contains("file2"),
				Contains("file3"),
				Contains("file4"),
			).
			NavigateToLine(Contains("file1")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file4")).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Confirmation().
			Title(Equals("Discard file changes")).
			Content(Equals("Are you sure you want to remove changes to the selected file(s) from this commit?\n\nThis action will start a rebase, reverting these file changes. Be aware that if subsequent commits depend on these changes, you may need to resolve conflicts.\nNote: This will also reset any active custom patches.")).
			Confirm()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("(none)"),
			).
			// for some reason I need to press escape twice. Seems like it happens every time
			// more than one file is removed from a commit
			PressEscape().
			PressEscape()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("third commit"),
				Contains("second commit to change"),
				Contains("first commit to change").IsSelected(),
				Contains("first commit"),
			).
			NavigateToLine(Contains("second commit to change")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("dir1").IsSelected(),
				Contains("dir2"),
				Contains("file2ToRemove"),
				Contains("fileToRemove"),
				Contains("multiLineFile"),
			).
			NavigateToLine(Contains("multiLineFile")).
			PressEnter()

		t.Views().PatchBuilding().
			IsFocused().
			SelectedLine(
				Contains("+this file has"),
			).
			PressPrimaryAction().
			PressEscape()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("dir1"),
				Contains("dir2"),
				Contains("file2ToRemove"),
				Contains("fileToRemove"),
				Contains("multiLineFile").IsSelected(),
			).
			NavigateToLine(Contains("dir1")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("dir2")).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Confirmation().
			Title(Equals("Discard file changes")).
			Content(Equals("Are you sure you want to remove changes to the selected file(s) from this commit?\n\nThis action will start a rebase, reverting these file changes. Be aware that if subsequent commits depend on these changes, you may need to resolve conflicts.\nNote: This will also reset any active custom patches.")).
			Confirm()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("(none)"),
			)
	},
})
