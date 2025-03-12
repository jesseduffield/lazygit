package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardUnstagedDirChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding unstaged changes in a directory",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("dir")
		shell.CreateFileAndAdd("dir/file-one", "original content\n")

		shell.Commit("first commit")

		shell.UpdateFileAndAdd("dir/file-one", "original content\nnew content\n")
		shell.UpdateFile("dir/file-one", "original content\nnew content\neven newer content\n")

		shell.CreateDir("dir/subdir")
		shell.CreateFile("dir/subdir/unstaged-file-one", "unstaged file")
		shell.CreateFile("dir/unstaged-file-two", "unstaged file")

		shell.CreateFile("unstaged-file-three", "unstaged file")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ dir"),
				Equals("    ▼ subdir"),
				Equals("      ?? unstaged-file-one"),
				Equals("    MM file-one"),
				Equals("    ?? unstaged-file-two"),
				Equals("  ?? unstaged-file-three"),
			).
			SelectNextItem().
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard unstaged changes")).
					Confirm()
			}).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ dir").IsSelected(),
				Equals("    M  file-one"),
				// this guy remains untouched because it wasn't inside the 'dir' directory
				Equals("  ?? unstaged-file-three"),
			)

		t.FileSystem().FileContent("dir/file-one", Equals("original content\nnew content\n"))
	},
})
