package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RenamedFiles = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Regression test for the display of renamed files in the file tree",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("dir")
		shell.CreateDir("dir/nested")
		shell.CreateFileAndAdd("file1", "file1 content\n")
		shell.CreateFileAndAdd("dir/file2", "file2 content\n")
		shell.CreateFileAndAdd("dir/nested/file3", "file3 content\n")
		shell.Commit("initial commit")
		shell.RunCommand([]string{"git", "mv", "file1", "dir/file1"})
		shell.RunCommand([]string{"git", "mv", "dir/file2", "dir/file2-renamed"})
		shell.RunCommand([]string{"git", "mv", "dir/nested/file3", "file3"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir"),
				Equals("    R  file1 → file1"),
				Equals("    R  file2 → file2-renamed"),
				Equals("  R  dir/nested/file3 → file3"),
			)
	},
})
