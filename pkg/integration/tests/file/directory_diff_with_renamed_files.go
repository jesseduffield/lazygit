package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DirectoryDiffWithRenamedFiles = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Selecting a directory in the files panel shows the renames of files that were moved into or out of it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("dir")
		shell.CreateDir("dir/nested")
		shell.CreateFileAndAdd("file1", "file1 content\n")
		shell.CreateFileAndAdd("dir/file2", "file2 content\n")
		shell.CreateFileAndAdd("dir/nested/file3", "file3 content\n")
		shell.Commit("initial commit")
		shell.RenameFileInGit("file1", "dir/file1")
		shell.RenameFileInGit("dir/file2", "dir/file2-renamed")
		shell.RenameFileInGit("dir/nested/file3", "file3")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ dir"),
				Equals("    R  file1 → file1"),
				Equals("    R  file2 → file2-renamed"),
				Equals("  R  dir/nested/file3 → file3"),
			)

		t.Views().Main().ContainsLines(
			Equals("diff --git a/file1 b/dir/file1"),
			Equals("similarity index 100%"),
			Equals("rename from file1"),
			Equals("rename to dir/file1"),
			Equals("diff --git a/dir/file2 b/dir/file2-renamed"),
			Equals("similarity index 100%"),
			Equals("rename from dir/file2"),
			Equals("rename to dir/file2-renamed"),
			Equals("diff --git a/dir/nested/file3 b/file3"),
			Equals("similarity index 100%"),
			Equals("rename from dir/nested/file3"),
			Equals("rename to file3"),
		)

		t.Views().Files().
			SelectNextItem().
			SelectedLine(Equals("  ▼ dir"))

		t.Views().Main().
			ContainsLines(
				Equals("diff --git a/file1 b/dir/file1"),
				Equals("similarity index 100%"),
				Equals("rename from file1"),
				Equals("rename to dir/file1"),
				Equals("diff --git a/dir/file2 b/dir/file2-renamed"),
				Equals("similarity index 100%"),
				Equals("rename from dir/file2"),
				Equals("rename to dir/file2-renamed"),
				Equals("diff --git a/dir/nested/file3 b/file3"),
				Equals("similarity index 100%"),
				Equals("rename from dir/nested/file3"),
				Equals("rename to file3"),
			)

		// The same applies when a filter reduces the directory to a single file
		t.Views().Files().
			FilterOrSearch("file1").
			Lines(
				Equals("▼ dir").IsSelected(),
				Equals("  R  file1 → file1"),
			)

		t.Views().Main().
			ContainsLines(
				Equals("diff --git a/file1 b/dir/file1"),
				Equals("similarity index 100%"),
				Equals("rename from file1"),
				Equals("rename to dir/file1"),
			)
	},
})
