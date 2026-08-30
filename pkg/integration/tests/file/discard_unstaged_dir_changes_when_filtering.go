package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardUnstagedDirChangesWhenFiltering = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding unstaged changes in a directory when filtering by path",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("dir")
		shell.CreateFileAndAdd("dir/file-one", "original content\n")
		shell.CreateFileAndAdd("dir/file-two", "original content\n")

		shell.Commit("first commit")

		shell.UpdateFileAndAdd("dir/file-one", "original content\nnew content\n")
		shell.UpdateFileAndAdd("dir/file-two", "original content\nnew content\n")
		shell.UpdateFile("dir/file-one", "original content\nnew content\neven newer content\n")
		shell.UpdateFile("dir/file-two", "original content\nnew content\neven newer content\n")

		shell.CreateFile("dir/unstaged-file-one", "unstaged file")
		shell.CreateFile("dir/unstaged-file-two", "unstaged file")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ dir").IsSelected(),
				Equals("  MM file-one"),
				Equals("  MM file-two"),
				Equals("  ?? unstaged-file-one"),
				Equals("  ?? unstaged-file-two"),
			).
			Press(keys.Universal.StartSearch).
			Tap(func() {
				t.ExpectSearch().
					Type("one").
					Confirm()
			}).
			Lines(
				Equals("▼ dir").IsSelected(),
				Equals("  MM file-one"),
				Equals("  ?? unstaged-file-one"),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard unstaged changes")).
					Confirm()
			}).
			Press(keys.Universal.Return). // Cancel filtering
			Lines(
				Equals("▼ dir").IsSelected(),
				Equals("  M  file-one"),
				Equals("  MM file-two"),
				Equals("  ?? unstaged-file-two"),
			)

		t.FileSystem().FileContent("dir/file-one", Equals("original content\nnew content\n"))
		t.FileSystem().FileContent("dir/file-two", Equals("original content\nnew content\neven newer content\n"))
		t.FileSystem().PathNotPresent("dir/unstaged-file-one")
		t.FileSystem().FileContent("dir/unstaged-file-two", Equals("unstaged file"))
	},
})
