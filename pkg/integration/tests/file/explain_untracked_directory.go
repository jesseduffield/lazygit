package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ExplainUntrackedDirectory = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "With status.showUntrackedFiles set to normal, the main view explains why an untracked directory is a single entry, and lists the files in it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.SetConfig("status.showUntrackedFiles", "normal")
		shell.CreateFile("dir/file1", "content")
		shell.CreateFile("dir/subdir/file2", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("?? dir/").IsSelected(),
			)

		t.Views().Main().
			Content(
				Contains("'dir' is an untracked directory").
					Contains("── Files in 'dir' ─").
					Contains("dir/file1").
					Contains("dir/subdir/file2"),
			)
	},
})
