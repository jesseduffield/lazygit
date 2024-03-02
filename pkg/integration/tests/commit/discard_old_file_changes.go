package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardOldFileChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding a range of files from an old commit.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("dir1/d1_file0", "file0\n")
		shell.CreateFileAndAdd("dir1/subd1/subfile0", "file1\n")
		shell.CreateFileAndAdd("dir2/d2_file1", "d2f1 content\n")
		shell.CreateFileAndAdd("dir2/d2_file2", "d2f4 content\n")
		shell.Commit("remove one file from this commit")

		shell.UpdateFileAndAdd("dir2/d2_file1", "d2f1 content\nsecond line\n")
		shell.DeleteFileAndAdd("dir2/d2_file2")
		shell.CreateFileAndAdd("dir2/d2_file3", "d2f3 content\n")
		shell.CreateFileAndAdd("dir2/d2_file4", "d2f2 content\n")
		shell.Commit("remove four files from this commit")

		shell.CreateFileAndAdd("dir1/fileToRemove", "file to remove content\n")
		shell.CreateFileAndAdd("dir1/multiLineFile", "this file has\ncontent on\nthree lines\n")
		shell.CreateFileAndAdd("dir1/subd1/file2ToRemove", "file2 to remove content\n")
		shell.Commit("remove changes in multiple dirs from this commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("remove changes in multiple dirs from this commit").IsSelected(),
				Contains("remove four files from this commit"),
				Contains("remove one file from this commit"),
			).
			NavigateToLine(Contains("remove one file from this commit")).
			PressEnter()

		// Check removing a single file from an old commit
		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("dir1").IsSelected(),
				Contains("subd1"),
				Contains("subfile0"),
				Contains("d1_file0"),
				Contains("dir2"),
				Contains("d2_file1"),
				Contains("d2_file2"),
			).
			NavigateToLine(Contains("d1_file0")).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Confirmation().
			Title(Equals("Discard file changes")).
			Content(Equals("Are you sure you want to remove changes to the selected file(s) from this commit?\n\nThis action will start a rebase, reverting these file changes. Be aware that if subsequent commits depend on these changes, you may need to resolve conflicts.\nNote: This will also reset any active custom patches.")).
			Confirm()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("dir1/subd1"),
				Contains("subfile0"),
				Contains("dir2"),
				Contains("d2_file1").IsSelected(),
				Contains("d2_file2"),
			).
			PressEscape()

		// Check removing 4 files in the same directory
		t.Views().Commits().
			Focus().
			Lines(
				Contains("remove changes in multiple dirs from this commit"),
				Contains("remove four files from this commit"),
				Contains("remove one file from this commit").IsSelected(),
			).
			NavigateToLine(Contains("remove four files from this commit")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("dir2").IsSelected(),
				Contains("d2_file1"),
				Contains("d2_file2"),
				Contains("d2_file3"),
				Contains("d2_file4"),
			).
			NavigateToLine(Contains("d2_file1")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("d2_file4")).
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
			PressEscape()

		// Check removing multiple files from 2 directories w/ a custom patch.
		// This checks node selection logic & if the custom patch is getting reset.
		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("remove changes in multiple dirs from this commit"),
				Contains("remove four files from this commit").IsSelected(),
				Contains("remove one file from this commit"),
			).
			NavigateToLine(Contains("remove changes in multiple dirs from this commit")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("dir1").IsSelected(),
				Contains("subd1"),
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
				Contains("subd1"),
				Contains("file2ToRemove"),
				Contains("fileToRemove"),
				Contains("multiLineFile").IsSelected(),
			).
			NavigateToLine(Contains("dir1")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("subd1")).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Confirmation().
			Title(Equals("Discard file changes")).
			Content(Equals("Are you sure you want to remove changes to the selected file(s) from this commit?\n\nThis action will start a rebase, reverting these file changes. Be aware that if subsequent commits depend on these changes, you may need to resolve conflicts.\nNote: This will also reset any active custom patches.")).
			Confirm()

		// "Building patch" will still be in this view if the patch isn't reset properly
		t.Views().Information().Content(DoesNotContain("Building patch"))

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("(none)"),
			)
	},
})
