package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveRangeToIndex = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Apply a custom patch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
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
			SelectNextItem().
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file2")).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		t.Views().Secondary().Content(Contains("second line"))
		t.Views().Secondary().Content(Contains("file two content"))

		t.Common().SelectPatchOption(MatchesRegexp(`Move patch out into index$`))

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file3").IsSelected(),
			).PressEscape()

		t.Views().Files().
			Focus().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  M  file1"),
				Equals("  A  file2"),
			)

		t.Views().Main().
			Content(Contains("second line"))

		t.Views().Files().Focus().NavigateToLine(Contains("file2"))

		t.Views().Main().
			Content(Contains("file two content"))
	},
})
