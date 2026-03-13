package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ToggleDirectory = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Toggle a directory for a custom patch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("dir1")
		shell.CreateFileAndAdd("dir1/file1", "file1 content\n")
		shell.CreateFileAndAdd("dir1/file2", "file2 content\n")
		shell.CreateFileAndAdd("dir1/file3", "file3 content\n")
		shell.CreateFileAndAdd("other-file", "other content\n")
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
				Equals("  ▼ dir1"),
				Equals("    A file1"),
				Equals("    A file2"),
				Equals("    A file3"),
				Equals("  A other-file"),
			).
			NavigateToLine(Contains("dir1")).
			PressPrimaryAction().
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir1").IsSelected(),
				Equals("    ● file1"),
				Equals("    ● file2"),
				Equals("    ● file3"),
				Equals("  A other-file"),
			)

		t.Views().Information().Content(Contains("Building patch"))

		// Toggle the directory again to remove all files from the patch
		t.Views().CommitFiles().
			PressPrimaryAction().
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir1").IsSelected(),
				Equals("    A file1"),
				Equals("    A file2"),
				Equals("    A file3"),
				Equals("  A other-file"),
			)
	},
})
