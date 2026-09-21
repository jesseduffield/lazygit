package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveMultiFileRangeToIndex = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a multi-file range from a commit to the index",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "first line\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "first line\nsecond line\n")
		shell.CreateFileAndAdd("file2", "file two content\n")
		shell.CreateFileAndAdd("file3", "file three content\n")
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("second commit").IsSelected(),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  M file1"),
				Equals("  A file2"),
				Equals("  A file3"),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+second line")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("+file two content")).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))
		t.Views().Secondary().
			Content(Contains("second line")).
			Content(Contains("file two content"))

		t.Common().SelectPatchOption(MatchesRegexp(`Move patch out into index$`))

		t.Views().Files().
			Focus().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  M  file1"),
				Equals("  A  file2"),
			)
		t.Views().Secondary().Content(Contains("second line"))
		t.Views().Files().NavigateToLine(Contains("file2"))
		t.Views().Secondary().Content(Contains("file two content"))
	},
})
