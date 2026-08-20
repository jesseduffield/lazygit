package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MovePatchToLaterCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a patch from a commit to a later commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("dir")
		shell.CreateFileAndAdd("dir/file1", "file1 content")
		shell.CreateFileAndAdd("dir/file2", "file2 content")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("dir/file1", "file1 content with old changes")
		shell.DeleteFileAndAdd("dir/file2")
		shell.CreateFileAndAdd("dir/file3", "file3 content")
		shell.Commit("commit to move from")

		shell.CreateFileAndAdd("unrelated-file", "")
		shell.Commit("destination commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("destination commit").IsSelected(),
				Contains("commit to move from"),
				Contains("first commit"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("dir").IsSelected(),
				Contains("  M file1"),
				Contains("  D file2"),
				Contains("  A file3"),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("-file1 content")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("+file3 content")).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))
		t.Views().Commits().Focus().SelectPreviousItem()
		t.Common().SelectPatchOption(Contains("Move patch to selected commit"))

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("destination commit").IsSelected(),
				Contains("commit to move from"),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ dir"),
				Equals("    M file1"),
				Equals("    D file2"),
				Equals("    A file3"),
				Equals("  A unrelated-file"),
			).
			PressEscape()

		t.Views().Commits().
			IsFocused().
			SelectNextItem().
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("(none)"),
			)
	},
})
