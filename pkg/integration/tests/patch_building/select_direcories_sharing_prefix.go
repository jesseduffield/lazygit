package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectDirecoriesSharingPrefix = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Select directories sharing a prefix in the commit files view and add them to a custom patch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("foo/file", "file1 content")
		shell.CreateFileAndAdd("foobar/file", "file2 content")
		shell.Commit("first commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("first commit").IsSelected(),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ foo"),
				Equals("    A file"),
				Equals("  ▼ foobar"),
				Equals("    A file"),
			).
			SelectNextItem().
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("foobar")).
			PressPrimaryAction().
			Lines(
				Equals("▼ /"),
				Equals("  ▼ foo").IsSelected(),
				Equals("    ● file").IsSelected(),
				Equals("  ▼ foobar").IsSelected(),
				Equals("    ● file"),
			)

		t.Views().Information().Content(Contains("Building patch"))

		t.Views().Secondary().Content(
			Contains("foo/file").Contains("foobar/file"),
		)
	},
})
